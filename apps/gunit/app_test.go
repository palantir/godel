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

package gunit_test

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"regexp"
	"strings"
	"testing"

	"github.com/nmiyake/pkg/dirs"
	"github.com/nmiyake/pkg/gofiles"
	"github.com/palantir/amalgomate/amalgomated"
	"github.com/palantir/pkg/pkgpath"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/palantir/godel/apps/gunit"
)

func TestMain(m *testing.M) {
	code := testHelper(m)
	os.Exit(code)
}

var supplier amalgomated.CmderSupplier

func testHelper(m *testing.M) int {
	tmpDir, cleanup, err := dirs.TempDir("", "")
	defer cleanup()
	if err != nil {
		panic(err)
	}

	libraries := map[string]string{
		"gocover":       "./vendor/github.com/nmiyake/gotest",
		"gojunitreport": "./vendor/github.com/jstemmer/go-junit-report",
		"gotest":        "./vendor/github.com/nmiyake/gotest",
		"gt":            "./vendor/rsc.io/gt",
	}

	supplier = temporaryBuildRunnerSupplier(tmpDir, libraries)
	return m.Run()
}

func TestRun(t *testing.T) {
	wd, err := os.Getwd()
	require.NoError(t, err)

	tmpDir, cleanup, err := dirs.TempDir(wd, "")
	defer cleanup()
	require.NoError(t, err)

	originalWd, err := os.Getwd()
	require.NoError(t, err)
	defer func() {
		if err := os.Chdir(originalWd); err != nil {
			fmt.Printf("Failed to set wd to %v: %v", originalWd, err)
		}
	}()

	for i, currCase := range []struct {
		name          string
		filesToCreate []gofiles.GoFileSpec
		config        string
		args          []string
		wantMatch     func(currCaseTmpDir string) string
		wantError     string
	}{
		{
			name: "passing tests",
			filesToCreate: []gofiles.GoFileSpec{
				{
					RelPath: "foo.go",
					Src: `package foo
					import "fmt"
					func Foo() {
						fmt.Println("Foo")
					}`,
				},
				{
					RelPath: "foo_test.go",
					Src: `package foo
					import "testing"
					func TestFoo(t *testing.T) {
						Foo()
					}`,
				},
				{
					RelPath: "vendor/bar/bar.go",
					Src: `package bar
					import "fmt"
					func Bar() {
						fmt.Println("Bar")
					}`,
				},
			},
			config: unindent(`exclude:
					  paths:
					    - "vendor"
					`),
			wantMatch: func(currCaseTmpDir string) string {
				return "ok  \t" + pkgName(t, currCaseTmpDir) + "\t[0-9.]+s"
			},
		},
		{
			name: "failing tests",
			filesToCreate: []gofiles.GoFileSpec{
				{
					RelPath: "foo.go",
					Src: `package foo
					import "fmt"
					func Foo() {
						fmt.Println("Foo")
					}`,
				},
				{
					RelPath: "foo_test.go",
					Src: `package foo
					import "testing"
					func TestFoo(t *testing.T) {
						Foo()
						t.Errorf("myFail")
					}`,
				},
			},
			wantMatch: func(currCaseTmpDir string) string {
				return `(?s)` +
					`Foo\n--- FAIL: TestFoo (.+)\n.+foo_test.go:[0-9]+: myFail.+FAIL\t` + pkgName(t, currCaseTmpDir) + "\t[0-9.]+s"
			},
			wantError: "(?s).+1 package had failing tests:.+",
		},
		{
			name: "test that does not compile fails",
			filesToCreate: []gofiles.GoFileSpec{
				{
					RelPath: "foo.go",
					Src: `package foo
					import "fmt"
					func Foo() {
						fmt.Println("Foo")
					}`,
				},
				{
					RelPath: "foo_test.go",
					Src: `package foo
					import "testing"
					import "github.com/palantir/godel/apps/gunit/blah/foo"
					func TestFoo(t *testing.T) {
						foo.Foo()
						t.Errorf("myFail")
					}`,
				},
			},
			wantMatch: func(currCaseTmpDir string) string {
				return `(?s)` +
					`foo_test.go:[0-9]+:[0-9]+: cannot find package.+\nFAIL\t` + pkgName(t, currCaseTmpDir) + `\t\[setup failed\]`
			},
			wantError: "(?s).+1 package had failing tests:.+",
		},
		{
			name: "running with a tag runs only tagged tests",
			filesToCreate: []gofiles.GoFileSpec{
				{
					RelPath: "foo_test.go",
					Src: `package foo
					import "testing"
					func TestFoo(t *testing.T) {
						t.Errorf("fooFail")
					}`,
				},
				{
					RelPath: "integration/bar_test.go",
					Src: `package bar
					import "testing"
					func TestBar(t *testing.T) {
						t.Errorf("barFail")
					}`,
				},
			},
			config: unindent(`tags:
					  integration:
					    names:
					      - "integration"
					exclude:
					  paths:
					    - "vendor"
					`),
			args: []string{
				"--tags", "integration",
			},
			wantMatch: func(currCaseTmpDir string) string {
				return `(?s)` +
					`--- FAIL: TestBar (.+)\n.+bar_test.go:[0-9]+: barFail.+FAIL\t` + pkgName(t, currCaseTmpDir) + "/integration\t[0-9.]+s"
			},
			wantError: "(?s).+1 package had failing tests:.+",
		},
		{
			name: "union of tags is run when multiple tags are specified",
			filesToCreate: []gofiles.GoFileSpec{
				{
					RelPath: "foo_test.go",
					Src: `package foo
					import "testing"
					func TestFoo(t *testing.T) {
						t.Errorf("fooFail")
					}`,
				},
				{
					RelPath: "bar/bar_test.go",
					Src: `package bar
					import "testing"
					func TestBar(t *testing.T) {
						t.Errorf("barFail")
					}`,
				},
				{
					RelPath: "baz/baz_test.go",
					Src: `package baz
					import "testing"
					func TestBaz(t *testing.T) {
						t.Errorf("bazFail")
					}`,
				},
			},
			config: unindent(`tags:
					  bar:
					    paths:
					      - "bar"
					  baz:
					    paths:
					      - "baz"
					exclude:
					  paths:
					    - "vendor"
					`),
			args: []string{
				"--tags", "bar,baz",
			},
			wantMatch: func(currCaseTmpDir string) string {
				return `(?s)` +
					`--- FAIL: TestBar (.+)\n.+bar_test.go:[0-9]+: barFail.+FAIL\t` + pkgName(t, currCaseTmpDir) + `/bar\t[0-9.]+s.+` +
					`--- FAIL: TestBaz (.+)\n.+baz_test.go:[0-9]+: bazFail.+FAIL\t` + pkgName(t, currCaseTmpDir) + `/baz\t[0-9.]+s.+`
			},
			wantError: "(?s).+2 packages had failing tests:.+",
		},
		{
			name: "only non-tagged tests are run if multiple tags are specified and tests are run with none argument",
			filesToCreate: []gofiles.GoFileSpec{
				{
					RelPath: "foo_test.go",
					Src: `package foo
					import "testing"
					func TestFoo(t *testing.T) {
						t.Errorf("fooFail")
					}`,
				},
				{
					RelPath: "bar/bar_test.go",
					Src: `package bar
					import "testing"
					func TestBar(t *testing.T) {
						t.Errorf("barFail")
					}`,
				},
				{
					RelPath: "baz/baz_test.go",
					Src: `package baz
					import "testing"
					func TestBaz(t *testing.T) {
						t.Errorf("bazFail")
					}`,
				},
			},
			config: unindent(`tags:
					  bar:
					    paths:
					      - "bar"
					  baz:
					    paths:
					      - "baz"
					exclude:
					  paths:
					    - "vendor"
					`),
			args: []string{
				"--tags", "none",
			},
			wantMatch: func(currCaseTmpDir string) string {
				return `(?s)` +
					`--- FAIL: TestFoo (.+)\n.+foo_test.go:[0-9]+: fooFail.+FAIL\t` + pkgName(t, currCaseTmpDir) + `\t[0-9.]+s.+`
			},
			wantError: "(?s).+1 package had failing tests:.+",
		},
		{
			name: "all tests are run if multiple tags are specified and tests are run without arguments",
			filesToCreate: []gofiles.GoFileSpec{
				{
					RelPath: "foo_test.go",
					Src: `package foo
					import "testing"
					func TestFoo(t *testing.T) {
						t.Errorf("fooFail")
					}`,
				},
				{
					RelPath: "bar/bar_test.go",
					Src: `package bar
					import "testing"
					func TestBar(t *testing.T) {
						t.Errorf("barFail")
					}`,
				},
				{
					RelPath: "baz/baz_test.go",
					Src: `package baz
					import "testing"
					func TestBaz(t *testing.T) {
						t.Errorf("bazFail")
					}`,
				},
			},
			config: unindent(`tags:
					  bar:
					    paths:
					      - "bar"
					  baz:
					    paths:
					      - "baz"
					exclude:
					  paths:
					    - "vendor"
					`),
			wantMatch: func(currCaseTmpDir string) string {
				return `(?s)` +
					`--- FAIL: TestFoo (.+)\n.+foo_test.go:[0-9]+: fooFail.+FAIL\s+` + pkgName(t, currCaseTmpDir) + `\s+[0-9.]+s.+` +
					`--- FAIL: TestBar (.+)\n.+bar_test.go:[0-9]+: barFail.+FAIL\s+` + pkgName(t, currCaseTmpDir) + `/bar\s+[0-9.]+s.+` +
					`--- FAIL: TestBaz (.+)\n.+baz_test.go:[0-9]+: bazFail.+FAIL\s+` + pkgName(t, currCaseTmpDir) + `/baz\s+[0-9.]+s.+`
			},
			wantError: "(?s).+3 packages had failing tests:.+",
		},
		{
			name: "fails if invalid tag is supplied",
			filesToCreate: []gofiles.GoFileSpec{
				{
					RelPath: "foo_test.go",
					Src: `package foo
					import "testing"
					func TestFoo(t *testing.T) {
						t.Errorf("fooFail")
					}`,
				},
			},
			config: unindent(`exclude:
					  paths:
					    - "vendor"
					`),
			args: []string{
				"--tags", "invalid,n!otvalid",
			},
			wantMatch: func(currCaseTmpDir string) string {
				return `Tags "invalid", "n!otvalid" not defined in configuration. No tags are defined.`
			},
			wantError: `Tags "invalid", "n!otvalid" not defined in configuration. No tags are defined.`,
		},
		{
			name: "fails if invalid tag is supplied and tag exists",
			filesToCreate: []gofiles.GoFileSpec{
				{
					RelPath: "foo_test.go",
					Src: `package foo
					import "testing"
					func TestFoo(t *testing.T) {
						t.Errorf("fooFail")
					}`,
				},
			},
			config: unindent(`tags:
					  bar:
					    paths:
					      - "bar"
					  exclude:
					    paths:
					      - "vendor"
					  other:
					    paths:
					      - "other"
					`),
			args: []string{
				"--tags", "invalid,n!otvalid",
			},
			wantMatch: func(currCaseTmpDir string) string {
				return `Tags "invalid", "n!otvalid" not defined in configuration. Valid tags: "bar", "exclude", "other"`
			},
			wantError: `Tags "invalid", "n!otvalid" not defined in configuration. Valid tags: "bar", "exclude", "other"`,
		},
	} {
		currCaseTmpDir, err := ioutil.TempDir(tmpDir, "")
		require.NoError(t, err, "Case %d: %s", i, currCase.name)

		_, err = gofiles.Write(currCaseTmpDir, currCase.filesToCreate)
		require.NoError(t, err, "Case %d: %s", i, currCase.name)

		configFile := path.Join(currCaseTmpDir, "config.yml")
		err = ioutil.WriteFile(configFile, []byte(currCase.config), 0644)
		require.NoError(t, err, "Case %d: %s", i, currCase.name)

		outBuf := bytes.Buffer{}

		err = os.Chdir(currCaseTmpDir)
		require.NoError(t, err)

		app := gunit.App(supplier)
		app.Stdout = &outBuf
		app.Stderr = &outBuf

		args := []string{"gunit", "--config", configFile}
		args = append(args, currCase.args...)
		args = append(args, "test")

		runTestCode := app.Run(args)
		output := outBuf.String()

		if runTestCode != 0 {
			if currCase.wantError == "" {
				t.Fatalf("Case %d: %s\nunexpected error:\n%v\nOutput: %v", i, currCase.name, err, output)
			} else if !regexp.MustCompile(currCase.wantError).MatchString(output) {
				t.Fatalf("Case %d: %s\nexpected error output to contain %v, but was %v", i, currCase.name, currCase.wantError, output)
			}
		} else if currCase.wantError != "" {
			t.Fatalf("Case %d: %s\nexpected error %v, but was none.\nOutput: %v", i, currCase.name, currCase.wantError, output)
		}

		expectedExpr := currCase.wantMatch(currCaseTmpDir)
		if !regexp.MustCompile(expectedExpr).MatchString(output) {
			t.Errorf("Case %d: %s\nOutput did not match expected expression.\nExpected:\n%v\nActual:\n%v", i, currCase.name, expectedExpr, output)
		}
	}
}

func TestClean(t *testing.T) {
	wd, err := os.Getwd()
	require.NoError(t, err)

	tmpDir, cleanup, err := dirs.TempDir(wd, "")
	defer cleanup()
	require.NoError(t, err)

	err = ioutil.WriteFile(path.Join(tmpDir, "tmp_placeholder_test.go"), []byte("package main"), 0644)
	require.NoError(t, err)

	originalWd, err := os.Getwd()
	require.NoError(t, err)
	defer func() {
		if err := os.Chdir(originalWd); err != nil {
			fmt.Printf("%+v\n", errors.Wrapf(err, "failed to restore working directory to %s", originalWd))
		}
	}()

	err = os.Chdir(tmpDir)
	require.NoError(t, err)

	app := gunit.App(supplier)
	exitCode := app.Run([]string{"gunit", "clean"})
	require.Equal(t, 0, exitCode)

	_, err = os.Stat(path.Join(tmpDir, "tmp_placeholder_test.go"))
	assert.True(t, os.IsNotExist(err))
}

func temporaryBuildRunnerSupplier(tmpDir string, libraries map[string]string) amalgomated.CmderSupplier {
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	// build executables for all commands once
	for currCmd, currLibrary := range libraries {
		executable := path.Join(tmpDir, currCmd)
		cmd := exec.Command("go", "build", "-o", executable)
		cmd.Dir = path.Join(wd, currLibrary)

		bytes, err := cmd.CombinedOutput()
		if err != nil {
			panic(fmt.Errorf("go build failed\nPath: %v\nArgs: %v\nDir: %v\nOutput: %v", cmd.Path, cmd.Args, cmd.Dir, string(bytes)))
		}
	}

	// return supplier that points to built executables
	return func(cmd amalgomated.Cmd) (amalgomated.Cmder, error) {
		return amalgomated.PathCmder(path.Join(tmpDir, cmd.Name())), nil
	}
}

func unindent(input string) string {
	return strings.Replace(input, "\n\t\t\t\t\t", "\n", -1)
}

func pkgName(t *testing.T, path string) string {
	pkgPath, err := pkgpath.NewAbsPkgPath(path).GoPathSrcRel()
	require.NoError(t, err)
	return pkgPath
}
