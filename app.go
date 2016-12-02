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

package godel

import (
	"github.com/nmiyake/pkg/errorstringer"
	"github.com/palantir/pkg/cli"

	"github.com/palantir/godel/cmd"
	"github.com/palantir/godel/cmd/checkpath"
	"github.com/palantir/godel/cmd/clicmds"
	"github.com/palantir/godel/cmd/githooks"
	"github.com/palantir/godel/cmd/githubwiki"
	"github.com/palantir/godel/cmd/godel"
	"github.com/palantir/godel/cmd/idea"
	"github.com/palantir/godel/cmd/packages"
	"github.com/palantir/godel/cmd/verify"
)

func App(gödelPath string) *cli.App {
	app := cli.NewApp(cli.DebugHandler(errorstringer.StackWithInterleavedMessages))
	app.Name = "godel"
	app.Usage = "Run tasks for coding, checking, formatting, testing, building and publishing Go code"
	app.Flags = append(app.Flags, cmd.GlobalCLIFlags()...)
	app.Version = godel.Version

	app.Subcommands = []cli.Command{
		godel.VersionCommand(),
		godel.InstallCommand(),
		godel.UpdateCommand(),
		checkpath.Command(),
		githooks.Command(),
		githubwiki.Command(),
		idea.Command(),
		packages.Command(),
		verify.Command(gödelPath),
	}
	app.Subcommands = append(app.Subcommands, clicmds.CfgCliCommands(gödelPath)...)

	return app
}
