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
	"path"
	"time"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/palantir/godel/framework/builtintasks/installupdate"
	"github.com/palantir/godel/framework/godellauncher"
)

func UpdateTask(wrapperPath string) godellauncher.Task {
	var (
		syncFlag          bool
		versionFlag       string
		cacheDurationFlag time.Duration
	)

	cmd := &cobra.Command{
		Use:   "update",
		Short: "Update gödel for project",
		RunE: func(cmd *cobra.Command, args []string) error {
			if wrapperPath == "" {
				return errors.Errorf("wrapper path not specified")
			}
			projectDir := path.Dir(wrapperPath)

			if syncFlag {
				// if sync flag is true, update version to what is specified in gödel.yml
				pkgSrc, err := installupdate.GodelPropsDistPkgInfo(projectDir)
				if err != nil {
					return err
				}
				return installupdate.Update(projectDir, pkgSrc, cmd.OutOrStdout())
			}
			return installupdate.InstallVersion(projectDir, versionFlag, cacheDurationFlag, false, cmd.OutOrStdout())
		},
	}
	cmd.Flags().BoolVar(&syncFlag, "sync", true, "use version and checksum specified in godel.properties")
	cmd.Flags().StringVar(&versionFlag, "version", "", "version to update (if blank, uses latest version)")
	cmd.Flags().DurationVar(&cacheDurationFlag, "cache-duration", time.Hour, "duration for which cache entries should be considered valid")

	return godellauncher.CobraCLITask(cmd)
}
