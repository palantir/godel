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

package dockerbuilder

import (
	"encoding/json"

	"github.com/palantir/godel/framework/pluginapi"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/palantir/distgo/assetapi"
	"github.com/palantir/distgo/distgo"
)

func AssetRootCmd(creator Creator, upgradeConfigFn pluginapi.UpgradeConfigFn, short string) *cobra.Command {
	name := creator.TypeName()
	rootCmd := &cobra.Command{
		Use:   name,
		Short: short,
	}

	creatorFn := creator.Creator()
	rootCmd.AddCommand(newNameCmd(name))
	rootCmd.AddCommand(newVerifyConfigCmd(creatorFn))
	rootCmd.AddCommand(assetapi.NewAssetTypeCmd(assetapi.DockerBuilder))
	rootCmd.AddCommand(newRunDockerBuildCmd(creatorFn))
	rootCmd.AddCommand(pluginapi.CobraUpgradeConfigCmd(upgradeConfigFn))

	return rootCmd
}

const nameCmdName = "name"

func newNameCmd(name string) *cobra.Command {
	return &cobra.Command{
		Use:   nameCmdName,
		Short: "Print the name of the DockerBuilder",
		RunE: func(cmd *cobra.Command, args []string) error {
			outputJSON, err := json.Marshal(name)
			if err != nil {
				return errors.Wrapf(err, "failed to marshal output as JSON")
			}
			cmd.Print(string(outputJSON))
			return nil
		},
	}
}

const commonCmdConfigYMLFlagName = "config-yml"

const (
	verifyConfigCmdName = "verify-config"
)

func newVerifyConfigCmd(creatorFn CreatorFunction) *cobra.Command {
	var configYMLFlagVal string
	verifyConfigCmd := &cobra.Command{
		Use:   verifyConfigCmdName,
		Short: "Verify that the provided input is valid configuration YML for this DockerBuilder",
		RunE: func(cmd *cobra.Command, args []string) error {
			_, err := creatorFn([]byte(configYMLFlagVal))
			return err
		},
	}
	verifyConfigCmd.Flags().StringVar(&configYMLFlagVal, commonCmdConfigYMLFlagName, "", "configuration YML to verify")
	mustMarkFlagsRequired(verifyConfigCmd, commonCmdConfigYMLFlagName)
	return verifyConfigCmd
}

const (
	runDockerBuildCmdName                          = "run-docker-build"
	runDockerBuildCmdDockerIDFlagName              = "docker-id"
	runDockerBuildCmdProductTaskOutputInfoFlagName = "product-task-output-info"
	runDockerBuildCmdVerboseFlagName               = "verbose"
	runDockerBuildCmdDryRunFlagName                = "dry-run"
)

func newRunDockerBuildCmd(creatorFn CreatorFunction) *cobra.Command {
	var (
		configYMLFlagVal             string
		dockerIDFlagVal              string
		productTaskOutputInfoFlagVal string
		verboseFlagVal               bool
		dryRunFlagVal                bool
	)
	runDockerBuildCmd := &cobra.Command{
		Use:   runDockerBuildCmdName,
		Short: "Runs the Docker build action",
		RunE: func(cmd *cobra.Command, args []string) error {
			dockerBuilder, err := creatorFn([]byte(configYMLFlagVal))
			if err != nil {
				return err
			}
			var productTaskOutputInfo distgo.ProductTaskOutputInfo
			if err := json.Unmarshal([]byte(productTaskOutputInfoFlagVal), &productTaskOutputInfo); err != nil {
				return errors.Wrapf(err, "failed to unmarshal JSON %s", productTaskOutputInfoFlagVal)
			}
			if err := dockerBuilder.RunDockerBuild(distgo.DockerID(dockerIDFlagVal), productTaskOutputInfo, verboseFlagVal, dryRunFlagVal, cmd.OutOrStdout()); err != nil {
				return err
			}
			return nil
		},
	}
	runDockerBuildCmd.Flags().StringVar(&configYMLFlagVal, commonCmdConfigYMLFlagName, "", "YML of DockerBuilder configuration")
	runDockerBuildCmd.Flags().StringVar(&dockerIDFlagVal, runDockerBuildCmdDockerIDFlagName, "", "DockerID for the current DockerBuilder task")
	runDockerBuildCmd.Flags().StringVar(&productTaskOutputInfoFlagVal, runDockerBuildCmdProductTaskOutputInfoFlagName, "", "JSON representation of distgo.ProductTaskOutputInfo")
	runDockerBuildCmd.Flags().BoolVar(&verboseFlagVal, runDockerBuildCmdVerboseFlagName, false, "print verbose output for build task")
	runDockerBuildCmd.Flags().BoolVar(&dryRunFlagVal, runDockerBuildCmdDryRunFlagName, false, "print the steps that would be taken for build without executing them")
	mustMarkFlagsRequired(
		runDockerBuildCmd,
		commonCmdConfigYMLFlagName,
		runDockerBuildCmdDockerIDFlagName,
		runDockerBuildCmdProductTaskOutputInfoFlagName,
		runDockerBuildCmdVerboseFlagName,
		runDockerBuildCmdDryRunFlagName,
	)
	return runDockerBuildCmd
}

func mustMarkFlagsRequired(cmd *cobra.Command, flagNames ...string) {
	for _, currFlagName := range flagNames {
		if err := cmd.MarkFlagRequired(currFlagName); err != nil {
			panic(err)
		}
	}
}
