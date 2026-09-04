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

const (
	defaultResolver     = "https://github.com/{{index GroupParts 1}}/{{index GroupParts 2}}/releases/download/v{{Version}}/{{Product}}-{{Version}}-{{OS}}-{{Arch}}.tgz"
	defaultResolverYAML = "https://github.com/{{index GroupParts 1}}/{{index GroupParts 2}}/releases/download/v{{Version}}/{{Product}}-{{Version}}.yml.tgz"
)

var defaultPluginsConfig = config.PluginsConfig{
	DefaultResolvers: []string{
		defaultResolver,
		defaultResolverYAML,
	},
	Plugins: config.ToSinglePluginConfigs([]config.SinglePluginConfig{
		{
			LocatorWithResolverConfig: config.ToLocatorWithResolverConfig(config.LocatorWithResolverConfig{
				Locator: config.ToLocatorConfig(config.LocatorConfig{
					ID: "com.palantir.distgo:dist-plugin:1.115.0",
					Checksums: map[string]string{
						"darwin-amd64": "c3af4854fa9aba352fd7715694b43c72531b0589a77178dd1154db371d596e3d",
						"darwin-arm64": "f176f1002b0a0a595d5821364372eb0a153bcbb3d4a3900782d07e69021ae832",
						"linux-amd64":  "22f60f06e7ea19d7786c2799889ca7305524feabcdbe79e0919b779a6e745e75",
						"linux-arm64":  "b13c97c797a1c5799dedf6daecbe17f91ebb38bdc50cc083439418240cf680e6",
					},
				}),
			}),
		},
		{
			LocatorWithResolverConfig: config.ToLocatorWithResolverConfig(config.LocatorWithResolverConfig{
				Locator: config.ToLocatorConfig(config.LocatorConfig{
					ID: "com.palantir.godel-format-plugin:format-plugin:1.63.0",
					Checksums: map[string]string{
						"darwin-amd64": "e38f1699edc8ac8f174eea4de7ce6be9c906a32cb3d3fc5d81cf350372330f01",
						"darwin-arm64": "63de0bd49defd48c439005fe05b4a64a844f9a919b26d1e1bd348792c71cc45f",
						"linux-amd64":  "feb4ebc2f19d93a08d7e431d7b45ddaad9d95f4936b64c8384771d1257f52a1a",
						"linux-arm64":  "b173f8748b334f25a211c91a90916eed5e24186218bd1a676cd6b4b77b0ec6e8",
					},
				}),
			}),
			Assets: config.ToLocatorWithResolverConfigs([]config.LocatorWithResolverConfig{
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-format-asset-ptimports:ptimports-asset:1.64.0",
						Checksums: map[string]string{
							"darwin-amd64": "6a74223775e037de44e784410d9d2c184805c8d609c5d4d27ab2283efd4e0ad3",
							"darwin-arm64": "2db933f316fa524e984190458501914e16e40fe542b9f31f004f69af663f6316",
							"linux-amd64":  "daa187c44f37193f83d33e00108b70c7080f7dad96bec97d6d9a38a02cd008a3",
							"linux-arm64":  "ed79e6beabbcfc5c12a5bf1d66c5c82a6cccb87affa3d8afd571fa60d004f124",
						},
					}),
				},
			}),
		},
		{
			LocatorWithResolverConfig: config.ToLocatorWithResolverConfig(config.LocatorWithResolverConfig{
				Locator: config.ToLocatorConfig(config.LocatorConfig{
					ID: "com.palantir.godel-goland-plugin:goland-plugin:1.57.0",
					Checksums: map[string]string{
						"darwin-amd64": "f4437383e0ee2e4802d5acc3e3fcb483800a6daa2bf43fbd1934191acbf8e7a2",
						"darwin-arm64": "125073255479a38097bf74ee64751b6385e6cd54399f6fceb877b11d8db1eb03",
						"linux-amd64":  "a5adc0566e4d5dbac77820d5331c5b00638459d0357276034aadf47a99c28863",
						"linux-arm64":  "b459649ca9e1adf0853a470e3a8a4f8bba9ff1dbfcfe7bd1c51eeec78a6a967b",
					},
				}),
			}),
		},
		{
			LocatorWithResolverConfig: config.ToLocatorWithResolverConfig(config.LocatorWithResolverConfig{
				Locator: config.ToLocatorConfig(config.LocatorConfig{
					ID: "com.palantir.godel-golangci-lint-plugin:golangci-lint-plugin:0.21.0",
					Checksums: map[string]string{
						"darwin-amd64": "c0aceb9d495b6749bef904aecc509f7412f1531281b0dc3cec66431b9ddc198b",
						"darwin-arm64": "a8b315645a7d2894731b9a0ad39bd088b366fa29909f0dfa3aeb1f60d8c64b16",
						"linux-amd64":  "012c4df80c7c10fd294e8701b6d2e8f5c7a8b85911d58c0d43e655e30b6dbc75",
						"linux-arm64":  "921c65faea50c73d67f9e3383b588af320dd2cd964d9dba11e4e09bb75a0f778",
					},
				}),
			}),
			Assets: config.ToLocatorWithResolverConfigs([]config.LocatorWithResolverConfig{
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.golangci-lint-palantir:golangci-lint-palantir:0.21.0",
						Checksums: map[string]string{
							"darwin-amd64": "677b808cd09e5e8b0e761c408e7a11379ffb4281b9a28683f8790c2ae0a24144",
							"darwin-arm64": "e4d67588a3d39ded6fe1b6ade12f4b0724293f2083d6433208cfd2e646652839",
							"linux-amd64":  "a4684d2f20a53782dc0cffbb66faffc17aaf36e1ebbeefc464991f4aa66714f4",
							"linux-arm64":  "e5309db06dfc93c5e1f4fce0f9f680a5a17d921b1c55d94eb65e0d43d770cec4",
						},
					}),
				},
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.golangci-lint-palantir:golangci-lint-palantir-config:0.21.0",
						Checksums: map[string]string{
							"darwin-amd64": "b297a9384c93a7f63bd1cda880e9827bd04e7cf12ef2dac1ba8cb6426e1a175c",
							"darwin-arm64": "b297a9384c93a7f63bd1cda880e9827bd04e7cf12ef2dac1ba8cb6426e1a175c",
							"linux-amd64":  "b297a9384c93a7f63bd1cda880e9827bd04e7cf12ef2dac1ba8cb6426e1a175c",
							"linux-arm64":  "b297a9384c93a7f63bd1cda880e9827bd04e7cf12ef2dac1ba8cb6426e1a175c",
						},
					}),
				},
			}),
		},
		{
			LocatorWithResolverConfig: config.ToLocatorWithResolverConfig(config.LocatorWithResolverConfig{
				Locator: config.ToLocatorConfig(config.LocatorConfig{
					ID: "com.palantir.godel-license-plugin:license-plugin:1.58.0",
					Checksums: map[string]string{
						"darwin-amd64": "f2e8463ba804f03ec99b40b806387ffe8554f99ee9df29de2782097e0c0e5367",
						"darwin-arm64": "532482c8b2901317807c2397fe4fe669422ee78363c222b71d1ef2cd39b8904f",
						"linux-amd64":  "1c83ea7c4db416a64625600e96fcc730f8cea4d7421692b66268bbe9fa70a079",
						"linux-arm64":  "814225ecf1cfb4bf07bd10e88452f2390c2c9780e937d49c5387a96136b07a1d",
					},
				}),
			}),
		},
		{
			LocatorWithResolverConfig: config.ToLocatorWithResolverConfig(config.LocatorWithResolverConfig{
				Locator: config.ToLocatorConfig(config.LocatorConfig{
					ID: "com.palantir.godel-test-plugin:test-plugin:1.58.0",
					Checksums: map[string]string{
						"darwin-amd64": "964307b324880c56ba9375069374ffeef82c2d313d3dcebc16a75b3e04192661",
						"darwin-arm64": "adc5d1be1b39183eae62cf0b3f12d6509fa7a9d6322fca3a0e149baf721d4e56",
						"linux-amd64":  "7f71fd3c3d0fe6d06526b7dec0cabdb8b83af4165d887498016423ffc24af315",
						"linux-arm64":  "3a188680941e3ca55e802b31e8e9b404f35ea71ac921c05210471d61d9a318ea",
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
	pluginsCfg.DefaultResolvers = pluginsinternal.Uniquify(append(cfg.DefaultResolvers, pluginsCfg.DefaultResolvers...))

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
