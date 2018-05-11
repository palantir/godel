// Copyright 2016 Palantir Technologies. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cobracli_test

import (
	"bytes"
	"regexp"
	"testing"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/palantir/pkg/cobracli"
)

func TestExecuteWithDefaultParams(t *testing.T) {
	for i, tc := range []struct {
		name       string
		runE       func(cmd *cobra.Command, args []string) error
		configure  func(cmd *cobra.Command)
		args       []string
		wantRV     int
		wantOutput interface{}
	}{
		{
			"standard output",
			func(cmd *cobra.Command, args []string) error {
				cmd.Println("hello, world!")
				return nil
			},
			nil,
			nil,
			0,
			"hello, world!\n",
		},
		{
			"invalid flag prints usage output",
			func(cmd *cobra.Command, args []string) error {
				cmd.Println("hello, world!")
				return nil
			},
			nil,
			[]string{"--invalid-flag"},
			1,
			`Error: unknown flag: --invalid-flag
Usage:
  my-app [flags]

Flags:
      --debug   run in debug mode
  -h, --help    help for my-app
`,
		},
		{
			"error prints output but without usage",
			func(cmd *cobra.Command, args []string) error {
				return errors.Errorf("hello-error")
			},
			nil,
			nil,
			1,
			"Error: hello-error\n",
		},
		{
			"error with debug flag prints full stack trace",
			func(cmd *cobra.Command, args []string) error {
				return errors.Errorf("hello-error")
			},
			nil,
			[]string{"--debug"},
			1,
			regexp.MustCompile(`(?s)^Error: hello-error
	github.com/palantir/pkg/cobracli_test.TestExecuteWithDefaultParams.+`),
		},
		{
			"debug flag is not added to CLI if it already exists",
			func(cmd *cobra.Command, args []string) error {
				return errors.Errorf("hello-error")
			},
			func(cmd *cobra.Command) {
				// add a debug flag outside of default execution
				cmd.Flags().Bool("debug", false, "some other debug flag")
			},
			[]string{"--debug"},
			1,
			// a "--debug" flag was already defined on the root command, so the default executor does not displace the
			// flag. Because that flag is not hooked up to the default Debug variable, no stack trace is printed.
			"Error: hello-error\n",
		},
		{
			"print usage when required flag is not provided",
			func(cmd *cobra.Command, args []string) error {
				cmd.Println(args)
				return nil
			},
			func(cmd *cobra.Command) {
				cmd.Flags().Bool("required-flag", false, "")
				_ = cmd.MarkFlagRequired("required-flag")
			},
			nil,
			1,
			`Error: required flag(s) "required-flag" not set
Usage:
  my-app [flags]

Flags:
      --debug           run in debug mode
  -h, --help            help for my-app
      --required-flag
`,
		},
		{
			"subcommand required flag error prints help for subcommand",
			nil,
			func(cmd *cobra.Command) {
				subCmd := &cobra.Command{
					Use: "subcmd",
					Run: func(cmd *cobra.Command, args []string) {
						cmd.Println("in subcommand")
					},
				}
				subCmd.Flags().Bool("sub-req-flag", false, "")
				_ = subCmd.MarkFlagRequired("sub-req-flag")
				cmd.AddCommand(subCmd)
			},
			[]string{"subcmd"},
			1,
			`Error: required flag(s) "sub-req-flag" not set
Usage:
  my-app subcmd [flags]

Flags:
  -h, --help           help for subcmd
      --sub-req-flag

Global Flags:
      --debug   run in debug mode
`,
		},
	} {
		func() {
			outBuf := &bytes.Buffer{}
			rootCmd := &cobra.Command{
				Use:  "my-app",
				RunE: tc.runE,
			}
			rootCmd.SetOutput(outBuf)
			rootCmd.SetArgs(tc.args)
			if tc.configure != nil {
				tc.configure(rootCmd)
			}

			rv := cobracli.ExecuteWithDefaultParams(rootCmd)
			require.Equal(t, tc.wantRV, rv, "Case %d: %s", i, tc.name)

			switch val := tc.wantOutput.(type) {
			case *regexp.Regexp:
				assert.Regexp(t, val, outBuf.String(), "Case %d: %s\nGot:\n%s", i, tc.name, outBuf.String())
			case string:
				assert.Equal(t, val, outBuf.String(), "Case %d: %s\nGot:\n%s", i, tc.name, outBuf.String())
			default:
				require.Fail(t, "unsupported type: %s. Case %d, %s", val, i, tc.name)
			}
		}()
	}
}
