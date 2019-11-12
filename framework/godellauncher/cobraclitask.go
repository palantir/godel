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
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/palantir/godel/v2/framework/godel"
)

// CobraCLITask creates a new Task that runs the provided *cobra.Command. The runner for the task does the following:
//
// * Creates a "dummy" root cobra.Command with "godel" as the command name
// * Adds the provided command as a subcommand of the dummy root
// * Executes the root command with the following "os.Args":
//     [executable] [task] [task args...]
//
// The second argument is an optional pointer. If the pointer is non-nil, then the value of the provided pointer will be
// set to the GlobalConfig provided when the task is run.
func CobraCLITask(cmd *cobra.Command, globalConfigPtr *GlobalConfig) Task {
	rootCmd := CobraCmdToRootCmd(cmd)
	return Task{
		Name:        cmd.Use,
		Description: cmd.Short,
		RunImpl: func(t *Task, global GlobalConfig, stdout io.Writer) error {
			if globalConfigPtr != nil {
				*globalConfigPtr = global
			}
			rootCmd.SetOutput(stdout)
			args := []string{global.Executable}
			args = append(args, global.Task)
			args = append(args, global.TaskArgs...)
			os.Args = args
			return rootCmd.Execute()
		},
	}
}

// CobraCmdToRootCmd takes the provided *cobra.Command and returns a new *cobra.Command that acts as its "root" command.
// The root command has "godel" as its command name and is configured to silence the built-in Cobra error printing.
// However, it has custom logic to match the standard Cobra error output for unrecognized flags.
func CobraCmdToRootCmd(cmd *cobra.Command) *cobra.Command {
	rootCmd := &cobra.Command{
		Use: godel.AppName,
	}
	rootCmd.AddCommand(cmd)

	// Set custom error behavior for flag errors. Usage should only be printed if the error is due to invalid flags. In
	// order to do this, set SilenceErrors and SilenceUsage to true and set flag error function that prints error and
	// usage and then returns an error with empty content (so that the top-level handler does not print it).
	rootCmd.SilenceErrors = true
	rootCmd.SilenceUsage = true
	rootCmd.SetFlagErrorFunc(func(c *cobra.Command, err error) error {
		return fmt.Errorf("%s\n%s", err.Error(), strings.TrimSuffix(c.UsageString(), "\n"))
	})
	return rootCmd
}

func UnknownCommandError(cmd *cobra.Command, args []string) error {
	errTmpl := `unknown command "%s" for "%s"
Run '%v --help' for usage.`
	return fmt.Errorf(errTmpl, args[0], cmd.CommandPath(), cmd.CommandPath())
}
