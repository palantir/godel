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
					ID: "com.palantir.distgo:dist-plugin:1.2.0",
					Checksums: map[string]string{
						"darwin-amd64": "dbc2a1678daa48ed458e2bf5509e8b402e7d339927b02f84c867c049c3cdc7c7",
						"linux-amd64":  "41e5c62d0082c3f1aa44758e2a69b653ca5039272b4a0e8afb79fd95b769d552",
					},
				}),
			}),
		},
		{
			LocatorWithResolverConfig: config.ToLocatorWithResolverConfig(config.LocatorWithResolverConfig{
				Locator: config.ToLocatorConfig(config.LocatorConfig{
					ID: "com.palantir.godel-format-plugin:format-plugin:1.0.0",
					Checksums: map[string]string{
						"darwin-amd64": "1c6c7e06226efab7213f0b7f2e8ee61f59c0198d7cb72fcee76b2c7d9abddf36",
						"linux-amd64":  "f29f608a13eb1991357a29e6576c4010738fcd43e0930f5e507c02b2b87dfb73",
					},
				}),
			}),
			Assets: config.ToLocatorWithResolverConfigs([]config.LocatorWithResolverConfig{
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-format-asset-ptimports:ptimports-asset:1.0.0",
						Checksums: map[string]string{
							"darwin-amd64": "44f7d2ee0392cb30174a189a03d537cbc24f80f521a1110f74e5d54345d66d93",
							"linux-amd64":  "94ef1e5f8260cc783396e4e72034cb1dba6b25b33b7b2a21371b06ae17376d60",
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
					ID: "com.palantir.okgo:check-plugin:1.0.0",
					Checksums: map[string]string{
						"darwin-amd64": "29c47ae4370c1a21f3b7832ec055669d102d4eb019c0d2c5785fcedfe688b867",
						"linux-amd64":  "80d4fa038326048b57507a56a6adbb3d492182fb38093bf9483b27d4c83428b0",
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
						ID: "com.palantir.godel-okgo-asset-errcheck:errcheck-asset:1.0.0",
						Checksums: map[string]string{
							"darwin-amd64": "52500c7d9abb4f44e3ff5baabe39926622ac383dd1186bffd9d1b45128850b0d",
							"linux-amd64":  "6f508aae2c45dce4d6cdc2bbe3c9e53db195639a63966c96b05215e4a1b900ed",
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
						ID: "com.palantir.godel-okgo-asset-outparamcheck:outparamcheck-asset:1.0.0",
						Checksums: map[string]string{
							"darwin-amd64": "9612037f770679e7deb9ee3f11f939a7e565193b6b558759199ba3a0635f97b0",
							"linux-amd64":  "3eec093c4c62a5a809ee58acf6a76efeb86a35bdf5416397ff7f5ddef1b8a35a",
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
