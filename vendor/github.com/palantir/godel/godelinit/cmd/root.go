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

package cmd

import (
	"os"
	"path"
	"time"

	"github.com/palantir/pkg/cobracli"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/palantir/godel/framework/builtintasks/installupdate"
	"github.com/palantir/godel/godelgetter"
)

var (
	localFlagVal             string
	versionFlagVal           string
	checksumFlagVal          string
	cacheDurationFlagVal     time.Duration
	skipUpgradeConfigFlagVal bool
)

var rootCmd = &cobra.Command{
	Use:   "godelinit",
	Short: "Add latest version of gödel to a project",
	Long: `godelinit adds godel to a project by adding the godelw script and godel configuration directory to it.
The default behavior adds the newest release of godel on GitHub (https://github.com/palantir/godel/releases)
to the project. If a specific version of godel is desired, it can be specified using the '--version' flag.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) > 0 {
			return cmd.Usage()
		}
		if localFlagVal != "" && versionFlagVal != "" {
			return errors.Errorf("cannot specify both '--local' and '--version' flag")
		}

		projectDir, err := os.Getwd()
		if err != nil {
			return errors.Wrapf(err, "failed to determine working directory")
		}

		// if current directory does not contain "godelw" wrapper, don't bother trying to upgrade configuration
		if _, err := os.Stat(path.Join(projectDir, "godelw")); err != nil {
			skipUpgradeConfigFlagVal = true
		}

		var runFn func() error
		if localFlagVal != "" {
			runFn = func() error {
				return installupdate.NewInstall(projectDir, godelgetter.NewPkgSrc(localFlagVal, checksumFlagVal), cmd.OutOrStdout())
			}
		} else {
			runFn = func() error {
				return installupdate.InstallVersion(projectDir, versionFlagVal, checksumFlagVal, cacheDurationFlagVal, true, cmd.OutOrStdout())
			}
		}
		return installupdate.RunActionAndUpgradeConfig(
			projectDir,
			skipUpgradeConfigFlagVal,
			runFn,
			cmd.OutOrStdout(),
			cmd.OutOrStderr(),
		)
	},
}

func Execute() int {
	return cobracli.ExecuteWithDefaultParams(rootCmd)
}

func init() {
	rootCmd.Flags().StringVar(&localFlagVal, "local", "", "path to local tgz file that should be used for installation")
	rootCmd.Flags().StringVar(&versionFlagVal, "version", "", "version to install (if unspecified, latest is used)")
	rootCmd.Flags().StringVar(&checksumFlagVal, "checksum", "", "expected checksum for package")
	rootCmd.Flags().DurationVar(&cacheDurationFlagVal, "cache-duration", time.Hour, "duration for which cache entries should be considered valid")
	rootCmd.Flags().BoolVar(&skipUpgradeConfigFlagVal, "skip-upgrade-config", false, "skips running configuration upgrade tasks after running update")
}
