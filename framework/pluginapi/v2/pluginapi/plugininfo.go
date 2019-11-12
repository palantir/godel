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
	"strings"

	"github.com/pkg/errors"

	"github.com/palantir/godel/v2/framework/godellauncher"
	v1 "github.com/palantir/godel/v2/framework/pluginapi"
)

const (
	CurrentSchemaVersion  = "2"
	PluginInfoCommandName = v1.PluginInfoCommandName
)

// PluginInfo specifies the information for a plugin and the tasks that it provides.
type PluginInfo interface {
	// PluginSchemaVersion returns the schema version for the plugin.
	PluginSchemaVersion() string

	// The Group, Product and Version functions return the "group", "product" and "version" components of the plugin.
	// When colon-delimited, they form a Maven identifier. All must be non-empty.

	// Group returns the group of the plugin.
	Group() string
	// Product returns the product name of the plugin. Must have the suffix "-plugin".
	Product() string
	// Version returns the version of hte plugin.
	Version() string

	// UsesConfig returns true if this plugin uses configuration, false otherwise.
	UsesConfig() bool

	// Tasks returns the tasks provided by the plugin. Requires the path to the plugin executable and assets as input.
	Tasks(pluginExecPath string, assets []string) []godellauncher.Task

	// UpgradeConfigTask returns the task that upgrades the configuration for this plugin. Returns nil if the plugin
	// does not support upgrading configuration.
	UpgradeConfigTask(pluginExecPath string, assets []string) *godellauncher.UpgradeConfigTask

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
	if err := validateComponent("group", group); err != nil {
		return nil, err
	}
	if err := validateComponent("product", product); err != nil {
		return nil, err
	}
	if !strings.HasSuffix(product, "-plugin") {
		return nil, errors.Errorf(`product must end with "-plugin", was %q`, product)
	}
	if err := validateComponent("version", version); err != nil {
		return nil, err
	}

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
	if builder.upgradeConfigTask != nil {
		builder.upgradeConfigTask.GroupID = group
		builder.upgradeConfigTask.ProductID = product
		builder.upgradeConfigTask.GlobalFlagOptionsVar = builder.globalFlagOpts
	}

	taskNameMap := make(map[string]struct{})
	for _, task := range builder.tasks {
		if _, ok := taskNameMap[task.Name()]; ok {
			return nil, errors.Errorf(`plugin %s specifies multiple tasks with name "%s"`, id, task.Name())
		}
		taskNameMap[task.Name()] = struct{}{}
	}

	if builder.upgradeConfigTask != nil && !builder.usesConfigFile {
		return nil, errors.Errorf(`plugin %s provides a configuration upgrade task but does not specify that it uses configuration`, id)
	}

	return pluginInfoImpl{
		PluginSchemaVersionVar: CurrentSchemaVersion,
		GroupVar:               group,
		ProductVar:             product,
		VersionVar:             version,
		UsesConfigVar:          builder.usesConfigFile,
		TasksVar:               builder.tasks,
		UpgradeConfigTaskVar:   builder.upgradeConfigTask,
	}, nil
}

func validateComponent(name, val string) error {
	if val == "" {
		return errors.Errorf("%s must be non-empty", name)
	}
	if strings.Contains(val, ":") {
		return errors.Errorf("%s cannot contain a ':', was %q", name, val)
	}
	return nil
}

type pluginInfoBuilder struct {
	usesConfigFile    bool
	tasks             []taskInfoImpl
	globalFlagOpts    *globalFlagOptionsImpl
	upgradeConfigTask *upgradeConfigTaskInfoImpl
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

func PluginInfoUpgradeConfigTaskInfo(params ...UpgradeConfigTaskInfoParam) PluginInfoParam {
	return pluginInfoParamFunc(func(impl *pluginInfoBuilder) error {
		if len(params) == 0 {
			return errors.Errorf("at least one param must be provided for upgrade configuration")
		}
		impl.upgradeConfigTask = newUpgradeConfigTaskInfoImpl(params...)
		return nil
	})
}

// pluginInfoImpl is a concrete implementation of Info. Note that the functions are defined on non-pointer receivers to reduce
// bugs in calling functions in closures.
type pluginInfoImpl struct {
	PluginSchemaVersionVar string `json:"pluginSchemaVersion"`
	// The Maven group identifier for the plugin. Must be non-empty.
	GroupVar string `json:"group"`
	// The Maven product identifier for the plugin. Must be non-empty and must have "-plugin" as its suffix.
	ProductVar string `json:"product"`
	// The Maven version identifier for the plugin. Must be non-empty.
	VersionVar string `json:"version"`
	// True if this plugin uses configuration, false otherwise.
	UsesConfigVar bool `json:"usesConfig"`
	// The tasks provided by the plugin.
	TasksVar []taskInfoImpl `json:"tasks"`
	// The configuration upgrade task provided by the plugin.
	UpgradeConfigTaskVar *upgradeConfigTaskInfoImpl `json:"upgradeTask"`
}

func (infoImpl pluginInfoImpl) PluginSchemaVersion() string {
	return infoImpl.PluginSchemaVersionVar
}

func (infoImpl pluginInfoImpl) Group() string {
	return infoImpl.GroupVar
}

func (infoImpl pluginInfoImpl) Product() string {
	return infoImpl.ProductVar
}

func (infoImpl pluginInfoImpl) Version() string {
	return infoImpl.VersionVar
}

func (infoImpl pluginInfoImpl) UsesConfig() bool {
	return infoImpl.UsesConfigVar
}

func (infoImpl pluginInfoImpl) configFileName() string {
	var configFileName string
	if infoImpl.UsesConfigVar {
		configFileName = infoImpl.ProductVar + ".yml"
	}
	return configFileName
}

func (infoImpl pluginInfoImpl) Tasks(pluginExecPath string, assets []string) []godellauncher.Task {
	var tasks []godellauncher.Task
	for _, ti := range infoImpl.TasksVar {
		tasks = append(tasks, ti.toTask(pluginExecPath, infoImpl.configFileName(), assets))
	}
	return tasks
}

func (infoImpl pluginInfoImpl) UpgradeConfigTask(pluginExecPath string, assets []string) *godellauncher.UpgradeConfigTask {
	if infoImpl.UpgradeConfigTaskVar == nil {
		return nil
	}
	taskVar := infoImpl.UpgradeConfigTaskVar.toTask(pluginExecPath, infoImpl.configFileName(), assets)
	return &taskVar
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

	var schemaVersionVar v1.SchemaVersion
	if err := json.Unmarshal(bytes, &schemaVersionVar); err != nil {
		return nil, errors.Wrapf(err, "failed to unmarshal plugin schema version for plugin %s from output %q", pluginPath, string(bytes))
	}

	var pluginInfo PluginInfo
	switch version := schemaVersionVar.PluginSchemaVersionVar; version {
	case v1.CurrentSchemaVersion:
		v1Info, err := v1.InfoFromBytes(bytes)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to unmarshal v1 plugin information")
		}

		// do baseline check/validation
		parts := strings.Split(v1Info.ID(), ":")
		if len(parts) != 3 {
			return nil, errors.Wrapf(err, "v1 plugin information provided invalid ID: %q", v1Info.ID())
		}
		if _, err := NewPluginInfo(parts[0], parts[1], parts[2]); err != nil {
			return nil, errors.Wrapf(err, "could not create v2 plugin info from v1 plugin info")
		}
		pluginInfo = wrappedV1PluginInfoImpl{
			v1PluginInfo: v1Info,
		}
	case CurrentSchemaVersion:
		var v2Info pluginInfoImpl
		if err := json.Unmarshal(bytes, &v2Info); err != nil {
			return nil, errors.Wrapf(err, "failed to unmarshal plugin information for plugin %s from output %q", pluginPath, string(bytes))
		}
		pluginInfo = v2Info
	default:
		return nil, errors.Errorf("unsupported plugin schema version: %s", version)
	}

	if !strings.HasSuffix(pluginInfo.Product(), "-plugin") {
		return nil, errors.Errorf(`plugin %s has an invalid product name: product names must have suffix "-plugin", but was %q`, pluginPath, pluginInfo.Product())
	}
	return pluginInfo, nil
}

type wrappedV1PluginInfoImpl struct {
	v1PluginInfo v1.PluginInfo
}

func (infoImpl wrappedV1PluginInfoImpl) PluginSchemaVersion() string {
	return CurrentSchemaVersion
}

func (infoImpl wrappedV1PluginInfoImpl) Group() string {
	return strings.Split(infoImpl.v1PluginInfo.ID(), ":")[0]
}

func (infoImpl wrappedV1PluginInfoImpl) Product() string {
	return strings.Split(infoImpl.v1PluginInfo.ID(), ":")[1]
}

func (infoImpl wrappedV1PluginInfoImpl) Version() string {
	return strings.Split(infoImpl.v1PluginInfo.ID(), ":")[2]
}

func (infoImpl wrappedV1PluginInfoImpl) UsesConfig() bool {
	return infoImpl.v1PluginInfo.ConfigFileName() != ""
}

func (infoImpl wrappedV1PluginInfoImpl) Tasks(pluginExecPath string, assets []string) []godellauncher.Task {
	return infoImpl.v1PluginInfo.Tasks(pluginExecPath, assets)
}

func (infoImpl wrappedV1PluginInfoImpl) UpgradeConfigTask(pluginExecPath string, assets []string) *godellauncher.UpgradeConfigTask {
	return nil
}

func (infoImpl wrappedV1PluginInfoImpl) private() {}
