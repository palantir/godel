package godellauncher

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/palantir/godel/framework/godel"
)

// CobraCLITask creates a new Task that runs the provided *cobra.Command. The runner for the task does the following:
//
// * Creates a "dummy" root cobra.Command with "godel" as the command name
// * Adds the provided command as a subcommand of the dummy root
// * Executes the root command with the following "os.Args":
//     [executable] [task] [task args...]
func CobraCLITask(cmd *cobra.Command) Task {
	rootCmd := CobraCmdToRootCmd(cmd)
	return Task{
		Name:        cmd.Use,
		Description: cmd.Short,
		RunImpl: func(t *Task, global GlobalConfig, stdout io.Writer) error {
			rootCmd.SetOutput(stdout)
			args := []string{global.Executable}
			args = append(args, global.Task)
			args = append(args, global.TaskArgs...)
			os.Args = args
			return rootCmd.Execute()
		},
	}
}

// CobraCmdToRootCmd takes the provided *cobra.Command and returns a new *cobra.Command that acts as its "root" command.
// The root command has "godel" as its command name and is configured to silence the built-in Cobra error printing.
// However, it has custom logic to match the standard Cobra error output for unrecognized flags.
func CobraCmdToRootCmd(cmd *cobra.Command) *cobra.Command {
	rootCmd := &cobra.Command{
		Use: godel.AppName,
	}
	rootCmd.AddCommand(cmd)

	// Set custom error behavior for flag errors. Usage should only be printed if the error is due to invalid flags. In
	// order to do this, set SilenceErrors and SilenceUsage to true and set flag error function that prints error and
	// usage and then returns an error with empty content (so that the top-level handler does not print it).
	rootCmd.SilenceErrors = true
	rootCmd.SilenceUsage = true
	rootCmd.SetFlagErrorFunc(func(c *cobra.Command, err error) error {
		return fmt.Errorf("%s\n%s", err.Error(), strings.TrimSuffix(c.UsageString(), "\n"))
	})
	return rootCmd
}

func UnknownCommandError(cmd *cobra.Command, args []string) error {
	return fmt.Errorf(`unknown command "%s" for "%s"
Run '%v --help' for usage.`, args[0], cmd.CommandPath(), cmd.CommandPath())
}
