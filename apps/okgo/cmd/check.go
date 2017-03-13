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

package cmd

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/nmiyake/pkg/dirs"
	"github.com/palantir/amalgomate/amalgomated"
	"github.com/palantir/pkg/cli"
	"github.com/palantir/pkg/cli/cfgcli"
	"github.com/palantir/pkg/cli/flag"
	"github.com/palantir/pkg/pkgpath"
	"github.com/pkg/errors"

	"github.com/palantir/godel/apps/okgo/checkoutput"
	"github.com/palantir/godel/apps/okgo/checks"
	"github.com/palantir/godel/apps/okgo/cmd/cmdlib"
	"github.com/palantir/godel/apps/okgo/config"
	"github.com/palantir/godel/apps/okgo/params"
)

const (
	// releaseTagEnvVar the environment variable used to override the latest release tag that should be used in the
	// default Go build context. See #72 for details.
	releaseTagEnvVar = "OKGO_RELEASE_TAG"
	packagesFlagName = "packages"
)

var packagesFlag = flag.StringSlice{
	Name:     packagesFlagName,
	Usage:    "Packages to check",
	Optional: true,
}

func SetReleaseTagEnvVar(releaseTag string) error {
	if releaseTag != "" {
		if err := os.Setenv(releaseTagEnvVar, releaseTag); err != nil {
			return err
		}
	}
	return nil
}

func GetReleaseTagEnvVar() string {
	return os.Getenv(releaseTagEnvVar)
}

func RunAllCommand(supplier amalgomated.CmderSupplier) cli.Command {
	return cli.Command{
		Name:  "runAll",
		Usage: "Run all checks",
		Flags: []flag.Flag{
			packagesFlag,
		},
		Action: func(ctx cli.Context) error {
			cfg, err := config.Load(cfgcli.ConfigPath, cfgcli.ConfigJSON)
			if err != nil {
				return err
			}
			wd, err := dirs.GetwdEvalSymLinks()
			if err != nil {
				return err
			}
			return DoRunAll(ctx.Slice(packagesFlagName), cfg, supplier, wd, ctx.App.Stdout)
		},
	}
}

func DoRunAll(pkgs []string, cfg params.OKGo, supplier amalgomated.CmderSupplier, wd string, stdout io.Writer) error {
	var checksWithOutput []amalgomated.Cmd
	for _, cmd := range cmdlib.Instance().Cmds() {
		// if "omit" is true, skip the check
		if cmdCfg, ok := cfg.Checks[cmd]; ok && cmdCfg.Skip {
			continue
		}

		cmder, err := supplier(cmd)
		if err != nil {
			return errors.Wrapf(err, "%s is not a valid command", cmd.Name())
		}

		producedOutput, err := executeSingleCheckWithOutput(cmd, cmder, cfg, pkgs, wd, stdout)
		if err != nil {
			// indicates unexpected hard failure -- check returning non-0 exit code will not trigger
			return errors.Wrapf(err, "check %s failed", cmd.Name())
		}

		if producedOutput {
			checksWithOutput = append(checksWithOutput, cmd)
		}
	}

	if len(checksWithOutput) != 0 {
		return errors.Errorf("Checks produced output: %v", checksWithOutput)
	}
	return nil
}

func SingleCheckCommand(cmd amalgomated.Cmd, supplier amalgomated.CmderSupplier) cli.Command {
	return cli.Command{
		Name:  cmd.Name(),
		Usage: "Run " + cmd.Name(),
		Flags: []flag.Flag{
			packagesFlag,
		},
		Action: func(ctx cli.Context) error {
			cfg, err := config.Load(cfgcli.ConfigPath, cfgcli.ConfigJSON)
			if err != nil {
				return err
			}
			wd, err := dirs.GetwdEvalSymLinks()
			if err != nil {
				return err
			}
			if err := SetReleaseTagEnvVar(cfg.ReleaseTag); err != nil {
				return err
			}

			cmder, err := supplier(cmd)
			if err != nil {
				return errors.Wrapf(err, "failed to create Cmder for %s", cmd.Name())
			}

			if producedOutput, err := executeSingleCheckWithOutput(cmd, cmder, cfg, ctx.Slice(packagesFlagName), wd, ctx.App.Stdout); producedOutput {
				return fmt.Errorf("")
			} else if err != nil {
				return err
			}
			return nil
		},
	}
}

// executeSingleCheckWithOutput runs the specified check and outputs the result to stdOut. Returns true if the check
// produced any output, false otherwise.
func executeSingleCheckWithOutput(cmd amalgomated.Cmd, cmder amalgomated.Cmder, cfg params.OKGo, pkgs []string, wd string, stdout io.Writer) (bool, error) {
	output, err := singleCheck(cmd, cmder, cfg, pkgs, wd, stdout)
	if err != nil {
		return false, err
	}

	producedOutput := len(output) != 0
	if producedOutput {
		outputLines := make([]string, len(output))
		for i, currLine := range output {
			outputLines[i] = currLine.String()
		}
		fmt.Fprintln(stdout, strings.Join(outputLines, "\n"))
	}
	return producedOutput, nil
}

func singleCheck(cmd amalgomated.Cmd, cmder amalgomated.Cmder, cfg params.OKGo, pkgs []string, cmdWd string, stdout io.Writer) ([]checkoutput.Issue, error) {
	checker, err := checks.GetChecker(cmd)
	if err != nil {
		return nil, err
	}

	fmt.Fprintf(stdout, "Running %v...\n", cmd.Name())

	if len(pkgs) == 0 {
		// if no arguments were provided, run check on "all"
		return checker.Check(cmder, cmdWd, cfg)
	}

	// convert arguments to packages
	packages, err := pkgpath.PackagesFromPaths(cmdWd, pkgs)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to convert arguments to packages: %v", pkgs)
	}

	// run check on specified packages
	return checker.CheckPackages(cmder, packages, cfg)
}
