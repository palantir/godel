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
	"strings"

	"github.com/pkg/errors"
)

const Unspecified = "unspecified"

// ProjectVersion returns the version string for the git repository that the provided directory is in. The output is the
// output of "git describe --tags --first-parent" followed by "-dirty" if the repository currently has any uncommitted
// changes (including untracked files) as determined by "git status --porcelain". Returns "unspecified" if the
// repository does not contain any tags.
func ProjectVersion(gitDir string) (string, error) {
	tags, err := Tags(gitDir)
	if err != nil {
		return "", err
	}

	// if no tags exist, return Unspecified as the version
	if tags == "" {
		return Unspecified, nil
	}

	result, err := CmdOutput(gitDir, "describe", "--tags", "--first-parent")
	if err != nil {
		return "", err
	}

	// use "git status --porcelain" rather than "git describe --dirty" to ensure that the existence of untracked files
	// will cause a repository to be considered dirty.
	dirtyFiles, err := CmdOutput(gitDir, "status", "--porcelain")
	if err != nil {
		return "", err
	}
	if dirtyFiles != "" {
		result += "-dirty"
	}
	return result, nil
}

func Tags(gitDir string) (string, error) {
	return CmdOutput(gitDir, "tag", "-l")
}

func CmdOutput(gitDir string, args ...string) (string, error) {
	cmd := exec.Command("git", args...)
	cmd.Dir = gitDir
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", errors.Wrapf(err, "command %v failed with output %v", cmd.Args, string(out))
	}
	return strings.TrimSpace(string(out)), err
}
