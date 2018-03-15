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
	"time"

	"github.com/spf13/cobra"

	"github.com/palantir/distgo/distgo"
	"github.com/palantir/distgo/distgo/dist"
)

var (
	distCmd = &cobra.Command{
		Use:   "dist [flags] [product-dist-ids]",
		Short: "Create distributions for products",
		RunE: func(cmd *cobra.Command, args []string) error {
			projectInfo, projectParam, err := distgoProjectParamFromFlags()
			if err != nil {
				return err
			}

			var configFileModTime *time.Time
			if !distForceFlagVal {
				// if force flag is false, use modification time of configuration file
				configFileModTime = distgoConfigModTime()
			}
			return dist.Products(projectInfo, projectParam, configFileModTime, distgo.ToProductDistIDs(args), distDryRunFlagVal, cmd.OutOrStdout())
		},
	}
)

var (
	distDryRunFlagVal bool
	distForceFlagVal  bool
)

func init() {
	distCmd.Flags().BoolVar(&distDryRunFlagVal, "dry-run", false, "print the operations that would be performed")
	distCmd.Flags().BoolVar(&distForceFlagVal, "force", false, "create distribution outputs even if they are considered up-to-date")

	RootCmd.AddCommand(distCmd)
}
