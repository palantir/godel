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
					ID: "com.palantir.distgo:dist-plugin:1.108.0",
					Checksums: map[string]string{
						"darwin-amd64": "9b4f45310264808b8d418dfb1a3f4c4069b0d82ba4e438078a5007dd8770df91",
						"darwin-arm64": "04339bc3937227abb1713e1f8522b5b2b85f77fbde3309897633c0f39be4b2d6",
						"linux-amd64":  "4b8fc88bfd9ce1d8bdfdf244816772f1cfb4e09a2502811bdc6d87af378b5608",
						"linux-arm64":  "414320b3870fb1c802cc1ec4f5101265c5c1f45c9e091a3d4adddd5acdd6563b",
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
						ID: "com.palantir.godel-format-asset-ptimports:ptimports-asset:1.59.0",
						Checksums: map[string]string{
							"darwin-amd64": "424c793f8aac1968dbd63ec02055cd27aa9941fe1a07bd87ced6fa546786ab5f",
							"darwin-arm64": "2450727fc27ac8692c2a796a3c0274362b9e3b70637d98c9b730923574b676b3",
							"linux-amd64":  "3a31bcae5cde3be115ed16f2f786716100684cdb070e2e60eb7bf02dcaf64694",
							"linux-arm64":  "cf62a5977bed49966ab9f3ddec668e447da261ec04d6c829d6df5c579545f294",
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
						ID: "com.palantir.golangci-lint-palantir:golangci-lint-palantir:0.13.0",
						Checksums: map[string]string{
							"darwin-amd64": "4c1cb95bfa6e27adb1a33423688002ecca7ed1a481ed34f7c19fcef5eda05518",
							"darwin-arm64": "d327786b4399db0cd6a1ff42667aaeab13244cdc84b32700646369e879461d73",
							"linux-amd64":  "9a192381f55786b064e6e27ecba6c9bec580e50d02bde305d83037222108b95b",
							"linux-arm64":  "f51ac0cd7b164d7dc7502e3582b566f974b68ee39f146885bf797b2a812e4118",
						},
					}),
				},
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.golangci-lint-palantir:golangci-lint-palantir-config:0.13.0",
						Checksums: map[string]string{
							"darwin-amd64": "",
							"darwin-arm64": "",
							"linux-amd64":  "",
							"linux-arm64":  "",
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
