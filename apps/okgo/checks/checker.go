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

package checks

import (
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/palantir/amalgomate/amalgomated"
	"github.com/palantir/pkg/pkgpath"
	"github.com/pkg/errors"

	"github.com/palantir/godel/apps/okgo/checkoutput"
	"github.com/palantir/godel/apps/okgo/params"
)

type Checker interface {
	// Cmd returns the command represented by this checker
	Cmd() amalgomated.Cmd

	// Check runs the current checker on "all" packages using the provided configuration and returns the output. The
	// provided Cmder is used to invoke the check. "rootDir" is used as the root directory to determine the packages
	// that the checker is run on. Conceptually, this method uses "config" and "rootDir" to generate the arguments
	// required to check "all" packages (excluding any that the config specifies should be excluded), invokes the
	// Cmder using those arguments in the given working directory, and returns the filtered output.
	Check(runner amalgomated.Cmder, rootDir string, config params.OKGo) ([]checkoutput.Issue, error)

	// CheckPackages runs the current checker on the specified packages using the provided configuration and returns
	// the output. The provided Cmder is used to invoke the check. The root directory contained in the "packages"
	// parameter is used as the working directory of the Cmder. Conceptually, this method uses "config" and
	// "packages" to generate the arguments required to check the specified packages (excluding any that the config
	// specifies should be excluded), invokes the Cmder using those arguments in the given working directory, and
	// returns the filtered output.
	CheckPackages(runner amalgomated.Cmder, packages pkgpath.Packages, config params.OKGo) ([]checkoutput.Issue, error)
}

const (
	packages = allArgType("")
	rootDir  = allArgType(".")
	splat    = allArgType("./...")
)

type checkerDefinition struct {
	cmd           amalgomated.Cmd
	lineParser    checkoutput.LineParser
	rawLineFilter func(line string) bool
	allArg        allArgType
	globalCheck   bool // true indicates that per-package checks are not supported (must be run on a root directory)
}

func (c *checkerDefinition) Cmd() amalgomated.Cmd {
	return c.cmd
}

func (c *checkerDefinition) GetParser(rootDir string) checkoutput.IssueParser {
	return &checkoutput.SingleLineIssueParser{
		LineParser: c.lineParser,
		RootDir:    rootDir,
	}
}

func (c *checkerDefinition) Check(cmder amalgomated.Cmder, rootDir string, config params.OKGo) ([]checkoutput.Issue, error) {
	if c.allArg == packages {
		// this Checker should handle "all" by calling "CheckPackages" with "all" packages
		packages, err := pkgpath.PackagesInDir(rootDir, config.Exclude)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to get all packages in directory %s", rootDir)
		}
		return c.CheckPackages(cmder, packages, config)
	}

	checkArgsFromConfig, _ := config.ArgsForCheck(c.cmd)
	allArg := string(c.allArg)
	return c.checkWithFilters(cmder, append(checkArgsFromConfig, allArg), rootDir, config.FiltersForCheck(c.cmd))
}

func (c *checkerDefinition) CheckPackages(cmder amalgomated.Cmder, packages pkgpath.Packages, config params.OKGo) ([]checkoutput.Issue, error) {
	if c.globalCheck {
		// this check does not support package-specific checks
		return nil, errors.Errorf("checker %s does not support specifying packages", c.cmd.Name())
	}

	// exclude packages specified to be excluded in config
	filteredPackages, err := packages.Filter(config.Exclude)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to filter packages %v specified in exclude configuration %v", packages, config.Exclude)
	}

	checkArgsFromConfig, _ := config.ArgsForCheck(c.cmd)

	// translate packages into proper argument format for this runner
	pkgArgs, err := filteredPackages.Paths(pkgpath.Relative)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get paths for packages")
	}
	return c.checkWithFilters(cmder, append(checkArgsFromConfig, pkgArgs...), packages.RootDir(), config.FiltersForCheck(c.cmd))
}

func (c *checkerDefinition) checkWithFilters(cmder amalgomated.Cmder, runnerArgs []string, cmdDir string, filters []checkoutput.Filterer) ([]checkoutput.Issue, error) {
	output, err := c.check(cmder, runnerArgs, cmdDir)
	if err != nil {
		return nil, err
	}
	filteredOutput, err := checkoutput.ApplyFilters(output, filters)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to apply filters to output %v", output)
	}
	return filteredOutput, nil
}

func (c *checkerDefinition) check(cmder amalgomated.Cmder, runnerArgs []string, cmdDir string) ([]checkoutput.Issue, error) {
	var err error
	if !filepath.IsAbs(cmdDir) {
		cmdDir, err = filepath.Abs(cmdDir)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to convert %s to an absolute path", cmdDir)
		}
	}

	cmd := cmder.Cmd(runnerArgs, cmdDir)
	rawOutput, err := cmd.CombinedOutput()
	if _, ok := err.(*exec.ExitError); err != nil && !ok {
		// only propagate error returned from runner if it is not an exec.ExitError -- ExitErrors are ignored
		// because many check programs will return a non-zero exit code if any violations are found even if the
		// check program itself ran successfully.
		return nil, errors.Wrapf(err, "running check with args %v in directory %s failed", runnerArgs, cmdDir)
	}

	reader := strings.NewReader(string(rawOutput))
	parser := c.GetParser(cmdDir)
	return checkoutput.ParseIssues(reader, parser, c.rawLineFilter)
}
