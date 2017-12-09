// Copyright 2016 Palantir Technologies. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cobracli

import (
	"github.com/nmiyake/pkg/errorstringer"
	"github.com/spf13/cobra"
)

var Version = "unspecified"

// ExecuteWithDefaultParams executes the provided root command using the parameters returned by DefaultParams. This
// function also adds a "version command that prints the value of the "version" variable of this package. The value of
// this version variable should be set using build flags. Typical usage is
// "os.Exit(cobracli.ExecuteWithDefaultParams(...))" in a main function.
func ExecuteWithDefaultParams(rootCmd *cobra.Command, debugVar *bool) int {
	return ExecuteWithDefaultParamsWithVersion(rootCmd, debugVar, Version)
}

// ExecuteWithDefaultParamsWithVersion executes the provided root command using the parameters returned by
// DefaultParams. Typical usage is "os.Exit(cobracli.ExecuteWithDefaultParamsWithVersion(...))" in a main function.
func ExecuteWithDefaultParamsWithVersion(rootCmd *cobra.Command, debugVar *bool, version string) int {
	return Execute(rootCmd, DefaultParams(debugVar, version)...)
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
// * If the provided version is non-empty, adds a "version" command that prints the version of the application in the
//   form "<rootCmd.Use> version <version>".
func DefaultParams(debugVar *bool, version string) []Param {
	params := []Param{
		// silence default error and usage printing provided by cobra CLI
		ConfigureCmdParam(SilenceErrorsConfigurer),
		// if error is encountered while parsing a flag, include the usage for the command as part of the error
		ConfigureCmdParam(FlagErrorsUsageErrorConfigurer),
		// set error handler that prints "Error: <error content>" (unless error content is empty, in which case nothing
		// is printed). If the value of the provided debug boolean pointer is true, then if the error is a pkg/errors
		// error, the full stack trace is printed.
		ErrorHandlerParam(PrintUsageOnRequiredFlagErrorHandlerDecorator(ErrorPrinterWithDebugHandler(debugVar, errorstringer.StackWithInterleavedMessages))),
	}
	if version != "" {
		params = append(params,
			// add the "version" command
			ConfigureCmdParam(func(cmd *cobra.Command) {
				cmd.AddCommand(VersionCmd(cmd.Use, version))
			}),
		)
	}
	return params
}
