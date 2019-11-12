// Copyright 2019 Palantir Technologies, Inc.
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
	"os"
	"os/exec"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/palantir/godel/v2/framework/godellauncher"
)

func ExecTask() godellauncher.Task {
	var globalCfg godellauncher.GlobalConfig
	return godellauncher.CobraCLITask(&cobra.Command{
		Use:   "exec",
		Short: "Executes given shell command using godel",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				return errors.Errorf("no command specified")
			}
			execCmd := exec.Command(args[0], args[1:]...)
			execCmd.Stdout = cmd.OutOrStdout()
			execCmd.Stderr = cmd.OutOrStderr()
			execCmd.Stdin = os.Stdin
			return execCmd.Run()
		},
	}, &globalCfg)
}
