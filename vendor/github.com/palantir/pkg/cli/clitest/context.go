// Copyright 2016 Palantir Technologies. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package clitest

import (
	"bytes"
	"fmt"

	"github.com/palantir/pkg/cli"
	"github.com/palantir/pkg/cli/flag"
)

func Context(flags map[string]interface{}) cli.Context {
	app := cli.NewApp()
	app.Stdout = new(bytes.Buffer)
	app.Stderr = new(bytes.Buffer)
	app.Flags = make([]flag.Flag, 0, len(flags))

	args := make([]string, 0, len(flags)*2+1)
	args = append(args, "dummyApp")
	for name, value := range flags {
		app.Flags = append(app.Flags, dummyFlag{Name: name, Value: value})
		args = append(args, fmt.Sprintf("--%v", name), "OVERWRITTEN_VALUE")
	}

	var theCtx cli.Context
	app.Action = func(ctx cli.Context) error {
		theCtx = ctx
		return nil
	}

	status := app.Run(args)
	if status != 0 {
		panic(status)
	}

	return theCtx
}

// Stdout returns the output printed to ctx.App.Stdout as a string. Assumes context was created by clitest.Context.
func Stdout(ctx cli.Context) string { return ctx.App.Stdout.(*bytes.Buffer).String() }

// Stderr returns the output printed to ctx.App.Stderr as a string. Assumes context was created by clitest.Context.
func Stderr(ctx cli.Context) string { return ctx.App.Stderr.(*bytes.Buffer).String() }
