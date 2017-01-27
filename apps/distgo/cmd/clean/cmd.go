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
	"github.com/palantir/pkg/cli"
	"github.com/palantir/pkg/cli/cfgcli"
	"github.com/palantir/pkg/cli/flag"

	"github.com/palantir/godel/apps/distgo/cmd"
	"github.com/palantir/godel/apps/distgo/config"
)

const (
	forceFlagName = "force"
)

var (
	forceFlag = flag.BoolFlag{
		Name:  forceFlagName,
		Usage: "force removal of files in build and dist directories that weren't produced in a build",
	}
)

func Command() cli.Command {
	return cli.Command{
		Name:  "clean",
		Usage: "Clean product build and distribution directories",
		Flags: []flag.Flag{
			cmd.ProductsParam,
			forceFlag,
		},
		Action: func(ctx cli.Context) error {
			cfg, err := config.Load(cfgcli.ConfigPath, cfgcli.ConfigJSON)
			if err != nil {
				return err
			}

			return Clean(ctx.Slice(cmd.ProductsParamName), cfg, ctx.Bool(forceFlagName), ctx.App.Stdout)
		},
	}
}
