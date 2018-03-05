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

package artifactresolver

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"path"
	"path/filepath"
	"testing"

	"github.com/nmiyake/pkg/dirs"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/palantir/godel/pkg/osarch"
)

func TestResolverLocal(t *testing.T) {
	tmpDir, cleanup, err := dirs.TempDir("", "")
	require.NoError(t, err)
	defer cleanup()

	const content = "file content"
	srcFile := path.Join(tmpDir, "srcFile")
	err = ioutil.WriteFile(srcFile, []byte(content), 0644)
	require.NoError(t, err)

	srcFileAbs, err := filepath.Abs(srcFile)
	require.NoError(t, err)

	r, err := NewTemplateResolver(srcFileAbs)
	require.NoError(t, err)

	dstFile := path.Join(tmpDir, "dstFile")
	dstFileAbs, err := filepath.Abs(dstFile)
	require.NoError(t, err)

	err = r.Resolve(LocatorParam{}, osarch.OSArch{}, dstFileAbs, ioutil.Discard)
	require.NoError(t, err)
	bytes, err := ioutil.ReadFile(dstFileAbs)
	require.NoError(t, err)
	assert.Equal(t, content, string(bytes))
}

func TestResolverURL(t *testing.T) {
	const content = "file content\n"
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, content)
	}))
	defer ts.Close()

	tmpDir, cleanup, err := dirs.TempDir("", "")
	require.NoError(t, err)
	defer cleanup()

	r, err := NewTemplateResolver(ts.URL)
	require.NoError(t, err)

	dstFile := path.Join(tmpDir, "dstFile")
	dstFileAbs, err := filepath.Abs(dstFile)
	require.NoError(t, err)

	err = r.Resolve(LocatorParam{}, osarch.OSArch{}, dstFileAbs, ioutil.Discard)
	require.NoError(t, err)
	bytes, err := ioutil.ReadFile(dstFileAbs)
	require.NoError(t, err)
	assert.Equal(t, content, string(bytes))
}

func TestRenderResolve(t *testing.T) {
	tmpDir, cleanup, err := dirs.TempDir("", "")
	require.NoError(t, err)
	defer cleanup()
	dstFile := path.Join(tmpDir, "dstFile")

	var got string
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		got = r.URL.String()
	}))
	defer ts.Close()

	for i, tc := range []struct {
		name     string
		template string
		locator  LocatorParam
		osArch   osarch.OSArch
		want     string
	}{
		{
			"literal template",
			ts.URL + "/foo",
			LocatorParam{},
			osarch.OSArch{},
			"/foo",
		},
		{
			"template with values",
			ts.URL + "/foo/{{Group}}/{{Product}}-{{OS}}-{{Arch}}-{{Version}}",
			LocatorParam{
				Locator: Locator{
					Group:   "Group",
					Product: "Product",
					Version: "Version",
				},
			},
			osarch.OSArch{
				OS:   "darwin",
				Arch: "amd64",
			},
			"/foo/Group/Product-darwin-amd64-Version",
		},
		{
			"group path",
			ts.URL + "/foo/{{GroupPath}}/{{Product}}-{{OS}}-{{Arch}}-{{Version}}",
			LocatorParam{
				Locator: Locator{
					Group:   "a.b.c",
					Product: "Product",
					Version: "Version",
				},
			},
			osarch.OSArch{
				OS:   "darwin",
				Arch: "amd64",
			},
			"/foo/a/b/c/Product-darwin-amd64-Version",
		},
	} {
		r, err := NewTemplateResolver(tc.template)
		require.NoError(t, err, "Case %d: %s", i, tc.name)
		buf := &bytes.Buffer{}
		err = r.Resolve(tc.locator, tc.osArch, dstFile, buf)
		require.NoError(t, err, "Case %d: %s", i, tc.name)
		assert.Equal(t, tc.want, got, "Case %d: %s", i, tc.name)
		fmt.Println(buf.String())
	}
}
