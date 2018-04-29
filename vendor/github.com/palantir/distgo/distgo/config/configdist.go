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
	"path"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"

	"github.com/palantir/distgo/distgo"
	"github.com/palantir/distgo/distgo/config/internal/v0"
)

type DistConfig v0.DistConfig

func ToDistConfig(in *DistConfig) *v0.DistConfig {
	return (*v0.DistConfig)(in)
}

// ToParam returns the DistParam represented by the receiver *DisterConfig and the provided default DisterConfig. If a
// config value is specified (non-nil) in the receiver config, it is used. If a config value is not specified in the
// receiver config but is specified in the default config, the default config value is used. If a value is not specified
// in either configuration, the program-specified default value (if any) is used.
func (cfg *DistConfig) ToParam(scriptIncludes string, defaultCfg DistConfig, disterFactory distgo.DisterFactory) (distgo.DistParam, error) {
	outputDir := getConfigStringValue(cfg.OutputDir, defaultCfg.OutputDir, "out/dist")
	if path.IsAbs(outputDir) {
		return distgo.DistParam{}, errors.Errorf("output-dir cannot be specified as an absolute path")
	}
	disters, err := (*DistersConfig)(cfg.Disters).ToParam((*DistersConfig)(defaultCfg.Disters), scriptIncludes, disterFactory)
	if err != nil {
		return distgo.DistParam{}, err
	}
	return distgo.DistParam{
		OutputDir:  outputDir,
		DistParams: disters,
	}, nil
}

type DisterConfig v0.DisterConfig

func ToDisterConfig(in DisterConfig) v0.DisterConfig {
	return (v0.DisterConfig)(in)
}

func (cfg *DisterConfig) ToParam(defaultCfg DisterConfig, scriptIncludes string, disterFactory distgo.DisterFactory) (distgo.DisterParam, error) {
	disterType := getConfigStringValue(cfg.Type, defaultCfg.Type, "")
	if disterType == "" {
		return distgo.DisterParam{}, errors.Errorf("dister type must be specified for DisterConfig")
	}
	dister, err := newDister(disterType, getConfigValue(cfg.Config, defaultCfg.Config, nil).(yaml.MapSlice), disterFactory)
	if err != nil {
		return distgo.DisterParam{}, err
	}

	return distgo.DisterParam{
		NameTemplate: getConfigStringValue(cfg.NameTemplate, defaultCfg.NameTemplate, "{{Product}}-{{Version}}"),
		Script:       distgo.CreateScriptContent(getConfigStringValue(cfg.Script, defaultCfg.Script, ""), scriptIncludes),
		Dister:       dister,
	}, nil
}

func newDister(disterType string, cfgYML yaml.MapSlice, disterFactory distgo.DisterFactory) (distgo.Dister, error) {
	if disterType == "" {
		return nil, errors.Errorf("dister type must be non-empty")
	}
	if disterFactory == nil {
		return nil, errors.Errorf("disterFactory must be provided")
	}
	cfgYMLBytes, err := yaml.Marshal(cfgYML)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to marshal configuration")
	}
	return disterFactory.NewDister(disterType, cfgYMLBytes)
}

type DistersConfig v0.DistersConfig

func ToDistersConfig(in *DistersConfig) *v0.DistersConfig {
	return (*v0.DistersConfig)(in)
}

func (cfgs *DistersConfig) ToParam(defaultCfg *DistersConfig, scriptIncludes string, disterFactory distgo.DisterFactory) (map[distgo.DistID]distgo.DisterParam, error) {
	// keys that exist either only in cfgs or only in defaultCfg
	distinctCfgs := make(map[distgo.DistID]DisterConfig)
	// keys that appear in both cfgs and defaultCfg
	commonCfgIDs := make(map[distgo.DistID]struct{})

	if cfgs != nil {
		for distID, distCfg := range *cfgs {
			if defaultCfg != nil {
				if _, ok := (*defaultCfg)[distID]; ok {
					commonCfgIDs[distID] = struct{}{}
					continue
				}
			}
			distinctCfgs[distID] = DisterConfig(distCfg)
		}
	}
	if defaultCfg != nil {
		for distID, distCfg := range *defaultCfg {
			if cfgs != nil {
				if _, ok := (*cfgs)[distID]; ok {
					commonCfgIDs[distID] = struct{}{}
					continue
				}
			}
			distinctCfgs[distID] = DisterConfig(distCfg)
		}
	}

	distParamsMap := make(map[distgo.DistID]distgo.DisterParam)
	// generate parameters for all of the distinct elements
	for distID, distCfg := range distinctCfgs {
		currParam, err := distCfg.ToParam(DisterConfig{}, scriptIncludes, disterFactory)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to generate parameter for dist configuration %s", distID)
		}
		distParamsMap[distID] = currParam
	}
	// merge keys that appear in both maps
	for distID := range commonCfgIDs {
		currCfg := (*cfgs)[distID]
		currParam, err := (*DisterConfig)(&currCfg).ToParam(DisterConfig((*defaultCfg)[distID]), scriptIncludes, disterFactory)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to generate parameter for dist configuration %s", distID)
		}
		distParamsMap[distID] = currParam
	}
	return distParamsMap, nil
}
