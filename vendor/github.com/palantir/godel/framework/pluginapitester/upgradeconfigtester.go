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

// Package pluginapitester provides functions that simulate invoking a plugin from gödel. Can be used to test plugin
// implementations in plugin projects.
package pluginapitester

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"sort"
	"testing"

	"github.com/nmiyake/pkg/dirs"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/palantir/godel/framework/builtintasks"
	"github.com/palantir/godel/framework/godellauncher"
	"github.com/palantir/godel/framework/godellauncher/defaulttasks"
	"github.com/palantir/godel/framework/pluginapi/v2/pluginapi"
)

// RunUpgradeConfig runs the "upgrade-config" task with the provided plugin and assets loaded. The plugin is loaded in
// the same manner that it would be for gödel itself. RunPlugin is mostly equivalent to calling
// "./godelw upgrade-config" where "godelw" is in "projectDir". The one difference is that the "upgrade-config" task
// will only run the upgrade task provided by the plugin (it will not run any builtin config upgrades). If "debug" is
// true, then it is the equivalent of calling "./godelw --debug upgrade-config". If the project directory does not
// contain a file named "godelw", this function creates the path. The returned "cleanup" function removes the "godelw"
// file if it was created by this function and is suitable to defer. Returns an error if the specified plugin does not
// provide a config upgrader.
func RunUpgradeConfig(
	pluginProvider PluginProvider,
	assetProviders []AssetProvider,
	legacy bool,
	projectDir string,
	debug bool,
	stdout io.Writer) (cleanup func(), rErr error) {

	var pluginPath string
	if pluginProvider != nil {
		pluginPath = pluginProvider.PluginFilePath()
	}
	var assets []string
	for _, asset := range assetProviders {
		assets = append(assets, asset.AssetFilePath())
	}

	cleanup = func() {}

	var upgradeTasks []godellauncher.UpgradeConfigTask
	if pluginPath != "" {
		info, err := pluginapi.InfoFromPlugin(pluginPath)
		if err != nil {
			return cleanup, err
		}
		upgradeTask := info.UpgradeConfigTask(pluginPath, assets)
		if upgradeTask == nil {
			return cleanup, errors.Errorf("plugin %s does not provide an upgrade task", pluginPath)
		}
		upgradeTasks = []godellauncher.UpgradeConfigTask{*upgradeTask}
	} else {
		// if no plugin was specified, add built-in upgraders
		upgradeTasks = defaulttasks.BuiltinUpgradeConfigTasks()
	}

	var taskArgs []string
	if legacy {
		taskArgs = append(taskArgs, "--legacy")
	}
	task := builtintasks.UpgradeConfigTask(upgradeTasks)
	globalConfig := godellauncher.GlobalConfig{
		Task:     task.Name,
		TaskArgs: taskArgs,
		Debug:    debug,
	}
	if projectDir != "" {
		if !filepath.IsAbs(projectDir) {
			wd, err := os.Getwd()
			if err != nil {
				return cleanup, errors.Wrapf(err, "failed to determine working directory")
			}
			projectDir = path.Join(wd, projectDir)
		}

		godelwPath := path.Join(projectDir, "godelw")
		if _, err := os.Stat(godelwPath); os.IsNotExist(err) {
			if err := ioutil.WriteFile(godelwPath, nil, 0644); err != nil {
				return cleanup, errors.Wrapf(err, "failed to create temporary godelw file")
			}
			cleanup = func() {
				if err := os.Remove(godelwPath); err != nil {
					fmt.Println(errors.Wrapf(err, "failed to remove temporary godelw file"))
				}
			}
		}
		globalConfig.Wrapper = path.Join(projectDir, "godelw")
	}
	return cleanup, task.Run(globalConfig, stdout)
}

type UpgradeConfigTestCase struct {
	Name        string
	ConfigFiles map[string]string
	Legacy      bool
	WantError   bool
	WantOutput  string
	WantFiles   map[string]string
}

// RunUpgradeConfigTest tests the "upgrade-config" operation using the provided plugin and assets. Resolves the plugin
// using the provided locator and resolver, provides it with the assets and invokes the "upgrade-config" command.
func RunUpgradeConfigTest(t *testing.T,
	pluginProvider PluginProvider,
	assetProviders []AssetProvider,
	testCases []UpgradeConfigTestCase,
) {
	tmpDir, cleanup, err := dirs.TempDir("", "")
	require.NoError(t, err)
	defer cleanup()

	tmpDir, err = filepath.EvalSymlinks(tmpDir)
	require.NoError(t, err)

	for i, tc := range testCases {
		projectDir, err := ioutil.TempDir(tmpDir, "")
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

		outputBuf := &bytes.Buffer{}
		func() {
			runPluginCleanup, err := RunUpgradeConfig(pluginProvider, assetProviders, tc.Legacy, projectDir, false, outputBuf)
			defer runPluginCleanup()
			if tc.WantError {
				require.EqualError(t, err, "", "Case %d: %s\nOutput: %s", i, tc.Name, outputBuf.String())
			} else {
				require.NoError(t, err, "Case %d: %s\nOutput: %s", i, tc.Name, outputBuf.String())
			}
			assert.Equal(t, tc.WantOutput, outputBuf.String(), "Case %d: %s", i, tc.Name)

			var sortedKeys []string
			for k := range tc.WantFiles {
				sortedKeys = append(sortedKeys, k)
			}
			sort.Strings(sortedKeys)
			for _, k := range sortedKeys {
				wantContent := tc.WantFiles[k]
				bytes, err := ioutil.ReadFile(path.Join(projectDir, k))
				require.NoError(t, err, "Case %d: %s", i, tc.Name)
				assert.Equal(t, wantContent, string(bytes), "Case %d: %s\nContent of file %s did not match expectation.\nActual:\n%s", i, tc.Name, k, string(bytes))
			}
		}()
	}
}
