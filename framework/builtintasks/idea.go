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

package builtintasks

import (
	"github.com/palantir/godel/v2/framework/builtintasks/idea"
	"github.com/palantir/godel/v2/framework/godellauncher"
	"github.com/spf13/cobra"
)

func IDEATask() godellauncher.Task {
	const intellijCmdUsage = "Create IntelliJ project files for this project"
	var globalCfg godellauncher.GlobalConfig

	ideaCmd := &cobra.Command{
		Use:   "idea",
		Short: intellijCmdUsage,
		RunE: func(cmd *cobra.Command, args []string) error {
			projectDir, err := globalCfg.ProjectDir()
			if err != nil {
				return err
			}

			// This command has subcommands, but does not accept any arguments itself. If the execution has reached this
			// point and "args" is non-empty, treat it as an unknown command (rather than just executing this command
			// and ignoring the extra arguments). Avoids executing the wrong command on a subcommand typo.
			if len(args) > 0 {
				return godellauncher.UnknownCommandError(cmd, args)
			}
			return idea.CreateIntelliJFiles(projectDir)
		},
	}
	goglandSubcommand := &cobra.Command{
		Use:   "gogland",
		Short: "Create Gogland project files for this project",
		RunE: func(cmd *cobra.Command, args []string) error {
			projectDir, err := globalCfg.ProjectDir()
			if err != nil {
				return err
			}
			return idea.CreateGoglandFiles(projectDir)
		},
	}
	intelliJSubcommand := &cobra.Command{
		Use:   "intellij",
		Short: intellijCmdUsage,
		RunE: func(cmd *cobra.Command, args []string) error {
			projectDir, err := globalCfg.ProjectDir()
			if err != nil {
				return err
			}
			return idea.CreateGoglandFiles(projectDir)
		},
	}
	cleanSubcommand := &cobra.Command{
		Use:   "clean",
		Short: "Remove the IDEA project files for this project",
		RunE: func(cmd *cobra.Command, args []string) error {
			projectDir, err := globalCfg.ProjectDir()
			if err != nil {
				return err
			}
			return idea.CleanIDEAFiles(projectDir)
		},
	}

	ideaCmd.AddCommand(
		goglandSubcommand,
		intelliJSubcommand,
		cleanSubcommand,
	)
	return godellauncher.CobraCLITask(ideaCmd, &globalCfg)
}
