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

package amalgomated

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/nmiyake/pkg/dirs"
	"github.com/nmiyake/pkg/errorstringer"
	"github.com/palantir/pkg/cli"
	"github.com/palantir/pkg/cli/flag"
	"github.com/palantir/pkg/pkgpath"
	"github.com/pkg/errors"

	"github.com/palantir/godel/apps/okgo/generated_src/internal/github.com/palantir/checks/nobadfuncs/nobadfuncs"
)

const (
	printAllFlagName	= "all"
	jsonConfigFlagName	= "config"
	pkgsFlagName		= "pkgs"
)

var (
	printAllFlag	= flag.BoolFlag{
		Name:	printAllFlagName,
		Usage:	"print all function references",
	}
	jsonFlag	= flag.StringFlag{
		Name:	jsonConfigFlagName,
		Usage: "JSON configuration specifying blacklisted functions. Must be a JSON map from string to string, " +
			"where the key is a function signature and the value is the failure message printed when a function" +
			"with that signature is found.",
	}
	pkgsFlag	= flag.StringSlice{
		Name:	pkgsFlagName,
		Usage:	"paths to the packages to check",
	}
)

func AmalgomatedMain() {
	app := cli.NewApp(cli.DebugHandler(errorstringer.SingleStack))
	app.Flags = append(
		app.Flags,
		printAllFlag,
		jsonFlag,
		pkgsFlag,
	)
	app.Action = func(ctx cli.Context) error {
		pkgPaths, err := getPkgPaths(ctx.Slice(pkgsFlagName))
		if err != nil {
			return errors.Wrapf(err, "failed to determine package paths")
		}

		if ctx.Bool(printAllFlagName) {
			if err := nobadfuncs.PrintAllFuncRefs(pkgPaths, ctx.App.Stdout); err != nil {
				return errors.Wrapf(err, "Failed to determine all function references")
			}
			return nil
		}

		var jsonConfig map[string]string
		if ctx.Has(jsonConfigFlagName) {
			if err := json.Unmarshal([]byte(ctx.String(jsonConfigFlagName)), &jsonConfig); err != nil {
				return errors.Wrapf(err, "failed to read configuration")
			}
		}
		ok, err := nobadfuncs.PrintBadFuncRefs(pkgPaths, jsonConfig, ctx.App.Stdout)
		if err != nil {
			return errors.Wrapf(err, "nobadfuncs failed")
		}
		if !ok {
			// if there was no error but bad references were found, return empty error
			return fmt.Errorf("")
		}
		return nil
	}
	os.Exit(app.Run(os.Args))
}

func getPkgPaths(relPaths []string) ([]string, error) {
	wd, err := dirs.GetwdEvalSymLinks()
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get working directory")
	}

	var goSrcPaths []string
	for _, currPkg := range relPaths {
		newPath, err := pkgpath.NewRelPkgPath(currPkg, wd).GoPathSrcRel()
		if err != nil {
			return nil, errors.Wrapf(err, "failed to determine package path")
		}
		goSrcPaths = append(goSrcPaths, newPath)
	}
	return goSrcPaths, nil
}
