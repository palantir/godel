// Copyright 2020 Palantir Technologies, Inc.
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

package installupdate

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Verifies that getLastLine properly splits on both \r and \n.
func TestGetLastLine(t *testing.T) {
	for i, tc := range []struct {
		in   string
		want string
	}{
		{
			"",
			"",
		},
		{
			`first line
second line`,
			"second line",
		},
		{
			fmt.Sprint("first line", "\r", "second line"),
			"second line",
		},
	} {
		t.Run(fmt.Sprintf("Case %d", i), func(t *testing.T) {
			got, err := getLastLine(tc.in)
			require.NoError(t, err)
			assert.Equal(t, tc.want, got)
		})
	}
}

// Verifies that getGodelVersion includes the wrapper's stderr output in the returned error.
func TestGetGodelVersionErrorIncludesStderr(t *testing.T) {
	projectDir := t.TempDir()
	script := "#!/usr/bin/env bash\necho 'SHA-256 checksum of downloaded file did not match expected value' >&2\nexit 1\n"
	require.NoError(t, os.WriteFile(filepath.Join(projectDir, "godelw"), []byte(script), 0755))
	_, err := getGodelVersion(projectDir)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "SHA-256 checksum of downloaded file did not match expected value")
}
