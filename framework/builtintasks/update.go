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
		syncFlag          bool
		versionFlag       string
		checksumFlag      string
		cacheDurationFlag time.Duration
		globalCfg         godellauncher.GlobalConfig
	)

	cmd := &cobra.Command{
		Use:   "update",
		Short: "Update gödel for project",
		RunE: func(cmd *cobra.Command, args []string) error {
			projectDir, err := globalCfg.ProjectDir()
			if err != nil {
				return err
			}
			if syncFlag {
				// if sync flag is true, update version to what is specified in gödel.yml
				pkgSrc, err := installupdate.GodelPropsDistPkgInfo(projectDir)
				if err != nil {
					return err
				}
				return installupdate.Update(projectDir, pkgSrc, cmd.OutOrStdout())
			}
			return installupdate.InstallVersion(projectDir, versionFlag, checksumFlag, cacheDurationFlag, false, cmd.OutOrStdout())
		},
	}
	cmd.Flags().BoolVar(&syncFlag, "sync", false, "use version and checksum specified in godel.properties (if true, all other flags are ignored)")
	cmd.Flags().StringVar(&versionFlag, "version", "", "version to update (if blank, uses latest version)")
	cmd.Flags().StringVar(&checksumFlag, "checksum", "", "expected checksum for package")
	cmd.Flags().DurationVar(&cacheDurationFlag, "cache-duration", time.Hour, "duration for which cache entries should be considered valid")

	return godellauncher.CobraCLITask(cmd, &globalCfg)
}
