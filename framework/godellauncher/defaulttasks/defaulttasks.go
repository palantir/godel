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
						ID: "com.palantir.godel-okgo-asset-compiles:compiles-asset:1.58.0",
						Checksums: map[string]string{
							"darwin-amd64": "91a866ffed7ce568ef935b59d14f66d9ce25d3415a7e1a2844ef9a617e29608a",
							"darwin-arm64": "04e5dcc0fb8a7c0495db3674eea156b5598eb7e0694867e139fbdb359f9aeada",
							"linux-amd64":  "470413a4775eede6bf9960ba6179b8c8de87867ecae0696f8ab0b6a4413f82e3",
							"linux-arm64":  "b5175186ba5a63860ab2920c48c0de62c3dc6ccd8e88d8057229c17e5acea33e",
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
						ID: "com.palantir.godel-okgo-asset-errcheck:errcheck-asset:1.53.0",
						Checksums: map[string]string{
							"darwin-amd64": "43fc2d104be4fa7b1d71bb96a5cfd94a7a81173e702a9981d615d2b52c721e3e",
							"darwin-arm64": "91204f92fc8a9d434fdfb8a8fc18e8f463053537198c821bb922245d3806c709",
							"linux-amd64":  "9a9ce5061717e539995df4808fa5434682ff8834e0767eb8d39f71d1481dae2f",
							"linux-arm64":  "b107e32e990e6f952ad011f3c461c1161058f73e9c84719d507ae15f8d2ad403",
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
						ID: "com.palantir.godel-okgo-asset-outparamcheck:outparamcheck-asset:1.50.0",
						Checksums: map[string]string{
							"darwin-amd64": "be9ccefa8e2fa1a0027955cfc2550a6f74024f523b02635251d0130325569cbe",
							"darwin-arm64": "f4e2a23d5f20919b4889b66af207224a6351b0f30857796fd9ed73f68261cac2",
							"linux-amd64":  "ca21a723c109332c3b88fb3df7f91399d58aebdba850c20c38ba4998504b9e1c",
							"linux-arm64":  "166bf165b65e1c281f7dac24470d6a11b984e61b75e36e045cd6335434172b75",
						},
					}),
				},
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-unconvert:unconvert-asset:1.51.0",
						Checksums: map[string]string{
							"darwin-amd64": "af36d123414d2727b54f239da05c31bd637315fb844007c6f2000f5299f993a9",
							"darwin-arm64": "bcc18abe74c9a1a877cc53f8c847b2108e298e490831331dfc6ee63513246195",
							"linux-amd64":  "d148829922d96e5788ea661ad6dd22251a68933dda3db7758c5b681c59154b19",
							"linux-arm64":  "af1a50768cbb3c498a02fb5972efd02b72ff4ed19832ea320a2e3c78dbe465f6",
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
