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
					ID: "com.palantir.distgo:dist-plugin:1.86.0",
					Checksums: map[string]string{
						"darwin-amd64": "7b830a3f16ecfe19738654c23a10ec0be6457f8d1f6de1dd1715e0c0e0640503",
						"darwin-arm64": "b6e3780885261403aafc2a7d592c19d354bcd5672f95a885d75162d0122742ab",
						"linux-amd64":  "2d007e1f9a123fa547f687960772ad9a93f9caa24cedb4450c9b3240893f4b8f",
						"linux-arm64":  "bfc720b3044f5c8ec7900ac9dc6b8044cceee33c1f9160b8ab3101379619c708",
					},
				}),
			}),
		},
		{
			LocatorWithResolverConfig: config.ToLocatorWithResolverConfig(config.LocatorWithResolverConfig{
				Locator: config.ToLocatorConfig(config.LocatorConfig{
					ID: "com.palantir.godel-format-plugin:format-plugin:1.53.0",
					Checksums: map[string]string{
						"darwin-amd64": "88a3c1666c2bddc350fc6601a32500d87886c3fcf95c477b31619666cd04b0a7",
						"darwin-arm64": "230cd3ad76d61b9b20899d70249e3a0873e92ed4d4a4d5667d5b2093f004a9d5",
						"linux-amd64":  "aa3135ae1131e4ad55c2813fd39da039e3c90e1fdc53338802ca6aeee6a1f62b",
						"linux-arm64":  "eb08998b772564284b1c8205643c2a81304c85e7cc8b433902743f6592dfd1f1",
					},
				}),
			}),
			Assets: config.ToLocatorWithResolverConfigs([]config.LocatorWithResolverConfig{
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-format-asset-ptimports:ptimports-asset:1.52.0",
						Checksums: map[string]string{
							"darwin-amd64": "9284c15341ad6d0eb0b7c5c7b01b584288d7f141e2ba5050d2f9f36c59d1e3e9",
							"darwin-arm64": "70a79e4e796ca89e856e574afde6dc1f1d45affc7279379120d47689821d1ad4",
							"linux-amd64":  "ebc76db1e196238dc55f15c438fe0a4aa1b2b9c52db03f8589a7519c506e71ab",
							"linux-arm64":  "a12f69cb2c616d49a13235ed0724246891ac767d9e6e738a196fc0317abf0507",
						},
					}),
				},
			}),
		},
		{
			LocatorWithResolverConfig: config.ToLocatorWithResolverConfig(config.LocatorWithResolverConfig{
				Locator: config.ToLocatorConfig(config.LocatorConfig{
					ID: "com.palantir.godel-goland-plugin:goland-plugin:1.47.0",
					Checksums: map[string]string{
						"darwin-amd64": "7dec66ae76c870e778a6ff79396e85986914fc4f7d5083915dada0c2c12e921c",
						"darwin-arm64": "012ee49427b0ba3498cc51bab65c65e32e03b0f1a4d5eab42780f43a1076c15d",
						"linux-amd64":  "69ada20846ed9fda2658f89c0248b6534d3d834ba0bf3a1fc6fe5c2bbe715eab",
						"linux-arm64":  "acefcca1716ad4b5279fb481fc273e296b8ff05508b270b114032fadeb7610d9",
					},
				}),
			}),
		},
		{
			LocatorWithResolverConfig: config.ToLocatorWithResolverConfig(config.LocatorWithResolverConfig{
				Locator: config.ToLocatorConfig(config.LocatorConfig{
					ID: "com.palantir.godel-golangci-lint-plugin:golangci-lint-plugin:0.4.0",
					Checksums: map[string]string{
						"darwin-amd64": "5f03f7105d59153a4c2ed9bdfdb051253bffc5e54d1771c29171db847b048a88",
						"darwin-arm64": "f66933b8319485f8866a6d569ab307327a7b4ec617bcf1413db40d31ef4c9ad0",
						"linux-amd64":  "12609f407d97b0cc9e02972c756a12cf2db2f2623cfa3962d8c29849e7c92523",
						"linux-arm64":  "b70dd8f9492f1eb2669ace4bcedf0f6e4bcdfaee46d65d9d5cda275ffe7a0b42",
					},
				}),
			}),
			Assets: config.ToLocatorWithResolverConfigs([]config.LocatorWithResolverConfig{
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.golangci-lint-palantir:golangci-lint-palantir:0.2.0",
						Checksums: map[string]string{
							"darwin-amd64": "b5c875cf7d3a56c05e99d3b940a7acdf5a5b7267c6f61b841f65809e1fb82b31",
							"darwin-arm64": "298a7cdecd31bfc1353dcd9bf870adec719bedbf4062a24bb3d083faa0142a7d",
							"linux-amd64":  "66f5aa346580c7a5c1b6c144b71119cbc04f0dd45e16d7fa772d8d30e8e56bb2",
							"linux-arm64":  "cda12309614a77de366cc17256ad6496c1a6fec03fb0827647406b1b0be9a564",
						},
					}),
				},
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.golangci-lint-palantir:golangci-lint-palantir-config:0.2.0",
						Checksums: map[string]string{
							"darwin-amd64": "b71b2325328e3ae12cbe05c2dbb1dcef0baf18ea167a4d6a5c6d61a0068676be",
							"darwin-arm64": "b71b2325328e3ae12cbe05c2dbb1dcef0baf18ea167a4d6a5c6d61a0068676be",
							"linux-amd64":  "b71b2325328e3ae12cbe05c2dbb1dcef0baf18ea167a4d6a5c6d61a0068676be",
							"linux-arm64":  "b71b2325328e3ae12cbe05c2dbb1dcef0baf18ea167a4d6a5c6d61a0068676be",
						},
					}),
				},
			}),
		},
		{
			LocatorWithResolverConfig: config.ToLocatorWithResolverConfig(config.LocatorWithResolverConfig{
				Locator: config.ToLocatorConfig(config.LocatorConfig{
					ID: "com.palantir.godel-license-plugin:license-plugin:1.49.0",
					Checksums: map[string]string{
						"darwin-amd64": "9bed959c0964a6a4051fb41f80af3c4fe7233cc128b125bc15a965ea9437e5f2",
						"darwin-arm64": "542c40727e3ab1c0964020187676486ef9e185c4e23f6482bf3fa02fb80865f7",
						"linux-amd64":  "63791eb2a7f3aecd85e1f3b0fdb1895b794a435c4eb44511160dca528be6a3d0",
						"linux-arm64":  "d03b47ddd347f2416614f0d87cfb0f735f0da312028946e840f3e85216db7531",
					},
				}),
			}),
		},
		{
			LocatorWithResolverConfig: config.ToLocatorWithResolverConfig(config.LocatorWithResolverConfig{
				Locator: config.ToLocatorConfig(config.LocatorConfig{
					ID: "com.palantir.godel-test-plugin:test-plugin:1.47.0",
					Checksums: map[string]string{
						"darwin-amd64": "03f8252e502fc6e573324ae3c2584d3d0a7028bcc2b541e55628a02506b4c09d",
						"darwin-arm64": "cea1eff96d37b7df8a110ee393cf9e8a2dcb8b532b6f28601bab8b30bfecbf6b",
						"linux-amd64":  "658be83977f879518ca7d533498de8354945e3f3477378b3c6605ce006a7c02b",
						"linux-arm64":  "704dbeb574d54af01f5c2772c8af2f164633963341b8a9c1bf0255e245eb2089",
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
