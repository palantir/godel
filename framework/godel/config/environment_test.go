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

package config_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/palantir/godel/v2/framework/godel/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v2"
)

func TestManageGoToolchainConfiguration(t *testing.T) {
	for _, tc := range []struct {
		name      string
		input     string
		specified bool
		enabled   bool
	}{
		{name: "omitted"},
		{name: "enabled", input: "manage-gotoolchain: true\n", specified: true, enabled: true},
		{name: "disabled", input: "manage-gotoolchain: false\n", specified: true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			path := filepath.Join(t.TempDir(), "godel.yml")
			require.NoError(t, os.WriteFile(path, []byte(tc.input), 0644))
			cfg, err := config.ReadGodelConfigFromFile(path)
			require.NoError(t, err)
			if tc.specified {
				require.NotNil(t, cfg.ManageGoToolchain)
				assert.Equal(t, tc.enabled, *cfg.ManageGoToolchain)
			} else {
				assert.Nil(t, cfg.ManageGoToolchain)
			}
			data, err := yaml.Marshal(cfg)
			require.NoError(t, err)
			var roundTrip config.GodelConfig
			require.NoError(t, yaml.Unmarshal(data, &roundTrip))
			assert.Equal(t, cfg.ManageGoToolchain, roundTrip.ManageGoToolchain)
			if tc.specified {
				assert.Contains(t, string(data), tc.input)
			}
		})
	}
}

func TestApplyEnvironment(t *testing.T) {
	for _, tc := range []struct {
		name             string
		setting          string
		extraEnvironment string
		goMod            string
		want             string
		wantErr          string
	}{
		{name: "omitted without module", want: "local"},
		{name: "disabled without module", setting: "manage-gotoolchain: false", want: "local"},
		{name: "disabled with invalid module", setting: "manage-gotoolchain: false", goMod: "invalid", want: "local"},
		{name: "unmanaged configured value", extraEnvironment: "  GOTOOLCHAIN: go1.26.5\n", want: "go1.26.5"},
		{name: "disabled configured value", setting: "manage-gotoolchain: false", extraEnvironment: "  GOTOOLCHAIN: go1.26.5\n", want: "go1.26.5"},
		{name: "override inherited value", setting: "manage-gotoolchain: true", goMod: "module example.com/test\ngo 1.27.0\ntoolchain go1.27.3\n", want: "go1.27.3"},
		{name: "fallback", setting: "manage-gotoolchain: true", goMod: "module example.com/test\ngo 1.27\n", want: "go1.27.0"},
		{name: "configured conflict", setting: "manage-gotoolchain: true", extraEnvironment: "  GOTOOLCHAIN: local\n", wantErr: "environment must not specify GOTOOLCHAIN"},
		{name: "empty configured conflict", setting: "manage-gotoolchain: true", extraEnvironment: "  GOTOOLCHAIN: \"\"\n", wantErr: "environment must not specify GOTOOLCHAIN"},
		{name: "missing module", setting: "manage-gotoolchain: true", wantErr: "read go.mod"},
		{name: "invalid module", setting: "manage-gotoolchain: true", goMod: "invalid", wantErr: "parse go.mod"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Setenv("GOTOOLCHAIN", "local")
			t.Setenv("GODEL_TEST_TOOLCHAIN_ENV", "original")
			projectDir := t.TempDir()
			if tc.goMod != "" {
				require.NoError(t, os.WriteFile(filepath.Join(projectDir, "go.mod"), []byte(tc.goMod), 0644))
			}
			cfgPath := filepath.Join(projectDir, "godel.yml")
			cfgYAML := tc.setting + "\nenvironment:\n  GODEL_TEST_TOOLCHAIN_ENV: configured\n" + tc.extraEnvironment
			require.NoError(t, os.WriteFile(cfgPath, []byte(cfgYAML), 0644))
			cfg, err := config.ReadGodelConfigFromFile(cfgPath)
			require.NoError(t, err)

			err = cfg.ApplyEnvironment(projectDir)
			if tc.wantErr != "" {
				require.ErrorContains(t, err, tc.wantErr)
				assert.Equal(t, "local", os.Getenv("GOTOOLCHAIN"))
				assert.Equal(t, "original", os.Getenv("GODEL_TEST_TOOLCHAIN_ENV"))
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tc.want, os.Getenv("GOTOOLCHAIN"))
			assert.Equal(t, "configured", os.Getenv("GODEL_TEST_TOOLCHAIN_ENV"))
		})
	}
}
