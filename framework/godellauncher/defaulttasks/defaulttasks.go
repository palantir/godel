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

const defaultResolver = "https://palantir.bintray.com/releases/{{GroupPath}}/{{Product}}/{{Version}}/{{Product}}-{{Version}}-{{OS}}-{{Arch}}.tgz"

var defaultPluginsConfig = config.PluginsConfig{
	DefaultResolvers: []string{
		defaultResolver,
	},
	Plugins: config.ToSinglePluginConfigs([]config.SinglePluginConfig{
		{
			LocatorWithResolverConfig: config.ToLocatorWithResolverConfig(config.LocatorWithResolverConfig{
				Locator: config.ToLocatorConfig(config.LocatorConfig{
					ID: "com.palantir.distgo:dist-plugin:1.20.3",
					Checksums: map[string]string{
						"darwin-amd64": "5740f3e75fd79ab423fffb0664fa2751060ddd970722cf40570738f1392406bf",
						"linux-amd64":  "06a3e38d1c92baf03cce637b70f4e47b47c7715372cea1dda8cb15a8f336f09b",
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
						ID: "com.palantir.godel-okgo-asset-compiles:compiles-asset:1.5.1",
						Checksums: map[string]string{
							"darwin-amd64": "474c251eac93810c78e7344d4039bc5756b208896fd06ca96b22b80eda9df236",
							"linux-amd64":  "12f1a7bf62f308d638d54a1e6e27c9506f25597dcd4cec1b7ad45333aba35340",
						},
					}),
				},
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-deadcode:deadcode-asset:1.4.1",
						Checksums: map[string]string{
							"darwin-amd64": "54ad4e88c01333702840f64b9d895382cb20931b07bca999a0e090527b019210",
							"linux-amd64":  "d49602063027ff6e7b49374df9ccaf77cfa3be2ae083eff19b178f923c090b3d",
						},
					}),
				},
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-errcheck:errcheck-asset:1.6.1",
						Checksums: map[string]string{
							"darwin-amd64": "8ba8b6f188b726d2d544ac06365ff07c4cdee5b67df94be05141e97eff7b58a7",
							"linux-amd64":  "19fd15d8032fcc653de350d560420ab7e7e56a23e494ff17ebd55092a7d5fb40",
						},
					}),
				},
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-golint:golint-asset:1.2.1",
						Checksums: map[string]string{
							"darwin-amd64": "d5292885ea4df5c4c7ea24f5642c33ff784efed87a8736cac8b4ac10ce3778d7",
							"linux-amd64":  "c9e4fd3885962b496cb167207c6b98d532ec71482324991c06578acc6aaa565d",
						},
					}),
				},
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-govet:govet-asset:1.2.1",
						Checksums: map[string]string{
							"darwin-amd64": "102171beba94b43b6540df36124c46dd8c1fa9b4de6989a10f5c827de7c9c599",
							"linux-amd64":  "509999f5a34aa72651e7d66769809c862069a34e55d5fa7d3aaa5fa924622abb",
						},
					}),
				},
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-importalias:importalias-asset:1.2.1",
						Checksums: map[string]string{
							"darwin-amd64": "f680f2b0bcb6a4dafc923fe4235c5208f889f383cb6122d19012315c914f47b6",
							"linux-amd64":  "e744c36a76eecb30e7e11bf2167dbf2ba64c80878ebe98657e9c1faf9cbb533f",
						},
					}),
				},
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-ineffassign:ineffassign-asset:1.2.1",
						Checksums: map[string]string{
							"darwin-amd64": "acc483e9c627ca2f94b1b2d07922bae7ad8661dcdfaaa8175b46a02430a01a31",
							"linux-amd64":  "7e24024dbcc29227acfd85b9d4e3bf8ae868d62c64dbe81e6e2c94dc1cd1bafa",
						},
					}),
				},
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-outparamcheck:outparamcheck-asset:1.6.1",
						Checksums: map[string]string{
							"darwin-amd64": "e45a688e5c47e67df9e4095e07c8691c77491e038c3ff1401832832b17873e21",
							"linux-amd64":  "a4bc0352f7519f91b118320e948b7304e986767ec8c54f847e01e2cc7520a7ef",
						},
					}),
				},
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-unconvert:unconvert-asset:1.5.1",
						Checksums: map[string]string{
							"darwin-amd64": "223ab0d196a3732090d689a72921a4349d80bf508aa7362cf9d237c9dc3356e6",
							"linux-amd64":  "6ec5cca8bee1659306bb943212cfbbe4c051f1d67208a114252fdb1b098b8c38",
						},
					}),
				},
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-varcheck:varcheck-asset:1.4.1",
						Checksums: map[string]string{
							"darwin-amd64": "ebeafdaba0c8ca665f6696df1b833f1bad4effef14e3c99bb234b17e5a8a4adc",
							"linux-amd64":  "936fb8180cd763246ce4b8451b621147cc93dc8d6364c8c873cca7087423216f",
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
