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
	"github.com/palantir/distgo/distgo/clean"
)

var (
	cleanCmd = &cobra.Command{
		Use:   "clean [flags] [product-ids]",
		Short: "Remove the build and dist outputs for products",
		RunE: func(cmd *cobra.Command, args []string) error {
			projectInfo, projectParam, err := distgoProjectParamFromFlags()
			if err != nil {
				return err
			}
			return clean.Products(projectInfo, projectParam, distgo.ToProductIDs(args), cleanDryRunFlagVal, cmd.OutOrStdout())
		},
	}

	cleanDryRunFlagVal bool
)

func init() {
	cleanCmd.Flags().BoolVar(&cleanDryRunFlagVal, "dry-run", false, "print the paths that would be removed by the operation without actually removing them")

	rootCmd.AddCommand(cleanCmd)
}
