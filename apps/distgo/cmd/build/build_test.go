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

package build_test

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"reflect"
	"regexp"
	"strings"
	"testing"

	"github.com/nmiyake/pkg/dirs"
	"github.com/palantir/pkg/pkgpath"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/palantir/godel/apps/distgo/cmd"
	"github.com/palantir/godel/apps/distgo/cmd/build"
	"github.com/palantir/godel/apps/distgo/params"
	"github.com/palantir/godel/apps/distgo/pkg/binspec"
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
	testCMain = `package main

import "C"
import "fmt"

func main() {
	fmt.Println("C")
}`
	testVersionValue = "1.0.1"
	testBuildScript  = `package main

import (
	"fmt"
	"./dependency" // written by the build script
)

func main() {
	fmt.Println(dependency.V)
}
`
	longCompileMain = `package main

import (
	"encoding/json"
	"net/http"
)

func main() {
	http.Get("")
	json.Marshal("")
}
`
)

func TestBuildAll(t *testing.T) {
	tmp, cleanup, err := dirs.TempDir("", "")
	defer cleanup()
	require.NoError(t, err)

	for i, currCase := range []struct {
		productName     string
		mainFileContent string
		mainFilePath    string
		params          params.Product
		wantError       bool
		runExecutable   bool
		wantOutput      string
	}{
		{
			productName:     "randomProduct",
			mainFileContent: testMain,
			mainFilePath:    "main.go",
			params: params.Product{
				Build: params.Build{
					MainPkg:    "./.",
					VersionVar: "main.testVersionVar",
					OSArchs: []osarch.OSArch{
						osarch.Current(),
					},
				},
			},
			runExecutable: true,
			wantOutput:    testVersionValue + ".dirty",
		},
		// building project that requires CGo succeeds if "CGO_ENABLED" environment variable is set to 1
		{
			productName:     "CProduct",
			mainFileContent: testCMain,
			mainFilePath:    "main.go",
			params: params.Product{
				Build: params.Build{
					MainPkg: "./.",
					Environment: map[string]string{
						"CGO_ENABLED": "1",
					},
					OSArchs: []osarch.OSArch{
						osarch.Current(),
					},
				},
			},
			runExecutable: true,
			wantOutput:    "C",
		},
		// building project that requires CGo fails if "CGO_ENABLED" environment variable is set to 0
		{
			productName:     "CProduct",
			mainFileContent: testCMain,
			mainFilePath:    "main.go",
			params: params.Product{
				Build: params.Build{
					MainPkg: "./.",
					Environment: map[string]string{
						"CGO_ENABLED": "0",
					},
					OSArchs: []osarch.OSArch{
						osarch.Current(),
					},
				},
			},
			wantError: true,
		},
		{
			productName:     "preBuildScript",
			mainFileContent: testBuildScript,
			mainFilePath:    "main.go",
			params: params.Product{
				Build: params.Build{
					Script: "" +
						"mkdir dependency\n" +
						"echo 'package dependency\n\nvar V = `success`\n' > dependency/lib.go\n",
					MainPkg: "./.",
					OSArchs: []osarch.OSArch{
						osarch.Current(),
					},
				},
			},
			wantOutput: "success",
		},
		{
			productName:     "customBuildScriptProduct",
			mainFileContent: testMain,
			mainFilePath:    "main.go",
			params: params.Product{
				Build: params.Build{
					MainPkg: "./.",
					BuildArgsScript: `set -eu pipefail
VALUE="foo bar"
echo "-ldflags"
echo "-X \"main.testVersionVar=$VALUE\""`,
					OSArchs: []osarch.OSArch{
						osarch.Current(),
					},
				},
			},
			runExecutable: true,
			wantOutput:    "foo bar",
		},
		{
			productName:     "foo",
			mainFileContent: testMain,
			mainFilePath:    "foo/main.go",
			params: params.Product{
				Build: params.Build{
					MainPkg: "./foo",
					OSArchs: []osarch.OSArch{
						{
							OS:   "darwin",
							Arch: "amd64",
						},
						{
							OS:   "linux",
							Arch: "amd64",
						},
						{
							OS:   "windows",
							Arch: "amd64",
						},
					},
				},
			},
			wantOutput: "defaultVersion",
		},
	} {
		currTmpDir, err := ioutil.TempDir(tmp, "")
		require.NoError(t, err)

		gittest.InitGitDir(t, currTmpDir)
		gittest.CreateGitTag(t, currTmpDir, testVersionValue)

		mainFilePath := path.Join(currTmpDir, currCase.mainFilePath)

		err = os.MkdirAll(path.Dir(mainFilePath), 0755)
		require.NoError(t, err)

		err = ioutil.WriteFile(mainFilePath, []byte(currCase.mainFileContent), 0644)
		require.NoError(t, err)

		binDir := path.Join(currTmpDir, "bin")
		err = os.Mkdir(binDir, 0755)
		require.NoError(t, err)

		pkgPath, err := pkgpath.NewAbsPkgPath(path.Dir(mainFilePath)).Rel(currTmpDir)
		require.NoError(t, err)

		spec := binspec.New(currCase.params.Build.OSArchs, path.Base(pkgPath))
		err = spec.CreateDirectoryStructure(binDir, nil, false)
		require.NoError(t, err)

		gitProductInfo, err := git.NewProjectInfo(currTmpDir)
		require.NoError(t, err)

		buildSpec := params.NewProductBuildSpec(
			currTmpDir,
			currCase.productName,
			gitProductInfo,
			currCase.params,
			params.Project{
				BuildOutputDir: "bin",
			},
		)

		foundExecForCurrOsArch := false

		err = build.Run([]params.ProductBuildSpec{buildSpec}, nil, build.Context{
			Parallel: false,
		}, ioutil.Discard)

		if currCase.wantError {
			assert.Error(t, err, fmt.Sprintf("Case %d", i))
		} else {
			assert.NoError(t, err, "Case %d", i)

			artifactPaths := build.ArtifactPaths(buildSpec)
			for _, currOSArch := range currCase.params.Build.OSArchs {
				pathToCurrExecutable, ok := artifactPaths[currOSArch]
				require.True(t, ok, "Case %d: could not find path for %s for %s", buildSpec.ProductName, currOSArch.String())
				fileInfo, err := os.Stat(pathToCurrExecutable)
				require.NoError(t, err, "Case %d", i)
				assert.False(t, fileInfo.IsDir())

				if reflect.DeepEqual(currOSArch, osarch.Current()) {
					foundExecForCurrOsArch = true
					output, err := exec.Command(pathToCurrExecutable).Output()
					require.NoError(t, err)
					assert.Equal(t, currCase.wantOutput, strings.TrimSpace(string(output)), "Case %d", i)
				}
			}

			if currCase.runExecutable {
				assert.True(t, foundExecForCurrOsArch, "Case %d: executable for current os/arch (%v) not found in %v", osarch.Current(), currCase.params.Build.OSArchs)
			}
		}
	}
}

func TestBuildOnlyDistinctSpecs(t *testing.T) {
	tmp, cleanup, err := dirs.TempDir("", "")
	defer cleanup()
	require.NoError(t, err)

	mainFilePath := path.Join(tmp, "foo/main.go")
	err = os.MkdirAll(path.Dir(mainFilePath), 0755)
	require.NoError(t, err)
	err = ioutil.WriteFile(mainFilePath, []byte(testMain), 0644)
	require.NoError(t, err)

	buildSpec := params.NewProductBuildSpec(
		tmp,
		"foo",
		git.ProjectInfo{},
		params.Product{
			Build: params.Build{
				MainPkg: "./foo",
			},
		},
		params.Project{
			BuildOutputDir: "bin",
		},
	)

	buf := &bytes.Buffer{}
	err = build.Run([]params.ProductBuildSpec{buildSpec, buildSpec}, nil, build.Context{
		Parallel: false,
	}, buf)
	require.NoError(t, err)

	assert.Equal(t, 1, strings.Count(buf.String(), "Finished building foo"))
}

func TestBuildOnlySpecifiedOSArchs(t *testing.T) {
	tmp, cleanup, err := dirs.TempDir("", "")
	defer cleanup()
	require.NoError(t, err)

	mainFilePath := path.Join(tmp, "foo/main.go")
	err = os.MkdirAll(path.Dir(mainFilePath), 0755)
	require.NoError(t, err)
	err = ioutil.WriteFile(mainFilePath, []byte(testMain), 0644)
	require.NoError(t, err)

	for i, currCase := range []struct {
		specOSArchs []osarch.OSArch
		osArchs     []osarch.OSArch
		want        []string
		notWant     []string
	}{
		// empty value for osArchs filter builds all
		{
			specOSArchs: []osarch.OSArch{{OS: "darwin", Arch: "amd64"}, {OS: "linux", Arch: "386"}},
			osArchs:     nil,
			want: []string{
				"Finished building foo for darwin-amd64",
				"Finished building foo for linux-386",
			},
		},
		// if non-empty filter is provided, only values matching filter are built
		{
			specOSArchs: []osarch.OSArch{{OS: "darwin", Arch: "amd64"}, {OS: "linux", Arch: "386"}},
			osArchs:     []osarch.OSArch{{OS: "linux", Arch: "386"}},
			want: []string{
				"Finished building foo for linux-386",
			},
			notWant: []string{
				"Finished building foo for darwin-amd64",
			},
		},
		// if no OS/arch values match filter, nothing is built
		{
			specOSArchs: []osarch.OSArch{{OS: "darwin", Arch: "amd64"}, {OS: "linux", Arch: "386"}},
			osArchs:     []osarch.OSArch{{OS: "windows", Arch: "386"}},
			want: []string{
				"$^",
			},
		},
	} {
		buildSpec := params.NewProductBuildSpec(
			tmp,
			"foo",
			git.ProjectInfo{},
			params.Product{
				Build: params.Build{
					MainPkg: "./foo",
					OSArchs: currCase.specOSArchs,
				},
			},
			params.Project{
				BuildOutputDir: "bin",
			},
		)

		buf := &bytes.Buffer{}
		err = build.Run([]params.ProductBuildSpec{buildSpec}, cmd.OSArchFilter(currCase.osArchs), build.Context{
			Parallel: false,
		}, buf)
		require.NoError(t, err)

		for _, want := range currCase.want {
			assert.Regexp(t, regexp.MustCompile(want), buf.String(), "Case %d", i)
		}

		for _, notWant := range currCase.notWant {
			assert.NotRegexp(t, regexp.MustCompile(notWant), buf.String(), "Case %d", i)
		}
	}
}

func TestBuildErrorMessage(t *testing.T) {
	tmp, cleanup, err := dirs.TempDir(".", "")
	defer cleanup()
	require.NoError(t, err)

	mainFilePath := path.Join(tmp, "foo/main.go")
	err = os.MkdirAll(path.Dir(mainFilePath), 0755)
	require.NoError(t, err)
	err = ioutil.WriteFile(mainFilePath, []byte(`package main; asdfa`), 0644)
	require.NoError(t, err)

	buildSpec := params.NewProductBuildSpec(
		tmp,
		"foo",
		git.ProjectInfo{},
		params.Product{
			Build: params.Build{
				MainPkg: "./foo",
			},
		},
		params.Project{
			BuildOutputDir: "bin",
		},
	)

	want := `(?s)^go install failed: build command \[.+go install ./foo\] run with additional environment variables \[GOOS=.+ GOARCH=.+\] failed with output:.+foo/main.go:1:15: syntax error: non-declaration statement outside function body$`

	buf := &bytes.Buffer{}
	err = build.Run([]params.ProductBuildSpec{buildSpec, buildSpec}, nil, build.Context{
		Install:  true,
		Parallel: false,
	}, buf)
	assert.Regexp(t, want, err.Error())
}

func TestBuildInstallErrorMessage(t *testing.T) {
	tmp, cleanup, err := dirs.TempDir(".", "")
	defer cleanup()
	require.NoError(t, err)

	goRoot, err := dirs.GoRoot()
	require.NoError(t, err)
	_, err = os.Stat(goRoot)
	require.NoError(t, err)

	pkgDir := path.Join(goRoot, "pkg")
	_, err = os.Stat(pkgDir)
	require.NoError(t, err)

	osArchPkgDir := path.Join(pkgDir, "dragonfly_amd64")
	_, err = os.Stat(osArchPkgDir)
	if os.IsNotExist(err) {
		// if directory does not exist, attempt to create it (and clean up afterwards)
		if err := os.Mkdir(osArchPkgDir, 0444); err == nil {
			defer func() {
				if err := os.RemoveAll(osArchPkgDir); err != nil {
					fmt.Printf("Failed to remove directory %v: %v\n", osArchPkgDir, err)
				}
			}()
		}
		// if creation failed, assume that write permissions do not exist, which is sufficient for the test
	}

	mainFilePath := path.Join(tmp, "foo/main.go")
	err = os.MkdirAll(path.Dir(mainFilePath), 0755)
	require.NoError(t, err)
	err = ioutil.WriteFile(mainFilePath, []byte(`package main`), 0644)
	require.NoError(t, err)

	buildSpec := params.NewProductBuildSpec(
		tmp,
		"foo",
		git.ProjectInfo{},
		params.Product{
			Build: params.Build{
				MainPkg: "./foo",
				OSArchs: []osarch.OSArch{
					{
						OS:   "dragonfly",
						Arch: "amd64",
					},
				},
			},
		},
		params.Project{
			BuildOutputDir: "bin",
		},
	)

	goBinary := "go"
	if output, err := exec.Command("command", "-v", "go").CombinedOutput(); err == nil {
		goBinary = strings.TrimSpace(string(output))
	}

	want := `(?s)go install failed: failed to install a Go standard library package due to insufficient permissions to create directory.\n` +
		`This typically means that the standard library for the OS/architecture combination have not been installed locally and the current user does not have write permissions to GOROOT/pkg.\n` +
		fmt.Sprintf("Run \"sudo env GOOS=dragonfly GOARCH=amd64 %s install std\" to install the standard packages for this combination as root and then try again.\n", goBinary) +
		`Full error: build command \[.+/go install ./foo\] run with additional environment variables \[GOOS=dragonfly GOARCH=amd64\] failed with output:\n` +
		`go install .+: mkdir .+: permission denied$`

	buf := &bytes.Buffer{}
	err = build.Run([]params.ProductBuildSpec{buildSpec, buildSpec}, nil, build.Context{
		Install:  true,
		Parallel: false,
	}, buf)
	assert.Regexp(t, want, err.Error())
}

func TestBuildAllParallel(t *testing.T) {
	tmp, cleanup, err := dirs.TempDir("", "")
	defer cleanup()
	require.NoError(t, err)

	for i, currCase := range []struct {
		mainFiles map[string]string
		specs     []params.ProductBuildSpec
	}{
		{
			mainFiles: map[string]string{
				"foo/main.go": longCompileMain,
				"bar/main.go": longCompileMain,
			},
			specs: []params.ProductBuildSpec{
				{
					ProductName: "foo",
					Product: params.Product{
						Build: params.Build{
							MainPkg: "./foo",
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
							OutputDir: "build",
						},
					},
					VersionInfo: git.ProjectInfo{
						Version: "0.1.0",
					},
				},
				{
					ProductName: "bar",
					Product: params.Product{
						Build: params.Build{
							MainPkg: "./bar",
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
							OutputDir: "build",
						},
					},
					VersionInfo: git.ProjectInfo{
						Version: "0.1.0",
					},
				},
			},
		},
	} {
		currTmpDir, err := ioutil.TempDir(tmp, "")
		require.NoError(t, err)

		for file, content := range currCase.mainFiles {
			err := os.MkdirAll(path.Join(currTmpDir, path.Dir(file)), 0755)
			require.NoError(t, err)
			err = ioutil.WriteFile(path.Join(currTmpDir, file), []byte(content), 0644)
			require.NoError(t, err)
		}

		for i := range currCase.specs {
			currCase.specs[i].ProjectDir = currTmpDir
		}

		err = build.Run(currCase.specs, nil, build.Context{
			Parallel: true,
		}, ioutil.Discard)
		assert.NoError(t, err, "Case %d", i)
	}
}
