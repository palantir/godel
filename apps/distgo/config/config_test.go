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
		want func() config.Project
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
			          reloadable: true
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
			want: func() config.Project {
				return config.Project{
					Products: map[string]config.Product{
						"test": {
							Build: config.Build{
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
							Dist: []config.Dist{{
								OutputDir: "dist",
								InputDir:  "resources/input",
								DistType: config.DistInfo{
									Type: string(params.SLSDistType),
									Info: config.SLSDist{
										ManifestTemplateFile: "resources/input/manifest.yml",
										ProductType:          "service.v1",
										Reloadable:           true,
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
			    build:
			      main-pkg: ./cmd/test
			    dist:
			      dist-type:
			        type: bin
			`,
			want: func() config.Project {
				return config.Project{
					Products: map[string]config.Product{
						"test": {
							Build: config.Build{
								MainPkg: "./cmd/test",
							},
							Dist: []config.Dist{{
								DistType: config.DistInfo{
									Type: string(params.BinDistType),
									Info: config.BinDist{
										OmitInitSh: nil,
									},
								},
							}},
						},
					},
					Exclude: matcher.NamesPathsCfg{},
				}
			},
		},
		{
			yml: `
			products:
			  test:
			    build:
			      main-pkg: ./cmd/test
			    dist:
			      dist-type:
			        type: os-arch-bin
			        info:
			          os-archs:
			            - os: "darwin"
			              arch: "amd64"
			            - os: "linux"
			              arch: "amd64"
			`,
			want: func() config.Project {
				return config.Project{
					Products: map[string]config.Product{
						"test": {
							Build: config.Build{
								MainPkg: "./cmd/test",
							},
							Dist: []config.Dist{{
								DistType: config.DistInfo{
									Type: string(params.OSArchBinDistType),
									Info: config.OSArchBinDist{
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
								},
							}},
						},
					},
					Exclude: matcher.NamesPathsCfg{},
				}
			},
		},
		{
			yml: `
			products:
			  test:
			    build:
			      skip: true
			    dist:
			      dist-type:
			        type: manual
			        info:
			          extension: tgz
			`,
			want: func() config.Project {
				return config.Project{
					Products: map[string]config.Product{
						"test": {
							Build: config.Build{
								Skip: true,
							},
							Dist: []config.Dist{{
								DistType: config.DistInfo{
									Type: string(params.ManualDistType),
									Info: config.ManualDist{
										Extension: "tgz",
									},
								},
							}},
						},
					},
					Exclude: matcher.NamesPathsCfg{},
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
			want: func() config.Project {
				return config.Project{
					Products: map[string]config.Product{
						"test": {
							Dist: []config.Dist{
								{
									DistType: config.DistInfo{
										Type: string(params.RPMDistType),
										Info: config.RPMDist{
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
			want: func() config.Project {
				return config.Project{
					Products: map[string]config.Product{
						"test": {
							Dist: []config.Dist{{
								DistType: config.DistInfo{
									Type: string(params.SLSDistType),
									Info: config.SLSDist{
										ManifestTemplateFile: "resources/input/manifest.yml",
									},
								},
							}, {
								DistType: config.DistInfo{
									Type: string(params.RPMDistType),
									Info: config.RPMDist{
										AfterInstallScript: "systemctl daemon-reload\n",
									},
								},
							}},
							DefaultPublish: config.Publish{
								GroupID: "com.palantir.pcloud",
								Almanac: config.Almanac{
									Metadata: map[string]string{"k": "v"},
									Tags:     []string{"borked"},
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
			    docker:
			      -
			        repository: docker.hub/test
			        tag: test
			        context-dir: context/dir/path
			        dependencies:
			          -
			            product: dep1
			            type: sls
			            target-file: dep1-sls.tgz
			          -
			            product: dep2
			            type: rpm
			            target-file: dep2-rpm.tgz
			      -
			        repository: docker.hub/test-alpine
			        tag: test
			        context-dir: context/dir/path-alpine
			        dependencies:
			          -
			            product: dep1
			            type: sls
			            target-file: dep1-sls.tgz
			`,
			want: func() config.Project {
				return config.Project{
					Products: map[string]config.Product{
						"test": {
							DockerImages: []config.DockerImage{
								{
									Repository: "docker.hub/test",
									Tag:        "test",
									ContextDir: "context/dir/path",
									Deps: []config.DockerDep{
										{
											Product:    "dep1",
											Type:       "sls",
											TargetFile: "dep1-sls.tgz",
										},
										{
											Product:    "dep2",
											Type:       "rpm",
											TargetFile: "dep2-rpm.tgz",
										},
									},
								},
								{
									Repository: "docker.hub/test-alpine",
									Tag:        "test",
									ContextDir: "context/dir/path-alpine",
									Deps: []config.DockerDep{
										{
											Product:    "dep1",
											Type:       "sls",
											TargetFile: "dep1-sls.tgz",
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
			    docker:
			      -
			        repository: docker.hub/test
			        tag: test
			        context-dir: context/dir/path
			        dependencies:
			          -
			            product: dep1
			            type: sls
			            target-file: dep1-sls.tgz
			          -
			            product: dep2
			            type: rpm
			            target-file: dep2-rpm.tgz
			        build-args-script: |
			                           echo "--label"
			                           echo "test_label=test_value"
			        info:
			          type: sls
			          data:
			            group-id: com.palantir.test
			            product-type: test-type
			            manifest-extensions:
			              test_key: test_value
			`,
			want: func() config.Project {
				return config.Project{
					Products: map[string]config.Product{
						"test": {
							DockerImages: []config.DockerImage{
								{
									Repository: "docker.hub/test",
									Tag:        "test",
									ContextDir: "context/dir/path",
									Deps: []config.DockerDep{
										{
											Product:    "dep1",
											Type:       "sls",
											TargetFile: "dep1-sls.tgz",
										},
										{
											Product:    "dep2",
											Type:       "rpm",
											TargetFile: "dep2-rpm.tgz",
										},
									},
									BuildArgsScript: `echo "--label"
echo "test_label=test_value"
`,
									Info: config.DockerImageInfo{
										Type: "sls",
										Data: config.SLSDockerImageInfo{
											GroupID:      "com.palantir.test",
											ProuductType: "test-type",
											Extensions: map[string]interface{}{
												"test_key": "test_value",
											},
										},
									},
								},
							},
						},
					},
				}
			},
		},
	} {
		// load config
		got, err := config.LoadRawConfig(unindent(currCase.yml), currCase.json)
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
