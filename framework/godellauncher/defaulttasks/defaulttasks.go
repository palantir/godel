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
					ID: "com.palantir.distgo:dist-plugin:1.67.0",
					Checksums: map[string]string{
						"darwin-amd64": "1f0d362dda2e099dd07452a706a2ee3dc6383c6efd3e65a7bdec02fefed7a944",
						"darwin-arm64": "7c5fe6bfae5edcf5a079014b16827ab7a782325b1c840e27ef913eca75505420",
						"linux-amd64":  "0469ecfe56bc491ed5e9a19bde5b7b646f0576ff6163595c940baa52b2575ad0",
						"linux-arm64":  "66a51ceb929c3f5faa851b632004184e89531f613688730ef2fb23bf7c30a8bf",
					},
				}),
			}),
		},
		{
			LocatorWithResolverConfig: config.ToLocatorWithResolverConfig(config.LocatorWithResolverConfig{
				Locator: config.ToLocatorConfig(config.LocatorConfig{
					ID: "com.palantir.godel-format-plugin:format-plugin:1.40.0",
					Checksums: map[string]string{
						"darwin-amd64": "e17a28aa753da06e9330744d8bb897e73e2a7049ef1d07f892fce449f1aaa68f",
						"darwin-arm64": "1af943226fb3745254677e321ac43aaa861cb54217fef3efb01af5a882024bf4",
						"linux-amd64":  "974d384374c0696df8569148647ae53fc1ad3e5ad2dd04c5a63bc798d7e8f4b2",
						"linux-arm64":  "69d4fc7228d36474b8f32caba8bc49e52124f658e85d4664fabfbf380bb98910",
					},
				}),
			}),
			Assets: config.ToLocatorWithResolverConfigs([]config.LocatorWithResolverConfig{
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-format-asset-ptimports:ptimports-asset:1.39.0",
						Checksums: map[string]string{
							"darwin-amd64": "85a48b1ab6bfb571c04fbf1f29a2b1d7fd6f8f3cc2baa1add7ab21b1e59ec0eb",
							"darwin-arm64": "4dc186f776a4a1df79d1384f0610c9d3022b7f736650439926a0b468148ad756",
							"linux-amd64":  "f2b4eb614ced30e833164ec98b147e61872c5631d67535b3d90949b086196d1b",
							"linux-arm64":  "b66199b3cc7f32f96affb13681f8bbf73ad0db0b7d2eb54b4551603b039f83a7",
						},
					}),
				},
			}),
		},
		{
			LocatorWithResolverConfig: config.ToLocatorWithResolverConfig(config.LocatorWithResolverConfig{
				Locator: config.ToLocatorConfig(config.LocatorConfig{
					ID: "com.palantir.godel-goland-plugin:goland-plugin:1.36.0",
					Checksums: map[string]string{
						"darwin-amd64": "d9743dfc65a3c28472147dee8c5a649741d5ad2c6be7f5601d3241826c681a33",
						"darwin-arm64": "13cedad519d41196a74d37d7de73c2fbf7405a53eed6850d15ef78275ba769d4",
						"linux-amd64":  "322e9830887e53930b5a6faf7dca2ad5d4a4ab36c5c9de8ca1ac312eb48b08ae",
						"linux-arm64":  "5572e239e35e9b2a7a4a53c9e5698f39eca5ec1d004097c516e5475b1dd366e5",
					},
				}),
			}),
		},
		{
			LocatorWithResolverConfig: config.ToLocatorWithResolverConfig(config.LocatorWithResolverConfig{
				Locator: config.ToLocatorConfig(config.LocatorConfig{
					ID: "com.palantir.okgo:check-plugin:1.52.0",
					Checksums: map[string]string{
						"darwin-amd64": "68bb0ab1aaf4bbabff7c838c18a478eaa29f07acee4b83e0e5d7c67fa9fba118",
						"darwin-arm64": "fe0793bcb7c40c5626e93dd9fd2cfdc29591de727923482b7c54a0ebf8c23a2c",
						"linux-amd64":  "b0d405694f97ad8ca2b7e953175d49005c8be7e88e41480b380837fbfeaf878d",
						"linux-arm64":  "d2d94c53865de85ddc32c3f81fd7bfc53bc2ca4b437939a1ce3d66e9562fcfe4",
					},
				}),
			}),
			Assets: config.ToLocatorWithResolverConfigs([]config.LocatorWithResolverConfig{
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-compiles:compiles-asset:1.45.0",
						Checksums: map[string]string{
							"darwin-amd64": "6c214fb54c319aefb84dcfce8a2d62d65dbc2eda655ae0996851aebdc277542d",
							"darwin-arm64": "f654bf5038f8c7b632d6136abbae6f1c5464b0f03fa8bbf6617e4d47fc56d5f7",
							"linux-amd64":  "f7abb928ff7e5cb857a6f00ef1e8fab45026dbc68fe86c434b08cc290dc8b784",
							"linux-arm64":  "75626060c3eeb8b96ce28b5395392e8c69d7c729b641b1d149bfc85ea45d8c25",
						},
					}),
				},
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-deadcode:deadcode-asset:1.41.0",
						Checksums: map[string]string{
							"darwin-amd64": "d23289b75cad1c4b4108602cefbe2688425d93e82a661c405b117ddc6f5bd18d",
							"darwin-arm64": "59e1cfa3bbe631d0b5dcab5de54c235a4dc4574d1ced6135909e6125dcb0edc0",
							"linux-amd64":  "681b85257385a66d1d65466b1e70593080aeaf8e9e65a029f2b383661a6edd37",
							"linux-arm64":  "fe8fa859039d509d82e1cdf1b2245614a2cd01e40159344f72a295b8136ed05e",
						},
					}),
				},
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-errcheck:errcheck-asset:1.43.0",
						Checksums: map[string]string{
							"darwin-amd64": "fa3ead0d14b269414db5051fb199ebbc44cded7ce36518a4eceb0939c1308ec5",
							"darwin-arm64": "948f266b226286b356792351bc592e4d406af9c6d6bff94ec3e5affcfe9308f6",
							"linux-amd64":  "e6a2df1665c77bbe397954bf3b579c3909abbe7439d020495d2fb8e3eab5926c",
							"linux-arm64":  "fe1cd6763be4d865d791718d30b19df3a51e87f037830791b37505108535f5c9",
						},
					}),
				},
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-golint:golint-asset:1.33.0",
						Checksums: map[string]string{
							"darwin-amd64": "a2570ee6b446ab97bbe3e51cb54ff60432a1d93f5c608e4bd28b8902663394d1",
							"darwin-arm64": "9d171731b84d6e6b7b6e4eb1ee9ad9272fa94a7334c271bef41ecc95bb083ee2",
							"linux-amd64":  "2f144a4ce2f1af204d071947b0662822de71d8103af0f4e231856d147d89df5a",
							"linux-arm64":  "53db72a692bc9cf8d501f591b329b19e187106cbdba451b5150ee535131e75bc",
						},
					}),
				},
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-govet:govet-asset:1.37.0",
						Checksums: map[string]string{
							"darwin-amd64": "bf736b38d8e158915edc568b49179c5b29da96c66a2cdb0549499b6ecfa57a84",
							"darwin-arm64": "037fdb95f1187adb38e048868358d211283485f43c9eb9ff4f649916583118aa",
							"linux-amd64":  "03930f7548ad6bc4c61f22af50a3e8c1a3e895d7ae1d9b9d9d9d74617ef58ddc",
							"linux-arm64":  "30a5323b98b5b3b18684f769f14058554ffc8c534cc727db4ea4b14e3be4444c",
						},
					}),
				},
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-importalias:importalias-asset:1.36.0",
						Checksums: map[string]string{
							"darwin-amd64": "792194a33f21d7b9d814297de5cbe76522ae7e2528edc653c1b977deb7cb88e6",
							"darwin-arm64": "e0996c773acfaf916be1034ce83944c1dd3dafc4f6776ea62bec224d89cb10e3",
							"linux-amd64":  "e4d6f3f55196e40c3e572e29196846b21b5f9110ad82fc3fcf05d514060d59f0",
							"linux-arm64":  "85e1f748dd829e608af9b043d7bb04124313ef7cddd358aa1beafdc394f0ea8d",
						},
					}),
				},
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-ineffassign:ineffassign-asset:1.39.0",
						Checksums: map[string]string{
							"darwin-amd64": "440ab24afd1fe8e529b97bbc66b183b9230d3888909bf99794ab19bf113f3c93",
							"darwin-arm64": "0cb6d801145b0c5f7c837bddf8d2609aae8d5477c99f4c4d8b51df1069e7eea0",
							"linux-amd64":  "049f653769fe9c51c2811b181c35e1045e00f12ab3efe3361bd26fe4845a0a2c",
							"linux-arm64":  "25ec7757442463579d4afdd0bee736269f63f2923649690d8a57f13c024d2515",
						},
					}),
				},
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-outparamcheck:outparamcheck-asset:1.41.0",
						Checksums: map[string]string{
							"darwin-amd64": "10f66d080e9d63e852bb5279340ee729c9ade792d127377e71211e711ae32dcf",
							"darwin-arm64": "5033a578733aa4d8b369dd50f0148fd1de6f6ef113ab7a024d977e5588ba8372",
							"linux-amd64":  "8b26ad008b90370fef219a1f377db55a77b9358fddacab4dcd83aed8d9176d33",
							"linux-arm64":  "f0a629aaeceb976d7a90ca6c7a2aa82c346e21acad38d83353cc1382ae55acad",
						},
					}),
				},
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-unconvert:unconvert-asset:1.41.0",
						Checksums: map[string]string{
							"darwin-amd64": "f91100ed466878c2d90ac0eb9015a7118c47851bc79da22ae501741188133457",
							"darwin-arm64": "5ffa8811a45c548d0dc6b2b564d03ea195af64d06320949c8fe7937a511e4cbc",
							"linux-amd64":  "d1bb7fa881c62c626d1c5004faf503f04af222507550251e9ac027a23683cf06",
							"linux-arm64":  "be61cb34fa0cb747e37432a078b20b2ad71b23c28ac82c77716016fa4cd18c9a",
						},
					}),
				},
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-varcheck:varcheck-asset:1.41.0",
						Checksums: map[string]string{
							"darwin-amd64": "696b297e23fda0ba14e5c17cc09e34bb43153b8e056157560e061f64617e976f",
							"darwin-arm64": "20d56d9bfbf9c8a80a3dddbcbc8097916ff9d783f370b988a4ec5da9ec3e4486",
							"linux-amd64":  "efff7e5d89425fb437bbb1b61eb78cfeb58c6b879001249af3636165e516c6ba",
							"linux-arm64":  "d37ca39359908896b6e01565ccac8e4bcfdc895529d46d7b5d8a551711e6ec03",
						},
					}),
				},
			}),
		},
		{
			LocatorWithResolverConfig: config.ToLocatorWithResolverConfig(config.LocatorWithResolverConfig{
				Locator: config.ToLocatorConfig(config.LocatorConfig{
					ID: "com.palantir.godel-license-plugin:license-plugin:1.39.0",
					Checksums: map[string]string{
						"darwin-amd64": "0bb1c73f447afcb7b20f7a1b24822816cd4dda83969d660b1d23bd6af722b247",
						"darwin-arm64": "c22150d2508585bf9c2314938b454c23edeb4aa65072d6a03765ce9527ea3fc4",
						"linux-amd64":  "4f6b20bd8bb366ea6297dfbd900c7d62a937403a6e861434af1d90e0d64eaeee",
						"linux-arm64":  "72800b7fcade406869e422d7c204537b541eb73424bf1d00e42639e7963b12e4",
					},
				}),
			}),
		},
		{
			LocatorWithResolverConfig: config.ToLocatorWithResolverConfig(config.LocatorWithResolverConfig{
				Locator: config.ToLocatorConfig(config.LocatorConfig{
					ID: "com.palantir.godel-test-plugin:test-plugin:1.38.0",
					Checksums: map[string]string{
						"darwin-amd64": "eb124ce49e48db5d53ed11cfb47e8a9e450f9dfe485368b222c2b94a6dec390b",
						"darwin-arm64": "275b9f409ead2bd1857602a975a5980fd125dcf33d6fb055ae1eb052433255f5",
						"linux-amd64":  "719fdbb28b2771c49658ee00dc17d339f998fad8415fdc430b72be371de02e7a",
						"linux-arm64":  "ae14f8616e10c1bf73a8a2cd6ae4993b6f00a7a70775cd43072b9a110343d081",
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
