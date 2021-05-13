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

package installupdate

import (
	"io"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
)

// RunUpgradeConfig runs the "upgrade-config" task by invoking "{{projectDir}}/godelw upgrade-config".
func RunUpgradeConfig(projectDir string, stdout, stderr io.Writer) error {
	return runUpgradeConfig(projectDir, nil, stdout, stderr)
}

// RunUpgradeLegacyConfig runs the "upgrade-config" task in legacy mode by invoking
// "{{projectDir}}/godelw upgrade-config --legacy".
func RunUpgradeLegacyConfig(projectDir string, stdout, stderr io.Writer) error {
	return runUpgradeConfig(projectDir, []string{"--legacy"}, stdout, stderr)
}

func runUpgradeConfig(projectDir string, args []string, stdout, stderr io.Writer) error {
	godelw := filepath.Join(projectDir, "godelw")
	cmd := exec.Command(godelw, append([]string{"upgrade-config"}, args...)...)
	cmd.Stdout = stdout
	cmd.Stderr = stderr
	return cmd.Run()
}

// getGodelVersion returns the Version returned by "{{projectDir}}/godelw version".
func getGodelVersion(projectDir string) (godelVersion, error) {
	godelw := filepath.Join(projectDir, "godelw")
	cmd := exec.Command(godelw, "version")
	output, err := cmd.Output()
	if err != nil {
		return godelVersion{}, errors.Wrapf(err, "failed to execute command %v: %s", cmd.Args, string(output))
	}

	// split input on line breaks and only consider final line. Do this in case invoking "godelw" causes assets to be
	// downloaded (in which case download messages will be in output before version is printed).
	outputString, err := getLastLine(string(output))
	if err != nil {
		return godelVersion{}, err
	}

	parts := strings.Split(outputString, " ")
	if len(parts) != 3 {
		return godelVersion{}, errors.Errorf(`expected output %q to have 3 parts when split by " ", but was %v`, outputString, parts)
	}
	v, err := newGodelVersion(parts[2])
	if err != nil {
		return godelVersion{}, errors.Wrapf(err, "failed to create version from output")
	}
	return v, nil
}

func getLastLine(in string) (string, error) {
	outputLines := strings.Split(strings.Replace(strings.TrimSpace(in), "\r", "\n", -1), "\n")
	if len(outputLines) == 0 {
		return "", errors.Errorf("no elements found when splitting output %q on newlines", in)

	}
	return outputLines[len(outputLines)-1], nil
}
