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
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"

	"github.com/palantir/godel/framework/godellauncher"
)

func InfoTask() godellauncher.Task {
	var globalCfg godellauncher.GlobalConfig
	cmd := &cobra.Command{
		Use:   "info",
		Short: "Print information regarding gödel",
	}
	cmd.AddCommand(&cobra.Command{
		Use:   "default-tasks",
		Short: "Print configuration for default tasks",
		RunE: func(cmd *cobra.Command, args []string) error {
			projectDir, err := globalCfg.ProjectDir()
			if err != nil {
				return err
			}
			godelCfg, err := godellauncher.ReadGodelConfigFromProjectDir(projectDir)
			if err != nil {
				return err
			}
			bytes, err := yaml.Marshal(godellauncher.DefaultTasksPluginsConfig(godelCfg.DefaultTasks))
			if err != nil {
				return errors.Wrapf(err, "failed to marshal default task config to JSON")
			}
			cmd.Print(string(bytes))
			return nil
		},
	})
	return godellauncher.CobraCLITask(cmd, &globalCfg)
}
