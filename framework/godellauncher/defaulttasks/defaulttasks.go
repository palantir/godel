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
					ID: "com.palantir.distgo:dist-plugin:1.39.0",
					Checksums: map[string]string{
						"darwin-amd64": "c6eb0411ef72d1812ae9d7e99371d24a51f013d1803be56079896457be3a5a52",
						"darwin-arm64": "ac7ab68c07f87bff152c3939516a9d3f6c49c316ebd61560f488565a98766157",
						"linux-amd64":  "b3318e85462a33e5b63b2ba12373aa95020272d518e6f3d63d6c9578acd441d7",
						"linux-arm64":  "674f933fc89efa0b84041a8baacdac23fdceb82f1de51582c2d8e54fa3ce8064",
					},
				}),
			}),
		},
		{
			LocatorWithResolverConfig: config.ToLocatorWithResolverConfig(config.LocatorWithResolverConfig{
				Locator: config.ToLocatorConfig(config.LocatorConfig{
					ID: "com.palantir.godel-format-plugin:format-plugin:1.15.0",
					Checksums: map[string]string{
						"darwin-amd64": "93a5ca0a53fee6f9a438e2117f3bf6c04b8fee986bcc8bf9151716ce1acc9476",
						"darwin-arm64": "893dc5db541a8f3702c18ff713dd12cc224e8205a08befe313493f5252db273b",
						"linux-amd64":  "768263777de3f29f37ebc68a03b77e43f3a0351c7147ed9bf1c0cbcf62ed2e97",
						"linux-arm64":  "ab87a566706b6ecb87a05a60e0d6db4ad5ad52f1ea6d02718e6a61e8f19897d0",
					},
				}),
			}),
			Assets: config.ToLocatorWithResolverConfigs([]config.LocatorWithResolverConfig{
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-format-asset-ptimports:ptimports-asset:1.15.0",
						Checksums: map[string]string{
							"darwin-amd64": "8945af05d9c9854b4e3b4df14f01099be2cb9934b2747a3bc6f0f94266e8db2e",
							"darwin-arm64": "9a0abaee149dd6b644faa38c588081c7a2e9a1d519a78196636b69c9742ececf",
							"linux-amd64":  "bc0372cb0475937180a4df985c0b1184a2bddfa85f6b6be4787a9168c24cb0b1",
							"linux-arm64":  "b0327405bfedad73ce25b7339a1ef677b862052a897ad7f202402fbcda47a44d",
						},
					}),
				},
			}),
		},
		{
			LocatorWithResolverConfig: config.ToLocatorWithResolverConfig(config.LocatorWithResolverConfig{
				Locator: config.ToLocatorConfig(config.LocatorConfig{
					ID: "com.palantir.godel-goland-plugin:goland-plugin:1.10.0",
					Checksums: map[string]string{
						"darwin-amd64": "9bab017aa1a684723c190356f9e5442a069f4635f4323f0dcf3deb2d2b8f1fa3",
						"darwin-arm64": "edd7db3df2fc454c8aedc4c6e94cadfe17657c756420beaf06468620c38dba10",
						"linux-amd64":  "af726c2babee7acba5b2b34554868b26f73234d6d9fa406d1f2f46a4562b1072",
						"linux-arm64":  "685f3ed174caa89a11ba29cd1398d7218bbaee022628b437d38ad5f9efef6c19",
					},
				}),
			}),
		},
		{
			LocatorWithResolverConfig: config.ToLocatorWithResolverConfig(config.LocatorWithResolverConfig{
				Locator: config.ToLocatorConfig(config.LocatorConfig{
					ID: "com.palantir.okgo:check-plugin:1.20.0",
					Checksums: map[string]string{
						"darwin-amd64": "f3af15ad6d691a65a50053631f1326b109aa5358801b40699e6d7350c4462400",
						"darwin-arm64": "5d85a4946847ea4d478233f5024aebae44c9b2ee2ac1c6b78558317c0adc23f6",
						"linux-amd64":  "543de30ed0f306a60b0623db754565593c82aab677a8d95a29a5c8400c744a97",
						"linux-arm64":  "b55641fc86669a53bfa0ee15321c22a8da5db33ebfe9e5c14af5cd21c8d144a4",
					},
				}),
			}),
			Assets: config.ToLocatorWithResolverConfigs([]config.LocatorWithResolverConfig{
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-compiles:compiles-asset:1.18.0",
						Checksums: map[string]string{
							"darwin-amd64": "33eae98ad2079629d4692178ad967d2d2311a50c5448f0ec2da8ea043d5c3fd9",
							"darwin-arm64": "241b88145f04064cc5dc3f9365b9098be46ddca865ccaa918124c5dbefda42e5",
							"linux-amd64":  "0dff84ce145642546570b1b7661d40d7b3e061baedcbc278e52e104247bcd484",
							"linux-arm64":  "0ffee4d11c6a78538234551f6bac0add3370831fc82370a634027e7646e04229",
						},
					}),
				},
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-deadcode:deadcode-asset:1.17.0",
						Checksums: map[string]string{
							"darwin-amd64": "f3459eb6e42eedb836670605168dbd082026473aae08a26c8f007f1e61ea5d39",
							"darwin-arm64": "972091d2bba0695b726c1e5b7084f1dcc69bfdabe46be0685df7e0f12120f31c",
							"linux-amd64":  "d6abc945801f52433e4f96afbee23c911fbe3ba2f92007c248e175713c396ada",
							"linux-arm64":  "d668c9fbe53a0c61b50f05f840b80791eddd6229dd29c090174197155cfaf235",
						},
					}),
				},
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-errcheck:errcheck-asset:1.17.0",
						Checksums: map[string]string{
							"darwin-amd64": "539eef97a2249b966d1c0d2a04abc7fbad965968dd14536a047066b4ebad8be9",
							"darwin-arm64": "5d1dc0a19e2d2fb389993b361607baa509f101454685f09632fec25208282478",
							"linux-amd64":  "0951193fc2715b3f7a10eabd1ae907c8522012e6ed6e5cb45a61511b098d9c59",
							"linux-arm64":  "d376bc8f65df42c1d0c6a4474bd3ec4252aad37103c52efcce430e2a996258e1",
						},
					}),
				},
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-golint:golint-asset:1.12.0",
						Checksums: map[string]string{
							"darwin-amd64": "28a2aab4033b33cd96a97248612ad0b9a27895eb0ba05a6803ce41fbdb5d85f7",
							"darwin-arm64": "c8488d13243dddfc13269a34ed59545601e4daa09791e31a5658e0977c42bb2c",
							"linux-amd64":  "5ed0fdf066849246547fbe831bfdaaefdea6a099bfdd93fe1718c3c57d264f35",
							"linux-arm64":  "e3d39c513de3d4c6d7d60041db2a257a3ae91f3f6b1dbbdd7958a301dce5b3aa",
						},
					}),
				},
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-govet:govet-asset:1.12.0",
						Checksums: map[string]string{
							"darwin-amd64": "cd6f424e841c099968f716a17eb4096e45ea99e0e70b338122f58a19da2fddd8",
							"darwin-arm64": "7866143e20a1abf12fe1cb49a7a295356a540dbf4206f8a6ce97285bb46c72ba",
							"linux-amd64":  "7fdaa69b78841af0bc3333a2c06d2cb443283f472c8138759cb1397bfb1be83f",
							"linux-arm64":  "e133604f61d9d229bf1cdba5f972992c128d5d3d5d3d11ead6304027fc3bd489",
						},
					}),
				},
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-importalias:importalias-asset:1.12.0",
						Checksums: map[string]string{
							"darwin-amd64": "e4eea32e70294e39193a05726be8f802ca1ec5588d4be96e31ef7b2e7144eefd",
							"darwin-arm64": "55daf57b0570b36d01fd8f1a604b225420da6afda608c7859332ebb170385766",
							"linux-amd64":  "e7a63cfcfd831848b3a3939592c041cf0b587525f575ec446c61b95207873121",
							"linux-arm64":  "e09652c92f3ab754c1b58ebdbd5846b7727c5996863215c4d9cc25f2c5b85d12",
						},
					}),
				},
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-ineffassign:ineffassign-asset:1.15.0",
						Checksums: map[string]string{
							"darwin-amd64": "d64384d5e72d80afa542beee43ab0c145d38c32c8b54646f52e925e9e64d9a05",
							"darwin-arm64": "5d690d98b4bbfba6e578328bd7ac0ff378dce0635224ad46b6528ed28345c88f",
							"linux-amd64":  "8c0831213f64e510f5ff0c9cae4e5dd300432b509bd679aca28b44e388a2fe8f",
							"linux-arm64":  "5e5f6a7cc94fb9e9286699bbb16ac843b3d7e46f83d3e38fd5a214c557949a4e",
						},
					}),
				},
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-outparamcheck:outparamcheck-asset:1.17.0",
						Checksums: map[string]string{
							"darwin-amd64": "334e476ae19fb746f71f59bac7c616f4fcf8253710f99f633989047bc974e4bf",
							"darwin-arm64": "5b71b8f5bf7ca5698899a7270f3182abe7ca87288200742c993cef2d7bb44e07",
							"linux-amd64":  "27a0cf1a3eb6f3aea1817bd2b8829dc74c9b7480de58805492ae02b3f943d3a5",
							"linux-arm64":  "20c188c1fc4ba64e77a2585347f04b0f5f778d2ed10c6b817eea1102b6eb5e96",
						},
					}),
				},
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-unconvert:unconvert-asset:1.18.0",
						Checksums: map[string]string{
							"darwin-amd64": "2edd22aa8667d1d3f84495494428bf4b8a27fb413ad50012d64bcaa579f1a215",
							"darwin-arm64": "906e5256343d64ef62048ae219a720268cf5a2470a4569e70a82e32e7f209fb4",
							"linux-amd64":  "364cde5f4910dc5e1aa76c200836cc34c3a8a662b26b2c2b518c98a9ac526e5f",
							"linux-arm64":  "8ec20673ff49c20f55c600bc1fd45a9dc1297d67907387c53def48ec5831d1e2",
						},
					}),
				},
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-varcheck:varcheck-asset:1.17.0",
						Checksums: map[string]string{
							"darwin-amd64": "56af3e6f97ddf9d8e387705e934974ccd0dfae494a5611aad5dc137c9e43df63",
							"darwin-arm64": "e47692e728c4fa1c13f8a65e41707d24eef3e5ca04c25b23ab4b0a18c750b260",
							"linux-amd64":  "8be2fbff924ed63b29fb314d5b972ee1e54f26495e576ada88a910578957bb7f",
							"linux-arm64":  "aedce11e14549a68a3c618b0c6463611695bcaac32695a032fd244f8b82d3ea6",
						},
					}),
				},
			}),
		},
		{
			LocatorWithResolverConfig: config.ToLocatorWithResolverConfig(config.LocatorWithResolverConfig{
				Locator: config.ToLocatorConfig(config.LocatorConfig{
					ID: "com.palantir.godel-license-plugin:license-plugin:1.13.0",
					Checksums: map[string]string{
						"darwin-amd64": "92820db8b7c25f2f833dbf4d63f1322640f109f46dcea99d1e33f4e6a8c790d1",
						"darwin-arm64": "be466f11e0aaea851489e260220d83dc8a92e4c759475dd4da09367ae25578eb",
						"linux-amd64":  "8ed97ef9e58df773051ee5c1fa486ad0934ab8b3a6d9946f3679fe4a25d41c73",
						"linux-arm64":  "5a5060446b68353e011b8326e04a1c1f5e3e66b2411620b4132259318bc6f6f9",
					},
				}),
			}),
		},
		{
			LocatorWithResolverConfig: config.ToLocatorWithResolverConfig(config.LocatorWithResolverConfig{
				Locator: config.ToLocatorConfig(config.LocatorConfig{
					ID: "com.palantir.godel-test-plugin:test-plugin:1.15.0",
					Checksums: map[string]string{
						"darwin-amd64": "4041c6b42a1dba444090e8e58ced4197870182f275deb4aea1f01e998bba0749",
						"darwin-arm64": "e40a858630e1bf9cfd8a30daa4c2802e78b5501394e80081c1dbb875a8b24a2f",
						"linux-amd64":  "47497de4a4a58203f9a5c52d7ff1a3427ae214791ec7dc9bf18ed58739ab2350",
						"linux-arm64":  "e9718b2819f88e84345f570709436c5c894ed6cf2f6bfc1e8eede02656e07294",
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
