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

package checks

import (
	"sort"
	"strings"

	"github.com/palantir/amalgomate/amalgomated"
	"github.com/palantir/pkg/pkgpath"
	"github.com/pkg/errors"

	"github.com/palantir/godel/apps/okgo/checkoutput"
	"github.com/palantir/godel/apps/okgo/cmd/cmdlib"
)

type byName []Checker

func (a byName) Len() int           { return len(a) }
func (a byName) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a byName) Less(i, j int) bool { return a[i].Cmd().Name() < a[j].Cmd().Name() }

type allArgType string

func GetChecker(cmd amalgomated.Cmd) (Checker, error) {
	for _, checker := range checkers() {
		if cmd == checker.Cmd() {
			return checker, nil
		}
	}
	return nil, errors.Errorf("no checker exists for %s", cmd.Name())
}

func checkers() []Checker {
	checkers := []Checker{
		&checkerDefinition{
			cmd:        cmdlib.Instance().MustNewCmd("compiles"),
			lineParser: checkoutput.MultiLineParser(pkgpath.Absolute),
			rawLineFilter: func(line string) bool {
				return strings.HasPrefix(line, "err: couldn't load packages due to errors:") || strings.HasPrefix(line, "-: ")
			},
		},
		&checkerDefinition{
			cmd:        cmdlib.Instance().MustNewCmd("deadcode"),
			lineParser: checkoutput.StartAfterFirstWhitespaceParser(pkgpath.Relative),
		},
		&checkerDefinition{
			cmd:        cmdlib.Instance().MustNewCmd("errcheck"),
			lineParser: checkoutput.DefaultParser(pkgpath.Relative),
		},
		&checkerDefinition{
			cmd:        cmdlib.Instance().MustNewCmd("extimport"),
			lineParser: checkoutput.DefaultParser(pkgpath.Absolute),
		},
		&checkerDefinition{
			cmd:        cmdlib.Instance().MustNewCmd("golint"),
			lineParser: checkoutput.DefaultParser(pkgpath.Relative),
			allArg:     splat,
		},
		&checkerDefinition{
			cmd:        cmdlib.Instance().MustNewCmd("govet"),
			lineParser: checkoutput.DefaultParser(pkgpath.Relative),
			rawLineFilter: func(line string) bool {
				return line == "exit status 1"
			},
		},
		&checkerDefinition{
			cmd:        cmdlib.Instance().MustNewCmd("importalias"),
			lineParser: checkoutput.DefaultParser(pkgpath.Relative),
		},
		&checkerDefinition{
			cmd:         cmdlib.Instance().MustNewCmd("ineffassign"),
			lineParser:  checkoutput.DefaultParser(pkgpath.Absolute),
			allArg:      rootDir,
			globalCheck: true,
		},
		&checkerDefinition{
			cmd:        cmdlib.Instance().MustNewCmd("nobadfuncs"),
			lineParser: checkoutput.DefaultParser(pkgpath.Absolute),
		},
		&checkerDefinition{
			cmd:        cmdlib.Instance().MustNewCmd("novendor"),
			lineParser: checkoutput.RawParser(),
		},
		&checkerDefinition{
			cmd:        cmdlib.Instance().MustNewCmd("outparamcheck"),
			lineParser: checkoutput.DefaultParser(pkgpath.Relative),
			rawLineFilter: func(line string) bool {
				return strings.HasSuffix(line, `the parameters listed above require the use of '&', for example f(&x) instead of f(x)`)
			},
		},
		&checkerDefinition{
			cmd:        cmdlib.Instance().MustNewCmd("unconvert"),
			lineParser: checkoutput.DefaultParser(pkgpath.Absolute),
		},
		&checkerDefinition{
			cmd:        cmdlib.Instance().MustNewCmd("varcheck"),
			lineParser: checkoutput.StartAfterFirstWhitespaceParser(pkgpath.Absolute),
		},
	}
	sort.Sort(byName(checkers))
	return checkers
}
