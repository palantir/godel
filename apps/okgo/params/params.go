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

package params

import (
	"github.com/palantir/amalgomate/amalgomated"
	"github.com/palantir/pkg/matcher"

	"github.com/palantir/godel/apps/okgo/checkoutput"
	"github.com/palantir/godel/apps/okgo/cmd/cmdlib"
)

type OKGo struct {
	// ReleaseTag specifies the latest supported Go release tag. Empty if default should be used.
	ReleaseTag string

	// Checks specifies the configuration used by the checks. The key is the name of the check and the value is the
	// custom configuration for that check.
	Checks map[amalgomated.Cmd]Checker

	// Exclude specifies the files that should be excluded from tests.
	Exclude matcher.Matcher
}

type Checker struct {
	// Skip specifies whether or not the check should be skipped entirely.
	Skip bool

	// Args specifies the command-line arguments provided to the check.
	Args []string

	// LineFilters specifies the filter definitions. Raw output lines that match the filter are excluded from
	// processing.
	LineFilters []checkoutput.Filterer
}

// ArgsForCheck returns the arguments for the requested check stored in the Config, or nil if no configuration for the
// specified check was present in the configuration. The second return value indicates whether or not configuration for
// the requested check was present.
func (p *OKGo) ArgsForCheck(check amalgomated.Cmd) ([]string, bool) {
	checkConfig, ok := p.Checks[check]
	if !ok {
		return nil, false
	}
	return checkConfig.Args, true
}

// FiltersForCheck returns the filters that should be used for the requested check. The returned slice is a
// concatenation of the global filters derived from the package excludes specified in the configuration followed by the
// filters specified for the provided check in the configuration. Returns an empty slice if no filters are present
// globally or for the specified check.The derivation from the global filters is done in case the packages can't be
// excluded before the check is run (can happen if the check only supports the "all" mode).
func (p *OKGo) FiltersForCheck(check amalgomated.Cmd) []checkoutput.Filterer {
	filters := append([]checkoutput.Filterer{}, checkoutput.MatcherFilter(p.Exclude))
	checkConfig, ok := p.Checks[check]
	if ok {
		filters = append(filters, checkConfig.LineFilters...)
	}
	return filters
}

func (p *OKGo) checkCommands() []amalgomated.Cmd {
	var cmds []amalgomated.Cmd
	for _, currCmd := range cmdlib.Instance().Cmds() {
		if _, ok := p.Checks[currCmd]; ok {
			cmds = append(cmds, currCmd)
		}
	}
	return cmds
}
