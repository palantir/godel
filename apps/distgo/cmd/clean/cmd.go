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

package clean

import (
	"github.com/nmiyake/pkg/dirs"
	"github.com/palantir/pkg/cli"
	"github.com/palantir/pkg/cli/cfgcli"
	"github.com/palantir/pkg/cli/flag"

	"github.com/palantir/godel/apps/distgo/cmd"
	"github.com/palantir/godel/apps/distgo/config"
)

const (
	dryRunFlagName = "dry-run"
)

func Command() cli.Command {
	return cli.Command{
		Name:  "clean",
		Usage: "Remove the build and dist outputs for products",
		Flags: []flag.Flag{
			cmd.ProductsParam,
			flag.BoolFlag{
				Name:  dryRunFlagName,
				Usage: "Print the paths that would be removed by the operation without actually removing them",
			},
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
			return Products(ctx.Slice(cmd.ProductsParamName), cfg, ctx.Bool(dryRunFlagName), wd, ctx.App.Stdout)
		},
	}
}
