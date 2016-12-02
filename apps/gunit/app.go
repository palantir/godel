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

package gunit

import (
	"github.com/nmiyake/pkg/errorstringer"
	"github.com/palantir/amalgomate/amalgomated"
	"github.com/palantir/pkg/cli"
	"github.com/palantir/pkg/cli/cfgcli"

	"github.com/palantir/godel/apps/gunit/cmd"
	"github.com/palantir/godel/apps/gunit/cmd/clean"
	"github.com/palantir/godel/apps/gunit/cmd/test"
)

func RunApp(args []string, supplier amalgomated.CmderSupplier) int {
	return amalgomated.RunApp(args, nil, cmd.Library, App(supplier).Run)
}

func App(supplier amalgomated.CmderSupplier) *cli.App {
	app := cli.NewApp(cfgcli.Handler(), cli.DebugHandler(errorstringer.StackWithInterleavedMessages))
	app.Name = "gunit"
	app.Usage = "Run test and coverage commands for Go code"
	app.Subcommands = []cli.Command{
		test.GoTestCommand(supplier),
		test.GoCoverCommand(supplier),
		test.GTCommand(supplier),
		clean.Command(),
	}
	app.Action = test.RunGoTestAction(supplier)
	app.Flags = append(app.Flags, cmd.GlobalFlags...)
	return app
}
