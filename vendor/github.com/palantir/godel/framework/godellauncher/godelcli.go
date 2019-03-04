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

package godellauncher

import (
	"fmt"
	"io"
	"strings"

	"github.com/pkg/errors"

	"github.com/palantir/godel/framework/godel"
)

// ParseAppArgs parses the arguments provided to the gödel launcher application into a GlobalConfig struct. Returns an
// error if the provided arguments are not legal. The provided args should be the arguments provided in os.Args. The
// arguments should match the following form:
//
// [executable] [<global flags>] [<task>] [<task flags/args>]
//
// <global flags> can be one of [--version], [--help|-h], [--debug] or [--wrapper <path>]. Note that, unlike the
// behavior of some other CLI programs, the flags can only be specified exactly as described: for example, inputs of the
// form "--version=true", "--debug false" and "--wrapper=<path>" are not valid.
func ParseAppArgs(args []string) (GlobalConfig, error) {
	// executable name must be specified
	if len(args) == 0 {
		return GlobalConfig{}, errors.Errorf("args cannot be empty")
	}

	var cfg GlobalConfig
	cfg.Executable = args[0]

	remainingArgs := args[1:]
	if len(remainingArgs) == 0 {
		// if only executable name is specified, treat as if help flag was specified
		cfg.Help = true
		return cfg, nil
	}

	for len(remainingArgs) > 0 {
		currArg := remainingArgs[0]
		remainingArgs = remainingArgs[1:]

		// treat "--" as ending flag interpretation
		if currArg == "--" {
			continue
		}

		// current argument is a flag
		if strings.HasPrefix(currArg, "-") {
			switch currArg {
			case "--version":
				cfg.Version = true
			case "--help", "-h":
				cfg.Help = true
			case "--debug":
				cfg.Debug = true
			case "--wrapper":
				if len(remainingArgs) == 0 {
					return GlobalConfig{}, errors.Errorf("flag '--wrapper' must specify a value")
				}
				currArg = remainingArgs[0]
				remainingArgs = remainingArgs[1:]
				cfg.Wrapper = currArg
			default:
				return GlobalConfig{}, errors.Errorf("unknown flag: %s", currArg)
			}
			continue
		}

		// first non-flag argument is the task, and everything following the task is treated as args for the task
		cfg.Task = currArg
		cfg.TaskArgs = remainingArgs
		break
	}

	return cfg, nil
}

// TaskForInput returns the Task that should be run based on the provided GlobalConfig.
//
// If the "Task" field of GlobalConfig is empty, it indicates that the launcher was run without any tasks specified. If
// that is the case, the following logic is used to determine the task to be returned:
//
// * If global.Help is true, the help output is printed
// * If global.Help is false and global.Version is true, the version is printed
// * If global.Help and global.Version are both false, the help output is printed
//
// If the "Task" field of GlobalConfig is non-empty, then the task in the provided "tasks" slice with the name that
// matches the "Task" field is returned.
//
// Returns an error if the "Task" field is non-empty but there is no corresponding task in the "tasks" slice, or if the
// provided "tasks" slice contains multiple entries with the same name.
func TaskForInput(global GlobalConfig, tasks []Task) (Task, error) {
	if global.Task == "" {
		if !global.Help && global.Version {
			return versionFlagTask(), nil
		}
		return helpFlagTask(tasks), nil
	}

	tasksMap := make(map[string]Task)
	for _, t := range tasks {
		if _, ok := tasksMap[t.Name]; ok {
			return Task{}, fmt.Errorf(`command "%s" defined multiple times`, t.Name)
		}
		tasksMap[t.Name] = t
	}

	task, ok := tasksMap[global.Task]
	if !ok {
		return Task{}, fmt.Errorf(`unknown command "%s" for "%s"`, global.Task, godel.AppName)
	}
	return task, nil
}

func versionFlagTask() Task {
	return Task{
		Name:        "version",
		Description: fmt.Sprintf("print %s version", godel.AppName),
		RunImpl: func(t *Task, global GlobalConfig, stdout io.Writer) error {
			fmt.Fprintln(stdout, godel.VersionOutput())
			return nil
		},
	}
}
