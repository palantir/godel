// Copyright 2016 Palantir Technologies, Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package amalgomated

import (
	"flag"
	"os"
	"strings"
)

const ProxyCmdPrefix = "__"

// RunApp runs an application using the provided arguments. The arguments should be of the same form as os.Args (the
// first element contains the invoking executable command and the rest contain the elements). If there are flags that
// can occur before the proxy command that should be ignored for the purposes of determining whether or not a command is
// a proxy command, they should be provided in "fset". If osArgs[1] of the non-flag arguments exists and is a proxy
// command, the corresponding command in cmdSet is run with the rest of the arguments and os.Exit is called. Otherwise,
// the provided "app" function is run and its return value is returned.
func RunApp(osArgs []string, fset *flag.FlagSet, cmdLibrary CmdLibrary, app func(osArgs []string) int) int {
	// process provided commands and process if it is a proxy command
	if processProxyCmd(osArgs, fset, cmdLibrary) {
		return 0
	}

	// otherwise, run the provided application
	return app(osArgs)
}

// processProxyCmd checks the second non-flag element of the provided osArgs slice (which is the first argument to the
// executable) to see if it is a proxy command. If it is, "os.Args" is set to be the non-proxy command arguments and the
// un-proxied command is run; otherwise, it is a no-op. Returns true if the command waas proxied and the proxied command
// was run; false otherwise. Note that, if a proxied command is run, it is possible that the proxied command may call
// some form of "os.Exit" itself. If this is the case, then this function will be terminal and will not return a value.
func processProxyCmd(osArgs []string, fset *flag.FlagSet, cmdLibrary CmdLibrary) bool {
	// if fset is provided and is able to parse the arguments, parse the provided arguments and only consider the
	// non-flag arguments when determining whether or not the arguments constitute a proxy command
	if fset != nil && len(osArgs) > 0 {
		if err := fset.Parse(osArgs[1:]); err == nil {
			osArgs = append([]string{osArgs[0]}, fset.Args()...)
		}
	}

	if len(osArgs) <= 1 || !isProxyCmd(cmdType(osArgs[1])) {
		// not a proxy command
		return false
	}

	// get un-proxied command
	rawCmd := unproxyCmd(cmdType(osArgs[1]))
	cmd := cmdLibrary.MustNewCmd(rawCmd.Name())

	// overwrite real "os.Args" with "osArgs" before calling the in-process runner
	os.Args = append(osArgs[:1], osArgs[2:]...)

	// run command in-process. Calls into the wrapped "main" function, so it is possible/likely that the call will
	// call os.Exit and terminate the program.
	cmdLibrary.Run(cmd)

	// if previous call completed, it means that it reached the end of the wrapped main method and os.Exit was not
	// called. This is assumed to mean successful execution.
	return true
}

type cmdType string

func (c cmdType) Name() string {
	return string(c)
}

func proxyCmd(cmd Cmd) Cmd {
	if isProxyCmd(cmd) {
		return cmd
	}
	return cmdType(ProxyCmdPrefix + cmd.Name())
}

func unproxyCmd(cmd Cmd) Cmd {
	if isProxyCmd(cmd) {
		return cmdType(cmd.Name()[len(ProxyCmdPrefix):])
	}
	return cmd
}

func isProxyCmd(cmd Cmd) bool {
	return strings.HasPrefix(cmd.Name(), ProxyCmdPrefix)
}
