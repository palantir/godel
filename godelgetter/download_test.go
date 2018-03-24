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

package godelgetter_test

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"path"
	"regexp"
	"testing"

	"github.com/mholt/archiver"
	"github.com/nmiyake/pkg/dirs"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/palantir/godel/godelgetter"
)

func TestDownloadIntoDirectory(t *testing.T) {
	for i, tc := range []struct {
		setup func(t *testing.T, repoTGZFile string) (srcPath string, cleanup func())
	}{
		{
			func(t *testing.T, repoTGZFile string) (srcPath string, cleanup func()) {
				ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					bytes, err := ioutil.ReadFile(repoTGZFile)
					require.NoError(t, err)
					_, err = w.Write(bytes)
					require.NoError(t, err)
				}))
				return ts.URL + "/test-on-server.tgz", ts.Close
			},
		},
		{
			func(t *testing.T, repoTGZFile string) (srcPath string, cleanup func()) {
				return repoTGZFile, nil
			},
		},
	} {
		func() {
			tmpDir, cleanup, err := dirs.TempDir("", "")
			defer cleanup()
			require.NoError(t, err)

			repoDir := path.Join(tmpDir, "repo")
			err = os.MkdirAll(repoDir, 0755)
			require.NoError(t, err, "Case %d", i)

			downloadsDir := path.Join(tmpDir, "downloads")
			err = os.MkdirAll(downloadsDir, 0755)
			require.NoError(t, err, "Case %d", i)

			repoTGZFile := path.Join(repoDir, "test.tgz")
			writeSimpleTestTgz(t, repoTGZFile)

			srcPath, cleanup := tc.setup(t, repoTGZFile)
			if cleanup != nil {
				defer cleanup()
			}

			outBytes := &bytes.Buffer{}
			fileName, err := godelgetter.DownloadIntoDirectory(godelgetter.NewPkgSrc(srcPath, ""), downloadsDir, outBytes)
			require.NoError(t, err, "Case %d", i)

			err = archiver.TarGz.Open(fileName, tmpDir)
			require.NoError(t, err, "Case %d", i)

			fileBytes, err := ioutil.ReadFile(path.Join(tmpDir, "test.txt"))
			require.NoError(t, err, "Case %d", i)

			assert.Equal(t, "Test file\n", string(fileBytes), "Case %d", i)
			assert.Regexp(t, fmt.Sprintf("(?s)Getting package from %s", srcPath)+regexp.QuoteMeta("...")+".+", outBytes.String(), "Case %d", i)
		}()
	}
}

func TestDownloadSameFileOK(t *testing.T) {
	tmpDir, cleanup, err := dirs.TempDir("", "")
	require.NoError(t, err)
	defer cleanup()

	testFile := path.Join(tmpDir, "test.txt")
	const content = "test content"
	err = ioutil.WriteFile(testFile, []byte(content), 0644)
	require.NoError(t, err)

	_, err = godelgetter.DownloadIntoDirectory(godelgetter.NewPkgSrc(testFile, ""), tmpDir, ioutil.Discard)
	require.NoError(t, err)

	gotContent, err := ioutil.ReadFile(testFile)
	require.NoError(t, err)

	assert.Equal(t, content, string(gotContent))
}

func TestFailedReaderDoesNotOverwriteDestinationFile(t *testing.T) {
	tmpDir, cleanup, err := dirs.TempDir("", "")
	require.NoError(t, err)
	defer cleanup()

	const fileName = "test.txt"

	srcDir := path.Join(tmpDir, "src")
	err = os.MkdirAll(srcDir, 0755)
	require.NoError(t, err)

	dstDir := path.Join(tmpDir, "dst")
	err = os.MkdirAll(dstDir, 0755)
	require.NoError(t, err)

	// write content in destination file
	dstFile := path.Join(dstDir, fileName)
	const content = "destination content"
	err = ioutil.WriteFile(dstFile, []byte(content), 0644)
	require.NoError(t, err)

	srcFile := path.Join(srcDir, fileName)
	_, err = godelgetter.DownloadIntoDirectory(godelgetter.NewPkgSrc(srcFile, ""), dstDir, os.Stdout)
	// download operation should fail
	assert.EqualError(t, err, fmt.Sprintf("%s does not exist", srcFile))

	// destination file should still exist (should not have been truncated by download with failed source reader)
	gotContent, err := ioutil.ReadFile(dstFile)
	require.NoError(t, err)
	assert.Equal(t, content, string(gotContent))
}

func writeSimpleTestTgz(t *testing.T, filePath string) {
	tmpDir, cleanup, err := dirs.TempDir("", "")
	defer cleanup()
	require.NoError(t, err)

	testFilePath := path.Join(tmpDir, "test.txt")
	err = ioutil.WriteFile(testFilePath, []byte("Test file\n"), 0644)
	require.NoError(t, err)

	err = archiver.TarGz.Make(filePath, []string{testFilePath})
	require.NoError(t, err)
}
