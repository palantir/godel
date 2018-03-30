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
		versionFlagVal           string
		checksumFlagVal          string
		cacheDurationFlagVal     time.Duration
		skipUpgradeConfigFlagVal bool
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
				skipUpgradeConfigFlagVal = true
			}
			return installupdate.RunActionAndUpgradeConfig(
				projectDir,
				skipUpgradeConfigFlagVal,
				func() error {
					return installupdate.InstallVersion(projectDir, versionFlagVal, checksumFlagVal, cacheDurationFlagVal, true, cmd.OutOrStdout())
				},
				cmd.OutOrStdout(),
				cmd.OutOrStderr(),
			)
		},
	}

	cmd.Flags().StringVar(&versionFlagVal, "version", "", "version to install (if unspecified, latest is used)")
	cmd.Flags().StringVar(&checksumFlagVal, "checksum", "", "expected checksum for package")
	cmd.Flags().DurationVar(&cacheDurationFlagVal, "cache-duration", time.Hour, "duration for which cache entries should be considered valid")
	cmd.Flags().BoolVar(&skipUpgradeConfigFlagVal, "skip-upgrade-config", false, "skips running configuration upgrade tasks after running update")
	return cmd
}
