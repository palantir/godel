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
	"testing"

	"github.com/nmiyake/pkg/dirs"
	"github.com/palantir/distgo/pkg/git"
	"github.com/palantir/pkg/gittest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

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
	assert.Equal(t, "testTag-dirty\n", string(output))
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
	assert.Regexp(t, `(?s)^`+regexp.QuoteMeta("Unformatted files exist -- run ./godelw format to format these files:")+"\n"+regexp.QuoteMeta("  helper.go")+"\n$", string(output))
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

	// write Go file to root directory of project
	badSrc := `package main`
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
