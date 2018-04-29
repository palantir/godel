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
	"testing"

	"github.com/nmiyake/pkg/gofiles"
	"github.com/palantir/godel/framework/pluginapitester"
	"github.com/palantir/godel/pkg/osarch"
	"github.com/palantir/godel/pkg/products"
	"github.com/stretchr/testify/require"

	"github.com/palantir/distgo/publisher/publishertester"
)

func TestGitHubPublish(t *testing.T) {
	const godelYML = `exclude:
  names:
    - "\\..+"
    - "vendor"
  paths:
    - "godel"
`

	pluginPath, err := products.Bin("dist-plugin")
	require.NoError(t, err)

	publishertester.RunAssetPublishTest(t,
		pluginapitester.NewPluginProvider(pluginPath),
		nil,
		"github",
		[]publishertester.TestCase{
			{
				Name: "publishes artifact to GitHub",
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
    dist:
      disters:
        type: os-arch-bin
    publish:
      group-id: com.test.group
      info:
        github:
          config:
            api-url: http://github.domain.com
            user: testUser
            token: testToken
            owner: testOwner
            repository: testRepo
`,
				},
				Args: []string{
					"--dry-run",
				},
				WantOutput: func(projectDir string) string {
					return fmt.Sprintf(`[DRY RUN] Creating GitHub release 1.0.0 for testOwner/testRepo...done
[DRY RUN] Uploading %s/out/dist/foo/1.0.0/os-arch-bin/foo-1.0.0-%s.tgz to GitHub (destination URL cannot be computed in dry run)
`, projectDir, osarch.Current().String())
				},
			},
		},
	)
}

func TestGitHubUpgradeConfig(t *testing.T) {
	pluginPath, err := products.Bin("dist-plugin")
	require.NoError(t, err)

	pluginapitester.RunUpgradeConfigTest(t,
		pluginapitester.NewPluginProvider(pluginPath),
		nil,
		[]pluginapitester.UpgradeConfigTestCase{
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
            - os: darwin
              arch: amd64
            - os: linux
              arch: amd64
    publish:
      group-id: com.test.group
      info:
        github:
          config:
            api-url: http://github.domain.com
            user: testUser
            token: testToken
            # comment
            owner: testOwner
            repository: testRepo
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
            - os: darwin
              arch: amd64
            - os: linux
              arch: amd64
    publish:
      group-id: com.test.group
      info:
        github:
          config:
            api-url: http://github.domain.com
            user: testUser
            token: testToken
            # comment
            owner: testOwner
            repository: testRepo
`,
				},
			},
		},
	)
}
