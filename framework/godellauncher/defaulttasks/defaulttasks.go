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
					ID: "com.palantir.distgo:dist-plugin:1.119.0",
					Checksums: map[string]string{
						"darwin-amd64": "17a3bc76d69383f35fa9df54ba45106c191683c7fb3a3df2b4aeca1c03631308",
						"darwin-arm64": "4cda95068f4be17c06554480df8d4e4e7626c0e1476a9c1fc6b92bbf85c7ec99",
						"linux-amd64":  "77cff2d0c86819f397a4cc9a81dc80d7a5204af8a79b9bb4dbe3eb982ec5c854",
						"linux-arm64":  "be45dddcf82751a2c7daee2ebf1665189c954192c0fa153647fee5d3bc7cc69b",
					},
				}),
			}),
		},
		{
			LocatorWithResolverConfig: config.ToLocatorWithResolverConfig(config.LocatorWithResolverConfig{
				Locator: config.ToLocatorConfig(config.LocatorConfig{
					ID: "com.palantir.godel-format-plugin:format-plugin:1.66.0",
					Checksums: map[string]string{
						"darwin-amd64": "091caf16269af22c74ad569675a60c121eb404b701b8c100634d2fbab3dfd60e",
						"darwin-arm64": "66415f7988eac1355d4c73b6009f3da9e61ca0a64370a6ea01d0543688beabc3",
						"linux-amd64":  "f6730a4c9a047a4fe8519444bb20b61b4f6413ceb02a99c57b3ed28dd3dd8b7f",
						"linux-arm64":  "6ccc872d676cf1663e332a759429b11fa2106d5ef720b374f362cda5cd2ba928",
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
					ID: "com.palantir.godel-goland-plugin:goland-plugin:1.60.0",
					Checksums: map[string]string{
						"darwin-amd64": "8fb2ff06e3510f539ad3a13d938acb741fc9cdcf6bbd72f84ae19b0a5e282e9e",
						"darwin-arm64": "9fbb8e5826bd6dfe2034732a760b8ca3acdf7defe21c64c60f5078a7c3ab97ba",
						"linux-amd64":  "bea858a6529962d1fd00388d326989aaf916789cef4832a0719b99d59da57aee",
						"linux-arm64":  "e13d1842bd88ae10bff7fdcb9691e54760ba854d95bc756e5dae4830368ed703",
					},
				}),
			}),
		},
		{
			LocatorWithResolverConfig: config.ToLocatorWithResolverConfig(config.LocatorWithResolverConfig{
				Locator: config.ToLocatorConfig(config.LocatorConfig{
					ID: "com.palantir.godel-golangci-lint-plugin:golangci-lint-plugin:0.46.0",
					Checksums: map[string]string{
						"darwin-amd64": "2f89f05fffcbe258ec906166c174bc4a878326c94ef1b28359e57b2496ff632d",
						"darwin-arm64": "3ad731688614035f3224ee65b4024f9ab80718b78615cafd37489f558d743a01",
						"linux-amd64":  "18cb1bf8f3237e78f0158984ca94e02227dd82669ec57e53900f73b05d30de83",
						"linux-arm64":  "f8ab9201e51fc2df924ed3c7abbb5bf4d6cedc6dfe45b847fa50fb218e5f39e6",
					},
				}),
			}),
			Assets: config.ToLocatorWithResolverConfigs([]config.LocatorWithResolverConfig{
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.golangci-lint-palantir:golangci-lint-palantir:0.63.0",
						Checksums: map[string]string{
							"darwin-amd64": "3f136e5344f9b19b15f93a71c91d331363d40058cec07324b9d7100ef251b923",
							"darwin-arm64": "3f43d1e5821f0ef809885cb2609e6c2875ad954e202f4b69fefc7c736283fcff",
							"linux-amd64":  "0e54d69a0083989e1bf43a60433a557d2b985ee232c5b3c40f85900ff23a8500",
							"linux-arm64":  "d6d70887b29f9a03c76c4cbb340b3b56568785555b9c042c8dc6606a36d165d3",
						},
					}),
				},
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.golangci-lint-palantir:golangci-lint-palantir-config:0.63.0",
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
					ID: "com.palantir.godel-license-plugin:license-plugin:1.62.0",
					Checksums: map[string]string{
						"darwin-amd64": "7e6362b6f079e951c4e45f0b7e1b27bf2592db34a1cc445ad20c8bdedbc7f205",
						"darwin-arm64": "8f1e8bde4175d1db68b012a0f4b33d743b62e8d70bc7e1a6cfe48142ca08c2d5",
						"linux-amd64":  "e3c7cc2d284220e4bb774a3fe14a8b802f9e7d19982f94b5b25239fe5bbe9295",
						"linux-arm64":  "ce9a53917616a7319fcbc598ca38891b7eb37eee2c949018dcb06184359846c3",
					},
				}),
			}),
		},
		{
			LocatorWithResolverConfig: config.ToLocatorWithResolverConfig(config.LocatorWithResolverConfig{
				Locator: config.ToLocatorConfig(config.LocatorConfig{
					ID: "com.palantir.godel-test-plugin:test-plugin:1.62.0",
					Checksums: map[string]string{
						"darwin-amd64": "ebe7a49039a26cd8792bdcd887b35077125a85e17ba3e1f1a4c36d4e3686a4af",
						"darwin-arm64": "38db323c3d3fd03be1e5a6b64c29aab5623470f27e07fb909e5ba8bb2ce1fc2c",
						"linux-amd64":  "7594440785026ad474e67e98911b416fafa21b715fdd1d7c9003587dc8c949f5",
						"linux-arm64":  "7db2b2373cc39737874ccf43b87af0407190e84c1f952181ee61793259665f05",
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
