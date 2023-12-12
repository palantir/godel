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
					ID: "com.palantir.distgo:dist-plugin:1.63.0",
					Checksums: map[string]string{
						"darwin-amd64": "10f0802ab2cfe7d2a44895bc1ae75a8b215fd7e7bbff8795dd4ddec408c5e11f",
						"darwin-arm64": "ea6b15994e2d26d988eda7ad270e677e34b4fe3b04e49138d9f7d77068387121",
						"linux-amd64":  "c43d28370ba090283adc3bba62acffb96b17232dc3aea970ea74a22b7b509bb3",
						"linux-arm64":  "07168d6cad6bdd19a64f28a521effe48488d2a70adb11a72a1cdbe11734fe87a",
					},
				}),
			}),
		},
		{
			LocatorWithResolverConfig: config.ToLocatorWithResolverConfig(config.LocatorWithResolverConfig{
				Locator: config.ToLocatorConfig(config.LocatorConfig{
					ID: "com.palantir.godel-format-plugin:format-plugin:1.37.0",
					Checksums: map[string]string{
						"darwin-amd64": "204408d8f28157c4f8c93b13d231c49bf4bf47104b09cb982042c0301be3e21a",
						"darwin-arm64": "a3c4305dc39a7435155bf64367499c31fd6e0a8deea418f7c6a2d1ed5ff3841c",
						"linux-amd64":  "7d58463864ad55e9bac03dacb4789170923d659ba1b1106dbf1bcc7773121c56",
						"linux-arm64":  "a77041223851e454738664e91c81ff74ee217bd3de8991f30bfe7c6961ab112c",
					},
				}),
			}),
			Assets: config.ToLocatorWithResolverConfigs([]config.LocatorWithResolverConfig{
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-format-asset-ptimports:ptimports-asset:1.36.0",
						Checksums: map[string]string{
							"darwin-amd64": "41188aa8a53397e0f967ca756aa44d657bc061056575dfc6d70aa0ea12953b14",
							"darwin-arm64": "7dc6e6567f2719b75a463568509c5cf2cbabc0c778a5235a0769375270329df7",
							"linux-amd64":  "367451ea2dbb0640b32ab7f1607d9fd10d191b2115a44d880bddfb9c2d2f36f9",
							"linux-arm64":  "f01b4d91018b68d778216b0c93197f229bc0ed2e7af179e0a3fd8b67426156c2",
						},
					}),
				},
			}),
		},
		{
			LocatorWithResolverConfig: config.ToLocatorWithResolverConfig(config.LocatorWithResolverConfig{
				Locator: config.ToLocatorConfig(config.LocatorConfig{
					ID: "com.palantir.godel-goland-plugin:goland-plugin:1.33.0",
					Checksums: map[string]string{
						"darwin-amd64": "ac141b8f463f1b505f82e1f3b1e4937205861829a778be9362eb972556e2601d",
						"darwin-arm64": "62f60bb74ddf1c8b49fd31165bd385ac0eb158e4b090860e4b90a6dfc69dfc1d",
						"linux-amd64":  "41380a5fd6ba6eedbdea36b898eab413afc5b65fdd17abf546f57e8ede4fa5fe",
						"linux-arm64":  "0b372f4da07ba257c7ed05d0993e9867212e06d58e3bdcfbce61088fbf79b2bd",
					},
				}),
			}),
		},
		{
			LocatorWithResolverConfig: config.ToLocatorWithResolverConfig(config.LocatorWithResolverConfig{
				Locator: config.ToLocatorConfig(config.LocatorConfig{
					ID: "com.palantir.okgo:check-plugin:1.42.0",
					Checksums: map[string]string{
						"darwin-amd64": "0f4c6f34f8ec9082ee654535aa47cbdc4e65c0536303c74cb0dbfb187131126a",
						"darwin-arm64": "557b873d672be367611e8ee025a982e479db530b3834524f5aba475db7f34821",
						"linux-amd64":  "776fc2f004c6952cc1c446aac5ac7909f25d59405a3c008b9e7d47f507f68c63",
						"linux-arm64":  "2127d6ae5460ed526431a7b7a4dcdcb22323e6aa832ee1903a3bf56a96fe8f60",
					},
				}),
			}),
			Assets: config.ToLocatorWithResolverConfigs([]config.LocatorWithResolverConfig{
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-compiles:compiles-asset:1.41.0",
						Checksums: map[string]string{
							"darwin-amd64": "07974a5ea35aeb1b29a3fd410bad6b93abff47357253d58c257182d5ccd40a87",
							"darwin-arm64": "564ff1eca8e9d162d3d6b57463970579d630d3f04345d377180b91cf4af49671",
							"linux-amd64":  "7c9fa4863eb24d58e1c0d70be85b4894669444ff2d7a224ab5b315e5c6a1d85f",
							"linux-arm64":  "e117e27c7ca64d0c855b2fa8e6748c5cab3f9ee7bfe480370f5e857e343e89c4",
						},
					}),
				},
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-deadcode:deadcode-asset:1.38.0",
						Checksums: map[string]string{
							"darwin-amd64": "3ad65c67135435b2a89fb126a93c7be540e7b1e062c701864a9f62476a5b939f",
							"darwin-arm64": "0be33adf893a3c67ff01e417103be0b942e75603e1c28090e437e17ee9b629fc",
							"linux-amd64":  "7bd4ad957359ea964e5731e1f90ee81b5fa51a74ed31cc215d813304618acdea",
							"linux-arm64":  "3a0c6bdd7b444ed2b5b868fb842c7580f8b46e0f1d1d163cb718520dc26ba394",
						},
					}),
				},
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-errcheck:errcheck-asset:1.39.0",
						Checksums: map[string]string{
							"darwin-amd64": "7de3245c9020e499a796170bcc4449c66b350b5d7ba0d3ac7c39ec2b96ab69da",
							"darwin-arm64": "7b29f051ac0a0c9f902cf43ffeba9199a47291e1dbae7805b3d2b4f92f0ffc8e",
							"linux-amd64":  "605f0ba695f5b75b260ec384da47ebaf866587754d12332276bc71d78f789b4c",
							"linux-arm64":  "53e999b8aaa594c3528044e0a7936427d3ea144f310261a853ca5cc1e9ba187a",
						},
					}),
				},
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-golint:golint-asset:1.30.0",
						Checksums: map[string]string{
							"darwin-amd64": "b107e5e519c71631fa065aa18925dba88066b07551dd116a603a82f93549fbbd",
							"darwin-arm64": "6f0501247f14d6f482c1f37431c503615adf832528cfe9e34d7d80e4d7b122aa",
							"linux-amd64":  "332c16ca538e309e54a236b69dfe51ca1c98ac949507ee1051c69412a4ee18a7",
							"linux-arm64":  "97288b6b74048c77091ddfcfc37d94e10b7cd4e0493d715e2c9279aacb5b2480",
						},
					}),
				},
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-govet:govet-asset:1.34.0",
						Checksums: map[string]string{
							"darwin-amd64": "2ff7a64feaf34376b9a6a58c4a310627d3f064989458dfbdae7e1bfa399d8016",
							"darwin-arm64": "53999c3617aa03eb1a7dd39b038d7f0f67dae0f23c85f595c57d0ec4e97a2fdc",
							"linux-amd64":  "43f8dd0e2e12ad624e9fe756495c4a35d555b253fc68dc55705bbce64a64a5cf",
							"linux-arm64":  "9ac50baf8b07182ca670a76e363fb26be23ffb5c24b3e18ac408c33ca8f53ecf",
						},
					}),
				},
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-importalias:importalias-asset:1.33.0",
						Checksums: map[string]string{
							"darwin-amd64": "dd69d1af4dcdf7f66cde993160c14d9cf7dc3ad5b0133d67a9d21606a51f39ac",
							"darwin-arm64": "6994f072962e532147427f687098f1ed35c5ead1f921653665162fe24353f855",
							"linux-amd64":  "d2fe02c1ece9f683a3fc0cade5f5642a222525eab81261e378f20ed12176ed42",
							"linux-arm64":  "08af97f20fcfa7bb3116b07c2d818efabe7720bb94a8c4b1d2dc9e8522dfd02e",
						},
					}),
				},
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-ineffassign:ineffassign-asset:1.36.0",
						Checksums: map[string]string{
							"darwin-amd64": "e269e73016972c1f844f28d01ce67400b4391f3865592ae63960c7b8651d3d2c",
							"darwin-arm64": "a8ff4a59e481688bacc1ad2604258b1ccc0cfdce12ecf863c0d08db9f27a9337",
							"linux-amd64":  "64ac95a9a0bc3be31b8bf2068881fd5574228c16904d2c24af4cdc1b45b8be46",
							"linux-arm64":  "2aaa6652417b765b16b2678db2e607960e919d25e03485ba884439f377ae9399",
						},
					}),
				},
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-outparamcheck:outparamcheck-asset:1.38.0",
						Checksums: map[string]string{
							"darwin-amd64": "876d24ee0952f4d16cf8d3b92234c3984c98d973dd78fcca3ff27769d54e6d47",
							"darwin-arm64": "456aaa0e5199d3a6bc668c1a9c77666c440b1cce04d10fc476d892e3def53403",
							"linux-amd64":  "60415a0bb9ce6acc94b70b9d879d0c527c7ac9a2227e75ad3ba82f9c803b1cf2",
							"linux-arm64":  "6f2d1842758b4c3935510a12330125833aa2ea3a6b8ce0a75e9a6b42af286fa0",
						},
					}),
				},
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-unconvert:unconvert-asset:1.38.0",
						Checksums: map[string]string{
							"darwin-amd64": "c4eb4f2786261f1005e6466ff8bc6948d4700d1ec1dadd7e458cb17de93a9c18",
							"darwin-arm64": "448553548ae4f2a5280c60a676362fa14dda48d9be1ad6883e011fa3ec2ff2ec",
							"linux-amd64":  "62a9ed7ed179d972c125b74ed819535a0f3e83ab439439c6bef945e1b58c9810",
							"linux-arm64":  "eb8bc381f38b230fd89907176bd7b7b918eb6fdbf582b1201162c0a95005fe09",
						},
					}),
				},
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-varcheck:varcheck-asset:1.38.0",
						Checksums: map[string]string{
							"darwin-amd64": "afb138dc4f8454aa1e482c70861817af3fe47a4f6de55870130d52f5ce6a13c7",
							"darwin-arm64": "f8bcac5517ec8f6776c810065c11a6903aecd87175f14944a7f14b2c7d7afce4",
							"linux-amd64":  "601cd2c8f06770a66474c137c971cb23e83778f719d5c8591051d7df209d6acc",
							"linux-arm64":  "dd6ffdbfbda245a001d4295f1b50e4dcd416a7d46a58bdb4690b9c99b2dfc406",
						},
					}),
				},
			}),
		},
		{
			LocatorWithResolverConfig: config.ToLocatorWithResolverConfig(config.LocatorWithResolverConfig{
				Locator: config.ToLocatorConfig(config.LocatorConfig{
					ID: "com.palantir.godel-license-plugin:license-plugin:1.36.0",
					Checksums: map[string]string{
						"darwin-amd64": "478bbdb55f295c8473b019e9a61c5adeedc82fe2a7eac5292ee9a3af6742a321",
						"darwin-arm64": "951a9dfe5afbc4e59b4f08598e80f9b3633ff74b7d2087afe2425815b1aafeda",
						"linux-amd64":  "7b94d677aaf53412d08c677c06878cc70083f69a6f9cfab3aa5e2caf010e40b3",
						"linux-arm64":  "5365076b81823818bab1f9605b7f5f1b9c20b74852285e3f2202a46cb94b60fb",
					},
				}),
			}),
		},
		{
			LocatorWithResolverConfig: config.ToLocatorWithResolverConfig(config.LocatorWithResolverConfig{
				Locator: config.ToLocatorConfig(config.LocatorConfig{
					ID: "com.palantir.godel-test-plugin:test-plugin:1.35.0",
					Checksums: map[string]string{
						"darwin-amd64": "82a1bd7dd4226fa856f960a699bc226b2b49b734b2d5bf27a0fcce002df9e5f7",
						"darwin-arm64": "f7169fc9fda2118486cbff07e816ae7a9399a93ca8d496811d492e79917cbcbb",
						"linux-amd64":  "5b211f6af7de6cc986e657c8273fa1c8bb614121f311de70b1b61ffb5ab35cc0",
						"linux-arm64":  "6fc773e9b1e21550736c06ad580e0cc1055a6c73737cb28b81ac48e1d36a1957",
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
