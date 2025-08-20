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
					ID: "com.palantir.distgo:dist-plugin:1.86.0",
					Checksums: map[string]string{
						"darwin-amd64": "7b830a3f16ecfe19738654c23a10ec0be6457f8d1f6de1dd1715e0c0e0640503",
						"darwin-arm64": "b6e3780885261403aafc2a7d592c19d354bcd5672f95a885d75162d0122742ab",
						"linux-amd64":  "2d007e1f9a123fa547f687960772ad9a93f9caa24cedb4450c9b3240893f4b8f",
						"linux-arm64":  "bfc720b3044f5c8ec7900ac9dc6b8044cceee33c1f9160b8ab3101379619c708",
					},
				}),
			}),
		},
		{
			LocatorWithResolverConfig: config.ToLocatorWithResolverConfig(config.LocatorWithResolverConfig{
				Locator: config.ToLocatorConfig(config.LocatorConfig{
					ID: "com.palantir.godel-format-plugin:format-plugin:1.53.0",
					Checksums: map[string]string{
						"darwin-amd64": "88a3c1666c2bddc350fc6601a32500d87886c3fcf95c477b31619666cd04b0a7",
						"darwin-arm64": "230cd3ad76d61b9b20899d70249e3a0873e92ed4d4a4d5667d5b2093f004a9d5",
						"linux-amd64":  "aa3135ae1131e4ad55c2813fd39da039e3c90e1fdc53338802ca6aeee6a1f62b",
						"linux-arm64":  "eb08998b772564284b1c8205643c2a81304c85e7cc8b433902743f6592dfd1f1",
					},
				}),
			}),
			Assets: config.ToLocatorWithResolverConfigs([]config.LocatorWithResolverConfig{
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-format-asset-ptimports:ptimports-asset:1.52.0",
						Checksums: map[string]string{
							"darwin-amd64": "9284c15341ad6d0eb0b7c5c7b01b584288d7f141e2ba5050d2f9f36c59d1e3e9",
							"darwin-arm64": "70a79e4e796ca89e856e574afde6dc1f1d45affc7279379120d47689821d1ad4",
							"linux-amd64":  "ebc76db1e196238dc55f15c438fe0a4aa1b2b9c52db03f8589a7519c506e71ab",
							"linux-arm64":  "a12f69cb2c616d49a13235ed0724246891ac767d9e6e738a196fc0317abf0507",
						},
					}),
				},
			}),
		},
		{
			LocatorWithResolverConfig: config.ToLocatorWithResolverConfig(config.LocatorWithResolverConfig{
				Locator: config.ToLocatorConfig(config.LocatorConfig{
					ID: "com.palantir.godel-goland-plugin:goland-plugin:1.47.0",
					Checksums: map[string]string{
						"darwin-amd64": "7dec66ae76c870e778a6ff79396e85986914fc4f7d5083915dada0c2c12e921c",
						"darwin-arm64": "012ee49427b0ba3498cc51bab65c65e32e03b0f1a4d5eab42780f43a1076c15d",
						"linux-amd64":  "69ada20846ed9fda2658f89c0248b6534d3d834ba0bf3a1fc6fe5c2bbe715eab",
						"linux-arm64":  "acefcca1716ad4b5279fb481fc273e296b8ff05508b270b114032fadeb7610d9",
					},
				}),
			}),
		},
		{
			LocatorWithResolverConfig: config.ToLocatorWithResolverConfig(config.LocatorWithResolverConfig{
				Locator: config.ToLocatorConfig(config.LocatorConfig{
					ID: "com.palantir.okgo:check-plugin:1.64.0",
					Checksums: map[string]string{
						"darwin-amd64": "941575c5f9e67ce98494b7f3a5d71a905bb421c4829c39d41b3484d673a67ffb",
						"darwin-arm64": "c7842bd41a89ebcaaca28882de50c97871c82f22ab1d3311f31c23e8349b5d6e",
						"linux-amd64":  "a179bc5d58e192e0a4de9366d87d22362193272da021d952efc8f71017ca37f4",
						"linux-arm64":  "3b073ed06f25201e4354c9e6afdf70a619e9787d42976a22e37c45ce6d6bc621",
					},
				}),
			}),
			Assets: config.ToLocatorWithResolverConfigs([]config.LocatorWithResolverConfig{
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-compiles:compiles-asset:1.58.0",
						Checksums: map[string]string{
							"darwin-amd64": "91a866ffed7ce568ef935b59d14f66d9ce25d3415a7e1a2844ef9a617e29608a",
							"darwin-arm64": "04e5dcc0fb8a7c0495db3674eea156b5598eb7e0694867e139fbdb359f9aeada",
							"linux-amd64":  "470413a4775eede6bf9960ba6179b8c8de87867ecae0696f8ab0b6a4413f82e3",
							"linux-arm64":  "b5175186ba5a63860ab2920c48c0de62c3dc6ccd8e88d8057229c17e5acea33e",
						},
					}),
				},
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-deadcode:deadcode-asset:1.51.0",
						Checksums: map[string]string{
							"darwin-amd64": "02854fa10a1348f8a69ed5a2066c67b4ac2d7ceea719669c141fe343eb95f318",
							"darwin-arm64": "b2d65ccbeee72826092434ae443a69682b011e4a0d3c0e65b51986c04bb66b4d",
							"linux-amd64":  "cdcb84923b8c518d785cd9c276f9291424d0709f4bc2d32e4a5812ca078c88a7",
							"linux-arm64":  "ec5dc7bed9bb9c844bb651c47c3534df83427d38a58fa64fdf0f699fe0af2917",
						},
					}),
				},
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-errcheck:errcheck-asset:1.53.0",
						Checksums: map[string]string{
							"darwin-amd64": "43fc2d104be4fa7b1d71bb96a5cfd94a7a81173e702a9981d615d2b52c721e3e",
							"darwin-arm64": "91204f92fc8a9d434fdfb8a8fc18e8f463053537198c821bb922245d3806c709",
							"linux-amd64":  "9a9ce5061717e539995df4808fa5434682ff8834e0767eb8d39f71d1481dae2f",
							"linux-arm64":  "b107e32e990e6f952ad011f3c461c1161058f73e9c84719d507ae15f8d2ad403",
						},
					}),
				},
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-golint:golint-asset:1.48.0",
						Checksums: map[string]string{
							"darwin-amd64": "525441c194316a941583fe71ea2e3df11cc0c24c8b003e486e46a282406d855c",
							"darwin-arm64": "488009eb4303300be99fb5503376b5f04d70be8e0885bcc4d7ac773d825fb7f1",
							"linux-amd64":  "0c313281951d4929b33371a54d21fd2a3c8338b970f46cc1f2d16a80aed99137",
							"linux-arm64":  "99ad3104330c5530325cc6666de1b1befd1b8362ace98aae83637a4ae1639a9e",
						},
					}),
				},
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-govet:govet-asset:1.47.0",
						Checksums: map[string]string{
							"darwin-amd64": "3a788fceade25928125724cca4711196f0150ddaae17c1ea57dda4fe4db0ea66",
							"darwin-arm64": "0a5c6e6eb669cf4b165b3448cb9a7f756af8f352f62081d71d0fa3d4b566d7d1",
							"linux-amd64":  "d42c55755c0c6e339d674eece2dbeeb68f7fd6148dc32dcfc5878c3593e06e9b",
							"linux-arm64":  "b2dffa6d454c6e253481f8b28c92d212b6b601c0157e6d0f9a6511cc2819ae29",
						},
					}),
				},
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-importalias:importalias-asset:1.47.0",
						Checksums: map[string]string{
							"darwin-amd64": "020ef0f22672643354c1268e14be6ce1512e83f2a0a51212abb6c36d79de67bb",
							"darwin-arm64": "c9c8ab2e2594e3c0afaa49dd37bf4de8c972edaca1548c8bdaa0e36c8bc5e083",
							"linux-amd64":  "6c18e67e02d027c510bf094ae88c40d266c13962ff8a7bc6887e06cdd556944e",
							"linux-arm64":  "86930e234c7c98dde5ae4d44b9fa670ec8e09972672afbf31cd5cdcd0c510169",
						},
					}),
				},
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-ineffassign:ineffassign-asset:1.47.0",
						Checksums: map[string]string{
							"darwin-amd64": "97e6d4e027b31bf7a7fd8d7d6178d467936398a437dcf71ec9c781684c934899",
							"darwin-arm64": "b7d862c943a44cb74a81acc48e050e7fc11040877a884aff6a8601d3b590da81",
							"linux-amd64":  "88bbc73338b5fa4cbf4bcf217555dc28bf622bdf819ebe4f92b9a35c456d593e",
							"linux-arm64":  "6dddc75e932ccf3b73367f496b8cc5d489c26d3def8dffc00dc813fc629533de",
						},
					}),
				},
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-outparamcheck:outparamcheck-asset:1.50.0",
						Checksums: map[string]string{
							"darwin-amd64": "be9ccefa8e2fa1a0027955cfc2550a6f74024f523b02635251d0130325569cbe",
							"darwin-arm64": "f4e2a23d5f20919b4889b66af207224a6351b0f30857796fd9ed73f68261cac2",
							"linux-amd64":  "ca21a723c109332c3b88fb3df7f91399d58aebdba850c20c38ba4998504b9e1c",
							"linux-arm64":  "166bf165b65e1c281f7dac24470d6a11b984e61b75e36e045cd6335434172b75",
						},
					}),
				},
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-unconvert:unconvert-asset:1.51.0",
						Checksums: map[string]string{
							"darwin-amd64": "af36d123414d2727b54f239da05c31bd637315fb844007c6f2000f5299f993a9",
							"darwin-arm64": "bcc18abe74c9a1a877cc53f8c847b2108e298e490831331dfc6ee63513246195",
							"linux-amd64":  "d148829922d96e5788ea661ad6dd22251a68933dda3db7758c5b681c59154b19",
							"linux-arm64":  "af1a50768cbb3c498a02fb5972efd02b72ff4ed19832ea320a2e3c78dbe465f6",
						},
					}),
				},
				{
					Locator: config.ToLocatorConfig(config.LocatorConfig{
						ID: "com.palantir.godel-okgo-asset-varcheck:varcheck-asset:1.50.0",
						Checksums: map[string]string{
							"darwin-amd64": "7774a8999eccfc41fbd952a4330fbbb750f4b24c7fac3dcbb0a3cfbb3d8d646b",
							"darwin-arm64": "b507356bef97acdc05fec6fda6c81c3d1c0dc228a5ac167260a245a553c94b98",
							"linux-amd64":  "5e155c1cbd6db37dc75e5c9fd4143a2de5fc00d9e72f71e2bfd506c4dea93ed8",
							"linux-arm64":  "14953896e2f28e4e495d6df7585cb1cdb33d4513be1ecf4885e6a479e7cb68f8",
						},
					}),
				},
			}),
		},
		{
			LocatorWithResolverConfig: config.ToLocatorWithResolverConfig(config.LocatorWithResolverConfig{
				Locator: config.ToLocatorConfig(config.LocatorConfig{
					ID: "com.palantir.godel-license-plugin:license-plugin:1.49.0",
					Checksums: map[string]string{
						"darwin-amd64": "9bed959c0964a6a4051fb41f80af3c4fe7233cc128b125bc15a965ea9437e5f2",
						"darwin-arm64": "542c40727e3ab1c0964020187676486ef9e185c4e23f6482bf3fa02fb80865f7",
						"linux-amd64":  "63791eb2a7f3aecd85e1f3b0fdb1895b794a435c4eb44511160dca528be6a3d0",
						"linux-arm64":  "d03b47ddd347f2416614f0d87cfb0f735f0da312028946e840f3e85216db7531",
					},
				}),
			}),
		},
		{
			LocatorWithResolverConfig: config.ToLocatorWithResolverConfig(config.LocatorWithResolverConfig{
				Locator: config.ToLocatorConfig(config.LocatorConfig{
					ID: "com.palantir.godel-test-plugin:test-plugin:1.47.0",
					Checksums: map[string]string{
						"darwin-amd64": "03f8252e502fc6e573324ae3c2584d3d0a7028bcc2b541e55628a02506b4c09d",
						"darwin-arm64": "cea1eff96d37b7df8a110ee393cf9e8a2dcb8b532b6f28601bab8b30bfecbf6b",
						"linux-amd64":  "658be83977f879518ca7d533498de8354945e3f3477378b3c6605ce006a7c02b",
						"linux-arm64":  "704dbeb574d54af01f5c2772c8af2f164633963341b8a9c1bf0255e245eb2089",
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
