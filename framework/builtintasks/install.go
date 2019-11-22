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
	"github.com/palantir/godel/v2/framework/builtintasks/installupdate"
	"github.com/palantir/godel/v2/framework/godellauncher"
	"github.com/palantir/godel/v2/godelgetter"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

func InstallTask() godellauncher.Task {
	var (
		globalCfg                godellauncher.GlobalConfig
		checksumFlagVal          string
		skipUpgradeConfigFlagVal bool
	)

	cmd := &cobra.Command{
		Use:   "install",
		Short: "Install gödel from a local tgz file",
		RunE: func(cmd *cobra.Command, args []string) error {
			projectDir, err := globalCfg.ProjectDir()
			if err != nil {
				return err
			}
			if len(args) == 0 {
				return errors.Errorf("path to package to install must be provided as an argument")
			}
			return installupdate.RunActionAndUpgradeConfig(
				projectDir,
				skipUpgradeConfigFlagVal,
				func() error {
					return installupdate.NewInstall(projectDir, godelgetter.NewPkgSrc(args[0], checksumFlagVal), cmd.OutOrStdout())
				},
				cmd.OutOrStdout(),
				cmd.OutOrStderr(),
			)
		},
	}
	cmd.Flags().BoolVar(&skipUpgradeConfigFlagVal, "skip-upgrade-config", false, "skips running configuration upgrade tasks after installation")
	cmd.Flags().StringVar(&checksumFlagVal, "checksum", "", "expected checksum for package")

	return godellauncher.CobraCLITask(cmd, &globalCfg)
}
