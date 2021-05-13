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
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"github.com/nmiyake/pkg/dirs"
	"github.com/palantir/distgo/pkg/git"
	"github.com/palantir/godel/pkg/products/v2"
	"github.com/palantir/pkg/gittest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	godelTGZ    string
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
	godelProjectDir := filepath.Join(wd, "..")
	version, err = git.ProjectVersion(godelProjectDir)
	if err != nil {
		panic(fmt.Sprintf("Failed to get version from directory %s: %v", godelProjectDir, err))
	}

	godelTGZ, err = products.Dist("godel")
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
	testProjectDir := setUpGodelTestAndDownload(t, testRootDir, godelTGZ, version)

	output := execCommand(t, testProjectDir, "./godelw", "--version")
	assert.Equal(t, fmt.Sprintf("godel version %v\n", version), output)
}

func TestProjectVersion(t *testing.T) {
	tmpDir, cleanup, err := dirs.TempDir("", "")
	defer cleanup()
	require.NoError(t, err)

	gittest.InitGitDir(t, tmpDir)
	gittest.CreateGitTag(t, tmpDir, "testTag")
	err = os.WriteFile(filepath.Join(tmpDir, "random.txt"), []byte(""), 0644)
	require.NoError(t, err)

	testProjectDir := setUpGodelTestAndDownload(t, tmpDir, godelTGZ, version)
	output := execCommand(t, testProjectDir, "./godelw", "project-version")
	assert.Equal(t, "testTag.dirty\n", output)
}

func TestGitHooksSuccess(t *testing.T) {
	// create project directory in temporary location so primary project's repository is not modified by test
	tmp, cleanup, err := dirs.TempDir("", "")
	defer cleanup()
	require.NoError(t, err)

	testProjectDir := setUpGodelTestAndDownload(t, tmp, godelTGZ, version)

	// initialize git repository
	gittest.InitGitDir(t, testProjectDir)

	// install commit hooks
	execCommand(t, testProjectDir, "./godelw", "git-hooks")

	// committing Go file that is properly formatted works
	formatted := `package main

func main() {
}
`
	err = os.WriteFile(filepath.Join(testProjectDir, "main.go"), []byte(formatted), 0644)
	require.NoError(t, err)
	execCommand(t, testProjectDir, "git", "add", ".")
	execCommand(t, testProjectDir, "git", "commit", "--author=testAuthor <test@author.com>", "-m", "Second commit")
}

func TestGitHooksFail(t *testing.T) {
	// create project directory in temporary location so primary project's repository is not modified by test
	tmp, cleanup, err := dirs.TempDir("", "")
	defer cleanup()
	require.NoError(t, err)

	testProjectDir := setUpGodelTestAndDownload(t, tmp, godelTGZ, version)

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
	err = os.WriteFile(filepath.Join(testProjectDir, "helper.go"), []byte(notFormatted), 0644)
	require.NoError(t, err)
	execCommand(t, testProjectDir, "git", "add", ".")

	cmd := exec.Command("git", "commit", "--author=testAuthor <test@author.com>", "-m", "Second commit")
	cmd.Dir = testProjectDir
	output, err := cmd.CombinedOutput()
	assert.Error(t, err, "exit status 1")
	assert.Regexp(t, `(?s)^`+regexp.QuoteMeta("Unformatted files exist -- run ./godelw format to format these files:")+"\n"+regexp.QuoteMeta("  helper.go")+"\n$", string(output))
}

func TestProducts(t *testing.T) {
	testProjectDir := setUpGodelTestAndDownload(t, testRootDir, godelTGZ, version)

	distYml := `
products:
  foo:
    build:
      main-pkg: ./foo
  bar:
    build:
      main-pkg: ./bar
`
	err := os.WriteFile(filepath.Join(testProjectDir, "godel", "config", "dist-plugin.yml"), []byte(distYml), 0644)
	require.NoError(t, err)

	src := `package main
	import "fmt"

	func main() {
		fmt.Println("hello, world!")
	}`
	err = os.MkdirAll(filepath.Join(testProjectDir, "foo"), 0755)
	require.NoError(t, err)
	err = os.WriteFile(filepath.Join(testProjectDir, "foo", "foo.go"), []byte(src), 0644)
	require.NoError(t, err)

	err = os.MkdirAll(filepath.Join(testProjectDir, "bar"), 0755)
	require.NoError(t, err)
	err = os.WriteFile(filepath.Join(testProjectDir, "bar", "bar.go"), []byte(src), 0644)
	require.NoError(t, err)

	execCommand(t, testProjectDir, "./godelw", "products")
}

func TestExec(t *testing.T) {
	testProjectDir := setUpGodelTestAndDownload(t, testRootDir, godelTGZ, version)

	currGodelYML, err := os.ReadFile(filepath.Join(testProjectDir, "godel", "config", "godel.yml"))
	require.NoError(t, err)

	updatedGodelYML := string(currGodelYML) + `
environment:
  MY_ENV_VAR: "FOO"
`
	err = os.WriteFile(filepath.Join(testProjectDir, "godel", "config", "godel.yml"), []byte(updatedGodelYML), 0644)
	require.NoError(t, err)

	out := execCommand(t, testProjectDir, "./godelw", "exec", "env")
	idx := strings.Index(out, "MY_ENV_VAR=FOO\n")

	// do not print content of "out" on failure, as it may contain sensitive environment variables
	assert.True(t, idx != -1, "did not find expected environment variable in output")
}

func TestTest(t *testing.T) {
	testProjectDir := setUpGodelTestAndDownload(t, testRootDir, godelTGZ, version)
	src := `package foo_test
	import "testing"

	func TestFoo(t *testing.T) {}`
	err := os.WriteFile(filepath.Join(testProjectDir, "foo_test.go"), []byte(src), 0644)
	require.NoError(t, err)

	execCommand(t, testProjectDir, "./godelw", "test")
}

// Run "../godelw check" and ensure that it works (command supports being invoked from subdirectory). The action should
// execute with the subdirectory as the working directory.
func TestCheckFromNestedDirectory(t *testing.T) {
	testProjectDir := setUpGodelTestAndDownload(t, testRootDir, godelTGZ, version)

	// write Go file to root directory of project
	badSrc := `package main`
	err := os.WriteFile(filepath.Join(testProjectDir, "main.go"), []byte(badSrc), 0644)
	require.NoError(t, err)

	// write valid Go file to child directory
	childDir := filepath.Join(testProjectDir, "childDir")
	err = os.MkdirAll(childDir, 0755)
	require.NoError(t, err)
	src := `package main
	import "fmt"

	func main() {
		fmt.Println("hello, world!")
	}`
	err = os.WriteFile(filepath.Join(childDir, "main.go"), []byte(src), 0644)
	require.NoError(t, err)

	execCommand(t, childDir, "../godelw", "check")
}

func TestDebugFlagPrintsStackTrace(t *testing.T) {
	testProjectDir := setUpGodelTestAndDownload(t, testRootDir, godelTGZ, version)

	cmd := exec.Command("./godelw", "install", "foo")
	cmd.Dir = testProjectDir
	output, err := cmd.CombinedOutput()
	require.Error(t, err)
	assert.Regexp(t, `^Error: failed to install from foo into .+: foo does not exist\n$`, string(output))

	cmd = exec.Command("./godelw", "--debug", "install", "foo")
	cmd.Dir = testProjectDir
	output, err = cmd.CombinedOutput()
	require.Error(t, err)
	assert.Regexp(t, `(?s)^Error: foo does not exist.+`+regexp.QuoteMeta(`github.com/palantir/godel/v2/godelgetter.(*localFilePkg).Reader`)+`.+failed to install from foo into .+`, string(output))
}

func execCommand(t *testing.T, dir, cmdName string, args ...string) string {
	cmd := exec.Command(cmdName, args...)
	cmd.Dir = dir
	output, err := cmd.CombinedOutput()
	require.NoError(t, err, "Command %v failed. Output:\n%v", cmd.Args, string(output))
	return string(output)
}

func execCommandExpectError(t *testing.T, dir, cmdName string, args ...string) string {
	cmd := exec.Command(cmdName, args...)
	cmd.Dir = dir
	output, err := cmd.CombinedOutput()
	require.Error(t, err)
	return string(output)
}
