// Copyright 2016 Palantir Technologies, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package plugins

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"testing"
	"time"

	"github.com/mholt/archiver"
	"github.com/nmiyake/pkg/dirs"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/palantir/godel/framework/artifactresolver"
	"github.com/palantir/godel/framework/godellauncher"
	"github.com/palantir/godel/framework/internal/pathsinternal"
	"github.com/palantir/godel/framework/pluginapi/v2/pluginapi"
	"github.com/palantir/godel/pkg/osarch"
)

var pluginScriptTmpl = fmt.Sprintf(`#!/usr/bin/env bash

if [ "$1" = "%s" ]; then
    echo '{"pluginSchemaVersion":"2","group":"com.palantir","product":"%s","version":"1.0.0","usesConfig":true,"tasks":[{"name":"fooTest","description":"","command":["foo"],"globalFlagOptions":null,"verifyOptions":null}],"upgradeTask":null}'
fi
`, pluginapi.PluginInfoCommandName, "%s")

// If the plugin JSON schema changes, uncomment the following and run it to generate the "echo" line in the pluginScriptTmpl above.

//func TestPrintPluginInfoJSON(t *testing.T) {
//	pluginInfo, err := pluginapi.NewPluginInfo("com.palantir", "placeholder-plugin", "1.0.0",
//		pluginapi.PluginInfoUsesConfigFile(),
//		pluginapi.PluginInfoTaskInfo(
//			"fooTest",
//			"",
//			pluginapi.TaskInfoCommand("foo"),
//		),
//	)
//	require.NoError(t, err)
//	bytes, err := json.Marshal(pluginInfo)
//	require.NoError(t, err)
//	fmt.Println(`echo '` + strings.Replace(string(bytes), "placeholder-plugin", `%s`, 1) + `'`)
//}

func TestInfoFromResolved(t *testing.T) {
	tmpDir, cleanup, err := dirs.TempDir("", "")
	require.NoError(t, err)
	defer cleanup()

	pluginName := newPluginName()
	pluginFile := path.Join(tmpDir, fmt.Sprintf("com.palantir-%s-1.0.0", pluginName))
	err = ioutil.WriteFile(pluginFile, []byte(fmt.Sprintf(pluginScriptTmpl, pluginName)), 0755)
	require.NoError(t, err)

	gotInfo, err := pluginapi.InfoFromPlugin(path.Join(tmpDir, pathsinternal.PluginFileName(artifactresolver.Locator{
		Group:   "com.palantir",
		Product: pluginName,
		Version: "1.0.0",
	})))
	require.NoError(t, err)

	wantInfo := pluginapi.MustNewPluginInfo(
		"com.palantir",
		pluginName,
		"1.0.0",
		pluginapi.PluginInfoUsesConfigFile(),
		pluginapi.PluginInfoTaskInfo("fooTest", "", pluginapi.TaskInfoCommand("foo")),
	)
	assert.Equal(t, wantInfo, gotInfo)
}

func TestInfoFromResolvedError(t *testing.T) {
	tmpDir, cleanup, err := dirs.TempDir("", "")
	require.NoError(t, err)
	defer cleanup()

	pluginName := newPluginName()
	pluginFile := path.Join(tmpDir, fmt.Sprintf("com.palantir-%s-1.0.0", pluginName))
	err = ioutil.WriteFile(pluginFile, []byte(`#!/usr/bin/env bash

exit 1
`), 0755)
	require.NoError(t, err)

	_, err = pluginapi.InfoFromPlugin(path.Join(tmpDir, pathsinternal.PluginFileName(artifactresolver.Locator{
		Group:   "com.palantir",
		Product: pluginName,
		Version: "1.0.0",
	})))
	require.Error(t, err)
	assert.Regexp(t, regexp.QuoteMeta("command [")+".+"+regexp.QuoteMeta(fmt.Sprintf("/com.palantir-%s-1.0.0 %s] failed.\nError:\nexit status 1\nOutput:\n\n: exit status 1", pluginName, pluginapi.PluginInfoCommandName)), err.Error())
}

func TestResolvePlugins(t *testing.T) {
	tmpDir, cleanup, err := dirs.TempDir("", "")
	require.NoError(t, err)
	defer cleanup()

	loc, resolver, osArch := createTestPlugin(t, tmpDir)

	pluginsDir := path.Join(tmpDir, "plugins")
	err = os.Mkdir(pluginsDir, 0755)
	require.NoError(t, err)

	assetsDir := path.Join(tmpDir, "assets")
	err = os.Mkdir(assetsDir, 0755)
	require.NoError(t, err)

	downloadsDir := path.Join(tmpDir, "downloads")
	err = os.Mkdir(downloadsDir, 0755)
	require.NoError(t, err)

	outBuf := &bytes.Buffer{}
	plugins, errs := resolvePlugins(pluginsDir, assetsDir, downloadsDir, osArch, godellauncher.PluginsParam{
		Plugins: []godellauncher.SinglePluginParam{
			{
				LocatorWithResolverParam: artifactresolver.LocatorWithResolverParam{
					LocatorWithChecksums: artifactresolver.LocatorParam{
						Locator: loc,
					},
					Resolver: resolver,
				},
			},
		},
	}, outBuf)
	assert.NoError(t, errs)

	wantPlugins := map[artifactresolver.Locator]pluginInfoWithAssets{
		loc: {
			PluginInfo: pluginapi.MustNewPluginInfo(
				"com.palantir",
				loc.Product,
				"1.0.0",
				pluginapi.PluginInfoUsesConfigFile(),
				pluginapi.PluginInfoTaskInfo("fooTest", "", pluginapi.TaskInfoCommand("foo")),
			),
		},
	}
	assert.Equal(t, wantPlugins, plugins)
}

func createTestPlugin(t *testing.T, tmpDir string) (artifactresolver.Locator, artifactresolver.Resolver, osarch.OSArch) {
	pluginName := newPluginName()
	testProductDir := path.Join(tmpDir, "repo", "com", "palantir", pluginName, "1.0.0")
	err := os.MkdirAll(testProductDir, 0755)
	require.NoError(t, err)

	testProductPath := path.Join(testProductDir, pluginName)
	err = ioutil.WriteFile(testProductPath, []byte(fmt.Sprintf(pluginScriptTmpl, pluginName)), 0755)
	require.NoError(t, err)

	testProductTGZPath := path.Join(testProductDir, pluginName+"-darwin-amd64-1.0.0.tgz")
	err = archiver.TarGz.Make(testProductTGZPath, []string{testProductPath})
	require.NoError(t, err)

	tmpDirAbs, err := filepath.Abs(tmpDir)
	require.NoError(t, err)

	testResolver, err := artifactresolver.NewTemplateResolver(tmpDirAbs + "/repo/{{GroupPath}}/{{Product}}/{{Version}}/{{Product}}-{{OS}}-{{Arch}}-{{Version}}.tgz")
	require.NoError(t, err)

	darwinOSArch, err := osarch.New("darwin-amd64")
	require.NoError(t, err)

	return artifactresolver.Locator{
		Group:   "com.palantir",
		Product: pluginName,
		Version: "1.0.0",
	}, testResolver, darwinOSArch
}

func TestVerifyPluginCompatibility(t *testing.T) {
	for i, tc := range []struct {
		name  string
		input map[artifactresolver.Locator]pluginInfoWithAssets
		want  string
	}{
		{
			"no plugin conflicts",
			map[artifactresolver.Locator]pluginInfoWithAssets{
				{
					Group:   "com.palantir",
					Product: "foo-plugin",
					Version: "1.0.0",
				}: {
					PluginInfo: pluginapi.MustNewPluginInfo("com.palantir", "foo-plugin", "1.0.0",
						pluginapi.PluginInfoUsesConfigFile(),
					),
				},
			},
			"",
		},
		{
			"verify catches plugins with same group and product but different version",
			map[artifactresolver.Locator]pluginInfoWithAssets{
				{
					Group:   "com.palantir",
					Product: "foo-plugin",
					Version: "1.0.0",
				}: {
					PluginInfo: pluginapi.MustNewPluginInfo("com.palantir", "foo-plugin", "1.0.0",
						pluginapi.PluginInfoUsesConfigFile(),
						pluginapi.PluginInfoTaskInfo("foo", ""),
					),
				},
				{
					Group:   "com.palantir",
					Product: "foo-plugin",
					Version: "2.0.0",
				}: {
					PluginInfo: pluginapi.MustNewPluginInfo("com.palantir", "foo-plugin", "1.0.0",
						pluginapi.PluginInfoUsesConfigFile(),
						pluginapi.PluginInfoTaskInfo("foo", ""),
					),
				},
			},
			`2 plugins had compatibility issues:
    com.palantir:foo-plugin:1.0.0:
        different version of the same plugin
    com.palantir:foo-plugin:2.0.0:
        different version of the same plugin`,
		},
		{
			"verify catches plugins with conflicting commands",
			map[artifactresolver.Locator]pluginInfoWithAssets{
				{
					Group:   "com.palantir",
					Product: "foo-plugin",
					Version: "1.0.0",
				}: {
					PluginInfo: pluginapi.MustNewPluginInfo("com.palantir", "foo-plugin", "1.0.0",
						pluginapi.PluginInfoUsesConfigFile(),
						pluginapi.PluginInfoTaskInfo("foo", ""),
					),
				},
				{
					Group:   "com.palantir",
					Product: "bar-plugin",
					Version: "2.0.0",
				}: {
					PluginInfo: pluginapi.MustNewPluginInfo("com.palantir", "bar-plugin", ".0.0",
						pluginapi.PluginInfoUsesConfigFile(),
						pluginapi.PluginInfoTaskInfo("foo", ""),
					),
				},
			},
			`2 plugins had compatibility issues:
    com.palantir:bar-plugin:2.0.0:
        provides conflicting tasks: [foo]
    com.palantir:foo-plugin:1.0.0:
        provides conflicting tasks: [foo]`,
		},
		{
			"verify catches plugins with same product name that both use config files",
			map[artifactresolver.Locator]pluginInfoWithAssets{
				{
					Group:   "com.palantir",
					Product: "foo-plugin",
					Version: "1.0.0",
				}: {
					PluginInfo: pluginapi.MustNewPluginInfo("com.palantir", "foo-plugin", "1.0.0",
						pluginapi.PluginInfoUsesConfigFile(),
						pluginapi.PluginInfoTaskInfo("foo", ""),
					),
				},
				{
					Group:   "com.zcorp",
					Product: "foo-plugin",
					Version: "2.0.0",
				}: {
					PluginInfo: pluginapi.MustNewPluginInfo("com.zcorp", "foo-plugin", "2.0.0",
						pluginapi.PluginInfoUsesConfigFile(),
						pluginapi.PluginInfoTaskInfo("bar", ""),
					),
				},
			},
			`2 plugins had compatibility issues:
    com.palantir:foo-plugin:1.0.0:
        plugins have the same product name and both use configuration (this not currently supported -- if this situation is encountered, please file an issue to flag it)
    com.zcorp:foo-plugin:2.0.0:
        plugins have the same product name and both use configuration (this not currently supported -- if this situation is encountered, please file an issue to flag it)`,
		},
	} {
		got := verifyPluginCompatibility(tc.input)
		if tc.want == "" {
			assert.NoError(t, got, "Case %d: %s", i, tc.name)
		} else {
			assert.EqualError(t, got, tc.want, "Case %d: %s", i, tc.name)
		}
	}
}

func newPluginName() string {
	return fmt.Sprintf("tester-%d-plugin", time.Now().Unix())
}
