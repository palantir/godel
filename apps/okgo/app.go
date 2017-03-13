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

package okgo

import (
	"go/build"

	"github.com/nmiyake/pkg/dirs"
	"github.com/nmiyake/pkg/errorstringer"
	"github.com/palantir/amalgomate/amalgomated"
	"github.com/palantir/pkg/cli"
	"github.com/palantir/pkg/cli/cfgcli"

	"github.com/palantir/godel/apps/okgo/cmd"
	"github.com/palantir/godel/apps/okgo/cmd/cmdlib"
	"github.com/palantir/godel/apps/okgo/config"
)

func RunApp(args []string, supplier amalgomated.CmderSupplier) int {
	return amalgomated.RunApp(args, nil, cmdlib.Instance(), App(supplier).Run)
}

func App(supplier amalgomated.CmderSupplier) *cli.App {
	if releaseTag := cmd.GetReleaseTagEnvVar(); releaseTag != "" {
		// if ReleaseTag is specified and is found in the default build config, update default build config to
		// only include tags that occur before the provided release tag.
		var newRelTags []string
		for _, currTag := range build.Default.ReleaseTags {
			newRelTags = append(newRelTags, currTag)
			if releaseTag == currTag {
				// if this tag matches, do not include any of the newer tags
				break
			}
		}
		build.Default.ReleaseTags = newRelTags
	}

	app := cli.NewApp(cfgcli.Handler(), cli.DebugHandler(errorstringer.StackWithInterleavedMessages))
	app.Name = "okgo"
	app.Usage = "Run static checks for project"

	// add runAll and individual checks
	runAllCmd := cmd.RunAllCommand(supplier)
	app.Subcommands = []cli.Command{
		runAllCmd,
	}
	for _, currCmd := range cmdlib.Instance().Cmds() {
		app.Subcommands = append(app.Subcommands, cmd.SingleCheckCommand(currCmd, supplier))
	}
	// set default action to runAll
	app.Action = func(ctx cli.Context) error {
		cfg, err := config.Load(cfgcli.ConfigPath, cfgcli.ConfigJSON)
		if err != nil {
			return err
		}
		wd, err := dirs.GetwdEvalSymLinks()
		if err != nil {
			return err
		}
		if err := cmd.SetReleaseTagEnvVar(cfg.ReleaseTag); err != nil {
			return err
		}
		return cmd.DoRunAll(nil, cfg, supplier, wd, ctx.App.Stdout)
	}
	return app
}
