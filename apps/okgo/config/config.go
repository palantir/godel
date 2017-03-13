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

type OKGo struct {
	// ReleaseTag specifies the newest Go release build tag supported by the codebase being checked. If this value
	// is not specified, it defaults to the Go release that was used to build the check tool. If the code being
	// checked is known to use a version of Go that is earlier than the version of Go used to build the check tool
	// and the codebase being checked contains build tags for the newer Go version, this value should be explicitly
	// set. For example, if the check tool was compiled using Go 1.8 but the codebase being checked uses Go 1.7 and
	// contains files that use the "// +build go1.8" build tag, then this should be set to "go1.7".
	ReleaseTag string `yaml:"release-tag" json:"release-tag"`

	// Checks specifies the configuration used by the checks. The key is the name of the check and the value is the
	// custom configuration for that check.
	Checks map[string]Checker `yaml:"checks" json:"checks"`

	// Exclude specifies the files that should be excluded from tests.
	Exclude matcher.NamesPathsCfg `yaml:"exclude" json:"exclude"`
}

type Checker struct {
	// Skip specifies whether or not the check should be skipped entirely.
	Skip bool `yaml:"skip" json:"skip"`

	// Args specifies the command-line arguments provided to the check.
	Args []string `yaml:"args" json:"args"`

	// Filters specifies the filter definitions. Raw output lines that match the filter are excluded from
	// processing.
	Filters []Filter `yaml:"filters" json:"filters"`
}

type Filter struct {
	// Type specifies the type of the filter: "message", "name" or "path". If blank, defaults to "message".
	Type string `yaml:"type" json:"type"`

	// The value of the filter.
	Value string `yaml:"value" json:"value"`
}

func (r *OKGo) ToParams() (params.OKGo, error) {
	checks := make(map[amalgomated.Cmd]params.Checker)
	for key, value := range r.Checks {
		singleParam, err := value.ToParam()
		if err != nil {
			return params.OKGo{}, err
		}
		cmd, err := cmdlib.Instance().NewCmd(key)
		if err != nil {
			return params.OKGo{}, errors.Wrapf(err, "unable to convert %s into a command", key)
		}
		checks[cmd] = singleParam
	}
	return params.OKGo{
		ReleaseTag: r.ReleaseTag,
		Checks:     checks,
		Exclude:    r.Exclude.Matcher(),
	}, nil
}

func (r *Checker) ToParam() (params.Checker, error) {
	var lineFilters []checkoutput.Filterer
	for _, cfg := range r.Filters {
		checkFilter, err := cfg.toFilter(checkoutput.MessageRegexpFilter)
		if err != nil {
			return params.Checker{}, errors.Wrapf(err, "failed to parse filter: %v", cfg)
		}
		lineFilters = append(lineFilters, checkFilter)
	}
	return params.Checker{
		Skip:        r.Skip,
		Args:        r.Args,
		LineFilters: lineFilters,
	}, nil
}

func (f *Filter) toFilter(filterForBlankType func(name string) checkoutput.Filterer) (checkoutput.Filterer, error) {
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

func Load(configPath, jsonContent string) (params.OKGo, error) {
	var yml []byte
	if configPath != "" {
		var err error
		yml, err = ioutil.ReadFile(configPath)
		if err != nil {
			return params.OKGo{}, errors.Wrapf(err, "failed to read file %s", configPath)
		}
	}
	cfg, err := LoadRawConfig(string(yml), jsonContent)
	if err != nil {
		return params.OKGo{}, err
	}
	return cfg.ToParams()
}

func LoadRawConfig(ymlContent, jsonContent string) (OKGo, error) {
	rawCfg := OKGo{}
	if ymlContent != "" {
		if err := yaml.Unmarshal([]byte(ymlContent), &rawCfg); err != nil {
			return OKGo{}, errors.Wrapf(err, "failed to unmarshal YML %s", ymlContent)
		}
	}
	if jsonContent != "" {
		jsonCfg := OKGo{}
		if err := json.Unmarshal([]byte(jsonContent), &jsonCfg); err != nil {
			return OKGo{}, errors.Wrapf(err, "failed to parse JSON %s", jsonContent)
		}
		rawCfg.Exclude.Add(jsonCfg.Exclude)
	}
	return rawCfg, nil
}
