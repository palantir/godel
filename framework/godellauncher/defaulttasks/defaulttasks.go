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
					ID: "com.palantir.distgo:dist-plugin:1.40.0",
					Checksums: map[string]string{
						"darwin-amd64": "363bbd86a15974d834eab8666d9ddc2341e59c38c400b49a5bf5c55a73930522",
						"darwin-arm64": "f9fd61ccc7846fcc15b1e1f823d2db9d1a2f4f6b2e7bd08126478d542a83cd43",
						"linux-amd64":  "1b5985e5925c7de80611ed4f4ca7b519d2cd8f50f9d29115abc5558db3e65ef4",
						"linux-arm64":  "f42ddbbdf4594753708b72820a8f87edf163f133d213802b662264d87fc41aea",
					},
				}),
			}),
		},
		{
			LocatorWithResolverConfig: config.ToLocatorWithResolverConfig(config.LocatorWithResolverConfig{
				Locator: config.ToLocatorConfig(config.LocatorConfig{
					ID: "com.palantir.godel-format-plugin:format-plugin:1.16.0",
					Checksums: map[string]string{
						"darwin-amd64": "7ed2effb5e38caf927b850bf2a74adc81a1344440d050eb107427a05b3bbe6b3",
						"darwin-arm64": "48b26b419c9105d9795c4197bcc0ba112b58bfc36d01d86b32ad651cc35ae689",
						"linux-amd64":  "36948ed98f9f6a90da8828404758629b605cf0f61cf516b8fae5632d6734dcc6",
						"linux-arm64":  "f858480228e210c6d873770eeef545a67594e0b21ff5fc6a2bf13731e18cec08",
					},
				}),
			}),
			Assets: config.ToLocatorWithResolverConfigs([]config.LocatorWithResolverConfig{
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-format-asset-ptimports:ptimports-asset:1.16.0",
						Checksums: map[string]string{
							"darwin-amd64": "036532d9145801720f8f795fcb90bf5a75a2ce69bf96e54961d3de4c7f777b9c",
							"darwin-arm64": "4c8a58e23deb66ea9dffc0f41a1ddeeb841b866237487fc47416aa9af64ed74c",
							"linux-amd64":  "d3bbf12a56677319e1bb26e7dd77fcc8f6360dcff940a8c3b996eb2b169a1a06",
							"linux-arm64":  "d20cb4b33df286719cd1ceb30eb56e22a0bf2514ede9e3f8b06a4a14a96c55be",
						},
					}),
				},
			}),
		},
		{
			LocatorWithResolverConfig: config.ToLocatorWithResolverConfig(config.LocatorWithResolverConfig{
				Locator: config.ToLocatorConfig(config.LocatorConfig{
					ID: "com.palantir.godel-goland-plugin:goland-plugin:1.11.0",
					Checksums: map[string]string{
						"darwin-amd64": "600a9a58290a4b7a07522a15c100a06adb7fcc258bc6848791126d21c30b6366",
						"darwin-arm64": "d99d5c98192185cca8e28b12285c9c6b224d1caf155e091c3a956555a21c550e",
						"linux-amd64":  "5d07db95432530c20caece26c3e6a69e63caa493f90f3e29408fb2ce085d1108",
						"linux-arm64":  "09a96cb88a6fd909e464175dec8f825ab84bd483f8048e737fe0d7f26527e3ab",
					},
				}),
			}),
		},
		{
			LocatorWithResolverConfig: config.ToLocatorWithResolverConfig(config.LocatorWithResolverConfig{
				Locator: config.ToLocatorConfig(config.LocatorConfig{
					ID: "com.palantir.okgo:check-plugin:1.21.0",
					Checksums: map[string]string{
						"darwin-amd64": "fb0bdf716c396c5de913b4855039e64aff00bf913a064ae876861e84bd84867a",
						"darwin-arm64": "38cbe1b2b04a0290d26b0a80dd68942a64627eea732b7994bc9107728cf44686",
						"linux-amd64":  "f74682d292c6978892b32c3c4071b1f7c9c925178dfd75fda6deae822a0da88f",
						"linux-arm64":  "00372a65c14b38b6214b15f5b6c06d950205dfbd94e06a082e3dcf6f38e4fc67",
					},
				}),
			}),
			Assets: config.ToLocatorWithResolverConfigs([]config.LocatorWithResolverConfig{
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-compiles:compiles-asset:1.19.0",
						Checksums: map[string]string{
							"darwin-amd64": "93fbf2e82652c818561bf99e3228d46cfcdf17fd0c9b5ba8e1cb05dbc7ef9afe",
							"darwin-arm64": "831f3adc7524e58bbbce6febdb97ad9bef33d6a981641fdbcca6292f7cb11197",
							"linux-amd64":  "6755afb555d11eba03e0d7456266efe2c75fb7494a1780b9fad29780dcca9d52",
							"linux-arm64":  "52f8d1fb08dbd4b15be10016bd445dc67913a8636108a1dad277a9699670bfcc",
						},
					}),
				},
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-deadcode:deadcode-asset:1.18.0",
						Checksums: map[string]string{
							"darwin-amd64": "8bb623e73264051c4bbe8a110a08bbe41bcf7d83a301e1f83a9f7252ebe7f6b6",
							"darwin-arm64": "68a8b5a03eadc2dfc4b71c574666e6e40d99186b2ed8f1597b03d031b375a85e",
							"linux-amd64":  "ba5b1754233f3609fe191ca0912bc0319b1f76728e1369c3a17153f083dccb09",
							"linux-arm64":  "3d1851d9a9d672b70811bd64ff002fa08e4c6b0cc88c2b0d61f63a192951db5c",
						},
					}),
				},
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-errcheck:errcheck-asset:1.18.0",
						Checksums: map[string]string{
							"darwin-amd64": "734135c4262bd88940cefa7fccb0ad2a458c3175fe4d05f476a7ab5df80f9ade",
							"darwin-arm64": "d9d58c7ba6cedce8ff70e13ad0ee5d6f05c14ff0524921ebf22894878affff29",
							"linux-amd64":  "92ffd2cae742149644ddfdf78872598752f23aea94332cf2c947567c31cef9a0",
							"linux-arm64":  "2f689cd9a372260dea34b158b37c41efd6db7aad7e64f64638fe821a0810daa7",
						},
					}),
				},
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-golint:golint-asset:1.13.0",
						Checksums: map[string]string{
							"darwin-amd64": "d07122fd864de29dcf4d991387a068e3625286861efe083857ed82e56fc4a418",
							"darwin-arm64": "1a6e8cde9b78a1146ed9495a5a030a10985fc434ce051f63894a637a1f228eda",
							"linux-amd64":  "eef4341e850e077a1781bc6d15a8a85ca759426eed1e178154d98d0a51af09bb",
							"linux-arm64":  "4f680546fc703877d89150b2527ceb181c8dbac360a17a1afd68bdb8fbf0428e",
						},
					}),
				},
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-govet:govet-asset:1.13.0",
						Checksums: map[string]string{
							"darwin-amd64": "23f78c5366acd44882850b886a5085869b94cc6d5dc66e49ef966dbeb7f88dc6",
							"darwin-arm64": "565605e0ddef093358d4a16ac5bf9ce85369e7ad4935883a4938bda2c60b9cfe",
							"linux-amd64":  "c37053dbd36219dc7537469f0dbbaa7bebe5196e83f8ff9ca1494df12e927655",
							"linux-arm64":  "dbdb87ecb6c8bdeb16329fdaf663ad2e31252e2a7fd4165240a88d7b97b6e561",
						},
					}),
				},
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-importalias:importalias-asset:1.13.0",
						Checksums: map[string]string{
							"darwin-amd64": "22383a62fb087bd48d7c76a5e1011c5f141fafc20a56f026c801b5fdc2894aa5",
							"darwin-arm64": "c7ec4b590eb5fa44f77b8ac8857d451c671490561e78e9d570c91644542820b5",
							"linux-amd64":  "f1b15b6dbea36ff185650279618f2a4898300c4f60117cab89331fc3408246db",
							"linux-arm64":  "303bcacd4c467c68329161e9b2938cb62ff96684aa104589764645e2c6bd8cb5",
						},
					}),
				},
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-ineffassign:ineffassign-asset:1.16.0",
						Checksums: map[string]string{
							"darwin-amd64": "a7588461e538f203056e4bfd5023895e330afc73e37c8a9c351146c7a281601c",
							"darwin-arm64": "f453f5ddc466fc11292c7a417473de37fbfd9901a18590d1225aef8c61006ff7",
							"linux-amd64":  "a02f88063dfd20ca2c838d4c3811329577c1dbc85b2d2d0825bc7eba69950105",
							"linux-arm64":  "4b4ae84c03204364f294237a77f8940e20749b09f4ad7f326751fa40398d29b9",
						},
					}),
				},
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-outparamcheck:outparamcheck-asset:1.18.0",
						Checksums: map[string]string{
							"darwin-amd64": "085506345c75098990ae5b90faaf8407c4eeb27a709db429c7f00b0bc03fcfb3",
							"darwin-arm64": "53e4d9907bcd94397b9a7f520ec928ac2343242453223904e9a7a1920b56e926",
							"linux-amd64":  "c4cf0582b2422dbf5abd4c31d26a20deeebc4ac9bc604c2539594d8b7536fb7e",
							"linux-arm64":  "3eb76f09c191a6bbd2db01ffa1af5a9e6a4f20b58549eca273cd9334ebc6c000",
						},
					}),
				},
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-unconvert:unconvert-asset:1.19.0",
						Checksums: map[string]string{
							"darwin-amd64": "6baa8c295997e741203dc02864bb7a0003f8b69dade8578b99e18a33c06164c8",
							"darwin-arm64": "8dbc618bd9bf9bc13dac9375c19d41ff033296fd2c19f0c5f97b46d526de1ebc",
							"linux-amd64":  "b849a18f240edac6ff7e109aa333d8db33f325b33eb6fbc75b2fbb655f1063d8",
							"linux-arm64":  "5c03338c119a787b522ef9bf0e9662f44fa66cb5651c763ae9a8b05242ae0aab",
						},
					}),
				},
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-varcheck:varcheck-asset:1.18.0",
						Checksums: map[string]string{
							"darwin-amd64": "07a9087ed2c8e8f007e4e1f2b4a2c325d6d75d53a2e74e5439eb500221f65348",
							"darwin-arm64": "cc5b7c2e058dc4ad5ff0667476afc04ecb93be5288e1e89303543ab5ab96621d",
							"linux-amd64":  "19fbf5d4d4bbd779f2bab87e0a639a10cb2803c289ac890c1bb81e46533e735f",
							"linux-arm64":  "4370d82ba571f5c7a74f5fea3376ee9d46c6add69fd6a2f47b38ad001126852d",
						},
					}),
				},
			}),
		},
		{
			LocatorWithResolverConfig: config.ToLocatorWithResolverConfig(config.LocatorWithResolverConfig{
				Locator: config.ToLocatorConfig(config.LocatorConfig{
					ID: "com.palantir.godel-license-plugin:license-plugin:1.14.0",
					Checksums: map[string]string{
						"darwin-amd64": "2f35eebc7800beb1e12d1e6eb60471aa1e9df44f1af4a9e85d098a2a10304b21",
						"darwin-arm64": "13d7f6331b169c060b67a442f8c65acbd9d28d735589e437e4524bc459d4a273",
						"linux-amd64":  "0e0c1ff2108e903c5029063bba773b7e7587878d28ada755b6a0ed35820a5f8e",
						"linux-arm64":  "d7e600462feee30a6ef7296ef5e0dba7264afc9313f736335fd6017bdcf94e2a",
					},
				}),
			}),
		},
		{
			LocatorWithResolverConfig: config.ToLocatorWithResolverConfig(config.LocatorWithResolverConfig{
				Locator: config.ToLocatorConfig(config.LocatorConfig{
					ID: "com.palantir.godel-test-plugin:test-plugin:1.16.0",
					Checksums: map[string]string{
						"darwin-amd64": "a05a9b9e96909e7e7eb26a5a34eaa4fc1b4d0c2b80637a8f7650a315481e52a8",
						"darwin-arm64": "15a106d89fcc514906c47cdeccd1631fb9a51494f878fd2dffcf94b425578b30",
						"linux-amd64":  "1ca8bb93069ad8fe95b2c2aee38e46840b5fcbc1f4ea6ca6c50308788630d0dd",
						"linux-arm64":  "e5246c4879548630124bc2892fb1a016611de76ce729dda638253124a4c58cbb",
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
