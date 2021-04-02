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
					ID: "com.palantir.distgo:dist-plugin:1.28.0",
					Checksums: map[string]string{
						"darwin-amd64": "db83a534814f33d43e477802fd40b73419d0ce62f7a8a70bea279180fc21e363",
						"linux-amd64":  "a86d4623dad3701954e14e07d0e364f55fe98291f700d6d873c0007f972e7d84",
					},
				}),
			}),
		},
		{
			LocatorWithResolverConfig: config.ToLocatorWithResolverConfig(config.LocatorWithResolverConfig{
				Locator: config.ToLocatorConfig(config.LocatorConfig{
					ID: "com.palantir.godel-format-plugin:format-plugin:1.7.0",
					Checksums: map[string]string{
						"darwin-amd64": "f083c0672fb80b1708bceb55a24a277bf15c54f56884296da59dc9bf1462c4a9",
						"linux-amd64":  "193f6d783c6de24017f56e725dfb51c2258e85fd48a244bf0a8b7108921497f5",
					},
				}),
			}),
			Assets: config.ToLocatorWithResolverConfigs([]config.LocatorWithResolverConfig{
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-format-asset-ptimports:ptimports-asset:1.6.0",
						Checksums: map[string]string{
							"darwin-amd64": "3ee53cf8ff82190d72e3c622be3698c03c8c899fcb3b67ec1a7c3f028d346eee",
							"linux-amd64":  "54fec6ab0376c5d69fdcfa0c5db86198ddb5f702ca0d3b112809238ca3534999",
						},
					}),
				},
			}),
		},
		{
			LocatorWithResolverConfig: config.ToLocatorWithResolverConfig(config.LocatorWithResolverConfig{
				Locator: config.ToLocatorConfig(config.LocatorConfig{
					ID: "com.palantir.godel-goland-plugin:goland-plugin:1.3.0",
					Checksums: map[string]string{
						"darwin-amd64": "d81be9d59cf4c63daeaad5f006e293a5fbcb5db280b240f8a21f9d3f3d71c355",
						"linux-amd64":  "f84147bef63aa2d84b3be8abc99294f90894e279992a3cdca99e419441fb6fde",
					},
				}),
			}),
		},
		{
			LocatorWithResolverConfig: config.ToLocatorWithResolverConfig(config.LocatorWithResolverConfig{
				Locator: config.ToLocatorConfig(config.LocatorConfig{
					ID: "com.palantir.okgo:check-plugin:1.10.0",
					Checksums: map[string]string{
						"darwin-amd64": "6ddff74f78e2689e81423f8383960e44d0fa543fdc5235cd63399d7e162fb29b",
						"linux-amd64":  "2805276b57dfe466ccd1b73e9f3114913683e5c64a98adaefd754d275f1f71d7",
					},
				}),
			}),
			Assets: config.ToLocatorWithResolverConfigs([]config.LocatorWithResolverConfig{
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-compiles:compiles-asset:1.7.0",
						Checksums: map[string]string{
							"darwin-amd64": "f2dcc3d7f946985f757ee3223928ed192c8a2ceb4ad2e975c1603d6a0911d424",
							"linux-amd64":  "6926a489032ba5b5ba5c32a5c43cdc331a8764b2aba7250d6a445adcdcaf74db",
						},
					}),
				},
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-deadcode:deadcode-asset:1.6.0",
						Checksums: map[string]string{
							"darwin-amd64": "864d58b5bf458a82d4d1895d3682f06c90dd26c6ac60e5f0b89f0b5fe804373c",
							"linux-amd64":  "5c5262b019c820debfe95328b89f0aac2275e4986a0fe5cf64365f2af2f916f1",
						},
					}),
				},
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-errcheck:errcheck-asset:1.8.0",
						Checksums: map[string]string{
							"darwin-amd64": "7b0e823d0e20d145fc146da1802b3b4032c60312f91313b7ad548b5bb177140a",
							"linux-amd64":  "2282506af22eea9396904aa20cfcfd36921e118096331777647f68be30b700ca",
						},
					}),
				},
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-golint:golint-asset:1.4.0",
						Checksums: map[string]string{
							"darwin-amd64": "52d37c0a02b87b2294adb24fe97715312a1875c0f235dc796ec307c20d7ed3b2",
							"linux-amd64":  "2455ba69528b81676a04677653ade2a28570b0ee8f816357858be355af6a74af",
						},
					}),
				},
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-govet:govet-asset:1.4.0",
						Checksums: map[string]string{
							"darwin-amd64": "2483818f56c194e8515a30a638aaf50db67e89801f8aa9c5f945bad471475462",
							"linux-amd64":  "6d9b6980565d2bfa25b2a9314c487545776ca2e3be01b5bf106fc6f3ac8db5d7",
						},
					}),
				},
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-importalias:importalias-asset:1.4.0",
						Checksums: map[string]string{
							"darwin-amd64": "e9179b43f589871e1532959649764124da7a3f6c0a1a9ea540d6c696e4d74ce9",
							"linux-amd64":  "ed041851159fa0e486e96202a24d6983e53c052d3c72f7308e695e032638e7d8",
						},
					}),
				},
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-ineffassign:ineffassign-asset:1.4.0",
						Checksums: map[string]string{
							"darwin-amd64": "f9f21c86bf23897eaf24ca970b798b06dee5f9be8d6c8a73ac699fd1c852fc6e",
							"linux-amd64":  "6c5e4f1400cd2517fbcbeb510980a94faf0417116b38d3610af37860701de993",
						},
					}),
				},
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-outparamcheck:outparamcheck-asset:1.8.0",
						Checksums: map[string]string{
							"darwin-amd64": "40ebed8bdab5db64d79b644da57a8817f25a8caedd1c1f5a3603ba22123a82c1",
							"linux-amd64":  "95778b5e6ecbeb47af3d9693a3366dd3eb3dee66eeb588175c94e402e58c4204",
						},
					}),
				},
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-unconvert:unconvert-asset:1.7.0",
						Checksums: map[string]string{
							"darwin-amd64": "f336b15a963c3b0205a3c2e0958c1fd34307673aa6a774e2098f533f06328d9d",
							"linux-amd64":  "27c4846c6e5752349179c5f4fae0b029e063f228370e62ac2e93daa5ef4e4c65",
						},
					}),
				},
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-varcheck:varcheck-asset:1.6.0",
						Checksums: map[string]string{
							"darwin-amd64": "2fe87fde123952c70bf4625bf691b982bc1f57e8e9eefc20e5fdcd3ead10b61f",
							"linux-amd64":  "9a40e061d3bf7a7bae1c406afceb1cde5ae85d71794f587dd3201b815f9427bd",
						},
					}),
				},
			}),
		},
		{
			LocatorWithResolverConfig: config.ToLocatorWithResolverConfig(config.LocatorWithResolverConfig{
				Locator: config.ToLocatorConfig(config.LocatorConfig{
					ID: "com.palantir.godel-license-plugin:license-plugin:1.5.0",
					Checksums: map[string]string{
						"darwin-amd64": "62c88e0b928ec801c552a77e30a011154c39ae8ccf1f5f54e8e393922ead0335",
						"linux-amd64":  "e50a141f97d3d447661e38b7acc8aea3c2d82673a9673fc4291a90a2b3d9a29c",
					},
				}),
			}),
		},
		{
			LocatorWithResolverConfig: config.ToLocatorWithResolverConfig(config.LocatorWithResolverConfig{
				Locator: config.ToLocatorConfig(config.LocatorConfig{
					ID: "com.palantir.godel-test-plugin:test-plugin:1.6.0",
					Checksums: map[string]string{
						"darwin-amd64": "63499d0fcf9b5fc9e62e159bf5bce9e390f739ce70e3114cfa0ef83971f3a2b0",
						"linux-amd64":  "0f3948fd76e27bbfeaaccab0784c737fce3ec1d9be4ae5f346bebedadbb42358",
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
