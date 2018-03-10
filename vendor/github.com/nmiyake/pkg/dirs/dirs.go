// MIT License
//
// Copyright (c) 2016 Nick Miyake
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

// Package dirs implements utility functions for creating, getting information on and removing directories.
package dirs

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// SetGoEnvVariables sets the values of the GOPATH and GOROOT environment variables to be the correct values with all
// symbolic links evaluated. The GOPATH to resolve is determined by using the original value of the GOPATH environment
// variable, while the GOROOT to resolve is determined using the GoRoot() function.
func SetGoEnvVariables() error {
	// set GOPATH environment variable to current value of GOPATH environment variable with symlinks resolved
	gopath := os.Getenv("GOPATH")
	if resolvedGoPath, err := filepath.EvalSymlinks(gopath); err == nil {
		gopath = resolvedGoPath
	}
	if err := os.Setenv("GOPATH", gopath); err != nil {
		return err
	}

	// set GOROOT environment variable to be current value of Go root determined by GoRoot() with symlinks resolved
	goroot, err := GoRoot()
	if err != nil {
		return err
	}
	if resolvedGoRoot, err := filepath.EvalSymlinks(goroot); err == nil {
		goroot = resolvedGoRoot
	}
	if err := os.Setenv("GOROOT", goroot); err != nil {
		return err
	}

	return nil
}

// GoRoot returns the value for GOROOT for the current system. Similar to runtime.GOROOT(), but if the GOROOT
// environment variable is not set, falls back on the value provided by the output of running "go env GOROOT" rather
// than using sys.DefaultGoroot. This approach is more portable in situations where the binary is being run in an
// environment that is different from the one in which it is compiled.
func GoRoot() (string, error) {
	if goroot := os.Getenv("GOROOT"); goroot != "" {
		return goroot, nil
	}
	if output, err := exec.Command("go", "env", "GOROOT").CombinedOutput(); err == nil {
		return strings.TrimSpace(string(output)), nil
	}
	return "", fmt.Errorf("unable to determine GOROOT")
}

// GetwdEvalSymLinks returns the working directory of the current process using os.GetWd(). Returns the path with all
// symbolic links resolved (equivalent to 'pwd -P'). Returns an error if os.Getwd or filepath.EvalSymlinks returns an
// error.
func GetwdEvalSymLinks() (string, error) {
	// get wd
	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	// resolve any symlinks in path
	physicalWd, err := filepath.EvalSymlinks(wd)
	if err != nil {
		return "", err
	}
	return physicalWd, nil
}

// MustGetwdEvalSymLinks returns the result of GetwdEvalSymLinks. If GetwdEvalSymLinks returns an error, panics.
func MustGetwdEvalSymLinks() string {
	physicalWd, err := GetwdEvalSymLinks()
	if err != nil {
		panic(fmt.Sprintf("failed to get real path to wd: %v", err))
	}
	return physicalWd
}
