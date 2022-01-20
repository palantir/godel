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
					ID: "com.palantir.distgo:dist-plugin:1.32.0",
					Checksums: map[string]string{
						"darwin-amd64": "d6116dd87cd7ba13846e59285546eaac674490c78bc7459566ec8f03bf79eec2",
						"darwin-arm64": "a02c361dc23ed8b4681c55574c42b57bf58df11d00aa61324532329f9d610de4",
						"linux-amd64":  "a37e792ad847c9abe2ef6362854c606b6a81ecc0e93f696cfbc23f6d98e28e51",
						"linux-arm64":  "e7cdd4dea613904c5876067998e32e96f6a23ce7108ddce687a0287e0ab27490",
					},
				}),
			}),
		},
		{
			LocatorWithResolverConfig: config.ToLocatorWithResolverConfig(config.LocatorWithResolverConfig{
				Locator: config.ToLocatorConfig(config.LocatorConfig{
					ID: "com.palantir.godel-format-plugin:format-plugin:1.10.0",
					Checksums: map[string]string{
						"darwin-amd64": "58287344b300fcfb64d651eaddff4dd4a9472e171acc50e108dc3c8688b5f2f8",
						"darwin-arm64": "804c6f529f6ef279ca852eaa62f057ceba22a09b37df22d3cbef38e93ccb1752",
						"linux-amd64":  "9875dc2d051223e689f0f574637bdc7d49c68d11aabf994936a04feb64794a58",
						"linux-arm64":  "cc0e0bc0acfa39caa70c1ca2dfc774f3dd8959f47457b553544544f60a101395",
					},
				}),
			}),
			Assets: config.ToLocatorWithResolverConfigs([]config.LocatorWithResolverConfig{
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-format-asset-ptimports:ptimports-asset:1.10.0",
						Checksums: map[string]string{
							"darwin-amd64": "acfac7e6e29e18e793ebc63d4256be050fe2cb4727332a69c78e01b598ec28f6",
							"darwin-arm64": "069837fe01d98d7f0da3a6f35199c5ef5b93bf9c2be54f6ddc63264f4fda8c22",
							"linux-amd64":  "2dd0a9077f23a110fd3ad51bb688d82758d0d114080de2169da82de06cbed810",
							"linux-arm64":  "975b4d54b30327e9ded4cdbdb891caf60804cc02ccc75d7c2c0bb28244a4cf71",
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
					ID: "com.palantir.okgo:check-plugin:1.15.0",
					Checksums: map[string]string{
						"darwin-amd64": "36a7a88bbbd99be1ad482f0e2a009e316012f9ecb670c2ca1c8cfe32c614b3bf",
						"darwin-arm64": "c308ba878dc58b5b511ee72c262ddc10517b330a248d4e6e3f8ebef86280fb8c",
						"linux-amd64":  "71a06ad46244f1fc69189538e1958b0b05440103c7e07600edede98b86de2be7",
						"linux-arm64":  "93a80e0a73e15afe6e5373e3415b67a10ac53ea070eaaccf602eed0c1aecb5aa",
					},
				}),
			}),
			Assets: config.ToLocatorWithResolverConfigs([]config.LocatorWithResolverConfig{
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-compiles:compiles-asset:1.10.0",
						Checksums: map[string]string{
							"darwin-amd64": "b5e9ba2afbf03201379bb09e0582d3ee945552514bf5ecbcc27b4225967a17f3",
							"darwin-arm64": "0ccd6d7d02bcc7327ced361039dd262de4d375a21377a279ebd6e3fb0b5342d0",
							"linux-amd64":  "9fd1a534b1be468bae81cefa63ce9f63c33771764ec1f08cd0466e23dab80a75",
							"linux-arm64":  "44f728bc63aef501eefb58335a273cc8d8e604654648b2cf594e5b57a94d83ef",
						},
					}),
				},
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-deadcode:deadcode-asset:1.9.0",
						Checksums: map[string]string{
							"darwin-amd64": "9c3e42d1f172cd1e52da2b127aab39d07444bf55d5ae5e149a4ad69f354eaac7",
							"darwin-arm64": "fc9bffe5f10a6829ebac598e37480746286bc63f843cc644b932918b5dcfd4b5",
							"linux-amd64":  "9cc164921ad65b727904647fbbcf2a2f47ce7a9eb652975e48c78717dffc3b94",
							"linux-arm64":  "aee698e396fa38c7127f4e69799698e8ed9931525e9ab0d3202c6a16dcc5072f",
						},
					}),
				},
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-errcheck:errcheck-asset:1.11.0",
						Checksums: map[string]string{
							"darwin-amd64": "6cbfdd09f9740902cab3ce70b7e48c498d4b827aaa09dce720e4c20fd9a87857",
							"darwin-arm64": "0d9aa82bcb91ec3f7b233e5f41c201bf8c73c0e3631d4b10343476a18e15e2d4",
							"linux-amd64":  "95db33766cb76d4b3545c5ace08d075b8fb757337b1925a7da0df0a7cfeb4c5a",
							"linux-arm64":  "482ede398a86d97bbfa503691072d68ff49714c547369d505908acb6b2e40b93",
						},
					}),
				},
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-golint:golint-asset:1.7.0",
						Checksums: map[string]string{
							"darwin-amd64": "499c304c426b759296270de17e3d8c28fd61e135b5f475357ae74033bb1319bb",
							"darwin-arm64": "d68d6195f41439c54c34b7346c3d4a3efac049efa01112cd12543af4fc82493e",
							"linux-amd64":  "b17108ee2d4106156e10cf46491b6f3ce9b999dbc08618ea610c75c91bdec3c6",
							"linux-arm64":  "9a857697ad745d0fa1623d195ac19ba60424283228f4534ea685d2089eebcb52",
						},
					}),
				},
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-govet:govet-asset:1.7.0",
						Checksums: map[string]string{
							"darwin-amd64": "d72cf3f8f8b0170b5f7adee8859e184f131daa7ba629000ebbf52c7c5a6af4e8",
							"darwin-arm64": "6585a26caf5c9cf2206587569832f18dd9343cd9712519ba1aa6112dba1263d3",
							"linux-amd64":  "9de684c62afeec85cc0b70d66c5b7282c96dc1130bccc8d1d8e4185a05d97506",
							"linux-arm64":  "b6a6c33cec2d01583204f790f1d53478ca437c2b3030ed7bcc2469aa4b09ba24",
						},
					}),
				},
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-importalias:importalias-asset:1.7.0",
						Checksums: map[string]string{
							"darwin-amd64": "589fee6854771e55e14e800933f9a68c66ed9af37fc64aa16aa11f4b6a307cca",
							"darwin-arm64": "6272badff04bf92c46972a3e4171138d4db81fc622a17bbd7ef4d903671813e9",
							"linux-amd64":  "30a8aed056d3ec86562ea22969e2030bec27ae96474f31edca5e2fee0d239c15",
							"linux-arm64":  "4ff9ef6c8219880b66b9d9e790bfb2717a158728caa443d611b86ab168c9b2b2",
						},
					}),
				},
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-ineffassign:ineffassign-asset:1.7.0",
						Checksums: map[string]string{
							"darwin-amd64": "9afd388e2552d2e58b76f8df9646c643feb49d87a1b116a7eb60e2829bc12344",
							"darwin-arm64": "cfc7ed9928a3ece7fbb7926054dfc9e2a8d90e1e60a1f838cd6599405ce188d8",
							"linux-amd64":  "77a6195313323de91d962a473424a88eec8279421ca799d14d3daf60cca6b394",
							"linux-arm64":  "0f2cc4160b2e637e885bb5cc6ce59e8869bbc802e7fa4fa44f2b0dc859743e7b",
						},
					}),
				},
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-outparamcheck:outparamcheck-asset:1.11.0",
						Checksums: map[string]string{
							"darwin-amd64": "668aa8f77b10113bd4bb6d89ce2f341f11e674e79e7e121caee57cb0005763da",
							"darwin-arm64": "2e076acd7d5043117a0328fd261eda605760f57068367cfd01ebdda74904ebab",
							"linux-amd64":  "00d99b1b167420a48a41b464a3a4e484cdfc253ed7c2290368987154acbd1707",
							"linux-arm64":  "821747ede86bca3e660b0da5c576007ef44f3ec0d47825bc6d6b7201481beeab",
						},
					}),
				},
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-unconvert:unconvert-asset:1.10.0",
						Checksums: map[string]string{
							"darwin-amd64": "f1fb10cf474c811bf30f7d014b89938835c8f27349311a995f2ce2bea241e89c",
							"darwin-arm64": "15d4ce0792be7f2495f8e1dd1fe0132110f07d5c9aa531873cc4276a0440b625",
							"linux-amd64":  "de39587afd51b71cf68749fae23ea187b447cc43a273e31a5d3d384781b852f9",
							"linux-arm64":  "0f08be0bee51f5775f2830b7a43f9a75f4ffc4600dbcafa241538b991ba0dbed",
						},
					}),
				},
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-varcheck:varcheck-asset:1.9.0",
						Checksums: map[string]string{
							"darwin-amd64": "890fdfb2a128da4d5e968dbeba7115ea6bde5427c03c9e2bff5c6d779cee4de7",
							"darwin-arm64": "121aa5b16fa31f08fd0c72bfc8f635d36005baa0831ab832776ecb98e8bd0710",
							"linux-amd64":  "357dfc62fb862612c365691973358fa0d4f8ee8673fae0fbc95913d0b4723d7a",
							"linux-arm64":  "a98a0aab24f87b5fa65bf73e73c01f63fa0a2fd627c2a31f93830613f2904595",
						},
					}),
				},
			}),
		},
		{
			LocatorWithResolverConfig: config.ToLocatorWithResolverConfig(config.LocatorWithResolverConfig{
				Locator: config.ToLocatorConfig(config.LocatorConfig{
					ID: "com.palantir.godel-license-plugin:license-plugin:1.8.0",
					Checksums: map[string]string{
						"darwin-amd64": "640bb33301c0ee3e7e776d92a8f9714c50a1e870d45c87aa6cc64357aa3431f8",
						"darwin-arm64": "489459976552c913b1ad69282cbbf56ad1b8d4c3b34387b7831fc9b679780a50",
						"linux-amd64":  "73363bb86d6daf6a0dd295d3dd5331a341b38be11166ba4bbd475cf429ca9c03",
						"linux-arm64":  "e6cea1ad5c53873ed0e2633d97ba48e7f33ae99ef827b3fa83b1314f87016ac6",
					},
				}),
			}),
		},
		{
			LocatorWithResolverConfig: config.ToLocatorWithResolverConfig(config.LocatorWithResolverConfig{
				Locator: config.ToLocatorConfig(config.LocatorConfig{
					ID: "com.palantir.godel-test-plugin:test-plugin:1.10.0",
					Checksums: map[string]string{
						"darwin-amd64": "7cbad403f3be70dd57c2d10c0915fec978be785f1faf97748a8d0096f4a81e69",
						"darwin-arm64": "c2ff5f64e7a7e32c14bdcc940604a6f7ef21811d8a40e170f9316dfdda9912d0",
						"linux-amd64":  "4648596ba6bc25e71f53fe1005049bef6a23ceba2731eb6537c9e7ad19c48b16",
						"linux-arm64":  "8fbf12cf640d117396002ae33f3fa75a11c1afb94b846de6e32532b126076df2",
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
