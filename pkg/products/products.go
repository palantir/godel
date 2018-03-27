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

package products

import (
	"fmt"
	"os/exec"
	"path"
	"runtime"
	"strconv"
	"strings"
)

// List returns a slice that contains all of the products in the project.
func List() ([]string, error) {
	gödelw, err := newGodelwRunner()
	if err != nil {
		return nil, err
	}
	products, err := gödelw.run("products")
	if err != nil {
		return nil, err
	}
	return strings.Split(products, "\n"), nil
}

// Bin returns the path to the executable for the given product for the current OS/Architecture, building the executable
// using "godelw build" if the executable does not already exist or is not up-to-date.
func Bin(product string) (string, error) {
	godelw, err := newGodelwRunner()
	if err != nil {
		return "", err
	}

	// return error if version is too new
	majorVersion, err := majorVersion(godelw)
	if err != nil {
		return "", err
	}
	if majorVersion >= 2 {
		return "", fmt.Errorf("this package does not support godel with major version >=2, but was %d: use v2 of the library instead", majorVersion)
	}

	currOSArchFlag := fmt.Sprintf("--os-arch=%s-%s", runtime.GOOS, runtime.GOARCH)
	requiresBuildOutput, err := godelw.run("artifacts", "build", "--absolute", currOSArchFlag, "--requires-build", product)
	if err != nil {
		return "", err
	}
	if requiresBuildOutput != "" {
		if _, err := godelw.run("build", currOSArchFlag, product); err != nil {
			return "", err
		}
	}
	binPath, err := godelw.run("artifacts", "build", "--absolute", currOSArchFlag, product)
	if err != nil {
		return "", err
	}
	if binPath == "" {
		return "", fmt.Errorf("no build artifact for product %s with GOOS %s and GOARCH %s", product, runtime.GOOS, runtime.GOARCH)
	}
	return binPath, nil
}

// Dist builds the distribution for the specified product using the "godelw dist" command and returns the path to the
// created distribution artifact.
func Dist(product string) (string, error) {
	godelw, err := newGodelwRunner()
	if err != nil {
		return "", err
	}
	if _, err := godelw.run("dist", product); err != nil {
		return "", err
	}
	return godelw.run("artifacts", "dist", "--absolute", product)
}

type godelwRunner interface {
	run(args ...string) (string, error)
}

type godelwRunnerStruct struct {
	path string
}

func (g *godelwRunnerStruct) run(args ...string) (string, error) {
	cmd := exec.Command(g.path, args...)
	output, err := cmd.CombinedOutput()
	outputStr := strings.TrimSpace(string(output))
	if err != nil {
		err = fmt.Errorf("command %v failed with output:\n%s\nError: %v", cmd.Args, outputStr, err)
	}
	return outputStr, err
}

func newGodelwRunner() (godelwRunner, error) {
	path, err := godelwPath()
	if err != nil {
		return nil, err
	}
	return &godelwRunnerStruct{
		path: path,
	}, nil
}

func majorVersion(r godelwRunner) (int, error) {
	versionOutput, err := r.run("version")
	if err != nil {
		return -1, err
	}
	parts := strings.Split(versionOutput, " ")
	if len(parts) < 3 {
		return -1, fmt.Errorf("output of version must have at least 3 ' '-separated parts, but was %q", versionOutput)
	}
	versionParts := strings.Split(parts[2], ".")
	if len(versionParts) < 3 {
		return -1, fmt.Errorf("version must have at least 3 '.'-separated parts, but was %q", parts[2])
	}
	majorVersion, err := strconv.Atoi(versionParts[0])
	if err != nil {
		return -1, fmt.Errorf("unable to parse %q as integer: %v", versionParts[0], err)
	}
	return majorVersion, nil
}

func godelwPath() (string, error) {
	projectDir, err := projectDir()
	if err != nil {
		return "", err
	}
	return path.Join(projectDir, "godelw"), nil
}

func projectDir() (string, error) {
	cmd := exec.Command("git", "rev-parse", "--show-toplevel")
	output, err := cmd.CombinedOutput()
	outputStr := strings.TrimSpace(string(output))
	if err != nil {
		err = fmt.Errorf("command %v failed with output:\n%s\nError: %v", cmd.Args, outputStr, err)
	}
	return outputStr, err
}
