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

package plugins

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v2"

	"github.com/palantir/godel/framework/godellauncher"
)

func TestConfigToParams(t *testing.T) {
	cfgContent := `
resolvers:
  - https://localhost:8080/repo/{{GroupPath}}/{{Product}}/{{Version}}/{{Product}}-{{OS}}-{{Arch}}-{{Version}}.tgz
plugins:
  - locator:
      id: "com.palantir:tester:1.0.0"
      checksums:
        darwin-amd64: d22c0ac9d3b65ebe5b830c1324f3d43e777ebc085c580af7c39fb1e5e3c909a7
`
	var cfg godellauncher.PluginsConfig
	err := yaml.Unmarshal([]byte(cfgContent), &cfg)
	require.NoError(t, err)
	_, err = projectParamsFromConfig(cfg)
	require.NoError(t, err)
}

func TestConfigToParamsInvalidLocator(t *testing.T) {
	cfgContent := `
plugins:
  - locator:
      id: "tester:1.0.0"
`
	var cfg godellauncher.PluginsConfig
	err := yaml.Unmarshal([]byte(cfgContent), &cfg)
	require.NoError(t, err)
	_, err = projectParamsFromConfig(cfg)
	assert.EqualError(t, err, `invalid locator: locator ID must consist of 3 colon-delimited components ([group]:[product]:[version]), but had 2: "tester:1.0.0"`)
}
