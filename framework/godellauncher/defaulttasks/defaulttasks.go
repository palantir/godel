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
						ID: "com.palantir.godel-okgo-asset-compiles:compiles-asset:1.51.0",
						Checksums: map[string]string{
							"darwin-amd64": "6c96de0277588b420bedd4666999f0a6458bf92aebd0b35c22a88d4799db0840",
							"darwin-arm64": "30fa2c3b4c2e901ebb0eb04887ead83da225fc1a14a5b5a386c8e8caca092fe9",
							"linux-amd64":  "725aedcf8024caf4c62cd16ee5e1b5cce504d5305439a9c6020be22d3308a328",
							"linux-arm64":  "e51388bd7be7afda9692c345bcb89633f50d92e3c887f66f901d0712e65d06a1",
						},
					}),
				},
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-deadcode:deadcode-asset:1.45.0",
						Checksums: map[string]string{
							"darwin-amd64": "d800ad95d45aa74d8137d7c78c649b7b9fe203f89c798b69e06380537c07be5b",
							"darwin-arm64": "c0ad7e7dd0437c82b98837305e7259dfd6db778260f6c8449bb548d6d90dc53a",
							"linux-amd64":  "41d1c816de8539bb76b76fd4a4f925469123d7fad8ce3ea4df46eea4259fa596",
							"linux-arm64":  "44586bcb2e6d102c6fe378bb0f6de2505e345811ca537b9c4d5e9b242c6416e6",
						},
					}),
				},
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-errcheck:errcheck-asset:1.47.0",
						Checksums: map[string]string{
							"darwin-amd64": "71d6999d7a90a01f2c7b6107fe7fd838a0606e9ab7fb5cbe5a62198c13f374d0",
							"darwin-arm64": "b0b300d5c7a578f13febbec5a6109aaf39544d306d05197a5a817ebec09eadba",
							"linux-amd64":  "20deb54278ab31d3ed442b0641fdc8f7e0318b27a4709e2419c9c54c2b8d0818",
							"linux-arm64":  "7c1d2c18cb574ade0e96998aeb69fcb5a9594be1cfee6a8b5b287b8c919d4ad6",
						},
					}),
				},
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-golint:golint-asset:1.37.0",
						Checksums: map[string]string{
							"darwin-amd64": "7e179bb625a2dc9aa0146a7d67bfa6c97073ad1123fd9d3afafa3f9a5473df97",
							"darwin-arm64": "87ff384bd0b8b6fa024ea53c139fbd6efc3034cfe39835fe26b5cce683696c98",
							"linux-amd64":  "454405768d2e33c1a9982bf6d27c3f09612fed2c44797830ec8ce0db26813dff",
							"linux-arm64":  "4c81969282accb04f7656e1caa866a51fdb904aa73a689fcd8af3248c48411e3",
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
						ID: "com.palantir.godel-okgo-asset-ineffassign:ineffassign-asset:1.43.0",
						Checksums: map[string]string{
							"darwin-amd64": "e12662997f3d6bc54144eb77ca8b31c5e97ace753e8f34b1c6c7804b8dc3ca56",
							"darwin-arm64": "946e970b8cd1b60f02eeacc49847ec3317990a161d07b9d927ae42549ade28d8",
							"linux-amd64":  "7640e562ee35ad65929522386ab7245d1a51cd68a75b2b4075a3fab992d34113",
							"linux-arm64":  "0edcbf537dfd5a2dd7478ee3cf439e50a8dc0512aaf503af976b47446f4458d9",
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
						ID: "com.palantir.godel-okgo-asset-varcheck:varcheck-asset:1.45.0",
						Checksums: map[string]string{
							"darwin-amd64": "da77d28095d5355efc80701489afe7388248069fb00a9574a5a066aaf741ef62",
							"darwin-arm64": "45a4bcb158284ad6fdacc6d0acdbc530b579e54755bb264544b1014ea9621ca4",
							"linux-amd64":  "abcb7cd750debcea5230307b0935cb850a6308f4ae316df006ac9dc2bd828713",
							"linux-arm64":  "ce010e1cc472efe891228457e0d086ea8ad26a6eb2b5b1fbb3f660ca5ccc25cc",
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
					ID: "com.palantir.godel-test-plugin:test-plugin:1.42.0",
					Checksums: map[string]string{
						"darwin-amd64": "cc7c3294e6dfd32f4f116f1f91b560dadb599b20e4a879ee4e73fa1bdbf3ca5e",
						"darwin-arm64": "0d373b10743095e90a52fa17ebe3f35cd79ca24709fca1f599eb93fe1b925912",
						"linux-amd64":  "f1be4fbfb707a61a7a6fad6cea0cec00db76c358e08851c25758fd115c030d47",
						"linux-arm64":  "512cde17fbba7070bd1dacc01a52ed72c0495c46c765695d53c7a4a0298b5e42",
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
