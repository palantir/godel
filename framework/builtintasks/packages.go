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

	"github.com/nmiyake/pkg/dirs"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/palantir/godel/framework/builtintasks/packages"
	"github.com/palantir/godel/framework/godellauncher"
)

func PackagesTask() godellauncher.Task {
	return godellauncher.CobraCLITask(&cobra.Command{
		Use:   "packages",
		Short: "Lists all of the packages in the project except those excluded by configuration",
		RunE: func(cmd *cobra.Command, args []string) error {
			wd, err := dirs.GetwdEvalSymLinks()
			if err != nil {
				return errors.Wrapf(err, "failed to determine working directory")
			}

			cfgDir, err := godellauncher.ConfigDirPath(wd)
			if err != nil {
				return err
			}
			cfg, err := godellauncher.ReadGodelConfig(cfgDir)
			if err != nil {
				return err
			}
			pkgs, err := packages.List(cfg.Exclude.Matcher(), wd)
			if err != nil {
				return err
			}
			fmt.Fprintln(cmd.OutOrStdout(), strings.Join(pkgs, "\n"))
			return nil
		},
	})
}
