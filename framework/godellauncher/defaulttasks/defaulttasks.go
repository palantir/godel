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

const defaultResolver = "https://github.com/{{index GroupParts 1}}/{{index GroupParts 2}}/releases/download/v{{Version}}/{{Product}}-{{Version}}-{{OS}}-{{Arch}}.tgz"

var defaultPluginsConfig = config.PluginsConfig{
	DefaultResolvers: []string{
		defaultResolver,
	},
	Plugins: config.ToSinglePluginConfigs([]config.SinglePluginConfig{
		{
			LocatorWithResolverConfig: config.ToLocatorWithResolverConfig(config.LocatorWithResolverConfig{
				Locator: config.ToLocatorConfig(config.LocatorConfig{
					ID: "com.palantir.distgo:dist-plugin:1.70.0",
					Checksums: map[string]string{
						"darwin-amd64": "1c6a043303e36bd6c8f5d276629f7de94ded07315822f898ae496d2a742e6ab1",
						"darwin-arm64": "fe4b451edb857d166113f8a235a7d4c63d986fb8b3e3f3c785837d612c367849",
						"linux-amd64":  "ad367c25a2c23270f7b96ff63d146336158ba4e9929003c6f0b0cb5fea4598dc",
						"linux-arm64":  "873f29cbbaea4349403d0ea3f3cc29a24063ff3eb0f5ee67b97ce0e34d4537aa",
					},
				}),
			}),
		},
		{
			LocatorWithResolverConfig: config.ToLocatorWithResolverConfig(config.LocatorWithResolverConfig{
				Locator: config.ToLocatorConfig(config.LocatorConfig{
					ID: "com.palantir.godel-format-plugin:format-plugin:1.42.0",
					Checksums: map[string]string{
						"darwin-amd64": "8a6f5c42b0ff8db52cfbc57389a20a4c9bebbd862923d3594e6f9595236228c1",
						"darwin-arm64": "1f8a3d926aea5c4f75d64cb024fb431ecfa5b14de33b0edb4dea59c31d9d82ab",
						"linux-amd64":  "ad02d84f4234bebe9efee5ed9f239462c9b370e11c61d887066aac16b0128fd2",
						"linux-arm64":  "c9ab9ce0a53432c33eb8b6ecdc8d8a84cac2b822abaf4e03991f179be4861d00",
					},
				}),
			}),
			Assets: config.ToLocatorWithResolverConfigs([]config.LocatorWithResolverConfig{
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-format-asset-ptimports:ptimports-asset:1.41.0",
						Checksums: map[string]string{
							"darwin-amd64": "d1dfa22c3d5ecf7d256e1b4daae7efa0037d144b5d659504e8ddbe837a4f96c8",
							"darwin-arm64": "4c1cb01e88c688b502c7f178f858c6517b226578d0a673ee9e28cbabbbc7afa4",
							"linux-amd64":  "5d5dcda7967126b22d63e0b900384532662de68655a090ee560f84b206ab073f",
							"linux-arm64":  "dc3352006de948694a3847cd90b6b28eabe5ba241523e832b1d2c5765dfeb8e1",
						},
					}),
				},
			}),
		},
		{
			LocatorWithResolverConfig: config.ToLocatorWithResolverConfig(config.LocatorWithResolverConfig{
				Locator: config.ToLocatorConfig(config.LocatorConfig{
					ID: "com.palantir.godel-goland-plugin:goland-plugin:1.38.0",
					Checksums: map[string]string{
						"darwin-amd64": "ce13efbfce7d0ecaee673ebdb6cdb43be114cae7c7c98c45ab53f070f053bbea",
						"darwin-arm64": "b7c3732c15796ca2894c7c73047fb573b5f65b4231b21d97e1654df2da33c51e",
						"linux-amd64":  "944ce2f4fead395ddcdc86baca986ee528097fcb3064c1210aa29297668168cf",
						"linux-arm64":  "4b16828ffccb796c6a597af3b902cadbf7bb13e8e011933a8e8b870356d8f47f",
					},
				}),
			}),
		},
		{
			LocatorWithResolverConfig: config.ToLocatorWithResolverConfig(config.LocatorWithResolverConfig{
				Locator: config.ToLocatorConfig(config.LocatorConfig{
					ID: "com.palantir.okgo:check-plugin:1.54.0",
					Checksums: map[string]string{
						"darwin-amd64": "0a424449cdd2a32e13eb24fd1c310da93407ae4ca112155e41fb2b163459f702",
						"darwin-arm64": "dd65b59a5285b12b910e785259e9da15ff0465727f860d824b8c64c15d127f21",
						"linux-amd64":  "4af8b81960f3930f34a1c1dcd788e1a72084b742f4b735af9abee22bb840d384",
						"linux-arm64":  "84b1fa7daf892c981179ef6ffd8b71a3ab02795c4b4049805dd4fc5d81869591",
					},
				}),
			}),
			Assets: config.ToLocatorWithResolverConfigs([]config.LocatorWithResolverConfig{
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-compiles:compiles-asset:1.47.0",
						Checksums: map[string]string{
							"darwin-amd64": "b4c431677bbee74cd5ea765010f222c0bc6be164e65131cda9e055eea2cc964d",
							"darwin-arm64": "41dc227e23df9574a52547e622141142bdac450d4e330b8bb3c6e6459032bbe7",
							"linux-amd64":  "c0041c43fd3545cf4a8b80fd4bcf5caceea0808706a13ced639cb7427170bab3",
							"linux-arm64":  "7f47fb00bedb3a4609d19b28e773807d6e901ff71157d5e9438d128e371a7ba3",
						},
					}),
				},
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-deadcode:deadcode-asset:1.43.0",
						Checksums: map[string]string{
							"darwin-amd64": "b2ff706012ddc080a59ff814a470b8e68773bde3eb12b4bec5b60e6e31c49926",
							"darwin-arm64": "58ce7986f9360c8c522e862f93b69440c6ee1d78843410a8733999b5d24ab981",
							"linux-amd64":  "cee551ce8fc314e1ac1dacef018e2e32b50d949268924429045819da161f7dd1",
							"linux-arm64":  "5b99a2567a83ee78bb0e909f9ebdbe39b7bde7c030edd44a7bc6573b79676030",
						},
					}),
				},
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-errcheck:errcheck-asset:1.45.0",
						Checksums: map[string]string{
							"darwin-amd64": "a7ff94cd449c5ac733cf3cba34479c52af123f8190ec8367919c3394ffbb0b0f",
							"darwin-arm64": "8a164e13260fa878f7a0dbffaf0c5de2c58ea4da8c7cedb9f57e6f12ae01288c",
							"linux-amd64":  "b1c9ae4a3005f4746c3b28e7dbaaded2106bb5c53e9a09708583c7fd34cdc3f2",
							"linux-arm64":  "f4208aa531adad2a5b035ec1c932797ae9aa8b69fbabb031cc03d42d638b1dbe",
						},
					}),
				},
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-golint:golint-asset:1.35.0",
						Checksums: map[string]string{
							"darwin-amd64": "28e9f90ed0feca152288a1ea0cde3e37abe61508f5889a2a99d48d9636ad2775",
							"darwin-arm64": "c34d45c514f3817117b5b6d72a15fd722d357595aa8faa444658c7a5155c9f6d",
							"linux-amd64":  "106392bdf6038160a04b49027696de28f0f0705f853814a42d8897d42337ac41",
							"linux-arm64":  "ceefa9f71fa2d3757f6c07d33ab58b0ab102d5df8d0a6e31a757e34dba044fd3",
						},
					}),
				},
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-govet:govet-asset:1.39.0",
						Checksums: map[string]string{
							"darwin-amd64": "bbc30142e8e7f3a50c7fd9ab9fec37e1b4f0fede194780dc31e0b99f98eb0be6",
							"darwin-arm64": "c30ef977acbcce03bbe28ddd65356099d8c4968acab9af43709cb08047197e77",
							"linux-amd64":  "105784cd6fdc1bad680c6e4573b00a2fe8baa5062301d95573d5e77ad15bc6a7",
							"linux-arm64":  "66a3d494e5b64865e1819b0c21c19543505012461172c7cd19344fc6ed96b7f8",
						},
					}),
				},
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-importalias:importalias-asset:1.38.0",
						Checksums: map[string]string{
							"darwin-amd64": "f9dc427f211422553959ab41b5cc30118ccbc77f87763791181f1ad069dff8c1",
							"darwin-arm64": "8e8031405ff1d277ad1d8173d3d2d11772944a32559c1856f136f4dab1eb76b7",
							"linux-amd64":  "0a72ba7f6c86b77149976bf44931c366e60daa450d76ef6ffc083a6975a3f1a1",
							"linux-arm64":  "f1926fb8edf6f8a91e598268b61d120aa14cb8378ba2ffecd9097553e647365c",
						},
					}),
				},
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-ineffassign:ineffassign-asset:1.41.0",
						Checksums: map[string]string{
							"darwin-amd64": "565c76eac920f4b25e5815410d537e520d929a61acf8f961419cb7f1905a00ca",
							"darwin-arm64": "310407f55de28142e780bde23dec748a7ad5b225bcda20a8899f8e2fe63d084e",
							"linux-amd64":  "f06800336614acce306110c0caa5391e641adff43d9cb98d1da7d9d047eb9878",
							"linux-arm64":  "031352197220d9c8ece866499b19518868eeed20cf1ef1cf355ba0f9b6e4cbb7",
						},
					}),
				},
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-outparamcheck:outparamcheck-asset:1.43.0",
						Checksums: map[string]string{
							"darwin-amd64": "a655a3d8e24a57be77f81fe2acc15a0b4d06ea47f99f16dbd364d88b5b03dc51",
							"darwin-arm64": "5ebddbf67f9766e50ef4119f0743a46f6f9a03d00608e2e20d066524b87b9aae",
							"linux-amd64":  "7e425c8aa7934f936efd8f7f4990f2e498bab5bd672b451b4e122e83ecb9439b",
							"linux-arm64":  "c48e729b31c9358e8a4741edc3f7f152841c9601424d9ebe85cd6fae24e1dd5a",
						},
					}),
				},
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-unconvert:unconvert-asset:1.43.0",
						Checksums: map[string]string{
							"darwin-amd64": "b3c113ff080f1cfd0362f05cb1e1e8d5deab3e5927299b83f18bc24434982462",
							"darwin-arm64": "9cad5c3b2b6251bc8ae072f9449860ab4379a8dec0221a873d55e0f544ca8694",
							"linux-amd64":  "a4e17a17b9d5d79f69a61127c4172f24fd2e94f91e77bc3501f575ee607f2aaf",
							"linux-arm64":  "55bd459d1a13589bf244613b9443efdface36f26d4402afb6cf58f985d1d38cc",
						},
					}),
				},
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-varcheck:varcheck-asset:1.43.0",
						Checksums: map[string]string{
							"darwin-amd64": "14aaeb4e3c7703f069c2bd09cb688d259250c24410a2338340c09ab2556adc39",
							"darwin-arm64": "895e984172fc2ec81f3bbe6972926ea9accbec259745cf82b2557297272e980a",
							"linux-amd64":  "375a52172f4a779f1e87a6ea658ad4a2fe169b42a3d9730c5c4cc58b02c5ae3a",
							"linux-arm64":  "0a1b0894e611d774fbae7af7c97324a28757982fe506b4019ca638930187632e",
						},
					}),
				},
			}),
		},
		{
			LocatorWithResolverConfig: config.ToLocatorWithResolverConfig(config.LocatorWithResolverConfig{
				Locator: config.ToLocatorConfig(config.LocatorConfig{
					ID: "com.palantir.godel-license-plugin:license-plugin:1.41.0",
					Checksums: map[string]string{
						"darwin-amd64": "aa280a42a88b38c29ffa446fb3cb46a8b77f508b28cb533460866785d61b8023",
						"darwin-arm64": "199ec09d2592a4fcedd6007dc64074fcaade91789ff1b8e6fe914a6de974395e",
						"linux-amd64":  "f746c89ca7a899ad7560a9cf6d42c088c2c9fcb012078c5c5e38af6fc019d12a",
						"linux-arm64":  "034f4e4e018960fbeaf58e69398cdf4069b8e2a4f191b74c1eb4bb6dac8a1502",
					},
				}),
			}),
		},
		{
			LocatorWithResolverConfig: config.ToLocatorWithResolverConfig(config.LocatorWithResolverConfig{
				Locator: config.ToLocatorConfig(config.LocatorConfig{
					ID: "com.palantir.godel-test-plugin:test-plugin:1.40.0",
					Checksums: map[string]string{
						"darwin-amd64": "a6651b144c7bf3d342566b223005224b3b16931596ebed233785a23b5ae6b88e",
						"darwin-arm64": "3a6b453a9b3aa46d20d619fc5b8591c718fb203c0d6c4128b88940f083bfef07",
						"linux-amd64":  "47321c14a7a8283d429072b8dc6798ba0e8d3549574baee56a7a27b177692783",
						"linux-arm64":  "a2f4d958066875b67f2d2ed505f410e883046795cf9359a91c01bc316acd2de5",
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
