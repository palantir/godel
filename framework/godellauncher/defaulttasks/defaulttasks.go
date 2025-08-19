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
					ID: "com.palantir.distgo:dist-plugin:1.85.0",
					Checksums: map[string]string{
						"darwin-amd64": "e2a123516e4916936438b49a1bd13a631a83ef6857dac97656587266c3525123",
						"darwin-arm64": "7918a6119f7bdfbff065b06d5ff1fddd99282f11eaac2c011361fc1a957c6dc2",
						"linux-amd64":  "054ba9eeb08436c8688d277bbe087bca0f4fd3572ded4d0c66dd3323231ee5e5",
						"linux-arm64":  "a295c355cc69a11076c1b6648fdb49bc4d0f92cb51623839d993ce1792fb767a",
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
						ID: "com.palantir.godel-okgo-asset-compiles:compiles-asset:1.57.0",
						Checksums: map[string]string{
							"darwin-amd64": "ca4cd2a0abb3de24db41f325979655431c2b97482a4677c60fe0249ddd8eea4a",
							"darwin-arm64": "d02462df9f13194bfa704307ae0a83dde8af67738293c9536065e999f7a9c368",
							"linux-amd64":  "4548840c06c64ded6f217c94bfcbd3dca77df1ac21e4d2aa47d80c3f82531c39",
							"linux-arm64":  "c58884ba4b4ecf0d45be8818ccd61e3b017c3c1e39b58385affab02c49c85157",
						},
					}),
				},
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-deadcode:deadcode-asset:1.50.0",
						Checksums: map[string]string{
							"darwin-amd64": "adc9ac7df1a7c5f3d3baf3ffc397948f87e027e740a2010c319addd0b9bfd5e8",
							"darwin-arm64": "5cea9971308be884c24df0caec946c5e069a9ba9fa612cac8c24c29307982ed4",
							"linux-amd64":  "3a4180c1bf870985a69bf2b9368c27dda14c346c81f86519fb8a511da4485264",
							"linux-arm64":  "56c6d0555ae944d32a97d565d8308a3125b7a689c0472dfe8e6c7cc4205fc2c0",
						},
					}),
				},
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-errcheck:errcheck-asset:1.52.0",
						Checksums: map[string]string{
							"darwin-amd64": "16d90ef7d2551916571ffee22cbda5c92dbf9087f3f82fbab2ee3cae579c15e5",
							"darwin-arm64": "edbc562c65c192c0104873a6914f23285462231e0f4b3444bdb51ae88fff9ee7",
							"linux-amd64":  "f8fbc0a5c0f4955cbb44514b8c49d6af1260242acee22e9b271ddbea932b0e7a",
							"linux-arm64":  "e13a2ed4ae923966a29733b81127006bb98f2574632a0743e35fac14d8972e23",
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
						ID: "com.palantir.godel-okgo-asset-outparamcheck:outparamcheck-asset:1.49.0",
						Checksums: map[string]string{
							"darwin-amd64": "4c70e575a8cc6189b8f58b794144d69ba46118e178f57807a27f28f8a6808a2f",
							"darwin-arm64": "7b23e07e39682dd9de7de7a0f0b09649af05c3e575b851fc67a5d4e03c749b05",
							"linux-amd64":  "ab7e6c45625757b5ca8be3399596f004f37c49745647edae91bdefea4385ec9e",
							"linux-arm64":  "3526832d049c76516b0b1246075454f5b0dfd941810edfa70a996496b02261a9",
						},
					}),
				},
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-unconvert:unconvert-asset:1.50.0",
						Checksums: map[string]string{
							"darwin-amd64": "dd429d8fb8eef0ad498523f00914844068b0d2198cca829e8a12f2872470f577",
							"darwin-arm64": "949530300f61e593b0c5e8353e2f7a1522d7fd9aaf4f1a42bc44ad5c2e17f0c7",
							"linux-amd64":  "c67dcf1cf30c0c2888a8a04e44ae53dfba02a61d4f7397e046042629edfb4636",
							"linux-arm64":  "72dc523547808317c7a01457ed2a6b34f99e5559718e0a548e70f4e0146a2332",
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
