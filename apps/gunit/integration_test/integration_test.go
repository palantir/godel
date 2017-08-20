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

package integration_test

import (
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"regexp"
	"strings"
	"testing"

	"github.com/nmiyake/pkg/dirs"
	"github.com/nmiyake/pkg/gofiles"
	"github.com/palantir/pkg/pkgpath"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/palantir/godel/pkg/products"
)

func TestRun(t *testing.T) {
	cli, err := products.Bin("gunit")
	require.NoError(t, err)

	wd, err := os.Getwd()
	require.NoError(t, err)

	tmpDir, cleanup, err := dirs.TempDir(wd, "")
	defer cleanup()
	require.NoError(t, err)

	originalWd, err := os.Getwd()
	require.NoError(t, err)
	defer func() {
		if err := os.Chdir(originalWd); err != nil {
			require.NoError(t, err)
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
			name: "fails if there are no packages to test (no buildable Go files)",
			wantMatch: func(currCaseTmpDir string) string {
				return "^no packages to test\n$"
			},
			wantError: "^no packages to test\n$",
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
			name: "tag skips excluded tests",
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
				{
					RelPath: "integration/baz/baz_test.go",
					Src: `package baz
					import "testing"
					func TestBaz(t *testing.T) {
						t.Errorf("bazFail")
					}`,
				},
				{
					RelPath: "integration/exclude/exclude_test.go",
					Src: `package exclude
					import "testing"
					func TestExclude(t *testing.T) {
						t.Errorf("exclude")
					}`,
				},
			},
			config: unindent(`tags:
					  integration:
					    names:
					      - "integration"
					    exclude:
					      paths:
					        - "integration/exclude"
					exclude:
					  paths:
					    - "vendor"
					`),
			args: []string{
				"--tags", "integration",
			},
			wantMatch: func(currCaseTmpDir string) string {
				return `(?s)` +
					`--- FAIL: TestBar (.+)\n.+bar_test.go:[0-9]+: barFail.+FAIL\t` + pkgName(t, currCaseTmpDir) + `/integration\s+[0-9.]+s.+` +
					`--- FAIL: TestBaz (.+)\n.+baz_test.go:[0-9]+: bazFail.+FAIL\t` + pkgName(t, currCaseTmpDir) + `/integration/baz\s+[0-9.]+s.+`
			},
			wantError: "(?s).+2 packages had failing tests:.+",
		},
		{
			name: "tags are case-insensitive",
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
				"--tags", "INTEGRATION",
			},
			wantMatch: func(currCaseTmpDir string) string {
				return `(?s)` +
					`--- FAIL: TestBar (.+)\n.+bar_test.go:[0-9]+: barFail.+FAIL\t` + pkgName(t, currCaseTmpDir) + "/integration\t[0-9.]+s"
			},
			wantError: "(?s).+1 package had failing tests:.+",
		},
		{
			name: "tags can be specified multiple times",
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
				"--tags", "INTEGRATION,integration",
			},
			wantMatch: func(currCaseTmpDir string) string {
				return `(?s)` +
					`--- FAIL: TestBar (.+)\n.+bar_test.go:[0-9]+: barFail.+FAIL\t` + pkgName(t, currCaseTmpDir) + "/integration\t[0-9.]+s"
			},
			wantError: "(?s).+1 package had failing tests:.+",
		},
		{
			name: "fails if tags do not match any packages",
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
					  integration:
					    names:
					      - "integration"
					`),
			args: []string{
				"--tags", "integration",
			},
			wantMatch: func(currCaseTmpDir string) string {
				return "^no packages to test\n$"
			},
			wantError: "^no packages to test\n$",
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
			name: "only non-tagged tests are run if none is specified as tag",
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
			name: "only tagged tests are run if all is specified as tag",
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
				"--tags", "all",
			},
			wantMatch: func(currCaseTmpDir string) string {
				return `(?s)` +
					`--- FAIL: TestBar (.+)\n.+bar_test.go:[0-9]+: barFail.+FAIL\t` + pkgName(t, currCaseTmpDir) + `/bar\s+[0-9.]+s.+` +
					`--- FAIL: TestBaz (.+)\n.+baz_test.go:[0-9]+: bazFail.+FAIL\t` + pkgName(t, currCaseTmpDir) + `/baz\s+[0-9.]+s.+`
			},
			wantError: "(?s).+2 packages had failing tests:.+",
		},
		{
			name: "tags specified with uppercase letters works",
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
			},
			config: unindent(`tags:
					  Bar:
					    paths:
					      - "bar"
					exclude:
					  paths:
					    - "vendor"
					`),
			args: []string{
				"--tags", "Bar",
			},
			wantMatch: func(currCaseTmpDir string) string {
				return `(?s)` +
					`--- FAIL: TestBar (.+)\n.+bar_test.go:[0-9]+: barFail.+FAIL\t` + pkgName(t, currCaseTmpDir) + `/bar\t[0-9.]+s.+`
			},
			wantError: "(?s).+1 package had failing tests:.+",
		},
		{
			name: "all tests (tagged and non-tagged) are run if tags are not specified",
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
					`--- FAIL: TestFoo (.+)\n.+foo_test.go:[0-9]+: fooFail.+FAIL\t` + pkgName(t, currCaseTmpDir) + `\s+[0-9.]+s.+` +
					`--- FAIL: TestBar (.+)\n.+bar_test.go:[0-9]+: barFail.+FAIL\t` + pkgName(t, currCaseTmpDir) + `/bar\s+[0-9.]+s.+` +
					`--- FAIL: TestBaz (.+)\n.+baz_test.go:[0-9]+: bazFail.+FAIL\t` + pkgName(t, currCaseTmpDir) + `/baz\s+[0-9.]+s.+`
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
			name: "fails if 'all' is supplied as a non-exclusive tag",
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
				"--tags", "integration,all",
			},
			wantError: regexp.QuoteMeta(`if "all" tag is specified, it must be the only tag specified`),
		},
		{
			name: "detects race conditions if race flag is supplied",
			filesToCreate: []gofiles.GoFileSpec{
				{
					RelPath: "foo_test.go",
					Src: `package foo
					import (
					  "fmt"
					  "math/rand"
					  "testing"
					  "time"
					)
					func TestFoo(t *testing.T) {
					    timeTest()
					}
					func timeTest() {
					    start := time.Now()
					     var t *time.Timer
					     t = time.AfterFunc(randomDuration(), func() {
						 fmt.Println(time.Now().Sub(start))
						 t.Reset(randomDuration())
					     })
					     time.Sleep(1 * time.Second)
					}
					func randomDuration() time.Duration {
					     return time.Duration(rand.Int63n(1e9))
					}`,
				},
			},
			args: []string{
				"--race",
			},
			wantMatch: func(currCaseTmpDir string) string {
				return `(?s).+WARNING: DATA RACE\n.+`
			},
			wantError: `(?s).+WARNING: DATA RACE\n.+`,
		},
	} {
		currCaseTmpDir, err := ioutil.TempDir(tmpDir, "")
		require.NoError(t, err, "Case %d: %s", i, currCase.name)

		_, err = gofiles.Write(currCaseTmpDir, currCase.filesToCreate)
		require.NoError(t, err, "Case %d: %s", i, currCase.name)

		configFile := path.Join(currCaseTmpDir, "config.yml")
		err = ioutil.WriteFile(configFile, []byte(currCase.config), 0644)
		require.NoError(t, err, "Case %d: %s", i, currCase.name)

		err = os.Chdir(currCaseTmpDir)
		require.NoError(t, err)

		args := []string{"--config", configFile}
		args = append(args, currCase.args...)
		args = append(args, "test")

		cmd := exec.Command(cli, args...)
		outputBytes, err := cmd.CombinedOutput()
		output := string(outputBytes)

		if err != nil {
			if currCase.wantError == "" {
				t.Fatalf("Case %d: %s\nunexpected error:\n%v\nOutput: %v", i, currCase.name, err, output)
			} else if !regexp.MustCompile(currCase.wantError).MatchString(output) {
				t.Fatalf("Case %d: %s\nexpected error output to contain %v, but was %v", i, currCase.name, currCase.wantError, output)
			}
		} else if currCase.wantError != "" {
			t.Fatalf("Case %d: %s\nexpected error %v, but was none.\nOutput: %v", i, currCase.name, currCase.wantError, output)
		}

		if currCase.wantMatch != nil {
			expectedExpr := currCase.wantMatch(currCaseTmpDir)
			if !regexp.MustCompile(expectedExpr).MatchString(output) {
				t.Errorf("Case %d: %s\nOutput did not match expected expression.\nExpected:\n%v\nActual:\n%v", i, currCase.name, expectedExpr, output)
			}
		}
	}
}

func TestClean(t *testing.T) {
	cli, err := products.Bin("gunit")
	require.NoError(t, err)

	wd, err := os.Getwd()
	require.NoError(t, err)

	tmpDir, cleanup, err := dirs.TempDir(wd, "")
	defer cleanup()
	require.NoError(t, err)

	files, err := gofiles.Write(tmpDir, []gofiles.GoFileSpec{
		{
			RelPath: "tmp_placeholder_test.go",
			Src:     "package main",
		},
	})
	require.NoError(t, err)

	originalWd, err := os.Getwd()
	require.NoError(t, err)
	defer func() {
		if err := os.Chdir(originalWd); err != nil {
			require.NoError(t, err)
		}
	}()

	err = os.Chdir(tmpDir)
	require.NoError(t, err)

	cmd := exec.Command(cli, "clean")
	output, err := cmd.CombinedOutput()
	require.NoError(t, err, "Command %v failed with output:\n%s", cmd.Args, string(output))

	_, err = os.Stat(files["tmp_placeholder_test.go"].Path)
	assert.True(t, os.IsNotExist(err))
}

func unindent(input string) string {
	return strings.Replace(input, "\n\t\t\t\t\t", "\n", -1)
}

func pkgName(t *testing.T, path string) string {
	pkgPath, err := pkgpath.NewAbsPkgPath(path).GoPathSrcRel()
	require.NoError(t, err)
	return pkgPath
}
