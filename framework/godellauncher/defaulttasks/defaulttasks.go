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

	"github.com/pkg/errors"

	"github.com/palantir/godel/v2/framework/godel/config"
	"github.com/palantir/godel/v2/framework/internal/pluginsinternal"
)

const defaultResolver = "https://palantir.bintray.com/releases/{{GroupPath}}/{{Product}}/{{Version}}/{{Product}}-{{Version}}-{{OS}}-{{Arch}}.tgz"

var defaultPluginsConfig = config.PluginsConfig{
	DefaultResolvers: []string{
		defaultResolver,
	},
	Plugins: config.ToSinglePluginConfigs([]config.SinglePluginConfig{
		{
			LocatorWithResolverConfig: config.ToLocatorWithResolverConfig(config.LocatorWithResolverConfig{
				Locator: config.ToLocatorConfig(config.LocatorConfig{
					ID: "com.palantir.distgo:dist-plugin:1.19.0",
					Checksums: map[string]string{
						"darwin-amd64": "85b160a5c765e984f0d8d21edbf2ac77b5162fe26ec618682271c65bdf522c24",
						"linux-amd64":  "544e01f473811891c3e3478d9f2045128c91839498a61387b107e4c56dbd5007",
					},
				}),
			}),
		},
		{
			LocatorWithResolverConfig: config.ToLocatorWithResolverConfig(config.LocatorWithResolverConfig{
				Locator: config.ToLocatorConfig(config.LocatorConfig{
					ID: "com.palantir.godel-format-plugin:format-plugin:1.3.0",
					Checksums: map[string]string{
						"darwin-amd64": "423079b4e5768ed6c396f6aa9fdc7992fe58dcb60cdc24306357a1a4ba0e2535",
						"linux-amd64":  "a44c905aa4b9e4e196bd784c33c2364713954685586a5a1054c7de568e547850",
					},
				}),
			}),
			Assets: config.ToLocatorWithResolverConfigs([]config.LocatorWithResolverConfig{
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-format-asset-ptimports:ptimports-asset:1.4.0",
						Checksums: map[string]string{
							"darwin-amd64": "f904088d8bb33ced244c2d36f98ef9d3d082439eb78dd07e009e41cfc77b39d8",
							"linux-amd64":  "14b1b01f593987e9f20c4ac0f5b033f1e888be874a896259305441424a2eca33",
						},
					}),
				},
			}),
		},
		{
			LocatorWithResolverConfig: config.ToLocatorWithResolverConfig(config.LocatorWithResolverConfig{
				Locator: config.ToLocatorConfig(config.LocatorConfig{
					ID: "com.palantir.godel-goland-plugin:goland-plugin:1.1.0",
					Checksums: map[string]string{
						"darwin-amd64": "e0b8ec0629bad270493501c0ffd92bf72f5ac592028c8150d6bc12a6716857f1",
						"linux-amd64":  "ad322d2dfef926edd03f28beef9b527726d11c68ab794a1f400164abce303084",
					},
				}),
			}),
		},
		{
			LocatorWithResolverConfig: config.ToLocatorWithResolverConfig(config.LocatorWithResolverConfig{
				Locator: config.ToLocatorConfig(config.LocatorConfig{
					ID: "com.palantir.okgo:check-plugin:1.7.0",
					Checksums: map[string]string{
						"darwin-amd64": "b90823ad40b29412540d370de2a83e845904951730817020e3e44619ba76af5b",
						"linux-amd64":  "0954b12695a4ba627195e651f37cb3c1230dc11ccbe8ed6876168d2d04d545fc",
					},
				}),
			}),
			Assets: config.ToLocatorWithResolverConfigs([]config.LocatorWithResolverConfig{
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-compiles:compiles-asset:1.4.0",
						Checksums: map[string]string{
							"darwin-amd64": "185967b490d92b8bc9e2cf037470a5874bef19567741cef0caa56b0cfa2fb26b",
							"linux-amd64":  "58f234b614442a1851e59818e64c31240332d10fdfc40def60962cf389590070",
						},
					}),
				},
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-deadcode:deadcode-asset:1.3.0",
						Checksums: map[string]string{
							"darwin-amd64": "73784237d1a0a6bc4eb2130a20adb27312cac7b507f9ceb1b1ecdb7dad9664be",
							"linux-amd64":  "9bbb6bd3c8a2be5dd73ae0d952025c9865f2527d805854fe35fc0d4fcb1d7c40",
						},
					}),
				},
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-errcheck:errcheck-asset:1.5.0",
						Checksums: map[string]string{
							"darwin-amd64": "b3d90c36b184178d3ed3e071e62541aa342ea503b6fb66b52db036a32afc97cb",
							"linux-amd64":  "1bc2107777bc99d4e83072c925a690b7bde32a01b2be3bfa3cdefe50d5ee295e",
						},
					}),
				},
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-golint:golint-asset:1.1.0",
						Checksums: map[string]string{
							"darwin-amd64": "2324fd71a531d898a196d44f9fd6838c1b09a7e19acd8ed96c27b42d7d9cc6ce",
							"linux-amd64":  "f3a10dfdf78fd184babdd84c862dd37845f564d427b8907b6cf65aa47b7ac2ff",
						},
					}),
				},
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-govet:govet-asset:1.1.0",
						Checksums: map[string]string{
							"darwin-amd64": "3f4e93b97eac531b44e7c57016803be45663df5212b4f70184b116a96e3c19f9",
							"linux-amd64":  "11775d26d8c71b883ae5e584279f3016978bd37ced2708bf4368f8fac5c169f1",
						},
					}),
				},
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-importalias:importalias-asset:1.1.0",
						Checksums: map[string]string{
							"darwin-amd64": "228bcc7060b7a6996e9d5c8a61a57553092332c907b4e13b0c97eb05c0db1ae9",
							"linux-amd64":  "437f6315905094b5bf78d473a849a5266331882a9a4f74732555d4f36a8332d6",
						},
					}),
				},
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-ineffassign:ineffassign-asset:1.1.0",
						Checksums: map[string]string{
							"darwin-amd64": "de143193ad6265e372c4f84645cc1d1f3dae5c23ccd4b9d7b9f89ecef6f3c49a",
							"linux-amd64":  "0298e35f311c3d0069044e446830b5636aab93fe38317f03a113e4edae67ce37",
						},
					}),
				},
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-outparamcheck:outparamcheck-asset:1.5.0",
						Checksums: map[string]string{
							"darwin-amd64": "bc87dd25cdebbeb0b7fa335d143c67a3a598e469002c2d59efaf88d8a5d422a8",
							"linux-amd64":  "0d1d849bac975ccc6ece6e1e8723cb748a46681300227b7725cee954164daab9",
						},
					}),
				},
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-unconvert:unconvert-asset:1.4.0",
						Checksums: map[string]string{
							"darwin-amd64": "7c4747af08edbe79737be38cebfc93e47132094add38ea71346f370adcf28431",
							"linux-amd64":  "dd9d0b1e42424d608581a1c935c8c29021e5591f21064c8aa48707997a85b7de",
						},
					}),
				},
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-varcheck:varcheck-asset:1.3.0",
						Checksums: map[string]string{
							"darwin-amd64": "1c920ee44db304dcb936677648246d22f4a547515eccc15952f1fc83c3df1a3c",
							"linux-amd64":  "2ff5dcbebd1ba2402ce691c2e5cd40865d66b4f6d42d35b1c399d74f9bcd354c",
						},
					}),
				},
			}),
		},
		{
			LocatorWithResolverConfig: config.ToLocatorWithResolverConfig(config.LocatorWithResolverConfig{
				Locator: config.ToLocatorConfig(config.LocatorConfig{
					ID: "com.palantir.godel-license-plugin:license-plugin:1.2.0",
					Checksums: map[string]string{
						"darwin-amd64": "339da7a86b948c06796784e0e362c83884ac80206caa5777ddb8aca4d0199235",
						"linux-amd64":  "68ec11b272a37172cc54eb0cf21a28a9cd3ad6d2afbb4cc04e2c231d10f34530",
					},
				}),
			}),
		},
		{
			LocatorWithResolverConfig: config.ToLocatorWithResolverConfig(config.LocatorWithResolverConfig{
				Locator: config.ToLocatorConfig(config.LocatorConfig{
					ID: "com.palantir.godel-test-plugin:test-plugin:1.3.0",
					Checksums: map[string]string{
						"darwin-amd64": "25474251b4253be163bdbfa79652d850395612508019c459bf6eb3d8aad1544c",
						"linux-amd64":  "0bba9663738714e68a684a91eac4ee67b7a0e0966b9712f153c7238deefce469",
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
	pluginsCfg.DefaultResolvers = pluginsinternal.Uniquify(append(pluginsCfg.DefaultResolvers, cfg.DefaultResolvers...))

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
