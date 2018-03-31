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

package config_test

import (
	"testing"

	"github.com/palantir/godel/framework/pluginapitester"
)

func TestUpgradeConfig(t *testing.T) {
	pluginapitester.RunUpgradeConfigTest(t,
		nil,
		nil,
		[]pluginapitester.UpgradeConfigTestCase{
			{
				Name: "legacy config is not upgraded if exclude.yml not present",
				ConfigFiles: map[string]string{
					"godel/config/godel.yml": `
exclude:
  # comment
  names:
    - "\\..+"
    - "vendor"
  paths:
    - "godel"
`,
				},
				Legacy:     true,
				WantOutput: "",
				WantFiles: map[string]string{
					"godel/config/godel.yml": `
exclude:
  # comment
  names:
    - "\\..+"
    - "vendor"
  paths:
    - "godel"
`,
				},
			},
			{
				Name: "legacy config is upgraded if exclude.yml is present",
				ConfigFiles: map[string]string{
					"godel/config/godel.yml": `
exclude:
  # comment
  names:
    - "\\..+"
    - "vendor"
  paths:
    - "godel"
`,
					"godel/config/exclude.yml": `
names:
  - "\\..+"
  - "mocks"
  - "vendor"
  - ".*\\.pb\\.go"
paths:
  - "godel"
  - "internal/conjure/sls"
`,
				},
				Legacy: true,
				WantOutput: `Upgraded configuration for godel.yml
`,
				WantFiles: map[string]string{
					"godel/config/godel.yml": `exclude:
  names:
  - \..+
  - vendor
  - mocks
  - .*\.pb\.go
  paths:
  - godel
  - internal/conjure/sls
`,
				},
			},
			{
				Name: "current config is unmodified",
				ConfigFiles: map[string]string{
					"godel/config/godel.yml": `
default-tasks:
  resolvers:
    - https://palantir.bintray.com/releases/{{GroupPath}}/{{Product}}/{{Version}}/{{Product}}-{{Version}}-{{OS}}-{{Arch}}.tgz
  tasks:
    # comment in config
    com.palantir.distgo:dist-plugin:
      exclude-all-default-assets: true
    com.palantir.okgo:check-plugin:
      locator:
        id: "com.palantir.okgo:check-plugin:1.0.0-rc4"
      assets:
      - locator:
          id: com.palantir.godel-okgo-asset-nobadfuncs:nobadfuncs-asset:1.0.0-rc2
exclude:
  names:
    - "\\..+"
    - "vendor"
  paths:
    - "godel"
`,
				},
				WantOutput: "",
				WantFiles: map[string]string{
					"godel/config/godel.yml": `
default-tasks:
  resolvers:
    - https://palantir.bintray.com/releases/{{GroupPath}}/{{Product}}/{{Version}}/{{Product}}-{{Version}}-{{OS}}-{{Arch}}.tgz
  tasks:
    # comment in config
    com.palantir.distgo:dist-plugin:
      exclude-all-default-assets: true
    com.palantir.okgo:check-plugin:
      locator:
        id: "com.palantir.okgo:check-plugin:1.0.0-rc4"
      assets:
      - locator:
          id: com.palantir.godel-okgo-asset-nobadfuncs:nobadfuncs-asset:1.0.0-rc2
exclude:
  names:
    - "\\..+"
    - "vendor"
  paths:
    - "godel"
`,
				},
			},
		},
	)
}
