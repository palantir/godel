// Copyright 2016 Palantir Technologies, Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package amalgomated

import (
	"fmt"
	"unicode"

	"github.com/pkg/errors"
)

// StringCmdSet is a set of commands that can be run that are represented as strings.
type StringCmdSet interface {
	Run(cmd string)
	Cmds() []string
}

// CmdWithRunner pairs a named command with the function for the command.
type CmdWithRunner struct {
	cmdName string
	runner  func()
}

func (c *CmdWithRunner) Name() string {
	return c.cmdName
}

// NewCmdWithRunner creates a new CmdWithRunner for the provided name and runner. Returns an error if the provided name
// is not a legal command name.
func NewCmdWithRunner(cmdName string, runner func()) (*CmdWithRunner, error) {
	if cmdName == "" {
		return nil, errors.New("cmdName cannot be blank")
	}

	for _, r := range cmdName {
		if unicode.IsSpace(r) {
			return nil, errors.Errorf("cmdName cannot contain whitespace: %q", cmdName)
		}
	}

	return &CmdWithRunner{
		cmdName: cmdName,
		runner:  runner,
	}, nil
}

// MustNewCmdWithRunner returns the result of NewCmdWithRunner and panics in cases where the function returns an error.
func MustNewCmdWithRunner(cmdName string, runner func()) *CmdWithRunner {
	cmdWithRunner, err := NewCmdWithRunner(cmdName, runner)
	if err != nil {
		panic(err)
	}
	return cmdWithRunner
}

type cmdWithRunnerCmdSet []*CmdWithRunner

func (s cmdWithRunnerCmdSet) Run(cmd string) {
	for _, curr := range s {
		if cmd == curr.cmdName {
			curr.runner()
			return
		}
	}
	panic(fmt.Sprintf("cmd %v not found in %v", cmd, s))
}

func (s cmdWithRunnerCmdSet) Cmds() []string {
	cmds := make([]string, len(s))
	for i := range s {
		cmds[i] = s[i].cmdName
	}
	return cmds
}

// NewStringCmdSetForRunners creates a new StringCmdSet from the provided cmds (all of which must be non-nil). Returns
// an error if any of the provided commands have the same name.
func NewStringCmdSetForRunners(cmds ...*CmdWithRunner) (StringCmdSet, error) {
	seenCmdNames := make(map[string]bool)
	var duplicates []string
	duplicatesSet := make(map[string]bool)
	for _, cmd := range cmds {
		cmdName := cmd.cmdName
		if seenCmdNames[cmdName] {
			if !duplicatesSet[cmdName] {
				duplicates = append(duplicates, cmdName)
			}
			duplicatesSet[cmdName] = true
		}
		seenCmdNames[cmdName] = true
	}
	if len(duplicates) > 0 {
		return nil, errors.Errorf("multiple runners provided for commands: %v", duplicates)
	}
	return cmdWithRunnerCmdSet(cmds), nil
}

// CmdLibrary represents a library of commands that can be run.
type CmdLibrary interface {
	// Run runs the provided command.
	Run(cmd Cmd)
	// Cmds returns the set of valid commands.
	Cmds() []Cmd
	// NewCmd creates a new Cmd for this library using the provided string. Returns an error if the provided string
	// does not correspond to a valid command for this library.
	NewCmd(cmd string) (Cmd, error)
	// MustNewCmd performs the same operation as NewCmd, but panics if NewCmd returns a non-nil error.
	MustNewCmd(cmd string) Cmd
}

type cmdLibraryImpl struct {
	cmdSet StringCmdSet
}

func NewCmdLibrary(cmdSet StringCmdSet) CmdLibrary {
	return &cmdLibraryImpl{
		cmdSet: cmdSet,
	}
}

func (c *cmdLibraryImpl) Run(cmd Cmd) {
	c.cmdSet.Run(cmd.Name())
}

func (c *cmdLibraryImpl) Cmds() []Cmd {
	stringCmds := c.cmdSet.Cmds()
	newCmds := make([]Cmd, len(stringCmds))
	for i := range stringCmds {
		newCmds[i] = c.MustNewCmd(stringCmds[i])
	}
	return newCmds
}

func (c *cmdLibraryImpl) NewCmd(cmd string) (Cmd, error) {
	found := false
	cmds := c.cmdSet.Cmds()
	for _, currCmd := range cmds {
		if cmd == currCmd {
			found = true
			break
		}
	}

	if !found {
		return nil, fmt.Errorf("invalid command %q (valid values: %v)", cmd, cmds)
	}

	return cmdImpl(cmd), nil
}

func (c *cmdLibraryImpl) MustNewCmd(cmd string) Cmd {
	newCmd, err := c.NewCmd(cmd)
	if err != nil {
		panic(err)
	}
	return newCmd
}
