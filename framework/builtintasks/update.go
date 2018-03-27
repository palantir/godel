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

	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/palantir/godel/framework/builtintasks/installupdate"
	"github.com/palantir/godel/framework/godellauncher"
)

func UpdateTask() godellauncher.Task {
	var (
		syncFlag              bool
		versionFlag           string
		checksumFlag          string
		cacheDurationFlag     time.Duration
		skipUpgradeConfigFlag bool
		globalCfg             godellauncher.GlobalConfig
	)

	cmd := &cobra.Command{
		Use:   "update",
		Short: "Update gödel for project",
		RunE: func(cmd *cobra.Command, args []string) error {
			projectDir, err := globalCfg.ProjectDir()
			if err != nil {
				return err
			}

			var godelVersionBeforeUpdate installupdate.Version
			if !skipUpgradeConfigFlag {
				versionBeforeUpdateVar, err := installupdate.GodelVersion(projectDir)
				if err != nil {
					return errors.Wrapf(err, "failed to determine version before update")
				}
				godelVersionBeforeUpdate = versionBeforeUpdateVar
			}

			if syncFlag {
				// if sync flag is true, update version to what is specified in gödel.yml
				pkgSrc, err := installupdate.GodelPropsDistPkgInfo(projectDir)
				if err != nil {
					return err
				}
				if err := installupdate.Update(projectDir, pkgSrc, cmd.OutOrStdout()); err != nil {
					return err
				}
			} else {
				if err := installupdate.InstallVersion(projectDir, versionFlag, checksumFlag, cacheDurationFlag, false, cmd.OutOrStdout()); err != nil {
					return err
				}
			}

			// run "upgrade-config" after upgrade if new version is greater than or equal to previous version.
			if !skipUpgradeConfigFlag {
				godelVersionAfterUpdate, err := installupdate.GodelVersion(projectDir)
				if err != nil {
					return errors.Wrapf(err, "failed to determine version after update")
				}
				if cmp, ok := godelVersionAfterUpdate.CompareTo(godelVersionBeforeUpdate); !ok || cmp >= 0 {
					if err := installupdate.RunUpgradeConfig(projectDir, cmd.OutOrStdout(), cmd.OutOrStderr()); err != nil {
						return err
					}
				}
			}
			return nil
		},
	}
	cmd.Flags().BoolVar(&syncFlag, "sync", false, "use version and checksum specified in godel.properties (if true, all other flags are ignored)")
	cmd.Flags().StringVar(&versionFlag, "version", "", "version to update (if blank, uses latest version)")
	cmd.Flags().StringVar(&checksumFlag, "checksum", "", "expected checksum for package")
	cmd.Flags().DurationVar(&cacheDurationFlag, "cache-duration", time.Hour, "duration for which cache entries should be considered valid")
	cmd.Flags().BoolVar(&skipUpgradeConfigFlag, "skip-upgrade-config", false, "skips running configuration upgrade tasks after running update")

	return godellauncher.CobraCLITask(cmd, &globalCfg)
}
