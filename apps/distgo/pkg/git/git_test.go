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

package git_test

import (
	"io/ioutil"
	"path"
	"testing"

	"github.com/nmiyake/pkg/dirs"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/palantir/godel/apps/distgo/pkg/git"
	"github.com/palantir/godel/apps/distgo/pkg/git/gittest"
)

func TestProjectInfo(t *testing.T) {
	tmp, cleanup, err := dirs.TempDir("", "")
	defer cleanup()
	require.NoError(t, err)

	for i, currCase := range []struct {
		gitOperations func(gitDir string)
		want          git.ProjectInfo
	}{
		{
			gitOperations: func(gitDir string) {
			},
			want: git.ProjectInfo{
				Version:  "unspecified$",
				Branch:   "unspecified$",
				Revision: "1$",
			},
		},
		{
			gitOperations: func(gitDir string) {
				gittest.CommitRandomFile(t, gitDir, "Second commit")
				err = ioutil.WriteFile(path.Join(gitDir, "foo"), []byte("foo"), 0644)
				require.NoError(t, err)
			},
			want: git.ProjectInfo{
				Version:  "unspecified$",
				Branch:   "unspecified$",
				Revision: "2$",
			},
		},
		{
			gitOperations: func(gitDir string) {
				gittest.CreateGitTag(t, gitDir, "0.0.1")
			},
			want: git.ProjectInfo{
				Version:  "0.0.1$",
				Branch:   "0.0.1$",
				Revision: "0$",
			},
		},
		{
			gitOperations: func(gitDir string) {
				gittest.CreateGitTag(t, gitDir, "v0.0.1")
			},
			want: git.ProjectInfo{
				Version:  "0.0.1$",
				Branch:   "0.0.1$",
				Revision: "0$",
			},
		},
		{
			gitOperations: func(gitDir string) {
				gittest.CreateGitTag(t, gitDir, "0.0.1")
				err = ioutil.WriteFile(path.Join(gitDir, "foo"), []byte("foo"), 0644)
				require.NoError(t, err)
			},
			want: git.ProjectInfo{
				Version:  "0.0.1.dirty$",
				Branch:   "0.0.1$",
				Revision: "0$",
			},
		},
		{
			gitOperations: func(gitDir string) {
				gittest.CreateGitTag(t, gitDir, "0.0.1")
				gittest.CommitRandomFile(t, gitDir, "Second commit")
			},
			want: git.ProjectInfo{
				Version:  "0.0.1-1-g[a-f0-9]{7}$",
				Branch:   "0.0.1$",
				Revision: "1$",
			},
		},
		{
			gitOperations: func(gitDir string) {
				gittest.CreateGitTag(t, gitDir, "0.0.1")
				gittest.CommitRandomFile(t, gitDir, "Second commit")
				err = ioutil.WriteFile(path.Join(gitDir, "foo"), []byte("foo"), 0644)
				require.NoError(t, err)
			},
			want: git.ProjectInfo{
				Version: "0.0.1-1-g[a-f0-9]{7}.dirty$",
			},
		},
		{
			gitOperations: func(gitDir string) {
				gittest.CreateGitTag(t, gitDir, "0.0.1")

				gittest.CreateBranch(t, gitDir, "hotfix-branch")
				gittest.CommitRandomFile(t, gitDir, "hotfix commit")
				gittest.CreateGitTag(t, gitDir, "0.0.1-hotfix")

				gittest.RunGitCommand(t, gitDir, "checkout", "master")
				gittest.Merge(t, gitDir, "hotfix-branch")
			},
			want: git.ProjectInfo{
				Version: "^0.0.1-1-g[a-f0-9]{7}$",
			},
		},
		{
			gitOperations: func(gitDir string) {
				gittest.CreateGitTag(t, gitDir, "0.0.1")

				gittest.CreateBranch(t, gitDir, "hotfix-branch")
				gittest.CommitRandomFile(t, gitDir, "hotfix commit")
				gittest.CreateGitTag(t, gitDir, "0.0.1-hotfix")

				gittest.RunGitCommand(t, gitDir, "checkout", "master")
				gittest.Merge(t, gitDir, "hotfix-branch")

				gittest.CreateGitTag(t, gitDir, "0.0.2")
			},
			want: git.ProjectInfo{
				Version: "^0.0.2$",
			},
		},
	} {
		currTmp, err := ioutil.TempDir(tmp, "")
		require.NoError(t, err)

		gittest.InitGitDir(t, currTmp)
		currCase.gitOperations(currTmp)

		got, err := git.NewProjectInfo(currTmp)
		require.NoError(t, err)

		assert.Regexp(t, currCase.want.Version, got.Version, "Case %d", i)
		assert.Regexp(t, currCase.want.Branch, got.Branch, "Case %d", i)
		assert.Regexp(t, currCase.want.Revision, got.Revision, "Case %d", i)
	}
}

func TestIsSnapshotVersion(t *testing.T) {
	for i, currCase := range []struct {
		version    string
		isSnapshot bool
	}{
		{"0.1.0-2-g0f9fa0a", true},
		{"0.1.0-rc1-2-g0f9fa0a", true},
		{"0.0.1", false},
		{"0.0.1-rc1", false},
		{"0.0.1-rc1.dirty", false},
	} {
		assert.Equal(t, currCase.isSnapshot, git.IsSnapshotVersion(currCase.version), "Case %d", i)
	}
}
