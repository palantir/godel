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

	"github.com/palantir/godel/docsgenerator/generator"
)

var (
	inputDirFlagVal  string
	outputDirFlagVal string
	baseImageFlagVal string
	tagPrefixFlagVal string

	runDockerBuildFlagVal       bool
	suppressDockerOutputFlagVal bool
	startStepFlagVal            int
	endStepFlagVal              int
	leaveGeneratedFilesFlagVal  bool
)

var RootCmd = &cobra.Command{
	Use: "docs-generator",
	RunE: func(cmd *cobra.Command, args []string) error {
		params := generator.Params{
			TagPrefix:            tagPrefixFlagVal,
			RunDockerBuild:       runDockerBuildFlagVal,
			SuppressDockerOutput: suppressDockerOutputFlagVal,
			StartStep:            startStepFlagVal,
			EndStep:              endStepFlagVal,
			LeaveGeneratedFiles:  leaveGeneratedFilesFlagVal,
		}
		return generator.Generate(inputDirFlagVal, outputDirFlagVal, baseImageFlagVal, params, cmd.OutOrStdout())
	},
}

func init() {
	RootCmd.Flags().StringVar(&inputDirFlagVal, "input-dir", "", "input directory")
	if err := RootCmd.MarkFlagRequired("input-dir"); err != nil {
		panic(err)
	}
	RootCmd.Flags().StringVar(&outputDirFlagVal, "output-dir", "", "output directory")
	if err := RootCmd.MarkFlagRequired("output-dir"); err != nil {
		panic(err)
	}
	RootCmd.Flags().StringVar(&baseImageFlagVal, "base-image", "", "the base image for the first Docker image")
	if err := RootCmd.MarkFlagRequired("base-image"); err != nil {
		panic(err)
	}

	RootCmd.Flags().StringVar(&tagPrefixFlagVal, "tag-prefix", "docsgenerator", "the prefix for the Docker tag used for the images")
	RootCmd.Flags().BoolVar(&runDockerBuildFlagVal, "run-docker-build", true, "run the 'docker build' actions for the templates")
	RootCmd.Flags().BoolVar(&suppressDockerOutputFlagVal, "suppress-docker-output", false, "suppress the output of the 'docker build' operation(s)")
	RootCmd.Flags().IntVar(&startStepFlagVal, "start-step", -1, "start step")
	RootCmd.Flags().IntVar(&endStepFlagVal, "end-step", -1, "end step")
	RootCmd.Flags().BoolVar(&leaveGeneratedFilesFlagVal, "leave-generated-files", false, "do not clean up the generated intermediate files")
}
