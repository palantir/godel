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
	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/palantir/distgo/distgo"
	"github.com/palantir/distgo/distgo/run"
)

var (
	runCmd = &cobra.Command{
		Use:   "run [product-id] [arguments...]",
		Short: "Run product",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				return errors.Errorf("a single product must be specified as the first argument")
			}
			projectInfo, projectParam, err := distgoProjectParamFromFlags()
			if err != nil {
				return err
			}
			productParams, err := distgo.ProductParamsForProductArgs(projectParam.Products, distgo.ProductID(args[0]))
			if err != nil {
				return err
			}
			return run.Product(projectInfo, productParams[0], args[1:], cmd.OutOrStdout(), cmd.OutOrStderr())
		},
	}
)

func init() {
	RootCmd.AddCommand(runCmd)
}
