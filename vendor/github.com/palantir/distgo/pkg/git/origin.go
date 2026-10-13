// Copyright 2026 Palantir Technologies, Inc.
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
	"net/url"
	"strings"

	giturls "github.com/chainguard-dev/git-urls"
)

// RepoOrigin returns the origin remote of the Git repository containing dir and its HTTPS web URL.
// Git failures or an empty origin return empty strings. If parsing fails, only the raw remote is returned.
// Local paths are converted using the same URL rules as other remotes, without requiring a host.
func RepoOrigin(dir string) (gitURL, webURL string) {
	remote, err := CmdOutput(dir, "remote", "get-url", "origin")
	if err != nil || remote == "" {
		return "", ""
	}
	return parseRepoOrigin(remote)
}

func parseRepoOrigin(remote string) (gitURL, webURL string) {
	u, err := giturls.Parse(remote)
	if err != nil {
		return remote, ""
	}
	path := strings.TrimSuffix(strings.TrimSuffix(u.Path, "/"), ".git")
	return remote, (&url.URL{Scheme: "https", Host: u.Host, Path: path}).String()
}
