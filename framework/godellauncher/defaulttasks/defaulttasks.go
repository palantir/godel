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
					ID: "com.palantir.distgo:dist-plugin:1.113.0",
					Checksums: map[string]string{
						"darwin-amd64": "ebfb146976133e3cd866122ca2d6c921218143ea46eec6f0a994af8dc0e17bbb",
						"darwin-arm64": "6e7d7846095bae2136022746bbda11508cfd82aef496372fc3891dac69067203",
						"linux-amd64":  "a2ec9bd51b4377eb3d170d830a5d1bc1bfc83c9fea109bafb8336185431095f8",
						"linux-arm64":  "a643061225f178f7de5a56b53fe63f4f5437da74cf057634787d08aced71e3a6",
					},
				}),
			}),
		},
		{
			LocatorWithResolverConfig: config.ToLocatorWithResolverConfig(config.LocatorWithResolverConfig{
				Locator: config.ToLocatorConfig(config.LocatorConfig{
					ID: "com.palantir.godel-format-plugin:format-plugin:1.62.0",
					Checksums: map[string]string{
						"darwin-amd64": "959d380220cae3aea8dddf73651a4c9a612fac2d1d2a90dd0e2b5615bc2be096",
						"darwin-arm64": "ad8b823d149b20f3fa6e733a5dae34fe750390bbaba01d3ce9292497e73f69df",
						"linux-amd64":  "19877f2eaf9a7a920a6cce9f3708f7272556cdabb1bdaa62512251744d083396",
						"linux-arm64":  "6fa8a924eab76b21c4a4c64c3349301059af94b6285c8a7990815050138dab1e",
					},
				}),
			}),
			Assets: config.ToLocatorWithResolverConfigs([]config.LocatorWithResolverConfig{
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-format-asset-ptimports:ptimports-asset:1.63.0",
						Checksums: map[string]string{
							"darwin-amd64": "d28b687b3df8acb33e6db5c3405056768181db12f3f47ed02eebff7f66b24561",
							"darwin-arm64": "0d6c16260afe631e8c64f103156acc9a2ab27931a3d494df3fcfbad22263ca8b",
							"linux-amd64":  "aabc460c7e944f49c8f047ca07dc8b46120d1dfc899b1d67ff91404d4d75345b",
							"linux-arm64":  "faa9f6a79b883d4d07f6168bd84766bfa2cb6c857cc87618519e4153d07eb721",
						},
					}),
				},
			}),
		},
		{
			LocatorWithResolverConfig: config.ToLocatorWithResolverConfig(config.LocatorWithResolverConfig{
				Locator: config.ToLocatorConfig(config.LocatorConfig{
					ID: "com.palantir.godel-goland-plugin:goland-plugin:1.56.0",
					Checksums: map[string]string{
						"darwin-amd64": "0e8526566b0b7f4581778767759741c2e71e43031b692e24750bd57e60433f9b",
						"darwin-arm64": "eb8cc3acac6d6c24a56901a7e63eb87209c16196411fb903ca084428f287c5c1",
						"linux-amd64":  "a12db1beca35ff40a7becc7847c0c10ecffbc221366b088d9aeeae37491142e0",
						"linux-arm64":  "6753b3880759cc450c75786b298652087538c244ef0df82e35e5fd8913e4ee57",
					},
				}),
			}),
		},
		{
			LocatorWithResolverConfig: config.ToLocatorWithResolverConfig(config.LocatorWithResolverConfig{
				Locator: config.ToLocatorConfig(config.LocatorConfig{
					ID: "com.palantir.godel-golangci-lint-plugin:golangci-lint-plugin:0.19.0",
					Checksums: map[string]string{
						"darwin-amd64": "8ce7ea6e97ebc3111ad2eefddae578187082a6750ed36ca74836d5c4abf9b298",
						"darwin-arm64": "a5e7007374d7ba15730660c13c2764544bacab74a1b7782a75d10416edaea2dc",
						"linux-amd64":  "52053648078705c97bfd11d1014342a49548d32834cf2b2c97f963573bde3e43",
						"linux-arm64":  "b816128f4f903e8750ef18345a47149e4fd3c32c4c3a56a89ee60665c606127f",
					},
				}),
			}),
			Assets: config.ToLocatorWithResolverConfigs([]config.LocatorWithResolverConfig{
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.golangci-lint-palantir:golangci-lint-palantir:0.18.0",
						Checksums: map[string]string{
							"darwin-amd64": "ba679089f2d67ff40928b47fa783b180db04c118169d555893e007c073b81116",
							"darwin-arm64": "7d6f6ca3bd72c69608397e3ac5766e6ab439b3cb938c292b5df6091edf2eacf1",
							"linux-amd64":  "f49e271748a40a309b6c6970a3ba8fc7be32cb01f949b017ebc21b1721d73315",
							"linux-arm64":  "2eec38e0806f61c87b2615890ffae14e87cfe8dfe2b14e12fab9aed22e16b8f7",
						},
					}),
				},
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.golangci-lint-palantir:golangci-lint-palantir-config:0.18.0",
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
					ID: "com.palantir.godel-license-plugin:license-plugin:1.57.0",
					Checksums: map[string]string{
						"darwin-amd64": "f3432678fcb1546f6339bf934634c02cd90564e4f9e5fb378b086740496e5474",
						"darwin-arm64": "76fc8588411179708d6c3cb4739fc68d60fe3d801f624cf8d051ccaba0dbbb6d",
						"linux-amd64":  "c7ef04f0b51141fdba79b53aceb04beb81d7a5238767491bb7a2b165fcd6839f",
						"linux-arm64":  "d8a24b05f3451f424ebcb47bafff1a322a541c30a45678b9cb3723cc8b102c69",
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
