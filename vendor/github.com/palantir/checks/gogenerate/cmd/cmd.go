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

package cmd

import (
	"github.com/nmiyake/pkg/dirs"
	"github.com/palantir/pkg/cli"
	"github.com/palantir/pkg/cli/cfgcli"
	"github.com/palantir/pkg/cli/flag"

	"github.com/palantir/checks/gogenerate/config"
	"github.com/palantir/checks/gogenerate/gogenerate"
)

const (
	verifyFlagName = "verify"
)

var flags = []flag.Flag{
	flag.BoolFlag{
		Name:  verifyFlagName,
		Usage: "verify that running generators does not change the current output",
	},
}

func Command() cli.Command {
	return cli.Command{
		Name:  "generate",
		Usage: "Run generators specified in configuration",
		Flags: flags,
		Action: func(ctx cli.Context) error {
			wd, err := dirs.GetwdEvalSymLinks()
			if err != nil {
				return err
			}

			cfg, err := config.Load(cfgcli.ConfigPath, cfgcli.ConfigJSON)
			if err != nil {
				return err
			}

			return gogenerate.Run(wd, cfg, ctx.Bool(verifyFlagName), ctx.App.Stdout)
		},
	}
}
