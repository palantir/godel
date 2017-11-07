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

package layout_test

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"runtime"
	"testing"

	"github.com/nmiyake/pkg/dirs"
	"github.com/palantir/pkg/specdir"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/palantir/godel/framework/builtintasks/installupdate/layout"
)

func TestAppLayoutNoValidation(t *testing.T) {
	specDir, err := specdir.New("testRoot", layout.AppSpec(), layout.AppSpecTemplate("0.0.1"), specdir.SpecOnly)
	require.NoError(t, err)

	for i, currCase := range []struct {
		aliasName string
		want      string
	}{
		{
			aliasName: layout.AppExecutable,
			want:      fmt.Sprintf("godel-0.0.1/bin/%v-%v/godel", runtime.GOOS, runtime.GOARCH),
		},
	} {
		actual := specDir.Path(currCase.aliasName)
		assert.Equal(t, currCase.want, actual, "Case %d", i)
	}
}

func TestAppSpecLayoutValidationFail(t *testing.T) {
	tmpDir, cleanup, err := dirs.TempDir("", "")
	defer cleanup()
	require.NoError(t, err)

	gödelDir := path.Join(tmpDir, "godel-0.0.1")
	err = os.Mkdir(gödelDir, 0755)
	require.NoError(t, err)

	for i, currCase := range []struct {
		rootName  string
		wantError string
	}{
		{
			rootName:  "testRoot",
			wantError: `testRoot is not a path to godel-0.0.1`,
		},
		{
			rootName:  gödelDir,
			wantError: `godel-0.0.1/bin does not exist`,
		},
	} {
		spec := layout.AppSpec()
		err = spec.Validate(currCase.rootName, layout.AppSpecTemplate("0.0.1"))
		assert.Error(t, err, fmt.Sprintf("Case %d", i))

		if err != nil {
			assert.EqualError(t, err, currCase.wantError, "Case %d", i)
		}
	}
}

func TestAppLayoutValidation(t *testing.T) {
	tmpDir, cleanup, err := dirs.TempDir("", "")
	defer cleanup()
	require.NoError(t, err)

	filesToCreate := map[string]string{
		path.Join("godel-0.0.1", "bin", "darwin-amd64", "godel"):          "godel",
		path.Join("godel-0.0.1", "bin", "linux-amd64", "godel"):           "godel",
		path.Join("godel-0.0.1", "wrapper", "godel", "bin", "godelw"):     "godelw",
		path.Join("godel-0.0.1", "wrapper", "godel", "config", "foo.yml"): "testconfig",
		path.Join("godel-0.0.1", "wrapper", "godelw"):                     "godelw",
	}

	createdFilesTmpDir := createFiles(t, tmpDir, filesToCreate)
	specDir, err := specdir.New(path.Join(createdFilesTmpDir, "godel-0.0.1"), layout.AppSpec(), layout.AppSpecTemplate("0.0.1"), specdir.Validate)
	require.NoError(t, err)

	for i, currCase := range []struct {
		aliasName string
		want      string
	}{
		{
			aliasName: layout.AppExecutable,
			want:      fmt.Sprintf("godel-0.0.1/bin/%v-%v/godel", runtime.GOOS, runtime.GOARCH),
		},
	} {
		expected := path.Join(createdFilesTmpDir, currCase.want)
		got := specDir.Path(currCase.aliasName)
		assert.Equal(t, expected, got, "Case %d", i)
	}
}

func createFiles(t *testing.T, tmpDir string, files map[string]string) string {
	currCaseTmpDir, err := ioutil.TempDir(tmpDir, "")
	require.NoError(t, err)

	for currFile, currContent := range files {
		err = os.MkdirAll(path.Join(currCaseTmpDir, path.Dir(currFile)), 0755)
		require.NoError(t, err)
		err = ioutil.WriteFile(path.Join(currCaseTmpDir, currFile), []byte(currContent), 0644)
		require.NoError(t, err)
	}

	return currCaseTmpDir
}
