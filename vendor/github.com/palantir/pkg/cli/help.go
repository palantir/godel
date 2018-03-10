// Copyright 2016 Palantir Technologies. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cli

import (
	"fmt"
	"io"
	"strings"
	"text/tabwriter"
	"text/template"

	"github.com/mitchellh/go-wordwrap"

	"github.com/palantir/pkg/cli/flag"
)

const appHelpTemplate = `NAME:
   {{.Name}} - {{.Usage}}

USAGE:
   {{.Name}}{{if .Flags}} [global flags]{{end}}{{if anyCommands .Subcommands}} command [command flags]{{end}}{{if .Description}}

DESCRIPTION:
   {{.Description}}{{end}}{{if anyCommands .Subcommands}}

COMMANDS:{{range .Subcommands}}{{if showCommand .}}
   {{commandHelp .}}{{end}}{{end}}{{end}}{{if .Flags}}

GLOBAL FLAGS:{{range .Flags}}{{if showFlag .}}
{{flagHelp .}}{{end}}{{end}}{{end}}{{if .Version}}

VERSION:
   {{.Version}}{{end}}

`

const commandHelpTemplate = `NAME:
   {{shortDescLine commandPath .Usage}}

USAGE:
   {{usageLine programName commandPath .Flags}}{{if .Description}}

DESCRIPTION:
   {{.Description}}{{end}}{{if anyCommands .Subcommands}}

SUBCOMMANDS:{{range .Subcommands}}{{if showCommand .}}
   {{commandHelp .}}{{end}}{{end}}{{end}}{{if .Flags}}

FLAGS:{{range .Flags}}{{if showFlag .}}
{{flagHelp .}}{{end}}{{end}}{{end}}

`

const decisionHelpTemplate = `NAME:
   {{shortDescLine commandPath .Usage}}

USAGE:{{range .Subcommands}}{{if showCommand .}}
   {{decisionLine programName commandPath decisionFlag .Name .Flags}}{{end}}{{end}}{{range .Subcommands}}{{if showCommand .}}

   {{.Usage}}

      USAGE:
         {{decisionLine programName commandPath decisionFlag .Name .Flags}}{{if .Description}}

      DESCRIPTION:
         {{.Description}}{{end}}{{if .Flags}}

      FLAGS:{{range .Flags}}{{if showFlag .}}
      {{flagHelp .}}{{end}}{{end}}{{end}}{{end}}{{end}}

`

func (ctx Context) PrintHelp(out io.Writer) {
	var tmpl string
	var data interface{}
	if len(ctx.Path) == 0 {
		tmpl = appHelpTemplate
		data = ctx.App
	} else if ctx.Command.DecisionFlag == "" {
		tmpl = commandHelpTemplate
		data = ctx.Command
	} else {
		tmpl = decisionHelpTemplate
		data = ctx.Command
	}

	funcMap := template.FuncMap{
		"usageLine":     usageLine,
		"shortDescLine": shortDescLine,
		"decisionLine":  decisionLine,
		"programName":   func() string { return ctx.App.Name },
		"commandPath":   func() string { return strings.Join(ctx.Path, " ") },
		"decisionFlag":  func() string { return ctx.Command.DecisionFlag },
		"anyCommands":   anyCommands,
		"showCommand":   showCommand,
		"commandHelp":   commandHelp,
		"showFlag":      showFlag,
		"flagHelp":      flagHelp,
	}

	minwidth := 0
	tabwidth := 8
	padding := 2
	padchar := byte(' ')
	flags := uint(0)
	w := tabwriter.NewWriter(out, minwidth, tabwidth, padding, padchar, flags)
	t := template.Must(template.New("help").Funcs(funcMap).Parse(tmpl))
	err := t.Execute(w, data)
	if err != nil {
		panic(err)
	}
	err = w.Flush()
	_ = err // ignore failure to flush
}

func usageLine(programName, commandPath string, flags []flag.Flag) string {
	str := programName + " " + commandPath
	str += flagsString(flags)

	width := uint(75 - len(programName))
	separator := "\\\n" + strings.Repeat(" ", len(programName)+4)

	return wrapWithSeparator(str, width, separator)
}

func decisionLine(programName, commandPath, decisionFlag, decisionName string, flags []flag.Flag) string {
	if decisionFlag != "" && decisionName != "" {
		commandPath += fmt.Sprintf(" %v %v", flag.WithPrefix(decisionFlag), decisionName)
	}
	return usageLine(programName, commandPath, flags)
}

// Create a string like " <servicename> --stack <stack>" containing all required flags.
func flagsString(flags []flag.Flag) string {
	str := ""
	hasOptional := false

	for _, f := range flags {
		if _, isSlice := f.(flag.StringSlice); isSlice {
			str += fmt.Sprintf(" <%s>...", f.MainName())
		} else if !f.HasLeader() {
			str += fmt.Sprintf(" <%s>", f.MainName())
		} else if f.IsRequired() {
			str += fmt.Sprintf(" %s <%s>", f.FullNames()[0], f.MainName())
		} else if f.MainName() != "help" {
			hasOptional = true
		}
	}

	if hasOptional {
		str += " [flags...]"
	}

	return str
}

func anyCommands(cmds []Command) bool {
	for _, cmd := range cmds {
		if showCommand(cmd) {
			return true
		}
	}
	return false
}

func showCommand(cmd Command) bool {
	return !strings.HasPrefix(cmd.Name, "_")
}

func commandHelp(cmd Command) string {
	names := strings.Join(cmd.Names(), ", ")

	// leave room within 80 characters for the name
	usage := wrapWithSeparator(cmd.Usage, uint(50), "\n\t")
	return names + "\t" + usage
}

func shortDescLine(usageLine string, usage string) string {
	usageShortDesc := fmt.Sprintf("%v - %v", usageLine, usage)
	width := uint(75)
	separator := "\n" + strings.Repeat(" ", 3)
	return wrapWithSeparator(usageShortDesc, width, separator)
}

func wrapWithSeparator(text string, width uint, separator string) string {
	return strings.Replace(wordwrap.WrapString(text, width), "\n", separator, -1)
}

func showFlag(f flag.Flag) bool {
	return f.DeprecationStr() == ""
}

func flagHelp(f flag.Flag) string {
	var names string
	if f.HasLeader() {
		names = strings.Join(f.FullNames(), ", ")
	} else {
		names = fmt.Sprintf("<%s>", f.MainName())
	}

	indent := "   "
	defaults := []string{}
	if f.IsRequired() {
		indent = " * "
	} else {
		if f.EnvVarStr() != "" {
			defaults = append(defaults, fmt.Sprintf("$%v", f.EnvVarStr()))
		}
		if f.DefaultStr() != "" {
			defaults = append(defaults, fmt.Sprintf("%q", f.DefaultStr()))
		}
	}

	usage := f.UsageStr()
	if len(defaults) > 0 {
		usage += " " + fmt.Sprintf(`(default=%v)`, strings.Join(defaults, ","))
	}

	usageWidth := uint(50) // leave room within 80 characters for the name
	usage = wordwrap.WrapString(usage, usageWidth)
	usage = strings.Replace(usage, "\n", "\n\t", -1)

	return indent + names + "\t" + usage
}
