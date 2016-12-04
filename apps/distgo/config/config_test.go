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

package config_test

import (
	"strings"
	"testing"

	"github.com/palantir/pkg/matcher"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/palantir/godel/apps/distgo/config"
	"github.com/palantir/godel/apps/distgo/params"
	"github.com/palantir/godel/apps/distgo/pkg/osarch"
)

func TestReadConfig(t *testing.T) {
	for i, currCase := range []struct {
		yml  string
		json string
		want func() config.RawProjectConfig
	}{
		{
			yml: `
			products:
			    test:
			        build:
			            main-pkg: ./cmd/test
			            output-dir: build
			            build-args-script: |
			                               YEAR=$(date +%Y)
			                               echo "-ldflags"
			                               echo "-X"
			                               echo "main.year=$YEAR"
			            version-var: main.version
			            environment:
			                foo: bar
			                baz: 1
			                bool: TRUE
			            os-archs:
			                - os: "darwin"
			                  arch: "amd64"
			                - os: "linux"
			                  arch: "amd64"
			        dist:
			            output-dir: dist
			            input-dir: resources/input
			            dist-type:
			                type: sls
			                info:
			                    manifest-template-file: resources/input/manifest.yml
			                    product-type: service.v1
			                    yml-validation-exclude:
			                      names:
			                        - foo
			                      paths:
			                        - bar
			exclude:
			  names:
			    - ".*test"
			  paths:
			    - "vendor"
			`,
			json: `{"exclude":{"names":["distgo"],"paths":["generated_src"]}}`,
			want: func() config.RawProjectConfig {
				return config.RawProjectConfig{
					Products: map[string]config.RawProductConfig{
						"test": {
							Build: config.RawBuildConfig{
								MainPkg:   "./cmd/test",
								OutputDir: "build",
								BuildArgsScript: `YEAR=$(date +%Y)
echo "-ldflags"
echo "-X"
echo "main.year=$YEAR"
`,
								VersionVar: "main.version",
								Environment: map[string]string{
									"foo":  "bar",
									"baz":  "1",
									"bool": "TRUE",
								},
								OSArchs: []osarch.OSArch{
									{
										OS:   "darwin",
										Arch: "amd64",
									},
									{
										OS:   "linux",
										Arch: "amd64",
									},
								},
							},
							Dist: []config.RawDistConfig{{
								OutputDir: "dist",
								InputDir:  "resources/input",
								DistType: config.RawDistInfoConfig{
									Type: string(params.SLSDistType),
									Info: config.RawSLSDistConfig{
										ManifestTemplateFile: "resources/input/manifest.yml",
										ProductType:          "service.v1",
										YMLValidationExclude: matcher.NamesPathsCfg{
											Names: []string{"foo"},
											Paths: []string{"bar"},
										},
									},
								},
							}},
						},
					},
					Exclude: matcher.NamesPathsCfg{
						Names: []string{`.*test`, `distgo`},
						Paths: []string{`vendor`, `generated_src`},
					},
				}
			},
		},
		{
			yml: `
			products:
			  test:
			    dist:
			      dist-type:
			        type: rpm
			        info:
			          config-files:
			            - /usr/lib/systemd/system/orchestrator.service
			          before-install-script: |
			              /usr/bin/getent group orchestrator || /usr/sbin/groupadd \
			                  -g 380 orchestrator
			              /usr/bin/getent passwd orchestrator || /usr/sbin/useradd -r \
			                  -d /var/lib/orchestrator -g orchestrator -u 380 -m \
			                  -s /sbin/nologin orchestrator
			          after-install-script: |
			              systemctl daemon-reload
			          after-remove-script: |
			              systemctl daemon-reload
			`,
			want: func() config.RawProjectConfig {
				return config.RawProjectConfig{
					Products: map[string]config.RawProductConfig{
						"test": {
							Dist: []config.RawDistConfig{
								{
									DistType: config.RawDistInfoConfig{
										Type: string(params.RPMDistType),
										Info: config.RawRPMDistConfig{
											ConfigFiles: []string{"/usr/lib/systemd/system/orchestrator.service"},
											BeforeInstallScript: "" +
												"/usr/bin/getent group orchestrator || /usr/sbin/groupadd \\\n" +
												"    -g 380 orchestrator\n" +
												"/usr/bin/getent passwd orchestrator || /usr/sbin/useradd -r \\\n" +
												"    -d /var/lib/orchestrator -g orchestrator -u 380 -m \\\n" +
												"    -s /sbin/nologin orchestrator\n",
											AfterInstallScript: "systemctl daemon-reload\n",
											AfterRemoveScript:  "systemctl daemon-reload\n",
										},
									},
								},
							},
						},
					},
				}
			},
		},
		{
			yml: `
			products:
			  test:
			    dist:
			      - dist-type:
			          type: sls
			          info:
			            manifest-template-file: resources/input/manifest.yml
			      - dist-type:
			          type: rpm
			          info:
			            after-install-script: |
			                systemctl daemon-reload
			    publish:
			      group-id: com.palantir.pcloud
			      almanac:
			        metadata:
			          k: "v"
			        tags:
			          - "borked"
			`,
			want: func() config.RawProjectConfig {
				return config.RawProjectConfig{
					Products: map[string]config.RawProductConfig{
						"test": {
							Dist: []config.RawDistConfig{{
								DistType: config.RawDistInfoConfig{
									Type: string(params.SLSDistType),
									Info: config.RawSLSDistConfig{
										ManifestTemplateFile: "resources/input/manifest.yml",
									},
								},
							}, {
								DistType: config.RawDistInfoConfig{
									Type: string(params.RPMDistType),
									Info: config.RawRPMDistConfig{
										AfterInstallScript: "systemctl daemon-reload\n",
									},
								},
							}},
							DefaultPublish: config.RawPublishConfig{
								GroupID: "com.palantir.pcloud",
								Almanac: config.RawAlmanacConfig{
									Metadata: map[string]string{"k": "v"},
									Tags:     []string{"borked"},
								},
							},
						},
					},
				}
			},
		},
	} {
		// load config
		got, err := config.Read(unindent(currCase.yml), currCase.json)
		require.NoError(t, err, "Case %d", i)

		// require that it is valid
		_, err = got.ToParams()
		require.NoError(t, err, "Case %d", i)

		assert.Equal(t, currCase.want(), got, "Case %d", i)
	}
}

func TestFilteredProducts(t *testing.T) {
	for i, currCase := range []struct {
		cfg  func() params.Project
		want map[string]params.Product
	}{
		{
			cfg: func() params.Project {
				excludeCfg := matcher.NamesPathsCfg{
					Paths: []string{"vendor"},
				}
				return params.Project{
					Products: map[string]params.Product{
						"test": {
							Build: params.Build{
								MainPkg: "./test/main",
							},
						},
						"vendored": {
							Build: params.Build{
								MainPkg: "./vendor/test/main",
							},
						},
					},
					Exclude: excludeCfg.Matcher(),
				}
			},
			want: map[string]params.Product{
				"test": {
					Build: params.Build{
						MainPkg: "./test/main",
					},
				},
			},
		},
	} {
		got := currCase.cfg().FilteredProducts()
		assert.Equal(t, currCase.want, got, "Case %d", i)
	}
}

func unindent(input string) string {
	return strings.Replace(input, "\n\t\t\t", "\n", -1)
}
