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
	"io"

	"github.com/pkg/errors"
	"github.com/spf13/pflag"
)

type Task struct {
	// The name of the command. This is used as the task command for invocation and should not contain any whitespace.
	Name string

	// The description for this task. Should be suitable to use as the command description in CLI help.
	Description string

	// The name of the configuration file for the task ("task.yml", etc.). Can be blank if the task does not require
	// file-based configuration.
	ConfigFile string

	// Verify stores the option for the "--verify" task. If non-nil, this command is run as part of the "verify" task.
	Verify *VerifyOptions

	// The runner that is invoked to run this task. Should be possible to run in-process (that is, this function should
	// not call os.Exit or equivalent).
	RunImpl func(t *Task, global GlobalConfig, stdout io.Writer) error
}

type VerifyOptions struct {
	// VerifyTaskFlags stores the task-specific flags supported by this verify task.
	VerifyTaskFlags []VerifyFlag
	// Ordering stores the weighting/ordering of the task as it will be run in the verify task.
	Ordering int
	// VerifyArgs specifies the arguments (typically flags) that will be provided to to this verify task when apply mode
	// is false. For example, []string{"--verify"} or []string{"-l"}.
	VerifyArgs []string
}

type VerifyFlag struct {
	Name        string
	Description string
	Type        FlagType
}

// AddFlag adds the flag represented by VerifyFlag to the specified pflag.FlagSet. Returns the pointer to the value that
// can be used to retrieve the value.
func (f VerifyFlag) AddFlag(fset *pflag.FlagSet) (interface{}, error) {
	switch f.Type {
	case StringFlag:
		return fset.String(f.Name, "", f.Description), nil
	default:
		return nil, errors.Errorf("unrecognized flag type: %v", f.Type)
	}
}

// ToFlagArgs takes the input parameter (which should be the value returned by calling AddFlag for the receiver) and
// returns a string slice that reconstructs the flag arguments for the given flag.
func (f VerifyFlag) ToFlagArgs(flagVal interface{}) ([]string, error) {
	switch f.Type {
	case StringFlag:
		flagValStr := flagVal.(*string)
		if flagValStr == nil || len(*flagValStr) == 0 {
			return nil, nil
		}
		return []string{"--" + f.Name, *flagValStr}, nil
	default:
		return nil, errors.Errorf("unrecognized flag type: %v", f.Type)
	}
}

// FlagType represents the type of a flag (string, boolean, etc.). Currently only string flags are supported.
type FlagType int

const (
	StringFlag FlagType = iota
)

func (t *Task) Run(global GlobalConfig, stdout io.Writer) error {
	return t.RunImpl(t, global, stdout)
}
