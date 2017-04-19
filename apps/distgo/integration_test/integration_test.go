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
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/palantir/godel/apps/distgo/pkg/git/gittest"
	"github.com/palantir/godel/pkg/products"
)

func TestRun(t *testing.T) {
	cli, err := products.Bin("distgo")
	require.NoError(t, err)

	wd, err := os.Getwd()
	require.NoError(t, err)

	tmpDir, cleanup, err := dirs.TempDir(wd, "")
	defer cleanup()
	require.NoError(t, err)

	for i, currCase := range []struct {
		name          string
		filesToCreate []gofiles.GoFileSpec
		config        string
		args          []string
		wantStdout    string
	}{
		{
			name: "Run runs program",
			filesToCreate: []gofiles.GoFileSpec{
				{
					RelPath: "main.go",
					Src: `package main

					import "fmt"

					func main() {
						fmt.Println("Hello, world!")
					}
					`,
				},
			},
			config: `
products:
  hello:
    build:
      main-pkg: .
`,
			args: []string{
				"--product",
				"hello",
			},
			wantStdout: "Hello, world!\n",
		},
		{
			name: "Run uses trailing arguments",
			filesToCreate: []gofiles.GoFileSpec{
				{
					RelPath: "main.go",
					Src: `package main

					import (
						"fmt"
						"os"
					)

					func main() {
						fmt.Println(os.Args[1:])
					}
					`,
				},
			},
			config: `
products:
  hello:
    build:
      main-pkg: .
`,
			args: []string{
				"--product",
				"hello",
				"arg1",
				"arg2",
				"arg3",
			},
			wantStdout: "[arg1 arg2 arg3]\n",
		},
		{
			name: "Run uses trailing arguments and supports flags",
			filesToCreate: []gofiles.GoFileSpec{
				{
					RelPath: "main.go",
					Src: `package main

					import (
						"fmt"
						"os"
					)

					func main() {
						fmt.Println(os.Args[1:])
					}
					`,
				},
			},
			config: `
products:
  hello:
    build:
      main-pkg: .
`,
			args: []string{
				`flag:--foo-arg`,
				"flag:",
				"flag:flag:",
				"arg3",
			},
			wantStdout: "[--foo-arg flag: flag: arg3]\n",
		},
	} {
		currCaseTmpDir, err := ioutil.TempDir(tmpDir, "")
		require.NoError(t, err)

		_, err = gofiles.Write(currCaseTmpDir, currCase.filesToCreate)
		require.NoError(t, err, "Case %d", i)

		configFile := path.Join(currCaseTmpDir, "config.yml")
		err = ioutil.WriteFile(configFile, []byte(currCase.config), 0644)
		require.NoError(t, err)

		var output []byte
		func() {
			err := os.Chdir(currCaseTmpDir)
			defer func() {
				err := os.Chdir(wd)
				require.NoError(t, err)
			}()
			require.NoError(t, err)

			args := []string{"--config", configFile, "run"}
			args = append(args, currCase.args...)
			cmd := exec.Command(cli, args...)
			output, err = cmd.CombinedOutput()
			require.NoError(t, err, "Case %d: %s\nOutput: %s", i, currCase.name, string(output))
		}()

		content := string(output)[strings.Index(string(output), "\n")+1:]
		assert.Equal(t, currCase.wantStdout, content, "Case %d: %s", i, currCase.name)
	}
}

func TestRunWithStdin(t *testing.T) {
	cli, err := products.Bin("distgo")
	require.NoError(t, err)

	wd, err := os.Getwd()
	require.NoError(t, err)

	tmpDir, cleanup, err := dirs.TempDir(wd, "")
	defer cleanup()
	require.NoError(t, err)

	currCaseTmpDir, err := ioutil.TempDir(tmpDir, "")
	require.NoError(t, err)

	filesToCreate := []gofiles.GoFileSpec{
		{
			RelPath: "main.go",
			Src: `package main

			import (
				"bufio"
				"fmt"
				"os"
			)

			func main() {
				reader := bufio.NewReader(os.Stdin)
				text, _ := reader.ReadString('\n')
				fmt.Printf("read: %q", text)
			}
			`,
		},
	}
	config := `
products:
  hello:
    build:
      main-pkg: .
`
	runArgs := []string{
		"--product",
		"hello",
	}

	stdInContent := "output passed to stdin\n"

	_, err = gofiles.Write(currCaseTmpDir, filesToCreate)
	require.NoError(t, err)

	configFile := path.Join(currCaseTmpDir, "config.yml")
	err = ioutil.WriteFile(configFile, []byte(config), 0644)
	require.NoError(t, err)

	var output []byte
	func() {
		err := os.Chdir(currCaseTmpDir)
		defer func() {
			err := os.Chdir(wd)
			require.NoError(t, err)
		}()
		require.NoError(t, err)

		args := []string{"--config", configFile, "run"}
		args = append(args, runArgs...)
		cmd := exec.Command(cli, args...)

		stdinPipe, err := cmd.StdinPipe()
		require.NoError(t, err)
		_, err = stdinPipe.Write([]byte(stdInContent))
		require.NoError(t, err)

		output, err = cmd.CombinedOutput()
		require.NoError(t, err, "Output: %s", string(output))
	}()

	content := string(output)[strings.Index(string(output), "\n")+1:]
	assert.Equal(t, fmt.Sprintf("read: %q", stdInContent), content)
}

func TestProjectVersion(t *testing.T) {
	cli, err := products.Bin("distgo")
	require.NoError(t, err)

	wd, err := os.Getwd()
	require.NoError(t, err)

	tmpDir, cleanup, err := dirs.TempDir("", "")
	defer cleanup()
	require.NoError(t, err)

	for i, currCase := range []struct {
		name       string
		gitOps     func(t *testing.T, caseNum int, caseName, currCaseDir string)
		wantStdout string
	}{
		{
			name: "prints unspecified when project has no git directory",
			gitOps: func(t *testing.T, caseNum int, caseName, currCaseDir string) {
			},
			wantStdout: "^unspecified\n$",
		},
		{
			name: "prints tag for tagged commit",
			gitOps: func(t *testing.T, caseNum int, caseName, currCaseDir string) {
				gittest.InitGitDir(t, currCaseDir)
				gittest.CreateGitTag(t, currCaseDir, "testCaseTag")
			},
			wantStdout: "^testCaseTag\n$",
		},
		{
			name: "prints tag.dirty for tagged commit with uncommitted files",
			gitOps: func(t *testing.T, caseNum int, caseName, currCaseDir string) {
				gittest.InitGitDir(t, currCaseDir)
				gittest.CreateGitTag(t, currCaseDir, "testCaseTag")
				err := ioutil.WriteFile(path.Join(currCaseDir, "random.txt"), []byte(""), 0644)
				require.NoError(t, err, "Case %d: %s", caseNum, caseName)
			},
			wantStdout: "^testCaseTag.dirty\n$",
		},
		{
			name: "prints version for non-tagged commit",
			gitOps: func(t *testing.T, caseNum int, caseName, currCaseDir string) {
				gittest.InitGitDir(t, currCaseDir)
				gittest.CreateGitTag(t, currCaseDir, "testCaseTag")
				gittest.CommitRandomFile(t, currCaseDir, "Test commit message")
			},
			wantStdout: "^testCaseTag-1-g[a-f0-9]{7}\n$",
		},
		{
			name: "prints version.dirty for non-tagged commit with uncommitted files",
			gitOps: func(t *testing.T, caseNum int, caseName, currCaseDir string) {
				gittest.InitGitDir(t, currCaseDir)
				gittest.CreateGitTag(t, currCaseDir, "testCaseTag")
				gittest.CommitRandomFile(t, currCaseDir, "Test commit message")
				err := ioutil.WriteFile(path.Join(currCaseDir, "random.txt"), []byte(""), 0644)
				require.NoError(t, err, "Case %d: %s", caseNum, caseName)
			},
			wantStdout: "^testCaseTag-1-g[a-f0-9]{7}.dirty\n$",
		},
	} {
		currCaseTmpDir, err := ioutil.TempDir(tmpDir, "")
		require.NoError(t, err)

		currCase.gitOps(t, i, currCase.name, currCaseTmpDir)

		var output []byte
		func() {
			err := os.Chdir(currCaseTmpDir)
			defer func() {
				err := os.Chdir(wd)
				require.NoError(t, err)
			}()
			require.NoError(t, err)

			cmd := exec.Command(cli, "project-version")
			output, err = cmd.CombinedOutput()
			require.NoError(t, err, "Case %d: %s\nOutput: %s", i, currCase.name, string(output))
		}()

		assert.Regexp(t, currCase.wantStdout, string(output))
	}
}
