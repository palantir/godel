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
	"fmt"
	"io"
	"os"
	"os/exec"
	"path"
	"strings"
	"unicode"

	"github.com/pkg/errors"

	"github.com/palantir/godel/v2/framework/godellauncher"
)

// TaskInfo is a JSON-serializable interface that can be translated into a godellauncher.Task. Refer to that struct for
// field documentation.
type TaskInfo interface {
	Name() string
	Description() string
	Command() []string
	GlobalFlagOptions() GlobalFlagOptions
	VerifyOptions() VerifyOptions

	toTask(pluginExecPath, cfgFileName string, assets []string) godellauncher.Task
}

// taskInfoImpl is a concrete implementation of TaskInfo. Note that the functions are defined on non-pointer receivers
// to reduce bugs in calling functions in closures.
type taskInfoImpl struct {
	NameVar              string                 `json:"name"`
	DescriptionVar       string                 `json:"description"`
	CommandVar           []string               `json:"command"`
	GlobalFlagOptionsVar *globalFlagOptionsImpl `json:"globalFlagOptions"`
	VerifyOptionsVar     *verifyOptionsImpl     `json:"verifyOptions"`
}

type TaskInfoParam interface {
	apply(*taskInfoImpl)
}

type taskInfoParamFunc func(*taskInfoImpl)

func (f taskInfoParamFunc) apply(impl *taskInfoImpl) {
	f(impl)
}

func TaskInfoCommand(command ...string) TaskInfoParam {
	return taskInfoParamFunc(func(impl *taskInfoImpl) {
		impl.CommandVar = command
	})
}

func TaskInfoVerifyOptions(params ...VerifyOptionsParam) TaskInfoParam {
	return taskInfoParamFunc(func(impl *taskInfoImpl) {
		verifyOpts := newVerifyOptionsImpl(params...)
		var verifyImpls []verifyFlagImpl
		for _, v := range verifyOpts.VerifyTaskFlags() {
			verifyImpls = append(verifyImpls, verifyFlagImpl{
				NameVar:        v.Name(),
				DescriptionVar: v.Description(),
				TypeVar:        v.Type(),
			})
		}
		impl.VerifyOptionsVar = &verifyOptionsImpl{
			VerifyTaskFlagsVar: verifyImpls,
			OrderingVar:        verifyOpts.Ordering(),
			ApplyTrueArgsVar:   verifyOpts.ApplyTrueArgs(),
			ApplyFalseArgsVar:  verifyOpts.ApplyFalseArgs(),
		}
	})
}

func newTaskInfoImpl(name, description string, params ...TaskInfoParam) (taskInfoImpl, error) {
	for _, r := range name {
		if unicode.IsSpace(r) {
			return taskInfoImpl{}, errors.Errorf("task name cannot contain whitespace: %q", name)
		}
	}
	impl := taskInfoImpl{
		NameVar:        name,
		DescriptionVar: description,
	}
	for _, p := range params {
		if p == nil {
			continue
		}
		p.apply(&impl)
	}
	return impl, nil
}

func (ti taskInfoImpl) Name() string {
	return ti.NameVar
}

func (ti taskInfoImpl) Description() string {
	return ti.DescriptionVar
}

func (ti taskInfoImpl) Command() []string {
	return ti.CommandVar
}

func (ti taskInfoImpl) VerifyOptions() VerifyOptions {
	if ti.VerifyOptionsVar == nil {
		return nil
	}
	return ti.VerifyOptionsVar
}

func (ti taskInfoImpl) GlobalFlagOptions() GlobalFlagOptions {
	if ti.GlobalFlagOptionsVar == nil {
		return nil
	}
	return ti.GlobalFlagOptionsVar
}

func (ti taskInfoImpl) toTask(pluginExecPath, cfgFileName string, assets []string) godellauncher.Task {
	var verifyOpts *godellauncher.VerifyOptions
	if ti.VerifyOptions() != nil {
		opts := ti.VerifyOptionsVar.toGodelVerifyOptions()
		verifyOpts = &opts
	}
	var globalFlagOpts godellauncher.GlobalFlagOptions
	if ti.GlobalFlagOptionsVar != nil {
		globalFlagOpts = ti.GlobalFlagOptionsVar.toGodelGlobalFlagOptions()
	}
	return godellauncher.Task{
		Name:           ti.NameVar,
		Description:    ti.DescriptionVar,
		ConfigFile:     cfgFileName,
		Verify:         verifyOpts,
		GlobalFlagOpts: globalFlagOpts,
		RunImpl: func(t *godellauncher.Task, global godellauncher.GlobalConfig, stdout io.Writer) error {
			cmdArgs, err := globalFlagArgs(t.GlobalFlagOpts, t.ConfigFile, global)
			if err != nil {
				return err
			}
			// if assets are specified, provide as slice argument
			if len(assets) > 0 {
				cmdArgs = append(cmdArgs, "--"+AssetsFlagName)
				cmdArgs = append(cmdArgs, strings.Join(assets, ","))
			}
			cmdArgs = append(cmdArgs, ti.CommandVar...)
			cmdArgs = append(cmdArgs, global.TaskArgs...)
			cmd := exec.Command(pluginExecPath, cmdArgs...)
			cmd.Stdout = stdout
			cmd.Stderr = stdout
			cmd.Stdin = os.Stdin
			if err := cmd.Run(); err != nil {
				if _, ok := err.(*exec.ExitError); ok {
					// create empty error because command will likely print its own error
					return fmt.Errorf("")
				}
				return errors.Wrapf(err, "plugin execution failed")
			}
			return nil
		},
	}
}

func globalFlagArgs(globalFlagOpts godellauncher.GlobalFlagOptions, configFileName string, global godellauncher.GlobalConfig) ([]string, error) {
	var args []string
	if global.Debug && globalFlagOpts.DebugFlag != "" {
		args = append(args, globalFlagOpts.DebugFlag)
	}

	// the rest of the arguments depend on "--wrapper" being specified in the global configuration
	if global.Wrapper == "" {
		return args, nil
	}

	projectDir := path.Dir(global.Wrapper)
	if globalFlagOpts.ProjectDirFlag != "" {
		args = append(args, globalFlagOpts.ProjectDirFlag, projectDir)
	}

	// if config dir flags were not specified, nothing more to do
	if globalFlagOpts.GodelConfigFlag == "" && (globalFlagOpts.ConfigFlag == "" || configFileName == "") {
		return args, nil
	}

	cfgDir, err := godellauncher.ConfigDirPath(projectDir)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to determine config directory path")
	}

	if globalFlagOpts.GodelConfigFlag != "" {
		args = append(args, globalFlagOpts.GodelConfigFlag, path.Join(cfgDir, godellauncher.GodelConfigYML))
	}

	if globalFlagOpts.ConfigFlag != "" && configFileName != "" {
		args = append(args, globalFlagOpts.ConfigFlag, path.Join(cfgDir, configFileName))
	}

	return args, nil
}
