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
	"github.com/palantir/godel/framework/pluginapitester"
	"github.com/palantir/godel/pkg/products"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRun(t *testing.T) {
	cli, err := products.Bin("dist-plugin")
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
	cli, err := products.Bin("dist-plugin")
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

func TestUpgradeConfig(t *testing.T) {
	pluginPath, err := products.Bin("dist-plugin")
	require.NoError(t, err)
	pluginProvider := pluginapitester.NewPluginProvider(pluginPath)

	pluginapitester.RunUpgradeConfigTest(t,
		pluginProvider,
		nil,
		[]pluginapitester.UpgradeConfigTestCase{
			{
				Name: "legacy configuration is upgraded",
				ConfigFiles: map[string]string{
					"godel/config/dist.yml": `
products:
  foo:
    build:
      main-pkg: ./foo/main/foo
      output-dir: foo/build/bin
      version-var: github.com/palantir/foo/main.version
      os-archs:
        - os: linux
          arch: amd64
    dist:
      input-dir: foo/dist/input
      output-dir: foo/build/distributions
      input-products:
        - bar
      dist-type:
        type: bin
        info:
          omit-init-sh: true
      script: |
               # move bin directory into service directory
               mkdir $DIST_DIR/service
               mv $DIST_DIR/bin $DIST_DIR/service/bin
    docker:
      - repository: test/foo
        tag: snapshot
        context-dir: foo/dist/docker
        dependencies:
         - product: foo
           type: bin
           target-file: foo-latest.tgz
      - repository: test/foo-other
        tag: snapshot
        context-dir: other/foo/dist/docker
        dependencies:
         - product: foo
           type: bin
           target-file: foo-latest.tgz
  bar:
    build:
      main-pkg: ./bar/main/bar
      output-dir: bar/build/bin
      version-var: github.com/palantir/bar/main.version
      os-archs:
        - os: darwin
          arch: amd64
        - os: linux
          arch: amd64
    dist:
      input-dir: bar/dist/bar
      output-dir: bar/build/distributions
      dist-type:
        type: bin
        info:
          omit-init-sh: true
      script: |
               if [ "$IS_SNAPSHOT" == "1" ]; then
                 echo "snapshot"
               fi
               # move bin directory into service directory
               mv $DIST_DIR/bin/darwin-amd64 $DIST_DIR/service/bin/darwin-amd64
               mv $DIST_DIR/bin/linux-amd64 $DIST_DIR/service/bin/linux-amd64
               rm -rf $DIST_DIR/bin
  baz:
    build:
      main-pkg: ./baz/main/baz
      os-archs:
        - os: darwin
          arch: amd64
        - os: linux
          arch: amd64
group-id: com.palantir.group
`,
				},
				Legacy:     true,
				WantOutput: "Upgraded configuration for dist-plugin.yml\n",
				WantFiles: map[string]string{
					"godel/config/dist-plugin.yml": `products:
  bar:
    build:
      output-dir: bar/build/bin
      main-pkg: ./bar/main/bar
      version-var: github.com/palantir/bar/main.version
      os-archs:
      - os: darwin
        arch: amd64
      - os: linux
        arch: amd64
    dist:
      output-dir: bar/build/distributions
      disters:
        bin:
          type: bin
          script: |
            #!/bin/bash
            ### START: auto-generated back-compat code for "input-dir" behavior ###
            cp -r "$PROJECT_DIR"/bar/dist/bar/. "$DIST_WORK_DIR"
            find "$DIST_WORK_DIR" -type f -name .gitkeep -exec rm '{}' \;
            ### END: auto-generated back-compat code for "input-dir" behavior ###
            ### START: auto-generated back-compat code for "IS_SNAPSHOT" variable ###
            IS_SNAPSHOT=0
            if [[ $VERSION =~ .+g[-+.]?[a-fA-F0-9]{3,}$ ]]; then IS_SNAPSHOT=1; fi
            ### END: auto-generated back-compat code for "IS_SNAPSHOT" variable ###
            if [ "$IS_SNAPSHOT" == "1" ]; then
              echo "snapshot"
            fi
            # move bin directory into service directory
            mv $DIST_WORK_DIR/bin/darwin-amd64 $DIST_WORK_DIR/service/bin/darwin-amd64
            mv $DIST_WORK_DIR/bin/linux-amd64 $DIST_WORK_DIR/service/bin/linux-amd64
            rm -rf $DIST_WORK_DIR/bin
    publish: {}
  baz:
    build:
      main-pkg: ./baz/main/baz
      os-archs:
      - os: darwin
        arch: amd64
      - os: linux
        arch: amd64
    dist:
      disters:
        os-arch-bin:
          type: os-arch-bin
          config:
            os-archs:
            - os: darwin
              arch: amd64
            - os: linux
              arch: amd64
    publish: {}
  foo:
    build:
      output-dir: foo/build/bin
      main-pkg: ./foo/main/foo
      version-var: github.com/palantir/foo/main.version
      os-archs:
      - os: linux
        arch: amd64
    dist:
      output-dir: foo/build/distributions
      disters:
        bin:
          type: bin
          script: |
            #!/bin/bash
            ### START: auto-generated back-compat code for "input-dir" behavior ###
            cp -r "$PROJECT_DIR"/foo/dist/input/. "$DIST_WORK_DIR"
            find "$DIST_WORK_DIR" -type f -name .gitkeep -exec rm '{}' \;
            ### END: auto-generated back-compat code for "input-dir" behavior ###
            # move bin directory into service directory
            mkdir $DIST_WORK_DIR/service
            mv $DIST_WORK_DIR/bin $DIST_WORK_DIR/service/bin
    publish: {}
    docker:
      docker-builders:
        docker-image-0:
          type: default
          context-dir: foo/dist/docker
          input-dists:
          - foo.bin
          tag-templates:
          - '{{Repository}}test/foo:snapshot'
        docker-image-1:
          type: default
          context-dir: other/foo/dist/docker
          input-dists:
          - foo.bin
          tag-templates:
          - '{{Repository}}test/foo-other:snapshot'
    dependencies:
    - bar
product-defaults:
  publish:
    group-id: com.palantir.group
`,
				},
			},
			{
				Name: "legacy configuration dist block is not upgraded if os-archs not specified for build",
				ConfigFiles: map[string]string{
					"godel/config/dist.yml": `
products:
  foo:
    build:
      main-pkg: ./foo/main/foo
      output-dir: foo/build/bin
      version-var: github.com/palantir/foo/main.version
`,
				},
				Legacy:     true,
				WantOutput: "Upgraded configuration for dist-plugin.yml\n",
				WantFiles: map[string]string{
					"godel/config/dist-plugin.yml": `products:
  foo:
    build:
      output-dir: foo/build/bin
      main-pkg: ./foo/main/foo
      version-var: github.com/palantir/foo/main.version
`,
				},
			},
			{
				Name: "legacy configuration with no Docker tag is upgraded",
				ConfigFiles: map[string]string{
					"godel/config/dist.yml": `
products:
  foo:
    build:
      main-pkg: ./foo
      os-archs:
        - os: linux
          arch: amd64
    dist:
      dist-type:
        type: bin
    docker:
    - repository: repo/foo
      context-dir: foo-docker
`,
				},
				Legacy:     true,
				WantOutput: "Upgraded configuration for dist-plugin.yml\n",
				WantFiles: map[string]string{
					"godel/config/dist-plugin.yml": `products:
  foo:
    build:
      main-pkg: ./foo
      os-archs:
      - os: linux
        arch: amd64
    dist:
      disters:
        bin:
          type: bin
    docker:
      docker-builders:
        docker-image-0:
          type: default
          context-dir: foo-docker
          tag-templates:
          - '{{Repository}}repo/foo:{{Version}}'
`,
				},
			},
			{
				Name: "valid v0 configuration is not modified",
				ConfigFiles: map[string]string{
					"godel/config/dist-plugin.yml": `
products:
  # comment
  test:
    build:
      main-pkg: ./cmd/test
      output-dir: build
      build-args-script: |
                         YEAR=$(date +%Y)
                         echo "-ldflags"
                         echo "-X"
                         echo "main.year=$YEAR"
      version-var: main.version
      environment:
        foo: bar
        baz: 1
        bool: TRUE
      os-archs:
        - os: "darwin"
          arch: "amd64"
        - os: "linux"
          arch: "amd64"
    dist:
      output-dir: dist
      disters:
        type: bin
    publish:
      group-id: com.test.foo
      info:
        bintray:
          config:
            username: username
            password: password
script-includes: |
                 #!/usr/bin/env bash
exclude:
  names:
    - ".*test"
  paths:
    - "vendor"
`,
				},
				WantOutput: ``,
				WantFiles: map[string]string{
					"godel/config/dist-plugin.yml": `
products:
  # comment
  test:
    build:
      main-pkg: ./cmd/test
      output-dir: build
      build-args-script: |
                         YEAR=$(date +%Y)
                         echo "-ldflags"
                         echo "-X"
                         echo "main.year=$YEAR"
      version-var: main.version
      environment:
        foo: bar
        baz: 1
        bool: TRUE
      os-archs:
        - os: "darwin"
          arch: "amd64"
        - os: "linux"
          arch: "amd64"
    dist:
      output-dir: dist
      disters:
        type: bin
    publish:
      group-id: com.test.foo
      info:
        bintray:
          config:
            username: username
            password: password
script-includes: |
                 #!/usr/bin/env bash
exclude:
  names:
    - ".*test"
  paths:
    - "vendor"
`,
				},
			},
		},
	)
}
