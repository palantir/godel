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
					ID: "com.palantir.distgo:dist-plugin:1.49.0",
					Checksums: map[string]string{
						"darwin-amd64": "293d732d3642bb6efbef12d223e8c55a8f8d07287b7ee517416e02b0b4b2787f",
						"darwin-arm64": "9ddb653038d0ff672a69ba0cc3eb57b1f4addcd96e722576138a9255094ae663",
						"linux-amd64":  "44c5b7d961bb6b73533d1a4cda294bc1bb0124780cbc37775469aebf7a7a1721",
						"linux-arm64":  "b4f8a5f2fdc02310c5e89dca0b01ce4ec4ab98244500d787025cd6d078b00a0d",
					},
				}),
			}),
		},
		{
			LocatorWithResolverConfig: config.ToLocatorWithResolverConfig(config.LocatorWithResolverConfig{
				Locator: config.ToLocatorConfig(config.LocatorConfig{
					ID: "com.palantir.godel-format-plugin:format-plugin:1.27.0",
					Checksums: map[string]string{
						"darwin-amd64": "68b78df9b53331a5d9e812f89207d389dfbdb2c8d5fd20b00a928472bffba663",
						"darwin-arm64": "ca9d81d370fd5aaedd6ec4cf3ed956fca87321eae4124c8b9df90e973be4947e",
						"linux-amd64":  "6c40aa420d60863ccb66e0c8f46c05f03ac32e749f8b6145a9cb3329909b3f1a",
						"linux-arm64":  "77d0f872ac4f82e06a9b7113dab9c013e4e979c022acd7fbe075c1dc300a5ee1",
					},
				}),
			}),
			Assets: config.ToLocatorWithResolverConfigs([]config.LocatorWithResolverConfig{
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-format-asset-ptimports:ptimports-asset:1.26.0",
						Checksums: map[string]string{
							"darwin-amd64": "cc5a7b9c78c1b383c45a6554423253eecfe705c2bc261f4e3d9ab05174ae7188",
							"darwin-arm64": "547c1bf33b0931c42f255ea304a054fcb6c9288e8dfdef4934c2037f9ac55970",
							"linux-amd64":  "d3bb948843282c221de56ed4db91410307e9135d958c77ea3c069f1fa3971139",
							"linux-arm64":  "9f67b59a8140a48f807d66b6eb5fcaa910770e005b8ae1c5ba03b986e4f8742e",
						},
					}),
				},
			}),
		},
		{
			LocatorWithResolverConfig: config.ToLocatorWithResolverConfig(config.LocatorWithResolverConfig{
				Locator: config.ToLocatorConfig(config.LocatorConfig{
					ID: "com.palantir.godel-goland-plugin:goland-plugin:1.23.0",
					Checksums: map[string]string{
						"darwin-amd64": "2c431867c27aca9e2241eb7c74a00fec964e2d7e51e4302c164707d4f69f4bab",
						"darwin-arm64": "b5699362c6fb9a293c84a2ac454ee81e138b31b56ab21c7ff64b8787574c4f19",
						"linux-amd64":  "7abae3d5ade7f14c0095f568121685a4bde3ecad0de1ca1177fe57ca92e4b0a9",
						"linux-arm64":  "0efb176aa5ec6d1f2f54c51a70b872888e32d3d79fa239de3faa353e4401d68c",
					},
				}),
			}),
		},
		{
			LocatorWithResolverConfig: config.ToLocatorWithResolverConfig(config.LocatorWithResolverConfig{
				Locator: config.ToLocatorConfig(config.LocatorConfig{
					ID: "com.palantir.okgo:check-plugin:1.32.0",
					Checksums: map[string]string{
						"darwin-amd64": "7e74365f148f0c2209f5ecbf73c308c533abcd9eee843141e35fb76dc9be6d33",
						"darwin-arm64": "1b6729041a8d24ad36af346956f41d772a9e22247e15f0c21847594e4c23b9dc",
						"linux-amd64":  "9ee58a3dfa794b7fc72242e09a3a7837d2b752a9fefb4780d5e03bb5a3ddd492",
						"linux-arm64":  "f1436797051521a2432ee4bc9cc0990d8fdebd278f6ed1efb14c664753dac9b5",
					},
				}),
			}),
			Assets: config.ToLocatorWithResolverConfigs([]config.LocatorWithResolverConfig{
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-compiles:compiles-asset:1.31.0",
						Checksums: map[string]string{
							"darwin-amd64": "0c7c23e0182d5dfb8c3e719960e2b19bdd14d540c1631cf1f9cd2ea81c563825",
							"darwin-arm64": "62a45c7567c90ad28f9d2379dbba6274d21a0eb09131645f0159ea9f12f6942f",
							"linux-amd64":  "f14154bc8481a8f0deb214f72a8711e6e5673bd343aafc8e9d33cbb2c3c2a019",
							"linux-arm64":  "49914ea77b1f31383d72e27c022c32c08e8aaf2001f47cc36bb9535ec9bc72d8",
						},
					}),
				},
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-deadcode:deadcode-asset:1.28.0",
						Checksums: map[string]string{
							"darwin-amd64": "f1a4be5cc28542ea83d78afacfd1de46790d9b722be0ca3ce9be327201d00e1c",
							"darwin-arm64": "4f2a1cdcc8f1c2ddf111b104743dfaf5f177f1cc55f7b1b77e92fdb4aca2c831",
							"linux-amd64":  "e412135f6116cab785ee274bbeb3fa72b54bd9f4b64047bb3012285bd0c48ab8",
							"linux-arm64":  "d0536cd428fa96c76a7881a16b1a51183f4e322d6cf2e2950aaa50c70ffc63f1",
						},
					}),
				},
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-errcheck:errcheck-asset:1.29.0",
						Checksums: map[string]string{
							"darwin-amd64": "4af30f65888071ce831c31c366fb5364de10bfb4a7606369d6fde4ce5bdf44ff",
							"darwin-arm64": "1fa8995ea1d08a401145dc8c04ebafe6799e23aa59d2d9bd776993a07437198e",
							"linux-amd64":  "a4d44e4def8d058c864d3274b0c69985f0e095a0561ee51de219214049fd2e1d",
							"linux-arm64":  "e293bf984b7af12f10b95f91610f05a6dcdd4c5f78da4bad921eb4d8834ed27e",
						},
					}),
				},
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-golint:golint-asset:1.20.0",
						Checksums: map[string]string{
							"darwin-amd64": "6a79384dbfd9735f653f9a4819e2fbcc6004945eb9090dbacc07e045b95b666a",
							"darwin-arm64": "ca5e8cdeb403f1b6a9bc312ec6fde27c66b66a99dc5ad6f2b8ae327b73731610",
							"linux-amd64":  "48d64300dc56569152cf467fde75fe4e73413cd55ea7c5dd6690d5a712507fee",
							"linux-arm64":  "40220268757b59c2af35f33b5d89fe9eae24ae2617aa9162ea37d07a17f984e5",
						},
					}),
				},
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-govet:govet-asset:1.24.0",
						Checksums: map[string]string{
							"darwin-amd64": "e69572c1fc10c72f1acb55ee5510f776404d6d2199494d6d919eb94228edc8f3",
							"darwin-arm64": "f4f7c0c2cd08abc8f97de4b212263fd20414ffe4e0596846b34d8dbe761f455e",
							"linux-amd64":  "c57266f472626ccf40c31a7c6e2228dbb66ec113da47d9a5436c74b566fd2030",
							"linux-arm64":  "0a223c98e54d0f63891cf80cf04973464669dac452584698a019ec3fafbcb5bf",
						},
					}),
				},
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-importalias:importalias-asset:1.24.0",
						Checksums: map[string]string{
							"darwin-amd64": "9561903f74928480ed72ab8be0e166d9574e049081e680cbe42130c0cc13a1a6",
							"darwin-arm64": "cebb2ac118a72f423fc9163b14a6f78323ee8078395e10119d37fa4905795f15",
							"linux-amd64":  "9d092318b895c778ada01d9a789791cc616b562a9a89b24e220320525ec16fac",
							"linux-arm64":  "b492e6e8e22e54e68f407800c79034cde0a33978d12f319a05224c5b9689702f",
						},
					}),
				},
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-ineffassign:ineffassign-asset:1.26.0",
						Checksums: map[string]string{
							"darwin-amd64": "1f0417f1c0ae29f872f1d75e3f66a63c7ef828c0e68d797bd03213fea51f0145",
							"darwin-arm64": "c3810f2f6d9c61a9caa74c3c0f83efbd53df94157d42e6fc0a97a31364d8af4f",
							"linux-amd64":  "ed5d9946b50f88f951739ae2ddcb8952b1199c580e3f47a4074d98ee5e531b65",
							"linux-arm64":  "14c605535e4e0876fa27ec404909b30c71c762b74f5303a16e248bf04d9dfb0b",
						},
					}),
				},
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-outparamcheck:outparamcheck-asset:1.28.0",
						Checksums: map[string]string{
							"darwin-amd64": "6216f85b614892b1a3e24615ed86309a27cdbb9736b0548d13207da44ce5a275",
							"darwin-arm64": "4fd32a4d25e1188688b7985280e43247f67ca96b013f730e38981f8f55bf6ba1",
							"linux-amd64":  "a58b691469785f5f62be4d09b227a05fc806c2ea0f8b8db09607c12a7f9ebb6f",
							"linux-arm64":  "b966ea828fc60836da6cc35a32f6a2b28377279d37e6b6c306a224b9d08da7d4",
						},
					}),
				},
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-unconvert:unconvert-asset:1.28.0",
						Checksums: map[string]string{
							"darwin-amd64": "b5c120f6f427ab698e575e627f7f614ee2f13519260acc30582f0e59309cb92c",
							"darwin-arm64": "65c6590f90d6ecac817a26db82eb61ab3a27548b280b4d9add21f9ec4a32c5d4",
							"linux-amd64":  "1702f73b8fa6cd3f9eac51aada8221f3bee3ce3819eaadf72f3eb2d360155a7c",
							"linux-arm64":  "0f408d1dff0b0ab561fc9a02673cc90496b39b1f875cc2421c1ee5b0220f862d",
						},
					}),
				},
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-varcheck:varcheck-asset:1.28.0",
						Checksums: map[string]string{
							"darwin-amd64": "3e8c0ddf62dfd03f0a402ff1802513d3beabfb4c221d1321686fba9e54ac2b10",
							"darwin-arm64": "7738ae5f0ad9b86a88a08d779f5f4fa398216b46a1e906d2e995bf069d460b71",
							"linux-amd64":  "21645981cd268dcc06f9da6dede3cbd4c79af16477e2667adc3ff86836b6b3f5",
							"linux-arm64":  "33de6f583dcd9c5402fa132c2d7dca2b1fc6bd00caf87e8672bec3a0b34fbd9e",
						},
					}),
				},
			}),
		},
		{
			LocatorWithResolverConfig: config.ToLocatorWithResolverConfig(config.LocatorWithResolverConfig{
				Locator: config.ToLocatorConfig(config.LocatorConfig{
					ID: "com.palantir.godel-license-plugin:license-plugin:1.26.0",
					Checksums: map[string]string{
						"darwin-amd64": "f78b2175049e6f2bd90b9fd34a7ce2413a4598b1a3464cbbb58304046b04fab4",
						"darwin-arm64": "9b0f6525655c4d088f737fcd721a004e326e2522d11f04901681a5739c678e14",
						"linux-amd64":  "3ab98f80b28f88a2c86a412e4b80ae0ba04174aca79f5782bf4f1a32d38291f5",
						"linux-arm64":  "7a7117a7480c55b03cc05f7f75212b4288c739848ace520c7e66e797aa09be50",
					},
				}),
			}),
		},
		{
			LocatorWithResolverConfig: config.ToLocatorWithResolverConfig(config.LocatorWithResolverConfig{
				Locator: config.ToLocatorConfig(config.LocatorConfig{
					ID: "com.palantir.godel-test-plugin:test-plugin:1.25.0",
					Checksums: map[string]string{
						"darwin-amd64": "fd09d3720ec3552a3067372f30e0995d06978edbd9e04f0b9d6938fb06f47a03",
						"darwin-arm64": "182c730f2ec9e5eca54980d78823873bbbe7390f8ba3450f9fbaa327ef2f9a9a",
						"linux-amd64":  "8f56cc31cad6fb232b9b9365fadde8212888573513f38d66b7561d25ef5f8ba9",
						"linux-arm64":  "0b6a0fafd89bf77cb5dc979d5ad35bcded91911f35bafed177ba3f2a8855c841",
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
