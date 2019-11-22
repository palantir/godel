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

	"github.com/palantir/godel/v2/framework/godel"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

type flagDescProvider interface {
	addFlag(fset *pflag.FlagSet)
}

type stringFlagDesc struct {
	name, usage string
}

func (f stringFlagDesc) addFlag(fset *pflag.FlagSet) {
	fset.String(f.name, "", f.usage)
}

type boolFlagDesc struct {
	name, usage string
}

func (f boolFlagDesc) addFlag(fset *pflag.FlagSet) {
	fset.Bool(f.name, false, f.usage)
}

// helpFlagTask returns the task that should be executed for the "--help" flag. The launcher is a custom CLI, but for
// consistency it uses the Cobra framework to print its help output. In order to do so, the task creates a new Cobra CLI
// that has the same structure (top-level commands and global flags) as the launcher CLI and runs it with the "--help"
// flag.
func helpFlagTask(tasks []Task) Task {
	return Task{
		Name:        "help",
		Description: fmt.Sprintf("help for %s", godel.AppName),
		RunImpl: func(t *Task, global GlobalConfig, stdout io.Writer) error {
			os.Args = []string{global.Executable, "--help"}
			return helpCobraCmd(tasks).Execute()
		},
	}
}

// UsageString returns the usage string for the launcher application with the specified tasks. The returned string does
// not have a trailing newline.
func UsageString(tasks []Task) string {
	return strings.TrimSuffix(helpCobraCmd(tasks).UsageString(), "\n")
}

func helpCobraCmd(tasks []Task) *cobra.Command {
	rootCmd := &cobra.Command{
		Use: godel.AppName,
	}
	for _, t := range tasks {
		rootCmd.AddCommand(&cobra.Command{
			Use:   t.Name,
			Short: t.Description,
			// provide a non-nil (but no-op) command so that the command is listed in help
			Run: func(cmd *cobra.Command, args []string) {},
		})
	}
	for _, f := range globalFlags() {
		f.addFlag(rootCmd.Flags())
	}

	// set help command to be empty hidden command to effectively remove the built-in help command
	rootCmd.SetHelpCommand(&cobra.Command{
		Hidden: true,
	})

	return rootCmd
}

func globalFlags() []flagDescProvider {
	return []flagDescProvider{
		boolFlagDesc{
			name:  versionFlagTask().Name,
			usage: versionFlagTask().Description,
		},
		boolFlagDesc{
			name:  "debug",
			usage: "run in debug mode (print full stack traces on failures and include other debugging output)",
		},
		stringFlagDesc{
			name:  "wrapper",
			usage: "path to the wrapper script for this invocation",
		},
	}
}
