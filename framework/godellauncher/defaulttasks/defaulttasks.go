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
					ID: "com.palantir.distgo:dist-plugin:1.45.0",
					Checksums: map[string]string{
						"darwin-amd64": "f1a75a7158248c7136fff97629cce389f8c639cfcd9ced6feb3a4c73d1845d0c",
						"darwin-arm64": "72e75636029f60d64ac8029792c92082c10d8c12fd04a823d58de9f9b871870d",
						"linux-amd64":  "84d27586b00cee4ea4c00ef85fa64d69baec00aecfc13c486d9a471460408ae6",
						"linux-arm64":  "750bffa429654c15ad9802c669d729fc0e2a6331861b15508ee3e7fa24e11d63",
					},
				}),
			}),
		},
		{
			LocatorWithResolverConfig: config.ToLocatorWithResolverConfig(config.LocatorWithResolverConfig{
				Locator: config.ToLocatorConfig(config.LocatorConfig{
					ID: "com.palantir.godel-format-plugin:format-plugin:1.23.0",
					Checksums: map[string]string{
						"darwin-amd64": "454676ce35e77307485869f5db14be8306ba158fb18c620e6ec61ef65e58e34e",
						"darwin-arm64": "e8357905248973b555239678b87b84b0152c0f449f0443790b5b4133d46f8e5a",
						"linux-amd64":  "b945673eee1bf311429afb7d12205ed485841d5c94c3521eb2240846a89a5794",
						"linux-arm64":  "259fb3d10f3b0a698625307c2f21e87e5f4796561e88cc213532bd5bddf1d679",
					},
				}),
			}),
			Assets: config.ToLocatorWithResolverConfigs([]config.LocatorWithResolverConfig{
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-format-asset-ptimports:ptimports-asset:1.22.0",
						Checksums: map[string]string{
							"darwin-amd64": "af2e1b4b9a1a96e247767d40485c66796bd6040e6571cc35db23a65bcb2c4bac",
							"darwin-arm64": "aa5275988bd9407609e8788eda5a462f9e2f3c11fb95d15da2faae1063f526dd",
							"linux-amd64":  "7bda7b26d32d8cabd5c040ef6040a486d5cd447c105d1f24c22ba00764cfea8c",
							"linux-arm64":  "434f30d81951ebf08f6eea1036621f37534ddc0a4df15718c20d6043785e3da5",
						},
					}),
				},
			}),
		},
		{
			LocatorWithResolverConfig: config.ToLocatorWithResolverConfig(config.LocatorWithResolverConfig{
				Locator: config.ToLocatorConfig(config.LocatorConfig{
					ID: "com.palantir.godel-goland-plugin:goland-plugin:1.19.0",
					Checksums: map[string]string{
						"darwin-amd64": "fdde289cdad15007fb4cb23e4adf870e0148ce1c7d5839d395d5e4f78ce0dfd4",
						"darwin-arm64": "27acd2fb7c44128e275f1e10e78fdb912e67e3a74ce00a0de61fe96d19a9a4e7",
						"linux-amd64":  "4f72b8066121fb5c238af08c3eba48097a2b406be94dc3cccc90a6c0845f5854",
						"linux-arm64":  "9da571ea244502624376cd210b88edff6db8e224d3c7467391ed6bd35785f6b7",
					},
				}),
			}),
		},
		{
			LocatorWithResolverConfig: config.ToLocatorWithResolverConfig(config.LocatorWithResolverConfig{
				Locator: config.ToLocatorConfig(config.LocatorConfig{
					ID: "com.palantir.okgo:check-plugin:1.28.0",
					Checksums: map[string]string{
						"darwin-amd64": "9180fd080deae0e0a3a9e5a93b302086436ebbba7780f35816e5eec02ceb4c16",
						"darwin-arm64": "68dbcf4cab49b52cef444d3ecdb6be908186c1b6c76dd9fc2be6a21c6370bf16",
						"linux-amd64":  "c615f530308d2aa502afdb3ff6608f899daa6c1e83e345cb1b9cb0509a1a9ab6",
						"linux-arm64":  "0fa1b12e36751064b9c5000a8b3cb196db36352efb8612cc4673268f86ccdb98",
					},
				}),
			}),
			Assets: config.ToLocatorWithResolverConfigs([]config.LocatorWithResolverConfig{
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-compiles:compiles-asset:1.26.0",
						Checksums: map[string]string{
							"darwin-amd64": "52866a61b4bbda13fbea2302a2d1b80dc1ba037f61c149c4106005d508b1d936",
							"darwin-arm64": "0e5a452dc5abd0d200b897894022a6f2c47dce7a10cdd35d02c0dd20a784dea8",
							"linux-amd64":  "d1a5223144ec8a371db5b69bff66c5b08505f7fdecd56499e59e57743ce06ae5",
							"linux-arm64":  "cf6a5bd92fc57a4d268c3287cd0aab25c1daba7fc56c14efaff1435bcdff8e9f",
						},
					}),
				},
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-deadcode:deadcode-asset:1.24.0",
						Checksums: map[string]string{
							"darwin-amd64": "967ef0a2e3b010690afbc2fdc3ee2d0aaec31a191ac3e3ea7d5c99e41a843ddc",
							"darwin-arm64": "025ac7142ff811cde8fb5910e1b388dc12a483506b680f3a1ce5f842bcac1d38",
							"linux-amd64":  "aa58a89243ebf3c8d3d8487314668b1ed70b7cc55f44d5529e2fd9b1d7fcae8a",
							"linux-arm64":  "9f24cc7607b06703efaff9dd587608a1ceb590993d1c5bbb2c36649215b631a9",
						},
					}),
				},
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-errcheck:errcheck-asset:1.25.0",
						Checksums: map[string]string{
							"darwin-amd64": "c28a356728eb9e369e8923d91326747ae25a51dee2ca9880b458ffc07bc8983b",
							"darwin-arm64": "0b0a2a690aece465f595d11bd13816bd7e57962f996965e653f1dc36eb0342c6",
							"linux-amd64":  "13bb2785ef5ef1ec687122a38645a87eba8af2be48ab2c1c28dbe2d33a4e2cf5",
							"linux-arm64":  "bb1e1f91f960da7c2cf787eec2b45cbe3c86c31decf0ecc3a6617602071f723d",
						},
					}),
				},
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-golint:golint-asset:1.16.0",
						Checksums: map[string]string{
							"darwin-amd64": "8a562f7ff158def91800fe615eaf08fe5bbf8636fe04b97062d127657d774d5a",
							"darwin-arm64": "44f007302ac3da832c4fe52a14cd88c46e1e5c2c589a557910bac3a3e506173a",
							"linux-amd64":  "6f5885157ca0bfbb7147ffbf6783519ef860f157d72ca9a5f08667f0f9530b29",
							"linux-arm64":  "1c6276b3c5f7b77d6eb089e7f878f29680e79c33a192671e3195865b27293ed5",
						},
					}),
				},
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-govet:govet-asset:1.20.0",
						Checksums: map[string]string{
							"darwin-amd64": "6d9f17f46f764ab87e3ddf8f0b8b8b1417eeb38f342b8a14ed47865825f66fd4",
							"darwin-arm64": "aa2ffc81fe9a2bc33fdd0aa160b11abcca01ba4dc7fd191895b1705332b8f9d8",
							"linux-amd64":  "614fa0d4450c74fb4d6d83b272f2e0d64d022cdbbe5f85e67ebcce5e2565d20a",
							"linux-arm64":  "8595f4ac59b78ea982e262ee1324cbf50127f4067fdf43b62f6a5eec715555c2",
						},
					}),
				},
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-importalias:importalias-asset:1.20.0",
						Checksums: map[string]string{
							"darwin-amd64": "bd52b96c52ac108559c6ae2092e931b8b3683987f17bddc50689993f9223bde2",
							"darwin-arm64": "4976cae57c4b7e981896d6d11969a5da28c8b5b62d0672d52a2342705a36a478",
							"linux-amd64":  "cda4f52c475abec727a39bd438c5daeae35dc6f8e8a529ad33ca26c2d6fad88d",
							"linux-arm64":  "45090820ed3f3802d049313fd559985e5eb985b081aed2d4b3c08971414659eb",
						},
					}),
				},
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-ineffassign:ineffassign-asset:1.22.0",
						Checksums: map[string]string{
							"darwin-amd64": "2bebfe874d6b557754b9ea60f2701005f01a79584c6922f84a36ee1f460438c7",
							"darwin-arm64": "f0df8f6d269ca8d94bda85117959d270fb0a74fbef802ece62ab8e386608235f",
							"linux-amd64":  "b405c94ff5cc0a455eb089f858e42387880aefbd4ca7af910aaea672d34bdbe0",
							"linux-arm64":  "b666ff16f64508788f2e090127713200df5d96d75eb214354942993d8e3821ec",
						},
					}),
				},
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-outparamcheck:outparamcheck-asset:1.23.0",
						Checksums: map[string]string{
							"darwin-amd64": "7d26732ec6064471588a36f3d1040b1253cb583aff68481258bd7f4311c7b013",
							"darwin-arm64": "6700f339f2cab04ce25e640e20ea3714614f950868e6e8f0c13a8b8317411889",
							"linux-amd64":  "338a9167bc825f6c39a9361f148614c41420933db99d8276c813645ad165f60d",
							"linux-arm64":  "d94f7308d9bb2326b5e98686ce9c6362a741c2edae30cac65b1323858e252920",
						},
					}),
				},
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-unconvert:unconvert-asset:1.24.0",
						Checksums: map[string]string{
							"darwin-amd64": "215c9af2344201a4d0fb43f69e8f1287637fb10a23e6ec2d6e73d64c0bfb28e3",
							"darwin-arm64": "a833a76abd1b7b50c11c83aa68130f2743b2b90060104f2eda63c5e114a19a0e",
							"linux-amd64":  "582b8b65cb621048575384056dfdca8d5af34776f5a6774f8f997fd1f2cc464e",
							"linux-arm64":  "793613985fd9e0fbcfad3414887b637f959df3849809016fee47f78ee76f6803",
						},
					}),
				},
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-varcheck:varcheck-asset:1.24.0",
						Checksums: map[string]string{
							"darwin-amd64": "ee9054270e5e78f6e1fb9ccb4177809b3b202488b5441adc9c2807a0aa266014",
							"darwin-arm64": "9ae74befe6a787706778d9c9af362ba92593328dd2bf7264a6d8aba482975f75",
							"linux-amd64":  "03f3486e81a0adc8389d0e4f52bf184f3928344e453589f141dacf7a5363c012",
							"linux-arm64":  "392fc620c3b30479124a61c8b14915954d88ac4e8255c7bcb29473db5ede1e3d",
						},
					}),
				},
			}),
		},
		{
			LocatorWithResolverConfig: config.ToLocatorWithResolverConfig(config.LocatorWithResolverConfig{
				Locator: config.ToLocatorConfig(config.LocatorConfig{
					ID: "com.palantir.godel-license-plugin:license-plugin:1.21.0",
					Checksums: map[string]string{
						"darwin-amd64": "cfd24be3a44545f87de0cb205e5facc869f8601930da60e4d319bb6ababed76d",
						"darwin-arm64": "49aaf7aed3ef65433506841737ab9d41ce59667f37426ca9a98a33db6cd19774",
						"linux-amd64":  "b473fbae0d410a7dc93501280366c7d144d5d38dfe2a5e848fcec2d38e99792a",
						"linux-arm64":  "1ea467bd02bd8e0abc06b729087b1b8808520b1685ea692663321ef26e8b1898",
					},
				}),
			}),
		},
		{
			LocatorWithResolverConfig: config.ToLocatorWithResolverConfig(config.LocatorWithResolverConfig{
				Locator: config.ToLocatorConfig(config.LocatorConfig{
					ID: "com.palantir.godel-test-plugin:test-plugin:1.21.0",
					Checksums: map[string]string{
						"darwin-amd64": "e8f8e7b34242bedf975740a41e2906184ba747c89d779116b4f04da5a1b11313",
						"darwin-arm64": "9e593c6151c7f397f6801d02deb413e349ecdb01454a40641902834653767957",
						"linux-amd64":  "f04b99d71cfa68ac10ef3de88ddf04026d01da1ef29199e15a1010cecfe34ef8",
						"linux-arm64":  "f4928c65da86ebb805ac720d9a9378b12cf8de4e33062d14f66e332a06c8e77d",
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
