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
					ID: "com.palantir.distgo:dist-plugin:1.61.0",
					Checksums: map[string]string{
						"darwin-amd64": "dbb2afc28fecd29d4a99f254281043664b16e58650f664d137bf2b6da23525eb",
						"darwin-arm64": "0bdd2de0f8862976c7b838b7862eb7b970b0ce1da47198d10b4069ed95af86bc",
						"linux-amd64":  "fc6a350d8ba57ef31b8b1eaf9e3f563d46c725a14d2b977bd378702a21e7b24f",
						"linux-arm64":  "64e0cd24350192577ebc3e92addaaff9ce95700ebdf81951402045b41cad8da1",
					},
				}),
			}),
		},
		{
			LocatorWithResolverConfig: config.ToLocatorWithResolverConfig(config.LocatorWithResolverConfig{
				Locator: config.ToLocatorConfig(config.LocatorConfig{
					ID: "com.palantir.godel-format-plugin:format-plugin:1.35.0",
					Checksums: map[string]string{
						"darwin-amd64": "7922741f475fd917402cd9baf6d85638a3007dcb88b1b7d7dd48f0b47a486b81",
						"darwin-arm64": "4b7527a452bbb44e01337306da93866ea4c5d70b290db43039cf834a3beffd21",
						"linux-amd64":  "f9a8d0dc067de76ec3e8ea2d4fcf1f5d0ec5f3d400172c8b4664284046bb8d00",
						"linux-arm64":  "c9fb3640492c43f93a4726f695423973d138da48a6cb91940ac5741919d6ded3",
					},
				}),
			}),
			Assets: config.ToLocatorWithResolverConfigs([]config.LocatorWithResolverConfig{
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-format-asset-ptimports:ptimports-asset:1.34.0",
						Checksums: map[string]string{
							"darwin-amd64": "6ba45995eb917b9b0bf628fdfc88ba4da71fd509946210a4a32df63a7e799c16",
							"darwin-arm64": "9ddd23c1637433dccee37f7600752f4f03714b7fdce0732fab7c29a0efa191b6",
							"linux-amd64":  "c9341f81cbfae683998be1e9a1f4967ca78a4889e2e8b1d4815ce9711ec4d73d",
							"linux-arm64":  "a694274bc29a5f0f501d3cbcdb563049d54895221ecf532f84046474a6d597bf",
						},
					}),
				},
			}),
		},
		{
			LocatorWithResolverConfig: config.ToLocatorWithResolverConfig(config.LocatorWithResolverConfig{
				Locator: config.ToLocatorConfig(config.LocatorConfig{
					ID: "com.palantir.godel-goland-plugin:goland-plugin:1.31.0",
					Checksums: map[string]string{
						"darwin-amd64": "60aac7dc77a8e658767c63e6d21f05c9def7e1b2196825577b5d8233d70c6166",
						"darwin-arm64": "c64edbd6438978652b84d624ec5e92919cf886f43c2465cf33663f7555e2bb16",
						"linux-amd64":  "58ae739469d342258c7e3c3706e0d2b83cecd5bca7fcd8c0056b1e0235b3846e",
						"linux-arm64":  "352f04dfb16458ecf5ce423e9454ed27ed575c21d30ac93a2638af454d6c7df7",
					},
				}),
			}),
		},
		{
			LocatorWithResolverConfig: config.ToLocatorWithResolverConfig(config.LocatorWithResolverConfig{
				Locator: config.ToLocatorConfig(config.LocatorConfig{
					ID: "com.palantir.okgo:check-plugin:1.40.0",
					Checksums: map[string]string{
						"darwin-amd64": "53739294fbdbee884c8e8c47b2823b08fc1b7feb35f5836ec8d7cff0ea448374",
						"darwin-arm64": "ee30b0f71e98120e42eee56efc8793638a57b12d2db81855f101d07c1bb7c298",
						"linux-amd64":  "099c7ae5538c024abed57fc59daafca1b493c01ec50880f3bd4b8d3f280e4bec",
						"linux-arm64":  "8e1d54a73e1146c0da3f0de100e34ee61f7bbfa2e1c7cfe1950cd03d022b58b0",
					},
				}),
			}),
			Assets: config.ToLocatorWithResolverConfigs([]config.LocatorWithResolverConfig{
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-compiles:compiles-asset:1.39.0",
						Checksums: map[string]string{
							"darwin-amd64": "d3503cf4d6042125427bec67cecd141e241b426a883f5d0bef3eeac6c1bc3a39",
							"darwin-arm64": "a18a9b67110bc308d150f6a89851710fbccf1c2f80804a9cc75fd82fbbbf1701",
							"linux-amd64":  "1fbf81561497face046eb45303a863a4466a628268bfc78d3e679dccd94b5c71",
							"linux-arm64":  "79446c11bcf59fba95118d981733a8d82da8d83e2e5eabaabfd73293c698edb5",
						},
					}),
				},
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-deadcode:deadcode-asset:1.36.0",
						Checksums: map[string]string{
							"darwin-amd64": "5c098d94195077c598dc812125fd3cdb0c49d96a389b1225f13a317018607e0f",
							"darwin-arm64": "122a79857547e5317f7efde4a8a53959655ca16b17e3bff3252764a8b77c4d37",
							"linux-amd64":  "2897d9c60db808a7c3791d3bcde22891a1e76eba802966784e79ac17fec6161f",
							"linux-arm64":  "7b7dad06d486484c846288626f61cc329c9c082dd86c3b5a5c8571f1c44ea7f7",
						},
					}),
				},
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-errcheck:errcheck-asset:1.37.0",
						Checksums: map[string]string{
							"darwin-amd64": "af81d9d714fe579874a7592e259f1824486bae14c34babf657829d83f932ffb2",
							"darwin-arm64": "8f61584d3aeefecbce1b307c34a59b035ca16e45fea88bece869e5693764ca9a",
							"linux-amd64":  "253e10d1745a019881e2483b08b328feca04bb2e088dde25ed40eb9005c9baba",
							"linux-arm64":  "041e0f5322e6bc646b0065cd34ac7c727d23611952f1de0fea201b83148edff0",
						},
					}),
				},
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-golint:golint-asset:1.28.0",
						Checksums: map[string]string{
							"darwin-amd64": "9b372ec74e3b98497ca1e221890ea0767f84e199a0845c8747840fd124380419",
							"darwin-arm64": "faa9da84e84f5fe227fc16e18353cd907504337fd7bf16d89282718971a12ce4",
							"linux-amd64":  "daaf51ffc89414a916dcc66b2528156dbe423eb07d7e0bf79695452d0e1e75c3",
							"linux-arm64":  "fd8b094d8b55f11ef23c3be32cd61593435c8968a1b0ae492d086fd79d1a36f1",
						},
					}),
				},
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-govet:govet-asset:1.32.0",
						Checksums: map[string]string{
							"darwin-amd64": "710fe8382dc9931c7f3fbd5b3482475c81193aeb58e7e7a92d52b86502f16977",
							"darwin-arm64": "aa5147f991433d3bddcc2ad9e8f4df672742f2b9e817cc559f24c8fd5683eb15",
							"linux-amd64":  "709805b5f4e31ba87c03dc7e673aa975ed1f092f42b9f17fe311643c4b705d7e",
							"linux-arm64":  "6ee72c59ac9417b500e7153c0722982af7b5807e8421f06ab47475dd90f5e251",
						},
					}),
				},
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-importalias:importalias-asset:1.31.0",
						Checksums: map[string]string{
							"darwin-amd64": "f9d7c54b05651ea021036a57f19b24b15e6a56f2b63ea55ff8eae57705164e08",
							"darwin-arm64": "659ca1f389a5742360dc6ac7c9250287f07a29cb1894b8740ef698eb963cab92",
							"linux-amd64":  "495b433da80e180500419201a240580f9f157aa90ab372f6633d68a6b6bd8c54",
							"linux-arm64":  "8f49287829c7a2adfec934a1f1b829e31d4aa30b10cee812131eed7bff17ca76",
						},
					}),
				},
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-ineffassign:ineffassign-asset:1.34.0",
						Checksums: map[string]string{
							"darwin-amd64": "239df01ce6af764ab503f96ae4a03a74928d8f7c5e088eb4f46b7fd68516770f",
							"darwin-arm64": "2493a5ac74a6fffab19a7570a4bb471cd497d11144af8fec0a1db2898513b468",
							"linux-amd64":  "563eef788caf4832b7559f4cbdc3b10e2661465cf1c6981e4536e462c2e602e8",
							"linux-arm64":  "a1694783d8f67ff99ff9d283b658b0fee5323d6ba23207d32eea7471e8fed9ff",
						},
					}),
				},
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-outparamcheck:outparamcheck-asset:1.36.0",
						Checksums: map[string]string{
							"darwin-amd64": "725e268fb66ec8911a63e8f1856c6e05456a62286483764f108109a90aa1878a",
							"darwin-arm64": "97a02229f735f5d6661f2875689061d03b3a3244d7c29f3174c61f5f0f9c6540",
							"linux-amd64":  "15117d7e0682e9adc121ddd48e8d59957cc04414f6fa268de5f2ca948181f66f",
							"linux-arm64":  "7b4ff5b81bdf8093450d81b8b454b7bd9968311cdb37a02b10ec59444fa2b266",
						},
					}),
				},
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-unconvert:unconvert-asset:1.36.0",
						Checksums: map[string]string{
							"darwin-amd64": "706fa703eab16b3bda9256ee9978389714600b614dd54e37fab4edae5805ca17",
							"darwin-arm64": "494595566a6c0febee80338f26f6b7cf9db50fb7c85ecf2d3fc558d038103beb",
							"linux-amd64":  "8f2a38fe0976ec1bae3fba590e625cc071bed2175956b6827992523abc13058c",
							"linux-arm64":  "793247a861c790f52e986adb1b929bb2b8a196271c9df29839fd7463bf9e9d62",
						},
					}),
				},
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-varcheck:varcheck-asset:1.36.0",
						Checksums: map[string]string{
							"darwin-amd64": "8c94135b0e5141db74c1deed64ebc78f1b8ddab5614d1edcc91fd99dd78893f4",
							"darwin-arm64": "f5bf634dab048a8ed464910a59d26c8afd1b48910cc25996192dfcb6467d493d",
							"linux-amd64":  "393d5fe017be3e4efa1415c390a1a2ba1734207a7126315fdbdadaaad6e12b9c",
							"linux-arm64":  "33c9d6a29b2e9bff27842e35dd317050c731f22dd813b9140e2752b9ca089b02",
						},
					}),
				},
			}),
		},
		{
			LocatorWithResolverConfig: config.ToLocatorWithResolverConfig(config.LocatorWithResolverConfig{
				Locator: config.ToLocatorConfig(config.LocatorConfig{
					ID: "com.palantir.godel-license-plugin:license-plugin:1.34.0",
					Checksums: map[string]string{
						"darwin-amd64": "8a1b8f03eead9fec6c6cd1aa533a353ab14a091d8f43adda72a5296ff4360ab6",
						"darwin-arm64": "d3e6de9fd074d74ac765bcc01aad570801a9a2c4a87e2962b9514a4a9a095803",
						"linux-amd64":  "ffa32f5cdcad147751aaed258da7edf8f5e98841e314288733afb51e87e8319b",
						"linux-arm64":  "95744b593ac583a38b410ebb8ae14c759fe6ecaf63695ab01f8f69adc9f0146a",
					},
				}),
			}),
		},
		{
			LocatorWithResolverConfig: config.ToLocatorWithResolverConfig(config.LocatorWithResolverConfig{
				Locator: config.ToLocatorConfig(config.LocatorConfig{
					ID: "com.palantir.godel-test-plugin:test-plugin:1.33.0",
					Checksums: map[string]string{
						"darwin-amd64": "40050b3709d2fe99667fd8244bf31c9ecd79535c8c0dc24866d005fc1d6c6001",
						"darwin-arm64": "b2edddf07e2e648515abbbf182577792f1093a4adcd188588f5afc7198a4c2f6",
						"linux-amd64":  "1be448d36b7ebc5fb234aefb2710792a2ca6daf61c7112ea0c7bdf01a0e3544f",
						"linux-arm64":  "86c72770e58efe677411df674d153983a55080d83f363ff444ba8533e7ea39ae",
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
