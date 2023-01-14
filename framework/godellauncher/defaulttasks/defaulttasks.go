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
					ID: "com.palantir.distgo:dist-plugin:1.48.0",
					Checksums: map[string]string{
						"darwin-amd64": "ecda8a7fe267fbf16924c4f1104a3d5022bb773d9d83359ad4d6ce1c3f4d5d2e",
						"darwin-arm64": "40b854820b2d834925f214509428b8ba9e35abeb37742fe0c1d70095e54eeaea",
						"linux-amd64":  "30597b0fd89680c232933b74b978b19383cdb6fb419b1a11b14a2659719290ae",
						"linux-arm64":  "210356287e2de55091a6cce9c61718dff671b41f6c0bbf92366d64e2e9d3b0de",
					},
				}),
			}),
		},
		{
			LocatorWithResolverConfig: config.ToLocatorWithResolverConfig(config.LocatorWithResolverConfig{
				Locator: config.ToLocatorConfig(config.LocatorConfig{
					ID: "com.palantir.godel-format-plugin:format-plugin:1.26.0",
					Checksums: map[string]string{
						"darwin-amd64": "27689e297e05617968659cfcc930e6e5ccd45f85adaf8ccf3399b41230a0f343",
						"darwin-arm64": "3aec315471258a92e5f71dfb7b9d40b39a5e0784ea00ca1825f4f11713f2ae19",
						"linux-amd64":  "b5f3ae6f77c9ff3f7fcedc885254b59fc3a1f33f199b0f4c174ea385feef8cd5",
						"linux-arm64":  "07f3447314711b268528c9ed85e337205c1cb1a90869d6e136c1d29b8264701e",
					},
				}),
			}),
			Assets: config.ToLocatorWithResolverConfigs([]config.LocatorWithResolverConfig{
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-format-asset-ptimports:ptimports-asset:1.25.0",
						Checksums: map[string]string{
							"darwin-amd64": "f399a948ec5218b40199c893c935e58e38b99869838c48b0db296b59af14f4f5",
							"darwin-arm64": "13da45fe6e4e8ccf811a12aade3093c86c614b04ecab66e694209c2f04fcbcb6",
							"linux-amd64":  "465db8a7d621dd61e74cf76d314b96d478eadbc0b511420e29ce70a17bd270c9",
							"linux-arm64":  "e4251215d7c8ec1e254dbfb7fffd409bbd8e580a7573d28abb38840956233a61",
						},
					}),
				},
			}),
		},
		{
			LocatorWithResolverConfig: config.ToLocatorWithResolverConfig(config.LocatorWithResolverConfig{
				Locator: config.ToLocatorConfig(config.LocatorConfig{
					ID: "com.palantir.godel-goland-plugin:goland-plugin:1.22.0",
					Checksums: map[string]string{
						"darwin-amd64": "8f593bb5980ad7172a43cd2474f3689bcd45b1d92727c3b647d6b2eaa2458cf3",
						"darwin-arm64": "6db17d9b7e1cf404b53c6221517d887a38732425aa7499e3d0da6bc7a40590ee",
						"linux-amd64":  "863fc0f1750081bff1a3d52d7a5bd7e105ce1bf29912d64458552c7d474c3e74",
						"linux-arm64":  "3d8e2263c0aa1d2b93a36b47b6691fa02c4d56f54d64c3515c30fee8345d5335",
					},
				}),
			}),
		},
		{
			LocatorWithResolverConfig: config.ToLocatorWithResolverConfig(config.LocatorWithResolverConfig{
				Locator: config.ToLocatorConfig(config.LocatorConfig{
					ID: "com.palantir.okgo:check-plugin:1.31.0",
					Checksums: map[string]string{
						"darwin-amd64": "f57a159c0c7b1e48ced2e2c1e96cafd81815ede6d49b81234647227e2c6caa19",
						"darwin-arm64": "7fcb573fc97d47305893f7aaff5cc809550c15f447ff276f494784916f719868",
						"linux-amd64":  "cf5eded38081c865a0a4ae4e6fe8d755f0c541b8ae981ca01b108d8d5f5cb382",
						"linux-arm64":  "a928573848687377953378a1bac512736fb96f19dc7dd491da4024d51b37b1e3",
					},
				}),
			}),
			Assets: config.ToLocatorWithResolverConfigs([]config.LocatorWithResolverConfig{
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-compiles:compiles-asset:1.29.0",
						Checksums: map[string]string{
							"darwin-amd64": "8d9f3ca7e138b9ee59ecb055318cd59d1c222e1978085b69b0e398c96be9e0d3",
							"darwin-arm64": "d1bd9c6395f8c6bc86fed7e6e6a5da634a4ab2d7956d674f841f2418c756775a",
							"linux-amd64":  "f70b8852e327f6d540a0ed1de63f9fa940c35178f9030f9f8ed6932d49db09e3",
							"linux-arm64":  "032264f3d169bc863acd03ef60ee5d83750b5b2ff4f6061072162c7d9f9e1dd5",
						},
					}),
				},
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-deadcode:deadcode-asset:1.27.0",
						Checksums: map[string]string{
							"darwin-amd64": "d1b0562c299776e34fc1f687bcbfceb1d883f68a2856935940b3081d5b8fdbbc",
							"darwin-arm64": "9980cd1338e4bb21a11bb79469590a8958e1c4093e2b29f65356d53f9aaaef58",
							"linux-amd64":  "ae0fc0576c4a53b0ade6480d39b8a3d6f023aeec1b1ab07a686696999bfd16f4",
							"linux-arm64":  "5c459eb0198666443ae9cad4fbec9913c91ad2774b0c927494f5219deff875f2",
						},
					}),
				},
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-errcheck:errcheck-asset:1.28.0",
						Checksums: map[string]string{
							"darwin-amd64": "248773205e4e8bbfd439d98e139252a4c4ef9495102d975a4a39678d96fb43f9",
							"darwin-arm64": "bfc3c18fb438fdf1e4b95af405508f8d5850119c6f6513357c6356f9b2c07169",
							"linux-amd64":  "2457940d98819eb6a511a6ec0f4d02bd6bfe8596f43d7194003682d49ea02167",
							"linux-arm64":  "1d35a48f3aa5399737374d1b7fd185e75a52710603ea7c144a4c76ca723d9553",
						},
					}),
				},
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-golint:golint-asset:1.19.0",
						Checksums: map[string]string{
							"darwin-amd64": "044ed2dec194a0290bb9fb347b555665c3d1d772b90431b52f274cd224023167",
							"darwin-arm64": "bcb7cfd6e486862cfffedc5a8be10bea4ae4e4be53ea1e4755332ab567499ca3",
							"linux-amd64":  "c577919ad85f09a3492b52e145c516bb6257c13a184b44e9e2c91993d9559fd8",
							"linux-arm64":  "77d11706fb47b91a640482007abe9e805c7c79c36a75245abe6acd92da78e801",
						},
					}),
				},
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-govet:govet-asset:1.23.0",
						Checksums: map[string]string{
							"darwin-amd64": "0bf95dd18b249dece6eba869cba4f4466f4689d3ef97eaf3ce30f6a754efc874",
							"darwin-arm64": "27f89fd333105f41b3834e3b9089532ba6085ca3fd1b1e356dc99eeaacd2ba59",
							"linux-amd64":  "3b5356b8c75f5608d70715bd2268297f7107e382b3a61de6a722ff39bd288a50",
							"linux-arm64":  "f8b7b35e71cbb38a9bf74a82d23d04476b4da8bfe0f8c8eed099de8d66e1b934",
						},
					}),
				},
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-importalias:importalias-asset:1.23.0",
						Checksums: map[string]string{
							"darwin-amd64": "9f3e04ff2e16666264ecefdb2bdd9f78237d2ee7fd8a7aa757bb317399d684f5",
							"darwin-arm64": "a5e8b0ac98a1d72dcd7640d560a8badc683819860c001b2f1cc93e9221a81ed8",
							"linux-amd64":  "413d29707bda97518ad2923bbfc0416322110b3040e24e8cab9da9f7017a1135",
							"linux-arm64":  "b45015538f452ec8f020173809b4b8bd9893e0c4b2381ed2cd0701333d757bdb",
						},
					}),
				},
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-ineffassign:ineffassign-asset:1.25.0",
						Checksums: map[string]string{
							"darwin-amd64": "02a6b0911a85a4e0e3b9774460ca4adcc48da6620e68bd43749cecfd911f920d",
							"darwin-arm64": "646814e9d29fee8173c2ecf736e7512f4c0479b969fd86114b681cf2541d607f",
							"linux-amd64":  "b605bd2da69f417a9349040c38bcb88106460a58e8f53327c640a6353a56cfcd",
							"linux-arm64":  "82ed72f8990f0f236dfab7a2877cc84cffa4649f87cc2688ef0de8c2ca744990",
						},
					}),
				},
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-outparamcheck:outparamcheck-asset:1.26.0",
						Checksums: map[string]string{
							"darwin-amd64": "a0ee405152af31ef3c3fb6f71bba72e9df15be2c0a2fe215beabb1ed58c152dc",
							"darwin-arm64": "05c8475a64d418e8851a26b0f03d308a66368b21ad3df17261db449228e7794b",
							"linux-amd64":  "c9329efda750382195ec2663eac083211dc37cca36beb1e0edb545c4f8faa676",
							"linux-arm64":  "71acdbf31fa3005d79cb3515ddc893affedbd0ac14fd732548ca4d1c1ddf919f",
						},
					}),
				},
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-unconvert:unconvert-asset:1.27.0",
						Checksums: map[string]string{
							"darwin-amd64": "959e0b51d52774240bd495dc64aef8c7de3b83cef05c29274f09b7512f9e5af3",
							"darwin-arm64": "4df1464dfe6398110ef6ad11359057ada2d13def86a714da5192c03d8c82e3d8",
							"linux-amd64":  "95217737ddd9efa2d9184a9f74ee7df6e2c628a26ee949560e0815bb67d1e1ec",
							"linux-arm64":  "3aedcdfbe99a5b3ea3a0e180db4b19e27e4bdd2f201ba5b45233ab1890808c6d",
						},
					}),
				},
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-varcheck:varcheck-asset:1.27.0",
						Checksums: map[string]string{
							"darwin-amd64": "885454c8b242da54a158148ac32d1a2bc0654304f9ac80cbf9936cb061b1b2d4",
							"darwin-arm64": "d4bd1e28cd7751180a23192c1bd239df874b19d9b9da0e875b4c117c243c3ada",
							"linux-amd64":  "14038403717afb4394b087f6e08b6dd4e68e80062742355193794ef0eb9f29d9",
							"linux-arm64":  "f98e0924282377fad81146310e27d04d6ad82b85e28aa9e8b73d73324e1490e2",
						},
					}),
				},
			}),
		},
		{
			LocatorWithResolverConfig: config.ToLocatorWithResolverConfig(config.LocatorWithResolverConfig{
				Locator: config.ToLocatorConfig(config.LocatorConfig{
					ID: "com.palantir.godel-license-plugin:license-plugin:1.25.0",
					Checksums: map[string]string{
						"darwin-amd64": "fc7fb391984eec1777aeb4a2f84920361c5f8a60f9182f702c3505593482a66a",
						"darwin-arm64": "5ef691019dd3979741045b92917b8c7fadf259e7222c7b3f3516b6167590f8bb",
						"linux-amd64":  "de5654265709ad5ae744ca57e1428f64f37a6522bc7efefd3b5b4bb9c5245703",
						"linux-arm64":  "d0d2a27fa25f806dad6045af8d7527a8bbb4d74dc3f397dbb16542cf0895de44",
					},
				}),
			}),
		},
		{
			LocatorWithResolverConfig: config.ToLocatorWithResolverConfig(config.LocatorWithResolverConfig{
				Locator: config.ToLocatorConfig(config.LocatorConfig{
					ID: "com.palantir.godel-test-plugin:test-plugin:1.24.0",
					Checksums: map[string]string{
						"darwin-amd64": "467d91fe1637b61286f1f10b5e23b0b4332c6dd95e86384e1693a8dc8b1b1989",
						"darwin-arm64": "fa26cb9d73b4899218a5ed4f5232556a3cc37a93b6e8fcb036333c04d492cb8d",
						"linux-amd64":  "33a5932ab063af73af95cb23be8f31e26313797305a7650df35f505bd7751b3c",
						"linux-arm64":  "df0a2fd3184005a757768ef289b5d81d7c9dd491a92098dc477def7a3960eb91",
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
