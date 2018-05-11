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
	"github.com/spf13/cobra"

	"github.com/palantir/distgo/distgo"
	"github.com/palantir/distgo/distgo/build"
)

var (
	buildCmd = &cobra.Command{
		Use:   "build [flags] [product-build-ids]",
		Short: "Build the executables for products",
		RunE: func(cmd *cobra.Command, args []string) error {
			projectInfo, projectParam, err := distgoProjectParamFromFlags()
			if err != nil {
				return err
			}
			return build.Products(projectInfo, projectParam, distgo.ToProductBuildIDs(args), build.Options{
				Parallel: buildParallelFlagVal,
				Install:  buildInstallFlagVal,
				DryRun:   buildDryRunFlagVal,
			}, cmd.OutOrStdout())
		},
	}
)

var (
	buildParallelFlagVal bool
	buildInstallFlagVal  bool
	buildOSArchsFlagVal  []string
	buildDryRunFlagVal   bool
)

func init() {
	buildCmd.Flags().BoolVar(&buildParallelFlagVal, "parallel", true, "build binaries in parallel")
	buildCmd.Flags().BoolVar(&buildInstallFlagVal, "install", false, "build products with the '-i' flag")
	buildCmd.Flags().StringSliceVar(&buildOSArchsFlagVal, "os-arch", nil, "if specified, only builds the binaries for the specified GOOS-GOARCH(s)")
	buildCmd.Flags().BoolVar(&buildDryRunFlagVal, "dry-run", false, "print the operations that would be performed")

	rootCmd.AddCommand(buildCmd)
}
