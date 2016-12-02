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

	"github.com/palantir/checks/golicense"
)

type RawConfig struct {
	Header        string                `yaml:"header" json:"header"`
	CustomHeaders []RawLicenseConfig    `yaml:"custom-headers" json:"custom-headers"`
	Exclude       matcher.NamesPathsCfg `yaml:"exclude" json:"exclude"`
}

func (r *RawConfig) ToParams() (golicense.LicenseParams, error) {
	customHeaders := make([]golicense.CustomLicenseParam, len(r.CustomHeaders))
	for i, v := range r.CustomHeaders {
		customHeaders[i] = v.ToParam()
	}
	customParams, err := golicense.NewCustomLicenseParams(customHeaders)
	if err != nil {
		return golicense.LicenseParams{}, err
	}
	return golicense.LicenseParams{
		Header:        r.Header,
		CustomHeaders: customParams,
		Exclude:       r.Exclude.Matcher(),
	}, nil
}

type RawLicenseConfig struct {
	Name   string   `yaml:"name" json:"name"`
	Header string   `yaml:"header" json:"header"`
	Paths  []string `yaml:"paths" json:"paths"`
}

func (r *RawLicenseConfig) ToParam() golicense.CustomLicenseParam {
	return golicense.CustomLicenseParam{
		Name:         r.Name,
		Header:       r.Header,
		IncludePaths: r.Paths,
	}
}

func Load(configPath, jsonContent string) (golicense.LicenseParams, error) {
	var yml []byte
	if configPath != "" {
		var err error
		yml, err = ioutil.ReadFile(configPath)
		if err != nil {
			return golicense.LicenseParams{}, errors.Wrapf(err, "failed to read file %s", configPath)
		}
	}
	cfg, err := LoadRawConfig(string(yml), jsonContent)
	if err != nil {
		return golicense.LicenseParams{}, err
	}
	return cfg.ToParams()
}

func LoadRawConfig(ymlContent string, jsonContent string) (RawConfig, error) {
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
