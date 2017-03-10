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

package gittest

import (
	"io/ioutil"
	"os/exec"
	"testing"

	"github.com/stretchr/testify/require"
)

func InitGitDir(t *testing.T, gitDir string) {
	RunGitCommand(t, gitDir, "init")
	RunGitCommand(t, gitDir, "config", "user.email", "test@author.com")
	RunGitCommand(t, gitDir, "config", "user.name", "testAuthor")
	CommitRandomFile(t, gitDir, "Initial commit")
}

func CreateGitTag(t *testing.T, gitDir, tagValue string) {
	RunGitCommand(t, gitDir, "tag", "-a", tagValue, "-m", "")
}

func CreateBranch(t *testing.T, gitDir, branch string) {
	RunGitCommand(t, gitDir, "checkout", "-b", branch)
}

func CommitAllFiles(t *testing.T, gitDir, commitMessage string) {
	RunGitCommand(t, gitDir, "add", ".")
	RunGitCommand(t, gitDir, "commit", "--author=testAuthor <test@author.com>", "-m", commitMessage)
}

func CommitRandomFile(t *testing.T, gitDir, commitMessage string) {
	file, err := ioutil.TempFile(gitDir, "random-file-")
	require.NoError(t, err)
	require.NoError(t, file.Close())
	CommitAllFiles(t, gitDir, commitMessage)
}

func Merge(t *testing.T, gitDir, branch string) {
	RunGitCommand(t, gitDir, "merge", "--no-ff", branch)
}

func RunGitCommand(t *testing.T, gitDir string, args ...string) {
	cmd := exec.Command("git", args...)
	cmd.Dir = gitDir
	output, err := cmd.CombinedOutput()
	require.NoError(t, err, string(output))
}
