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
					ID: "com.palantir.distgo:dist-plugin:1.94.0",
					Checksums: map[string]string{
						"darwin-amd64": "5038dca1c1af984c621f3a7fc6cd4346157f4c6b8a9b45d368db7b94b0ba5d53",
						"darwin-arm64": "5100166d6ec819280026770bca0a589d657ec3e96b3eae29e206df53431018ac",
						"linux-amd64":  "44064688246214e9a1344a5f6f6b07fc572533e9afeb871e9d111f3cc758745c",
						"linux-arm64":  "9e48d3a14a44a48b669c33ad8fd4ab63cd7d9ae612d7f169ad4f7f554cc2aa9a",
					},
				}),
			}),
		},
		{
			LocatorWithResolverConfig: config.ToLocatorWithResolverConfig(config.LocatorWithResolverConfig{
				Locator: config.ToLocatorConfig(config.LocatorConfig{
					ID: "com.palantir.godel-format-plugin:format-plugin:1.54.0",
					Checksums: map[string]string{
						"darwin-amd64": "072589605b666f29a777f42b5a33c37bc636df8169b413e8a125fe5b3701aa5d",
						"darwin-arm64": "290b90b547b5d293ee6d431576270d17d46a90a60d3f5ac1149bb37ab05b45d5",
						"linux-amd64":  "d2fe1cb63623a0947d3a83a5cb4bdaf1dbc71ea73c15497f8556fe2979cbd399",
						"linux-arm64":  "aa4f6d1c923023ce2446d2f497f0ba85c5e603aa441410bda5deaad46dca251e",
					},
				}),
			}),
			Assets: config.ToLocatorWithResolverConfigs([]config.LocatorWithResolverConfig{
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-format-asset-ptimports:ptimports-asset:1.53.0",
						Checksums: map[string]string{
							"darwin-amd64": "f34997c27bb7e5bd48cb3783d41f1c10b5a41fc511014d6f8af0462c527b514b",
							"darwin-arm64": "2dfb9a5f35cacf9f2151b7553c169bac0b59e0bd0d46f83637412710b5328955",
							"linux-amd64":  "a904aff8ba4094cb45ce0ede61d20553e393c38987efc308808e2c846d6aafd6",
							"linux-arm64":  "0d0c12e7d04b76fee2688a2ad98ea7a5d490ade465bc1133eeda001b2073d638",
						},
					}),
				},
			}),
		},
		{
			LocatorWithResolverConfig: config.ToLocatorWithResolverConfig(config.LocatorWithResolverConfig{
				Locator: config.ToLocatorConfig(config.LocatorConfig{
					ID: "com.palantir.godel-goland-plugin:goland-plugin:1.48.0",
					Checksums: map[string]string{
						"darwin-amd64": "dc5867fcddf3dc8813edfc9bb3f23c0b5af1495c63e69170f12f552299baa45d",
						"darwin-arm64": "f1d40fff2aa0c043233b0fd7da23f178d29df73fe07b3c40b92c448db8cf5eeb",
						"linux-amd64":  "89892eee20bdc786d87dd845a2753bd53fdb8823d5ddd42d2955b7ec661f94e1",
						"linux-arm64":  "c637351ca4a7afe446fce58347f66e7797ec38f9704da773010b514e0f347740",
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
					ID: "com.palantir.godel-license-plugin:license-plugin:1.50.0",
					Checksums: map[string]string{
						"darwin-amd64": "16bb46a4202e081af1614ab586c8050290873c6671af305564f2f5346fc64a6b",
						"darwin-arm64": "fcd54d83fdc3def6e0534aa3d264ac984f04803ac22dd8fd9dc9d5ae9d89003c",
						"linux-amd64":  "d49fbcceba85dcb451b9c264f61ef1ef1eb47bbbed8a7575fb4e10905e1f1a56",
						"linux-arm64":  "79d5ab0b47cc57ffb77e43ef0342ac3fb435a4f6863087265dfbb73c13a81481",
					},
				}),
			}),
		},
		{
			LocatorWithResolverConfig: config.ToLocatorWithResolverConfig(config.LocatorWithResolverConfig{
				Locator: config.ToLocatorConfig(config.LocatorConfig{
					ID: "com.palantir.godel-test-plugin:test-plugin:1.50.0",
					Checksums: map[string]string{
						"darwin-amd64": "bfed52eb653de39b5386922ccefb4bfbfea593d90c68c8ab6beb03d7e1b1054f",
						"darwin-arm64": "dac670f98b3ba299b4340a75ecdd9c51b7ac2de1290c71960bff764d1720fb0a",
						"linux-amd64":  "9cf623ae7b9b905c516811f4ff8c435d5196fb9dc4b28479afd4d289ca6a4d22",
						"linux-arm64":  "56d7905c2b21474ef2e2ab32cce0ba9736a7afe67eaad7b4a38c2e8908b21025",
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
