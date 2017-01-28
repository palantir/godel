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
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"github.com/nmiyake/pkg/dirs"
	"github.com/nmiyake/pkg/gofiles"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/palantir/godel/apps/distgo/pkg/git"
	"github.com/palantir/godel/apps/distgo/pkg/git/gittest"
	"github.com/palantir/godel/pkg/products"
)

var (
	gödelTGZ    string
	testRootDir string
	version     string
)

func TestMain(m *testing.M) {
	os.Exit(runTestMain(m))
}

func runTestMain(m *testing.M) int {
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	gödelProjectDir := path.Join(wd, "..")
	version, err = git.ProjectVersion(gödelProjectDir)
	if err != nil {
		panic(fmt.Sprintf("Failed to get version from directory %s: %v", gödelProjectDir, err))
	}

	gödelTGZ, err = products.Dist("godel")
	if err != nil {
		panic(fmt.Sprintf("Failed create distribution: %v", err))
	}

	var cleanup func()
	testRootDir, cleanup, err = dirs.TempDir(wd, "")
	defer cleanup()
	if err != nil {
		panic(fmt.Sprintf("Failed to create temporary directory in %s: %v", wd, err))
	}

	return m.Run()
}

func TestVersion(t *testing.T) {
	testProjectDir := setUpGödelTestAndDownload(t, testRootDir, gödelTGZ, version)

	output := execCommand(t, testProjectDir, "./godelw", "--version")
	assert.Equal(t, fmt.Sprintf("godel version %v\n", version), string(output))
}

func TestGitHooksSuccess(t *testing.T) {
	// create project directory in temporary location so primary project's repository is not modified by test
	tmp, cleanup, err := dirs.TempDir("", "")
	defer cleanup()
	require.NoError(t, err)

	testProjectDir := setUpGödelTestAndDownload(t, tmp, gödelTGZ, version)

	// initialize git repository
	gittest.InitGitDir(t, testProjectDir)

	// install commit hooks
	execCommand(t, testProjectDir, "./godelw", "git-hooks")

	// committing Go file that is properly formatted works
	formatted := `package main

func main() {
}
`
	err = ioutil.WriteFile(path.Join(testProjectDir, "main.go"), []byte(formatted), 0644)
	require.NoError(t, err)
	execCommand(t, testProjectDir, "git", "add", ".")
	execCommand(t, testProjectDir, "git", "commit", "--author=testAuthor <test@author.com>", "-m", "Second commit")
}

func TestGitHooksFail(t *testing.T) {
	// create project directory in temporary location so primary project's repository is not modified by test
	tmp, cleanup, err := dirs.TempDir("", "")
	defer cleanup()
	require.NoError(t, err)

	testProjectDir := setUpGödelTestAndDownload(t, tmp, gödelTGZ, version)

	// initialize git repository
	gittest.InitGitDir(t, testProjectDir)

	// install commit hooks
	execCommand(t, testProjectDir, "./godelw", "git-hooks")

	// committing Go file that is not properly formatted causes error
	notFormatted := `package main
import "fmt"

func Foo() {
fmt.Println("foo")
}`
	err = ioutil.WriteFile(path.Join(testProjectDir, "helper.go"), []byte(notFormatted), 0644)
	require.NoError(t, err)
	execCommand(t, testProjectDir, "git", "add", ".")

	cmd := exec.Command("git", "commit", "--author=testAuthor <test@author.com>", "-m", "Second commit")
	cmd.Dir = testProjectDir
	output, err := cmd.CombinedOutput()
	assert.Error(t, err, "exit status 1")
	assert.Regexp(t, `(?s)^Unformatted files exist -- run ./godelw format to format these files:\n  .+/helper.go\n$`, string(output))
}

func TestFormat(t *testing.T) {
	testProjectDir := setUpGödelTestAndDownload(t, testRootDir, gödelTGZ, version)

	src := `package main
		import "fmt"

	func main() {
	fmt.Println("hello, world!")
	}`

	formattedSrc := `package main

import (
	"fmt"
)

func main() {
	fmt.Println("hello, world!")
}
`
	err := ioutil.WriteFile(path.Join(testProjectDir, "main.go"), []byte(src), 0644)
	require.NoError(t, err)

	execCommand(t, testProjectDir, "./godelw", "format")

	content, err := ioutil.ReadFile(path.Join(testProjectDir, "main.go"))
	require.NoError(t, err)
	assert.Equal(t, formattedSrc, string(content))
}

func TestImports(t *testing.T) {
	testProjectDir := setUpGödelTestAndDownload(t, testRootDir, gödelTGZ, version)

	const importsYML = `root-dirs:
  - pkg`
	err := ioutil.WriteFile(path.Join(testProjectDir, "godel", "config", "imports.yml"), []byte(importsYML), 0644)
	require.NoError(t, err)

	specs := []gofiles.GoFileSpec{
		{
			RelPath: "pkg/foo/foo.go",
			Src:     `package foo; import _ "{{index . "bar.go"}}";`,
		},
		{
			RelPath: "bar.go",
			Src:     "package bar",
		},
	}

	files, err := gofiles.Write(testProjectDir, specs)
	require.NoError(t, err)

	want := fmt.Sprintf(`{
    "imports": [
        {
            "path": "%s",
            "numGoFiles": 1,
            "numImportedGoFiles": 0,
            "importedFrom": [
                "%s"
            ]
        }
    ],
    "mainOnlyImports": [],
    "testOnlyImports": []
}`, files["bar.go"].ImportPath, files["pkg/foo/foo.go"].ImportPath)

	execCommand(t, testProjectDir, "./godelw", "imports")

	content, err := ioutil.ReadFile(path.Join(testProjectDir, "pkg", "gocd_imports.json"))
	require.NoError(t, err)
	assert.Equal(t, want, string(content))
}

func TestImportsVerify(t *testing.T) {
	testProjectDir := setUpGödelTestAndDownload(t, testRootDir, gödelTGZ, version)

	const importsYML = `root-dirs:
  - pkg`
	err := ioutil.WriteFile(path.Join(testProjectDir, "godel", "config", "imports.yml"), []byte(importsYML), 0644)
	require.NoError(t, err)

	specs := []gofiles.GoFileSpec{
		{
			RelPath: "pkg/foo/foo.go",
			Src:     `package foo; import _ "{{index . "bar.go"}}";`,
		},
		{
			RelPath: "bar.go",
			Src:     "package bar",
		},
	}

	_, err = gofiles.Write(testProjectDir, specs)
	require.NoError(t, err)

	cmd := exec.Command("./godelw", "imports", "--verify")
	cmd.Dir = testProjectDir
	output, err := cmd.CombinedOutput()
	require.Error(t, err)
	assert.Equal(t, "gocd_imports.json out of date for 1 directory:\n\tpkg: gocd_imports.json does not exist\n", string(output))
}

func TestLicense(t *testing.T) {
	testProjectDir := setUpGödelTestAndDownload(t, testRootDir, gödelTGZ, version)

	const licenseYML = `header: |
  /*
  Copyright 2016 Palantir Technologies, Inc.

  Licensed under the Apache License, Version 2.0 (the "License");
  you may not use this file except in compliance with the License.
  You may obtain a copy of the License at

  http://www.apache.org/licenses/LICENSE-2.0

  Unless required by applicable law or agreed to in writing, software
  distributed under the License is distributed on an "AS IS" BASIS,
  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
  See the License for the specific language governing permissions and
  limitations under the License.
  */
`
	err := ioutil.WriteFile(path.Join(testProjectDir, "godel", "config", "license.yml"), []byte(licenseYML), 0644)
	require.NoError(t, err)

	specs := []gofiles.GoFileSpec{
		{
			RelPath: "foo.go",
			Src:     "package foo",
		},
		{
			RelPath: "vendor/github.com/bar.go",
			Src:     "package bar",
		},
	}

	files, err := gofiles.Write(testProjectDir, specs)
	require.NoError(t, err)

	want := `/*
Copyright 2016 Palantir Technologies, Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package foo`

	execCommand(t, testProjectDir, "./godelw", "license")

	content, err := ioutil.ReadFile(files["foo.go"].Path)
	require.NoError(t, err)
	assert.Equal(t, want, string(content))

	want = `package bar`
	content, err = ioutil.ReadFile(files["vendor/github.com/bar.go"].Path)
	require.NoError(t, err)
	assert.Equal(t, want, string(content))
}

func TestLicenseVerify(t *testing.T) {
	testProjectDir := setUpGödelTestAndDownload(t, testRootDir, gödelTGZ, version)

	const licenseYML = `header: |
  /*
  Copyright 2016 Palantir Technologies, Inc.

  Licensed under the Apache License, Version 2.0 (the "License");
  you may not use this file except in compliance with the License.
  You may obtain a copy of the License at

  http://www.apache.org/licenses/LICENSE-2.0

  Unless required by applicable law or agreed to in writing, software
  distributed under the License is distributed on an "AS IS" BASIS,
  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
  See the License for the specific language governing permissions and
  limitations under the License.
  */
`
	err := ioutil.WriteFile(path.Join(testProjectDir, "godel", "config", "license.yml"), []byte(licenseYML), 0644)
	require.NoError(t, err)

	specs := []gofiles.GoFileSpec{
		{
			RelPath: "foo.go",
			Src:     "package foo",
		},
		{
			RelPath: "bar/bar.go",
			Src: `/*
Copyright 2016 Palantir Technologies, Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package bar`,
		},
		{
			RelPath: "vendor/github.com/baz.go",
			Src:     "package baz",
		},
	}

	_, err = gofiles.Write(testProjectDir, specs)
	require.NoError(t, err)

	cmd := exec.Command("./godelw", "license", "--verify")
	cmd.Dir = testProjectDir
	output, err := cmd.CombinedOutput()
	require.Error(t, err)
	assert.Equal(t, "1 file does not have the correct license header:\n\tfoo.go\n", string(output))
}

func TestCheck(t *testing.T) {
	testProjectDir := setUpGödelTestAndDownload(t, testRootDir, gödelTGZ, version)
	src := `package main
	import "fmt"

	func main() {
		fmt.Println("hello, world!")
	}`
	err := ioutil.WriteFile(path.Join(testProjectDir, "main.go"), []byte(src), 0644)
	require.NoError(t, err)

	execCommand(t, testProjectDir, "./godelw", "check")
}

func TestProducts(t *testing.T) {
	testProjectDir := setUpGödelTestAndDownload(t, testRootDir, gödelTGZ, version)

	distYml := `
products:
  foo:
    build:
      main-pkg: ./foo
  bar:
    build:
      main-pkg: ./bar
`
	err := ioutil.WriteFile(path.Join(testProjectDir, "godel", "config", "dist.yml"), []byte(distYml), 0644)
	require.NoError(t, err)

	src := `package main
	import "fmt"

	func main() {
		fmt.Println("hello, world!")
	}`
	err = os.MkdirAll(path.Join(testProjectDir, "foo"), 0755)
	require.NoError(t, err)
	err = ioutil.WriteFile(path.Join(testProjectDir, "foo", "foo.go"), []byte(src), 0644)
	require.NoError(t, err)

	err = os.MkdirAll(path.Join(testProjectDir, "bar"), 0755)
	require.NoError(t, err)
	err = ioutil.WriteFile(path.Join(testProjectDir, "bar", "bar.go"), []byte(src), 0644)
	require.NoError(t, err)

	execCommand(t, testProjectDir, "./godelw", "products")
}

func TestArtifactsBuild(t *testing.T) {
	testProjectDir := setUpGödelTestAndDownload(t, testRootDir, gödelTGZ, version)
	gittest.InitGitDir(t, testProjectDir)

	distYml := `
products:
  foo:
    build:
      main-pkg: ./foo
      os-archs:
        - os: darwin
          arch: amd64
        - os: linux
          arch: amd64
  bar:
    build:
      main-pkg: ./bar
      os-archs:
        - os: windows
          arch: amd64
`
	err := ioutil.WriteFile(path.Join(testProjectDir, "godel", "config", "dist.yml"), []byte(distYml), 0644)
	require.NoError(t, err)

	src := `package main
	import "fmt"

	func main() {
		fmt.Println("hello, world!")
	}`
	err = os.MkdirAll(path.Join(testProjectDir, "foo"), 0755)
	require.NoError(t, err)
	err = ioutil.WriteFile(path.Join(testProjectDir, "foo", "foo.go"), []byte(src), 0644)
	require.NoError(t, err)

	err = os.MkdirAll(path.Join(testProjectDir, "bar"), 0755)
	require.NoError(t, err)
	err = ioutil.WriteFile(path.Join(testProjectDir, "bar", "bar.go"), []byte(src), 0644)
	require.NoError(t, err)

	gittest.CommitAllFiles(t, testProjectDir, "Commit files")
	gittest.CreateGitTag(t, testProjectDir, "0.1.0")

	output := execCommand(t, testProjectDir, "./godelw", "artifacts", "build")

	want := `build/0.1.0/windows-amd64/bar.exe
build/0.1.0/darwin-amd64/foo
build/0.1.0/linux-amd64/foo
`
	assert.Equal(t, want, output)
}

func TestArtifactsDist(t *testing.T) {
	testProjectDir := setUpGödelTestAndDownload(t, testRootDir, gödelTGZ, version)
	gittest.InitGitDir(t, testProjectDir)

	distYml := `
products:
  foo:
    build:
      main-pkg: ./foo
  bar:
    build:
      main-pkg: ./bar
`
	err := ioutil.WriteFile(path.Join(testProjectDir, "godel", "config", "dist.yml"), []byte(distYml), 0644)
	require.NoError(t, err)

	src := `package main
	import "fmt"

	func main() {
		fmt.Println("hello, world!")
	}`
	err = os.MkdirAll(path.Join(testProjectDir, "foo"), 0755)
	require.NoError(t, err)
	err = ioutil.WriteFile(path.Join(testProjectDir, "foo", "foo.go"), []byte(src), 0644)
	require.NoError(t, err)

	err = os.MkdirAll(path.Join(testProjectDir, "bar"), 0755)
	require.NoError(t, err)
	err = ioutil.WriteFile(path.Join(testProjectDir, "bar", "bar.go"), []byte(src), 0644)
	require.NoError(t, err)

	gittest.CommitAllFiles(t, testProjectDir, "Commit files")
	gittest.CreateGitTag(t, testProjectDir, "0.1.0")

	output := execCommand(t, testProjectDir, "./godelw", "artifacts", "dist")
	assert.Equal(t, "dist/bar-0.1.0.sls.tgz\ndist/foo-0.1.0.sls.tgz\n", output)
}

func TestArtifactsBuildClean(t *testing.T) {
	testProjectDir := setUpGödelTestAndDownload(t, testRootDir, gödelTGZ, version)
	gittest.InitGitDir(t, testProjectDir)

	distYml := `
products:
  foo:
    build:
      output-dir: ./foo/build
      main-pkg: ./foo
      os-archs:
        - os: darwin
          arch: amd64
        - os: linux
          arch: amd64
  bar:
    build:
      main-pkg: ./bar
      os-archs:
        - os: linux
          arch: amd64
`
	err := ioutil.WriteFile(path.Join(testProjectDir, "godel", "config", "dist.yml"), []byte(distYml), 0644)
	require.NoError(t, err)

	src := `package main
	import "fmt"

	func main() {
		fmt.Println("hello, world!")
	}`
	err = os.MkdirAll(path.Join(testProjectDir, "foo"), 0755)
	require.NoError(t, err)
	err = ioutil.WriteFile(path.Join(testProjectDir, "foo", "foo.go"), []byte(src), 0644)
	require.NoError(t, err)
	err = os.MkdirAll(path.Join(testProjectDir, "foo", "build"), 0755)
	require.NoError(t, err)
	err = ioutil.WriteFile(path.Join(testProjectDir, "foo", "build", "test-file"), []byte("test"), 0644)
	require.NoError(t, err)

	err = os.MkdirAll(path.Join(testProjectDir, "bar"), 0755)
	require.NoError(t, err)
	err = ioutil.WriteFile(path.Join(testProjectDir, "bar", "bar.go"), []byte(src), 0644)
	require.NoError(t, err)

	gittest.CommitAllFiles(t, testProjectDir, "Commit files")
	gittest.CreateGitTag(t, testProjectDir, "0.1.0")

	execCommand(t, testProjectDir, "./godelw", "build")
	fooBuildDir := path.Join(testProjectDir, "foo")
	barBuildDir := path.Join(testProjectDir)

	// Builds exist for both products
	fooBuildFiles, err := listRecursive(fooBuildDir)
	require.NoError(t, err)
	barBuildFiles, err := listRecursive(barBuildDir)
	require.NoError(t, err)
	assert.Contains(t, fooBuildFiles, "build/0.1.0/darwin-amd64/foo")
	assert.Contains(t, fooBuildFiles, "build/0.1.0/linux-amd64/foo")
	assert.Contains(t, barBuildFiles, "build/0.1.0/linux-amd64/bar")

	// Cleaning foo only removes foo builds
	execCommand(t, testProjectDir, "./godelw", "clean", "foo")
	fooBuildFiles, err = listRecursive(fooBuildDir)
	require.NoError(t, err)
	barBuildFiles, err = listRecursive(barBuildDir)
	require.NoError(t, err)
	assert.Contains(t, fooBuildFiles, "build/test-file", "Non-build files in build dir should not be removed")
	assert.NotContains(t, fooBuildFiles, "build/0.1.0/darwin-amd64/foo", "Build for foo should have been removed")
	assert.NotContains(t, fooBuildFiles, "build/0.1.0/linux-amd64/foo", "Build for foo should have been removed")
	assert.Contains(t, barBuildFiles, "build/0.1.0/linux-amd64/bar", "Build for bar should exist")

	// Cleaning foo with force removes everything from the build dir
	execCommand(t, testProjectDir, "./godelw", "clean", "foo", "--force")
	fooBuildFiles, err = listRecursive(fooBuildDir)
	require.NoError(t, err)
	barBuildFiles, err = listRecursive(barBuildDir)
	require.NoError(t, err)
	assert.NotContains(t, fooBuildFiles, "build/test-file", "Non-build files in build dir should get removed with force")
	assert.Contains(t, barBuildFiles, "build/0.1.0/linux-amd64/bar", "Build for bar should exist")

	// Clean with no products removes default builds
	execCommand(t, testProjectDir, "./godelw", "clean")
	barBuildFiles, err = listRecursive(barBuildDir)
	require.NoError(t, err)
	assert.NotContains(t, barBuildFiles, "build/0.1.0/linux-amd64/bar", "Build for bar should have been removed")
}

func TestArtifactsDistClean(t *testing.T) {
	testProjectDir := setUpGödelTestAndDownload(t, testRootDir, gödelTGZ, version)
	gittest.InitGitDir(t, testProjectDir)

	distYml := `
group-id: com.palantir.godel
products:
  foo:
    build:
      main-pkg: ./foo
      os-archs:
        - os: linux
          arch: amd64
    dist:
      output-dir: ./foo/dist
  bar:
    build:
      main-pkg: ./bar
      os-archs:
        - os: darwin
          arch: amd64
`
	err := ioutil.WriteFile(path.Join(testProjectDir, "godel", "config", "dist.yml"), []byte(distYml), 0644)
	require.NoError(t, err)

	src := `package main
	import "fmt"

	func main() {
		fmt.Println("hello, world!")
	}`
	err = os.MkdirAll(path.Join(testProjectDir, "foo"), 0755)
	require.NoError(t, err)
	err = ioutil.WriteFile(path.Join(testProjectDir, "foo", "foo.go"), []byte(src), 0644)
	require.NoError(t, err)

	err = os.MkdirAll(path.Join(testProjectDir, "bar"), 0755)
	require.NoError(t, err)
	err = ioutil.WriteFile(path.Join(testProjectDir, "bar", "bar.go"), []byte(src), 0644)
	require.NoError(t, err)

	gittest.CommitAllFiles(t, testProjectDir, "Commit files")
	gittest.CreateGitTag(t, testProjectDir, "0.1.0")

	execCommand(t, testProjectDir, "./godelw", "dist")
	fooDistDir := path.Join(testProjectDir, "foo")
	barDistDir := path.Join(testProjectDir)

	// Dists exist for both products
	fooDistFiles, err := listRecursive(fooDistDir)
	require.NoError(t, err)
	barDistFiles, err := listRecursive(barDistDir)
	require.NoError(t, err)
	assert.Contains(t, fooDistFiles, "dist/foo-0.1.0.sls.tgz")
	assert.Contains(t, barDistFiles, "dist/bar-0.1.0.sls.tgz")

	// Cleaning foo only removes foo dists
	execCommand(t, testProjectDir, "./godelw", "clean", "foo")
	fooDistFiles, err = listRecursive(fooDistDir)
	require.NoError(t, err)
	barDistFiles, err = listRecursive(barDistDir)
	require.NoError(t, err)
	assert.NotContains(t, fooDistFiles, "dist/foo-0.1.0.sls.tgz", "Distribution for foo should have been removed")
	assert.Contains(t, barDistFiles, "dist/bar-0.1.0.sls.tgz", "Distribution for bar should exist")

	// Clean with no products removes default dists
	execCommand(t, testProjectDir, "./godelw", "clean")
	barDistFiles, err = listRecursive(barDistDir)
	require.NoError(t, err)
	assert.NotContains(t, barDistFiles, "dist/bar-0.1.0.sls.tgz", "Distribution for bar should have been removed")
}

func TestTest(t *testing.T) {
	testProjectDir := setUpGödelTestAndDownload(t, testRootDir, gödelTGZ, version)
	src := `package foo_test
	import "testing"

	func TestFoo(t *testing.T) {}`
	err := ioutil.WriteFile(path.Join(testProjectDir, "foo_test.go"), []byte(src), 0644)
	require.NoError(t, err)

	execCommand(t, testProjectDir, "./godelw", "test")
}

// Run "../godelw check" and ensure that it works (command supports being invoked from subdirectory). The action should
// execute with the subdirectory as the working directory.
func TestCheckFromNestedDirectory(t *testing.T) {
	testProjectDir := setUpGödelTestAndDownload(t, testRootDir, gödelTGZ, version)

	// write invalid Go file to root directory of project
	badSrc := `badContentForGoFile`
	err := ioutil.WriteFile(path.Join(testProjectDir, "main.go"), []byte(badSrc), 0644)
	require.NoError(t, err)

	// write valid Go file to child directory
	childDir := path.Join(testProjectDir, "childDir")
	err = os.MkdirAll(childDir, 0755)
	require.NoError(t, err)
	src := `package main
	import "fmt"

	func main() {
		fmt.Println("hello, world!")
	}`
	err = ioutil.WriteFile(path.Join(childDir, "main.go"), []byte(src), 0644)
	require.NoError(t, err)

	execCommand(t, childDir, "../godelw", "check")
}

// Run "../godelw build" and ensure that it works (command supports being invoked from sub-directory). The build action
// should execute with the root project directory as the working directory. Verifies #235.
func TestBuildFromNestedDirectory(t *testing.T) {
	testProjectDir := setUpGödelTestAndDownload(t, testRootDir, gödelTGZ, version)
	src := `package main
	import "fmt"

	func main() {
		fmt.Println("hello, world!")
	}`
	err := ioutil.WriteFile(path.Join(testProjectDir, "main.go"), []byte(src), 0644)
	require.NoError(t, err)

	childDir := path.Join(testProjectDir, "childDir")
	err = os.MkdirAll(childDir, 0755)
	require.NoError(t, err)

	execCommand(t, childDir, "../godelw", "build")

	info, err := os.Stat(path.Join(testProjectDir, "build"))
	require.NoError(t, err)
	assert.True(t, info.IsDir())
}

// Run "./godelw publish" and verify that it prints a help message and exits with a non-0 exit code. Verifies #243.
func TestPublishWithNoAction(t *testing.T) {
	testProjectDir := setUpGödelTestAndDownload(t, testRootDir, gödelTGZ, version)
	src := `package main
	import "fmt"

	func main() {
		fmt.Println("hello, world!")
	}`
	err := ioutil.WriteFile(path.Join(testProjectDir, "main.go"), []byte(src), 0644)
	require.NoError(t, err)

	cmd := exec.Command("./godelw", "publish")
	cmd.Dir = testProjectDir
	output, err := cmd.CombinedOutput()
	require.Error(t, err)
	assert.Regexp(t, regexp.MustCompile(`(?s)NAME:.+publish - Publish product distributions.+USAGE:.+godel publish.+SUBCOMMANDS:.+FLAGS:.+`), string(output))
}

func TestVerify(t *testing.T) {
	testProjectDir := setUpGödelTestAndDownload(t, testRootDir, gödelTGZ, version)
	const (
		src = `package main
	import "fmt"

	func main() {
		fmt.Println("hello, world!")
	}`
		testSrc = `package main_test
	import "testing"

	func TestFoo(t *testing.T) {
		t=t
		t.Fail()
	}`
		importsYML = `root-dirs:
  - .`
		licenseYML = `header: |
  /*
  Copyright 2016 Palantir Technologies, Inc.

  Licensed under the Apache License, Version 2.0 (the "License");
  you may not use this file except in compliance with the License.
  You may obtain a copy of the License at

  http://www.apache.org/licenses/LICENSE-2.0

  Unless required by applicable law or agreed to in writing, software
  distributed under the License is distributed on an "AS IS" BASIS,
  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
  See the License for the specific language governing permissions and
  limitations under the License.
  */
`
	)
	err := ioutil.WriteFile(path.Join(testProjectDir, "main.go"), []byte(src), 0644)
	require.NoError(t, err)
	err = ioutil.WriteFile(path.Join(testProjectDir, "main_test.go"), []byte(testSrc), 0644)
	require.NoError(t, err)
	err = ioutil.WriteFile(path.Join(testProjectDir, "godel", "config", "imports.yml"), []byte(importsYML), 0644)
	require.NoError(t, err)
	err = ioutil.WriteFile(path.Join(testProjectDir, "godel", "config", "license.yml"), []byte(licenseYML), 0644)
	require.NoError(t, err)

	for i, currCase := range []struct {
		args []string
		want string
	}{
		{want: `(?s).+Failed tasks:\n\tformat -v -l\n\timports --verify\n\tlicense --verify\n\tcheck\n\ttest`},
		{args: []string{"--skip-format"}, want: `(?s).+Failed tasks:\n\timports --verify\n\tlicense --verify\n\tcheck\n\ttest`},
		{args: []string{"--skip-check"}, want: `(?s).+Failed tasks:\n\tformat -v -l\n\timports --verify\n\tlicense --verify\n\ttest`},
		{args: []string{"--skip-imports"}, want: `(?s).+Failed tasks:\n\tformat -v -l\n\tlicense --verify\n\tcheck\n\ttest`},
		{args: []string{"--skip-license"}, want: `(?s).+Failed tasks:\n\tformat -v -l\n\timports --verify\n\tcheck\n\ttest`},
		{args: []string{"--skip-test"}, want: `(?s).+Failed tasks:\n\tformat -v -l\n\timports --verify\n\tlicense --verify\n\tcheck`},
	} {
		cmd := exec.Command("./godelw", append([]string{"verify", "--apply=false"}, currCase.args...)...)
		cmd.Dir = testProjectDir
		output, err := cmd.CombinedOutput()
		require.Error(t, err)
		assert.Regexp(t, regexp.MustCompile(currCase.want), string(output), "Case %d", i)
	}
}

func TestVerifyApply(t *testing.T) {
	testProjectDir := setUpGödelTestAndDownload(t, testRootDir, gödelTGZ, version)
	const (
		src = `package main
	import "fmt"

	func main() {
		fmt.Println("hello, world!")
	}`
		testSrc = `package main_test
	import "testing"

	func TestFoo(t *testing.T) {
		t=t
		t.Fail()
	}`
		formattedTestSrc = `package main_test

import (
	"testing"
)

func TestFoo(t *testing.T) {
	t = t
	t.Fail()
}
`
		licensedTestSrc = `/*
Copyright 2016 Palantir Technologies, Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main_test
	import "testing"

	func TestFoo(t *testing.T) {
		t=t
		t.Fail()
	}`
		licensedAndFormattedTestSrc = `/*
Copyright 2016 Palantir Technologies, Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main_test

import (
	"testing"
)

func TestFoo(t *testing.T) {
	t = t
	t.Fail()
}
`
		importsYML = `root-dirs:
  - .`
		importsJSON = `{
    "imports": [],
    "mainOnlyImports": [],
    "testOnlyImports": []
}`

		licenseYML = `header: |
  /*
  Copyright 2016 Palantir Technologies, Inc.

  Licensed under the Apache License, Version 2.0 (the "License");
  you may not use this file except in compliance with the License.
  You may obtain a copy of the License at

  http://www.apache.org/licenses/LICENSE-2.0

  Unless required by applicable law or agreed to in writing, software
  distributed under the License is distributed on an "AS IS" BASIS,
  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
  See the License for the specific language governing permissions and
  limitations under the License.
  */
`
	)
	err := ioutil.WriteFile(path.Join(testProjectDir, "main.go"), []byte(src), 0644)
	require.NoError(t, err)
	err = ioutil.WriteFile(path.Join(testProjectDir, "godel", "config", "imports.yml"), []byte(importsYML), 0644)
	require.NoError(t, err)
	err = ioutil.WriteFile(path.Join(testProjectDir, "godel", "config", "license.yml"), []byte(licenseYML), 0644)
	require.NoError(t, err)

	for i, currCase := range []struct {
		args            []string
		want            string
		wantTestSrc     string
		wantImportsJSON string
	}{
		{want: `(?s).+Failed tasks:\n\tcheck\n\ttest`, wantTestSrc: licensedAndFormattedTestSrc, wantImportsJSON: importsJSON},
		{args: []string{"--skip-format"}, want: `(?s).+Failed tasks:\n\tcheck\n\ttest`, wantTestSrc: licensedTestSrc, wantImportsJSON: importsJSON},
		{args: []string{"--skip-imports"}, want: `(?s).+Failed tasks:\n\tcheck\n\ttest`, wantTestSrc: licensedAndFormattedTestSrc},
		{args: []string{"--skip-license"}, want: `(?s).+Failed tasks:\n\tcheck\n\ttest`, wantTestSrc: formattedTestSrc, wantImportsJSON: importsJSON},
		{args: []string{"--skip-check"}, want: `(?s).+Failed tasks:\n\ttest`, wantTestSrc: licensedAndFormattedTestSrc, wantImportsJSON: importsJSON},
		{args: []string{"--skip-test"}, want: `(?s).+Failed tasks:\n\tcheck`, wantTestSrc: licensedAndFormattedTestSrc, wantImportsJSON: importsJSON},
	} {
		err = ioutil.WriteFile(path.Join(testProjectDir, "main_test.go"), []byte(testSrc), 0644)
		require.NoError(t, err, "Case %d", i)

		cmd := exec.Command("./godelw", append([]string{"verify"}, currCase.args...)...)
		cmd.Dir = testProjectDir
		output, err := cmd.CombinedOutput()
		require.Error(t, err, fmt.Sprintf("Case %d", i))
		assert.Regexp(t, regexp.MustCompile(currCase.want), string(output), "Case %d", i)

		bytes, err := ioutil.ReadFile(path.Join(testProjectDir, "main_test.go"))
		require.NoError(t, err, "Case %d", i)
		assert.Equal(t, currCase.wantTestSrc, string(bytes), "Case %d", i)

		importsJSONPath := path.Join(testProjectDir, "gocd_imports.json")
		if currCase.wantImportsJSON == "" {
			_, err = os.Stat(importsJSONPath)
			assert.True(t, os.IsNotExist(err), "Case %d: gocd_imports.json should not exist", i)
		} else {
			bytes, err = ioutil.ReadFile(importsJSONPath)
			require.NoError(t, err, "Case %d", i)
			assert.Equal(t, currCase.wantImportsJSON, string(bytes), "Case %d", i)
			err = os.Remove(importsJSONPath)
			require.NoError(t, err, "Case %d", i)
		}
	}
}

func TestVerifyWithJUnitOutput(t *testing.T) {
	testProjectDir := setUpGödelTestAndDownload(t, testRootDir, gödelTGZ, version)
	src := `package main
	import "fmt"
	func main() {
		fmt.Println("hello, world!")
	}`
	err := ioutil.WriteFile(path.Join(testProjectDir, "main.go"), []byte(src), 0644)
	require.NoError(t, err)
	testSrc := `package main_test
	import "testing"
	func TestFoo(t *testing.T) {
	}`
	err = ioutil.WriteFile(path.Join(testProjectDir, "main_test.go"), []byte(testSrc), 0644)
	require.NoError(t, err)

	junitOutputFile := "test-output.xml"
	cmd := exec.Command("./godelw", "verify", "--apply=false", "--junit-output", junitOutputFile)
	cmd.Dir = testProjectDir
	err = cmd.Run()
	require.Error(t, err)

	fi, err := os.Stat(path.Join(testProjectDir, junitOutputFile))
	require.NoError(t, err)

	assert.False(t, fi.IsDir())
}

func TestDebugFlagPrintsStackTrace(t *testing.T) {
	testProjectDir := setUpGödelTestAndDownload(t, testRootDir, gödelTGZ, version)

	cmd := exec.Command("./godelw", "install", "foo")
	cmd.Dir = testProjectDir
	output, err := cmd.CombinedOutput()
	require.Error(t, err)
	assert.Regexp(t, `^Failed to install from foo into .+: foo does not exist\n$`, string(output))

	cmd = exec.Command("./godelw", "--debug", "install", "foo")
	cmd.Dir = testProjectDir
	output, err = cmd.CombinedOutput()
	require.Error(t, err)
	assert.Regexp(t, `(?s)^foo does not exist.+cmd/godel.localPkg.getPkg.+Failed to install from foo into .+`, string(output))
}

// * Symlink "test-go" -> $GOPATH
// * Set current directory to test project inside the symlink
// * Verify that "./godelw check" works in sym-linked path
func TestCheckInGoPathSymLink(t *testing.T) {
	testProjectDir := setUpGödelTestAndDownload(t, testRootDir, gödelTGZ, version)
	src := `package foo_test
	import "testing"

	func TestFoo(t *testing.T) {}`
	err := ioutil.WriteFile(path.Join(testProjectDir, "foo_test.go"), []byte(src), 0644)
	require.NoError(t, err)

	symLinkParentDir, cleanup, err := dirs.TempDir("", "")
	defer cleanup()
	require.NoError(t, err)
	symLinkPath := path.Join(symLinkParentDir, "test-go")

	originalGoPath := os.Getenv("GOPATH")
	err = os.Symlink(originalGoPath, symLinkPath)
	require.NoError(t, err)

	testProjectRelPath, err := filepath.Rel(originalGoPath, testProjectDir)
	require.NoError(t, err)

	// use script to set cd because setting wd on exec.Command does not work for symlinks
	projectPathInSymLink := path.Join(symLinkPath, testProjectRelPath)
	scriptTemplate := `#!/bin/bash
cd %v
pwd
`
	scriptFilePath := path.Join(symLinkParentDir, "script.sh")
	err = ioutil.WriteFile(scriptFilePath, []byte(fmt.Sprintf(scriptTemplate, projectPathInSymLink)), 0755)
	require.NoError(t, err)

	cmd := exec.Command(scriptFilePath)
	output, err := cmd.CombinedOutput()
	require.NoError(t, err, "Command %v failed. Output:\n%v", cmd.Args, string(output))
	assert.Equal(t, projectPathInSymLink, strings.TrimSpace(string(output)))

	scriptTemplate = `#!/bin/bash
cd %v
./godelw check
`
	err = ioutil.WriteFile(scriptFilePath, []byte(fmt.Sprintf(scriptTemplate, projectPathInSymLink)), 0755)
	require.NoError(t, err)

	cmd = exec.Command(scriptFilePath)
	output, err = cmd.CombinedOutput()
	require.NoError(t, err, "Command %v failed. Output:\n%v", cmd.Args, string(output))
}

// * Symlink "test-go" -> $GOPATH
// * Set $GOPATH to be the symlink ("test-go")
// * Set current directory to test project inside the symlink
// * Verify that "./godelw check" works in sym-linked path
// * Restore $GOPATH to original value
func TestCheckInGoPathSymLinkGoPathSymLink(t *testing.T) {
	testProjectDir := setUpGödelTestAndDownload(t, testRootDir, gödelTGZ, version)
	src := `package foo_test
	import "testing"

	func TestFoo(t *testing.T) {}`
	err := ioutil.WriteFile(path.Join(testProjectDir, "foo_test.go"), []byte(src), 0644)
	require.NoError(t, err)

	symLinkParentDir, cleanup, err := dirs.TempDir("", "")
	defer cleanup()
	require.NoError(t, err)
	symLinkPath := path.Join(symLinkParentDir, "test-go")

	originalGoPath := os.Getenv("GOPATH")
	err = os.Symlink(originalGoPath, symLinkPath)
	require.NoError(t, err)

	err = os.Setenv("GOPATH", symLinkPath)
	require.NoError(t, err)
	defer func() {
		if err := os.Setenv("GOPATH", originalGoPath); err != nil {
			require.NoError(t, err, "failed to restore GOPATH environment variable in defer")
		}
	}()

	testProjectRelPath, err := filepath.Rel(originalGoPath, testProjectDir)
	require.NoError(t, err)

	// use script to set cd because setting wd on exec.Command does not work for symlinks
	projectPathInSymLink := path.Join(symLinkPath, testProjectRelPath)
	scriptTemplate := `#!/bin/bash
cd %v
pwd
`
	scriptFilePath := path.Join(symLinkParentDir, "script.sh")
	err = ioutil.WriteFile(scriptFilePath, []byte(fmt.Sprintf(scriptTemplate, projectPathInSymLink)), 0755)
	require.NoError(t, err)

	cmd := exec.Command(scriptFilePath)
	output, err := cmd.CombinedOutput()
	require.NoError(t, err, "Command %v failed. Output:\n%v", cmd.Args, string(output))
	assert.Equal(t, projectPathInSymLink, strings.TrimSpace(string(output)))

	scriptTemplate = `#!/bin/bash
cd %v
./godelw check
`
	err = ioutil.WriteFile(scriptFilePath, []byte(fmt.Sprintf(scriptTemplate, projectPathInSymLink)), 0755)
	require.NoError(t, err)

	cmd = exec.Command(scriptFilePath)
	output, err = cmd.CombinedOutput()
	require.NoError(t, err, "Command %v failed. Output:\n%v", cmd.Args, string(output))
}

// * Symlink "test-go" -> $GOPATH
// * Set $GOPATH to be the symlink ("test-go")
// * Set current directory to real project (not inside symlink)
// * Verify that "./godelw check" works in real path
// * Restore $GOPATH to original value
func TestCheckInGoPathNonSymLinkWhenGoPathIsSymLink(t *testing.T) {
	testProjectDir := setUpGödelTestAndDownload(t, testRootDir, gödelTGZ, version)
	src := `package foo_test
	import "testing"

	func TestFoo(t *testing.T) {}`
	err := ioutil.WriteFile(path.Join(testProjectDir, "foo_test.go"), []byte(src), 0644)
	require.NoError(t, err)

	symLinkParentDir, cleanup, err := dirs.TempDir("", "")
	defer cleanup()
	require.NoError(t, err)
	symLinkPath := path.Join(symLinkParentDir, "test-go")

	originalGoPath := os.Getenv("GOPATH")
	err = os.Symlink(originalGoPath, symLinkPath)
	require.NoError(t, err)

	err = os.Setenv("GOPATH", symLinkPath)
	require.NoError(t, err)
	defer func() {
		if err := os.Setenv("GOPATH", originalGoPath); err != nil {
			require.NoError(t, err, "failed to restore GOPATH environment variable in defer")
		}
	}()

	cmd := exec.Command("./godelw", "check")
	cmd.Dir = testProjectDir
	output, err := cmd.CombinedOutput()
	require.NoError(t, err, "Command %v failed. Output:\n%v", cmd.Args, string(output))
}

func execCommand(t *testing.T, dir, cmdName string, args ...string) string {
	cmd := exec.Command(cmdName, args...)
	cmd.Dir = dir
	output, err := cmd.CombinedOutput()
	require.NoError(t, err, "Command %v failed. Output:\n%v", cmd.Args, string(output))
	return string(output)
}

func listRecursive(dir string) ([]string, error) {
	fileList := []string{}
	err := filepath.Walk(dir, func(path string, f os.FileInfo, err error) error {
		fileList = append(fileList, strings.TrimPrefix(strings.TrimPrefix(path, dir), "/"))
		return nil
	})
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return fileList, nil
}
