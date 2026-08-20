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
					ID: "com.palantir.distgo:dist-plugin:1.112.0",
					Checksums: map[string]string{
						"darwin-amd64": "e2b4f548cb3cec1ba0ba11af70928a30fb67ab6992e67499eb7dfb38d44b226a",
						"darwin-arm64": "5c0c1c2cf7985bee67e85ec75181e2df39d4bce5d771c2b9a8dc90bb2fdb5243",
						"linux-amd64":  "b602fc173070628fa1d7117e8f85018a0723fa9b1a7d1495f77578b73e37ff04",
						"linux-arm64":  "8e895fc6fdeaad8037791d55dfffd51642094db17d8813264ffe9bd7b078f8c6",
					},
				}),
			}),
		},
		{
			LocatorWithResolverConfig: config.ToLocatorWithResolverConfig(config.LocatorWithResolverConfig{
				Locator: config.ToLocatorConfig(config.LocatorConfig{
					ID: "com.palantir.godel-format-plugin:format-plugin:1.60.0",
					Checksums: map[string]string{
						"darwin-amd64": "aaa6fe9fc2888ba9aa369cf948eec6c37280a717220183b5b7a0448d0963e927",
						"darwin-arm64": "c78773a137b9002ee371e3149d87fce2466e768302b8bbea78df9b0125352a10",
						"linux-amd64":  "7449c50666714a5c73df6d8194dffe7f8348a42f8a1752e03ea7de878f3ab483",
						"linux-arm64":  "865ae218aaf6cda1ad083e2ea3690292f471a09233c240e86e3e93bbcc7299b1",
					},
				}),
			}),
			Assets: config.ToLocatorWithResolverConfigs([]config.LocatorWithResolverConfig{
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-format-asset-ptimports:ptimports-asset:1.61.0",
						Checksums: map[string]string{
							"darwin-amd64": "8112c28051db3bb84824a0b8404ca731e3ba8d1ef9cc197bc865de1311cdd923",
							"darwin-arm64": "bf38e5f243671ae026d821c2844237ae51d324503b8c5a1962d393afceb7a50a",
							"linux-amd64":  "b108e0c2d403f3b757068341fd962a24b572fe7ae0550631dfabf39877665d2d",
							"linux-arm64":  "24769cbb1010f2c5775b3d3987cab54eb838890f1e6813aa616567a6aff24063",
						},
					}),
				},
			}),
		},
		{
			LocatorWithResolverConfig: config.ToLocatorWithResolverConfig(config.LocatorWithResolverConfig{
				Locator: config.ToLocatorConfig(config.LocatorConfig{
					ID: "com.palantir.godel-goland-plugin:goland-plugin:1.54.0",
					Checksums: map[string]string{
						"darwin-amd64": "e82d7bc30f698eb601b039abb1fa451bc3dfc3c074a2bda893aadf980508fbcf",
						"darwin-arm64": "3f7ac9bf68e4de72072ae42eca353248066db19846a0d45375ef452a116262a7",
						"linux-amd64":  "f64700e213e3770fd61cccfe5d0076d797d82f6a0f44f3736bff89c89a09adfc",
						"linux-arm64":  "5d5b0d92de2b68b5f101a3416aed21b990d5843d1cf3e24ea1872a93dbcd74e2",
					},
				}),
			}),
		},
		{
			LocatorWithResolverConfig: config.ToLocatorWithResolverConfig(config.LocatorWithResolverConfig{
				Locator: config.ToLocatorConfig(config.LocatorConfig{
					ID: "com.palantir.godel-golangci-lint-plugin:golangci-lint-plugin:0.18.0",
					Checksums: map[string]string{
						"darwin-amd64": "11ed1bff4b34bfe383d1b3c9ac8afd7fe43ba45053335c62b67de22489ef0c49",
						"darwin-arm64": "96f4d00b3f3c3c480e5e0dd62c890f4809e8b5f4cb0f15a530713847b9a11994",
						"linux-amd64":  "a8ed430c45f5d2d7493495377c8faae5378f2afb2f6c3b7c2e8a60ccfd361838",
						"linux-arm64":  "f007a95a8d9332d0cbdf376fd10fff23b973921a6accac7173279af9118e170f",
					},
				}),
			}),
			Assets: config.ToLocatorWithResolverConfigs([]config.LocatorWithResolverConfig{
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.golangci-lint-palantir:golangci-lint-palantir:0.17.0",
						Checksums: map[string]string{
							"darwin-amd64": "8e4420871cb96e179a844dfa97d13e8cc913d7e6783361992729a2ed2a86bd44",
							"darwin-arm64": "ab7614c493520369e327e9b0a95033947c05681c3da1c4273425cbb472192a50",
							"linux-amd64":  "929ca574b5143c3fa6bda7372e1896139674549f301aff77e182cb47b3519666",
							"linux-arm64":  "daab6f1fb8a6f06401d54d9ef46e38830cb3d3feb16a7000d23cb2f22850af9c",
						},
					}),
				},
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.golangci-lint-palantir:golangci-lint-palantir-config:0.17.0",
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
					ID: "com.palantir.godel-license-plugin:license-plugin:1.56.0",
					Checksums: map[string]string{
						"darwin-amd64": "a986de6123706392dae9114921b32189f78811be2abd07219a1c378323d901df",
						"darwin-arm64": "0f9d6ffa72e2148d7fcd84ffece619d8e7ce5552a7174572dfdbba2f682f3dc3",
						"linux-amd64":  "8e0ccc82bfecd074b936dca4c5b2209125871185a7ba04a8c566689eea6d043c",
						"linux-arm64":  "9b9df26ce2c5d9943af3150a7582918f4e8117d423a72910774493870e0b29a5",
					},
				}),
			}),
		},
		{
			LocatorWithResolverConfig: config.ToLocatorWithResolverConfig(config.LocatorWithResolverConfig{
				Locator: config.ToLocatorConfig(config.LocatorConfig{
					ID: "com.palantir.godel-test-plugin:test-plugin:1.56.0",
					Checksums: map[string]string{
						"darwin-amd64": "6f1f36251080ba076d0c55ec96320a8852953912d8785e1ff6df110e7ff501a8",
						"darwin-arm64": "60c2bcfe340f15d9b64d048c9dac08d1c129e98d0f8b4b338af7c8bcefc91dc1",
						"linux-amd64":  "28b1acfa040433880b4dca7faabf94168f8a38e073b556e7d0d80993e3f40110",
						"linux-arm64":  "1ba0ed0da8194e7fdfe87336dad3a408d61ec0b0d54fc5a47816c5541a47a497",
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
