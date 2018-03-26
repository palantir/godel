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
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v2"

	"github.com/palantir/godel/framework/godel/config"
)

func TestMarshalConfig(t *testing.T) {
	cfg := config.GodelConfig{
		TasksConfig: config.ToTasksConfig(config.TasksConfig{
			Plugins: config.ToPluginsConfig(config.PluginsConfig{
				Plugins: config.ToSinglePluginConfigs([]config.SinglePluginConfig{
					{
						LocatorWithResolverConfig: config.ToLocatorWithResolverConfig(config.LocatorWithResolverConfig{
							Locator: config.ToLocatorConfig(config.LocatorConfig{
								ID: "com.palantir:plugin:1.0.0",
							}),
						}),
						Assets: config.ToLocatorWithResolverConfigs([]config.LocatorWithResolverConfig{
							{
								Locator: config.ToLocatorConfig(config.LocatorConfig{
									ID: "com.palantir:asset:1.0.0",
								}),
							},
						}),
					},
				}),
			}),
		}),
	}
	got, err := yaml.Marshal(cfg)
	require.NoError(t, err)

	want := `version: ""
tasks-config-providers:
  resolvers: []
  providers: []
default-tasks:
  resolvers: []
  tasks: {}
plugins:
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
	var got config.GodelConfig
	err := yaml.Unmarshal([]byte(cfgYAML), &got)
	require.NoError(t, err)

	want := config.GodelConfig{
		TasksConfig: config.ToTasksConfig(config.TasksConfig{
			Plugins: config.ToPluginsConfig(config.PluginsConfig{
				DefaultResolvers: []string{
					"foo/repo/{{GroupPath}}/{{Product}}/{{Version}}/{{Product}}-{{OS}}-{{Arch}}-{{Version}}.tgz",
				},
				Plugins: config.ToSinglePluginConfigs([]config.SinglePluginConfig{
					{
						LocatorWithResolverConfig: config.ToLocatorWithResolverConfig(config.LocatorWithResolverConfig{
							Locator: config.ToLocatorConfig(config.LocatorConfig{
								ID: "com.palantir:plugin:1.0.0",
							}),
						}),
						Assets: config.ToLocatorWithResolverConfigs([]config.LocatorWithResolverConfig{
							{
								Locator: config.ToLocatorConfig(config.LocatorConfig{
									ID: "com.palantir:asset:1.0.0",
								}),
							},
						}),
					},
				}),
			}),
		}),
	}
	assert.Equal(t, want, got)
}

func TestUnmarshalConfigWithDefaults(t *testing.T) {
	cfgYAML := `
default-tasks:
  resolvers:
    - default/repo/{{GroupPath}}/{{Product}}/{{Version}}/{{Product}}-{{OS}}-{{Arch}}-{{Version}}.tgz
  tasks:
    com.palantir.godel:format:
      exclude-default-assets:
        - com.palantir.godel:foo-asset
        - com.palantir.godel:bar-asset
      assets:
        - locator:
            id: "com.palantir.godel:bar-asset:1.0.0"
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
	var got config.GodelConfig
	err := yaml.Unmarshal([]byte(cfgYAML), &got)
	require.NoError(t, err)

	want := config.GodelConfig{
		TasksConfig: config.ToTasksConfig(config.TasksConfig{
			DefaultTasks: config.ToDefaultTasksConfig(config.DefaultTasksConfig{
				DefaultResolvers: []string{
					"default/repo/{{GroupPath}}/{{Product}}/{{Version}}/{{Product}}-{{OS}}-{{Arch}}-{{Version}}.tgz",
				},
				Tasks: config.ToTasks(map[string]config.SingleDefaultTaskConfig{
					"com.palantir.godel:format": {
						DefaultAssetsToExclude: []string{
							"com.palantir.godel:foo-asset",
							"com.palantir.godel:bar-asset",
						},
						Assets: config.ToLocatorWithResolverConfigs([]config.LocatorWithResolverConfig{
							{
								Locator: config.ToLocatorConfig(config.LocatorConfig{
									ID: "com.palantir.godel:bar-asset:1.0.0",
								}),
							},
						}),
					},
				}),
			}),
			Plugins: config.ToPluginsConfig(config.PluginsConfig{
				DefaultResolvers: []string{
					"foo/repo/{{GroupPath}}/{{Product}}/{{Version}}/{{Product}}-{{OS}}-{{Arch}}-{{Version}}.tgz",
				},
				Plugins: config.ToSinglePluginConfigs([]config.SinglePluginConfig{
					{
						LocatorWithResolverConfig: config.ToLocatorWithResolverConfig(config.LocatorWithResolverConfig{
							Locator: config.ToLocatorConfig(config.LocatorConfig{
								ID: "com.palantir:plugin:1.0.0",
							}),
						}),
						Assets: config.ToLocatorWithResolverConfigs([]config.LocatorWithResolverConfig{
							{
								Locator: config.ToLocatorConfig(config.LocatorConfig{
									ID: "com.palantir:asset:1.0.0",
								}),
							},
						}),
					},
				}),
			}),
		}),
	}
	assert.Equal(t, want, got)
}

func TestPluginsConfig_ToParam(t *testing.T) {
	cfgContent := `
resolvers:
  - https://localhost:8080/repo/{{GroupPath}}/{{Product}}/{{Version}}/{{Product}}-{{OS}}-{{Arch}}-{{Version}}.tgz
plugins:
  - locator:
      id: "com.palantir:tester:1.0.0"
      checksums:
        darwin-amd64: d22c0ac9d3b65ebe5b830c1324f3d43e777ebc085c580af7c39fb1e5e3c909a7
`
	var cfg config.PluginsConfig
	err := yaml.Unmarshal([]byte(cfgContent), &cfg)
	require.NoError(t, err)
	_, err = cfg.ToParam()
	require.NoError(t, err)
}

func TestPluginsConfig_ToParam_InvalidLocator(t *testing.T) {
	cfgContent := `
plugins:
  - locator:
      id: "tester:1.0.0"
`
	var cfg config.PluginsConfig
	err := yaml.Unmarshal([]byte(cfgContent), &cfg)
	require.NoError(t, err)
	_, err = cfg.ToParam()
	assert.EqualError(t, err, `invalid locator: locator ID must consist of 3 colon-delimited components ([group]:[product]:[version]), but had 2: "tester:1.0.0"`)
}
