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

package script

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"

	"github.com/pkg/errors"

	"github.com/palantir/godel/apps/distgo/params"
)

func Write(buildSpec params.ProductBuildSpec, script string) (name string, cleanup func() error, rErr error) {
	tmpFile, err := ioutil.TempFile(buildSpec.ProjectDir, "")
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

func WriteAndExecute(buildSpec params.ProductBuildSpec, script string, stdOut, stdErr io.Writer, additionalEnvVars map[string]string) (rErr error) {
	// if script exists, write it as a temporary file and execute it
	if script != "" {
		tmpFile, cleanup, err := Write(buildSpec, script)
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
		cmd.Env = env
		cmd.Dir = buildSpec.ProjectDir
		cmd.Stdout = stdOut
		cmd.Stderr = stdErr
		if err := cmd.Run(); err != nil {
			return errors.Wrapf(err, "Dist script for %v failed", buildSpec.ProductName)
		}
	}
	return nil
}

func GetBuildArgs(buildSpec params.ProductBuildSpec, script string) ([]string, error) {
	stdoutBuf := bytes.Buffer{}
	stderrBuf := bytes.Buffer{}
	combinedBuf := bytes.Buffer{}
	stdoutMW := io.MultiWriter(&stdoutBuf, &combinedBuf)
	stderrMW := io.MultiWriter(&stderrBuf, &combinedBuf)
	if err := WriteAndExecute(buildSpec, script, stdoutMW, stderrMW, nil); err != nil || stderrBuf.String() != "" {
		return nil, errors.Wrapf(err, "failed to execute build args script for %v: %v", buildSpec.ProductName, combinedBuf.String())
	}

	buildArgsString := strings.TrimSpace(stdoutBuf.String())
	if buildArgsString == "" {
		return nil, nil
	}
	return strings.Split(buildArgsString, "\n"), nil
}
