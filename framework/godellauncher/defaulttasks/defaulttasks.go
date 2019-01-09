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
						ID: "com.palantir.godel-format-asset-ptimports:ptimports-asset:1.2.0",
						Checksums: map[string]string{
							"darwin-amd64": "1baaffabcf85bbc76b342b2d5e58a644de88276e0787bdbb1a7223fe96f199a7",
							"linux-amd64":  "faa42718dc4ec49a957e9b6f1ae1cd25f7b6dfc7fed1bb1b17d72be9d9c912f6",
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
					ID: "com.palantir.okgo:check-plugin:1.4.0",
					Checksums: map[string]string{
						"darwin-amd64": "f68e3d4971a85ce11bb77eff25e9974b99f494cdb55119339f998632d14a7dba",
						"linux-amd64":  "07353abf2342e2c930942c6cea42bf41ca45c2b84d287a91305bac53889de619",
					},
				}),
			}),
			Assets: config.ToLocatorWithResolverConfigs([]config.LocatorWithResolverConfig{
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-compiles:compiles-asset:1.2.2",
						Checksums: map[string]string{
							"darwin-amd64": "a0a8aac0ea80ac012242255f82ad9975f9beed5bf6b4f9d8c4a14a0972592b1d",
							"linux-amd64":  "f2089613c3f2561c330d753534de369178ca802540ec7d4df19d5e62cc56e677",
						},
					}),
				},
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-deadcode:deadcode-asset:1.2.1",
						Checksums: map[string]string{
							"darwin-amd64": "982287591aa0fc5ff8085299602dc8cbeb7cbbc6df49932b93bb014293156645",
							"linux-amd64":  "cf52e9da69c92f7d4a5e68649458201250d2ac59fef31bff671c11e01b3d3d74",
						},
					}),
				},
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-errcheck:errcheck-asset:1.3.1",
						Checksums: map[string]string{
							"darwin-amd64": "9848a7668cfc9fa39b0017e2714abc6a3a8a29d10972f121f1a70f8bf217c63f",
							"linux-amd64":  "e297a4c8e564ef9ce25dd3292dc984214ba10991636a0ddc5858b0eee5f3c3f5",
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
						ID: "com.palantir.godel-okgo-asset-outparamcheck:outparamcheck-asset:1.3.1",
						Checksums: map[string]string{
							"darwin-amd64": "683088d5216a01fb6a3a66c6c0c0af7b176e018af42a0ccdb24840411546a47f",
							"linux-amd64":  "cfe47fbc3d19c2f64340cd946918dad83dec0876519ce120b847ffc534a587ca",
						},
					}),
				},
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-unconvert:unconvert-asset:1.3.0",
						Checksums: map[string]string{
							"darwin-amd64": "4e26ec04ff39fc07ce994dd3446e14acd7b1b69d56ab33c3a80a7f44362bad1e",
							"linux-amd64":  "d7587376b26c35655fdb421b9208ae86896d0c2e2fde6a3d34a87bf7e04b7c2f",
						},
					}),
				},
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-varcheck:varcheck-asset:1.2.0",
						Checksums: map[string]string{
							"darwin-amd64": "594d801fcb90174a9dc6959835699569ee4977112137249d8195d1b944a36f0c",
							"linux-amd64":  "00e239bb6c20588c0f900bac4a4e353b4ad9f7bed54505c756b5a398096ce209",
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
