// Copyright 2016 Palantir Technologies. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gittest

import (
	"io/ioutil"
	"os/exec"
	"runtime/debug"
	"testing"
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
	requireNoError(t, err, "failed to create temporary file")
	requireNoError(t, file.Close(), "failed to close temporary file")
	CommitAllFiles(t, gitDir, commitMessage)
}

func Merge(t *testing.T, gitDir, branch string) {
	RunGitCommand(t, gitDir, "merge", "--no-ff", branch)
}

func RunGitCommand(t *testing.T, gitDir string, args ...string) string {
	cmd := exec.Command("git", args...)
	cmd.Dir = gitDir
	output, err := cmd.CombinedOutput()
	requireNoError(t, err, string(output))
	return string(output)
}

func requireNoError(t *testing.T, err error, msg string) {
	if err == nil {
		return
	}
	t.Errorf("unexpected error: %v: %s%s", err, msg, string(debug.Stack()))
	t.FailNow()
}
