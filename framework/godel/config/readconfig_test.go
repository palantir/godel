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
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/nmiyake/pkg/dirs"
	"github.com/palantir/godel/v2/framework/godel/config"
	v0 "github.com/palantir/godel/v2/framework/godel/config/internal/v0"
	"github.com/palantir/pkg/matcher"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
				Plugins: v0.PluginsConfig{
					DefaultResolvers: []string{
						"https://palantir.bintray.com/releases/{{GroupPath}}/{{Product}}/{{Version}}/{{Product}}-{{Version}}-{{OS}}-{{Arch}}.tgz",
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
		inputFile := filepath.Join(testDir, fmt.Sprintf("test_%d.yml", i))
		err := os.WriteFile(inputFile, []byte(tc.ymlInput), 0644)
		require.NoError(t, err, "Case %d")

		gotCfg, err := config.ReadGodelConfigFromFile(inputFile)
		require.NoError(t, err, "Case %d")

		assert.Equal(t, tc.want, gotCfg)
	}
}

func TestReadGodelConfigFromFileDuplicateKeys(t *testing.T) {
	testDir, cleanup, err := dirs.TempDir("", "")
	require.NoError(t, err)
	defer cleanup()

	inputFile := filepath.Join(testDir, "godel.yml")
	err = os.WriteFile(inputFile, []byte(`plugins:
  plugins:
    - locator:
        id: com.palantir.godel-mod-plugin:mod-plugin:1.0.0
        id: com.palantir.godel-mod-plugin:mod-plugin:2.0.0
        checksums:
          darwin-amd64: abc
          linux-amd64: def
          linux-amd64: ghi
exclude:
  names:
    - vendor
exclude:
  names:
    - godel
`), 0644)
	require.NoError(t, err)

	_, err = config.ReadGodelConfigFromFile(inputFile)
	require.Error(t, err)
	assert.Equal(t, fmt.Sprintf(`invalid gödel configuration in %s: configuration contains duplicate keys:
    line 5: key "id" already defined at line 4
    line 9: key "linux-amd64" already defined at line 8
    line 13: key "exclude" already defined at line 10`, inputFile), err.Error())
}

func TestReadGodelConfigFromFileAllowsUnrecognizedKeys(t *testing.T) {
	testDir, cleanup, err := dirs.TempDir("", "")
	require.NoError(t, err)
	defer cleanup()

	inputFile := filepath.Join(testDir, "godel.yml")
	err = os.WriteFile(inputFile, []byte(`unrecognized-key: value
exclude:
  names:
    - vendor
`), 0644)
	require.NoError(t, err)

	gotCfg, err := config.ReadGodelConfigFromFile(inputFile)
	require.NoError(t, err)
	assert.Equal(t, config.GodelConfig(v0.GodelConfig{
		Exclude: matcher.NamesPathsCfg{
			Names: []string{
				"vendor",
			},
		},
	}), gotCfg)
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
		inputFile := filepath.Join(testDir, fmt.Sprintf("test_%d.yml", i))
		err := os.WriteFile(inputFile, []byte(tc.ymlInput), 0644)
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
