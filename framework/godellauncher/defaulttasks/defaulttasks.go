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
					ID: "com.palantir.godel-golangci-lint-plugin:golangci-lint-plugin:0.17.0",
					Checksums: map[string]string{
						"darwin-amd64": "e7ea386d0d18c884504c9a11c969b108073b70c4ba5a2c7b8846db0e2467fa5f",
						"darwin-arm64": "6799952ea20b6697cf880fa4e1c25ce7ecd6a8658995403fb774c9cf5e19721b",
						"linux-amd64":  "850bb89422de320a17fa826290430edd4abecf031eaf03599e3c703a1784e481",
						"linux-arm64":  "f54d84c4800b175db96d28948f3ba02ca5efbeceb3f67e36de08d20e5ea6525f",
					},
				}),
			}),
			Assets: config.ToLocatorWithResolverConfigs([]config.LocatorWithResolverConfig{
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.golangci-lint-palantir:golangci-lint-palantir:0.16.0",
						Checksums: map[string]string{
							"darwin-amd64": "8dd8b2b1d9566a640ed71f9037ad810712aa54f35b5ea1c7238f85f0cef2edd6",
							"darwin-arm64": "771c5b10164df868a25505692ff41c2a876bff14797dafa520c43ff23b54b4c7",
							"linux-amd64":  "a30f7f7e513444717265720dff8eefa5a56b3d6d7ee27ace88453448f19953bc",
							"linux-arm64":  "110a9e40d0741767ab8d25f5be48df01cd731463712314fce6fda2e04a6d2408",
						},
					}),
				},
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.golangci-lint-palantir:golangci-lint-palantir-config:0.16.0",
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
