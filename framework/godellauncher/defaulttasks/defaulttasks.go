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
					ID: "com.palantir.distgo:dist-plugin:1.36.0",
					Checksums: map[string]string{
						"darwin-amd64": "18330d074a5dca280908ffed9ee5b9d834693ccfd086f8a84a1665d346a3b433",
						"darwin-arm64": "7797d9c893fb9b868185b8abddcbffc87e7015f6d89d9090ee4710e678ecc69e",
						"linux-amd64":  "a3165bd1b73f546b33425b8646283f0f1609978de6f5eb14236ee62ebd439465",
						"linux-arm64":  "52941e58f35fd441b8770957e74c92cd2c9a63d17cbf33132c736c67c1c86e0a",
					},
				}),
			}),
		},
		{
			LocatorWithResolverConfig: config.ToLocatorWithResolverConfig(config.LocatorWithResolverConfig{
				Locator: config.ToLocatorConfig(config.LocatorConfig{
					ID: "com.palantir.godel-format-plugin:format-plugin:1.12.0",
					Checksums: map[string]string{
						"darwin-amd64": "95f0e01a31700f70547e2cb8d95402191ce3e20b0c7a9bd110d2201bf7271b33",
						"darwin-arm64": "4c1c01cc9611a6f2816241601215c4d76b1957a7a0442e316360990048a31a41",
						"linux-amd64":  "d6a4897b365d5bdf8bf4de924279bb4522f487ea630001309b48b4452179596d",
						"linux-arm64":  "5f4e5ab36243f1b94c01daac607ad42ce574f34258185c355f9ab508e383121d",
					},
				}),
			}),
			Assets: config.ToLocatorWithResolverConfigs([]config.LocatorWithResolverConfig{
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-format-asset-ptimports:ptimports-asset:1.12.0",
						Checksums: map[string]string{
							"darwin-amd64": "94a3cf390cbcf423a14d9157636020d23857efa1661e006215bdcdc0ac99e70c",
							"darwin-arm64": "bf23b400ab3999cab7712bd92f86e9a7a0304584b0b6e434eb3acf90f1334b79",
							"linux-amd64":  "a67513e739645948223fdde85bb5e3db78fb09cf92ed808f4e39eb1bd2a34c75",
							"linux-arm64":  "22935b6e24957f1cef0a403007b8d354b9f70e05417db97dd77536cfd9dc2e0f",
						},
					}),
				},
			}),
		},
		{
			LocatorWithResolverConfig: config.ToLocatorWithResolverConfig(config.LocatorWithResolverConfig{
				Locator: config.ToLocatorConfig(config.LocatorConfig{
					ID: "com.palantir.godel-goland-plugin:goland-plugin:1.7.0",
					Checksums: map[string]string{
						"darwin-amd64": "04949463dfd5967d0fbbf6c2c99bbe784b4c6ec9c8b01de6361e75c673115cd5",
						"darwin-arm64": "3c9c575c3a249ac3f3632354d09966d4c72a092625619cebb706f81013e0d762",
						"linux-amd64":  "cc324e7f7eb7e8fe755b610ac61a5f57aca527eb071d2c1786d08b64fc9adb68",
						"linux-arm64":  "40a4e3cd613100f6371f5ea754e0ad5f5713497d7a03ab3a286a81943de9e974",
					},
				}),
			}),
		},
		{
			LocatorWithResolverConfig: config.ToLocatorWithResolverConfig(config.LocatorWithResolverConfig{
				Locator: config.ToLocatorConfig(config.LocatorConfig{
					ID: "com.palantir.okgo:check-plugin:1.17.0",
					Checksums: map[string]string{
						"darwin-amd64": "b80cf6782d40b11f1776421c02ad55e4b4d4b5c675dd534b912981ecc9d38f6c",
						"darwin-arm64": "34a3a481f5a69fa7522f5c2152c83a88d288a6f248d1084b142a9ed53f07aa01",
						"linux-amd64":  "17a3e465e9165e56af963d1d09cab2f2cd81ea6f841bd97057a2946b3a78bf82",
						"linux-arm64":  "67a10b950912c2a3eb68a687613a84f6d6cf0c8653be4bef6ea2130e039ce6fa",
					},
				}),
			}),
			Assets: config.ToLocatorWithResolverConfigs([]config.LocatorWithResolverConfig{
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-compiles:compiles-asset:1.15.0",
						Checksums: map[string]string{
							"darwin-amd64": "a27529cfbede306286a7a92ed8f5444bab64734260422d4d8de8735d987dfed4",
							"darwin-arm64": "2ed1a3b5fdc1a28154cf3b71f2098a1f3ff9ab98ec9062647bf9821f09220a7f",
							"linux-amd64":  "dd30ce7c09da78771d9f1cba8293573ac6ccb2cd41bcbb691eff2fbfb2c3fce0",
							"linux-arm64":  "c6bb6ed3168f8dc4a3b3efc94c005a018246f68b1e913ad1e8c17d26557d5070",
						},
					}),
				},
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-deadcode:deadcode-asset:1.13.0",
						Checksums: map[string]string{
							"darwin-amd64": "dda3c41dd2b45d8a14c730f3c9038df7496eb1607e01790d7102553a2d7c4ed1",
							"darwin-arm64": "d8d1f80661aaf562293fd2492125365fb9af8a257d2a8f194d3a9e88366755ed",
							"linux-amd64":  "cb6f78bf027f2a1a07941bdd05b2b07419ffa40ffafe1aa61a6d8f9df07f3509",
							"linux-arm64":  "094b0854a03e68e8cbf1880456e2af72337d3d7704e7ac8076dbda074e79a9df",
						},
					}),
				},
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-errcheck:errcheck-asset:1.14.0",
						Checksums: map[string]string{
							"darwin-amd64": "9a7efd3c0f25b7b949c69c7dbd4b6bb545a53f1a9cb915d1b2ef535203b7a101",
							"darwin-arm64": "735baf11b19a9c3a383be65079cc9fb488c33332067db0854e4bb172cd006d84",
							"linux-amd64":  "b9a78b08145e6b09cec68abdab5acd2bfa975d0b8e24ea3a194ada32bd2f0fe0",
							"linux-arm64":  "abc98656da71473c420da1e2082ac85276393c8cd6bb6b71850fe226374dc52e",
						},
					}),
				},
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-golint:golint-asset:1.9.0",
						Checksums: map[string]string{
							"darwin-amd64": "470a5b54027420c41f8c62048d21ea3f2e030191c86f5e4cc4ce4580162fd5fb",
							"darwin-arm64": "7c4d3d9070de05dd486e6089e6982beb0fc82f6d200eb7b74c1e4a2ec7f73af4",
							"linux-amd64":  "59684cfbd3d9567eeb8286c3e6cd0ec54b05c0b2f9f3404e634bbbb2a322bb92",
							"linux-arm64":  "5762a6bf366eb41ef1307792c6c71136daf57986cff696a1d0949f01e7fc6e17",
						},
					}),
				},
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-govet:govet-asset:1.9.0",
						Checksums: map[string]string{
							"darwin-amd64": "e0c00216ee12f3ec9f9ae5204a1e6087d5fbf62c20c7030475c6f88c9c325e8f",
							"darwin-arm64": "53e12111172062f1550ea952f824afb602373acf470b4ea20f7985ed66b0f411",
							"linux-amd64":  "72491bab336d150d806fcabbfc8029f6d49106a631d14cd2f98d5315a0e6734d",
							"linux-arm64":  "93c92cbf378e9dbe9f957a0c1e8aa11dcd6d2593d9e5610a1371d09fd0b1ba50",
						},
					}),
				},
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-importalias:importalias-asset:1.9.0",
						Checksums: map[string]string{
							"darwin-amd64": "0683fce7451d789e00da205903fa3972a42a223c3316a860c72e5efa74d3e2c3",
							"darwin-arm64": "08a2c5cd55d3dc2bb3ae38701626d78ee01975d979bf5f2ed32d46e8c0ed3cf8",
							"linux-amd64":  "612a77f5caaeb6a9aa8475e2337b6e6598a06226d2e78f83d23a045c94f8ac71",
							"linux-arm64":  "17b6c074025ab54f0ab0198eeee4aa6cfa1729a528b658449472fe59db915727",
						},
					}),
				},
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-ineffassign:ineffassign-asset:1.11.0",
						Checksums: map[string]string{
							"darwin-amd64": "e849964f88b5160b79f5e1384bcd19c8895a4fc1e8aa0e89d6c68088253ca48a",
							"darwin-arm64": "5cb3c535520f1c41f16510f6f20d19dc6249a8b4afddaf08a1b47472daa5c6c4",
							"linux-amd64":  "6ff1f2b38b627e28d01807255a6af1349c3c6acca006d68c847c1e7e1bfd6c61",
							"linux-arm64":  "7196a589470f3908d35f35e40c62bf08a73f5349763da0b3c0fb0f7463afd507",
						},
					}),
				},
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-outparamcheck:outparamcheck-asset:1.14.0",
						Checksums: map[string]string{
							"darwin-amd64": "fe920812a8130b4d3fcad586c789154dd6e302f182de98f2d6aa66e266665d7b",
							"darwin-arm64": "45e65b8f640d531346fb7fee378e609676e0eee1cd5c970bfbdf3942e76fe6bb",
							"linux-amd64":  "c0192442dd3637a60c00f20e614dd292afcdc8986df360771e031e1ad4819bdd",
							"linux-arm64":  "af40c5738335ab727059814aa1a5950e5923b6fa46c469a660feca4ffc9089c9",
						},
					}),
				},
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-unconvert:unconvert-asset:1.14.0",
						Checksums: map[string]string{
							"darwin-amd64": "23930f79ce58e5ee3a6304e8b71507506c366c1a9dcbda23388fe3861e114195",
							"darwin-arm64": "6a2a4116d1e86bfd52faddeddbc774bbe85ff459f0058ddcb9e74b602141ab7d",
							"linux-amd64":  "8149b3fa820c8690ab51cf0a3ca7608d5932d681e86c976ddbd7391afec8c031",
							"linux-arm64":  "a81eafe48c2040eff475aafb2b781bd3f742f84f9a646fb78855169eb0b408d0",
						},
					}),
				},
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-varcheck:varcheck-asset:1.13.0",
						Checksums: map[string]string{
							"darwin-amd64": "bec7706c5bc47465d87472b8b8add4f6e1f06a49981894c6304dae511f1d85a4",
							"darwin-arm64": "9330fc9a64d402d5bc8b7937290bd61b09a506bba8dede717d6517ee9a5d8321",
							"linux-amd64":  "5cbf434c7c7f513ef5fae431d0936bdf8bb4e1547fed78dc94274128484079f2",
							"linux-arm64":  "fb0e3d51d719dc8d1f0daa20e424f79fcdb7d33496bc50e9420af066d013ed9c",
						},
					}),
				},
			}),
		},
		{
			LocatorWithResolverConfig: config.ToLocatorWithResolverConfig(config.LocatorWithResolverConfig{
				Locator: config.ToLocatorConfig(config.LocatorConfig{
					ID: "com.palantir.godel-license-plugin:license-plugin:1.10.0",
					Checksums: map[string]string{
						"darwin-amd64": "b8838b00417b3bdf9a04780bb8dac9f77de8d3d79288ad628c35a8f5801f0e6a",
						"darwin-arm64": "0b1113ad67b77ed7108593734b04c425fc68fda88580cba14ec04e72e8ea523d",
						"linux-amd64":  "ed26c20ecf6ab625d8e13ac7f8e7f65fc4c54bbed73f9cecd7c11a0adedce389",
						"linux-arm64":  "ea593c4dc026c5b3af0fec69fffa75244ba21f6ead866d8dbab6d092cfc444e4",
					},
				}),
			}),
		},
		{
			LocatorWithResolverConfig: config.ToLocatorWithResolverConfig(config.LocatorWithResolverConfig{
				Locator: config.ToLocatorConfig(config.LocatorConfig{
					ID: "com.palantir.godel-test-plugin:test-plugin:1.12.0",
					Checksums: map[string]string{
						"darwin-amd64": "b492291fcf5cf1b55de5f8f5ed141fbbc14bd30223fa506fa6adaab19a35c746",
						"darwin-arm64": "0e1004b04bdb024f0d6aa86b130d9bf1e153bc8ca12005e9eb3002aa6027c20f",
						"linux-amd64":  "db36ac64dea5652f721ec2c892574b14c360bb4aa2ba25321b912e1fd06337f2",
						"linux-arm64":  "b59b66f5faae55492fc801a6abfa1115f1277ee131797243e1535c9bf029d7e5",
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
