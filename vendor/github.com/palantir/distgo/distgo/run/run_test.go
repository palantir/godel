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

package run_test

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"regexp"
	"strings"
	"testing"

	"github.com/nmiyake/pkg/dirs"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/palantir/distgo/dister/disterfactory"
	"github.com/palantir/distgo/distgo"
	distgoconfig "github.com/palantir/distgo/distgo/config"
	"github.com/palantir/distgo/distgo/run"
	"github.com/palantir/distgo/dockerbuilder/dockerbuilderfactory"
)

const (
	runTestMain = `package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
)

func main() {
	fmt.Println("testMainOutput")
	ioutil.WriteFile(path.Join("{{OUTPUT_PATH}}", "runTestMainOutput.txt"), []byte(fmt.Sprintf("%v", os.Args[1:])), 0644)
}
`
	runTestMainBar = `package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
)

func main() {
	bar("testMainOutput")
	ioutil.WriteFile(path.Join("{{OUTPUT_PATH}}", "runTestMainOutput.txt"), []byte(fmt.Sprintf("%v", os.Args[1:])), 0644)
}
`
)

func TestRun(t *testing.T) {
	tmp, cleanup, err := dirs.TempDir("", "")
	defer cleanup()
	require.NoError(t, err)

	for i, tc := range []struct {
		name          string
		productConfig distgoconfig.ProductConfig
		runArgs       []string
		preRunAction  func(projectDir string)
		validate      func(runErr error, caseNum int, projectDir string)
	}{
		{
			`"run" runs main file`,
			distgoconfig.ProductConfig{
				Build: distgoconfig.ToBuildConfig(&distgoconfig.BuildConfig{
					MainPkg: stringPtr("."),
				}),
			},
			nil,
			func(projectDir string) {
				err := ioutil.WriteFile(path.Join(projectDir, "main.go"), []byte(strings.Replace(runTestMain, "{{OUTPUT_PATH}}", projectDir, -1)), 0644)
				require.NoError(t, err)
			},
			func(runErr error, caseNum int, projectDir string) {
				assert.NoError(t, runErr, "Case %d", caseNum)
				bytes, err := ioutil.ReadFile(path.Join(projectDir, "runTestMainOutput.txt"))
				require.NoError(t, err, "Case %d", caseNum)
				assert.Equal(t, "[]", string(bytes))
			},
		},
		{
			`"run" uses arguments provided in configuration (but does not evaluate them)`,
			distgoconfig.ProductConfig{
				Build: distgoconfig.ToBuildConfig(&distgoconfig.BuildConfig{
					MainPkg: stringPtr("."),
				}),
				Run: distgoconfig.ToRunConfig(&distgoconfig.RunConfig{
					Args: &[]string{
						"foo",
						"bar",
						"$GOPATH",
					},
				}),
			},
			nil,
			func(projectDir string) {
				err := ioutil.WriteFile(path.Join(projectDir, "main.go"), []byte(strings.Replace(runTestMain, "{{OUTPUT_PATH}}", projectDir, -1)), 0644)
				require.NoError(t, err)
			},
			func(runErr error, caseNum int, projectDir string) {
				assert.NoError(t, runErr, "Case %d", caseNum)
				bytes, err := ioutil.ReadFile(path.Join(projectDir, "runTestMainOutput.txt"))
				require.NoError(t, err, "Case %d", caseNum)
				assert.Equal(t, "[foo bar $GOPATH]", string(bytes))
			},
		},
		{
			`"run" uses arguments provided in slice`,
			distgoconfig.ProductConfig{
				Build: distgoconfig.ToBuildConfig(&distgoconfig.BuildConfig{
					MainPkg: stringPtr("."),
				}),
			},
			[]string{"foo", "bar", "$GOPATH"},
			func(projectDir string) {
				err := ioutil.WriteFile(path.Join(projectDir, "main.go"), []byte(strings.Replace(runTestMain, "{{OUTPUT_PATH}}", projectDir, -1)), 0644)
				require.NoError(t, err)
			},
			func(runErr error, caseNum int, projectDir string) {
				assert.NoError(t, runErr, "Case %d", caseNum)
				bytes, err := ioutil.ReadFile(path.Join(projectDir, "runTestMainOutput.txt"))
				require.NoError(t, err, "Case %d", caseNum)
				assert.Equal(t, "[foo bar $GOPATH]", string(bytes))
			},
		},
		{
			`"run" combines arguments in configuration with provided arguments`,
			distgoconfig.ProductConfig{
				Build: distgoconfig.ToBuildConfig(&distgoconfig.BuildConfig{
					MainPkg: stringPtr("."),
				}),
				Run: distgoconfig.ToRunConfig(&distgoconfig.RunConfig{
					Args: &[]string{
						"cfgArg_foo",
						"cfgArg_bar",
						"$cfgArg",
					},
				}),
			},
			[]string{"runArg_foo", "runArg_bar", "$runArg"},
			func(projectDir string) {
				err := ioutil.WriteFile(path.Join(projectDir, "main.go"), []byte(strings.Replace(runTestMain, "{{OUTPUT_PATH}}", projectDir, -1)), 0644)
				require.NoError(t, err)
			},
			func(runErr error, caseNum int, projectDir string) {
				assert.NoError(t, runErr, "Case %d", caseNum)
				bytes, err := ioutil.ReadFile(path.Join(projectDir, "runTestMainOutput.txt"))
				require.NoError(t, err, "Case %d", caseNum)
				assert.Equal(t, "[cfgArg_foo cfgArg_bar $cfgArg runArg_foo runArg_bar $runArg]", string(bytes))
			},
		},
		{
			`"run" uses build arguments specified in build configuration`,
			distgoconfig.ProductConfig{
				Build: distgoconfig.ToBuildConfig(&distgoconfig.BuildConfig{
					MainPkg:    stringPtr("."),
					VersionVar: stringPtr("main.testVersionVar"),
				}),
			},
			nil,
			func(projectDir string) {
				currMainContent := `package main

import (
	"io/ioutil"
	"path"
)

var testVersionVar = "defaultVersion"

func main() {
	ioutil.WriteFile(path.Join("{{OUTPUT_PATH}}", "runTestMainOutput.txt"), []byte(testVersionVar), 0644)
}
`
				err := ioutil.WriteFile(path.Join(projectDir, "main.go"), []byte(strings.Replace(currMainContent, "{{OUTPUT_PATH}}", projectDir, -1)), 0644)
				require.NoError(t, err)
			},
			func(runErr error, caseNum int, projectDir string) {
				assert.NoError(t, runErr, "Case %d", caseNum)
				bytes, err := ioutil.ReadFile(path.Join(projectDir, "runTestMainOutput.txt"))
				require.NoError(t, err, "Case %d", caseNum)
				assert.Equal(t, "0.1.0", string(bytes))
			},
		},
		{
			`"run" works with multiple main package files as long as there is a single main function`,
			distgoconfig.ProductConfig{
				Build: distgoconfig.ToBuildConfig(&distgoconfig.BuildConfig{
					MainPkg: stringPtr("./foo"),
				}),
			},
			nil,
			func(projectDir string) {
				err := os.MkdirAll(path.Join(projectDir, "foo"), 0755)
				require.NoError(t, err)

				err = ioutil.WriteFile(path.Join(projectDir, "foo", "main_file.go"), []byte(strings.Replace(runTestMainBar, "{{OUTPUT_PATH}}", projectDir, -1)), 0644)
				require.NoError(t, err)
				err = ioutil.WriteFile(path.Join(projectDir, "foo", "other_main_file.go"), []byte(`package main
import "fmt"
func bar(a ...interface{}) (n int, err error) {
	return fmt.Println(a...)
}
`), 0644)
				require.NoError(t, err)
				err = ioutil.WriteFile(path.Join(projectDir, "foo", "main_test.go"), []byte(`package main_test
func Bar() string {
	return "bar"
}
`), 0644)
				require.NoError(t, err)
			},
			func(runErr error, caseNum int, projectDir string) {
				assert.NoError(t, runErr, "Case %d", caseNum)
				bytes, err := ioutil.ReadFile(path.Join(projectDir, "runTestMainOutput.txt"))
				require.NoError(t, err, "Case %d", caseNum)
				assert.Equal(t, "[]", string(bytes))
			},
		},
		{
			`"run" works with multiple main package files with tests`,
			distgoconfig.ProductConfig{
				Build: distgoconfig.ToBuildConfig(&distgoconfig.BuildConfig{
					MainPkg: stringPtr("./foo"),
				}),
			},
			nil,
			func(projectDir string) {
				err := os.MkdirAll(path.Join(projectDir, "foo"), 0755)
				require.NoError(t, err)

				err = ioutil.WriteFile(path.Join(projectDir, "foo", "main_file.go"), []byte(strings.Replace(runTestMain, "{{OUTPUT_PATH}}", projectDir, -1)), 0644)
				require.NoError(t, err)
				err = ioutil.WriteFile(path.Join(projectDir, "foo", "main_test.go"), []byte(`package main
import "testing"
func TestBar(t *testing.T) {
}
`), 0644)
				require.NoError(t, err)
			},
			func(runErr error, caseNum int, projectDir string) {
				assert.NoError(t, runErr, "Case %d", caseNum)
				bytes, err := ioutil.ReadFile(path.Join(projectDir, "runTestMainOutput.txt"))
				require.NoError(t, err, "Case %d", caseNum)
				assert.Equal(t, "[]", string(bytes))
			},
		},
		{
			`"run" fails if a main package does not exist`,
			distgoconfig.ProductConfig{
				Build: distgoconfig.ToBuildConfig(&distgoconfig.BuildConfig{
					MainPkg: stringPtr("./foo"),
				}),
			},
			nil,
			func(projectDir string) {
				err := os.MkdirAll(path.Join(projectDir, "foo"), 0755)
				require.NoError(t, err)

				err = ioutil.WriteFile(path.Join(projectDir, "foo", "not_main_pkg.go"), []byte(`package foo
func main() {
}
`), 0644)
				require.NoError(t, err)
			},
			func(runErr error, caseNum int, projectDir string) {
				assert.Error(t, runErr, fmt.Sprintf("Case %d", caseNum))
				assert.Regexp(t, regexp.MustCompile(`^failed to find Go files for main package: no go file with main package and main function exists in directory .+/foo$`), runErr.Error(), "Case %d", caseNum)
			},
		},
		{
			`"run" fails if main function does not exist in a main pkg`,
			distgoconfig.ProductConfig{
				Build: distgoconfig.ToBuildConfig(&distgoconfig.BuildConfig{
					MainPkg: stringPtr("./foo"),
				}),
			},
			nil,
			func(projectDir string) {
				err := os.MkdirAll(path.Join(projectDir, "foo"), 0755)
				require.NoError(t, err)

				err = ioutil.WriteFile(path.Join(projectDir, "foo", "no_main_func.go"), []byte(`package main
func Foo() string {
	return "foo"
}
`), 0644)
				require.NoError(t, err)

				err = ioutil.WriteFile(path.Join(projectDir, "foo", "main_func_not_main_pkg.go"), []byte(`package main_test
func main() {
}
`), 0644)
				require.NoError(t, err)
			},
			func(runErr error, caseNum int, projectDir string) {
				assert.Error(t, runErr, fmt.Sprintf("Case %d", caseNum))
				assert.Regexp(t, regexp.MustCompile(`^failed to find Go files for main package: no go file with main package and main function exists in directory .+/foo$`), runErr.Error(), "Case %d", caseNum)
			},
		},
	} {
		projectDir, err := ioutil.TempDir(tmp, "")
		require.NoError(t, err, "Case %d: %s", i, tc.name)

		if tc.preRunAction != nil {
			tc.preRunAction(projectDir)
		}

		projectInfo := distgo.ProjectInfo{
			ProjectDir: projectDir,
			Version:    "0.1.0",
		}

		disterFactory, err := disterfactory.New(nil, nil)
		require.NoError(t, err, "Case %d: %s", i, tc.name)
		dockerBuilderFactory, err := dockerbuilderfactory.New(nil, nil)
		require.NoError(t, err, "Case %d: %s", i, tc.name)

		productParam, err := tc.productConfig.ToParam("foo", "", distgoconfig.ProductConfig{}, disterFactory, dockerBuilderFactory)
		require.NoError(t, err, "Case %d: %s", i, tc.name)

		err = run.Product(projectInfo, productParam, tc.runArgs, ioutil.Discard, ioutil.Discard)
		if tc.validate != nil {
			tc.validate(err, i, projectDir)
		}
	}
}

func stringPtr(in string) *string {
	return &in
}
