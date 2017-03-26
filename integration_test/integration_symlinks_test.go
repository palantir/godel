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
	"strings"
	"testing"

	"github.com/nmiyake/pkg/dirs"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

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
