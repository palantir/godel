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

type Gonform struct {
	// Formatters specifies the configuration used by the formatters. The key is the name of the formatter and the
	// value is the custom configuration for that formatter.
	Formatters map[string]Formatter `yaml:"formatters" json:"formatters"`

	// Exclude specifies the files that should be excluded from formatting.
	Exclude matcher.NamesPathsCfg `yaml:"exclude" json:"exclude"`
}

type Formatter struct {
	// Args specifies the command-line arguments that are provided to the formatter.
	Args []string `yaml:"args" json:"args"`
}

func (r *Gonform) ToParams() params.Formatters {
	m := make(map[string]params.Formatter, len(r.Formatters))
	for k, v := range r.Formatters {
		m[k] = v.ToParams()
	}
	return params.Formatters{
		Formatters: m,
		Exclude:    r.Exclude.Matcher(),
	}
}

func (r *Formatter) ToParams() params.Formatter {
	return params.Formatter{
		Args: r.Args,
	}
}

func Load(cfgPath, jsonContent string) (params.Formatters, error) {
	var yml []byte
	if cfgPath != "" {
		var err error
		yml, err = ioutil.ReadFile(cfgPath)
		if err != nil {
			return params.Formatters{}, errors.Wrapf(err, "failed to read file %s", cfgPath)
		}
	}
	cfg, err := LoadRawConfig(string(yml), jsonContent)
	if err != nil {
		return params.Formatters{}, err
	}
	return cfg.ToParams(), nil
}

func LoadRawConfig(ymlContent, jsonContent string) (Gonform, error) {
	cfg := Gonform{}
	if ymlContent != "" {
		if err := yaml.Unmarshal([]byte(ymlContent), &cfg); err != nil {
			return Gonform{}, errors.Wrapf(err, "failed to unmarshal YML %s", ymlContent)
		}
	}
	if jsonContent != "" {
		jsonCfg := Gonform{}
		if err := json.Unmarshal([]byte(jsonContent), &jsonCfg); err != nil {
			return Gonform{}, err
		}
		cfg.Exclude.Add(jsonCfg.Exclude)
	}
	return cfg, nil
}
