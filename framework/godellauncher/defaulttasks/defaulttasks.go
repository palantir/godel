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

const (
	defaultResolver     = "https://github.com/{{index GroupParts 1}}/{{index GroupParts 2}}/releases/download/v{{Version}}/{{Product}}-{{Version}}-{{OS}}-{{Arch}}.tgz"
	defaultResolverYAML = "https://github.com/{{index GroupParts 1}}/{{index GroupParts 2}}/releases/download/v{{Version}}/{{Product}}-{{Version}}.yml.tgz"
)

var defaultPluginsConfig = config.PluginsConfig{
	DefaultResolvers: []string{
		defaultResolver,
		defaultResolverYAML,
	},
	Plugins: config.ToSinglePluginConfigs([]config.SinglePluginConfig{
		{
			LocatorWithResolverConfig: config.ToLocatorWithResolverConfig(config.LocatorWithResolverConfig{
				Locator: config.ToLocatorConfig(config.LocatorConfig{
					ID: "com.palantir.distgo:dist-plugin:1.118.0",
					Checksums: map[string]string{
						"darwin-amd64": "c5149135fd164f3633f5898cbf073a053a6230617b5555fef970da40f6239141",
						"darwin-arm64": "38b7164b3954c6f4b40e21a17f0bfff6e6356da2e7818ba409ab1416fe725a4e",
						"linux-amd64":  "f510e1771c3b96824448371d7649ffe9bcb23789e50b8ad58b724fc07c78a8a0",
						"linux-arm64":  "810e3a39b9ec176bd53693c1ef93bf6003698140516add76baa1ac96c4f5b541",
					},
				}),
			}),
		},
		{
			LocatorWithResolverConfig: config.ToLocatorWithResolverConfig(config.LocatorWithResolverConfig{
				Locator: config.ToLocatorConfig(config.LocatorConfig{
					ID: "com.palantir.godel-format-plugin:format-plugin:1.65.0",
					Checksums: map[string]string{
						"darwin-amd64": "7dc153d1678e484f38819b97017a135113e10e6c75980908b318fdc470556545",
						"darwin-arm64": "5d50435a3b96346940328f9e63c697f37c7d38e943c6a697453aa7115b11a7cf",
						"linux-amd64":  "31f4bda07205d32b763d36ee8926302cab5ff83f2fd2a7fcf82fd22e6cbe4a93",
						"linux-arm64":  "f7ab4ea96c45d05548bd1f0b3fb4c15d17a51352a57d49f17041ddffcdd17f3b",
					},
				}),
			}),
			Assets: config.ToLocatorWithResolverConfigs([]config.LocatorWithResolverConfig{
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-format-asset-ptimports:ptimports-asset:1.66.0",
						Checksums: map[string]string{
							"darwin-amd64": "0521c89bfd544385a02a8437088c1b3990600deae6e6b1edd6bc5cb3ee762c5e",
							"darwin-arm64": "a32ca179733ea6071d59b4eaec820b4119770f4f0b8feb21f7d3c381122ec9d8",
							"linux-amd64":  "24aa41eba11520a3452b65b584b1e7c3b3f73efc9bcea43bf4fd2fdc8634eda8",
							"linux-arm64":  "2d2ad56789f1675f95f3a8ea412e18e2723ecdaf8598d5e48042c0b940f924e5",
						},
					}),
				},
			}),
		},
		{
			LocatorWithResolverConfig: config.ToLocatorWithResolverConfig(config.LocatorWithResolverConfig{
				Locator: config.ToLocatorConfig(config.LocatorConfig{
					ID: "com.palantir.godel-goland-plugin:goland-plugin:1.59.0",
					Checksums: map[string]string{
						"darwin-amd64": "168ace054cc0df2c85e1e2bdd924325ae0f1bcf2ec9b2f94ec3214d9a59586b0",
						"darwin-arm64": "4dbc8f7ecfe0efd3145c31b63267d60a9afb1b19873d58f895d1ca6a57cc57d0",
						"linux-amd64":  "d68cf0bbb81990270ccd09c2610eeca354b5ef17423a307cc00a8d16c16511ab",
						"linux-arm64":  "3ca8092fe0132113ccadd5a9b89438c25bb12e91f8a7fa7bd554366d4565e4ce",
					},
				}),
			}),
		},
		{
			LocatorWithResolverConfig: config.ToLocatorWithResolverConfig(config.LocatorWithResolverConfig{
				Locator: config.ToLocatorConfig(config.LocatorConfig{
					ID: "com.palantir.godel-golangci-lint-plugin:golangci-lint-plugin:0.45.0",
					Checksums: map[string]string{
						"darwin-amd64": "e1b6bad0e63f1cf8ded256ee1573de9c6c45085846acab663a57b03d7bf49b38",
						"darwin-arm64": "6bc3abfd515a002318be0b850ea5c68f7a8a9f1caee5f9864b6c6868f7690e5a",
						"linux-amd64":  "98676338aceb8b897718872b19b729fd810d8fc72f4b91963e45d59b544c222a",
						"linux-arm64":  "b6db1adb0916b6b96dd6302c27279c000ace276d1e9c391eeeb4ae68c803c67f",
					},
				}),
			}),
			Assets: config.ToLocatorWithResolverConfigs([]config.LocatorWithResolverConfig{
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.golangci-lint-palantir:golangci-lint-palantir:0.62.0",
						Checksums: map[string]string{
							"darwin-amd64": "3cc0a7d2c3d388e6ced3e7d476cd4b261f2f6f815dcde02c22159f778f080f36",
							"darwin-arm64": "03e8a6051adcdd7521c2b2c2662f106099901eb88019bf2e5366c9f8945826af",
							"linux-amd64":  "cff8e4bcc8740959d5794fa95ba2184029a58347b98a0ee0a8908f71daa67fc0",
							"linux-arm64":  "47cdd552565a052ea0e524a0be7a1228455dcf662ec87f778d7792712654efe5",
						},
					}),
				},
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.golangci-lint-palantir:golangci-lint-palantir-config:0.62.0",
						Checksums: map[string]string{
							"darwin-amd64": "b297a9384c93a7f63bd1cda880e9827bd04e7cf12ef2dac1ba8cb6426e1a175c",
							"darwin-arm64": "b297a9384c93a7f63bd1cda880e9827bd04e7cf12ef2dac1ba8cb6426e1a175c",
							"linux-amd64":  "b297a9384c93a7f63bd1cda880e9827bd04e7cf12ef2dac1ba8cb6426e1a175c",
							"linux-arm64":  "b297a9384c93a7f63bd1cda880e9827bd04e7cf12ef2dac1ba8cb6426e1a175c",
						},
					}),
				},
			}),
		},
		{
			LocatorWithResolverConfig: config.ToLocatorWithResolverConfig(config.LocatorWithResolverConfig{
				Locator: config.ToLocatorConfig(config.LocatorConfig{
					ID: "com.palantir.godel-license-plugin:license-plugin:1.61.0",
					Checksums: map[string]string{
						"darwin-amd64": "01da27466dbbd3c92490905b3de1971ef89c5d9d9a31b281b51b431e32701348",
						"darwin-arm64": "c1f58d3612ed23278b494282a7c6455415fa1b05752ef37414d0ac236fee5631",
						"linux-amd64":  "2072e2ec5a88ec380b420c3d58a7fa25656ab6b75a954a4712ebc5998cc9c470",
						"linux-arm64":  "c455d28c4ab834896e107afa7a49e08720b12213cafb04e7412de1b0297a2936",
					},
				}),
			}),
		},
		{
			LocatorWithResolverConfig: config.ToLocatorWithResolverConfig(config.LocatorWithResolverConfig{
				Locator: config.ToLocatorConfig(config.LocatorConfig{
					ID: "com.palantir.godel-test-plugin:test-plugin:1.61.0",
					Checksums: map[string]string{
						"darwin-amd64": "5e3bff96c8a9f20adfba63dfe682cb4783a8dbfbde308465475635439fb855fb",
						"darwin-arm64": "d994f8dc62798e66881a9687f3c851d7f4d1fab23cf99565f4a3f2e0f3f61467",
						"linux-amd64":  "7dedf61715b247c85ce651d5e67da98bf64b4631f8b4f36c75b5e5ccffa1377e",
						"linux-arm64":  "4822731c3ffc0a3cfe32a47369cdf5c383cef819f2ce51766338fa606150623f",
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
