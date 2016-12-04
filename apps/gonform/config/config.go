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

	"github.com/palantir/godel/apps/gonform/params"
)

type RawConfig struct {
	Formatters map[string]RawFormatterConfig `yaml:"formatters" json:"formatters"` // custom configuration provided to formatters
	Exclude    matcher.NamesPathsCfg         `yaml:"exclude" json:"exclude"`
}

func (r *RawConfig) ToParams() params.Params {
	m := make(map[string]params.FormatterParams, len(r.Formatters))
	for k, v := range r.Formatters {
		m[k] = v.ToParams()
	}
	return params.Params{
		Formatters: m,
		Exclude:    r.Exclude.Matcher(),
	}
}

type RawFormatterConfig struct {
	Args []string `yaml:"args" json:"args"`
}

func (r *RawFormatterConfig) ToParams() params.FormatterParams {
	return params.FormatterParams{
		Args: r.Args,
	}
}

func Load(cfgPath, jsonContent string) (params.Params, error) {
	var yml []byte
	if cfgPath != "" {
		var err error
		yml, err = ioutil.ReadFile(cfgPath)
		if err != nil {
			return params.Params{}, errors.Wrapf(err, "failed to read file %s", cfgPath)
		}
	}
	cfg, err := LoadRawConfig(string(yml), jsonContent)
	if err != nil {
		return params.Params{}, err
	}
	return cfg.ToParams(), nil
}

func LoadRawConfig(ymlContent, jsonContent string) (RawConfig, error) {
	cfg := RawConfig{}
	if ymlContent != "" {
		if err := yaml.Unmarshal([]byte(ymlContent), &cfg); err != nil {
			return RawConfig{}, errors.Wrapf(err, "failed to unmarshal YML %s", ymlContent)
		}
	}
	if jsonContent != "" {
		jsonCfg := RawConfig{}
		if err := json.Unmarshal([]byte(jsonContent), &jsonCfg); err != nil {
			return RawConfig{}, err
		}
		cfg.Exclude.Add(jsonCfg.Exclude)
	}
	return cfg, nil
}
