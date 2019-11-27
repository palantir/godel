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
	"fmt"
	"os/exec"
	"regexp"
	"strings"

	"github.com/pkg/errors"
)

const Unspecified = "unspecified"

var fullDescribeRegexp = regexp.MustCompile(`^(.+)-([0-9]+)-g([0-9a-f]{40})$`)

// ProjectVersion returns the version string for the git repository that the provided directory is in. The output is the
// output of "git describe --tags --first-parent" followed by ".dirty" if the repository currently has any uncommitted
// changes (including untracked files) as determined by "git status --porcelain". If the "git describe" output includes
// a trailing commit hash ("-g[0-9a-f]+", where [0-9a-f]+ is the commit hash), then the commit hash will be 7
// characters long even if the "git describe" operation returns a longer hash. If the output starts with the character
// 'v' followed by a digit (0-9), then the leading 'v' is trimmed. Returns "unspecified" if the current commit cannot be
// described.
func ProjectVersion(gitDir string) (string, error) {
	return ProjectVersionWithPrefix(gitDir, "")
}

// ProjectVersionWithPrefix works in the same manner as ProjectVersion, but only matches tags that begin with the
// provided tagPrefix. This can be useful in scenarios where a single repository contains multiple projects and tag
// prefixes (such as "@org/product@") are used to distinguish between releases of different products. The returned
// version includes the prefix.
func ProjectVersionWithPrefix(gitDir, tagPrefix string) (string, error) {
	// use "--long" and "--abbrev=40" to ensure that output is always of the form [tag]-[0-9]+-g[0-9a-f]{40}
	result, err := CmdOutput(gitDir, "describe", "--tags", "--first-parent", "--long", "--abbrev=40", fmt.Sprintf("--match=%s*", tagPrefix))
	if err != nil {
		if strings.HasPrefix(strings.TrimSpace(result), "fatal:") {
			// if output starts with "fatal: ", treat as a Git error ("fatal: No names found, cannot describe anything.",
			// "fatal: No tags can describe '[0-9a-f]{40}'.", etc.).
			return Unspecified, nil
		}
		return "", err
	}

	matchParts := fullDescribeRegexp.FindStringSubmatch(result)
	if matchParts == nil {
		return "", errors.Errorf("output %q does not match regexp %s", result, fullDescribeRegexp.String())
	}

	result = matchParts[1]
	if matchParts[2] != "0" {
		// use only the first 7 characters of the hash to ensure that output is deterministic
		result += fmt.Sprintf("-%s-g%s", matchParts[2], matchParts[3][:7])
	}

	// if tag name starts with "v#", strip the leading 'v'.
	if len(result) >= 2 && result[0] == 'v' && result[1] >= '0' && result[1] <= '9' {
		result = result[1:]
	}

	// use "git status --porcelain" rather than "git describe --dirty" to ensure that the existence of untracked files
	// will cause a repository to be considered dirty.
	dirtyFiles, err := CmdOutput(gitDir, "status", "--porcelain")
	if err != nil {
		return "", err
	}
	if dirtyFiles != "" {
		result += ".dirty"
	}
	return result, nil
}

func CmdOutput(gitDir string, args ...string) (string, error) {
	cmd := exec.Command("git", args...)
	cmd.Dir = gitDir
	outBytes, err := cmd.CombinedOutput()
	out := string(outBytes)
	if err != nil {
		return out, errors.Wrapf(err, "command %v failed with output %v", cmd.Args, out)
	}
	return strings.TrimSpace(out), nil
}
