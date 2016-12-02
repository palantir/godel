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

package packages

import (
	"fmt"
	"path"
	"strings"

	"github.com/nmiyake/pkg/dirs"
	"github.com/palantir/pkg/cli"
	"github.com/palantir/pkg/matcher"

	"github.com/palantir/godel/cmd"
	"github.com/palantir/godel/config"
)

func Command() cli.Command {
	return cli.Command{
		Name:  "packages",
		Usage: "Lists all of the packages in the project except those excluded by exclude.yml",
		Action: func(ctx cli.Context) error {
			excludeCfg := matcher.NamesPathsCfg{}
			if cfgDir, _ := cmd.ConfigDirPath(ctx); cfgDir != "" {
				var err error
				excludeCfg, err = config.GetExcludeCfgFromYML(path.Join(cfgDir, config.ExcludeYML))
				if err != nil {
					return err
				}
			}
			wd, err := dirs.GetwdEvalSymLinks()
			if err != nil {
				return err
			}
			pkgs, err := List(excludeCfg.Matcher(), wd)
			if err != nil {
				return err
			}
			fmt.Println(strings.Join(pkgs, "\n"))
			return nil
		},
	}
}
