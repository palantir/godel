// Copyright (c) 2016 Palantir Technologies. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cobracli

import (
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/spf13/cobra"
)

// Execute executes the provided root command configured with the provided parameters. Returns an integer that should be
// used as the exit code for the application. Typical usage is "os.Exit(cobracli.Execute(...))" in a main function.
func Execute(rootCmd *cobra.Command, params ...Param) int {
	executor := &executor{}
	for _, p := range params {
		if p == nil {
			continue
		}
		p.apply(executor)
	}

	for _, configureCmd := range executor.rootCmdConfigurers {
		configureCmd(rootCmd)
	}

	executedCmd, err := rootCmd.ExecuteC()
	if err == nil {
		// command ran successfully: return 0
		return 0
	}

	// print error if error-printing function is defined
	if executor.errorHandler != nil {
		executor.errorHandler(executedCmd, err)
	}

	// extract custom exit code if exit code extractor is defined
	if executor.exitCodeExtractor != nil {
		return executor.exitCodeExtractor(err)
	}

	return 1
}

type executor struct {
	rootCmdConfigurers []func(*cobra.Command)
	errorHandler       func(*cobra.Command, error)
	exitCodeExtractor  func(error) int
}

type Param interface {
	apply(*executor)
}

type paramFunc func(*executor)

func (f paramFunc) apply(e *executor) {
	f(e)
}

// ExitCodeExtractorParam sets the exit code extractor function for the executor. If executing the root command returns
// an error, the error is provided to the function and the code returned by the extractor is used as the exit code.
func ExitCodeExtractorParam(extractor func(error) int) Param {
	return paramFunc(func(executor *executor) {
		executor.exitCodeExtractor = extractor
	})
}

// ErrorHandlerParam sets the error handler for the command. If executing the root command returns an error, the command
// that was executed is provided to the error handler.
func ErrorHandlerParam(handler func(*cobra.Command, error)) Param {
	return paramFunc(func(executor *executor) {
		executor.errorHandler = handler
	})
}

// ErrorPrinterWithDebugHandler returns an error handler that prints the provided error as "Error: <error.Error()>"
// unless "error.Error()" is empty, in which case nothing is printed. If the provided boolean variable pointer is
// non-nil and the value is true, then the error output is provided to the specified error transform function before
// being printed.
func ErrorPrinterWithDebugHandler(debugVar *bool, debugErrTransform func(error) string) func(*cobra.Command, error) {
	return func(command *cobra.Command, err error) {
		errStr := err.Error()
		if errStr == "" {
			return
		}
		if debugVar != nil && *debugVar && debugErrTransform != nil {
			errStr = debugErrTransform(err)
		}
		command.Println("Error:", errStr)
	}
}

// PrintUsageOnRequiredFlagErrorHandlerDecorator decorates the provided error handler to add functionality that prints
// the command usage if the error that occurred was due to a required flag not being specified. This handler first
// processes the error using the provided handler. Then, it examines the string returned by the Error() function of the
// error to determine if it matches the form of an error that indicates that a required flag was missing. If so, the
// usage string of the command is printed.
func PrintUsageOnRequiredFlagErrorHandlerDecorator(fn func(*cobra.Command, error)) func(*cobra.Command, error) {
	return func(command *cobra.Command, err error) {
		// allow inner handler to process first
		fn(command, err)

		if !isRequiredFlagError(err) {
			return
		}
		// if error was a required flags error, print usage
		command.Println(strings.TrimSuffix(command.UsageString(), "\n"))
	}
}

// isRequiredFlagError returns true if the provided error is of the form returned when a required flag is not specified,
// false otherwise.
func isRequiredFlagError(inErr error) bool {
	if inErr == nil || inErr.Error() == "" {
		return false
	}

	// create a dummy command, set it to require a flag, execute it without a flag and parse the error output
	cmd := &cobra.Command{
		Run:           func(cmd *cobra.Command, args []string) {},
		SilenceErrors: true,
		SilenceUsage:  true,
	}
	cmd.SetOutput(ioutil.Discard)
	const dummyFlagName = "dummy-flag-name"
	cmd.Flags().Bool(dummyFlagName, false, "")
	_ = cmd.MarkFlagRequired(dummyFlagName)
	cmd.SetArgs([]string{})
	err := cmd.Execute()
	flagErrStr := err.Error()
	idx := strings.Index(flagErrStr, dummyFlagName)
	if idx == -1 {
		return false
	}

	// determine prefix and suffix of missing required flag error
	prefix := flagErrStr[:idx]
	suffix := flagErrStr[idx+len(dummyFlagName):]

	// if provided error has same prefix and suffix as missing flag error, treat it as a missing flag error
	inErrStr := inErr.Error()
	return strings.HasPrefix(inErrStr, prefix) && strings.HasSuffix(inErrStr, suffix)
}

// ConfigureCmdParam adds the provided configuration function to the executor. All of the configuration functions on the
// executor are applied to the root command before it is executed.
func ConfigureCmdParam(configureCmd func(*cobra.Command)) Param {
	return paramFunc(func(executor *executor) {
		executor.rootCmdConfigurers = append(executor.rootCmdConfigurers, configureCmd)
	})
}

// RemoveHelpCommandConfigurer removes the "help" subcommand from the provided command.
func RemoveHelpCommandConfigurer(command *cobra.Command) {
	// set help command to be empty hidden command to effectively remove the built-in help command. Needs to be done in
	// this manner rather than by removing it because otherwise the default "Execute" logic will re-add the default
	// help command implementation.
	command.SetHelpCommand(&cobra.Command{
		Hidden: true,
	})
}

// AddDebugPersistentFlagParam adds "--debug" as a boolean persistent flag that sets the value of the provided *bool.
func AddDebugPersistentFlagParam(debug *bool) Param {
	return ConfigureCmdParam(func(cmd *cobra.Command) {
		cmd.PersistentFlags().BoolVar(debug, "debug", false, "run in debug mode")
	})
}

// VersionFlagParam configures a command so that its "Version" field has the provided value. If it is non-empty, this
// will add a top-level "--version" flag that prints the version for that command.
func VersionFlagParam(version string) Param {
	return ConfigureCmdParam(func(cmd *cobra.Command) {
		cmd.Version = version
	})
}

// VersionCmdParam configures a command so that it has a "version" subcommand that prints the value of the provided
// version. Is a noop if the provided version is empty.
func VersionCmdParam(version string) Param {
	return ConfigureCmdParam(func(cmd *cobra.Command) {
		if version == "" {
			return
		}
		cmd.AddCommand(VersionCmd(cmd.Use, version))
	})
}

// SilenceErrorsConfigurer configures the provided command to silence the default behavior of printing errors and
// printing command usage on errors.
func SilenceErrorsConfigurer(command *cobra.Command) {
	command.SilenceErrors = true
	command.SilenceUsage = true
}

// FlagErrorsUsageErrorConfigurer configures the provided command such that, when it encounters an error processing a
// flag, the returned error includes the usage string for the command.
func FlagErrorsUsageErrorConfigurer(command *cobra.Command) {
	command.SetFlagErrorFunc(func(c *cobra.Command, err error) error {
		return fmt.Errorf("%s\n%s", err.Error(), strings.TrimSuffix(c.UsageString(), "\n"))
	})
}
