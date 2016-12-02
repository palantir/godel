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

package binspec_test

import (
	"io/ioutil"
	"os"
	"path"
	"regexp"
	"testing"

	"github.com/nmiyake/pkg/dirs"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/palantir/godel/apps/distgo/pkg/binspec"
	"github.com/palantir/godel/apps/distgo/pkg/osarch"
)

func TestBinSpecCreateDirectoryStructureFailBadRootDirectory(t *testing.T) {
	tmp, cleanup, err := dirs.TempDir("", "")
	defer cleanup()
	require.NoError(t, err)

	spec := binspec.New([]osarch.OSArch{
		{
			OS:   "darwin",
			Arch: "amd-64",
		},
	}, "testExecutable")
	err = spec.CreateDirectoryStructure(tmp, nil, false)
	require.Error(t, err)

	assert.Regexp(t, regexp.MustCompile(`^.+ is not a path to bin$`), err.Error())
}

func TestBinSpecCreateStructure(t *testing.T) {
	tmp, cleanup, err := dirs.TempDir("", "")
	defer cleanup()
	require.NoError(t, err)

	for i, currCase := range []struct {
		input         []osarch.OSArch
		expectedPaths []string
	}{
		{
			input:         []osarch.OSArch{{OS: "darwin", Arch: "amd64"}},
			expectedPaths: []string{"bin/darwin-amd64"},
		},
		{
			input:         []osarch.OSArch{{OS: "darwin", Arch: "amd64"}, {OS: "linux", Arch: "amd64"}},
			expectedPaths: []string{"bin/darwin-amd64", "bin/linux-amd64"},
		},
	} {
		currTmpDir, err := ioutil.TempDir(tmp, "")
		require.NoError(t, err)

		currRoot := path.Join(currTmpDir, "bin")
		err = os.Mkdir(currRoot, 0755)
		require.NoError(t, err)

		spec := binspec.New(currCase.input, "testExecutable")
		err = spec.CreateDirectoryStructure(currRoot, nil, false)
		require.NoError(t, err)

		for _, currPath := range currCase.expectedPaths {
			_, err := os.Stat(path.Join(currTmpDir, currPath))
			assert.NoError(t, err, "Case %d", i)
		}
	}
}
