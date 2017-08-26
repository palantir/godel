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
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"strings"
	"testing"

	"github.com/nmiyake/pkg/dirs"
	"github.com/nmiyake/pkg/gofiles"
	"github.com/palantir/amalgomate/amalgomated"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/palantir/godel/apps/okgo/checkoutput"
	"github.com/palantir/godel/apps/okgo/checks"
	"github.com/palantir/godel/apps/okgo/cmd/cmdlib"
	"github.com/palantir/godel/apps/okgo/config"
	"github.com/palantir/godel/apps/okgo/params"
	"github.com/palantir/godel/pkg/products"
)

func TestCheckers(t *testing.T) {
	cli, err := products.Bin("okgo")
	require.NoError(t, err)

	for i, currCase := range []struct {
		check amalgomated.Cmd
		want  []string
	}{
		{
			check: cmdlib.Instance().MustNewCmd("deadcode"),
			want: []string{
				"pkg1/bad.go:27:1: deadcode is unused",
				"pkg1/bad.go:40:1: varcheck is unused",
				"pkg2/bad2.go:27:1: deadcode is unused",
				"pkg2/bad2.go:40:1: varcheck is unused",
			},
		},
		{
			check: cmdlib.Instance().MustNewCmd("errcheck"),
			want: []string{
				"pkg1/bad.go:11:8: helper()",
				"pkg2/bad2.go:11:8: helper()",
			},
		},
		{
			check: cmdlib.Instance().MustNewCmd("golint"),
			want: []string{
				`pkg1/bad.go:49:1: comment on exported function Lint should be of the form "Lint ..."`,
				`pkg2/bad2.go:49:1: comment on exported function Lint should be of the form "Lint ..."`,
			},
		},
		{
			check: cmdlib.Instance().MustNewCmd("govet"),
			want: []string{
				"pkg1/bad.go:23: self-assignment of foo to foo",
				"pkg2/bad2.go:23: self-assignment of foo to foo",
			},
		},
		{
			check: cmdlib.Instance().MustNewCmd("importalias"),
			want: []string{
				`pkg1/bad.go:3:8: uses alias "myjson" to import package "encoding/json". No consensus alias exists for this import in the project ("ejson" and "myjson" are both used once each).`,
				`pkg2/bad2.go:3:8: uses alias "ejson" to import package "encoding/json". No consensus alias exists for this import in the project ("ejson" and "myjson" are both used once each).`,
			},
		},
		{
			check: cmdlib.Instance().MustNewCmd("ineffassign"),
			want: []string{
				"pkg1/bad.go:34:2: ineffectual assignment to kvs",
				"pkg1/bad.go:36:2: ineffectual assignment to kvs",
				"pkg2/bad2.go:34:2: ineffectual assignment to kvs",
				"pkg2/bad2.go:36:2: ineffectual assignment to kvs",
			},
		},
		{
			check: cmdlib.Instance().MustNewCmd("outparamcheck"),
			want: []string{
				`github.com/palantir/godel/apps/okgo/integration_test/testdata/standard/pkg1/bad.go:16:28: _ = myjson.Unmarshal(nil, "")  // 2nd argument of 'Unmarshal' requires '&'`,
				`github.com/palantir/godel/apps/okgo/integration_test/testdata/standard/pkg2/bad2.go:16:27: _ = ejson.Unmarshal(nil, "")  // 2nd argument of 'Unmarshal' requires '&'`,
			},
		},
		{
			check: cmdlib.Instance().MustNewCmd("unconvert"),
			want: []string{
				"pkg1/bad.go:45:14: unnecessary conversion",
				"pkg2/bad2.go:45:14: unnecessary conversion",
			},
		},
		{
			check: cmdlib.Instance().MustNewCmd("varcheck"),
			want: []string{
				"pkg1/bad.go:40:7: varcheck",
				"pkg2/bad2.go:40:7: varcheck",
			},
		},
	} {
		checker, err := checks.GetChecker(currCase.check)
		require.NoError(t, err)

		runner := amalgomated.PathCmder(cli, amalgomated.ProxyCmdPrefix+currCase.check.Name())
		lineInfo, err := checker.Check(runner, "./testdata/standard", params.OKGo{})
		require.NoError(t, err, "Case %d", i)

		assert.Equal(t, currCase.want, toStringSlice(lineInfo), "Case %d", i)
	}
}

func TestCompilesChecker(t *testing.T) {
	cli, err := products.Bin("okgo")
	require.NoError(t, err)

	wd, err := os.Getwd()
	require.NoError(t, err)

	tmpDir, cleanup, err := dirs.TempDir(wd, "")
	defer cleanup()
	require.NoError(t, err)

	for i, currCase := range []struct {
		check         amalgomated.Cmd
		filesToWrite  []gofiles.GoFileSpec
		pathToCheck   func(projectDir string) string
		want          func(files map[string]gofiles.GoFile) []string
		customMatcher func(caseNum int, expected, actual []string)
	}{
		{
			check: cmdlib.Instance().MustNewCmd("compiles"),
			filesToWrite: []gofiles.GoFileSpec{
				{
					RelPath: "foo/foo.go",
					Src: `package foo
func Foo() int {
	return "foo"
}`,
				},
			},
			pathToCheck: func(projectDir string) string {
				return path.Join(projectDir, "foo")
			},
			want: func(files map[string]gofiles.GoFile) []string {
				return []string{
					`foo.go:3:9: cannot convert "foo" (untyped string constant) to int`,
				}
			},
		},
		{
			check: cmdlib.Instance().MustNewCmd("compiles"),
			filesToWrite: []gofiles.GoFileSpec{
				{
					RelPath: "foo/foo.go",
					Src: `package foo
import "bar"
func Foo() {
	bar.Bar()
}`,
				},
			},
			pathToCheck: func(projectDir string) string {
				return path.Join(projectDir, "foo")
			},
			want: func(files map[string]gofiles.GoFile) []string {
				return []string{
					`foo.go:2:8: could not import bar \(cannot find package "bar" in any of:
.+ \(vendor tree\)
.+
.+ \(from \$GOROOT\)
.+ \(from \$GOPATH\)\)`,
				}
			},
			customMatcher: func(caseNum int, want, got []string) {
				ok := assert.Equal(t, len(want), len(got), "Case %d: number of output lines do not match", caseNum)
				if ok {
					for i := range want {
						assert.Regexp(t, want[i], got[i], "Case %d, want case %d", caseNum, i)
					}
				}
			},
		},
		{
			check: cmdlib.Instance().MustNewCmd("compiles"),
			filesToWrite: []gofiles.GoFileSpec{
				{
					RelPath: "foo/foo.go",
					Src: `package foo
func Foo() {
	bar.Bar()
	baz.Baz()
}`,
				},
			},
			pathToCheck: func(projectDir string) string {
				return path.Join(projectDir, "foo")
			},
			want: func(files map[string]gofiles.GoFile) []string {
				return []string{
					`foo.go:3:2: undeclared name: bar`,
					`foo.go:4:2: undeclared name: baz`,
				}
			},
			customMatcher: func(caseNum int, want, got []string) {
				ok := assert.Equal(t, len(want), len(got), "Case %d: number of output lines do not match", caseNum)
				if ok {
					for i := range want {
						assert.Regexp(t, want[i], got[i], "Case %d, want case %d", caseNum, i)
					}
				}
			},
		},
		{
			check: cmdlib.Instance().MustNewCmd("extimport"),
			filesToWrite: []gofiles.GoFileSpec{
				{
					RelPath: "foo/foo.go",
					Src: `package foo
import "{{index . "bar/bar.go"}}"
func Foo() {
	bar.Bar()
}
`,
				},
				{
					RelPath: "bar/bar.go",
					Src:     `package bar; func Bar() {}`,
				},
			},
			pathToCheck: func(projectDir string) string {
				return path.Join(projectDir, "foo")
			},
			want: func(files map[string]gofiles.GoFile) []string {
				return []string{
					fmt.Sprintf(`foo.go:2:8: imports external package %s`, files["bar/bar.go"].ImportPath),
				}
			},
		},
	} {
		currCaseProjectDir, err := ioutil.TempDir(tmpDir, "")
		require.NoError(t, err, "Case %d", i)

		files, err := gofiles.Write(currCaseProjectDir, currCase.filesToWrite)
		require.NoError(t, err, "Case %d", i)

		checker, err := checks.GetChecker(currCase.check)
		require.NoError(t, err, "Case %d", i)

		runner := amalgomated.PathCmder(cli, amalgomated.ProxyCmdPrefix+currCase.check.Name())
		lineInfo, err := checker.Check(runner, currCase.pathToCheck(currCaseProjectDir), params.OKGo{})
		require.NoError(t, err, "Case %d", i)

		want := currCase.want(files)
		got := toStringSlice(lineInfo)
		if currCase.customMatcher == nil {
			assert.Equal(t, want, got, "Case %d", i)
		} else {
			currCase.customMatcher(i, want, got)
		}
	}
}

func TestFilters(t *testing.T) {
	cli, err := products.Bin("okgo")
	require.NoError(t, err)

	cmd := cmdlib.Instance().MustNewCmd("golint")
	checker, err := checks.GetChecker(cmd)
	require.NoError(t, err)
	runner := amalgomated.PathCmder(cli, amalgomated.ProxyCmdPrefix+cmd.Name())

	for i, currCase := range []struct {
		filters []checkoutput.Filterer
		want    []string
	}{
		{
			filters: nil,
			want: []string{
				"bad.go:3:1: exported function Bad should have comment or be unexported",
				"mock/mock.go:3:1: exported function Mock should have comment or be unexported",
				"nested/mock/nestedmock.go:3:1: exported function NestedMock should have comment or be unexported",
			},
		},
		{
			filters: []checkoutput.Filterer{
				checkoutput.RelativePathFilter("mock"),
			},
			want: []string{
				"bad.go:3:1: exported function Bad should have comment or be unexported",
				"nested/mock/nestedmock.go:3:1: exported function NestedMock should have comment or be unexported",
			},
		},
		{
			filters: []checkoutput.Filterer{
				checkoutput.NamePathFilter("mock"),
			},
			want: []string{
				"bad.go:3:1: exported function Bad should have comment or be unexported",
			},
		},
		{
			filters: []checkoutput.Filterer{
				checkoutput.MessageRegexpFilter("should have comment or be unexported"),
			},
			want: []string{},
		},
	} {
		lineInfo, err := checker.Check(runner, "./testdata/filter", params.OKGo{})
		require.NoError(t, err, "Case %d", i)

		filteredLines, err := checkoutput.ApplyFilters(lineInfo, currCase.filters)
		require.NoError(t, err, "Case %d", i)

		assert.Equal(t, currCase.want, toStringSlice(filteredLines), "Case %d", i)
	}
}

func TestCheckerUsesConfig(t *testing.T) {
	cli, err := products.Bin("okgo")
	require.NoError(t, err)

	tmpDir, cleanup, err := dirs.TempDir("", "")
	defer cleanup()
	require.NoError(t, err)

	for i, currCase := range []struct {
		config string
		want   []string
	}{
		{
			config: "",
			want: []string{
				"bad.go:3:1: exported function Bad should have comment or be unexported",
				"mock/mock.go:3:1: exported function Mock should have comment or be unexported",
				"nested/mock/nestedmock.go:3:1: exported function NestedMock should have comment or be unexported",
			},
		},
		{
			config: `
			exclude:
			  paths:
			    - "mock"
			`,
			want: []string{
				"bad.go:3:1: exported function Bad should have comment or be unexported",
				"nested/mock/nestedmock.go:3:1: exported function NestedMock should have comment or be unexported",
			},
		},
		{
			config: `
			exclude:
			  names:
			    - "m.ck"
			`,
			want: []string{
				"bad.go:3:1: exported function Bad should have comment or be unexported",
			},
		},
		{
			config: `
			checks:
			  golint:
			    filters:
			      - type: "message"
			        value: "should have comment or be unexported"
			`,
			want: []string{},
		},
	} {
		tmpFile, err := ioutil.TempFile(tmpDir, "")
		require.NoError(t, err, "Case %d", i)
		tmpFilePath := tmpFile.Name()
		err = tmpFile.Close()
		require.NoError(t, err, "Case %d", i)
		err = ioutil.WriteFile(tmpFilePath, []byte(unindent(currCase.config)), 0644)
		require.NoError(t, err, "Case %d", i)

		cfg, err := config.Load(tmpFilePath, "")
		require.NoError(t, err, "Case %d", i)

		cmd := cmdlib.Instance().MustNewCmd("golint")
		checker, err := checks.GetChecker(cmd)
		require.NoError(t, err, "Case %d", i)

		runner := amalgomated.PathCmder(cli, amalgomated.ProxyCmdPrefix+cmd.Name())
		lineInfo, err := checker.Check(runner, "./testdata/filter", cfg)
		require.NoError(t, err, "Case %d", i)

		assert.Equal(t, currCase.want, toStringSlice(lineInfo), "Case %d", i)
	}
}

func TestCheckerUsesReleaseTagConfig(t *testing.T) {
	cli, err := products.Bin("okgo")
	require.NoError(t, err)

	tmpDir, cleanup, err := dirs.TempDir("", "")
	defer cleanup()
	require.NoError(t, err)

	for i, currCase := range []struct {
		name   string
		files  []gofiles.GoFileSpec
		config string
		want   []string
	}{
		{
			name: "file with go1.7 build tag processed",
			files: []gofiles.GoFileSpec{
				{
					RelPath: "foo.go",
					Src: `// +build go1.7

					package foo
					import "os"
					func Foo() {
						os.Setenv("foo", "bar")
					}
					`,
				},
				{
					RelPath: "bar.go",
					Src: `package foo
					import "os"
					func Bar() {
						os.Setenv("foo", "bar")
					}
					`,
				},
			},
			config: "",
			want: []string{
				"Running errcheck...",
				"bar.go:4:16: os.Setenv(\"foo\", \"bar\")",
				"foo.go:6:16: os.Setenv(\"foo\", \"bar\")",
				"",
			},
		},
		{
			name: "file with go1.7 build tag ignored if release-tag set to go1.6",
			files: []gofiles.GoFileSpec{
				{
					RelPath: "foo.go",
					Src: `// +build go1.7

					package foo
					import "os"
					func Foo() {
						os.Setenv("foo", "bar")
					}
					`,
				},
				{
					RelPath: "bar.go",
					Src: `package foo
					import "os"
					func Bar() {
						os.Setenv("foo", "bar")
					}
					`,
				},
			},
			config: `release-tag: go1.6`,
			want: []string{
				"Running errcheck...",
				"bar.go:4:16: os.Setenv(\"foo\", \"bar\")",
				"",
			},
		},
		{
			name: "ignoring a returned error is flagged",
			files: []gofiles.GoFileSpec{
				{
					RelPath: "foo.go",
					Src: `package foo
					import "os"
					func Foo() {
						os.Open("/")
						os.Pipe()
					}
					`,
				},
			},
			config: ``,
			want: []string{
				"Running errcheck...",
				"foo.go:4:14: os.Open(\"/\")",
				"foo.go:5:14: os.Pipe()",
				"",
			},
		},
		{
			name: "ignoring a returned error is not flagged if referenced from an exclude list",
			files: []gofiles.GoFileSpec{
				{
					RelPath: "foo.go",
					Src: `package foo
					import "os"
					func Foo() {
						os.Open("/")
						os.Pipe()
					}
					`,
				},
				{
					RelPath: "exclude.txt",
					Src: `os.Open
					`,
				},
			},
			config: `
			checks:
			  errcheck:
			    args:
			      - "-exclude"
			      - "exclude.txt"
			`,
			want: []string{
				"Running errcheck...",
				"foo.go:5:14: os.Pipe()",
				"",
			},
		},
	} {
		currCaseDir, err := ioutil.TempDir(tmpDir, "")
		require.NoError(t, err, "Case %d: %s", i, currCase.name)

		_, err = gofiles.Write(currCaseDir, currCase.files)
		require.NoError(t, err, "Case %d: %s", i, currCase.name)

		tmpFile, err := ioutil.TempFile(currCaseDir, "")
		require.NoError(t, err, "Case %d: %s", i, currCase.name)
		cfgFilePath := tmpFile.Name()
		err = tmpFile.Close()
		require.NoError(t, err, "Case %d: %s", i, currCase.name)
		err = ioutil.WriteFile(cfgFilePath, []byte(unindent(currCase.config)), 0644)
		require.NoError(t, err, "Case %d: %s", i, currCase.name)

		cmd := exec.Command(cli, "--config", cfgFilePath, "errcheck", ".")
		cmd.Dir = currCaseDir
		output, err := cmd.CombinedOutput()
		require.Error(t, err, fmt.Errorf("Expected command %v to fail. Output:\n%v", cmd.Args, string(output)))

		assert.Equal(t, currCase.want, strings.Split(string(output), "\n"), "Case %d: %s", i, currCase.name)
	}
}

func toStringSlice(input []checkoutput.Issue) []string {
	output := make([]string, len(input))
	for i, curr := range input {
		output[i] = curr.String()
	}
	return output
}

func unindent(input string) string {
	return strings.Replace(input, "\n\t\t\t", "\n", -1)
}
