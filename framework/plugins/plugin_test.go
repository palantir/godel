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

	"github.com/nmiyake/archiver"
	"github.com/nmiyake/pkg/dirs"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/palantir/godel/apps/distgo/pkg/osarch"
	"github.com/palantir/godel/framework/pluginapi"
)

var pluginScriptTmpl = fmt.Sprintf(`#!/usr/bin/env bash

if [ "$1" = "%s" ]; then
    echo '{"pluginSchemaVersion":"1","id":"com.palantir:%s:1.0.0","configFileName":"tester.yml","tasks":[{"name":"fooTest","description":"","command":["foo"],"globalFlagOptions":null,"verifyOptions":null}]}'
fi
`, pluginapi.InfoCommandName, "%s")

func TestInfoFromResolved(t *testing.T) {
	tmpDir, cleanup, err := dirs.TempDir("", "")
	require.NoError(t, err)
	defer cleanup()

	pluginName := newPluginName()
	pluginFile := path.Join(tmpDir, fmt.Sprintf("com.palantir-%s-1.0.0", pluginName))
	err = ioutil.WriteFile(pluginFile, []byte(fmt.Sprintf(pluginScriptTmpl, pluginName)), 0755)
	require.NoError(t, err)

	gotInfo, err := pluginapi.InfoFromPlugin(path.Join(tmpDir, pluginFileName(locator{
		Group:   "com.palantir",
		Product: pluginName,
		Version: "1.0.0",
	})))
	require.NoError(t, err)

	wantInfo := pluginapi.MustNewInfo(
		"com.palantir",
		pluginName,
		"1.0.0",
		"tester.yml",
		pluginapi.MustNewTaskInfo("fooTest", "", pluginapi.TaskInfoCommand("foo")),
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

	_, err = pluginapi.InfoFromPlugin(path.Join(tmpDir, pluginFileName(locator{
		Group:   "com.palantir",
		Product: pluginName,
		Version: "1.0.0",
	})))
	require.Error(t, err)
	assert.Regexp(t, regexp.QuoteMeta("command [")+".+"+regexp.QuoteMeta(fmt.Sprintf("/com.palantir-%s-1.0.0 %s] failed.\nError:\nexit status 1\nOutput:\n\n: exit status 1", pluginName, pluginapi.InfoCommandName)), err.Error())
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
	plugins, errs := resolvePlugins(pluginsDir, assetsDir, downloadsDir, osArch, projectParams{
		Plugins: []singlePluginParam{
			{
				locatorWithResolverParam: locatorWithResolverParam{
					LocatorWithChecksums: locatorWithChecksumsParam{
						locator: loc,
					},
					Resolver: resolver,
				},
			},
		},
	}, outBuf)
	assert.NoError(t, errs)

	wantPlugins := map[locator]pluginInfoWithAssets{
		loc: {
			PluginInfo: pluginapi.MustNewInfo(
				"com.palantir",
				loc.Product,
				"1.0.0",
				"tester.yml",
				pluginapi.MustNewTaskInfo("fooTest", "", pluginapi.TaskInfoCommand("foo")),
			),
		},
	}
	assert.Equal(t, wantPlugins, plugins)
}

func createTestPlugin(t *testing.T, tmpDir string) (locator, resolver, osarch.OSArch) {
	pluginName := newPluginName()
	testProductDir := path.Join(tmpDir, "repo", "com", "palantir", pluginName, "1.0.0")
	err := os.MkdirAll(testProductDir, 0755)
	require.NoError(t, err)

	testProductPath := path.Join(testProductDir, pluginName)
	err = ioutil.WriteFile(testProductPath, []byte(fmt.Sprintf(pluginScriptTmpl, pluginName)), 0755)
	require.NoError(t, err)

	testProductTGZPath := path.Join(testProductDir, pluginName+"-darwin-amd64-1.0.0.tgz")
	err = archiver.TarGz(testProductTGZPath, []string{testProductPath})
	require.NoError(t, err)

	tmpDirAbs, err := filepath.Abs(tmpDir)
	require.NoError(t, err)

	testResolver, err := newTemplateResolver(tmpDirAbs + "/repo/{{GroupPath}}/{{Product}}/{{Version}}/{{Product}}-{{OS}}-{{Arch}}-{{Version}}.tgz")
	require.NoError(t, err)

	darwinOSArch, err := osarch.New("darwin-amd64")
	require.NoError(t, err)

	return locator{
		Group:   "com.palantir",
		Product: pluginName,
		Version: "1.0.0",
	}, testResolver, darwinOSArch
}

func TestVerifyPluginCompatibility(t *testing.T) {
	for i, tc := range []struct {
		name  string
		input map[locator]pluginInfoWithAssets
		want  string
	}{
		{
			"no plugin conflicts",
			map[locator]pluginInfoWithAssets{
				locator{
					Group:   "com.palantir",
					Product: "foo",
					Version: "1.0.0",
				}: {
					PluginInfo: pluginapi.MustNewInfo("com.palantir", "foo", "1.0.0", "foo.yml"),
				},
			},
			"",
		},
		{
			"verify catches plugins with same group and product but different version",
			map[locator]pluginInfoWithAssets{
				locator{
					Group:   "com.palantir",
					Product: "foo",
					Version: "1.0.0",
				}: {
					PluginInfo: pluginapi.MustNewInfo("com.palantir", "foo", "1.0.0", "foo.yml",
						pluginapi.MustNewTaskInfo("foo", "", nil, nil, nil),
					),
				},
				locator{
					Group:   "com.palantir",
					Product: "foo",
					Version: "2.0.0",
				}: {
					PluginInfo: pluginapi.MustNewInfo("com.palantir", "foo", "1.0.0", "foo.yml",
						pluginapi.MustNewTaskInfo("foo", "", nil, nil, nil),
					),
				},
			},
			`2 plugins had compatibility issues:
    com.palantir:foo:1.0.0:
        different version of the same plugin
    com.palantir:foo:2.0.0:
        different version of the same plugin`,
		},
		{
			"verify catches plugins with conflicting commands",
			map[locator]pluginInfoWithAssets{
				locator{
					Group:   "com.palantir",
					Product: "foo",
					Version: "1.0.0",
				}: {
					PluginInfo: pluginapi.MustNewInfo("com.palantir", "foo", "1.0.0", "foo.yml",
						pluginapi.MustNewTaskInfo("foo", "", nil, nil, nil),
					),
				},
				locator{
					Group:   "com.palantir",
					Product: "bar",
					Version: "2.0.0",
				}: {
					PluginInfo: pluginapi.MustNewInfo("com.palantir", "bar", "1.0.0", "foo.yml",
						pluginapi.MustNewTaskInfo("foo", "", nil, nil),
					),
				},
			},
			`2 plugins had compatibility issues:
    com.palantir:bar:2.0.0:
        provides conflicting tasks: [foo]
    com.palantir:foo:1.0.0:
        provides conflicting tasks: [foo]`,
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
	return fmt.Sprintf("tester-%d", time.Now().Unix())
}
