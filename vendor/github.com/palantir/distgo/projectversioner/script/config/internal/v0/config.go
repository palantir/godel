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

package v0

import (
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

type Config struct {
	// Script specifies the content of the script that is run to determine the version for the project. This script is
	// written to a temporary directory with executable permissions (0755) and run. The absolute path to the project
	// directory is set as the environment variable "PROJECT_DIR". The script is written exactly as provided, so any
	// necessary headers (#! etc.) should be included. If the script exits with an exit code of 0, the result of calling
	// strings.TrimSpace on the output produced by the script (STDOUT and STDERR) is returned as the version. If the
	// script exist with a non-0 exit code, that is treated as an error.
	Script string `yaml:"script,omitempty"`
}

func UpgradeConfig(cfgBytes []byte) ([]byte, error) {
	var cfg Config
	if err := yaml.UnmarshalStrict(cfgBytes, &cfg); err != nil {
		return nil, errors.Wrapf(err, "failed to unmarshal script project versioner v0 configuration")
	}
	return cfgBytes, nil
}
