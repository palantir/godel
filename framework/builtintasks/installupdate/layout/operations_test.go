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
	"io/ioutil"
	"os"
	"path"
	"testing"

	"github.com/nmiyake/pkg/dirs"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/palantir/godel/framework/builtintasks/installupdate/layout"
)

type Spec struct {
	Path    string
	IsDir   bool
	Content string
}

type Specs []Spec

func (s Specs) AllPaths() map[string]bool {
	m := make(map[string]bool, len(s))
	for _, c := range s {
		m[c.Path] = c.IsDir
		// if path is a file, populate entries for all parent directories
		if !c.IsDir {
			for d := path.Dir(c.Path); d != "."; d = path.Dir(d) {
				m[d] = true
			}
		}
	}
	return m
}

func Dir(path string) Spec {
	return Spec{
		Path: path,
	}
}

func File(path, content string) Spec {
	return Spec{
		Path:    path,
		Content: content,
	}
}

func TestSyncDir(t *testing.T) {
	tmpDir, cleanup, err := dirs.TempDir("", "")
	defer cleanup()
	require.NoError(t, err)

	for i, currCase := range []struct {
		srcDirLayout  Specs
		dstDirLayout  Specs
		skip          []string
		wantModified  bool
		wantDirLayout Specs
	}{
		// no modification occurs if src and dst the same
		{
			srcDirLayout: []Spec{
				File("both.txt", "src"),
				Dir("both"),
			},
			dstDirLayout: []Spec{
				File("both.txt", "src"),
				Dir("both"),
			},
			wantModified: false,
		},
		// no modification occurs if paths to sync are skipped
		{
			srcDirLayout: []Spec{
				Dir("src-only"),
			},
			dstDirLayout: []Spec{
				File("dst-only.txt", "dst"),
			},
			skip: []string{
				"src-only",
				"dst-only.txt",
			},
			wantModified: false,
		},
		// sync files
		{
			srcDirLayout: []Spec{
				File("src-only.txt", "src-only"),
				File("both.txt", "src"),
			},
			dstDirLayout: []Spec{
				File("dst-only.txt", "dst-only"),
				File("both.txt", "dst"),
			},
			wantDirLayout: []Spec{
				File("src-only.txt", "src-only"),
				File("both.txt", "src"),
			},
			wantModified: true,
		},
		// sync directories
		{
			srcDirLayout: []Spec{
				Dir("src-only"),
				File("both/both.txt", "src"),
			},
			dstDirLayout: []Spec{
				Dir("dst-only"),
				File("both/both.txt", "dst"),
			},
			wantDirLayout: []Spec{
				Dir("src-only"),
				File("both/both.txt", "src"),
			},
			wantModified: true,
		},
		// sync paths with different types
		{
			srcDirLayout: []Spec{
				Dir("dir-in-src"),
				File("file-in-src", "src"),
			},
			dstDirLayout: []Spec{
				File("dir-in-src", "dst"),
				Dir("file-in-src"),
			},
			wantDirLayout: []Spec{
				Dir("dir-in-src"),
				File("file-in-src", "src"),
			},
			wantModified: true,
		},
		// sync operation with multiple actions
		{
			srcDirLayout: []Spec{
				File("foo.txt", "foo"),
				File("baz.txt", "src-baz"),
				Dir("dir-in-src"),
			},
			dstDirLayout: []Spec{
				File("bar.txt", "bar"),
				File("baz.txt", "dst-baz"),
				File("dir-in-src", "file"),
			},
			wantDirLayout: []Spec{
				File("foo.txt", "foo"),
				File("baz.txt", "src-baz"),
				Dir("dir-in-src"),
			},
			wantModified: true,
		},
	} {
		srcDir, err := ioutil.TempDir(tmpDir, "src")
		require.NoError(t, err, "Case %d", i)
		writeLayout(t, srcDir, currCase.srcDirLayout)

		dstDir, err := ioutil.TempDir(tmpDir, "dst")
		require.NoError(t, err, "Case %d", i)
		writeLayout(t, dstDir, currCase.dstDirLayout)

		modified, err := layout.SyncDir(srcDir, dstDir, currCase.skip)
		require.NoError(t, err, "Case %d", i)

		assert.Equal(t, currCase.wantModified, modified, "Case %d", i)
		if currCase.wantModified {
			// if modification is expected, verify result
			assertLayoutEqual(t, i, currCase.wantDirLayout, dstDir)
		} else {
			// if modification is not expected,
			assertLayoutEqual(t, i, currCase.dstDirLayout, dstDir)
		}
	}

}

func assertLayoutEqual(t *testing.T, caseNum int, want Specs, got string) {
	// verify that paths are correct (catches case where path not in spec exists)
	gotPaths, err := layout.AllPaths(got)
	require.NoError(t, err, "Case %d", caseNum)
	assert.Equal(t, want.AllPaths(), gotPaths)

	// verify that provided directory matches all provided specs
	for _, curr := range want {
		p := path.Join(got, curr.Path)

		fi, err := os.Stat(p)
		assert.NoError(t, err, "Case %d", caseNum)
		assert.Equal(t, curr.IsDir, fi.IsDir(), "Case %d", caseNum)

		if !curr.IsDir {
			content, err := ioutil.ReadFile(p)
			require.NoError(t, err, "Case %d", caseNum)
			assert.Equal(t, curr.Content, string(content), "Case %d", caseNum)
		}
	}
}

func writeLayout(t *testing.T, dir string, specs []Spec) {
	for _, curr := range specs {
		p := path.Join(dir, curr.Path)

		dir := p
		if !curr.IsDir {
			dir = path.Dir(dir)
		}
		err := os.MkdirAll(dir, 0755)
		require.NoError(t, err, "Failed to create directory %v", dir)

		if !curr.IsDir {
			err = ioutil.WriteFile(p, []byte(curr.Content), 0644)
			require.NoError(t, err, "Failed to write file %v", p)
		}
	}
}
