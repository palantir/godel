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
					ID: "com.palantir.distgo:dist-plugin:1.69.0",
					Checksums: map[string]string{
						"darwin-amd64": "bcb9e268c5afc352a40b38bb4734a35c2c0b5518c70b8238d049bfff4413e37e",
						"darwin-arm64": "6938178d9307c7d7cef79b1ae52c6ef1eba575de4e5589e8a4a36b3824f8a0b0",
						"linux-amd64":  "cc9f224c0f4160bce3fb91a14af7f9bd214a6551b6e53fc3f6c9779fab2eee43",
						"linux-arm64":  "39e55226c3b19a3f69b9d39830e7d47352f8491e33851fc2e94de2aaf4a4cb3a",
					},
				}),
			}),
		},
		{
			LocatorWithResolverConfig: config.ToLocatorWithResolverConfig(config.LocatorWithResolverConfig{
				Locator: config.ToLocatorConfig(config.LocatorConfig{
					ID: "com.palantir.godel-format-plugin:format-plugin:1.41.0",
					Checksums: map[string]string{
						"darwin-amd64": "825adbc5fb6d24ab6e85b1c67cb0cc9070a4e020297414737e92ee59dff191db",
						"darwin-arm64": "48a112e38598d43b589a1a58a587dd24d0227bf898c9b539d6be6da2d6b227b8",
						"linux-amd64":  "05c1dd1c5b1336013388531790b6f9b66c3b5ab1e1bd0349fea7572983ae34dc",
						"linux-arm64":  "767368b694130e3233e67f6073e792303ee076650ebfa672d86de514744eb4b7",
					},
				}),
			}),
			Assets: config.ToLocatorWithResolverConfigs([]config.LocatorWithResolverConfig{
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-format-asset-ptimports:ptimports-asset:1.40.0",
						Checksums: map[string]string{
							"darwin-amd64": "67ca0c688aae13b32130b63b5a7ff777f3ecaba46a2b858196c70194c520e28d",
							"darwin-arm64": "6ddbe29b00c116c925defd3e0b3f43e1099b58348d6fbf48e7185f8bb84a640b",
							"linux-amd64":  "1ef858557293ad51da80198e4bc167e7d2ee2bfa5d688cdfb2c98bdf7663d6a8",
							"linux-arm64":  "3f669242fc882f3c4e7cddb1c4fe933287c8ef6cb0b1b7281edea56253101acd",
						},
					}),
				},
			}),
		},
		{
			LocatorWithResolverConfig: config.ToLocatorWithResolverConfig(config.LocatorWithResolverConfig{
				Locator: config.ToLocatorConfig(config.LocatorConfig{
					ID: "com.palantir.godel-goland-plugin:goland-plugin:1.37.0",
					Checksums: map[string]string{
						"darwin-amd64": "95f80c56311a5a3e66316c5684a58a0396ccd3f0ce1da0df9e54879cd7e27f8e",
						"darwin-arm64": "13f6ada8fa93680e9701cdd248f9333424b39a1b0e3c49f16bb6bdc62839c8a5",
						"linux-amd64":  "e1f367cd5d11ea928de59036d479a274ba936d9c83406463affaa1d72f3c135a",
						"linux-arm64":  "240b5e37f6dca1f3b93693997e44680f78328bade9ed9593c07655801bf55293",
					},
				}),
			}),
		},
		{
			LocatorWithResolverConfig: config.ToLocatorWithResolverConfig(config.LocatorWithResolverConfig{
				Locator: config.ToLocatorConfig(config.LocatorConfig{
					ID: "com.palantir.okgo:check-plugin:1.53.0",
					Checksums: map[string]string{
						"darwin-amd64": "233c541e953b27563263ae75d1ecb6f6b748abece8ad947768107f4706c1b56d",
						"darwin-arm64": "59accf179240e805a7b86fe18475d1ac9f663bc00dba2c0e698101c20c5850d9",
						"linux-amd64":  "5fc8768356389e5b54b136d1474ded7cfaef87ea098696ca30cdfe66373d0292",
						"linux-arm64":  "2bb513845c4aaf90fcb065be658d1fb7e05cf5ed3317cfa9f04db880da0d1ccb",
					},
				}),
			}),
			Assets: config.ToLocatorWithResolverConfigs([]config.LocatorWithResolverConfig{
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-compiles:compiles-asset:1.46.0",
						Checksums: map[string]string{
							"darwin-amd64": "1c7499d5138617ae6021ff4131440667605a5f452bf537d4d4237235c7676474",
							"darwin-arm64": "b15ffb89219918431b12ebfb622160b4f0ee9b10150cc9d06293926b9d8e2663",
							"linux-amd64":  "ee9bae69cdb0ebe70b0222cbce15c169caeed957166f939414b96b5832e86ca8",
							"linux-arm64":  "cd1616c33925f5d965654914187aeadd104b2e273fa5abfd4f485324fd3183b4",
						},
					}),
				},
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-deadcode:deadcode-asset:1.42.0",
						Checksums: map[string]string{
							"darwin-amd64": "99187c761a14acb8050602aefd62352ffe73fe2653657f9c491c4cb11e55fda8",
							"darwin-arm64": "83e50a5889866ff925f4be43785e36f67d87f3a76778f8f60ce8502c19bb9f27",
							"linux-amd64":  "5e4d53a2229d4ef971778d140bd7c99e7cedaf9e4aa0679182a1496393c109f7",
							"linux-arm64":  "fecfb6594c69840ef4be74f35f8be7a5fda30297222ee02b73d4d43e5d35809d",
						},
					}),
				},
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-errcheck:errcheck-asset:1.44.0",
						Checksums: map[string]string{
							"darwin-amd64": "970bcc3a4a9ade4f2bd08652f079d45c4fe516d15baf965115c797d0798a27f2",
							"darwin-arm64": "c11e02f97d2b23f3497a735f0eb4ee7d258dc5bc6b96541753b6b33014d0d7d8",
							"linux-amd64":  "27eeee1c1e7d097def17036796bf84c95ff659538e2f6127182f2007995001a8",
							"linux-arm64":  "5143c33b182c563e447c11516a788065ab937c57f318a0e3193e39befa72e028",
						},
					}),
				},
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-golint:golint-asset:1.34.0",
						Checksums: map[string]string{
							"darwin-amd64": "e2c9b682a65c16eec45eb807736f61c293712dd8990997f31542c0702a97b19b",
							"darwin-arm64": "5d261372ce39c9dd817d64c16d0b522d7cc680aab3d109ebbfa2d9ccc06ea8fe",
							"linux-amd64":  "3fb3545f907fb242448f8061b01ea239710fe9913a273416f44457c22655c822",
							"linux-arm64":  "bc301cbf71f6816c8ea89bc78a9a504457c8d24ec334864a099f255bfd009d2c",
						},
					}),
				},
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-govet:govet-asset:1.38.0",
						Checksums: map[string]string{
							"darwin-amd64": "812d344989f5cdae39c8d075f4a30b8e0bf4e4831934045ccba96f2a2ee0d2f7",
							"darwin-arm64": "a4898777ebdfb09abdbdff3b3041d81fbf956ab0fe4928dd10b086619277adb5",
							"linux-amd64":  "28e45341264003e05d4ee48422917ab09fda61aed561bbc0939a9102fe9f5e28",
							"linux-arm64":  "90d5bc0e60964ee43df8bf6d0699be81bad0cccd4291aa593e32926cf12e74aa",
						},
					}),
				},
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-importalias:importalias-asset:1.37.0",
						Checksums: map[string]string{
							"darwin-amd64": "faff5dc22996909b025ff76a35e191bda7d9c0476df84b476fd129871a7ce360",
							"darwin-arm64": "b429a44f616e14de63e6fe72f9e143662284fe19759e69be5cf05cabca9e408d",
							"linux-amd64":  "7fea9c09fff27f26745a649208e128b198ce4e55d4655a8c52d3a569e771c390",
							"linux-arm64":  "47885c5189af09579e8b9e6c0508c8e165ce21e6dff7311210ddca385454d1e4",
						},
					}),
				},
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-ineffassign:ineffassign-asset:1.40.0",
						Checksums: map[string]string{
							"darwin-amd64": "704a334dfeec15b72a1b17b7a50ed7d9bc1aaa4aa82fde9c7a3acb5b6c77a46f",
							"darwin-arm64": "3d7a915624e8d056178f9a4de47afbd5141ed509651af27c60440c6942d222fc",
							"linux-amd64":  "8e72cfa2c8f8127fb1053c4ea6f56561b3acc4d3900be501053026ea7e330ce0",
							"linux-arm64":  "1d6f7048fd779f06c40ef7fd0a17f1ddfec3e58181f3033a38eaad84621a0695",
						},
					}),
				},
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-outparamcheck:outparamcheck-asset:1.42.0",
						Checksums: map[string]string{
							"darwin-amd64": "730d5525797d2ab94449c21effd2fa6084c10f590537c65c84d42da4eceba159",
							"darwin-arm64": "cd5db2f0cedebe07ad92d794e5fbe33b9e6058ebd6a486af7a1946186ea5852e",
							"linux-amd64":  "3d0ed6e9e84be293f669cae007d05cb07a6a8729de4b1cf5c57301a44aba6a0e",
							"linux-arm64":  "00cf0b0d343b56d095a650f76b5430cee1eddc4dfebb3ddd298529f14fe8b62d",
						},
					}),
				},
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-unconvert:unconvert-asset:1.42.0",
						Checksums: map[string]string{
							"darwin-amd64": "717b41cf322ae5c922d456bad2ae4003c0afaef4545db7f256a96ece39c25cea",
							"darwin-arm64": "e38575e0dad93e02bb48d2a58c29ccce21ef7d579fb29935d69c77e8c51acbdc",
							"linux-amd64":  "509bad4a27b2396112c29dc8c083a15ddf7c6293bc92bc22bf984b747a18d251",
							"linux-arm64":  "f7435a0c5be6f1983db3711ae51cd45dc0f4c6993076b5ca35d72ae5ff7a769f",
						},
					}),
				},
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-varcheck:varcheck-asset:1.42.0",
						Checksums: map[string]string{
							"darwin-amd64": "ed3877e0fb841611223ca8c4b73fb52e85049fee3838138068ca21d3da5b52c5",
							"darwin-arm64": "f3b3fe2822de941de21c6f8e16159f2701c661af718e649c521f4c3e3c1bcdb4",
							"linux-amd64":  "64f1fc9dcb62d873b86bdb002827417af5f800313e02beba996bd5ee42e586f9",
							"linux-arm64":  "7a3020f35493517b30daa4ce5b63b1351643a3fefa52786905f55e64df4ee86b",
						},
					}),
				},
			}),
		},
		{
			LocatorWithResolverConfig: config.ToLocatorWithResolverConfig(config.LocatorWithResolverConfig{
				Locator: config.ToLocatorConfig(config.LocatorConfig{
					ID: "com.palantir.godel-license-plugin:license-plugin:1.40.0",
					Checksums: map[string]string{
						"darwin-amd64": "692abd2ed3ff1842222eb0c3719077334d72ba25ea63fba1e48fe7941e0138d3",
						"darwin-arm64": "e391290432cb444ebd7fa005e64b79b9ca3a8c97532e8e91a3d128eee69f715d",
						"linux-amd64":  "52cb742f284b14dc922c5110bcb9c885af0da7250b694c44e7af2bd397231f65",
						"linux-arm64":  "da74a3167eb135410e80fb42a9a981fa3bf2f43231d9afbfc835ad5bad9f0488",
					},
				}),
			}),
		},
		{
			LocatorWithResolverConfig: config.ToLocatorWithResolverConfig(config.LocatorWithResolverConfig{
				Locator: config.ToLocatorConfig(config.LocatorConfig{
					ID: "com.palantir.godel-test-plugin:test-plugin:1.39.0",
					Checksums: map[string]string{
						"darwin-amd64": "f5b5deb5851810b09a693d8d7bee585f53c29b071640709aad74654129e5eddb",
						"darwin-arm64": "d8ef0ab8963eb1af3fc1e13ca7a9955a47e6dd726aacb168125a7fb1b45bedc8",
						"linux-amd64":  "759352abd0622563cf44ca640f5c01122388f660e5ddec01ad00392b064b7d1a",
						"linux-arm64":  "fca750c0d7a92236d63a29289eab3ac27052b00ddffb833dfaae037ac8a36fbe",
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
