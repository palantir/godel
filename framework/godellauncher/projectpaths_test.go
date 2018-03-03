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

package godellauncher_test

import (
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"testing"

	"github.com/nmiyake/pkg/dirs"
	"github.com/nmiyake/pkg/gofiles"
	"github.com/palantir/pkg/matcher"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/palantir/godel/framework/godellauncher"
)

func TestListProjectPaths(t *testing.T) {
	tmpDir, cleanup, err := dirs.TempDir("", "")
	require.NoError(t, err)
	defer cleanup()

	origWd, err := os.Getwd()
	require.NoError(t, err)

	for i, tc := range []struct {
		name    string
		files   []gofiles.GoFileSpec
		wd      string
		include matcher.Matcher
		exclude matcher.Matcher
		want    []string
	}{
		{
			"empty matcher matches nothing",
			[]gofiles.GoFileSpec{
				{
					RelPath: "foo.go",
				},
			},
			".",
			nil,
			nil,
			nil,
		},
		{
			"matcher matches files and directories",
			[]gofiles.GoFileSpec{
				{
					RelPath: "foo.go",
				},
				{
					RelPath: "bar/bar.go",
				},
			},
			".",
			matcher.Name(`.+`),
			nil,
			[]string{
				"bar",
				"bar/bar.go",
				"foo.go",
			},
		},
		{
			"matcher returns relative paths",
			[]gofiles.GoFileSpec{
				{
					RelPath: "foo.go",
				},
				{
					RelPath: "bar/bar.go",
				},
			},
			"bar",
			matcher.Name(`.+`),
			nil,
			[]string{
				"../bar",
				"../bar/bar.go",
				"../foo.go",
			},
		},
		{
			"exclude matcher is used",
			[]gofiles.GoFileSpec{
				{
					RelPath: "foo.go",
				},
				{
					RelPath: "bar/bar.go",
				},
			},
			"bar",
			matcher.Name(`.+`),
			matcher.Name(`bar.go`),
			[]string{
				"../bar",
				"../foo.go",
			},
		},
	} {
		projectDir, err := ioutil.TempDir(tmpDir, "project")
		require.NoError(t, err)
		projectDir, err = filepath.EvalSymlinks(projectDir)
		require.NoError(t, err)

		_, err = gofiles.Write(projectDir, tc.files)
		require.NoError(t, err)

		func() {
			err = os.Chdir(path.Join(projectDir, tc.wd))
			require.NoError(t, err)
			defer func() {
				err = os.Chdir(origWd)
				require.NoError(t, err)
			}()

			got, err := godellauncher.ListProjectPaths(projectDir, tc.include, tc.exclude)
			require.NoError(t, err)
			assert.Equal(t, tc.want, got, "Case %d: %s", i, tc.name)
		}()
	}
}
