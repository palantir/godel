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
