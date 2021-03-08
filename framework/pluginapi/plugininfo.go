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

	"github.com/palantir/godel/v2/framework/godellauncher"
	"github.com/pkg/errors"
)

const (
	CurrentSchemaVersion  = "1"
	PluginInfoCommandName = "_godelPluginInfo"
)

// PluginInfo specifies the information for a plugin and the tasks that it provides.
type PluginInfo interface {
	// PluginSchemaVersion returns the schema version for the plugin.
	PluginSchemaVersion() string
	// ID returns the identifier for a plugin and is a string of the form "group:product:version".
	ID() string
	// ConfigFileName returns the name of the configuration file used by the plugin. Empty string means that no
	// configuration file is used.
	ConfigFileName() string
	// Tasks returns the tasks provided by the plugin. Requires the path to the plugin executable and assets as input.
	Tasks(pluginExecPath string, assets []string) []godellauncher.Task

	// private function on interface to keep implementation private to package.
	private()
}

// MustNewPluginInfo returns the result of calling NewInfo with the provided parameters. Panics if the call to
// NewPluginInfo returns an error, so this function should only be used when the inputs are static and known to be valid.
func MustNewPluginInfo(group, product, version string, params ...PluginInfoParam) PluginInfo {
	cmd, err := NewPluginInfo(group, product, version, params...)
	if err != nil {
		panic(errors.Wrapf(err, "failed to create plugin info"))
	}
	return cmd
}

// NewPluginInfo creates a new PluginInfo for the plugin using the provided configuration. Returns an error if the
// provided configuration is not valid (can occur if multiple tasks have the same name or if a task is not valid).
func NewPluginInfo(group, product, version string, params ...PluginInfoParam) (PluginInfo, error) {
	id := fmt.Sprintf("%s:%s:%s", group, product, version)

	builder := pluginInfoBuilder{}
	for _, param := range params {
		if err := param.apply(&builder); err != nil {
			return nil, err
		}
	}

	// set global flag options based on builder
	for i := range builder.tasks {
		builder.tasks[i].GlobalFlagOptionsVar = builder.globalFlagOpts
	}

	taskNameMap := make(map[string]struct{})
	for _, task := range builder.tasks {
		if _, ok := taskNameMap[task.Name()]; ok {
			return nil, errors.Errorf(`plugin %s specifies multiple tasks with name "%s"`, id, task.Name())
		}
		taskNameMap[task.Name()] = struct{}{}
	}

	var configFileName string
	if builder.usesConfigFile {
		configFileName = product + ".yml"
	}

	return pluginInfoImpl{
		PluginSchemaVersionVar: CurrentSchemaVersion,
		IDVar:                  id,
		ConfigFileNameVar:      configFileName,
		TasksVar:               builder.tasks,
	}, nil
}

type pluginInfoBuilder struct {
	usesConfigFile bool
	tasks          []taskInfoImpl
	globalFlagOpts *globalFlagOptionsImpl
}

type PluginInfoParam interface {
	apply(*pluginInfoBuilder) error
}

type pluginInfoParamFunc func(*pluginInfoBuilder) error

func (f pluginInfoParamFunc) apply(impl *pluginInfoBuilder) error {
	return f(impl)
}

func PluginInfoUsesConfigFile() PluginInfoParam {
	return pluginInfoParamFunc(func(impl *pluginInfoBuilder) error {
		impl.usesConfigFile = true
		return nil
	})
}

func PluginInfoTaskInfo(name, description string, params ...TaskInfoParam) PluginInfoParam {
	return pluginInfoParamFunc(func(impl *pluginInfoBuilder) error {
		taskInfoImpl, err := newTaskInfoImpl(name, description, params...)
		if err != nil {
			return err
		}
		impl.tasks = append(impl.tasks, taskInfoImpl)
		return nil
	})
}

func PluginInfoGlobalFlagOptions(params ...GlobalFlagOptionsParam) PluginInfoParam {
	return pluginInfoParamFunc(func(impl *pluginInfoBuilder) error {
		if len(params) == 0 {
			return errors.Errorf("at least one param must be provided for global flag options configuration")
		}
		impl.globalFlagOpts = newGlobalFlagOptionsImpl(params...)
		return nil
	})
}

// pluginInfoImpl is a concrete implementation of Info. Note that the functions are defined on non-pointer receivers to reduce
// bugs in calling functions in closures.
type pluginInfoImpl struct {
	PluginSchemaVersionVar string `json:"pluginSchemaVersion"`
	// ID is the identifier for a plugin and is a string of the form "group:product:version".
	IDVar string `json:"id"`
	// The name of the configuration file used by the plugin.
	ConfigFileNameVar string `json:"configFileName"`
	// The tasks provided by the plugin.
	TasksVar []taskInfoImpl `json:"tasks"`
}

func (infoImpl pluginInfoImpl) PluginSchemaVersion() string {
	return infoImpl.PluginSchemaVersionVar
}

func (infoImpl pluginInfoImpl) ID() string {
	return infoImpl.IDVar
}

func (infoImpl pluginInfoImpl) ConfigFileName() string {
	return infoImpl.ConfigFileNameVar
}

func (infoImpl pluginInfoImpl) Tasks(pluginExecPath string, assets []string) []godellauncher.Task {
	var tasks []godellauncher.Task
	for _, ti := range infoImpl.TasksVar {
		tasks = append(tasks, ti.toTask(pluginExecPath, infoImpl.ConfigFileNameVar, assets))
	}
	return tasks
}

func (infoImpl pluginInfoImpl) private() {}

// InfoFromPlugin returns the Info for the plugin at the specified path. Does so by invoking the InfoCommand on the
// plugin and parsing the output.
func InfoFromPlugin(pluginPath string) (PluginInfo, error) {
	cmd := exec.Command(pluginPath, PluginInfoCommandName)
	bytes, err := cmd.CombinedOutput()
	if err != nil {
		return nil, errors.Wrapf(err, "command %v failed.\nError:\n%v\nOutput:\n%s\n", cmd.Args, err, string(bytes))
	}

	info, err := InfoFromBytes(bytes)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to unmarshal plugin information for plugin %s", pluginPath)
	}
	return info, nil
}

func InfoFromBytes(infoBytes []byte) (PluginInfo, error) {
	var info pluginInfoImpl
	if err := json.Unmarshal(infoBytes, &info); err != nil {
		return nil, errors.Wrapf(err, "failed to unmarshal plugin information from output %q", string(infoBytes))
	}
	return info, nil
}
