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

package versionedconfig_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/palantir/godel/pkg/versionedconfig"
)

func TestConfigNotSupported(t *testing.T) {
	for i, tc := range []struct {
		name      string
		cfgName   string
		yaml      string
		wantError string
		wantBytes []byte
	}{
		{
			"empty configuration valid",
			"test",
			``,
			"",
			[]byte{},
		},
		{
			"comment-only configuration valid",
			"test",
			`
# only a comment
`,
			"",
			[]byte(`
# only a comment
`),
		},
		{
			"non-empty configuration invalid",
			"test",
			`key: value
`,
			"test does not currently support configuration",
			nil,
		},
	} {
		got, err := versionedconfig.ConfigNotSupported(tc.cfgName, []byte(tc.yaml))
		if tc.wantError != "" {
			require.Error(t, err, "Case %d: %s", i, tc.name)
			assert.EqualError(t, err, tc.wantError, "Case %d: %s", i, tc.name)
		} else {
			require.NoError(t, err, "Case %d: %s", i, tc.name)
		}
		assert.Equal(t, tc.wantBytes, got, "Case %d: %s", i, tc.name)
	}
}
