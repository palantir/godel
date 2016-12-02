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

package golicense

import (
	"github.com/nmiyake/pkg/errorstringer"
	"github.com/palantir/pkg/cli"
	"github.com/palantir/pkg/cli/cfgcli"

	"github.com/palantir/checks/golicense/cmd"
)

func App() *cli.App {
	app := cli.NewApp(cli.DebugHandler(errorstringer.SingleStack), cfgcli.Handler())
	flags := app.Flags
	app.Command = cmd.Command()
	app.Flags = append(flags, app.Flags...)
	return app
}
