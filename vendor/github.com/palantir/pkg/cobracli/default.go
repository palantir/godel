// Copyright (c) 2016 Palantir Technologies. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cobracli

import (
	"github.com/nmiyake/pkg/errorstringer"
	"github.com/spf13/cobra"
)

// ExecuteWithDefaultParams executes the provided root command using the parameters returned by DefaultParams and the
// provided params. Also adds a "--debug" boolean flag as a persistent flag on the command (unless the command already
// has a flag with that name). If the "--debug" flag is added by this function, then invoking the command with the flag
// will make it such that, if the command exits with an error, full stack traces will be printed if available as part of
// the error.
func ExecuteWithDefaultParams(rootCmd *cobra.Command, params ...Param) int {
	debug := false
	return ExecuteWithDebugVarAndDefaultParams(rootCmd, &debug, params...)
}

// ExecuteWithDebugVarAndDefaultParams executes the provided root command using the parameters returned by DefaultParams
// and the provided params. If the provided debugVar pointer is non-nil, adds a "--debug" boolean flag as a persistent
// flag on the command (unless the command already has a flag with that name). If the "--debug" flag is added by this
// function, then invoking the command with the flag will set the value of the variable pointed to by debugVar. If that
// variable is true, then if the command exits with an error, full stack traces will be printed if available as part of
// the error.
func ExecuteWithDebugVarAndDefaultParams(rootCmd *cobra.Command, debugVar *bool, params ...Param) int {
	defaultParams := DefaultParams(debugVar)
	if debugVar != nil {
		if debugFlag := rootCmd.Flag("debug"); debugFlag == nil {
			defaultParams = append(defaultParams, AddDebugPersistentFlagParam(debugVar))
		}
	}
	defaultParams = append(defaultParams, params...)
	return Execute(rootCmd, defaultParams...)
}

// DefaultParams returns a slice of Params that configures Cobra CLI execution with specific opinionated default
// behavior:
//
// * Sets SilenceErrors and SilenceUsage to true, which disables Cobra's built-in error and usage behavior. This
//   prevents the behavior where usage is printed on any error returned by the command.
// * Registers a custom flag usage error handler that appends Cobra's command usage string to the errors encountered
//   while parsing flags. This makes it such that errors that occur due to invalid flags do print the usage.
// * Registers an error printer that prints top-level errors as "Error: <error.Error()>" unless <error.Error()> is the
//   empty string, in which case no error is printed. If the "debugVar" pointer is non-nil and its underlying value is
//   true, then <error.Error()> is printed as a full verbose stack trace if it is a pkg/errors error. This printer is
//   also configured to print the usage output for a command if the command returns an error that indicates that a
//   required flag was not provided.
func DefaultParams(debugVar *bool) []Param {
	return []Param{
		// silence default error and usage printing provided by cobra CLI
		ConfigureCmdParam(SilenceErrorsConfigurer),
		// if error is encountered while parsing a flag, include the usage for the command as part of the error
		ConfigureCmdParam(FlagErrorsUsageErrorConfigurer),
		// set error handler that prints "Error: <error content>" (unless error content is empty, in which case nothing
		// is printed). If the value of the provided debug boolean pointer is true, then if the error is a pkg/errors
		// error, the full stack trace is printed.
		ErrorHandlerParam(PrintUsageOnRequiredFlagErrorHandlerDecorator(ErrorPrinterWithDebugHandler(debugVar, errorstringer.StackWithInterleavedMessages))),
	}
}
