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

package docker

import (
	"github.com/nmiyake/pkg/dirs"
	"github.com/palantir/pkg/cli"
	"github.com/palantir/pkg/cli/cfgcli"
	"github.com/palantir/pkg/cli/flag"

	"github.com/palantir/godel/apps/distgo/cmd"
	"github.com/palantir/godel/apps/distgo/config"
)

const (
	baseRepoFlagName = "base-repo"
	verboseFlagName  = "verbose"
)

var (
	baseRepoFlag = flag.StringFlag{
		Name:  baseRepoFlagName,
		Usage: "This is joined with per image repository path while building/publishing images",
		Value: "",
	}
	verboseFlag = flag.BoolFlag{
		Name:  verboseFlagName,
		Alias: "v",
		Usage: "Print the output from the Docker commands",
	}
)

func Command() cli.Command {
	build := cli.Command{
		Name: "build",
		Flags: []flag.Flag{
			verboseFlag,
			baseRepoFlag,
			cmd.ProductsParam,
		},
		Action: func(ctx cli.Context) error {
			cfg, err := config.Load(cfgcli.ConfigPath, cfgcli.ConfigJSON)
			if err != nil {
				return err
			}
			wd, err := dirs.GetwdEvalSymLinks()
			if err != nil {
				return err
			}
			return Build(
				ctx.Slice(cmd.ProductsParamName),
				cfg,
				wd,
				ctx.String(baseRepoFlagName),
				ctx.Bool(verboseFlagName),
				ctx.App.Stdout,
			)
		},
	}

	publish := cli.Command{
		Name: "publish",
		Flags: []flag.Flag{
			verboseFlag,
			baseRepoFlag,
			cmd.ProductsParam,
		},
		Action: func(ctx cli.Context) error {
			cfg, err := config.Load(cfgcli.ConfigPath, cfgcli.ConfigJSON)
			if err != nil {
				return err
			}
			wd, err := dirs.GetwdEvalSymLinks()
			if err != nil {
				return err
			}
			return Publish(
				ctx.Slice(cmd.ProductsParamName),
				cfg,
				wd,
				ctx.String(baseRepoFlagName),
				ctx.Bool(verboseFlagName),
				ctx.App.Stdout,
			)
		},
	}

	return cli.Command{
		Name:  "docker",
		Usage: "Runs docker tasks",
		Subcommands: []cli.Command{
			build,
			publish,
		},
	}
}
