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
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"

	"github.com/pkg/errors"

	"github.com/palantir/godel/v2/framework/godellauncher"
	"github.com/palantir/godel/v2/framework/pluginapi/v2/pluginapi"
)

// RunPlugin runs a plugin with the specified arguments. The plugin is loaded in the same manner that it would be for
// gödel itself. RunPlugin is equivalent to calling "./godelw [taskName] [args]" where "godelw" is in "projectDir". Note
// that this call does not change the working directory, so any arguments that take relative paths should take that into
// account. If "debug" is true, then it is the equivalent of calling "./godelw --debug [taskName] [args]". If the
// project directory does not contain a file named "godelw", this function creates the path. The returned "cleanup"
// function removes the "godelw" file if it was created by this function and is suitable to defer.
func RunPlugin(
	pluginProvider PluginProvider,
	assetProviders []AssetProvider,
	taskName string,
	args []string,
	projectDir string,
	debug bool,
	stdout io.Writer) (cleanup func(), rErr error) {

	pluginPath := pluginProvider.PluginFilePath()
	var assets []string
	for _, asset := range assetProviders {
		assets = append(assets, asset.AssetFilePath())
	}

	cleanup = func() {}
	info, err := pluginapi.InfoFromPlugin(pluginPath)
	if err != nil {
		return cleanup, err
	}
	taskMap := make(map[string]godellauncher.Task)
	var taskNames []string
	for _, task := range info.Tasks(pluginPath, assets) {
		taskMap[task.Name] = task
		taskNames = append(taskNames, task.Name)
	}
	task, ok := taskMap[taskName]
	if !ok {
		return cleanup, errors.Errorf("task %s does not exist. Valid tasks: %v", taskName, taskNames)
	}

	globalConfig := godellauncher.GlobalConfig{
		Task:     taskName,
		TaskArgs: args,
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
