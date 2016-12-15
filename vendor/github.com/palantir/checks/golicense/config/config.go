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

	"github.com/palantir/checks/golicense/golicense"
)

type GoLicense struct {
	// Header is the expected license header. All applicable files are expected to start with this header followed
	// by a newline.
	Header string `yaml:"header" json:"header"`

	// CustomHeaders specifies the custom header parameters. Custom header parameters can be used to specify that
	// certain directories or files in the project should use a header that is different from "Header".
	CustomHeaders []License `yaml:"custom-headers" json:"custom-headers"`

	// Exclude matches the files and directories that should be excluded from consideration for verifying or
	// applying licenses.
	Exclude matcher.NamesPathsCfg `yaml:"exclude" json:"exclude"`
}

type License struct {
	// Name is the identifier used to identify this custom license parameter. Must be unique.
	Name string `yaml:"name" json:"name"`

	// Header is the expected license header. All applicable files are expected to start with this header followed
	// by a newline.
	Header string `yaml:"header" json:"header"`

	// Paths specifies the paths for which this custom license is applicable. If multiple custom parameters match a
	// file or directory, the parameter with the longest path match is used. If multiple custom parameters match a
	// file or directory exactly (match length is equal), it is treated as an error.
	Paths []string `yaml:"paths" json:"paths"`
}

func (l *GoLicense) ToParams() (golicense.LicenseParams, error) {
	customHeaders := make([]golicense.CustomLicenseParam, len(l.CustomHeaders))
	for i, v := range l.CustomHeaders {
		customHeaders[i] = v.ToParam()
	}
	customParams, err := golicense.NewCustomLicenseParams(customHeaders)
	if err != nil {
		return golicense.LicenseParams{}, err
	}
	return golicense.LicenseParams{
		Header:        l.Header,
		CustomHeaders: customParams,
		Exclude:       l.Exclude.Matcher(),
	}, nil
}

func (l *License) ToParam() golicense.CustomLicenseParam {
	return golicense.CustomLicenseParam{
		Name:         l.Name,
		Header:       l.Header,
		IncludePaths: l.Paths,
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
	cfg, err := LoadFromStrings(string(yml), jsonContent)
	if err != nil {
		return golicense.LicenseParams{}, err
	}
	return cfg.ToParams()
}

func LoadFromStrings(ymlContent, jsonContent string) (GoLicense, error) {
	cfg := GoLicense{}
	if ymlContent != "" {
		if err := yaml.Unmarshal([]byte(ymlContent), &cfg); err != nil {
			return GoLicense{}, errors.Wrapf(err, "failed to unmarshal YML %s", ymlContent)
		}
	}
	if jsonContent != "" {
		jsonCfg := GoLicense{}
		if err := json.Unmarshal([]byte(jsonContent), &jsonCfg); err != nil {
			return GoLicense{}, err
		}
		cfg.Exclude.Add(jsonCfg.Exclude)
	}
	return cfg, nil
}
