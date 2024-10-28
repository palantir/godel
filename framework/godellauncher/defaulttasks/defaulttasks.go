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
					ID: "com.palantir.distgo:dist-plugin:1.77.0",
					Checksums: map[string]string{
						"darwin-amd64": "86ac12f8b12f07a339a1bf6ac3209af331e8bfd1bb173dd0cfb96f4424ecea0d",
						"darwin-arm64": "a1a24549f9b23571288e2732903db65b50939d41c7e0396f3bb023f4c5874146",
						"linux-amd64":  "7a750ef417d639dbfd795fd636cf567166d78004badc43e39d129dd698fe9de9",
						"linux-arm64":  "f1d5f98dda6d28077b257e6dd658980ba70a2a285972b007464267e78310ea8f",
					},
				}),
			}),
		},
		{
			LocatorWithResolverConfig: config.ToLocatorWithResolverConfig(config.LocatorWithResolverConfig{
				Locator: config.ToLocatorConfig(config.LocatorConfig{
					ID: "com.palantir.godel-format-plugin:format-plugin:1.48.0",
					Checksums: map[string]string{
						"darwin-amd64": "9fc2cf05fa6faf56b7395df381058103e8f98496cd64c81940fb8322114f3258",
						"darwin-arm64": "b503c88d5fdea2ca8efcafbfd4a65a910394b094e11d9c9d5235218ee3a98d86",
						"linux-amd64":  "f903bda68be17b29517f9b1aaeffec55527923309d69c782f882456494b8e7cb",
						"linux-arm64":  "a1f4ddae43c50e59c215a31c24bda9afef5c3b0be6fb81b5d79c60b9aa718cee",
					},
				}),
			}),
			Assets: config.ToLocatorWithResolverConfigs([]config.LocatorWithResolverConfig{
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-format-asset-ptimports:ptimports-asset:1.47.0",
						Checksums: map[string]string{
							"darwin-amd64": "90d027bf67e7a84ba3704076254f365f69f47b6073d44b2223b73dc3f0bd4a4b",
							"darwin-arm64": "61617acf7b6ef5db673a946803d78b3584610c66e2b44d2bcd1c59108c9ef1db",
							"linux-amd64":  "7e83ca228ea055d48428c54389fbee7bd9b3c3d75bdbac9a0da9dee4b111388a",
							"linux-arm64":  "5de739b1b901fdfbd706fdb620fe5564acd53067eff734012165833b0d698ef3",
						},
					}),
				},
			}),
		},
		{
			LocatorWithResolverConfig: config.ToLocatorWithResolverConfig(config.LocatorWithResolverConfig{
				Locator: config.ToLocatorConfig(config.LocatorConfig{
					ID: "com.palantir.godel-goland-plugin:goland-plugin:1.44.0",
					Checksums: map[string]string{
						"darwin-amd64": "80b87316f83d9f8ae1f22daef451832edec69757b90bd0907e97a811bb46d0a1",
						"darwin-arm64": "5900332c767b512c89e3eb5b03c2f854ac8ba50041561e3dee666328b80032f8",
						"linux-amd64":  "139ad36ab3abcac01f988acb0b3e09456e305eb178d829fec7b209ec3f5683ff",
						"linux-arm64":  "a6ee226cf69dd669cfc60cd74e16c8312858489592451977306d176ecd8920e8",
					},
				}),
			}),
		},
		{
			LocatorWithResolverConfig: config.ToLocatorWithResolverConfig(config.LocatorWithResolverConfig{
				Locator: config.ToLocatorConfig(config.LocatorConfig{
					ID: "com.palantir.okgo:check-plugin:1.61.0",
					Checksums: map[string]string{
						"darwin-amd64": "780ec7bdc240d0176895011fb32aa64a964c7c92d47e315dc0a0f17930d85b51",
						"darwin-arm64": "c9314361dc09684bbb5d74e174821256e7f3524a840ee7e60e3b673c3b449d49",
						"linux-amd64":  "91a8e4f00dadbf8016b38f93875735e3d0765594db7c4f640b4ed82b3ae867dd",
						"linux-arm64":  "62697b7292de224ceb596c1044295fb2e32de34d74b1ca87f10806d1596a3c6b",
					},
				}),
			}),
			Assets: config.ToLocatorWithResolverConfigs([]config.LocatorWithResolverConfig{
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-compiles:compiles-asset:1.53.0",
						Checksums: map[string]string{
							"darwin-amd64": "40fce2ba1199c2ed65d9f9b67c23fcec069f57ebbdb8776b430b65ea76a8bc8d",
							"darwin-arm64": "9895a8512d6d43fd45add72fbe07c5fd615699cd0f9c12852fd42e6647166bab",
							"linux-amd64":  "d5eb65089b6abb6d13a94093bc557e03d1b9360edbb5cc312026cc00b8494d4b",
							"linux-arm64":  "6ef9decc9d98076f85bb9827044e8803773550bd9c806bc801f047c3ca924155",
						},
					}),
				},
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-deadcode:deadcode-asset:1.47.0",
						Checksums: map[string]string{
							"darwin-amd64": "fe9704d35abe10f1157970c928c4204e3e35320837c0f7515d7bd409c051710b",
							"darwin-arm64": "d6eea1aad8697f57dea36a1d6450999f909404a2c16c9716c5783915991dae65",
							"linux-amd64":  "c095abf5c7ab052b7a1a636d9f1cf9099c18a3e6a9f5f8ab274588545c31577d",
							"linux-arm64":  "76996a063adc00e9f1f012b612727fd9641c0d4907008208ef845c18756e6b10",
						},
					}),
				},
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-errcheck:errcheck-asset:1.49.0",
						Checksums: map[string]string{
							"darwin-amd64": "8a95a2bf3062ab69a7022b6f9a7d8b43f54974b767e5bad36eeb051465d95615",
							"darwin-arm64": "ae82d45cb91fc8f6041fd34913a67555ecf712bd92ff4564705c1a4460e19533",
							"linux-amd64":  "41a53515a531ffc533d3208021d427769a0dcc5390e510e5d04bbbf5fadcaf4d",
							"linux-arm64":  "5823057b181870f3eaf0e26fc55f30187ae669133148ec1306b6efa8f117ba43",
						},
					}),
				},
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-golint:golint-asset:1.40.0",
						Checksums: map[string]string{
							"darwin-amd64": "0613f98b068a319ccfbc83b6d2e4b64bac2e584dd701cf43191cda17a5330f39",
							"darwin-arm64": "8a9e5a9598df3fec4b8e6e7e5c15c9e712951705371ca85fd3d56a0aad258c1d",
							"linux-amd64":  "875226e2a9c906f8a493ac3aa526f3acd2c743f49a3bfbdcc47e663717b7a951",
							"linux-arm64":  "a65cd31b39b1168946d5ec466c16862e3255b647c686f2e780053324d2d9423f",
						},
					}),
				},
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-govet:govet-asset:1.45.0",
						Checksums: map[string]string{
							"darwin-amd64": "f7c1cd91c7ae1b68855cbd4b4a32c00dc6de5637518ee2964a4582afab6c1ce9",
							"darwin-arm64": "853e082633385ee081534d287032d91f51e41efce84cfd1dcb50dd2ffc0839ce",
							"linux-amd64":  "e88e1f596588dff02399d508450966e7f428f2731fe960d8937fd6d02fd5fb60",
							"linux-arm64":  "a41d77386677d206a341cd3fe7c08f6aa2399fbe02a2def4254f00b6e099844b",
						},
					}),
				},
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-importalias:importalias-asset:1.44.0",
						Checksums: map[string]string{
							"darwin-amd64": "e43472367103fc6df1af98630d1f11fb1866154d463f925f82d5f222cb3b4bd8",
							"darwin-arm64": "1d6641ebc38a0b25ad31ff2b8689fb64571fd97d6c25a39908ef3c84754053bb",
							"linux-amd64":  "5edbe4a0b640af14db717c68228499b60df8af0be3c04b99fb619a1f393ddae4",
							"linux-arm64":  "20db1615b9fd0f57a5ff0d368ba9842eb69ecd7f1aa201c4985d536ac81394c4",
						},
					}),
				},
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-ineffassign:ineffassign-asset:1.45.0",
						Checksums: map[string]string{
							"darwin-amd64": "8552ce61d63357dbea1c74422153a8fea076608740a0a3d549172b939d4f6bef",
							"darwin-arm64": "2453122947c66d3b3b7a70c4ceafbb0b888f3e1cc864fa9a789288d12e41f06c",
							"linux-amd64":  "24b9d11cec00c32a262a08ced0129c130915ea3ead74f100e2e41f53dfe3f74a",
							"linux-arm64":  "66109c6b18efb567a53cc1c81e9e1a73c2f6e98105939f08c8abe275eda297d2",
						},
					}),
				},
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-outparamcheck:outparamcheck-asset:1.46.0",
						Checksums: map[string]string{
							"darwin-amd64": "1d7029dac07ee74c1dc1b893bee9ca39b2b0424df717da21327676253351bc72",
							"darwin-arm64": "fd869946e5f5b3b7375d418c6a3ce79d899dcfafae51118040329f2b5501f212",
							"linux-amd64":  "1c7f88d8554233c85ae74219150abdc885230f1185921c2a112bcc4f1c772ffd",
							"linux-arm64":  "8abd5b0415e6075eaa0e6e74682eb28f0247c2987c0fdae43a4e0950be098156",
						},
					}),
				},
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-unconvert:unconvert-asset:1.47.0",
						Checksums: map[string]string{
							"darwin-amd64": "7de6459fc2202e870c83a582ee9b9219c95927f86b1c268e2342cdaf1f41a8e3",
							"darwin-arm64": "4ea35e55588aebba8cb62a2fc4da31c2deae3edcf88c82a3612a60e2cbb96b4d",
							"linux-amd64":  "487be2026b4e272959d083e7edff24f7e6d55eb86f2b300ece0573409e7cdf37",
							"linux-arm64":  "0d949df4d369eb475ac76d6dfd172ae8988965414131bf858b6e359978f50c1d",
						},
					}),
				},
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-varcheck:varcheck-asset:1.47.0",
						Checksums: map[string]string{
							"darwin-amd64": "f73a9bb75656260a32847435771d09e3ac48302698f768b5c4cfe0b509270f4c",
							"darwin-arm64": "4261a7dd22a10d37fb03ca1b4714334ec64a610004b0c0222969b66eecd848e7",
							"linux-amd64":  "0624edf56d167f3706e615d5c964b095c1675f83c7c8f39bc2ba352fc7af6360",
							"linux-arm64":  "ae7493d83ef79eeae7498b72fb14d774261078bb6a44de8d0bd78e02270ef9c1",
						},
					}),
				},
			}),
		},
		{
			LocatorWithResolverConfig: config.ToLocatorWithResolverConfig(config.LocatorWithResolverConfig{
				Locator: config.ToLocatorConfig(config.LocatorConfig{
					ID: "com.palantir.godel-license-plugin:license-plugin:1.47.0",
					Checksums: map[string]string{
						"darwin-amd64": "788dcdead84c7e46b8a29ae2ad490b6e7d13e29cb0e9bb34c3d67836bf365379",
						"darwin-arm64": "fb85b3cea5392e8f6d2c1b25357773a80ae927bc619b5e303cff8bf9478a2066",
						"linux-amd64":  "090f206f78eeef072a55e85649125562ef875d155891a79a81c4911d9b2b6dd9",
						"linux-arm64":  "a7d9217a79c94e0caa28852f00d1f2037dec6bd5aa92c85997d73f774156a6c3",
					},
				}),
			}),
		},
		{
			LocatorWithResolverConfig: config.ToLocatorWithResolverConfig(config.LocatorWithResolverConfig{
				Locator: config.ToLocatorConfig(config.LocatorConfig{
					ID: "com.palantir.godel-test-plugin:test-plugin:1.45.0",
					Checksums: map[string]string{
						"darwin-amd64": "b8a81a75306440a33a2e8b45b3b4172bd67a701601227da5a53abe2d8bfdf276",
						"darwin-arm64": "8f06c34f7d1b3c09a090243c9131815e3b1a31dfd0339076a59c95062dc4ded2",
						"linux-amd64":  "18a973aab52d1374caccf92b36a5593b7b8456f87da5084e47edf56afbf198c9",
						"linux-arm64":  "28d9ba5e59e1d3935ddd300e2524946793164b26ad5ec4c17ec7f3f7f04ae77a",
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
