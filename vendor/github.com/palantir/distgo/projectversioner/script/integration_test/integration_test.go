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
	"regexp"
	"testing"

	"github.com/palantir/godel/framework/pluginapitester"
	"github.com/palantir/godel/pkg/products"
	"github.com/stretchr/testify/require"

	"github.com/palantir/distgo/projectversioner/projectversiontester"
)

func TestScriptProjectVersioner(t *testing.T) {
	const godelYML = `exclude:
  names:
    - "\\..+"
    - "vendor"
  paths:
    - "godel"
`

	pluginPath, err := products.Bin("dist-plugin")
	require.NoError(t, err)

	projectversiontester.RunAssetProjectVersionTest(t,
		pluginapitester.NewPluginProvider(pluginPath),
		nil,
		[]projectversiontester.TestCase{
			{
				Name: "version of project is output of script",
				ConfigFiles: map[string]string{
					"godel/config/godel.yml": godelYML,
					"godel/config/dist-plugin.yml": `
project-versioner:
  type: script
  config:
    # comment
    script: |
            #!/usr/bin/env bash
            echo "1.0.0"
`,
				},
				WantOutput: func(projectDir string) *regexp.Regexp {
					return regexp.MustCompile("^" + regexp.QuoteMeta("1.0.0") + "\n$")
				},
			},
			{
				Name: "project directory is available as environment variable",
				ConfigFiles: map[string]string{
					"godel/config/godel.yml": godelYML,
					"godel/config/dist-plugin.yml": `
project-versioner:
  type: script
  config:
    # comment
    script: |
            #!/usr/bin/env bash
            echo "$PROJECT_DIR"
`,
				},
				WantOutput: func(projectDir string) *regexp.Regexp {
					return regexp.MustCompile("^" + regexp.QuoteMeta(projectDir) + "\n$")
				},
			},
		},
	)
}

func TestScriptUpgradeConfig(t *testing.T) {
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
project-versioner:
  type: script
  config:
    # comment
    script: |
            #!/usr/bin/env bash
            echo "1.0.0"
`,
				},
				WantOutput: ``,
				WantFiles: map[string]string{
					"godel/config/dist-plugin.yml": `
project-versioner:
  type: script
  config:
    # comment
    script: |
            #!/usr/bin/env bash
            echo "1.0.0"
`,
				},
			},
		},
	)
}
