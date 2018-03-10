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

package docker

import (
	"fmt"
	"io"
	"os/exec"
	"sort"

	"github.com/palantir/distgo/distgo"
)

func PushProducts(projectInfo distgo.ProjectInfo, projectParam distgo.ProjectParam, productDockerIDs []distgo.ProductDockerID, dryRun bool, stdout io.Writer) error {
	// determine products that match specified productDockerIDs
	productParams, err := distgo.ProductParamsForDockerProductArgs(projectParam.Products, productDockerIDs...)
	if err != nil {
		return err
	}

	// run push only for specified products
	for _, currParam := range productParams {
		if err := RunPush(projectInfo, currParam, dryRun, stdout); err != nil {
			return err
		}
	}
	return nil
}

func RunPush(projectInfo distgo.ProjectInfo, productParam distgo.ProductParam, dryRun bool, stdout io.Writer) error {
	if productParam.Docker == nil {
		distgo.PrintlnOrDryRunPrintln(stdout, fmt.Sprintf("%s does not have Docker outputs; skipping build", productParam.ID), dryRun)
		return nil
	}

	var dockerIDs []distgo.DockerID
	for k := range productParam.Docker.DockerBuilderParams {
		dockerIDs = append(dockerIDs, k)
	}
	sort.Sort(distgo.ByDockerID(dockerIDs))

	productTaskOutputInfo, err := distgo.ToProductTaskOutputInfo(projectInfo, productParam)
	if err != nil {
		return err
	}

	for _, dockerID := range dockerIDs {
		if err := runSingleDockerPush(
			productParam.ID,
			dockerID,
			productTaskOutputInfo,
			dryRun,
			stdout,
		); err != nil {
			return err
		}
	}
	return nil
}

func runSingleDockerPush(
	productID distgo.ProductID,
	dockerID distgo.DockerID,
	productTaskOutputInfo distgo.ProductTaskOutputInfo,
	dryRun bool,
	stdout io.Writer) (rErr error) {

	distgo.PrintlnOrDryRunPrintln(stdout, fmt.Sprintf("Running Docker push for configuration %s of product %s...", dockerID, productID), dryRun)
	for _, tag := range productTaskOutputInfo.Product.DockerOutputInfos.DockerBuilderOutputInfos[dockerID].RenderedTags {
		cmd := exec.Command("docker", "push", tag)
		if err := distgo.RunCommandWithVerboseOption(cmd, true, dryRun, stdout); err != nil {
			return err
		}
	}
	return nil
}
