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

	"github.com/palantir/godel/apps/okgo/config"
)

func TestLoadRawConfig(t *testing.T) {
	for i, currCase := range []struct {
		yml  string
		json string
		want config.OKGo
	}{
		{
			yml: `
			checks:
			  errcheck:
			    args:
			      - "-ignore"
			      - "github.com/seelog:(Info|Warn|Error|Critical)f?"
			    filters:
			      - type: "message"
			        value: "\\w+"
			exclude:
			  names:
			    - "m?cks"
			  paths:
			    - "vendor"
			`,
			json: `{"exclude":{"names":["exclude.*"],"paths":["generated_src"]}}`,
			want: config.OKGo{
				Checks: map[string]config.Checker{
					"errcheck": {
						Args: []string{
							"-ignore",
							"github.com/seelog:(Info|Warn|Error|Critical)f?",
						},
						Filters: []config.Filter{
							{
								Type:  "message",
								Value: `\w+`,
							},
						},
					},
				},
				Exclude: matcher.NamesPathsCfg{
					Names: []string{
						"m?cks",
						"exclude.*",
					},
					Paths: []string{
						"vendor",
						"generated_src",
					},
				},
			},
		},
		{
			yml: `
			checks:
			  errcheck:
			    filters:
			      - value: "\\w+"
			`,
			want: config.OKGo{
				Checks: map[string]config.Checker{
					"errcheck": {
						Filters: []config.Filter{
							{
								Value: `\w+`,
							},
						},
					},
				},
			},
		},
	} {
		got, err := config.LoadRawConfig(unindent(currCase.yml), currCase.json)
		assert.Equal(t, currCase.want, got, "Case %d", i)
		require.NoError(t, err, "Case %d", i)
	}
}

func TestLoadBadConfig(t *testing.T) {
	tmpDir, cleanup, err := dirs.TempDir("", "")
	defer cleanup()
	require.NoError(t, err)

	for i, currCase := range []struct {
		yml       string
		json      string
		wantError string
	}{
		{
			yml: `
			checks:
			  errcheck:
			    filters:
			      - type: "unknown"
			`,
			wantError: `failed to parse filter: {unknown }: unknown filter type: unknown`,
		},
		{
			json:      `[}`,
			wantError: `failed to parse JSON [}: invalid character '}' looking for beginning of value`,
		},
	} {
		var path string
		if currCase.yml != "" {
			tmpFile, err := ioutil.TempFile(tmpDir, "")
			require.NoError(t, err, "Case %d", i)
			err = ioutil.WriteFile(tmpFile.Name(), []byte(unindent(currCase.yml)), 0644)
			require.NoError(t, err, "Case %d", i)
			path = tmpFile.Name()
		}

		_, err = config.Load(path, currCase.json)
		require.Error(t, err, fmt.Sprintf("Case %d", i))
		require.EqualError(t, err, currCase.wantError, "Case %d", i)
	}
}

func unindent(input string) string {
	return strings.Replace(input, "\n\t\t\t", "\n", -1)
}
