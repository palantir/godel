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

package godellauncher

import (
	"io"
)

type UpgradeConfigTask struct {
	// The ID for the task. Should be of the form "groupID:productID" ("com.palantir.godel:godel",
	// "com.palantir.format-plugin:format-plugin", etc.).
	ID string

	// The name of the configuration file for the task ("godel.yml", "check-plugin.yml" etc.).
	ConfigFile string

	// The name of the legacy configuration file for the task. Blank if none exists.
	LegacyConfigFile string

	// Configures the manner in which the global flags are processed.
	GlobalFlagOpts GlobalFlagOptions

	// The runner that is invoked to run this config upgrade task. Takes the provided input config bytes and returns the
	// upgraded configuration bytes. If the provided config bytes are a valid representation of configuration for the
	// most recent version, the returned bytes should be the same as the input. Should be possible to run in-process
	// (that is, this function should not call os.Exit or equivalent).
	RunImpl func(t *UpgradeConfigTask, global GlobalConfig, configBytes []byte, stdout io.Writer) ([]byte, error)
}

func (t *UpgradeConfigTask) Run(configBytes []byte, global GlobalConfig, stdout io.Writer) ([]byte, error) {
	return t.RunImpl(t, global, configBytes, stdout)
}
