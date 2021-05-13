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
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"regexp"
	"testing"
	"time"

	"github.com/mholt/archiver/v3"
	"github.com/palantir/godel/v2/framework/builtintasks/installupdate/layout"
	"github.com/palantir/godel/v2/framework/godel/config"
	"github.com/palantir/godel/v2/framework/godellauncher"
	"github.com/palantir/godel/v2/framework/pluginapi/v2/pluginapi"
	"github.com/palantir/godel/v2/pkg/osarch"
	"github.com/palantir/pkg/specdir"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v2"
)

var echoPluginTmpl = fmt.Sprintf(`#!/bin/sh
if [ "$1" = "%s" ]; then
    echo '%s'
    exit 0
fi

echo $@
`, pluginapi.PluginInfoCommandName, `%s`)

func TestPlugins(t *testing.T) {
	pluginName := fmt.Sprintf("tester-integration-%d-%d-plugin", time.Now().Unix(), rand.Int())

	testProjectDir := setUpGodelTestAndDownload(t, testRootDir, godelTGZ, version)
	writeMainFile(t, testProjectDir)

	cfg, err := config.ReadGodelConfigFromProjectDir(testProjectDir)
	require.NoError(t, err)

	cfgContent := fmt.Sprintf(`
plugins:
  resolvers:
    - %s/repo/{{GroupPath}}/{{Product}}/{{Version}}/{{Product}}-{{OS}}-{{Arch}}-{{Version}}.tgz
  plugins:
    - locator:
        id: "com.palantir:%s:1.0.0"
`, testProjectDir, pluginName)
	err = yaml.Unmarshal([]byte(cfgContent), &cfg)
	require.NoError(t, err)

	pluginDir := filepath.Join(testProjectDir, "repo", "com", "palantir", pluginName, "1.0.0")
	err = os.MkdirAll(pluginDir, 0755)
	require.NoError(t, err)

	writeDefaultPlugin(t, testProjectDir, pluginName, "1.0.0")

	cfgBytes, err := yaml.Marshal(cfg)
	require.NoError(t, err)
	cfgDir, err := godellauncher.ConfigDirPath(testProjectDir)
	require.NoError(t, err)
	err = os.WriteFile(filepath.Join(cfgDir, godellauncher.GodelConfigYML), cfgBytes, 0644)
	require.NoError(t, err)

	// plugin is resolved on first run
	gotOutput := execCommand(t, testProjectDir, "./godelw", "version")
	wantOutput := "(?s)" + regexp.QuoteMeta(fmt.Sprintf(`Getting package from %s/repo/com/palantir/%s/1.0.0/%s-%s-1.0.0.tgz...`, testProjectDir, pluginName, pluginName, osarch.Current())) + ".+"
	assert.Regexp(t, wantOutput, gotOutput)

	gotOutput = execCommand(t, testProjectDir, "./godelw", "echo-task", "foo", "--bar", "baz")
	wantOutput = fmt.Sprintf("--project-dir %s --godel-config %s/godel/config/godel.yml --config %s/godel/config/%s.yml echo foo --bar baz\n", testProjectDir, testProjectDir, testProjectDir, pluginName)
	assert.Equal(t, wantOutput, gotOutput)

	gotOutput = execCommand(t, testProjectDir, "./godelw", "verify", "--skip-check", "--skip-license", "--skip-test")
	wantOutput = fmt.Sprintf(`Running format...
Running echo-task...
--project-dir %s --godel-config %s/godel/config/godel.yml --config %s/godel/config/%s.yml echo
`, testProjectDir, testProjectDir, testProjectDir, pluginName)
	assert.Equal(t, wantOutput, gotOutput)

	gotOutput = execCommand(t, testProjectDir, "./godelw", "verify", "--skip-check", "--skip-license", "--skip-test", "--apply=false")
	wantOutput = fmt.Sprintf(`Running format...
Running echo-task...
--project-dir %s --godel-config %s/godel/config/godel.yml --config %s/godel/config/%s.yml echo --verify
`, testProjectDir, testProjectDir, testProjectDir, pluginName)
	assert.Equal(t, wantOutput, gotOutput)
}

func TestPluginsWithAssets(t *testing.T) {
	pluginName := fmt.Sprintf("tester-integration-%d-%d-plugin", time.Now().Unix(), rand.Int())
	assetName := pluginName + "-asset"

	testProjectDir := setUpGodelTestAndDownload(t, testRootDir, godelTGZ, version)
	writeMainFile(t, testProjectDir)

	cfg, err := config.ReadGodelConfigFromProjectDir(testProjectDir)
	require.NoError(t, err)

	cfgContent := fmt.Sprintf(`
plugins:
  resolvers:
    - %s/repo/{{GroupPath}}/{{Product}}/{{Version}}/{{Product}}-{{OS}}-{{Arch}}-{{Version}}.tgz
  plugins:
    - locator:
        id: "com.palantir:%s:1.0.0"
      assets:
        - locator:
            id: "com.palantir:%s:1.0.0"
`, testProjectDir, pluginName, assetName)
	err = yaml.Unmarshal([]byte(cfgContent), &cfg)
	require.NoError(t, err)

	pluginDir := filepath.Join(testProjectDir, "repo", "com", "palantir", pluginName, "1.0.0")
	err = os.MkdirAll(pluginDir, 0755)
	require.NoError(t, err)

	assetDir := filepath.Join(testProjectDir, "repo", "com", "palantir", assetName, "1.0.0")
	err = os.MkdirAll(assetDir, 0755)
	require.NoError(t, err)

	writeDefaultPlugin(t, testProjectDir, pluginName, "1.0.0")

	assetFile := filepath.Join(assetDir, assetName+"-1.0.0")
	err = os.WriteFile(assetFile, []byte("asset content"), 0644)
	require.NoError(t, err)

	assetTGZPath := filepath.Join(assetDir, fmt.Sprintf("%s-%s-1.0.0.tgz", assetName, osarch.Current()))
	err = archiver.DefaultTarGz.Archive([]string{assetFile}, assetTGZPath)
	require.NoError(t, err)

	cfgBytes, err := yaml.Marshal(cfg)
	require.NoError(t, err)
	cfgDir, err := godellauncher.ConfigDirPath(testProjectDir)
	require.NoError(t, err)
	err = os.WriteFile(filepath.Join(cfgDir, godellauncher.GodelConfigYML), cfgBytes, 0644)
	require.NoError(t, err)

	// plugin and asset is resolved on first run
	gotOutput := execCommand(t, testProjectDir, "./godelw", "version")
	wantOutput := "(?s)" +
		regexp.QuoteMeta(fmt.Sprintf(`Getting package from %s/repo/com/palantir/%s/1.0.0/%s-%s-1.0.0.tgz...`, testProjectDir, pluginName, pluginName, osarch.Current())) +
		".+" +
		regexp.QuoteMeta(fmt.Sprintf(`Getting package from %s/repo/com/palantir/%s/1.0.0/%s-%s-1.0.0.tgz...`, testProjectDir, assetName, assetName, osarch.Current()))
	assert.Regexp(t, wantOutput, gotOutput)

	godelHomeSpecDir, err := layout.GodelHomeSpecDir(specdir.SpecOnly)
	require.NoError(t, err)
	assetsDir := godelHomeSpecDir.Path(layout.AssetsDir)
	assetPath := filepath.Join(assetsDir, "com.palantir-"+assetName+"-1.0.0")

	gotOutput = execCommand(t, testProjectDir, "./godelw", "echo-task", "foo", "--bar", "baz")
	wantOutput = fmt.Sprintf("--project-dir %s --godel-config %s/godel/config/godel.yml --config %s/godel/config/%s.yml --assets %s echo foo --bar baz\n", testProjectDir, testProjectDir, testProjectDir, pluginName, assetPath)
	assert.Equal(t, wantOutput, gotOutput)

	gotOutput = execCommand(t, testProjectDir, "./godelw", "verify", "--skip-check", "--skip-license", "--skip-test")
	wantOutput = fmt.Sprintf(`Running format...
Running echo-task...
--project-dir %s --godel-config %s/godel/config/godel.yml --config %s/godel/config/%s.yml --assets %s echo
`, testProjectDir, testProjectDir, testProjectDir, pluginName, assetPath)
	assert.Equal(t, wantOutput, gotOutput)

	gotOutput = execCommand(t, testProjectDir, "./godelw", "verify", "--skip-check", "--skip-license", "--skip-test", "--apply=false")
	wantOutput = fmt.Sprintf(`Running format...
Running echo-task...
--project-dir %s --godel-config %s/godel/config/godel.yml --config %s/godel/config/%s.yml --assets %s echo --verify
`, testProjectDir, testProjectDir, testProjectDir, pluginName, assetPath)
	assert.Equal(t, wantOutput, gotOutput)
}

func TestConfigProvider(t *testing.T) {
	pluginName := fmt.Sprintf("tester-integration-%d-%d-plugin", time.Now().Unix(), rand.Int())

	testProjectDir := setUpGodelTestAndDownload(t, testRootDir, godelTGZ, version)
	writeMainFile(t, testProjectDir)

	cfg, err := config.ReadGodelConfigFromProjectDir(testProjectDir)
	require.NoError(t, err)

	cfgProviderContent := fmt.Sprintf(`
plugins:
  resolvers:
    - %s/repo/{{GroupPath}}/{{Product}}/{{Version}}/{{Product}}-{{OS}}-{{Arch}}-{{Version}}.tgz
  plugins:
    - locator:
        id: "com.palantir:%s:1.0.0"
`, testProjectDir, pluginName)

	configProviderName := fmt.Sprintf("tester-integration-config-provider-%d-%d", time.Now().Unix(), rand.Int())
	err = os.MkdirAll(filepath.Join(testProjectDir, "com", "palantir", configProviderName), os.ModePerm)
	assert.NoError(t, err)
	resolverLocation := filepath.Join(testProjectDir, "com", "palantir", configProviderName, "1.0.0.yml")
	err = os.WriteFile(resolverLocation, []byte(cfgProviderContent), 0755)
	assert.NoError(t, err)

	cfgContent := fmt.Sprintf(`
tasks-config-providers:
  resolvers:
    - %s/{{GroupPath}}/{{Product}}/{{Version}}.yml
  providers:
    - locator:
        id: "com.palantir:%s:1.0.0"
`, testProjectDir, configProviderName)
	err = yaml.Unmarshal([]byte(cfgContent), &cfg)
	require.NoError(t, err)

	writeDefaultPlugin(t, testProjectDir, pluginName, "1.0.0")

	cfgBytes, err := yaml.Marshal(cfg)
	require.NoError(t, err)
	cfgDir, err := godellauncher.ConfigDirPath(testProjectDir)
	require.NoError(t, err)
	err = os.WriteFile(filepath.Join(cfgDir, godellauncher.GodelConfigYML), cfgBytes, 0644)
	require.NoError(t, err)

	// plugin is resolved on first run
	gotOutput := execCommand(t, testProjectDir, "./godelw", "version")
	wantOutput := "(?s)" + regexp.QuoteMeta(fmt.Sprintf(`Getting package from %s/repo/com/palantir/%s/1.0.0/%s-%s-1.0.0.tgz...`, testProjectDir, pluginName, pluginName, osarch.Current())) + ".+"
	assert.Regexp(t, wantOutput, gotOutput)

	gotOutput = execCommand(t, testProjectDir, "./godelw", "echo-task", "foo", "--bar", "baz")
	wantOutput = fmt.Sprintf("--project-dir %s --godel-config %s/godel/config/godel.yml --config %s/godel/config/%s.yml echo foo --bar baz\n", testProjectDir, testProjectDir, testProjectDir, pluginName)
	assert.Equal(t, wantOutput, gotOutput)

	gotOutput = execCommand(t, testProjectDir, "./godelw", "verify", "--skip-check", "--skip-license", "--skip-test")
	wantOutput = fmt.Sprintf(`Running format...
Running echo-task...
--project-dir %s --godel-config %s/godel/config/godel.yml --config %s/godel/config/%s.yml echo
`, testProjectDir, testProjectDir, testProjectDir, pluginName)
	assert.Equal(t, wantOutput, gotOutput)

	gotOutput = execCommand(t, testProjectDir, "./godelw", "verify", "--skip-check", "--skip-license", "--skip-test", "--apply=false")
	wantOutput = fmt.Sprintf(`Running format...
Running echo-task...
--project-dir %s --godel-config %s/godel/config/godel.yml --config %s/godel/config/%s.yml echo --verify
`, testProjectDir, testProjectDir, testProjectDir, pluginName)
	assert.Equal(t, wantOutput, gotOutput)
}

func TestConfigProviderCannotSpecifyOverride(t *testing.T) {
	pluginName := fmt.Sprintf("tester-integration-%d-%d-plugin", time.Now().Unix(), rand.Int())

	testProjectDir := setUpGodelTestAndDownload(t, testRootDir, godelTGZ, version)
	writeMainFile(t, testProjectDir)

	cfg, err := config.ReadGodelConfigFromProjectDir(testProjectDir)
	require.NoError(t, err)

	cfgProviderContent := fmt.Sprintf(`
plugins:
  resolvers:
    - %s/repo/{{GroupPath}}/{{Product}}/{{Version}}/{{Product}}-{{OS}}-{{Arch}}-{{Version}}.tgz
  plugins:
    - locator:
        id: "com.palantir:%s:1.0.0"
      override: true
`, testProjectDir, pluginName)

	configProviderName := fmt.Sprintf("tester-integration-config-provider-%d-%d", time.Now().Unix(), rand.Int())
	err = os.MkdirAll(filepath.Join(testProjectDir, "com", "palantir", configProviderName), os.ModePerm)
	assert.NoError(t, err)
	resolverLocation := filepath.Join(testProjectDir, "com", "palantir", configProviderName, "1.0.0.yml")
	err = os.WriteFile(resolverLocation, []byte(cfgProviderContent), 0755)
	assert.NoError(t, err)

	cfgContent := fmt.Sprintf(`
tasks-config-providers:
  resolvers:
    - %s/{{GroupPath}}/{{Product}}/{{Version}}.yml
  providers:
    - locator:
        id: "com.palantir:%s:1.0.0"
`, testProjectDir, configProviderName)
	err = yaml.Unmarshal([]byte(cfgContent), &cfg)
	require.NoError(t, err)

	writeDefaultPlugin(t, testProjectDir, pluginName, "1.0.0")

	cfgBytes, err := yaml.Marshal(cfg)
	require.NoError(t, err)
	cfgDir, err := godellauncher.ConfigDirPath(testProjectDir)
	require.NoError(t, err)
	err = os.WriteFile(filepath.Join(cfgDir, godellauncher.GodelConfigYML), cfgBytes, 0644)
	require.NoError(t, err)

	// configuration provider should fail to load because it has a plugin that specifies an "override" property
	gotOutput := execCommandExpectError(t, testProjectDir, "./godelw", "version")
	wantOutput := "(?s).+" + regexp.QuoteMeta(fmt.Sprintf(`Error: failed to resolve 1 configuration provider(s):`)) + ".+" + regexp.QuoteMeta(`plugins specify override property as 'true', which is not supported in config providers`) + ".+"
	assert.Regexp(t, wantOutput, gotOutput)
}

// TestOverrideResolverPlugin tests that plugins provided by a config provider can be overridden by local plugin
// configuration using an "override" property.
func TestOverrideResolverPlugin(t *testing.T) {
	pluginName := fmt.Sprintf("tester-integration-%d-%d-plugin", time.Now().Unix(), rand.Int())
	testProjectDir := setUpGodelTestAndDownload(t, testRootDir, godelTGZ, version)
	writeMainFile(t, testProjectDir)

	cfg, err := config.ReadGodelConfigFromProjectDir(testProjectDir)
	require.NoError(t, err)

	cfgProviderContent := fmt.Sprintf(`
plugins:
  resolvers:
    - %s/repo/{{GroupPath}}/{{Product}}/{{Version}}/{{Product}}-{{OS}}-{{Arch}}-{{Version}}.tgz
  plugins:
    - locator:
        id: "com.palantir:%s:2.0.0"
`, testProjectDir, pluginName)

	configProviderName := fmt.Sprintf("tester-integration-config-provider-%d-%d", time.Now().Unix(), rand.Int())
	err = os.MkdirAll(filepath.Join(testProjectDir, "com", "palantir", configProviderName), os.ModePerm)
	assert.NoError(t, err)
	resolverLocation := filepath.Join(testProjectDir, "com", "palantir", configProviderName, "1.0.0.yml")
	err = os.WriteFile(resolverLocation, []byte(cfgProviderContent), 0755)
	assert.NoError(t, err)
	cfgContent := fmt.Sprintf(`
tasks-config-providers:
  resolvers:
    - %s/{{GroupPath}}/{{Product}}/{{Version}}.yml
  providers:
    - locator:
        id: "com.palantir:%s:1.0.0"
plugins:
  resolvers:
    - %s/repo/{{GroupPath}}/{{Product}}/{{Version}}/{{Product}}-{{OS}}-{{Arch}}-{{Version}}.tgz
  plugins:
    - locator:
        id: "com.palantir:%s:1.0.0"
      override: true
`, testProjectDir, configProviderName, testProjectDir, pluginName)
	err = yaml.Unmarshal([]byte(cfgContent), &cfg)
	require.NoError(t, err)

	// write version 1.0.0 of plugin
	writePlugin(t, testProjectDir, pluginName, "1.0.0", fmt.Sprintf(`#!/bin/sh
if [ "$1" = "%s" ]; then
    echo '%s'
    exit 0
fi

echo "1.0.0: $@"
`, pluginapi.PluginInfoCommandName, `%s`))

	// write version 2.0.0 of plugin
	writePlugin(t, testProjectDir, pluginName, "2.0.0", fmt.Sprintf(`#!/bin/sh
if [ "$1" = "%s" ]; then
    echo '%s'
    exit 0
fi

echo "2.0.0: $@"
`, pluginapi.PluginInfoCommandName, `%s`))

	writeDefaultPlugin(t, testProjectDir, pluginName, "2.0.0")
	cfgBytes, err := yaml.Marshal(cfg)
	require.NoError(t, err)
	cfgDir, err := godellauncher.ConfigDirPath(testProjectDir)
	require.NoError(t, err)
	err = os.WriteFile(filepath.Join(cfgDir, godellauncher.GodelConfigYML), cfgBytes, 0644)
	require.NoError(t, err)

	// plugin is resolved on first run
	gotOutput := execCommand(t, testProjectDir, "./godelw", "version")
	wantOutput := "(?s)" + regexp.QuoteMeta(fmt.Sprintf(`Getting package from %s/repo/com/palantir/%s/1.0.0/%s-%s-1.0.0.tgz...`, testProjectDir, pluginName, pluginName, osarch.Current())) + ".+"
	assert.Regexp(t, wantOutput, gotOutput)

	// verify that overridden version (version 1.0.0) is used
	gotOutput = execCommand(t, testProjectDir, "./godelw", "echo-task", "foo", "--bar", "baz")
	wantOutput = fmt.Sprintf("1.0.0: --project-dir %s --godel-config %s/godel/config/godel.yml --config %s/godel/config/%s.yml echo foo --bar baz\n", testProjectDir, testProjectDir, testProjectDir, pluginName)
	assert.Equal(t, wantOutput, gotOutput)
}

func writeMainFile(t *testing.T, testProjectDir string) {
	src := `package main

import "fmt"

func main() {
	fmt.Println("hello, world!")
}
`
	err := os.WriteFile(filepath.Join(testProjectDir, "main.go"), []byte(src), 0644)
	require.NoError(t, err)
}

func writeDefaultPlugin(t *testing.T, testProjectDir, pluginName, pluginVersion string) {
	writePlugin(t, testProjectDir, pluginName, pluginVersion, echoPluginTmpl)
}

func writePlugin(t *testing.T, testProjectDir, pluginName, pluginVersion, pluginContent string) {
	pluginDir := filepath.Join(testProjectDir, "repo", "com", "palantir", pluginName, pluginVersion)
	err := os.MkdirAll(pluginDir, 0755)
	require.NoError(t, err)
	pluginInfoJSON := getDefaultPluginInfoJSON(t, pluginName, pluginVersion)

	pluginScript := filepath.Join(pluginDir, fmt.Sprintf("%s-%s", pluginName, pluginVersion))
	err = os.WriteFile(pluginScript, []byte(fmt.Sprintf(pluginContent, string(pluginInfoJSON))), 0755)
	require.NoError(t, err)

	pluginTGZPath := filepath.Join(pluginDir, fmt.Sprintf("%s-%s-%s.tgz", pluginName, osarch.Current(), pluginVersion))
	if _, err := os.Stat(pluginTGZPath); os.IsNotExist(err) {
		err = archiver.DefaultTarGz.Archive([]string{pluginScript}, pluginTGZPath)
		require.NoError(t, err)
	}
}

func getDefaultPluginInfoJSON(t *testing.T, pluginName, pluginVersion string) []byte {
	pluginInfo := pluginapi.MustNewPluginInfo("com.palantir", pluginName, pluginVersion,
		pluginapi.PluginInfoUsesConfigFile(),
		pluginapi.PluginInfoGlobalFlagOptions(
			pluginapi.GlobalFlagOptionsParamDebugFlag("--debug"),
			pluginapi.GlobalFlagOptionsParamProjectDirFlag("--project-dir"),
			pluginapi.GlobalFlagOptionsParamGodelConfigFlag("--godel-config"),
			pluginapi.GlobalFlagOptionsParamConfigFlag("--config"),
		),
		pluginapi.PluginInfoTaskInfo(
			"echo-task",
			"Echoes input",
			pluginapi.TaskInfoCommand("echo"),
			pluginapi.TaskInfoVerifyOptions(
				pluginapi.VerifyOptionsApplyFalseArgs("--verify"),
			),
		),
	)

	pluginInfoJSON, err := json.Marshal(pluginInfo)
	require.NoError(t, err)
	return pluginInfoJSON
}
