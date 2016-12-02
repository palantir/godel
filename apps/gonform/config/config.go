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
)

type rawConfig struct {
	Formatters map[string]FormatterConfig `yaml:"formatters" json:"formatters"` // custom configuration provided to formatters
	Exclude    matcher.NamesPathsCfg      `yaml:"exclude" json:"exclude"`
}

type FormatterConfig struct {
	// Args is a slice of the arguments provided to the formatter. Each element in the slice is provided as a
	// separate argument for the formatter.
	Args []string `yaml:"args" json:"args"`
}

type Config struct {
	Formatters map[string]FormatterConfig
	Exclude    matcher.Matcher
}

func Load(cfgPath, jsonContent string) (Config, error) {
	var ymlContent string
	if cfgPath != "" {
		content, err := ioutil.ReadFile(cfgPath)
		if err != nil {
			return Config{}, errors.Wrapf(err, "failed to read file %s", cfgPath)
		}
		ymlContent = string(content)
	}
	return LoadFromString(ymlContent, jsonContent)
}

func LoadFromString(ymlContent, jsonContent string) (Config, error) {
	cfg := rawConfig{}
	if ymlContent != "" {
		if err := yaml.Unmarshal([]byte(ymlContent), &cfg); err != nil {
			return Config{}, errors.Wrapf(err, "failed to unmarshal YML %s", ymlContent)
		}
	}

	if jsonContent != "" {
		jsonCfg := rawConfig{}
		if err := json.Unmarshal([]byte(jsonContent), &jsonCfg); err != nil {
			return Config{}, err
		}
		cfg.Exclude.Add(jsonCfg.Exclude)
	}

	return Config{
		Formatters: cfg.Formatters,
		Exclude:    cfg.Exclude.Matcher(),
	}, nil
}
