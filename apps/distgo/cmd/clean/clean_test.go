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

package clean_test

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"testing"

	"github.com/nmiyake/pkg/dirs"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/palantir/godel/apps/distgo/cmd/build"
	"github.com/palantir/godel/apps/distgo/cmd/clean"
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
)

func TestDist(t *testing.T) {
	tmp, cleanup, err := dirs.TempDir("", "")
	defer cleanup()
	require.NoError(t, err)

	for i, currCase := range []struct {
		name           string
		spec           func(projectDir string) params.ProductBuildSpecWithDeps
		preDistAction  func(projectDir string, buildSpec params.ProductBuildSpec)
		preCleanAction func(caseNum int, name string, projectDir string)
		postValidate   func(caseNum int, name string, projectDir string)
	}{
		{
			name: "cleans default distribution",
			spec: func(projectDir string) params.ProductBuildSpecWithDeps {
				specWithDeps, err := params.NewProductBuildSpecWithDeps(params.NewProductBuildSpec(
					projectDir,
					"foo",
					git.ProjectInfo{
						Version: "0.1.0",
					},
					params.Product{
						Build: params.Build{
							MainPkg: "./.",
						},
					},
					params.Project{
						GroupID: "com.test.group",
					},
				), nil)
				require.NoError(t, err)
				return specWithDeps
			},
			preDistAction: func(projectDir string, buildSpec params.ProductBuildSpec) {
				gittest.CreateGitTag(t, projectDir, "0.1.0")
			},
			preCleanAction: func(caseNum int, name string, projectDir string) {
				info, err := os.Stat(path.Join(projectDir, "build", "0.1.0", osarch.Current().String(), "foo"))
				require.NoError(t, err)
				assert.False(t, info.IsDir(), "Case %d: %s", caseNum, name)

				info, err = os.Stat(path.Join(projectDir, "dist", fmt.Sprintf("foo-0.1.0-%s.tgz", osarch.Current().String())))
				require.NoError(t, err)
				assert.False(t, info.IsDir(), "Case %d: %s", caseNum, name)
			},
			postValidate: func(caseNum int, name string, projectDir string) {
				_, err := os.Stat(path.Join(projectDir, "build", "0.1.0", osarch.Current().String(), "foo"))
				assert.True(t, os.IsNotExist(err), "Case %d: %s", caseNum, name)
				_, err = os.Stat(path.Join(projectDir, "build"))
				assert.True(t, os.IsNotExist(err), "Case %d: %s", caseNum, name)

				_, err = os.Stat(path.Join(projectDir, "dist", fmt.Sprintf("foo-0.1.0-%s.tgz", osarch.Current().String())))
				assert.True(t, os.IsNotExist(err), "Case %d: %s", caseNum, name)
				_, err = os.Stat(path.Join(projectDir, "dist"))
				assert.True(t, os.IsNotExist(err), "Case %d: %s", caseNum, name)
			},
		},
		{
			name: "cleans works if output does not exist",
			spec: func(projectDir string) params.ProductBuildSpecWithDeps {
				specWithDeps, err := params.NewProductBuildSpecWithDeps(params.NewProductBuildSpec(
					projectDir,
					"foo",
					git.ProjectInfo{
						Version: "0.1.0",
					},
					params.Product{
						Build: params.Build{
							MainPkg: "./.",
						},
					},
					params.Project{
						GroupID: "com.test.group",
					},
				), nil)
				require.NoError(t, err)
				return specWithDeps
			},
			preDistAction: func(projectDir string, buildSpec params.ProductBuildSpec) {
				gittest.CreateGitTag(t, projectDir, "0.1.0")
			},
			preCleanAction: func(caseNum int, name string, projectDir string) {
				err := os.RemoveAll(path.Join(projectDir, "build"))
				require.NoError(t, err)

				err = os.RemoveAll(path.Join(projectDir, "dist"))
				require.NoError(t, err)
			},
			postValidate: func(caseNum int, name string, projectDir string) {
				_, err := os.Stat(path.Join(projectDir, "build", "0.1.0", osarch.Current().String(), "foo"))
				assert.True(t, os.IsNotExist(err), "Case %d: %s", caseNum, name)
				_, err = os.Stat(path.Join(projectDir, "build"))
				assert.True(t, os.IsNotExist(err), "Case %d: %s", caseNum, name)

				_, err = os.Stat(path.Join(projectDir, "dist", fmt.Sprintf("foo-0.1.0-%s.tgz", osarch.Current().String())))
				assert.True(t, os.IsNotExist(err), "Case %d: %s", caseNum, name)
				_, err = os.Stat(path.Join(projectDir, "dist"))
				assert.True(t, os.IsNotExist(err), "Case %d: %s", caseNum, name)
			},
		},
	} {
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
		err = build.Run(build.RequiresBuild(currSpecWithDeps, nil).Specs(), nil, build.Context{}, ioutil.Discard)
		require.NoError(t, err, "Case %d: %s", i, currCase.name)

		err = dist.Run(currSpecWithDeps, ioutil.Discard)
		require.NoError(t, err, "Case %d: %s", i, currCase.name)
		if currCase.preCleanAction != nil {
			currCase.preCleanAction(i, currCase.name, currTmpDir)
		}

		err = clean.Run(currSpecWithDeps, false, ioutil.Discard)
		require.NoError(t, err, "Case %d: %s", i, currCase.name)
		if currCase.postValidate != nil {
			currCase.postValidate(i, currCase.name, currTmpDir)
		}
	}
}
