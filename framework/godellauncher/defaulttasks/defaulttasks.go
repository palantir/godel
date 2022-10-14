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
					ID: "com.palantir.distgo:dist-plugin:1.46.0",
					Checksums: map[string]string{
						"darwin-amd64": "66de1f4bc8465cfe56f659a0f83de8fb81003cd578dc8d8cd47623c729426e48",
						"darwin-arm64": "24a54779f23330396b21372f38b20205446b93c16553ac4e206f84f8adc8a675",
						"linux-amd64":  "27970437f0f7a8d44aff717102e424d48fe0ea842b4e845c744186384af71719",
						"linux-arm64":  "2bcca14c0ae3a1c3cbf02ade850d94834eff77f817ec1c80621061b2be227b60",
					},
				}),
			}),
		},
		{
			LocatorWithResolverConfig: config.ToLocatorWithResolverConfig(config.LocatorWithResolverConfig{
				Locator: config.ToLocatorConfig(config.LocatorConfig{
					ID: "com.palantir.godel-format-plugin:format-plugin:1.24.0",
					Checksums: map[string]string{
						"darwin-amd64": "b57707792f38f19fe7eda72d26ec545524304a9b32412804f467533846c9eca9",
						"darwin-arm64": "e5b82a166fed70cdf0c6ceae61dfc07f18dd2c1247b1228f960a88b2a585a174",
						"linux-amd64":  "32357db8d7d6f452d409c74e06d1383d1b91a406bfad726784e2abe0ae7d09fd",
						"linux-arm64":  "314c45c754266ea4a6e2d2bb3ce6422fa7d03a41a86c27e35480b00211d895ed",
					},
				}),
			}),
			Assets: config.ToLocatorWithResolverConfigs([]config.LocatorWithResolverConfig{
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-format-asset-ptimports:ptimports-asset:1.23.0",
						Checksums: map[string]string{
							"darwin-amd64": "821cd3facb0a01225ad87fc06f9426d2860121618fbfe079584125bb2c876459",
							"darwin-arm64": "06da4107fa1a9285400084a26bae367bc6acefbcba0c1b1f3b1ac55e24bf120c",
							"linux-amd64":  "16ef60e954e42656e56d29533cc8cffb5a1005ff772e3448a9667efc3824225b",
							"linux-arm64":  "d0076549d838ea8acb7d4909c6d403609e568b8957ed896ef69a3b34b438d2b6",
						},
					}),
				},
			}),
		},
		{
			LocatorWithResolverConfig: config.ToLocatorWithResolverConfig(config.LocatorWithResolverConfig{
				Locator: config.ToLocatorConfig(config.LocatorConfig{
					ID: "com.palantir.godel-goland-plugin:goland-plugin:1.20.0",
					Checksums: map[string]string{
						"darwin-amd64": "6d69d6314184fb5a0a069ed375a246cee7f60f6fa280e1684e372d35bdb2c37d",
						"darwin-arm64": "2e4b5ea8394de5030752679f2aa351b75395f0a04c658bb46d2cae758c7a8238",
						"linux-amd64":  "877fe9de1c5885ff4464ace5cd334e9bd0a1443d0604c2a07003ac6f2d8eaf6e",
						"linux-arm64":  "bd9e991e7d317c8e628a8271b80490c91803a00aed41f3d4c363223dca98edd3",
					},
				}),
			}),
		},
		{
			LocatorWithResolverConfig: config.ToLocatorWithResolverConfig(config.LocatorWithResolverConfig{
				Locator: config.ToLocatorConfig(config.LocatorConfig{
					ID: "com.palantir.okgo:check-plugin:1.29.0",
					Checksums: map[string]string{
						"darwin-amd64": "3cc008be5407852d09340b567ea56098615e21efa65930c6fb653605818dbe5a",
						"darwin-arm64": "6a022ba2d21a28c1899806802ca37ac5c8259e5f1f86c33d92cdca142b5e3e15",
						"linux-amd64":  "2a187a627de7eb80c839f52c9af1849278990e2ca3f557d39de774d3e443016a",
						"linux-arm64":  "b9a5cc35a3dd4e3a03f5f846c4e5e3d00e11bafe036674e2d74f3e74bda8a14d",
					},
				}),
			}),
			Assets: config.ToLocatorWithResolverConfigs([]config.LocatorWithResolverConfig{
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-compiles:compiles-asset:1.27.0",
						Checksums: map[string]string{
							"darwin-amd64": "b3caac26a740f929e35a4e7327851eaa5cd6f816e996e36bc14d8dc18beb27b5",
							"darwin-arm64": "52e95c4c4adefc5b15a7c39e9125a67b08b01bc7a14b5283266b9e212a701323",
							"linux-amd64":  "c6166b9affab5fbef6e4fb2895b0a0c2158c4379051994eaf4a36924df1083c7",
							"linux-arm64":  "efa5f189147193c5cca5604cfb8cdd7fede3907121f290c717f10148de618a7f",
						},
					}),
				},
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-deadcode:deadcode-asset:1.25.0",
						Checksums: map[string]string{
							"darwin-amd64": "d98a36d0257e148075ef8deef6a5c901ec59c1e6a092d6f96a7226879bcd4b99",
							"darwin-arm64": "29fcc50b3adf016f9707a640f58e2a9f08ed3b1bdce6624c6c68d703f01e292f",
							"linux-amd64":  "10d0f899e4e6cc9655dc3ea59dc79b12a108830df92c19b882d441ec2ece17f0",
							"linux-arm64":  "8a573f6b2e8fc937857cbb0c47d2bff75eadb2af63ab485878c1dcae56f5b305",
						},
					}),
				},
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-errcheck:errcheck-asset:1.26.0",
						Checksums: map[string]string{
							"darwin-amd64": "52c4d5e53da6c8d14deaf0ebf3a29fbb19cf766ce9ca38b4870af75804cf2261",
							"darwin-arm64": "9d7c6a0154d7af0fe1b7d37d75712993da174fd90dc16a0a45c2258085d44a9a",
							"linux-amd64":  "dffcb6b87f9f7da09375d4375b5bc77095621d7ea11d5b0c8f81af6623da2a59",
							"linux-arm64":  "804d93cec6985613c645c9a379cb84f048ed571695a106446996189e374da71c",
						},
					}),
				},
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-golint:golint-asset:1.17.0",
						Checksums: map[string]string{
							"darwin-amd64": "62d3d7afa129b51bd483dbede6952ea995561b9f721c645a513e4baafeab733d",
							"darwin-arm64": "9a08abc04db991e5bbf2b6f1dd572085e184d23357d97d3299ba9666b543e0bf",
							"linux-amd64":  "9a6be8f98f8f915f439de7888c3b652d6032b26ac1dbbe410874c642c213b2f0",
							"linux-arm64":  "6bde03a9954400c43b9c6256e8c6d4a2ac2245e0cf92e35b1842328e5fce209b",
						},
					}),
				},
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-govet:govet-asset:1.21.0",
						Checksums: map[string]string{
							"darwin-amd64": "3d414eb556cba0fe3a5f9a051cace4dde18a034c364c9158eb2418af5ebddc71",
							"darwin-arm64": "5d16860fcecb011cd47a67aa5e9209dd0e8c79e1b404d4ff276e123707133c7a",
							"linux-amd64":  "ad76d231fe4f6ca1e8fa8f940971d8798388e11f13b5f2c29a8be8e9044ad784",
							"linux-arm64":  "add1a1d9a5e289a9fac8fc6ac3530f8aed186c35441e363c9964403567de9ecb",
						},
					}),
				},
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-importalias:importalias-asset:1.21.0",
						Checksums: map[string]string{
							"darwin-amd64": "1c3f49c70465b2676e9964c443590c890bcdd776fb8f7cba2a7a008571ade4ee",
							"darwin-arm64": "292c09b4afc641bcf1d7f16f118902f12d2203c0c43efac0d977f363f0398c58",
							"linux-amd64":  "daeadafe2fb82793d2a1bcf398415e5715746127476f07383d3b896cbd611ba1",
							"linux-arm64":  "9bff97fa68557a7684203bbf9b450e16269ecd3f2a1a229ac065fe80eb0e7a76",
						},
					}),
				},
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-ineffassign:ineffassign-asset:1.23.0",
						Checksums: map[string]string{
							"darwin-amd64": "1a2288a0d346740744f9c926e0b6f7d18e97f74ccfbc4c14b31bc810c4fee0cb",
							"darwin-arm64": "c7ca2479e72becd2d96852a790cfcb135ef1e27628931caa979bf83d5fa9c1e0",
							"linux-amd64":  "11bd688f4717dab066bb87f817211b6e739b148166cc184a1ab04209b47f503d",
							"linux-arm64":  "f434de4bffa83fe6d221c2d47d88be843e6e441170f94ed9d9575bec3651e25e",
						},
					}),
				},
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-outparamcheck:outparamcheck-asset:1.24.0",
						Checksums: map[string]string{
							"darwin-amd64": "b6d254fa04fb42ae0784c7cb219970dc50bd18941dc70193b922c23c2f188f38",
							"darwin-arm64": "6fc92797f00b6e69035b6da4a9eb90f972b386eae963167777b852fa66997774",
							"linux-amd64":  "e4cbc4408269250b546eb1f0289d182c00ff97ed6ab6788afeb060aab83a334d",
							"linux-arm64":  "47bd0254d8d669591a974a25443dc59a94e0122656b86cde1ec8b680b73e133e",
						},
					}),
				},
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-unconvert:unconvert-asset:1.25.0",
						Checksums: map[string]string{
							"darwin-amd64": "001783ed100dcb97773f5ebdbeb6b4538ab6742c83fc551ad9c0a70cbd226219",
							"darwin-arm64": "78978dae2a7636adbcb900407b07aff84ea968919c3df0cf361bb21c24e3444c",
							"linux-amd64":  "4a6e2b3771529091f26c6f86fc7f6fef64691ffff2a6e37b191933611e558f24",
							"linux-arm64":  "8cf7874f7d61a9fbe5411d96150abc0e61f518b99059debb47ef5c3d37528f1c",
						},
					}),
				},
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-varcheck:varcheck-asset:1.25.0",
						Checksums: map[string]string{
							"darwin-amd64": "7cc81530e858c6af38063419ce157bad28f6ae90d0636dfd1db8b1a0edad5aa5",
							"darwin-arm64": "71040a5f2cc9a8a887ccc9c7c5c27b8acf111a74daf64f1dfa9686dcea3fb10e",
							"linux-amd64":  "e6212aaa6144dce271ca7176a0816385c9dc676e3782509f46ce3ebca8e2d437",
							"linux-arm64":  "e0f6b4ed9486ef665dd31c926e4f547ded40d54347727f1e29d9f4f0059c200e",
						},
					}),
				},
			}),
		},
		{
			LocatorWithResolverConfig: config.ToLocatorWithResolverConfig(config.LocatorWithResolverConfig{
				Locator: config.ToLocatorConfig(config.LocatorConfig{
					ID: "com.palantir.godel-license-plugin:license-plugin:1.23.0",
					Checksums: map[string]string{
						"darwin-amd64": "248517f5b2a5d1a2f134d60a6230ea96e1c108d9aae2981c65d7c9e8ceabc383",
						"darwin-arm64": "bc618152fdadde8a90ef2d8cbe74a5ac9ee8ac2953a2c3813451b35a25378191",
						"linux-amd64":  "9f59190e16d1f461a1431522478d62601c148a4ab2edb65f7cc0290794a632c0",
						"linux-arm64":  "ca5dfdde630052ec506b031fe90c7273ef58bac40a21fe61f2ea0b468d26e2da",
					},
				}),
			}),
		},
		{
			LocatorWithResolverConfig: config.ToLocatorWithResolverConfig(config.LocatorWithResolverConfig{
				Locator: config.ToLocatorConfig(config.LocatorConfig{
					ID: "com.palantir.godel-test-plugin:test-plugin:1.22.0",
					Checksums: map[string]string{
						"darwin-amd64": "8d9e041b15b62080ebc1ef4e21600b2da7301f5657100afb2cb76fa230b127bb",
						"darwin-arm64": "2e52e235066968b6ef8f8465edd4ac1bd63d78caeb35c1fae477f255cf36e135",
						"linux-amd64":  "38aeee03263717a08e8418484fe54b6b6b089d9aaf4ae5f18cb5cd45834dab24",
						"linux-arm64":  "e78bf944fab6d0225d93465b07a31ae32cd2bb4dda0cb3797c2b8431aea18657",
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
