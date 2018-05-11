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
	"path"
	"testing"

	"github.com/nmiyake/pkg/gofiles"
	"github.com/palantir/godel/framework/pluginapitester"
	"github.com/palantir/godel/pkg/products"
	"github.com/palantir/pkg/specdir"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/palantir/distgo/dister/distertester"
)

func TestOSArchBinDist(t *testing.T) {
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
				Name: "os-arch-bin creates expected output",
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
        type: os-arch-bin
        config:
          os-archs:
            - os: darwin
              arch: amd64
            - os: linux
              arch: amd64
`,
				},
				WantOutput: func(projectDir string) string {
					return `Creating distribution for foo at out/dist/foo/1.0.0/os-arch-bin/foo-1.0.0-darwin-amd64.tgz, out/dist/foo/1.0.0/os-arch-bin/foo-1.0.0-linux-amd64.tgz
Finished creating os-arch-bin distribution for foo
`
				},
				Validate: func(projectDir string) {
					wantLayout := specdir.NewLayoutSpec(
						specdir.Dir(specdir.LiteralName("1.0.0"), "",
							specdir.Dir(specdir.LiteralName("os-arch-bin"), "",
								specdir.Dir(specdir.LiteralName("foo-1.0.0"), "",
									specdir.Dir(specdir.LiteralName("darwin-amd64"), "",
										specdir.File(specdir.LiteralName("foo"), ""),
									),
									specdir.Dir(specdir.LiteralName("linux-amd64"), "",
										specdir.File(specdir.LiteralName("foo"), ""),
									),
								),
								specdir.File(specdir.LiteralName("foo-1.0.0-darwin-amd64.tgz"), ""),
								specdir.File(specdir.LiteralName("foo-1.0.0-linux-amd64.tgz"), ""),
							),
						), true,
					)
					assert.NoError(t, wantLayout.Validate(path.Join(projectDir, "out", "dist", "foo", "1.0.0"), nil))
				},
			},
		},
	)
}

func TestOSArchBinUpgradeConfig(t *testing.T) {
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
        type: os-arch-bin
        info:
          os-archs:
            - os: darwin
              arch: amd64
            - os: linux
              arch: amd64
`,
				},
				Legacy: true,
				WantOutput: `Upgraded configuration for dist-plugin.yml
`,
				WantFiles: map[string]string{
					"godel/config/dist-plugin.yml": `products:
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
        os-arch-bin:
          type: os-arch-bin
          config:
            os-archs:
            - os: darwin
              arch: amd64
            - os: linux
              arch: amd64
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
        type: os-arch-bin
        config:
          os-archs:
            # comment
            - os: darwin
              arch: amd64
            - os: linux
              arch: amd64
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
        type: os-arch-bin
        config:
          os-archs:
            # comment
            - os: darwin
              arch: amd64
            - os: linux
              arch: amd64
`,
				},
			},
		},
	)
}
