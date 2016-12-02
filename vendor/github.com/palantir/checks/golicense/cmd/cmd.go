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
	"fmt"
	"strings"

	"github.com/nmiyake/pkg/dirs"
	"github.com/palantir/pkg/cli"
	"github.com/palantir/pkg/cli/cfgcli"
	"github.com/palantir/pkg/cli/flag"
	"github.com/palantir/pkg/matcher"
	"github.com/pkg/errors"

	"github.com/palantir/checks/golicense"
)

const (
	filesFlagName  = "files"
	verifyFlagName = "verify"
	removeFlagName = "remove"
)

var flags = []flag.Flag{
	flag.BoolFlag{
		Name:  verifyFlagName,
		Usage: "verify that files have proper license headers applied",
	},
	flag.BoolFlag{
		Name:  removeFlagName,
		Usage: "remove the license header from files (no-op if verify is true)",
	},
	flag.StringSlice{
		Name:     filesFlagName,
		Usage:    "files on which to perform operation (if they are not excluded by configuration)",
		Optional: true,
	},
}

func Command() cli.Command {
	return cli.Command{
		Name:  "license",
		Usage: "Write or verify license headers for Go files",
		Flags: flags,
		Action: func(ctx cli.Context) error {
			wd, err := dirs.GetwdEvalSymLinks()
			if err != nil {
				return err
			}

			cfg, err := golicense.Load(cfgcli.ConfigPath, cfgcli.ConfigJSON)
			if err != nil {
				return err
			}

			// if header and matchers do not exist, return (nothing to check)
			if cfg.Header == "" && len(cfg.CustomHeaders) == 0 {
				return nil
			}

			var files []string
			if ctx.Has(filesFlagName) {
				files = ctx.Slice(filesFlagName)
			} else {
				files, err = matcher.ListFiles(wd, matcher.Name(`.+`), nil)
				if err != nil {
					return err
				}
			}

			verify := false
			if ctx.Has(verifyFlagName) {
				verify = ctx.Bool(verifyFlagName)
			}

			switch {
			case verify:
				// run verify
				modified, err := golicense.LicenseFiles(files, cfg, !verify)
				if err != nil {
					return err
				}
				if len(modified) > 0 {
					var plural string
					if len(modified) == 1 {
						plural = "file does"
					} else {
						plural = "files do"
					}

					parts := append([]string{fmt.Sprintf("%d %s not have the correct license header:", len(modified), plural)}, modified...)
					return errors.New(strings.Join(parts, "\n\t"))
				}
			case ctx.Has(removeFlagName) && ctx.Bool(removeFlagName):
				// run unlicense
				if _, err := golicense.UnlicenseFiles(files, cfg, true); err != nil {
					return err
				}
			default:
				// run license
				if _, err := golicense.LicenseFiles(files, cfg, !verify); err != nil {
					return err
				}
			}
			return nil
		},
	}
}
