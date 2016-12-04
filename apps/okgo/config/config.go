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

	"github.com/palantir/amalgomate/amalgomated"
	"github.com/palantir/pkg/matcher"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"

	"github.com/palantir/godel/apps/okgo/checkoutput"
	"github.com/palantir/godel/apps/okgo/cmd/cmdlib"
	"github.com/palantir/godel/apps/okgo/params"
)

type RawConfig struct {
	Checks  map[string]RawSingleCheckerConfig `yaml:"checks" json:"checks"` // configuration for checkers
	Exclude matcher.NamesPathsCfg             `yaml:"exclude" json:"exclude"`
}

func (r *RawConfig) ToParams() (params.Params, error) {
	checks := make(map[amalgomated.Cmd]params.SingleCheckerParam)
	for key, value := range r.Checks {
		singleParam, err := value.ToParam()
		if err != nil {
			return params.Params{}, err
		}
		cmd, err := cmdlib.Instance().NewCmd(key)
		if err != nil {
			return params.Params{}, errors.Wrapf(err, "unable to convert %s into a command", key)
		}
		checks[cmd] = singleParam
	}
	return params.Params{
		Checks:  checks,
		Exclude: r.Exclude.Matcher(),
	}, nil
}

type RawSingleCheckerConfig struct {
	Skip    bool              `yaml:"skip" json:"skip"`       // skip this check if true
	Args    []string          `yaml:"args" json:"args"`       // arguments provided to the check
	Filters []RawFilterConfig `yaml:"filters" json:"filters"` // defines filters that filters out raw output lines that match the filters from consideration
}

func (r *RawSingleCheckerConfig) ToParam() (params.SingleCheckerParam, error) {
	var lineFilters []checkoutput.Filterer
	for _, cfg := range r.Filters {
		checkFilter, err := cfg.toFilter(checkoutput.MessageRegexpFilter)
		if err != nil {
			return params.SingleCheckerParam{}, errors.Wrapf(err, "failed to parse filter: %v", cfg)
		}
		lineFilters = append(lineFilters, checkFilter)
	}
	return params.SingleCheckerParam{
		Skip:        r.Skip,
		Args:        r.Args,
		LineFilters: lineFilters,
	}, nil
}

type RawFilterConfig struct {
	Type  string `yaml:"type" json:"type"`   // type of filter: "message", "name" or "path"
	Value string `yaml:"value" json:"value"` // value of the filter
}

func (f *RawFilterConfig) toFilter(filterForBlankType func(name string) checkoutput.Filterer) (checkoutput.Filterer, error) {
	switch f.Type {
	case "message":
		return checkoutput.MessageRegexpFilter(f.Value), nil
	case "name":
		return checkoutput.NamePathFilter(f.Value), nil
	case "path":
		return checkoutput.RelativePathFilter(f.Value), nil
	case "":
		if filterForBlankType != nil {
			return filterForBlankType(f.Value), nil
		}
		fallthrough
	default:
		return nil, errors.Errorf("unknown filter type: %v", f.Type)
	}
}

func Load(configPath, jsonContent string) (params.Params, error) {
	var yml []byte
	if configPath != "" {
		var err error
		yml, err = ioutil.ReadFile(configPath)
		if err != nil {
			return params.Params{}, errors.Wrapf(err, "failed to read file %s", configPath)
		}
	}
	cfg, err := LoadRawConfig(string(yml), jsonContent)
	if err != nil {
		return params.Params{}, err
	}
	return cfg.ToParams()
}

func LoadRawConfig(ymlContent, jsonContent string) (RawConfig, error) {
	rawCfg := RawConfig{}
	if ymlContent != "" {
		if err := yaml.Unmarshal([]byte(ymlContent), &rawCfg); err != nil {
			return RawConfig{}, errors.Wrapf(err, "failed to unmarshal YML %s", ymlContent)
		}
	}
	if jsonContent != "" {
		jsonCfg := RawConfig{}
		if err := json.Unmarshal([]byte(jsonContent), &jsonCfg); err != nil {
			return RawConfig{}, errors.Wrapf(err, "failed to parse JSON %s", jsonContent)
		}
		rawCfg.Exclude.Add(jsonCfg.Exclude)
	}
	return rawCfg, nil
}
