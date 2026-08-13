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
					ID: "com.palantir.distgo:dist-plugin:1.111.0",
					Checksums: map[string]string{
						"darwin-amd64": "95261c895a01139bf1134801b0546a5dd99eb8738f4c42f402c98976acc4cd79",
						"darwin-arm64": "2712a4736d5d8b78e2420955a2021f2eb127cd94cffe21553c9041dd926654cd",
						"linux-amd64":  "9f2019625682fb895e7688e01bb75ba039248b25c48e09430c249732bb1e934d",
						"linux-arm64":  "c04adce2ce2604846e7fb2f920e90f34f1ddd2c696cb544c5bd81aa3c69ddad8",
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
						ID: "com.palantir.godel-format-asset-ptimports:ptimports-asset:1.60.0",
						Checksums: map[string]string{
							"darwin-amd64": "b422cafb714147bb7806a5b6d22d3fdfcd3679440578c4d7fc01805553d33ee6",
							"darwin-arm64": "e1f4fd787c47afa220b7a29184cbe12efbf34593ea43465c2871cd7f5d086307",
							"linux-amd64":  "baea9f1675817631ab92e2b717a677be34a2fa09218456608371f6b4aa940f86",
							"linux-arm64":  "eec0256fb407fa101398c1947e6df4835b45ebedb2095c6cbf662127d61f6234",
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
					ID: "com.palantir.godel-golangci-lint-plugin:golangci-lint-plugin:0.16.0",
					Checksums: map[string]string{
						"darwin-amd64": "b7acb7fe33c844374962e10a6d91d3ec6d271d2ad5520a6b38b3e2a5381aaf59",
						"darwin-arm64": "0f08b1e8df5108afc6ef25f28965c7f0a26fcf726babcefdda85c79b68e01989",
						"linux-amd64":  "53a6bfa27ce92267b117114439db3a6269f2634244d135fc2b4f20bb0b8d3cd8",
						"linux-arm64":  "0ab6c2d81b680ea2716d150799a86d43bff4dd883bf456be6bcebf4343b5416d",
					},
				}),
			}),
			Assets: config.ToLocatorWithResolverConfigs([]config.LocatorWithResolverConfig{
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.golangci-lint-palantir:golangci-lint-palantir:0.15.0",
						Checksums: map[string]string{
							"darwin-amd64": "5370487294c61f53c246872ad37e104d933cb22f7ddec0b576f46446704cc842",
							"darwin-arm64": "a365143bd5066ab0655f86069b90df9a4f3975aaabcc9eb5cad439697d18a4ed",
							"linux-amd64":  "b4dbecc8da635f1a2e5726e8e25e59958e0cce49e5ed5fa360149204a921f273",
							"linux-arm64":  "b8d8b172e4364fd02e5dc175fa871dce33ec2a6efa83d7a0504a421a225b757f",
						},
					}),
				},
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.golangci-lint-palantir:golangci-lint-palantir-config:0.15.0",
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
