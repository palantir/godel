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
	"encoding/base64"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	"github.com/pkg/errors"

	"github.com/palantir/godel/framework/godellauncher"
)

// UpgradeConfigTaskInfo is a JSON-serializable interface that can be translated into a godellauncher.UpgradeConfigTask.
// Refer to that struct for field documentation.
type UpgradeConfigTaskInfo interface {
	Command() []string
	GlobalFlagOptions() GlobalFlagOptions
	// LegacyConfigFile returns the name of the legacy configuration file (the name of the configuration file used in
	// version 1 of gödel). Blank if this plugin does not have a legacy file or does not support upgrading legacy files.
	LegacyConfigFile() string
	toTask(pluginExecPath, cfgFileName string, assets []string) godellauncher.UpgradeConfigTask
}

type UpgradeConfigTaskInfoParam interface {
	apply(*upgradeConfigTaskInfoImpl)
}

type upgradeConfigTaskInfoParamFunc func(*upgradeConfigTaskInfoImpl)

func (f upgradeConfigTaskInfoParamFunc) apply(impl *upgradeConfigTaskInfoImpl) {
	f(impl)
}

func UpgradeConfigTaskInfoCommand(command ...string) UpgradeConfigTaskInfoParam {
	return upgradeConfigTaskInfoParamFunc(func(impl *upgradeConfigTaskInfoImpl) {
		impl.CommandVar = command
	})
}

func LegacyConfigFile(legacyConfigFile string) UpgradeConfigTaskInfoParam {
	return upgradeConfigTaskInfoParamFunc(func(impl *upgradeConfigTaskInfoImpl) {
		impl.LegacyConfigFile = legacyConfigFile
	})
}

// upgradeConfigTaskInfoImpl is a concrete implementation of UpgradeConfigTaskInfo. Note that the functions are defined
// on non-pointer receivers to reduce bugs in calling functions in closures.
type upgradeConfigTaskInfoImpl struct {
	// GroupID is the group ID of the plugin.
	GroupID string `json:"groupId"`
	// ProductID is the product ID of the plugin.
	ProductID string `json:"productId"`
	// CommandVar stores the commands to invoke to run the "upgrade-config" task
	CommandVar           []string               `json:"command"`
	LegacyConfigFile     string                 `json:"legacyConfigFile"`
	GlobalFlagOptionsVar *globalFlagOptionsImpl `json:"globalFlagOptions"`
}

func newUpgradeConfigTaskInfoImpl(params ...UpgradeConfigTaskInfoParam) *upgradeConfigTaskInfoImpl {
	if len(params) == 0 {
		return nil
	}
	impl := &upgradeConfigTaskInfoImpl{}
	for _, p := range params {
		if p == nil {
			continue
		}
		p.apply(impl)
	}
	return impl
}

func (ti upgradeConfigTaskInfoImpl) Command() []string {
	return ti.CommandVar
}

func (ti upgradeConfigTaskInfoImpl) GlobalFlagOptions() GlobalFlagOptions {
	if ti.GlobalFlagOptionsVar == nil {
		return nil
	}
	return ti.GlobalFlagOptionsVar
}

func (ti upgradeConfigTaskInfoImpl) toTask(pluginExecPath, cfgFileName string, assets []string) godellauncher.UpgradeConfigTask {
	var globalFlagOpts godellauncher.GlobalFlagOptions
	if ti.GlobalFlagOptionsVar != nil {
		globalFlagOpts = ti.GlobalFlagOptionsVar.toGodelGlobalFlagOptions()
	}
	return godellauncher.UpgradeConfigTask{
		ID:               ti.GroupID + ":" + ti.ProductID,
		ConfigFile:       cfgFileName,
		LegacyConfigFile: ti.LegacyConfigFile,
		GlobalFlagOpts:   globalFlagOpts,
		RunImpl: func(t *godellauncher.UpgradeConfigTask, global godellauncher.GlobalConfig, configBytes []byte, stdout io.Writer) ([]byte, error) {
			cmdArgs, err := globalFlagArgs(t.GlobalFlagOpts, t.ConfigFile, global)
			if err != nil {
				return nil, err
			}
			// if assets are specified, provide as slice argument
			if len(assets) > 0 {
				cmdArgs = append(cmdArgs, "--"+AssetsFlagName)
				cmdArgs = append(cmdArgs, strings.Join(assets, ","))
			}
			cmdArgs = append(cmdArgs, ti.CommandVar...)

			// add base64-encoded config as argument
			cmdArgs = append(cmdArgs, base64.StdEncoding.EncodeToString(configBytes))

			cmd := exec.Command(pluginExecPath, cmdArgs...)
			cmd.Stdin = os.Stdin
			outputBytes, err := cmd.CombinedOutput()
			if err != nil {
				output := string(outputBytes)
				if output == "" {
					output = fmt.Sprintf("command %v failed: %v", cmd.Args, err)
				} else {
					// clean up output for underlying command
					output = strings.TrimPrefix(output, "Error: ")
					output = strings.TrimSuffix(output, "\n")
				}
				if _, ok := err.(*exec.ExitError); ok {
					// if error was an exit error, don't bother wrapping because it's probably just "exit 1"
					return nil, errors.Errorf(output)
				}
				return nil, errors.Wrapf(err, output)
			}

			// valid output bytes are encoded as base64, so decode
			decoded, err := base64.StdEncoding.DecodeString(string(outputBytes))
			if err != nil {
				return nil, errors.Wrapf(err, "failed to decode base64 output")
			}
			return decoded, nil
		},
	}
}
