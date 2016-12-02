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

package config_test

import (
	"fmt"
	"io/ioutil"
	"strings"
	"testing"

	"github.com/nmiyake/pkg/dirs"
	"github.com/palantir/pkg/matcher"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/palantir/godel/apps/gunit/config"
)

func TestLoadConfig(t *testing.T) {
	tmpDir, cleanup, err := dirs.TempDir("", "")
	defer cleanup()
	require.NoError(t, err)

	for i, currCase := range []struct {
		yml  string
		json string
		want func() config.Config
	}{
		{
			yml: `
			tags:
			  integration:
			    names:
			      - "integration_tests"
			    paths:
			      - "test"
			exclude:
			  names:
			    - ".*test"
			    - "m?cks"
			  paths:
			    - "vendor"
			`,
			json: `{"exclude":{"names":["gunit"],"paths":["generated_src"]}}`,
			want: func() config.Config {
				includeCfg := matcher.NamesPathsCfg{
					Names: []string{`integration_tests`},
					Paths: []string{`test`},
				}
				excludeCfg := matcher.NamesPathsCfg{
					Names: []string{`.*test`, `m?cks`, `gunit`},
					Paths: []string{`vendor`, `generated_src`},
				}
				return config.Config{
					Tags: map[string]matcher.Matcher{
						"integration": includeCfg.Matcher(),
					},
					Exclude: excludeCfg.Matcher(),
				}
			},
		},
	} {
		path, err := ioutil.TempFile(tmpDir, "")
		require.NoError(t, err, "Case %d", i)
		err = ioutil.WriteFile(path.Name(), []byte(unindent(currCase.yml)), 0644)
		require.NoError(t, err, "Case %d", i)

		got, err := config.Load(path.Name(), currCase.json)
		require.NoError(t, err, "Case %d", i)
		assert.Equal(t, currCase.want(), got, "Case %d", i)

		_, err = config.Load(path.Name(), currCase.json)
		assert.NoError(t, err, "Case %d", i)
	}
}

func TestLoadInvalidConfig(t *testing.T) {
	tmpDir, cleanup, err := dirs.TempDir("", "")
	defer cleanup()
	require.NoError(t, err)

	for i, currCase := range []struct {
		yml       string
		wantError string
	}{
		{
			yml: `
			tags:
			  integration:
			    names:
			      - "integration_tests"
			  foo-bar:
			    paths:
			     - "foo-bar"
			  foo_bar:
			    names:
			      - "foo_bar"
			  "invalid,entry":
			    names:
			      - "invalid"
			  "another bad":
			    names:
			      - "another bad"
			`,
			wantError: "invalid tag names: [another bad invalid,entry]",
		},
	} {
		path, err := ioutil.TempFile(tmpDir, "")
		require.NoError(t, err, "Case %d", i)
		err = ioutil.WriteFile(path.Name(), []byte(unindent(currCase.yml)), 0644)
		require.NoError(t, err, "Case %d", i)

		_, err = config.Load(path.Name(), "")
		require.Error(t, err, fmt.Sprintf("Case %d", i))
		assert.Equal(t, currCase.wantError, err.Error(), "Case %d", i)
	}
}

func unindent(input string) string {
	return strings.Replace(input, "\n\t\t\t", "\n", -1)
}
