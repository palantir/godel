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
					ID: "com.palantir.distgo:dist-plugin:1.74.0",
					Checksums: map[string]string{
						"darwin-amd64": "2266b6bc8b215784ecc688460efb7afd30dfdb7848c800ac2ce4e3127968721a",
						"darwin-arm64": "42fbcd47c4694f5800373baa59873b96362637e1e5b106582c0fa11b9ab83c1f",
						"linux-amd64":  "af2a232915ae95470013d11768c6f42973afef20ee12b38056fa8d1f978e342c",
						"linux-arm64":  "019004657ed33cea6f6e5de47c36e2940dd6b8ea50ec5f21a16b09b3c71f605b",
					},
				}),
			}),
		},
		{
			LocatorWithResolverConfig: config.ToLocatorWithResolverConfig(config.LocatorWithResolverConfig{
				Locator: config.ToLocatorConfig(config.LocatorConfig{
					ID: "com.palantir.godel-format-plugin:format-plugin:1.45.0",
					Checksums: map[string]string{
						"darwin-amd64": "88f79932969417dcf283a9f0bcb98d497f02d243440652616cc754aa8efac112",
						"darwin-arm64": "eaab96481f950c6f085291edd40ab4ccc432a853c8d6017d02fe904bd89e4d72",
						"linux-amd64":  "65e87dd9e6c2fc2958ebcbc8f74a7e0b7f3e616b6b23b2a500020734e3e33e2e",
						"linux-arm64":  "168a41118d29295f14259db235db72d9938b0cb8c0d3c3a16dd37acf139c7803",
					},
				}),
			}),
			Assets: config.ToLocatorWithResolverConfigs([]config.LocatorWithResolverConfig{
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-format-asset-ptimports:ptimports-asset:1.44.0",
						Checksums: map[string]string{
							"darwin-amd64": "f19fd767d1fe09a6562e797c4b2e083d5b58a9c5b0d52842e8349f053e1ecf01",
							"darwin-arm64": "a5a1131b55d1aa9fffc5e4cd62edce8daab175d2ef83bafde357a0448ea41529",
							"linux-amd64":  "1e2f3814c64c62c6150c4af9c7ea57631799d737e6954663d8cc3e4e8111984e",
							"linux-arm64":  "2a785c016dcba98181be38bcd15474404b8f187776f32c3b6c9e7cc97e35c8b7",
						},
					}),
				},
			}),
		},
		{
			LocatorWithResolverConfig: config.ToLocatorWithResolverConfig(config.LocatorWithResolverConfig{
				Locator: config.ToLocatorConfig(config.LocatorConfig{
					ID: "com.palantir.godel-goland-plugin:goland-plugin:1.41.0",
					Checksums: map[string]string{
						"darwin-amd64": "2d19f895bf76934d186483ec3e80e6f536e99442f87ffae7d634ce162b782df8",
						"darwin-arm64": "a027b4a159e2ff76595eb21488404d6e85bf1c9faf4cfdbbd4c6f8769baaa899",
						"linux-amd64":  "926cfba60749301854cc70961dd91661430de0060a1a16d71c6a24328e7ecce3",
						"linux-arm64":  "65642339a22e3b24390bbad684b7d925b3c97c4291f5dc9395a6093893cdd1c2",
					},
				}),
			}),
		},
		{
			LocatorWithResolverConfig: config.ToLocatorWithResolverConfig(config.LocatorWithResolverConfig{
				Locator: config.ToLocatorConfig(config.LocatorConfig{
					ID: "com.palantir.okgo:check-plugin:1.57.0",
					Checksums: map[string]string{
						"darwin-amd64": "2177e105c18fb44e3824880dffc037048b9cc1ebe9368c756f65671371de49cf",
						"darwin-arm64": "4a97176d10672370416defe00219ad444bd94ea044066b0d8e2e57db927959de",
						"linux-amd64":  "4710feb8412e37c310ac7fff9993ccbe9b36d4a6a34caf767d9483efceb9e207",
						"linux-arm64":  "b4163d42be4e88ea9b25cf3ec5a4aa76c3c3697a2967dcb0624cba0b895a2152",
					},
				}),
			}),
			Assets: config.ToLocatorWithResolverConfigs([]config.LocatorWithResolverConfig{
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-compiles:compiles-asset:1.50.0",
						Checksums: map[string]string{
							"darwin-amd64": "cefc635e61885d664df7b13dbffaa80c9cc9a89a5040fa85e2c484bde90f7119",
							"darwin-arm64": "7523423a4fb9abd36cc8518d4224c74e083125f3cba1e618983c8ec2d6b82f21",
							"linux-amd64":  "2cf7045c6eeb94db38b0448753fcee4b37b07e6f61c83f8af43b6dec8e5e32bc",
							"linux-arm64":  "e065a5e95a28e6e8c608b6aa2a4ae3c6f5e23326a150de7e9e224e59f593b836",
						},
					}),
				},
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-deadcode:deadcode-asset:1.45.0",
						Checksums: map[string]string{
							"darwin-amd64": "d800ad95d45aa74d8137d7c78c649b7b9fe203f89c798b69e06380537c07be5b",
							"darwin-arm64": "c0ad7e7dd0437c82b98837305e7259dfd6db778260f6c8449bb548d6d90dc53a",
							"linux-amd64":  "41d1c816de8539bb76b76fd4a4f925469123d7fad8ce3ea4df46eea4259fa596",
							"linux-arm64":  "44586bcb2e6d102c6fe378bb0f6de2505e345811ca537b9c4d5e9b242c6416e6",
						},
					}),
				},
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-errcheck:errcheck-asset:1.47.0",
						Checksums: map[string]string{
							"darwin-amd64": "71d6999d7a90a01f2c7b6107fe7fd838a0606e9ab7fb5cbe5a62198c13f374d0",
							"darwin-arm64": "b0b300d5c7a578f13febbec5a6109aaf39544d306d05197a5a817ebec09eadba",
							"linux-amd64":  "20deb54278ab31d3ed442b0641fdc8f7e0318b27a4709e2419c9c54c2b8d0818",
							"linux-arm64":  "7c1d2c18cb574ade0e96998aeb69fcb5a9594be1cfee6a8b5b287b8c919d4ad6",
						},
					}),
				},
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-golint:golint-asset:1.37.0",
						Checksums: map[string]string{
							"darwin-amd64": "7e179bb625a2dc9aa0146a7d67bfa6c97073ad1123fd9d3afafa3f9a5473df97",
							"darwin-arm64": "87ff384bd0b8b6fa024ea53c139fbd6efc3034cfe39835fe26b5cce683696c98",
							"linux-amd64":  "454405768d2e33c1a9982bf6d27c3f09612fed2c44797830ec8ce0db26813dff",
							"linux-arm64":  "4c81969282accb04f7656e1caa866a51fdb904aa73a689fcd8af3248c48411e3",
						},
					}),
				},
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-govet:govet-asset:1.42.0",
						Checksums: map[string]string{
							"darwin-amd64": "4dee1e72eadf3fefaf1d3ae8e329409a24bb4ad900d03e30821a02342c18f536",
							"darwin-arm64": "68241a33e17072dab1865fb80b82db8aa2f216a8b01f5afb3c0b0e3293989acb",
							"linux-amd64":  "771b5fd84eba2ea808312e342ee738d950b2408306b9b1bc5655207ed9d0e6a0",
							"linux-arm64":  "3117e7ec26abc45b2f6d97cbe2961eeceae4d2f2f2ebda6f7cf6bba9a3f6268b",
						},
					}),
				},
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-importalias:importalias-asset:1.41.0",
						Checksums: map[string]string{
							"darwin-amd64": "6434938bcc49d9887ca1e5b992eb36a07739319c026b720ecd1d2d9f53d323b9",
							"darwin-arm64": "d79cbf1eb482f4f7b603f8a3ea812efaed3dc3f55d8ae8cca20443735b65661a",
							"linux-amd64":  "730ce0807b9216b3e32e900da1019e6377df67815e1413197e842d36f9bba32f",
							"linux-arm64":  "672311af6027883fad84ad0d5f4fd554d04ebe0bedf2c931050d5789799e3657",
						},
					}),
				},
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-ineffassign:ineffassign-asset:1.43.0",
						Checksums: map[string]string{
							"darwin-amd64": "e12662997f3d6bc54144eb77ca8b31c5e97ace753e8f34b1c6c7804b8dc3ca56",
							"darwin-arm64": "946e970b8cd1b60f02eeacc49847ec3317990a161d07b9d927ae42549ade28d8",
							"linux-amd64":  "7640e562ee35ad65929522386ab7245d1a51cd68a75b2b4075a3fab992d34113",
							"linux-arm64":  "0edcbf537dfd5a2dd7478ee3cf439e50a8dc0512aaf503af976b47446f4458d9",
						},
					}),
				},
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-outparamcheck:outparamcheck-asset:1.45.0",
						Checksums: map[string]string{
							"darwin-amd64": "0c90072217bf31ccab8972fe0684e8e01c458673075d6b18be1fc157550b2838",
							"darwin-arm64": "7b7f84628805e0585da79f88511fa1f054e622e5b5de956cd068ae1cb6b1f2d1",
							"linux-amd64":  "17adcd7698e90d7fe3063e3c9243586abe470e29c7a9495c4ba44e3a0117426f",
							"linux-arm64":  "6f988514015cfac7be674f88407867e8d970c7ec1d4446c06ca8554bfb4f42b2",
						},
					}),
				},
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-unconvert:unconvert-asset:1.45.0",
						Checksums: map[string]string{
							"darwin-amd64": "c06116cfd0b15bc10eef4682557740eaa800e73acf632982b6b1ea22864db813",
							"darwin-arm64": "4c34c468ebc656c268adf8a5095885b89ee9241cc6f1d96a8dced59ea12f247c",
							"linux-amd64":  "986e7aece23a2b0d3f81818b024608ada3352425fa261bbf0d2804d2a9ce5e6f",
							"linux-arm64":  "dc3921e6a83cd765ca515a3e160231383dc8e4aa0136a1db76b647aeccd5fd49",
						},
					}),
				},
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-varcheck:varcheck-asset:1.45.0",
						Checksums: map[string]string{
							"darwin-amd64": "da77d28095d5355efc80701489afe7388248069fb00a9574a5a066aaf741ef62",
							"darwin-arm64": "45a4bcb158284ad6fdacc6d0acdbc530b579e54755bb264544b1014ea9621ca4",
							"linux-amd64":  "abcb7cd750debcea5230307b0935cb850a6308f4ae316df006ac9dc2bd828713",
							"linux-arm64":  "ce010e1cc472efe891228457e0d086ea8ad26a6eb2b5b1fbb3f660ca5ccc25cc",
						},
					}),
				},
			}),
		},
		{
			LocatorWithResolverConfig: config.ToLocatorWithResolverConfig(config.LocatorWithResolverConfig{
				Locator: config.ToLocatorConfig(config.LocatorConfig{
					ID: "com.palantir.godel-license-plugin:license-plugin:1.44.0",
					Checksums: map[string]string{
						"darwin-amd64": "f720037145ef4ef9610903f65dad6c942a6189d69acacdd14a5ed6c5a0ec9f59",
						"darwin-arm64": "52509f76f6b72d2d4448d2fb7e459f7214902ec2715c95f8d32facc70d78bbfb",
						"linux-amd64":  "c09b4bdc5b4f1edd5e372d725cc490fc014ee200d4f7c209a5a4cda75e1f045b",
						"linux-arm64":  "5979516512609076957892cd37eff61705faaa8ced5a708954cb351c72f56971",
					},
				}),
			}),
		},
		{
			LocatorWithResolverConfig: config.ToLocatorWithResolverConfig(config.LocatorWithResolverConfig{
				Locator: config.ToLocatorConfig(config.LocatorConfig{
					ID: "com.palantir.godel-test-plugin:test-plugin:1.42.0",
					Checksums: map[string]string{
						"darwin-amd64": "cc7c3294e6dfd32f4f116f1f91b560dadb599b20e4a879ee4e73fa1bdbf3ca5e",
						"darwin-arm64": "0d373b10743095e90a52fa17ebe3f35cd79ca24709fca1f599eb93fe1b925912",
						"linux-amd64":  "f1be4fbfb707a61a7a6fad6cea0cec00db76c358e08851c25758fd115c030d47",
						"linux-arm64":  "512cde17fbba7070bd1dacc01a52ed72c0495c46c765695d53c7a4a0298b5e42",
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
