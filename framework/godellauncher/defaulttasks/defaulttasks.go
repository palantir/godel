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
	"sort"
	"strings"

	"github.com/palantir/godel/v2/framework/godel/config"
	"github.com/palantir/godel/v2/framework/internal/pluginsinternal"
	"github.com/pkg/errors"
)

const defaultResolver = "https://palantir.bintray.com/releases/{{GroupPath}}/{{Product}}/{{Version}}/{{Product}}-{{Version}}-{{OS}}-{{Arch}}.tgz"

var defaultPluginsConfig = config.PluginsConfig{
	DefaultResolvers: []string{
		defaultResolver,
	},
	Plugins: config.ToSinglePluginConfigs([]config.SinglePluginConfig{
		{
			LocatorWithResolverConfig: config.ToLocatorWithResolverConfig(config.LocatorWithResolverConfig{
				Locator: config.ToLocatorConfig(config.LocatorConfig{
					ID: "com.palantir.distgo:dist-plugin:1.21.0",
					Checksums: map[string]string{
						"darwin-amd64": "a6d97b81e21c7b66cfbaacb9dcf32c1a9b239195558a370d5a14e1ae29c3b3f4",
						"linux-amd64":  "e0785618eabbdd37bc594d47f8b3ab2bba397a67fe9f05acd7849f15742fc0d9",
					},
				}),
			}),
		},
		{
			LocatorWithResolverConfig: config.ToLocatorWithResolverConfig(config.LocatorWithResolverConfig{
				Locator: config.ToLocatorConfig(config.LocatorConfig{
					ID: "com.palantir.godel-format-plugin:format-plugin:1.6.0",
					Checksums: map[string]string{
						"darwin-amd64": "65462d693c3d0dda205bee1a10af8149e084909f2a0f2076d93820d759c33188",
						"linux-amd64":  "efac3b8f8e3fee0a88e0d9dc93a82dd68df92dcef8b7ab22aef5e6e982ad723a",
					},
				}),
			}),
			Assets: config.ToLocatorWithResolverConfigs([]config.LocatorWithResolverConfig{
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-format-asset-ptimports:ptimports-asset:1.5.0",
						Checksums: map[string]string{
							"darwin-amd64": "1089054a9fb08263c24b2fbf5623112c3a5e7ee5aff10798ce8d92a7987e0e8b",
							"linux-amd64":  "f8962fb0f71a49b7fdc33ad0c937e5f2b824c7a27eb815850ff315ffce2ff204",
						},
					}),
				},
			}),
		},
		{
			LocatorWithResolverConfig: config.ToLocatorWithResolverConfig(config.LocatorWithResolverConfig{
				Locator: config.ToLocatorConfig(config.LocatorConfig{
					ID: "com.palantir.godel-goland-plugin:goland-plugin:1.2.0",
					Checksums: map[string]string{
						"darwin-amd64": "a318362fce2c67a8f27bb66a4e95799902b331f064b2302eb53adc028f8f13b6",
						"linux-amd64":  "6f7a10415f54110d41df093259463f0a80ad8020206b002ddd5b43733232f19e",
					},
				}),
			}),
		},
		{
			LocatorWithResolverConfig: config.ToLocatorWithResolverConfig(config.LocatorWithResolverConfig{
				Locator: config.ToLocatorConfig(config.LocatorConfig{
					ID: "com.palantir.okgo:check-plugin:1.9.0",
					Checksums: map[string]string{
						"darwin-amd64": "ef5420e43a1671cf55acf3d613e23b0aae65f2f9eacef1e433e4614bee5eec81",
						"linux-amd64":  "d579c956bc693599cf3eea26fbd80c1172efd46478b8e73873c19c72fdccc3a3",
					},
				}),
			}),
			Assets: config.ToLocatorWithResolverConfigs([]config.LocatorWithResolverConfig{
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-compiles:compiles-asset:1.6.0",
						Checksums: map[string]string{
							"darwin-amd64": "fe61d4fcec960e8e078167ccb4edc46840441ece06305f4b3494d9450cdcfe1a",
							"linux-amd64":  "5f5a6a6793d8177111555c3e87ee5a9d87eb2efc5ff52fc45b42cb02df00d7df",
						},
					}),
				},
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-deadcode:deadcode-asset:1.5.0",
						Checksums: map[string]string{
							"darwin-amd64": "be7796e7e3fab271b9c8e402a0bbaf2e9ea00b70d6317589db59a2e30c941cfd",
							"linux-amd64":  "695b1aaa2030390eaab5ae92ed6f836e8f9f79f16b03cc812e03776de2ee5043",
						},
					}),
				},
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-errcheck:errcheck-asset:1.7.0",
						Checksums: map[string]string{
							"darwin-amd64": "f928215dbf87200aa52e8c92ea93cec5e02da7aa2269aa989703d4d2fb356861",
							"linux-amd64":  "5d7f3d15d5379827ea202939a8db22e42f21bbd12eb8ed31742b709a3932c8fa",
						},
					}),
				},
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-golint:golint-asset:1.3.0",
						Checksums: map[string]string{
							"darwin-amd64": "d0e3d94b0ac6c269df8ace9defc6f997faffee102ab888f1cf03fb764c5f1ccf",
							"linux-amd64":  "887991e20847265e2bed5dbfa5c14bf81f0c3012e8f0e46cde1d293c879e8c1f",
						},
					}),
				},
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-govet:govet-asset:1.3.0",
						Checksums: map[string]string{
							"darwin-amd64": "c8ebe89c9ec922a571d84cc743b2324646b6d6c22508f4305438410322235a43",
							"linux-amd64":  "613b0321b8a3a08363eb58fd8ac079e0dd272494e06492b4c74e9c3c2d8a4a92",
						},
					}),
				},
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-importalias:importalias-asset:1.3.0",
						Checksums: map[string]string{
							"darwin-amd64": "1be88e3cee403cf2388cd0e0bd28a7a705af6153546baa1cc220983c750833cc",
							"linux-amd64":  "9d578fa8d74902c1000e06d38eeee32a6972129fcdb195f95e65df6c8ef6b6e6",
						},
					}),
				},
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-ineffassign:ineffassign-asset:1.3.0",
						Checksums: map[string]string{
							"darwin-amd64": "128f177a019d165b24e0202897981fb9b8a5e0c5188c491c71bf50266c342e52",
							"linux-amd64":  "ef733f085ed9ef59397df977eb2c17d5378de79b5c7a31fa325607be9da11290",
						},
					}),
				},
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-outparamcheck:outparamcheck-asset:1.7.0",
						Checksums: map[string]string{
							"darwin-amd64": "a4c6119ac4892e9c1f89032f8201e17639e6ab83e888d3a77ec724d53a7f0ef9",
							"linux-amd64":  "b484c3045993ad3db1b16e55eec2202d0a053e3a8a21629731baf3d5027e84e1",
						},
					}),
				},
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-unconvert:unconvert-asset:1.6.0",
						Checksums: map[string]string{
							"darwin-amd64": "fc004ea7a2d478c78bad745873855815c09e307097b5b18b63c7bcab4d0fd7bf",
							"linux-amd64":  "7058bf5c44899fb5a5115761cd7c02ae7c005d203383ccb939c513f19d341c7c",
						},
					}),
				},
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-varcheck:varcheck-asset:1.5.0",
						Checksums: map[string]string{
							"darwin-amd64": "4cd5d945d91ec78ee15b4ab87b45240de94e0c852ba079ae085deeff922a35c5",
							"linux-amd64":  "311399a476831cbe818d0f539d8a852b19d9a3247b3e17e1456f15f561d490ff",
						},
					}),
				},
			}),
		},
		{
			LocatorWithResolverConfig: config.ToLocatorWithResolverConfig(config.LocatorWithResolverConfig{
				Locator: config.ToLocatorConfig(config.LocatorConfig{
					ID: "com.palantir.godel-license-plugin:license-plugin:1.4.0",
					Checksums: map[string]string{
						"darwin-amd64": "3cc8a0d6ddf23e995246762ed8c26b12f64b268252ce916478a9973f4a5a8f62",
						"linux-amd64":  "460fe99400293bfc92b27a2645e41ed15651117c4148275614f2c964c031a375",
					},
				}),
			}),
		},
		{
			LocatorWithResolverConfig: config.ToLocatorWithResolverConfig(config.LocatorWithResolverConfig{
				Locator: config.ToLocatorConfig(config.LocatorConfig{
					ID: "com.palantir.godel-test-plugin:test-plugin:1.5.0",
					Checksums: map[string]string{
						"darwin-amd64": "6f09d24295fe54e101739ad53e1bece7856614e15f9bcefb8443265276d63dc5",
						"linux-amd64":  "669a78db753419f96c2626a31aba5e035764aa2a397e19f28243203d2ce3f0f0",
					},
				}),
			}),
		},
	}),
}

func BuiltinPluginsConfig() config.PluginsConfig {
	return defaultPluginsConfig
}

func PluginsConfig(cfg config.DefaultTasksConfig) (config.PluginsConfig, error) {
	// start with configuration that uses default resolver
	pluginsCfg := config.PluginsConfig{
		DefaultResolvers: defaultPluginsConfig.DefaultResolvers,
	}
	// append default resolvers provided by the configuration and uniquify
	pluginsCfg.DefaultResolvers = pluginsinternal.Uniquify(append(pluginsCfg.DefaultResolvers, cfg.DefaultResolvers...))

	defaultPluginKeys := make(map[string]struct{})
	for _, currPlugin := range defaultPluginsConfig.Plugins {
		currKey := locatorIDWithoutVersion(currPlugin.Locator.ID)
		defaultPluginKeys[currKey] = struct{}{}

		var assets []config.LocatorWithResolverConfig
		for _, asset := range currPlugin.Assets {
			assets = append(assets, config.LocatorWithResolverConfig(asset))
		}
		taskCfgV0, ok := cfg.Tasks[currKey]
		if !ok {
			// if custom configuration is not specified, use default and continue
			pluginsCfg.Plugins = append(pluginsCfg.Plugins, currPlugin)
			continue
		}
		taskCfg := config.SingleDefaultTaskConfig(taskCfgV0)

		// custom configuration was non-empty: start it with default LocatorWithResolver configuration
		currCfg := config.SinglePluginConfig{
			LocatorWithResolverConfig: currPlugin.LocatorWithResolverConfig,
		}
		if taskCfg.Locator.ID != "" {
			currCfg.Locator = taskCfg.Locator
		}
		if taskCfg.Resolver != "" {
			currCfg.Resolver = taskCfg.Resolver
		}

		currCfg.Assets = append(currCfg.Assets, config.ToLocatorWithResolverConfigs(assetConfigFromDefault(assets, taskCfg))...)
		currCfg.Assets = append(currCfg.Assets, taskCfg.Assets...)
		pluginsCfg.Plugins = append(pluginsCfg.Plugins, config.ToSinglePluginConfig(currCfg))
	}

	var invalidKeys []string
	for providedDefaultCfgKey := range cfg.Tasks {
		if _, ok := defaultPluginKeys[providedDefaultCfgKey]; ok {
			continue
		}
		invalidKeys = append(invalidKeys, providedDefaultCfgKey)
	}
	sort.Strings(invalidKeys)

	if len(invalidKeys) > 0 {
		var validKeys []string
		for k := range defaultPluginKeys {
			validKeys = append(validKeys, k)
		}
		sort.Strings(validKeys)
		return config.PluginsConfig{}, errors.Errorf("default-task key(s) specified but are not valid: %v. Valid values: %v", invalidKeys, validKeys)
	}

	return pluginsCfg, nil
}

func assetConfigFromDefault(baseCfg []config.LocatorWithResolverConfig, cfg config.SingleDefaultTaskConfig) []config.LocatorWithResolverConfig {
	if cfg.ExcludeAllDefaultAssets {
		return nil
	}
	exclude := make(map[string]struct{})
	for _, currExclude := range cfg.DefaultAssetsToExclude {
		exclude[currExclude] = struct{}{}
	}
	var out []config.LocatorWithResolverConfig
	for _, asset := range baseCfg {
		if _, ok := exclude[locatorIDWithoutVersion(asset.Locator.ID)]; ok {
			continue
		}
		out = append(out, asset)
	}
	return out
}

func locatorIDWithoutVersion(locatorID string) string {
	parts := strings.Split(locatorID, ":")
	return strings.Join(parts[:2], ":")
}
