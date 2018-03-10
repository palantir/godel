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
	"io"
	"os/exec"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"

	"github.com/palantir/distgo/distgo"
)

const DefaultBuilderTypeName = "default"

type DefaultDockerBuilderConfig struct {
	BuildArgs []string `yaml:"build-args"`
}

func (cfg *DefaultDockerBuilderConfig) ToDockerBuilder() distgo.DockerBuilder {
	return &defaultDockerBuilder{
		buildArgs: cfg.BuildArgs,
	}
}

type defaultDockerBuilder struct {
	buildArgs []string
}

func NewDefaultDockerBuilder(buildArgs []string) distgo.DockerBuilder {
	return &defaultDockerBuilder{
		buildArgs: buildArgs,
	}
}

func NewDefaultDockerBuilderFromConfig(cfgYML []byte) (distgo.DockerBuilder, error) {
	var dockerBuilderCfg DefaultDockerBuilderConfig
	if err := yaml.Unmarshal(cfgYML, &dockerBuilderCfg); err != nil {
		return nil, errors.Wrapf(err, "failed to unmarshal YAML")
	}
	return dockerBuilderCfg.ToDockerBuilder(), nil
}

func (d *defaultDockerBuilder) TypeName() (string, error) {
	return DefaultBuilderTypeName, nil
}

func (d *defaultDockerBuilder) RunDockerBuild(dockerID distgo.DockerID, productTaskOutputInfo distgo.ProductTaskOutputInfo, verbose, dryRun bool, stdout io.Writer) error {
	dockerBuilderOutputInfo := productTaskOutputInfo.Product.DockerOutputInfos.DockerBuilderOutputInfos[dockerID]
	args := []string{
		"build",
		"--file", dockerBuilderOutputInfo.DockerfilePath,
	}
	for _, tag := range dockerBuilderOutputInfo.RenderedTags {
		args = append(args,
			"-t", tag,
		)
	}
	args = append(args, dockerBuilderOutputInfo.ContextDir)
	args = append(args, d.buildArgs...)

	cmd := exec.Command("docker", args...)
	if err := distgo.RunCommandWithVerboseOption(cmd, verbose, dryRun, stdout); err != nil {
		return err
	}
	return nil
}
