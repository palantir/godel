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

package cli

import (
	"github.com/palantir/pkg/cli/flag"
)

type Command struct {
	Name         string
	Alias        string
	Usage        string
	Description  string // prose
	Flags        []flag.Flag
	Subcommands  []Command
	DecisionFlag string
	Action       func(ctx Context) error
}

const DefaultDecision string = ""

func (cmd Command) Names() []string {
	names := []string{}
	if cmd.Name != DefaultDecision {
		names = append(names, cmd.Name)
	}
	if cmd.Alias != "" {
		names = append(names, cmd.Alias)
	}
	return names
}
