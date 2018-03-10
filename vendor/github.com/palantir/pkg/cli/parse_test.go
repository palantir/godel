// Copyright 2016 Palantir Technologies. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cli_test

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"regexp"
	"testing"

	"github.com/palantir/pkg/cli"
	"github.com/palantir/pkg/cli/flag"
)

func TestParseFlags(t *testing.T) {
	cases := []struct {
		name           string
		flags          []flag.Flag
		args           []string
		expectedOutput string
		expectedError  string
	}{
		{
			name: "optional string flag with default value has value without flag",
			flags: []flag.Flag{
				flag.StringFlag{
					Name:  "name",
					Value: "default",
				},
			},
			args: []string{
				"./test",
				"test-cmd",
			},
			expectedOutput: "name: default name: [default]",
		},
		{
			name: "bool flag with default value has value without flag",
			flags: []flag.Flag{
				flag.BoolFlag{
					Name:  "bool",
					Value: true,
				},
			},
			args: []string{
				"./test",
				"test-cmd",
			},
			expectedOutput: "bool: true",
		},
		{
			name: "string flag with space parses correctly",
			flags: []flag.Flag{
				flag.StringFlag{
					Name: "name",
				},
			},
			args: []string{
				"./test",
				"test-cmd",
				"--name",
				"foo",
			},
			expectedOutput: "name: foo name: [foo]",
		},
		{
			name: "string flag with '=' parses correctly",
			flags: []flag.Flag{
				flag.StringFlag{
					Name: "name",
				},
			},
			args: []string{
				"./test",
				"test-cmd",
				"--name=foo",
			},
			expectedOutput: "name: foo name: [foo]",
		},
		{
			name: "string flag with empty value after '=' does not parse (interpreted as missing flag value)",
			flags: []flag.Flag{
				flag.StringFlag{
					Name: "name",
				},
			},
			args: []string{
				"./test",
				"test-cmd",
				"--name=",
			},
			expectedError: "Missing value for flag --name",
		},
		{
			name: "string flag with empty value after '=' does not parse (interpreted as missing flag value)",
			flags: []flag.Flag{
				flag.StringFlag{
					Name: "name",
				},
			},
			args: []string{
				"./test",
				"test-cmd",
				"--name=",
				"subcommand",
			},
			expectedError: "Missing value for flag --name",
		},
		{
			name: "string flag with multiple values parses correctly",
			flags: []flag.Flag{
				flag.StringFlag{
					Name: "name",
				},
			},
			args: []string{
				"./test",
				"test-cmd",
				"--name",
				"foo",
				"--name",
				"bar",
				"--name",
				"bar",
			},
			expectedOutput: "name: bar name: [foo bar bar]",
		},
		{
			name: "parameters are not parsed as flags",
			flags: []flag.Flag{
				flag.StringSlice{
					Name: "args",
				},
			},
			args: []string{
				"./test",
				"test-cmd",
				"foo=1",
				"bar=2",
			},
			expectedOutput: "args: [foo=1 bar=2]",
		},
		{
			name: "'=' is a legal character in a value",
			flags: []flag.Flag{
				flag.StringFlag{
					Name: "name",
				},
			},
			args: []string{
				"./test",
				"test-cmd",
				"--name",
				"foo=bar",
			},
			expectedOutput: "name: foo=bar name: [foo=bar]",
		},
		{
			name: "only first '=' in a flag is considered",
			flags: []flag.Flag{
				flag.StringFlag{
					Name: "name",
				},
			},
			args: []string{
				"./test",
				"test-cmd",
				"--name=foo=bar",
			},
			expectedOutput: "name: foo=bar name: [foo=bar]",
		},
		{
			name: "flag name can contain '='",
			flags: []flag.Flag{
				flag.StringFlag{
					Name: "name=foo",
				},
			},
			args: []string{
				"./test",
				"test-cmd",
				"--name=foo=bar",
			},
			expectedOutput: "name=foo: bar name=foo: [bar]",
		},
		{
			name: "bool flag with no value parses as 'true'",
			flags: []flag.Flag{
				flag.BoolFlag{
					Name: "bool",
				},
			},
			args: []string{
				"./test",
				"test-cmd",
				"--bool",
			},
			expectedOutput: "bool: true",
		},
		{
			name: "bool flag can be set to false using '--flag=' syntax",
			flags: []flag.Flag{
				flag.BoolFlag{
					Name:  "bool",
					Value: true,
				},
			},
			args: []string{
				"./test",
				"test-cmd",
				"--bool=false",
			},
			expectedOutput: "bool: false",
		},
		{
			name: "bool flag with missing value after '=' is invalid",
			flags: []flag.Flag{
				flag.BoolFlag{
					Name: "bool",
				},
			},
			args: []string{
				"./test",
				"test-cmd",
				"--bool=",
			},
			expectedError: "Missing value for flag --bool",
		},
		{
			name: "bool flag with invalid value is invalid",
			flags: []flag.Flag{
				flag.BoolFlag{
					Name: "bool",
				},
			},
			args: []string{
				"./test",
				"test-cmd",
				"--bool=NOT_VALID",
			},
			expectedError: `--bool: strconv.ParseBool: parsing "NOT_VALID": invalid syntax`,
		},
		{
			name: "int flag parses correctly",
			flags: []flag.Flag{
				flag.IntFlag{
					Name: "int",
				},
			},
			args: []string{
				"./test",
				"test-cmd",
				"--int=5",
			},
			expectedOutput: "int: 5 int: [5]",
		},
		{
			name: "int flag with multiple values parses correctly",
			flags: []flag.Flag{
				flag.IntFlag{
					Name: "int",
				},
			},
			args: []string{
				"./test",
				"test-cmd",
				"--int=5",
				"--int=13",
				"--int=13",
			},
			expectedOutput: "int: 13 int: [5 13 13]",
		},
		{
			name: "int flag with invalid value is invalid",
			flags: []flag.Flag{
				flag.IntFlag{
					Name: "int",
				},
			},
			args: []string{
				"./test",
				"test-cmd",
				"--int=NOT_VALID",
			},
			expectedError: `--int: strconv.ParseInt: parsing "NOT_VALID": invalid syntax`,
		},
	}

	for i, currCase := range cases {
		t.Run(currCase.name, func(t *testing.T) {
			app := cli.NewApp()
			app.Name = "test"

			output := &bytes.Buffer{}
			app.Subcommands = []cli.Command{
				{
					Name:  "test-cmd",
					Flags: currCase.flags,
					Action: func(ctx cli.Context) error {
						printFlags(output, ctx, currCase.flags)
						return nil
					},
				},
			}

			app.Stdout = ioutil.Discard

			stdErr := &bytes.Buffer{}
			app.Stderr = stdErr

			exitStatus := app.Run(currCase.args)
			expectedExitStatus := 0
			if currCase.expectedError != "" {
				expectedExitStatus = 1
			}
			if expectedExitStatus != exitStatus {
				t.Errorf("Case %d:\nExpected: %d\nActual:   %d", i, expectedExitStatus, exitStatus)
			}

			if currCase.expectedOutput != output.String() {
				t.Errorf("Case %d:\nExpected: %q\nActual:   %q", i, currCase.expectedOutput, output.String())
			}

			if currCase.expectedError != "" && !regexp.MustCompile(currCase.expectedError).MatchString(stdErr.String()) {
				t.Errorf("Case %d: regexp did not match\nExpected: %v\nActual:   %v", i, currCase.expectedError, stdErr.String())
			}
		})
	}
}

func printFlags(w io.Writer, ctx cli.Context, flags []flag.Flag) {
	for _, currFlag := range flags {
		switch currFlag := currFlag.(type) {
		case flag.StringFlag:
			fmt.Fprintf(w, "%v: %v", currFlag.Name, ctx.String(currFlag.Name))
			fmt.Fprintf(w, " %v: %v", currFlag.Name, ctx.StringSlice(currFlag.Name))
		case flag.StringParam:
			fmt.Fprintf(w, "%v: %v ", currFlag.Name, ctx.String(currFlag.Name))
			fmt.Fprintf(w, "%v: %v", currFlag.Name, ctx.StringSlice(currFlag.Name))
		case flag.BoolFlag:
			fmt.Fprintf(w, "%v: %v", currFlag.Name, ctx.Bool(currFlag.Name))
		case flag.StringSlice:
			fmt.Fprintf(w, "%v: %v", currFlag.Name, ctx.Slice(currFlag.Name))
		case flag.IntFlag:
			fmt.Fprintf(w, "%v: %v", currFlag.Name, ctx.Int(currFlag.Name))
			fmt.Fprintf(w, " %v: %v", currFlag.Name, ctx.IntSlice(currFlag.Name))
		default:
			panic(fmt.Sprintf("unsupported type: %v", currFlag))
		}
	}
}
