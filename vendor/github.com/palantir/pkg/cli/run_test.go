// Copyright 2016 Palantir Technologies. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cli_test

import (
	"bytes"
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/palantir/pkg/cli"
	"github.com/palantir/pkg/cli/flag"
)

func TestRunErrorOutput(t *testing.T) {
	cases := []struct {
		err            error
		expectedOutput string
	}{
		// empty error should not produce any output
		{
			err:            fmt.Errorf(""),
			expectedOutput: "",
		},
		// non-empty error is printed with newline appended
		{
			err:            fmt.Errorf("foo"),
			expectedOutput: "foo\n",
		},
	}

	for i, currCase := range cases {
		app := cli.NewApp()
		app.Action = func(ctx cli.Context) error {
			return currCase.err
		}

		stderr := &bytes.Buffer{}
		app.Stderr = stderr

		exitStatus := app.Run([]string{"testApp"})
		if exitStatus == 0 {
			t.Errorf("Case %d: expected exitStatus to be non-0, was: %d", i, exitStatus)
		}

		if stderr.String() != currCase.expectedOutput {
			t.Errorf("Case %d:\nExpected: %q\nActual:   %q", i, currCase.expectedOutput, stderr.String())
		}
	}
}

func TestRunErrorHandler(t *testing.T) {
	cases := []struct {
		err              error
		handler          func(ctx cli.Context, err error) int
		expectedExitCode int
		expectedOutput   string
	}{
		// custom error handler is invoked if provided
		{
			err: fmt.Errorf(""),
			handler: func(ctx cli.Context, err error) int {
				ctx.Errorf("Error: %v\n", err)
				return 13
			},
			expectedExitCode: 13,
			expectedOutput:   "Error: \n",
		},
		// default behavior is used if custom error hander is nil
		{
			err:              fmt.Errorf("foo"),
			expectedExitCode: 1,
			expectedOutput:   "foo\n",
		},
	}

	for i, currCase := range cases {
		app := cli.NewApp()
		app.ErrorHandler = currCase.handler
		app.Action = func(ctx cli.Context) error {
			return currCase.err
		}

		stderr := &bytes.Buffer{}
		app.Stderr = stderr

		exitCode := app.Run([]string{"testApp"})
		assert.Equal(t, currCase.expectedExitCode, exitCode, "Case %d", i)
		assert.Equal(t, currCase.expectedOutput, stderr.String(), "Case %d", i)
	}
}

func TestRunContext(t *testing.T) {
	testFlagName := "test-flag"
	testFlagValue := "foobar"

	var customContextFunc = func(cliCtx cli.Context, ctx context.Context) context.Context {
		for _, flag := range cliCtx.Command.Flags {
			flagVal := cliCtx.FlagValue(flag.MainName())
			ctx = context.WithValue(ctx, flagVal.Name(), flagVal.ValueString())
		}
		return ctx
	}

	var testFunc = func(ctx cli.Context) error {
		assert.Equal(t, testFlagValue, ctx.Context().Value(testFlagName))
		return nil
	}

	cases := []struct {
		name  string
		check func(*testing.T)
	}{
		{
			name: "check that context is propagated to app action",
			check: func(t *testing.T) {
				app := cli.NewApp()

				app.Command.Flags = []flag.Flag{
					flag.StringFlag{
						Name: testFlagName,
					},
				}

				app.ContextConfig = customContextFunc

				app.Action = testFunc

				assert.Equal(t, 0, app.Run([]string{"testApp", "--" + testFlagName, testFlagValue}))
			},
		},
		{
			name: "check that context is propagated to app error handler",
			check: func(t *testing.T) {
				app := cli.NewApp()

				app.Command.Flags = []flag.Flag{
					flag.StringFlag{
						Name:  testFlagName,
						Value: "value",
					},
				}

				app.ContextConfig = customContextFunc

				app.ErrorHandler = func(ctx cli.Context, _ error) int {
					assert.NoError(t, testFunc(ctx))
					return 0
				}

				app.Action = func(ctx cli.Context) error {
					return fmt.Errorf("an error occured")
				}

				assert.Equal(t, 0, app.Run([]string{"testApp", "--" + testFlagName, testFlagValue}))
			},
		},
		{
			name: "check that context is propagated to app subcommand",
			check: func(t *testing.T) {
				app := cli.NewApp()

				app.Command.Flags = []flag.Flag{
					flag.StringFlag{
						Name: testFlagName,
					},
				}

				app.ContextConfig = customContextFunc

				app.Subcommands = []cli.Command{
					{
						Name: "subCommand",

						Action: testFunc,

						Flags: []flag.Flag{
							flag.StringFlag{
								Name: testFlagName,
							},
						},
					},
				}

				assert.Equal(t, 0, app.Run([]string{"testApp", "subCommand", "--" + testFlagName, testFlagValue}))
			},
		},
	}

	for i, currCase := range cases {

		name := fmt.Sprintf("Case %d - %s", i, currCase.name)

		t.Run(name, currCase.check)
	}
}
