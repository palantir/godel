// Copyright 2016 Palantir Technologies, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package clicmds

import (
	"fmt"
	"os"
	"path"

	"github.com/nmiyake/pkg/dirs"
	"github.com/palantir/amalgomate/amalgomated"
	"github.com/palantir/checks/gocd/cmd/gocd"
	"github.com/palantir/checks/golicense/cmd/golicense"
	"github.com/palantir/pkg/cli"
	"github.com/palantir/pkg/cli/cfgcli"
	"github.com/pkg/errors"

	"github.com/palantir/godel/apps/distgo"
	"github.com/palantir/godel/apps/gonform"
	"github.com/palantir/godel/apps/gunit"
	"github.com/palantir/godel/apps/okgo"
	"github.com/palantir/godel/cmd"
)

func CfgCliCmdSet(gödelPath string) (amalgomated.CmdLibrary, error) {
	cmds := make([]*amalgomated.CmdWithRunner, len(cfgCliCmds))
	for i := range cfgCliCmds {
		supplier := gödelRunnerSupplier(gödelPath, cfgCliCmds[i].name)
		cmd, err := cfgCliCmds[i].namedCmd(supplier)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to create command")
		}
		cmds[i] = cmd
	}

	cmdSet, err := amalgomated.NewStringCmdSetForRunners(cmds...)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to create StringCmdSet for runners")
	}

	return amalgomated.NewCmdLibrary(cmdSet), nil
}

func CfgCliCommands(gödelPath string) []cli.Command {
	var cmds []cli.Command
	for _, cmd := range cfgCliCmds {
		supplier := gödelRunnerSupplier(gödelPath, cmd.name)
		cmds = append(cmds, cmd.cliCmd(supplier))
	}
	return cmds
}

func gödelRunnerSupplier(gödelPath, cmdName string) amalgomated.CmderSupplier {
	return func(cmd amalgomated.Cmd) (amalgomated.Cmder, error) {
		// first underscore indicates to gödel that it is running in impersonation mode, while second underscore
		// signals this to the command itself (handled by "processHiddenCommand" in
		// amalgomated.cmdSetApp.RunApp)
		return amalgomated.PathCmder(gödelPath, amalgomated.ProxyCmdPrefix+cmdName, amalgomated.ProxyCmdPrefix+cmd.Name()), nil
	}
}

var (
	distgoCreator = func(supplier amalgomated.CmderSupplier) *cli.App {
		return distgo.App()
	}
	distgoDecorator = func(f func(ctx cli.Context) error) func(ctx cli.Context) error {
		return func(ctx cli.Context) error {
			cfgDirPath, err := cmd.ConfigDirPath(ctx)
			if err != nil {
				return err
			}
			// use the project directory for the configuration as the working directory for distgo commands
			projectPath := path.Join(cfgDirPath, "..", "..")
			if err := os.Chdir(projectPath); err != nil {
				return err
			}
			return f(ctx)
		}
	}
	cfgCliCmds = []*goRunnerCmd{
		{
			name:   "format",
			app:    gonform.App,
			runApp: gonform.RunApp,
		},
		{
			name:   "check",
			app:    okgo.App,
			runApp: okgo.RunApp,
		},
		{
			name:   "test",
			app:    gunit.App,
			runApp: gunit.RunApp,
		},
		{
			name: "imports",
			app: func(supplier amalgomated.CmderSupplier) *cli.App {
				return gocd.App()
			},
			pathToCfg: []string{"imports.yml"},
		},
		{
			name: "license",
			app: func(supplier amalgomated.CmderSupplier) *cli.App {
				return golicense.App()
			},
			pathToCfg: []string{"license.yml"},
		},
		{
			name:           "run",
			app:            distgoCreator,
			decorator:      distgoDecorator,
			subcommandPath: []string{"run"},
			pathToCfg:      []string{"dist.yml"},
		},
		{
			name:           "build",
			app:            distgoCreator,
			decorator:      distgoDecorator,
			subcommandPath: []string{"build"},
			pathToCfg:      []string{"dist.yml"},
		},
		{
			name:           "dist",
			app:            distgoCreator,
			decorator:      distgoDecorator,
			subcommandPath: []string{"dist"},
			pathToCfg:      []string{"dist.yml"},
		},
		{
			name:           "artifacts",
			app:            distgoCreator,
			decorator:      distgoDecorator,
			subcommandPath: []string{"artifacts"},
			pathToCfg:      []string{"dist.yml"},
		},
		{
			name:           "products",
			app:            distgoCreator,
			decorator:      distgoDecorator,
			subcommandPath: []string{"products"},
			pathToCfg:      []string{"dist.yml"},
		},
		{
			name:           "publish",
			app:            distgoCreator,
			decorator:      distgoDecorator,
			subcommandPath: []string{"publish"},
			pathToCfg:      []string{"dist.yml"},
		},
	}
)

type appSupplier func(supplier amalgomated.CmderSupplier) *cli.App

// locator that returns the subcommand
func mustAppSubcommand(app *cli.App, subcommands ...string) cli.Command {
	currCmd := app.Command
	for _, currDesiredSubcmd := range subcommands {
		found := false
		for _, currActualSubcmd := range currCmd.Subcommands {
			if currActualSubcmd.Name == currDesiredSubcmd {
				currCmd = currActualSubcmd
				found = true
				break
			}
		}
		if !found {
			panic(fmt.Sprintf("subcommand %v not present in command %v (full path: %v, top-level command: %v)", currDesiredSubcmd, currCmd.Subcommands, subcommands, app))
		}
	}
	return currCmd
}

type goRunnerCmd struct {
	name           string
	app            appSupplier
	runApp         func(args []string, supplier amalgomated.CmderSupplier) int
	subcommandPath []string
	pathToCfg      []string
	// optional decorator for actions
	decorator func(func(ctx cli.Context) error) func(ctx cli.Context) error
}

func (c *goRunnerCmd) cliCmd(supplier amalgomated.CmderSupplier) cli.Command {
	// create the cliCmd app for this command
	app := c.app(supplier)

	// get the app command
	cmdToRun := mustAppSubcommand(app, c.subcommandPath...)

	// update name of command to be specified name
	cmdToRun.Name = c.name

	// decorate all actions so that it sets cfgcli config variables before command is run
	c.decorateAllActions(&cmdToRun)

	return cmdToRun
}

func (c *goRunnerCmd) decorateAllActions(cmd *cli.Command) {
	if cmd != nil {
		if cmd.Action != nil {
			cmd.Action = c.decorateAction(cmd.Action)
		}
		for i := range cmd.Subcommands {
			c.decorateAllActions(&cmd.Subcommands[i])
		}
	}
}

func (c *goRunnerCmd) namedCmd(supplier amalgomated.CmderSupplier) (*amalgomated.CmdWithRunner, error) {
	return amalgomated.NewCmdWithRunner(c.name, func() {
		c.runApp(os.Args, supplier)
	})
}

func (c *goRunnerCmd) decorateAction(f func(ctx cli.Context) error) func(ctx cli.Context) error {
	return func(ctx cli.Context) error {
		wd, err := dirs.GetwdEvalSymLinks()
		if err != nil {
			return err
		}
		pathToCfg := c.name + ".yml"
		if len(c.pathToCfg) != 0 {
			pathToCfg = path.Join(c.pathToCfg...)
		}
		cfgFilePath, cfgJSON, err := cmd.Config(ctx, pathToCfg, wd)
		if err != nil {
			return err
		}
		cfgcli.ConfigPath = cfgFilePath
		cfgcli.ConfigJSON = cfgJSON

		if c.decorator != nil {
			// if goRunnerCmd declares a decorator, decorate action with it before invoking
			f = c.decorator(f)
		}
		return f(ctx)
	}
}
