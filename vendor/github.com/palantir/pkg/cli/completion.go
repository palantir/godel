// Copyright 2016 Palantir Technologies. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cli

import (
	"fmt"
	"strings"

	"github.com/palantir/pkg/cli/completion"
	"github.com/palantir/pkg/cli/flag"
)

const CompletionFlag = "--generate-bash-completion"

var none = []string{}

// Returns exit code if completion happened, or -1 to continue app.
func (app *App) doCompletion(args []string) (exitCode int) {
	if hasCompletionFlag(args) {
		completions, exitCode := doCompletion(app, args[1:])
		for _, completion := range completions {
			fmt.Fprintln(app.Stdout, completion)
		}
		return int(exitCode)
	}
	app.Subcommands = append(app.Subcommands, completionCommand(app.Name))
	return -1
}

// More testable than the function above.
func doCompletion(app *App, args []string) (completions []string, exitCode uint8) {
	defer func() {
		// Ignore panics
		if r := recover(); r != nil {
			completions = none
		}
	}()
	args = stripSlashes(args[:len(args)-1])
	return generateCompletions(app, args)
}

func hasCompletionFlag(args []string) bool {
	if len(args) < 3 {
		// first argument is the binary
		// second-last argument must be partial word
		// last argument must be "--generate-bash-completion"
		return false
	}
	return args[len(args)-1] == CompletionFlag
}

func stripSlashes(args []string) []string {
	// TODO: interpret escape sequences and quotation
	for i := range args {
		args[i] = strings.Replace(args[i], "\\", "", -1)
	}
	return args
}

type completionState struct {
	path        []string
	command     *Command
	flags       map[string]string
	currentFlag flag.Flag
	isDecision  bool
	providers   map[string]completion.Provider
	exitCode    uint8
}

func generateCompletions(app *App, args []string) (completions []string, exitCode uint8) {
	state := initialState(app)
	for _, arg := range args[:len(args)-1] {
		if strings.HasPrefix(arg, "#") {
			return none, 0
		}
		if state.currentFlag != nil && strings.HasPrefix(arg, "-") {
			if _, isSlice := state.currentFlag.(flag.StringSlice); isSlice {
				state.currentFlag = nil // end of values for this slice
			} else {
				return failCompletion(args), 0
			}
		}
		if state.handleCurrentFlag(arg) {
			continue
		}
		if state.handleCommand(arg) {
			continue
		}
		if state.isDecision { // should have been handled in handleCommand
			return failCompletion(args), 0
		}
		if state.handleDecision(arg) {
			continue
		}
		if state.handleFlag(arg) {
			continue
		}
		return failCompletion(args), 0
	}
	completions = state.generateCompletionsWithPrefix(args[len(args)-1])
	return completions, state.exitCode
}

func initialState(app *App) *completionState {
	return &completionState{
		path:        []string{},
		command:     &app.Command,
		flags:       map[string]string{},
		currentFlag: nil,
		isDecision:  false,
		providers:   app.Completion,
		exitCode:    0,
	}
}

func failCompletion(args []string) (completions []string) {
	if args[len(args)-2] != "#-[" {
		return []string{"#-["}
	}
	return none
}

func (state *completionState) handleCurrentFlag(arg string) bool {
	if state.currentFlag == nil {
		return false
	}
	if _, isSlice := state.currentFlag.(flag.StringSlice); isSlice {
		state.flags[state.currentFlag.MainName()] += "\x00" + arg
	} else {
		state.flags[state.currentFlag.MainName()] = arg
		state.currentFlag = nil
	}
	return true
}

func (state *completionState) handleCommand(arg string) bool {
	for _, cmd := range state.command.Subcommands {
		for _, name := range cmd.Names() {
			if arg == name {
				if state.isDecision {
					state.flags[state.command.DecisionFlag] = cmd.Name
					state.isDecision = false
				} else {
					state.path = append(state.path, cmd.Name)
				}
				state.command = &cmd
				return true
			}
		}
	}
	return false
}

func (state *completionState) handleDecision(arg string) bool {
	if state.command.DecisionFlag == "" {
		return false
	}
	if arg == flag.WithPrefix(state.command.DecisionFlag) {
		state.isDecision = true
		return true
	}
	state.command = &state.command.Subcommands[0]
	return false
}

func (state *completionState) handleFlag(arg string) bool {
	for _, f := range state.command.Flags {
		if !f.HasLeader() && !strings.HasPrefix(arg, "-") {
			if _, ok := state.flags[f.MainName()]; ok {
				continue // already have a value for this flag
			}
			state.flags[f.MainName()] = arg
			if _, isSlice := f.(flag.StringSlice); isSlice {
				state.currentFlag = f // there may be more values
			}
			return true
		}
		for _, name := range f.FullNames() {
			if arg == name {
				if bf, isbool := f.(flag.BoolFlag); isbool {
					state.flags[bf.MainName()] = "true"
				} else {
					state.currentFlag = f
				}
				return true
			}
		}
	}
	return false
}

func (state *completionState) generateCompletionsWithPrefix(prefix string) (completions []string) {
	var allCompletions []string
	if _, isSlice := state.currentFlag.(flag.StringSlice); isSlice {
		// can either be another element in slice, or the next flag
		slice := state.completeFlagValues(state.currentFlag, prefix)
		flags := state.generateFlagsWithPrefix(prefix)
		if len(flags) > 0 {
			state.exitCode = 0 // show other options even if slice is trying to complete filepaths
		}
		allCompletions = append(slice, flags...)
	} else if state.currentFlag != nil {
		allCompletions = state.completeFlagValues(state.currentFlag, prefix)
	} else if state.isDecision {
		allCompletions = state.generateCommandsWithPrefix(prefix)
	} else if state.isHelpOrVersion() {
		return []string{"#-]"}
	} else if state.command.DecisionFlag != "" {
		allCompletions = state.generateDecisionsWithPrefix(prefix)
	} else {
		commands := state.generateCommandsWithPrefix(prefix)
		flags := state.generateFlagsWithPrefix(prefix)
		allCompletions = append(commands, flags...)
	}
	return withPrefix(allCompletions, prefix)
}

func (state *completionState) isHelpOrVersion() bool {
	if state.flags[helpFlag.MainName()] == "true" {
		return true // --help
	}
	if len(state.path) == 0 && state.flags[versionFlag.MainName()] == "true" {
		return true // --version
	}
	return false
}

func (state *completionState) generateCommandsWithPrefix(prefix string) (completions []string) {
	for _, cmd := range state.command.Subcommands {
		names := cmd.Names()
		if len(names) == 0 {
			continue
		} else if strings.HasPrefix(names[0], "_") && prefix == "" {
			// do not complete commands starting with underscore unless user started typing it
			continue
		} else if strings.HasPrefix(names[0], prefix) {
			// complete full name only unless prefix matches alias but not full name
			completions = append(completions, names[0])
		} else {
			// prefix is nonempty and does not match full name, so complete aliases
			for _, alias := range names[1:] {
				if strings.HasPrefix(alias, prefix) {
					completions = append(completions, alias)
				}
			}
		}
	}
	return completions
}

func (state *completionState) generateDecisionsWithPrefix(prefix string) (completions []string) {
	var defaultCmd *Command
	for _, cmd := range state.command.Subcommands {
		if cmd.Name == DefaultDecision {
			defaultCmd = &cmd
			break
		}
	}
	decisionFlag := flag.WithPrefix(state.command.DecisionFlag)
	if defaultCmd == nil {
		// decision flag is required
		return []string{decisionFlag}
	}
	// get completions for default command which is when decision flag is not given
	state.command = defaultCmd
	completions = state.generateCompletionsWithPrefix(prefix)
	if !strings.HasPrefix(decisionFlag, prefix) {
		// it has been narrowed down to the default command
		return completions
	}
	// may be either default command or decision flag
	state.exitCode = 0 // complete decisionFlag even if defaultCmd is trying to complete filepaths
	return append(completions, decisionFlag)
}

func (state *completionState) generateFlagsWithPrefix(prefix string) (completions []string) {
	if prefix == "" {
		firstRequired := state.generateFirstRequiredFlag(prefix)
		if firstRequired != nil {
			return firstRequired
		}
		if len(state.command.Subcommands) == 0 {
			completions = []string{"#-]"}
		}
	}
	for _, f := range state.missingFlags() {
		if !f.HasLeader() {
			completions = append(completions, state.completeFlagValues(f, prefix)...)
		} else if strings.HasPrefix(f.FullNames()[0], prefix) {
			completions = append(completions, f.FullNames()[0])
		} else {
			// Generate aliases only if the user has started typing it
			for _, alias := range f.FullNames() {
				if strings.HasPrefix(alias, prefix) {
					completions = append(completions, alias)
				}
			}
		}
	}
	// Special case for help flag - complete only if user has started typing
	// something that matches nothing else
	if len(completions) == 0 && prefix != "" && state.flags[helpFlag.MainName()] == "" {
		if strings.HasPrefix(helpFlag.FullNames()[0], prefix) {
			completions = []string{helpFlag.FullNames()[0]}
		} else {
			completions = []string{flag.WithPrefix(helpFlag.Alias)}
		}
	}
	return withPrefix(completions, prefix)
}

func (state *completionState) generateFirstRequiredFlag(prefix string) (completions []string) {
	for _, f := range state.missingFlags() {
		if !f.IsRequired() {
			continue
		}
		if f.HasLeader() {
			return []string{f.FullNames()[0]}
		}
		return state.completeFlagValues(f, prefix)
	}
	return nil
}

func (state *completionState) missingFlags() []flag.Flag {
	missing := []flag.Flag{}
	for _, f := range state.command.Flags {
		if f == helpFlag {
			continue
		}
		if !showFlag(f) {
			continue
		}
		if _, ok := state.flags[f.MainName()]; !ok {
			missing = append(missing, f)
		}
	}
	return missing
}

func (state *completionState) completeFlagValues(f flag.Flag, prefix string) (completions []string) {
	provider := state.providers[f.MainName()]
	if provider == nil {
		return completeFlagPlaceholder(f)
	}
	ctx := &completion.ProviderCtx{
		Flags:    state.flags,
		Command:  state.path,
		Partial:  prefix,
		ExitCode: &state.exitCode,
	}
	completions = provider(ctx)
	if len(completions) == 0 {
		return completeFlagPlaceholder(f)
	}
	return completions
}

const nbsp = "\u00A0"

func completeFlagPlaceholder(f flag.Flag) (completions []string) {
	// Also return nbsp, an invisible completion that sorts after the flag name
	// but prevents  bash from autofilling the completion
	return []string{f.PlaceholderStr(), nbsp}
}

func withPrefix(completions []string, prefix string) []string {
	var completionsWithPrefix []string
	for _, completion := range completions {
		if strings.HasPrefix(completion, prefix) {
			completionsWithPrefix = append(completionsWithPrefix, completion)
		}
	}
	return completionsWithPrefix
}
