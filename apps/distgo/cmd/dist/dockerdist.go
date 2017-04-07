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
	"os"
	"os/exec"
	"path"

	"github.com/pkg/errors"

	"github.com/palantir/godel/apps/distgo/params"
)

func dockerDist(buildSpecWithDeps params.ProductBuildSpecWithDeps, distCfg params.Dist) (Packager, error) {
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
	dockerDeps := dockerDistInfo.DistDeps
	for depProduct := range dockerDeps {
		for depDistType := range dockerDeps[depProduct] {
			if depDistType == params.DockerDistType {
				// no need to copy docker artifacts
				continue
			}
			depProductSpec := buildSpecWithDeps.DistDeps[depProduct]
			for _, depDist := range depProductSpec.Dist {
				if depDist.Info.Type() != depDistType {
					continue
				}
				artifactLocation := ArtifactPath(depProductSpec, depDist)
				outputLink := dockerDeps[depProduct][depDistType]
				if outputLink == "" {
					outputLink = path.Base(artifactLocation)
				}
				if err := os.Link(artifactLocation, path.Join(contextDir, outputLink)); err != nil {
					return nil, err
				}
			}
		}
	}

	return packager(func() error {
		if err := buildWithCmd(completeTag, contextDir); err != nil {
			return errors.WithStack(err)
		}
		return nil
	}), nil
}

func buildWithCmd(tag, contextDir string) error {
	var args []string
	args = append(args, "build")
	args = append(args, "--tag", tag)
	args = append(args, contextDir)

	dockerBuild := exec.Command("docker", args...)
	if output, err := dockerBuild.CombinedOutput(); err != nil {
		return errors.Wrap(err, fmt.Sprintf("docker build failed with error:\n%s\n", string(output)))
	}
	return nil
}
