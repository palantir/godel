// Copyright 2016 Palantir Technologies. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package matcher_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v2"

	"github.com/palantir/pkg/matcher"
)

func TestNamesPathsCfgMatcher(t *testing.T) {
	for i, currCase := range []struct {
		cfg         matcher.NamesPathsCfg
		wantMatch   []string
		wantNoMatch []string
	}{
		{
			cfg: matcher.NamesPathsCfg{
				Names: []string{
					"foo",
				},
				Paths: []string{
					"bar",
				},
			},
			wantMatch: []string{
				"foo",
				"baz/bing/foo",
				"bar",
				"bar/baz",
			},
			wantNoMatch: []string{
				"baz",
				"fooz/bar",
			},
		},
	} {
		m := currCase.cfg.Matcher()
		for _, curr := range currCase.wantMatch {
			assert.True(t, m.Match(curr), "expected %q to match matcher in case %d", curr, i)
		}
		for _, curr := range currCase.wantNoMatch {
			assert.False(t, m.Match(curr), "expected %q to not match matcher in case %d", curr, i)
		}
	}
}

func TestLoadNamesPathsCfg(t *testing.T) {
	for i, currCase := range []struct {
		yml     string
		json    string
		wantCfg matcher.NamesPathsCfg
	}{
		{
			yml: `
names:
  - foo
  - bar
paths:
  - baz/bing
`,
			json: `{"names":["foo","bar"],"paths":["baz/bing"]}`,
			wantCfg: matcher.NamesPathsCfg{
				Names: []string{
					"foo",
					"bar",
				},
				Paths: []string{
					"baz/bing",
				},
			},
		},
	} {
		var gotCfgFromJSON matcher.NamesPathsCfg
		err := json.Unmarshal([]byte(currCase.json), &gotCfgFromJSON)
		require.NoError(t, err, "Case %d", i)
		assert.Equal(t, currCase.wantCfg, gotCfgFromJSON, "Case %d", i)

		var gotCfgFromYML matcher.NamesPathsCfg
		err = yaml.Unmarshal([]byte(currCase.yml), &gotCfgFromYML)
		require.NoError(t, err, "Case %d", i)
		assert.Equal(t, currCase.wantCfg, gotCfgFromYML, "Case %d", i)
	}
}

func TestNamesPathsWithExcludeCfgMatcher(t *testing.T) {
	for i, currCase := range []struct {
		cfg         matcher.NamesPathsWithExcludeCfg
		wantMatch   []string
		wantNoMatch []string
	}{
		{
			cfg: matcher.NamesPathsWithExcludeCfg{
				NamesPathsCfg: matcher.NamesPathsCfg{
					Names: []string{
						"foo",
					},
					Paths: []string{
						"bar",
					},
				},
				Exclude: matcher.NamesPathsCfg{
					Paths: []string{
						"bar/baz",
					},
				},
			},
			wantMatch: []string{
				"foo",
				"baz/bing/foo",
				"bar",
			},
			wantNoMatch: []string{
				"bar/baz",
				"baz",
				"fooz/bar",
			},
		},
	} {
		m := currCase.cfg.Matcher()
		for _, curr := range currCase.wantMatch {
			assert.True(t, m.Match(curr), "expected %q to match matcher in case %d", curr, i)
		}
		for _, curr := range currCase.wantNoMatch {
			assert.False(t, m.Match(curr), "expected %q to not match matcher in case %d", curr, i)
		}
	}
}

func TestLoadNamesPathsWithExcludeCfg(t *testing.T) {
	for i, currCase := range []struct {
		yml     string
		json    string
		wantCfg matcher.NamesPathsWithExcludeCfg
	}{
		{
			yml: `
names:
  - foo
  - bar
paths:
  - baz/bing
  - abc/def
exclude:
  names:
    - foo
  paths:
    - baz/bing
`,
			json: `{"names":["foo","bar"],"paths":["baz/bing","abc/def"],"exclude":{"names":["foo"],"paths":["baz/bing"]}}`,
			wantCfg: matcher.NamesPathsWithExcludeCfg{
				NamesPathsCfg: matcher.NamesPathsCfg{
					Names: []string{
						"foo",
						"bar",
					},
					Paths: []string{
						"baz/bing",
						"abc/def",
					},
				},
				Exclude: matcher.NamesPathsCfg{
					Names: []string{
						"foo",
					},
					Paths: []string{
						"baz/bing",
					},
				},
			},
		},
	} {
		var gotCfgFromJSON matcher.NamesPathsWithExcludeCfg
		err := json.Unmarshal([]byte(currCase.json), &gotCfgFromJSON)
		require.NoError(t, err, "Case %d", i)
		assert.Equal(t, currCase.wantCfg, gotCfgFromJSON, "Case %d", i)

		var gotCfg matcher.NamesPathsWithExcludeCfg
		err = yaml.Unmarshal([]byte(currCase.yml), &gotCfg)
		require.NoError(t, err, "Case %d", i)
		assert.Equal(t, currCase.wantCfg, gotCfg, "Case %d", i)
	}
}
