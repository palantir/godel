// Copyright 2016 Palantir Technologies, Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package amalgomated

import (
	"os/exec"

	"github.com/kardianos/osext"
	"github.com/pkg/errors"
)

// Cmder creates an *exec.Cmd that can be run.
type Cmder interface {
	// Cmd returns the *exec.Cmd object configured with the provided arguments and with possible
	// implementation-specific augmentations. The returned command has not yet been executed or run, and the caller
	// may change the configuration as desired before executing the command. The returned command is ready to call,
	// and a standard use case is to call CombinedOutput() on the returned command to execute it and retrieve the
	// output generated to stdOut and stdErr.
	Cmd(args []string, cmdWd string) *exec.Cmd
}

// CmderSupplier returns the Cmder that runs the specified command. Returns an error if a Runner cannot be created for
// the requested Cmd.
type CmderSupplier func(cmd Cmd) (Cmder, error)

// Cmd represents a command that should be run. It is used as a key provided to a CmderSupplier to generate a Cmder.
type Cmd interface {
	Name() string
}

type cmdImpl string

func (c cmdImpl) Name() string {
	return string(c)
}

// SelfProxyCmderSupplier returns a supplier that, given a command, re-invokes the current executable with a proxy
// version of the provided command.
func SelfProxyCmderSupplier() CmderSupplier {
	return func(cmd Cmd) (Cmder, error) {
		selfCmder, err := selfCmder()
		if err != nil {
			return nil, err
		}
		return CmderWithPrependedArgs(selfCmder, proxyCmd(cmd).Name()), nil
	}
}

// SupplierWithPrependedArgs returns a new Supplier that invokes the provided supplier and returns the result of calling
// RunnerWithPrependedArgs on the returned runner with the result of applying the provided "argsFunc" function to the
// provided command.
func SupplierWithPrependedArgs(s CmderSupplier, argsFunc func(cmd Cmd) []string) CmderSupplier {
	return func(cmd Cmd) (Cmder, error) {
		r, err := s(cmd)
		if err != nil {
			return nil, err
		}
		return CmderWithPrependedArgs(r, argsFunc(cmd)...), nil
	}
}

type runner struct {
	pathToExecutable string
	prependedArgs    []string
}

func (r *runner) Cmd(args []string, cmdWd string) *exec.Cmd {
	cmd := exec.Command(r.pathToExecutable, append(r.prependedArgs, args...)...)
	cmd.Dir = cmdWd
	return cmd
}

type wrappedCmder struct {
	inner         Cmder
	prependedArgs []string
}

func (r *wrappedCmder) Cmd(args []string, cmdWd string) *exec.Cmd {
	combinedArgs := append(append(make([]string, 0, len(r.prependedArgs)+len(args)), r.prependedArgs...), args...)
	return r.inner.Cmd(combinedArgs, cmdWd)
}

// selfCmder returns a Cmder that creates a command that re-invokes the currently running executable.
func selfCmder() (Cmder, error) {
	pathToSelf, err := osext.Executable()
	if err != nil {
		return nil, errors.Wrapf(err, "failed to determine path for current executable")
	}
	return PathCmder(pathToSelf), nil
}

// CmderWithPrependedArgs returns a new Cmder that invokes the provided Cmder, but always adds the provided
// "prependedArgs" before any user-supplied arguments. Note that if the runner being wrapped has a notion of
// "prependedArgs" itself, those arguments will precede the "prependedArgs" provided in this method.
func CmderWithPrependedArgs(r Cmder, prependedArgs ...string) Cmder {
	return &wrappedCmder{
		inner:         r,
		prependedArgs: prependedArgs,
	}
}

// PathCmder returns a Cmder that runs the command at the supplied path with the specified "prependedArgs" provided as
// arguments to the executable before all of the other arguments that are added. The path should resolve to the
// executable that should be run: for example, "/usr/bin/git". The "prependedArgs" will be combined with the arguments
// provided to the Run method for any given execution.
func PathCmder(pathToExecutable string, prependedArgs ...string) Cmder {
	return &runner{
		pathToExecutable: pathToExecutable,
		prependedArgs:    prependedArgs,
	}
}
