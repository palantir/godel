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

package printproducts_test

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"path"
	"testing"

	"github.com/nmiyake/pkg/dirs"
	"github.com/nmiyake/pkg/gofiles"
	"github.com/palantir/pkg/gittest"
	"github.com/palantir/pkg/matcher"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/palantir/distgo/dister"
	"github.com/palantir/distgo/distgo"
	"github.com/palantir/distgo/distgo/printproducts"
	"github.com/palantir/distgo/dockerbuilder"
)

func TestProducts(t *testing.T) {
	rootDir, cleanup, err := dirs.TempDir("", "")
	require.NoError(t, err)
	defer cleanup()

	for i, tc := range []struct {
		name            string
		projectConfig   distgo.ProjectConfig
		setupProjectDir func(projectDir string)
		want            func(projectDir string) string
	}{
		{
			"prints products defined in param",
			distgo.ProjectConfig{
				Products: map[distgo.ProductID]distgo.ProductConfig{
					"foo": {},
					"bar": {},
				},
			},
			func(projectDir string) {},
			func(projectDir string) string {
				return `bar
foo
`
			},
		},
		{
			"if param is empty, prints main packages",
			distgo.ProjectConfig{},
			func(projectDir string) {
				_, err := gofiles.Write(projectDir, []gofiles.GoFileSpec{
					{
						RelPath: "main.go",
						Src:     `package main`,
					},
					{
						RelPath: "bar/bar.go",
						Src:     `package bar`,
					},
					{
						RelPath: "foo/foo.go",
						Src:     `package main`,
					},
				})
				require.NoError(t, err)
			},
			func(projectDir string) string {
				return fmt.Sprintf(`%s
foo
`, path.Base(projectDir))
			},
		},
		{
			"if param is empty, prints main packages and uses exclude param",
			distgo.ProjectConfig{
				Exclude: matcher.NamesPathsCfg{
					Paths: []string{
						"foo",
					},
				},
			},
			func(projectDir string) {
				_, err := gofiles.Write(projectDir, []gofiles.GoFileSpec{
					{
						RelPath: "main.go",
						Src:     `package main`,
					},
					{
						RelPath: "bar/bar.go",
						Src:     `package bar`,
					},
					{
						RelPath: "foo/foo.go",
						Src:     `package main`,
					},
				})
				require.NoError(t, err)
			},
			func(projectDir string) string {
				return fmt.Sprintf(`%s
`, path.Base(projectDir))
			},
		},
	} {
		projectDir, err := ioutil.TempDir(rootDir, "")
		require.NoError(t, err, "Case %d: %s", i, tc.name)

		gittest.InitGitDir(t, projectDir)
		tc.setupProjectDir(projectDir)

		disterFactory, err := dister.NewDisterFactory()
		require.NoError(t, err, "Case %d: %s", i, tc.name)
		defaultDisterCfg, err := dister.DefaultConfig()
		require.NoError(t, err, "Case %d: %s", i, tc.name)
		dockerBuilderFactory, err := dockerbuilder.NewDockerBuilderFactory()
		require.NoError(t, err, "Case %d: %s", i, tc.name)

		projectParam, err := tc.projectConfig.ToParam(projectDir, disterFactory, defaultDisterCfg, dockerBuilderFactory)
		require.NoError(t, err, "Case %d: %s", i, tc.name)

		buf := &bytes.Buffer{}
		err = printproducts.Run(projectParam, buf)
		require.NoError(t, err, "Case %d: %s", i, tc.name)
		assert.Equal(t, tc.want(projectDir), buf.String(), "Case %d: %s", i, tc.name)
	}
}
