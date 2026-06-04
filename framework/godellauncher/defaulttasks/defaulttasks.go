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
					ID: "com.palantir.distgo:dist-plugin:1.97.0",
					Checksums: map[string]string{
						"darwin-amd64": "",
						"darwin-arm64": "",
						"linux-amd64":  "",
						"linux-arm64":  "",
					},
				}),
			}),
		},
		{
			LocatorWithResolverConfig: config.ToLocatorWithResolverConfig(config.LocatorWithResolverConfig{
				Locator: config.ToLocatorConfig(config.LocatorConfig{
					ID: "com.palantir.godel-format-plugin:format-plugin:1.58.0",
					Checksums: map[string]string{
						"darwin-amd64": "ecd6259ddd51a799d0d65f4400a3c89c6cc40a06086860c36dca2f30e8a4e84c",
						"darwin-arm64": "7a06bc9816d2749719569ff4dc51b4da798e0dfecfc2457a43c41193be402555",
						"linux-amd64":  "7e21f96b42ef5ebe116d1c1a564f12460ab79d0b952d448c17bc58741ba02652",
						"linux-arm64":  "d38c80027dc2e4abb3ac6349d2905a21f4d996d1003071ac09e4b1a3248785c3",
					},
				}),
			}),
			Assets: config.ToLocatorWithResolverConfigs([]config.LocatorWithResolverConfig{
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-format-asset-ptimports:ptimports-asset:1.57.0",
						Checksums: map[string]string{
							"darwin-amd64": "7c240f3a1edfcdf1cf8526864a9fc9a7606a17e78cd142903ae7e80cc610d753",
							"darwin-arm64": "40fa613d5eb56c16300fc631e1cf9d560d2923143064759772c6eb3f301488ce",
							"linux-amd64":  "067d6a445bda373fc49b56d6534c1b2ebb72c19b32e8bdbf2ca1c7f9243eb0dd",
							"linux-arm64":  "2d5af4cae254409f203a80cfc31a00bf78d6751fe9e71ba22b45944ab8b11d8f",
						},
					}),
				},
			}),
		},
		{
			LocatorWithResolverConfig: config.ToLocatorWithResolverConfig(config.LocatorWithResolverConfig{
				Locator: config.ToLocatorConfig(config.LocatorConfig{
					ID: "com.palantir.godel-goland-plugin:goland-plugin:1.52.0",
					Checksums: map[string]string{
						"darwin-amd64": "708df25a3bfe47146b49bf273e22248e2ba35813b10e62ba26e8b26003f01ca7",
						"darwin-arm64": "528223169180f658c2bb636e0437b6b43d41ca59ec1b61e9b331d43c1dba321b",
						"linux-amd64":  "7427345326407a89feaf195cdf1cc863e85a5086a1b23ec0d4c991a833434ac4",
						"linux-arm64":  "70fd2a14534edf90509be6f57bc41697b42c09d489f423564f54cacd0a27e5ed",
					},
				}),
			}),
		},
		{
			LocatorWithResolverConfig: config.ToLocatorWithResolverConfig(config.LocatorWithResolverConfig{
				Locator: config.ToLocatorConfig(config.LocatorConfig{
					ID: "com.palantir.godel-golangci-lint-plugin:golangci-lint-plugin:0.14.0",
					Checksums: map[string]string{
						"darwin-amd64": "",
						"darwin-arm64": "",
						"linux-amd64":  "",
						"linux-arm64":  "",
					},
				}),
			}),
			Assets: config.ToLocatorWithResolverConfigs([]config.LocatorWithResolverConfig{
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.golangci-lint-palantir:golangci-lint-palantir:0.10.0",
						Checksums: map[string]string{
							"darwin-amd64": "",
							"darwin-arm64": "",
							"linux-amd64":  "",
							"linux-arm64":  "",
						},
					}),
				},
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.golangci-lint-palantir:golangci-lint-palantir-config:0.10.0",
						Checksums: map[string]string{
							"darwin-amd64": "",
							"darwin-arm64": "",
							"linux-amd64":  "",
							"linux-arm64":  "",
						},
					}),
				},
			}),
		},
		{
			LocatorWithResolverConfig: config.ToLocatorWithResolverConfig(config.LocatorWithResolverConfig{
				Locator: config.ToLocatorConfig(config.LocatorConfig{
					ID: "com.palantir.godel-license-plugin:license-plugin:1.54.0",
					Checksums: map[string]string{
						"darwin-amd64": "4e776941f721f6f91320a402f0ddb0dfa98f3050f51ab321135dad07cfb73038",
						"darwin-arm64": "e8c3de31eda2628db25bbbfafe3205ba3d199c0de414c2e4742fc823cdf5f62b",
						"linux-amd64":  "9048ec2a006d2f412e02575e654ee81fa08765c0d71535d0f3f747e808c73b01",
						"linux-arm64":  "5b243d45c4f068ff7940bc47451ca5a42c3211bb14ee9874b72d4d8a600154c8",
					},
				}),
			}),
		},
		{
			LocatorWithResolverConfig: config.ToLocatorWithResolverConfig(config.LocatorWithResolverConfig{
				Locator: config.ToLocatorConfig(config.LocatorConfig{
					ID: "com.palantir.godel-test-plugin:test-plugin:1.54.0",
					Checksums: map[string]string{
						"darwin-amd64": "f066fcf9d23973f3388f856c0c404635f7b1a3006b01f93040ab6a54d6a68719",
						"darwin-arm64": "a430dfb2aee19fc4eb4330983fc3068823676e636c9d8e844bbab3e72bf30dd4",
						"linux-amd64":  "2305f9540796f1935f2c0d456d37d0dc7f076426bb6a20f7b02d78a01e4771d6",
						"linux-arm64":  "3060bfbc1c35ef18072ebbf13beb4aa8e27df9b064e7f4f62c28c9a67747e188",
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
