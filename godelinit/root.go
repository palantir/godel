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
	"time"

	"github.com/nmiyake/pkg/dirs"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/palantir/godel/framework/builtintasks/installupdate"
)

func rootCmd() *cobra.Command {
	var (
		versionFlag       string
		cacheDurationFlag time.Duration
	)

	cmd := &cobra.Command{
		Use:   "godelinit",
		Short: "Add latest version of gödel to a project",
		Long: `godelinit adds godel to a project by adding the godelw script and godel configuration directory to it.
The default behavior adds the newest release of godel on GitHub (https://github.com/palantir/godel/releases)
to the project. If a specific version of godel is desired, it can be specified using the '--version' flag.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			wd, err := dirs.GetwdEvalSymLinks()
			if err != nil {
				return errors.Wrapf(err, "failed to determine working directory")
			}
			return installupdate.InstallVersion(wd, versionFlag, cacheDurationFlag, true, cmd.OutOrStdout())
		},
	}

	cmd.Flags().StringVar(&versionFlag, "version", "", "version to install (if unspecified, latest is used)")
	cmd.Flags().DurationVar(&cacheDurationFlag, "cache-duration", time.Hour, "duration for which cache entries should be considered valid")
	return cmd
}
