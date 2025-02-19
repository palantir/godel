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
					ID: "com.palantir.distgo:dist-plugin:1.80.0",
					Checksums: map[string]string{
						"darwin-amd64": "d30be4b41aa55f85c5b8f369a7b284623bc477ed3e0df828b5cd46a70c692c58",
						"darwin-arm64": "7d1dc07a909274d55bbdbf5faa354f42ee7130ba33037ea5e996d67ad5edb8ee",
						"linux-amd64":  "7c60a0891e7c93c7bc50a8e07324ad783563fcc706100856f8cffc2e650f9c3a",
						"linux-arm64":  "7046744ae01f8d653b300cba3412a411f81489757527492cd803e9ba8b402765",
					},
				}),
			}),
		},
		{
			LocatorWithResolverConfig: config.ToLocatorWithResolverConfig(config.LocatorWithResolverConfig{
				Locator: config.ToLocatorConfig(config.LocatorConfig{
					ID: "com.palantir.godel-format-plugin:format-plugin:1.50.0",
					Checksums: map[string]string{
						"darwin-amd64": "2b96ec2f6572bb751f70452bfb37e28e6364243bfe8e0e5fbf013792e3e72d38",
						"darwin-arm64": "ab01ec0793744fcf5a5deb30778a0b7005de4bc1d3bb0a7289b5ced048f17262",
						"linux-amd64":  "2368756b2d477aba1bb637b92e17eab433659a7c7ef8ebe8217bd1b66755cebc",
						"linux-arm64":  "51afec2fcf90a551b1077c9f74e060f269405a2409931cf8de9e94e683e7fc2b",
					},
				}),
			}),
			Assets: config.ToLocatorWithResolverConfigs([]config.LocatorWithResolverConfig{
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-format-asset-ptimports:ptimports-asset:1.49.0",
						Checksums: map[string]string{
							"darwin-amd64": "7a7a3121e3388c4f5183442d6b9ba3d34333a15a60ebcf254047d6b78a04b5df",
							"darwin-arm64": "87bb103a8d009fbe7d5ff16269a767826280953d9e732282d7012b2648393dde",
							"linux-amd64":  "eabe4c5de843b3fb1aa5d64d6b54f8b8dafb0bb1d31d54f5aee9b94e76d08d06",
							"linux-arm64":  "2374491ac3069284130b22e14f81e9f68e941a1976222d76c5bfa5bc6f0b5547",
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
					ID: "com.palantir.okgo:check-plugin:1.61.0",
					Checksums: map[string]string{
						"darwin-amd64": "780ec7bdc240d0176895011fb32aa64a964c7c92d47e315dc0a0f17930d85b51",
						"darwin-arm64": "c9314361dc09684bbb5d74e174821256e7f3524a840ee7e60e3b673c3b449d49",
						"linux-amd64":  "91a8e4f00dadbf8016b38f93875735e3d0765594db7c4f640b4ed82b3ae867dd",
						"linux-arm64":  "62697b7292de224ceb596c1044295fb2e32de34d74b1ca87f10806d1596a3c6b",
					},
				}),
			}),
			Assets: config.ToLocatorWithResolverConfigs([]config.LocatorWithResolverConfig{
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-compiles:compiles-asset:1.55.0",
						Checksums: map[string]string{
							"darwin-amd64": "a52d48a4618e9e615c5223bd40384dcc46fcfeae5e2e6277831033f6bec501be",
							"darwin-arm64": "a4daf46e0afd926841583ac1a47b6bb4131ce6f77ce9b619b2e67f88749bbde7",
							"linux-amd64":  "20601839070092a55f7a735c58138c3da7609bd3954e0fc3c9794ac3aced5d5d",
							"linux-arm64":  "0831267fad42d9f90445f6f7079e6e181b9156b0e1b5c9fc5b484da961be275b",
						},
					}),
				},
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-deadcode:deadcode-asset:1.48.0",
						Checksums: map[string]string{
							"darwin-amd64": "117bd696fdf8bf210c06a090cdf232efee2c0c268507f620619dbe7d2c0b5d46",
							"darwin-arm64": "6c86d4ec40effd94397ac221f4c1eb34358dffde020076a8130b6b4d590e1008",
							"linux-amd64":  "1cdcd8fe6a6073e8b5c2a6c0fa93ae4410f74cbb94de598b2f9a33ab70913a75",
							"linux-arm64":  "0978e13be0cee8290bf1ce8a3225d7aa67859b9c47267925f1b4e7e198155a0a",
						},
					}),
				},
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-errcheck:errcheck-asset:1.50.0",
						Checksums: map[string]string{
							"darwin-amd64": "2caf963dee72914cd5ad847b88f52c54ac720262a46e61e5e9deb2e3fda930a7",
							"darwin-arm64": "1702fd0474439f3e3a515f5cf8b49e8744356ff93a5129d8c3a2a54cd473a76a",
							"linux-amd64":  "2b05e827b48ff688cbf4ea5fc272cef632889f6b7e5ebb83e58b682381881917",
							"linux-arm64":  "2a1bbbfaf1c3f5a2e147eb36861a7ea69f94574a28ef6ef971f9fbe69376471b",
						},
					}),
				},
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-golint:golint-asset:1.40.0",
						Checksums: map[string]string{
							"darwin-amd64": "0613f98b068a319ccfbc83b6d2e4b64bac2e584dd701cf43191cda17a5330f39",
							"darwin-arm64": "8a9e5a9598df3fec4b8e6e7e5c15c9e712951705371ca85fd3d56a0aad258c1d",
							"linux-amd64":  "875226e2a9c906f8a493ac3aa526f3acd2c743f49a3bfbdcc47e663717b7a951",
							"linux-arm64":  "a65cd31b39b1168946d5ec466c16862e3255b647c686f2e780053324d2d9423f",
						},
					}),
				},
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-govet:govet-asset:1.45.0",
						Checksums: map[string]string{
							"darwin-amd64": "f7c1cd91c7ae1b68855cbd4b4a32c00dc6de5637518ee2964a4582afab6c1ce9",
							"darwin-arm64": "853e082633385ee081534d287032d91f51e41efce84cfd1dcb50dd2ffc0839ce",
							"linux-amd64":  "e88e1f596588dff02399d508450966e7f428f2731fe960d8937fd6d02fd5fb60",
							"linux-arm64":  "a41d77386677d206a341cd3fe7c08f6aa2399fbe02a2def4254f00b6e099844b",
						},
					}),
				},
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-importalias:importalias-asset:1.44.0",
						Checksums: map[string]string{
							"darwin-amd64": "e43472367103fc6df1af98630d1f11fb1866154d463f925f82d5f222cb3b4bd8",
							"darwin-arm64": "1d6641ebc38a0b25ad31ff2b8689fb64571fd97d6c25a39908ef3c84754053bb",
							"linux-amd64":  "5edbe4a0b640af14db717c68228499b60df8af0be3c04b99fb619a1f393ddae4",
							"linux-arm64":  "20db1615b9fd0f57a5ff0d368ba9842eb69ecd7f1aa201c4985d536ac81394c4",
						},
					}),
				},
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-ineffassign:ineffassign-asset:1.45.0",
						Checksums: map[string]string{
							"darwin-amd64": "8552ce61d63357dbea1c74422153a8fea076608740a0a3d549172b939d4f6bef",
							"darwin-arm64": "2453122947c66d3b3b7a70c4ceafbb0b888f3e1cc864fa9a789288d12e41f06c",
							"linux-amd64":  "24b9d11cec00c32a262a08ced0129c130915ea3ead74f100e2e41f53dfe3f74a",
							"linux-arm64":  "66109c6b18efb567a53cc1c81e9e1a73c2f6e98105939f08c8abe275eda297d2",
						},
					}),
				},
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-outparamcheck:outparamcheck-asset:1.47.0",
						Checksums: map[string]string{
							"darwin-amd64": "61649761bdcd59dc2ea79405f3a5017d9bf6ee3ca54ae86d2ebc7dbea69a3517",
							"darwin-arm64": "257b7651f76670e635255a7368bd3193a5db45365338dc2785cb57ed6c11f751",
							"linux-amd64":  "5ba5beed5e36724a98632aed5ceb35499c94faddde7f645000fa95888e4bfcd5",
							"linux-arm64":  "080a9647724f6f896be5f479663c78180a1d1df0c3a94a984f1fd5473726cca6",
						},
					}),
				},
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-unconvert:unconvert-asset:1.48.0",
						Checksums: map[string]string{
							"darwin-amd64": "514ca8f2b815b4f684e54a752591d77bc9786a3dd7f41a5654c84f6d79b18e20",
							"darwin-arm64": "c55a32581a311a44b58270407abc2452a2f4a8eaf33857133ccd888fd2d3951d",
							"linux-amd64":  "9cc14c36aee19dc7de877b0e3c163f101e4f6727f6d6c5c16e531301c95f4d4a",
							"linux-arm64":  "d9623b3911013d23f82b784ee0432973d0863736faec5a9e4d00a306dd4eceda",
						},
					}),
				},
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-varcheck:varcheck-asset:1.48.0",
						Checksums: map[string]string{
							"darwin-amd64": "a9370fce64016639e146ae26cfef1563095349d1962bd148c9ec012f0731b02d",
							"darwin-arm64": "9ed6e5c4419a84dde638a2f9bf383c62eff2c4bf248ef501273f9321614541d5",
							"linux-amd64":  "7279f86c9db0ad9ba01002829e7b39704877a08f5fbb736650e14a41a6070687",
							"linux-arm64":  "a56c371f13e3c327e882b99bf6ab78ef5617b0dab3302e9913df1344525fc30f",
						},
					}),
				},
			}),
		},
		{
			LocatorWithResolverConfig: config.ToLocatorWithResolverConfig(config.LocatorWithResolverConfig{
				Locator: config.ToLocatorConfig(config.LocatorConfig{
					ID: "com.palantir.godel-license-plugin:license-plugin:1.47.0",
					Checksums: map[string]string{
						"darwin-amd64": "788dcdead84c7e46b8a29ae2ad490b6e7d13e29cb0e9bb34c3d67836bf365379",
						"darwin-arm64": "fb85b3cea5392e8f6d2c1b25357773a80ae927bc619b5e303cff8bf9478a2066",
						"linux-amd64":  "090f206f78eeef072a55e85649125562ef875d155891a79a81c4911d9b2b6dd9",
						"linux-arm64":  "a7d9217a79c94e0caa28852f00d1f2037dec6bd5aa92c85997d73f774156a6c3",
					},
				}),
			}),
		},
		{
			LocatorWithResolverConfig: config.ToLocatorWithResolverConfig(config.LocatorWithResolverConfig{
				Locator: config.ToLocatorConfig(config.LocatorConfig{
					ID: "com.palantir.godel-test-plugin:test-plugin:1.45.0",
					Checksums: map[string]string{
						"darwin-amd64": "b8a81a75306440a33a2e8b45b3b4172bd67a701601227da5a53abe2d8bfdf276",
						"darwin-arm64": "8f06c34f7d1b3c09a090243c9131815e3b1a31dfd0339076a59c95062dc4ded2",
						"linux-amd64":  "18a973aab52d1374caccf92b36a5593b7b8456f87da5084e47edf56afbf198c9",
						"linux-arm64":  "28d9ba5e59e1d3935ddd300e2524946793164b26ad5ec4c17ec7f3f7f04ae77a",
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
