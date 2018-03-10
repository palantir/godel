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
		debugVar   *bool
		version    string
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
			nil,
			"",
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
			nil,
			"",
			1,
			"Error: unknown flag: --invalid-flag\nUsage:\n  my-app [flags]\n\nFlags:\n  -h, --help   help for my-app\n",
		},
		{
			"error prints output but without usage",
			func(cmd *cobra.Command, args []string) error {
				return errors.Errorf("hello-error")
			},
			nil,
			nil,
			nil,
			"",
			1,
			"Error: hello-error\n",
		},
		{
			"error with debug variable prints full stack trace",
			func(cmd *cobra.Command, args []string) error {
				return errors.Errorf("hello-error")
			},
			nil,
			nil,
			boolVar(true),
			"",
			1,
			regexp.MustCompile("^Error: hello-error\n\tgithub.com/palantir/pkg/cobracli_test.TestExecuteWithDefaultParams.+"),
		},
		{
			"version command prints version",
			func(cmd *cobra.Command, args []string) error {
				return errors.Errorf("hello-error")
			},
			nil,
			[]string{"version"},
			nil,
			"1.0.0",
			0,
			"my-app version 1.0.0\n",
		},
		{
			"version command does not exist if version is empty",
			func(cmd *cobra.Command, args []string) error {
				cmd.Println(args)
				return nil
			},
			nil,
			[]string{"version"},
			nil,
			"",
			0,
			"[version]\n",
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
			nil,
			"",
			1,
			regexp.MustCompile(`(?s)^Error: .+` + "\nUsage:\n  my-app " + regexp.QuoteMeta(`[flags]`) + "\n\nFlags:\n  -h, --help            help for my-app\n      --required-flag\n$"),
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
			nil,
			"",
			1,
			regexp.MustCompile(`(?s)^Error: .+` + "\nUsage:\n  my-app subcmd " + regexp.QuoteMeta(`[flags]`) + "\n\nFlags:\n  -h, --help           help for subcmd\n      --sub-req-flag\n$"),
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

			rv := cobracli.ExecuteWithDefaultParamsWithVersion(rootCmd, tc.debugVar, tc.version)
			require.Equal(t, tc.wantRV, rv, "Case %d: %s", i, tc.name)

			switch val := tc.wantOutput.(type) {
			case *regexp.Regexp:
				assert.Regexp(t, val, outBuf.String(), "Case %d: %s", i, tc.name)
			case string:
				assert.Equal(t, val, outBuf.String(), "Case %d: %s", i, tc.name)
			default:
				require.Fail(t, "unsupported type: %s. Case %d, %s", val, i, tc.name)
			}
		}()
	}
}

func boolVar(b bool) *bool {
	return &b
}
