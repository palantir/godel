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
					ID: "com.palantir.distgo:dist-plugin:1.62.0",
					Checksums: map[string]string{
						"darwin-amd64": "d54d3bf32941dddf5150e8a86e52bbcc84e3254633d8736c62336189f5ed74a7",
						"darwin-arm64": "835c2dc36cc364c323c3df10ef1cea545b2e4248721302afda927fdea1eef260",
						"linux-amd64":  "cffbd63d4ce66f15a40220f58b85c0acf0aceca91d8c2a93088fc1b9611c447f",
						"linux-arm64":  "dcabb1eed25438a730c2ba2a0585b4bd248b8c97c9a3b535f48a43b28388fc61",
					},
				}),
			}),
		},
		{
			LocatorWithResolverConfig: config.ToLocatorWithResolverConfig(config.LocatorWithResolverConfig{
				Locator: config.ToLocatorConfig(config.LocatorConfig{
					ID: "com.palantir.godel-format-plugin:format-plugin:1.36.0",
					Checksums: map[string]string{
						"darwin-amd64": "6c6166daeafa129f75cb96150c7b91f652439b01d2fe94409c68cd0cc8abd15c",
						"darwin-arm64": "57e7c2b0c63434894d3da7b7a626508cc91b1f52db89bd36ba74f29a96865be5",
						"linux-amd64":  "ebff909f9e2cf450337487f6101945fd331569d86968eb139c42c1f479408815",
						"linux-arm64":  "5e9673a0117b80de37d09c653071d1436e545354a10a7456d7537c4205873eb2",
					},
				}),
			}),
			Assets: config.ToLocatorWithResolverConfigs([]config.LocatorWithResolverConfig{
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-format-asset-ptimports:ptimports-asset:1.35.0",
						Checksums: map[string]string{
							"darwin-amd64": "02b4df7534cc3b7357705cc62551fa11762e0ae2f3a0dc3bc3becca71180d898",
							"darwin-arm64": "329d968018a302f7b25ecedae51ecfe5f116990fae7c9b329c32bf6d17ffec9d",
							"linux-amd64":  "cfe76d907dbdd067820d00dc95d3336b2e22446f1bfaf817919c5342a84197a8",
							"linux-arm64":  "39cd502f0282bd5e1bafef79704bfbbf79e798f8a20988226a955298f2872bd1",
						},
					}),
				},
			}),
		},
		{
			LocatorWithResolverConfig: config.ToLocatorWithResolverConfig(config.LocatorWithResolverConfig{
				Locator: config.ToLocatorConfig(config.LocatorConfig{
					ID: "com.palantir.godel-goland-plugin:goland-plugin:1.32.0",
					Checksums: map[string]string{
						"darwin-amd64": "cb591dcfa23198c8c3ea0660ecbedd8cc09fb0c43e73c5940b8a73fc4228cce4",
						"darwin-arm64": "269b069cc00e446e84f92c8dcae7e25734baf78c981ad4e636eb7135cf4c9e3d",
						"linux-amd64":  "8ab9a21b1e209ffec930641f2e3f931d5143cc4bcd3735ab643859d9bcc93920",
						"linux-arm64":  "9c895312a61ba937c58b0fae94bf31453e712393989e721593d01437c5555c52",
					},
				}),
			}),
		},
		{
			LocatorWithResolverConfig: config.ToLocatorWithResolverConfig(config.LocatorWithResolverConfig{
				Locator: config.ToLocatorConfig(config.LocatorConfig{
					ID: "com.palantir.okgo:check-plugin:1.41.0",
					Checksums: map[string]string{
						"darwin-amd64": "f5fe94fd1ae450fb1db81f7849eb2a310c9f078ab8b0a00cbaab5fd76c4e3ce9",
						"darwin-arm64": "2afad9e31824620fe6951b50681065584aeff571cbc093d8fe8ddbf237fe70c9",
						"linux-amd64":  "182377cad7b730f4ffdc3a34d32eea72cbc4cae7209149314ccc19f27d401606",
						"linux-arm64":  "96c84efddb2af682f69715362e5df3d413792a1b222d6e03acbfbf7311c25a23",
					},
				}),
			}),
			Assets: config.ToLocatorWithResolverConfigs([]config.LocatorWithResolverConfig{
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-compiles:compiles-asset:1.40.0",
						Checksums: map[string]string{
							"darwin-amd64": "26d3e84750cf4888712f94085294d39eb2ff4ba26eb9ae196f2966611b064175",
							"darwin-arm64": "140e42ecdce73b0ee34da126d8aea49c58086ba0441b1c4704266cbf89c4cf71",
							"linux-amd64":  "627c5b8c51020094df9b8912212564788fabb466567aa5cfbd0892fe42a832b5",
							"linux-arm64":  "d2a01ad09e53bab6c69a253fec91f6ecdd9a895df6273883379f4ab69939b9ed",
						},
					}),
				},
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-deadcode:deadcode-asset:1.37.0",
						Checksums: map[string]string{
							"darwin-amd64": "8db3e1350025e3ffd6a5bee21219cc88e512b3cd9acb25d68db301b398afb036",
							"darwin-arm64": "c4931fd33ef48bb7b85ee64a65a0dc3444955b790a653970b0c3fbfdb9d7da10",
							"linux-amd64":  "aff4a16e7a3fa1c84c051f248f20be7ae59917298c13b1755763407d05cab9e0",
							"linux-arm64":  "88ee6fd9ddf678616ec1fa3bc9c7ed328cfe1bbd99b028436e48b9381905969f",
						},
					}),
				},
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-errcheck:errcheck-asset:1.38.0",
						Checksums: map[string]string{
							"darwin-amd64": "04600bce714e59ebbb68db4885654f37f4bb17f4ae2721d58a27173ba31d244f",
							"darwin-arm64": "0f4c7756ea60912da7f824e1561e10c8cb8ad169bbc8ebb67db9f7b35e6fa4c1",
							"linux-amd64":  "c47fbe54825e9616c44a71b59b20922cf0c64655c1c6735c50dd0955549b927a",
							"linux-arm64":  "81aa2280321085665b62d1b68720f6d651bedcea4d046cc7e81b6509b1d0af05",
						},
					}),
				},
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-golint:golint-asset:1.29.0",
						Checksums: map[string]string{
							"darwin-amd64": "9de0f221f92290735822c0ed7d8cd83f0fab1e7fd713cdb4c9f394d6ecfc7385",
							"darwin-arm64": "0967dc54a3cc4ac57e4fcd958b562cca4dcae89bc74674761af86e05e7424093",
							"linux-amd64":  "7987637f724be4384e18b1a923d2378df8ce8299b5739a0cf6fdc21f87472b5e",
							"linux-arm64":  "12f421427aa92a1812b1341d7bc6757b526ba117410fcbba6ba872e67ee05ded",
						},
					}),
				},
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-govet:govet-asset:1.33.0",
						Checksums: map[string]string{
							"darwin-amd64": "8c51ce914ec335c9da31a7694eedb51d5dc43ef7920539068572c32b368a78e8",
							"darwin-arm64": "c10521ba9393ddecddbf4368383af859a7af1e28d57a6fb34410fecd112c0d98",
							"linux-amd64":  "37f82ba2beb4ac09bf89788d16ffec809e715c3e4d71819e17bc61f7480b80bf",
							"linux-arm64":  "8ddd0c9044a2c7cc09c4bb53a908c36c063c52a42867612342c253d47c826806",
						},
					}),
				},
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-importalias:importalias-asset:1.32.0",
						Checksums: map[string]string{
							"darwin-amd64": "a53f2ccc12f842d1e47d03470429b091005069c56c13683a9e29570bb290afac",
							"darwin-arm64": "aa9520b42262763a8996d2eac0908eb684e3c3d73b7ec881484e7f4691d0885b",
							"linux-amd64":  "27a6bdcd4fda1b9ef6cb318d6f8d1f3491ff5dd129e97e2f2fec85a0d61bdb84",
							"linux-arm64":  "80d1f2767a3ed53cbd1826ef1b51448aadf050eb4d3e6231615fc618114950b3",
						},
					}),
				},
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-ineffassign:ineffassign-asset:1.35.0",
						Checksums: map[string]string{
							"darwin-amd64": "663598457d89adea9955f670cd9c392cc90035567d8aa438c6c3ac92d40e3880",
							"darwin-arm64": "648a9f3ea8e4b294cb5d8a390f76378dcc9b7ad918e56a3ef365fa96a5c88593",
							"linux-amd64":  "ab32449200fae64833a7114bbac40eb474bd8837b130aa3d44510b56322f7bc4",
							"linux-arm64":  "da121ee39befe7abd7be18213d1fe63b737f4747db7dca1c608e9a4636fd6674",
						},
					}),
				},
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-outparamcheck:outparamcheck-asset:1.37.0",
						Checksums: map[string]string{
							"darwin-amd64": "795a9fc0ccdc3284409c08a891d114409df84a8e4b4ac128297b5a35e5b07dd4",
							"darwin-arm64": "349bc910d69d289d7e47f57e6b0ab9ed393304731eb9e944d23039a973d752e1",
							"linux-amd64":  "ffbdef3f8f07537261d628ede229a4616dc859e55f3dcd19de5320f120506ec2",
							"linux-arm64":  "c39501b9032dc92d860f466f7969e9187297be915eddc5eca164db80c23ec17a",
						},
					}),
				},
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-unconvert:unconvert-asset:1.37.0",
						Checksums: map[string]string{
							"darwin-amd64": "36e3b278d5ac6ae7f3ff13b4440cda911df1be4fe173d6570eadd6cf65bba066",
							"darwin-arm64": "9110166e07cee143df0bcf2d9208fbfcaf4565539d765f18bfde2d3cd6cfe0cb",
							"linux-amd64":  "42219354a3a25bcfdb7ab607e16a0a01bfbe92822bec9c4327d7b7557e42622e",
							"linux-arm64":  "0ec080d2940d4cc14313d51775551f72f7eb7521ed06d24bb374733e7674609b",
						},
					}),
				},
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-varcheck:varcheck-asset:1.37.0",
						Checksums: map[string]string{
							"darwin-amd64": "0ea858951e4feda105075e9d47f74461e5d6acc34e1b968c2e50d70b2d0643c4",
							"darwin-arm64": "a0e31886537735bb0abf85249a719ebcadaf64fe66f69dac1993c860cfdd0d3c",
							"linux-amd64":  "aae8cdb9212f035a6d040f92ce25e0d9e0215ab6918c696e965539c5a2f09039",
							"linux-arm64":  "03b4b68ec5c2906cef2e3770a0edafe29efe409e1630f6c83175365944d59a16",
						},
					}),
				},
			}),
		},
		{
			LocatorWithResolverConfig: config.ToLocatorWithResolverConfig(config.LocatorWithResolverConfig{
				Locator: config.ToLocatorConfig(config.LocatorConfig{
					ID: "com.palantir.godel-license-plugin:license-plugin:1.35.0",
					Checksums: map[string]string{
						"darwin-amd64": "ca86c8def75d2363d1d40e889595ed8ee1d42125be5b62446faca87e5ba81cdf",
						"darwin-arm64": "d143f11871ee27f9fae24b594761d0598ca03ebbe12710effb6f84563499055f",
						"linux-amd64":  "ae0b90f74ca5bae8f0bbbe785979beaf439ad92d543bf33816abea646264fa58",
						"linux-arm64":  "67ad6794be9eed70fb75e6bf4f5041f46ccfde0c4973cf4b5fbbb1a8abd213f4",
					},
				}),
			}),
		},
		{
			LocatorWithResolverConfig: config.ToLocatorWithResolverConfig(config.LocatorWithResolverConfig{
				Locator: config.ToLocatorConfig(config.LocatorConfig{
					ID: "com.palantir.godel-test-plugin:test-plugin:1.34.0",
					Checksums: map[string]string{
						"darwin-amd64": "baa0d5527eef7a0b0a0c3b2ad35defbf67fdeb78ae1a98ae027643825fc208b9",
						"darwin-arm64": "8a3adcdca760747ce939ce39af72675873fdc79130182206ffc02bfaf92fdf98",
						"linux-amd64":  "730216b6fb50f53f06bac963cb96c6f8413f28a48b00905c42629e0d2a61c4ed",
						"linux-arm64":  "b7dba9e03b38dd1c6ff49b527e754c4bf1eae46e461d6b62aec7bf8ee78fd173",
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
