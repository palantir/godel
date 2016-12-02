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

package githubwiki

import (
	"github.com/palantir/pkg/cli"
	"github.com/palantir/pkg/cli/flag"
)

const (
	docsDirFlag        = "docs-dir"
	repoFlag           = "repository"
	authorNameFlag     = "author-name"
	authorEmailFlag    = "author-email"
	committerNameFlag  = "committer-name"
	committerEmailFlag = "committer-email"
	msgFlag            = "message"
)

func Command() cli.Command {
	return cli.Command{
		Name:  "github-wiki",
		Usage: "Push contents of a documents directory to a GitHub Wiki repository",
		Flags: []flag.Flag{
			flag.StringFlag{Name: docsDirFlag, Usage: "Directory whose contents should be pushed to the GitHub Wiki repository", Required: true},
			flag.StringFlag{Name: repoFlag, Usage: "GitHub wiki repository address (for example, git@github.com:org/project.wiki.git)", Required: true},
			flag.StringFlag{Name: authorNameFlag, Usage: "Author name to use for commit (if blank, uses value of last commit in current project)"},
			flag.StringFlag{Name: authorEmailFlag, Usage: "Author email to use for commit (if blank, uses value of last commit in current project)"},
			flag.StringFlag{Name: committerNameFlag, Usage: "Committer name to use for commit (if blank, uses value of last commit in current project)"},
			flag.StringFlag{Name: committerEmailFlag, Usage: "Committer email to use for commit (if blank, uses value of last commit in current project)"},
			flag.StringFlag{Name: msgFlag, Usage: "Commit message to use for commit in GitHub Wiki repository", Value: `Sync documentation using godel github-wiki task ({{ printf "%.7s" .CommitID}})`},
		},
		Action: func(ctx cli.Context) error {
			return SyncGitHubWiki(Params{
				DocsDir:        ctx.String(docsDirFlag),
				Repo:           ctx.String(repoFlag),
				AuthorName:     ctx.String(authorNameFlag),
				AuthorEmail:    ctx.String(authorEmailFlag),
				CommitterName:  ctx.String(committerNameFlag),
				CommitterEmail: ctx.String(committerEmailFlag),
				Msg:            ctx.String(msgFlag),
			}, ctx.App.Stdout)
		},
	}
}
