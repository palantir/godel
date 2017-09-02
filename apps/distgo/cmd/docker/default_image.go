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
	"bytes"
	"fmt"
	"io"
	"os/exec"
	"path"

	"github.com/pkg/errors"

	"github.com/palantir/godel/apps/distgo/params"
	"github.com/palantir/godel/apps/distgo/pkg/script"
)

type defaultImageBuilder struct {
	image *params.DockerImage
}

func (di *defaultImageBuilder) build(buildSpec params.ProductBuildSpecWithDeps, stdout io.Writer) error {
	contextDir := path.Join(buildSpec.Spec.ProjectDir, di.image.ContextDir)
	var args []string
	args = append(args, "build")
	args = append(args, "--tag", fmt.Sprintf("%s:%s", di.image.Repository, di.image.Tag))
	buildArgs, err := script.GetBuildArgs(buildSpec.Spec, di.image.BuildArgsScript)
	if err != nil {
		return err
	}
	args = append(args, buildArgs...)
	args = append(args, contextDir)

	buildCmd := exec.Command("docker", args...)
	bufWriter := &bytes.Buffer{}
	buildCmd.Stdout = io.MultiWriter(stdout, bufWriter)
	buildCmd.Stderr = bufWriter
	if err := buildCmd.Run(); err != nil {
		return errors.Wrap(err, fmt.Sprintf("docker build failed with error:\n%s\n", bufWriter.String()))
	}
	return nil
}
