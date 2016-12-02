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
	"fmt"
	"io/ioutil"
	"regexp"
	"sort"

	"github.com/palantir/pkg/matcher"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

type Config struct {
	Tags    map[string]matcher.Matcher
	Exclude matcher.Matcher
}

type rawConfig struct {
	// Tags group tests into different sets. The key is the name of the tag and the value is a matcher.NamesPathsCfg
	// that specifies the rules for matching the tests that are part of the tag. Any test that matches the provided
	// matcher is considered part of the tag.
	Tags    map[string]matcher.NamesPathsCfg `yaml:"tags" json:"tags"`
	Exclude matcher.NamesPathsCfg            `yaml:"exclude" json:"exclude"`
}

func Load(cfgPath, jsonContent string) (Config, error) {
	var ymlContent string
	if cfgPath != "" {
		file, err := ioutil.ReadFile(cfgPath)
		if err != nil {
			return Config{}, errors.Wrapf(err, "failed to read file %s", cfgPath)
		}
		ymlContent = string(file)
	}
	return LoadFromString(ymlContent, jsonContent)
}

func LoadFromString(ymlContent, jsonContent string) (Config, error) {
	rawCfg := rawConfig{}
	if ymlContent != "" {
		if err := yaml.Unmarshal([]byte(ymlContent), &rawCfg); err != nil {
			return Config{}, errors.Wrapf(err, "failed to unmarshal YML %s", ymlContent)
		}
	}

	if jsonContent != "" {
		jsonCfg := rawConfig{}
		if err := json.Unmarshal([]byte(jsonContent), &jsonCfg); err != nil {
			return Config{}, err
		}
		rawCfg.Exclude.Add(jsonCfg.Exclude)
	}

	cfg := Config{
		Exclude: rawCfg.Exclude.Matcher(),
	}

	if len(rawCfg.Tags) > 0 {
		var invalidTagNames []string

		cfg.Tags = make(map[string]matcher.Matcher, len(rawCfg.Tags))
		for k, v := range rawCfg.Tags {
			if !validTagName(k) {
				invalidTagNames = append(invalidTagNames, k)
			}
			cfg.Tags[k] = v.Matcher()
		}

		if len(invalidTagNames) > 0 {
			sort.Strings(invalidTagNames)
			return Config{}, fmt.Errorf("invalid tag names: %v", invalidTagNames)
		}
	}
	return cfg, nil
}

var tagRegExp = regexp.MustCompile(`[A-Za-z0-9_-]+`)

func validTagName(tag string) bool {
	return len(tagRegExp.ReplaceAllString(tag, "")) == 0
}
