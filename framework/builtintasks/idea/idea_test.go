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

package idea_test

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"testing"

	"github.com/nmiyake/pkg/dirs"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/palantir/godel/framework/builtintasks/idea"
)

func TestCreateIdeaFiles(t *testing.T) {
	tmpDir, cleanup, err := dirs.TempDir("", "")
	defer cleanup()
	require.NoError(t, err)

	err = idea.CreateIntelliJFiles(tmpDir)
	assert.NoError(t, err)

	verifyXMLHelper(t, ideaFilePath(tmpDir, "iml"))
	verifyXMLHelper(t, ideaFilePath(tmpDir, "ipr"))
}

func TestCleanIdeaFiles(t *testing.T) {
	tmpDir, cleanup, err := dirs.TempDir("", "")
	defer cleanup()
	require.NoError(t, err)

	for i, currCase := range [][]string{
		{"iml", "ipr", "iws"},
		// test case where not all files to be cleaned exist
		{"iml", "ipr"},
	} {
		currDir, err := ioutil.TempDir(tmpDir, "")
		require.NoError(t, err)

		for _, ext := range currCase {
			currPath := ideaFilePath(currDir, ext)
			err := ioutil.WriteFile(currPath, []byte(ext), 0644)
			require.NoError(t, err, "Case %d: failed to write %v", i, currPath)
		}

		err = idea.CleanIDEAFiles(currDir)
		require.NoError(t, err)

		for _, ext := range currCase {
			currPath := ideaFilePath(currDir, ext)
			_, err = os.Stat(currPath)
			assert.True(t, os.IsNotExist(err), "Case %d: did not expect %v to exist", i, currPath)
		}
	}
}

func verifyXMLHelper(t *testing.T, fPath string) {
	fInfo, err := os.Stat(fPath)
	assert.NoError(t, err)
	assert.False(t, fInfo.IsDir())
}

func ideaFilePath(dir, ext string) string {
	return path.Join(dir, fmt.Sprintf("%v.%v", path.Base(dir), ext))
}
