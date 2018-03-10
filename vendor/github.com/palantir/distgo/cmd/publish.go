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
	"fmt"
	"sort"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/palantir/distgo/distgo"
	"github.com/palantir/distgo/distgo/publish"
	"github.com/palantir/distgo/publisher"
)

var (
	publishCmd = &cobra.Command{
		Use:   "publish [action] [flags] [product-dist-ids]",
		Short: "Publish products",
	}
)

var (
	publishDryRunFlagVal bool
)

func init() {
	RootCmd.AddCommand(publishCmd)
}

func addPublishSubcommands() {
	publishers := publisher.Publishers()
	var sortedPublisherKeys []string
	for k := range publishers {
		sortedPublisherKeys = append(sortedPublisherKeys, k)
	}
	sort.Strings(sortedPublisherKeys)

	for _, k := range sortedPublisherKeys {
		publisher := publishers[k]
		currFlags, err := publisher.Flags()
		if err != nil {
			panic(errors.Wrapf(err, "failed to get flags for publisher %s", k))
		}
		currPublisherSubCmd := &cobra.Command{
			Use: fmt.Sprintf("%s [flags] [products]", k),
			RunE: func(cmd *cobra.Command, args []string) error {
				projectInfo, projectParam, err := distgoProjectParamFromFlags()
				if err != nil {
					return err
				}
				flagVals := make(map[distgo.PublisherFlagName]interface{})
				for _, currFlag := range currFlags {
					// if flag was not explicitly provided, don't add it to the flagVals map
					if !cmd.Flags().Changed(string(currFlag.Name)) {
						continue
					}
					val, err := currFlag.GetFlagValue(cmd.Flags())
					if err != nil {
						return err
					}
					flagVals[currFlag.Name] = val
				}
				return publish.Products(projectInfo, projectParam, distgo.ToProductDistIDs(args), publisher, flagVals, publishDryRunFlagVal, cmd.OutOrStdout())
			},
		}
		for _, currFlag := range currFlags {
			if _, err := currFlag.AddFlag(currPublisherSubCmd.Flags()); err != nil {
				panic(errors.Wrapf(err, "failed to add flag %v for publisher %s", currFlag, k))
			}
		}
		currPublisherSubCmd.Flags().BoolVar(&publishDryRunFlagVal, "dry-run", false, "print the operations that would be performed")
		publishCmd.AddCommand(currPublisherSubCmd)
	}
}
