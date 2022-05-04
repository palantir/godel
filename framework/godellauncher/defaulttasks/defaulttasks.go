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
					ID: "com.palantir.distgo:dist-plugin:1.38.0",
					Checksums: map[string]string{
						"darwin-amd64": "6732bda697ea8159af9d7264f01d37715990598cf46b5129a4bd5609ab0f8c03",
						"darwin-arm64": "4b488d43223d18ca4cec753f42518c8887be32472b476fd941bbff412db72c5e",
						"linux-amd64":  "dca6491dee6187ee5929e2d01f87a6f0ffb6e0235bcf60ee616de6b9e1a4459c",
						"linux-arm64":  "2bb637f96502d4dbfc06683c89eeb823556908dcb6a034f4f3f5f7cdfb45cce1",
					},
				}),
			}),
		},
		{
			LocatorWithResolverConfig: config.ToLocatorWithResolverConfig(config.LocatorWithResolverConfig{
				Locator: config.ToLocatorConfig(config.LocatorConfig{
					ID: "com.palantir.godel-format-plugin:format-plugin:1.14.0",
					Checksums: map[string]string{
						"darwin-amd64": "880733098cec14f3532dc1e395a17019978901495084f0de125fb2ca6dc438fc",
						"darwin-arm64": "ccd5fe06bf39a05eae06f9c377a53b3dbf0239e8ae4930637b7e42e67dc9e576",
						"linux-amd64":  "b03b0d433efcaa989bf54237bdc4d4d25f063e8d72ee1aaf53cfd8779dd7dcdd",
						"linux-arm64":  "1d44e25f8b96ca774a315ef55e689a07a6c3d328c8cf8ab7d3c2419e42fae30b",
					},
				}),
			}),
			Assets: config.ToLocatorWithResolverConfigs([]config.LocatorWithResolverConfig{
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-format-asset-ptimports:ptimports-asset:1.14.0",
						Checksums: map[string]string{
							"darwin-amd64": "6748ae95531b758630a1391b06f1414e155da61fdbd0f2179bc198f5d086848e",
							"darwin-arm64": "c13565a53cced426a1a8473e70ecab8a0cd8ca98b78f859c4b3e40739c3c67f4",
							"linux-amd64":  "f70e9de55badad374399758aa7135be78b45bf794f04e17be79ce2605a8cc615",
							"linux-arm64":  "28a204ccc76d349a96338bb8eea63a3fe5fd387e4d14c855f906d7bb3fc6dd57",
						},
					}),
				},
			}),
		},
		{
			LocatorWithResolverConfig: config.ToLocatorWithResolverConfig(config.LocatorWithResolverConfig{
				Locator: config.ToLocatorConfig(config.LocatorConfig{
					ID: "com.palantir.godel-goland-plugin:goland-plugin:1.9.0",
					Checksums: map[string]string{
						"darwin-amd64": "25270d66ee8a7e6a6ee1026d870c64b45db48bfa10a17cae133eebcb13926206",
						"darwin-arm64": "06e9436d390611c22ebedba879ae4c2a17952b3faf6e0ac7c3222bf10574c3d1",
						"linux-amd64":  "44764fb3baf52364d1338c0b482001f9f687ad6695635db9bcfcaaa9b87ba9c6",
						"linux-arm64":  "0fc949f9d9cacb5aa52596fac2b7790b987af0e6893af278a146711a48334d59",
					},
				}),
			}),
		},
		{
			LocatorWithResolverConfig: config.ToLocatorWithResolverConfig(config.LocatorWithResolverConfig{
				Locator: config.ToLocatorConfig(config.LocatorConfig{
					ID: "com.palantir.okgo:check-plugin:1.19.0",
					Checksums: map[string]string{
						"darwin-amd64": "35e4694c175de3e8f796036925aa51b3ece4c1868f693f89de77f04b8a3a9fe1",
						"darwin-arm64": "5e46bd08d1bce7bb924d9c110f120dc0fb5a481bdacd7b5df032ba555ce11af5",
						"linux-amd64":  "d1abc80bfdde1d5b37a55408913010c9bafe4350e6d6f0e292469cc85991e316",
						"linux-arm64":  "0f85beecd827ff7ab9b38c902fdfe2d2c74dff15bc48fc3605c8b017ef6685b2",
					},
				}),
			}),
			Assets: config.ToLocatorWithResolverConfigs([]config.LocatorWithResolverConfig{
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-compiles:compiles-asset:1.17.0",
						Checksums: map[string]string{
							"darwin-amd64": "1d02d48f581438bb14e07f425e7c1827b145ce1b8143cdb0359dea15439ac39f",
							"darwin-arm64": "b19e72bcbe7fc6bf90dacbb7745a147983a7919b1b8200297c32c056a394f243",
							"linux-amd64":  "fdfed742020b10b2a157b16d2c33e50827f5f0ede22b370380f3965846f16e0f",
							"linux-arm64":  "72f5e1282e5b62d802bff02759dc4d04eef20af633f8bd54af37ab5b1228bac0",
						},
					}),
				},
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-deadcode:deadcode-asset:1.15.0",
						Checksums: map[string]string{
							"darwin-amd64": "e53166a3d38cbb17b7468e9e53dc13ac469ed2d8b9285a2d28326191be983611",
							"darwin-arm64": "2f70ce2cf8babff177d628a41918f734d8bca0d0f7aa6f64a29225acf5303077",
							"linux-amd64":  "81024e6605dab33292d1e698bb00ac767d5079e76b721a21bff907b15034b8ee",
							"linux-arm64":  "e7b8d452924ff8615e781f241869a73846ed8fcf9bbbc1d990e44e53a83758a5",
						},
					}),
				},
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-errcheck:errcheck-asset:1.16.0",
						Checksums: map[string]string{
							"darwin-amd64": "824dc9eeae1b65173f50cda2b3cbd8f1923c96ff980d182c23ff5d40d23426f0",
							"darwin-arm64": "64a38ecdcc1becda2165c0e437c61de2b772eb7dc607d1a04f6f0178a92a3dd0",
							"linux-amd64":  "d71c0002e85c2b070402c8e677d3859b7575c498d38b64e9dbc795d5e83a42ae",
							"linux-arm64":  "7b81cdbde8162d36a9c29c6ff76fa7522e2b3a4c5a455168e59e96e303300f9b",
						},
					}),
				},
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-golint:golint-asset:1.11.0",
						Checksums: map[string]string{
							"darwin-amd64": "2defb276e659a64cb7cb28b6158d7f9344746ba15c292ac796740a6a07dd21e9",
							"darwin-arm64": "a43f9fe282d7d1e92575f94fa8e7620ea334f0a5f8019abf4a217446a7045a5d",
							"linux-amd64":  "13bd1c5125970f7e3acc043dbd626ce6163e394bca9ecc5c41fb7bfa1ba0d052",
							"linux-arm64":  "f112627e4aef16117e81ea63258ebc0d10569f2855695e9ee3f7ef478c73aa18",
						},
					}),
				},
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-govet:govet-asset:1.11.0",
						Checksums: map[string]string{
							"darwin-amd64": "9415f325f94de174c27d8b189259db2d3e1fdfcacb6be59a34240a99e2118c12",
							"darwin-arm64": "0cfa613e974186cee9e0341f91edbaa32e9434b626a130d06ad5ab65e4e8f206",
							"linux-amd64":  "c34efc0b4494ec149e60d55da6d4464dcb020820a746bd7a1b656ee92c3de497",
							"linux-arm64":  "4959f81a241b2f54695b5b7ed05cf4c9978bb12e7dcd619fd6d9addfd877be0b",
						},
					}),
				},
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-importalias:importalias-asset:1.11.0",
						Checksums: map[string]string{
							"darwin-amd64": "a11ce9ec2781ec16edb174f48e2e1a0a371d956556c8111e106b72f71ab5ed64",
							"darwin-arm64": "aac1df750ee1986a02b44e16f8c9f4707d212bfab32255165c6222e43d959291",
							"linux-amd64":  "602823a470657946fbae7649f165d5f847274b89afed9947417a74bddf77f06a",
							"linux-arm64":  "cf84413491383c51c1294bb83eee6fe2b98d5c336b7d03666c2adfa9859b77a6",
						},
					}),
				},
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-ineffassign:ineffassign-asset:1.14.0",
						Checksums: map[string]string{
							"darwin-amd64": "448f20a23394d219a02abbe64a59e42ea599e0f521aacc68b8240300dbea9b88",
							"darwin-arm64": "abe81f53ddd68ceedc0805cb6e765a029d9b9d39de072e011155b4b093e7dc22",
							"linux-amd64":  "b92a60599ffb8a272dd78ddc92ea835b5a844eb03421a5cda19001e74949c144",
							"linux-arm64":  "ed7246342d6ada0fc71b0a374cf8edfa1a27e1559a29f8e5ba545579a8a2ab52",
						},
					}),
				},
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-outparamcheck:outparamcheck-asset:1.16.0",
						Checksums: map[string]string{
							"darwin-amd64": "51b1c4359a7cc3e99ad80dff44702b0d1791845967a13ede3a8100c9b933b81c",
							"darwin-arm64": "3f68b61d2134290b94d2566147678667a8919f1f2b7ca5060a258f52ce03afca",
							"linux-amd64":  "ab3f6c1abbdcd88665ac3b5bd7a9982df95d29ebbbb11aae29731c8ba47add95",
							"linux-arm64":  "027893d94c9b2d0bb5e297b1a263b432df9adfb045303fe379a8cb266706a708",
						},
					}),
				},
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-unconvert:unconvert-asset:1.16.0",
						Checksums: map[string]string{
							"darwin-amd64": "3bd04bbe868bbb80b02ddf4048d71141043df11967f5d495aa67b99b825c7ec8",
							"darwin-arm64": "8c6ff4a324678d4453f0dd0897590972e05a95d74b376a067c74c6360a343fd4",
							"linux-amd64":  "36e21a0b4485537ef1308e26d19b4108e7ec1f84135d49c3605d86594bfac761",
							"linux-arm64":  "26433682862f13b3b449ab048c75cd9945973cc2a1c2202ec1a4752d4ca21843",
						},
					}),
				},
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-varcheck:varcheck-asset:1.15.0",
						Checksums: map[string]string{
							"darwin-amd64": "0edbb3b4046b75df32b9f01580dad46c41e9bc226c2708dd943f3b945f47fb33",
							"darwin-arm64": "21ddd95ccb38a8ff728b98f03d09da19c7f4e9c604d9a1837949de86e19e4ec3",
							"linux-amd64":  "c2d7198aab8c479995876c5c464af32513f050b3dc65df8e4e367c284d142964",
							"linux-arm64":  "9bb25424e3cdc1de66e4e3d9b495f5e1ec383327ff8c16ced13c897109c11835",
						},
					}),
				},
			}),
		},
		{
			LocatorWithResolverConfig: config.ToLocatorWithResolverConfig(config.LocatorWithResolverConfig{
				Locator: config.ToLocatorConfig(config.LocatorConfig{
					ID: "com.palantir.godel-license-plugin:license-plugin:1.12.0",
					Checksums: map[string]string{
						"darwin-amd64": "a229f8a73a61d0967da5103e4e38d0b08e4467cd1d79e237a3a2ad5366479249",
						"darwin-arm64": "93428d345111606d94bae2afa0126758b7b50f4b10f6424fc9acea355a48b597",
						"linux-amd64":  "1700fb4e42bad7f3b2cf1b8b49518fc9de2b8f034666aa41687a1756527044e5",
						"linux-arm64":  "2de63fd02f506e555467445d520db41d41ee4e11be6bc95f0692e963f54bd998",
					},
				}),
			}),
		},
		{
			LocatorWithResolverConfig: config.ToLocatorWithResolverConfig(config.LocatorWithResolverConfig{
				Locator: config.ToLocatorConfig(config.LocatorConfig{
					ID: "com.palantir.godel-test-plugin:test-plugin:1.14.0",
					Checksums: map[string]string{
						"darwin-amd64": "92590940e09c3a9618395dd4333b441b236acc485485dfb6ff7e73cfb362a2c6",
						"darwin-arm64": "742e65175c0339d3c3249e9ba6b44c0112e89dcaa370324f98fce5fcdd246066",
						"linux-amd64":  "52cfa41dc3522c8c324ebf0d73331dff01842a8a9c414aa5d0042abeebfc1b65",
						"linux-arm64":  "be48eb17b89c0ea25a552a0a04b27027c8493a4d909151bbe66c657550e42c6a",
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
