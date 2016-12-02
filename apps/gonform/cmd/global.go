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
	"github.com/palantir/amalgomate/amalgomated"
	"github.com/palantir/pkg/cli/flag"

	"github.com/palantir/godel/apps/gonform/generated_src"
)

const (
	listFlagName    = "list"
	filesParamName  = "files"
	verboseFlagName = "verbose"
)

var (
	Library     = amalgomated.NewCmdLibrary(amalgomatedformatters.Instance())
	VerboseFlag = flag.BoolFlag{
		Name:  verboseFlagName,
		Alias: "v",
		Usage: "Print formatters as they are run",
	}
	ListFlag = flag.BoolFlag{
		Name:  listFlagName,
		Alias: "l",
		Usage: "Print the files that would change if command is run",
	}
	FilesParam = flag.StringSlice{
		Name:     filesParamName,
		Usage:    "Files to format (defaults to all project .go files)",
		Optional: true,
	}
)
