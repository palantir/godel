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
					ID: "com.palantir.distgo:dist-plugin:1.58.0",
					Checksums: map[string]string{
						"darwin-amd64": "0b9628dfb9784e722fe84ae70d8bbdbd4344731246d85233f66dfb24c872fe74",
						"darwin-arm64": "f608bcc80c8b253850a6fcb32f5abcea3bb176cfa8a95b74ecf896e45e68e999",
						"linux-amd64":  "1f79b1033cb07392fb31cf35f38c38e8d0f4967f3684f2dcb2447a8be9a08b0c",
						"linux-arm64":  "27a485b03583279b54e9af0b23f4c46854fe89826484ed4d7367338f15389732",
					},
				}),
			}),
		},
		{
			LocatorWithResolverConfig: config.ToLocatorWithResolverConfig(config.LocatorWithResolverConfig{
				Locator: config.ToLocatorConfig(config.LocatorConfig{
					ID: "com.palantir.godel-format-plugin:format-plugin:1.32.0",
					Checksums: map[string]string{
						"darwin-amd64": "ea3245cea0158407d8b55860e9a4e6f93469ae497e86091cb3c7d3a25891360b",
						"darwin-arm64": "a951afa5c8d28d60a817b429c34f43f1a194e3b3ba8ca3a4dd3342cfbba34f1a",
						"linux-amd64":  "acd388c61b886dfbc743abca70214b3ea40340ce939bc0da2b48a598d1e7151d",
						"linux-arm64":  "2d61251c053f32dd25484bd2a610aa68e80737508b97f55ce4c71d210ad440d6",
					},
				}),
			}),
			Assets: config.ToLocatorWithResolverConfigs([]config.LocatorWithResolverConfig{
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-format-asset-ptimports:ptimports-asset:1.31.0",
						Checksums: map[string]string{
							"darwin-amd64": "e5d0c0070ced6dcdb23d7a31695e2be97d12d1bf985f14c7ac32907cebb2ad00",
							"darwin-arm64": "fa5b46ba4575f1d386aff0af3a527baef944f9d775d038014f82798994aa877e",
							"linux-amd64":  "f11c716b212f7e6a3da1eb0c834e3be7baa4de985c2cbcc37513501c8b45c9dd",
							"linux-arm64":  "5b148ccd2104156d91b95ea8a9ef8b2994ad0f292a561acfef766454302bb4df",
						},
					}),
				},
			}),
		},
		{
			LocatorWithResolverConfig: config.ToLocatorWithResolverConfig(config.LocatorWithResolverConfig{
				Locator: config.ToLocatorConfig(config.LocatorConfig{
					ID: "com.palantir.godel-goland-plugin:goland-plugin:1.28.0",
					Checksums: map[string]string{
						"darwin-amd64": "82030cd37e6f3938b2004bc5babd06503ab8f9fae9c14b0080bba8460db8023f",
						"darwin-arm64": "767ec8c98fa7e20cb811465c392a796716b5847bc458703516ab582a69036e25",
						"linux-amd64":  "a367de0e396c17b6f0f79b48dad0d982ccd50b4e2994a05f9142771dac9e4492",
						"linux-arm64":  "57404e113eb524284f64b4ac2a43364141841e5c3ae7281ef6f22e7a86a3eaa7",
					},
				}),
			}),
		},
		{
			LocatorWithResolverConfig: config.ToLocatorWithResolverConfig(config.LocatorWithResolverConfig{
				Locator: config.ToLocatorConfig(config.LocatorConfig{
					ID: "com.palantir.okgo:check-plugin:1.37.0",
					Checksums: map[string]string{
						"darwin-amd64": "237cb0d84c8ec0bead697a1870301220cf8f05e3d53c52a8a967a6a0549110be",
						"darwin-arm64": "17c113e8d47a07974cddc9c945c53df78255fa020c5688334adbb9f68950d17e",
						"linux-amd64":  "1a65d3637d2986d67f6dd325e847b13c65ff87d56da31345baa81e608f001497",
						"linux-arm64":  "61d510df33662ea2a4b63329adf45afe6b792f37ea8b89dd650155ce1c901dc3",
					},
				}),
			}),
			Assets: config.ToLocatorWithResolverConfigs([]config.LocatorWithResolverConfig{
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-compiles:compiles-asset:1.36.0",
						Checksums: map[string]string{
							"darwin-amd64": "270531f74a36429f1806a5514c7f6662a5f64f2e491b7763e9b42be319d67076",
							"darwin-arm64": "a3da7131bee5b3b33ac6717ff2a95ce715e69e574920115b90f63a266aed7ca8",
							"linux-amd64":  "9c4cf6bf61f54786e417b18141bfe04f057c96e7ec41836c9d8ad73cec5bbf3c",
							"linux-arm64":  "a321c6d20ed26b35fd351b654a631eef44d508de6a2b37c093cd799af9b077b9",
						},
					}),
				},
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-deadcode:deadcode-asset:1.33.0",
						Checksums: map[string]string{
							"darwin-amd64": "c296f355ad1e6bab37e00ec0adda73d28beaa1fd531af4004b71af9659fab961",
							"darwin-arm64": "4f0f42d0aac2fe62487bae7090c2e7f4084dd7ba652147aebf293ffe5ee88c5a",
							"linux-amd64":  "c1009d69d2202c52dc10d3ec14db730e3bfa78bc21e5dca7a33616e22293c0fa",
							"linux-arm64":  "447f1d17b5493c17c21aff10b6cf65a10be54adf7f441b551fc9fd4d423322c6",
						},
					}),
				},
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-errcheck:errcheck-asset:1.34.0",
						Checksums: map[string]string{
							"darwin-amd64": "2283d569a3b0f27d156698bd19fd36620e2e4121f3d185a3d7a01b1543084fa4",
							"darwin-arm64": "0ce654249fdc3bae3b87d98b1a8fdc3a4ed9e6f83b0cad1192dc3ade2894c20f",
							"linux-amd64":  "b10a36f091205e53760bb308ded6f30096f5857783b1c3a73344aca2decb107c",
							"linux-arm64":  "ec6c96ff5ca912577bdb159bd2ace3dec27be66cf0f66bd5ae92157f553d3e16",
						},
					}),
				},
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-golint:golint-asset:1.25.0",
						Checksums: map[string]string{
							"darwin-amd64": "003ca59f23185bbf90abb14d06bb74ad7828f047b87c44c9878f4fffc074e23f",
							"darwin-arm64": "2856541f340d664241e8b04c57566034443391ba780bb25e6929f6de640c54a8",
							"linux-amd64":  "3b376b24802ee0363c57bc1c88147c0cf88472e775dc1d9996890acef08ecc7d",
							"linux-arm64":  "504032909cb007e823b132fa34ed23dd3ac619fc6c4e35674254ddda1d4622f1",
						},
					}),
				},
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-govet:govet-asset:1.29.0",
						Checksums: map[string]string{
							"darwin-amd64": "aa973a8051deb70bb022d3facba8a3702b48e0a83f5d062ff7835abb4d0c38a4",
							"darwin-arm64": "8e1e18fdbb0f2859737b8b4edf795ab8d29b7df82fe39c864b6ea79b4ce488ca",
							"linux-amd64":  "31a86fe02412c00e68eb25c6c8624b50ec7def66ab7e707b8b46c82970fce465",
							"linux-arm64":  "f688dcaf88bb8b21a653fc3c2c5ac808b6902778aa122c04fec1b482400f68ec",
						},
					}),
				},
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-importalias:importalias-asset:1.28.0",
						Checksums: map[string]string{
							"darwin-amd64": "2f77edf72fe337440e1d0bf2a71dbfdce973656ce7e1f40fbbba5e7dc67057b2",
							"darwin-arm64": "1db9fb6dfd3fef20d9db058fa05a26a3c0e98ad7f9ac7685fe09dddfa736c21f",
							"linux-amd64":  "9233d7cd047e2090a13d90bd1179965729647882306e054c83a4fa9a64353320",
							"linux-arm64":  "2d50c7f628bc27d97bcdb2cddd81b4cf8ae744254629ab43621748a11c9fe89d",
						},
					}),
				},
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-ineffassign:ineffassign-asset:1.31.0",
						Checksums: map[string]string{
							"darwin-amd64": "62aca3a6fee773bb1533e99deac5da3c842433dee1d9001fc1625c05660509bf",
							"darwin-arm64": "475898cea4a450e6dddc3481042f5f5e374b07bca629607373bd3e61f76394cd",
							"linux-amd64":  "8cfaf5dafb742508cfaa713fc7732f9508e44dc2ce9cf1e7d50fb75ba9849d07",
							"linux-arm64":  "90f69150df5e20af5a34ec74ae9b1038f223cc87b3f6d137b50554740112dbcc",
						},
					}),
				},
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-outparamcheck:outparamcheck-asset:1.33.0",
						Checksums: map[string]string{
							"darwin-amd64": "d20eeec09c8324b4fef98d8608ac13415cc9b040ed794c1c775512b71c5024b8",
							"darwin-arm64": "80f9668a3e86648c2a5a44f8cd90f9c309f911f3ecf457de077680896528b350",
							"linux-amd64":  "3572313b590e79aadd4be82886ea9b62feae5e84206e51c4531d6c6e6c4edfd0",
							"linux-arm64":  "3f74584dba834d549f402f9e43bdce2e71368f5ce7f8c0ff6fec9bfe02416816",
						},
					}),
				},
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-unconvert:unconvert-asset:1.33.0",
						Checksums: map[string]string{
							"darwin-amd64": "ce7875349cd9f3c9a22244857e2f7b91159f48e7bb1b509842c8f3711a1e4b44",
							"darwin-arm64": "b6c40ce5088f8461f44e4adbfcb58f5fc7af2339b3b42f04778caa8ac98a5467",
							"linux-amd64":  "1b9f4f751d971d5e93378647f2af9f2342f871f479d11f3cd491846a0bda3a80",
							"linux-arm64":  "5c7232e4c2a049f7f38dcc3a52f35542d74e505579c3b4b86c1a1e7c9561461c",
						},
					}),
				},
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-varcheck:varcheck-asset:1.33.0",
						Checksums: map[string]string{
							"darwin-amd64": "f135aa4c27bf9fb8607473335a3b8c71a33e50a531d5df8376cd71ff95d6e3d7",
							"darwin-arm64": "42b474b4de6944984e3bf9a45ae0f5b46feee7c32123a09ea4654d124850a9a1",
							"linux-amd64":  "c8aeedbf2f4c62f5ab3e70ab241025a7c96e10c305e36b8fd56af44f99b1c3a9",
							"linux-arm64":  "68ad88012a32a3820e699e306c8c7ea505cdd6dd52e717bd4022193f08d452db",
						},
					}),
				},
			}),
		},
		{
			LocatorWithResolverConfig: config.ToLocatorWithResolverConfig(config.LocatorWithResolverConfig{
				Locator: config.ToLocatorConfig(config.LocatorConfig{
					ID: "com.palantir.godel-license-plugin:license-plugin:1.31.0",
					Checksums: map[string]string{
						"darwin-amd64": "4c08fd21a5816277bb69cd4f9403614caf1f5f082c71be39f855e56a0f301ae6",
						"darwin-arm64": "6828134d820d33aa9898f0293f384a6b4291469917202c0adf88d189f203c821",
						"linux-amd64":  "84f3ccff60b7dd4f4588a1227f1412613f278e7abdc795961b9309432833e937",
						"linux-arm64":  "18c08a2801971f12fba53c72b3436cf7af1ec164fc94fc1087fe508600e08843",
					},
				}),
			}),
		},
		{
			LocatorWithResolverConfig: config.ToLocatorWithResolverConfig(config.LocatorWithResolverConfig{
				Locator: config.ToLocatorConfig(config.LocatorConfig{
					ID: "com.palantir.godel-test-plugin:test-plugin:1.30.0",
					Checksums: map[string]string{
						"darwin-amd64": "0ba3cf62ea7bb500fc6085cf2af87a8c0d8abd893722fa0ac3a322334c5ea0f2",
						"darwin-arm64": "489426e18288d5d13c1900e034615d144a5a2dfa866e3b96bb6fa8511be8d878",
						"linux-amd64":  "d5615566b4dd29cbf401f135e0733e526a9840e9a8041dbb8a51e33e85dfc5de",
						"linux-arm64":  "7af030b6c1d43d1aa8f458aaeb203c27c9b93181a8809469425d706ba1dcca39",
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
