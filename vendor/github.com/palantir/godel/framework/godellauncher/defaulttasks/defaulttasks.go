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
					ID: "com.palantir.distgo:dist-plugin:1.16.0",
					Checksums: map[string]string{
						"darwin-amd64": "d653ca7f15ab383eb9c5080cda1f58d0e2e5c7757535971b692c4e3a97a9bea4",
						"linux-amd64":  "7921569175a339dad5b89fb4122235ada883db125315f7ba0ffac6e0c9c301e8",
					},
				}),
			}),
		},
		{
			LocatorWithResolverConfig: config.ToLocatorWithResolverConfig(config.LocatorWithResolverConfig{
				Locator: config.ToLocatorConfig(config.LocatorConfig{
					ID: "com.palantir.godel-format-plugin:format-plugin:1.2.0",
					Checksums: map[string]string{
						"darwin-amd64": "177b8c9d7323d762a2c350a6413f043b231480a47afc8e3d9efbfad04e41085b",
						"linux-amd64":  "1e8fbde9b7ab84731407a2cb7b0ba18beb1050383a235b4c53655d364795f646",
					},
				}),
			}),
			Assets: config.ToLocatorWithResolverConfigs([]config.LocatorWithResolverConfig{
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-format-asset-ptimports:ptimports-asset:1.3.0",
						Checksums: map[string]string{
							"darwin-amd64": "06c365d2487ecf86042d3bb590cee30d0c2cd839d2c40a03e1cadff4f030e65e",
							"linux-amd64":  "5af61a197fe26156fc0b969799e13a9158de3d7e051c1f4bfc88938250ca19a1",
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
					ID: "com.palantir.okgo:check-plugin:1.5.0",
					Checksums: map[string]string{
						"darwin-amd64": "ff92a09dcefc9f49d17e852389cffee9046f8d96356f80b5b0b40ef5d5b1c7a4",
						"linux-amd64":  "cfbaca6ddce977131e9b294d3fc1d02e6fceb2bfae510088996346d837f08d92",
					},
				}),
			}),
			Assets: config.ToLocatorWithResolverConfigs([]config.LocatorWithResolverConfig{
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-compiles:compiles-asset:1.3.0",
						Checksums: map[string]string{
							"darwin-amd64": "a42a0461c56f7c5de60ad7d3c85f8d32a3817ed3a31d0db74edfb01144ee6528",
							"linux-amd64":  "bfdd670706468b7d3746d1c2f4a6a01699617dd4320734b9ba293c38f55c26ed",
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
						ID: "com.palantir.godel-okgo-asset-errcheck:errcheck-asset:1.4.0",
						Checksums: map[string]string{
							"darwin-amd64": "b919f328e339328df732d6ffd8be9516b25c70e1760302a5ee0fbddd3ace55fa",
							"linux-amd64":  "8955b27817f08d31d27ff8575fc7008f384408220e56cc963fc0f83bdc261949",
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
						ID: "com.palantir.godel-okgo-asset-outparamcheck:outparamcheck-asset:1.4.0",
						Checksums: map[string]string{
							"darwin-amd64": "4f14a0d36d0c5ae110be378c91d4169b2d92cc09d3ba0dedd09e74c2d92dfdb1",
							"linux-amd64":  "736b8964286f53e6c6650a1c7d4dcf77d666ffb6b11876bbeae0840a0c3921ef",
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
					ID: "com.palantir.godel-license-plugin:license-plugin:1.1.0",
					Checksums: map[string]string{
						"darwin-amd64": "25234b18d1090e60af5ba05575f677e6bc8e66bd95500083537bd8cca7135e99",
						"linux-amd64":  "2911aa673a349d72deacf7452fc9694886e51d996b5288006d661cc6bd9f1f76",
					},
				}),
			}),
		},
		{
			LocatorWithResolverConfig: config.ToLocatorWithResolverConfig(config.LocatorWithResolverConfig{
				Locator: config.ToLocatorConfig(config.LocatorConfig{
					ID: "com.palantir.godel-test-plugin:test-plugin:1.1.0",
					Checksums: map[string]string{
						"darwin-amd64": "f1575cef034e9f8b7440f925cc00237845080b9ce4aa7dd6c9594d99059a7c88",
						"linux-amd64":  "93c1bc4b3087b5cf88be928058f29d922261b0d9ff5d8118cd3b8f2ec9050598",
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
