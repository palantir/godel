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

package main

import (
	"os"
	"path"
	"time"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/palantir/godel/framework/builtintasks/installupdate"
)

func rootCmd() *cobra.Command {
	var (
		versionFlag           string
		checksumFlag          string
		cacheDurationFlag     time.Duration
		skipUpgradeConfigFlag bool
	)

	cmd := &cobra.Command{
		Use:   "godelinit",
		Short: "Add latest version of gödel to a project",
		Long: `godelinit adds godel to a project by adding the godelw script and godel configuration directory to it.
The default behavior adds the newest release of godel on GitHub (https://github.com/palantir/godel/releases)
to the project. If a specific version of godel is desired, it can be specified using the '--version' flag.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			projectDir, err := os.Getwd()
			if err != nil {
				return errors.Wrapf(err, "failed to determine working directory")
			}

			// if current directory does not contain "godelw" wrapper, don't bother trying to upgrade configuration
			if _, err := os.Stat(path.Join(projectDir, "godelw")); err != nil {
				skipUpgradeConfigFlag = true
			}

			// determine version before install
			var godelVersionBeforeUpdate installupdate.Version
			if !skipUpgradeConfigFlag {
				versionBeforeUpdateVar, err := installupdate.GodelVersion(projectDir)
				if err != nil {
					return errors.Wrapf(err, "failed to determine version before update")
				}
				godelVersionBeforeUpdate = versionBeforeUpdateVar
			}

			// perform install
			if err := installupdate.InstallVersion(projectDir, versionFlag, checksumFlag, cacheDurationFlag, true, cmd.OutOrStdout()); err != nil {
				return err
			}

			// run configuration upgrade if needed
			if !skipUpgradeConfigFlag {
				godelVersionAfterUpdate, err := installupdate.GodelVersion(projectDir)
				if err != nil {
					return errors.Wrapf(err, "failed to determine version after update")
				}

				if godelVersionBeforeUpdate.MajorVersionNum() <= 1 && godelVersionAfterUpdate.MajorVersionNum() >= 2 {
					// if going from <=1 to >=2, run "upgrade-config --legacy" task to upgrade configuration
					if err := installupdate.RunUpgradeLegacyConfig(projectDir, cmd.OutOrStdout(), cmd.OutOrStderr()); err != nil {
						return err
					}
				} else if godelVersionBeforeUpdate.MajorVersionNum() >= 2 {
					// if previous version is >=2 and new version is >= previous version, run "upgrade-config"
					if cmp, ok := godelVersionAfterUpdate.CompareTo(godelVersionBeforeUpdate); !ok || cmp >= 0 {
						if err := installupdate.RunUpgradeConfig(projectDir, cmd.OutOrStdout(), cmd.OutOrStderr()); err != nil {
							return err
						}
					}
				}
			}
			return nil
		},
	}

	cmd.Flags().StringVar(&versionFlag, "version", "", "version to install (if unspecified, latest is used)")
	cmd.Flags().StringVar(&checksumFlag, "checksum", "", "expected checksum for package")
	cmd.Flags().DurationVar(&cacheDurationFlag, "cache-duration", time.Hour, "duration for which cache entries should be considered valid")
	cmd.Flags().BoolVar(&skipUpgradeConfigFlag, "skip-upgrade-config", false, "skips running configuration upgrade tasks after running update")
	return cmd
}
