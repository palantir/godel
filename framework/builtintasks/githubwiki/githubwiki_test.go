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
	"bytes"
	"fmt"
	"io/ioutil"
	"os/exec"
	"path"
	"regexp"
	"strings"
	"testing"

	"github.com/nmiyake/pkg/dirs"
	"github.com/palantir/pkg/gittest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSyncGitHubWiki(t *testing.T) {
	tmpDir, cleanup, err := dirs.TempDir("", "")
	defer cleanup()
	require.NoError(t, err)

	srcRepoParams := commitUserParam{
		authorName:     "src-author",
		authorEmail:    "src-author@email.com",
		committerName:  "src-committer",
		committerEmail: "src-committer@email.com",
	}

	for i, currCase := range []struct {
		params      commitUserParam
		msg         string
		want        commitUserParam
		wantMessage func(docsDir string) string
	}{
		// provided parameters are used
		{
			params: commitUserParam{
				authorName:     "Author Name",
				authorEmail:    "author@email.com",
				committerName:  "Committer Name",
				committerEmail: "committer@email.com",
			},
			want: commitUserParam{
				authorName:     "Author Name",
				authorEmail:    "author@email.com",
				committerName:  "Committer Name",
				committerEmail: "committer@email.com",
			},
			msg: "Unit test message",
			wantMessage: func(docsDir string) string {
				return "Unit test message"
			},
		},
		// if parameters are empty, values from source commit are used
		{
			params: commitUserParam{},
			want:   srcRepoParams,
			msg:    "Unit test message",
			wantMessage: func(docsDir string) string {
				return "Unit test message"
			},
		},
		// templating commit message works
		{
			params: commitUserParam{},
			want:   srcRepoParams,
			msg:    `CommitID: {{.CommitID}}, CommitTime: {{.CommitTime.Unix}}`,
			wantMessage: func(docsDir string) string {
				commitID := gitCommitID(t, docsDir)
				commitTime := gitCommitTime(t, docsDir)
				return fmt.Sprintf("CommitID: %s, CommitTime: %s", commitID, commitTime)
			},
		},
		// invalid template message is used as string literal
		{
			params: commitUserParam{},
			want:   srcRepoParams,
			msg:    "CommitID: {{.CommitID}",
			wantMessage: func(docsDir string) string {
				return "CommitID: {{.CommitID}"
			},
		},
	} {
		currCaseTmpDir, err := ioutil.TempDir(tmpDir, fmt.Sprintf("case-%d-", i))
		require.NoError(t, err)

		// src repo
		githubWikiSrcRepo, err := ioutil.TempDir(currCaseTmpDir, "github-wiki-")
		require.NoError(t, err)
		gittest.InitGitDir(t, githubWikiSrcRepo)

		// add commit with specific user and committer
		err = ioutil.WriteFile(path.Join(githubWikiSrcRepo, "foo.txt"), []byte("foo"), 0644)
		require.NoError(t, err)
		err = git(githubWikiSrcRepo).commitAll("Original commit message", srcRepoParams)
		require.NoError(t, err)

		// bare version of source repo
		githubWikiBareRepo := gitCloneBare(t, githubWikiSrcRepo, currCaseTmpDir)

		// docs directory
		docsDir, err := ioutil.TempDir(tmpDir, "docs-")
		require.NoError(t, err)
		gittest.InitGitDir(t, docsDir)
		err = ioutil.WriteFile(path.Join(docsDir, "page.md"), []byte("Test page"), 0644)
		require.NoError(t, err)
		gittest.CommitAllFiles(t, docsDir, "Initial docs directory commit")

		err = SyncGitHubWiki(Params{
			DocsDir:        docsDir,
			Repo:           githubWikiBareRepo,
			AuthorName:     currCase.params.authorName,
			AuthorEmail:    currCase.params.authorEmail,
			CommitterName:  currCase.params.committerName,
			CommitterEmail: currCase.params.committerEmail,
			Msg:            currCase.msg,
		}, ioutil.Discard)
		require.NoError(t, err, "Case %d", i)

		got := commitUserParamForRepo(t, git(githubWikiBareRepo))
		assert.Equal(t, currCase.wantMessage(docsDir), gitMessage(t, githubWikiBareRepo), "Case %d", i)
		assert.Equal(t, currCase.want, got, "Case %d", i)
	}
}

// Tests operations when documents directory being published is not in a Git repository.
func TestSyncGitHubWikiNonGitDocsDir(t *testing.T) {
	tmpDir, cleanup, err := dirs.TempDir("", "")
	defer cleanup()
	require.NoError(t, err)

	srcRepoParams := commitUserParam{
		authorName:     "src-author",
		authorEmail:    "src-author@email.com",
		committerName:  "src-committer",
		committerEmail: "src-committer@email.com",
	}

	for i, currCase := range []struct {
		params           commitUserParam
		msg              string
		want             commitUserParam
		wantMessage      func(docsDir string) string
		wantStdoutRegexp string
	}{
		// if input directory is not in a Git repository, templating will not work. Warning is printed to stdout, but operation still completes.
		{
			params: commitUserParam{},
			want:   srcRepoParams,
			msg:    `CommitID: {{.CommitID}}, CommitTime: {{.CommitTime.Unix}}`,
			wantMessage: func(docsDir string) string {
				return `CommitID: {{.CommitID}}, CommitTime: {{.CommitTime.Unix}}`
			},
			wantStdoutRegexp: "(?s).+Continuing with templating disabled. To fix this issue, ensure that the directory is in a Git repository.",
		},
	} {
		currCaseTmpDir, err := ioutil.TempDir(tmpDir, fmt.Sprintf("case-%d-", i))
		require.NoError(t, err)

		// src repo
		githubWikiSrcRepo, err := ioutil.TempDir(currCaseTmpDir, "github-wiki-")
		require.NoError(t, err)
		gittest.InitGitDir(t, githubWikiSrcRepo)

		// add commit with specific user and committer
		err = ioutil.WriteFile(path.Join(githubWikiSrcRepo, "foo.txt"), []byte("foo"), 0644)
		require.NoError(t, err)
		err = git(githubWikiSrcRepo).commitAll("Original commit message", srcRepoParams)
		require.NoError(t, err)

		// bare version of source repo
		githubWikiBareRepo := gitCloneBare(t, githubWikiSrcRepo, currCaseTmpDir)

		// docs directory
		docsDir, err := ioutil.TempDir(tmpDir, "docs-")
		require.NoError(t, err)
		err = ioutil.WriteFile(path.Join(docsDir, "page.md"), []byte("Test page"), 0644)
		require.NoError(t, err)

		buf := &bytes.Buffer{}
		err = SyncGitHubWiki(Params{
			DocsDir:        docsDir,
			Repo:           githubWikiBareRepo,
			AuthorName:     currCase.params.authorName,
			AuthorEmail:    currCase.params.authorEmail,
			CommitterName:  currCase.params.committerName,
			CommitterEmail: currCase.params.committerEmail,
			Msg:            currCase.msg,
		}, buf)
		require.NoError(t, err, "Case %d", i)

		got := commitUserParamForRepo(t, git(githubWikiBareRepo))
		assert.Equal(t, currCase.wantMessage(docsDir), gitMessage(t, githubWikiBareRepo), "Case %d", i)
		assert.Equal(t, currCase.want, got, "Case %d", i)
		assert.Regexp(t, regexp.MustCompile(currCase.wantStdoutRegexp), buf.String(), "Case %d", i)
	}
}

func commitUserParamForRepo(t *testing.T, g git) commitUserParam {
	authorName, err := g.valueOr("", authorNameParam)
	require.NoError(t, err)
	authorEmail, err := g.valueOr("", authorEmailParam)
	require.NoError(t, err)
	committerName, err := g.valueOr("", committerNameParam)
	require.NoError(t, err)
	committerEmail, err := g.valueOr("", committerEmailParam)
	require.NoError(t, err)
	return commitUserParam{
		authorName:     authorName,
		authorEmail:    authorEmail,
		committerName:  committerName,
		committerEmail: committerEmail,
	}
}

func gitCloneBare(t *testing.T, srcGitRepo, tmpDir string) string {
	bareRepo, err := ioutil.TempDir(tmpDir, "bare-")
	require.NoError(t, err)
	cmd := exec.Command("git", "clone", "--bare", srcGitRepo, bareRepo)
	output, err := cmd.CombinedOutput()
	require.NoError(t, err, "Failed to execute %v: %v", cmd.Args, string(output))
	return bareRepo
}

func gitMessage(t *testing.T, gitRepo string) string {
	return gitCommand(t, gitRepo, "show", "-s", "--format=%B")
}

func gitCommitID(t *testing.T, gitRepo string) string {
	return gitCommand(t, gitRepo, "rev-parse", "HEAD")
}

func gitCommitTime(t *testing.T, gitRepo string) string {
	return gitCommand(t, gitRepo, "show", "-s", "--format=%ct")
}

func gitCommand(t *testing.T, gitRepo string, args ...string) string {
	cmd := exec.Command("git", args...)
	cmd.Dir = gitRepo
	output, err := cmd.CombinedOutput()
	require.NoError(t, err, "Failed to execute %v: %v", cmd.Args, string(output))
	return strings.TrimSpace(string(output))
}
