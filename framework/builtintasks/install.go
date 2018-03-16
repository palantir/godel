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
	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/palantir/godel/framework/builtintasks/installupdate"
	"github.com/palantir/godel/framework/godellauncher"
	"github.com/palantir/godel/godelgetter"
)

func InstallTask() godellauncher.Task {
	var globalCfg godellauncher.GlobalConfig
	return godellauncher.CobraCLITask(&cobra.Command{
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
			return installupdate.NewInstall(projectDir, godelgetter.NewPkgSrc(args[0], ""), cmd.OutOrStdout())
		},
	}, &globalCfg)
}
