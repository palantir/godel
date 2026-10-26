// Copyright 2026 Palantir Technologies, Inc.
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

package gotoolchain_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/palantir/godel/v2/pkg/gotoolchain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFromGoMod(t *testing.T) {
	for _, tc := range []struct {
		name       string
		directives string
		want       string
		wantErr    string
	}{
		{name: "toolchain", directives: "go 1.27.0\ntoolchain go1.27.3\n", want: "go1.27.3"},
		{name: "patch", directives: "go 1.27.2\n", want: "go1.27.2"},
		{name: "language", directives: "go 1.27\n", want: "go1.27.0"},
		{name: "default", directives: "go 1.27\ntoolchain default\n", want: "go1.27.0"},
		{name: "rc", directives: "go 1.28rc1\ntoolchain go1.28rc2\n", want: "go1.28rc2"},
		{name: "beta", directives: "go 1.28beta1\ntoolchain go1.28beta2\n", want: "go1.28beta2"},
		{name: "go prerelease", directives: "go 1.28rc1\n", want: "go1.28rc1"},
		{name: "lower toolchain", directives: "go 1.27.3\ntoolchain go1.27.1\n", want: "go1.27.1"},
		{name: "comments", directives: "go 1.27 // language\ntoolchain go1.27.3 // exact\n", want: "go1.27.3"},
		{name: "pre 1.21", directives: "go 1.20\n", want: "go1.20.0"},
		{name: "missing go", directives: "toolchain go1.27.3\n", wantErr: "must contain a go directive"},
		{name: "custom", directives: "go 1.27.0\ntoolchain go1.27.3-custom\n", wantErr: "unsupported non-standard toolchain"},
		{name: "auto", directives: "go 1.27.0\ntoolchain go1.27.3+auto\n", wantErr: "invalid toolchain"},
		{name: "invalid go", directives: "go invalid\n", wantErr: "parse go.mod"},
		{name: "duplicate", directives: "go 1.27\ngo 1.27.0\n", wantErr: "parse go.mod"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			path := filepath.Join(t.TempDir(), "go.mod")
			require.NoError(t, os.WriteFile(path, []byte("module example.com/test\n"+tc.directives), 0644))
			got, err := gotoolchain.FromGoMod(path)
			if tc.wantErr != "" {
				require.ErrorContains(t, err, tc.wantErr)
				assert.Empty(t, got)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestFromGoModMissingFile(t *testing.T) {
	_, err := gotoolchain.FromGoMod(filepath.Join(t.TempDir(), "go.mod"))
	require.ErrorContains(t, err, "read go.mod")
}
