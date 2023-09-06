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
					ID: "com.palantir.distgo:dist-plugin:1.60.0",
					Checksums: map[string]string{
						"darwin-amd64": "614dab54eaa0ff18ec827e4691bb2357f4593424ee3d304e0e9cce2b44e3aebf",
						"darwin-arm64": "c3bafee1c94ce33ed8b7720b8e12c55202f61d752f9c7af704630483b8a17521",
						"linux-amd64":  "1981f6281ca8e0820145fa534b70c00b5f5f8b3ef004506fdb99831af50099ee",
						"linux-arm64":  "10904ff30512945f5d9ae531ce9904603fdc175bb8ff1703f28219dfede5c1d6",
					},
				}),
			}),
		},
		{
			LocatorWithResolverConfig: config.ToLocatorWithResolverConfig(config.LocatorWithResolverConfig{
				Locator: config.ToLocatorConfig(config.LocatorConfig{
					ID: "com.palantir.godel-format-plugin:format-plugin:1.34.0",
					Checksums: map[string]string{
						"darwin-amd64": "a09b1cb2f8465c46af050c7c343287dad8bc52e3ab217d42a9d0cc6f19ff8cde",
						"darwin-arm64": "62fe9eda40acd391a04b40555119ab8cceef2d6446e46190f7485431db1da0af",
						"linux-amd64":  "96968633783184b7f8eb1bfd2ef9c6bfad329981d6b7b9a956b297bd0dde72d3",
						"linux-arm64":  "10887e19732692e9a5f8a403dd36d75143108f9dcf71813226019c69683ec080",
					},
				}),
			}),
			Assets: config.ToLocatorWithResolverConfigs([]config.LocatorWithResolverConfig{
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-format-asset-ptimports:ptimports-asset:1.33.0",
						Checksums: map[string]string{
							"darwin-amd64": "2071092eae166b7126c5cefd0d18ec4dd8320cb26621d144a9d29f088eeda723",
							"darwin-arm64": "42b9293c357340aa5740278fe5ef88be7f5be045b0040b69e68c02c0c81a98ad",
							"linux-amd64":  "ed9a45f2f280be6ff1495dce7b0e263a317739ec67405b8522019b90cf4dac25",
							"linux-arm64":  "a2b7e4fc8660fa7a84fa703ab5b2caf68dd4b1747d3a87588d59c7fab95186cf",
						},
					}),
				},
			}),
		},
		{
			LocatorWithResolverConfig: config.ToLocatorWithResolverConfig(config.LocatorWithResolverConfig{
				Locator: config.ToLocatorConfig(config.LocatorConfig{
					ID: "com.palantir.godel-goland-plugin:goland-plugin:1.30.0",
					Checksums: map[string]string{
						"darwin-amd64": "944d290fb017e6bbd9b374381af130e63d6dc355d0a58e8618cb42dbd46256de",
						"darwin-arm64": "00585389eb1377cbfae90b537ee555f1c5ee10f341b86cada112c77879f18c7c",
						"linux-amd64":  "3132bfbc8454cbc9ccec0a9e66edb2bb4c527bd08e6074422269eb1421fc1ca3",
						"linux-arm64":  "1823e50a6d84deb62064390040ddf56863965244243cec0a4c6cb1f1ec6804c7",
					},
				}),
			}),
		},
		{
			LocatorWithResolverConfig: config.ToLocatorWithResolverConfig(config.LocatorWithResolverConfig{
				Locator: config.ToLocatorConfig(config.LocatorConfig{
					ID: "com.palantir.okgo:check-plugin:1.39.0",
					Checksums: map[string]string{
						"darwin-amd64": "b0686225268a08c91d75354ea4b55af6bac50bd12d93410aafaa2e134bf3116b",
						"darwin-arm64": "1f57f118669d9fb1369927ca4f5fbfcc7ac8731f8f021bd7e789f2994a296d3b",
						"linux-amd64":  "3c5f8219aecad5e2be0115cca635190a9ce3f777f76459de375359e7d878dca3",
						"linux-arm64":  "12df17d53d9a9738070c01b12570404b4e8fa079a556223b016a2178c80dbb8b",
					},
				}),
			}),
			Assets: config.ToLocatorWithResolverConfigs([]config.LocatorWithResolverConfig{
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-compiles:compiles-asset:1.38.0",
						Checksums: map[string]string{
							"darwin-amd64": "15815626dde743cf4c3f7579f386cd8c1c2fc078fcaf2a8918ac76891c0f0fc2",
							"darwin-arm64": "b99200fd5c5d3af1a0e6333c8e2c043410a08d915bc99652a0a40f49ebdf84d3",
							"linux-amd64":  "e99e9decf3702d341edcf796db2cdbb570d436e5476fa1e61d7f13960574ad0f",
							"linux-arm64":  "e80ecc856c364fa6c439ba5c9cb4a54249a75334416418e5a89956c543cbfa24",
						},
					}),
				},
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-deadcode:deadcode-asset:1.35.0",
						Checksums: map[string]string{
							"darwin-amd64": "09a172dfaf734cf82e5233e0d1565698d02bd4e4e35176509135ccae9901b554",
							"darwin-arm64": "cb47cd490225f021f9f5ffb6f44b3dcf9ebbd0d11bf72e5e364f6d592feae819",
							"linux-amd64":  "4e93274cd2cead2b87e616f3d0587cbe727944542da3f94e0677036073700114",
							"linux-arm64":  "fdb2a279626e0fa94b33f39570d388fa2940531ec6d5f86999375b5f9b566de4",
						},
					}),
				},
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-errcheck:errcheck-asset:1.36.0",
						Checksums: map[string]string{
							"darwin-amd64": "6ed7a0389d14075bd29a617ee2773d79855f2a1ef14812977367f7db7db79ce9",
							"darwin-arm64": "ea7caacb10fe2f61975da3550d4ba32a83e3cf27499459f26c3e3c4aeb0cf2cb",
							"linux-amd64":  "cc47642b18e6e2e94faf9a6899ed42a07a091e013563c4ce7f23cd26262902af",
							"linux-arm64":  "defd3db3b0ee942b42520ec80776644c9913d623511f740577925d15c27ebce8",
						},
					}),
				},
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-golint:golint-asset:1.27.0",
						Checksums: map[string]string{
							"darwin-amd64": "bbb6fc667cc2c9834fda0575d580dc4dfc213c132339cb21ba8600e267511e7f",
							"darwin-arm64": "3d25d87c2da8ba245a69ad2269abe6e0b39bb77f4fedd10edd11294e7486fd1b",
							"linux-amd64":  "f3d318b3aa075e97f4e433988e13e7e2d10e0f89a43509c53e2c975a310cb039",
							"linux-arm64":  "b207f5b8e6968e7f1e90155e409bdda310e788f9e2aeed5f59fdcdb153cfa459",
						},
					}),
				},
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-govet:govet-asset:1.31.0",
						Checksums: map[string]string{
							"darwin-amd64": "e55e81d83e9d8bc843ffd7c8b7f5b1b7f4ae4297187a1f507274da2c0fb14ca5",
							"darwin-arm64": "341dcc196df7e92923bb97a18d55e301e8890aaff067235355982a885b6b9a0e",
							"linux-amd64":  "c6a47eb28eebd751e7959b13faeed452713dd28c64695b305cdb526c462083c2",
							"linux-arm64":  "4f601be53eb9c24be3e957b0cf61b8bb86b329f9de0148a0a7814a54fb057269",
						},
					}),
				},
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-importalias:importalias-asset:1.30.0",
						Checksums: map[string]string{
							"darwin-amd64": "2f40d0fa0eceacf8307df5240269d31f9ccf495e8c356dde04d033ac6f0aff73",
							"darwin-arm64": "0dd2948b40507aa5b9709c2af36aaf7bce6e3477ee0cf42aed10d2899ab23ee0",
							"linux-amd64":  "f2f9bd802f6667e61f531e2090666b4106dd21c85a4c2ff4462accbb7622704c",
							"linux-arm64":  "090fc9d6872ea8bf30ceb4d837b048f83ca39103622c03326f692e4bec103211",
						},
					}),
				},
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-ineffassign:ineffassign-asset:1.33.0",
						Checksums: map[string]string{
							"darwin-amd64": "58e874681498577cd7fa6147f44d30ac0dda3cdd72b767f748a6c4095351b175",
							"darwin-arm64": "ddd65fa8240fc395830ad233995b2ea8ef08a52c9ba7a638c46beff80c5bad65",
							"linux-amd64":  "4bfe7ac4eefad77841090ad172e57746ec1cd99f45520ee99b00361f3acc7245",
							"linux-arm64":  "311eaf3140116cbc57ea6f9ac08e07c76a85f832ef8bd450542d68030627b284",
						},
					}),
				},
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-outparamcheck:outparamcheck-asset:1.35.0",
						Checksums: map[string]string{
							"darwin-amd64": "498fbaf9099558ccd139ce09a0e4e0770e9a28ab6d6e3c1650819f1d8ad7c2c8",
							"darwin-arm64": "251a07ca3fa3fedfe6e59a372d52c9ee4da67e937ee9db1d1ca8e87680292ac4",
							"linux-amd64":  "bf3866013ecc23bc4105a5376f9d104cb33ea7218dfd5fcb713cbd670bdaf475",
							"linux-arm64":  "ff9dd77ded62551ab6c2c8b9ec1076e5a29ec4d4eff36b15b0e13d895b99cfff",
						},
					}),
				},
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-unconvert:unconvert-asset:1.35.0",
						Checksums: map[string]string{
							"darwin-amd64": "3da533871d4cf2a2a3316c52a8ffbdebb76fabde41e958a9e2d72fb96b35802c",
							"darwin-arm64": "23719bf526d65c8ec1c4b2af2964942835ee82efce16486249d4f9cf94b3166b",
							"linux-amd64":  "6906c4775c07a96a70b780382b4cb9f493af48e2e25dc6a78dbff4d663096868",
							"linux-arm64":  "e0befbfdbf9e865bc078d3ac25e6c8ffffbe86fff43def22d5b6710f403cbe7a",
						},
					}),
				},
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-varcheck:varcheck-asset:1.35.0",
						Checksums: map[string]string{
							"darwin-amd64": "9d183f075388a9bbe3b17b646a374be0546acce3a4d4316fac2e2fc0948d47ba",
							"darwin-arm64": "9a1d6cb8cc1e65ef5ccd42fb6264869cfdb9270a233c5ebbe034b513fde1df7b",
							"linux-amd64":  "a309e586f2bdcec915935c6abd01cb796d5da4fcb6ad9ac00fbadc440a481566",
							"linux-arm64":  "d7ffe78f37ff0bbb91c1aa4dde6cdca421c42f2d9b69aded6a1b21894f599698",
						},
					}),
				},
			}),
		},
		{
			LocatorWithResolverConfig: config.ToLocatorWithResolverConfig(config.LocatorWithResolverConfig{
				Locator: config.ToLocatorConfig(config.LocatorConfig{
					ID: "com.palantir.godel-license-plugin:license-plugin:1.33.0",
					Checksums: map[string]string{
						"darwin-amd64": "425d92ba373b6cf070102ca78496712a5a66ef7d1fa67a46f664ed12a2fe41c1",
						"darwin-arm64": "1103d7ce515d966905e04636c4fc36abc1b84c2e03b651a5bc2e5b37d429dcc5",
						"linux-amd64":  "38e7fdf87384756a4bd0fdbc8c249a9f8f812fe4258bbda7ae80eebf953f4d7c",
						"linux-arm64":  "31315c43329af68018f88f07eb0456c49388b005f016011944997bee1fef0d4c",
					},
				}),
			}),
		},
		{
			LocatorWithResolverConfig: config.ToLocatorWithResolverConfig(config.LocatorWithResolverConfig{
				Locator: config.ToLocatorConfig(config.LocatorConfig{
					ID: "com.palantir.godel-test-plugin:test-plugin:1.32.0",
					Checksums: map[string]string{
						"darwin-amd64": "0163fa650280e888fb3dc089db6a918ed7cdd752947794d8645aa7738df60ad2",
						"darwin-arm64": "4c9e6e7b931e02fcc381ac99170d6e8b645fa1eeca847db66bbd4c2b0782b4e6",
						"linux-amd64":  "9a2d0179c49da1034de5d6d4cb77ae3a335d3f6eef92f33c388f29d124fa7f12",
						"linux-arm64":  "06503caab50770c9daa00fd70119d14fe103d8de5810d6526fce75ea51a81be1",
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
