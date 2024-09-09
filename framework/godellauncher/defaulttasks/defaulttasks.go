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
					ID: "com.palantir.distgo:dist-plugin:1.76.0",
					Checksums: map[string]string{
						"darwin-amd64": "78be7634f0f5b32d76aa7361c4a48cf53321c598beea69defd1628e8fd738382",
						"darwin-arm64": "46d70c56e4fcb1fd1e747d9599b0c5ad3e70b40bcd1ba1ff4d0e973f1bd767e9",
						"linux-amd64":  "524d46233357b5c4f0f043aaa9822c6f735916c25289f2395dde6451d6442263",
						"linux-arm64":  "e46b1d6e5676a3d56ae0884b1faf915454618d233bd3b6017c577b7f71cedaf2",
					},
				}),
			}),
		},
		{
			LocatorWithResolverConfig: config.ToLocatorWithResolverConfig(config.LocatorWithResolverConfig{
				Locator: config.ToLocatorConfig(config.LocatorConfig{
					ID: "com.palantir.godel-format-plugin:format-plugin:1.46.0",
					Checksums: map[string]string{
						"darwin-amd64": "b5210fc5e6cc2a456fb8f42e42de2622e6f216af5981708be02bdba1c0774563",
						"darwin-arm64": "67f01069e73011130ed942ddee9ac3f3d6ca0eb1f8f924c48ee0f78657283b8a",
						"linux-amd64":  "694be1e75e492e8ba3ab1348ece542729fd1b08c6b743594ed3dc579ce1a80ab",
						"linux-arm64":  "e874fad3de2f51ba9a4541c61fd0d1e82919f8434c471fd4e8f1e8361cfc5043",
					},
				}),
			}),
			Assets: config.ToLocatorWithResolverConfigs([]config.LocatorWithResolverConfig{
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-format-asset-ptimports:ptimports-asset:1.45.0",
						Checksums: map[string]string{
							"darwin-amd64": "749c7fe7af60ce3b8308e96875d1855047bef588199679175d2fe904d529af54",
							"darwin-arm64": "79928dbbe091ce3a53c38fdfa758b2205bf44b2628e2e571dc7bd6220b439c1e",
							"linux-amd64":  "7418e6576e0611aac590af1f75556657b1a83eafa23ce8dd36953bbd91252092",
							"linux-arm64":  "b988c6a4658b825a95359767e47033e1e72b7245e3d443681330a84dd33b13e1",
						},
					}),
				},
			}),
		},
		{
			LocatorWithResolverConfig: config.ToLocatorWithResolverConfig(config.LocatorWithResolverConfig{
				Locator: config.ToLocatorConfig(config.LocatorConfig{
					ID: "com.palantir.godel-goland-plugin:goland-plugin:1.43.0",
					Checksums: map[string]string{
						"darwin-amd64": "7a6e85ecb8d8be9db012d58cd033cf6ec9a1f092f18aa56e6ed7756ac23552a3",
						"darwin-arm64": "793ed13aab10ef0c9d48e2707e5714dc94b10765ffd000c91316d2387b427ead",
						"linux-amd64":  "31afde7c76bbfc58db4b1bd5b8aa24219b41b9e31411342a9a3dcb6624d34f2b",
						"linux-arm64":  "44afff8b5d59a89a073498685f62f3434ecff5c3664520cc71512fc1ab39ce2e",
					},
				}),
			}),
		},
		{
			LocatorWithResolverConfig: config.ToLocatorWithResolverConfig(config.LocatorWithResolverConfig{
				Locator: config.ToLocatorConfig(config.LocatorConfig{
					ID: "com.palantir.okgo:check-plugin:1.59.0",
					Checksums: map[string]string{
						"darwin-amd64": "97d2080f8cb039f46de1474452c0abf3a4e5becb37485aca50878ecb68ad1dc8",
						"darwin-arm64": "0d54796ece319fe57f8974c723796aec87cd32c5544bbd486bf1a223ec48b629",
						"linux-amd64":  "debfddf54935bd09db8dc0686b7d71ad85ce8af4a42a21daced36064fcf5062d",
						"linux-arm64":  "cbafbf938a87750162977fa37255a9bfba7aa6aa2e2e3f980d8937eac8a11201",
					},
				}),
			}),
			Assets: config.ToLocatorWithResolverConfigs([]config.LocatorWithResolverConfig{
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-compiles:compiles-asset:1.52.0",
						Checksums: map[string]string{
							"darwin-amd64": "368cec200439e2e075b32c2291478a164e2aad07e228e3f98649cf447e58c3db",
							"darwin-arm64": "9d543ae5b3411e5ba869a80f5c42b2c2fda10dc116bfc0000f052aa013f2fb16",
							"linux-amd64":  "58f39229bd1d6b8962daf5afeeb9c11aaf7f5036f6eb22a6788b2d3275dd2520",
							"linux-arm64":  "b73e7e34a58234bac533a2ce0b58ba946d1a368e136347996ddc4b6c2d81c41e",
						},
					}),
				},
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-deadcode:deadcode-asset:1.46.0",
						Checksums: map[string]string{
							"darwin-amd64": "2b260327ae7104f5ee2d5735751723a46d46e8d2ad7cfe6cbe953bbb0d901d68",
							"darwin-arm64": "208b26eb71a170c9ef5b0a2869dd7502fb67db38acdc0041bd2a5102886061e4",
							"linux-amd64":  "54e4004130c6ee0ba9924342f7e1531d8a6b3cfb2f6f26e4f2765d1ce7040435",
							"linux-arm64":  "b8afb2dbb8a90b3843fd37c465fc6b5d2d17cd38262ac64c7b8c4bd69bd640c3",
						},
					}),
				},
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-errcheck:errcheck-asset:1.48.0",
						Checksums: map[string]string{
							"darwin-amd64": "83eabe44b6f8ebaa5530a8e290f8e2c4277d56ec01b86cfa9d76201569c94249",
							"darwin-arm64": "4ff6cf82578d2152514b353e1cde83e1a1a76d3630c840ee70aeea93a2062eea",
							"linux-amd64":  "3cab2bfb07302d808478d806d1f46c1e32a87ea832dac61f458140c2a3b6b81e",
							"linux-arm64":  "63f97ac1d6aceb7a94f2cb38fbd034e1ae76ba6f1f9a005acc854f54dc549ceb",
						},
					}),
				},
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-golint:golint-asset:1.38.0",
						Checksums: map[string]string{
							"darwin-amd64": "9ad7e4265142d362cf54734ac203216270ec3c7a310accfd1e25259b10867916",
							"darwin-arm64": "beb75f96cc064f95fe19832bc58e6d8eeb89192d8683d19e4687f1b3906a116a",
							"linux-amd64":  "e84a94ac7f211042eeb37e2371f39a24c30935446a82a3993e177f3e7f06c6de",
							"linux-arm64":  "dd79fb3168fd5604b5da38c3008331ea5828a04226ed260d42650fd78603b444",
						},
					}),
				},
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-govet:govet-asset:1.43.0",
						Checksums: map[string]string{
							"darwin-amd64": "179f0dcf82af53ed7539774e7388efa2dc744798a0d1127e318df5bc775e64c5",
							"darwin-arm64": "631e1ecbc08a50c4c38ea9f079f1861e47b80049961d1ebe1411f67c0d25c13c",
							"linux-amd64":  "998fef33755658bdcb4a3fb2ca6dc9a71f13582cedf0177857add6070e440ae9",
							"linux-arm64":  "20f21ac0b20586fdf65a6fd52b291302d276e300db200117ae59d86c1afa0e1f",
						},
					}),
				},
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-importalias:importalias-asset:1.42.0",
						Checksums: map[string]string{
							"darwin-amd64": "8e02ddbbf805ecc7e3b8a377628766816e7748771f21b83688c2778f0af59bb5",
							"darwin-arm64": "d19989a6fd9e14972faffa7006ef1dcacbcf0283ea0620f666178fa0b8c8ef09",
							"linux-amd64":  "0cfdcfedb86479d380206a07155df6838d59a9008f28e448890b1e8104a115e9",
							"linux-arm64":  "9074da12d06f7e333e1adb78cd15e524773dcc55967545bdae797badf29dab20",
						},
					}),
				},
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-ineffassign:ineffassign-asset:1.44.0",
						Checksums: map[string]string{
							"darwin-amd64": "876585ef870d0b5b158be54ad3f1887ec75f07681ed0d0644c0f3b719ab7ce73",
							"darwin-arm64": "53928f0c3165bc89662a517927c769cd200d26309fea5b3d908f79fd279355ea",
							"linux-amd64":  "c2b42fe9d835dbad81de0481240080c474e45d9d73e97b1319fe64687987af16",
							"linux-arm64":  "251a2b0c48dc5ad7e005cbd7d17ef9be4be87f715b41ae15d6fd5bc5b55e24d6",
						},
					}),
				},
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-outparamcheck:outparamcheck-asset:1.45.0",
						Checksums: map[string]string{
							"darwin-amd64": "0c90072217bf31ccab8972fe0684e8e01c458673075d6b18be1fc157550b2838",
							"darwin-arm64": "7b7f84628805e0585da79f88511fa1f054e622e5b5de956cd068ae1cb6b1f2d1",
							"linux-amd64":  "17adcd7698e90d7fe3063e3c9243586abe470e29c7a9495c4ba44e3a0117426f",
							"linux-arm64":  "6f988514015cfac7be674f88407867e8d970c7ec1d4446c06ca8554bfb4f42b2",
						},
					}),
				},
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-unconvert:unconvert-asset:1.46.0",
						Checksums: map[string]string{
							"darwin-amd64": "37a9200cce87723bd8c8f921898393d83f0fbe1d41970d5d10575a1883565039",
							"darwin-arm64": "4ff5c078412ba52645e43dbbc8f3adf5257fbce096dfc2e60fdedc12ac245b42",
							"linux-amd64":  "35ad1b373f118584d186c72143c427d49e09ef4e97f4879a57f656f2109953fd",
							"linux-arm64":  "5be71b66d0f910fcbf8fb591a24df86990a7c175198d49027a5024559cde2339",
						},
					}),
				},
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-varcheck:varcheck-asset:1.46.0",
						Checksums: map[string]string{
							"darwin-amd64": "fb865167b2653cd70034cff57181cfa6a0e3dbc345fbe9b89a733f5f1b35afde",
							"darwin-arm64": "b82b9cf72d070ee22d04e5b52e6352534eeaf07f45b927d5136bd551518db911",
							"linux-amd64":  "ac4483f5b70ec6a28b8644f84e479c7b2de420953c3d3288c0a733621028ed02",
							"linux-arm64":  "89de846a649a17968f0863d15a8057716eaa1e42e369ba2af263b308c723272c",
						},
					}),
				},
			}),
		},
		{
			LocatorWithResolverConfig: config.ToLocatorWithResolverConfig(config.LocatorWithResolverConfig{
				Locator: config.ToLocatorConfig(config.LocatorConfig{
					ID: "com.palantir.godel-license-plugin:license-plugin:1.45.0",
					Checksums: map[string]string{
						"darwin-amd64": "f789038f283664077009d3d291b5a2e00e434401c9d0ec0ebdb2a60a6d574a99",
						"darwin-arm64": "907dfcae5a8f96eb3edc9277e02311da33e4845a0d270b00b86aecdc6a28a349",
						"linux-amd64":  "0dd92fd3f74c3d7c48be1cbf86de2b885535a8a784a72a383ca0bfec9517f212",
						"linux-arm64":  "a3dd2fd1bbaa9304d8b3d11a141b418400cd4ff00522e09249d22853ee787efe",
					},
				}),
			}),
		},
		{
			LocatorWithResolverConfig: config.ToLocatorWithResolverConfig(config.LocatorWithResolverConfig{
				Locator: config.ToLocatorConfig(config.LocatorConfig{
					ID: "com.palantir.godel-test-plugin:test-plugin:1.43.0",
					Checksums: map[string]string{
						"darwin-amd64": "1d801a56057a6085fcf07248ad29941162655d6179791492e69cb9aca517b9de",
						"darwin-arm64": "e55f82c626743646d12cd39e32cf1a3ded24df4254a3d7517a47abf28f2835ec",
						"linux-amd64":  "80e5c101c9384f9bebdaf01e5e4db3719f256fd01190286cc95b0b56f07e25f8",
						"linux-arm64":  "4a0fbdc75a3be5c52b46aa43c4eb57ef0efd78ae209fde7260f0e7cafccb31d1",
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
