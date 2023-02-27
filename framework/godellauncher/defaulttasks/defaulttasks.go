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
					ID: "com.palantir.distgo:dist-plugin:1.51.0",
					Checksums: map[string]string{
						"darwin-amd64": "166a60a9bdb6584f20b715afee9cee4a8b9534c36c53264f0edd15e06dc3d17b",
						"darwin-arm64": "73fc820d6c2babe69047da4f0080a23422559c18338a2a134b123c2186793a33",
						"linux-amd64":  "a4528b89029fef91399fbd861c2c995bbd53c3491816590bc8b29ff765fc970b",
						"linux-arm64":  "f41cd0b7a57d7b7e948623f672e4b58300deab00d993660787647d113db73856",
					},
				}),
			}),
		},
		{
			LocatorWithResolverConfig: config.ToLocatorWithResolverConfig(config.LocatorWithResolverConfig{
				Locator: config.ToLocatorConfig(config.LocatorConfig{
					ID: "com.palantir.godel-format-plugin:format-plugin:1.28.0",
					Checksums: map[string]string{
						"darwin-amd64": "ed05ef0516cebfe73556601a648a504fed5edeb239794797731c6f36125ec944",
						"darwin-arm64": "0c88a85ece1b4e2d7295f8b5ac2ced0d342ad39a883f4f5ff5af63dee1e3e4ce",
						"linux-amd64":  "16019a59c63c5394ff512cd4dd6e778993b54a7dc36bc7c40f0d13d65c9280f7",
						"linux-arm64":  "6db8a041bd1569104c9a60543e6579bd26dd5d053a3ad1a9c7b654fd7c6aa834",
					},
				}),
			}),
			Assets: config.ToLocatorWithResolverConfigs([]config.LocatorWithResolverConfig{
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-format-asset-ptimports:ptimports-asset:1.27.0",
						Checksums: map[string]string{
							"darwin-amd64": "ee7fd82780f1f84bce77698a40863a8e85be247ff5bb59c99d69d62322e3cc93",
							"darwin-arm64": "2910db17283e7a64c3ce7b0934a7ad1b925abf359511c5fbbc1c62b53caa508c",
							"linux-amd64":  "b227aabf22862cfe3f74031c7e125f039332a1bfa236c868e440d77a5e17a0a3",
							"linux-arm64":  "4cad7308e139cacb03fc9e45b3abd597af976bb21a563e67448d881ab61c52b0",
						},
					}),
				},
			}),
		},
		{
			LocatorWithResolverConfig: config.ToLocatorWithResolverConfig(config.LocatorWithResolverConfig{
				Locator: config.ToLocatorConfig(config.LocatorConfig{
					ID: "com.palantir.godel-goland-plugin:goland-plugin:1.24.0",
					Checksums: map[string]string{
						"darwin-amd64": "b2b9bc71bfd76e7e944d0ff0003dad9953758485793cc13d3a6cd2c978eef85b",
						"darwin-arm64": "39e098cf0cfa465b59ec3544fb8de19bd823bbe8ad38926d1c11ac820179e94c",
						"linux-amd64":  "dfae364d77a75818352404f4dd3802f887a0f5224bfb84293716d37d37e3e01b",
						"linux-arm64":  "a5cd855d3517c4c25634d6e14ecd8bacdca6119ba8a4772c9268c6f50245459a",
					},
				}),
			}),
		},
		{
			LocatorWithResolverConfig: config.ToLocatorWithResolverConfig(config.LocatorWithResolverConfig{
				Locator: config.ToLocatorConfig(config.LocatorConfig{
					ID: "com.palantir.okgo:check-plugin:1.33.0",
					Checksums: map[string]string{
						"darwin-amd64": "19ba39a731a61e57f3ce2df9246106953792fd58c123f8a94733b89be2785a94",
						"darwin-arm64": "767fe00e680c16514aca4069eda8d8a95c9fd9d5e25d3c413234d7c9deb03f32",
						"linux-amd64":  "b46f7aba06bab5c6a948660b79d84dd19e0b4a95a9fcfc62751dffafa880dc6c",
						"linux-arm64":  "dfe255559791c21efca7aff0e8d3e4cbf282d445d28f8ec72d9035aa2f649ec7",
					},
				}),
			}),
			Assets: config.ToLocatorWithResolverConfigs([]config.LocatorWithResolverConfig{
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-compiles:compiles-asset:1.32.0",
						Checksums: map[string]string{
							"darwin-amd64": "32e0ef56f8b2a56a4f66515a122baf60b07041eb26e343f565a8c177406da7e1",
							"darwin-arm64": "47393853d0374dec54593078e8f3bb3c165199d65aae20661aa41acf69766514",
							"linux-amd64":  "71742f35127d4ccede01e4f778a62ee1e8a6f0dcbaa036b7f0492df73fa4f383",
							"linux-arm64":  "1a715a34eb7d28bfc3b18ab05fedcf66689a93076669ad4af1a3122392f3f5ff",
						},
					}),
				},
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-deadcode:deadcode-asset:1.29.0",
						Checksums: map[string]string{
							"darwin-amd64": "96bc5515f368764aeade4fcd6dc79b475e2ce6f17f4e4008d6566b959247dfcd",
							"darwin-arm64": "92d50093e1191281a98a7fa506d3260c124486660c67b0c774386596afe99133",
							"linux-amd64":  "a5662ab2ffd2d5062f6966cd9085aa6a4ed613f9716a946b5d88410398fe3a7d",
							"linux-arm64":  "d63e019e4240220e5e6d692454eecc0898a7bcb5e4319f5235a2c257b7a631ed",
						},
					}),
				},
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-errcheck:errcheck-asset:1.30.0",
						Checksums: map[string]string{
							"darwin-amd64": "0560c23aeff52d3ab25a79b7e8c585fe22249b38beb46e746ba00348e16fd867",
							"darwin-arm64": "733c9fa0656e40780ebde69d83c4254e81f96d31ed9ebbf3e2b647e459d39784",
							"linux-amd64":  "5aff5f871ae47c5b3df1a0bc5c8aecaed4c49ca92bbb836a6b106f05d0c48562",
							"linux-arm64":  "f3acf0a9c5113f24059bf8531751f73acae11fce69c9cb2cd7178740b490c71b",
						},
					}),
				},
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-golint:golint-asset:1.21.0",
						Checksums: map[string]string{
							"darwin-amd64": "213b0e45b56a6962de85b8f7da7628fd4c7367186ec3a5e210c0f0e5af3aaf24",
							"darwin-arm64": "5731ab7f3e6f395c8b53935c5d20ac3497430e6ac8d4f7a60a6f5f3687de5d45",
							"linux-amd64":  "d6e7d8e587e53bf5471b124c9dc187ef9c4282bddc250e0b3232eccd8423b7cd",
							"linux-arm64":  "25e503095eadc849f94b219b596a7b672962234c16453847e13951af14003062",
						},
					}),
				},
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-govet:govet-asset:1.25.0",
						Checksums: map[string]string{
							"darwin-amd64": "ea3210f3162f4aa56cc84d63b6f0798edd137cf5672d48aeb8dfc9388bf4760c",
							"darwin-arm64": "e80aa4c40312616662b23489d3b601024ceda70a198bdd0c4a3357f8cb7d9ea8",
							"linux-amd64":  "9d70270d734041ad10732a78da4d0c2e515ceff5926b954b2f07fe068b556260",
							"linux-arm64":  "ad959a1f287403410691f21df7a8d5bc75480e319b033164b5ef3a52cc27ebb9",
						},
					}),
				},
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-importalias:importalias-asset:1.25.0",
						Checksums: map[string]string{
							"darwin-amd64": "d05d30e70c6e79dac641312cbd88e7d82771a7cb6ee1914e67c5bf9f805eb02a",
							"darwin-arm64": "d04ea79416aa9a0f601500e2505a29cc597ccae42ff1cdbbeb613ca6c7d334f5",
							"linux-amd64":  "4a5b16b1871246a08b19b13f1ef5866be2474090fe2038af2757f7e552dfacc2",
							"linux-arm64":  "f01e4c921fe5480aee16acbe2a1d31a90ec827a6a57ae90293f5e47c2d09505c",
						},
					}),
				},
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-ineffassign:ineffassign-asset:1.27.0",
						Checksums: map[string]string{
							"darwin-amd64": "bf8ed28c0e49766e142cd3c04ada9a1221a889a056339c769073dd7177e21d01",
							"darwin-arm64": "b0809168edcf90741e0b5d94a5dcf52be9cbf3816f7893a18259f9e4d8a36241",
							"linux-amd64":  "13b4819f9652210f1c38db175be22bea87b53618fcb5ab45a96300daa9a1e59e",
							"linux-arm64":  "28f9ee1026ee2f9e410c8829c427b9ae03f244aad6ea7d22e051b894b8555cb4",
						},
					}),
				},
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-outparamcheck:outparamcheck-asset:1.29.0",
						Checksums: map[string]string{
							"darwin-amd64": "613a53cb7a159ffea6d78725b52781bee94bdfebe27341635795159ecaad84a5",
							"darwin-arm64": "8e9e4aa1ec1565eb40f31ae3357e301823ef74e6caba760cb46e61bb3454e15e",
							"linux-amd64":  "4be9b3dbd51e73e721f1f8b864fd04413c137635939f5da1d11dd38cb88f78d8",
							"linux-arm64":  "f52e7d0cedd5d0a7ed16d6c19d045034d245d76b638d2cf87fb4eec04321e640",
						},
					}),
				},
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-unconvert:unconvert-asset:1.29.0",
						Checksums: map[string]string{
							"darwin-amd64": "a24ae452b1fd4a9d29e14834ce7506e476945a217c8474a9fdc05678c07b57bf",
							"darwin-arm64": "f6bb41ac380d3cbcb2fd65f33a6e287f61a77ebb5bc11d22ef1ec9db75795fdd",
							"linux-amd64":  "439566f00268eaf34f5a73155ff5127f411bab547af86bca881c41b373bd34d0",
							"linux-arm64":  "0b2a4c387c00cd651424138ab871848332a6ce8dd4181e978123a54e0efd6306",
						},
					}),
				},
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-varcheck:varcheck-asset:1.29.0",
						Checksums: map[string]string{
							"darwin-amd64": "6261476e0a4a0d684258272a4dd9c86a993bea0c713df7e090046c5ee43fac82",
							"darwin-arm64": "15f70e85c3ca777573e49e6ce41496964185dcce06b5a65cf058b67779934ed6",
							"linux-amd64":  "6247c13f7b2a333db2992470870ba751ff527f1552345666e83ff1dc8f0f9a1b",
							"linux-arm64":  "646ce35c96e578ec709baf1647feae31b9e100b80c54d249014bf99b71728987",
						},
					}),
				},
			}),
		},
		{
			LocatorWithResolverConfig: config.ToLocatorWithResolverConfig(config.LocatorWithResolverConfig{
				Locator: config.ToLocatorConfig(config.LocatorConfig{
					ID: "com.palantir.godel-license-plugin:license-plugin:1.27.0",
					Checksums: map[string]string{
						"darwin-amd64": "78be53000ce9d7d3d3990d4de91d52d2872d86fb60d58474c5cb962335cbf91c",
						"darwin-arm64": "60ed3e502c5179d3c692543d4f50386f1cd2956151a320db833c4a70e7a7c667",
						"linux-amd64":  "62a708dfd65211e0535fa2a76717cfa63cfa63cece8d537906c135cded95a375",
						"linux-arm64":  "19c73f9da67a1274f94ef4bf625d1767ed6aed94726749026dc562a7f1189240",
					},
				}),
			}),
		},
		{
			LocatorWithResolverConfig: config.ToLocatorWithResolverConfig(config.LocatorWithResolverConfig{
				Locator: config.ToLocatorConfig(config.LocatorConfig{
					ID: "com.palantir.godel-test-plugin:test-plugin:1.26.0",
					Checksums: map[string]string{
						"darwin-amd64": "369d9d30ca2356c92ff33e2d8b27ebd457ca179c5861750a6289491bb9ec2712",
						"darwin-arm64": "f9e3186456559bccdcfcec92c550fafee100dfbf2bbc6eb65e338a28176456d1",
						"linux-amd64":  "48a0d90e6ee3c9502de33b94b6569d1491e70925698dcaddce511ae2cce0d9ee",
						"linux-arm64":  "e132b3de24de2671741db628242a84b0c55f63d3e13d680c0455a77f8422fa01",
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
