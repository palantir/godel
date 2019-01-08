// Copyright (c) 2016 Palantir Technologies. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cobracli

import (
	"fmt"

	"github.com/spf13/cobra"
)

// VersionCmd returns a command that prints the version of the application with the given name and given version to the
// Stdout of the command.
func VersionCmd(appName, version string) *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: fmt.Sprintf("Print %s version", appName),
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Printf("%s version %s\n", appName, version)
		},
	}
}
