// Copyright 2016 Palantir Technologies. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cli

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"syscall"
	"text/template"

	"github.com/palantir/pkg/cli/flag"
)

const descriptionText = `
Register completion for {{appName}} by running one of the following in a
terminal, replacing [program] with the command you use to run {{appName}}{{if ne appName program}}
(eg. {{program}}){{end}}.

# Bash ~4.x
source <([program] _completion)

# Bash ~3.x
[program] _completion | source /dev/stdin

# Bash (any version)
eval "$([program] _completion)"

By default this registers completion for the absolute path to {{appName}},
which will work if {{appName}} is accessible on your PATH. You can specify
a program name to complete for instead using the --{{progFlag}} flag, which is
required if you are using an alias or wrapper script to run {{appName}}.
If you are using an alias, you also need to provide --{{aliasFlag}} with
the full contents of the alias.
`

const (
	progFlagName  = "prog"
	aliasFlagName = "alias"
	bashFlagName  = "bash"
	zshFlagName   = "zsh"
)

func completionCommand(appName string) Command {
	return Command{
		Name:        "_completion",
		Usage:       "Enable command line completion",
		Description: renderDescription(appName),
		Action:      printCompletionScript,
		Flags: []flag.Flag{
			flag.StringFlag{Name: progFlagName, Value: os.Args[0], Usage: "program to complete"},
			flag.StringFlag{Name: aliasFlagName, Usage: "if prog is an alias, the contents of the alias"},
			flag.BoolFlag{Name: bashFlagName, Usage: "assume shell is Bash"},
			flag.BoolFlag{Name: zshFlagName, Usage: "assume shell is Zsh"},
		},
	}
}

func renderDescription(appName string) string {
	funcMap := template.FuncMap{
		"appName":   func() string { return appName },
		"program":   func() string { return filepath.Base(os.Args[0]) },
		"progFlag":  func() string { return progFlagName },
		"aliasFlag": func() string { return aliasFlagName },
	}
	t := template.Must(template.New("description").Funcs(funcMap).Parse(descriptionText))

	var descriptionBuffer bytes.Buffer
	err := t.Execute(&descriptionBuffer, nil)
	if err != nil {
		panic(err)
	}
	description := descriptionBuffer.String()

	// Indent
	description = strings.Replace(description, "\n", "\n   ", -1)
	description = strings.Trim(description, "\n ")
	return description
}

func printCompletionScript(ctx Context) error {
	if ctx.Bool(bashFlagName) {
		printProvidedCompletionScript(ctx, bashCompletionScript)
		return nil
	}
	if ctx.Bool(zshFlagName) {
		printProvidedCompletionScript(ctx, zshCompletionScript)
		return nil
	}

	if ctx.IsTerminal() {
		ctx.PrintHelp(ctx.App.Stdout)
		return nil
	}

	shell, ok := syscall.Getenv("SHELL")
	if !ok {
		command := filepath.Base(os.Args[0])
		return fmt.Errorf(
			"Unable to detect shell. Use `%s _completion --%s` to assume Bash, or `%s _completion --%s` to assume Zsh.",
			command, bashFlagName, command, zshFlagName)
	}
	shell = filepath.Base(shell)
	shell = strings.ToLower(shell)

	switch shell {
	case "bash", "logging_bash":
		printProvidedCompletionScript(ctx, bashCompletionScript)
	case "zsh":
		printProvidedCompletionScript(ctx, zshCompletionScript)
	default:
		command := filepath.Base(os.Args[0])
		return fmt.Errorf(
			"Unsupported shell %q. Use `%s _completion --%s` to assume Bash, or `%s _completion --%s` to assume Zsh.",
			shell, command, bashFlagName, command, zshFlagName)
	}

	return nil
}

func printProvidedCompletionScript(ctx Context, script string) {
	random := fmt.Sprint(rand.Uint32())
	funcMap := template.FuncMap{
		"osArgs": func() string {
			return strings.Join(os.Args, " ")
		},
		"random": func() string {
			return random
		},
		"prog": func() string {
			return fmt.Sprintf("%q", ctx.String(progFlagName))
		},
		"execute": func() string {
			if ctx.String(aliasFlagName) != "" {
				return ctx.String(aliasFlagName)
			}
			return fmt.Sprintf("%q", ctx.String(progFlagName))
		},
	}
	t := template.Must(template.New("completion").Funcs(funcMap).Parse(script))
	err := t.Execute(ctx.App.Stdout, nil)
	if err != nil {
		panic(err)
	}
}

const bashCompletionScript = `_completion_{{random}}() {
	local cur opts
	COMPREPLY=()
	cur="${COMP_WORDS[COMP_CWORD]}"
	opts=$({{execute}} "${COMP_WORDS[@]:1:COMP_CWORD}" --generate-bash-completion)
	case $? in
	99) # complete file paths
		type compopt &>/dev/null && compopt -o filenames
		local IFS=$'\n'
		COMPREPLY=($(compgen -f -- "${cur}"))
		;;
	98) # complete directories
		type compopt &>/dev/null && compopt -o filenames
		local IFS=$'\n'
		COMPREPLY=($(compgen -d -- "${cur}"))
		;;
	97) # custom filepaths
		type compopt &>/dev/null && compopt -o filenames
		local IFS=$'\n'
		read -rd '' -a COMPREPLY <<< "${opts}"
		if [[ ${#COMPREPLY[@]} -eq 1 && "${COMPREPLY[0]: -1}" == "/" ]]
		then
			type compopt &>/dev/null && compopt -o nospace
		fi
		;;
	*)
		local IFS=$'\n'
		COMPREPLY=($(compgen -W "${opts}" -- "${cur}"))
		;;
	esac
	return 0
}
complete -F _completion_{{random}} {{prog}}
# Run 'eval "$({{osArgs}})"' to apply
`
const zshCompletionScript = `_completion_{{random}}() {
    local curcontext="$curcontext" state line
    typeset -A opt_args
    opts=$({{execute}} "${words[@]:1:$CURRENT-1}" --generate-bash-completion)
    case $? in
      99) # complete file paths
        _path_files
        ;;
      98) # complete directories
        _path_files -/
        ;;
      97) # custom filepaths
        # ignore '@' for custom filepath completion
        if compset -P '@'; then
          _files
        fi
        ;;
      *)
        args=("${(f)opts}")
        if [[ ${#args[@]} -eq 2 && "${args[1]}" == "@" ]]; then
          # if completing '@', use "-S ''" to ensure that space is not
          # added after '@' is chosen so that file path completion can occur
          compadd -S '' -Q -a args
        elif [[ ${#args[@]} -eq 1 && "${args[1]}" == "" ]]; then
          # don't add any completions if args contains only one
          # element that is the empty string
          :
        else
          compadd -Q -a args
        fi
        ;;
    esac
    return 0
}
compdef _completion_{{random}} {{prog}}
# Run 'eval "$({{osArgs}})"' to apply
`
