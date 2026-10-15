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
					ID: "com.palantir.godel-format-plugin:format-plugin:1.64.0",
					Checksums: map[string]string{
						"darwin-amd64": "2694f518a0c10dbf629975351e94c08a53d984ff0a30661f82963b32bb0cf823",
						"darwin-arm64": "702e6df6b50ceb071e7730e63adab156fba34ba85abc2e0268aec106a571b38b",
						"linux-amd64":  "5e8298ba96d6afe224d43ad5a3d4e23c6378b424770e8042dc635fa8b06ed499",
						"linux-arm64":  "f0123b67deccd001c9629c9763acfe8af25f7cce585777db1a9964fbb683fa34",
					},
				}),
			}),
			Assets: config.ToLocatorWithResolverConfigs([]config.LocatorWithResolverConfig{
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-format-asset-ptimports:ptimports-asset:1.65.0",
						Checksums: map[string]string{
							"darwin-amd64": "c1ff760876f7773a6042e83a8070df89e3f484e1408f9ba8f0144e3342d21480",
							"darwin-arm64": "e335fe28295dedf40f22f9f26006081eb3920743831676382e04741f503b2078",
							"linux-amd64":  "5f5ba855b690d3e93c423db0036d86f79817033845daf154febaae3e4393666b",
							"linux-arm64":  "b2239b784fae571ca8aef4c54e2875d7669f0fc89a321b3b46118203dea1f043",
						},
					}),
				},
			}),
		},
		{
			LocatorWithResolverConfig: config.ToLocatorWithResolverConfig(config.LocatorWithResolverConfig{
				Locator: config.ToLocatorConfig(config.LocatorConfig{
					ID: "com.palantir.godel-goland-plugin:goland-plugin:1.58.0",
					Checksums: map[string]string{
						"darwin-amd64": "15d8f74162774031f6a99af80ae7359fb6ec948a78f5956f34b739433c460f2a",
						"darwin-arm64": "70be03ba54d05dcbef6591d41c3ed4d347d26d4772bbe2bc273c0f566c157568",
						"linux-amd64":  "2aad17e318e612c874282baafc507cb35392c42cde1e57606349ed8c4d6ecdb8",
						"linux-arm64":  "70d29de842231f3baa4e044e5f4136b551b8f2da02832785347e76dff3220ab7",
					},
				}),
			}),
		},
		{
			LocatorWithResolverConfig: config.ToLocatorWithResolverConfig(config.LocatorWithResolverConfig{
				Locator: config.ToLocatorConfig(config.LocatorConfig{
					ID: "com.palantir.godel-golangci-lint-plugin:golangci-lint-plugin:0.44.0",
					Checksums: map[string]string{
						"darwin-amd64": "d522454e9923b3524db04da57de003e7cdb809c3c5d7a9ac409f371cb08b163a",
						"darwin-arm64": "ee573996f30ceb083b6144ce2d31c860fe132b1fcee16e5cf2bcc052bb30e991",
						"linux-amd64":  "fcf6eba55e9e9c238c555d0c05c3bb6c186380f42ea01aabc17d26ca796de5dc",
						"linux-arm64":  "a6d61de36e9795a276253d825225530137b5e6a15fe7df0a60930c6c077125bc",
					},
				}),
			}),
			Assets: config.ToLocatorWithResolverConfigs([]config.LocatorWithResolverConfig{
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.golangci-lint-palantir:golangci-lint-palantir:0.52.0",
						Checksums: map[string]string{
							"darwin-amd64": "30d93cc9c6572908498d06d83859d98318708c81db4cf2b1175c3b5e2fc530d3",
							"darwin-arm64": "ef40c581ed522c767efa7ebd05c5ca4146dcf97d6d6ffd0fb191c85c36c53fdf",
							"linux-amd64":  "7a33c930eb115a433374097fd6d34c3e02f8efd3c8b656dd6ee1c268665b29e7",
							"linux-arm64":  "8a62430c3b236ba4a4220373da1cd994e274a8a3eaafcbcc97ced24ba7508c0a",
						},
					}),
				},
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.golangci-lint-palantir:golangci-lint-palantir-config:0.52.0",
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
					ID: "com.palantir.godel-license-plugin:license-plugin:1.60.0",
					Checksums: map[string]string{
						"darwin-amd64": "22e506c13e89567f98dd5419eba492cad7ea4f46b422eb1d6fe8c3f51c615234",
						"darwin-arm64": "3eb8d07cf7273307e2d3d170c4bfbb8ccd774155c4378704a04757c9909987e7",
						"linux-amd64":  "cdb5e319077c9510a68692183d5c72b1a8318eb8fa96b90cd4ce885d754165bf",
						"linux-arm64":  "b59dac2645c24d10f545594c978006ca0706d93f29538be21d22306ef4b091fb",
					},
				}),
			}),
		},
		{
			LocatorWithResolverConfig: config.ToLocatorWithResolverConfig(config.LocatorWithResolverConfig{
				Locator: config.ToLocatorConfig(config.LocatorConfig{
					ID: "com.palantir.godel-test-plugin:test-plugin:1.60.0",
					Checksums: map[string]string{
						"darwin-amd64": "252ea7b8f3584dfa4414f1b4bf5816072798eb4392415afb11735f9d2d3f3419",
						"darwin-arm64": "e1ba5e2a8e0fc3a1a2302448e68e5916af14aaa317f7b6374310bc7d3af6d8a6",
						"linux-amd64":  "fd439a436da483b15a43993607c27e4da29cf3b82c345f15390eed4bdacc2377",
						"linux-arm64":  "ddca45cb34d0abfd8af6f2ec004c7c5359d5b6004e8e121cad7edc9d56224571",
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
