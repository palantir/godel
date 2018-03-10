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
	"github.com/palantir/godel/pkg/products"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRun(t *testing.T) {
	cli, err := products.Bin("distgo-plugin")
	require.NoError(t, err)

	tmpDir, cleanup, err := dirs.TempDir(".", "")
	defer cleanup()
	require.NoError(t, err)

	wd, err := os.Getwd()
	require.NoError(t, err)

	for i, tc := range []struct {
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
				"hello",
				"--",
				"--foo-arg",
				"flag:",
				"arg3",
			},
			wantStdout: "[--foo-arg flag: arg3]\n",
		},
	} {
		projectDir, err := ioutil.TempDir(tmpDir, "")
		require.NoError(t, err)

		_, err = gofiles.Write(projectDir, tc.filesToCreate)
		require.NoError(t, err, "Case %d", i)

		configFile := path.Join(projectDir, "config.yml")
		err = ioutil.WriteFile(configFile, []byte(tc.config), 0644)
		require.NoError(t, err)

		var output []byte
		func() {
			err := os.Chdir(projectDir)
			defer func() {
				err := os.Chdir(wd)
				require.NoError(t, err)
			}()
			require.NoError(t, err)

			args := []string{"--config", "config.yml", "run"}
			args = append(args, tc.args...)
			cmd := exec.Command(cli, args...)
			output, err = cmd.CombinedOutput()
			require.NoError(t, err, "Case %d: %s\nOutput: %s", i, tc.name, string(output))
		}()

		content := string(output)[strings.Index(string(output), "\n")+1:]
		assert.Equal(t, tc.wantStdout, content, "Case %d: %s", i, tc.name)
	}
}

func TestRunWithStdin(t *testing.T) {
	cli, err := products.Bin("distgo-plugin")
	require.NoError(t, err)

	wd, err := os.Getwd()
	require.NoError(t, err)

	tmpDir, cleanup, err := dirs.TempDir(wd, "")
	defer cleanup()
	require.NoError(t, err)

	projectDir, err := ioutil.TempDir(tmpDir, "")
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
		"hello",
	}

	stdInContent := "output passed to stdin\n"

	_, err = gofiles.Write(projectDir, filesToCreate)
	require.NoError(t, err)

	configFile := path.Join(projectDir, "config.yml")
	err = ioutil.WriteFile(configFile, []byte(config), 0644)
	require.NoError(t, err)

	var output []byte
	func() {
		err := os.Chdir(projectDir)
		defer func() {
			err := os.Chdir(wd)
			require.NoError(t, err)
		}()
		require.NoError(t, err)

		args := []string{"--config", "config.yml", "run"}
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
