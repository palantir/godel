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

package publish_test

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"runtime"
	"testing"

	"github.com/nmiyake/pkg/dirs"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/palantir/godel/apps/distgo/cmd/build"
	"github.com/palantir/godel/apps/distgo/cmd/dist"
	"github.com/palantir/godel/apps/distgo/cmd/publish"
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
)

func TestPublishLocal(t *testing.T) {
	tmp, cleanup, err := dirs.TempDir("", "")
	defer cleanup()
	require.NoError(t, err)

	for i, currCase := range []struct {
		name        string
		buildSpec   func(projectDir string) params.ProductBuildSpecWithDeps
		skip        func() bool
		wantPaths   []string
		wantContent map[string]string
	}{
		{
			name: "local publish for SLS product",
			buildSpec: func(projectDir string) params.ProductBuildSpecWithDeps {
				specWithDeps, err := params.NewProductBuildSpecWithDeps(params.NewProductBuildSpec(projectDir, "publish-test-service", git.ProjectInfo{
					Version:  "0.0.1",
					Branch:   "0.0.1",
					Revision: "0",
				}, params.Product{
					Build: params.Build{
						MainPkg: "./.",
					},
					Dist: []params.Dist{{
						Info: &params.SLSDistInfo{},
					}},
					Publish: params.Publish{
						GroupID: "com.palantir.distgo-publish-test",
					},
				}, params.Project{}), nil)
				require.NoError(t, err)
				return specWithDeps
			},
			wantPaths: []string{
				"com/palantir/distgo-publish-test/publish-test-service/0.0.1/publish-test-service-0.0.1.pom",
				"com/palantir/distgo-publish-test/publish-test-service/0.0.1/publish-test-service-0.0.1.sls.tgz",
			},
			wantContent: map[string]string{
				"com/palantir/distgo-publish-test/publish-test-service/0.0.1/publish-test-service-0.0.1.pom": "<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n" +
					"<project xsi:schemaLocation=\"http://maven.apache.org/POM/4.0.0 http://maven.apache.org/xsd/maven-4.0.0.xsd\" xmlns=\"http://maven.apache.org/POM/4.0.0\"\nxmlns:xsi=\"http://www.w3.org/2001/XMLSchema-instance\">\n" +
					"<modelVersion>4.0.0</modelVersion>\n" +
					"<groupId>com.palantir.distgo-publish-test</groupId>\n" +
					"<artifactId>publish-test-service</artifactId>\n" +
					"<version>0.0.1</version>\n" +
					"<packaging>sls.tgz</packaging>\n" +
					"</project>\n",
			},
		},
		{
			name: "local publish for bin product",
			buildSpec: func(projectDir string) params.ProductBuildSpecWithDeps {
				specWithDeps, err := params.NewProductBuildSpecWithDeps(params.NewProductBuildSpec(projectDir, "publish-test-service", git.ProjectInfo{
					Version:  "0.0.1",
					Branch:   "0.0.1",
					Revision: "0",
				}, params.Product{
					Build: params.Build{
						MainPkg: "./.",
					},
					Dist: []params.Dist{{
						Info: &params.BinDistInfo{},
					}},
					Publish: params.Publish{
						GroupID: "com.palantir.distgo-publish-test",
					},
				}, params.Project{}), nil)
				require.NoError(t, err)
				return specWithDeps
			},
			wantPaths: []string{
				"com/palantir/distgo-publish-test/publish-test-service/0.0.1/publish-test-service-0.0.1.pom",
				"com/palantir/distgo-publish-test/publish-test-service/0.0.1/publish-test-service-0.0.1.tgz",
			},
			wantContent: map[string]string{
				"com/palantir/distgo-publish-test/publish-test-service/0.0.1/publish-test-service-0.0.1.pom": "<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n" +
					"<project xsi:schemaLocation=\"http://maven.apache.org/POM/4.0.0 http://maven.apache.org/xsd/maven-4.0.0.xsd\" xmlns=\"http://maven.apache.org/POM/4.0.0\"\n" +
					"xmlns:xsi=\"http://www.w3.org/2001/XMLSchema-instance\">\n" +
					"<modelVersion>4.0.0</modelVersion>\n" +
					"<groupId>com.palantir.distgo-publish-test</groupId>\n" +
					"<artifactId>publish-test-service</artifactId>\n" +
					"<version>0.0.1</version>\n" +
					"<packaging>tgz</packaging>\n" +
					"</project>\n",
			},
		},
		{
			name: "local publish for OSArch-bin product",
			buildSpec: func(projectDir string) params.ProductBuildSpecWithDeps {
				specWithDeps, err := params.NewProductBuildSpecWithDeps(params.NewProductBuildSpec(projectDir, "publish-test-service", git.ProjectInfo{
					Version:  "0.0.1",
					Branch:   "0.0.1",
					Revision: "0",
				}, params.Product{
					Build: params.Build{
						MainPkg: "./.",
						OSArchs: []osarch.OSArch{
							{
								OS:   "darwin",
								Arch: "amd64",
							},
						},
					},
					Dist: []params.Dist{{
						Info: &params.OSArchsBinDistInfo{
							OSArchs: []osarch.OSArch{
								{
									OS:   "darwin",
									Arch: "amd64",
								},
							},
						},
					}},
					Publish: params.Publish{
						GroupID: "com.palantir.distgo-publish-test",
					},
				}, params.Project{}), nil)
				require.NoError(t, err)
				return specWithDeps
			},
			wantPaths: []string{
				"com/palantir/distgo-publish-test/publish-test-service/0.0.1/publish-test-service-0.0.1.pom",
				"com/palantir/distgo-publish-test/publish-test-service/0.0.1/publish-test-service-0.0.1-darwin-amd64.tgz",
			},
			wantContent: map[string]string{
				"com/palantir/distgo-publish-test/publish-test-service/0.0.1/publish-test-service-0.0.1.pom": "<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n" +
					"<project xsi:schemaLocation=\"http://maven.apache.org/POM/4.0.0 http://maven.apache.org/xsd/maven-4.0.0.xsd\" xmlns=\"http://maven.apache.org/POM/4.0.0\"\n" +
					"xmlns:xsi=\"http://www.w3.org/2001/XMLSchema-instance\">\n" +
					"<modelVersion>4.0.0</modelVersion>\n" +
					"<groupId>com.palantir.distgo-publish-test</groupId>\n" +
					"<artifactId>publish-test-service</artifactId>\n" +
					"<version>0.0.1</version>\n" +
					"<packaging>tgz</packaging>\n" +
					"</project>\n",
			},
		},
		{
			name: "local publish for OSArch-bin product with multiple targets creates multiple artifacts but single POM",
			buildSpec: func(projectDir string) params.ProductBuildSpecWithDeps {
				specWithDeps, err := params.NewProductBuildSpecWithDeps(params.NewProductBuildSpec(projectDir, "publish-test-service", git.ProjectInfo{
					Version:  "0.0.1",
					Branch:   "0.0.1",
					Revision: "0",
				}, params.Product{
					Build: params.Build{
						MainPkg: "./.",
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
					Dist: []params.Dist{{
						Info: &params.OSArchsBinDistInfo{
							OSArchs: []osarch.OSArch{
								{
									OS:   "darwin",
									Arch: "amd64",
								},
							},
						},
					}, {
						Info: &params.OSArchsBinDistInfo{
							OSArchs: []osarch.OSArch{
								{
									OS:   "linux",
									Arch: "amd64",
								},
							},
						},
					}},
					Publish: params.Publish{
						GroupID: "com.palantir.distgo-publish-test",
					},
				}, params.Project{}), nil)
				require.NoError(t, err)
				return specWithDeps
			},
			wantPaths: []string{
				"com/palantir/distgo-publish-test/publish-test-service/0.0.1/publish-test-service-0.0.1.pom",
				"com/palantir/distgo-publish-test/publish-test-service/0.0.1/publish-test-service-0.0.1-darwin-amd64.tgz",
				"com/palantir/distgo-publish-test/publish-test-service/0.0.1/publish-test-service-0.0.1-linux-amd64.tgz",
			},
			wantContent: map[string]string{
				"com/palantir/distgo-publish-test/publish-test-service/0.0.1/publish-test-service-0.0.1.pom": "<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n" +
					"<project xsi:schemaLocation=\"http://maven.apache.org/POM/4.0.0 http://maven.apache.org/xsd/maven-4.0.0.xsd\" xmlns=\"http://maven.apache.org/POM/4.0.0\"\n" +
					"xmlns:xsi=\"http://www.w3.org/2001/XMLSchema-instance\">\n" +
					"<modelVersion>4.0.0</modelVersion>\n" +
					"<groupId>com.palantir.distgo-publish-test</groupId>\n" +
					"<artifactId>publish-test-service</artifactId>\n" +
					"<version>0.0.1</version>\n" +
					"<packaging>tgz</packaging>\n" +
					"</project>\n",
			},
		},
		{
			name: "local publish for OSArch-bin product with dist with multiple OS/Archs creates multiple artifacts but single POM",
			buildSpec: func(projectDir string) params.ProductBuildSpecWithDeps {
				specWithDeps, err := params.NewProductBuildSpecWithDeps(params.NewProductBuildSpec(projectDir, "publish-test-service", git.ProjectInfo{
					Version:  "0.0.1",
					Branch:   "0.0.1",
					Revision: "0",
				}, params.Product{
					Build: params.Build{
						MainPkg: "./.",
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
					Dist: []params.Dist{{
						Info: &params.OSArchsBinDistInfo{
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
					}},
					Publish: params.Publish{
						GroupID: "com.palantir.distgo-publish-test",
					},
				}, params.Project{}), nil)
				require.NoError(t, err)
				return specWithDeps
			},
			wantPaths: []string{
				"com/palantir/distgo-publish-test/publish-test-service/0.0.1/publish-test-service-0.0.1.pom",
				"com/palantir/distgo-publish-test/publish-test-service/0.0.1/publish-test-service-0.0.1-darwin-amd64.tgz",
				"com/palantir/distgo-publish-test/publish-test-service/0.0.1/publish-test-service-0.0.1-linux-amd64.tgz",
			},
			wantContent: map[string]string{
				"com/palantir/distgo-publish-test/publish-test-service/0.0.1/publish-test-service-0.0.1.pom": "<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n" +
					"<project xsi:schemaLocation=\"http://maven.apache.org/POM/4.0.0 http://maven.apache.org/xsd/maven-4.0.0.xsd\" xmlns=\"http://maven.apache.org/POM/4.0.0\"\n" +
					"xmlns:xsi=\"http://www.w3.org/2001/XMLSchema-instance\">\n" +
					"<modelVersion>4.0.0</modelVersion>\n" +
					"<groupId>com.palantir.distgo-publish-test</groupId>\n" +
					"<artifactId>publish-test-service</artifactId>\n" +
					"<version>0.0.1</version>\n" +
					"<packaging>tgz</packaging>\n" +
					"</project>\n",
			},
		},
		{
			name: "local publish for manual dist product",
			buildSpec: func(projectDir string) params.ProductBuildSpecWithDeps {
				specWithDeps, err := params.NewProductBuildSpecWithDeps(params.NewProductBuildSpec(projectDir, "publish-test-service", git.ProjectInfo{
					Version:  "0.0.1",
					Branch:   "0.0.1",
					Revision: "0",
				}, params.Product{
					Build: params.Build{
						Skip: true,
					},
					Dist: []params.Dist{{
						Script: `
echo "test-dist-contents" > "$DIST_DIR/$PRODUCT-$VERSION.tgz"
`,
						Info: &params.ManualDistInfo{
							Extension: "tgz",
						},
					}},
					Publish: params.Publish{
						GroupID: "com.palantir.distgo-publish-test",
					},
				}, params.Project{}), nil)
				require.NoError(t, err)
				return specWithDeps
			},
			wantPaths: []string{
				"com/palantir/distgo-publish-test/publish-test-service/0.0.1/publish-test-service-0.0.1.pom",
				"com/palantir/distgo-publish-test/publish-test-service/0.0.1/publish-test-service-0.0.1.tgz",
			},
			wantContent: map[string]string{
				"com/palantir/distgo-publish-test/publish-test-service/0.0.1/publish-test-service-0.0.1.pom": "<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n" +
					"<project xsi:schemaLocation=\"http://maven.apache.org/POM/4.0.0 http://maven.apache.org/xsd/maven-4.0.0.xsd\" xmlns=\"http://maven.apache.org/POM/4.0.0\"\n" +
					"xmlns:xsi=\"http://www.w3.org/2001/XMLSchema-instance\">\n" +
					"<modelVersion>4.0.0</modelVersion>\n" +
					"<groupId>com.palantir.distgo-publish-test</groupId>\n" +
					"<artifactId>publish-test-service</artifactId>\n" +
					"<version>0.0.1</version>\n" +
					"<packaging>tgz</packaging>\n" +
					"</project>\n",
			},
		},
		{
			name: "local publish for product with no distribution specified creates os-arch-bin",
			buildSpec: func(projectDir string) params.ProductBuildSpecWithDeps {
				specWithDeps, err := params.NewProductBuildSpecWithDeps(params.NewProductBuildSpec(projectDir, "publish-test-service", git.ProjectInfo{
					Version:  "0.0.1",
					Branch:   "0.0.1",
					Revision: "0",
				}, params.Product{
					Build: params.Build{
						MainPkg: "./.",
						OSArchs: []osarch.OSArch{
							{OS: "darwin", Arch: "amd64"},
						},
					},
					Publish: params.Publish{
						GroupID: "com.palantir.distgo-publish-test",
					},
				}, params.Project{}), nil)
				require.NoError(t, err)
				return specWithDeps
			},
			wantPaths: []string{
				"com/palantir/distgo-publish-test/publish-test-service/0.0.1/publish-test-service-0.0.1.pom",
				"com/palantir/distgo-publish-test/publish-test-service/0.0.1/publish-test-service-0.0.1-darwin-amd64.tgz",
			},
			wantContent: map[string]string{
				"com/palantir/distgo-publish-test/publish-test-service/0.0.1/publish-test-service-0.0.1.pom": "<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n" +
					"<project xsi:schemaLocation=\"http://maven.apache.org/POM/4.0.0 http://maven.apache.org/xsd/maven-4.0.0.xsd\" xmlns=\"http://maven.apache.org/POM/4.0.0\"\n" +
					"xmlns:xsi=\"http://www.w3.org/2001/XMLSchema-instance\">\n" +
					"<modelVersion>4.0.0</modelVersion>\n" +
					"<groupId>com.palantir.distgo-publish-test</groupId>\n" +
					"<artifactId>publish-test-service</artifactId>\n" +
					"<version>0.0.1</version>\n" +
					"<packaging>tgz</packaging>\n" +
					"</project>\n",
			},
		},
		{
			name: "local publish for product with bin and RPM distributions",
			buildSpec: func(projectDir string) params.ProductBuildSpecWithDeps {
				specWithDeps, err := params.NewProductBuildSpecWithDeps(params.NewProductBuildSpec(projectDir, "test", git.ProjectInfo{
					Version:  "0.0.1",
					Branch:   "0.0.1",
					Revision: "0",
				}, params.Product{
					Build: params.Build{
						MainPkg: "./.",
					},
					Dist: []params.Dist{{
						Info:      &params.BinDistInfo{},
						OutputDir: "dist/bin",
						Script:    "touch $DIST_DIR/dist-1.txt",
					}, {
						Info:      &params.RPMDistInfo{},
						OutputDir: "dist/rpm",
						Script:    "touch $DIST_DIR/dist-2.txt",
						Publish: params.Publish{
							GroupID: "com.palantir.pcloud-rpm",
						},
					}},
					Publish: params.Publish{
						GroupID: "com.palantir.pcloud-bin",
					},
				}, params.Project{}), nil)
				require.NoError(t, err)
				return specWithDeps
			},
			skip: func() bool {
				// rpm dist type is currently only supported for linux-amd64
				return !(runtime.GOOS == "linux" && runtime.GOARCH == "amd64")
			},
			wantPaths: []string{
				"com/palantir/pcloud-bin/test/0.0.1/test-0.0.1.pom",
				"com/palantir/pcloud-bin/test/0.0.1/test-0.0.1.tgz",
				"com/palantir/pcloud-rpm/test/0.0.1/test-0.0.1.pom",
				"com/palantir/pcloud-rpm/test/0.0.1/test-0.0.1-1.x86_64.rpm",
			},
			wantContent: map[string]string{
				"com/palantir/pcloud-bin/test/0.0.1/test-0.0.1.pom": "<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n<project xsi:schemaLocation=\"http://maven.apache.org/POM/4.0.0 http://maven.apache.org/xsd/maven-4.0.0.xsd\" xmlns=\"http://maven.apache.org/POM/4.0.0\"\nxmlns:xsi=\"http://www.w3.org/2001/XMLSchema-instance\">\n<modelVersion>4.0.0</modelVersion>\n<groupId>com.palantir.pcloud-bin</groupId>\n<artifactId>test</artifactId>\n<version>0.0.1</version>\n<packaging>tgz</packaging>\n</project>\n",
				"com/palantir/pcloud-rpm/test/0.0.1/test-0.0.1.pom": "<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n<project xsi:schemaLocation=\"http://maven.apache.org/POM/4.0.0 http://maven.apache.org/xsd/maven-4.0.0.xsd\" xmlns=\"http://maven.apache.org/POM/4.0.0\"\nxmlns:xsi=\"http://www.w3.org/2001/XMLSchema-instance\">\n<modelVersion>4.0.0</modelVersion>\n<groupId>com.palantir.pcloud-rpm</groupId>\n<artifactId>test</artifactId>\n<version>0.0.1</version>\n<packaging>rpm</packaging>\n</project>\n",
			},
		},
	} {
		if currCase.skip != nil && currCase.skip() {
			fmt.Printf("Skipping case %d\n", i)
			continue
		}

		currTmp, err := ioutil.TempDir(tmp, "")
		require.NoError(t, err)

		gittest.InitGitDir(t, currTmp)

		err = ioutil.WriteFile(path.Join(currTmp, "main.go"), []byte(testMain), 0644)
		require.NoError(t, err)

		currSpecWithDeps := currCase.buildSpec(currTmp)

		err = build.Run(build.RequiresBuild(currSpecWithDeps, nil).Specs(), nil, build.Context{}, ioutil.Discard)
		require.NoError(t, err, "Case %d: %s", i, currCase.name)

		err = dist.Run(currSpecWithDeps, ioutil.Discard)
		require.NoError(t, err, "Case %d: %s", i, currCase.name)

		repo := path.Join(currTmp, "repository")
		err = publish.Run(currSpecWithDeps, &publish.LocalPublishInfo{
			Path: repo,
		}, nil, ioutil.Discard)
		require.NoError(t, err, "Case %d: %s", i, currCase.name)

		for _, currPath := range currCase.wantPaths {
			info, err := os.Stat(path.Join(repo, currPath))
			require.NoError(t, err, "Case %d: %s", i, currCase.name)
			assert.False(t, info.IsDir(), "Case %d: %s", i, currCase.name)
		}

		for k, v := range currCase.wantContent {
			bytes, err := ioutil.ReadFile(path.Join(repo, k))
			require.NoError(t, err)
			assert.Equal(t, v, string(bytes), "Case %d", i)
		}
	}
}
