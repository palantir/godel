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
					ID: "com.palantir.godel-format-plugin:format-plugin:1.63.0",
					Checksums: map[string]string{
						"darwin-amd64": "e38f1699edc8ac8f174eea4de7ce6be9c906a32cb3d3fc5d81cf350372330f01",
						"darwin-arm64": "63de0bd49defd48c439005fe05b4a64a844f9a919b26d1e1bd348792c71cc45f",
						"linux-amd64":  "feb4ebc2f19d93a08d7e431d7b45ddaad9d95f4936b64c8384771d1257f52a1a",
						"linux-arm64":  "b173f8748b334f25a211c91a90916eed5e24186218bd1a676cd6b4b77b0ec6e8",
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
					ID: "com.palantir.godel-goland-plugin:goland-plugin:1.57.0",
					Checksums: map[string]string{
						"darwin-amd64": "f4437383e0ee2e4802d5acc3e3fcb483800a6daa2bf43fbd1934191acbf8e7a2",
						"darwin-arm64": "125073255479a38097bf74ee64751b6385e6cd54399f6fceb877b11d8db1eb03",
						"linux-amd64":  "a5adc0566e4d5dbac77820d5331c5b00638459d0357276034aadf47a99c28863",
						"linux-arm64":  "b459649ca9e1adf0853a470e3a8a4f8bba9ff1dbfcfe7bd1c51eeec78a6a967b",
					},
				}),
			}),
		},
		{
			LocatorWithResolverConfig: config.ToLocatorWithResolverConfig(config.LocatorWithResolverConfig{
				Locator: config.ToLocatorConfig(config.LocatorConfig{
					ID: "com.palantir.godel-golangci-lint-plugin:golangci-lint-plugin:0.23.0",
					Checksums: map[string]string{
						"darwin-amd64": "e1d0df155d68a45d232877715b162ee9d94dc5076995af2ccc6cc9fd1a709aaf",
						"darwin-arm64": "e36fe99bbdebb205f5b1f4b46e3e02b08bd08164cb5499e104299b164014194d",
						"linux-amd64":  "0eca088efe5c9c5f2ce1f317ac8406aad38a5f526ee033dda9e19ea960d9c87f",
						"linux-arm64":  "2e66f8f6a984be457f56746428bdcbba6c2ab6cf36a84b7b5623950a4d926d81",
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
