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

	"github.com/palantir/pkg/matcher"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"

	"github.com/palantir/godel/v2/framework/godellauncher"
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
	if _, err := os.Stat(cfgFile); os.IsNotExist(err) {
		return GodelConfig{}, nil
	}

	var godelCfg GodelConfig
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
	return godelCfg, nil
}

// ReadGodelConfigExcludesFromFile reads the excludes specified in the gödel godel configuration from the provided file
// and returns the loaded configuration. Returns an empty configuration if the file does not exist. Callers that only
// require the exclude configuration should prefer this function to using the ReadGodelConfigFrom* functions and
// accessing the exclude configuration there, as this function is more robust to configuration changes (for example, the
// ReadGodelConfigFrom* functions will return an error if the configuration has an unrecognized key, while this function
// only considers the "exclude" portion of configuration).
func ReadGodelConfigExcludesFromFile(cfgFilePath string) (matcher.NamesPathsCfg, error) {
	if _, err := os.Stat(cfgFilePath); os.IsNotExist(err) {
		return matcher.NamesPathsCfg{}, nil
	}

	type excludeConfig struct {
		Exclude matcher.NamesPathsCfg `yaml:"exclude,omitempty"`
	}
	cfgBytes, err := ioutil.ReadFile(cfgFilePath)
	if err != nil {
		return matcher.NamesPathsCfg{}, errors.WithStack(err)
	}
	var exclude excludeConfig
	if err := yaml.Unmarshal(cfgBytes, &exclude); err != nil {
		return matcher.NamesPathsCfg{}, errors.WithStack(err)
	}
	return exclude.Exclude, nil
}
