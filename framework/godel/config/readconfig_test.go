// Copyright 2019 Palantir Technologies, Inc.
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

package config_test

import (
	"fmt"
	"io/ioutil"
	"path"
	"testing"
	"time"

	"github.com/nmiyake/pkg/dirs"
	"github.com/palantir/pkg/matcher"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/palantir/godel/v2/framework/godel/config"
	v0 "github.com/palantir/godel/v2/framework/godel/config/internal/v0"
)

func TestReadGodelConfigFromFile(t *testing.T) {
	testDir, cleanup, err := dirs.TempDir("", "")
	require.NoError(t, err)
	defer cleanup()

	for i, tc := range []struct {
		ymlInput string
		want     config.GodelConfig
	}{
		{
			ymlInput: `
plugins:
  resolvers:
    - "https://palantir.bintray.com/releases/{{GroupPath}}/{{Product}}/{{Version}}/{{Product}}-{{Version}}-{{OS}}-{{Arch}}.tgz"
exclude:
  names:
    - "vendor"
  paths:
    - "godel"
`,
			want: config.GodelConfig(v0.GodelConfig{
				TasksConfig: v0.TasksConfig{
					Plugins: v0.PluginsConfig{
						DefaultResolvers: []string{
							"https://palantir.bintray.com/releases/{{GroupPath}}/{{Product}}/{{Version}}/{{Product}}-{{Version}}-{{OS}}-{{Arch}}.tgz",
						},
					},
				},
				Exclude: matcher.NamesPathsCfg{
					Names: []string{
						"vendor",
					},
					Paths: []string{
						"godel",
					},
				},
			}),
		},
	} {
		inputFile := path.Join(testDir, fmt.Sprintf("test_%d.yml", i))
		err := ioutil.WriteFile(inputFile, []byte(tc.ymlInput), 0644)
		require.NoError(t, err, "Case %d")

		gotCfg, err := config.ReadGodelConfigFromFile(inputFile)
		require.NoError(t, err, "Case %d")

		assert.Equal(t, tc.want, gotCfg)
	}
}

func TestReadGodelConfigFromFileError(t *testing.T) {
	nonexistentFile := fmt.Sprintf("TestReadGodelConfigFromFileError-%d.yml", time.Now().Unix())

	wantCfg := config.GodelConfig{}
	gotCfg, err := config.ReadGodelConfigFromFile(nonexistentFile)
	require.NoError(t, err)

	assert.Equal(t, wantCfg, gotCfg)
}

func TestReadGodelConfigExcludesFromFile(t *testing.T) {
	testDir, cleanup, err := dirs.TempDir("", "")
	require.NoError(t, err)
	defer cleanup()

	for i, tc := range []struct {
		ymlInput string
		want     matcher.NamesPathsCfg
	}{
		{
			ymlInput: `
invalid-top-level-key: value
plugins:
  resolvers:
    - "https://palantir.bintray.com/releases/{{GroupPath}}/{{Product}}/{{Version}}/{{Product}}-{{Version}}-{{OS}}-{{Arch}}.tgz"
exclude:
  names:
    - "vendor"
  paths:
    - "godel"
`,
			want: matcher.NamesPathsCfg{
				Names: []string{
					"vendor",
				},
				Paths: []string{
					"godel",
				},
			},
		},
	} {
		inputFile := path.Join(testDir, fmt.Sprintf("test_%d.yml", i))
		err := ioutil.WriteFile(inputFile, []byte(tc.ymlInput), 0644)
		require.NoError(t, err, "Case %d")

		gotExcludes, err := config.ReadGodelConfigExcludesFromFile(inputFile)
		require.NoError(t, err, "Case %d")

		assert.Equal(t, tc.want, gotExcludes)
	}
}

func TestReadGodelConfigExcludesFromFileError(t *testing.T) {
	nonexistentFile := fmt.Sprintf("TestReadGodelConfigExcludesFromFileError-%d.yml", time.Now().Unix())

	wantCfg := matcher.NamesPathsCfg{}
	gotCfg, err := config.ReadGodelConfigExcludesFromFile(nonexistentFile)
	require.NoError(t, err)

	assert.Equal(t, wantCfg, gotCfg)
}
