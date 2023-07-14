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
					ID: "com.palantir.distgo:dist-plugin:1.59.0",
					Checksums: map[string]string{
						"darwin-amd64": "f0c937d2957c6b2bd1d9b1672044a51fbe1c8aa60ff693869193406cace856cf",
						"darwin-arm64": "5c35c8fcc8a8673a2f7c8c6865b2598a896502ef06f5b5fd3ef8829d748ad54e",
						"linux-amd64":  "8384b242bf8be81bc228f5c4ced273d80fc31d3d5cf08a5e16f662f62f2d48cb",
						"linux-arm64":  "dbdf9afcfda13222365d82d41dead6ea246c97b8990d1548fd2103b7190f969a",
					},
				}),
			}),
		},
		{
			LocatorWithResolverConfig: config.ToLocatorWithResolverConfig(config.LocatorWithResolverConfig{
				Locator: config.ToLocatorConfig(config.LocatorConfig{
					ID: "com.palantir.godel-format-plugin:format-plugin:1.33.0",
					Checksums: map[string]string{
						"darwin-amd64": "4fb7c39eade1272d9004c7e61041a0539258a5654c4af2e89c4a47abd5d714bc",
						"darwin-arm64": "bd684c2bff7a2ceaf128f0310b2546dd5038e6e6089853b380ed4314ae2f59d1",
						"linux-amd64":  "a13a1b5e95dee81bd3bdbd40dc1a836b9d7d104fa94bc434ade4a96446de30b6",
						"linux-arm64":  "7f4bf807cf44d5755fcd6fa9271f7db7472f3741d21f23df2b6f4f5936eacd25",
					},
				}),
			}),
			Assets: config.ToLocatorWithResolverConfigs([]config.LocatorWithResolverConfig{
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-format-asset-ptimports:ptimports-asset:1.32.0",
						Checksums: map[string]string{
							"darwin-amd64": "3df9ced9e5cf3fa8a0a6179dfe81c56df0cef3cb537fa65d622c2ab1ec961b39",
							"darwin-arm64": "c4a3b7f515a1695126312d8299a84ecc487f08c4dbdfc05ad58e4a3934f9d33e",
							"linux-amd64":  "cd34aef7616ccba90c4fb3982500a09454f5300efa2c1f8d6630027dc190a1f0",
							"linux-arm64":  "99d24725332d7650f147cc8ada02768355c5b6a86d9364b7b5cacc2e3af75abc",
						},
					}),
				},
			}),
		},
		{
			LocatorWithResolverConfig: config.ToLocatorWithResolverConfig(config.LocatorWithResolverConfig{
				Locator: config.ToLocatorConfig(config.LocatorConfig{
					ID: "com.palantir.godel-goland-plugin:goland-plugin:1.29.0",
					Checksums: map[string]string{
						"darwin-amd64": "e6d51b4ba28f3c7f6c015296f0db29e7b33dfc2d311ec3be071fea21516e4841",
						"darwin-arm64": "5b96a528a8a47e84da54eb64f470f9c92fe1cd093a908712bb19e84d63db3588",
						"linux-amd64":  "f30829903fd9fa996352ac7924cef2f3ab10822423c4ea686ff569e64972dad2",
						"linux-arm64":  "d3fb3b9370c2715eda6d32218ab4135a9b98734d5d1081e907a7fecbff97db25",
					},
				}),
			}),
		},
		{
			LocatorWithResolverConfig: config.ToLocatorWithResolverConfig(config.LocatorWithResolverConfig{
				Locator: config.ToLocatorConfig(config.LocatorConfig{
					ID: "com.palantir.okgo:check-plugin:1.38.0",
					Checksums: map[string]string{
						"darwin-amd64": "97767341aad49c8854874ae6f93532f2f14420979ec82fed5a65379fab068a29",
						"darwin-arm64": "e911f00e2f1b21749b7ad29ad46577f7d1f4ab73e220c51798ebd7bbf21d27ba",
						"linux-amd64":  "e8e10a5c7dd9ac17922513b9b9fa0f63e2d0d4686126166b3ad7822d7bd9c767",
						"linux-arm64":  "c8870bcc6e2bff501a3e9e0d2f0dd721255ff2a42c09abcccfb774122d4e05ce",
					},
				}),
			}),
			Assets: config.ToLocatorWithResolverConfigs([]config.LocatorWithResolverConfig{
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-compiles:compiles-asset:1.37.0",
						Checksums: map[string]string{
							"darwin-amd64": "7d0838309244422dbfd53eb308cf1ecd11047d492de9f7e38161c92e60030002",
							"darwin-arm64": "f60f571e17cfe47f546dc825b805b0f740ce66b7840cf900a3f55f8127f787c3",
							"linux-amd64":  "6a82aa0b225184867cf4ec0f128e4af0ea89c14851a5844428d4960b2767c8d5",
							"linux-arm64":  "4ba92767ba24c93c35586283acd6f32e26eed067dc8f838508debab95894592d",
						},
					}),
				},
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-deadcode:deadcode-asset:1.34.0",
						Checksums: map[string]string{
							"darwin-amd64": "37aa2a1ad1adfea0618c55e3c1cf8ae2d9c30c12f63c4970ff36168c53d6a712",
							"darwin-arm64": "ab48a15ee8d53579afeb51973cf527fe22dc01cc41f364b1d97d7541bc2df871",
							"linux-amd64":  "dc58f28488973e23694c5285045f68c97cc8d2ad6900766380f709cec65e2aac",
							"linux-arm64":  "a53f8cf95e5305c890a79e72f2252d49926ff76704982dc389a357887bb06b20",
						},
					}),
				},
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-errcheck:errcheck-asset:1.35.0",
						Checksums: map[string]string{
							"darwin-amd64": "ea93eb92599eebd5056a0dc9cda7360a08900f251de88c5029828a1a357f1098",
							"darwin-arm64": "5f592b89684b98e2e386597375859bf44f390aaac2fd682ae6d8394be4f201d9",
							"linux-amd64":  "90dbce1be1e264482ff758943b6b0f1161f68d346e8b12b7eb2a5f5f6cf751a3",
							"linux-arm64":  "d9be9860fb722f017603616468233f93ad4f64dfb8495b40cf2bc94e07620462",
						},
					}),
				},
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-golint:golint-asset:1.26.0",
						Checksums: map[string]string{
							"darwin-amd64": "5caa6aeedc8a0ed12078f9f97d0ec96716ce0300b2536dd3972357676fa198db",
							"darwin-arm64": "13a95bb657001d5029388f4a95e10dc5343122a21cec1564c37cfa112f14119b",
							"linux-amd64":  "9749b3727b1147e4ae3836eccaf2df777bdcb15f8da01a2af2729379c63b1923",
							"linux-arm64":  "c27c1b40b46666f9b698fddb970acddbfb2c044b99405bf6ee9af6b323734882",
						},
					}),
				},
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-govet:govet-asset:1.30.0",
						Checksums: map[string]string{
							"darwin-amd64": "a98225ddc2064215b97074e2a2f77a70d1d0ef6b38eca348fa000797a6fe882e",
							"darwin-arm64": "9812261e8268d7567ac4aeb81cc7c611fdd1fe1df4285536b2f3097db244a355",
							"linux-amd64":  "0d341b3bbca56f908ba448a341599857a6234b0465b6cbfdad5fca9a259f62d0",
							"linux-arm64":  "423e6444c7885c7e48988702e87c328ac70c95c0917032b738730300a908d5c8",
						},
					}),
				},
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-importalias:importalias-asset:1.29.0",
						Checksums: map[string]string{
							"darwin-amd64": "af07aacc83db5ac65f234190ec0f75c92762b3ea5df15fbf298bcb2f63817d7e",
							"darwin-arm64": "4114ecee571eb3c5aabd1f03e38d4c69e319f1ae11fe676702e4293b4463fbe9",
							"linux-amd64":  "819696a0132a69c66434662cb66f41d1cf36e40abbd71b62a4b29ffd133de22f",
							"linux-arm64":  "b3b484bc15d69893bb01d7b8a56458055619a365b8f2d27001cdd73a0f4d704c",
						},
					}),
				},
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-ineffassign:ineffassign-asset:1.32.0",
						Checksums: map[string]string{
							"darwin-amd64": "b66aed6d626887edb89a660d0618b769f35c936326f34de2c9ff47cd97c680d0",
							"darwin-arm64": "12c1f3ea04d2816dd0d899062798e78b88d2169dc9111fa44fdbeaad5c1de6a1",
							"linux-amd64":  "9a6c902e7d3a8c28d84c8cf72abe8e4391098f9a642ee0e6d1a60f01bbd533a0",
							"linux-arm64":  "0ee91389c5465050819294f2b4fbdb94ebf1df8df9c8c2f492eb3985006c41c7",
						},
					}),
				},
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-outparamcheck:outparamcheck-asset:1.34.0",
						Checksums: map[string]string{
							"darwin-amd64": "904b3bcf1abc090aaa9de24e1bb8efa1cf003ae3de75326baf2792c8759fdf82",
							"darwin-arm64": "18c05212fcf258f169188637059948709fc8a2b4a977e712f72dafcac4a364ba",
							"linux-amd64":  "6ad77dafaf62060620294b89a96cfaa9b0935d5012f66ae69fbd2c31c96023c8",
							"linux-arm64":  "804d8a1b8955235e04b4acf72921a2da81be512cfb04d94f3ff130f8cf279d3d",
						},
					}),
				},
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-unconvert:unconvert-asset:1.34.0",
						Checksums: map[string]string{
							"darwin-amd64": "c32421b4f4cdb670769d2067e03fa57619cfc5fd35f328805bebb4416fda8561",
							"darwin-arm64": "a71b2692ce0146d84e4d9b27e3e7a2e1a00335ab2da8f4959bcb7c4f7110b2d1",
							"linux-amd64":  "b22379d34f2cd17e045a349ec1a6fb8b95efd0ca1e87cb3ae430993d953274e9",
							"linux-arm64":  "1ebd3b63ea0b2699afe89283ab5cdba30818a9bca3f03f5a1ca270e0bda5a90f",
						},
					}),
				},
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-varcheck:varcheck-asset:1.34.0",
						Checksums: map[string]string{
							"darwin-amd64": "c2f67382c8eb2b0345176bc2fcbe5c560d7a8f8c6ebf80411bfaea784257ea0c",
							"darwin-arm64": "317c6e49daabe7126d530aa52ad1216e6ffbc5e6a6e1c3d27b10703a582eb4c2",
							"linux-amd64":  "f493ac2fac24b1f8f67cc4e603511fceca4b39c9269fa1d964db5412cbbc30f7",
							"linux-arm64":  "51f504c95d23d3cb7d23e363ff962ab31568823f30be428e96e7b339b9f4b5d2",
						},
					}),
				},
			}),
		},
		{
			LocatorWithResolverConfig: config.ToLocatorWithResolverConfig(config.LocatorWithResolverConfig{
				Locator: config.ToLocatorConfig(config.LocatorConfig{
					ID: "com.palantir.godel-license-plugin:license-plugin:1.32.0",
					Checksums: map[string]string{
						"darwin-amd64": "39d373a5eaf46cf91cd1e3907db0844052ac96549f49b71dfdaefc29ed1e0225",
						"darwin-arm64": "3339f1ff670666d1bd9294d3f7943afa922c49383d24b9e156b61e3413ce7ad6",
						"linux-amd64":  "4494c191b25a851c5a0cb6a5cc1285ad1a7c7ea5c45b2ee1bc4573250bd63e6b",
						"linux-arm64":  "aeb51fdcf6c323591a640816b29ba1f9587f5856bc74befb5365c79548f5a4a3",
					},
				}),
			}),
		},
		{
			LocatorWithResolverConfig: config.ToLocatorWithResolverConfig(config.LocatorWithResolverConfig{
				Locator: config.ToLocatorConfig(config.LocatorConfig{
					ID: "com.palantir.godel-test-plugin:test-plugin:1.31.0",
					Checksums: map[string]string{
						"darwin-amd64": "198494ec151cd5933b9f73468b8b65293db9ed33253c01dca426a07f5ca84999",
						"darwin-arm64": "34b2067c086a6e4396d4edb39afc5065fbd253ca4db44de654d98851c56365ad",
						"linux-amd64":  "ec34e19f4c7530f0dc6e7ef1df74d19beeb3964eb36b48d6cd54b7859ac0a2b6",
						"linux-arm64":  "2dc24810c44eb9eb80bad2f9d572929807f462ab70d29fe4bcbc0778d15d01c2",
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
