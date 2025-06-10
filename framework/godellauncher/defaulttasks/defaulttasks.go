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
					ID: "com.palantir.distgo:dist-plugin:1.84.0",
					Checksums: map[string]string{
						"darwin-amd64": "cdbb61ae72d9e86c6bc2b41c473e7dfb257308d9702c31e5391bc74bcffbc2e0",
						"darwin-arm64": "e101120a50862b82435dfd2068f125fa80b76f927930ed285dd3c489007b6d35",
						"linux-amd64":  "30066ebc2f628a302648c4507d720a652aefc4ea6ef6694f50509bc3309afa47",
						"linux-arm64":  "9506ff5d898ff4f91a8ece0d643e8cbc1488300546532eca34e2c00254ee4f96",
					},
				}),
			}),
		},
		{
			LocatorWithResolverConfig: config.ToLocatorWithResolverConfig(config.LocatorWithResolverConfig{
				Locator: config.ToLocatorConfig(config.LocatorConfig{
					ID: "com.palantir.godel-format-plugin:format-plugin:1.52.0",
					Checksums: map[string]string{
						"darwin-amd64": "f6ddadd74d9faa407c7a40195986b632bd69c2ef75db7af9ab21f8a8fbca8032",
						"darwin-arm64": "a5877536ccc8b30cef2e51e207e9239e5b3623c77a2ac4377b9305e05df015f7",
						"linux-amd64":  "961051d2f1a9824247a640f92c6d61bef203d650df40c4f3f5aa2689560c80a1",
						"linux-arm64":  "9b6e4da44ca36bda389e527b8981ef1a2dee892718d6521526923180661e79ef",
					},
				}),
			}),
			Assets: config.ToLocatorWithResolverConfigs([]config.LocatorWithResolverConfig{
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-format-asset-ptimports:ptimports-asset:1.51.0",
						Checksums: map[string]string{
							"darwin-amd64": "a20cfea8eac2475ac258814ef33ecc8f992264a5683de54746607486e4e5c9cc",
							"darwin-arm64": "35d1909924aaaae72912bf3916addefff56b82c3c74f693d6deb3ec6d82520ca",
							"linux-amd64":  "d5abdede79f1acad49f3f6f81750028059a284fb59b617a7d2d0bb2a1bfb21f4",
							"linux-arm64":  "d177f3667b426946f21c4a60cd31e8b59af8e26b8ec73064c93e505136338dc7",
						},
					}),
				},
			}),
		},
		{
			LocatorWithResolverConfig: config.ToLocatorWithResolverConfig(config.LocatorWithResolverConfig{
				Locator: config.ToLocatorConfig(config.LocatorConfig{
					ID: "com.palantir.godel-goland-plugin:goland-plugin:1.46.0",
					Checksums: map[string]string{
						"darwin-amd64": "77034a1f7fb8fa07ee66889ddd3dbf9e1cfdc0493ea9e3aeb7b62ad85d4eac08",
						"darwin-arm64": "839da1df87d8271d86b4bdb254b6bb496b207432c3ae881e27f8caa1d9202a83",
						"linux-amd64":  "3505c9ddbb00135cacd2c24257e753b2ba53a46fc03f0de79d10dad4f08cfe2d",
						"linux-arm64":  "933331716d0ddaa002b6125b34d9295b6c35e366611d085b0d9b3cbd34468b0a",
					},
				}),
			}),
		},
		{
			LocatorWithResolverConfig: config.ToLocatorWithResolverConfig(config.LocatorWithResolverConfig{
				Locator: config.ToLocatorConfig(config.LocatorConfig{
					ID: "com.palantir.okgo:check-plugin:1.63.0",
					Checksums: map[string]string{
						"darwin-amd64": "62272db913c3a8cda19e018ac7e11356564f221e037fda705f85f3ec74e66b2c",
						"darwin-arm64": "7e2fd26a1c41e0ef25617dcf2c6d9052280feae33c4b0301871cfc277ec73dde",
						"linux-amd64":  "8c3e9a130823c2c76c69572f7201ee131cd2a1631c35f5373dfe1a2767992319",
						"linux-arm64":  "31459e95bf1b06bf8b8b87d7f4ffffeaa2380cf3af84ebd734b9b98dff2188af",
					},
				}),
			}),
			Assets: config.ToLocatorWithResolverConfigs([]config.LocatorWithResolverConfig{
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-compiles:compiles-asset:1.56.0",
						Checksums: map[string]string{
							"darwin-amd64": "8600fd7c591e7242e45f942d1fceddc668ea22245dd94fe6d2e79cc03ca72df8",
							"darwin-arm64": "938d2e301b2c1b912b22ca783a8f55b8b83e8167dd6af607d171d6952b1db5ea",
							"linux-amd64":  "dda85d839f1d19d2789da56f2f99f9d9d8a2a2a96a9a98cf5696b7ee27f90f20",
							"linux-arm64":  "0c177c8a86ff1a05b7f3fbcf96d63051bfec20a310a69238b22dc69018e14425",
						},
					}),
				},
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-deadcode:deadcode-asset:1.49.0",
						Checksums: map[string]string{
							"darwin-amd64": "80c8588d17eb5072ca305fc21fc2ada48868daf55bcef21954a7b8e6131d5f6d",
							"darwin-arm64": "112013f6ca169bfb7e7dfd7540e6c0c072b258dd14d06274edb5e8798bc7d9c0",
							"linux-amd64":  "57d14b464632314221feeec6bdee75de74be69e7591b7f10961a8054e1c80414",
							"linux-arm64":  "d470baa25c5c965bc4a8a7c142fcafd36886a1ee9d45d6e13412fbd56a93b2cd",
						},
					}),
				},
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-errcheck:errcheck-asset:1.51.0",
						Checksums: map[string]string{
							"darwin-amd64": "599df7ca82ab6c805511df2277954c6bd186cf8742fb3a46157a67ebef810019",
							"darwin-arm64": "d457c3f89718d09a12ed0fe177ad2a2b09e076cf22dfe9ffb61d6b171ba6c81b",
							"linux-amd64":  "0e4bd4fccba4e3e62d1b34361d1daf3cc89409f8ef015c25fad3b2925364994b",
							"linux-arm64":  "0a3bf1501d64b49a819c1e3f5fe59015d0f1fe25238e4541dd63958445a65fc7",
						},
					}),
				},
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-golint:golint-asset:1.47.0",
						Checksums: map[string]string{
							"darwin-amd64": "a126f0eec1a4da17c8dc8bda0ba04275160e262cf4a401bee003fe1e7b465bf5",
							"darwin-arm64": "eaad88da5f85fbc56c850452250fbecd8b43eca163159ca97cc0568a2384592e",
							"linux-amd64":  "262cf63cad62fab17a247cf450721ee6d32e0c67e9709f91ca6afea8cd152642",
							"linux-arm64":  "9be0380f112545407d87f55f443d63533c1f81dd5a38f3866fbac9807a9bde20",
						},
					}),
				},
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-govet:govet-asset:1.46.0",
						Checksums: map[string]string{
							"darwin-amd64": "2344637eef7e2cba4363e34a9bbc6c7951a63230d2df75906fc47da7c6a0cab1",
							"darwin-arm64": "46b3c0956c48298942e7729e3e4e1bf6ac173ff5fae3146af7e7081ee4aa69ad",
							"linux-amd64":  "0579f00282287cd9d75e0255185b285f0b079f37eed8b09c4d8020771194e79d",
							"linux-arm64":  "111b64cfcdd097f62d4b396ccf4c4d1111318101373d45608122ecf60f477cb8",
						},
					}),
				},
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-importalias:importalias-asset:1.46.0",
						Checksums: map[string]string{
							"darwin-amd64": "157253e4fd0c2400221c8b5460e8a54b8dc16210f21e5efdf9d6ed4ce0dce872",
							"darwin-arm64": "236d53bfdbdb080693bd661f67a867e30b2b682c717a4e722ffe13a8c527ad20",
							"linux-amd64":  "a25ffe877f15ac54dd3b23be1a01b02f8a37a7fb8b3bc41f3a9cbbb4b4303a31",
							"linux-arm64":  "7c1a483742fc64a8bfdf3c3a07f8a79312b7563f47822988ccb8a24c15d4ea70",
						},
					}),
				},
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-ineffassign:ineffassign-asset:1.46.0",
						Checksums: map[string]string{
							"darwin-amd64": "e75af5b82698da5956284c22e54cd0c80d19d2b1f8467e067d9ddb5668ba5888",
							"darwin-arm64": "6421ba2148f303c109989d0621b57c4f485d7e475c37d93b2cc74c0860af75e8",
							"linux-amd64":  "5fb155c6a64d07738c6a0f62d1e6e1e7ac4c88cb118489fd31b16ad7577a2461",
							"linux-arm64":  "060a2fcbd86e51c703f51c2c5a4eb45aa9861b873fda0a691ab35e5f340a251e",
						},
					}),
				},
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-outparamcheck:outparamcheck-asset:1.48.0",
						Checksums: map[string]string{
							"darwin-amd64": "9adee2605720bb15a3c01678d9f3f7cf33ca10c55f0c0fd0e4601c533aaf8956",
							"darwin-arm64": "d9275cddd5508eb25343447f977b430d3d9f9d18a6f352eb2130e00eb9731293",
							"linux-amd64":  "5cd24bf11d0c23520e7f736875286e80ddb2ab3a68c79061b27b51cb3bf36ac1",
							"linux-arm64":  "a0e35fa675b1fe770d2a76fd8ea927e6e36efab71f9214b879bd8dea522a8f3b",
						},
					}),
				},
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-unconvert:unconvert-asset:1.49.0",
						Checksums: map[string]string{
							"darwin-amd64": "585c805298335e454af5371cffe22cf68ba44ee856bc10f2c0ea82676a8a588e",
							"darwin-arm64": "05cdfd1d5b346d971f5d56f12cb7025bceb97dc4dec6754c35aba930b2c0b42c",
							"linux-amd64":  "97ff7355d7cd4fa26636f33045fcca953c2a99791183888d159bb8343dac657a",
							"linux-arm64":  "c391630a4dea43de52a5a4749633907fedd5e145585a5f00bf6dbb09811c4fff",
						},
					}),
				},
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-varcheck:varcheck-asset:1.49.0",
						Checksums: map[string]string{
							"darwin-amd64": "3f3005a3bad62cd7a27292d14b42a0955eee066d931653fbb7b4836a13a66e78",
							"darwin-arm64": "7c0f7d6c808595795e9c7077bdee0115e4bf89b4520c0e1487c59dbe2ba42822",
							"linux-amd64":  "ac942afa4e81fbe0082414dc0b3d73bc03938729c2a29223d4aaac02d6eb94c2",
							"linux-arm64":  "b3f8a75a026e4683f72e3eff0e5ad4f22f8ab825cbb661276c933f34deb4de6d",
						},
					}),
				},
			}),
		},
		{
			LocatorWithResolverConfig: config.ToLocatorWithResolverConfig(config.LocatorWithResolverConfig{
				Locator: config.ToLocatorConfig(config.LocatorConfig{
					ID: "com.palantir.godel-license-plugin:license-plugin:1.48.0",
					Checksums: map[string]string{
						"darwin-amd64": "c691fc6c5fdb12adb141eddb8c1f647c451f3b8f2da58f1d034bc7f174fdd4b6",
						"darwin-arm64": "4bd15a222980cd207b56264395c2542848cf6554bc4721145ca86759d494a987",
						"linux-amd64":  "4b5d308f0566cf995bf6ff1d85a246169cb57a51875fb7be099eb6433cf3a3df",
						"linux-arm64":  "e74b7731f8bb38878d1982ada3b5abd13604f865197e874ea2dc40e944f6f896",
					},
				}),
			}),
		},
		{
			LocatorWithResolverConfig: config.ToLocatorWithResolverConfig(config.LocatorWithResolverConfig{
				Locator: config.ToLocatorConfig(config.LocatorConfig{
					ID: "com.palantir.godel-test-plugin:test-plugin:1.46.0",
					Checksums: map[string]string{
						"darwin-amd64": "531dd9e60fa72e664285f476ee1f22c221cbbe73838a73e1ad12517b50cc1b3e",
						"darwin-arm64": "e1ffb0ecf27a84e2aba2ae2e6a8a3ba86c55ed4531db2320d27366eece7d5cff",
						"linux-amd64":  "294d46c1a95100bb583b1637e4719cb5c8867bdf2a7c70b2ea4b33095cabc037",
						"linux-arm64":  "1692ff583cbcbb2156bcad5c94881c857773ba21e8a5d609e361ddcfdbb1a268",
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
