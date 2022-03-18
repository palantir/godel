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
					ID: "com.palantir.distgo:dist-plugin:1.37.0",
					Checksums: map[string]string{
						"darwin-amd64": "d2165a0c005a73848e3d000f341bbe7ba8aad42077696196019d4de8e4b0356d",
						"darwin-arm64": "d68af31478abb86cbf4a723017225dcd8595cce433604b3d56fc9f8fc979383a",
						"linux-amd64":  "24d97c3ddef384ee18cde41615735c8b1dea0b1045ff9bd9f06acc6b0ecbd793",
						"linux-arm64":  "ca711ada25780110a74a4dbe4bb834c7910ef3d063922c12a113ba989e9f279d",
					},
				}),
			}),
		},
		{
			LocatorWithResolverConfig: config.ToLocatorWithResolverConfig(config.LocatorWithResolverConfig{
				Locator: config.ToLocatorConfig(config.LocatorConfig{
					ID: "com.palantir.godel-format-plugin:format-plugin:1.13.0",
					Checksums: map[string]string{
						"darwin-amd64": "b019c8f8abeff39b25f24d41878e23d1844d50a1c88b08f9072c88ff0b6dad86",
						"darwin-arm64": "1805c59c83a790031aea2b8b7a7af2f09dd2f82884b55242a12c5cbc9b66c1c2",
						"linux-amd64":  "0e01b3f57f71a93720d1950bf43d95d3fbe7663945e2a53d9c6933b8f5fc015e",
						"linux-arm64":  "3f8959700b3ba29888bec0223b6db7a50489e6c6cd282d38b80f9565838f37da",
					},
				}),
			}),
			Assets: config.ToLocatorWithResolverConfigs([]config.LocatorWithResolverConfig{
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-format-asset-ptimports:ptimports-asset:1.13.0",
						Checksums: map[string]string{
							"darwin-amd64": "be9ae48587570614c69aeb8689322e37348384daabe1a326742e19c0d8b357dc",
							"darwin-arm64": "65715c67c07aefe65bf67230bc681ea75c6651afd4b3767d5079b49ea8a173e0",
							"linux-amd64":  "73462bff7af1950a4923fc50fe2a3e436e116f3060a39e9cc9fddbbf0c3217e7",
							"linux-arm64":  "4324af89cbbd9ebb9e3132ce1404f8be2317c26112fbaef5c8251e1f53b36c1c",
						},
					}),
				},
			}),
		},
		{
			LocatorWithResolverConfig: config.ToLocatorWithResolverConfig(config.LocatorWithResolverConfig{
				Locator: config.ToLocatorConfig(config.LocatorConfig{
					ID: "com.palantir.godel-goland-plugin:goland-plugin:1.8.0",
					Checksums: map[string]string{
						"darwin-amd64": "1d89d845a7b1d745db0fd72a11004fa8636cd26627a98109fda19933eb6f34e1",
						"darwin-arm64": "76ed801a1e45a6ba0cb9c7c611095dbc324387a56d251f39b654c4475bd4199a",
						"linux-amd64":  "4aab4a4496cfa9de6ec3d7f7af79eb95d66cdb3eb47dceaf73c3e842c3eecf29",
						"linux-arm64":  "aad1fee52d66ae34ca7fb6b7c79757d25b41ef2ae98ea51e610ce50a04c6bfd8",
					},
				}),
			}),
		},
		{
			LocatorWithResolverConfig: config.ToLocatorWithResolverConfig(config.LocatorWithResolverConfig{
				Locator: config.ToLocatorConfig(config.LocatorConfig{
					ID: "com.palantir.okgo:check-plugin:1.18.0",
					Checksums: map[string]string{
						"darwin-amd64": "ece5ce450507714919a14522fe015db97bc0d696ee14a050fa78227793af7cfc",
						"darwin-arm64": "0b978f5fffb498d0167f1e99b5681ad1a0e5a8f1dfd5a24d843fbb4bd2e1aedd",
						"linux-amd64":  "3780d7dc058b30fe2c2caeb3258bc26efaa292b6316db50c880b1732cfccf477",
						"linux-arm64":  "dc709f28622ce8866e3f65342f49d176f8e9892962e175dd3ef9f89e61426c3e",
					},
				}),
			}),
			Assets: config.ToLocatorWithResolverConfigs([]config.LocatorWithResolverConfig{
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-compiles:compiles-asset:1.16.0",
						Checksums: map[string]string{
							"darwin-amd64": "deb36ca4a9b2a66a58846c7ab88408f478dc9f6aec035eebdc4408287cd35692",
							"darwin-arm64": "9ae6f4bf20ae88e1932fb15a12bef86b5192f08f25b1942ef5840f2deb40f607",
							"linux-amd64":  "833f09e3d0517a8576f6fcd3ac42c6b684033ce6618d96fd1a9f96a21d820a99",
							"linux-arm64":  "63ea946004d7b535c247814bd13f23af72fb5d353d8a2d0aac01bfe85409402f",
						},
					}),
				},
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-deadcode:deadcode-asset:1.14.0",
						Checksums: map[string]string{
							"darwin-amd64": "efe3f1c304ced0de5d4e8fcf51ad49d0989536d527bb7cfa193801d4bae88b0a",
							"darwin-arm64": "d1eddf11a81841b38c04168040847c3f044297f79c904527c0ad6239307d0e71",
							"linux-amd64":  "fa78e382130d7a8bf0032b612af02601d1c7f9e1c7c3b6b589ecb30b8c8b0494",
							"linux-arm64":  "917080582b0dd8764d218dd86d2b977238152830a4006d7fc726169599c76a4e",
						},
					}),
				},
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-errcheck:errcheck-asset:1.15.0",
						Checksums: map[string]string{
							"darwin-amd64": "ada80094dbcb2cb98dc749e5281708fbdaf86a2249374c2ec4e051f1fe4bb7f3",
							"darwin-arm64": "0950439307606d9c0acb2bbef3bad4b77070cc7a4815667b0e35db8782819782",
							"linux-amd64":  "74319b117fe19df8ffd679c64642bb896e1a10613f012c2634a6b4abdf6476de",
							"linux-arm64":  "200e255982326747a34ef03126cd12aae33ef1f8dacba40cd6bfd6799241eded",
						},
					}),
				},
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-golint:golint-asset:1.10.0",
						Checksums: map[string]string{
							"darwin-amd64": "e9f0f2e22f5831cdcba9556fd51a5a46937b13a29de4e2ade755c6afc90c7061",
							"darwin-arm64": "1491520d1d8cd6c5e6335771ccd6c9f096613cf13a74ad051c9c3bdf2613c7ed",
							"linux-amd64":  "efd00bca07d15e7e95946381f8f8140082152ac59d9bf3550976c97b2ab66d3a",
							"linux-arm64":  "98f2ad9beed6bb8951c61dac1cc4d5936c4c973fc3e7ed297a8c8bfbef178b34",
						},
					}),
				},
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-govet:govet-asset:1.10.0",
						Checksums: map[string]string{
							"darwin-amd64": "baf0b0f8a938b2ae2bcce27d5e4d5c02e8e6b249feae16ee15c399ac728d6ad4",
							"darwin-arm64": "0333fb2aeafbfa7386d7ae5ef2dd2b4a26949f4db3aea06e3348f7db340f3a5c",
							"linux-amd64":  "754b47727cb740f62d50e262854d090d45e7e79c731b52a8a30e595aa13d98ea",
							"linux-arm64":  "906f2489b9452db0eb65fb047f3e675e24f7d9dc4510ee12fae3b3b7f58e43f3",
						},
					}),
				},
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-importalias:importalias-asset:1.10.0",
						Checksums: map[string]string{
							"darwin-amd64": "dd0d02b2bce7854a86f900352539ac79c202b63d074f3a80d30b312f759a11b7",
							"darwin-arm64": "3251b7bec1f2b8acc5555a8e13b1c2100a0489b4ea10cc405b3e5fbcb6646e32",
							"linux-amd64":  "51cc43e305be2070e4c0cc8067c217b97dc07b2542d773fb9027c8af8aea3343",
							"linux-arm64":  "d7b9c14acab970ce076b4cdc50a519d446644b41cc14100b74256ca7afcf0402",
						},
					}),
				},
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-ineffassign:ineffassign-asset:1.12.0",
						Checksums: map[string]string{
							"darwin-amd64": "97bf584c929e1604185041f054cc5132ca6080053b15bd6f68c2fee7daf84518",
							"darwin-arm64": "24a80dcf9de7022a02a246d8e81c7f96a2ce60846f971c94f5bdd0e9671d06d6",
							"linux-amd64":  "cc05558caff5c25f4e2f819003a554d62a354e457e3f533cc48d3655138fdfba",
							"linux-arm64":  "cfa9ec7bd3dbc94becad928efa14726916cfa78c0d2c92c04fe91eeb75fbf807",
						},
					}),
				},
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-outparamcheck:outparamcheck-asset:1.15.0",
						Checksums: map[string]string{
							"darwin-amd64": "980e00197bf4637d595d55137579521323cf292f9c66c0e953b29bea228e585f",
							"darwin-arm64": "54798b5748b1ba4f4fc2b9f57f4c8cf670f709b4c9fab2e6ae634420f90e932b",
							"linux-amd64":  "a23c9c0fcf285ef08227f9b9188962f6830c0ab361a54717a380f114b48b8bf0",
							"linux-arm64":  "b5e70c3e7b8b31cdbb6f0f48eb6c02b88a4f8d0f2c3dd028f54515297b26775f",
						},
					}),
				},
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-unconvert:unconvert-asset:1.15.0",
						Checksums: map[string]string{
							"darwin-amd64": "99c0a729db3eba73746c53d3c5edf1dc41c582e5ae32e6cde8e8851ae6a63a58",
							"darwin-arm64": "41de2eccd7e4791e46169036c0f748e31a578552c3a14da7fe22ae40a3d5cd75",
							"linux-amd64":  "c54d80bd9aaeef8b686bb61aa419e3961baff2e92f22e4cb3616d05821761c54",
							"linux-arm64":  "3b054304bb8491f0d8aca18873d0df1890f1cc6702f0533b4c9be8f1a2c6edc3",
						},
					}),
				},
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-varcheck:varcheck-asset:1.14.0",
						Checksums: map[string]string{
							"darwin-amd64": "34e16f85175889d0b3e9b4515dc83dd44671381558ebc32a24a417863b7e66b5",
							"darwin-arm64": "589beeb5beabcf0fc9e29be83f7b3289ba13a4c29083fe381a7c3ac3d2cbd747",
							"linux-amd64":  "9a0045a7a98e6a1475c2c30bafb3cb6da96424077d570bc957f553e2268c35c8",
							"linux-arm64":  "00631b498c94b3c1ad630f93e2b53ab62651dd9bc8a39d846030591e2b129cd0",
						},
					}),
				},
			}),
		},
		{
			LocatorWithResolverConfig: config.ToLocatorWithResolverConfig(config.LocatorWithResolverConfig{
				Locator: config.ToLocatorConfig(config.LocatorConfig{
					ID: "com.palantir.godel-license-plugin:license-plugin:1.11.0",
					Checksums: map[string]string{
						"darwin-amd64": "caee8215e6b1c9d3b3841a5494a4c8966bccf780a49995cac23541fd309361d1",
						"darwin-arm64": "99b3ff08cddcd6a9f5ec0324d073e7e7108cdf0a046628ba8a51d23272491582",
						"linux-amd64":  "d31cc808f419e7879d2a96bd1080c03bcae3141242c2b3f8c477410a32b81d43",
						"linux-arm64":  "1c699fd1c99d64ccbfe6fba5946529141280c32fb93e9d4b392f04f6bf6fdf57",
					},
				}),
			}),
		},
		{
			LocatorWithResolverConfig: config.ToLocatorWithResolverConfig(config.LocatorWithResolverConfig{
				Locator: config.ToLocatorConfig(config.LocatorConfig{
					ID: "com.palantir.godel-test-plugin:test-plugin:1.13.0",
					Checksums: map[string]string{
						"darwin-amd64": "698e5c1a8fc047e5c213f5b4c1e7991197e1562b6ae6eace928d25fb1219cfb0",
						"darwin-arm64": "b81b0a04ccc9a9a7ccd957dcf8ea2fe265e53473534006acff7812eaab38c348",
						"linux-amd64":  "6fb8e8d277e7b8baecb54524da738ab03bf5c1b355da6f1e9c3a1996c325fdb8",
						"linux-arm64":  "3679e63570692864f35931153746e6ff23181194bfb1c29d613ee3caf9f09bab",
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
