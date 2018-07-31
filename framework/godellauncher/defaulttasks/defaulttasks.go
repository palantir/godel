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

	"github.com/palantir/godel/framework/godel/config"
	"github.com/palantir/godel/framework/internal/pluginsinternal"
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
					ID: "com.palantir.distgo:dist-plugin:1.13.1",
					Checksums: map[string]string{
						"darwin-amd64": "e9376b9beb4a02fa9b205df881bf32f9b943bfc28e9cc04fe592a257ed614e46",
						"linux-amd64":  "9ccb951367234a3126a9dacd45f2123aa787f3911ec0d5b131884454952fa1ed",
					},
				}),
			}),
		},
		{
			LocatorWithResolverConfig: config.ToLocatorWithResolverConfig(config.LocatorWithResolverConfig{
				Locator: config.ToLocatorConfig(config.LocatorConfig{
					ID: "com.palantir.godel-format-plugin:format-plugin:1.1.1",
					Checksums: map[string]string{
						"darwin-amd64": "30848937399398b7fd6acf206112b0fce129d3d79450f1aacdc6d77f91146001",
						"linux-amd64":  "6700496ced47596e802ca2016f03c890c3e685b17f559d98397533e433b5556c",
					},
				}),
			}),
			Assets: config.ToLocatorWithResolverConfigs([]config.LocatorWithResolverConfig{
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-format-asset-ptimports:ptimports-asset:1.1.0",
						Checksums: map[string]string{
							"darwin-amd64": "d5968b3c17ce1b2d83960d769ee5b14e5122cf37100c5ab12e4426787d07a8c5",
							"linux-amd64":  "47e2aec9306e6c6eb66e4e72dfd54f635d90b6bf13da6cdb683686700dbaca9d",
						},
					}),
				},
			}),
		},
		{
			LocatorWithResolverConfig: config.ToLocatorWithResolverConfig(config.LocatorWithResolverConfig{
				Locator: config.ToLocatorConfig(config.LocatorConfig{
					ID: "com.palantir.godel-goland-plugin:goland-plugin:1.0.0",
					Checksums: map[string]string{
						"darwin-amd64": "5b518708e5c41d81545d89d7224d2b61bf56d953eb560513ad047903eaa11b12",
						"linux-amd64":  "a000f7cd87f878d4c2e51e74f6015beb8fe48ea242c45f7731c7435a93f5a419",
					},
				}),
			}),
		},
		{
			LocatorWithResolverConfig: config.ToLocatorWithResolverConfig(config.LocatorWithResolverConfig{
				Locator: config.ToLocatorConfig(config.LocatorConfig{
					ID: "com.palantir.okgo:check-plugin:1.1.1",
					Checksums: map[string]string{
						"darwin-amd64": "ac6d56640e587d3ed6778a8591d4ab8749555537d6d6fbda97b2c386e66dfb2a",
						"linux-amd64":  "f03b1e1f1ed680ec004db3edf3eed74df1e87720910cc679516f8515d8345d03",
					},
				}),
			}),
			Assets: config.ToLocatorWithResolverConfigs([]config.LocatorWithResolverConfig{
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-compiles:compiles-asset:1.0.0",
						Checksums: map[string]string{
							"darwin-amd64": "d740b5ddf4befcb553ee90ec07a62ce322943069d48688cf48eefae6df45bd62",
							"linux-amd64":  "a01ecaeeb9093b5ba2838bde7ee041c3f181d3f286aeca740c3b39659a04e12f",
						},
					}),
				},
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-deadcode:deadcode-asset:1.0.0",
						Checksums: map[string]string{
							"darwin-amd64": "385437dd9e424b7372934a009cb13d2cba5442d9f6f5ce4d77b7831e789cca10",
							"linux-amd64":  "9d51dd2942a5b2758246ce9ddb7baecf097a9dc83e90aff8ccbd847b92fcc932",
						},
					}),
				},
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-errcheck:errcheck-asset:1.1.1",
						Checksums: map[string]string{
							"darwin-amd64": "58430b7ed1cfdc50d4450f421b07abd2539bbefa23a0f4c9b8cdef924db77b78",
							"linux-amd64":  "8df5114ca5086a5e0e1d86dfbc92d016c17bccd6a1be197f308a8d8eea6ae0a9",
						},
					}),
				},
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-extimport:extimport-asset:1.1.0",
						Checksums: map[string]string{
							"darwin-amd64": "af3757a10ef206ad2f21d1dba3cca0c798d4cd9c3a903e588c1589d93c3ad407",
							"linux-amd64":  "290008db565ad95924cc0cb07740f944bc84a1f0dcdff221f028df2a0eecd7e5",
						},
					}),
				},
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-golint:golint-asset:1.0.0",
						Checksums: map[string]string{
							"darwin-amd64": "23105fb709f06f241be59c8441afa3125191a7c4c3b6e68beabea56955b0a410",
							"linux-amd64":  "b29cdb47fba9127166306ff072285e2c71753b98c8e17fc4646379dc58dab17d",
						},
					}),
				},
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-govet:govet-asset:1.0.0",
						Checksums: map[string]string{
							"darwin-amd64": "8ce27439c844e3617dceeb9b6cb99218f73b17ef56d0eda8d6ea4434825b78ce",
							"linux-amd64":  "43a526f8800dbb6802755747c5a133010e4accaa1afcef39e2d14ffa39bb0dd8",
						},
					}),
				},
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-importalias:importalias-asset:1.0.0",
						Checksums: map[string]string{
							"darwin-amd64": "2453a34101380c82a02427cf04ebe0431ed14c37c87d6d5811adfc268b75f391",
							"linux-amd64":  "8737bf922f7771df845a92339f3479d53a509bab4ebb5b0893899738e800d21a",
						},
					}),
				},
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-ineffassign:ineffassign-asset:1.0.0",
						Checksums: map[string]string{
							"darwin-amd64": "a5722caf601170bf25a6a4d95b75016fd76272413b73541847df47c26c431662",
							"linux-amd64":  "66afe8103ccdd8c35f8f05a8200e5a258ea0cd35762dc3a8f1b7fcfb466e837d",
						},
					}),
				},
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-novendor:novendor-asset:1.0.0",
						Checksums: map[string]string{
							"darwin-amd64": "94e90c93487985b033382f6a92cd1bad410172dcbfa41b336189125be988d771",
							"linux-amd64":  "37b3d2ef4aae0e4a05bdf88062d0359474de81eb54d8542f169ba2377ca3369a",
						},
					}),
				},
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-outparamcheck:outparamcheck-asset:1.1.1",
						Checksums: map[string]string{
							"darwin-amd64": "ee30e5a6f703aafcafa87c02b8f2bd3841942c4727115322cae121bca3cb9fbb",
							"linux-amd64":  "eade1327d5e46bdb7e4e4de6adc5976726fab5e36b090b72d4a37d29a4474fc1",
						},
					}),
				},
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-unconvert:unconvert-asset:1.0.0",
						Checksums: map[string]string{
							"darwin-amd64": "a6c92dd8659d7b0824d4ef4e653002299a2e3323c1d441f1f9899737da81270d",
							"linux-amd64":  "f429c2361af530d835451edf8447393f221f21b4ad054cc4e44bf4eb51c4de75",
						},
					}),
				},
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-varcheck:varcheck-asset:1.0.0",
						Checksums: map[string]string{
							"darwin-amd64": "3a58db2a6c810c45985970cd85c15a529080b55e2f43c54bd777b4715abb2a69",
							"linux-amd64":  "ab07f4bcb182412696523b252f721c4ba502670cadb2bcf8fc2bef1abd0ed29e",
						},
					}),
				},
			}),
		},
		{
			LocatorWithResolverConfig: config.ToLocatorWithResolverConfig(config.LocatorWithResolverConfig{
				Locator: config.ToLocatorConfig(config.LocatorConfig{
					ID: "com.palantir.godel-license-plugin:license-plugin:1.0.0",
					Checksums: map[string]string{
						"darwin-amd64": "9b3b464ebfcad71718fe31dac47b192457d705e679d05fae564caa92661533f5",
						"linux-amd64":  "5e51f6df8bb3ee77ac25c1c9f9bf84498a34139590efee7ce9deda28d96d52cd",
					},
				}),
			}),
		},
		{
			LocatorWithResolverConfig: config.ToLocatorWithResolverConfig(config.LocatorWithResolverConfig{
				Locator: config.ToLocatorConfig(config.LocatorConfig{
					ID: "com.palantir.godel-test-plugin:test-plugin:1.0.0",
					Checksums: map[string]string{
						"darwin-amd64": "6a40ddbc0c6d1c0f705f455c8efcb76d05766c71c058caa938ead942d60bf190",
						"linux-amd64":  "95716e73a388865fd04881ed2cb7c521ffb9a0ffeb25a3090bd543d34c4983e5",
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
