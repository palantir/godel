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
	"os/exec"
	"path"

	"github.com/pkg/errors"

	"github.com/palantir/godel/apps/distgo/params"
	"github.com/palantir/godel/apps/distgo/pkg/script"
)

type DefaultDockerImage struct {
	// Repository and Tag are the part of the image coordinates.
	// For example, in alpine:latest, alpine is the repository
	// and the latest is the tag
	Repository string
	Tag        string
	// ContextDir is the directory in which the docker build task is executed.
	ContextDir      string
	BuildArgsScript string
	// DistDeps is a slice of DockerDistDep.
	// DockerDistDep contains a product, dist type and target file.
	// For a particular product's dist type, we create a link from its output
	// inside the ContextDir with the name specified in target file.
	// This will be used to order the dist tasks such that all the dependent
	// products' dist tasks will be executed first, after which the dist tasks for the
	// current product are executed.
	Deps []params.DockerDep
}

func (di *DefaultDockerImage) ContextDirectory() string {
	return di.ContextDir
}

func (di *DefaultDockerImage) Dependencies() []params.DockerDep {
	return di.Deps
}

func (di *DefaultDockerImage) Build(buildSpec params.ProductBuildSpecWithDeps) error {
	contextDir := path.Join(buildSpec.Spec.ProjectDir, di.ContextDir)
	var args []string
	args = append(args, "build")
	args = append(args, "--tag", fmt.Sprintf("%s:%s", di.Repository, di.Tag))
	buildArgs, err := script.GetBuildArgs(buildSpec.Spec, di.BuildArgsScript)
	if err != nil {
		return errors.Wrap(err, "")
	}
	if len(buildArgs) > 0 {
		args = append(args, buildArgs...)
	}
	args = append(args, contextDir)

	buildCmd := exec.Command("docker", args...)
	if output, err := buildCmd.CombinedOutput(); err != nil {
		return errors.Wrap(err, fmt.Sprintf("docker build failed with error:\n%s\n", string(output)))
	}
	return nil
}

func (di *DefaultDockerImage) Coordinates() (string, string) {
	return di.Repository, di.Tag
}

func (di *DefaultDockerImage) SetDefaults(repo, tag string) {
	if di.Repository == "" {
		di.Repository = repo
	}
	if di.Tag == "" {
		di.Tag = tag
	}
}

func (di *DefaultDockerImage) SetRepository(repo string) {
	di.Repository = repo
}
