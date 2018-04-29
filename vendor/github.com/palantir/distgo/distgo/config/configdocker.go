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
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"

	"github.com/palantir/distgo/distgo"
	"github.com/palantir/distgo/distgo/config/internal/v0"
)

type DockerConfig v0.DockerConfig

func ToDockerConfig(in *DockerConfig) *v0.DockerConfig {
	return (*v0.DockerConfig)(in)
}

func (cfg *DockerConfig) ToParam(defaultCfg DockerConfig, dockerBuilderFactory distgo.DockerBuilderFactory) (distgo.DockerParam, error) {
	dockerBuilderParams, err := (*DockerBuildersConfig)(cfg.DockerBuildersConfig).ToParam((*DockerBuildersConfig)(cfg.DockerBuildersConfig), dockerBuilderFactory)
	if err != nil {
		return distgo.DockerParam{}, err
	}
	return distgo.DockerParam{
		Repository:          getConfigStringValue(cfg.Repository, defaultCfg.Repository, ""),
		DockerBuilderParams: dockerBuilderParams,
	}, nil
}

type DockerBuildersConfig v0.DockerBuildersConfig

func ToDockerBuildersConfig(in *DockerBuildersConfig) *v0.DockerBuildersConfig {
	return (*v0.DockerBuildersConfig)(in)
}

func (cfgs *DockerBuildersConfig) ToParam(defaultCfg *DockerBuildersConfig, dockerBuilderFactory distgo.DockerBuilderFactory) (map[distgo.DockerID]distgo.DockerBuilderParam, error) {
	// keys that exist either only in cfgs or only in defaultCfg
	distinctCfgs := make(map[distgo.DockerID]DockerBuilderConfig)
	// keys that appear in both cfgs and defaultCfg
	commonCfgIDs := make(map[distgo.DockerID]struct{})

	if cfgs != nil {
		for dockerID, dockerCfg := range *cfgs {
			if defaultCfg != nil {
				if _, ok := (*defaultCfg)[dockerID]; ok {
					commonCfgIDs[dockerID] = struct{}{}
					continue
				}
			}
			distinctCfgs[dockerID] = DockerBuilderConfig(dockerCfg)
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
			distinctCfgs[distID] = DockerBuilderConfig(distCfg)
		}
	}

	dockerBuilderParamsMap := make(map[distgo.DockerID]distgo.DockerBuilderParam)
	// generate parameters for all of the distinct elements
	for dockerID, dockerBuilderCfg := range distinctCfgs {
		currParam, err := dockerBuilderCfg.ToParam(DockerBuilderConfig{}, dockerBuilderFactory)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to generate parameter for Docker configuration %s", dockerID)
		}
		dockerBuilderParamsMap[dockerID] = currParam
	}
	// merge keys that appear in both maps
	for dockerID := range commonCfgIDs {
		currCfg := (*cfgs)[dockerID]
		currParam, err := (*DockerBuilderConfig)(&currCfg).ToParam(DockerBuilderConfig((*defaultCfg)[dockerID]), dockerBuilderFactory)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to generate parameter for dist configuration %s", dockerID)
		}
		dockerBuilderParamsMap[dockerID] = currParam
	}
	return dockerBuilderParamsMap, nil
}

type DockerBuilderConfig v0.DockerBuilderConfig

func ToDockerBuilderConfig(in DockerBuilderConfig) v0.DockerBuilderConfig {
	return (v0.DockerBuilderConfig)(in)
}

func (cfg *DockerBuilderConfig) ToParam(defaultCfg DockerBuilderConfig, dockerBuilderFactory distgo.DockerBuilderFactory) (distgo.DockerBuilderParam, error) {
	dockerBuilderType := getConfigStringValue(cfg.Type, defaultCfg.Type, "")
	if dockerBuilderType == "" {
		return distgo.DockerBuilderParam{}, errors.Errorf("type must be non-empty")
	}
	dockerBuilder, err := newDockerBuilder(dockerBuilderType, getConfigValue(cfg.Config, defaultCfg.Config, nil).(yaml.MapSlice), dockerBuilderFactory)
	if err != nil {
		return distgo.DockerBuilderParam{}, err
	}

	contextDir := getConfigStringValue(cfg.ContextDir, defaultCfg.ContextDir, "")
	if contextDir == "" {
		return distgo.DockerBuilderParam{}, errors.Errorf("context-dir must be non-empty")
	}
	tagTemplates := getConfigValue(cfg.TagTemplates, defaultCfg.TagTemplates, nil).([]string)
	if len(tagTemplates) == 0 {
		return distgo.DockerBuilderParam{}, errors.Errorf("tag-templates must be non-empty")
	}

	return distgo.DockerBuilderParam{
		DockerBuilder:    dockerBuilder,
		DockerfilePath:   getConfigStringValue(cfg.DockerfilePath, defaultCfg.DockerfilePath, "Dockerfile"),
		ContextDir:       contextDir,
		InputProductsDir: getConfigStringValue(cfg.InputProductsDir, defaultCfg.InputProductsDir, ""),
		InputBuilds:      getConfigValue(cfg.InputBuilds, defaultCfg.InputBuilds, nil).([]distgo.ProductBuildID),
		InputDists:       getConfigValue(cfg.InputDists, defaultCfg.InputDists, nil).([]distgo.ProductDistID),
		TagTemplates:     tagTemplates,
	}, nil
}

func newDockerBuilder(dockerBuilderType string, cfgYML yaml.MapSlice, dockerBuilderFactory distgo.DockerBuilderFactory) (distgo.DockerBuilder, error) {
	if dockerBuilderType == "" {
		return nil, errors.Errorf("type must be non-empty")
	}
	if dockerBuilderFactory == nil {
		return nil, errors.Errorf("dockerBuilderFactory must be provided")
	}
	cfgYMLBytes, err := yaml.Marshal(cfgYML)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to marshal configuration")
	}
	return dockerBuilderFactory.NewDockerBuilder(dockerBuilderType, cfgYMLBytes)
}
