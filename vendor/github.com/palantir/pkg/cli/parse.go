// Copyright 2016 Palantir Technologies. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cli

import (
	"fmt"
	"strings"
	"syscall"

	"golang.org/x/crypto/ssh/terminal"

	"github.com/palantir/pkg/cli/flag"
)

func (app *App) parse(args []string) (Context, error) {
	ctx := Context{
		App:     app,
		Command: &app.Command,
		IsTerminal: func() bool {
			return terminal.IsTerminal(int(uintptr(syscall.Stdout)))
		},
		defaults:  map[string]interface{}{},
		specified: map[string]interface{}{},
		allVals:   map[string][]interface{}{},
	}

	args = args[1:] // skip name of binary
	fillDefaults(&ctx)
	for {
		applyBackcompat(&ctx, &args)
		more, err := parseNextFlag(&ctx, &args)
		if err != nil {
			return ctx, err
		}
		if len(ctx.Path) == 0 && ctx.Bool(versionFlag.MainName()) {
			return ctx, nil // --version, stop parsing
		}
		if !more {
			return ctx, checkRequiredFlags(&ctx)
		}
	}
}

func applyBackcompat(ctx *Context, args *[]string) {
	if len(ctx.Path) > 0 {
		return
	}
	for _, backcompat := range ctx.App.Backcompat {
		if len(backcompat.Path) > len(*args) {
			continue
		}
		startsWithPath := true
		for i := range backcompat.Path {
			if backcompat.Path[i] != (*args)[i] {
				startsWithPath = false
				break
			}
		}
		if startsWithPath {
			*args = (*args)[len(backcompat.Path):]
			ctx.Command = &backcompat.Command
			ctx.Path = backcompat.Path
			fillDefaults(ctx)
			return
		}
	}
}

func fillDefaults(ctx *Context) {
	for _, f := range ctx.Command.Flags {
		if !f.IsRequired() {
			ctx.defaults[f.MainName()] = f.Default()
		}
	}
}

func parseNextFlag(ctx *Context, args *[]string) (more bool, err error) {
	// TODO: split up enormous function
	if len(*args) == 0 {
		return false, nil
	}
	first := pop(args)
	hasHyphen := strings.HasPrefix(first, "-")

	// look for a flag
	for _, f := range ctx.Command.Flags {
		if !f.HasLeader() {
			if hasHyphen || ctx.Has(f.MainName()) {
				continue
			}
			if _, isSlice := f.(flag.StringSlice); isSlice {
				strs := []string{first}
				for len(*args) > 0 && !strings.HasPrefix((*args)[0], "-") {
					strs = append(strs, pop(args))
				}
				ctx.specified[f.MainName()] = strs
			} else {
				ctx.specified[f.MainName()], err = f.Parse(first)
				if err != nil {
					return false, fmt.Errorf("%v: %v", f.MainName(), err)
				}
			}
			return true, nil
		}
		for _, name := range f.FullNames() {
			equalsIndex := -1
			// if argument matches "flag=", parse it as a flag
			if strings.HasPrefix(first, name+"=") {
				equalsIndex = len(name)
				// if '=' occurs in the flag, split on the first occurrence and use the value after the '=' as the next argument
				next := first[equalsIndex+1:]
				first = first[:equalsIndex]
				if next == "" {
					// if value after "=" is empty, treat it as a missing value
					return false, fmt.Errorf("Missing value for flag %v", first)
				}
				push(next, args)
			}

			if first != name {
				continue
			}
			if f, ok := f.(flag.BoolFlag); ok && equalsIndex == -1 {
				// if the flag is a boolean flag without an '=', interpret as "true"
				ctx.specified[f.MainName()] = true
				return true, nil
			}
			if len(*args) == 0 {
				return false, fmt.Errorf("Missing value for flag %v", first)
			}
			ctx.specified[f.MainName()], err = f.Parse(pop(args))
			if err != nil {
				return false, fmt.Errorf("%v: %v", name, err)
			}
			ctx.allVals[f.MainName()] = append(ctx.allVals[f.MainName()], ctx.specified[f.MainName()])
			return true, nil
		}
	}

	decisionFlag := ctx.Command.DecisionFlag
	isDecision := decisionFlag != "" && first == flag.WithPrefix(decisionFlag)
	if isDecision {
		if len(*args) == 0 {
			// be consistent with missing flag message above
			return false, fmt.Errorf("Missing value for flag %v", first)
		}
		ctx.Path = append(ctx.Path, first)
		first = pop(args)
	} else if decisionFlag != "" {
		push(first, args)
		return defaultDecision(ctx, args)
	} else if hasHyphen {
		return false, fmt.Errorf("Unknown flag %v", first)
	}

	// look for a subcommand
	for _, cmd := range ctx.Command.Subcommands {
		for _, name := range cmd.Names() {
			if first != name {
				continue
			}
			ctx.Command = &cmd
			ctx.Path = append(ctx.Path, cmd.Names()[0])
			fillDefaults(ctx)
			return true, nil
		}
	}

	push(first, args)
	if isDecision {
		return badDecision(ctx, *args)
	}
	return matchedNothing(ctx, *args)
}

func badDecision(ctx *Context, args []string) (more bool, err error) {
	var cmdNames []string
	for _, cmd := range ctx.Command.Subcommands {
		if len(cmd.Names()) > 0 {
			cmdNames = append(cmdNames, cmd.Names()[0])
		}
	}
	var expected string
	if len(cmdNames) == 1 {
		expected = fmt.Sprintf("%q", cmdNames[0])
	} else {
		expected = fmt.Sprintf("one of [%v]", strings.Join(cmdNames, ", "))
	}
	return false, fmt.Errorf("Value of %v must be %v, not %q",
		flag.WithPrefix(ctx.Command.DecisionFlag), expected, args[0])
}

func matchedNothing(ctx *Context, args []string) (more bool, err error) {
	if len(ctx.Command.Subcommands) > 0 {
		return false, fmt.Errorf("Unknown command %q", args[0])
	}
	// TODO: conditionally quote the arguments like in cli/cli.go
	return false, fmt.Errorf("Too many arguments: %v", strings.Join(args, " "))
}

func defaultDecision(ctx *Context, args *[]string) (more bool, err error) {
	for _, cmd := range ctx.Command.Subcommands {
		if cmd.Name == DefaultDecision {
			ctx.Command = &cmd
			fillDefaults(ctx)
			return true, nil
		}
	}
	// there is no default, decision flag is required
	return false, fmt.Errorf("Expected %v after %q",
		flag.WithPrefix(ctx.Command.DecisionFlag),
		strings.Join(ctx.Path, " "))
}

func checkRequiredFlags(ctx *Context) error {
	var missing []string
	for _, f := range ctx.Command.Flags {
		if f.IsRequired() && !ctx.Has(f.MainName()) {
			missing = append(missing, f.FullNames()[0])
		}
	}
	if len(missing) == 1 {
		return fmt.Errorf("Flag %s is required.", missing[0])
	} else if len(missing) > 1 {
		return fmt.Errorf("Flags %s are required.", strings.Join(missing, ", "))
	}
	return nil
}

func pop(args *[]string) string {
	first := (*args)[0]
	*args = (*args)[1:]
	return first
}

func push(arg string, args *[]string) {
	*args = append([]string{arg}, *args...)
}
