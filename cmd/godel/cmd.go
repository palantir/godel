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
	"github.com/nmiyake/pkg/dirs"
	"github.com/palantir/pkg/cli"
	"github.com/palantir/pkg/cli/flag"

	"github.com/palantir/godel/cmd"
	"github.com/palantir/godel/layout"
)

const packageParam = "package"

var Version = "unspecified"

func VersionCommand() cli.Command {
	return cli.Command{
		Name:  "version",
		Usage: "Print the version",
		Action: func(ctx cli.Context) error {
			ctx.Println(layout.AppName, "version", Version)
			return nil
		},
	}
}

func InstallCommand() cli.Command {
	return cli.Command{
		Name:  "install",
		Usage: "Install gödel from a local tgz file",
		Flags: []flag.Flag{
			flag.StringParam{
				Name:  packageParam,
				Usage: "path to tgz of gödel distribution to install",
			},
		},
		Action: func(ctx cli.Context) error {
			wd, err := dirs.GetwdEvalSymLinks()
			if err != nil {
				return err
			}
			return NewInstall(wd, ctx.String(packageParam), ctx.App.Stdout)
		},
	}
}

func UpdateCommand() cli.Command {
	return cli.Command{
		Name:  "update",
		Usage: "Download and install the version of gödel specified in the godel.properties file",
		Action: func(ctx cli.Context) error {
			return Update(cmd.WrapperFlagValue(ctx), ctx.App.Stdout)
		},
	}
}
