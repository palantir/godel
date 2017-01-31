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

package run

import (
	"github.com/nmiyake/pkg/dirs"
	"github.com/palantir/pkg/cli"
	"github.com/palantir/pkg/cli/cfgcli"
	"github.com/palantir/pkg/cli/flag"
	"github.com/pkg/errors"

	"github.com/palantir/godel/apps/distgo/cmd/build"
	"github.com/palantir/godel/apps/distgo/config"
)

const (
	productFlagName = "product"
	argsFlagName    = "args"
)

func Command() cli.Command {
	return cli.Command{
		Name:  "run",
		Usage: "Run a product in the project",
		Flags: []flag.Flag{
			flag.StringFlag{
				Name:     productFlagName,
				Usage:    "Product to run",
				Required: false,
			},
			flag.StringSlice{
				Name:     argsFlagName,
				Usage:    "arguments to pass to product",
				Optional: true,
			},
		},
		Action: func(ctx cli.Context) error {
			cfg, err := config.Load(cfgcli.ConfigPath, cfgcli.ConfigJSON)
			if err != nil {
				return err
			}

			var products []string
			product := ctx.String(productFlagName)
			if product != "" {
				products = append(products, product)
			}

			wd, err := dirs.GetwdEvalSymLinks()
			if err != nil {
				return err
			}

			buildSpecsWithDeps, err := build.SpecsWithDepsForArgs(cfg, products, wd)
			if err != nil {
				return err
			}

			if len(buildSpecsWithDeps) > 1 {
				programNames := make([]string, len(buildSpecsWithDeps))
				for i, currSpec := range buildSpecsWithDeps {
					programNames[i] = currSpec.Spec.ProductName
				}
				return errors.Errorf("more than one product exists, so product to run must be specified using the '--%v' flag: %v", productFlagName, programNames)
			}

			return DoRun(buildSpecsWithDeps[0].Spec, ctx.Slice(argsFlagName), ctx.App.Stdout, ctx.App.Stderr)
		},
	}
}
