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
	goflag "flag"
	"io/ioutil"
	"path"

	"github.com/palantir/pkg/cli"
	"github.com/palantir/pkg/cli/flag"
	"github.com/pkg/errors"

	"github.com/palantir/godel/config"
)

const (
	configFlag     = "config"
	wrapperFlag    = "wrapper"
	verifyPathFlag = "verify-path"
)

var flagCreators = []flagCreator{
	&stringFlagCreator{
		name:  configFlag,
		usage: "path to the gödel config directory",
	},
	&stringFlagCreator{
		name:  wrapperFlag,
		usage: "path to the wrapper script for this invocation",
	},
	&boolFlagCreator{
		name:  verifyPathFlag,
		usage: "verify project path on startup",
	},
}

func GlobalCLIFlags() []flag.Flag {
	flags := make([]flag.Flag, len(flagCreators))
	for i, creator := range flagCreators {
		flags[i] = creator.cliFlag()
	}
	return flags
}

func GlobalFlagSet() *goflag.FlagSet {
	fset := goflag.NewFlagSet("gödel", goflag.ContinueOnError)
	for _, creator := range flagCreators {
		creator.goFlag(fset)
	}
	fset.SetOutput(ioutil.Discard)
	return fset
}

func ConfigDirPath(ctx cli.Context) (string, error) {
	return config.GetCfgDirPath(ctx.String(configFlag), ctx.String(wrapperFlag))
}

func WrapperFlagValue(ctx cli.Context) string {
	return ctx.String(wrapperFlag)
}

func Config(ctx cli.Context, cmdCfgPath, wd string) (cmdCfgFilePath string, excludeJSON string, err error) {
	gödelCfgDirPath, err := ConfigDirPath(ctx)
	if err != nil {
		return "", "", errors.Wrapf(err, "failed to determine configuration directory path")
	}

	if !path.IsAbs(gödelCfgDirPath) {
		gödelCfgDirPath = path.Join(wd, gödelCfgDirPath)
	}
	cmdCfgFilePath = path.Join(gödelCfgDirPath, cmdCfgPath)

	// get exclude JSON from exclude config file if it exists
	excludeCfgPath := path.Join(gödelCfgDirPath, config.ExcludeYML)
	excludeJSONBytes, err := config.ReadExcludeJSONFromYML(excludeCfgPath)
	if err != nil {
		return "", "", errors.Wrapf(err, "failed to read exclude config from %s", excludeCfgPath)
	}
	excludeJSON = string(excludeJSONBytes)

	return cmdCfgFilePath, excludeJSON, nil
}

type flagCreator interface {
	cliFlag() flag.Flag
	goFlag(fset *goflag.FlagSet)
}

type stringFlagCreator struct {
	name  string
	usage string
	value string
}

func (c *stringFlagCreator) cliFlag() flag.Flag {
	return flag.StringFlag{
		Name:  c.name,
		Usage: c.usage,
		Value: c.value,
	}
}

func (c *stringFlagCreator) goFlag(fset *goflag.FlagSet) {
	_ = fset.String(c.name, c.value, c.usage)
}

type boolFlagCreator struct {
	name  string
	usage string
	value bool
}

func (c *boolFlagCreator) cliFlag() flag.Flag {
	return flag.BoolFlag{
		Name:  c.name,
		Usage: c.usage,
		Value: c.value,
	}
}

func (c *boolFlagCreator) goFlag(fset *goflag.FlagSet) {
	_ = fset.Bool(c.name, c.value, c.usage)
}
