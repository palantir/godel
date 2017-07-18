// Copyright 2016 Palantir Technologies. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cli

import (
	"fmt"
	"strconv"
	"strings"
	"text/template"
	"time"

	"github.com/palantir/pkg/cli/flag"
)

func addManpageCommand(app *App) {
	if app.Manpage != nil {
		app.Subcommands = append(app.Subcommands, Command{
			Name:   "_manpage",
			Usage:  "Generate manpage",
			Action: printManpage,
		})
	}
}

const manpageHeader = `\
.TH {{quote (upper .Name)}} 1 {{quote date}} {{quote .Manpage.Source}} {{quote .Manpage.Manual}}
.nh
`

const manpageName = `\
.SH NAME
{{.Name}} \- {{lower .Usage}}
`

const manpageSynopsis = `\
.SH SYNOPSIS
.B {{.Name}}
{{if .Flags}}\
	.RI [ options ]
{{end}}\
{{if .Subcommands}}\
	.I subcommand
	.RI [ arguments ]
{{end}}\
`

const manpageDescription = `\
{{if .Description}}\
	.SH DESCRIPTION
	{{.Description}}
{{end}}\
`

const manpageSubcommands = `\
{{if .Subcommands}}\
	.SH SUBCOMMANDS
	{{range allSubcommands .}}\
		.TP
		{{.Path}}{{allFlags .Sub.Flags}}
		.RS
		{{.Sub.Usage}}
		{{if .Sub.Description}}\
			.PP
			{{.Sub.Description}}
		{{end}}\
		{{range .Sub.Flags}}\
			{{if and (ne .MainName "help") (showFlag .)}}\
				.TP
				{{flag .}}
				.RS
				{{flagUsage .}}
				.RE
			{{end}}\
		{{end}}\
		.RE
	{{end}}\
{{end}}\
`

const manpageOptions = `\
{{if .Flags}}\
	.SH OPTIONS
	{{range .Flags}}\
		{{if showFlag .}}\
			.TP
			{{flag .}}
			.RS
			{{flagUsage .}}
			.RE
		{{end}}\
	{{end}}\
{{end}}\
`

const manpageSeeAlso = `\
{{if .Manpage.SeeAlso}}\
	.SH SEE ALSO
	{{range $i, $m := .Manpage.SeeAlso}}\
		{{if $i}}\
			{{",\n"}}\
		{{end}}\
		.BR {{quote $m.Name}} ({{$m.Section}})\
	{{end}}
{{end}}\
`

const manpageBugTracker = `\
{{if .Manpage.BugTracker}}\
	.SH BUGS
	.ad l
	Presumably. Report them or discuss them at
	\%{{.Manpage.BugTracker}}
{{end}}\
`

func printManpage(ctx Context) error {
	funcMap := template.FuncMap{
		"quote":          strconv.Quote,
		"upper":          strings.ToUpper,
		"lower":          strings.ToLower,
		"date":           func() string { return time.Now().UTC().Format("2006-01-02") },
		"allSubcommands": allSubcommands,
		"allFlags":       renderAllFlags,
		"flag":           renderFlag,
		"flagUsage":      renderFlagUsage,
		"showFlag":       showFlag,
	}

	fullTemplate := strings.Join([]string{
		manpageHeader,
		manpageName,
		manpageSynopsis,
		manpageDescription,
		manpageSubcommands,
		manpageOptions,
		manpageSeeAlso,
		manpageBugTracker,
	}, "")
	fullTemplate = strings.Replace(fullTemplate, "\t", "", -1)
	fullTemplate = strings.Replace(fullTemplate, "\\\n", "", -1)

	t := template.Must(template.New("manpage").Funcs(funcMap).Parse(fullTemplate))
	if err := t.Execute(ctx.App.Stdout, ctx.App); err != nil {
		panic(err)
	}
	return nil
}

type subcommandAt struct {
	Path string
	Sub  Command
}

func allSubcommands(app *App) []subcommandAt {
	var cmds []subcommandAt
	dfsSubcommands(app.Command, nil, &cmds)
	return cmds
}

func dfsSubcommands(cmd Command, path []string, cmds *[]subcommandAt) {
	if len(cmd.Subcommands) == 0 {
		*cmds = append(*cmds, subcommandAt{
			Path: strings.Join(path, " "),
			Sub:  cmd,
		})
	} else {
		for _, sub := range cmd.Subcommands {
			if showCommand(sub) {
				subpath := append(path, renderSubcommandName(cmd, sub)...)
				dfsSubcommands(sub, subpath, cmds)
			}
		}
	}
}

func renderSubcommandName(parent, sub Command) []string {
	if parent.DecisionFlag != "" {
		df := flag.WithPrefix(parent.DecisionFlag)
		if sub.Name != "" {
			return []string{fmt.Sprintf(`\fB%s %s\fP`, df, sub.Name)}
		}
		if sub.Alias != "" {
			return []string{fmt.Sprintf(`[\fB%s %s\fP]`, df, sub.Alias)}
		}
		return []string{}
	}
	if sub.Alias != "" {
		return []string{fmt.Sprintf(`(\fB%s\fP | \fB%s\fP)`, sub.Name, sub.Alias)}
	}
	return []string{fmt.Sprintf(`\fB%s\fP`, sub.Name)}
}

func renderAllFlags(flags []flag.Flag) string {
	str := ""
	var optionalFlags []flag.Flag
	var onlyFlags []flag.Flag

	for _, f := range flags {
		if !showFlag(f) {
			continue
		} else if _, isSlice := f.(flag.StringSlice); isSlice {
			str += fmt.Sprintf(` \fI%s\fP...`, f.MainName())
		} else if !f.HasLeader() {
			str += fmt.Sprintf(` \fI%s\fP`, f.MainName())
		} else if f.IsRequired() {
			placeholder := strings.ToLower(f.PlaceholderStr())
			str += fmt.Sprintf(` \fB%s\fP \fI%s\fP`, f.FullNames()[0], placeholder)
		} else if _, isbool := f.(flag.BoolFlag); isbool && strings.HasSuffix(f.MainName(), "-only") {
			onlyFlags = append(onlyFlags, f)
		} else if f.MainName() != "help" {
			optionalFlags = append(optionalFlags, f)
		}
	}

	for _, f := range optionalFlags {
		if _, isbool := f.(flag.BoolFlag); isbool {
			str += fmt.Sprintf(` [\fB%s\fP]`, f.FullNames()[0])
		} else {
			placeholder := strings.ToLower(f.PlaceholderStr())
			str += fmt.Sprintf(` [\fB%s\fP \fI%s\fP]`, f.FullNames()[0], placeholder)
		}
	}

	if len(onlyFlags) > 0 {
		boldOnlyFlags := make([]string, 0, len(onlyFlags))
		for _, onlyFlag := range onlyFlags {
			boldOnlyFlags = append(boldOnlyFlags, fmt.Sprintf(`\fB%s\fP`, onlyFlag.FullNames()[0]))
		}
		str += fmt.Sprintf(` [%s]`, strings.Join(boldOnlyFlags, " | "))
	}

	return str
}

func renderFlag(f flag.Flag) string {
	if !f.HasLeader() {
		return fmt.Sprintf(`<\fI%s\fP>`, f.MainName())
	}

	boldNames := make([]string, 0, len(f.FullNames()))
	for _, name := range f.FullNames() {
		boldNames = append(boldNames, fmt.Sprintf(`\fB%s\fP`, name))
	}

	str := strings.Join(boldNames, ", ")
	if _, isbool := f.(flag.BoolFlag); !isbool {
		placeholder := strings.ToLower(f.PlaceholderStr())
		str += fmt.Sprintf(` \fI%s\fP`, placeholder)
	}
	return str
}

func renderFlagUsage(f flag.Flag) string {
	usage := f.UsageStr()
	if f.IsRequired() {
		return usage
	}

	if f.EnvVarStr() != "" {
		usage += "\n.PP\nEnvironment variable: " + fmt.Sprintf(`\fC$%s\fP`, f.EnvVarStr())
	}
	if f.EnvVarStr() != "" && f.DefaultStr() != "" {
		usage += "\n.PD 0"
	}
	if f.DefaultStr() != "" {
		usage += "\n.PP\nDefault value: " + fmt.Sprintf(`\fC%s\fP`, f.DefaultStr())
	}
	if f.EnvVarStr() != "" && f.DefaultStr() != "" {
		usage += "\n.PD"
	}
	return usage
}
