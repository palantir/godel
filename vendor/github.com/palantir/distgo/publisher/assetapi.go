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

package publisher

import (
	"encoding/json"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/palantir/distgo/assetapi"
	"github.com/palantir/distgo/distgo"
)

func AssetRootCmd(creator Creator, short string) *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   creator.TypeName(),
		Short: short,
	}

	publisher := creator.Publisher()
	rootCmd.AddCommand(newNameCmd(publisher))
	rootCmd.AddCommand(assetapi.NewAssetTypeCmd(assetapi.Publisher))
	rootCmd.AddCommand(newFlagsCmd(publisher))
	rootCmd.AddCommand(newRunPublishCmd(publisher))

	return rootCmd
}

const nameCmdName = "name"

func newNameCmd(publisher distgo.Publisher) *cobra.Command {
	return &cobra.Command{
		Use:   nameCmdName,
		Short: "Print the name of the publisher",
		RunE: func(cmd *cobra.Command, args []string) error {
			name, err := publisher.TypeName()
			if err != nil {
				return err
			}
			outputJSON, err := json.Marshal(name)
			if err != nil {
				return errors.Wrapf(err, "failed to marshal output as JSON")
			}
			cmd.Print(string(outputJSON))
			return nil
		},
	}
}

const flagsCmdName = "flags"

func newFlagsCmd(publisher distgo.Publisher) *cobra.Command {
	flagsCmd := &cobra.Command{
		Use:   flagsCmdName,
		Short: "Prints the specifications for the flags supported by this publish operation",
		RunE: func(cmd *cobra.Command, args []string) error {
			flags, err := publisher.Flags()
			if err != nil {
				return err
			}
			outputJSON, err := json.Marshal(flags)
			if err != nil {
				return errors.Wrapf(err, "failed to marshal output as JSON")
			}
			cmd.Print(string(outputJSON))
			return nil
		},
	}
	return flagsCmd
}

const (
	runPublishCmdName                          = "run-publish"
	runPublishCmdProductTaskOutputInfoFlagName = "product-task-output-info"
	runPublishCmdConfigYMLFlagName             = "config-yml"
	runPublishCmdFlagValsFlagName              = "flag-vals"
	runPublishCmdDryRunFlagName                = "dry-run"
)

func newRunPublishCmd(publisher distgo.Publisher) *cobra.Command {
	var (
		productTaskOutputInfoFlagVal string
		configYMLFlagVal             string
		flagValsFlagVal              string
		dryRunFlagVal                bool
	)
	runDistCmd := &cobra.Command{
		Use:   runPublishCmdName,
		Short: "Runs the publish action",
		RunE: func(cmd *cobra.Command, args []string) error {
			var productTaskOutputInfo distgo.ProductTaskOutputInfo
			if err := json.Unmarshal([]byte(productTaskOutputInfoFlagVal), &productTaskOutputInfo); err != nil {
				return errors.Wrapf(err, "failed to unmarshal JSON %s", productTaskOutputInfoFlagVal)
			}
			var flagVals map[distgo.PublisherFlagName]interface{}
			if err := json.Unmarshal([]byte(flagValsFlagVal), &flagVals); err != nil {
				return errors.Wrapf(err, "failed to unmarshal JSON %s", flagValsFlagVal)
			}
			if err := publisher.RunPublish(productTaskOutputInfo, []byte(configYMLFlagVal), flagVals, dryRunFlagVal, cmd.OutOrStdout()); err != nil {
				return err
			}
			return nil
		},
	}
	runDistCmd.Flags().StringVar(&productTaskOutputInfoFlagVal, runPublishCmdProductTaskOutputInfoFlagName, "", "JSON representation of distgo.ProductTaskOutputInfo")
	runDistCmd.Flags().StringVar(&configYMLFlagVal, runPublishCmdConfigYMLFlagName, "", "the configuration YML for this publish operation")
	runDistCmd.Flags().StringVar(&flagValsFlagVal, runPublishCmdFlagValsFlagName, "", "JSON representation of map[distgo.PublisherFlag]interface{}")
	runDistCmd.Flags().BoolVar(&dryRunFlagVal, runPublishCmdDryRunFlagName, false, "true if the operation should be run as a dry run")
	mustMarkFlagsRequired(runDistCmd,
		runPublishCmdProductTaskOutputInfoFlagName,
		runPublishCmdConfigYMLFlagName,
		runPublishCmdFlagValsFlagName,
		runPublishCmdDryRunFlagName,
	)
	return runDistCmd
}

func mustMarkFlagsRequired(cmd *cobra.Command, flagNames ...string) {
	for _, currFlagName := range flagNames {
		if err := cmd.MarkFlagRequired(currFlagName); err != nil {
			panic(err)
		}
	}
}
