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
	"github.com/palantir/godel/pkg/osarch"
	"github.com/palantir/pkg/gittest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/palantir/distgo/distgo"
	"github.com/palantir/distgo/distgo/build"
	"github.com/palantir/distgo/pkg/git"
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
	longCompileMain  = `package main

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

	for i, tc := range []struct {
		productName     string
		mainFileContent string
		mainFilePath    string
		productParam    distgo.ProductParam
		wantError       bool
		runExecutable   bool
		wantOutput      string
	}{
		{
			productName:     "randomProduct",
			mainFileContent: testMain,
			mainFilePath:    "main.go",
			productParam: createBuildProductParam(func(param *distgo.ProductParam) {
				param.Build.VersionVar = "main.testVersionVar"
			}),
			runExecutable: true,
			wantOutput:    testVersionValue + ".dirty",
		},
		// building project that requires CGo succeeds if "CGO_ENABLED" environment variable is set to 1
		{
			productName:     "CProduct",
			mainFileContent: testCMain,
			mainFilePath:    "main.go",
			productParam: createBuildProductParam(func(param *distgo.ProductParam) {
				param.Build.Environment = map[string]string{
					"CGO_ENABLED": "1",
				}
			}),
			runExecutable: true,
			wantOutput:    "C",
		},
		// building project that requires CGo fails if "CGO_ENABLED" environment variable is set to 0
		{
			productName:     "CProduct",
			mainFileContent: testCMain,
			mainFilePath:    "main.go",
			productParam: createBuildProductParam(func(param *distgo.ProductParam) {
				param.Build.Environment = map[string]string{
					"CGO_ENABLED": "0",
				}
			}),
			wantError: true,
		},
		{
			productName:     "customBuildScriptProduct",
			mainFileContent: testMain,
			mainFilePath:    "main.go",
			productParam: createBuildProductParam(func(param *distgo.ProductParam) {
				param.Build.BuildArgsScript = `#!/usr/bin/env bash
set -eu pipefail
VALUE="foo bar"
echo "-ldflags"
echo "-X \"main.testVersionVar=$VALUE\""`
			}),
			runExecutable: true,
			wantOutput:    "foo bar",
		},
		{
			productName:     "foo",
			mainFileContent: testMain,
			mainFilePath:    "foo/main.go",
			productParam: createBuildProductParam(func(param *distgo.ProductParam) {
				param.Build.MainPkg = "./foo"
				param.Build.OSArchs = []osarch.OSArch{
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
				}
			}),
			wantOutput: "defaultVersion",
		},
	} {
		currTmpDir, err := ioutil.TempDir(tmp, "")
		require.NoError(t, err, "Case %d", i)

		gittest.InitGitDir(t, currTmpDir)
		gittest.CreateGitTag(t, currTmpDir, testVersionValue)

		mainFilePath := path.Join(currTmpDir, tc.mainFilePath)

		err = os.MkdirAll(path.Dir(mainFilePath), 0755)
		require.NoError(t, err, "Case %d", i)

		err = ioutil.WriteFile(mainFilePath, []byte(tc.mainFileContent), 0644)
		require.NoError(t, err, "Case %d", i)

		foundExecForCurrOSArch := false

		version, err := git.ProjectVersion(currTmpDir)
		require.NoError(t, err, "Case %d", i)

		projectInfo := distgo.ProjectInfo{
			ProjectDir: currTmpDir,
			Version:    version,
		}
		productOutputInfo, err := tc.productParam.ToProductOutputInfo(projectInfo.Version)
		require.NoError(t, err, "Case %d", i)

		err = build.Run(projectInfo, []distgo.ProductParam{
			tc.productParam,
		}, build.Options{
			Parallel: false,
		}, ioutil.Discard)

		if tc.wantError {
			assert.Error(t, err, fmt.Sprintf("Case %d", i))
		} else {
			assert.NoError(t, err, "Case %d", i)
			artifactPaths := distgo.ProductBuildArtifactPaths(projectInfo, productOutputInfo)
			for _, currOSArch := range tc.productParam.Build.OSArchs {
				pathToCurrExecutable, ok := artifactPaths[currOSArch]
				require.True(t, ok, "Case %d: could not find path for %s for %s", tc.productName, currOSArch.String())
				fileInfo, err := os.Stat(pathToCurrExecutable)
				require.NoError(t, err, "Case %d", i)
				assert.False(t, fileInfo.IsDir())

				if reflect.DeepEqual(currOSArch, osarch.Current()) {
					foundExecForCurrOSArch = true
					output, err := exec.Command(pathToCurrExecutable).Output()
					require.NoError(t, err, "Case %d", i)
					assert.Equal(t, tc.wantOutput, strings.TrimSpace(string(output)), "Case %d", i)
				}
			}

			if tc.runExecutable {
				assert.True(t, foundExecForCurrOSArch, "Case %d: executable for current os/arch (%v) not found in %v", osarch.Current(), tc.productParam.Build.OSArchs)
			}
		}
	}
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

	for i, tc := range []struct {
		specOSArchs []osarch.OSArch
		want        []string
		notWant     []string
	}{
		// empty value for osArchs filter builds all
		{
			specOSArchs: []osarch.OSArch{{OS: "darwin", Arch: "amd64"}, {OS: "linux", Arch: "386"}},
			want: []string{
				"Finished building testProduct for darwin-amd64",
				"Finished building testProduct for linux-386",
			},
		},
		// if non-empty filter is provided, only values matching filter are built
		{
			specOSArchs: []osarch.OSArch{{OS: "linux", Arch: "386"}},
			want: []string{
				"Finished building testProduct for linux-386",
			},
			notWant: []string{
				"Finished building testProduct for darwin-amd64",
			},
		},
		// if no OS/arch values match filter, nothing is built
		{
			specOSArchs: []osarch.OSArch{},
			want: []string{
				"$^",
			},
		},
	} {
		projectInfo := distgo.ProjectInfo{
			ProjectDir: tmp,
		}
		productParam := createBuildProductParam(func(param *distgo.ProductParam) {
			param.Build.MainPkg = "./foo"
			param.Build.OSArchs = tc.specOSArchs
		})

		buf := &bytes.Buffer{}
		err = build.Run(projectInfo, []distgo.ProductParam{productParam}, build.Options{
			Parallel: false,
		}, buf)
		require.NoError(t, err, "Case %d", i)

		for _, want := range tc.want {
			assert.Regexp(t, regexp.MustCompile(want), buf.String(), "Case %d", i)
		}

		for _, notWant := range tc.notWant {
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

	projectInfo := distgo.ProjectInfo{
		ProjectDir: tmp,
	}
	productParam := createBuildProductParam(func(param *distgo.ProductParam) {
		param.Build.MainPkg = "./foo"
	})

	want := fmt.Sprintf(`(?s)^go build failed: build command \[.+go build -i -o out/build/testProduct/%v/testProduct ./foo\] run in directory %s with additional environment variables \[GOOS=.+ GOARCH=.+\] failed with output:.+foo/main.go:1:15: syntax error: non-declaration statement outside function body$`,
		osarch.Current(), tmp)

	buf := &bytes.Buffer{}
	err = build.Run(projectInfo, []distgo.ProductParam{productParam}, build.Options{
		Install:  true,
		Parallel: false,
	}, buf)
	assert.Regexp(t, want, err.Error())
}

// TODO: run test in environment where current user is not root and re-enable
//func TestBuildInstallErrorMessage(t *testing.T) {
//	tmp, cleanup, err := dirs.TempDir(".", "")
//	defer cleanup()
//	require.NoError(t, err)
//
//	goRoot, err := dirs.GoRoot()
//	require.NoError(t, err)
//	_, err = os.Stat(goRoot)
//	require.NoError(t, err)
//
//	pkgDir := path.Join(goRoot, "pkg")
//	_, err = os.Stat(pkgDir)
//	require.NoError(t, err)
//
//	osArchPkgDir := path.Join(pkgDir, "dragonfly_amd64")
//	_, err = os.Stat(osArchPkgDir)
//	if os.IsNotExist(err) {
//		// if directory does not exist, attempt to create it (and clean up afterwards)
//		if err := os.Mkdir(osArchPkgDir, 0444); err == nil {
//			defer func() {
//				if err := os.RemoveAll(osArchPkgDir); err != nil {
//					fmt.Printf("Failed to remove directory %v: %v\n", osArchPkgDir, err)
//				}
//			}()
//		}
//		// if creation failed, assume that write permissions do not exist, which is sufficient for the test
//	}
//
//	mainFilePath := path.Join(tmp, "foo/main.go")
//	err = os.MkdirAll(path.Dir(mainFilePath), 0755)
//	require.NoError(t, err)
//	err = ioutil.WriteFile(mainFilePath, []byte(`package main`), 0644)
//	require.NoError(t, err)
//
//	projectInfo := distgo.ProjectInfo{
//		ProjectDir: tmp,
//	}
//	productParam := createBuildProductParam(func(param *distgo.ProductParam) {
//		param.Build.MainPkg = "./foo"
//		param.Build.OSArchs = []osarch.OSArch{
//			{
//				OS:   "dragonfly",
//				Arch: "amd64",
//			},
//		}
//	})
//
//	goBinary := "go"
//	if output, err := exec.Command("command", "-v", "go").CombinedOutput(); err == nil {
//		goBinary = strings.TrimSpace(string(output))
//	}
//
//	wantRegexps := []*regexp.Regexp{
//		regexp.MustCompile(`^` + regexp.QuoteMeta(`go build failed: failed to install a Go standard library package due to insufficient permissions to create directory.`) + `$`),
//		regexp.MustCompile(`^` + regexp.QuoteMeta(`This typically means that the standard library for the OS/architecture combination have not been installed locally and the current user does not have write permissions to GOROOT/pkg.`) + `$`),
//		regexp.MustCompile(`^` + regexp.QuoteMeta(fmt.Sprintf(`Run "sudo env GOOS=dragonfly GOARCH=amd64 %s install std" to install the standard packages for this combination as root and then try again.`, goBinary)) + `$`),
//		regexp.MustCompile(fmt.Sprintf(`Full error: build command \[.+/go build -i -o out/build/testProduct/dragonfly-amd64/testProduct ./foo\] run in directory %s with additional environment variables \[GOOS=dragonfly GOARCH=amd64\] failed with output:`, tmp)),
//		regexp.MustCompile(`go build [^:]+: mkdir [^:]+: permission denied`),
//	}
//
//	buf := &bytes.Buffer{}
//	err = build.Run(projectInfo, []distgo.ProductParam{productParam}, build.Options{
//		Install:  true,
//		Parallel: false,
//	}, buf)
//
//	parts := strings.Split(err.Error(), "\n")
//	for i := range parts {
//		if i >= len(wantRegexps) {
//			break
//		}
//		assert.Regexp(t, wantRegexps[i], parts[i])
//	}
//}

func TestBuildAllParallel(t *testing.T) {
	tmp, cleanup, err := dirs.TempDir("", "")
	defer cleanup()
	require.NoError(t, err)

	for i, tc := range []struct {
		mainFiles     map[string]string
		productParams []distgo.ProductParam
	}{
		{
			mainFiles: map[string]string{
				"foo/main.go": longCompileMain,
				"bar/main.go": longCompileMain,
			},
			productParams: []distgo.ProductParam{
				createBuildProductParam(func(param *distgo.ProductParam) {
					param.Build.MainPkg = "./foo"
					param.Build.OSArchs = []osarch.OSArch{
						{
							OS:   "darwin",
							Arch: "amd64",
						},
						{
							OS:   "linux",
							Arch: "amd64",
						},
					}
				}),
				createBuildProductParam(func(param *distgo.ProductParam) {
					param.Build.MainPkg = "./bar"
					param.Build.OSArchs = []osarch.OSArch{
						{
							OS:   "darwin",
							Arch: "amd64",
						},
						{
							OS:   "linux",
							Arch: "amd64",
						},
					}
				}),
			},
		},
	} {
		currTmpDir, err := ioutil.TempDir(tmp, "")
		require.NoError(t, err)

		for file, content := range tc.mainFiles {
			err := os.MkdirAll(path.Join(currTmpDir, path.Dir(file)), 0755)
			require.NoError(t, err)
			err = ioutil.WriteFile(path.Join(currTmpDir, file), []byte(content), 0644)
			require.NoError(t, err)
		}

		projectInfo := distgo.ProjectInfo{
			ProjectDir: currTmpDir,
			Version:    "0.1.0",
		}
		err = build.Run(projectInfo, tc.productParams, build.Options{
			Parallel: true,
		}, ioutil.Discard)
		assert.NoError(t, err, "Case %d", i)
	}
}

func createBuildProductParam(fn func(*distgo.ProductParam)) distgo.ProductParam {
	param := distgo.ProductParam{
		ID: "testProduct",
		Build: &distgo.BuildParam{
			NameTemplate: "{{Product}}",
			MainPkg:      ".",
			OutputDir:    "out/build",
			OSArchs: []osarch.OSArch{
				osarch.Current(),
			},
		},
	}
	if fn != nil {
		fn(&param)
	}
	return param
}
