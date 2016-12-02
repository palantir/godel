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
	"sort"

	"github.com/palantir/pkg/cli"
	"github.com/palantir/pkg/cli/cfgcli"
	"github.com/palantir/pkg/cli/flag"
	"github.com/pkg/errors"

	"github.com/palantir/checks/gocd/config"
)

const (
	inputDirsParamName = "dirs"
	verifyFlagName     = "verify"
)

var flags = []flag.Flag{
	flag.BoolFlag{
		Name:  verifyFlagName,
		Usage: "verify that imports file exists and is up-to-date",
	},
	flag.StringSlice{
		Name:     inputDirsParamName,
		Usage:    "directories for which to perform operation",
		Optional: true,
	},
}

func Command() cli.Command {
	return cli.Command{
		Name:  "gocd",
		Usage: "Write or verify package import information for Go packages",
		Flags: flags,
		Action: func(ctx cli.Context) error {
			params, err := config.Load(cfgcli.ConfigPath, cfgcli.ConfigJSON)
			if err != nil {
				return err
			}

			var dirs []string
			if len(params.RootDirs) > 0 {
				cfgDirs := make(map[string]struct{})
				for _, dir := range params.RootDirs {
					cfgDirs[dir] = struct{}{}
				}

				// if dirs argument was provided, use it as a filter
				if ctx.Has(inputDirsParamName) {
					inputDirsSlice := ctx.Slice(inputDirsParamName)
					for _, dir := range inputDirsSlice {
						if _, ok := cfgDirs[dir]; ok {
							dirs = append(dirs, dir)
						}
					}
					if len(dirs) == 0 {
						return errors.Errorf("specified directories %v did not match any directories in configuration: %v", inputDirsSlice, sortedKeys(cfgDirs))
					}
				} else {
					dirs = params.RootDirs
					if len(dirs) == 0 {
						return errors.New("no input directories were specified and none were found in the configuration")
					}
				}
			} else if ctx.Has(inputDirsParamName) {
				dirs = ctx.Slice(inputDirsParamName)
				if len(dirs) == 0 {
					return errors.New("no input directories specified")
				}
			}

			if ctx.Bool(verifyFlagName) {
				return DoVerify(dirs)
			}

			return DoWriteImportsJSON(dirs)
		},
	}
}

func sortedKeys(in map[string]struct{}) []string {
	out := make([]string, 0, len(in))
	for k := range in {
		out = append(out, k)
	}
	sort.Strings(out)
	return out
}
