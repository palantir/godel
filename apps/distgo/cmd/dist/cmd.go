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

package dist

import (
	"github.com/nmiyake/pkg/dirs"
	"github.com/palantir/pkg/cli"
	"github.com/palantir/pkg/cli/cfgcli"
	"github.com/palantir/pkg/cli/flag"

	"github.com/palantir/godel/apps/distgo/cmd"
	"github.com/palantir/godel/apps/distgo/config"
)

const (
	forceBuildFlagName = "force-build"
)

var (
	forceBuildFlag = flag.BoolFlag{
		Name:  forceBuildFlagName,
		Usage: "Build all input build specs for distribution",
	}
)

func Command() cli.Command {
	return cli.Command{
		Name:  "dist",
		Usage: "Create a distribution for one or more products in the project",
		Flags: []flag.Flag{
			cmd.ProductsParam,
			forceBuildFlag,
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

			return Products(ctx.Slice(cmd.ProductsParamName), cfg, ctx.Bool(forceBuildFlagName), wd, ctx.App.Stdout)
		},
	}
}
