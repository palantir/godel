// Copyright 2016 Palantir Technologies. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cli

import (
	"fmt"
	"os"
	"regexp"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/palantir/pkg/cli/flag"
)

func TestHasReturnsTrueForPresentValue(t *testing.T) {
	ctx := createTestContext(t, flag.StringFlag{Name: "testFlag"}, []string{"--testFlag", "testValue"})
	assert.True(t, ctx.Has("testFlag"))
}

func TestHasReturnsFalseForMissingValue(t *testing.T) {
	ctx := createTestContext(t, flag.StringFlag{Name: "testFlag"}, []string{})
	assert.False(t, ctx.Has("notPresentFlag"))
}

func TestHasReturnsFalseForDefaultValue(t *testing.T) {
	ctx := createTestContext(t, flag.StringFlag{Name: "testFlag", Value: "testFlagDefault"}, []string{})
	assert.False(t, ctx.Has("testFlag"))
}

func TestParseBoolFlagWithEnvVar(t *testing.T) {
	const envVarName = "TestParseBoolFlagWithEnvVar"

	for i, currCase := range []struct {
		defaultVal  bool
		envVal      string
		providedVal string
		want        bool
	}{
		// default value is used
		{defaultVal: true, want: true},
		{defaultVal: false, want: false},
		// if value is provided, it overrides default
		{defaultVal: true, providedVal: "0", want: false},
		{defaultVal: false, providedVal: "1", want: true},
		// if environment variable is set, it overrides default
		{defaultVal: true, envVal: "f", want: false},
		{defaultVal: false, envVal: "t", want: true},
		// if environment variable cannot be parsed, it is interpreted as "false"
		{defaultVal: true, envVal: "notABool", want: false},
		// if environment variable is set and value is provided, provided value is used
		{defaultVal: true, envVal: "true", providedVal: "false", want: false},
		{defaultVal: false, envVal: "false", providedVal: "true", want: true},
	} {
		boolFlag := flag.BoolFlag{Name: "testFlag", Value: currCase.defaultVal}

		// set environment variable if currCase.envVal is non-empty
		if currCase.envVal != "" {
			boolFlag.EnvVar = envVarName
			err := os.Setenv(envVarName, currCase.envVal)
			require.NoError(t, err)
		}

		var args []string
		if currCase.providedVal != "" {
			args = append(args, "--testFlag="+currCase.providedVal)
		}
		ctx := createTestContext(t, boolFlag, args)
		got := ctx.Bool("testFlag")
		assert.Equal(t, currCase.want, got, "Case %d", i)

		err := os.Unsetenv(envVarName)
		require.NoError(t, err)
	}
}

func TestTypedGetFunctionsReturnValue(t *testing.T) {
	flagName := "testFlag"

	cases := []struct {
		getFunction func(ctx Context) interface{}
		testFlag    flag.Flag
		flagValue   []string
		value       interface{}
	}{
		{
			getFunction: func(ctx Context) interface{} {
				return ctx.String(flagName)
			},
			testFlag:  flag.StringFlag{Name: flagName, Value: "testDefault"},
			flagValue: []string{"testValue"},
			value:     "testValue",
		},
		{
			getFunction: func(ctx Context) interface{} {
				return ctx.Bool(flagName)
			},
			testFlag:  flag.BoolFlag{Name: flagName, Value: false},
			flagValue: []string{},
			value:     true,
		},
		{
			getFunction: func(ctx Context) interface{} {
				return ctx.Duration(flagName)
			},
			testFlag:  flag.DurationFlag{Name: flagName, Value: "120s"},
			flagValue: []string{time.Minute.String()},
			value:     time.Minute,
		},
		{
			getFunction: func(ctx Context) interface{} {
				return ctx.Int(flagName)
			},
			testFlag:  flag.IntFlag{Name: flagName, Value: 3},
			flagValue: []string{"5"},
			value:     5,
		},
	}

	for _, c := range cases {
		ctx := createTestContext(t, c.testFlag, append([]string{"--" + flagName}, c.flagValue...))
		assert.Equal(t, c.value, c.getFunction(ctx))
	}
}

func TestTypedGetFunctionsPanicOnMissingFlags(t *testing.T) {
	flagName := "testFlag"
	missingFlagName := "invalidFlag"

	cases := []struct {
		getFunction func(ctx Context) interface{}
		testFlag    flag.Flag
	}{
		{
			getFunction: func(ctx Context) interface{} {
				return ctx.String(missingFlagName)
			},
			testFlag: flag.StringFlag{Name: flagName, Value: "testDefault"},
		},
		{
			getFunction: func(ctx Context) interface{} {
				return ctx.Bool(missingFlagName)
			},
			testFlag: flag.BoolFlag{Name: flagName, Value: true},
		},
		{
			getFunction: func(ctx Context) interface{} {
				return ctx.Duration(missingFlagName)
			},
			testFlag: flag.DurationFlag{Name: flagName, Value: "120s"},
		},
		{
			getFunction: func(ctx Context) interface{} {
				return ctx.Int(missingFlagName)
			},
			testFlag: flag.IntFlag{Name: flagName, Value: 3},
		},
	}

	for _, c := range cases {
		testFunction := func() {
			ctx := createTestContext(t, c.testFlag, []string{})
			c.getFunction(ctx)
		}
		assertPanic(t, testFunction, fmt.Sprintf("command \"testCommand\" does not have a flag named \"%s\"", missingFlagName))
	}
}

func TestTypedGetFunctionsPanicOnValueOfWrongType(t *testing.T) {
	flagName := "testFlag"

	cases := []struct {
		getFunction      func(ctx Context) interface{}
		incompatibleFlag flag.Flag
		incompatibleType string
		typeName         string
	}{
		{
			getFunction: func(ctx Context) interface{} {
				return ctx.String(flagName)
			},
			incompatibleFlag: flag.BoolFlag{Name: flagName, Value: true},
			incompatibleType: "bool",
			typeName:         "string",
		},
		{
			getFunction: func(ctx Context) interface{} {
				return ctx.Bool(flagName)
			},
			incompatibleFlag: flag.StringFlag{Name: flagName, Value: "testDefault"},
			incompatibleType: "string",
			typeName:         "bool",
		},
		{
			getFunction: func(ctx Context) interface{} {
				return ctx.Duration(flagName)
			},
			incompatibleFlag: flag.StringFlag{Name: flagName, Value: "testDefault"},
			incompatibleType: "string",
			typeName:         "time.Duration",
		},
		{
			getFunction: func(ctx Context) interface{} {
				return ctx.Int(flagName)
			},
			incompatibleFlag: flag.StringFlag{Name: flagName, Value: "testDefault"},
			incompatibleType: "string",
			typeName:         "int",
		},
	}

	for _, c := range cases {
		testFunction := func() {
			ctx := createTestContext(t, c.incompatibleFlag, []string{})
			c.getFunction(ctx)
		}
		// conditional regexp causes panic message to match for both Go 1.8 and 1.7
		assertPanic(t, testFunction, fmt.Sprintf(`interface conversion: interface (\{\} )?is %s, not %s`, c.incompatibleType, c.typeName))
	}
}

func createTestContext(t *testing.T, f flag.Flag, args []string) Context {
	app := NewApp()
	app.Subcommands = []Command{
		{
			Name: "testCommand",
			Flags: []flag.Flag{
				f,
			},
		},
	}
	app.Flags = append(app.Flags, versionFlag)
	ctx, err := app.parse(append([]string{"deployctl", "testCommand"}, args...))
	if err != nil {
		assert.Fail(t, "error creating context", "%v", err)
	}
	return ctx
}

// assert that provided function panics and that the value
// of panic is an error whose "Error()" method produces a
// string that is equal to the provided string
func assertPanic(t *testing.T, testFunction func(), expected string) {
	defer func() {
		if panicValue := recover(); panicValue == nil {
			assert.Fail(t, "expected panic")
		} else if err, ok := panicValue.(error); !ok {
			assert.Fail(t, "value of panic was not of type error: was type %T with value %#v", panicValue, panicValue)
		} else {
			assert.Regexp(t, regexp.MustCompile(expected), err.Error())
		}
	}()

	testFunction()
}
