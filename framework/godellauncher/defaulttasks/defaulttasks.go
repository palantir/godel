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
					ID: "com.palantir.distgo:dist-plugin:1.31.0",
					Checksums: map[string]string{
						"darwin-amd64": "e22750a11b9b0d1003ab92477fa436d0c692fd8b927131905445fdec20a253ac",
						"darwin-arm64": "f184713c28e002decc8691cb9418fac2c3531619521f87a1cb1ef987726e0c12",
						"linux-amd64":  "5f0c5aec6587880541a5ea7469a44a88cebe8f048acf54f2a4d9f45f14976d19",
						"linux-arm64":  "af0bb8ccc182c76870201cac7466b0433d624298dc0fed6a075b25a98e3726c2",
					},
				}),
			}),
		},
		{
			LocatorWithResolverConfig: config.ToLocatorWithResolverConfig(config.LocatorWithResolverConfig{
				Locator: config.ToLocatorConfig(config.LocatorConfig{
					ID: "com.palantir.godel-format-plugin:format-plugin:1.9.0",
					Checksums: map[string]string{
						"darwin-amd64": "b3bd44f450549df9a06ab9f024923f3f78befe878974d859d03622f01cf8b75a",
						"darwin-arm64": "ee14e8ea89b91a62849fc347f32d2c308805f53c8e715232ad3b0f6fe6d42274",
						"linux-amd64":  "6ad7dbd06e0c05c6383c1fe6f680e7b74208d098b267ca8d63be1c5fc3a1fc0c",
						"linux-arm64":  "51d90722459e93a0569d5e7116161bc70b1447e8df42d07755f2ab02da9c1f6b",
					},
				}),
			}),
			Assets: config.ToLocatorWithResolverConfigs([]config.LocatorWithResolverConfig{
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-format-asset-ptimports:ptimports-asset:1.9.0",
						Checksums: map[string]string{
							"darwin-amd64": "af4be7fbe9a282ae3c59a07df26999adcf7a80d0186e0e78d88b9da40337ec47",
							"darwin-arm64": "3431d5bdbb495376c27b7a62f6afcdaaf8c348977c10a70cf370155721d7a905",
							"linux-amd64":  "9bfd1c17ac6e17acf19e5be3abf0f4c6e92fcbeff07c8cec809537a298abe228",
							"linux-arm64":  "b70eeda355f6598c567e1525349bd5d1bddabb83b72760e3e5b7ff7b7f505b52",
						},
					}),
				},
			}),
		},
		{
			LocatorWithResolverConfig: config.ToLocatorWithResolverConfig(config.LocatorWithResolverConfig{
				Locator: config.ToLocatorConfig(config.LocatorConfig{
					ID: "com.palantir.godel-goland-plugin:goland-plugin:1.5.0",
					Checksums: map[string]string{
						"darwin-amd64": "deb203a1b28c26c7191535341e302189c35ce191f69998a8b36bc3f0711ecad9",
						"darwin-arm64": "2978c8ccd2404e4bedcead4722ef48f1c9e12a4d6db5484f86e7a07b61ec088d",
						"linux-amd64":  "708147a09ddcaad98e887f875a8b3ab0b59cd8709be2de8c25453d1a4296de87",
						"linux-arm64":  "e6141b8cb6fd8768deba7b51b650e04f296f0affacee15c4a2b88f16ae5d7e60",
					},
				}),
			}),
		},
		{
			LocatorWithResolverConfig: config.ToLocatorWithResolverConfig(config.LocatorWithResolverConfig{
				Locator: config.ToLocatorConfig(config.LocatorConfig{
					ID: "com.palantir.okgo:check-plugin:1.14.0",
					Checksums: map[string]string{
						"darwin-amd64": "91a9f11d595038a253561771c569f07d19588dd5b868b1141778d02b499fca4c",
						"darwin-arm64": "729716eeff3592bfda7287225d8301909beef684a140f5ca414ed7ed80aa0127",
						"linux-amd64":  "ca5298cfc14fb9a34514f9e187525f2959d7b3755801473439b5ca4d823e409c",
						"linux-arm64":  "1336d11f3d1b9a0047bb3f41c1e0f1cec2bd04b7d0f0d9992a004d0645f69fea",
					},
				}),
			}),
			Assets: config.ToLocatorWithResolverConfigs([]config.LocatorWithResolverConfig{
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-compiles:compiles-asset:1.9.0",
						Checksums: map[string]string{
							"darwin-amd64": "fc331fc4cada70187b2509264f7bda5f71e813b3093259e180048f4597f25cbf",
							"darwin-arm64": "d64c78c52f5c16e269ead00b3b7277d02c00045582f701a3f5d025d71bbd01a4",
							"linux-amd64":  "a7f13e223ccfe16e3888d8adf68ae5009227ba3fcc00af3b90f0bd8e440020fd",
							"linux-arm64":  "5184d8ec40a434e7497dabe3556638ed79eb7ca1e3f8ad4800ad90cd128cfd56",
						},
					}),
				},
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-deadcode:deadcode-asset:1.8.0",
						Checksums: map[string]string{
							"darwin-amd64": "fdea969f302e2f7ce51536ddc3503a20a85704504cad054a4c7d778e30ce90d2",
							"darwin-arm64": "e76c811165c3abf1d5546bdf128565849bfee97cbf34179679d78f1221d4ca40",
							"linux-amd64":  "c30c4befd5877697da0122b340e8cc01146ad938c8f89a40d629228ee8529b68",
							"linux-arm64":  "aceee326d172477c5945707b2726cf614ec8296ffb41e151592b9ed5a816dd67",
						},
					}),
				},
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-errcheck:errcheck-asset:1.10.0",
						Checksums: map[string]string{
							"darwin-amd64": "d70043bc37501130f69ea24eeba05cbbfbb76a703b8a3c73fffada052c57302b",
							"darwin-arm64": "f7237c7b9f7438801eaa209ca872e7e396e232ca49128bf4210dc73866e4bab6",
							"linux-amd64":  "60e49db00b32acc3ceca58c9b28581a937bc0c4fba2322faeece37e2f02b12c1",
							"linux-arm64":  "fddb985eb8d43656afb8f191e9c979b0658f286195e49f55cea0767663bb4566",
						},
					}),
				},
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-golint:golint-asset:1.6.0",
						Checksums: map[string]string{
							"darwin-amd64": "190c84657c563e1025a1931c46545d8240dcbb2575e8a6e346c58dc82d43f997",
							"darwin-arm64": "25859106cb097fbd90f5cfc1be0a3d4b139373d2b22687959db2830fdd198e49",
							"linux-amd64":  "a25d20a17bfae6a691fdfb07b1e667c3a4d0dc945a9559f4769b30a8698ff562",
							"linux-arm64":  "32f286924c9d9b1a2b39b7542c2fa29b794360efa945787ef29f23046ec19fe6",
						},
					}),
				},
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-govet:govet-asset:1.6.0",
						Checksums: map[string]string{
							"darwin-amd64": "72f811228b0cc835a8dda00db54d7cbc9dedcf9c3430c18f3a3d81d918cc8eae",
							"darwin-arm64": "d4a68d905f908e52dccad5040ea5e8a517a6fca35eb9edcce6e756c6337a7dde",
							"linux-amd64":  "8eee7f19963db8c5930bd364f2cb6ddd28c1e212dd9d7582acbb4d8a89d51b5d",
							"linux-arm64":  "72ac305f3f7cca977cd0337c5113eeda4897173bed77dcd33b100d4ddf6c0798",
						},
					}),
				},
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-importalias:importalias-asset:1.6.0",
						Checksums: map[string]string{
							"darwin-amd64": "fee1586c590779a031a3db548192498547338dda1afa9a7b96f00ec1447ccd7f",
							"darwin-arm64": "08eece28998bb4265de0a6910b159090a8b60e0f71e7870be8d03ffe318ec04f",
							"linux-amd64":  "1bf65145a5c723a63ea6720c9ac68ef481036924d245519113e96bf7504ffd58",
							"linux-arm64":  "932717db2c7e3f8691efc0e900ecf50247a54a747fbb94c4f1d9d7918b4d5851",
						},
					}),
				},
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-ineffassign:ineffassign-asset:1.6.0",
						Checksums: map[string]string{
							"darwin-amd64": "e3935b34023d9769109502a9e3df31c2842e86b1aa3af043180cb7e8fcb88fcd",
							"darwin-arm64": "e206ecb747059a42242187d9805285049a346603518b53c986f7078f88aedf52",
							"linux-amd64":  "b7cedd4c06a9bf2bb13faa03b25dda1f941877e15a9cb7f37d1094869bb98b09",
							"linux-arm64":  "9f5267deddba748a997c2cf17218e9afe8f7349cf5450dee820d44cb3ffc7aac",
						},
					}),
				},
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-outparamcheck:outparamcheck-asset:1.10.0",
						Checksums: map[string]string{
							"darwin-amd64": "6838b157fd680d4d7270288e663c833ebe2d24af8ce3cae47d1505acc0cacc6d",
							"darwin-arm64": "4a90596a3084cc9bf4671bf1e20bae36f1385aa220174d8525a501719b1e18a8",
							"linux-amd64":  "54e9de32f4d116c924009f6aac73360587eb03fea1725daa24063d624ad605ca",
							"linux-arm64":  "58d4986701fd2e3b9d5c48ebbfff4bcb3842602e6ff4cd4a2949d6773f6411f9",
						},
					}),
				},
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-unconvert:unconvert-asset:1.9.0",
						Checksums: map[string]string{
							"darwin-amd64": "46472ba758f74506c061a94d67b7397ea7737e18ba4207f87f452b620d949664",
							"darwin-arm64": "106c3519492ace227de4610559b0b97f1e8924e7349a65f6415f97c6ba3258ce",
							"linux-amd64":  "68156b768470230f5c7e150490f15f7b78b6e5e678858a9d61819695446dadb1",
							"linux-arm64":  "699be45a8fb5b0c2ac0ae54a44c2b57e5395f79a34a147f59f8aebc6a90c9e5d",
						},
					}),
				},
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-varcheck:varcheck-asset:1.8.0",
						Checksums: map[string]string{
							"darwin-amd64": "9d57c7b86e469dde6d3d60f0151724053f4e618c69ec4a4e19496625f394ae98",
							"darwin-arm64": "31e97dae606281c9d42b9190caed67b1866a4d537f071c1f251410a7f428ad0f",
							"linux-amd64":  "63cb0a807e54c415a9d8e47cf10cca5588109f2f963c8bdcd690c11b7e6d6f54",
							"linux-arm64":  "f2d64393a5c8844d8d252e56374b79e43f90ff3bda98343d8d38e29ad3803ead",
						},
					}),
				},
			}),
		},
		{
			LocatorWithResolverConfig: config.ToLocatorWithResolverConfig(config.LocatorWithResolverConfig{
				Locator: config.ToLocatorConfig(config.LocatorConfig{
					ID: "com.palantir.godel-license-plugin:license-plugin:1.7.0",
					Checksums: map[string]string{
						"darwin-amd64": "521b4029cb33360edb97d40326ff519031da94f57f15d5c1d26fb7c06419ea38",
						"darwin-arm64": "207c9d3ddf5fa70d7bfa558e58f9af4fde7b3d52b4a4a15ced1ba3668ba78861",
						"linux-amd64":  "1f4f7fbdb0c04af4d0247d8a8aeddf05b2d51a9e914672748273a3313e0a94ad",
						"linux-arm64":  "9a8e5a1874d4564a3ea63600ee8879d56edfc1652d88b04f13e9cce32533bfe4",
					},
				}),
			}),
		},
		{
			LocatorWithResolverConfig: config.ToLocatorWithResolverConfig(config.LocatorWithResolverConfig{
				Locator: config.ToLocatorConfig(config.LocatorConfig{
					ID: "com.palantir.godel-test-plugin:test-plugin:1.9.0",
					Checksums: map[string]string{
						"darwin-amd64": "4d77d8f0251c7ab0149463d978112caf19986bbc4ac1004f74c0483999e2ca96",
						"darwin-arm64": "e2610529c84f8e9ca6a196e2081769c6c7b558154dc04e3d62676c17152680c7",
						"linux-amd64":  "e4a613ec4006bdcf8448924303f9a422090b18d9edb0ea7f768a15b2b662f25d",
						"linux-arm64":  "f7c6a3bd51959a7c748a6f54e58f04172e88b2767ffba76ae431c3893ccf0712",
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
