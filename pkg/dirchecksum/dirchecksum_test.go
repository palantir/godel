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

package dirchecksum_test

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"testing"

	"github.com/nmiyake/pkg/dirs"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/palantir/godel/pkg/dirchecksum"
)

func TestChecksumsForMatchingPaths(t *testing.T) {
	tmpDir, cleanup, err := dirs.TempDir(".", "")
	require.NoError(t, err)
	defer cleanup()

	dir1 := path.Join(tmpDir, "dir1")
	createTestFiles(t, dir1)
	checksums1, err := dirchecksum.ChecksumsForMatchingPaths(dir1, nil)
	require.NoError(t, err)

	dir2 := path.Join(tmpDir, "dir2")
	createTestFiles(t, dir2)
	checksums2, err := dirchecksum.ChecksumsForMatchingPaths(dir2, nil)
	require.NoError(t, err)

	assert.Equal(t, checksums1.Checksums, checksums2.Checksums)

	diff := checksums1.Diff(checksums2)
	assert.Equal(t, 0, len(diff.Diffs))
}

func TestChecksumsDiffString(t *testing.T) {
	tmpDir, cleanup, err := dirs.TempDir(".", "")
	require.NoError(t, err)
	defer cleanup()

	dir1 := path.Join(tmpDir, "dir1")
	createTestFiles(t, dir1)

	const originalFileName = "original.txt"
	missingFile := path.Join(dir1, originalFileName)
	err = ioutil.WriteFile(missingFile, []byte("original"), 0644)
	require.NoError(t, err)

	checksumDiffFile := path.Join(dir1, "checksumdiff.txt")
	err = ioutil.WriteFile(checksumDiffFile, []byte("original"), 0644)
	require.NoError(t, err)

	checksums1, err := dirchecksum.ChecksumsForMatchingPaths(dir1, nil)
	require.NoError(t, err)

	dir2 := path.Join(tmpDir, "dir2")
	createTestFiles(t, dir2)

	const newFileName = "new.txt"
	extraFile := path.Join(dir2, newFileName)
	err = ioutil.WriteFile(extraFile, []byte("new"), 0644)
	require.NoError(t, err)

	const diffFileName = "checksumdiff.txt"
	checksumDiffFile = path.Join(dir2, diffFileName)
	err = ioutil.WriteFile(checksumDiffFile, []byte("new"), 0644)
	require.NoError(t, err)

	checksums2, err := dirchecksum.ChecksumsForMatchingPaths(dir2, nil)
	require.NoError(t, err)

	diff := checksums1.Diff(checksums2)

	want := fmt.Sprintf(`%s/%s: checksum changed from 0682c5f2076f099c34cfdd15a9e063849ed437a49677e6fcc5b4198c76575be5 to 11507a0e2f5e69d5dfa40a62a1bd7b6ee57e6bcd85c67c9b8431b36fff21c437
%s/%s: extra
%s/%s: missing`, dir1, diffFileName, dir1, newFileName, dir1, originalFileName)
	assert.Equal(t, want, diff.String())
}

func createTestFiles(t *testing.T, rootDir string) (string, string) {
	err := os.MkdirAll(rootDir, 0755)
	require.NoError(t, err)

	testFile := path.Join(rootDir, "testfile.txt")
	err = ioutil.WriteFile(testFile, []byte("foo"), 0644)
	require.NoError(t, err)

	testInnerFile := path.Join(rootDir, "dir", "innerfile.txt")
	err = os.MkdirAll(path.Dir(testInnerFile), 0755)
	require.NoError(t, err)
	err = ioutil.WriteFile(testInnerFile, []byte("bar"), 0644)
	require.NoError(t, err)

	return testFile, testInnerFile
}
