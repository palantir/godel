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
	"os"
	"path"

	"github.com/nmiyake/pkg/dirs"
	"github.com/palantir/pkg/cli"
	"github.com/palantir/pkg/cli/cfgcli"
	"github.com/palantir/pkg/cli/flag"
	"github.com/pkg/errors"

	"github.com/palantir/godel/apps/gunit/cmd"
	"github.com/palantir/godel/apps/gunit/config"
)

var pkgsParamName = "packages"

func Command() cli.Command {
	return cli.Command{
		Name:  "clean",
		Usage: "Remove any 'tmp_placeholder_test.go' files in the project",
		Flags: []flag.Flag{
			flag.StringSlice{
				Name:     pkgsParamName,
				Usage:    "Packages for which 'tmp_placeholder_test.go' files should be cleaned",
				Optional: true,
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
			pkgs, err := cmd.PkgPaths(ctx.Slice(pkgsParamName), wd, cfg.Exclude)
			if err != nil {
				return err
			}
			return run(pkgs, wd)
		},
	}
}

func run(pkgDirs []string, wd string) error {
	for _, currPkg := range pkgDirs {
		tmpPlaceholder := path.Join(wd, currPkg, "tmp_placeholder_test.go")
		if err := os.Remove(tmpPlaceholder); err != nil && !os.IsNotExist(err) {
			return errors.Wrapf(err, "failed to delete placeholder file %s", tmpPlaceholder)
		}
	}
	return nil
}
