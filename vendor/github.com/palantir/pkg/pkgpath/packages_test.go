// Copyright 2016 Palantir Technologies. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package pkgpath_test

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"testing"

	"github.com/nmiyake/pkg/dirs"
	"github.com/nmiyake/pkg/gofiles"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/palantir/pkg/matcher"
	"github.com/palantir/pkg/pkgpath"
)

func TestNewPackages(t *testing.T) {
	wd, err := os.Getwd()
	require.NoError(t, err)

	tmpDir, cleanup, err := dirs.TempDir(wd, "")
	defer cleanup()
	require.NoError(t, err)

	for i, currCase := range []struct {
		filesToCreate []gofiles.GoFileSpec
		args          []string
		want          map[string]string
	}{
		{
			filesToCreate: []gofiles.GoFileSpec{
				{RelPath: "foo/main.go", Src: "package main"},
				{RelPath: "foo/bar/bar.go", Src: "package bar"},
				{RelPath: "baz.go", Src: "package baz"},
			},
			args: []string{
				"./...",
			},
			want: map[string]string{
				"./foo":     "main",
				"./foo/bar": "bar",
				"./.":       "baz",
			},
		},
		{
			filesToCreate: []gofiles.GoFileSpec{
				{RelPath: "foo/main.go", Src: "package main"},
				{RelPath: "foo/bar/bar.go", Src: "package bar"},
				{RelPath: "baz.go", Src: "package baz"},
			},
			args: []string{
				"./foo/...",
			},
			want: map[string]string{
				"./foo":     "main",
				"./foo/bar": "bar",
			},
		},
	} {
		currCaseDir, err := ioutil.TempDir(tmpDir, "")
		require.NoError(t, err)

		_, err = gofiles.Write(currCaseDir, currCase.filesToCreate)
		require.NoError(t, err)

		pkgs, err := pkgpath.PackagesFromPaths(currCaseDir, currCase.args)
		require.NoError(t, err, "Case %d", i)

		got, err := pkgs.Packages(pkgpath.Relative)
		require.NoError(t, err, "Case %d", i)

		assert.Equal(t, currCase.want, got)
	}
}

func TestListPackages(t *testing.T) {
	wd, err := os.Getwd()
	require.NoError(t, err)

	testPkgPath, err := filepath.Rel(path.Join(os.Getenv("GOPATH"), "src"), wd)
	require.NoError(t, err)

	tmpDir, cleanup, err := dirs.TempDir(wd, "")
	defer cleanup()
	require.NoError(t, err)

	for i, currCase := range []struct {
		filesToCreate []gofiles.GoFileSpec
		exclude       matcher.Matcher
		want          map[string]string
	}{
		{
			filesToCreate: []gofiles.GoFileSpec{
				{RelPath: "matchers.go", Src: "package matchers"},
			},
			want: map[string]string{
				".": "matchers",
			},
		},
		{
			filesToCreate: []gofiles.GoFileSpec{
				{RelPath: "main.go", Src: "package main"},
			},
			want: map[string]string{
				".": "main",
			},
		},
		// build tags are taken into consideration
		{
			filesToCreate: []gofiles.GoFileSpec{
				{RelPath: "main.go", Src: "package main"},
				{RelPath: "no_build.go", Src: `// +build ignore

package different`},
			},
			want: map[string]string{
				".": "main",
			},
		},
		{
			filesToCreate: []gofiles.GoFileSpec{
				{RelPath: "foo/main.go", Src: "package main"},
				{RelPath: "nosource/notgo.txt", Src: "package notgo"},
				{RelPath: "pkg/barpkg.go", Src: "package bar"},
				{RelPath: "regular/regular.go", Src: "package regular"},
				{RelPath: ".hidden/hidden.go", Src: "package hidden"},
			},
			want: map[string]string{
				".hidden": "hidden",
				"foo":     "main",
				"pkg":     "bar",
				"regular": "regular",
			},
		},
		{
			filesToCreate: []gofiles.GoFileSpec{
				{RelPath: "foo/main.go", Src: "package main"},
				{RelPath: "nosource/notgo.txt", Src: "package notgo"},
				{RelPath: "pkg/barpkg.go", Src: "package bar"},
				{RelPath: "regular/regular.go", Src: "package regular"},
				{RelPath: ".hidden/hidden.go", Src: "package hidden"},
			},
			exclude: matcher.Hidden(),
			want: map[string]string{
				"foo":     "main",
				"pkg":     "bar",
				"regular": "regular",
			},
		},
		{
			filesToCreate: []gofiles.GoFileSpec{
				{RelPath: "foo/main.go", Src: "package main"},
				{RelPath: "pkg/.barpkg.go", Src: "package bar"},
				{RelPath: ".hidden/hidden.go", Src: "package hidden"},
			},
			exclude: matcher.Hidden(),
			want: map[string]string{
				"foo": "main",
			},
		},
		{
			filesToCreate: []gofiles.GoFileSpec{
				{RelPath: "foo/main.go", Src: "package main"},
				{RelPath: "foo/main_test.go", Src: "package main_test"},
			},
			want: map[string]string{
				"foo": "main",
			},
		},
		{
			filesToCreate: []gofiles.GoFileSpec{
				{RelPath: "foo/main.go", Src: "package main"},
				{RelPath: "bar/bar.go", Src: "package bar"},
				{RelPath: "foo/integration_tests/main_test.go", Src: "package main_test"},
			},
			exclude: matcher.Not(matcher.Name("integration_tests")),
			want: map[string]string{
				"foo/integration_tests": "main_test",
			},
		},
		{
			filesToCreate: []gofiles.GoFileSpec{
				{RelPath: "foo/main.go", Src: "package main"},
				{RelPath: "_bar/bar.go", Src: "package bar"},
				{RelPath: "baz/_baz.go", Src: "package baz"},
				{RelPath: ".hiddendir/hiddendir.go", Src: "package hiddendir"},
				{RelPath: "hidden/.filehidden.go", Src: "package filehidden"},
				{RelPath: "testdata/tester.go", Src: "package tester"},
				{RelPath: "nottestdata/testdata.go", Src: "package testdata"},
			},
			exclude: pkgpath.DefaultGoPkgExcludeMatcher(),
			want: map[string]string{
				"foo":         "main",
				"nottestdata": "testdata",
			},
		},
	} {
		currCaseDir, err := ioutil.TempDir(tmpDir, "")
		require.NoError(t, err, "Case %d", i)

		currCaseDirRelPath, err := filepath.Rel(wd, currCaseDir)
		require.NoError(t, err, "Case %d", i)

		_, err = gofiles.Write(currCaseDir, currCase.filesToCreate)
		require.NoError(t, err, "Case %d", i)

		pkgs, err := pkgpath.PackagesInDir(currCaseDir, currCase.exclude)
		require.NoError(t, err, "Case %d", i)

		for _, mode := range []pkgpath.Type{
			pkgpath.Relative,
			pkgpath.GoPathSrcRelative,
			pkgpath.Absolute,
		} {
			want := make(map[string]string)
			for k, v := range currCase.want {
				switch mode {
				case pkgpath.Relative:
					k = "./" + k
				case pkgpath.GoPathSrcRelative:
					k = path.Join(testPkgPath, currCaseDirRelPath, k)
				case pkgpath.Absolute:
					k = path.Join(wd, currCaseDirRelPath, k)
				default:
					require.Fail(t, "Unhandled case: %v", mode)
				}
				want[k] = v
			}

			got, err := pkgs.Packages(mode)
			require.NoError(t, err, "Case %d", i)
			assert.Equal(t, want, got, "Case %d, mode %v", i, mode)
		}
	}
}

func TestListPackagesFailsWithMultiplePackages(t *testing.T) {
	wd, err := os.Getwd()
	require.NoError(t, err)

	tmpDir, cleanup, err := dirs.TempDir(wd, "")
	defer cleanup()
	require.NoError(t, err)

	for i, currCase := range []struct {
		filesToCreate []gofiles.GoFileSpec
		errorMessage  string
	}{
		{
			filesToCreate: []gofiles.GoFileSpec{
				{RelPath: "foo/main.go", Src: "package main"},
				{RelPath: "foo/foo.go", Src: "package foo"},
			},
			errorMessage: `.+contains more than 1 package: \[foo main\]`,
		},
		{
			filesToCreate: []gofiles.GoFileSpec{
				{RelPath: "foo/main.go", Src: "package main"},
				{RelPath: "foo/foo_test.go", Src: "package foo_test"},
			},
			errorMessage: `.+contains more than 1 package: \[foo main\]`,
		},
	} {
		currCaseDir, err := ioutil.TempDir(tmpDir, "")
		require.NoError(t, err, "Case %d", i)

		_, err = gofiles.Write(currCaseDir, currCase.filesToCreate)
		require.NoError(t, err, "Case %d", i)

		_, err = pkgpath.PackagesInDir(currCaseDir, nil)
		require.Error(t, err, fmt.Sprintf("Case %d", i))

		assert.Regexp(t, regexp.MustCompile(currCase.errorMessage), err.Error())
	}
}

// Verify that ListPackages uses the current value of the GOPATH environment variable to determine the package paths.
func TestListPackagesSetGoPath(t *testing.T) {
	tmpDir, cleanup, err := dirs.TempDir("", "")
	defer cleanup()
	require.NoError(t, err)

	// get original value of GOPATH and restore after test
	origGoPath := os.Getenv("GOPATH")
	defer func() {
		if err := os.Setenv("GOPATH", origGoPath); err != nil {
			fmt.Println("Failed to restore value of GOPATH enviornment variable:", err)
		}
	}()

	err = os.Setenv("GOPATH", tmpDir)
	require.NoError(t, err)

	projectDir := path.Join(tmpDir, "src", "github.com", "test")
	err = os.MkdirAll(projectDir, 0755)
	require.NoError(t, err)

	goFiles := []gofiles.GoFileSpec{
		{RelPath: "foo/main.go", Src: "package main"},
		{RelPath: "foo/main_test.go", Src: "package main_test"},
	}
	_, err = gofiles.Write(projectDir, goFiles)
	require.NoError(t, err)

	pkgs, err := pkgpath.PackagesInDir(projectDir, nil)
	require.NoError(t, err)

	got, err := pkgs.Packages(pkgpath.GoPathSrcRelative)
	require.NoError(t, err)
	want := map[string]string{
		"github.com/test/foo": "main",
	}

	assert.Equal(t, want, got)
}

func TestPkgPathOutsideGoPathFails(t *testing.T) {
	goPathSrc := path.Join(os.Getenv("GOPATH"), "src")
	msg := fmt.Sprintf(`^resolving /foo against base %s produced relative path starting with ../: .+/foo$`, goPathSrc)

	_, err := pkgpath.NewAbsPkgPath("/foo").GoPathSrcRel()
	require.Regexp(t, msg, err.Error())
}
