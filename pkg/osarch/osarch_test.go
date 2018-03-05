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

package osarch_test

import (
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/palantir/godel/pkg/osarch"
)

func TestOSArchNew(t *testing.T) {
	for i, currCase := range []struct {
		input     string
		want      osarch.OSArch
		wantError bool
	}{
		{input: "darwin-amd64", want: osarch.OSArch{OS: "darwin", Arch: "amd64"}},
		{input: "foo-bar", want: osarch.OSArch{OS: "foo", Arch: "bar"}},
		{input: "foo-bar-baz", wantError: true},
		{input: "foo", wantError: true},
		{input: "foo-", wantError: true},
		{input: "-bar", wantError: true},
		{input: "foo-b@r", wantError: true},
		{input: "f?o-bar", wantError: true},
		{input: "", wantError: true},
	} {
		got, err := osarch.New(currCase.input)
		if currCase.wantError {
			assert.EqualError(t, err, "not a valid OSArch value: "+currCase.input, "Case %d", i)
		} else {
			require.NoError(t, err, "Case %d", i)
			assert.Equal(t, currCase.want, got, "Case %d", i)
		}
	}
}

func TestOSArchCurrent(t *testing.T) {
	want := osarch.OSArch{OS: runtime.GOOS, Arch: runtime.GOARCH}
	assert.Equal(t, want, osarch.Current())
}
