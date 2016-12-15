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

	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"

	"github.com/palantir/checks/gocd/gocd"
)

type GoCD struct {
	RootDirs []string `yaml:"root-dirs"`
}

func (r *GoCD) ToParams() gocd.Params {
	return gocd.Params{
		RootDirs: r.RootDirs,
	}
}

func Load(configPath, _ string) (gocd.Params, error) {
	var yml []byte
	if configPath != "" {
		var err error
		yml, err = ioutil.ReadFile(configPath)
		if err != nil {
			return gocd.Params{}, errors.Wrapf(err, "failed to read file %s", configPath)
		}
	}
	cfg, err := LoadFromYML(string(yml))
	if err != nil {
		return gocd.Params{}, err
	}
	return cfg.ToParams(), nil
}

func LoadFromYML(yml string) (GoCD, error) {
	cfg := GoCD{}
	if yml != "" {
		if err := yaml.Unmarshal([]byte(yml), &cfg); err != nil {
			return GoCD{}, errors.Wrapf(err, "failed to unmarshal YML %s", yml)
		}
	}
	return cfg, nil
}
