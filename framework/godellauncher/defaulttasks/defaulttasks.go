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
					ID: "com.palantir.distgo:dist-plugin:1.116.0",
					Checksums: map[string]string{
						"darwin-amd64": "d3b891189e7a2985495fc1c4b58d285140fbed0d8293fc65354e360514f754de",
						"darwin-arm64": "f97133130e9d5e513d29e0e7dd1753ba9f7c31df16992ce83d06376d84bb54bb",
						"linux-amd64":  "24005b9ee7b4a6f51191dd03d33ebfcbee8bcd5f370e71eda5a446e5cdcac11e",
						"linux-arm64":  "23513d17345ff20a8a752a12330b74586b07861a1528d4f969c372321bf4e6a9",
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
						ID: "com.palantir.godel-format-asset-ptimports:ptimports-asset:1.64.0",
						Checksums: map[string]string{
							"darwin-amd64": "6a74223775e037de44e784410d9d2c184805c8d609c5d4d27ab2283efd4e0ad3",
							"darwin-arm64": "2db933f316fa524e984190458501914e16e40fe542b9f31f004f69af663f6316",
							"linux-amd64":  "daa187c44f37193f83d33e00108b70c7080f7dad96bec97d6d9a38a02cd008a3",
							"linux-arm64":  "ed79e6beabbcfc5c12a5bf1d66c5c82a6cccb87affa3d8afd571fa60d004f124",
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
					ID: "com.palantir.godel-golangci-lint-plugin:golangci-lint-plugin:0.24.0",
					Checksums: map[string]string{
						"darwin-amd64": "592f8101674373f705eca0a960dc75dd617bacbb12d824b1d511ac16b0dce885",
						"darwin-arm64": "f94c241138fc7598bc0e671dd62f3e81cb3320f718cda92f4b6b3d118a4ccc3a",
						"linux-amd64":  "aea69648605861a3d86326d110868084e5396d4b0c102d3c71a3f9a96a67b6b5",
						"linux-arm64":  "a053e88c2d5e2ab13e4b14622ed76cd2d50d4da55612c00780d42c58009f75c6",
					},
				}),
			}),
			Assets: config.ToLocatorWithResolverConfigs([]config.LocatorWithResolverConfig{
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.golangci-lint-palantir:golangci-lint-palantir:0.24.0",
						Checksums: map[string]string{
							"darwin-amd64": "5d41e8e0113273ad0611b65625622caa6341601dc9e52dd4f771e3ee57a46b42",
							"darwin-arm64": "61d76c113254da5d1f3207f64466a2629f875f16a5b604c777c920740c159f45",
							"linux-amd64":  "b3d2965d49f01b2c4d6c9dc0493437834d83ae66e696cc31c27c8fa879c6feec",
							"linux-arm64":  "7778a8c30a9ab291d909ad7cd5aa18d587b4270439fa3e35f31b796e047dd1f2",
						},
					}),
				},
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.golangci-lint-palantir:golangci-lint-palantir-config:0.24.0",
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
					ID: "com.palantir.godel-license-plugin:license-plugin:1.59.0",
					Checksums: map[string]string{
						"darwin-amd64": "2fef6a54638dde3887d1b6612f3df866eabb1d31a8c8641f186010cefe6275c5",
						"darwin-arm64": "4292c8ff5b94bd23962a8d727b00df70ce3e961dc298e6576f4ce3e11157b808",
						"linux-amd64":  "efb4ed4fa920e54b6a4db74fc243099dda0ec8b1befdfac5e8a6682792bc72f0",
						"linux-arm64":  "9d26b52482e1cd38564e686b690b057087cbebab9869f47db4a50d178c818bdc",
					},
				}),
			}),
		},
		{
			LocatorWithResolverConfig: config.ToLocatorWithResolverConfig(config.LocatorWithResolverConfig{
				Locator: config.ToLocatorConfig(config.LocatorConfig{
					ID: "com.palantir.godel-test-plugin:test-plugin:1.59.0",
					Checksums: map[string]string{
						"darwin-amd64": "53d43a98d0575acaeb38f9cd9350c6a701b0901c71500e79f4e32734b98bd7eb",
						"darwin-arm64": "1bd6a890aaed601c81c20a1c568bd4e0453398e389850fe54320aa2b3f2206dc",
						"linux-amd64":  "93f322a89b5d4ed3725792037eade94f3c3ac9554915807a982ea0fcc86299d9",
						"linux-arm64":  "4989f2e3c417cab116e0a2c7945f38252b7b9868107582ef878800ced20ee744",
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
