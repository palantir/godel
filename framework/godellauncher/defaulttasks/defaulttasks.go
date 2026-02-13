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
					ID: "com.palantir.distgo:dist-plugin:1.93.0",
					Checksums: map[string]string{
						"darwin-amd64": "e79748ec8d4b558fc1c353b450f4a3e1e39c13be5e907463a8ae6caf3e09929d",
						"darwin-arm64": "8ec11e6cc480f893962f30edeb72a3c2920a23e764266b08584865cbc74b69e7",
						"linux-amd64":  "0089bdd3873b800ee9e5210cfd17b174f0d113144f203eeafcefa34341bc638f",
						"linux-arm64":  "f2edd94e9f8951fd1b26b3e1f560c7b103c124c59f5a83786777126dc426ea36",
					},
				}),
			}),
		},
		{
			LocatorWithResolverConfig: config.ToLocatorWithResolverConfig(config.LocatorWithResolverConfig{
				Locator: config.ToLocatorConfig(config.LocatorConfig{
					ID: "com.palantir.godel-format-plugin:format-plugin:1.53.0",
					Checksums: map[string]string{
						"darwin-amd64": "88a3c1666c2bddc350fc6601a32500d87886c3fcf95c477b31619666cd04b0a7",
						"darwin-arm64": "230cd3ad76d61b9b20899d70249e3a0873e92ed4d4a4d5667d5b2093f004a9d5",
						"linux-amd64":  "aa3135ae1131e4ad55c2813fd39da039e3c90e1fdc53338802ca6aeee6a1f62b",
						"linux-arm64":  "eb08998b772564284b1c8205643c2a81304c85e7cc8b433902743f6592dfd1f1",
					},
				}),
			}),
			Assets: config.ToLocatorWithResolverConfigs([]config.LocatorWithResolverConfig{
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-format-asset-ptimports:ptimports-asset:1.52.0",
						Checksums: map[string]string{
							"darwin-amd64": "9284c15341ad6d0eb0b7c5c7b01b584288d7f141e2ba5050d2f9f36c59d1e3e9",
							"darwin-arm64": "70a79e4e796ca89e856e574afde6dc1f1d45affc7279379120d47689821d1ad4",
							"linux-amd64":  "ebc76db1e196238dc55f15c438fe0a4aa1b2b9c52db03f8589a7519c506e71ab",
							"linux-arm64":  "a12f69cb2c616d49a13235ed0724246891ac767d9e6e738a196fc0317abf0507",
						},
					}),
				},
			}),
		},
		{
			LocatorWithResolverConfig: config.ToLocatorWithResolverConfig(config.LocatorWithResolverConfig{
				Locator: config.ToLocatorConfig(config.LocatorConfig{
					ID: "com.palantir.godel-goland-plugin:goland-plugin:1.47.0",
					Checksums: map[string]string{
						"darwin-amd64": "7dec66ae76c870e778a6ff79396e85986914fc4f7d5083915dada0c2c12e921c",
						"darwin-arm64": "012ee49427b0ba3498cc51bab65c65e32e03b0f1a4d5eab42780f43a1076c15d",
						"linux-amd64":  "69ada20846ed9fda2658f89c0248b6534d3d834ba0bf3a1fc6fe5c2bbe715eab",
						"linux-arm64":  "acefcca1716ad4b5279fb481fc273e296b8ff05508b270b114032fadeb7610d9",
					},
				}),
			}),
		},
		{
			LocatorWithResolverConfig: config.ToLocatorWithResolverConfig(config.LocatorWithResolverConfig{
				Locator: config.ToLocatorConfig(config.LocatorConfig{
					ID: "com.palantir.godel-golangci-lint-plugin:golangci-lint-plugin:0.9.0",
					Checksums: map[string]string{
						"darwin-amd64": "2bed1794aade4883747125b48ab3f60b414538e3b64f311eb860866c77851b28",
						"darwin-arm64": "684a7c140df92b66b7d3429d171eb3b88e6fd93bcaf59629d8862519f9b9d860",
						"linux-amd64":  "56193ba96d065d5a17dba127086d37951e818a2546a5a5ea38565ffe0a4cdb5c",
						"linux-arm64":  "5a65dd18429c661fc840b2778575ca8b97c133e698c0560dd67b89428e93ef7f",
					},
				}),
			}),
			Assets: config.ToLocatorWithResolverConfigs([]config.LocatorWithResolverConfig{
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.golangci-lint-palantir:golangci-lint-palantir:0.7.0",
						Checksums: map[string]string{
							"darwin-amd64": "44c6982c63e4ec69550403460252d20edb21e20c2ee08a5dc517ff5b84c4f4c6",
							"darwin-arm64": "bcd4447f4a84ed8e73630f2dee664bf08d61a5a79625a3033e5fd6ca58efe02d",
							"linux-amd64":  "1bf4767a613b285ccf33936ff26e28201864152bd6d381180bfef85ca529bca2",
							"linux-arm64":  "53934ee01b44dc10fe51f1e8afbd48b82c67b524c80f7e1aaec584f9fbe4fab2",
						},
					}),
				},
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.golangci-lint-palantir:golangci-lint-palantir-config:0.7.0",
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
					ID: "com.palantir.godel-license-plugin:license-plugin:1.49.0",
					Checksums: map[string]string{
						"darwin-amd64": "9bed959c0964a6a4051fb41f80af3c4fe7233cc128b125bc15a965ea9437e5f2",
						"darwin-arm64": "542c40727e3ab1c0964020187676486ef9e185c4e23f6482bf3fa02fb80865f7",
						"linux-amd64":  "63791eb2a7f3aecd85e1f3b0fdb1895b794a435c4eb44511160dca528be6a3d0",
						"linux-arm64":  "d03b47ddd347f2416614f0d87cfb0f735f0da312028946e840f3e85216db7531",
					},
				}),
			}),
		},
		{
			LocatorWithResolverConfig: config.ToLocatorWithResolverConfig(config.LocatorWithResolverConfig{
				Locator: config.ToLocatorConfig(config.LocatorConfig{
					ID: "com.palantir.godel-test-plugin:test-plugin:1.49.0",
					Checksums: map[string]string{
						"darwin-amd64": "8204bae08325eb3bcd96be7665abe1f2325c98fbbb173f447a82ca83762455ab",
						"darwin-arm64": "57f4ca069e9d7b6087ee2b011aa9a64f6cf1d9469daea3b850c445ed20384e28",
						"linux-amd64":  "f8f82e02babdb135179ca00aafbc3a471d8784b72c38a1306bb783de21a03356",
						"linux-arm64":  "1e67c9f82c5a40bfd0897127c91e6425e3bd7b337fada24b3f44aaaff6d85187",
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
