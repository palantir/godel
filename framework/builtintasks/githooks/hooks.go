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

package githooks

import (
	"fmt"
	"io/ioutil"
	"path"

	"github.com/pkg/errors"

	"github.com/palantir/godel/v2/framework/builtintasks/installupdate/layout"
)

var hooks = map[string]string{
	"pre-commit": `#!/bin/bash
gofiles=$(git diff --cached --name-only --diff-filter=ACM | grep '\.go$')
[ -z "$gofiles" ] && exit 0

unformatted=$(./godelw format --verify $gofiles)
exitCode=$?
[ "$exitCode" -eq "0" ] && exit 0

if [ -n "$unformatted" ]; then
  echo "Unformatted files exist -- run ./godelw format to format these files:"
  for file in $unformatted; do
    echo "  $file"
  done
fi

exit $exitCode
`,
}

func InstallGitHooks(rootDir string) error {
	gitDir := path.Join(rootDir, ".git")
	if err := layout.VerifyDirExists(gitDir); err != nil {
		return fmt.Errorf(".git directory does not exist at %v", gitDir)
	}

	for hook, contents := range hooks {
		hookPath := path.Join(gitDir, "hooks", hook)
		if err := ioutil.WriteFile(hookPath, []byte(contents), 0755); err != nil {
			return errors.Wrapf(err, "failed to write %s hook to %s", hook, hookPath)
		}
	}

	return nil
}
