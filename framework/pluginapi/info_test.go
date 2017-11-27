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
	"github.com/palantir/godel/framework/pluginapi"
)

var echoPluginTmpl = fmt.Sprintf(`#!/bin/sh
if [ "$1" = "%s" ]; then
    echo '%s'
    exit 0
fi

echo $@
`, pluginapi.InfoCommandName, `%s`)

func TestNewInfo(t *testing.T) {
	info, err := pluginapi.NewInfo("group", "product", "1.0.0", "foo.yml")
	require.NoError(t, err)

	assert.Equal(t, pluginapi.CurrentSchemaVersion, info.PluginSchemaVersion())
	assert.Equal(t, "group:product:1.0.0", info.ID())
	assert.Equal(t, "foo.yml", info.ConfigFileName())
	assert.Nil(t, info.Tasks("", nil))
}

func TestInfoJSONMarshal(t *testing.T) {
	for i, tc := range []struct {
		group, product, version, configFileName string
		taskInfos                               []pluginapi.TaskInfo
		want                                    string
	}{
		{
			"group", "product", "1.0.0", "foo.yml", nil,
			`{"pluginSchemaVersion":"1","id":"group:product:1.0.0","configFileName":"foo.yml","tasks":null}`,
		},
		{
			"group", "product", "1.0.0", "foo.yml",
			[]pluginapi.TaskInfo{
				pluginapi.MustNewTaskInfo("foo", "does foo things"),
			},
			`{"pluginSchemaVersion":"1","id":"group:product:1.0.0","configFileName":"foo.yml","tasks":[{"name":"foo","description":"does foo things","command":null,"globalFlagOptions":null,"verifyOptions":null}]}`,
		},
		{
			"group", "product", "1.0.0", "foo.yml",
			[]pluginapi.TaskInfo{
				pluginapi.MustNewTaskInfo("foo", "does foo things", pluginapi.TaskInfoCommand("foo"), pluginapi.TaskInfoVerifyOptions(pluginapi.NewVerifyOptions())),
			},
			`{"pluginSchemaVersion":"1","id":"group:product:1.0.0","configFileName":"foo.yml","tasks":[{"name":"foo","description":"does foo things","command":["foo"],"globalFlagOptions":null,"verifyOptions":{"verifyTaskFlags":null,"ordering":null,"applyTrueArgs":null,"applyFalseArgs":null}}]}`,
		},
		{
			"group", "product", "1.0.0", "foo.yml",
			[]pluginapi.TaskInfo{
				pluginapi.MustNewTaskInfo("foo", "does foo things",
					pluginapi.TaskInfoCommand("foo"),
					pluginapi.TaskInfoGlobalFlagOptions(pluginapi.NewGlobalFlagOptions(pluginapi.GlobalFlagOptionsParamProjectDirFlag("--project-dir"))),
					pluginapi.TaskInfoVerifyOptions(pluginapi.NewVerifyOptions()),
				)},
			`{"pluginSchemaVersion":"1","id":"group:product:1.0.0","configFileName":"foo.yml","tasks":[{"name":"foo","description":"does foo things","command":["foo"],"globalFlagOptions":{"debugFlag":"","projectDirFlag":"--project-dir","godelConfigFlag":"","configFlag":""},"verifyOptions":{"verifyTaskFlags":null,"ordering":null,"applyTrueArgs":null,"applyFalseArgs":null}}]}`,
		},
	} {
		info, err := pluginapi.NewInfo(tc.group, tc.product, tc.version, tc.configFileName, tc.taskInfos...)
		require.NoError(t, err, "Case %d", i)

		bytes, err := json.Marshal(info)
		require.NoError(t, err, "Case %d", i)

		assert.Equal(t, tc.want, string(bytes), "Case %d", i)
	}
}

func TestNewInfoError(t *testing.T) {
	ti1, err := pluginapi.NewTaskInfo("name", "description", nil, nil)
	require.NoError(t, err)

	ti2, err := pluginapi.NewTaskInfo("name", "description-2", nil, nil)
	require.NoError(t, err)

	tasks := []pluginapi.TaskInfo{
		ti1,
		ti2,
	}
	_, err = pluginapi.NewInfo("group", "product", "1.0.0", "foo.yml", tasks...)
	require.Error(t, err)
	assert.EqualError(t, err, `plugin group:product:1.0.0 specifies multiple tasks with name "name"`)
}

func TestRunPluginFromInfo(t *testing.T) {
	tmpDir, cleanup, err := dirs.TempDir("", "")
	require.NoError(t, err)
	defer cleanup()

	for i, tc := range []struct {
		name         string
		params       pluginapi.TaskInfoParam
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
			pluginapi.TaskInfoGlobalFlagOptions(pluginapi.NewGlobalFlagOptions(pluginapi.GlobalFlagOptionsParamDebugFlag("--debug-flag"))),
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
			pluginapi.TaskInfoGlobalFlagOptions(pluginapi.NewGlobalFlagOptions(
				pluginapi.GlobalFlagOptionsParamDebugFlag("--debug-flag"),
				pluginapi.GlobalFlagOptionsParamProjectDirFlag("--project-dir"),
				pluginapi.GlobalFlagOptionsParamGodelConfigFlag("--godel-config"),
				pluginapi.GlobalFlagOptionsParamConfigFlag("--config"),
			)),
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
			pluginapi.TaskInfoGlobalFlagOptions(pluginapi.NewGlobalFlagOptions(
				pluginapi.GlobalFlagOptionsParamDebugFlag("--debug-flag"),
				pluginapi.GlobalFlagOptionsParamProjectDirFlag("--project-dir"),
				pluginapi.GlobalFlagOptionsParamGodelConfigFlag("--godel-config"),
				pluginapi.GlobalFlagOptionsParamConfigFlag("--config"),
			)),
			godellauncher.GlobalConfig{
				Wrapper:  "../../godelw",
				TaskArgs: []string{"--echo-bool-flag", "-f", "echo-str-flag-val", "echo-arg"},
			},
			0,
			func(assetDir string) string {
				return "--project-dir ../.. --godel-config ../../godel/config/godel.yml --config ../../godel/config/echo.yml --echo-bool-flag -f echo-str-flag-val echo-arg\n"
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
		pluginInfo, err := pluginapi.NewInfo("group", "echo", "1.0.0", "echo.yml",
			pluginapi.MustNewTaskInfo("echo", "echoes the provided input", tc.params))
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
		name         string
		params       []pluginapi.TaskInfoParam
		globalConfig godellauncher.GlobalConfig
		want         string
	}{
		{
			"verify task",
			[]pluginapi.TaskInfoParam{
				pluginapi.TaskInfoCommand("verify-subcmd"),
				pluginapi.TaskInfoVerifyOptions(pluginapi.NewVerifyOptions(
					pluginapi.VerifyOptionsApplyFalseArgs("--no-apply"),
					pluginapi.VerifyOptionsApplyTrueArgs("--apply"),
				)),
			},
			godellauncher.GlobalConfig{
				TaskArgs: []string{"verify"},
			},
			"Running echo...\nverify-subcmd --apply\n",
		},
		{
			"verify task with apply=false",
			[]pluginapi.TaskInfoParam{
				pluginapi.TaskInfoCommand("verify-subcmd"),
				pluginapi.TaskInfoVerifyOptions(pluginapi.NewVerifyOptions(
					pluginapi.VerifyOptionsApplyFalseArgs("--no-apply"),
					pluginapi.VerifyOptionsApplyTrueArgs("--apply"),
				)),
			},
			godellauncher.GlobalConfig{
				TaskArgs: []string{"verify", "--apply=false"},
			},
			"Running echo...\nverify-subcmd --no-apply\n",
		},
		{
			"verify task with global flag options",
			[]pluginapi.TaskInfoParam{
				pluginapi.TaskInfoCommand("verify-subcmd"),
				pluginapi.TaskInfoGlobalFlagOptions(pluginapi.NewGlobalFlagOptions(
					pluginapi.GlobalFlagOptionsParamProjectDirFlag("--project-dir"),
				)),
				pluginapi.TaskInfoVerifyOptions(pluginapi.NewVerifyOptions(
					pluginapi.VerifyOptionsApplyFalseArgs("--no-apply"),
					pluginapi.VerifyOptionsApplyTrueArgs("--apply"),
				)),
			},
			godellauncher.GlobalConfig{
				TaskArgs: []string{"verify"},
				Wrapper:  "godelw",
			},
			"Running echo...\n--project-dir . verify-subcmd --apply\n",
		},
	} {
		pluginInfo, err := pluginapi.NewInfo("group", "echo", "1.0.0", "echo.yml",
			pluginapi.MustNewTaskInfo("echo", "echoes the provided input", tc.params...))
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
