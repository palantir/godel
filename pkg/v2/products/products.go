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
	"strings"
)

// List returns a slice that contains all of the products in the project.
func List() ([]string, error) {
	gödelw, err := newGödelwRunner()
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
	gödelw, err := newGödelwRunner()
	if err != nil {
		return "", err
	}
	productBuildID := product + "." + runtime.GOOS + "-" + runtime.GOARCH

	requiresBuildOutput, err := gödelw.run("artifacts", "build", "--absolute", "--requires-build", productBuildID)
	if err != nil {
		return "", err
	}
	if requiresBuildOutput != "" {
		if _, err := gödelw.run("build", productBuildID); err != nil {
			return "", err
		}
	}
	binPath, err := gödelw.run("artifacts", "build", "--absolute", productBuildID)
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
	gödelw, err := newGödelwRunner()
	if err != nil {
		return "", err
	}
	if _, err := gödelw.run("dist", product); err != nil {
		return "", err
	}
	return gödelw.run("artifacts", "dist", "--absolute", product)
}

type gödelwRunner interface {
	run(args ...string) (string, error)
}

type gödelwRunnerStruct struct {
	path string
}

func (g *gödelwRunnerStruct) run(args ...string) (string, error) {
	cmd := exec.Command(g.path, args...)
	output, err := cmd.CombinedOutput()
	outputStr := strings.TrimSpace(string(output))
	if err != nil {
		err = fmt.Errorf("command %v failed with output:\n%s\nError: %v", cmd.Args, outputStr, err)
	}
	return outputStr, err
}

func newGödelwRunner() (gödelwRunner, error) {
	path, err := gödelwPath()
	if err != nil {
		return nil, err
	}
	return &gödelwRunnerStruct{
		path: path,
	}, nil
}

func gödelwPath() (string, error) {
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
