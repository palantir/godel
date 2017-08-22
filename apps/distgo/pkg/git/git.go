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

package git

import (
	"os/exec"
	"regexp"
	"strings"

	"github.com/pkg/errors"
)

const Unspecified = "unspecified"

type ProjectInfo struct {
	Version  string
	Branch   string
	Revision string
}

const snapshotRegexp = `.+g[-+.]?[a-fA-F0-9]{3,}$`

func IsSnapshotVersion(version string) bool {
	return regexp.MustCompile(snapshotRegexp).MatchString(version)
}

func NewProjectInfo(gitDir string) (ProjectInfo, error) {
	version, err := ProjectVersion(gitDir)
	if err != nil {
		return ProjectInfo{}, err
	}

	branch, err := ProjectBranch(gitDir)
	if err != nil {
		return ProjectInfo{}, err
	}

	revision, err := ProjectRevision(gitDir)
	if err != nil {
		return ProjectInfo{}, err
	}

	return ProjectInfo{
		Version:  version,
		Branch:   branch,
		Revision: revision,
	}, nil
}

// ProjectVersion returns the version string for the git repository that the provided directory is in. The output is the output
// of "git describe --tags" followed by ".dirty" if the repository currently has any uncommitted changes. Returns
// an error if the provided path is not in a git root or if the git repository has no commits or no tags.
func ProjectVersion(gitDir string) (string, error) {
	tags, err := tags(gitDir)
	if err != nil {
		return "", err
	}

	// if no tags exist, return Unspecified as the version
	if tags == "" {
		return Unspecified, nil
	}

	result, err := trimmedCombinedGitCmdOutput(gitDir, "describe", "--tags", "--first-parent")
	if err != nil {
		return "", err
	}

	// trim "v" prefix in tags
	if strings.HasPrefix(result, "v") {
		result = result[1:]
	}

	// handle untracked files as well as "actual" dirtiness
	dirtyFiles, err := trimmedCombinedGitCmdOutput(gitDir, "status", "--porcelain")
	if err != nil {
		return "", err
	}
	if dirtyFiles != "" {
		result += ".dirty"
	}
	return result, nil
}

func ProjectBranch(gitDir string) (string, error) {
	tags, err := tags(gitDir)
	if err != nil {
		return "", err
	}

	// if no tags exist, return Unspecified as the branch
	if tags == "" {
		return Unspecified, nil
	}

	branch, err := branch(gitDir)
	if err != nil {
		return "", err
	}

	if strings.HasPrefix(branch, "v") {
		branch = branch[1:]
	}

	return branch, nil
}

func ProjectRevision(gitDir string) (string, error) {
	tags, err := tags(gitDir)
	if err != nil {
		return "", err
	}

	// if no tags exist, return revision count from first commit
	if tags == "" {
		return trimmedCombinedGitCmdOutput(gitDir, "rev-list", "HEAD", "--count")
	}

	branch, err := branch(gitDir)
	if err != nil {
		return "", err
	}
	return trimmedCombinedGitCmdOutput(gitDir, "rev-list", branch+"..HEAD", "--count")
}

func tags(gitDir string) (string, error) {
	return trimmedCombinedGitCmdOutput(gitDir, "tag", "-l")
}

func branch(gitDir string) (string, error) {
	return trimmedCombinedGitCmdOutput(gitDir, "describe", "--abbrev=0", "--tags", "--first-parent")
}

func trimmedCombinedGitCmdOutput(gitDir string, args ...string) (string, error) {
	cmd := exec.Command("git", args...)
	cmd.Dir = gitDir
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", errors.Wrapf(err, "Command %v failed. Output: %v", cmd.Args, string(out))
	}
	return strings.TrimSpace(string(out)), err
}
