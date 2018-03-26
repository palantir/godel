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

package pluginapi_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path"
	"testing"

	"github.com/nmiyake/pkg/dirs"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/palantir/godel/framework/builtintasks"
	"github.com/palantir/godel/framework/godellauncher"
	"github.com/palantir/godel/framework/pluginapi/v2/pluginapi"
)

var echoPluginTmpl = fmt.Sprintf(`#!/bin/sh
if [ "$1" = "%s" ]; then
    echo '%s'
    exit 0
fi

echo $@
`, pluginapi.PluginInfoCommandName, `%s`)

func TestNewPluginInfo(t *testing.T) {
	info, err := pluginapi.NewPluginInfo("group", "product-plugin", "1.0.0",
		pluginapi.PluginInfoUsesConfigFile(),
	)
	require.NoError(t, err)

	assert.Equal(t, pluginapi.CurrentSchemaVersion, info.PluginSchemaVersion())
	assert.Equal(t, "group", info.Group())
	assert.Equal(t, "product-plugin", info.Product())
	assert.Equal(t, "1.0.0", info.Version())
	assert.Nil(t, info.Tasks("", nil))
	assert.Nil(t, info.UpgradeConfigTask("", nil))
}

func TestPluginInfoJSONMarshal(t *testing.T) {
	for i, tc := range []struct {
		group, product, version string
		params                  []pluginapi.PluginInfoParam
		want                    string
	}{
		{
			"group", "product-plugin", "1.0.0",
			[]pluginapi.PluginInfoParam{
				pluginapi.PluginInfoUsesConfigFile(),
			},
			`{"pluginSchemaVersion":"2","group":"group","product":"product-plugin","version":"1.0.0","usesConfig":true,"tasks":null,"upgradeTask":null}`,
		},
		{
			"group", "product-plugin", "1.0.0",
			[]pluginapi.PluginInfoParam{
				pluginapi.PluginInfoUsesConfigFile(),
				pluginapi.PluginInfoTaskInfo("foo", "does foo things"),
			},
			`{"pluginSchemaVersion":"2","group":"group","product":"product-plugin","version":"1.0.0","usesConfig":true,"tasks":[{"name":"foo","description":"does foo things","command":null,"globalFlagOptions":null,"verifyOptions":null}],"upgradeTask":null}`,
		},
		{
			"group", "product-plugin", "1.0.0",
			[]pluginapi.PluginInfoParam{
				pluginapi.PluginInfoUsesConfigFile(),
				pluginapi.PluginInfoTaskInfo("foo", "does foo things",
					pluginapi.TaskInfoCommand("foo"),
					pluginapi.TaskInfoVerifyOptions(),
				),
			},
			`{"pluginSchemaVersion":"2","group":"group","product":"product-plugin","version":"1.0.0","usesConfig":true,"tasks":[{"name":"foo","description":"does foo things","command":["foo"],"globalFlagOptions":null,"verifyOptions":{"verifyTaskFlags":null,"ordering":null,"applyTrueArgs":null,"applyFalseArgs":null}}],"upgradeTask":null}`,
		},
		{
			"group", "product-plugin", "1.0.0",
			[]pluginapi.PluginInfoParam{
				pluginapi.PluginInfoUsesConfigFile(),
				pluginapi.PluginInfoGlobalFlagOptions(
					pluginapi.GlobalFlagOptionsParamProjectDirFlag("--project-dir"),
				),
				pluginapi.PluginInfoTaskInfo("foo", "does foo things",
					pluginapi.TaskInfoCommand("foo"),
					pluginapi.TaskInfoVerifyOptions(),
				),
			},
			`{"pluginSchemaVersion":"2","group":"group","product":"product-plugin","version":"1.0.0","usesConfig":true,"tasks":[{"name":"foo","description":"does foo things","command":["foo"],"globalFlagOptions":{"debugFlag":"","projectDirFlag":"--project-dir","godelConfigFlag":"","configFlag":""},"verifyOptions":{"verifyTaskFlags":null,"ordering":null,"applyTrueArgs":null,"applyFalseArgs":null}}],"upgradeTask":null}`,
		},
	} {
		info, err := pluginapi.NewPluginInfo(tc.group, tc.product, tc.version, tc.params...)
		require.NoError(t, err, "Case %d", i)

		bytes, err := json.Marshal(info)
		require.NoError(t, err, "Case %d", i)

		assert.Equal(t, tc.want, string(bytes), "Case %d\nGot:\n%s", i, string(bytes))
	}
}

func TestNewPluginInfoError(t *testing.T) {
	for i, tc := range []struct {
		name                    string
		group, product, version string
		params                  []pluginapi.PluginInfoParam
		wantError               string
	}{
		{
			"plugins cannot provide multiple tasks with the same name",
			"group", "product-plugin", "1.0.0",
			[]pluginapi.PluginInfoParam{
				pluginapi.PluginInfoTaskInfo("name", "description"),
				pluginapi.PluginInfoTaskInfo("name", "description-2"),
			},
			`plugin group:product-plugin:1.0.0 specifies multiple tasks with name "name"`,
		},
		{
			"plugin cannot provide upgrade task if it does not use configuration",
			"group", "product-plugin", "1.0.0",
			[]pluginapi.PluginInfoParam{
				pluginapi.PluginInfoUpgradeConfigTaskInfo(pluginapi.UpgradeConfigTaskInfoCommand("upgrade-task")),
			},
			`plugin group:product-plugin:1.0.0 provides a configuration upgrade task but does not specify that it uses configuration`,
		},
	} {
		_, err := pluginapi.NewPluginInfo(tc.group, tc.product, tc.version, tc.params...)
		require.Error(t, err, "Case %d: %s", i, tc.name)
		assert.EqualError(t, err, tc.wantError, "Case %d: %s", i, tc.name)
	}
}

func TestRunPluginFromInfo(t *testing.T) {
	tmpDir, cleanup, err := dirs.TempDir("", "")
	require.NoError(t, err)
	defer cleanup()

	for i, tc := range []struct {
		name         string
		params       []pluginapi.PluginInfoParam
		globalConfig godellauncher.GlobalConfig
		numAssets    int
		want         func(assetDir string) string
	}{
		{
			"no global flags",
			nil,
			godellauncher.GlobalConfig{
				TaskArgs: []string{"--echo-bool-flag", "-f", "echo-str-flag-val", "echo-arg"},
			},
			0,
			func(assetDir string) string {
				return "--echo-bool-flag -f echo-str-flag-val echo-arg\n"
			},
		},
		{
			"global debug flag support",
			[]pluginapi.PluginInfoParam{
				pluginapi.PluginInfoGlobalFlagOptions(
					pluginapi.GlobalFlagOptionsParamDebugFlag("--debug-flag"),
				),
			},
			godellauncher.GlobalConfig{
				TaskArgs: []string{"--echo-bool-flag", "-f", "echo-str-flag-val", "echo-arg"},
				Debug:    true,
			},
			0,
			func(assetDir string) string {
				return "--debug-flag --echo-bool-flag -f echo-str-flag-val echo-arg\n"
			},
		},
		{
			"all global flag support without project dir configured",
			[]pluginapi.PluginInfoParam{
				pluginapi.PluginInfoGlobalFlagOptions(
					pluginapi.GlobalFlagOptionsParamDebugFlag("--debug-flag"),
					pluginapi.GlobalFlagOptionsParamProjectDirFlag("--project-dir"),
					pluginapi.GlobalFlagOptionsParamGodelConfigFlag("--godel-config"),
					pluginapi.GlobalFlagOptionsParamConfigFlag("--config"),
				),
			},
			godellauncher.GlobalConfig{
				TaskArgs: []string{"--echo-bool-flag", "-f", "echo-str-flag-val", "echo-arg"},
			},
			0,
			func(assetDir string) string {
				return "--echo-bool-flag -f echo-str-flag-val echo-arg\n"
			},
		},
		{
			"all global flag support with project dir configured",
			[]pluginapi.PluginInfoParam{
				pluginapi.PluginInfoGlobalFlagOptions(
					pluginapi.GlobalFlagOptionsParamDebugFlag("--debug-flag"),
					pluginapi.GlobalFlagOptionsParamProjectDirFlag("--project-dir"),
					pluginapi.GlobalFlagOptionsParamGodelConfigFlag("--godel-config"),
					pluginapi.GlobalFlagOptionsParamConfigFlag("--config"),
				),
			},
			godellauncher.GlobalConfig{
				Wrapper:  "../../godelw",
				TaskArgs: []string{"--echo-bool-flag", "-f", "echo-str-flag-val", "echo-arg"},
			},
			0,
			func(assetDir string) string {
				return "--project-dir ../.. --godel-config ../../godel/config/godel.yml --config ../../godel/config/echo-plugin.yml --echo-bool-flag -f echo-str-flag-val echo-arg\n"
			},
		},
		{
			"no global flags with asset",
			nil,
			godellauncher.GlobalConfig{
				TaskArgs: []string{"--echo-bool-flag", "-f", "echo-str-flag-val", "echo-arg"},
			},
			1,
			func(assetDir string) string {
				return fmt.Sprintf("--assets %s/echo-4-asset-0 --echo-bool-flag -f echo-str-flag-val echo-arg\n", assetDir)
			},
		},
		{
			"no global flags with multiple assets",
			nil,
			godellauncher.GlobalConfig{
				TaskArgs: []string{"--echo-bool-flag", "-f", "echo-str-flag-val", "echo-arg"},
			},
			2,
			func(assetDir string) string {
				return fmt.Sprintf("--assets %s/echo-5-asset-0,%s/echo-5-asset-1 --echo-bool-flag -f echo-str-flag-val echo-arg\n", assetDir, assetDir)
			},
		},
	} {
		pluginInfo, err := pluginapi.NewPluginInfo("group", "echo-plugin", "1.0.0",
			append([]pluginapi.PluginInfoParam{
				pluginapi.PluginInfoUsesConfigFile(),
				pluginapi.PluginInfoTaskInfo("echo", "echoes the provided input"),
			}, tc.params...)...,
		)
		require.NoError(t, err, "Case %d: %s", i, tc.name)

		pluginInfoJSON, err := json.Marshal(pluginInfo)
		require.NoError(t, err)
		pluginExecPath := path.Join(tmpDir, fmt.Sprintf("echo-%d.sh", i))

		err = ioutil.WriteFile(pluginExecPath, []byte(fmt.Sprintf(echoPluginTmpl, string(pluginInfoJSON))), 0755)
		require.NoError(t, err)

		var assets []string
		for assetNum := 0; assetNum < tc.numAssets; assetNum++ {
			assetPath := path.Join(tmpDir, fmt.Sprintf("echo-%d-asset-%d", i, assetNum))

			err = ioutil.WriteFile(assetPath, []byte(fmt.Sprintf("asset %d", assetNum)), 0755)
			require.NoError(t, err)

			assets = append(assets, assetPath)
		}

		tasks := pluginInfo.Tasks(pluginExecPath, assets)
		require.Equal(t, 1, len(tasks), "Case %d: %s", i, tc.name)

		outBuf := &bytes.Buffer{}
		gc := tc.globalConfig
		gc.Executable = pluginExecPath

		err = tasks[0].Run(tc.globalConfig, outBuf)
		require.NoError(t, err, "Case %d: %s", i, tc.name)

		assert.Equal(t, tc.want(tmpDir), outBuf.String(), "Case %d: %s", i, tc.name)
	}
}

func TestRunPluginVerify(t *testing.T) {
	tmpDir, cleanup, err := dirs.TempDir("", "")
	require.NoError(t, err)
	defer cleanup()

	for i, tc := range []struct {
		name             string
		pluginInfoParams []pluginapi.PluginInfoParam
		taskInfoParams   []pluginapi.TaskInfoParam
		globalConfig     godellauncher.GlobalConfig
		want             string
	}{
		{
			"verify task",
			nil,
			[]pluginapi.TaskInfoParam{
				pluginapi.TaskInfoCommand("verify-subcmd"),
				pluginapi.TaskInfoVerifyOptions(
					pluginapi.VerifyOptionsApplyFalseArgs("--no-apply"),
					pluginapi.VerifyOptionsApplyTrueArgs("--apply"),
				),
			},
			godellauncher.GlobalConfig{
				TaskArgs: []string{"verify"},
			},
			"Running echo...\nverify-subcmd --apply\n",
		},
		{
			"verify task with apply=false",
			nil,
			[]pluginapi.TaskInfoParam{
				pluginapi.TaskInfoCommand("verify-subcmd"),
				pluginapi.TaskInfoVerifyOptions(
					pluginapi.VerifyOptionsApplyFalseArgs("--no-apply"),
					pluginapi.VerifyOptionsApplyTrueArgs("--apply"),
				),
			},
			godellauncher.GlobalConfig{
				TaskArgs: []string{"verify", "--apply=false"},
			},
			"Running echo...\nverify-subcmd --no-apply\n",
		},
		{
			"verify task with global flag options",
			[]pluginapi.PluginInfoParam{
				pluginapi.PluginInfoGlobalFlagOptions(
					pluginapi.GlobalFlagOptionsParamProjectDirFlag("--project-dir"),
				),
			},
			[]pluginapi.TaskInfoParam{
				pluginapi.TaskInfoCommand("verify-subcmd"),
				pluginapi.TaskInfoVerifyOptions(
					pluginapi.VerifyOptionsApplyFalseArgs("--no-apply"),
					pluginapi.VerifyOptionsApplyTrueArgs("--apply"),
				),
			},
			godellauncher.GlobalConfig{
				TaskArgs: []string{"verify"},
				Wrapper:  "godelw",
			},
			"Running echo...\n--project-dir . verify-subcmd --apply\n",
		},
	} {
		pluginInfo, err := pluginapi.NewPluginInfo("group", "echo-plugin", "1.0.0",
			append([]pluginapi.PluginInfoParam{
				pluginapi.PluginInfoUsesConfigFile(),
				pluginapi.PluginInfoTaskInfo("echo", "echoes the provided input", tc.taskInfoParams...),
			}, tc.pluginInfoParams...)...,
		)
		require.NoError(t, err, "Case %d: %s", i, tc.name)

		pluginInfoJSON, err := json.Marshal(pluginInfo)
		require.NoError(t, err)
		pluginExecPath := path.Join(tmpDir, fmt.Sprintf("echo-%d.sh", i))

		err = ioutil.WriteFile(pluginExecPath, []byte(fmt.Sprintf(echoPluginTmpl, string(pluginInfoJSON))), 0755)
		require.NoError(t, err)

		tasks := pluginInfo.Tasks(pluginExecPath, nil)
		require.Equal(t, 1, len(tasks), "Case %d: %s", i, tc.name)

		outBuf := &bytes.Buffer{}
		gc := tc.globalConfig
		gc.Executable = pluginExecPath

		vTask := builtintasks.VerifyTask(tasks)
		err = vTask.Run(tc.globalConfig, outBuf)
		require.NoError(t, err, "Case %d: %s", i, tc.name)

		assert.Equal(t, tc.want, outBuf.String(), "Case %d: %s", i, tc.name)
	}
}
