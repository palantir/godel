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
					ID: "com.palantir.distgo:dist-plugin:1.34.0",
					Checksums: map[string]string{
						"darwin-amd64": "9aebeb0e322f0e3c429cf2563761c6c6032217c07d8c422a8f5e75ef21146647",
						"darwin-arm64": "e2bbbd4275e2d780c245c57d6f29c5f7d6f468b151012f431a87a27500a220ae",
						"linux-amd64":  "6da4d0ff179fdd26915b87cb451263041a6c3a489a5862f5fbf1549acb0005b5",
						"linux-arm64":  "9755ee89635ad9210b91f37c552f67f69403e6bc04df15e6c095a81e98eff924",
					},
				}),
			}),
		},
		{
			LocatorWithResolverConfig: config.ToLocatorWithResolverConfig(config.LocatorWithResolverConfig{
				Locator: config.ToLocatorConfig(config.LocatorConfig{
					ID: "com.palantir.godel-format-plugin:format-plugin:1.11.0",
					Checksums: map[string]string{
						"darwin-amd64": "2da0528eb2a906a0c65225a3cda967758ab32e6837615e7519cecc63e3cbdb8b",
						"darwin-arm64": "b3085a791f511b29251242ca40391255bbed54246f0096c58b33877ab11fc8a4",
						"linux-amd64":  "750cf16ddf757d8fb1649107b06199c50d7d694d3b7e2d101b1e356b5a7fab88",
						"linux-arm64":  "4616df0fe80aa6e3a373e44195dfc63eb05883d30f8683c976e76cbe1f7dcdfb",
					},
				}),
			}),
			Assets: config.ToLocatorWithResolverConfigs([]config.LocatorWithResolverConfig{
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-format-asset-ptimports:ptimports-asset:1.11.0",
						Checksums: map[string]string{
							"darwin-amd64": "e1ba7ee0a799aa6c29485e47ba2a3a5d4ff71c65cfc8f2a7f9d627df9852c3e4",
							"darwin-arm64": "c54ebd0238c84bff6def7c24875674406a391fcfeda238ff4cd44fd2dfac19a8",
							"linux-amd64":  "a0f814cd0656322f54152d7d0838539b502c9b2c1a5431a193ac8af38fce1954",
							"linux-arm64":  "5ef830fce0ace03bdf2da3661a0921c1d5d8cdda66bcabe284d248c38d28c08f",
						},
					}),
				},
			}),
		},
		{
			LocatorWithResolverConfig: config.ToLocatorWithResolverConfig(config.LocatorWithResolverConfig{
				Locator: config.ToLocatorConfig(config.LocatorConfig{
					ID: "com.palantir.godel-goland-plugin:goland-plugin:1.6.0",
					Checksums: map[string]string{
						"darwin-amd64": "2501993b871903c8a1ee995375e8a7d16da78ac7c7840305bb0368d9d65e429b",
						"darwin-arm64": "ffaff49f8fca4e62f2a369c1c643f9d2705a6ba75b2a5ddcd185d1f421428712",
						"linux-amd64":  "349d96dd7fa3c34e04e31121839504974d493f6a433c657f5818b28bb543e661",
						"linux-arm64":  "11bd067ee4b6fef19f840e10ca3a77022777d18610044feed2b43009d1e243f7",
					},
				}),
			}),
		},
		{
			LocatorWithResolverConfig: config.ToLocatorWithResolverConfig(config.LocatorWithResolverConfig{
				Locator: config.ToLocatorConfig(config.LocatorConfig{
					ID: "com.palantir.okgo:check-plugin:1.16.0",
					Checksums: map[string]string{
						"darwin-amd64": "c57a11f5c60143ce63aaa301e5fead36d3abd3f1c9a0771b56149aa26d2f5991",
						"darwin-arm64": "226764180ed323748b92ec68943df89e883fd867ab7f1a92e0af30b1bdb54a1b",
						"linux-amd64":  "c27bc8645525a76a2d6e3b46522a7bfbb0f8907092a3bccd6def835ab331a5ce",
						"linux-arm64":  "d728ed40b7242ec8168c8fa8ec06ef41d9d37802125dea28760499e0edf67c2f",
					},
				}),
			}),
			Assets: config.ToLocatorWithResolverConfigs([]config.LocatorWithResolverConfig{
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-compiles:compiles-asset:1.13.0",
						Checksums: map[string]string{
							"darwin-amd64": "3bd03f96c4ee8338ca9fc2d254d37b571da7555c3a840352d7ff192958cf20d4",
							"darwin-arm64": "c378de4befe67cf0c235c9e48d71fcd183244c7b09938b2f5f54fcd2f30b3a7a",
							"linux-amd64":  "41f7ea0c15c36ffdcab7aa1424ddabef2487b95c896b759ad515870c69b9cab9",
							"linux-arm64":  "a85076666f835dc3aa9b3676f0a6fe9f7309b5a938819a5394bbcd7c5dfc94cf",
						},
					}),
				},
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-deadcode:deadcode-asset:1.10.0",
						Checksums: map[string]string{
							"darwin-amd64": "286058f26f238bbfbe71987f21475734b22671f65019d1235530e9e7954a50d5",
							"darwin-arm64": "6244a8513f6d079649cd71616dcf9a634d6c2265e85d160b4cfdc1f43e2d0aea",
							"linux-amd64":  "e81bd8ba28a559162ccaa51eadd99e5635e2c1d3e26438f07b1a804f86b5909e",
							"linux-arm64":  "853d31720eab730a373f77fb84c5715bd4770bbe4db22fde8257cd10a63ded51",
						},
					}),
				},
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-errcheck:errcheck-asset:1.12.0",
						Checksums: map[string]string{
							"darwin-amd64": "f8297e2d93f9ce7b6b9b9a373d00bdb8b038ea7b45f93415c158c1749892c7e4",
							"darwin-arm64": "8e461d49ba6be74d50fa0c0a1a74cc07d01f7db84900e047ecf32b07823db474",
							"linux-amd64":  "621d6a2922644d2a114268482f2766ac1c3d282be26212daa10a28678a6cf220",
							"linux-arm64":  "0d583cf1775941ef6fc68d181a3b966a0d921d07cb1b2b2a33bac6e9c6e90e8c",
						},
					}),
				},
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-golint:golint-asset:1.8.0",
						Checksums: map[string]string{
							"darwin-amd64": "64e66ff7ab9208a4048e99ee4dbcafd914eb448faa356f3e4d87e95329a413c8",
							"darwin-arm64": "9181af1ff2d85fe426bc4d44250d3683d6526b81184aa4515dbb22dc9b175e00",
							"linux-amd64":  "1f7cf26e019e75719c90b724811a1887e610438b5f00c244b05590e8d1547086",
							"linux-arm64":  "73541e680d92cb3ea75ef0f8f07962d7a124b6704ff345ee070146eb6a7d9bde",
						},
					}),
				},
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-govet:govet-asset:1.8.0",
						Checksums: map[string]string{
							"darwin-amd64": "d9be0574228c79d4854724be4c9a8766545ee9fd8acba318bfebae0f588d32f5",
							"darwin-arm64": "70625d164e3a25c806de9a251aad0f07ca34161b998293afd814b83fc0c43826",
							"linux-amd64":  "515678dfe23b6ccfa492e1413413e5a8c79b86b9665a7e4d7a85e05f651e5e18",
							"linux-arm64":  "8de47f4cfa5b38474b47fd686324a89bebb28775bd0863a852803088ae1f83b2",
						},
					}),
				},
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-importalias:importalias-asset:1.8.0",
						Checksums: map[string]string{
							"darwin-amd64": "e4cf3fd52bb8da2a51f05785b24a6ce34dbbf69136c59dfad77b4feddd8ca22d",
							"darwin-arm64": "c51beb01f83702484670753ea501d72dc1662dbdb455c5bb58e5c7eaaf49e9d1",
							"linux-amd64":  "2431d5bbf19fa8c252e258951b7fcbc23477cfd8a16d1e119bd4139afd4a9e6e",
							"linux-arm64":  "75f64f469a23abecd27a60baa1b746c2264b414667346e0f7f7773a8b2c1144b",
						},
					}),
				},
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-ineffassign:ineffassign-asset:1.8.0",
						Checksums: map[string]string{
							"darwin-amd64": "90deabba319cd56a944c363af3f51b93b9f32c4d326c4ec3745821bc3209a8dc",
							"darwin-arm64": "dc9118ed8b80fd74c7d49d91b794a959d798042eaed586183ac519b04c5822c0",
							"linux-amd64":  "d166e99a6b06eaaad98827cfeb40b6040bb0e72b169ebb8fc6218018af1589a8",
							"linux-arm64":  "51eaa52a66d5a4017e573853ec16d2a109b9c54cf9edcf99da9ba406d2d82db6",
						},
					}),
				},
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-outparamcheck:outparamcheck-asset:1.12.0",
						Checksums: map[string]string{
							"darwin-amd64": "f44abb668d28190ede5371502552c1662e0eed5f1b18102c80b168ceb135c100",
							"darwin-arm64": "d5000ec5abce437e4c148e5f98f82714c17780ea21328b9a8b708d1808865791",
							"linux-amd64":  "5cb79ae9491cbb511cd86fbe89adf8a2ab01c628548e1398aa7103626d5decc9",
							"linux-arm64":  "3ee42881b02bf64b87cf8ca1529f46eb87866db5c7c267c84bad2e835b64c525",
						},
					}),
				},
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-unconvert:unconvert-asset:1.11.0",
						Checksums: map[string]string{
							"darwin-amd64": "331b172e72ad898a2618b12c812d8a186cbf0742134c04a4d85148cfdfbf9571",
							"darwin-arm64": "cd6131ba41366c32e3170fa2db949edd2366cd4593c9273cddb932715445a32d",
							"linux-amd64":  "58a8af003b4bea9e2151e9233fd4ffb4081bef503c4206d5df22be0d01689aa0",
							"linux-arm64":  "02e0f6dd449042128371761c1ce59b066a3cbb3fca1717072e41fbc4af523eeb",
						},
					}),
				},
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-varcheck:varcheck-asset:1.10.0",
						Checksums: map[string]string{
							"darwin-amd64": "60185cde0327f22c09a945db0fb3d96325741f30e59538bf5679ae0e3bab25da",
							"darwin-arm64": "4145d4e2c6459edbffbf6caa6814322950e487dc64a5053a3b2d65f3b9a6dbef",
							"linux-amd64":  "5bce749c402864849a00b17765607503368e0c217c2afa7f556a1ff83ba7df1a",
							"linux-arm64":  "b6fec13fc753f504d9389b94ae5f6a6ff48abc3279f04da1d6f7d44867648484",
						},
					}),
				},
			}),
		},
		{
			LocatorWithResolverConfig: config.ToLocatorWithResolverConfig(config.LocatorWithResolverConfig{
				Locator: config.ToLocatorConfig(config.LocatorConfig{
					ID: "com.palantir.godel-license-plugin:license-plugin:1.9.0",
					Checksums: map[string]string{
						"darwin-amd64": "81a1f9d70263b5b334da2eedd8935dcca8522f6d4bcf3a4e512273a3479a98d3",
						"darwin-arm64": "d2ceebad31952aae1f697e1919a4f228d300e8b0eed283e8b95b5ca66b03c53b",
						"linux-amd64":  "3135a8843b4b5af8908f2903ebf2b6e91fb4833a646790f0e2bbcd337a56bede",
						"linux-arm64":  "1f5db134062b805d209c503f0f46dac9787b41210e9324f5e48cd1421a586881",
					},
				}),
			}),
		},
		{
			LocatorWithResolverConfig: config.ToLocatorWithResolverConfig(config.LocatorWithResolverConfig{
				Locator: config.ToLocatorConfig(config.LocatorConfig{
					ID: "com.palantir.godel-test-plugin:test-plugin:1.11.0",
					Checksums: map[string]string{
						"darwin-amd64": "b5d8d235bed77e8f2d63ca9ca89d38d4d2e7d08187a628f7d75fc1c11f4ac1c8",
						"darwin-arm64": "e0a8d1e6fba6ce5ab2d0680fe5f10aba2c2dae150362829286a14e42f5237646",
						"linux-amd64":  "4ec565984c233f535d1a95341c2c6743cd91ce5edf9a32556a8f79aa08054ccf",
						"linux-arm64":  "08ddf0348f5e3870bf0b39536c026964d816420e70a30ff4a21b052efce3dd5b",
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
