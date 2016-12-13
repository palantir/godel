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
	"path"

	"github.com/palantir/pkg/matcher"
	"github.com/palantir/pkg/specdir"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"

	"github.com/palantir/godel/layout"
)

const (
	ExcludeYML = "exclude.yml"
)

type Exclude struct {
	// Exclude specifies the files and directories that should be excluded from gödel operations. This parameter is
	// also passed to subtasks to augment their notion of included and excluded files and directories.
	Exclude matcher.NamesPathsCfg `json:"exclude"`
}

func LoadFromFile(cfgPath string) (matcher.NamesPathsCfg, error) {
	fileBytes, err := ioutil.ReadFile(cfgPath)
	if err != nil {
		return matcher.NamesPathsCfg{}, errors.Wrapf(err, "Failed to read file %s", cfgPath)
	}
	return LoadFromYML(string(fileBytes))
}

func LoadFromYML(ymlContent string) (matcher.NamesPathsCfg, error) {
	excludeCfg := matcher.NamesPathsCfg{}
	if err := yaml.Unmarshal([]byte(ymlContent), &excludeCfg); err != nil {
		return matcher.NamesPathsCfg{}, errors.Wrapf(err, "Failed to unmarshal YML %s", ymlContent)
	}
	return excludeCfg, nil
}

func ReadExcludeJSONFromYML(cfgPath string) ([]byte, error) {
	excludeCfg, err := LoadFromFile(cfgPath)
	if err != nil {
		return nil, err
	}
	return json.Marshal(Exclude{
		Exclude: excludeCfg,
	})
}

func GetCfgDirPath(cfgDirPath, wrapperScriptPath string) (string, error) {
	if cfgDirPath != "" {
		return cfgDirPath, nil
	}
	if wrapperScriptPath == "" {
		return "", nil
	}
	wrapper, err := specdir.New(path.Dir(wrapperScriptPath), layout.WrapperSpec(), nil, specdir.Validate)
	if err != nil {
		return "", err
	}
	return wrapper.Path(layout.WrapperConfigDir), nil
}
