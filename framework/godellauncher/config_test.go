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

package godellauncher_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v2"

	"github.com/palantir/godel/framework/godellauncher"
)

func TestMarshalConfig(t *testing.T) {
	cfg := godellauncher.GodelConfig{
		Plugins: godellauncher.PluginsConfig{
			Plugins: []godellauncher.SinglePluginConfig{
				{
					LocatorWithResolverConfig: godellauncher.LocatorWithResolverConfig{
						Locator: godellauncher.LocatorConfig{
							ID: "com.palantir:plugin:1.0.0",
						},
					},
					Assets: []godellauncher.LocatorWithResolverConfig{
						{
							Locator: godellauncher.LocatorConfig{
								ID: "com.palantir:asset:1.0.0",
							},
						},
					},
				},
			},
		},
	}
	got, err := yaml.Marshal(cfg)
	require.NoError(t, err)

	want := `plugins:
  resolvers: []
  plugins:
  - locator:
      id: com.palantir:plugin:1.0.0
      checksums: {}
    resolver: ""
    assets:
    - locator:
        id: com.palantir:asset:1.0.0
        checksums: {}
      resolver: ""
exclude:
  names: []
  paths: []
`
	assert.Equal(t, want, string(got))
}

func TestUnmarshalConfig(t *testing.T) {
	cfgYAML := `
plugins:
  resolvers:
    - foo/repo/{{GroupPath}}/{{Product}}/{{Version}}/{{Product}}-{{OS}}-{{Arch}}-{{Version}}.tgz
  plugins:
    - locator:
        id: "com.palantir:plugin:1.0.0"
      assets:
        - locator:
            id: "com.palantir:asset:1.0.0"
`
	var got godellauncher.GodelConfig
	err := yaml.Unmarshal([]byte(cfgYAML), &got)
	require.NoError(t, err)

	want := godellauncher.GodelConfig{
		Plugins: godellauncher.PluginsConfig{
			DefaultResolvers: []string{
				"foo/repo/{{GroupPath}}/{{Product}}/{{Version}}/{{Product}}-{{OS}}-{{Arch}}-{{Version}}.tgz",
			},
			Plugins: []godellauncher.SinglePluginConfig{
				{
					LocatorWithResolverConfig: godellauncher.LocatorWithResolverConfig{
						Locator: godellauncher.LocatorConfig{
							ID: "com.palantir:plugin:1.0.0",
						},
					},
					Assets: []godellauncher.LocatorWithResolverConfig{
						{
							Locator: godellauncher.LocatorConfig{
								ID: "com.palantir:asset:1.0.0",
							},
						},
					},
				},
			},
		},
	}
	assert.Equal(t, want, got)
}
