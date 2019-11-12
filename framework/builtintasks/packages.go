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
	"strings"

	"github.com/spf13/cobra"

	"github.com/palantir/godel/v2/framework/builtintasks/packages"
	"github.com/palantir/godel/v2/framework/godel/config"
	"github.com/palantir/godel/v2/framework/godellauncher"
)

func PackagesTask() godellauncher.Task {
	var globalCfg godellauncher.GlobalConfig
	return godellauncher.CobraCLITask(&cobra.Command{
		Use:   "packages",
		Short: "Lists all of the packages in the project except those excluded by configuration",
		RunE: func(cmd *cobra.Command, args []string) error {
			projectDir, err := globalCfg.ProjectDir()
			if err != nil {
				return err
			}

			cfg, err := config.ReadGodelConfigFromProjectDir(projectDir)
			if err != nil {
				return err
			}
			pkgs, err := packages.List(cfg.Exclude.Matcher(), projectDir)
			if err != nil {
				return err
			}
			_, _ = fmt.Fprintln(cmd.OutOrStdout(), strings.Join(pkgs, "\n"))
			return nil
		},
	}, &globalCfg)
}
