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
)

func testDefaultPluginsConfig() PluginsConfig {
	return PluginsConfig{
		DefaultResolvers: []string{defaultResolver},
		Plugins: []SinglePluginConfig{
			{
				LocatorWithResolverConfig: LocatorWithResolverConfig{
					Locator: LocatorConfig{
						ID: "com.palantir.test:test-plugin:1.2.3",
					},
				},
				Assets: []LocatorWithResolverConfig{
					{
						Locator: LocatorConfig{
							ID: "com.palantir.test:test-asset-1:2.3.4",
						},
					},
					{
						Locator: LocatorConfig{
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
						LocatorWithResolverConfig: LocatorWithResolverConfig{
							Locator: LocatorConfig{
								ID: "com.palantir.test:test-plugin:1.2.3",
							},
						},
						Assets: []LocatorWithResolverConfig{
							{
								Locator: LocatorConfig{
									ID: "com.palantir.test:test-asset-1:2.3.4",
								},
							},
							{
								Locator: LocatorConfig{
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
				"com.palantir.test:test-plugin": {
					LocatorWithResolverConfig: LocatorWithResolverConfig{
						Resolver: "custom-resolver",
					},
				},
			},
			PluginsConfig{
				DefaultResolvers: []string{
					defaultResolver,
				},
				Plugins: []SinglePluginConfig{
					{
						LocatorWithResolverConfig: LocatorWithResolverConfig{
							Locator: LocatorConfig{
								ID: "com.palantir.test:test-plugin:1.2.3",
							},
							Resolver: "custom-resolver",
						},
						Assets: []LocatorWithResolverConfig{
							{
								Locator: LocatorConfig{
									ID: "com.palantir.test:test-asset-1:2.3.4",
								},
							},
							{
								Locator: LocatorConfig{
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
				"com.palantir.test:test-plugin": {
					LocatorWithResolverConfig: LocatorWithResolverConfig{
						Locator: LocatorConfig{
							ID: "com.palantir.godel:override:1.2.3",
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
						LocatorWithResolverConfig: LocatorWithResolverConfig{
							Locator: LocatorConfig{
								ID: "com.palantir.godel:override:1.2.3",
							},
						},
						Assets: []LocatorWithResolverConfig{
							{
								Locator: LocatorConfig{
									ID: "com.palantir.test:test-asset-1:2.3.4",
								},
							},
							{
								Locator: LocatorConfig{
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
				"com.palantir.test:test-plugin": {
					Assets: []LocatorWithResolverConfig{
						{
							Locator: LocatorConfig{
								ID: "com.palantir.godel:custom-asset:1.2.3",
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
						LocatorWithResolverConfig: LocatorWithResolverConfig{
							Locator: LocatorConfig{
								ID: "com.palantir.test:test-plugin:1.2.3",
							},
						},
						Assets: []LocatorWithResolverConfig{
							{
								Locator: LocatorConfig{
									ID: "com.palantir.test:test-asset-1:2.3.4",
								},
							},
							{
								Locator: LocatorConfig{
									ID: "com.palantir.test:test-asset-2:3.4.5",
								},
							},
							{
								Locator: LocatorConfig{
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
				"com.palantir.test:test-plugin": {
					ExcludeAllDefaultAssets: true,
					Assets: []LocatorWithResolverConfig{
						{
							Locator: LocatorConfig{
								ID: "com.palantir.godel:custom-asset:1.2.3",
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
						LocatorWithResolverConfig: LocatorWithResolverConfig{
							Locator: LocatorConfig{
								ID: "com.palantir.test:test-plugin:1.2.3",
							},
						},
						Assets: []LocatorWithResolverConfig{
							{
								Locator: LocatorConfig{
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
				"com.palantir.test:test-plugin": {
					DefaultAssetsToExclude: []string{
						"com.palantir.test:test-asset-2",
					},
					Assets: []LocatorWithResolverConfig{
						{
							Locator: LocatorConfig{
								ID: "com.palantir.godel:custom-asset:1.2.3",
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
						LocatorWithResolverConfig: LocatorWithResolverConfig{
							Locator: LocatorConfig{
								ID: "com.palantir.test:test-plugin:1.2.3",
							},
						},
						Assets: []LocatorWithResolverConfig{
							{
								Locator: LocatorConfig{
									ID: "com.palantir.test:test-asset-1:2.3.4",
								},
							},
							{
								Locator: LocatorConfig{
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
