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

package gocd

import (
	"io/ioutil"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

type Config struct {
	RootDirs []string
}

type rawConfig struct {
	RootDirs []string `yaml:"root-dirs"`
}

func Load(configPath, _ string) (Config, error) {
	rawCfg := rawConfig{}
	if configPath != "" {
		var err error
		rawCfg, err = loadRawConfig(configPath)
		if err != nil {
			return Config{}, errors.Wrapf(err, "failed to read YML configuration")
		}
	}

	return Config{
		RootDirs: rawCfg.RootDirs,
	}, nil
}

func loadRawConfig(configPath string) (rawConfig, error) {
	file, err := ioutil.ReadFile(configPath)
	if err != nil {
		return rawConfig{}, errors.Wrapf(err, "failed to read file %s", configPath)
	}
	cfg := rawConfig{}
	if err := yaml.Unmarshal(file, &cfg); err != nil {
		return rawConfig{}, errors.Wrapf(err, "failed to unmarshal file %s", configPath)
	}
	return cfg, nil
}
