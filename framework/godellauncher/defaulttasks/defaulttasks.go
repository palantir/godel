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
					ID: "com.palantir.distgo:dist-plugin:1.44.0",
					Checksums: map[string]string{
						"darwin-amd64": "b23bf12945521e17e60780c37cca6d2d936a993eececf564c8cb38ebfc392b64",
						"darwin-arm64": "3420d96e4a2a9a285136099513048d9a4f8f28a3355390cfb4969eb71797f4ba",
						"linux-amd64":  "7c8b8061ccc39e80b805c2b27d0b458ab25b3d2f171d9d3493fa9cae2835dbd1",
						"linux-arm64":  "4592fb2a4693dfb6f3bb850a90bfe31b7e81f84a062bd2f6d4d3cf668e10460b",
					},
				}),
			}),
		},
		{
			LocatorWithResolverConfig: config.ToLocatorWithResolverConfig(config.LocatorWithResolverConfig{
				Locator: config.ToLocatorConfig(config.LocatorConfig{
					ID: "com.palantir.godel-format-plugin:format-plugin:1.21.0",
					Checksums: map[string]string{
						"darwin-amd64": "cd2b2392c146d5abe35330bdd12e26dae5f8c818cc756eafb3e84e1a099fac03",
						"darwin-arm64": "4c74d9d1cd268189df775817db271edda0709401439f6bd966f2eb7ff3b4f174",
						"linux-amd64":  "f550fd9612145c4ee878ee23993130b6662ef96d3a8797d4f389c357be21eb55",
						"linux-arm64":  "eadbfcf46e338c6b2b035d3766fe4ad78f036ab9d24644ee07a4be67151c5d65",
					},
				}),
			}),
			Assets: config.ToLocatorWithResolverConfigs([]config.LocatorWithResolverConfig{
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-format-asset-ptimports:ptimports-asset:1.20.0",
						Checksums: map[string]string{
							"darwin-amd64": "1615ed67a3aef7fb1988e786b74c2caf0826363c4745b681b1842a49f43f0d38",
							"darwin-arm64": "a01efb568c4e094ed33c18c2e9efe18ca527f96c906d2c578d438bd061cb8d1d",
							"linux-amd64":  "f19c315812afcda845c528e6d0d5723255fb275f615bf88e4a7a30402b6bfb8d",
							"linux-arm64":  "2ff21b8cbcd3054262d1f319e796311816a031a642e07f19f4faa6177059751b",
						},
					}),
				},
			}),
		},
		{
			LocatorWithResolverConfig: config.ToLocatorWithResolverConfig(config.LocatorWithResolverConfig{
				Locator: config.ToLocatorConfig(config.LocatorConfig{
					ID: "com.palantir.godel-goland-plugin:goland-plugin:1.17.0",
					Checksums: map[string]string{
						"darwin-amd64": "7febedb3d7c3d7e5b8afbd7bd1416c106dd8c4533f830fe4283505db9e5dea4f",
						"darwin-arm64": "a3179a7263c18ee15bfa678f7e382b82e4343806ace97460e62222309a8f4a22",
						"linux-amd64":  "17c01fb1d487968793d61df23c48e16cb2c2f08502b44b55d2ce4fcbad441381",
						"linux-arm64":  "cdfe3e1313ff6dda98f52c5191297852dfb2f31e5c2d8ca7909fa6fa786a14d3",
					},
				}),
			}),
		},
		{
			LocatorWithResolverConfig: config.ToLocatorWithResolverConfig(config.LocatorWithResolverConfig{
				Locator: config.ToLocatorConfig(config.LocatorConfig{
					ID: "com.palantir.okgo:check-plugin:1.26.0",
					Checksums: map[string]string{
						"darwin-amd64": "5ccc998936923e8027a8831d25fa6dd6c1630c96103c93f1da5ebc1f03a8222d",
						"darwin-arm64": "12b0744ac109d1945821b0f95440c11e4e7d05237aae1e120069bfb39a149981",
						"linux-amd64":  "5281c46fdadcb09c7c724a165f3abdec5c0730dc2522ef7d6631955768dbec22",
						"linux-arm64":  "b80c5b42414c41588c9b8b707614f2efd7556b4837a881d19eee54de1515911e",
					},
				}),
			}),
			Assets: config.ToLocatorWithResolverConfigs([]config.LocatorWithResolverConfig{
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-compiles:compiles-asset:1.24.0",
						Checksums: map[string]string{
							"darwin-amd64": "509fa8b2529353e30be4c9a9a8ce45b86a7ee8821e0de4e3f35a60082441143d",
							"darwin-arm64": "8c2b6e5130b9b0d1d6c3c44475e014e3bafd7fbeb47a91594f1f2a2e88136481",
							"linux-amd64":  "ce591b2d7b102f3da17e1802c74348db4ff9a26d2699574a57fe5ba0a3c27b90",
							"linux-arm64":  "ba9df79f2acaa806b9bdd3d759cb7bf37af6d0c2eb2b0af54317ce252a952f9f",
						},
					}),
				},
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-deadcode:deadcode-asset:1.23.0",
						Checksums: map[string]string{
							"darwin-amd64": "04dfa687edb8943d5106848cf945ed9bc41c011425158c0ca4d75cbadaf43680",
							"darwin-arm64": "e4087508263440dfbe1d085e92e1e02e6218c9624214ea6cc8d601e9abd68a77",
							"linux-amd64":  "ff8c7d12d2cf4fa1681ca774691f1a053387a80cf5e3a9360b2bfd952980a6d6",
							"linux-arm64":  "fcfbac7a879de24c1be5bb391caa9caceff923f5ee645fd0f4070753c3fcc84c",
						},
					}),
				},
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-errcheck:errcheck-asset:1.23.0",
						Checksums: map[string]string{
							"darwin-amd64": "aa865c7f727186b7f670cce55e58e55d7f7b2a2e465f81b4ec00aa00f043eede",
							"darwin-arm64": "5787bd47c09955c5f10974d98ec9ef9fc015f1cf1de9a14173d181ab6e34024d",
							"linux-amd64":  "4bb419b1b0d620ae575f503128a6c448ec305a3b9f171656bf0bb74959731784",
							"linux-arm64":  "e707d039406dd2d34a7813997505e3402e81a23615de9b694f2349231a7e59fa",
						},
					}),
				},
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-golint:golint-asset:1.15.0",
						Checksums: map[string]string{
							"darwin-amd64": "445113bd71e7e8dd629a5a2fd3e103edc9d885a9c55b63a9c14a2c7b378a6308",
							"darwin-arm64": "7aeadb76032d459864cb1fa0bf91cd2fba64a7602a5624cedb52d1b6e3cb6a53",
							"linux-amd64":  "fd311686957f4854317dd9dedfb40b460704a61b5e0eaa031669954458b95bb7",
							"linux-arm64":  "0cc0187b90d5f369cda2cdcde11f48ec6c7728c61f3454a30f2ae2a55bbc37dc",
						},
					}),
				},
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-govet:govet-asset:1.18.0",
						Checksums: map[string]string{
							"darwin-amd64": "2d52cc6c28e63f8d7c03dd2cf2bc224bf0b4d637570e371e5c8f1d20727324f5",
							"darwin-arm64": "3c4a63b16c43c64b54fef73abad42ce52cdf4d54503183469affe7ceb4c66e93",
							"linux-amd64":  "9400e10f2fcbd1230f65010f9d5bd9855811a9af0c88cc98c63802f0ecaabca5",
							"linux-arm64":  "39fd1f5acfdbbabc3b6a48ff052b475ee9bba2bdd7d592cba0650f32565a4675",
						},
					}),
				},
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-importalias:importalias-asset:1.18.0",
						Checksums: map[string]string{
							"darwin-amd64": "206f3989d52e0263efce3e8128f504f6e322fc5f6be29fd799836bee322e65bc",
							"darwin-arm64": "84a6a8557834f077f8b169780dc62040b1925984fb3f3275b818a37ecb46dfcb",
							"linux-amd64":  "d8674d5816113942b4d147ac6ef4d44d56e78295681be5d01d77b3f263623c86",
							"linux-arm64":  "c994489a0ddc7e23fc22e458401280ea48765e233a29bfffa3dc387265b2e131",
						},
					}),
				},
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-ineffassign:ineffassign-asset:1.21.0",
						Checksums: map[string]string{
							"darwin-amd64": "7d54106f2366d038806a30da785e178cb7c20521f81907f8a435fcf3bea3941f",
							"darwin-arm64": "761397fc33503a73353269e08a02ebc7bb1d6151bde322ce61b1e3f01f67c86b",
							"linux-amd64":  "49561ab99ff9bb5512acf24cc93ddcc05d8a5b6f5aa4199b2589b5b2e52644d5",
							"linux-arm64":  "ee8166f3613b683c747a671d2fd6429d1f1db4564826e89fa5bef2519917b676",
						},
					}),
				},
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-outparamcheck:outparamcheck-asset:1.22.0",
						Checksums: map[string]string{
							"darwin-amd64": "ef761f99ae2702b9897368b6c4edee44bd9e7b939336f11876144bc6caa40a91",
							"darwin-arm64": "330bb38608c844669cda204844752c391718271ede7e848bb7a3fa173ad5ace1",
							"linux-amd64":  "323259f474a2f86d20d00cac7ea89f7373e7128bf2ccdb2eb59f519fc49a3c2c",
							"linux-arm64":  "8412f3d7028e48ce6f3843d30cc23f77d5fb3dd70b6aa3570510a36dbc248cb6",
						},
					}),
				},
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-unconvert:unconvert-asset:1.23.0",
						Checksums: map[string]string{
							"darwin-amd64": "689db3ab7eefbb69dbf8a56c1ec12bf3ea5d39f9e55afc92d9bbcc83cad0e47c",
							"darwin-arm64": "f2be8a87a6f4f27c587ba52323c60dea16b541036395e72e91556a437742558c",
							"linux-amd64":  "ea0389231519904f191de9384261e922adf37b24fe556f21d75841fe52cab63c",
							"linux-arm64":  "86426ff232217412af420dba56e4af961478ae3698c2272f5dc8403356382e0a",
						},
					}),
				},
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-varcheck:varcheck-asset:1.22.0",
						Checksums: map[string]string{
							"darwin-amd64": "451d9ab782129a7971385154ae55be810dff98a36fa95896ed46896ce854f4f2",
							"darwin-arm64": "5cdfe2fad0955ec4a1a9829dff4d70b3dd353e96aa9c672c4fccf8d112f5f8bd",
							"linux-amd64":  "1ed01ab3bc7ae3df9dd2469cd891924daf9acec81a97a55089a6ac734274f00d",
							"linux-arm64":  "c87831c4911ebc679878c033f2ce368b2b053ef50aaeb9dd047cd815585df49f",
						},
					}),
				},
			}),
		},
		{
			LocatorWithResolverConfig: config.ToLocatorWithResolverConfig(config.LocatorWithResolverConfig{
				Locator: config.ToLocatorConfig(config.LocatorConfig{
					ID: "com.palantir.godel-license-plugin:license-plugin:1.19.0",
					Checksums: map[string]string{
						"darwin-amd64": "dbd6ec4417d06c61c1968004a1d3e28cb66ada8ad6046ce2ac68db4ab701e003",
						"darwin-arm64": "9a40253650e23b9b063ad51ea03aa9f1422c67e7e1863bcfc71f578daae908c2",
						"linux-amd64":  "74be1ae70359227e9382d76d72d17912d06ff1ee92db320c3ddf7372b134f3d7",
						"linux-arm64":  "8bdc84009f051ec0e463a5417f5b3b19fa3d2a590be405c63e11b984255d2292",
					},
				}),
			}),
		},
		{
			LocatorWithResolverConfig: config.ToLocatorWithResolverConfig(config.LocatorWithResolverConfig{
				Locator: config.ToLocatorConfig(config.LocatorConfig{
					ID: "com.palantir.godel-test-plugin:test-plugin:1.20.0",
					Checksums: map[string]string{
						"darwin-amd64": "a3d3d5e31d3681359e7aac6c6662b4d44c0a2ef311287ddda96cdd6956a29152",
						"darwin-arm64": "c7cdbd012f2a27c910cf6418563c8fa330522aa77a5c9cecdce061c889045923",
						"linux-amd64":  "6c3972952b9daf58bf1578c9d06d325a5ac35e28276033b437d2acd66a9754da",
						"linux-arm64":  "4ece024e07cc92930129d476c89444cadad5f76adb62d1862a6c7c2ed001a926",
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
