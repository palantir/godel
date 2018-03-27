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

package defaulttasks

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/palantir/godel/framework/godel/config"
)

func testDefaultPluginsConfig() config.PluginsConfig {
	return config.PluginsConfig{
		DefaultResolvers: []string{defaultResolver},
		Plugins: config.ToSinglePluginConfigs([]config.SinglePluginConfig{
			{
				LocatorWithResolverConfig: config.ToLocatorWithResolverConfig(config.LocatorWithResolverConfig{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.test:test-plugin:1.2.3",
					}),
				}),
				Assets: config.ToLocatorWithResolverConfigs([]config.LocatorWithResolverConfig{
					{
						Locator: config.ToLocatorConfig(config.LocatorConfig{
							ID: "com.palantir.test:test-asset-1:2.3.4",
						}),
					},
					{
						Locator: config.ToLocatorConfig(config.LocatorConfig{
							ID: "com.palantir.test:test-asset-2:3.4.5",
						}),
					},
				}),
			},
		}),
	}
}

func TestDefaultTasksPluginsConfig(t *testing.T) {
	original := defaultPluginsConfig
	defer func() {
		defaultPluginsConfig = original
	}()
	defaultPluginsConfig = testDefaultPluginsConfig()

	for i, tc := range []struct {
		name      string
		in        config.DefaultTasksConfig
		wantError string
		want      config.PluginsConfig
	}{
		{
			"empty task param results in default configuration",
			config.DefaultTasksConfig{},
			"",
			config.PluginsConfig{
				DefaultResolvers: []string{
					defaultResolver,
				},
				Plugins: config.ToSinglePluginConfigs([]config.SinglePluginConfig{
					{
						LocatorWithResolverConfig: config.ToLocatorWithResolverConfig(config.LocatorWithResolverConfig{
							Locator: config.ToLocatorConfig(config.LocatorConfig{
								ID: "com.palantir.test:test-plugin:1.2.3",
							}),
						}),
						Assets: config.ToLocatorWithResolverConfigs([]config.LocatorWithResolverConfig{
							{
								Locator: config.ToLocatorConfig(config.LocatorConfig{
									ID: "com.palantir.test:test-asset-1:2.3.4",
								}),
							},
							{
								Locator: config.ToLocatorConfig(config.LocatorConfig{
									ID: "com.palantir.test:test-asset-2:3.4.5",
								}),
							},
						}),
					},
				}),
			},
		},
		{
			"specifying custom resolver overrides resolver",
			config.DefaultTasksConfig{
				Tasks: config.ToTasks(map[string]config.SingleDefaultTaskConfig{
					"com.palantir.test:test-plugin": {
						LocatorWithResolverConfig: config.ToLocatorWithResolverConfig(config.LocatorWithResolverConfig{
							Resolver: "custom-resolver",
						}),
					},
				}),
			},
			"",
			config.PluginsConfig{
				DefaultResolvers: []string{
					defaultResolver,
				},
				Plugins: config.ToSinglePluginConfigs([]config.SinglePluginConfig{
					{
						LocatorWithResolverConfig: config.ToLocatorWithResolverConfig(config.LocatorWithResolverConfig{
							Locator: config.ToLocatorConfig(config.LocatorConfig{
								ID: "com.palantir.test:test-plugin:1.2.3",
							}),
							Resolver: "custom-resolver",
						}),
						Assets: config.ToLocatorWithResolverConfigs([]config.LocatorWithResolverConfig{
							{
								Locator: config.ToLocatorConfig(config.LocatorConfig{
									ID: "com.palantir.test:test-asset-1:2.3.4",
								}),
							},
							{
								Locator: config.ToLocatorConfig(config.LocatorConfig{
									ID: "com.palantir.test:test-asset-2:3.4.5",
								}),
							},
						}),
					},
				}),
			},
		},
		{
			"specifying custom locator overrides locator",
			config.DefaultTasksConfig{
				Tasks: config.ToTasks(map[string]config.SingleDefaultTaskConfig{
					"com.palantir.test:test-plugin": {
						LocatorWithResolverConfig: config.ToLocatorWithResolverConfig(config.LocatorWithResolverConfig{
							Locator: config.ToLocatorConfig(config.LocatorConfig{
								ID: "com.palantir.godel:override:1.2.3",
							}),
						}),
					},
				}),
			},
			"",
			config.PluginsConfig{
				DefaultResolvers: []string{
					defaultResolver,
				},
				Plugins: config.ToSinglePluginConfigs([]config.SinglePluginConfig{
					{
						LocatorWithResolverConfig: config.ToLocatorWithResolverConfig(config.LocatorWithResolverConfig{
							Locator: config.ToLocatorConfig(config.LocatorConfig{
								ID: "com.palantir.godel:override:1.2.3",
							}),
						}),
						Assets: config.ToLocatorWithResolverConfigs([]config.LocatorWithResolverConfig{
							{
								Locator: config.ToLocatorConfig(config.LocatorConfig{
									ID: "com.palantir.test:test-asset-1:2.3.4",
								}),
							},
							{
								Locator: config.ToLocatorConfig(config.LocatorConfig{
									ID: "com.palantir.test:test-asset-2:3.4.5",
								}),
							},
						}),
					},
				}),
			},
		},
		{
			"specifying default resolver appends default resolver",
			config.DefaultTasksConfig{
				DefaultResolvers: []string{
					"default/repo/{{GroupPath}}/{{Product}}/{{Version}}/{{Product}}-{{OS}}-{{Arch}}-{{Version}}.tgz",
				},
				Tasks: config.ToTasks(map[string]config.SingleDefaultTaskConfig{
					"com.palantir.test:test-plugin": {
						LocatorWithResolverConfig: config.ToLocatorWithResolverConfig(config.LocatorWithResolverConfig{
							Locator: config.ToLocatorConfig(config.LocatorConfig{
								ID: "com.palantir.godel:override:1.2.3",
							}),
						}),
					},
				}),
			},
			"",
			config.PluginsConfig{
				DefaultResolvers: []string{
					defaultResolver,
					"default/repo/{{GroupPath}}/{{Product}}/{{Version}}/{{Product}}-{{OS}}-{{Arch}}-{{Version}}.tgz",
				},
				Plugins: config.ToSinglePluginConfigs([]config.SinglePluginConfig{
					{
						LocatorWithResolverConfig: config.ToLocatorWithResolverConfig(config.LocatorWithResolverConfig{
							Locator: config.ToLocatorConfig(config.LocatorConfig{
								ID: "com.palantir.godel:override:1.2.3",
							}),
						}),
						Assets: config.ToLocatorWithResolverConfigs([]config.LocatorWithResolverConfig{
							{
								Locator: config.ToLocatorConfig(config.LocatorConfig{
									ID: "com.palantir.test:test-asset-1:2.3.4",
								}),
							},
							{
								Locator: config.ToLocatorConfig(config.LocatorConfig{
									ID: "com.palantir.test:test-asset-2:3.4.5",
								}),
							},
						}),
					},
				}),
			},
		},
		{
			"specifying custom asset adds only that asset",
			config.DefaultTasksConfig{
				Tasks: config.ToTasks(map[string]config.SingleDefaultTaskConfig{
					"com.palantir.test:test-plugin": {
						Assets: config.ToLocatorWithResolverConfigs([]config.LocatorWithResolverConfig{
							{
								Locator: config.ToLocatorConfig(config.LocatorConfig{
									ID: "com.palantir.godel:custom-asset:1.2.3",
								}),
							},
						}),
					},
				}),
			},
			"",
			config.PluginsConfig{
				DefaultResolvers: []string{
					defaultResolver,
				},
				Plugins: config.ToSinglePluginConfigs([]config.SinglePluginConfig{
					{
						LocatorWithResolverConfig: config.ToLocatorWithResolverConfig(config.LocatorWithResolverConfig{
							Locator: config.ToLocatorConfig(config.LocatorConfig{
								ID: "com.palantir.test:test-plugin:1.2.3",
							}),
						}),
						Assets: config.ToLocatorWithResolverConfigs([]config.LocatorWithResolverConfig{
							{
								Locator: config.ToLocatorConfig(config.LocatorConfig{
									ID: "com.palantir.test:test-asset-1:2.3.4",
								}),
							},
							{
								Locator: config.ToLocatorConfig(config.LocatorConfig{
									ID: "com.palantir.test:test-asset-2:3.4.5",
								}),
							},
							{
								Locator: config.ToLocatorConfig(config.LocatorConfig{
									ID: "com.palantir.godel:custom-asset:1.2.3",
								}),
							},
						}),
					},
				}),
			},
		},
		{
			"setting exclude all and specifying custom asset adds asset to default",
			config.DefaultTasksConfig{
				Tasks: config.ToTasks(map[string]config.SingleDefaultTaskConfig{
					"com.palantir.test:test-plugin": {
						ExcludeAllDefaultAssets: true,
						Assets: config.ToLocatorWithResolverConfigs([]config.LocatorWithResolverConfig{
							{
								Locator: config.ToLocatorConfig(config.LocatorConfig{
									ID: "com.palantir.godel:custom-asset:1.2.3",
								}),
							},
						}),
					},
				}),
			},
			"",
			config.PluginsConfig{
				DefaultResolvers: []string{
					defaultResolver,
				},
				Plugins: config.ToSinglePluginConfigs([]config.SinglePluginConfig{
					{
						LocatorWithResolverConfig: config.ToLocatorWithResolverConfig(config.LocatorWithResolverConfig{
							Locator: config.ToLocatorConfig(config.LocatorConfig{
								ID: "com.palantir.test:test-plugin:1.2.3",
							}),
						}),
						Assets: config.ToLocatorWithResolverConfigs([]config.LocatorWithResolverConfig{
							{
								Locator: config.ToLocatorConfig(config.LocatorConfig{
									ID: "com.palantir.godel:custom-asset:1.2.3",
								}),
							},
						}),
					},
				}),
			},
		},
		{
			"specifying default asset with exclude and custom asset adds asset",
			config.DefaultTasksConfig{
				Tasks: config.ToTasks(map[string]config.SingleDefaultTaskConfig{
					"com.palantir.test:test-plugin": {
						DefaultAssetsToExclude: []string{
							"com.palantir.test:test-asset-2",
						},
						Assets: config.ToLocatorWithResolverConfigs([]config.LocatorWithResolverConfig{
							{
								Locator: config.ToLocatorConfig(config.LocatorConfig{
									ID: "com.palantir.godel:custom-asset:1.2.3",
								}),
							},
						}),
					},
				}),
			},
			"",
			config.PluginsConfig{
				DefaultResolvers: []string{
					defaultResolver,
				},
				Plugins: config.ToSinglePluginConfigs([]config.SinglePluginConfig{
					{
						LocatorWithResolverConfig: config.ToLocatorWithResolverConfig(config.LocatorWithResolverConfig{
							Locator: config.ToLocatorConfig(config.LocatorConfig{
								ID: "com.palantir.test:test-plugin:1.2.3",
							}),
						}),
						Assets: config.ToLocatorWithResolverConfigs([]config.LocatorWithResolverConfig{
							{
								Locator: config.ToLocatorConfig(config.LocatorConfig{
									ID: "com.palantir.test:test-asset-1:2.3.4",
								}),
							},
							{
								Locator: config.ToLocatorConfig(config.LocatorConfig{
									ID: "com.palantir.godel:custom-asset:1.2.3",
								}),
							},
						}),
					},
				}),
			},
		},
		{
			"specifying invalid key results in error",
			config.DefaultTasksConfig{
				Tasks: config.ToTasks(map[string]config.SingleDefaultTaskConfig{
					"com.palantir.test:test": {
						LocatorWithResolverConfig: config.ToLocatorWithResolverConfig(config.LocatorWithResolverConfig{
							Locator: config.ToLocatorConfig(config.LocatorConfig{
								ID: "com.palantir.godel:override:1.2.3",
							}),
						}),
					},
				}),
			},
			`default-task key(s) specified but are not valid: [com.palantir.test:test]. Valid values: [com.palantir.test:test-plugin]`,
			config.PluginsConfig{},
		},
	} {
		got, err := PluginsConfig(tc.in)
		if tc.wantError == "" {
			require.NoError(t, err)
			assert.Equal(t, tc.want, got, "Case %d: %s", i, tc.name)
		} else {
			assert.EqualError(t, err, tc.wantError, "Case %d: %s", i, tc.name)
		}
	}
}
