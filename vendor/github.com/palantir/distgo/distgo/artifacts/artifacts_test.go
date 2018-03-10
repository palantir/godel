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

package artifacts_test

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"path"
	"path/filepath"
	"testing"
	"time"

	"github.com/nmiyake/pkg/dirs"
	"github.com/nmiyake/pkg/gofiles"
	"github.com/palantir/godel/pkg/osarch"
	"github.com/palantir/pkg/gittest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/palantir/distgo/dister"
	"github.com/palantir/distgo/distgo"
	"github.com/palantir/distgo/distgo/artifacts"
	"github.com/palantir/distgo/distgo/build"
	"github.com/palantir/distgo/dockerbuilder"
)

func TestBuildArtifactsDefaultOutput(t *testing.T) {
	tmpDir, cleanup, err := dirs.TempDir("", "")
	defer cleanup()
	require.NoError(t, err)

	for i, tc := range []struct {
		name            string
		projectConfig   distgo.ProjectConfig
		setupProjectDir func(projectDir string)
		want            func(projectDir string) string
	}{
		{
			"if param is empty, prints main packages in build output directory",
			distgo.ProjectConfig{},
			func(projectDir string) {
				_, err := gofiles.Write(projectDir, []gofiles.GoFileSpec{
					{
						RelPath: "main.go",
						Src:     `package main`,
					},
					{
						RelPath: "bar/bar.go",
						Src:     `package bar`,
					},
					{
						RelPath: "foo/foo.go",
						Src:     `package main`,
					},
				})
				require.NoError(t, err)
			},
			func(projectDir string) string {
				return fmt.Sprintf(`%s/out/build/%s/unspecified/%v/%s
%s/out/build/foo/unspecified/%v/foo
`, projectDir, path.Base(projectDir), osarch.Current(), path.Base(projectDir), projectDir, osarch.Current())
			},
		},
		{
			"output directory specified in param is used",
			distgo.ProjectConfig{
				Products: map[distgo.ProductID]distgo.ProductConfig{
					"foo": {
						Build: &distgo.BuildConfig{
							OutputDir: stringPtr("build-output"),
							OSArchs: &[]osarch.OSArch{
								osarch.Current(),
							},
						},
					},
				},
			},
			nil,
			func(projectDir string) string {
				return fmt.Sprintf(`%s/build-output/foo/unspecified/%v/foo
`, projectDir, osarch.Current())
			},
		},
	} {
		projectDir, err := ioutil.TempDir(tmpDir, "")
		require.NoError(t, err, "Case %d: %s", i, tc.name)

		gittest.InitGitDir(t, projectDir)
		if tc.setupProjectDir != nil {
			tc.setupProjectDir(projectDir)
		}

		disterFactory, err := dister.NewDisterFactory()
		require.NoError(t, err, "Case %d: %s", i, tc.name)

		defaultDisterCfg, err := dister.DefaultConfig()
		require.NoError(t, err, "Case %d: %s", i, tc.name)

		dockerBuilderFactory, err := dockerbuilder.NewDockerBuilderFactory()
		require.NoError(t, err, "Case %d: %s", i, tc.name)

		projectParam, err := tc.projectConfig.ToParam(projectDir, disterFactory, defaultDisterCfg, dockerBuilderFactory)
		require.NoError(t, err, "Case %d: %s", i, tc.name)

		projectInfo, err := projectParam.ProjectInfo(projectDir)
		require.NoError(t, err, "Case %d: %s", i, tc.name)

		buf := &bytes.Buffer{}
		err = artifacts.PrintBuildArtifacts(projectInfo, projectParam, nil, false, false, buf)
		require.NoError(t, err, "Case %d: %s", i, tc.name)
		assert.Equal(t, tc.want(projectDir), buf.String(), "Case %d: %s", i, tc.name)
	}
}

func TestBuildArtifacts(t *testing.T) {
	tmpDir, cleanup, err := dirs.TempDir("", "")
	defer cleanup()
	require.NoError(t, err)

	for i, tc := range []struct {
		params []distgo.ProductParam
		want   func(projectDir string) map[distgo.ProductID][]string
	}{
		// empty spec
		{
			params: []distgo.ProductParam{},
			want: func(projectDir string) map[distgo.ProductID][]string {
				return map[distgo.ProductID][]string{}
			},
		},
		// returns paths for all OS/arch combinations if requested osArchs is empty
		{
			params: []distgo.ProductParam{
				createBuildSpec("foo", []osarch.OSArch{
					{OS: "darwin", Arch: "amd64"},
					{OS: "darwin", Arch: "386"},
					{OS: "linux", Arch: "amd64"},
				}),
			},
			want: func(projectDir string) map[distgo.ProductID][]string {
				return map[distgo.ProductID][]string{
					"foo": {
						path.Join(projectDir, "out", "build", "foo", "0.1.0", "darwin-386", "foo"),
						path.Join(projectDir, "out", "build", "foo", "0.1.0", "darwin-amd64", "foo"),
						path.Join(projectDir, "out", "build", "foo", "0.1.0", "linux-amd64", "foo"),
					},
				}
			},
		},
		// path to windows executable includes ".exe"
		{
			params: []distgo.ProductParam{
				createBuildSpec("foo", []osarch.OSArch{
					{OS: "windows", Arch: "amd64"},
				}),
			},
			want: func(projectDir string) map[distgo.ProductID][]string {
				return map[distgo.ProductID][]string{
					"foo": {
						path.Join(projectDir, "out", "build", "foo", "0.1.0", "windows-amd64", "foo.exe"),
					},
				}
			},
		},
	} {
		currProjectDir, err := ioutil.TempDir(tmpDir, "")
		require.NoError(t, err)

		projectInfo := distgo.ProjectInfo{
			ProjectDir: currProjectDir,
			Version:    "0.1.0",
		}
		got, err := artifacts.Build(projectInfo, tc.params, false)
		require.NoError(t, err, "Case %d", i)
		assert.Equal(t, tc.want(currProjectDir), got, "Case %d", i)
	}
}

func TestBuildArtifactsRequiresBuild(t *testing.T) {
	tmpDir, cleanup, err := dirs.TempDir(".", "")
	defer cleanup()
	require.NoError(t, err)

	tmpDir, err = filepath.Abs(tmpDir)
	require.NoError(t, err)

	for i, tc := range []struct {
		params        []distgo.ProductParam
		requiresBuild bool
		beforeAction  func(projectInfo distgo.ProjectInfo, productParams []distgo.ProductParam)
		want          func(projectDir string) map[distgo.ProductID][]string
	}{
		// returns paths to all artifacts if build has not happened
		{
			params: []distgo.ProductParam{
				createBuildSpec("foo", []osarch.OSArch{
					{OS: "darwin", Arch: "amd64"},
					{OS: "darwin", Arch: "386"},
					{OS: "linux", Arch: "amd64"},
				}),
			},
			want: func(projectDir string) map[distgo.ProductID][]string {
				return map[distgo.ProductID][]string{
					"foo": {
						path.Join(projectDir, "out", "build", "foo", "0.1.0", "darwin-386", "foo"),
						path.Join(projectDir, "out", "build", "foo", "0.1.0", "darwin-amd64", "foo"),
						path.Join(projectDir, "out", "build", "foo", "0.1.0", "linux-amd64", "foo"),
					},
				}
			},
		},
		// returns empty if all artifacts exist and are up-to-date
		{
			params: []distgo.ProductParam{
				createBuildSpec("foo", []osarch.OSArch{
					{OS: "darwin", Arch: "amd64"},
					{OS: "darwin", Arch: "386"},
					{OS: "linux", Arch: "amd64"},
				}),
			},
			beforeAction: func(projectInfo distgo.ProjectInfo, productParams []distgo.ProductParam) {
				// build products
				err := build.Run(projectInfo, productParams, build.Options{
					Parallel: false,
				}, ioutil.Discard)
				require.NoError(t, err)
			},
			want: func(projectDir string) map[distgo.ProductID][]string {
				return map[distgo.ProductID][]string{}
			},
		},
		// returns paths to all artifacts if input source file has been modified
		{
			params: []distgo.ProductParam{
				createBuildSpec("foo", []osarch.OSArch{
					{OS: "darwin", Arch: "amd64"},
					{OS: "darwin", Arch: "386"},
					{OS: "linux", Arch: "amd64"},
				}),
			},
			beforeAction: func(projectInfo distgo.ProjectInfo, params []distgo.ProductParam) {
				// build products
				err := build.Run(projectInfo, params, build.Options{
					Parallel: false,
				}, ioutil.Discard)
				require.NoError(t, err)

				// sleep to ensure that modification time will differ
				time.Sleep(time.Second)

				// update source file
				err = ioutil.WriteFile(path.Join(projectInfo.ProjectDir, "main.go"), []byte("package main; func main(){}"), 0644)
				require.NoError(t, err)
			},
			want: func(projectDir string) map[distgo.ProductID][]string {
				return map[distgo.ProductID][]string{
					"foo": {
						path.Join(projectDir, "out", "build", "foo", "0.1.0", "darwin-386", "foo"),
						path.Join(projectDir, "out", "build", "foo", "0.1.0", "darwin-amd64", "foo"),
						path.Join(projectDir, "out", "build", "foo", "0.1.0", "linux-amd64", "foo"),
					},
				}
			},
		},
	} {
		currProjectDir, err := ioutil.TempDir(tmpDir, "")
		require.NoError(t, err)

		err = ioutil.WriteFile(path.Join(currProjectDir, "main.go"), []byte("package main; func main(){}"), 0644)
		require.NoError(t, err)

		projectInfo := distgo.ProjectInfo{
			ProjectDir: currProjectDir,
			Version:    "0.1.0",
		}
		if tc.beforeAction != nil {
			tc.beforeAction(projectInfo, tc.params)
		}

		got, err := artifacts.Build(projectInfo, tc.params, true)
		require.NoError(t, err, "Case %d", i)
		assert.Equal(t, tc.want(currProjectDir), got, "Case %d", i)
	}
}

func TestDistArtifacts(t *testing.T) {
	tmpDir, cleanup, err := dirs.TempDir("", "")
	defer cleanup()
	require.NoError(t, err)

	for i, tc := range []struct {
		params []distgo.ProductParam
		want   func(projectDir string) map[distgo.ProductID][]string
	}{
		// empty spec
		{
			params: []distgo.ProductParam{},
			want: func(projectDir string) map[distgo.ProductID][]string {
				return map[distgo.ProductID][]string{}
			},
		},
		// returns dist artifact outputs
		{
			params: []distgo.ProductParam{
				createDistSpec("foo", dister.NewOSArchBinDister(
					osarch.OSArch{OS: "darwin", Arch: "amd64"},
					osarch.OSArch{OS: "linux", Arch: "amd64"},
				)),
			},
			want: func(projectDir string) map[distgo.ProductID][]string {
				return map[distgo.ProductID][]string{
					"foo": {
						path.Join(projectDir, "out", "dist", "foo", "0.1.0", "os-arch-bin", "foo-0.1.0-darwin-amd64.tgz"),
						path.Join(projectDir, "out", "dist", "foo", "0.1.0", "os-arch-bin", "foo-0.1.0-linux-amd64.tgz"),
					},
				}
			},
		},
	} {
		currProjectDir, err := ioutil.TempDir(tmpDir, "")
		require.NoError(t, err)

		projectInfo := distgo.ProjectInfo{
			ProjectDir: currProjectDir,
			Version:    "0.1.0",
		}
		got, err := artifacts.Dist(projectInfo, tc.params)
		require.NoError(t, err, "Case %d", i)
		assert.Equal(t, tc.want(currProjectDir), got, "Case %d", i)
	}
}

func TestDockerArtifacts(t *testing.T) {
	tmpDir, cleanup, err := dirs.TempDir("", "")
	defer cleanup()
	require.NoError(t, err)

	for i, tc := range []struct {
		name string
		cfg  distgo.ProjectConfig
		want map[distgo.ProductID][]string
	}{
		{
			"prints docker artifacts",
			distgo.ProjectConfig{
				Products: map[distgo.ProductID]distgo.ProductConfig{
					"foo": {
						Docker: &distgo.DockerConfig{
							Repository: stringPtr("repo"),
							DockerBuildersConfig: &distgo.DockerBuildersConfig{
								dockerbuilder.DefaultBuilderTypeName: distgo.DockerBuilderConfig{
									Type:       stringPtr(dockerbuilder.DefaultBuilderTypeName),
									ContextDir: stringPtr("dockerContextDir"),
									TagTemplates: &[]string{
										"{{Repository}}/foo:latest",
									},
								},
							},
						},
					},
				},
			},
			map[distgo.ProductID][]string{
				"foo": {
					"repo/foo:latest",
				},
			},
		},
	} {
		projectDir, err := ioutil.TempDir(tmpDir, "")
		require.NoError(t, err)
		gittest.InitGitDir(t, projectDir)
		gittest.CreateGitTag(t, projectDir, "0.1.0")

		disterFactory, err := dister.NewDisterFactory()
		require.NoError(t, err, "Case %d: %s", i, tc.name)
		defaultDisterCfg, err := dister.DefaultConfig()
		require.NoError(t, err, "Case %d: %s", i, tc.name)
		dockerBuilderFactory, err := dockerbuilder.NewDockerBuilderFactory()
		require.NoError(t, err, "Case %d: %s", i, tc.name)

		projectParam, err := tc.cfg.ToParam(projectDir, disterFactory, defaultDisterCfg, dockerBuilderFactory)
		require.NoError(t, err, "Case %d: %s", i, tc.name)

		projectInfo, err := projectParam.ProjectInfo(projectDir)
		require.NoError(t, err, "Case %d: %s", i, tc.name)

		products, err := distgo.ProductParamsForProductArgs(projectParam.Products)
		require.NoError(t, err, "Case %d: %s", i, tc.name)

		dockerArtifacts, err := artifacts.Docker(projectInfo, products)
		require.NoError(t, err, "Case %d: %s", i, tc.name)

		assert.Equal(t, tc.want, dockerArtifacts, "Case %d: %s", i, tc.name)
	}
}

func createBuildSpec(productName string, osArchs []osarch.OSArch) distgo.ProductParam {
	return distgo.ProductParam{
		ID: distgo.ProductID(productName),
		Build: &distgo.BuildParam{
			NameTemplate: "{{Product}}",
			OutputDir:    "out/build",
			MainPkg:      ".",
			OSArchs:      osArchs,
		},
	}
}

func createDistSpec(productName string, dister distgo.Dister) distgo.ProductParam {
	disterName, err := dister.TypeName()
	if err != nil {
		panic(err)
	}

	return distgo.ProductParam{
		ID: distgo.ProductID(productName),
		Dist: &distgo.DistParam{
			OutputDir: "out/dist",
			DistParams: map[distgo.DistID]distgo.DisterParam{
				distgo.DistID(disterName): {
					NameTemplate: "{{Product}}-{{Version}}",
					Dister:       dister,
				},
			},
		},
	}
}

func stringPtr(in string) *string {
	return &in
}
