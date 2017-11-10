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

package installupdate

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"path"
	"testing"

	"github.com/nmiyake/archiver"
	"github.com/nmiyake/pkg/dirs"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDownload(t *testing.T) {
	tmpDir, cleanup, err := dirs.TempDir("", "")
	defer cleanup()
	require.NoError(t, err)

	tgzFile, err := ioutil.TempFile(tmpDir, "")
	require.NoError(t, err)
	err = tgzFile.Close()
	require.NoError(t, err)
	err = writeSimpleTestTgz(tgzFile.Name())
	require.NoError(t, err)

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		bytes, err := ioutil.ReadFile(tgzFile.Name())
		require.NoError(t, err)
		_, err = w.Write(bytes)
		require.NoError(t, err)
	}))
	defer ts.Close()

	fileName, err := getPkg(remotePkg{
		Pkg: ts.URL,
	}, tmpDir, ioutil.Discard)
	require.NoError(t, err)

	err = archiver.UntarGz(fileName, tmpDir)
	require.NoError(t, err)

	fileBytes, err := ioutil.ReadFile(path.Join(tmpDir, "test.txt"))
	require.NoError(t, err)

	assert.Equal(t, "Test file\n", string(fileBytes))
}

func writeSimpleTestTgz(filePath string) error {
	tmpDir, cleanup, err := dirs.TempDir("", "")
	defer cleanup()
	if err != nil {
		return errors.Wrapf(err, "failed to create temp directory")
	}

	testFilePath := path.Join(tmpDir, "test.txt")
	err = ioutil.WriteFile(testFilePath, []byte("Test file\n"), 0644)
	if err != nil {
		return errors.Wrapf(err, "failed to write file %s", testFilePath)
	}

	err = archiver.TarGz(filePath, []string{testFilePath})
	if err != nil {
		return errors.Wrapf(err, "failed to compress tgz")
	}

	return nil
}
