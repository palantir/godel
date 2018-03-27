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

package config

import (
	"io/ioutil"
	"os"
	"path"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"

	"github.com/palantir/godel/framework/godellauncher"
)

// ReadGodelConfigFromProjectDir reads the gödel configuration from the "godel.yml" file in the configuration directory
// for the gödel project with the specified project directory and returns it. Returns an empty configuration if the
// configuration file does not exist.
func ReadGodelConfigFromProjectDir(projectDir string) (GodelConfig, error) {
	cfgDir, err := godellauncher.ConfigDirPath(projectDir)
	if err != nil {
		return GodelConfig{}, err
	}
	return ReadGodelConfigFromFile(path.Join(cfgDir, godellauncher.GodelConfigYML))
}

// ReadGodelConfigFromFile reads the gödel configuration from the provided file and returns the loaded configuration.
// Returns an empty configuration if the file does not exist.
func ReadGodelConfigFromFile(cfgFile string) (GodelConfig, error) {
	var godelCfg GodelConfig
	if _, err := os.Stat(cfgFile); err == nil {
		bytes, err := ioutil.ReadFile(cfgFile)
		if err != nil {
			return GodelConfig{}, errors.Wrapf(err, "failed to read file %s", cfgFile)
		}
		upgradedBytes, err := UpgradeConfig(bytes)
		if err != nil {
			return GodelConfig{}, errors.Wrapf(err, "failed to upgrade configuration")
		}
		if err := yaml.Unmarshal(upgradedBytes, &godelCfg); err != nil {
			return GodelConfig{}, errors.Wrapf(err, "failed to unmarshal gödel config YAML")
		}
	}
	return godelCfg, nil
}
