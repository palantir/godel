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
	"encoding/json"
	"io/ioutil"

	"github.com/palantir/pkg/matcher"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"

	"github.com/palantir/godel/apps/gunit/params"
)

type GUnit struct {
	// Tags group tests into different sets. The key is the name of the tag and the value is a
	// matcher.NamesPathsWithExcludeCfg that specifies the rules for matching the tests that are part of the tag.
	// Any test that matches the provided matcher is considered part of the tag.
	Tags map[string]matcher.NamesPathsWithExcludeCfg `yaml:"tags" json:"tags"`

	// Exclude specifies the files that should be excluded from tests.
	Exclude matcher.NamesPathsCfg `yaml:"exclude" json:"exclude"`
}

func (r *GUnit) ToParams() params.GUnit {
	m := make(map[string]matcher.Matcher, len(r.Tags))
	for k, v := range r.Tags {
		m[k] = v.Matcher()
	}
	return params.GUnit{
		Tags:    m,
		Exclude: r.Exclude.Matcher(),
	}
}

func Load(cfgPath, jsonContent string) (params.GUnit, error) {
	var yml []byte
	if cfgPath != "" {
		var err error
		yml, err = ioutil.ReadFile(cfgPath)
		if err != nil {
			return params.GUnit{}, errors.Wrapf(err, "failed to read file %s", cfgPath)
		}
	}
	cfg, err := LoadRawConfig(string(yml), jsonContent)
	if err != nil {
		return params.GUnit{}, err
	}
	return cfg.ToParams(), nil
}

func LoadRawConfig(ymlContent, jsonContent string) (GUnit, error) {
	cfg := GUnit{}
	if ymlContent != "" {
		if err := yaml.Unmarshal([]byte(ymlContent), &cfg); err != nil {
			return GUnit{}, errors.Wrapf(err, "failed to unmarshal YML %s", ymlContent)
		}
	}
	if jsonContent != "" {
		jsonCfg := GUnit{}
		if err := json.Unmarshal([]byte(jsonContent), &jsonCfg); err != nil {
			return GUnit{}, err
		}
		cfg.Exclude.Add(jsonCfg.Exclude)
	}
	return cfg, nil
}
