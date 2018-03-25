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
	"gopkg.in/yaml.v2"

	"github.com/palantir/godel/pkg/versionedconfig"
)

func TestTrimLegacyPrefix(t *testing.T) {
	type Foo struct {
		versionedconfig.ConfigWithLegacy `yaml:",inline"`
		Foo                              string `yaml:"foo"`
	}
	fooBytes, err := yaml.Marshal(Foo{
		ConfigWithLegacy: versionedconfig.ConfigWithLegacy{
			Legacy: true,
		},
		Foo: "foo-value",
	})
	require.NoError(t, err)

	out, trimmed := versionedconfig.TrimLegacyPrefix(fooBytes)
	assert.True(t, trimmed, "expected prefix to be trimmed")
	assert.Equal(t, `foo: foo-value
`, string(out))
}

func TestTrimLegacyPrefixDoesNotTrimSuffix(t *testing.T) {
	type Foo struct {
		Foo                              string `yaml:"foo"`
		versionedconfig.ConfigWithLegacy `yaml:",inline"`
	}
	fooBytes, err := yaml.Marshal(Foo{
		ConfigWithLegacy: versionedconfig.ConfigWithLegacy{
			Legacy: true,
		},
		Foo: "foo-value",
	})
	require.NoError(t, err)

	out, trimmed := versionedconfig.TrimLegacyPrefix(fooBytes)
	assert.False(t, trimmed, "expected prefix to be trimmed")
	assert.Equal(t, `foo: foo-value
legacy-config: true
`, string(out))
}
