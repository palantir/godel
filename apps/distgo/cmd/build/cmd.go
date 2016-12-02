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

package build

import (
	"github.com/nmiyake/pkg/dirs"
	"github.com/palantir/pkg/cli"
	"github.com/palantir/pkg/cli/cfgcli"
	"github.com/palantir/pkg/cli/flag"

	"github.com/palantir/godel/apps/distgo/cmd"
	"github.com/palantir/godel/apps/distgo/config"
)

const (
	parallelFlagName = "parallel"
	installFlagName  = "install"
	pkgDirFlagName   = "pkgdir"
)

var (
	parallelFlag = flag.BoolFlag{
		Name:  parallelFlagName,
		Usage: "Build binaries in parallel",
		Value: true,
	}
	installFlag = flag.BoolFlag{
		Name:  installFlagName,
		Usage: "Run 'install' before 'build'",
		Value: true,
	}
	pkgDirFlag = flag.BoolFlag{
		Name:  pkgDirFlagName,
		Usage: "Use a custom 'pkg' directory for 'install' action (only takes effect if 'install' is true)",
	}
)

func DefaultContext() Context {
	return Context{
		Parallel: parallelFlag.Value,
		Install:  installFlag.Value,
		Pkgdir:   pkgDirFlag.Value,
	}
}

func Command() cli.Command {
	return cli.Command{
		Name:  "build",
		Usage: "Build products",
		Flags: []flag.Flag{
			cmd.ProductsParam,
			parallelFlag,
			installFlag,
			pkgDirFlag,
			cmd.OSArchFlag,
		},
		Action: func(ctx cli.Context) error {
			buildCtx := Context{
				Parallel: ctx.Bool(parallelFlagName),
				Install:  ctx.Bool(installFlagName),
				Pkgdir:   ctx.Bool(pkgDirFlagName),
			}

			cfg, err := config.Load(cfgcli.ConfigPath, cfgcli.ConfigJSON)
			if err != nil {
				return err
			}

			wd, err := dirs.GetwdEvalSymLinks()
			if err != nil {
				return err
			}

			osArchs, err := cmd.NewOSArchFilter(ctx.String(cmd.OSArchFlagName))
			if err != nil {
				return err
			}
			return Products(ctx.Slice(cmd.ProductsParamName), osArchs, buildCtx, cfg, wd, ctx.App.Stdout)
		},
	}
}
