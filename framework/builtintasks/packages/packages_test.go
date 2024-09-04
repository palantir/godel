// Copyright 2024 Palantir Technologies, Inc.
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

package packages

import (
	"testing"

	"github.com/nmiyake/pkg/gofiles"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestListNoExcluders(t *testing.T) {
	for _, tc := range []struct {
		name  string
		files []gofiles.GoFileSpec
		want  []string
	}{
		{
			name: "Basic project with single directory",
			files: []gofiles.GoFileSpec{
				{
					RelPath: "go.mod",
					Src:     "module github.com/godel/testmodule\n",
				},
				{
					RelPath: "src.go",
					Src:     "package foo\n",
				},
			},
			want: []string{
				"./.",
			},
		},
		{
			name: "Basic project with multiple directories",
			files: []gofiles.GoFileSpec{
				{
					RelPath: "go.mod",
					Src:     "module github.com/godel/testmodule\n",
				},
				{
					RelPath: "src.go",
					Src:     "package foo\n",
				},
				{
					RelPath: "bar/bar.go",
					Src:     "package bar\n",
				},
			},
			want: []string{
				"./.",
				"./bar",
			},
		},
		{
			name: "Multi-module project excludes submodules",
			files: []gofiles.GoFileSpec{
				{
					RelPath: "go.mod",
					Src:     "module github.com/godel/testmodule\n",
				},
				{
					RelPath: "src.go",
					Src:     "package foo\n",
				},
				{
					RelPath: "bar/go.mod",
					Src:     "module github.com/godel/bar\n",
				},
				{
					RelPath: "bar/bar.go",
					Src:     "package bar\n",
				},
			},
			want: []string{
				"./.",
			},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			rootDir := t.TempDir()
			_, err := gofiles.Write(rootDir, tc.files)
			require.NoError(t, err)

			matchedPkgs, err := List(nil, rootDir)
			require.NoError(t, err)

			assert.Equal(t, tc.want, matchedPkgs)
		})
	}
}
