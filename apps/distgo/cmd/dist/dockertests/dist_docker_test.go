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

package dockertests

import (
	"context"
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
	dockerRepoPrefix = "test-docker-dist"
)

func TestDockerDist(t *testing.T) {
	tmp, cleanup, err := dirs.TempDir("", "")
	defer cleanup()
	require.NoError(t, err)

	for i, currCase := range []struct {
		name          string
		spec          func(projectDir string, randomPad string) []params.ProductBuildSpecWithDeps
		preDistAction func(projectDir string, buildSpec []params.ProductBuildSpecWithDeps)
		validate      func(caseNum int, name string, pad string, cli *dockercli.Client)
	}{
		{
			name: "docker dist",
			spec: func(projectDir string, pad string) []params.ProductBuildSpecWithDeps {
				fooSpecWithDeps, err := params.NewProductBuildSpecWithDeps(params.NewProductBuildSpec(
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
								Info: &params.DockerDistInfo{
									Repository: fullRepoName("foo", pad),
									Tag:        "0.1.0",
									ContextDir: "foo/docker",
									DistDeps: map[string][]params.DistInfoType{
										"bar": {
											params.DockerDistType,
										},
									},
								},
							}},
					},
					params.Project{
						GroupID: "com.test.group",
					},
				), nil)
				require.NoError(t, err)
				barSpecWithDeps, err := params.NewProductBuildSpecWithDeps(params.NewProductBuildSpec(
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
						Dist: []params.Dist{
							{
								Info: &params.DockerDistInfo{
									Repository: fullRepoName("bar", pad),
									Tag:        "0.1.0",
									ContextDir: "bar/docker",
								},
							},
						},
					}, params.Project{
						GroupID: "com.test.group",
					},
				), nil)
				require.NoError(t, err)
				return []params.ProductBuildSpecWithDeps{fooSpecWithDeps, barSpecWithDeps}
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
	} {
		cli, err := dockercli.NewEnvClient()
		require.NoError(t, err)

		currTmpDir, err := ioutil.TempDir(tmp, "")
		require.NoError(t, err, "Case %d: %s", i, currCase.name)
		pad := randomPad(8)
		spec := currCase.spec(currTmpDir, pad)

		gittest.InitGitDir(t, currTmpDir)
		// initialize foo
		fooDir := path.Join(currTmpDir, "foo")
		err = os.Mkdir(fooDir, 0777)
		require.NoError(t, err, "Case %d: %s", i, currCase.name)
		err = ioutil.WriteFile(path.Join(fooDir, "main.go"), []byte(testMain), 0644)
		require.NoError(t, err, "Case %d: %s", i, currCase.name)
		fooDockerDir := path.Join(fooDir, "docker")
		err = os.Mkdir(fooDockerDir, 0777)
		require.NoError(t, err, "Case %d: %s", i, currCase.name)
		fooDockerFile := fmt.Sprintf("FROM %v:0.1.0\n", fullRepoName("bar", pad))
		err = ioutil.WriteFile(path.Join(fooDockerDir, "Dockerfile"), []byte(fooDockerFile), 0777)
		require.NoError(t, err, "Case %d: %s", i, currCase.name)

		// initialize bar
		barDir := path.Join(currTmpDir, "bar")
		err = os.Mkdir(barDir, 0777)
		require.NoError(t, err, "Case %d: %s", i, currCase.name)
		err = ioutil.WriteFile(path.Join(barDir, "main.go"), []byte(testMain), 0644)
		require.NoError(t, err, "Case %d: %s", i, currCase.name)
		barDockerDir := path.Join(barDir, "docker")
		err = os.Mkdir(barDockerDir, 0777)
		require.NoError(t, err, "Case %d: %s", i, currCase.name)
		err = ioutil.WriteFile(path.Join(barDockerDir, "Dockerfile"), []byte(dockerfile), 0777)
		require.NoError(t, err, "Case %d: %s", i, currCase.name)

		// commit
		gittest.CommitAllFiles(t, currTmpDir, "Commit")

		if currCase.preDistAction != nil {
			currCase.preDistAction(currTmpDir, spec)
		}

		// clean up docker images
		defer func(pad string) {
			images := []string{fmt.Sprintf("%v:0.1.0", fullRepoName("foo", pad)),
				fmt.Sprintf("%v:0.1.0", fullRepoName("bar", pad))}
			for _, image := range images {
				_, err := cli.ImageRemove(context.Background(), image, types.ImageRemoveOptions{})
				if err != nil {
					t.Fatalf("Error while pruning the image %v: %v\n", image, err.Error())
				}
			}

		}(pad)

		orderedSpecs, err := dist.OrderBuildSpecs(spec)
		require.NoError(t, err, "Case %d: %s", i, currCase.name)
		for _, currSpecWithDeps := range orderedSpecs {
			err = build.Run(build.RequiresBuild(currSpecWithDeps, nil).Specs(), nil, build.Context{}, ioutil.Discard)
			require.NoError(t, err, "Case %d: %s", i, currCase.name)
			err = dist.Run(currSpecWithDeps, ioutil.Discard)
			require.NoError(t, err, "Case %d: %s", i, currCase.name)
		}

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
