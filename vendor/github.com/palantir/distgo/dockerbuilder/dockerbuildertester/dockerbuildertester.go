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

package dockerbuildertester

import (
	"bytes"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"sort"
	"testing"
	"time"

	"github.com/nmiyake/pkg/dirs"
	"github.com/nmiyake/pkg/gofiles"
	"github.com/palantir/godel/framework/pluginapitester"
	"github.com/palantir/pkg/gittest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type TestCase struct {
	Name        string
	Specs       []gofiles.GoFileSpec
	ConfigFiles map[string]string
	Args        []string
	WantError   bool
	WantOutput  func(projectDir string) string
	Validate    func(projectDir string)
}

var builtinSpecs = []gofiles.GoFileSpec{
	{
		RelPath: "godelw",
		Src:     `// placeholder`,
	},
	{
		RelPath: ".gitignore",
		Src: `/out
`,
	},
}

// RunAssetDockerBuilderTest tests the "docker" operation using the provided asset. Uses the provided plugin provider
// and asset provider to resolve the plugin and asset and invokes the "docker" command.
func RunAssetDockerBuilderTest(t *testing.T,
	pluginProvider pluginapitester.PluginProvider,
	assetProvider pluginapitester.AssetProvider,
	testCases []TestCase,
) {
	wd, err := os.Getwd()
	require.NoError(t, err)

	tmpDir, cleanup, err := dirs.TempDir("", "")
	require.NoError(t, err)
	if !filepath.IsAbs(tmpDir) {
		tmpDir = path.Join(wd, tmpDir)
	}
	defer cleanup()

	tmpDir, err = filepath.EvalSymlinks(tmpDir)
	require.NoError(t, err)

	for i, tc := range testCases {
		projectDir, err := ioutil.TempDir(tmpDir, "")
		require.NoError(t, err)

		gittest.InitGitDir(t, projectDir)
		require.NoError(t, err)

		var sortedKeys []string
		for k := range tc.ConfigFiles {
			sortedKeys = append(sortedKeys, k)
		}
		sort.Strings(sortedKeys)

		for _, k := range sortedKeys {
			err = os.MkdirAll(path.Dir(path.Join(projectDir, k)), 0755)
			require.NoError(t, err)
			err = ioutil.WriteFile(path.Join(projectDir, k), []byte(tc.ConfigFiles[k]), 0644)
			require.NoError(t, err)
		}

		// write files required for test framework
		_, err = gofiles.Write(projectDir, builtinSpecs)
		require.NoError(t, err)
		// write provided specs
		_, err = gofiles.Write(projectDir, tc.Specs)
		require.NoError(t, err)

		// commit all files and tag project as v1.0.0
		gittest.CommitAllFiles(t, projectDir, "Commit all files")
		gittest.CreateGitTag(t, projectDir, "v1.0.0")

		outputBuf := &bytes.Buffer{}
		func() {
			wantWd := projectDir
			err = os.Chdir(wantWd)
			require.NoError(t, err)
			defer func() {
				err = os.Chdir(wd)
				require.NoError(t, err)
			}()

			var assetProviders []pluginapitester.AssetProvider
			if assetProvider != nil {
				assetProviders = append(assetProviders, assetProvider)
			}

			// dist artifacts are considered stale if the generation time of the oldest artifact matches the
			// modification time of the configuration file (at second granularity). Wait until values are different to
			// ensure that cache will be valid (will take at most 1 second).
			distPluginConfigYML := path.Join(projectDir, "godel", "config", "dist-plugin.yml")
			if fi, err := os.Stat(distPluginConfigYML); err == nil {
				modTimeToSecond := fi.ModTime().Truncate(time.Second)
				for {
					currTimeToSecond := time.Now().Truncate(time.Second)
					if currTimeToSecond != modTimeToSecond {
						break
					}
					time.Sleep(100 * time.Millisecond)
				}
			}

			// run dist task first
			func() {
				runPluginCleanup, err := pluginapitester.RunPlugin(
					pluginProvider,
					assetProviders,
					"dist", nil,
					projectDir, false, outputBuf)
				defer runPluginCleanup()
				require.NoError(t, err, "Case %d: %s\nDist operation failed with output:\n%s", i, tc.Name, outputBuf.String())
				outputBuf = &bytes.Buffer{}
			}()

			var args []string
			args = append(args, tc.Args...)

			runPluginCleanup, err := pluginapitester.RunPlugin(
				pluginProvider,
				assetProviders,
				"docker", args,
				projectDir, false, outputBuf)
			defer runPluginCleanup()
			if tc.WantError {
				require.EqualError(t, err, "", "Case %d: %s", i, tc.Name)
			} else {
				require.NoError(t, err, "Case %d: %s\nOutput:\n%s", i, tc.Name, outputBuf.String())
			}
			if tc.WantOutput != nil {
				assert.Equal(t, tc.WantOutput(projectDir), outputBuf.String(), "Case %d: %s", i, tc.Name)
			}
			if tc.Validate != nil {
				tc.Validate(projectDir)
			}
		}()
	}
}
