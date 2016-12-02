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

package golicense

import (
	"encoding/json"
	"io/ioutil"

	"github.com/palantir/pkg/matcher"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

type rawConfig struct {
	Header        string                `yaml:"header" json:"header"`
	CustomHeaders []rawLicenseConfig    `yaml:"custom-headers" json:"custom-headers"`
	Exclude       matcher.NamesPathsCfg `yaml:"exclude" json:"exclude"`
}

type rawLicenseConfig struct {
	Name    string                `yaml:"name" json:"name"`
	Header  string                `yaml:"header" json:"header"`
	Include matcher.NamesPathsCfg `yaml:"matchers" json:"matchers"`
}

func Load(configPath, jsonContent string) (LicenseParams, error) {
	rawCfg := rawConfig{}
	if configPath != "" {
		var err error
		rawCfg, err = loadRawConfig(configPath)
		if err != nil {
			return LicenseParams{}, errors.Wrapf(err, "failed to read YML configuration")
		}
	}

	if jsonContent != "" {
		jsonCfg := rawConfig{}
		if err := json.Unmarshal([]byte(jsonContent), &jsonCfg); err != nil {
			return LicenseParams{}, err
		}
		rawCfg.Exclude.Add(jsonCfg.Exclude)
	}

	var emptyNameParams []CustomLicenseParam
	nameToParams := make(map[string][]CustomLicenseParam)

	customHeaders := make([]CustomLicenseParam, len(rawCfg.CustomHeaders))
	for i, v := range rawCfg.CustomHeaders {
		p := CustomLicenseParam{
			Name:    v.Name,
			Header:  v.Header,
			Include: v.Include.Matcher(),
		}
		customHeaders[i] = p

		if p.Name == "" {
			emptyNameParams = append(emptyNameParams, p)
		}
		nameToParams[p.Name] = append(nameToParams[p.Name], p)
	}

	p := LicenseParams{
		Header:        rawCfg.Header,
		CustomHeaders: customHeaders,
		Exclude:       rawCfg.Exclude.Matcher(),
	}

	if err := p.validate(); err != nil {
		return LicenseParams{}, err
	}
	return p, nil
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
