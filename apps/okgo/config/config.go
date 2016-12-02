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
)

type Config struct {
	Checks  map[amalgomated.Cmd]SingleCheckerConfig
	Exclude matcher.Matcher
}

type SingleCheckerConfig struct {
	Skip        bool
	Args        []string
	LineFilters []checkoutput.Filterer
}

type rawConfig struct {
	Checks  map[string]rawCheckConfig `yaml:"checks" json:"checks"` // configuration for checkers
	Exclude matcher.NamesPathsCfg     `yaml:"exclude" json:"exclude"`
}

type rawCheckConfig struct {
	Skip    bool        `yaml:"skip" json:"skip"`       // skip this check if true
	Args    []string    `yaml:"args" json:"args"`       // arguments provided to the check
	Filters []rawFilter `yaml:"filters" json:"filters"` // defines filters that filters out raw output lines that match the filters from consideration
}

type rawFilter struct {
	Type  string `yaml:"type" json:"type"`   // type of filter: "message", "name" or "path"
	Value string `yaml:"value" json:"value"` // value of the filter
}

func (f *rawFilter) toFilter(filterForBlankType func(name string) checkoutput.Filterer) (checkoutput.Filterer, error) {
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

// ArgsForCheck returns the arguments for the requested check stored in the Config, or nil if no configuration for the
// specified check was present in the configuration. The second return value indicates whether or not configuration for
// the requested check was present.
func (c *Config) ArgsForCheck(check amalgomated.Cmd) ([]string, bool) {
	checkConfig, ok := c.Checks[check]
	if !ok {
		return nil, false
	}
	return checkConfig.Args, true
}

// FiltersForCheck returns the filters that should be used for the requested check. The returned slice is a
// concatenation of the global filters derived from the package excludes specified in the configuration followed by the
// filters specified for the provided check in the configuration. Returns an empty slice if no filters are present
// globally or for the specified check.The derivation from the global filters is done in case the packages can't be
// excluded before the check is run (can happen if the check only supports the "all" mode).
func (c *Config) FiltersForCheck(check amalgomated.Cmd) []checkoutput.Filterer {
	filters := append([]checkoutput.Filterer{}, checkoutput.MatcherFilter(c.Exclude))
	checkConfig, ok := c.Checks[check]
	if ok {
		filters = append(filters, checkConfig.LineFilters...)
	}
	return filters
}

func (c *Config) checkCommands() []amalgomated.Cmd {
	var cmds []amalgomated.Cmd
	for _, currCmd := range cmdlib.Instance().Cmds() {
		if _, ok := c.Checks[currCmd]; ok {
			cmds = append(cmds, currCmd)
		}
	}
	return cmds
}

func Load(configPath, jsonContent string) (Config, error) {
	var ymlContent string
	if configPath != "" {
		file, err := ioutil.ReadFile(configPath)
		if err != nil {
			return Config{}, errors.Wrapf(err, "failed to read file %s", configPath)
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
			return Config{}, errors.Wrapf(err, "failed to parse JSON %s", jsonContent)
		}
		rawCfg.Exclude.Add(jsonCfg.Exclude)
	}

	checks := make(map[amalgomated.Cmd]SingleCheckerConfig)
	for key, value := range rawCfg.Checks {
		checkFilters := make([]checkoutput.Filterer, len(value.Filters))
		for i, rawCheckFilter := range value.Filters {
			checkFilter, err := rawCheckFilter.toFilter(checkoutput.MessageRegexpFilter)
			if err != nil {
				return Config{}, errors.Wrapf(err, "failed to parse filter: %v", rawCheckFilter)
			}
			checkFilters[i] = checkFilter
		}

		cmd, err := cmdlib.Instance().NewCmd(key)
		if err != nil {
			return Config{}, errors.Wrapf(err, "unable to convert %s into a command", key)
		}

		checks[cmd] = SingleCheckerConfig{
			Skip:        value.Skip,
			Args:        value.Args,
			LineFilters: checkFilters,
		}
	}

	return Config{
		Checks:  checks,
		Exclude: rawCfg.Exclude.Matcher(),
	}, nil
}
