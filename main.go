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

package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"

	"github.com/nmiyake/pkg/dirs"
	"github.com/nmiyake/pkg/errorstringer"
	"github.com/pkg/errors"

	"github.com/palantir/godel/framework/builtintasks"
	"github.com/palantir/godel/framework/godel"
	"github.com/palantir/godel/framework/godel/config"
	"github.com/palantir/godel/framework/godellauncher"
	"github.com/palantir/godel/framework/godellauncher/defaulttasks"
	"github.com/palantir/godel/framework/plugins"
)

func main() {
	if err := dirs.SetGoEnvVariables(); err != nil {
		printErrAndExit(errors.Wrapf(err, "failed to set Go environment variables"), false)
	}
	os.Exit(runGodelApp(os.Args))
}

func runGodelApp(osArgs []string) int {
	os.Args = osArgs

	global, err := godellauncher.ParseAppArgs(os.Args)
	tasksCfgInfo := config.TasksConfigInfo{
		BuiltinPluginsConfig: defaulttasks.BuiltinPluginsConfig(),
	}
	if err != nil {
		// match invalid flag output with that provided by Cobra CLI
		printErrAndExit(fmt.Errorf(err.Error()+"\n"+godellauncher.UsageString(createTasks(nil, nil, nil, tasksCfgInfo))), false)
	}

	var allUpgradeConfigTasks []godellauncher.UpgradeConfigTask
	var defaultTasks, pluginTasks []godellauncher.Task
	if global.Wrapper != "" {
		godelCfg, err := config.ReadGodelConfigFromProjectDir(path.Dir(global.Wrapper))
		if err != nil {
			printErrAndExit(err, global.Debug)
		}

		taskCfgProviders := config.TasksConfigProvidersConfig(godelCfg.TasksConfigProviders)
		configProvidersParam, err := taskCfgProviders.ToParam()
		if err != nil {
			printErrAndExit(err, global.Debug)
		}
		providedConfigs, err := plugins.LoadProvidedConfigurations(configProvidersParam, os.Stdout)
		if err != nil {
			printErrAndExit(err, global.Debug)
		}
		// combine base configuration with resolved configurations
		tasksConfig := config.TasksConfig(godelCfg.TasksConfig)
		tasksConfig.Combine(providedConfigs...)
		tasksCfgInfo.TasksConfig = tasksConfig

		// add default tasks
		defaultTasksCfg, err := defaulttasks.PluginsConfig(config.DefaultTasksConfig(tasksConfig.DefaultTasks))
		if err != nil {
			printErrAndExit(err, global.Debug)
		}
		defaultTasksParam, err := defaultTasksCfg.ToParam()
		if err != nil {
			printErrAndExit(err, global.Debug)
		}

		var defaultUpgradeConfigTasks, pluginUpgradeConfigTasks []godellauncher.UpgradeConfigTask

		tasksCfgInfo.DefaultTasksPluginsConfig = defaultTasksCfg
		defaultTasks, defaultUpgradeConfigTasks, err = plugins.LoadPluginsTasks(defaultTasksParam, os.Stdout)
		if err != nil {
			printErrAndExit(err, global.Debug)
		}

		// add tasks provided by plugins
		pluginsCfg := config.PluginsConfig(tasksConfig.Plugins)
		pluginsParam, err := pluginsCfg.ToParam()
		if err != nil {
			printErrAndExit(err, global.Debug)
		}
		pluginTasks, pluginUpgradeConfigTasks, err = plugins.LoadPluginsTasks(pluginsParam, os.Stdout)
		if err != nil {
			printErrAndExit(err, global.Debug)
		}

		if len(defaultTasksCfg.Plugins) != 0 && len(tasksConfig.Plugins.Plugins) != 0 {
			// verify that there are no conflicts
			combinedCfg := config.PluginsConfig(tasksConfig.Plugins)
			combinedCfg.DefaultResolvers = append(combinedCfg.DefaultResolvers, tasksConfig.Plugins.DefaultResolvers...)
			combinedCfg.Plugins = append(combinedCfg.Plugins, tasksConfig.Plugins.Plugins...)
			combinedParam, err := combinedCfg.ToParam()
			if err != nil {
				printErrAndExit(err, global.Debug)
			}
			if _, _, err := plugins.LoadPluginsTasks(combinedParam, ioutil.Discard); err != nil {
				printErrAndExit(err, global.Debug)
			}
		}

		// add all upgrade tasks
		allUpgradeConfigTasks = append(allUpgradeConfigTasks, defaulttasks.BuiltinUpgradeConfigTasks()...)
		allUpgradeConfigTasks = append(allUpgradeConfigTasks, defaultUpgradeConfigTasks...)
		allUpgradeConfigTasks = append(allUpgradeConfigTasks, pluginUpgradeConfigTasks...)
	}
	task, err := godellauncher.TaskForInput(global, createTasks(defaultTasks, pluginTasks, allUpgradeConfigTasks, tasksCfgInfo))
	if err != nil {
		// match missing command output with that provided by Cobra CLI
		errTmpl := "%s\nRun '%s --help' for usage."
		printErrAndExit(fmt.Errorf(errTmpl, err.Error(), godel.AppName), false)
	}

	if err := task.Run(global, os.Stdout); err != nil {
		// note that only app/amalgomated tasks will never reach this point, as they return an exit code and then
		// pass through an empty error. Those tasks are expected to handle their own error output.
		printErrAndExit(err, global.Debug)
	}
	return 0
}

func createTasks(defaultTasks, pluginTasks []godellauncher.Task, upgradeConfigTasks []godellauncher.UpgradeConfigTask, tasksCfgInfo config.TasksConfigInfo) []godellauncher.Task {
	var allTasks []godellauncher.Task
	allTasks = append(allTasks, builtintasks.Tasks(tasksCfgInfo)...)
	allTasks = append(allTasks, defaultTasks...)
	allTasks = append(allTasks, builtintasks.VerifyTask(append(allTasks, pluginTasks...)))
	allTasks = append(allTasks, builtintasks.UpgradeConfigTask(upgradeConfigTasks))
	allTasks = append(allTasks, pluginTasks...)
	return allTasks
}

func printErrAndExit(err error, debug bool) {
	if errStr := err.Error(); errStr != "" {
		if debug {
			errStr = errorstringer.StackWithInterleavedMessages(err)
		}
		fmt.Println("Error:", errStr)
	}
	os.Exit(1)
}
