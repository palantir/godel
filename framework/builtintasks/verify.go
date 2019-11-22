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
	"fmt"
	"io"
	"os"
	"sort"
	"strings"

	"github.com/palantir/godel/v2/framework/godellauncher"
	"github.com/spf13/cobra"
)

func VerifyTask(tasks []godellauncher.Task) godellauncher.Task {
	const (
		verifyCmdName = "verify"
		apply         = "apply"
	)

	verifyTasks, verifyTaskFlags := extractVerifyTasks(tasks)
	skipVerifyTasks := make(map[string]*bool)
	verifyTaskFlagVals := make(map[string]map[godellauncher.VerifyFlag]interface{})

	cmd := &cobra.Command{
		Use:   verifyCmdName,
		Short: "Run verify tasks for project",
	}

	applyVar := cmd.Flags().Bool(apply, true, "apply changes when possible")
	for _, task := range verifyTasks {
		skipVerifyTasks[task.Name] = cmd.Flags().Bool("skip-"+task.Name, false, fmt.Sprintf("skip '%s' task", task.Name))

		flags, ok := verifyTaskFlags[task.Name]
		if !ok {
			continue
		}
		if len(flags) == 0 {
			continue
		}

		// hook up task-specific flags
		verifyTaskFlagVals[task.Name] = make(map[godellauncher.VerifyFlag]interface{})
		for _, f := range flags {
			flagVal, err := f.AddFlag(cmd.Flags())
			if err != nil {
				panic(err)
			}
			verifyTaskFlagVals[task.Name][f] = flagVal
		}
	}

	cmd.SilenceErrors = true
	cmd.SilenceUsage = true
	return godellauncher.Task{
		Name:        cmd.Use,
		Description: cmd.Short,
		RunImpl: func(t *godellauncher.Task, global godellauncher.GlobalConfig, stdout io.Writer) error {
			args := []string{global.Executable}
			args = append(args, global.Task)
			args = append(args, global.TaskArgs...)
			os.Args = args

			cmd.RunE = func(cmd *cobra.Command, args []string) error {
				var failedChecks []string
				for _, task := range verifyTasks {
					// skip the task
					if *skipVerifyTasks[task.Name] {
						continue
					}

					var taskFlagArgs []string
					if *applyVar {
						taskFlagArgs = append(taskFlagArgs, task.Verify.ApplyTrueArgs...)
					} else {
						taskFlagArgs = append(taskFlagArgs, task.Verify.ApplyFalseArgs...)
					}

					// get task-specific flag values
					for _, f := range verifyTaskFlags[task.Name] {
						flagArgs, err := f.ToFlagArgs(verifyTaskFlagVals[task.Name][f])
						if err != nil {
							panic(err)
						}
						taskFlagArgs = append(taskFlagArgs, flagArgs...)
					}

					taskGlobal := global
					taskGlobal.Task = task.Name
					taskGlobal.TaskArgs = taskFlagArgs

					_, _ = fmt.Fprintf(stdout, "Running %s...\n", task.Name)
					if err := task.Run(taskGlobal, stdout); err != nil {
						var applyArgs []string
						if *applyVar {
							applyArgs = task.Verify.ApplyTrueArgs
						} else {
							applyArgs = task.Verify.ApplyFalseArgs
						}
						nameWithFlag := strings.Join(append([]string{task.Name}, applyArgs...), " ")
						failedChecks = append(failedChecks, nameWithFlag)
					}
				}

				if len(failedChecks) != 0 {
					msgParts := []string{"Failed tasks:"}
					for _, check := range failedChecks {
						msgParts = append(msgParts, "\t"+check)
					}
					_, _ = fmt.Fprintln(stdout, strings.Join(msgParts, "\n"))
					return fmt.Errorf("")
				}
				return nil
			}

			rootCmd := godellauncher.CobraCmdToRootCmd(cmd)
			rootCmd.SetOutput(stdout)
			return rootCmd.Execute()
		},
	}
}

func extractVerifyTasks(tasks []godellauncher.Task) ([]godellauncher.Task, map[string][]godellauncher.VerifyFlag) {
	var verifyTasks []godellauncher.Task
	verifyTaskFlags := make(map[string][]godellauncher.VerifyFlag)
	for _, task := range tasks {
		if task.Verify == nil {
			continue
		}
		verifyTasks = append(verifyTasks, task)
		verifyTaskFlags[task.Name] = task.Verify.VerifyTaskFlags
	}
	sort.SliceStable(verifyTasks, func(i, j int) bool {
		return verifyTasks[i].Verify.Ordering < verifyTasks[j].Verify.Ordering
	})
	return verifyTasks, verifyTaskFlags
}
