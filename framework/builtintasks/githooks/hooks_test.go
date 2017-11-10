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

package githooks_test

import (
	"fmt"
	"io/ioutil"
	"os/exec"
	"path"
	"regexp"
	"testing"

	"github.com/nmiyake/pkg/dirs"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/palantir/godel/framework/builtintasks/githooks"
)

func TestInstallGitHooks(t *testing.T) {
	tmpDir, cleanup, err := dirs.TempDir("", "")
	defer cleanup()
	require.NoError(t, err)

	for i, currCase := range []struct {
		setup    func(projectDir string)
		validate func(projectDir string, caseNum int, err error)
	}{
		{
			setup: func(projectDir string) {
				cmd := exec.Command("git", "init")
				cmd.Dir = projectDir
				err := cmd.Run()
				require.NoError(t, err)
			},
			validate: func(projectDir string, caseNum int, err error) {
				require.NoError(t, err, "Case %d", caseNum)
				bytes, err := ioutil.ReadFile(path.Join(projectDir, ".git/hooks/pre-commit"))
				require.NoError(t, err, "Case %d", caseNum)
				assert.Regexp(t, regexp.MustCompile(`(?s).+\./godelw format -l runAll \$gofiles.+`), string(bytes), "Case %d", caseNum)
			},
		},
		{
			validate: func(projectDir string, caseNum int, err error) {
				expectedErr := fmt.Sprintf(".git directory does not exist at %v", path.Join(projectDir, ".git"))
				assert.EqualError(t, err, expectedErr, "Case %d", caseNum)
			},
		},
	} {
		projectDir, err := ioutil.TempDir(tmpDir, "")
		require.NoError(t, err, "Case %d", i)

		if currCase.setup != nil {
			currCase.setup(projectDir)
		}
		currCase.validate(projectDir, i, githooks.InstallGitHooks(projectDir))
	}
}
