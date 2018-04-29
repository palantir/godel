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

package distgo

import (
	"bytes"
	"fmt"
	"io"
	"os/exec"

	"github.com/pkg/errors"
)

type DockerBuilder interface {
	// TypeName returns the type of this DockerBuilder.
	TypeName() (string, error)

	// RunDockerBuild runs the Docker build task.
	RunDockerBuild(dockerID DockerID, productTaskOutputInfo ProductTaskOutputInfo, verbose, dryRun bool, stdout io.Writer) error
}

type DockerBuilderFactory interface {
	NewDockerBuilder(typeName string, cfgYMLBytes []byte) (DockerBuilder, error)
	ConfigUpgrader(typeName string) (ConfigUpgrader, error)
	Types() []string
}

func RunCommandWithVerboseOption(cmd *exec.Cmd, verbose, dryRun bool, stdout io.Writer) error {
	if dryRun {
		DryRunPrintln(stdout, fmt.Sprintf("Run %v", cmd.Args))
	} else {
		buffer := &bytes.Buffer{}
		cmd.Stdout = buffer
		cmd.Stderr = buffer
		if verbose {
			cmd.Stdout = stdout
			cmd.Stderr = stdout
		}
		if err := cmd.Run(); err != nil {
			output := fmt.Sprintf("command %v failed", cmd.Args)
			if !verbose {
				output += fmt.Sprintf(" with output:\n%s", buffer.String())
			}
			return errors.Wrapf(err, output)
		}
	}
	return nil
}
