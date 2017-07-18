// Copyright 2016 Palantir Technologies. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cli

import (
	"github.com/palantir/pkg/cli/flag"
)

const (
	DebugFlagName   = "debug"
	debugEnvVarName = "DEBUG_CLI"
)

var debugFlag = flag.BoolFlag{
	Name:   DebugFlagName,
	EnvVar: debugEnvVarName,
	Usage:  "Run in debug mode (print full stack traces on failures and include other debugging output)",
}

type ErrorStringer func(err error) string

// DebugHandler returns an Option function that configures a provided *App with debug handler functionality by adding a
// global "debug" boolean flag and a custom error handler. If the debug flag is true, the error handler prints the
// representation of the error returned by errorStringer to the context's Stderr writer; otherwise, the output of the
// error's "Error" function is printed. If the error implements ExitCoder, the handler returns the value returned by the
// error's "ExitCode" function; otherwise, it returns 1.
func DebugHandler(errorStringer ErrorStringer) Option {
	return func(app *App) {
		app.Flags = append(app.Flags, debugFlag)
		app.ErrorHandler = func(ctx Context, err error) int {
			if ctx.Bool(DebugFlagName) && errorStringer != nil {
				if msg := errorStringer(err); msg != "" {
					ctx.Errorln(msg)
				}
			} else if err.Error() != "" {
				ctx.Errorln(err)
			}
			if exitCoder, ok := err.(ExitCoder); ok {
				return exitCoder.ExitCode()
			}
			return 1
		}
	}
}
