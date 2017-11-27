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

package pluginapi

import (
	"encoding/json"
	"fmt"
	"os/exec"

	"github.com/pkg/errors"

	"github.com/palantir/godel/framework/godellauncher"
)

const (
	CurrentSchemaVersion = "1"
	InfoCommandName      = "_godelPluginInfo"
)

// Info specifies the information for a plugin and the tasks that it provides.
type Info interface {
	// PluginSchemaVersion returns the schema version for the plugin.
	PluginSchemaVersion() string
	// ID returns the identifier for a plugin and is a string of the form "group:product:version".
	ID() string
	// ConfigFileName returns the name of the configuration file used by the plugin.
	ConfigFileName() string
	// Tasks returns the tasks provided by the plugin. Requires the path to the plugin executable and assets as input.
	Tasks(pluginExecPath string, assets []string) []godellauncher.Task

	// private function on interface to keep implementation private to package.
	private()
}

// MustNewInfo returns the result of calling NewInfo with the provided parameters. Panics if the call to NewInfo returns
// an error, so this function should only be used when the inputs are static and known to be valid.
func MustNewInfo(group, product, version, configFileName string, tasks ...TaskInfo) Info {
	cmd, err := NewInfo(group, product, version, configFileName, tasks...)
	if err != nil {
		panic(errors.Wrapf(err, "failed to create Info for plugin"))
	}
	return cmd
}

// NewInfo creates a new Info for the plugin using the provided configuration. Returns an error if the provided
// configuration is not valid (can occur if multiple tasks have the same name).
func NewInfo(group, product, version, configFileName string, tasks ...TaskInfo) (Info, error) {
	id := fmt.Sprintf("%s:%s:%s", group, product, version)
	taskNameMap := make(map[string]struct{})
	for _, task := range tasks {
		if _, ok := taskNameMap[task.Name()]; ok {
			return nil, errors.Errorf(`plugin %s specifies multiple tasks with name "%s"`, id, task.Name())
		}
		taskNameMap[task.Name()] = struct{}{}
	}

	var taskInfoImpls []taskInfoImpl
	for _, task := range tasks {
		var globalFlagOpts *globalFlagOptionsImpl
		if task.GlobalFlagOptions() != nil {
			globalFlagOpts = &globalFlagOptionsImpl{
				DebugFlagVar:       task.GlobalFlagOptions().DebugFlag(),
				ProjectDirFlagVar:  task.GlobalFlagOptions().ProjectDirFlag(),
				GodelConfigFlagVar: task.GlobalFlagOptions().GodelConfigFlag(),
				ConfigFlagVar:      task.GlobalFlagOptions().ConfigFlag(),
			}
		}
		var verifyOptsImpl *verifyOptionsImpl
		if task.VerifyOptions() != nil {
			var verifyFlags []verifyFlagImpl
			for _, f := range task.VerifyOptions().VerifyTaskFlags() {
				verifyFlags = append(verifyFlags, verifyFlagImpl{
					NameVar:        f.Name(),
					DescriptionVar: f.Description(),
					TypeVar:        f.Type(),
				})
			}
			verifyOptsImpl = &verifyOptionsImpl{
				VerifyTaskFlagsVar: verifyFlags,
				OrderingVar:        task.VerifyOptions().Ordering(),
				ApplyTrueArgsVar:   task.VerifyOptions().ApplyTrueArgs(),
				ApplyFalseArgsVar:  task.VerifyOptions().ApplyFalseArgs(),
			}
		}
		taskInfoImpls = append(taskInfoImpls, taskInfoImpl{
			NameVar:              task.Name(),
			DescriptionVar:       task.Description(),
			CommandVar:           task.Command(),
			GlobalFlagOptionsVar: globalFlagOpts,
			VerifyOptionsVar:     verifyOptsImpl,
		})
	}

	return &infoImpl{
		PluginSchemaVersionVar: CurrentSchemaVersion,
		IDVar:             id,
		ConfigFileNameVar: configFileName,
		TasksVar:          taskInfoImpls,
	}, nil
}

type infoImpl struct {
	PluginSchemaVersionVar string `json:"pluginSchemaVersion"`
	// ID is the identifier for a plugin and is a string of the form "group:product:version".
	IDVar string `json:"id"`
	// The name of the configuration file used by the plugin.
	ConfigFileNameVar string `json:"configFileName"`
	// The tasks provided by the plugin.
	TasksVar []taskInfoImpl `json:"tasks"`
}

func (infoImpl *infoImpl) PluginSchemaVersion() string {
	return infoImpl.PluginSchemaVersionVar
}

func (infoImpl *infoImpl) ID() string {
	return infoImpl.IDVar
}

func (infoImpl *infoImpl) ConfigFileName() string {
	return infoImpl.ConfigFileNameVar
}

func (infoImpl *infoImpl) Tasks(pluginExecPath string, assets []string) []godellauncher.Task {
	var tasks []godellauncher.Task
	for _, ti := range infoImpl.TasksVar {
		tasks = append(tasks, ti.toTask(pluginExecPath, infoImpl.ConfigFileNameVar, assets))
	}
	return tasks
}

func (infoImpl *infoImpl) private() {}

// InfoFromPlugin returns the Info for the plugin at the specified path. Does so by invoking the InfoCommand on the
// plugin and parsing the output.
func InfoFromPlugin(pluginPath string) (Info, error) {
	cmd := exec.Command(pluginPath, InfoCommandName)
	bytes, err := cmd.CombinedOutput()
	if err != nil {
		return nil, errors.Wrapf(err, "command %v failed.\nError:\n%v\nOutput:\n%s\n", cmd.Args, err, string(bytes))
	}

	var info infoImpl
	if err := json.Unmarshal(bytes, &info); err != nil {
		return nil, errors.Wrapf(err, "failed to unmarshal plugin information for plugin %s from output %q", pluginPath, string(bytes))
	}
	return &info, nil
}
