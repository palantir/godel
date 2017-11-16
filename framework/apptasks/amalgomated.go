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

package apptasks

import (
	"fmt"
	"io"
	"os"

	"github.com/palantir/amalgomate/amalgomated"
	"github.com/palantir/pkg/cli"
	"github.com/pkg/errors"

	"github.com/palantir/godel/apps/gonform"
	gonformcmd "github.com/palantir/godel/apps/gonform/cmd"
	"github.com/palantir/godel/apps/gunit"
	gunitcmd "github.com/palantir/godel/apps/gunit/cmd"
	"github.com/palantir/godel/apps/okgo"
	okgocmdlib "github.com/palantir/godel/apps/okgo/cmd/cmdlib"
	"github.com/palantir/godel/framework/godel"
	"github.com/palantir/godel/framework/godellauncher"
	"github.com/palantir/godel/framework/verifyorder"
)

var amalgomatedCmds = []amalgCmd{
	{
		Name:       "format",
		ConfigFile: "format.yml",
		CmdLib:     gonformcmd.Library,
		AppCreator: gonform.App,
		Verify: &godellauncher.VerifyOptions{
			ApplyFalseArgs: []string{"-l"},
			Ordering:       verifyorder.Format,
		},
	},
	{
		Name:       "check",
		ConfigFile: "check.yml",
		CmdLib:     okgocmdlib.Instance(),
		AppCreator: okgo.App,
		Verify: &godellauncher.VerifyOptions{
			Ordering: verifyorder.Check,
		},
	},
	{
		Name:       "test",
		ConfigFile: "test.yml",
		CmdLib:     gunitcmd.Library,
		AppCreator: gunit.App,
		Verify: &godellauncher.VerifyOptions{
			Ordering: verifyorder.Test,
			VerifyTaskFlags: []godellauncher.VerifyFlag{
				{
					Name:        "junit-output",
					Description: "Path to JUnit XML output (only used if 'test' task is run)",
					Type:        godellauncher.StringFlag,
				},
				{
					Name:        "tags",
					Description: "Specify tags that should be used for tests (only used if 'test' task is run)",
					Type:        godellauncher.StringFlag,
				},
			},
		},
	},
}

func AmalgomatedCmdLib(gödelPath string) (amalgomated.CmdLibrary, error) {
	var cmds []*amalgomated.CmdWithRunner
	for _, amalgCmd := range amalgomatedCmds {
		namedCmd, err := amalgCmd.namedCmd(gödelPath)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to create command")
		}
		cmds = append(cmds, namedCmd)
	}

	cmdSet, err := amalgomated.NewStringCmdSetForRunners(cmds...)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to create StringCmdSet for runners")
	}
	return amalgomated.NewCmdLibrary(cmdSet), nil
}

func AmalgomatedTasks() []godellauncher.Task {
	var tasks []godellauncher.Task
	for _, cmd := range amalgomatedCmds {
		tasks = append(tasks, cmd.toTask())
	}
	return tasks
}

func gödelRunnerSupplier(gödelPath, cmdName string) amalgomated.CmderSupplier {
	return func(cmd amalgomated.Cmd) (amalgomated.Cmder, error) {
		// first underscore indicates to gödel that it is running in impersonation mode, while second underscore
		// signals this to the command itself (handled by "processHiddenCommand" in
		// amalgomated.cmdSetApp.RunApp)
		return amalgomated.PathCmder(gödelPath, amalgomated.ProxyCmdPrefix+cmdName, amalgomated.ProxyCmdPrefix+cmd.Name()), nil
	}
}

type amalgCmd struct {
	Name       string
	ConfigFile string
	Verify     *godellauncher.VerifyOptions
	CmdLib     amalgomated.CmdLibrary
	AppCreator func(supplier amalgomated.CmderSupplier) *cli.App
}

func (c amalgCmd) toTask() godellauncher.Task {
	return godellauncher.Task{
		Name:        c.Name,
		Description: c.AppCreator(nil).Usage,
		ConfigFile:  c.ConfigFile,
		Verify:      c.Verify,
		RunImpl: func(t *godellauncher.Task, global godellauncher.GlobalConfig, stdout io.Writer) error {
			gödelPath := global.Executable
			args, err := cfgCLIArgs(global, nil, t.ConfigFile)
			if err != nil {
				return err
			}
			os.Args = args

			if exitCode := c.runApp(args, gödelRunnerSupplier(gödelPath, c.Name)); exitCode != 0 {
				return fmt.Errorf("")
			}
			return nil
		},
	}
}

func (c amalgCmd) namedCmd(gödelPath string) (*amalgomated.CmdWithRunner, error) {
	return amalgomated.NewCmdWithRunner(c.Name, func() {
		c.runApp(os.Args, gödelRunnerSupplier(gödelPath, c.Name))
	})
}

func (c amalgCmd) runApp(args []string, supplier amalgomated.CmderSupplier) int {
	app := c.AppCreator(supplier)
	app.Name = godel.AppName
	return amalgomated.RunApp(args, nil, c.CmdLib, app.Run)
}
