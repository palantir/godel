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

package builtintasks

import (
	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/palantir/godel/v2/framework/builtintasks/githubwiki"
	"github.com/palantir/godel/v2/framework/godellauncher"
)

func GitHubWikiTask() godellauncher.Task {
	var githubWikiParams githubwiki.Params

	const (
		docsDirFlag        = "docs-dir"
		repoFlag           = "repository"
		authorNameFlag     = "author-name"
		authorEmailFlag    = "author-email"
		committerNameFlag  = "committer-name"
		committerEmailFlag = "committer-email"
		msgFlag            = "message"
	)

	cmd := &cobra.Command{
		Use:   "github-wiki",
		Short: "Push contents of a documents directory to a GitHub Wiki repository",
		RunE: func(cmd *cobra.Command, args []string) error {
			return githubwiki.SyncGitHubWiki(githubWikiParams, cmd.OutOrStdout())
		},
	}

	cmd.Flags().StringVar(&githubWikiParams.DocsDir, docsDirFlag, "", "Directory whose contents should be pushed to the GitHub Wiki repository")
	cmd.Flags().StringVar(&githubWikiParams.Repo, repoFlag, "", "GitHub wiki repository address (for example, git@github.com:org/project.wiki.git)")
	cmd.Flags().StringVar(&githubWikiParams.AuthorName, authorNameFlag, "", "Author name to use for commit (if blank, uses value of last commit in current project)")
	cmd.Flags().StringVar(&githubWikiParams.AuthorEmail, authorEmailFlag, "", "Author email to use for commit (if blank, uses value of last commit in current project)")
	cmd.Flags().StringVar(&githubWikiParams.CommitterName, committerNameFlag, "", "Committer name to use for commit (if blank, uses value of last commit in current project)")
	cmd.Flags().StringVar(&githubWikiParams.CommitterEmail, committerEmailFlag, "", "Committer email to use for commit (if blank, uses value of last commit in current project)")
	cmd.Flags().StringVar(&githubWikiParams.Msg, msgFlag, `Sync documentation using godel github-wiki task ({{ printf "%.7s" .CommitID}})`, "Commit message to use for commit in GitHub Wiki repository")

	if err := cmd.MarkFlagRequired(docsDirFlag); err != nil {
		panic(errors.Wrapf(err, "failed to mark flag %s as required", docsDirFlag))
	}
	if err := cmd.MarkFlagRequired(repoFlag); err != nil {
		panic(errors.Wrapf(err, "failed to mark flag %s as required", repoFlag))
	}

	return godellauncher.CobraCLITask(cmd, nil)
}
