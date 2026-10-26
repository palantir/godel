// Copyright 2026 Palantir Technologies, Inc.
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

package pluginapitester_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/palantir/godel/v2/framework/godellauncher"
	"github.com/palantir/godel/v2/framework/pluginapi/v2/pluginapi"
	"github.com/palantir/godel/v2/framework/pluginapitester"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRunPluginScopesManagedGoToolchain(t *testing.T) {
	t.Setenv("GOTOOLCHAIN", "local")
	pluginInfo := pluginapi.MustNewPluginInfo(
		"com.palantir",
		"environment-plugin",
		"1.0.0",
		pluginapi.PluginInfoTaskInfo("environment", "Prints GOTOOLCHAIN", pluginapi.TaskInfoCommand("environment")),
	)
	pluginInfoJSON, err := json.Marshal(pluginInfo)
	require.NoError(t, err)

	for _, tc := range []struct {
		name       string
		config     string
		goMod      string
		wantPlugin string
		parentSet  bool
	}{
		{
			name:       "managed",
			config:     "manage-gotoolchain: true\n",
			goMod:      "module example.com/test\ngo 1.27.0\ntoolchain go1.27.1\n",
			wantPlugin: "go1.27.1",
			parentSet:  true,
		},
		{
			name:       "unmanaged after managed",
			config:     "manage-gotoolchain: false\n",
			wantPlugin: "local",
			parentSet:  true,
		},
		{
			name:       "managed restores unset value",
			config:     "manage-gotoolchain: true\n",
			goMod:      "module example.com/test\ngo 1.27.0\ntoolchain go1.27.1\n",
			wantPlugin: "go1.27.1",
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if !tc.parentSet {
				require.NoError(t, os.Unsetenv("GOTOOLCHAIN"))
			}
			projectDir := t.TempDir()
			configDir := filepath.Join(projectDir, "godel", "config")
			require.NoError(t, os.MkdirAll(configDir, 0755))
			require.NoError(t, os.WriteFile(filepath.Join(configDir, godellauncher.GodelConfigYML), []byte(tc.config), 0644))
			if tc.goMod != "" {
				require.NoError(t, os.WriteFile(filepath.Join(projectDir, "go.mod"), []byte(tc.goMod), 0644))
			}

			metadataEnvironment := filepath.Join(projectDir, "metadata-environment")
			pluginPath := filepath.Join(projectDir, "plugin")
			plugin := fmt.Sprintf(`#!/bin/sh
if [ "$1" = %q ]; then
    printf '%%s' "$GOTOOLCHAIN" > %q
    printf '%%s' %q
    exit 0
fi
printf '%%s' "$GOTOOLCHAIN"
`, pluginapi.PluginInfoCommandName, metadataEnvironment, string(pluginInfoJSON))
			require.NoError(t, os.WriteFile(pluginPath, []byte(plugin), 0755))

			var stdout bytes.Buffer
			cleanup, err := pluginapitester.RunPlugin(
				pluginapitester.NewPluginProvider(pluginPath),
				nil,
				"environment",
				nil,
				projectDir,
				false,
				&stdout,
			)
			require.NoError(t, err)
			wrapperPath := filepath.Join(projectDir, "godelw")
			require.FileExists(t, wrapperPath)
			assert.Equal(t, tc.wantPlugin, stdout.String())
			metadataValue, err := os.ReadFile(metadataEnvironment)
			require.NoError(t, err)
			assert.Equal(t, tc.wantPlugin, string(metadataValue))
			parentValue, parentSet := os.LookupEnv("GOTOOLCHAIN")
			assert.Equal(t, tc.parentSet, parentSet)
			if tc.parentSet {
				assert.Equal(t, "local", parentValue)
			}
			cleanup()
			assert.NoFileExists(t, wrapperPath)
		})
	}
}
