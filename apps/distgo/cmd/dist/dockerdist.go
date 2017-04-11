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

package dist

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path"

	"github.com/pkg/errors"

	"github.com/palantir/godel/apps/distgo/params"
)

func dockerDist(buildSpecWithDeps params.ProductBuildSpecWithDeps, distCfg params.Dist, stdout io.Writer) (Packager, error) {
	fmt.Fprintf(stdout, "Creating docker distribution for %v\n", buildSpecWithDeps.Spec.ProductName)
	if _, ok := distCfg.Info.(*params.DockerDistInfo); !ok {
		return nil, errors.New("Dist info provided is not of type docker info")
	}
	dockerDistInfo := *distCfg.Info.(*params.DockerDistInfo)
	if dockerDistInfo.Tag == "" {
		dockerDistInfo.Tag = buildSpecWithDeps.Spec.ProductVersion
	}
	if dockerDistInfo.Repository == "" {
		dockerDistInfo.Repository = buildSpecWithDeps.Spec.ProductName
	}
	completeTag := fmt.Sprintf("%s:%s", dockerDistInfo.Repository, dockerDistInfo.Tag)
	contextDir := path.Join(buildSpecWithDeps.Spec.ProjectDir, dockerDistInfo.ContextDir)

	// link dependent artifacts into the context directory
	for depProduct, distTypes := range dockerDistInfo.DistDeps.ToMap() {
		for distType := range distTypes {
			if distType == params.DockerDistType {
				// no need to copy docker artifacts
				continue
			}
			depProductSpec := buildSpecWithDeps.DistDeps[depProduct]
			matches := 0
			for _, depDist := range depProductSpec.Dist {
				if depDist.Info.Type() != distType {
					continue
				}
				matches++
				artifactLocation := ArtifactPath(depProductSpec, depDist)
				targetFile := distTypes[distType]
				if targetFile == "" {
					targetFile = path.Base(artifactLocation)
				}
				target := path.Join(contextDir, targetFile)
				if _, err := os.Stat(target); err == nil {
					// ensure the target does not exists before creating a new one
					if err := os.Remove(target); err != nil {
						return nil, err
					}
				}
				if err := os.Link(artifactLocation, target); err != nil {
					return nil, err
				}
			}
			if matches == 0 {
				return nil, errors.Errorf("Failed to build docker dist for %v. The dependent dist type %v does not exist on the product: %v\n",
					buildSpecWithDeps.Spec.ProductName, distType, depProduct)
			}
		}
	}

	return packager(func() error {
		if err := buildWithCmd(completeTag, contextDir, stdout); err != nil {
			return errors.WithStack(err)
		}
		return nil
	}), nil
}

func buildWithCmd(tag, contextDir string, stdout io.Writer) error {
	var args []string
	args = append(args, "build")
	args = append(args, "--tag", tag)
	args = append(args, contextDir)
	fmt.Fprintf(stdout, "Building docker image %v\n", tag)

	dockerBuild := exec.Command("docker", args...)
	if output, err := dockerBuild.CombinedOutput(); err != nil {
		return errors.Wrap(err, fmt.Sprintf("docker build failed with error:\n%s\n", string(output)))
	}
	return nil
}
