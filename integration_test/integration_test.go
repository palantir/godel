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
	"regexp"
	"strings"
	"testing"

	"github.com/nmiyake/pkg/dirs"
	"github.com/nmiyake/pkg/gofiles"
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

func TestProjectVersion(t *testing.T) {
	tmpDir, cleanup, err := dirs.TempDir("", "")
	defer cleanup()
	require.NoError(t, err)

	gittest.InitGitDir(t, tmpDir)
	gittest.CreateGitTag(t, tmpDir, "testTag")
	err = ioutil.WriteFile(path.Join(tmpDir, "random.txt"), []byte(""), 0644)
	require.NoError(t, err)

	testProjectDir := setUpGödelTestAndDownload(t, tmpDir, gödelTGZ, version)
	output := execCommand(t, testProjectDir, "./godelw", "project-version")
	assert.Equal(t, "testTag.dirty\n", string(output))
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

func TestGenerate(t *testing.T) {
	testProjectDir := setUpGödelTestAndDownload(t, testRootDir, gödelTGZ, version)

	const generateYML = `
generators:
  foo:
    go-generate-dir: gen
    gen-paths:
      paths:
        - "gen/output.txt"
`
	err := ioutil.WriteFile(path.Join(testProjectDir, "godel", "config", "generate.yml"), []byte(generateYML), 0644)
	require.NoError(t, err)

	specs := []gofiles.GoFileSpec{
		{
			RelPath: "gen/testbar.go",
			Src: `package testbar

//go:generate go run generator_main.go
`,
		},
		{
			RelPath: "gen/generator_main.go",
			Src: `// +build ignore

package main

import (
	"io/ioutil"
)

func main() {
	if err := ioutil.WriteFile("output.txt", []byte("foo-output"), 0644); err != nil {
		panic(err)
	}
}
`,
		},
	}

	_, err = gofiles.Write(testProjectDir, specs)
	require.NoError(t, err)

	execCommand(t, testProjectDir, "./godelw", "generate")

	content, err := ioutil.ReadFile(path.Join(testProjectDir, "gen", "output.txt"))
	require.NoError(t, err)
	assert.Equal(t, "foo-output", string(content))
}

func TestGenerateVerify(t *testing.T) {
	testProjectDir := setUpGödelTestAndDownload(t, testRootDir, gödelTGZ, version)

	const generateYML = `
generators:
  foo:
    go-generate-dir: gen
    gen-paths:
      paths:
        - "gen/output.txt"
`
	err := ioutil.WriteFile(path.Join(testProjectDir, "godel", "config", "generate.yml"), []byte(generateYML), 0644)
	require.NoError(t, err)

	specs := []gofiles.GoFileSpec{
		{
			RelPath: "gen/testbar.go",
			Src: `package testbar

//go:generate go run generator_main.go
`,
		},
		{
			RelPath: "gen/generator_main.go",
			Src: `// +build ignore

package main

import (
	"io/ioutil"
)

func main() {
	if err := ioutil.WriteFile("output.txt", []byte("foo-output"), 0644); err != nil {
		panic(err)
	}
}
`,
		},
	}

	_, err = gofiles.Write(testProjectDir, specs)
	require.NoError(t, err)

	err = ioutil.WriteFile(path.Join(testProjectDir, "gen", "output.txt"), []byte("original"), 0644)
	require.NoError(t, err)

	cmd := exec.Command("./godelw", "generate", "--verify")
	cmd.Dir = testProjectDir
	output, err := cmd.CombinedOutput()
	require.Error(t, err)

	want := "Generators produced output that differed from what already exists: [foo]\n  foo:\n    gen/output.txt: previously had checksum 0682c5f2076f099c34cfdd15a9e063849ed437a49677e6fcc5b4198c76575be5, now has checksum 380a300b764683667309818ff127a401c6ea6ab1959f386fe0f05505d660ba37\n"
	assert.Equal(t, want, string(output))
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
    dist:
      dist-type:
        type: sls
  bar:
    build:
      main-pkg: ./bar
    dist:
      dist-type:
        type: sls
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

func TestRun(t *testing.T) {
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

	fooSrc := `package main
	import (
		"fmt"
		"os"
	)

	func main() {
		fmt.Println("foo:", os.Args[1:])
	}`
	err = os.MkdirAll(path.Join(testProjectDir, "foo"), 0755)
	require.NoError(t, err)
	err = ioutil.WriteFile(path.Join(testProjectDir, "foo", "foo.go"), []byte(fooSrc), 0644)
	require.NoError(t, err)

	barSrc := `package main
	import (
		"fmt"
	)

	func main() {
		fmt.Println("bar")
	}`
	err = os.MkdirAll(path.Join(testProjectDir, "bar"), 0755)
	require.NoError(t, err)
	err = ioutil.WriteFile(path.Join(testProjectDir, "bar", "bar.go"), []byte(barSrc), 0644)
	require.NoError(t, err)

	gittest.CommitAllFiles(t, testProjectDir, "Commit files")
	gittest.CreateGitTag(t, testProjectDir, "0.1.0")

	output := execCommand(t, testProjectDir, "./godelw", "run", "--product", "foo", "arg1", "arg2")
	output = output[strings.Index(output, "\n")+1:]
	assert.Equal(t, "foo: [arg1 arg2]\n", output)
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
	assert.Regexp(t, `(?s)NAME:.+publish - Publish product distributions.+USAGE:.+godel publish.+SUBCOMMANDS:.+FLAGS:.+`, string(output))
}

func TestDebugFlagPrintsStackTrace(t *testing.T) {
	testProjectDir := setUpGödelTestAndDownload(t, testRootDir, gödelTGZ, version)

	cmd := exec.Command("./godelw", "install", "foo")
	cmd.Dir = testProjectDir
	output, err := cmd.CombinedOutput()
	require.Error(t, err)
	assert.Regexp(t, `^Error: failed to install from foo into .+: foo does not exist\n$`, string(output))

	cmd = exec.Command("./godelw", "--debug", "install", "foo")
	cmd.Dir = testProjectDir
	output, err = cmd.CombinedOutput()
	require.Error(t, err)
	assert.Regexp(t, `(?s)^Error: foo does not exist.+`+regexp.QuoteMeta(`github.com/palantir/godel/godelgetter.(*localFilePkg).Reader`)+`.+failed to install from foo into .+`, string(output))
}

func execCommand(t *testing.T, dir, cmdName string, args ...string) string {
	cmd := exec.Command(cmdName, args...)
	cmd.Dir = dir
	output, err := cmd.CombinedOutput()
	require.NoError(t, err, "Command %v failed. Output:\n%v", cmd.Args, string(output))
	return string(output)
}
