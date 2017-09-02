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

package docker_test

import (
	"context"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"path"
	"testing"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	dockercli "github.com/docker/docker/client"
	"github.com/nmiyake/pkg/dirs"
	"github.com/stretchr/testify/require"

	"github.com/palantir/godel/apps/distgo/cmd/build"
	"github.com/palantir/godel/apps/distgo/cmd/dist"
	"github.com/palantir/godel/apps/distgo/cmd/docker"
	"github.com/palantir/godel/apps/distgo/params"
	"github.com/palantir/godel/apps/distgo/pkg/git"
	"github.com/palantir/godel/apps/distgo/pkg/git/gittest"
	"github.com/palantir/godel/apps/distgo/pkg/osarch"
)

const (
	testMain = `package main

import "fmt"

var testVersionVar = "defaultVersion"

func main() {
	fmt.Println(testVersionVar)
}
`
	dockerfile = `FROM alpine:3.5
`
	configFile = `test_key: 'test_value'
test_key_2: test_value_2
`
	slsDepDockerFile = `FROM alpine:3.5

COPY foo-sls.tgz .
COPY foo-bin.tgz .
COPY bar-sls.tgz .
`
	dockerRepoPrefix = "test-docker-dist"
)

func TestDockerDist(t *testing.T) {
	tmp, cleanup, err := dirs.TempDir("", "")
	defer cleanup()
	require.NoError(t, err)

	for i, currCase := range []struct {
		name          string
		spec          func(projectDir string, randomPad string) []params.ProductBuildSpecWithDeps
		setupProject  func(projectDir, pad string) error
		cleanup       func(cli *dockercli.Client, projectDir, pad string)
		preDistAction func(projectDir string, buildSpec []params.ProductBuildSpecWithDeps)
		validate      func(caseNum int, name string, pad string, cli *dockercli.Client)
	}{
		{
			name: "docker dist with dependent images",
			setupProject: func(projectDir, pad string) error {
				gittest.InitGitDir(t, projectDir)
				// initialize foo
				fooDir := path.Join(projectDir, "foo")
				if err := os.Mkdir(fooDir, 0777); err != nil {
					return err
				}
				if err := ioutil.WriteFile(path.Join(fooDir, "main.go"), []byte(testMain), 0644); err != nil {
					return err
				}
				fooDockerDir := path.Join(fooDir, "docker")
				if err = os.Mkdir(fooDockerDir, 0777); err != nil {
					return err
				}
				fooDockerFile := fmt.Sprintf("FROM %v:0.1.0\n", fullRepoName("bar", pad))
				if err = ioutil.WriteFile(path.Join(fooDockerDir, "Dockerfile"), []byte(fooDockerFile), 0777); err != nil {
					return err
				}

				// initialize bar
				barDir := path.Join(projectDir, "bar")
				if err := os.Mkdir(barDir, 0777); err != nil {
					return err
				}
				if err = ioutil.WriteFile(path.Join(barDir, "main.go"), []byte(testMain), 0644); err != nil {
					return err
				}
				barDockerDir := path.Join(barDir, "docker")
				if err = os.Mkdir(barDockerDir, 0777); err != nil {
					return err
				}
				if err = ioutil.WriteFile(path.Join(barDockerDir, "Dockerfile"), []byte(dockerfile), 0777); err != nil {
					return err
				}

				// commit
				gittest.CommitAllFiles(t, projectDir, "Commit")
				return nil
			},
			spec: func(projectDir string, pad string) []params.ProductBuildSpecWithDeps {
				allSpec := make(map[string]params.ProductBuildSpec)
				barSpec := params.NewProductBuildSpec(
					projectDir,
					"bar",
					git.ProjectInfo{
						Version: "0.1.0",
					},
					params.Product{
						Build: params.Build{
							MainPkg: "./bar",
							OSArchs: []osarch.OSArch{
								{
									OS:   "linux",
									Arch: "amd64",
								},
							},
						},
						DockerImages: []params.DockerImage{
							{
								Repository: fullRepoName("bar", pad),
								Tag:        "0.1.0",
								ContextDir: "bar/docker",
								Info:       &params.DefaultDockerImageInfo{},
							},
						},
					}, params.Project{
						GroupID: "com.test.group",
					},
				)
				allSpec["bar"] = barSpec
				barSpecWithDeps, err := params.NewProductBuildSpecWithDeps(barSpec, allSpec)
				require.NoError(t, err)
				fooSpec := params.NewProductBuildSpec(
					projectDir,
					"foo",
					git.ProjectInfo{
						Version: "0.1.0",
					},
					params.Product{
						Build: params.Build{
							MainPkg: "./foo",
							OSArchs: []osarch.OSArch{
								{
									OS:   "linux",
									Arch: "amd64",
								},
							},
						},
						DockerImages: []params.DockerImage{
							{
								Repository: fullRepoName("foo", pad),
								Tag:        "0.1.0",
								ContextDir: "foo/docker",
								Deps: []params.DockerDep{
									{
										Product:    "bar",
										Type:       params.DockerDepDocker,
										TargetFile: "",
									},
								},
								Info: &params.DefaultDockerImageInfo{},
							},
						},
					},
					params.Project{
						GroupID: "com.test.group",
					},
				)
				allSpec["foo"] = fooSpec
				fooSpecWithDeps, err := params.NewProductBuildSpecWithDeps(fooSpec, allSpec)
				require.NoError(t, err)

				return []params.ProductBuildSpecWithDeps{fooSpecWithDeps, barSpecWithDeps}
			},
			cleanup: func(cli *dockercli.Client, projectDir, pad string) {
				images := []string{fmt.Sprintf("%v:0.1.0", fullRepoName("foo", pad)),
					fmt.Sprintf("%v:0.1.0", fullRepoName("bar", pad))}
				err := removeImages(cli, images)
				if err != nil {
					t.Logf("Failed to remove images: %v", err)
				}
			},
			preDistAction: func(projectDir string, buildSpec []params.ProductBuildSpecWithDeps) {
				gittest.CreateGitTag(t, projectDir, "0.1.0")
			},
			validate: func(caseNum int, name string, pad string, cli *dockercli.Client) {
				filter := filters.NewArgs()
				filter.Add("reference", fmt.Sprintf("%v:0.1.0", fullRepoName("foo", pad)))
				images, err := cli.ImageList(context.Background(), types.ImageListOptions{Filters: filter})
				require.NoError(t, err, "Case %d: %s", caseNum, name)
				require.True(t, len(images) > 0, "Case %d: %s", caseNum, name)
				filter = filters.NewArgs()
				filter.Add("reference", fmt.Sprintf("%v:0.1.0", fullRepoName("bar", pad)))
				images, err = cli.ImageList(context.Background(), types.ImageListOptions{Filters: filter})
				require.NoError(t, err, "Case %d: %s", caseNum, name)
				require.True(t, len(images) > 0, "Case %d: %s", caseNum, name)
			},
		},
		{
			name: "docker dist with build script args",
			setupProject: func(projectDir, pad string) error {
				gittest.InitGitDir(t, projectDir)
				// initialize foo
				fooDir := path.Join(projectDir, "foo")
				if err := os.Mkdir(fooDir, 0777); err != nil {
					return err
				}
				if err := ioutil.WriteFile(path.Join(fooDir, "main.go"), []byte(testMain), 0644); err != nil {
					return err
				}
				fooDockerDir := path.Join(fooDir, "docker")
				if err = os.Mkdir(fooDockerDir, 0777); err != nil {
					return err
				}
				if err = ioutil.WriteFile(path.Join(fooDockerDir, "Dockerfile"), []byte(dockerfile), 0777); err != nil {
					return err
				}

				// commit
				gittest.CommitAllFiles(t, projectDir, "Commit")
				return nil
			},
			spec: func(projectDir string, pad string) []params.ProductBuildSpecWithDeps {
				allSpec := make(map[string]params.ProductBuildSpec)
				fooSpec := params.NewProductBuildSpec(
					projectDir,
					"foo",
					git.ProjectInfo{
						Version: "0.1.0",
					},
					params.Product{
						Build: params.Build{
							MainPkg: "./foo",
							OSArchs: []osarch.OSArch{
								{
									OS:   "linux",
									Arch: "amd64",
								},
							},
						},
						DockerImages: []params.DockerImage{
							{
								Repository: fullRepoName("foo", pad),
								Tag:        "0.1.0",
								ContextDir: "foo/docker",
								BuildArgsScript: `echo "--label"
echo "test_label=test_value"`,
								Info: &params.DefaultDockerImageInfo{},
							},
						},
					},
					params.Project{
						GroupID: "com.test.group",
					},
				)
				allSpec["foo"] = fooSpec
				fooSpecWithDeps, err := params.NewProductBuildSpecWithDeps(fooSpec, allSpec)
				require.NoError(t, err)

				return []params.ProductBuildSpecWithDeps{fooSpecWithDeps}
			},
			cleanup: func(cli *dockercli.Client, projectDir, pad string) {
				images := []string{fmt.Sprintf("%v:0.1.0", fullRepoName("foo", pad))}
				err := removeImages(cli, images)
				if err != nil {
					t.Logf("Failed to remove images: %v", err)
				}
			},
			preDistAction: func(projectDir string, buildSpec []params.ProductBuildSpecWithDeps) {
				gittest.CreateGitTag(t, projectDir, "0.1.0")
			},
			validate: func(caseNum int, name string, pad string, cli *dockercli.Client) {
				filter := filters.NewArgs()
				filter.Add("reference", fmt.Sprintf("%v:0.1.0", fullRepoName("foo", pad)))
				images, err := cli.ImageList(context.Background(), types.ImageListOptions{Filters: filter})
				require.NoError(t, err, "Case %d: %s", caseNum, name)
				require.True(t, len(images) > 0, "Case %d: %s", caseNum, name)
				image := images[0]
				inspect, _, err := cli.ImageInspectWithRaw(context.Background(), image.ID)
				require.NoError(t, err, "Case %d: %s", caseNum, name)
				actualValue, ok := inspect.Config.Labels["test_label"]
				require.True(t, ok, "Case %d: %s", caseNum, name)
				require.Equal(t, "test_value", actualValue, "Case%d: %s", caseNum, name)
			},
		},
		{
			name: "sls docker dist",
			setupProject: func(projectDir, pad string) error {
				gittest.InitGitDir(t, projectDir)
				// initialize foo
				fooDir := path.Join(projectDir, "foo")
				if err := os.Mkdir(fooDir, 0777); err != nil {
					return err
				}
				if err := ioutil.WriteFile(path.Join(fooDir, "main.go"), []byte(testMain), 0644); err != nil {
					return err
				}
				fooDockerDir := path.Join(fooDir, "docker")
				if err = os.Mkdir(fooDockerDir, 0777); err != nil {
					return err
				}
				if err = ioutil.WriteFile(path.Join(fooDockerDir, "Dockerfile"), []byte(dockerfile), 0777); err != nil {
					return err
				}
				if err = ioutil.WriteFile(path.Join(fooDockerDir, docker.ConfigurationFileName), []byte(configFile), 0777); err != nil {
					return err
				}

				// commit
				gittest.CommitAllFiles(t, projectDir, "Commit")
				return nil
			},
			spec: func(projectDir string, pad string) []params.ProductBuildSpecWithDeps {
				allSpec := make(map[string]params.ProductBuildSpec)
				fooSpec := params.NewProductBuildSpec(
					projectDir,
					"foo",
					git.ProjectInfo{
						Version: "0.1.0",
					},
					params.Product{
						Build: params.Build{
							MainPkg: "./foo",
							OSArchs: []osarch.OSArch{
								{
									OS:   "linux",
									Arch: "amd64",
								},
							},
						},
						DockerImages: []params.DockerImage{
							{
								Repository: fullRepoName("foo", pad),
								Tag:        "0.1.0",
								ContextDir: "foo/docker",
								Deps: []params.DockerDep{
									{
										Product:    "bar",
										Type:       params.DockerDepDocker,
										TargetFile: "",
									},
								},
								Info: &params.SLSDockerImageInfo{
									ProuductType: "test_type",
									GroupID:      "com.palantir.godel",
									Extensions: map[string]interface{}{
										"test_key": "test_value",
									},
								},
							},
						},
					},
					params.Project{
						GroupID: "com.test.group",
					},
				)
				allSpec["foo"] = fooSpec
				fooSpecWithDeps, err := params.NewProductBuildSpecWithDeps(fooSpec, allSpec)
				require.NoError(t, err)

				return []params.ProductBuildSpecWithDeps{fooSpecWithDeps}
			},
			cleanup: func(cli *dockercli.Client, projectDir, pad string) {
				images := []string{fmt.Sprintf("%v:0.1.0", fullRepoName("foo", pad))}
				err := removeImages(cli, images)
				if err != nil {
					t.Logf("Failed to remove images: %v", err)
				}
			},
			preDistAction: func(projectDir string, buildSpec []params.ProductBuildSpecWithDeps) {
				gittest.CreateGitTag(t, projectDir, "0.1.0")
			},
			validate: func(caseNum int, name string, pad string, cli *dockercli.Client) {
				filter := filters.NewArgs()
				filter.Add("reference", fmt.Sprintf("%v:0.1.0", fullRepoName("foo", pad)))
				images, err := cli.ImageList(context.Background(), types.ImageListOptions{Filters: filter})
				require.NoError(t, err, "Case %d: %s", caseNum, name)
				require.True(t, len(images) > 0, "Case %d: %s", caseNum, name)
				image := images[0]
				inspect, _, err := cli.ImageInspectWithRaw(context.Background(), image.ID)
				require.NoError(t, err, "Case %d: %s", caseNum, name)
				actualManifestEncoded, ok := inspect.Config.Labels[docker.ManifestLabel]
				require.True(t, ok, "Case %d: %s", caseNum, name)
				actualManifest, err := base64.StdEncoding.DecodeString(actualManifestEncoded)
				require.NoError(t, err, "Case %d: %s", caseNum, name)
				expectedManifest := "manifest-version: \"1.0\"\n" +
					"product-group: com.palantir.godel\n" +
					"product-name: foo\n" +
					"product-version: 0.1.0\n" +
					"product-type: test_type\n" +
					"extensions:\n  test_key: test_value\n"
				require.Equal(t, expectedManifest, string(actualManifest), "Case %d: %s", caseNum, name)
				actualConfigEncoded, ok := inspect.Config.Labels[docker.ConfigurationLabel]
				require.True(t, ok, "Case %d: %s", caseNum, name)
				actualConfig, err := base64.StdEncoding.DecodeString(actualConfigEncoded)
				require.NoError(t, err, "Case %d: %s", caseNum, name)
				require.Equal(t, configFile, string(actualConfig), "Case %d: %s", caseNum, name)
			},
		},
		{
			name: "docker dist with dependent sls dist",
			setupProject: func(projectDir, pad string) error {
				gittest.InitGitDir(t, projectDir)
				// initialize foo
				fooDir := path.Join(projectDir, "foo")
				if err := os.Mkdir(fooDir, 0777); err != nil {
					return err
				}
				if err := ioutil.WriteFile(path.Join(fooDir, "main.go"), []byte(testMain), 0644); err != nil {
					return err
				}

				// initialize bar
				barDir := path.Join(projectDir, "bar")
				if err := os.Mkdir(barDir, 0777); err != nil {
					return err
				}
				if err = ioutil.WriteFile(path.Join(barDir, "main.go"), []byte(testMain), 0644); err != nil {
					return err
				}
				barDockerDir := path.Join(barDir, "docker")
				if err = os.Mkdir(barDockerDir, 0777); err != nil {
					return err
				}
				if err = ioutil.WriteFile(path.Join(barDockerDir, "Dockerfile"), []byte(slsDepDockerFile), 0777); err != nil {
					return err
				}

				// commit
				gittest.CommitAllFiles(t, projectDir, "Commit")
				return nil
			},
			spec: func(projectDir string, pad string) []params.ProductBuildSpecWithDeps {
				allSpec := make(map[string]params.ProductBuildSpec)
				fooSpec := params.NewProductBuildSpec(
					projectDir,
					"foo",
					git.ProjectInfo{
						Version: "0.1.0",
					},
					params.Product{
						Build: params.Build{
							MainPkg: "./foo",
							OSArchs: []osarch.OSArch{
								{
									OS:   "linux",
									Arch: "amd64",
								},
							},
						},
						Dist: []params.Dist{
							{
								Info: &params.SLSDistInfo{},
							},
							{
								Info: &params.BinDistInfo{},
							},
						},
					},
					params.Project{
						GroupID: "com.test.group",
					},
				)
				allSpec["foo"] = fooSpec
				fooSpecWithDeps, err := params.NewProductBuildSpecWithDeps(fooSpec, allSpec)
				require.NoError(t, err)
				barSpec := params.NewProductBuildSpec(
					projectDir,
					"bar",
					git.ProjectInfo{
						Version: "0.1.0",
					},
					params.Product{
						Build: params.Build{
							MainPkg: "./bar",
							OSArchs: []osarch.OSArch{
								{
									OS:   "linux",
									Arch: "amd64",
								},
							},
						},
						DockerImages: []params.DockerImage{
							{
								Repository: fullRepoName("bar", pad),
								Tag:        "0.1.0",
								ContextDir: "bar/docker",
								Info:       &params.DefaultDockerImageInfo{},
								Deps: []params.DockerDep{
									{
										Product:    "bar",
										Type:       params.DockerDepSLS,
										TargetFile: "bar-sls.tgz",
									},
									{
										Product:    "foo",
										Type:       params.DockerDepBin,
										TargetFile: "foo-bin.tgz",
									},
									{
										Product:    "foo",
										Type:       params.DockerDepSLS,
										TargetFile: "foo-sls.tgz",
									},
								},
							},
						},
						Dist: []params.Dist{
							{
								Info: &params.SLSDistInfo{},
							},
						},
					}, params.Project{
						GroupID: "com.test.group",
					},
				)
				allSpec["bar"] = barSpec
				barSpecWithDeps, err := params.NewProductBuildSpecWithDeps(barSpec, allSpec)
				require.NoError(t, err)
				return []params.ProductBuildSpecWithDeps{fooSpecWithDeps, barSpecWithDeps}
			},
			cleanup: func(cli *dockercli.Client, projectDir, pad string) {
				images := []string{fmt.Sprintf("%v:0.1.0", fullRepoName("bar", pad))}
				err := removeImages(cli, images)
				if err != nil {
					t.Logf("Failed to remove images: %v", err)
				}
			},
			preDistAction: func(projectDir string, buildSpec []params.ProductBuildSpecWithDeps) {
				gittest.CreateGitTag(t, projectDir, "0.1.0")
			},
			validate: func(caseNum int, name string, pad string, cli *dockercli.Client) {
				filter := filters.NewArgs()
				filter.Add("reference", fmt.Sprintf("%v:0.1.0", fullRepoName("bar", pad)))
				images, err := cli.ImageList(context.Background(), types.ImageListOptions{Filters: filter})
				require.NoError(t, err, "Case %d: %s", caseNum, name)
				require.True(t, len(images) > 0, "Case %d: %s", caseNum, name)
			},
		},
	} {
		cli, err := dockercli.NewEnvClient()
		require.NoError(t, err)

		currTmpDir, err := ioutil.TempDir(tmp, "")
		require.NoError(t, err, "Case %d: %s", i, currCase.name)
		pad := randomPad(8)

		err = currCase.setupProject(currTmpDir, pad)
		require.NoError(t, err, "Case %d: %s", i, currCase.name)
		spec := currCase.spec(currTmpDir, pad)

		if currCase.preDistAction != nil {
			currCase.preDistAction(currTmpDir, spec)
		}

		if currCase.cleanup != nil {
			defer currCase.cleanup(cli, currTmpDir, pad)
		}

		for _, currSpecWithDeps := range spec {
			err = build.Run(build.RequiresBuild(currSpecWithDeps, nil).Specs(), nil, build.Context{}, ioutil.Discard)
			require.NoError(t, err, "Case %d: %s", i, currCase.name)
			err = dist.Run(currSpecWithDeps, ioutil.Discard)
			require.NoError(t, err, "Case %d: %s", i, currCase.name)
		}

		orderedSpecs, err := docker.OrderBuildSpecs(spec)
		require.NoError(t, err, "Case %d: %s", i, currCase.name)
		err = docker.RunBuild(orderedSpecs, false, os.Stdout)
		require.NoError(t, err, "Case %d: %s", i, currCase.name)

		if currCase.validate != nil {
			currCase.validate(i, currCase.name, pad, cli)
		}
	}
}

func fullRepoName(product string, pad string) string {
	return fmt.Sprintf("%v-%v-%v", dockerRepoPrefix, product, pad)
}

func randomPad(n int) string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyz1234567890")
	b := make([]rune, n)
	rand.Seed(time.Now().UTC().UnixNano())
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func removeImages(cli *dockercli.Client, images []string) error {
	for _, image := range images {
		_, err := cli.ImageRemove(context.Background(), image, types.ImageRemoveOptions{})
		if err != nil {
			return err
		}
	}
	return nil
}
