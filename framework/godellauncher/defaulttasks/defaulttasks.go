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
					ID: "com.palantir.distgo:dist-plugin:1.29.0",
					Checksums: map[string]string{
						"darwin-amd64": "2fc8ce6cc29882eae8f04e3fb6626c4fa6672552228a52fc7345177f4a79a2e1",
						"darwin-arm64": "8fac6ca01af68b210cf6b65bd74df20cc797f9dfb36506ea8024c97a3a2953d5",
						"linux-amd64":  "d710353fd59c331ab3db29ed711defdb25e1dff26c072ca21f10446313e9c117",
						"linux-arm64":  "0f704a470c0b285f34475bcb0adbcac62f694a9f2aadb8ce8e4353dd0f211916",
					},
				}),
			}),
		},
		{
			LocatorWithResolverConfig: config.ToLocatorWithResolverConfig(config.LocatorWithResolverConfig{
				Locator: config.ToLocatorConfig(config.LocatorConfig{
					ID: "com.palantir.godel-format-plugin:format-plugin:1.8.0",
					Checksums: map[string]string{
						"darwin-amd64": "6e68aecc285fc8edabb38d74bf160693f358f4700b50170bd230d9d95a85ff5b",
						"darwin-arm64": "d7c79dca4f200ddd9290eaadb734e6ed083b0f1bfb80ee0377e9bf2be2e7d25a",
						"linux-amd64":  "d8b9221dbdd26142f9cfad6130f516333cb9e2cd9600dd5b12ee35c4ee8c04f7",
						"linux-arm64":  "942c1a0836f609080ed87165b3eeb2ad2c6de1d23e3d56c7ec7a24d225cda230",
					},
				}),
			}),
			Assets: config.ToLocatorWithResolverConfigs([]config.LocatorWithResolverConfig{
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-format-asset-ptimports:ptimports-asset:1.8.0",
						Checksums: map[string]string{
							"darwin-amd64": "f3cdba70c23e2bf7b48d306805734f292b66628913ed251eaf74134268ff2fdf",
							"darwin-arm64": "77c7cecf7ac599e1dc97c70768d663437bbb11a4874459234391c5713b36b4e5",
							"linux-amd64":  "e7d8d625507e414a21036315f1a4b8139ff495ab0d5706766436ea63060c0f4a",
							"linux-arm64":  "658bb38525d274e8e599145fdb87b44381d160c33b6a138374cef97b3973dc96",
						},
					}),
				},
			}),
		},
		{
			LocatorWithResolverConfig: config.ToLocatorWithResolverConfig(config.LocatorWithResolverConfig{
				Locator: config.ToLocatorConfig(config.LocatorConfig{
					ID: "com.palantir.godel-goland-plugin:goland-plugin:1.4.0",
					Checksums: map[string]string{
						"darwin-amd64": "c31715f1175d90ac4f473dc7c208cff9e6187468e08563c2aab577ad60e36bc7",
						"darwin-arm64": "9a8cacc99df2a48a82e4ed59576ad0caf094e78547008fc1d5d5605b9c52be3c",
						"linux-amd64":  "77d36a6bed2a54b027fa511709be4df1719cdd9138d9ac6c0f79451755689697",
						"linux-arm64":  "9d35255178798d7574b92d74972106af31182c105cd2d489230bd51c83af49c3",
					},
				}),
			}),
		},
		{
			LocatorWithResolverConfig: config.ToLocatorWithResolverConfig(config.LocatorWithResolverConfig{
				Locator: config.ToLocatorConfig(config.LocatorConfig{
					ID: "com.palantir.okgo:check-plugin:1.13.0",
					Checksums: map[string]string{
						"darwin-amd64": "b36351676db4240c1928bbc0081635e1009fcc7cad78d9dc3be9e834a73b8f19",
						"darwin-arm64": "15933b50a1c5881a085f9124cecc069a84ebad3f1d6a1c8791dfbbd4f01ce0ed",
						"linux-amd64":  "c4fdd86d20cb6e50ea929509246a2f0046708abb9f02cd6442e409be7eb272ec",
						"linux-arm64":  "ac832c738cc9f8e222a2b3d96e44a20a962746d0c3936df213948d50e50e2ab0",
					},
				}),
			}),
			Assets: config.ToLocatorWithResolverConfigs([]config.LocatorWithResolverConfig{
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-compiles:compiles-asset:1.8.0",
						Checksums: map[string]string{
							"darwin-amd64": "6cdd3a3a37f90df38d69ffc6ccd2966ff68c0b22f9ea63ff2035ca6fc40144bc",
							"darwin-arm64": "014e4dd66a672e172d8f5c2b1fa519b53a7ad7878bb79fb3db34be741599eafa",
							"linux-amd64":  "0f00836f1846cca9318159cb630be30570b011c3178b43a717a3fb0e68c5fb73",
							"linux-arm64":  "561249b4b16d934eb92ff747f6c98a3937114c8c6c2832b5fd2cbeb6516734a0",
						},
					}),
				},
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-deadcode:deadcode-asset:1.7.0",
						Checksums: map[string]string{
							"darwin-amd64": "d9443e65e682abc86fb8e926f419d9b9caa7a53dd3231126568beeb21acc9553",
							"darwin-arm64": "e244b052839a7cb5f7e4d516abf40b5e193163ff1cd2c8f12beea32e4ffde670",
							"linux-amd64":  "a2b593426f98826af87f3cde2054cb1594a0b3588816c11b260c77a93b660262",
							"linux-arm64":  "94e3d040acc5d72a24ffea57f1c7875a574c6f41435a5a5af5f18a5f33d7182d",
						},
					}),
				},
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-errcheck:errcheck-asset:1.9.0",
						Checksums: map[string]string{
							"darwin-amd64": "66701e8f34d7e5075afa312fc141770885282ed5e898ac3d51e9b99a3d426cf5",
							"darwin-arm64": "930ceb7486ebe933011cdda2f1186a7ed65924e683fb1f74a607df57f5299a7a",
							"linux-amd64":  "4f99f6edfa3c1c34eb004a14d2060d3effb577caeb4fc5ec65d7a77d92f8b78e",
							"linux-arm64":  "554a24307940f39e619eccf81f098cbf341124238ee88ad2e4d08018d0b7ccd2",
						},
					}),
				},
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-golint:golint-asset:1.5.0",
						Checksums: map[string]string{
							"darwin-amd64": "86c90ed0bc9236e439997bd4a3028f63d72686daf4538f5d3f3db11ea002d79c",
							"darwin-arm64": "b53f4c3fd72289715204de17dee1334189ff2992d9032b262dfc63304e41b87d",
							"linux-amd64":  "4f51d6ba3d1963e35d110bd8200f1c18f76b501e7fc4c83e4a8f498cad191613",
							"linux-arm64":  "c816de967f6719d69d85ea75ac2671ba3ff02c5e192280a95c0e4d7c46165fab",
						},
					}),
				},
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-govet:govet-asset:1.5.0",
						Checksums: map[string]string{
							"darwin-amd64": "2483818f56c194e8515a30a638aaf50db67e89801f8aa9c5f945bad471475462",
							"darwin-arm64": "24a54b263687d43be6ac014e2ea98eeb345260b4211b502f43fe1d2b97b51034",
							"linux-amd64":  "6d9b6980565d2bfa25b2a9314c487545776ca2e3be01b5bf106fc6f3ac8db5d7",
							"linux-arm64":  "c2a123b09706ea9ee2f03f66d211e9e83d46f1064470f070edcd35d1fb00b42f",
						},
					}),
				},
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-importalias:importalias-asset:1.5.0",
						Checksums: map[string]string{
							"darwin-amd64": "e5b4b2cbaf306bfa20bc20280014f94ddcb40983e39966e7b3b09bd898dfd5ba",
							"darwin-arm64": "d67283f57f5f6e84c16ab09cbc7d9ef883cd89fd76eb9aa4f63ea6380fce0ffb",
							"linux-amd64":  "904d302c4acf35ee158f322b782d0499a82228c70b1c971b958bd36bed8ac9b7",
							"linux-arm64":  "622aa0e3bc5ba9aa7f3c7ca0710df7831573930c0f0fd38b5beb2d9ffa15d048",
						},
					}),
				},
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-ineffassign:ineffassign-asset:1.5.0",
						Checksums: map[string]string{
							"darwin-amd64": "87a47565d5c69513f46239cbae9570808290e538e748fe41a251f32b0920718f",
							"darwin-arm64": "ac8d5e38212d2e8eab22e0d44a1c583d578312af82d9587e40000d096d76ebfa",
							"linux-amd64":  "ba1d7accb22771476f8649a55936209e4f62b9977b3511fc310bdf8d329f781b",
							"linux-arm64":  "d38f618c388890c92c9320023153adde7c07575e18ac3d57fa7228511e0a3a62",
						},
					}),
				},
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-outparamcheck:outparamcheck-asset:1.9.0",
						Checksums: map[string]string{
							"darwin-amd64": "4e0050e0b39444fd3ceea1b99344a679d1127eed4f541236bceb6597ed2fbabc",
							"darwin-arm64": "a66ecec9eea5d2018aa2b2b856e4f0ad9bfe807da85d481a8dc605e7687eed01",
							"linux-amd64":  "7492c5de491ec50d30bec50a63d928216f8dd990b2565874e3668fae98472bc3",
							"linux-arm64":  "926ae44b73b99571e0e02801e1ddd5e51fdac7bcc4f72268832dbb8507f3e89f",
						},
					}),
				},
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-unconvert:unconvert-asset:1.8.0",
						Checksums: map[string]string{
							"darwin-amd64": "e06380ba081ab2bc86e15b842e5a11a769a7cecfb573bbfd4e8a316bf2effddb",
							"darwin-arm64": "dfa5aefbfeb069ba4c024cb7f7c36de88df50c0aabd8de997d05ad8e303e765f",
							"linux-amd64":  "09c584451f15995927e1c4f2acb9c4bb66f18e65c243d592c897d79e721da120",
							"linux-arm64":  "6dc3847873c8f421a10b8241d2937196793db04e30f2f602cde59e90be072490",
						},
					}),
				},
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-varcheck:varcheck-asset:1.7.0",
						Checksums: map[string]string{
							"darwin-amd64": "d90568b8b11c587394e92691a0154c1807e335fb018af81fdbc328378d1955eb",
							"darwin-arm64": "850b4474b18dbb8f7d184f25da2b7d726bf4f89dd7585a03e116b6bbbb07d97a",
							"linux-amd64":  "ea8cc874e2d0861bcd080c6b53a3910608c5a23de0134db985be46a8739e7659",
							"linux-arm64":  "426fdb0ed9a0395917f9c3646265812d8e28f6b3928c2091f372b225de0caafc",
						},
					}),
				},
			}),
		},
		{
			LocatorWithResolverConfig: config.ToLocatorWithResolverConfig(config.LocatorWithResolverConfig{
				Locator: config.ToLocatorConfig(config.LocatorConfig{
					ID: "com.palantir.godel-license-plugin:license-plugin:1.6.0",
					Checksums: map[string]string{
						"darwin-amd64": "7b4a17b9d208c0b3aecdfb053bb39acedf6d25bb347845863dd54f3d480073d8",
						"darwin-arm64": "d9b0a2b33a7bbf4098b6929ab5842977409efc53b56698e3cdcf1b336cc5709f",
						"linux-amd64":  "b4fa87f82d573022921cf935f4c13b63fa1a5808d341a5433d6b3d895292cecd",
						"linux-arm64":  "8045ddf167d3ffe5c5240048b67f9f5427a5e9acf1ead5145cb15a1b63a64841",
					},
				}),
			}),
		},
		{
			LocatorWithResolverConfig: config.ToLocatorWithResolverConfig(config.LocatorWithResolverConfig{
				Locator: config.ToLocatorConfig(config.LocatorConfig{
					ID: "com.palantir.godel-test-plugin:test-plugin:1.8.0",
					Checksums: map[string]string{
						"darwin-amd64": "812ef8496dae13c38716593c566e237d8588f9a16ca5492f1949ba7c4c800e71",
						"darwin-arm64": "061019dc30533938e9e8d97e7d8efb4c8616445fadb29b3c9835443154f533cd",
						"linux-amd64":  "744715a3cb4a6b51e40599129d36d2014b80167897d857dacc91e233dcaf51f5",
						"linux-arm64":  "31e5b95404d696f758ee0a2bfbfd0749cf31bba7228397253a0278758af7e86f",
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
