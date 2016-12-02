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

package gonform

import (
	"github.com/nmiyake/pkg/dirs"
	"github.com/nmiyake/pkg/errorstringer"
	"github.com/palantir/amalgomate/amalgomated"
	"github.com/palantir/pkg/cli"
	"github.com/palantir/pkg/cli/cfgcli"

	"github.com/palantir/godel/apps/gonform/cmd"
)

func RunApp(args []string, supplier amalgomated.CmderSupplier) int {
	return amalgomated.RunApp(args, nil, cmd.Library, App(supplier).Run)
}

func App(supplier amalgomated.CmderSupplier) *cli.App {
	app := cli.NewApp(cfgcli.Handler(), cli.DebugHandler(errorstringer.StackWithInterleavedMessages))
	app.Name = "gonform"
	app.Usage = "Format Go code"
	app.Subcommands = []cli.Command{
		cmd.RunAllCommand(supplier),
		cmd.GoFmtCommand(supplier),
		cmd.PTImportsCommand(supplier),
	}
	app.Flags = append(
		app.Flags,
		cmd.VerboseFlag,
		cmd.ListFlag,
	)
	app.Action = func(ctx cli.Context) error {
		wd, err := dirs.GetwdEvalSymLinks()
		if err != nil {
			return err
		}
		// no-arg invocation runs "runAll" command on all files
		return cmd.DoRunAll(nil, ctx, supplier, wd)
	}
	return app
}
