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
					ID: "com.palantir.distgo:dist-plugin:1.96.0",
					Checksums: map[string]string{
						"darwin-amd64": "7d34a594922a00f56de475c45dc653986982534807f286f307fa075b88b1a390",
						"darwin-arm64": "743aac99dacbdc0bfe5cc64fd75abeafc2510e884d28674654e3c251b458bf33",
						"linux-amd64":  "6ff6ace3172162981ba980aa91f8f407c3805870967a01d5d6d80b330526112d",
						"linux-arm64":  "425823cc490f9030f60f198f895845e752fb50ae8cab30e90d1408cead86abb6",
					},
				}),
			}),
		},
		{
			LocatorWithResolverConfig: config.ToLocatorWithResolverConfig(config.LocatorWithResolverConfig{
				Locator: config.ToLocatorConfig(config.LocatorConfig{
					ID: "com.palantir.godel-format-plugin:format-plugin:1.55.0",
					Checksums: map[string]string{
						"darwin-amd64": "024a7d0052b6e25c1b68c697eb2a130a05eae9be9c53fc659d67526e84a3940e",
						"darwin-arm64": "ac2b78a494b63e4a1d970072d736885af1469c894742c4b3ed5f73fb154e1a59",
						"linux-amd64":  "0c9b0c1341ed92da2a271f9a435ffdea1e7365df05421ffdf3d2c25a71d496b1",
						"linux-arm64":  "f017b58dbd1fafbbbda2670810fa4c3f3d342bf3084f8730c40e577a391c64e8",
					},
				}),
			}),
			Assets: config.ToLocatorWithResolverConfigs([]config.LocatorWithResolverConfig{
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-format-asset-ptimports:ptimports-asset:1.54.0",
						Checksums: map[string]string{
							"darwin-amd64": "3c61cf1d2688b886fd05cea9f9805cc397f1a6aa081d92cde20a67a28a7823b6",
							"darwin-arm64": "6d8192276cc0eb7e7a691d1e157657a683aee66d6e5e30fb72161591be4de78d",
							"linux-amd64":  "94b8a036df29dc83af23b6eb4d6534eeb40492f36e8a4455c3a7da1261708a04",
							"linux-arm64":  "b942edbd967590ea49bd467ac2bd6763ed2cbb71d27d620ba1896336545cc070",
						},
					}),
				},
			}),
		},
		{
			LocatorWithResolverConfig: config.ToLocatorWithResolverConfig(config.LocatorWithResolverConfig{
				Locator: config.ToLocatorConfig(config.LocatorConfig{
					ID: "com.palantir.godel-goland-plugin:goland-plugin:1.49.0",
					Checksums: map[string]string{
						"darwin-amd64": "50011d651bd676002147b1e6567dd585fc4ff13ef9926aea6a2a1e4eaaa705d8",
						"darwin-arm64": "1c14f6facac476eeee93140a62d628780112127ae158c4a209d811aa5aa7aeb9",
						"linux-amd64":  "ead0a86f03ad377eea17a1b904bc4bfd82e29b041f77f1b03ebf9e0caf111671",
						"linux-arm64":  "d6f845ab1015679e4841a731af8acd5c0d309b44c7996946aa942247312eb95b",
					},
				}),
			}),
		},
		{
			LocatorWithResolverConfig: config.ToLocatorWithResolverConfig(config.LocatorWithResolverConfig{
				Locator: config.ToLocatorConfig(config.LocatorConfig{
					ID: "com.palantir.godel-golangci-lint-plugin:golangci-lint-plugin:0.10.0",
					Checksums: map[string]string{
						"darwin-amd64": "93b5583305778c690f6ac98737d5a2e1417cc3fa91c4296d016dc30584dad887",
						"darwin-arm64": "97f82f400642da529dd110c47bef87034e1c04c0a170a558357790ae041020c8",
						"linux-amd64":  "d8a3850412f34c45061fd2482df8ee962d6bca6ad6423a91c1de7f71ef4d7527",
						"linux-arm64":  "0688eaab226c752779a52b44848e2ae5b617501c29fd91d67053153b4d9957b4",
					},
				}),
			}),
			Assets: config.ToLocatorWithResolverConfigs([]config.LocatorWithResolverConfig{
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.golangci-lint-palantir:golangci-lint-palantir:0.8.0",
						Checksums: map[string]string{
							"darwin-amd64": "86c7ee33171bf33367e39af75bee807d89cbec906142b1ed5d1e6c20568b54c7",
							"darwin-arm64": "b4eccb02c926910ddd7ad7ea57e4645962960607f5a2c3a761d46254f40ee22a",
							"linux-amd64":  "f1355d76d7b61316c78d2c61b675589515c6bdc6c3e1b2634dfec829af7d318a",
							"linux-arm64":  "9528b1886eb9dcd06299d543a2fce1aa395c2a13553e7305a9d04a3fbde04773",
						},
					}),
				},
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.golangci-lint-palantir:golangci-lint-palantir-config:0.8.0",
						Checksums: map[string]string{
							"darwin-amd64": "ec55705e29ca8bff760f1b78fa482b914fa8944ca5d16fb63546de97730105d8",
							"darwin-arm64": "ec55705e29ca8bff760f1b78fa482b914fa8944ca5d16fb63546de97730105d8",
							"linux-amd64":  "ec55705e29ca8bff760f1b78fa482b914fa8944ca5d16fb63546de97730105d8",
							"linux-arm64":  "ec55705e29ca8bff760f1b78fa482b914fa8944ca5d16fb63546de97730105d8",
						},
					}),
				},
			}),
		},
		{
			LocatorWithResolverConfig: config.ToLocatorWithResolverConfig(config.LocatorWithResolverConfig{
				Locator: config.ToLocatorConfig(config.LocatorConfig{
					ID: "com.palantir.godel-license-plugin:license-plugin:1.51.0",
					Checksums: map[string]string{
						"darwin-amd64": "fd53d3f9242ba7225fdddb2d6274bff1f29da0e9cc0cd17d19a5b660d731e678",
						"darwin-arm64": "4a4778dd5e256966c04fa2954bf5594862e71701f8cb79393eb37082c1b729a0",
						"linux-amd64":  "a2d8048e01750cda025b279a1c8626ce6a939fd400b6436ca42e774ff369cc9d",
						"linux-arm64":  "fef4474fba49c3ea29aaf2d4df81c56cc96967393f5fea0da9f1757abba1956b",
					},
				}),
			}),
		},
		{
			LocatorWithResolverConfig: config.ToLocatorWithResolverConfig(config.LocatorWithResolverConfig{
				Locator: config.ToLocatorConfig(config.LocatorConfig{
					ID: "com.palantir.godel-test-plugin:test-plugin:1.51.0",
					Checksums: map[string]string{
						"darwin-amd64": "83076332c4b5b8a0942969ccf8a6ff8493afeb47f2c1d4a4c01524139e259b12",
						"darwin-arm64": "b0c53bf9744168b4a9093a12d62612a1b05d5eb1ab93625a7b69f325c5c9f755",
						"linux-amd64":  "b269ef2a35cf8e4dd6ee19b2069153a2bb887c1e9a7f74dcc238f19e36fdd605",
						"linux-arm64":  "86d1e45acd9aea3ffb09313db25e5e2b9f17607b21b498c15672e2ede6008591",
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
