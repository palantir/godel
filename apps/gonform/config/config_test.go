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
	"strings"
	"testing"

	"github.com/palantir/pkg/matcher"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/palantir/godel/apps/gonform/config"
)

func TestLoad(t *testing.T) {
	for i, currCase := range []struct {
		yml  string
		json string
		want func() config.Gonform
	}{
		{
			yml: `
			formatters:
			  gofmt:
			    args:
			      - "-s"
			exclude:
			  names:
			    - ".*test"
			    - "m?cks"
			  paths:
			    - "vendor"
			`,
			json: `{"exclude":{"names":["gonform"],"paths":["generated_src"]}}`,
			want: func() config.Gonform {
				return config.Gonform{
					Formatters: map[string]config.Formatter{
						"gofmt": {
							Args: []string{
								"-s",
							},
						},
					},
					Exclude: matcher.NamesPathsCfg{
						Names: []string{`.*test`, `m?cks`, `gonform`},
						Paths: []string{`vendor`, `generated_src`},
					},
				}
			},
		},
	} {
		got, err := config.LoadRawConfig(unindent(currCase.yml), currCase.json)
		require.NoError(t, err, "Case %d", i)
		assert.Equal(t, currCase.want(), got, "Case %d", i)
	}
}

func unindent(input string) string {
	return strings.Replace(input, "\n\t\t\t", "\n", -1)
}
