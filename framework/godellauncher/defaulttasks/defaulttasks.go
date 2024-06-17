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
					ID: "com.palantir.distgo:dist-plugin:1.72.0",
					Checksums: map[string]string{
						"darwin-amd64": "17e411fd5a531adc7047f714ddd09c21a0422e9f467744f73ad5923af0981a8e",
						"darwin-arm64": "5e5f8ab2c5163ebdc088a859e57ae07d1a505bb9c9dcba0865007aa1249593ed",
						"linux-amd64":  "e06272934464b37c59231115d987b3dc0832bfeb643f8ad49b9ab0a8e0dae577",
						"linux-arm64":  "a29c7d241d822a571bc519b5c355a5c29a76a9b0bfcf11b9f3b224a3d1b0023f",
					},
				}),
			}),
		},
		{
			LocatorWithResolverConfig: config.ToLocatorWithResolverConfig(config.LocatorWithResolverConfig{
				Locator: config.ToLocatorConfig(config.LocatorConfig{
					ID: "com.palantir.godel-format-plugin:format-plugin:1.43.0",
					Checksums: map[string]string{
						"darwin-amd64": "cd6fb6072644fb13df9526ebe5270d1258b5b4b6ea8ece6b3db79bf1f3ec5d59",
						"darwin-arm64": "32502e22cac72729a00ef1ead23cbd21d19cf5d0cda857c9ccdf68150ae5211e",
						"linux-amd64":  "d6e499b75d82eb50033377928e5013ccd20df673b48a13213a5ab248602812ff",
						"linux-arm64":  "77b486e78dcde9974318ae5843b086c235e7522a61207a43c8b1a8ce8dce317c",
					},
				}),
			}),
			Assets: config.ToLocatorWithResolverConfigs([]config.LocatorWithResolverConfig{
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-format-asset-ptimports:ptimports-asset:1.42.0",
						Checksums: map[string]string{
							"darwin-amd64": "aa9fa7dbd5020dfaaea0a50323e3093dfacf692ed44aa56ba6753c43f5894d59",
							"darwin-arm64": "4e501945aa56c9ef070899c0cb57ceb488ea2a9e84fee42cd6128641e576f154",
							"linux-amd64":  "1a380e6ca046cb4325cc8ed2e200df01bae3193a662635f00e1a5141d27a1c6d",
							"linux-arm64":  "235c21a088e890e71c7e17747bfaad25626888430826a21d797961867244aa96",
						},
					}),
				},
			}),
		},
		{
			LocatorWithResolverConfig: config.ToLocatorWithResolverConfig(config.LocatorWithResolverConfig{
				Locator: config.ToLocatorConfig(config.LocatorConfig{
					ID: "com.palantir.godel-goland-plugin:goland-plugin:1.39.0",
					Checksums: map[string]string{
						"darwin-amd64": "f0519faf444385ea41366cd37f4f55311544d3a238f5f6f24ce903136d36334c",
						"darwin-arm64": "843f443ce7c90ed98d9ce6842ed396e9a2f6d626237e5ff9ddaebff18971ff25",
						"linux-amd64":  "65328df241c072f41f48c9b47e830a7fd1a267861fc7b8a5e36121a242a2f216",
						"linux-arm64":  "3542cc4d5804a456cd1a3ddaec19157e6e34032291d6b27fd7ad7feb741b3ac6",
					},
				}),
			}),
		},
		{
			LocatorWithResolverConfig: config.ToLocatorWithResolverConfig(config.LocatorWithResolverConfig{
				Locator: config.ToLocatorConfig(config.LocatorConfig{
					ID: "com.palantir.okgo:check-plugin:1.55.0",
					Checksums: map[string]string{
						"darwin-amd64": "e19d78628839691cb3f15879f61c55236f41cea4ab8af03aedfce4dd817aa159",
						"darwin-arm64": "230a5db3b3f2477eff3b6a3a530be8b0c7d57f6de47e944407f46d21d6212fe3",
						"linux-amd64":  "859a75f462df56fe981dd38461bd120473bf534859d04ee0eb55a2abb2919dab",
						"linux-arm64":  "69dd7344729d82082646399fd92fffc7b3ae10e87cd055e3edcf98af7916ea78",
					},
				}),
			}),
			Assets: config.ToLocatorWithResolverConfigs([]config.LocatorWithResolverConfig{
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-compiles:compiles-asset:1.48.0",
						Checksums: map[string]string{
							"darwin-amd64": "f04cdaf4018e89ad0dec55441d0a21f384c8a71c28e59c24bd586d3eeaaf6a91",
							"darwin-arm64": "ece907c67d5cb38d7fadb96a3a0ddea4122b9e2cc60d37025b410bf4d8d3e6de",
							"linux-amd64":  "1bc96bc770a4edf807751e1c200a10993b5133cd3fae112a58c6279ec8d5f2a1",
							"linux-arm64":  "2ac5b1f6cfadcb3586db2eb3212544b22f3494b10de845b1a016d6ef729f1657",
						},
					}),
				},
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-deadcode:deadcode-asset:1.44.0",
						Checksums: map[string]string{
							"darwin-amd64": "50a6b9efccddd903cd0fd2ce5a5fd5ffd8f2ad97c1eecdd6c2c114b5f1d9e673",
							"darwin-arm64": "e922001004e5994ae4c811a23461dda014d2b3197e31787b9081c274852edb34",
							"linux-amd64":  "5cbab8995703205a451c78499e0751e7de608e048e343884580f01176d180501",
							"linux-arm64":  "0291d4f0de9c5ed8569446f33247ed37fae347b2028898c2a9f19022d50c87a7",
						},
					}),
				},
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-errcheck:errcheck-asset:1.46.0",
						Checksums: map[string]string{
							"darwin-amd64": "d852baa3801024b571862be5f4da311ced4b78c50618a7a08de151ade071cea8",
							"darwin-arm64": "4c4d8e23d2493a1c0108412ecd6c95a9a28cb0589e8106acf0ac52511f6407c7",
							"linux-amd64":  "ca465fec823fe9b917292b7b7496cf726e00a9d08ce57ac271df80e0e8a6c32c",
							"linux-arm64":  "b51593f664ab212d916141dda736a9cfdbc063fbe992f396fe19cd4a7ffa86cb",
						},
					}),
				},
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-golint:golint-asset:1.36.0",
						Checksums: map[string]string{
							"darwin-amd64": "4b5c75d393eac156568bbeefa5ac40254ed43d2f4f147c3f00fa6620c392c776",
							"darwin-arm64": "df0060796a977ba65aa5a40e49b9e598be18f02f73e5ce54fd8e0d8093e31195",
							"linux-amd64":  "ca07c808072916ad71fc82bfc6d2f439488b592b8ca2097e007b0d9108e00d6b",
							"linux-arm64":  "1a540079e6194a1a53f4cb37b128ac46fb7437a3f74105c3fc9fb4e10677a205",
						},
					}),
				},
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-govet:govet-asset:1.40.0",
						Checksums: map[string]string{
							"darwin-amd64": "db417d5eb5af7cce228a72745ab120d839f2b3e4eca8158762b36e4da6ed565d",
							"darwin-arm64": "85a6e962c3373397b6a4ec2f86941753d95a28b642035422c8bca8587268458c",
							"linux-amd64":  "a42c031fd45b29cbde7b73bea28e61dbcc86e7aac57f373422659e0e20a854ea",
							"linux-arm64":  "fdfc2b48e2c7dd8009d5fc66379d4d43b93908d5c2b129d2f4871fa0786f289d",
						},
					}),
				},
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-importalias:importalias-asset:1.39.0",
						Checksums: map[string]string{
							"darwin-amd64": "62f83447d5960e41371b368463c28eeaa1d1ce99f526a5ab0fc8f14d2288dd2a",
							"darwin-arm64": "9dbe43492be8b66e81be9b907be8bb8fc43645678321e45e5b9069c0aae7c8c5",
							"linux-amd64":  "226a7c1d3371d553fb222aa810b10e35652d898a9d1e5bed25a7693b0359b608",
							"linux-arm64":  "b15939e9048ab525aa4c42a5487637c1c04d04e07f0b9d4e5132338c837c2ba0",
						},
					}),
				},
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-ineffassign:ineffassign-asset:1.42.0",
						Checksums: map[string]string{
							"darwin-amd64": "8dbbd625f12b3b589d1bd8279c2d04cc66ecc7b52652a295d062ff68bd019ca0",
							"darwin-arm64": "41bcd63ea66dc93be7bc20102caa91ce1255e4d472462c038865816ebe5e0d54",
							"linux-amd64":  "61bb90a53a11baf786297f50b70d5b1f64388013f9c54568f93234151188a0d0",
							"linux-arm64":  "f0a99ade6961a94b2c8dcc49b6fe069015920aee4cf9db3b6c739b1ad26d4950",
						},
					}),
				},
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-outparamcheck:outparamcheck-asset:1.44.0",
						Checksums: map[string]string{
							"darwin-amd64": "7411464066d4926e2357e69d33d1aea8bff70bd471a757b2d06fb94814898ab6",
							"darwin-arm64": "62273de861d1c1b47d3144e63a3197a3e735f7d14d5f3af42e4e866d98f4a26b",
							"linux-amd64":  "45c154199163842e2f0eb62961b22be0161a29f89e803f9331595736cefc7943",
							"linux-arm64":  "55c4b4a196051a8655cd65ec1773f0f03747bad451ad3854644d6a3c69d11c58",
						},
					}),
				},
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-unconvert:unconvert-asset:1.44.0",
						Checksums: map[string]string{
							"darwin-amd64": "eddf44f95ad102199463eec1955690047fc374d865b6202d6756d7209af57315",
							"darwin-arm64": "476795fbf7354a99453b012302f08d9f3e5e2aaf4fe02f1be366b4ed25f72250",
							"linux-amd64":  "56a16bb8bb5573c7df0202731bb3302efa01bf7cb29c6bfb0d2b9270755e4bca",
							"linux-arm64":  "00f2bb745a3e344054003d0a0001465dfa1a7cf000f4cd6483160899c780cb56",
						},
					}),
				},
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-varcheck:varcheck-asset:1.44.0",
						Checksums: map[string]string{
							"darwin-amd64": "4016bba7b72598e46e898ab17e8a1a891b760f229b494c3e89a13b6c932c931e",
							"darwin-arm64": "16dcb961026dcd6ec0a96980db3797bd0f260e64ffd89832c400df06394a4313",
							"linux-amd64":  "53a832512b31d58408c6a06378e17a10c51382d53dca585c7c139637b8f41542",
							"linux-arm64":  "e977325bc418ce9b48da189950ea425d0e14528c7e99de071e81621c7798089a",
						},
					}),
				},
			}),
		},
		{
			LocatorWithResolverConfig: config.ToLocatorWithResolverConfig(config.LocatorWithResolverConfig{
				Locator: config.ToLocatorConfig(config.LocatorConfig{
					ID: "com.palantir.godel-license-plugin:license-plugin:1.42.0",
					Checksums: map[string]string{
						"darwin-amd64": "1b45b215ffbf4c202ee0feca039fff1fe2faf99151d393c30906d0c7e4356f51",
						"darwin-arm64": "2a89fc240ca55f2c9d37ceb3b1fe38ef0d9ed4dcc59838fdd02de3ce2cba2cd0",
						"linux-amd64":  "7abac4848e74d92fa341e971a3242d25bab15eea088ffc784f3490aa45871bb6",
						"linux-arm64":  "330c327913b937f404c49ec81f7f60e2b2fa54d40620170332c337ff11f84cd0",
					},
				}),
			}),
		},
		{
			LocatorWithResolverConfig: config.ToLocatorWithResolverConfig(config.LocatorWithResolverConfig{
				Locator: config.ToLocatorConfig(config.LocatorConfig{
					ID: "com.palantir.godel-test-plugin:test-plugin:1.41.0",
					Checksums: map[string]string{
						"darwin-amd64": "a8d8ed537fb632bc80c1c13e16e240e41461a32bb50a1902cfc06df1b10d019e",
						"darwin-arm64": "e93a58a510bd16f19bc223e242334fc0ded842f6ce813b209b028962539e7a8c",
						"linux-amd64":  "461ba23c82cb7a98392b20e87fdca8efbdbb4af0178a2f71a0f485a9bbafc647",
						"linux-arm64":  "56b1433a9529791834ccb736550f9b70e0928d72774c5a9c0848c7c92a333025",
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
