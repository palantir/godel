// Copyright 2019 Palantir Technologies, Inc.
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
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"path"
	"regexp"
	"testing"
	"time"

	"github.com/mholt/archiver"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v2"

	"github.com/palantir/godel/v2/framework/godel/config"
	"github.com/palantir/godel/v2/framework/godellauncher"
	"github.com/palantir/godel/v2/framework/pluginapi/v2/pluginapi"
	"github.com/palantir/godel/v2/pkg/osarch"
)

// TestEnvironmentConfig verifies that the "environment" value in the godel configuration is used.
func TestEnvironmentConfig(t *testing.T) {
	pluginName := fmt.Sprintf("test-env-var-integration-%d-%d-plugin", time.Now().Unix(), rand.Int())

	testProjectDir := setUpGodelTestAndDownload(t, testRootDir, godelTGZ, version)
	src := `package main

import "fmt"

func main() {
	fmt.Println("hello, world!")
}
`
	err := ioutil.WriteFile(path.Join(testProjectDir, "main.go"), []byte(src), 0644)
	require.NoError(t, err)

	cfg, err := config.ReadGodelConfigFromProjectDir(testProjectDir)
	require.NoError(t, err)

	cfgContent := fmt.Sprintf(`
environment:
  GODEL_TEST_ENV_VAR: test-var-content
plugins:
  resolvers:
    - %s/repo/{{GroupPath}}/{{Product}}/{{Version}}/{{Product}}-{{OS}}-{{Arch}}-{{Version}}.tgz
  plugins:
    - locator:
        id: "com.palantir:%s:1.0.0"
`, testProjectDir, pluginName)
	err = yaml.Unmarshal([]byte(cfgContent), &cfg)
	require.NoError(t, err)

	pluginDir := path.Join(testProjectDir, "repo", "com", "palantir", pluginName, "1.0.0")
	err = os.MkdirAll(pluginDir, 0755)
	require.NoError(t, err)

	pluginInfo := pluginapi.MustNewPluginInfo("com.palantir", pluginName, "1.0.0",
		pluginapi.PluginInfoTaskInfo(
			"env-var",
			"Prints value of the GODEL_TEST_ENV_VAR variable",
			pluginapi.TaskInfoCommand("env-var"),
		),
	)
	pluginInfoJSON, err := json.Marshal(pluginInfo)
	require.NoError(t, err)

	pluginScript := path.Join(pluginDir, pluginName+"-1.0.0")
	err = ioutil.WriteFile(pluginScript, []byte(fmt.Sprintf(fmt.Sprintf(`#!/bin/sh
if [ "$1" = "%s" ]; then
    echo '%s'
    exit 0
fi

echo ${GODEL_TEST_ENV_VAR}
`, pluginapi.PluginInfoCommandName, `%s`), string(pluginInfoJSON))), 0755)
	require.NoError(t, err)

	pluginTGZPath := path.Join(pluginDir, fmt.Sprintf("%s-%s-1.0.0.tgz", pluginName, osarch.Current()))
	err = archiver.TarGz.Make(pluginTGZPath, []string{pluginScript})
	require.NoError(t, err)

	cfgBytes, err := yaml.Marshal(cfg)
	require.NoError(t, err)
	cfgDir, err := godellauncher.ConfigDirPath(testProjectDir)
	require.NoError(t, err)
	err = ioutil.WriteFile(path.Join(cfgDir, godellauncher.GodelConfigYML), cfgBytes, 0644)
	require.NoError(t, err)

	// plugin is resolved on first run
	gotOutput := execCommand(t, testProjectDir, "./godelw", "version")
	wantOutput := "(?s)" + regexp.QuoteMeta(fmt.Sprintf(`Getting package from %s/repo/com/palantir/%s/1.0.0/%s-%s-1.0.0.tgz...`, testProjectDir, pluginName, pluginName, osarch.Current())) + ".+"
	assert.Regexp(t, wantOutput, gotOutput)

	gotOutput = execCommand(t, testProjectDir, "./godelw", "env-var")
	assert.Equal(t, "test-var-content\n", gotOutput)
}
