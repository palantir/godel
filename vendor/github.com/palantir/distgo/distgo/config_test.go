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

package distgo_test

import (
	"io/ioutil"
	"os"
	"path"
	"testing"

	"github.com/nmiyake/pkg/dirs"
	"github.com/nmiyake/pkg/gofiles"
	"github.com/palantir/godel/pkg/osarch"
	"github.com/palantir/pkg/gittest"
	"github.com/palantir/pkg/matcher"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v2"

	"github.com/palantir/distgo/dister"
	"github.com/palantir/distgo/distgo"
	"github.com/palantir/distgo/dockerbuilder"
)

var (
	disterFactory        distgo.DisterFactory
	defaultDisterCfg     distgo.DisterConfig
	dockerBuilderFactory distgo.DockerBuilderFactory
)

func init() {
	var err error
	disterFactory, err = dister.NewDisterFactory()
	if err != nil {
		panic(err)
	}
	defaultDisterCfg, err = dister.DefaultConfig()
	if err != nil {
		panic(err)
	}
	dockerBuilderFactory, err = dockerbuilder.NewDockerBuilderFactory()
	if err != nil {
		panic(err)
	}
}

func TestLoadConfig(t *testing.T) {
	for i, tc := range []struct {
		yml  string
		want distgo.ProjectConfig
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
      disters:
        type: sls
        config:
          manifest-template-file: resources/input/manifest.yml
          product-type: service.v1
          reloadable: true
          yml-validation-exclude:
            names:
              - foo
            paths:
              - bar
    publish:
      group-id: com.test.foo
      info:
        bintray:
          username: username
          password: password
script-includes: |
                 #!/usr/bin/env bash
version-script: |
                git describe
exclude:
  names:
    - ".*test"
  paths:
    - "vendor"
`,
			want: distgo.ProjectConfig{
				Products: map[distgo.ProductID]distgo.ProductConfig{
					"test": {
						Build: &distgo.BuildConfig{
							OutputDir: stringPtr("build"),
							MainPkg:   stringPtr("./cmd/test"),
							BuildArgsScript: stringPtr(`YEAR=$(date +%Y)
echo "-ldflags"
echo "-X"
echo "main.year=$YEAR"
`),
							VersionVar: stringPtr("main.version"),
							Environment: &map[string]string{
								"foo":  "bar",
								"baz":  "1",
								"bool": "TRUE",
							},
							OSArchs: &[]osarch.OSArch{
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
						Dist: &distgo.DistConfig{
							OutputDir: stringPtr("dist"),
							Disters: &distgo.DistersConfig{
								"sls": {
									Type: stringPtr("sls"),
									Config: &yaml.MapSlice{
										{
											Key:   "manifest-template-file",
											Value: "resources/input/manifest.yml",
										},
										{
											Key:   "product-type",
											Value: "service.v1",
										},
										{
											Key:   "reloadable",
											Value: true,
										},
										{
											Key: "yml-validation-exclude",
											Value: yaml.MapSlice{
												{
													Key:   "names",
													Value: []interface{}{"foo"},
												},
												{
													Key:   "paths",
													Value: []interface{}{"bar"},
												},
											},
										},
									},
								},
							},
						},
						Publish: &distgo.PublishConfig{
							GroupID: stringPtr("com.test.foo"),
							PublishInfo: &map[distgo.PublishID]yaml.MapSlice{
								"bintray": []yaml.MapItem{
									{
										Key:   "username",
										Value: "username",
									},
									{
										Key:   "password",
										Value: "password",
									},
								},
							},
						},
					},
				},
				ScriptIncludes: `#!/usr/bin/env bash
`,
				VersionScript: `git describe
`,
				Exclude: matcher.NamesPathsCfg{
					Names: []string{`.*test`},
					Paths: []string{`vendor`},
				},
			},
		},
		{
			yml: `
products:
  test:
    build:
      main-pkg: ./cmd/test
    dist:
      disters:
        type: bin
`,
			want: distgo.ProjectConfig{
				Products: map[distgo.ProductID]distgo.ProductConfig{
					"test": {
						Build: &distgo.BuildConfig{
							MainPkg: stringPtr("./cmd/test"),
						},
						Dist: &distgo.DistConfig{
							Disters: &distgo.DistersConfig{
								"bin": {
									Type: stringPtr("bin"),
								},
							},
						},
					},
				},
				Exclude: matcher.NamesPathsCfg{},
			},
		},
		{
			yml: `
products:
  test:
    build:
      main-pkg: ./cmd/test
    dist:
      disters:
        type: os-arch-bin
        config:
          os-archs:
            - os: "darwin"
              arch: "amd64"
            - os: "linux"
              arch: "amd64"
`,
			want: distgo.ProjectConfig{
				Products: map[distgo.ProductID]distgo.ProductConfig{
					"test": {
						Build: &distgo.BuildConfig{
							MainPkg: stringPtr("./cmd/test"),
						},
						Dist: &distgo.DistConfig{
							Disters: &distgo.DistersConfig{
								dister.OSArchBinDistTypeName: {
									Type: stringPtr(dister.OSArchBinDistTypeName),
									Config: &yaml.MapSlice{
										{
											Key: "os-archs",
											Value: []interface{}{
												yaml.MapSlice{
													{
														Key:   "os",
														Value: "darwin",
													},
													{
														Key:   "arch",
														Value: "amd64",
													},
												},
												yaml.MapSlice{
													{
														Key:   "os",
														Value: "linux",
													},
													{
														Key:   "arch",
														Value: "amd64",
													},
												},
											},
										},
									},
								},
							},
						},
					},
				},
				Exclude: matcher.NamesPathsCfg{},
			},
		},
		{
			yml: `
products:
  test:
    dist:
      disters:
        type: manual
        config:
          extension: tgz
`,
			want: distgo.ProjectConfig{
				Products: map[distgo.ProductID]distgo.ProductConfig{
					"test": {
						Build: nil,
						Dist: &distgo.DistConfig{
							Disters: &distgo.DistersConfig{
								dister.ManualDistTypeName: {
									Type: stringPtr(dister.ManualDistTypeName),
									Config: &yaml.MapSlice{
										{
											Key:   "extension",
											Value: "tgz",
										},
									},
								},
							},
						},
					},
				},
				Exclude: matcher.NamesPathsCfg{},
			},
		},
	} {
		got, err := distgo.LoadConfig([]byte(tc.yml))
		require.NoError(t, err, "Case %d", i)
		assert.Equal(t, tc.want, got, "Case %d", i)
	}
}

func TestProjectConfig_ToParam(t *testing.T) {
	for i, tc := range []struct {
		name string
		yml  string
		want distgo.ProjectParam
	}{
		{
			"parameters populated from defaults",
			`
product-defaults:
  build:
    output-dir: default-output
products:
  test-1:
    build:
      output-dir: test1-output
  test-2:
    build:
      version-var: main.version
  test-3:
    build:
      output-dir: ""
  test-4:
`,
			distgo.ProjectParam{
				Products: map[distgo.ProductID]distgo.ProductParam{
					"test-1": {
						ID: "test-1",
						Build: &distgo.BuildParam{
							NameTemplate: "{{Product}}",
							OutputDir:    "test1-output",
							OSArchs: []osarch.OSArch{
								osarch.Current(),
							},
						},
					},
					"test-2": {
						ID: "test-2",
						Build: &distgo.BuildParam{
							NameTemplate: "{{Product}}",
							VersionVar:   "main.version",
							OutputDir:    "default-output",
							OSArchs: []osarch.OSArch{
								osarch.Current(),
							},
						},
					},
					"test-3": {
						ID: "test-3",
						Build: &distgo.BuildParam{
							NameTemplate: "{{Product}}",
							OutputDir:    "out/build",
							OSArchs: []osarch.OSArch{
								osarch.Current(),
							},
						},
					},
					"test-4": {
						ID: "test-4",
					},
				},
			},
		},
		{
			"dependencies populated properly",
			`
products:
  test-1:
    dependencies:
      - test-2
  test-2:
    dependencies:
      - test-3
  test-3:
`,
			distgo.ProjectParam{
				Products: map[distgo.ProductID]distgo.ProductParam{
					"test-1": {
						ID: "test-1",
						FirstLevelDependencies: []distgo.ProductID{
							"test-2",
						},
						AllDependencies: map[distgo.ProductID]distgo.ProductParam{
							"test-2": {
								ID: "test-2",
								FirstLevelDependencies: []distgo.ProductID{
									"test-3",
								},
							},
							"test-3": {
								ID: "test-3",
							},
						},
					},
					"test-2": {
						ID: "test-2",
						FirstLevelDependencies: []distgo.ProductID{
							"test-3",
						},
						AllDependencies: map[distgo.ProductID]distgo.ProductParam{
							"test-3": {
								ID: "test-3",
							},
						},
					},
					"test-3": {
						ID: "test-3",
					},
				},
			},
		},
		{
			"Docker build dependencies expanded",
			`
products:
  test-1:
    docker:
      docker-builders:
        default:
          type: default
          context-dir: docker
          input-builds:
            - test-2
            - test-2.linux-amd64
            - test-3.darwin-amd64
          tag-templates:
            - foo:latest
    dependencies:
      - test-2
  test-2:
    build:
      main-pkg: ./test-2
      os-archs:
        - os: "darwin"
          arch: "amd64"
        - os: "linux"
          arch: "amd64"
    dependencies:
      - test-3
  test-3:
    build:
      main-pkg: ./test-3
      os-archs:
        - os: "darwin"
          arch: "amd64"
        - os: "linux"
          arch: "amd64"
`,
			distgo.ProjectParam{
				Products: map[distgo.ProductID]distgo.ProductParam{
					"test-1": {
						ID: "test-1",
						Docker: &distgo.DockerParam{
							DockerBuilderParams: map[distgo.DockerID]distgo.DockerBuilderParam{
								"default": {
									DockerBuilder:  dockerbuilder.NewDefaultDockerBuilder(nil),
									DockerfilePath: "Dockerfile",
									ContextDir:     "docker",
									InputBuilds: []distgo.ProductBuildID{
										"test-2.darwin-amd64",
										"test-2.linux-amd64",
										"test-3.darwin-amd64",
									},
									TagTemplates: []string{
										"foo:latest",
									},
								},
							},
						},
						FirstLevelDependencies: []distgo.ProductID{
							"test-2",
						},
						AllDependencies: map[distgo.ProductID]distgo.ProductParam{
							"test-2": {
								ID: "test-2",
								Build: &distgo.BuildParam{
									NameTemplate: "{{Product}}",
									OutputDir:    "out/build",
									MainPkg:      "./test-2",
									OSArchs: []osarch.OSArch{
										mustOSArch("darwin-amd64"),
										mustOSArch("linux-amd64"),
									},
								},
								FirstLevelDependencies: []distgo.ProductID{
									"test-3",
								},
							},
							"test-3": {
								ID: "test-3",
								Build: &distgo.BuildParam{
									NameTemplate: "{{Product}}",
									OutputDir:    "out/build",
									MainPkg:      "./test-3",
									OSArchs: []osarch.OSArch{
										mustOSArch("darwin-amd64"),
										mustOSArch("linux-amd64"),
									},
								},
							},
						},
					},
					"test-2": {
						ID: "test-2",
						Build: &distgo.BuildParam{
							NameTemplate: "{{Product}}",
							OutputDir:    "out/build",
							MainPkg:      "./test-2",
							OSArchs: []osarch.OSArch{
								mustOSArch("darwin-amd64"),
								mustOSArch("linux-amd64"),
							},
						},
						FirstLevelDependencies: []distgo.ProductID{
							"test-3",
						},
						AllDependencies: map[distgo.ProductID]distgo.ProductParam{
							"test-3": {
								ID: "test-3",
								Build: &distgo.BuildParam{
									NameTemplate: "{{Product}}",
									OutputDir:    "out/build",
									MainPkg:      "./test-3",
									OSArchs: []osarch.OSArch{
										mustOSArch("darwin-amd64"),
										mustOSArch("linux-amd64"),
									},
								},
							},
						},
					},
					"test-3": {
						ID: "test-3",
						Build: &distgo.BuildParam{
							NameTemplate: "{{Product}}",
							OutputDir:    "out/build",
							MainPkg:      "./test-3",
							OSArchs: []osarch.OSArch{
								mustOSArch("darwin-amd64"),
								mustOSArch("linux-amd64"),
							},
						},
					},
				},
			},
		},
		{
			"Docker dist dependencies expanded",
			`
products:
  test-1:
    docker:
      docker-builders:
        default:
          type: default
          context-dir: docker
          input-dists:
            - test-2
            - test-2.bar
            - test-3.os-arch-bin
          tag-templates:
            - foo:latest
    dependencies:
      - test-2
  test-2:
    dist:
      disters:
        foo:
          type: os-arch-bin
        bar:
          type: os-arch-bin
    dependencies:
      - test-3
  test-3:
    dist:
      disters:
        type: os-arch-bin
`,
			distgo.ProjectParam{
				Products: map[distgo.ProductID]distgo.ProductParam{
					"test-1": {
						ID: "test-1",
						Docker: &distgo.DockerParam{
							DockerBuilderParams: map[distgo.DockerID]distgo.DockerBuilderParam{
								"default": {
									DockerBuilder:  dockerbuilder.NewDefaultDockerBuilder(nil),
									DockerfilePath: "Dockerfile",
									ContextDir:     "docker",
									InputDists: []distgo.ProductDistID{
										"test-2.bar",
										"test-2.foo",
										"test-3.os-arch-bin",
									},
									TagTemplates: []string{
										"foo:latest",
									},
								},
							},
						},
						FirstLevelDependencies: []distgo.ProductID{
							"test-2",
						},
						AllDependencies: map[distgo.ProductID]distgo.ProductParam{
							"test-2": {
								ID: "test-2",
								Dist: &distgo.DistParam{
									OutputDir: "out/dist",
									DistParams: map[distgo.DistID]distgo.DisterParam{
										"bar": {
											NameTemplate: "{{Product}}-{{Version}}",
											Dister:       dister.NewOSArchBinDister(osarch.Current()),
										},
										"foo": {
											NameTemplate: "{{Product}}-{{Version}}",
											Dister:       dister.NewOSArchBinDister(osarch.Current()),
										},
									},
								},
								FirstLevelDependencies: []distgo.ProductID{
									"test-3",
								},
							},
							"test-3": {
								ID: "test-3",
								Dist: &distgo.DistParam{
									OutputDir: "out/dist",
									DistParams: map[distgo.DistID]distgo.DisterParam{
										"os-arch-bin": {
											NameTemplate: "{{Product}}-{{Version}}",
											Dister:       dister.NewOSArchBinDister(osarch.Current()),
										},
									},
								},
							},
						},
					},
					"test-2": {
						ID: "test-2",
						Dist: &distgo.DistParam{
							OutputDir: "out/dist",
							DistParams: map[distgo.DistID]distgo.DisterParam{
								"bar": {
									NameTemplate: "{{Product}}-{{Version}}",
									Dister:       dister.NewOSArchBinDister(osarch.Current()),
								},
								"foo": {
									NameTemplate: "{{Product}}-{{Version}}",
									Dister:       dister.NewOSArchBinDister(osarch.Current()),
								},
							},
						},
						FirstLevelDependencies: []distgo.ProductID{
							"test-3",
						},
						AllDependencies: map[distgo.ProductID]distgo.ProductParam{
							"test-3": {
								ID: "test-3",
								Dist: &distgo.DistParam{
									OutputDir: "out/dist",
									DistParams: map[distgo.DistID]distgo.DisterParam{
										"os-arch-bin": {
											NameTemplate: "{{Product}}-{{Version}}",
											Dister:       dister.NewOSArchBinDister(osarch.Current()),
										},
									},
								},
							},
						},
					},
					"test-3": {
						ID: "test-3",
						Dist: &distgo.DistParam{
							OutputDir: "out/dist",
							DistParams: map[distgo.DistID]distgo.DisterParam{
								"os-arch-bin": {
									NameTemplate: "{{Product}}-{{Version}}",
									Dister:       dister.NewOSArchBinDister(osarch.Current()),
								},
							},
						},
					},
				},
			},
		},
	} {
		gotCfg, err := distgo.LoadConfig([]byte(tc.yml))
		require.NoError(t, err, "Case %d: %s", i, tc.name)
		gotParam, err := gotCfg.ToParam("", disterFactory, defaultDisterCfg, dockerBuilderFactory)
		require.NoError(t, err, "Case %d: %s", i, tc.name)
		assert.Equal(t, tc.want, gotParam, "Case %d: %s", i, tc.name)
	}
}

func TestProjectConfig_DefaultProducts(t *testing.T) {
	tmpDir, cleanup, err := dirs.TempDir("", "")
	require.NoError(t, err)
	defer cleanup()

	defaultProductParam := func(name, mainPkgDir string, modifyFns ...func(*distgo.ProductParam)) distgo.ProductParam {
		param := distgo.ProductParam{
			ID: distgo.ProductID(name),
			Build: &distgo.BuildParam{
				NameTemplate: "{{Product}}",
				OutputDir:    "out/build",
				MainPkg:      mainPkgDir,
				OSArchs: []osarch.OSArch{
					osarch.Current(),
				},
			},
			Run:     &distgo.RunParam{},
			Publish: &distgo.PublishParam{},
			Dist: &distgo.DistParam{
				OutputDir: "out/dist",
				DistParams: map[distgo.DistID]distgo.DisterParam{
					dister.OSArchBinDistTypeName: {
						NameTemplate: "{{Product}}-{{Version}}",
						Dister:       dister.NewOSArchBinDister(osarch.Current()),
					},
				},
			},
			Docker: &distgo.DockerParam{
				DockerBuilderParams: map[distgo.DockerID]distgo.DockerBuilderParam{},
			},
		}
		for _, currFn := range modifyFns {
			currFn(&param)
		}
		return param
	}

	for i, tc := range []struct {
		name    string
		yml     string
		gofiles []gofiles.GoFileSpec
		want    distgo.ProjectParam
	}{
		{
			"configuration created for main packages",
			`
`,
			[]gofiles.GoFileSpec{
				{
					RelPath: "foo/main.go",
					Src:     `package main`,
				},
				{
					RelPath: "bar/main.go",
					Src:     `package main`,
				},
			},
			distgo.ProjectParam{
				Products: map[distgo.ProductID]distgo.ProductParam{
					"foo": defaultProductParam("foo", "./foo"),
					"bar": defaultProductParam("bar", "./bar"),
				},
			},
		},
		{
			"generated configuration inherits defaults",
			`
product-defaults:
  build:
    output-dir: default-output
`,
			[]gofiles.GoFileSpec{
				{
					RelPath: "foo/main.go",
					Src:     `package main`,
				},
				{
					RelPath: "bar/main.go",
					Src:     `package main`,
				},
			},
			distgo.ProjectParam{
				Products: map[distgo.ProductID]distgo.ProductParam{
					"foo": defaultProductParam("foo", "./foo", func(param *distgo.ProductParam) {
						param.Build.OutputDir = "default-output"
					}),
					"bar": defaultProductParam("bar", "./bar", func(param *distgo.ProductParam) {
						param.Build.OutputDir = "default-output"
					}),
				},
			},
		},
		{
			"configuration created for main package names are disambiguated",
			`
`,
			[]gofiles.GoFileSpec{
				{
					RelPath: "bar/foo/main.go",
					Src:     `package main`,
				},
				{
					RelPath: "bar/baz-a/main.go",
					Src:     `package main`,
				},
				{
					RelPath: "baz-a/main.go",
					Src:     `package main`,
				},
				// will become "foo-2" because "foo-1" is already taken by a primary package
				{
					RelPath: "foo/main.go",
					Src:     `package main`,
				},
				{
					RelPath: "foo-1/main.go",
					Src:     `package main`,
				},
			},
			distgo.ProjectParam{
				Products: map[distgo.ProductID]distgo.ProductParam{
					"baz-a":   defaultProductParam("baz-a", "./bar/baz-a"),
					"baz-a-1": defaultProductParam("baz-a-1", "./baz-a"),
					"foo":     defaultProductParam("foo", "./bar/foo"),
					"foo-2":   defaultProductParam("foo-2", "./foo"),
					"foo-1":   defaultProductParam("foo-1", "./foo-1"),
				},
			},
		},
	} {
		currProjectDir, err := ioutil.TempDir(tmpDir, "")
		require.NoError(t, err, "Case %d: %s", i, tc.name)

		_, err = gofiles.Write(currProjectDir, tc.gofiles)
		require.NoError(t, err, "Case %d: %s", i, tc.name)

		gotCfg, err := distgo.LoadConfig([]byte(tc.yml))
		require.NoError(t, err, "Case %d: %s", i, tc.name)

		gotParam, err := gotCfg.ToParam(currProjectDir, disterFactory, defaultDisterCfg, dockerBuilderFactory)
		require.NoError(t, err, "Case %d: %s", i, tc.name)
		assert.Equal(t, tc.want, gotParam, "Case %d: %s", i, tc.name)
	}
}

func TestProjectConfig_DockerBuildDepForNonDependentProduct(t *testing.T) {
	for i, tc := range []struct {
		name      string
		yml       string
		wantError string
	}{
		{
			"Docker build dep cannot depend on product that exists but is not a dependent product",
			`
products:
  test-1:
    docker:
      docker-builders:
        default:
          type: default
          context-dir: docker
          input-builds:
            - test-3
          tag-templates:
            - foo:latest
    dependencies:
      - test-2
  test-2:
  test-3:
`,
			`invalid Docker input build(s) specified for DockerBuilderParam "default" for product "test-1"`,
		},
		{
			"Docker dist dep cannot depend on product that exists but is not a dependent product",
			`
products:
  test-1:
    docker:
      docker-builders:
        default:
          type: default
          context-dir: docker
          input-dists:
            - test-3
          tag-templates:
            - foo:latest
    dependencies:
      - test-2
  test-2:
  test-3:
`,
			`invalid Docker input dist(s) specified for DockerBuilderParam "default" for product "test-1"`,
		},
	} {
		gotCfg, err := distgo.LoadConfig([]byte(tc.yml))
		require.NoError(t, err, "Case %d: %s", i, tc.name)
		_, err = gotCfg.ToParam("", disterFactory, defaultDisterCfg, dockerBuilderFactory)
		assert.EqualError(t, err, tc.wantError, "Case %d: %s", i, tc.name)
	}
}

func TestProjectConfig_InvalidDependencies(t *testing.T) {
	for i, tc := range []struct {
		name      string
		yml       string
		wantError string
	}{
		{
			"dependent product that does not exist",
			`
products:
  test-1:
    dependencies:
      - nonexistent-product
`,
			`invalid dependencies for product(s) [test-1]:
  test-1: "nonexistent-product" is not a valid product`,
		},
		{
			"dependent product that is self-referential",
			`
products:
  test-1:
    dependencies:
      - test-1
`,
			`invalid dependencies for product(s) [test-1]:
  test-1: cycle exists: test-1 -> test-1`,
		},
		{
			"product dependencies that generate a cycle",
			`
products:
  test-1:
    dependencies:
      - test-2
  test-2:
    dependencies:
      - test-3
  test-3:
    dependencies:
      - test-1
`,
			`invalid dependencies for product(s) [test-1 test-2 test-3]:
  test-1: cycle exists: test-1 -> test-2 -> test-3 -> test-1
  test-2: cycle exists: test-2 -> test-3 -> test-1 -> test-2
  test-3: cycle exists: test-3 -> test-1 -> test-2 -> test-3`,
		},
	} {
		gotCfg, err := distgo.LoadConfig([]byte(tc.yml))
		require.NoError(t, err, "Case %d: %s", i, tc.name)
		_, err = gotCfg.ToParam("", disterFactory, defaultDisterCfg, dockerBuilderFactory)
		assert.EqualError(t, err, tc.wantError, "Case %d: %s", i, tc.name)
	}
}

func TestProductTaskParam_ToProductTaskOutputInfo(t *testing.T) {
	tmpDir, cleanup, err := dirs.TempDir("", "")
	require.NoError(t, err)
	defer cleanup()

	projectDir := path.Join(tmpDir, "project")
	err = os.Mkdir(projectDir, 0755)
	require.NoError(t, err)

	gittest.InitGitDir(t, projectDir)
	gittest.CreateGitTag(t, projectDir, "1.0.0")

	for i, tc := range []struct {
		name string
		yml  string
		want map[distgo.ProductID]distgo.ProductTaskOutputInfo
	}{
		{
			"name template rendered",
			`
products:
  test-one:
    build:
      name-template: "{{Product}}-{{Version}}-cli"
`,
			map[distgo.ProductID]distgo.ProductTaskOutputInfo{
				"test-one": {
					Project: distgo.ProjectInfo{
						ProjectDir: projectDir,
						Version:    "1.0.0",
					},
					Product: distgo.ProductOutputInfo{
						ID: "test-one",
						BuildOutputInfo: &distgo.BuildOutputInfo{
							BuildNameTemplateRendered: "test-one-1.0.0-cli",
							BuildOutputDir:            "out/build",
							OSArchs: []osarch.OSArch{
								osarch.Current(),
							},
						},
					},
				},
			},
		},
		{
			"custom version script executed",
			`
version-script: |
                 #!/usr/bin/env bash
                 echo "custom-version"
products:
  test-one:
    name: test-product-name
    build:
      name-template: "{{Product}}-{{Version}}"
`,
			map[distgo.ProductID]distgo.ProductTaskOutputInfo{
				"test-one": {
					Project: distgo.ProjectInfo{
						ProjectDir: projectDir,
						Version:    "custom-version",
					},
					Product: distgo.ProductOutputInfo{
						ID: "test-one",
						BuildOutputInfo: &distgo.BuildOutputInfo{
							BuildNameTemplateRendered: "test-one-custom-version",
							BuildOutputDir:            "out/build",
							OSArchs: []osarch.OSArch{
								osarch.Current(),
							},
						},
					},
				},
			},
		},
	} {
		gotCfg, err := distgo.LoadConfig([]byte(tc.yml))
		require.NoError(t, err, "Case %d: %s", i, tc.name)
		gotParam, err := gotCfg.ToParam("", disterFactory, defaultDisterCfg, dockerBuilderFactory)
		require.NoError(t, err, "Case %d: %s", i, tc.name)
		projectInfo, err := gotParam.ProjectInfo(projectDir)
		require.NoError(t, err, "Case %d: %s", i, tc.name)

		require.NoError(t, err, "Case %d: %s", i, tc.name)
		got := make(map[distgo.ProductID]distgo.ProductTaskOutputInfo)
		for k, v := range gotParam.Products {
			currInfo, err := distgo.ToProductTaskOutputInfo(projectInfo, v)
			require.NoError(t, err, "Case %d: %s", i, tc.name)
			got[k] = currInfo
		}
		assert.Equal(t, tc.want, got, "Case %d: %s", i, tc.name)
	}
}

func stringPtr(val string) *string {
	return &val
}

func mustOSArch(in string) osarch.OSArch {
	osArch, err := osarch.New(in)
	if err != nil {
		panic(err)
	}
	return osArch
}
