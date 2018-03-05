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

package godellauncher

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/palantir/godel/framework/artifactresolver"
)

func testDefaultPluginsConfig() PluginsConfig {
	return PluginsConfig{
		DefaultResolvers: []string{defaultResolver},
		Plugins: []SinglePluginConfig{
			{
				LocatorWithResolverConfig: artifactresolver.LocatorWithResolverConfig{
					Locator: artifactresolver.LocatorConfig{
						ID: "com.palantir.test:test-plugin:1.2.3",
					},
				},
				Assets: []artifactresolver.LocatorWithResolverConfig{
					{
						Locator: artifactresolver.LocatorConfig{
							ID: "com.palantir.test:test-asset-1:2.3.4",
						},
					},
					{
						Locator: artifactresolver.LocatorConfig{
							ID: "com.palantir.test:test-asset-2:3.4.5",
						},
					},
				},
			},
		},
	}
}

func TestDefaultTasksPluginsConfig(t *testing.T) {
	original := defaultPluginsConfig
	defer func() {
		defaultPluginsConfig = original
	}()
	defaultPluginsConfig = testDefaultPluginsConfig()

	for i, tc := range []struct {
		name string
		in   DefaultTasksConfig
		want PluginsConfig
	}{
		{
			"empty task param results in default configuration",
			DefaultTasksConfig{},
			PluginsConfig{
				DefaultResolvers: []string{
					defaultResolver,
				},
				Plugins: []SinglePluginConfig{
					{
						LocatorWithResolverConfig: artifactresolver.LocatorWithResolverConfig{
							Locator: artifactresolver.LocatorConfig{
								ID: "com.palantir.test:test-plugin:1.2.3",
							},
						},
						Assets: []artifactresolver.LocatorWithResolverConfig{
							{
								Locator: artifactresolver.LocatorConfig{
									ID: "com.palantir.test:test-asset-1:2.3.4",
								},
							},
							{
								Locator: artifactresolver.LocatorConfig{
									ID: "com.palantir.test:test-asset-2:3.4.5",
								},
							},
						},
					},
				},
			},
		},
		{
			"specifying custom resolver overrides resolver",
			DefaultTasksConfig{
				Tasks: map[string]SingleDefaultTaskConfig{
					"com.palantir.test:test-plugin": {
						LocatorWithResolverConfig: artifactresolver.LocatorWithResolverConfig{
							Resolver: "custom-resolver",
						},
					},
				},
			},
			PluginsConfig{
				DefaultResolvers: []string{
					defaultResolver,
				},
				Plugins: []SinglePluginConfig{
					{
						LocatorWithResolverConfig: artifactresolver.LocatorWithResolverConfig{
							Locator: artifactresolver.LocatorConfig{
								ID: "com.palantir.test:test-plugin:1.2.3",
							},
							Resolver: "custom-resolver",
						},
						Assets: []artifactresolver.LocatorWithResolverConfig{
							{
								Locator: artifactresolver.LocatorConfig{
									ID: "com.palantir.test:test-asset-1:2.3.4",
								},
							},
							{
								Locator: artifactresolver.LocatorConfig{
									ID: "com.palantir.test:test-asset-2:3.4.5",
								},
							},
						},
					},
				},
			},
		},
		{
			"specifying custom locator overrides locator",
			DefaultTasksConfig{
				Tasks: map[string]SingleDefaultTaskConfig{
					"com.palantir.test:test-plugin": {
						LocatorWithResolverConfig: artifactresolver.LocatorWithResolverConfig{
							Locator: artifactresolver.LocatorConfig{
								ID: "com.palantir.godel:override:1.2.3",
							},
						},
					},
				},
			},
			PluginsConfig{
				DefaultResolvers: []string{
					defaultResolver,
				},
				Plugins: []SinglePluginConfig{
					{
						LocatorWithResolverConfig: artifactresolver.LocatorWithResolverConfig{
							Locator: artifactresolver.LocatorConfig{
								ID: "com.palantir.godel:override:1.2.3",
							},
						},
						Assets: []artifactresolver.LocatorWithResolverConfig{
							{
								Locator: artifactresolver.LocatorConfig{
									ID: "com.palantir.test:test-asset-1:2.3.4",
								},
							},
							{
								Locator: artifactresolver.LocatorConfig{
									ID: "com.palantir.test:test-asset-2:3.4.5",
								},
							},
						},
					},
				},
			},
		},
		{
			"specifying default resolver appends default resolver",
			DefaultTasksConfig{
				DefaultResolvers: []string{
					"default/repo/{{GroupPath}}/{{Product}}/{{Version}}/{{Product}}-{{OS}}-{{Arch}}-{{Version}}.tgz",
				},
				Tasks: map[string]SingleDefaultTaskConfig{
					"com.palantir.test:test-plugin": {
						LocatorWithResolverConfig: artifactresolver.LocatorWithResolverConfig{
							Locator: artifactresolver.LocatorConfig{
								ID: "com.palantir.godel:override:1.2.3",
							},
						},
					},
				},
			},
			PluginsConfig{
				DefaultResolvers: []string{
					defaultResolver,
					"default/repo/{{GroupPath}}/{{Product}}/{{Version}}/{{Product}}-{{OS}}-{{Arch}}-{{Version}}.tgz",
				},
				Plugins: []SinglePluginConfig{
					{
						LocatorWithResolverConfig: artifactresolver.LocatorWithResolverConfig{
							Locator: artifactresolver.LocatorConfig{
								ID: "com.palantir.godel:override:1.2.3",
							},
						},
						Assets: []artifactresolver.LocatorWithResolverConfig{
							{
								Locator: artifactresolver.LocatorConfig{
									ID: "com.palantir.test:test-asset-1:2.3.4",
								},
							},
							{
								Locator: artifactresolver.LocatorConfig{
									ID: "com.palantir.test:test-asset-2:3.4.5",
								},
							},
						},
					},
				},
			},
		},
		{
			"specifying custom asset adds only that asset",
			DefaultTasksConfig{
				Tasks: map[string]SingleDefaultTaskConfig{
					"com.palantir.test:test-plugin": {
						Assets: []artifactresolver.LocatorWithResolverConfig{
							{
								Locator: artifactresolver.LocatorConfig{
									ID: "com.palantir.godel:custom-asset:1.2.3",
								},
							},
						},
					},
				},
			},
			PluginsConfig{
				DefaultResolvers: []string{
					defaultResolver,
				},
				Plugins: []SinglePluginConfig{
					{
						LocatorWithResolverConfig: artifactresolver.LocatorWithResolverConfig{
							Locator: artifactresolver.LocatorConfig{
								ID: "com.palantir.test:test-plugin:1.2.3",
							},
						},
						Assets: []artifactresolver.LocatorWithResolverConfig{
							{
								Locator: artifactresolver.LocatorConfig{
									ID: "com.palantir.test:test-asset-1:2.3.4",
								},
							},
							{
								Locator: artifactresolver.LocatorConfig{
									ID: "com.palantir.test:test-asset-2:3.4.5",
								},
							},
							{
								Locator: artifactresolver.LocatorConfig{
									ID: "com.palantir.godel:custom-asset:1.2.3",
								},
							},
						},
					},
				},
			},
		},
		{
			"setting exclude all and specifying custom asset adds asset to default",
			DefaultTasksConfig{
				Tasks: map[string]SingleDefaultTaskConfig{
					"com.palantir.test:test-plugin": {
						ExcludeAllDefaultAssets: true,
						Assets: []artifactresolver.LocatorWithResolverConfig{
							{
								Locator: artifactresolver.LocatorConfig{
									ID: "com.palantir.godel:custom-asset:1.2.3",
								},
							},
						},
					},
				},
			},
			PluginsConfig{
				DefaultResolvers: []string{
					defaultResolver,
				},
				Plugins: []SinglePluginConfig{
					{
						LocatorWithResolverConfig: artifactresolver.LocatorWithResolverConfig{
							Locator: artifactresolver.LocatorConfig{
								ID: "com.palantir.test:test-plugin:1.2.3",
							},
						},
						Assets: []artifactresolver.LocatorWithResolverConfig{
							{
								Locator: artifactresolver.LocatorConfig{
									ID: "com.palantir.godel:custom-asset:1.2.3",
								},
							},
						},
					},
				},
			},
		},
		{
			"specifying default asset with exclude and custom asset adds asset",
			DefaultTasksConfig{
				Tasks: map[string]SingleDefaultTaskConfig{
					"com.palantir.test:test-plugin": {
						DefaultAssetsToExclude: []string{
							"com.palantir.test:test-asset-2",
						},
						Assets: []artifactresolver.LocatorWithResolverConfig{
							{
								Locator: artifactresolver.LocatorConfig{
									ID: "com.palantir.godel:custom-asset:1.2.3",
								},
							},
						},
					},
				},
			},
			PluginsConfig{
				DefaultResolvers: []string{
					defaultResolver,
				},
				Plugins: []SinglePluginConfig{
					{
						LocatorWithResolverConfig: artifactresolver.LocatorWithResolverConfig{
							Locator: artifactresolver.LocatorConfig{
								ID: "com.palantir.test:test-plugin:1.2.3",
							},
						},
						Assets: []artifactresolver.LocatorWithResolverConfig{
							{
								Locator: artifactresolver.LocatorConfig{
									ID: "com.palantir.test:test-asset-1:2.3.4",
								},
							},
							{
								Locator: artifactresolver.LocatorConfig{
									ID: "com.palantir.godel:custom-asset:1.2.3",
								},
							},
						},
					},
				},
			},
		},
	} {
		got := DefaultTasksPluginsConfig(tc.in)
		assert.Equal(t, tc.want, got, "Case %d: %s", i, tc.name)
	}
}
