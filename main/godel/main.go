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

package main

import (
	"fmt"
	"os"

	"github.com/kardianos/osext"
	"github.com/nmiyake/pkg/dirs"
	"github.com/palantir/amalgomate/amalgomated"
	"github.com/pkg/errors"

	"github.com/palantir/godel"
	"github.com/palantir/godel/cmd"
	"github.com/palantir/godel/cmd/clicmds"
)

func main() {
	gödelPath, err := osext.Executable()
	if err != nil {
		fmt.Printf("%+v\n", errors.Wrapf(err, "failed to determine path for current executable"))
		os.Exit(1)
	}

	if err := dirs.SetGoEnvVariables(); err != nil {
		fmt.Printf("%+v\n", errors.Wrapf(err, "failed to set Go environment variables"))
		os.Exit(1)
	}

	cmdLibrary, err := clicmds.CfgCliCmdSet(gödelPath)
	if err != nil {
		fmt.Printf("%+v\n", errors.Wrapf(err, "failed to create CfgCliCmdSet"))
		os.Exit(1)
	}
	os.Exit(amalgomated.RunApp(os.Args, cmd.GlobalFlagSet(), cmdLibrary, godel.App(gödelPath).Run))
}
