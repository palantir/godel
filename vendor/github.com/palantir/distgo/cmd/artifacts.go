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
	"github.com/palantir/distgo/distgo/artifacts"
	"github.com/palantir/distgo/distgo/docker"
)

var (
	artifactsCmd = &cobra.Command{
		Use:   "artifacts",
		Short: "Print the artifacts for products",
	}
	artifactsBuildSubcmd = &cobra.Command{
		Use:   "build [flags] [product-build-ids]",
		Short: "Print the paths to the build artifacts for products",
		RunE: func(cmd *cobra.Command, args []string) error {
			projectInfo, projectParam, err := distgoProjectParamFromFlags()
			if err != nil {
				return err
			}
			return artifacts.PrintBuildArtifacts(projectInfo, projectParam, distgo.ToProductBuildIDs(args), artifactsAbsPathFlagVal, artifactsRequiresBuildFlagVal, cmd.OutOrStdout())
		},
	}
	artifactsDistSubcmd = &cobra.Command{
		Use:   "dist [flags] [product-dist-ids]",
		Short: "Print the paths to the distribution artifacts for products",
		RunE: func(cmd *cobra.Command, args []string) error {
			projectInfo, projectParam, err := distgoProjectParamFromFlags()
			if err != nil {
				return err
			}
			return artifacts.PrintDistArtifacts(projectInfo, projectParam, distgo.ToProductDistIDs(args), artifactsAbsPathFlagVal, cmd.OutOrStdout())
		},
	}
	artifactsDockerSubcmd = &cobra.Command{
		Use:   "docker [flags] [product-docker-ids]",
		Short: "Print the tags for the Docker images for products",
		RunE: func(cmd *cobra.Command, args []string) error {
			projectInfo, projectParam, err := distgoProjectParamFromFlags()
			if err != nil {
				return err
			}
			if artifactsDockerRepositoryFlagVal != "" {
				docker.SetDockerRepository(projectParam, artifactsDockerRepositoryFlagVal)
			}
			return artifacts.PrintDockerArtifacts(projectInfo, projectParam, distgo.ToProductDockerIDs(args), cmd.OutOrStdout())
		},
	}
)

var (
	artifactsAbsPathFlagVal          bool
	artifactsRequiresBuildFlagVal    bool
	artifactsDockerRepositoryFlagVal string
)

func init() {
	artifactsBuildSubcmd.Flags().BoolVar(&artifactsAbsPathFlagVal, "absolute", false, "print the absolute path for artifacts")
	artifactsBuildSubcmd.Flags().BoolVar(&artifactsRequiresBuildFlagVal, "requires-build", false, "only prints the artifacts that require building (omits artifacts that are already built and are up-to-date)")
	artifactsCmd.AddCommand(artifactsBuildSubcmd)

	artifactsDistSubcmd.Flags().BoolVar(&artifactsAbsPathFlagVal, "absolute", false, "print the absolute path for artifacts")
	artifactsCmd.AddCommand(artifactsDistSubcmd)

	artifactsDockerSubcmd.Flags().StringVar(&artifactsDockerRepositoryFlagVal, "repository", "", "specifies the value that should be used for the Docker repository (overrides any value(s) specified in configuration)")
	artifactsCmd.AddCommand(artifactsDockerSubcmd)

	RootCmd.AddCommand(artifactsCmd)
}
