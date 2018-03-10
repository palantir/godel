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
	"github.com/palantir/distgo/distgo/docker"
)

var (
	dockerCmd = &cobra.Command{
		Use:   "docker",
		Short: "Create or push Docker images for products",
	}
	dockerBuildSubCmd = &cobra.Command{
		Use:   "build [flags] [product-docker-ids]",
		Short: "Create Docker images for products",
		RunE: func(cmd *cobra.Command, args []string) error {
			projectInfo, projectParam, err := distgoProjectParamFromFlags()
			if err != nil {
				return err
			}
			if dockerBuildRepositoryFlagVal != "" {
				docker.SetDockerRepository(projectParam, dockerBuildRepositoryFlagVal)
			}
			return docker.BuildProducts(projectInfo, projectParam, distgo.ToProductDockerIDs(args), dockerBuildVerboseFlagVal, dockerBuildDryRunFlagVal, cmd.OutOrStdout())
		},
	}
	dockerPushSubCmd = &cobra.Command{
		Use:   "push [flags] [product-docker-ids]",
		Short: "Push Docker images for products",
		RunE: func(cmd *cobra.Command, args []string) error {
			projectInfo, projectParam, err := distgoProjectParamFromFlags()
			if err != nil {
				return err
			}
			if dockerPushRepositoryFlagVal != "" {
				docker.SetDockerRepository(projectParam, dockerPushRepositoryFlagVal)
			}
			return docker.PushProducts(projectInfo, projectParam, distgo.ToProductDockerIDs(args), dockerPushDryRunFlagVal, cmd.OutOrStdout())
		},
	}
)

var (
	dockerBuildRepositoryFlagVal string
	dockerBuildVerboseFlagVal    bool
	dockerBuildDryRunFlagVal     bool

	dockerPushRepositoryFlagVal string
	dockerPushDryRunFlagVal     bool
)

func init() {
	dockerBuildSubCmd.Flags().StringVar(&dockerBuildRepositoryFlagVal, "repository", "", "specifies the value that should be used for the Docker repository (overrides any value(s) specified in configuration)")
	dockerBuildSubCmd.Flags().BoolVar(&dockerBuildVerboseFlagVal, "verbose", false, "print verbose output for the operation")
	dockerBuildSubCmd.Flags().BoolVar(&dockerBuildDryRunFlagVal, "dry-run", false, "print the operations that would be performed")
	dockerCmd.AddCommand(dockerBuildSubCmd)

	dockerPushSubCmd.Flags().StringVar(&dockerPushRepositoryFlagVal, "repository", "", "specifies the value that should be used for the Docker repository (overrides any value(s) specified in configuration)")
	dockerPushSubCmd.Flags().BoolVar(&dockerPushDryRunFlagVal, "dry-run", false, "print the operations that would be performed")
	dockerCmd.AddCommand(dockerPushSubCmd)

	RootCmd.AddCommand(dockerCmd)
}
