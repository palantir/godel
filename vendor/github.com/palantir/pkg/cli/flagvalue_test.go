// Copyright 2016 Palantir Technologies. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cli_test

import (
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/palantir/pkg/cli"
	"github.com/palantir/pkg/cli/cliviper"
	"github.com/palantir/pkg/cli/flag"
)

func TestBindFlagValues(t *testing.T) {
	const (
		enabledFlagName  = "enabled"
		greetingFlagName = "greeting"
		sizeFlagName     = "size"
		durationFlagName = "duration"
		sliceFlagName    = "slice"
	)

	for i, tc := range []struct {
		name         string
		bindValsFunc func(cli.Context)
	}{
		{
			"bind flag values",
			func(ctx cli.Context) {
				err := viper.BindFlagValues(cliviper.FlagValueSet(&ctx))
				require.NoError(t, err)
			},
		},
		{
			"bind individual flag values",
			func(ctx cli.Context) {
				err := viper.BindFlagValue(enabledFlagName, ctx.FlagValue(enabledFlagName))
				require.NoError(t, err)
				err = viper.BindFlagValue(greetingFlagName, ctx.FlagValue(greetingFlagName))
				require.NoError(t, err)
				err = viper.BindFlagValue(sizeFlagName, ctx.FlagValue(sizeFlagName))
				require.NoError(t, err)
				err = viper.BindFlagValue(durationFlagName, ctx.FlagValue(durationFlagName))
				require.NoError(t, err)
				err = viper.BindFlagValue(sliceFlagName, ctx.FlagValue(sliceFlagName))
				require.NoError(t, err)
			},
		},
	} {
		app := cli.NewApp()
		app.Command = cli.Command{
			Name: "foo",
			Flags: []flag.Flag{
				flag.BoolFlag{Name: enabledFlagName},
				flag.StringFlag{Name: greetingFlagName},
				flag.IntFlag{Name: sizeFlagName},
				flag.DurationFlag{Name: durationFlagName, Value: "0s"},
				flag.StringSlice{Name: sliceFlagName},
			},
			Action: func(ctx cli.Context) error {
				tc.bindValsFunc(ctx)
				return nil
			},
		}

		app.Run([]string{
			"appName",
			"--enabled",
			"--greeting=hello",
			"--size=13",
			"--duration=300ms",
			"a",
			"b",
			"c",
		})

		assert.Equal(t, true, viper.GetBool(enabledFlagName), "Case %d: %s", i, tc.name)
		assert.Equal(t, "hello", viper.GetString(greetingFlagName), "Case %d: %s", i, tc.name)
		assert.Equal(t, 13, viper.GetInt(sizeFlagName), "Case %d: %s", i, tc.name)
		assert.Equal(t, "300ms", viper.GetString(durationFlagName), "Case %d: %s", i, tc.name)
		assert.Equal(t, []string{"a", "b", "c"}, viper.GetStringSlice(sliceFlagName), "Case %d: %s", i, tc.name)
	}
}

func TestBindFlagValuesStringParam(t *testing.T) {
	const (
		stringParamName = "stringParam"
	)

	for i, tc := range []struct {
		name         string
		bindValsFunc func(cli.Context)
	}{
		{
			"bind flag values",
			func(ctx cli.Context) {
				err := viper.BindFlagValues(cliviper.FlagValueSet(&ctx))
				require.NoError(t, err)
			},
		},
		{
			"bind individual flag values",
			func(ctx cli.Context) {
				err := viper.BindFlagValue(stringParamName, ctx.FlagValue(stringParamName))
				require.NoError(t, err)
			},
		},
	} {
		app := cli.NewApp()
		app.Command = cli.Command{
			Name: "foo",
			Flags: []flag.Flag{
				flag.StringParam{Name: stringParamName},
			},
			Action: func(ctx cli.Context) error {
				tc.bindValsFunc(ctx)
				return nil
			},
		}

		app.Run([]string{
			"appName",
			"stringVal",
		})

		assert.Equal(t, "stringVal", viper.GetString(stringParamName), "Case %d: %s", i, tc.name)
	}
}
