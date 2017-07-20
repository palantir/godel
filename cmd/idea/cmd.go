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

package idea

import (
	"github.com/nmiyake/pkg/dirs"
	"github.com/palantir/pkg/cli"
)

func Command() cli.Command {
	run := func(f func(wd string) error) error {
		wd, err := dirs.GetwdEvalSymLinks()
		if err != nil {
			return err
		}
		return f(wd)
	}

	const createIntelliJUsage = "Create IntelliJ project files for this project"
	createIntelliJFiles := func(ctx cli.Context) error {
		return run(func(wd string) error {
			return CreateIdeaFiles(wd, imlIntelliJTemplateContent, iprIntelliJTemplateContent)
		})
	}

	return cli.Command{
		Name:  "idea",
		Usage: createIntelliJUsage,
		Subcommands: []cli.Command{
			{
				Name:  "gogland",
				Usage: "Create Gogland project files for this project",
				Action: func(ctx cli.Context) error {
					return run(func(wd string) error {
						return CreateIdeaFiles(wd, imlGoglandTemplateContent, iprGoglandTemplateContent)
					})
				},
			},
			{
				Name:   "intellij",
				Usage:  createIntelliJUsage,
				Action: createIntelliJFiles,
			},
			{
				Name:  "clean",
				Usage: "Remove the IntelliJ IDEA project files for this project",
				Action: func(ctx cli.Context) error {
					return run(CleanIdeaFiles)
				},
			},
		},
		Action: createIntelliJFiles,
	}
}
