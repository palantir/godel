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
	"time"

	"github.com/spf13/cobra"

	"github.com/palantir/godel/framework/builtintasks/installupdate"
	"github.com/palantir/godel/framework/godellauncher"
)

func UpdateTask() godellauncher.Task {
	var (
		syncFlagVal              bool
		versionFlagVal           string
		checksumFlagVal          string
		cacheDurationFlagVal     time.Duration
		skipUpgradeConfigFlagVal bool
		globalCfg                godellauncher.GlobalConfig
	)

	cmd := &cobra.Command{
		Use:   "update",
		Short: "Update gödel for project",
		RunE: func(cmd *cobra.Command, args []string) error {
			projectDir, err := globalCfg.ProjectDir()
			if err != nil {
				return err
			}

			action := func() error {
				if syncFlagVal {
					// if sync flag is true, update version to what is specified in gödel.yml
					pkgSrc, err := installupdate.GodelPropsDistPkgInfo(projectDir)
					if err != nil {
						return err
					}
					if err := installupdate.Update(projectDir, pkgSrc, cmd.OutOrStdout()); err != nil {
						return err
					}
				} else {
					if err := installupdate.InstallVersion(projectDir, versionFlagVal, checksumFlagVal, cacheDurationFlagVal, false, cmd.OutOrStdout()); err != nil {
						return err
					}
				}
				return nil
			}
			return installupdate.RunActionAndUpgradeConfig(
				projectDir,
				skipUpgradeConfigFlagVal,
				action,
				cmd.OutOrStdout(),
				cmd.OutOrStderr(),
			)
		},
	}
	cmd.Flags().BoolVar(&syncFlagVal, "sync", false, "use version and checksum specified in godel.properties (if true, all other flags are ignored)")
	cmd.Flags().StringVar(&versionFlagVal, "version", "", "version to update (if blank, uses latest version)")
	cmd.Flags().StringVar(&checksumFlagVal, "checksum", "", "expected checksum for package")
	cmd.Flags().DurationVar(&cacheDurationFlagVal, "cache-duration", time.Hour, "duration for which cache entries should be considered valid")
	cmd.Flags().BoolVar(&skipUpgradeConfigFlagVal, "skip-upgrade-config", false, "skips running configuration upgrade tasks after running update")

	return godellauncher.CobraCLITask(cmd, &globalCfg)
}
