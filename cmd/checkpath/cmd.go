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

package checkpath

import (
	"github.com/nmiyake/pkg/dirs"
	"github.com/palantir/pkg/cli"
	"github.com/palantir/pkg/cli/flag"
)

const (
	cmd      = "check-path"
	infoFlag = "info"
)

func Command() cli.Command {
	return cli.Command{
		Name:  cmd,
		Usage: "Verify that the Go environment is set up properly and that the project is in the proper location",
		Flags: []flag.Flag{
			flag.BoolFlag{Name: infoFlag, Usage: "Provide information on current state and suggest fixes without executing them"},
		},
		Action: func(ctx cli.Context) error {
			wd, err := dirs.GetwdEvalSymLinks()
			if err != nil {
				return err
			}
			return VerifyProject(wd, ctx.Bool(infoFlag))
		},
	}
}
