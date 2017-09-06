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
	"io/ioutil"
	"path"
	"path/filepath"
	"testing"
	"time"

	"github.com/nmiyake/pkg/dirs"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/palantir/godel/apps/distgo/cmd/artifacts"
	"github.com/palantir/godel/apps/distgo/cmd/build"
	"github.com/palantir/godel/apps/distgo/params"
	"github.com/palantir/godel/apps/distgo/pkg/osarch"
)

func TestBuildArtifacts(t *testing.T) {
	tmpDir, cleanup, err := dirs.TempDir("", "")
	defer cleanup()
	require.NoError(t, err)

	for i, currCase := range []struct {
		specs   func(projectDir string) []params.ProductBuildSpecWithDeps
		osArchs []osarch.OSArch
		want    map[string][]string
	}{
		// empty spec
		{
			specs: func(projectDir string) []params.ProductBuildSpecWithDeps {
				return []params.ProductBuildSpecWithDeps{}
			},
			want: map[string][]string{},
		},
		// returns paths for all OS/arch combinations if requested osArchs is empty
		{
			specs: func(projectDir string) []params.ProductBuildSpecWithDeps {
				return []params.ProductBuildSpecWithDeps{
					createSpec(projectDir, "foo", "0.1.0", []osarch.OSArch{
						{OS: "darwin", Arch: "amd64"},
						{OS: "darwin", Arch: "386"},
						{OS: "linux", Arch: "amd64"},
					}, &params.SLSDistInfo{}),
				}
			},
			want: map[string][]string{
				"foo": {
					path.Join("build", "darwin-amd64", "foo"),
					path.Join("build", "darwin-386", "foo"),
					path.Join("build", "linux-amd64", "foo"),
				},
			},
		},
		// returns only path to requested OS/arch
		{
			specs: func(projectDir string) []params.ProductBuildSpecWithDeps {
				return []params.ProductBuildSpecWithDeps{
					createSpec(projectDir, "foo", "0.1.0", []osarch.OSArch{
						{OS: "darwin", Arch: "amd64"},
						{OS: "linux", Arch: "amd64"},
					}, &params.SLSDistInfo{}),
				}
			},
			osArchs: []osarch.OSArch{{OS: "darwin", Arch: "amd64"}},
			want: map[string][]string{
				"foo": {
					path.Join("build", "darwin-amd64", "foo"),
				},
			},
		},
		// path to windows executable includes ".exe"
		{
			specs: func(projectDir string) []params.ProductBuildSpecWithDeps {
				return []params.ProductBuildSpecWithDeps{
					createSpec(projectDir, "foo", "0.1.0", []osarch.OSArch{
						{OS: "windows", Arch: "amd64"},
					}, &params.SLSDistInfo{}),
				}
			},
			want: map[string][]string{
				"foo": {
					path.Join("build", "windows-amd64", "foo.exe"),
				},
			},
		},
		// returns empty if os/arch that is not part of the spec is requested
		{
			specs: func(projectDir string) []params.ProductBuildSpecWithDeps {
				return []params.ProductBuildSpecWithDeps{
					createSpec(projectDir, "foo", "0.1.0", []osarch.OSArch{
						{OS: "darwin", Arch: "amd64"},
						{OS: "linux", Arch: "amd64"},
					}, &params.SLSDistInfo{}),
				}
			},
			osArchs: []osarch.OSArch{{OS: "windows", Arch: "amd64"}},
			want:    map[string][]string{},
		},
	} {
		currProjectDir, err := ioutil.TempDir(tmpDir, "")
		require.NoError(t, err)

		// relative path
		got, err := artifacts.BuildArtifacts(currCase.specs(currProjectDir), artifacts.BuildArtifactsParams{
			OSArchs: currCase.osArchs,
		})
		require.NoError(t, err, "Case %d", i)
		assert.Equal(t, currCase.want, toMap(got), "Case %d", i)

		// absolute path
		got, err = artifacts.BuildArtifacts(currCase.specs(currProjectDir), artifacts.BuildArtifactsParams{
			AbsPath: true,
			OSArchs: currCase.osArchs,
		})
		require.NoError(t, err, "Case %d", i)
		assert.Equal(t, toAbs(currCase.want, currProjectDir), toMap(got), "Case %d", i)
	}
}

func TestBuildArtifactsRequiresBuild(t *testing.T) {
	tmpDir, cleanup, err := dirs.TempDir(".", "")
	defer cleanup()
	require.NoError(t, err)

	tmpDir, err = filepath.Abs(tmpDir)
	require.NoError(t, err)

	for i, currCase := range []struct {
		specs         func(projectDir string) params.ProductBuildSpecWithDeps
		osArchs       []osarch.OSArch
		requiresBuild bool
		beforeAction  func(projectDir string, specs []params.ProductBuildSpec)
		want          map[string][]string
	}{
		// returns paths to all artifacts if build has not happened
		{
			specs: func(projectDir string) params.ProductBuildSpecWithDeps {
				return createSpec(projectDir, "foo", "0.1.0", []osarch.OSArch{
					{OS: "darwin", Arch: "amd64"},
					{OS: "darwin", Arch: "386"},
					{OS: "linux", Arch: "amd64"},
				}, &params.SLSDistInfo{})
			},
			want: map[string][]string{
				"foo": {
					path.Join("build", "darwin-amd64", "foo"),
					path.Join("build", "darwin-386", "foo"),
					path.Join("build", "linux-amd64", "foo"),
				},
			},
		},
		// returns empty if all artifacts exist and are up-to-date
		{
			specs: func(projectDir string) params.ProductBuildSpecWithDeps {
				return createSpec(projectDir, "foo", "0.1.0", []osarch.OSArch{
					{OS: "darwin", Arch: "amd64"},
					{OS: "darwin", Arch: "386"},
					{OS: "linux", Arch: "amd64"},
				}, &params.SLSDistInfo{})
			},
			beforeAction: func(projectDir string, specs []params.ProductBuildSpec) {
				// build products
				err = build.Run(specs, nil, build.Context{
					Parallel: false,
				}, ioutil.Discard)
				require.NoError(t, err)
			},
			want: map[string][]string{},
		},
		// returns paths to all artifacts if input source file has been modified
		{
			specs: func(projectDir string) params.ProductBuildSpecWithDeps {
				return createSpec(projectDir, "foo", "0.1.0", []osarch.OSArch{
					{OS: "darwin", Arch: "amd64"},
					{OS: "darwin", Arch: "386"},
					{OS: "linux", Arch: "amd64"},
				}, &params.SLSDistInfo{})
			},
			beforeAction: func(projectDir string, specs []params.ProductBuildSpec) {
				// build products
				err := build.Run(specs, nil, build.Context{
					Parallel: false,
				}, ioutil.Discard)
				require.NoError(t, err)

				// sleep to ensure that modification time will differ
				time.Sleep(time.Second)

				// update source file
				err = ioutil.WriteFile(path.Join(projectDir, "main.go"), []byte("package main; func main(){}"), 0644)
				require.NoError(t, err)
			},
			want: map[string][]string{
				"foo": {
					path.Join("build", "darwin-amd64", "foo"),
					path.Join("build", "darwin-386", "foo"),
					path.Join("build", "linux-amd64", "foo"),
				},
			},
		},
		// if OS/Archs are specified, results are filtered base on that
		{
			specs: func(projectDir string) params.ProductBuildSpecWithDeps {
				return createSpec(projectDir, "foo", "0.1.0", []osarch.OSArch{
					{OS: "darwin", Arch: "amd64"},
					{OS: "darwin", Arch: "386"},
					{OS: "linux", Arch: "amd64"},
				}, &params.SLSDistInfo{})
			},
			osArchs: []osarch.OSArch{
				{OS: "windows", Arch: "amd64"},
			},
			want: map[string][]string{},
		},
	} {
		currProjectDir, err := ioutil.TempDir(tmpDir, "")
		require.NoError(t, err)

		err = ioutil.WriteFile(path.Join(currProjectDir, "main.go"), []byte("package main; func main(){}"), 0644)
		require.NoError(t, err)

		specWithDeps := currCase.specs(currProjectDir)
		if currCase.beforeAction != nil {
			currCase.beforeAction(currProjectDir, specWithDeps.AllSpecs())
		}

		got, err := artifacts.BuildArtifacts([]params.ProductBuildSpecWithDeps{specWithDeps}, artifacts.BuildArtifactsParams{
			RequiresBuild: true,
			OSArchs:       currCase.osArchs,
		})
		require.NoError(t, err, "Case %d", i)
		assert.Equal(t, currCase.want, toMap(got), "Case %d", i)
	}
}

func TestDistArtifacts(t *testing.T) {
	tmpDir, cleanup, err := dirs.TempDir("", "")
	defer cleanup()
	require.NoError(t, err)

	for i, currCase := range []struct {
		specs func(projectDir string) []params.ProductBuildSpecWithDeps
		want  map[string][]string
	}{
		{
			specs: func(projectDir string) []params.ProductBuildSpecWithDeps {
				return []params.ProductBuildSpecWithDeps{}
			},
			want: map[string][]string{},
		},
		{
			specs: func(projectDir string) []params.ProductBuildSpecWithDeps {
				return []params.ProductBuildSpecWithDeps{
					createSpec(projectDir, "foo", "0.1.0", nil, &params.SLSDistInfo{}),
				}
			},
			want: map[string][]string{
				"foo": {"foo-0.1.0.sls.tgz"},
			},
		},
		{
			specs: func(projectDir string) []params.ProductBuildSpecWithDeps {
				return []params.ProductBuildSpecWithDeps{
					createSpec(projectDir, "foo", "0.1.0", nil, &params.SLSDistInfo{}),
					createSpec(projectDir, "bar", "unspecified", nil, &params.SLSDistInfo{}),
				}
			},
			want: map[string][]string{
				"foo": {"foo-0.1.0.sls.tgz"},
				"bar": {"bar-unspecified.sls.tgz"},
			},
		},
		{
			specs: func(projectDir string) []params.ProductBuildSpecWithDeps {
				return []params.ProductBuildSpecWithDeps{
					createSpecWithDists(projectDir, "foo", "0.1.0", nil, params.Dist{Info: &params.SLSDistInfo{}}, params.Dist{Info: &params.BinDistInfo{}}),
				}
			},
			want: map[string][]string{
				"foo": {
					"foo-0.1.0.sls.tgz",
					"foo-0.1.0.tgz",
				},
			},
		},
		{
			specs: func(projectDir string) []params.ProductBuildSpecWithDeps {
				return []params.ProductBuildSpecWithDeps{
					createSpecWithDists(projectDir, "foo", "0.1.0", nil, params.Dist{Info: &params.OSArchsBinDistInfo{
						OSArchs: []osarch.OSArch{
							{
								OS:   "darwin",
								Arch: "amd64",
							},
						},
					}}, params.Dist{Info: &params.OSArchsBinDistInfo{
						OSArchs: []osarch.OSArch{
							{
								OS:   "linux",
								Arch: "amd64",
							},
						},
					}}),
				}
			},
			want: map[string][]string{
				"foo": {
					"foo-0.1.0-darwin-amd64.tgz",
					"foo-0.1.0-linux-amd64.tgz",
				},
			},
		},
		{
			specs: func(projectDir string) []params.ProductBuildSpecWithDeps {
				return []params.ProductBuildSpecWithDeps{
					createSpec(projectDir, "foo", "0.1.0", nil, &params.OSArchsBinDistInfo{
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
					}),
				}
			},
			want: map[string][]string{
				"foo": {
					"foo-0.1.0-darwin-amd64.tgz",
					"foo-0.1.0-linux-amd64.tgz",
				},
			},
		},
	} {
		currProjectDir, err := ioutil.TempDir(tmpDir, "")
		require.NoError(t, err)

		// relative path
		got, err := artifacts.DistArtifacts(currCase.specs(currProjectDir), false)
		require.NoError(t, err, "Case %d", i)
		assert.Equal(t, currCase.want, toMap(got), "Case %d", i)

		// absolute path
		got, err = artifacts.DistArtifacts(currCase.specs(currProjectDir), true)
		require.NoError(t, err, "Case %d", i)
		assert.Equal(t, toAbs(currCase.want, currProjectDir), toMap(got), "Case %d", i)
	}
}

func TestDockerArtifacts(t *testing.T) {
	tmpDir, cleanup, err := dirs.TempDir("", "")
	defer cleanup()
	require.NoError(t, err)

	for i, currCase := range []struct {
		specs func(projectDir string) []params.ProductBuildSpecWithDeps
		want  map[string][]string
	}{
		{
			specs: func(projectDir string) []params.ProductBuildSpecWithDeps {
				return []params.ProductBuildSpecWithDeps{}
			},
			want: map[string][]string{},
		},
		{
			specs: func(projectDir string) []params.ProductBuildSpecWithDeps {
				return []params.ProductBuildSpecWithDeps{
					createSpecWithDockerImages(projectDir, "foo", "0.1.0", nil,
						params.DockerImage{
							Repository: "foo/foo",
							Tag:        "snapshot",
						},
					),
				}
			},
			want: map[string][]string{
				"foo": {"foo/foo:snapshot"},
			},
		},
		{
			specs: func(projectDir string) []params.ProductBuildSpecWithDeps {
				return []params.ProductBuildSpecWithDeps{
					createSpecWithDockerImages(projectDir, "foo", "0.1.0", nil,
						params.DockerImage{
							Repository: "foo/foo",
							Tag:        "snapshot",
						},
					),
					createSpecWithDockerImages(projectDir, "bar", "0.1.0", nil,
						params.DockerImage{
							Repository: "bar/bar",
							Tag:        "snapshot",
						},
					),
				}
			},
			want: map[string][]string{
				"foo": {"foo/foo:snapshot"},
				"bar": {"bar/bar:snapshot"},
			},
		},
		{
			specs: func(projectDir string) []params.ProductBuildSpecWithDeps {
				return []params.ProductBuildSpecWithDeps{
					createSpecWithDockerImages(projectDir, "foo", "0.1.0", nil,
						params.DockerImage{
							Repository: "foo/foo",
							Tag:        "snapshot",
						},
						params.DockerImage{
							Repository: "foo/foo",
							Tag:        "release",
						},
					),
				}
			},
			want: map[string][]string{
				"foo": {
					"foo/foo:snapshot",
					"foo/foo:release",
				},
			},
		},
	} {
		currProjectDir, err := ioutil.TempDir(tmpDir, "")
		require.NoError(t, err)

		got := artifacts.DockerArtifacts(currCase.specs(currProjectDir))
		require.NoError(t, err, "Case %d", i)
		assert.Equal(t, currCase.want, got, "Case %d", i)
	}
}

func toAbs(input map[string][]string, baseDir string) map[string][]string {
	absWant := make(map[string][]string, len(input))
	for k, v := range input {
		absWant[k] = make([]string, len(v))
		for i := range v {
			absWant[k][i] = path.Join(baseDir, v[i])
		}
	}
	return absWant
}

func createSpec(projectDir, productName, productVersion string, osArchs []osarch.OSArch, distInfo params.DistInfo) params.ProductBuildSpecWithDeps {
	return createSpecWithDists(projectDir, productName, productVersion, osArchs, params.Dist{Info: distInfo})
}

func createSpecWithDists(projectDir, productName, productVersion string, osArchs []osarch.OSArch, dists ...params.Dist) params.ProductBuildSpecWithDeps {
	return params.ProductBuildSpecWithDeps{
		Spec: params.ProductBuildSpec{
			Product: params.Product{
				Build: params.Build{
					OutputDir: "build",
					OSArchs:   osArchs,
				},
				Dist: dists,
			},
			ProjectDir:     projectDir,
			ProductName:    productName,
			ProductVersion: productVersion,
		},
	}
}

func createSpecWithDockerImages(projectDir, productName, productVersion string, osArchs []osarch.OSArch, images ...params.DockerImage) params.ProductBuildSpecWithDeps {
	return params.ProductBuildSpecWithDeps{
		Spec: params.ProductBuildSpec{
			Product: params.Product{
				Build: params.Build{
					OutputDir: "build",
					OSArchs:   osArchs,
				},
				DockerImages: images,
			},
			ProjectDir:     projectDir,
			ProductName:    productName,
			ProductVersion: productVersion,
		},
	}
}

func toMap(input map[string]artifacts.OrderedStringSliceMap) map[string][]string {
	output := make(map[string][]string, len(input))
	for product, m := range input {
		keys := m.Keys()
		var values []string
		for _, k := range keys {
			values = append(values, m.Get(k)...)
		}
		output[product] = values
	}
	return output
}
