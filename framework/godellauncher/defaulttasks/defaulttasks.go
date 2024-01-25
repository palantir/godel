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
					ID: "com.palantir.distgo:dist-plugin:1.64.0",
					Checksums: map[string]string{
						"darwin-amd64": "289b2c3ecf908459dc921faa4d2882fd8a7ae385739e1bf745c12e55c3e33636",
						"darwin-arm64": "ecc134609be86a63b5b7c45c74ade1ca0b1746516fc237c5c56353400a241bd4",
						"linux-amd64":  "53dc0f1c27acee5416c9d5c122b554dd1d0dafa743b49162ff4d114c637b04cd",
						"linux-arm64":  "613795593bca99c4a465b09aa1ed36864043782a54565928442da1b98716d4fb",
					},
				}),
			}),
		},
		{
			LocatorWithResolverConfig: config.ToLocatorWithResolverConfig(config.LocatorWithResolverConfig{
				Locator: config.ToLocatorConfig(config.LocatorConfig{
					ID: "com.palantir.godel-format-plugin:format-plugin:1.39.0",
					Checksums: map[string]string{
						"darwin-amd64": "52801a11c2f8ded19d1fcbe100cd67c4e64ec6d91252dd0f82a75f61a4282b19",
						"darwin-arm64": "9ba299db7171f7ef05aee6fb5ef208516b217ee70c2abea0f162faaeb96bc509",
						"linux-amd64":  "4f478728396ec9a79c0b4ed64fb6bf8f69b93ddcee5053f7aa40e997af838661",
						"linux-arm64":  "9e2695bfb7b5d011838649e255066800ba11b1b5d2496dfe60099d2338ab9693",
					},
				}),
			}),
			Assets: config.ToLocatorWithResolverConfigs([]config.LocatorWithResolverConfig{
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-format-asset-ptimports:ptimports-asset:1.38.0",
						Checksums: map[string]string{
							"darwin-amd64": "3cb9246ec47975ad68ce3b7c910c8b44563016d1ae7597a011fc8e964d630c9c",
							"darwin-arm64": "ec9c7e8aa7c69d862841ff2e6e52bb9fbfc5436dc45f79de8274f54471d79be3",
							"linux-amd64":  "c0755db015762ed0f96c768b10220eb6b0ab91e4e26991ff216cf4c6cca0ba3d",
							"linux-arm64":  "7ef7773fa13c9e28b72f798457d803c05a4720dc98d44fe42ee1071faac57aee",
						},
					}),
				},
			}),
		},
		{
			LocatorWithResolverConfig: config.ToLocatorWithResolverConfig(config.LocatorWithResolverConfig{
				Locator: config.ToLocatorConfig(config.LocatorConfig{
					ID: "com.palantir.godel-goland-plugin:goland-plugin:1.35.0",
					Checksums: map[string]string{
						"darwin-amd64": "462402be17f9050aa240b6713afa250eafbdc4a9cfee9a6106107fe958dc559d",
						"darwin-arm64": "8735698576484726e844729ba45e12037cdf821eed42987181f07a0215880b98",
						"linux-amd64":  "0cc1a0bd3af7028564aa84fe127d42b7d878acd6007ac1db63e7a00d2212fa84",
						"linux-arm64":  "8f28ceaf4756becf5b4821d19f6ddf13bcb810c5cb39ae7d9138c0a3a15ef1c7",
					},
				}),
			}),
		},
		{
			LocatorWithResolverConfig: config.ToLocatorWithResolverConfig(config.LocatorWithResolverConfig{
				Locator: config.ToLocatorConfig(config.LocatorConfig{
					ID: "com.palantir.okgo:check-plugin:1.46.0",
					Checksums: map[string]string{
						"darwin-amd64": "7d71d986ca0aaf41b6aa48907c86cff50dc7a53724f0825d9ae40fb1fd830061",
						"darwin-arm64": "7be8d9a513f3432160efe07ff0af448af7127dc451af7b62fd1b98a55fe3cb03",
						"linux-amd64":  "f5474584b0b72bfc932c69ced525966cd2f868d065e3e73eff8c7e9153d039c8",
						"linux-arm64":  "661b67317083de6b1a7726f50defe4b5d2dad0a3c02079b7b31648c4a3902720",
					},
				}),
			}),
			Assets: config.ToLocatorWithResolverConfigs([]config.LocatorWithResolverConfig{
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-compiles:compiles-asset:1.43.0",
						Checksums: map[string]string{
							"darwin-amd64": "b1e45402b4a016d3caa53e833b868c4734de9ba27743bb99f3106715dd85bc20",
							"darwin-arm64": "ab5d1a52915f14d547a1ce69a1e3601dff8b64455815dfde85e4374c26e4d277",
							"linux-amd64":  "1b50705954df66412a357d05f76fdd3b336ef0ff751c55373ee4c93b51af6069",
							"linux-arm64":  "9f1c6c2ce00fa791df0cf401726c139648c313a843913f74979eb258bbb8a9cf",
						},
					}),
				},
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-deadcode:deadcode-asset:1.40.0",
						Checksums: map[string]string{
							"darwin-amd64": "83d717828c4f861385673e97c3227fbea7d221d9b24d8cd24bb83c7d28775aa1",
							"darwin-arm64": "9f189354f83f8c3439dce540b1377034f01d1d86f8997a3734ce6d1ddc145f3a",
							"linux-amd64":  "27a560e59d474148794b7666f018f4f5135e926ac6b35bc24a262093d9abc055",
							"linux-arm64":  "b5f843d628701e9ec789719281e5e0a4021dcde5f47d2655e9191c046c5ad435",
						},
					}),
				},
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-errcheck:errcheck-asset:1.41.0",
						Checksums: map[string]string{
							"darwin-amd64": "0e052c45d40bd3e6e649f41070d6b72f3998922acb524986796d339f4a7fc855",
							"darwin-arm64": "e898522ed2b50f9e2abe9332624e88f1a6d08f13de7d1380a3f5f3086f697302",
							"linux-amd64":  "f8cc6b57f7d75b88711477970a010a5bf646a5fdc257cc964b3b6fb66d3dfe20",
							"linux-arm64":  "ff2c78b6cce84fb2f4dd6a031672b907737ddbfaf07b4067ded617deb5b383a1",
						},
					}),
				},
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-golint:golint-asset:1.32.0",
						Checksums: map[string]string{
							"darwin-amd64": "f3567207c40f03ee6cf1339a4731f2afb258a80472b687bd193785f8c82c0601",
							"darwin-arm64": "dac535e2e7faccb4001b4a1c264138a66f79af2a3af7be5121af4f415302b96f",
							"linux-amd64":  "c524532e88a02d289c960afe12df99aaa21eec59976b5021fa281326939c6625",
							"linux-arm64":  "8ac4bd05f120575a8a457c7a570d3f28867395067b6c1ebe6925c182765c6edd",
						},
					}),
				},
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-govet:govet-asset:1.36.0",
						Checksums: map[string]string{
							"darwin-amd64": "640387d7e365c9075b276870eb8be5c4a437481235078b698b32785aca34ae54",
							"darwin-arm64": "5629cb23a5101d503a0be3bd24e8ccf7be35a9193eb86644904b25701fd29003",
							"linux-amd64":  "7202999b45f19beaa1b11d576899e51c1af6971d11f44e9add11caa8a6cc216e",
							"linux-arm64":  "579580b912ced2d820e11088d6e53ecdd7ba336e70d1185c18e7b5a7b9d6a25f",
						},
					}),
				},
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-importalias:importalias-asset:1.35.0",
						Checksums: map[string]string{
							"darwin-amd64": "d6b7101740bfde5ef1078fc32bbd15cde5c46c69ae2cbb10e476330e5cd9be40",
							"darwin-arm64": "6319ebbdb7ade67edd756e9d8585992c0fb3f921a88a2fd7009b0f50d20c85b6",
							"linux-amd64":  "1711a12ce61c7ef144438ae88dd8c6b34affd757f7b736b979ccbe077cae0624",
							"linux-arm64":  "681090bdf05dacd52a205229dcccdaf66706252766de519468341aa07569a360",
						},
					}),
				},
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-ineffassign:ineffassign-asset:1.38.0",
						Checksums: map[string]string{
							"darwin-amd64": "d63798eafff597c16fed78c516db32708df51764208f5e23bc11cafedc60fc4e",
							"darwin-arm64": "0f6bcee95f8bf76254b2bd9bdcae49921d6f6df9ed683fd2fbd794b88c34d37a",
							"linux-amd64":  "dc6150bba1aeab0aa8607188ddc0d400ff773fcb08c62d26e1667edaddec4a8b",
							"linux-arm64":  "7f24e2b0d4a253ee9758875939631f03bf4247a295e0406e05fe925dee7fe2d8",
						},
					}),
				},
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-outparamcheck:outparamcheck-asset:1.40.0",
						Checksums: map[string]string{
							"darwin-amd64": "5db0c4016057b73b7a7b46749d8a871c573f93ee84f4d5c82e749b00706c4062",
							"darwin-arm64": "2ab510e06f0001412fe2173db9bcea00a4507fc6083a1dacbe81de1b6a38de41",
							"linux-amd64":  "713b1a9afc487713f1552a4d5c465925e2eddb6a1e16e275cf1498799d5fb7c6",
							"linux-arm64":  "0b885c855d0c0e27e126345f7c9234786b71e54627b460d3f2872a5e9e423772",
						},
					}),
				},
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-unconvert:unconvert-asset:1.40.0",
						Checksums: map[string]string{
							"darwin-amd64": "aeec2fdc406c9253203e744dbb21120a60c414efbd4f8de26f5d77751b502d41",
							"darwin-arm64": "028b6affa80132a0f10c6a7c342e3dcaf8d558cbba964890f01c0d1bf70ffe3e",
							"linux-amd64":  "325668a9e0cbcccd08bfaf08c71ec3a9ce38010568985861b08c838202cb5096",
							"linux-arm64":  "7cc116c07a3c0b7e3697f3d269cb5711ab4fe99b9dd7657c33b704f1584258c6",
						},
					}),
				},
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-varcheck:varcheck-asset:1.40.0",
						Checksums: map[string]string{
							"darwin-amd64": "ec5f08c6acb63b5d4a2f8446f9de018a07f601a6573098931d96857767670805",
							"darwin-arm64": "0000790e319b1f2e53af19f5f940b2943caf50ffe730d62c2b391cedf15086a2",
							"linux-amd64":  "5f889b5bf0d0c8555aeca0d48a641ee940e20b77b37ee4399ea135a619bc7326",
							"linux-arm64":  "554f30e46a0426c2d499773aa060b417d485e68163b5022e788b713fd2ab8dce",
						},
					}),
				},
			}),
		},
		{
			LocatorWithResolverConfig: config.ToLocatorWithResolverConfig(config.LocatorWithResolverConfig{
				Locator: config.ToLocatorConfig(config.LocatorConfig{
					ID: "com.palantir.godel-license-plugin:license-plugin:1.38.0",
					Checksums: map[string]string{
						"darwin-amd64": "7f15c466731e9bce6c0b37fc18258f04fdd53482530ea2c0a339aad9d69ec98d",
						"darwin-arm64": "caf1c24d9e36950dffddf35ae424a6316e49528520b6666df7f62bdecc1e363d",
						"linux-amd64":  "45417098a6c0600256b0e509560c74c82b27926640f7702669e55139e5edbd65",
						"linux-arm64":  "68c1fa9b7762b3bee9a217a3330c97228a0336a88dbb378e9716b92f362c0daa",
					},
				}),
			}),
		},
		{
			LocatorWithResolverConfig: config.ToLocatorWithResolverConfig(config.LocatorWithResolverConfig{
				Locator: config.ToLocatorConfig(config.LocatorConfig{
					ID: "com.palantir.godel-test-plugin:test-plugin:1.37.0",
					Checksums: map[string]string{
						"darwin-amd64": "0284fa8d905dee67a3d574941b861b9652ac9b54bed96c1fa71d776b86d89899",
						"darwin-arm64": "60592772ed461c694d9835a96d3a7f1cfc3ede55284ef377c47e6fc66c40c3fe",
						"linux-amd64":  "8ef1033a373a70fe762be547a96ee8d1dec6900b52a4d9ad8323b95bbffa3369",
						"linux-arm64":  "2cd9673532c4ff17cca7e2c5d6d3dc80c3253ba4e315b3072ba0eda77af86583",
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
