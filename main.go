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
	"io"
	"os"
	"path"

	"github.com/kardianos/osext"
	"github.com/nmiyake/pkg/dirs"
	"github.com/nmiyake/pkg/errorstringer"
	"github.com/palantir/amalgomate/amalgomated"
	"github.com/pkg/errors"

	"github.com/palantir/godel/framework/apptasks"
	"github.com/palantir/godel/framework/builtintasks"
	"github.com/palantir/godel/framework/godel"
	"github.com/palantir/godel/framework/godellauncher"
	"github.com/palantir/godel/framework/plugins"
)

func main() {
	gödelPath, err := osext.Executable()
	if err != nil {
		printErrAndExit(errors.Wrapf(err, "failed to determine path for current executable"), false)
	}

	if err := dirs.SetGoEnvVariables(); err != nil {
		printErrAndExit(errors.Wrapf(err, "failed to set Go environment variables"), false)
	}

	cmdLib, err := apptasks.AmalgomatedCmdLib(gödelPath)
	if err != nil {
		printErrAndExit(errors.Wrapf(err, "failed to create amalgomated CmdLib"), false)
	}
	os.Exit(amalgomated.RunApp(os.Args, nil, cmdLib, runGodelApp))
}

func runGodelApp(osArgs []string) int {
	os.Args = osArgs

	global, err := godellauncher.ParseAppArgs(os.Args)
	if err != nil {
		// match invalid flag output with that provided by Cobra CLI
		printErrAndExit(fmt.Errorf(err.Error()+"\n"+godellauncher.UsageString(createTasks("", nil))), false)
	}

	var pluginTasks []godellauncher.Task
	if global.Wrapper != "" {
		// add tasks provided by plugins
		pluginTasks, err = createPluginTasks(global.Wrapper, os.Stdout)
		if err != nil {
			printErrAndExit(err, global.Debug)
		}
	}
	task, err := godellauncher.TaskForInput(global, createTasks(global.Wrapper, pluginTasks))
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

func createTasks(wrapperPath string, pluginTasks []godellauncher.Task) []godellauncher.Task {
	var allTasks []godellauncher.Task
	allTasks = append(allTasks, builtintasks.Tasks(wrapperPath)...)
	allTasks = append(allTasks, apptasks.AmalgomatedTasks()...)
	allTasks = append(allTasks, apptasks.AppTasks()...)
	allTasks = append(allTasks, builtintasks.VerifyTask(append(allTasks, pluginTasks...)))
	allTasks = append(allTasks, pluginTasks...)
	return allTasks
}

func createPluginTasks(wrapperPath string, stdout io.Writer) ([]godellauncher.Task, error) {
	cfgDir, err := godellauncher.ConfigDirPath(path.Dir(wrapperPath))
	if err != nil {
		return nil, err
	}
	cfg, err := godellauncher.ReadGodelConfig(cfgDir)
	if err != nil {
		return nil, err
	}
	pluginTasks, err := plugins.LoadPluginsTasks(cfg.Plugins, stdout)
	if err != nil {
		return nil, err
	}
	return pluginTasks, nil
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
