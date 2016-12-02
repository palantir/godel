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

package verify

import (
	"fmt"
	"strings"

	"github.com/nmiyake/pkg/dirs"
	"github.com/palantir/amalgomate/amalgomated"
	"github.com/palantir/pkg/cli"
	"github.com/palantir/pkg/cli/flag"
	"github.com/pkg/errors"

	"github.com/palantir/godel/cmd"
)

const (
	cmdName         = "verify"
	apply           = "apply"
	skipFormat      = "skip-format"
	skipImports     = "skip-imports"
	skipLicense     = "skip-license"
	skipCheck       = "skip-check"
	skipTest        = "skip-test"
	junitOutputPath = "junit-output"
)

func Command(gödelPath string) cli.Command {
	return cli.Command{
		Name:  cmdName,
		Usage: "Run format, check and test tasks",
		Flags: []flag.Flag{
			flag.BoolFlag{Name: apply, Usage: "Apply changes when possible", Value: true},
			flag.BoolFlag{Name: skipFormat, Usage: "Skip 'format' task"},
			flag.BoolFlag{Name: skipImports, Usage: "Skip 'imports' task"},
			flag.BoolFlag{Name: skipLicense, Usage: "Skip 'license' task"},
			flag.BoolFlag{Name: skipCheck, Usage: "Skip 'check' task"},
			flag.BoolFlag{Name: skipTest, Usage: "Skip 'test' task"},
			flag.StringFlag{Name: junitOutputPath, Usage: "Path to JUnit XML output (only used if 'test' task is run)"},
		},
		Action: func(ctx cli.Context) error {
			wd, err := dirs.GetwdEvalSymLinks()
			if err != nil {
				return err
			}

			globalFlags, err := globalFlags(ctx)
			if err != nil {
				return err
			}
			cmder := amalgomated.PathCmder(gödelPath, globalFlags...)

			var failedChecks []string

			if !ctx.Bool(skipFormat) {
				args := []string{"format", "-v"}
				if !ctx.Bool(apply) {
					args = append(args, "-l")
				}
				if err := runCmd(cmder, args, wd, ctx); err != nil {
					failedChecks = append(failedChecks, strings.Join(args, " "))
				}
			}

			if !ctx.Bool(skipImports) {
				args := []string{"imports"}
				if !ctx.Bool(apply) {
					args = append(args, "--verify")
				}
				ctx.Println("Running gocd...")
				if err := runCmd(cmder, args, wd, ctx); err != nil {
					failedChecks = append(failedChecks, strings.Join(args, " "))
				}
			}

			if !ctx.Bool(skipLicense) {
				args := []string{"license"}
				if !ctx.Bool(apply) {
					args = append(args, "--verify")
				}
				ctx.Println("Running golicense...")
				if err := runCmd(cmder, args, wd, ctx); err != nil {
					failedChecks = append(failedChecks, strings.Join(args, " "))
				}
			}

			if !ctx.Bool(skipCheck) {
				if err := runCmd(cmder, []string{"check"}, wd, ctx); err != nil {
					failedChecks = append(failedChecks, "check")
				}
			}

			if !ctx.Bool(skipTest) {
				args := []string{"test"}
				if ctx.Has(junitOutputPath) {
					args = append(args, "--"+junitOutputPath, ctx.String(junitOutputPath))
				}
				if err := runCmd(cmder, args, wd, ctx); err != nil {
					failedChecks = append(failedChecks, "test")
				}
			}

			if len(failedChecks) != 0 {
				msgParts := []string{"Failed tasks:"}
				for _, check := range failedChecks {
					msgParts = append(msgParts, "\t"+check)
				}
				return fmt.Errorf(strings.Join(msgParts, "\n"))
			}

			return nil
		},
	}
}

func globalFlags(ctx cli.Context) ([]string, error) {
	var globalArgs []string
	for _, f := range cmd.GlobalCLIFlags() {
		if ctx.Has(f.MainName()) {
			var flagValue string
			switch f.(type) {
			case flag.BoolFlag:
				flagValue = fmt.Sprintf("%v", ctx.Bool(f.MainName()))
			case flag.StringFlag:
				flagValue = ctx.String(f.MainName())
			default:
				return nil, errors.Errorf("Unhandled flag type %T for flag %v", f, f)
			}
			globalArgs = append(globalArgs, "--"+f.MainName(), flagValue)
		}
	}
	return globalArgs, nil
}

func runCmd(cmder amalgomated.Cmder, args []string, wd string, ctx cli.Context) error {
	cmd := cmder.Cmd(args, wd)
	cmd.Stdout = ctx.App.Stdout
	cmd.Stderr = ctx.App.Stderr
	return cmd.Run()
}
