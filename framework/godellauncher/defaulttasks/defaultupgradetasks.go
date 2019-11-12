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

package defaulttasks

import (
	"io"

	"github.com/palantir/godel/v2/framework/godel/config"
	"github.com/palantir/godel/v2/framework/godellauncher"
)

func BuiltinUpgradeConfigTasks() []godellauncher.UpgradeConfigTask {
	return []godellauncher.UpgradeConfigTask{
		upgradeGodelConfigTask(),
	}
}

func upgradeGodelConfigTask() godellauncher.UpgradeConfigTask {
	return godellauncher.UpgradeConfigTask{
		ID:         "com.palantir.godel:godel",
		ConfigFile: "godel.yml",
		// Not setting the "LegacyConfigFile" is intentional: legacy configuration is serialization-compatible so no
		// need to run the upgrader. godel.yml is also special because it is loaded before tasks.
		RunImpl: func(t *godellauncher.UpgradeConfigTask, global godellauncher.GlobalConfig, configBytes []byte, stdout io.Writer) ([]byte, error) {
			return config.UpgradeConfig(configBytes)
		},
	}
}
