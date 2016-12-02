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

package dist_test

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"regexp"
	"testing"
	"time"

	"github.com/nmiyake/pkg/dirs"
	"github.com/palantir/pkg/matcher"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/palantir/godel/apps/distgo/cmd/build"
	"github.com/palantir/godel/apps/distgo/cmd/dist"
	"github.com/palantir/godel/apps/distgo/config"
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
	expectManifest = `manifest-version: "1.0"
product-group: com.test.group
product-name: foo
product-version: 0.1.0
`
	expectManifestWithOptionalFields = `manifest-version: "1.0"
product-group: com.test.group
product-name: foo
product-version: 0.1.0
product-type: service.v1
extensions:
  bool-ext: true
  map-ext:
    hello: world
`
)

func TestDist(t *testing.T) {
	tmp, cleanup, err := dirs.TempDir("", "")
	defer cleanup()
	require.NoError(t, err)

	for i, currCase := range []struct {
		name            string
		skip            func() bool
		spec            func(projectDir string) config.ProductBuildSpecWithDeps
		preDistAction   func(projectDir string, buildSpec config.ProductBuildSpec)
		skipBuild       bool
		wantErrorRegexp string
		validate        func(caseNum int, name string, projectDir string)
	}{
		{
			name: "builds product and creates distribution directory and tgz",
			spec: func(projectDir string) config.ProductBuildSpecWithDeps {
				specWithDeps, err := config.NewProductBuildSpecWithDeps(config.NewProductBuildSpec(
					projectDir,
					"foo",
					git.ProjectInfo{
						Version: "0.1.0",
					},
					config.ProductConfig{
						Build: config.BuildConfig{
							MainPkg: "./.",
						},
					},
					config.ProjectConfig{
						GroupID: "com.test.group",
					},
				), nil)
				require.NoError(t, err)
				return specWithDeps
			},
			preDistAction: func(projectDir string, buildSpec config.ProductBuildSpec) {
				gittest.CreateGitTag(t, projectDir, "0.1.0")
			},
			validate: func(caseNum int, name string, projectDir string) {
				bytes, err := ioutil.ReadFile(path.Join(projectDir, "dist", "foo-0.1.0", "deployment", "manifest.yml"))
				require.NoError(t, err)
				assert.Equal(t, expectManifest, string(bytes), "Case %d: %s", caseNum, name)

				info, err := os.Stat(path.Join(projectDir, "dist", "foo-0.1.0.sls.tgz"))
				require.NoError(t, err)
				assert.False(t, info.IsDir(), "Case %d: %s", caseNum, name)

				info, err = os.Stat(path.Join(projectDir, "dist", "foo-0.1.0", "service", "bin", "init.sh"))
				require.NoError(t, err)
				assert.False(t, info.IsDir(), "Case %d: %s", caseNum, name)

				info, err = os.Stat(path.Join(projectDir, "build", "0.1.0", osarch.Current().String(), "foo"))
				require.NoError(t, err)
				assert.False(t, info.IsDir(), "Case %d: %s", caseNum, name)

				info, err = os.Stat(path.Join(projectDir, "dist", "foo-0.1.0", "service", "bin", osarch.Current().String(), "foo"))
				require.NoError(t, err)
				assert.False(t, info.IsDir(), "Case %d: %s", caseNum, name)
			},
		},
		{
			name: "SLS fails if GroupID is not specified",
			spec: func(projectDir string) config.ProductBuildSpecWithDeps {
				specWithDeps, err := config.NewProductBuildSpecWithDeps(config.NewProductBuildSpec(
					projectDir,
					"foo",
					git.ProjectInfo{
						Version: "0.1.0",
					},
					config.ProductConfig{
						Build: config.BuildConfig{
							MainPkg: "./.",
						},
					},
					config.ProjectConfig{},
				), nil)
				require.NoError(t, err)
				return specWithDeps
			},
			preDistAction: func(projectDir string, buildSpec config.ProductBuildSpec) {
				gittest.CreateGitTag(t, projectDir, "0.1.0")
			},
			wantErrorRegexp: "^failed to create manifest for SLS distribution: required properties were missing: group-id$",
		},
		{
			name: "SLS fails if generated artifact does not conform to SLS specification (missing manifest.yml)",
			spec: func(projectDir string) config.ProductBuildSpecWithDeps {
				specWithDeps, err := config.NewProductBuildSpecWithDeps(config.NewProductBuildSpec(
					projectDir,
					"foo",
					git.ProjectInfo{
						Version: "0.1.0",
					},
					config.ProductConfig{
						Build: config.BuildConfig{
							MainPkg: "./.",
						},
						Dist: []config.DistConfig{{
							Script: "rm $DIST_DIR/deployment/manifest.yml",
						}},
					},
					config.ProjectConfig{
						GroupID: "com.test.group",
					},
				), nil)
				require.NoError(t, err)
				return specWithDeps
			},
			preDistAction: func(projectDir string, buildSpec config.ProductBuildSpec) {
				gittest.CreateGitTag(t, projectDir, "0.1.0")
			},
			wantErrorRegexp: `(?s).+distribution directory failed SLS validation: foo-0.1.0/deployment/manifest.yml does not exist$`,
		},
		{
			name: "SLS fails if configuration.yml contains invalid YML",
			spec: func(projectDir string) config.ProductBuildSpecWithDeps {
				specWithDeps, err := config.NewProductBuildSpecWithDeps(config.NewProductBuildSpec(
					projectDir,
					"foo",
					git.ProjectInfo{
						Version: "0.1.0",
					},
					config.ProductConfig{
						Build: config.BuildConfig{
							MainPkg: "./.",
						},
						Dist: []config.DistConfig{{
							Script: `echo "{788=fads\n\tthis is invalid YML" > $DIST_DIR/deployment/configuration.yml`,
						}},
					},
					config.ProjectConfig{
						GroupID: "com.test.group",
					},
				), nil)
				require.NoError(t, err)
				return specWithDeps
			},
			preDistAction: func(projectDir string, buildSpec config.ProductBuildSpec) {
				gittest.CreateGitTag(t, projectDir, "0.1.0")
			},
			wantErrorRegexp: `(?s).+distribution directory failed SLS validation: invalid YML files: \[foo-0.1.0/deployment/configuration.yml\]
If these files are known to be correct, exclude them from validation using the SLS YML validation exclude matcher.$`,
		},
		{
			name: "SLS succeeds with invalid YML if it is excluded by matcher",
			spec: func(projectDir string) config.ProductBuildSpecWithDeps {
				specWithDeps, err := config.NewProductBuildSpecWithDeps(config.NewProductBuildSpec(
					projectDir,
					"foo",
					git.ProjectInfo{
						Version: "0.1.0",
					},
					config.ProductConfig{
						Build: config.BuildConfig{
							MainPkg: "./.",
						},
						Dist: []config.DistConfig{{
							Script: `echo "{788=fads\n\tthis is invalid YML" > $DIST_DIR/deployment/configuration.yml`,
							DistType: config.DistTypeConfig{
								Type: config.SLSDistType,
								Info: config.SLSDistInfo{
									YMLValidationExclude: matcher.NamesPathsCfg{
										Paths: []string{"deployment"},
									},
								},
							},
						}},
					},
					config.ProjectConfig{
						GroupID: "com.test.group",
					},
				), nil)
				require.NoError(t, err)
				return specWithDeps
			},
			preDistAction: func(projectDir string, buildSpec config.ProductBuildSpec) {
				gittest.CreateGitTag(t, projectDir, "0.1.0")
			},
		},
		{
			name: "copies executable from build location if it already exists",
			spec: func(projectDir string) config.ProductBuildSpecWithDeps {
				specWithDeps, err := config.NewProductBuildSpecWithDeps(config.NewProductBuildSpec(
					projectDir,
					"foo",
					git.ProjectInfo{
						Version: "0.1.0",
					},
					config.ProductConfig{
						Build: config.BuildConfig{
							MainPkg: "./.",
							OSArchs: []osarch.OSArch{
								{
									OS:   "fake",
									Arch: "fake",
								},
							},
						},
					},
					config.ProjectConfig{
						GroupID: "com.test.group",
					},
				), nil)
				require.NoError(t, err)
				return specWithDeps
			},
			preDistAction: func(projectDir string, buildSpec config.ProductBuildSpec) {
				gittest.CreateGitTag(t, projectDir, "0.1.0")

				// write fake executable
				artifactPath, ok := build.ArtifactPaths(buildSpec)[osarch.OSArch{OS: "fake", Arch: "fake"}]
				require.True(t, ok)

				err := os.MkdirAll(path.Dir(artifactPath), 0755)
				require.NoError(t, err)

				err = ioutil.WriteFile(artifactPath, []byte("test-content"), 0755)
				require.NoError(t, err)
			},
			validate: func(caseNum int, name string, projectDir string) {
				bytes, err := ioutil.ReadFile(path.Join(projectDir, "build", "0.1.0", "fake-fake", "foo"))
				require.NoError(t, err)
				assert.Equal(t, "test-content", string(bytes), "Case %d: %s", caseNum, name)
			},
		},
		{
			name: "re-builds executable if source files are newer than executable",
			spec: func(projectDir string) config.ProductBuildSpecWithDeps {
				specWithDeps, err := config.NewProductBuildSpecWithDeps(config.NewProductBuildSpec(
					projectDir,
					"foo",
					git.ProjectInfo{
						Version: "0.1.0",
					},
					config.ProductConfig{
						Build: config.BuildConfig{
							MainPkg: "./.",
						},
					},
					config.ProjectConfig{
						GroupID: "com.test.group",
					},
				), nil)
				require.NoError(t, err)
				return specWithDeps
			},
			preDistAction: func(projectDir string, buildSpec config.ProductBuildSpec) {
				gittest.CreateGitTag(t, projectDir, "0.1.0")

				// write fake executable
				artifactPath, ok := build.ArtifactPaths(buildSpec)[osarch.Current()]
				require.True(t, ok)

				err := os.MkdirAll(path.Dir(artifactPath), 0755)
				require.NoError(t, err)
				err = ioutil.WriteFile(artifactPath, []byte("test-content"), 0755)
				require.NoError(t, err)

				// write newer version of source file (sleep to ensure timestamp is later)
				time.Sleep(time.Second)
				err = ioutil.WriteFile(path.Join(projectDir, "main.go"), []byte(testMain+"\n"), 0644)
				require.NoError(t, err)
			},
			validate: func(caseNum int, name string, projectDir string) {
				// content should not be fake executable (build should be executed and overwrite content)
				bytes, err := ioutil.ReadFile(path.Join(projectDir, "build", "0.1.0", osarch.Current().String(), "foo"))
				require.NoError(t, err)
				assert.NotEqual(t, "test-content", string(bytes), "Case %d: %s", caseNum, name)
			},
		},
		{
			name: "copies layout from specified SLS input directory and ignores .gitkeep files",
			spec: func(projectDir string) config.ProductBuildSpecWithDeps {
				specWithDeps, err := config.NewProductBuildSpecWithDeps(config.NewProductBuildSpec(
					projectDir,
					"foo",
					git.ProjectInfo{
						Version: "0.1.0",
					},
					config.ProductConfig{
						Build: config.BuildConfig{
							MainPkg: "./.",
						},
						Dist: []config.DistConfig{{
							InputDir: "sls",
						}},
						DefaultPublish: config.PublishConfig{
							GroupID: "com.test.group",
						},
					},
					config.ProjectConfig{},
				), nil)
				require.NoError(t, err)
				return specWithDeps
			},
			preDistAction: func(projectDir string, buildSpec config.ProductBuildSpec) {
				gittest.CreateGitTag(t, projectDir, "0.1.0")

				// write manifest file that will be overwritten
				err := os.MkdirAll(path.Join(projectDir, "sls", "deployment"), 0755)
				require.NoError(t, err)
				err = ioutil.WriteFile(path.Join(projectDir, "sls", "deployment", "manifest.yml"), []byte("test-content"), 0644)
				require.NoError(t, err)

				// write .gitkeep file that should be ignored in top-level directory
				err = ioutil.WriteFile(path.Join(projectDir, "sls", ".gitkeep"), []byte(""), 0644)
				require.NoError(t, err)

				// write .gitkeep file that should be ignored in child directory
				err = os.MkdirAll(path.Join(projectDir, "sls", "empty"), 0755)
				require.NoError(t, err)
				err = ioutil.WriteFile(path.Join(projectDir, "sls", "empty", ".gitkeep"), []byte(""), 0644)
				require.NoError(t, err)

				// write test file that will be copied
				err = os.MkdirAll(path.Join(projectDir, "sls", "other"), 0755)
				require.NoError(t, err)
				err = ioutil.WriteFile(path.Join(projectDir, "sls", "other", "testfile"), []byte("test-content"), 0644)
				require.NoError(t, err)
			},
			validate: func(caseNum int, name string, projectDir string) {
				// manifest should be overwritten by dist
				bytes, err := ioutil.ReadFile(path.Join(projectDir, "dist", "foo-0.1.0", "deployment", "manifest.yml"))
				require.NoError(t, err)
				assert.Equal(t, expectManifest, string(bytes), "Case %d: %s", caseNum, name)

				// top-level .gitkeep should not exist
				fileInfo, err := os.Stat(path.Join(projectDir, "dist", "foo-0.1.0", ".gitkeep"))
				assert.True(t, os.IsNotExist(err), "Case %d: %s", caseNum, name)

				// empty directory should exist, but .gitkeep should not
				fileInfo, err = os.Stat(path.Join(projectDir, "dist", "foo-0.1.0", "empty"))
				assert.NoError(t, err, "Case %d: %s", caseNum, name)
				assert.True(t, fileInfo.IsDir(), "Case %d: %s", caseNum, name)
				fileInfo, err = os.Stat(path.Join(projectDir, "dist", "foo-0.1.0", "empty", ".gitkeep"))
				assert.True(t, os.IsNotExist(err), "Case %d: %s", caseNum, name)

				// test file should exist
				bytes, err = ioutil.ReadFile(path.Join(projectDir, "dist", "foo-0.1.0", "other", "testfile"))
				require.NoError(t, err)
				assert.Equal(t, "test-content", string(bytes), "Case %d: %s", caseNum, name)
			},
		},
		{
			name: "writes full SLS manifest with optional fields",
			spec: func(projectDir string) config.ProductBuildSpecWithDeps {
				specWithDeps, err := config.NewProductBuildSpecWithDeps(config.NewProductBuildSpec(
					projectDir,
					"foo",
					git.ProjectInfo{
						Version: "0.1.0",
					},
					config.ProductConfig{
						Build: config.BuildConfig{
							MainPkg: "./.",
						},
						Dist: []config.DistConfig{{
							DistType: config.DistTypeConfig{
								Type: config.SLSDistType,
								Info: config.SLSDistInfo{
									ProductType: "service.v1",
									ManifestExtensions: map[string]interface{}{
										"bool-ext": true,
										"map-ext": map[string]string{
											"hello": "world",
										},
									},
								},
							},
						}},
						DefaultPublish: config.PublishConfig{
							GroupID: "com.test.group",
						},
					},
					config.ProjectConfig{},
				), nil)
				require.NoError(t, err)
				return specWithDeps
			},
			preDistAction: func(projectDir string, buildSpec config.ProductBuildSpec) {
				gittest.CreateGitTag(t, projectDir, "0.1.0")
			},
			validate: func(caseNum int, name string, projectDir string) {
				// manifest should be overwritten by dist
				bytes, err := ioutil.ReadFile(path.Join(projectDir, "dist", "foo-0.1.0", "deployment", "manifest.yml"))
				require.NoError(t, err)
				assert.Equal(t, expectManifestWithOptionalFields, string(bytes), "Case %d: %s", caseNum, name)
			},
		},
		{
			name: "copies Windows executables",
			spec: func(projectDir string) config.ProductBuildSpecWithDeps {
				specWithDeps, err := config.NewProductBuildSpecWithDeps(config.NewProductBuildSpec(
					projectDir,
					"foo",
					git.ProjectInfo{
						Version: "0.1.0",
					},
					config.ProductConfig{
						Build: config.BuildConfig{
							MainPkg: "./.",
							OSArchs: []osarch.OSArch{
								{
									OS:   "windows",
									Arch: "amd64",
								},
							},
						},
					},
					config.ProjectConfig{
						GroupID: "com.test.group",
					},
				), nil)
				require.NoError(t, err)
				return specWithDeps
			},
			preDistAction: func(projectDir string, buildSpec config.ProductBuildSpec) {
				gittest.CreateGitTag(t, projectDir, "0.1.0")
			},
			validate: func(caseNum int, name string, projectDir string) {
				info, err := os.Stat(path.Join(projectDir, "dist", "foo-0.1.0", "service", "bin", "windows-amd64", "foo.exe"))
				require.NoError(t, err)
				assert.False(t, info.IsDir(), "Case %d: %s", caseNum, name)
			},
		},
		{
			name: "runs custom dist script",
			spec: func(projectDir string) config.ProductBuildSpecWithDeps {
				specWithDeps, err := config.NewProductBuildSpecWithDeps(config.NewProductBuildSpec(
					projectDir,
					"foo",
					git.ProjectInfo{
						Version: "0.1.0",
					},
					config.ProductConfig{
						Build: config.BuildConfig{
							MainPkg: "./.",
						},
						Dist: []config.DistConfig{{
							Script: "touch $DIST_DIR/test-file.txt",
						}},
					},
					config.ProjectConfig{
						GroupID: "com.test.group",
					},
				), nil)
				require.NoError(t, err)
				return specWithDeps
			},
			preDistAction: func(projectDir string, buildSpec config.ProductBuildSpec) {
				gittest.CreateGitTag(t, projectDir, "0.1.0")
			},
			validate: func(caseNum int, name string, projectDir string) {
				info, err := os.Stat(path.Join(projectDir, "dist", "foo-0.1.0", "test-file.txt"))
				require.NoError(t, err)
				assert.False(t, info.IsDir(), "Case %d: %s", caseNum, name)
			},
		},
		{
			name: "supports creating TGZ files that contain long paths",
			spec: func(projectDir string) config.ProductBuildSpecWithDeps {
				specWithDeps, err := config.NewProductBuildSpecWithDeps(config.NewProductBuildSpec(
					projectDir,
					"foo",
					git.ProjectInfo{
						Version: "0.1.0",
					},
					config.ProductConfig{
						Build: config.BuildConfig{
							MainPkg: "./.",
						},
						Dist: []config.DistConfig{{
							Script: `
							mkdir -p $DIST_DIR/0/1/2/3/4/5/6/7/8/9/10/11/12/13/14/15/16/17/18/19/20/21/22/23/24/25/26/27/28/29/30/31/32/33/
							touch $DIST_DIR/0/1/2/3/4/5/6/7/8/9/10/11/12/13/14/15/16/17/18/19/20/21/22/23/24/25/26/27/28/29/30/31/32/33/file.txt`,
						}},
					},
					config.ProjectConfig{
						GroupID: "com.test.group",
					},
				), nil)
				require.NoError(t, err)
				return specWithDeps
			},
			preDistAction: func(projectDir string, buildSpec config.ProductBuildSpec) {
				gittest.CreateGitTag(t, projectDir, "0.1.0")
			},
			validate: func(caseNum int, name string, projectDir string) {
				dst, err := ioutil.TempDir(projectDir, "expandedTGZDir")
				require.NoError(t, err)

				cmd := exec.Command("tar", "-C", dst, "-xzvf", path.Join(projectDir, "dist", "foo-0.1.0.sls.tgz"))
				output, err := cmd.CombinedOutput()
				require.NoError(t, err, "Command %v failed: %v", cmd.Args, string(output))

				// long file in tgz should be expanded properly
				_, err = os.Stat(path.Join(dst, "foo-0.1.0", "0/1/2/3/4/5/6/7/8/9/10/11/12/13/14/15/16/17/18/19/20/21/22/23/24/25/26/27/28/29/30/31/32/33/file.txt"))
				require.NoError(t, err, "Case %d: %s", caseNum, name)

				// stray file should not exist
				_, err = os.Stat(path.Join(dst, "file.txt"))
				require.Error(t, err, fmt.Sprintf("Case %d: %s", caseNum, name))
			},
		},
		{
			name: "custom dist script inherits process environment variables",
			spec: func(projectDir string) config.ProductBuildSpecWithDeps {
				err := os.Setenv("DIST_TEST_KEY", "distTestVal")
				require.NoError(t, err)
				err = os.Setenv("DIST_DIR", projectDir)
				require.NoError(t, err)

				specWithDeps, err := config.NewProductBuildSpecWithDeps(config.NewProductBuildSpec(
					projectDir,
					"foo",
					git.ProjectInfo{
						Version: "0.1.0",
					},
					config.ProductConfig{
						Build: config.BuildConfig{
							MainPkg: "./.",
						},
						Dist: []config.DistConfig{{
							Script: `touch $DIST_DIR/$DIST_TEST_KEY.txt
							touch $DIST_DIR/product:$PRODUCT
							touch $DIST_DIR/version:$VERSION
							touch $DIST_DIR/snapshot:$IS_SNAPSHOT`,
						}},
					},
					config.ProjectConfig{
						GroupID: "com.test.group",
					},
				), nil)
				require.NoError(t, err)
				return specWithDeps
			},
			preDistAction: func(projectDir string, buildSpec config.ProductBuildSpec) {
				gittest.CreateGitTag(t, projectDir, "0.1.0")
			},
			validate: func(caseNum int, name string, projectDir string) {
				info, err := os.Stat(path.Join(projectDir, "dist", "foo-0.1.0", "distTestVal.txt"))
				require.NoError(t, err)
				assert.False(t, info.IsDir(), "Case %d: %s", caseNum, name)

				info, err = os.Stat(path.Join(projectDir, "dist", "foo-0.1.0", "product:foo"))
				require.NoError(t, err)
				assert.False(t, info.IsDir(), "Case %d: %s", caseNum, name)

				info, err = os.Stat(path.Join(projectDir, "dist", "foo-0.1.0", "version:0.1.0"))
				require.NoError(t, err)
				assert.False(t, info.IsDir(), "Case %d: %s", caseNum, name)

				info, err = os.Stat(path.Join(projectDir, "dist", "foo-0.1.0", "snapshot:0"))
				require.NoError(t, err)
				assert.False(t, info.IsDir(), "Case %d: %s", caseNum, name)

				err = os.Unsetenv("DIST_TEST_KEY")
				require.NoError(t, err)
			},
		},
		{
			name: "custom dist script inherits dist script include",
			spec: func(projectDir string) config.ProductBuildSpecWithDeps {
				specWithDeps, err := config.NewProductBuildSpecWithDeps(config.NewProductBuildSpec(
					projectDir,
					"foo",
					git.ProjectInfo{
						Version: "0.1.0",
					},
					config.ProductConfig{
						Build: config.BuildConfig{
							MainPkg: "./.",
						},
						Dist: []config.DistConfig{{
							Script: `touch $DIST_DIR/$VERSION
							helper_func`,
						}},
					},
					config.ProjectConfig{
						DistScriptInclude: `touch $DIST_DIR/foo.txt
						helper_func() {
							touch $DIST_DIR/$IS_SNAPSHOT
						}`,
						GroupID: "com.test.group",
					},
				), nil)
				require.NoError(t, err)
				return specWithDeps
			},
			preDistAction: func(projectDir string, buildSpec config.ProductBuildSpec) {
				gittest.CreateGitTag(t, projectDir, "0.1.0")
			},
			validate: func(caseNum int, name string, projectDir string) {
				info, err := os.Stat(path.Join(projectDir, "dist", "foo-0.1.0", "foo.txt"))
				require.NoError(t, err)
				assert.False(t, info.IsDir(), "Case %d: %s", caseNum, name)

				info, err = os.Stat(path.Join(projectDir, "dist", "foo-0.1.0", "0.1.0"))
				require.NoError(t, err)
				assert.False(t, info.IsDir(), "Case %d: %s", caseNum, name)

				info, err = os.Stat(path.Join(projectDir, "dist", "foo-0.1.0", "0"))
				require.NoError(t, err)
				assert.False(t, info.IsDir(), "Case %d: %s", caseNum, name)
			},
		},
		{
			name: "custom dist script include does not run if script is not provided",
			spec: func(projectDir string) config.ProductBuildSpecWithDeps {
				specWithDeps, err := config.NewProductBuildSpecWithDeps(config.NewProductBuildSpec(
					projectDir,
					"foo",
					git.ProjectInfo{
						Version: "0.1.0",
					},
					config.ProductConfig{
						Build: config.BuildConfig{
							MainPkg: "./.",
						},
					},
					config.ProjectConfig{
						DistScriptInclude: "touch $DIST_DIR/foo.txt",
						GroupID:           "com.test.group",
					},
				), nil)
				require.NoError(t, err)
				return specWithDeps
			},
			preDistAction: func(projectDir string, buildSpec config.ProductBuildSpec) {
				gittest.CreateGitTag(t, projectDir, "0.1.0")
			},
			validate: func(caseNum int, name string, projectDir string) {
				_, err := os.Stat(path.Join(projectDir, "dist", "foo-0.1.0", "foo.txt"))
				assert.True(t, os.IsNotExist(err), "Case %d: %s", caseNum, name)
			},
		},
		{
			name: "copies dependent products",
			spec: func(projectDir string) config.ProductBuildSpecWithDeps {
				osArchsMap := make(map[osarch.OSArch]bool)
				osArchsMap[osarch.OSArch{
					OS:   "darwin",
					Arch: "amd64",
				}] = true
				osArchsMap[osarch.OSArch{
					OS:   "linux",
					Arch: "amd64",
				}] = true
				osArchsMap[osarch.Current()] = true

				var osArchsSlice []osarch.OSArch
				for osArch := range osArchsMap {
					osArchsSlice = append(osArchsSlice, osArch)
				}

				barSpec := config.NewProductBuildSpec(
					projectDir,
					"bar",
					git.ProjectInfo{
						Version: "0.1.0",
					},
					config.ProductConfig{
						Build: config.BuildConfig{
							MainPkg: "./.",
							OSArchs: osArchsSlice,
						},
					},
					config.ProjectConfig{
						GroupID: "com.test.group",
					},
				)

				specWithDeps, err := config.NewProductBuildSpecWithDeps(config.NewProductBuildSpec(
					projectDir,
					"foo",
					git.ProjectInfo{
						Version: "0.1.0",
					},
					config.ProductConfig{
						Build: config.BuildConfig{
							MainPkg: "./.",
						},
						Dist: []config.DistConfig{{
							InputProducts: []string{
								"bar",
							},
						}},
					},
					config.ProjectConfig{
						GroupID: "com.test.group",
					},
				), map[string]config.ProductBuildSpec{
					"bar": barSpec,
				})
				require.NoError(t, err)
				return specWithDeps
			},
			preDistAction: func(projectDir string, buildSpec config.ProductBuildSpec) {
				gittest.CreateGitTag(t, projectDir, "0.1.0")
			},
			validate: func(caseNum int, name string, projectDir string) {
				info, err := os.Stat(path.Join(projectDir, "dist", "foo-0.1.0", "service", "bin", osarch.Current().String(), "foo"))
				require.NoError(t, err)
				assert.False(t, info.IsDir(), "Case %d: %s", caseNum, name)

				info, err = os.Stat(path.Join(projectDir, "dist", "foo-0.1.0", "service", "bin", osarch.Current().String(), "bar"))
				require.NoError(t, err)
				assert.False(t, info.IsDir(), "Case %d: %s", caseNum, name)
			},
		},
		{
			name: "uses custom manifest when provided",
			spec: func(projectDir string) config.ProductBuildSpecWithDeps {
				manifestName := "test-manifest.yml"
				err := ioutil.WriteFile(path.Join(projectDir, manifestName), []byte(`---
manifestVersion: 1.0.0-alpha
productGroup: {{.Publish.GroupID}}
productName: {{.ProductName}}
productVersion: {{.ProductVersion}}
daemon: true
`), 0644)
				require.NoError(t, err)

				specWithDeps, err := config.NewProductBuildSpecWithDeps(config.NewProductBuildSpec(
					projectDir,
					"foo",
					git.ProjectInfo{
						Version: "0.1.0",
					},
					config.ProductConfig{
						Build: config.BuildConfig{
							MainPkg: "./.",
						},
						Dist: []config.DistConfig{{
							DistType: config.DistTypeConfig{
								Type: config.SLSDistType,
								Info: config.SLSDistInfo{
									ManifestTemplateFile: manifestName,
								},
							},
						}},
					},
					config.ProjectConfig{
						GroupID: "com.test.group",
					},
				), nil)
				require.NoError(t, err)
				return specWithDeps
			},
			preDistAction: func(projectDir string, buildSpec config.ProductBuildSpec) {
				gittest.CreateGitTag(t, projectDir, "0.1.0")
			},
			validate: func(caseNum int, name string, projectDir string) {
				bytes, err := ioutil.ReadFile(path.Join(projectDir, "dist", "foo-0.1.0", "deployment", "manifest.yml"))
				require.NoError(t, err)
				assert.Equal(t, `---
manifestVersion: 1.0.0-alpha
productGroup: com.test.group
productName: foo
productVersion: 0.1.0
daemon: true
`, string(bytes), "Case %d: %s", caseNum, name)
			},
		},
		{
			name: "uses custom init.sh when provided",
			spec: func(projectDir string) config.ProductBuildSpecWithDeps {
				initShName := "test-init.sh"
				err := ioutil.WriteFile(path.Join(projectDir, initShName), []byte(`init {{.ProductName}} {{.ProductVersion}} {{.Publish.GroupID}}`), 0644)
				require.NoError(t, err)

				specWithDeps, err := config.NewProductBuildSpecWithDeps(config.NewProductBuildSpec(
					projectDir,
					"foo",
					git.ProjectInfo{
						Version: "0.1.0",
					},
					config.ProductConfig{
						Build: config.BuildConfig{
							MainPkg: "./.",
						},
						Dist: []config.DistConfig{{
							DistType: config.DistTypeConfig{
								Type: config.SLSDistType,
								Info: config.SLSDistInfo{
									InitShTemplateFile: initShName,
								},
							},
						}},
					},
					config.ProjectConfig{
						GroupID: "com.test.group",
					},
				), nil)
				require.NoError(t, err)
				return specWithDeps
			},
			preDistAction: func(projectDir string, buildSpec config.ProductBuildSpec) {
				gittest.CreateGitTag(t, projectDir, "0.1.0")
			},
			validate: func(caseNum int, name string, projectDir string) {
				bytes, err := ioutil.ReadFile(path.Join(projectDir, "dist", "foo-0.1.0", "service", "bin", "init.sh"))
				require.NoError(t, err)
				assert.Equal(t, `init foo 0.1.0 com.test.group`, string(bytes), "Case %d: %s", caseNum, name)
			},
		},
		{
			name: "properly templatizes init.sh when ServiceArgs is empty",
			spec: func(projectDir string) config.ProductBuildSpecWithDeps {
				specWithDeps, err := config.NewProductBuildSpecWithDeps(config.NewProductBuildSpec(
					projectDir,
					"foo",
					git.ProjectInfo{
						Version: "0.1.0",
					},
					config.ProductConfig{
						Build: config.BuildConfig{
							MainPkg: "./.",
						},
						Dist: []config.DistConfig{{
							DistType: config.DistTypeConfig{
								Type: config.SLSDistType,
							},
						}},
					},
					config.ProjectConfig{
						GroupID: "com.test.group",
					},
				), nil)
				require.NoError(t, err)
				return specWithDeps
			},
			preDistAction: func(projectDir string, buildSpec config.ProductBuildSpec) {
				gittest.CreateGitTag(t, projectDir, "0.1.0")
			},
			validate: func(caseNum int, name string, projectDir string) {
				bytes, err := ioutil.ReadFile(path.Join(projectDir, "dist", "foo-0.1.0", "service", "bin", "init.sh"))
				require.NoError(t, err)
				assert.Regexp(t, `SERVICE_CMD="\$SERVICE_HOME/service/bin/\$OS_ARCH/\$SERVICE "\n`, string(bytes), "Case %d: %s", caseNum, name)
			},
		},
		{
			name: "properly templatizes init.sh with ServiceArgs",
			spec: func(projectDir string) config.ProductBuildSpecWithDeps {
				specWithDeps, err := config.NewProductBuildSpecWithDeps(config.NewProductBuildSpec(
					projectDir,
					"foo",
					git.ProjectInfo{
						Version: "0.1.0",
					},
					config.ProductConfig{
						Build: config.BuildConfig{
							MainPkg: "./.",
						},
						Dist: []config.DistConfig{{
							DistType: config.DistTypeConfig{
								Type: config.SLSDistType,
								Info: config.SLSDistInfo{
									ServiceArgs: "providedArgs arg2",
								},
							},
						}},
					},
					config.ProjectConfig{
						GroupID: "com.test.group",
					},
				), nil)
				require.NoError(t, err)
				return specWithDeps
			},
			preDistAction: func(projectDir string, buildSpec config.ProductBuildSpec) {
				gittest.CreateGitTag(t, projectDir, "0.1.0")
			},
			validate: func(caseNum int, name string, projectDir string) {
				bytes, err := ioutil.ReadFile(path.Join(projectDir, "dist", "foo-0.1.0", "service", "bin", "init.sh"))
				require.NoError(t, err)
				assert.Regexp(t, `SERVICE_CMD="\$SERVICE_HOME/service/bin/\$OS_ARCH/\$SERVICE providedArgs arg2"\n`, string(bytes), "Case %d: %s", caseNum, name)
			},
		},
		{
			name: "creates outputs using bin mode",
			spec: func(projectDir string) config.ProductBuildSpecWithDeps {
				specWithDeps, err := config.NewProductBuildSpecWithDeps(config.NewProductBuildSpec(
					projectDir,
					"foo",
					git.ProjectInfo{
						Version: "0.1.0",
					},
					config.ProductConfig{
						Build: config.BuildConfig{
							MainPkg: "./.",
						},
						Dist: []config.DistConfig{{
							DistType: config.DistTypeConfig{
								Type: config.BinDistType,
							},
						}},
					},
					config.ProjectConfig{},
				), nil)
				require.NoError(t, err)
				return specWithDeps
			},
			preDistAction: func(projectDir string, buildSpec config.ProductBuildSpec) {
				gittest.CreateGitTag(t, projectDir, "0.1.0")
			},
			validate: func(caseNum int, name string, projectDir string) {
				// bin directory exists in top-level directory
				fileInfo, err := os.Stat(path.Join(projectDir, "dist", "foo-0.1.0", "bin"))
				require.NoError(t, err)
				assert.True(t, fileInfo.IsDir(), "Case %d: %s", caseNum, name)

				// executable should exist in os-arch directory
				info, err := os.Stat(path.Join(projectDir, "dist", "foo-0.1.0", "bin", osarch.Current().String(), "foo"))
				require.NoError(t, err)
				assert.False(t, info.IsDir(), "Case %d: %s", caseNum, name)
			},
		},
		{
			name: "builds rpm",
			skip: func() bool {
				// Run this case only if both fpm and rpmbuild are available
				_, fpmErr := exec.LookPath("fpm")
				_, rpmbuildErr := exec.LookPath("rpmbuild")
				return fpmErr != nil || rpmbuildErr != nil
			},
			spec: func(projectDir string) config.ProductBuildSpecWithDeps {
				specWithDeps, err := config.NewProductBuildSpecWithDeps(config.NewProductBuildSpec(
					projectDir,
					"foo",
					git.ProjectInfo{
						Version: "0.1.0",
					},
					config.ProductConfig{
						Build: config.BuildConfig{
							MainPkg: "./.",
							OSArchs: []osarch.OSArch{
								{
									OS:   "linux",
									Arch: "amd64",
								},
							},
						},
						Dist: []config.DistConfig{{
							InputDir: "root",
							DistType: config.DistTypeConfig{
								Type: config.RPMDistType,
								Info: config.RPMDistInfo{
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
						}},
					},
					config.ProjectConfig{
						GroupID: "com.test.group",
					},
				), nil)
				require.NoError(t, err)
				return specWithDeps
			},
			preDistAction: func(projectDir string, buildSpec config.ProductBuildSpec) {
				gittest.CreateGitTag(t, projectDir, "0.1.0")

				// write fake systemd config file
				err := os.MkdirAll(path.Join(projectDir, "root", "usr", "lib", "systemd", "system"), 0755)
				require.NoError(t, err)

				err = ioutil.WriteFile(path.Join(projectDir, "root", "usr", "lib", "systemd", "system", "orchestrator.service"), []byte("configured"), 0600)
				require.NoError(t, err)
			},
			validate: func(caseNum int, name string, projectDir string) {
				info, err := os.Stat(path.Join(projectDir, "dist", "foo-0.1.0-1.x86_64.rpm"))
				require.NoError(t, err)
				assert.False(t, info.IsDir(), "Case %d: %s", caseNum, name)
			},
		},
		{
			name: "more than one dist",
			spec: func(projectDir string) config.ProductBuildSpecWithDeps {
				specWithDeps, err := config.NewProductBuildSpecWithDeps(config.NewProductBuildSpec(
					projectDir,
					"foo",
					git.ProjectInfo{
						Version: "0.1.0",
					},
					config.ProductConfig{
						Build: config.BuildConfig{
							MainPkg: "./.",
							OSArchs: []osarch.OSArch{
								{
									OS:   "linux",
									Arch: "amd64",
								},
							},
						},
						Dist: []config.DistConfig{{
							DistType: config.DistTypeConfig{
								Type: config.BinDistType,
							},
							OutputDir: "dist/bin",
							Script:    "touch $DIST_DIR/dist-1.txt",
						}, {
							DistType: config.DistTypeConfig{
								Type: config.RPMDistType,
							},
							OutputDir: "dist/rpm",
							Script:    "touch $DIST_DIR/dist-2.txt",
						}},
					},
					config.ProjectConfig{
						GroupID: "com.test.group",
					},
				), nil)
				require.NoError(t, err)
				return specWithDeps
			},
			preDistAction: func(projectDir string, buildSpec config.ProductBuildSpec) {
				gittest.CreateGitTag(t, projectDir, "0.1.0")
			},
			validate: func(caseNum int, name string, projectDir string) {
				info, err := os.Stat(path.Join(projectDir, "dist", "bin", "foo-0.1.0", "dist-1.txt"))
				require.NoError(t, err)
				assert.False(t, info.IsDir(), "Case %d: %s", caseNum, name)

				info, err = os.Stat(path.Join(projectDir, "dist", "rpm", "foo-0.1.0", "dist-2.txt"))
				require.NoError(t, err)
				assert.False(t, info.IsDir(), "Case %d: %s", caseNum, name)
			},
		},
	} {
		if currCase.skip != nil && currCase.skip() {
			fmt.Fprintln(os.Stderr, "SKIPPING CASE", i)
			continue
		}

		currTmpDir, err := ioutil.TempDir(tmp, "")
		require.NoError(t, err, "Case %d: %s", i, currCase.name)

		gittest.InitGitDir(t, currTmpDir)
		err = ioutil.WriteFile(path.Join(currTmpDir, "main.go"), []byte(testMain), 0644)
		require.NoError(t, err, "Case %d: %s", i, currCase.name)
		gittest.CommitAllFiles(t, currTmpDir, "Commit")

		if currCase.preDistAction != nil {
			currCase.preDistAction(currTmpDir, currCase.spec(currTmpDir).Spec)
		}

		currSpecWithDeps := currCase.spec(currTmpDir)
		if !currCase.skipBuild {
			err = build.Run(build.RequiresBuild(currSpecWithDeps, nil).Specs(), nil, build.Context{}, ioutil.Discard)
			require.NoError(t, err, "Case %d: %s", i, currCase.name)
		}

		err = dist.Run(currSpecWithDeps, ioutil.Discard)
		if currCase.wantErrorRegexp == "" {
			require.NoError(t, err, "Case %d: %s", i, currCase.name)
		} else {
			require.Error(t, err, fmt.Sprintf("Case %d: %s", i, currCase.name))
			assert.Regexp(t, regexp.MustCompile(currCase.wantErrorRegexp), err.Error(), "Case %d: %s", i, currCase.name)
		}

		if currCase.validate != nil {
			currCase.validate(i, currCase.name, currTmpDir)
		}
	}
}
