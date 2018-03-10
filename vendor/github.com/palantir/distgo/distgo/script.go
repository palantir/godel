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
	"io/ioutil"
	"os"
	"os/exec"
	"strings"

	"github.com/pkg/errors"
)

func WriteScript(projectInfo ProjectInfo, script string) (name string, cleanup func() error, rErr error) {
	tmpFile, err := ioutil.TempFile(projectInfo.ProjectDir, "")
	if err != nil {
		return "", nil, errors.Wrapf(err, "Failed to create script file")
	}
	cleanup = func() error {
		return os.Remove(tmpFile.Name())
	}
	// clean up unless everything below succeeds
	runCleanup := true
	defer func() {
		if runCleanup {
			if err := cleanup(); err != nil && rErr == nil {
				rErr = errors.Wrapf(err, "failed to remove script file %s", tmpFile.Name())
			}
		}
	}()

	if _, err := tmpFile.WriteString(fmt.Sprintf("#!/bin/bash\n%v", script)); err != nil {
		return "", nil, errors.Wrapf(err, "Failed to write script file")
	}

	if err := tmpFile.Close(); err != nil {
		return "", nil, errors.Wrapf(err, "Failed to close %v", tmpFile.Name())
	}

	if err := os.Chmod(tmpFile.Name(), 0755); err != nil {
		return "", nil, errors.Wrapf(err, "Failed to set file mode of %v to 0755", tmpFile.Name())
	}

	runCleanup = false
	return tmpFile.Name(), cleanup, nil
}

func WriteAndExecuteScript(projectInfo ProjectInfo, script string, additionalEnvVars map[string]string, stdOut io.Writer) (rErr error) {
	// if script exists, write it as a temporary file and execute it
	if script != "" {
		tmpFile, cleanup, err := WriteScript(projectInfo, script)
		if err != nil {
			return err
		}
		defer func() {
			if err := cleanup(); err != nil && rErr == nil {
				rErr = errors.Wrapf(err, "failed to remove script file %s", tmpFile)
			}
		}()

		currEnv := os.Environ()
		distEnvVars := additionalEnvVars
		env := make([]string, len(currEnv), len(currEnv)+len(distEnvVars))
		copy(env, currEnv)
		for k, v := range distEnvVars {
			env = append(env, fmt.Sprintf("%v=%v", k, v))
		}

		cmd := exec.Command(tmpFile)
		cmd.Dir = projectInfo.ProjectDir
		cmd.Env = env
		cmd.Stdout = stdOut
		cmd.Stderr = stdOut
		if err := cmd.Run(); err != nil {
			return errors.Wrapf(err, "script execution failed")
		}
	}
	return nil
}

func BuildArgsFromScript(productTaskOutputInfo ProductTaskOutputInfo, buildArgsScript string) ([]string, error) {
	outputBuf := &bytes.Buffer{}
	if err := WriteAndExecuteScript(productTaskOutputInfo.Project, buildArgsScript, BuildScriptEnvVariables(productTaskOutputInfo), outputBuf); err != nil {
		return nil, errors.Wrapf(err, "failed to execute build args script for %s: %s", productTaskOutputInfo.Product.ID, outputBuf.String())
	}

	buildArgsString := strings.TrimSpace(outputBuf.String())
	if buildArgsString == "" {
		return nil, nil
	}
	return strings.Split(buildArgsString, "\n"), nil
}
