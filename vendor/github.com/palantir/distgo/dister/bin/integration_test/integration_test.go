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

package integration

import (
	"io/ioutil"
	"path"
	"testing"

	"github.com/mholt/archiver"
	"github.com/nmiyake/pkg/gofiles"
	"github.com/palantir/godel/framework/pluginapitester"
	"github.com/palantir/godel/pkg/products"
	"github.com/palantir/pkg/specdir"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/palantir/distgo/dister/distertester"
)

func TestBinDist(t *testing.T) {
	const godelYML = `exclude:
  names:
    - "\\..+"
    - "vendor"
  paths:
    - "godel"
`

	pluginPath, err := products.Bin("dist-plugin")
	require.NoError(t, err)

	distertester.RunAssetDistTest(t,
		pluginapitester.NewPluginProvider(pluginPath),
		nil,
		[]distertester.TestCase{
			{
				Name: "bin creates expected output",
				Specs: []gofiles.GoFileSpec{
					{
						RelPath: "foo/foo.go",
						Src:     `package main; func main() {}`,
					},
				},
				ConfigFiles: map[string]string{
					"godel/config/godel.yml": godelYML,
					"godel/config/dist-plugin.yml": `
products:
  foo:
    build:
      main-pkg: ./foo
      os-archs:
        - os: darwin
          arch: amd64
        - os: linux
          arch: amd64
    dist:
      disters:
        type: bin
`,
				},
				WantOutput: func(projectDir string) string {
					return `Creating distribution for foo at out/dist/foo/1.0.0/bin/foo-1.0.0.tgz
Finished creating bin distribution for foo
`
				},
				Validate: func(projectDir string) {
					// layout for the contents of the directory (working directory and in TGZ)
					wantInnerLayout := specdir.NewLayoutSpec(
						specdir.Dir(specdir.LiteralName("foo-1.0.0"), "",
							specdir.Dir(specdir.LiteralName("bin"), "",
								specdir.Dir(specdir.LiteralName("darwin-amd64"), "",
									specdir.File(specdir.LiteralName("foo"), ""),
								),
								specdir.Dir(specdir.LiteralName("linux-amd64"), "",
									specdir.File(specdir.LiteralName("foo"), ""),
								),
							),
						), true,
					)

					// verify that work directory and output TGZ was created
					wantOuterLayout := specdir.NewLayoutSpec(
						specdir.Dir(specdir.LiteralName("1.0.0"), "",
							specdir.Dir(specdir.LiteralName("bin"), "",
								wantInnerLayout,
								specdir.File(specdir.LiteralName("foo-1.0.0.tgz"), ""),
							),
						), true,
					)
					assert.NoError(t, wantOuterLayout.Validate(path.Join(projectDir, "out", "dist", "foo", "1.0.0"), nil))

					// verify layout of the work directory and TGZ
					verifyLayoutAndTGZ(t, projectDir, wantInnerLayout)
				},
			},
			{
				Name: "bin compresses work directory and includes output created by script",
				Specs: []gofiles.GoFileSpec{
					{
						RelPath: "foo/foo.go",
						Src:     `package main; func main() {}`,
					},
				},
				ConfigFiles: map[string]string{
					"godel/config/godel.yml": godelYML,
					"godel/config/dist-plugin.yml": `
products:
  foo:
    build:
      main-pkg: ./foo
      os-archs:
        - os: darwin
          arch: amd64
        - os: linux
          arch: amd64
    dist:
      disters:
        type: bin
        script: |
                #!/usr/bin/env bash
                # move bin directory into service directory
                mkdir $DIST_WORK_DIR/service
                mv $DIST_WORK_DIR/bin $DIST_WORK_DIR/service/bin
                echo "hello" > $DIST_WORK_DIR/foo.txt
`,
				},
				WantOutput: func(projectDir string) string {
					return `Creating distribution for foo at out/dist/foo/1.0.0/bin/foo-1.0.0.tgz
Finished creating bin distribution for foo
`
				},
				Validate: func(projectDir string) {
					// layout for the contents of the directory (working directory and in TGZ)
					wantInnerLayout := specdir.NewLayoutSpec(
						specdir.Dir(specdir.LiteralName("foo-1.0.0"), "",
							specdir.Dir(specdir.LiteralName("service"), "",
								specdir.Dir(specdir.LiteralName("bin"), "",
									specdir.Dir(specdir.LiteralName("darwin-amd64"), "",
										specdir.File(specdir.LiteralName("foo"), ""),
									),
									specdir.Dir(specdir.LiteralName("linux-amd64"), "",
										specdir.File(specdir.LiteralName("foo"), ""),
									),
								),
							),
							specdir.File(specdir.LiteralName("foo.txt"), ""),
						), true)

					// verify that work directory and output TGZ was created
					wantOuterLayout := specdir.NewLayoutSpec(
						specdir.Dir(specdir.LiteralName("1.0.0"), "",
							specdir.Dir(specdir.LiteralName("bin"), "",
								wantInnerLayout,
								specdir.File(specdir.LiteralName("foo-1.0.0.tgz"), ""),
							),
						), true,
					)
					assert.NoError(t, wantOuterLayout.Validate(path.Join(projectDir, "out", "dist", "foo", "1.0.0"), nil))

					// verify layout of the work directory and TGZ
					verifyLayoutAndTGZ(t, projectDir, wantInnerLayout)
				},
			},
			{
				Name: "bin is able to create a valid TGZ archive containing long paths",
				Specs: []gofiles.GoFileSpec{
					{
						RelPath: "foo/foo.go",
						Src:     `package main; func main() {}`,
					},
				},
				ConfigFiles: map[string]string{
					"godel/config/godel.yml": godelYML,
					"godel/config/dist-plugin.yml": `
products:
  foo:
    build:
      main-pkg: ./foo
      os-archs:
        - os: darwin
          arch: amd64
        - os: linux
          arch: amd64
    dist:
      disters:
        type: bin
        script: |
                #!/usr/bin/env bash
                mkdir -p $DIST_WORK_DIR/0/1/2/3/4/5/6/7/8/9/10/11/12/13/14/15/16/17/18/19/20/21/22/23/24/25/26/27/28/29/30/31/32/33/
                touch $DIST_WORK_DIR/0/1/2/3/4/5/6/7/8/9/10/11/12/13/14/15/16/17/18/19/20/21/22/23/24/25/26/27/28/29/30/31/32/33/file.txt
`,
				},
				WantOutput: func(projectDir string) string {
					return `Creating distribution for foo at out/dist/foo/1.0.0/bin/foo-1.0.0.tgz
Finished creating bin distribution for foo
`
				},
				Validate: func(projectDir string) {
					// layout for the contents of the directory (working directory and in TGZ)
					wantInnerLayout := specdir.NewLayoutSpec(
						specdir.Dir(specdir.LiteralName("foo-1.0.0"), "",
							specdir.Dir(specdir.LiteralName("bin"), "",
								specdir.Dir(specdir.LiteralName("darwin-amd64"), "",
									specdir.File(specdir.LiteralName("foo"), ""),
								),
								specdir.Dir(specdir.LiteralName("linux-amd64"), "",
									specdir.File(specdir.LiteralName("foo"), ""),
								),
							),
							specdir.Dir(specdir.LiteralName("0"), "",
								specdir.Dir(specdir.LiteralName("1"), "",
									specdir.Dir(specdir.LiteralName("2"), "",
										specdir.Dir(specdir.LiteralName("3"), "",
											specdir.Dir(specdir.LiteralName("4"), "",
												specdir.Dir(specdir.LiteralName("5"), "",
													specdir.Dir(specdir.LiteralName("6"), "",
														specdir.Dir(specdir.LiteralName("7"), "",
															specdir.Dir(specdir.LiteralName("8"), "",
																specdir.Dir(specdir.LiteralName("9"), "",
																	specdir.Dir(specdir.LiteralName("10"), "",
																		specdir.Dir(specdir.LiteralName("11"), "",
																			specdir.Dir(specdir.LiteralName("12"), "",
																				specdir.Dir(specdir.LiteralName("13"), "",
																					specdir.Dir(specdir.LiteralName("14"), "",
																						specdir.Dir(specdir.LiteralName("15"), "",
																							specdir.Dir(specdir.LiteralName("16"), "",
																								specdir.Dir(specdir.LiteralName("17"), "",
																									specdir.Dir(specdir.LiteralName("18"), "",
																										specdir.Dir(specdir.LiteralName("19"), "",
																											specdir.Dir(specdir.LiteralName("20"), "",
																												specdir.Dir(specdir.LiteralName("21"), "",
																													specdir.Dir(specdir.LiteralName("22"), "",
																														specdir.Dir(specdir.LiteralName("23"), "",
																															specdir.Dir(specdir.LiteralName("24"), "",
																																specdir.Dir(specdir.LiteralName("25"), "",
																																	specdir.Dir(specdir.LiteralName("26"), "",
																																		specdir.Dir(specdir.LiteralName("27"), "",
																																			specdir.Dir(specdir.LiteralName("28"), "",
																																				specdir.Dir(specdir.LiteralName("29"), "",
																																					specdir.Dir(specdir.LiteralName("30"), "",
																																						specdir.Dir(specdir.LiteralName("31"), "",
																																							specdir.Dir(specdir.LiteralName("32"), "",
																																								specdir.Dir(specdir.LiteralName("33"), "",
																																									specdir.File(specdir.LiteralName("file.txt"), ""),
																																								),
																																							),
																																						),
																																					),
																																				),
																																			),
																																		),
																																	),
																																),
																															),
																														),
																													),
																												),
																											),
																										),
																									),
																								),
																							),
																						),
																					),
																				),
																			),
																		),
																	),
																),
															),
														),
													),
												),
											),
										),
									),
								),
							),
						), true)

					// verify that work directory and output TGZ was created
					wantOuterLayout := specdir.NewLayoutSpec(
						specdir.Dir(specdir.LiteralName("1.0.0"), "",
							specdir.Dir(specdir.LiteralName("bin"), "",
								wantInnerLayout,
								specdir.File(specdir.LiteralName("foo-1.0.0.tgz"), ""),
							),
						), true,
					)
					assert.NoError(t, wantOuterLayout.Validate(path.Join(projectDir, "out", "dist", "foo", "1.0.0"), nil))

					// verify layout of the work directory and TGZ
					verifyLayoutAndTGZ(t, projectDir, wantInnerLayout)
				},
			},
		},
	)
}

func verifyLayoutAndTGZ(t *testing.T, projectDir string, wantLayout specdir.LayoutSpec) {
	// validate directory layout
	assert.NoError(t, wantLayout.Validate(path.Join(projectDir, "out", "dist", "foo", "1.0.0", "bin", "foo-1.0.0"), nil))

	// expand tgz and validate directory layout of expanded tgz
	tmpDir, err := ioutil.TempDir(projectDir, "expanded")
	require.NoError(t, err)
	require.NoError(t, archiver.TarGz.Open(path.Join(projectDir, "out", "dist", "foo", "1.0.0", "bin", "foo-1.0.0.tgz"), tmpDir))
	assert.NoError(t, wantLayout.Validate(path.Join(tmpDir, "foo-1.0.0"), nil))
}

func TestBinUpgradeConfig(t *testing.T) {
	pluginPath, err := products.Bin("dist-plugin")
	require.NoError(t, err)

	pluginapitester.RunUpgradeConfigTest(t,
		pluginapitester.NewPluginProvider(pluginPath),
		nil,
		[]pluginapitester.UpgradeConfigTestCase{
			{
				Name: `legacy configuration is upgraded`,
				ConfigFiles: map[string]string{
					"godel/config/dist.yml": `products:
  foo:
    build:
      main-pkg: ./foo
      os-archs:
        - os: darwin
          arch: amd64
        - os: linux
          arch: amd64
    dist:
      dist-type:
        type: bin
        info:
          omit-init-sh: false
`,
				},
				Legacy: true,
				WantOutput: `Upgraded configuration for dist-plugin.yml
`,
				WantFiles: map[string]string{
					"godel/config/dist-plugin.yml": `products:
  foo:
    build:
      name-template: null
      output-dir: null
      main-pkg: ./foo
      build-args-script: null
      version-var: null
      environment: null
      os-archs:
      - os: darwin
        arch: amd64
      - os: linux
        arch: amd64
    run: null
    dist:
      output-dir: null
      disters:
        bin:
          type: bin
          config: null
          name-template: null
          script: |
            #!/bin/bash
            ### START: auto-generated back-compat code for "omit-init-sh: false" behavior for bin dist ###
            read -d '' GODELUPGRADED_scriptContent <<"EOF"
            #!/bin/bash
            set -euo pipefail
            BIN_DIR="$(cd "$(dirname "$0")" && pwd)"
            # determine OS
            OS=""
            case "$(uname)" in
              Darwin*)
                OS=darwin
                ;;
              Linux*)
                OS=linux
                ;;
              *)
                echo "Unsupported operating system: $(uname)"
                exit 1
                ;;
            esac
            # determine executable location based on OS
            CMD=$BIN_DIR/$OS-amd64/{{.ProductName}}
            # verify that executable exists
            if [ ! -e "$CMD" ]; then
                echo "Executable $CMD does not exist"
                exit 1
            fi
            # invoke appropriate executable
            $CMD "$@"
            EOF
            GODELUPGRADED_templated=${GODELUPGRADED_scriptContent//\{\{.ProductName\}\}/$PRODUCT}
            echo "$GODELUPGRADED_templated" > "$DIST_WORK_DIR"/bin/"$PRODUCT".sh
            chmod 755 "$DIST_WORK_DIR"/bin/"$PRODUCT".sh
            ### END: auto-generated back-compat code for "omit-init-sh: false" behavior for bin dist ###
    publish: null
    docker: null
    dependencies: null
product-defaults:
  build: null
  run: null
  dist: null
  publish: null
  docker: null
  dependencies: null
script-includes: ""
exclude:
  names: []
  paths: []
`,
				},
			},
			{
				Name: `valid v0 config works`,
				ConfigFiles: map[string]string{
					"godel/config/dist-plugin.yml": `
products:
  foo:
    build:
      main-pkg: ./foo
      os-archs:
        - os: darwin
          arch: amd64
        - os: linux
          arch: amd64
    dist:
      disters:
        type: bin
        config:
          # comment
`,
				},
				WantOutput: ``,
				WantFiles: map[string]string{
					"godel/config/dist-plugin.yml": `
products:
  foo:
    build:
      main-pkg: ./foo
      os-archs:
        - os: darwin
          arch: amd64
        - os: linux
          arch: amd64
    dist:
      disters:
        type: bin
        config:
          # comment
`,
				},
			},
		},
	)
}
