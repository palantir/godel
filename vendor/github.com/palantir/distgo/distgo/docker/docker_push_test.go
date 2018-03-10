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
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"regexp"
	"testing"

	"github.com/nmiyake/pkg/dirs"
	"github.com/palantir/pkg/gittest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/palantir/distgo/dister"
	"github.com/palantir/distgo/distgo"
	"github.com/palantir/distgo/distgo/docker"
	"github.com/palantir/distgo/dockerbuilder"
)

func TestDockerPublish(t *testing.T) {
	tmp, cleanup, err := dirs.TempDir("", "")
	defer cleanup()
	require.NoError(t, err)

	for i, tc := range []struct {
		name            string
		projectCfg      distgo.ProjectConfig
		dockerIDs       []distgo.ProductDockerID
		preDockerAction func(projectDir string, projectCfg distgo.ProjectConfig)
		wantErrorRegexp string
		wantStdout      string
	}{
		{
			"publish pushes Docker images",
			distgo.ProjectConfig{
				Products: map[distgo.ProductID]distgo.ProductConfig{
					"foo": {
						Build: &distgo.BuildConfig{
							MainPkg: stringPtr("./foo"),
						},
						Dist: &distgo.DistConfig{
							Disters: &distgo.DistersConfig{
								dister.OSArchBinDistTypeName: distgo.DisterConfig{
									Type: stringPtr(dister.OSArchBinDistTypeName),
								},
							},
						},
						Docker: &distgo.DockerConfig{
							DockerBuildersConfig: &distgo.DockerBuildersConfig{
								printDockerfileDockerBuilderTypeName: distgo.DockerBuilderConfig{
									Type:       stringPtr(printDockerfileDockerBuilderTypeName),
									ContextDir: stringPtr("docker-context-dir"),
									InputBuilds: &[]distgo.ProductBuildID{
										"foo",
									},
									InputDists: &[]distgo.ProductDistID{
										"foo",
									},
									TagTemplates: &[]string{
										"foo:latest",
									},
								},
							},
						},
					},
				},
			},
			nil,
			func(projectDir string, projectCfg distgo.ProjectConfig) {
				contextDir := path.Join(projectDir, "docker-context-dir")
				err := os.Mkdir(contextDir, 0755)
				require.NoError(t, err)
				dockerfile := path.Join(contextDir, "Dockerfile")
				err = ioutil.WriteFile(dockerfile, []byte(testDockerfile), 0644)
				require.NoError(t, err)
				gittest.CommitAllFiles(t, projectDir, "Commit files")
				gittest.CreateGitTag(t, projectDir, "0.1.0")
			},
			"",
			`[DRY RUN] Running Docker push for configuration print-dockerfile of product foo...
[DRY RUN] Run [docker push foo:latest]
`,
		},
		{
			"publish pushes all Docker images for a product",
			distgo.ProjectConfig{
				Products: map[distgo.ProductID]distgo.ProductConfig{
					"foo": {
						Build: &distgo.BuildConfig{
							MainPkg: stringPtr("./foo"),
						},
						Dist: &distgo.DistConfig{
							Disters: &distgo.DistersConfig{
								dister.OSArchBinDistTypeName: distgo.DisterConfig{
									Type: stringPtr(dister.OSArchBinDistTypeName),
								},
							},
						},
						Docker: &distgo.DockerConfig{
							DockerBuildersConfig: &distgo.DockerBuildersConfig{
								printDockerfileDockerBuilderTypeName: distgo.DockerBuilderConfig{
									Type:       stringPtr(printDockerfileDockerBuilderTypeName),
									ContextDir: stringPtr("docker-context-dir"),
									InputBuilds: &[]distgo.ProductBuildID{
										"foo",
									},
									InputDists: &[]distgo.ProductDistID{
										"foo",
									},
									TagTemplates: &[]string{
										"foo:latest",
										"foo:{{Version}}",
									},
								},
							},
						},
					},
				},
			},
			nil,
			func(projectDir string, projectCfg distgo.ProjectConfig) {
				contextDir := path.Join(projectDir, "docker-context-dir")
				err := os.Mkdir(contextDir, 0755)
				require.NoError(t, err)
				dockerfile := path.Join(contextDir, "Dockerfile")
				err = ioutil.WriteFile(dockerfile, []byte(testDockerfile), 0644)
				require.NoError(t, err)
				gittest.CommitAllFiles(t, projectDir, "Commit files")
				gittest.CreateGitTag(t, projectDir, "0.1.0")
			},
			"",
			`[DRY RUN] Running Docker push for configuration print-dockerfile of product foo...
[DRY RUN] Run [docker push foo:latest]
[DRY RUN] Run [docker push foo:0.1.0]
`,
		},
		{
			"publish pushes Docker images for a product but not for its dependencies",
			distgo.ProjectConfig{
				Products: map[distgo.ProductID]distgo.ProductConfig{
					"foo": {
						Build: &distgo.BuildConfig{
							MainPkg: stringPtr("./foo"),
						},
						Dist: &distgo.DistConfig{
							Disters: &distgo.DistersConfig{
								dister.OSArchBinDistTypeName: distgo.DisterConfig{
									Type: stringPtr(dister.OSArchBinDistTypeName),
								},
							},
						},
						Docker: &distgo.DockerConfig{
							DockerBuildersConfig: &distgo.DockerBuildersConfig{
								printDockerfileDockerBuilderTypeName: distgo.DockerBuilderConfig{
									Type:       stringPtr(printDockerfileDockerBuilderTypeName),
									ContextDir: stringPtr("docker-context-dir"),
									InputBuilds: &[]distgo.ProductBuildID{
										"foo",
									},
									InputDists: &[]distgo.ProductDistID{
										"foo",
									},
									TagTemplates: &[]string{
										"foo:latest",
										"foo:{{Version}}",
									},
								},
							},
						},
						Dependencies: &[]distgo.ProductID{
							"bar",
						},
					},
					"bar": {
						Docker: &distgo.DockerConfig{
							DockerBuildersConfig: &distgo.DockerBuildersConfig{
								printDockerfileDockerBuilderTypeName: distgo.DockerBuilderConfig{
									Type:       stringPtr(printDockerfileDockerBuilderTypeName),
									ContextDir: stringPtr("docker-context-dir"),
									TagTemplates: &[]string{
										"bar:latest",
									},
								},
							},
						},
					},
				},
			},
			[]distgo.ProductDockerID{
				"foo",
			},
			func(projectDir string, projectCfg distgo.ProjectConfig) {
				contextDir := path.Join(projectDir, "docker-context-dir")
				err := os.Mkdir(contextDir, 0755)
				require.NoError(t, err)
				dockerfile := path.Join(contextDir, "Dockerfile")
				err = ioutil.WriteFile(dockerfile, []byte(testDockerfile), 0644)
				require.NoError(t, err)
				gittest.CommitAllFiles(t, projectDir, "Commit files")
				gittest.CreateGitTag(t, projectDir, "0.1.0")
			},
			"",
			`[DRY RUN] Running Docker push for configuration print-dockerfile of product foo...
[DRY RUN] Run [docker push foo:latest]
[DRY RUN] Run [docker push foo:0.1.0]
`,
		},
	} {
		projectDir, err := ioutil.TempDir(tmp, "")
		require.NoError(t, err, "Case %d: %s", i, tc.name)

		gittest.InitGitDir(t, projectDir)
		err = os.MkdirAll(path.Join(projectDir, "foo"), 0755)
		require.NoError(t, err, "Case %d: %s", i, tc.name)
		err = ioutil.WriteFile(path.Join(projectDir, "foo", "main.go"), []byte(testMain), 0644)
		require.NoError(t, err, "Case %d: %s", i, tc.name)
		gittest.CommitAllFiles(t, projectDir, "Commit")

		if tc.preDockerAction != nil {
			tc.preDockerAction(projectDir, tc.projectCfg)
		}

		disterFactory, err := dister.NewDisterFactory()
		require.NoError(t, err, "Case %d: %s", i, tc.name)
		defaultDisterCfg, err := dister.DefaultConfig()
		require.NoError(t, err, "Case %d: %s", i, tc.name)
		dockerBuilderFactory, err := dockerbuilder.NewDockerBuilderFactory(dockerbuilder.NewCreator(printDockerfileDockerBuilderTypeName, newPrintDockerfileBuilder))
		require.NoError(t, err, "Case %d: %s", i, tc.name)

		projectParam, err := tc.projectCfg.ToParam(projectDir, disterFactory, defaultDisterCfg, dockerBuilderFactory)
		require.NoError(t, err, "Case %d: %s", i, tc.name)

		projectInfo, err := projectParam.ProjectInfo(projectDir)
		require.NoError(t, err, "Case %d: %s", i, tc.name)

		buffer := &bytes.Buffer{}
		err = docker.PushProducts(projectInfo, projectParam, tc.dockerIDs, true, buffer)
		if tc.wantErrorRegexp == "" {
			require.NoError(t, err, "Case %d: %s", i, tc.name)
		} else {
			require.Error(t, err, fmt.Sprintf("Case %d: %s", i, tc.name))
			assert.Regexp(t, regexp.MustCompile(tc.wantErrorRegexp), err.Error(), "Case %d: %s", i, tc.name)
		}

		if tc.wantStdout != "" {
			assert.Equal(t, tc.wantStdout, buffer.String(), "Case %d: %s", i, tc.name)
		}
	}
}

func exactMatchRegexp(in string) string {
	return "^" + regexp.QuoteMeta(in) + "$"
}
