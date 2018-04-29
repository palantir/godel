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
	"github.com/palantir/distgo/distgo"
	"github.com/palantir/distgo/distgo/config/internal/v0"
)

type ProductConfig v0.ProductConfig

func ToProductConfig(in *ProductConfig) *v0.ProductConfig {
	return (*v0.ProductConfig)(in)
}

func (cfg *ProductConfig) ToParam(productID distgo.ProductID, scriptIncludes string, defaultCfg ProductConfig, disterFactory distgo.DisterFactory, dockerBuilderFactory distgo.DockerBuilderFactory) (distgo.ProductParam, error) {
	var buildParam *distgo.BuildParam
	if cfg.Build != nil {
		defaultBuildCfg := BuildConfig{}
		if defaultCfg.Build != nil {
			defaultBuildCfg = BuildConfig(*defaultCfg.Build)
		}
		buildParamVar, err := (*BuildConfig)(cfg.Build).ToParam(scriptIncludes, defaultBuildCfg)
		if err != nil {
			return distgo.ProductParam{}, err
		}
		buildParam = &buildParamVar
	}

	var runParam *distgo.RunParam
	if cfg.Run != nil {
		defaultRunCfg := RunConfig{}
		if defaultCfg.Run != nil {
			defaultRunCfg = RunConfig(*defaultCfg.Run)
		}
		runParamVar := (*RunConfig)(cfg.Run).ToParam(defaultRunCfg)
		runParam = &runParamVar
	}

	var distParam *distgo.DistParam
	if cfg.Dist != nil {
		var defaultDistConfig DistConfig
		if defaultCfg.Dist != nil {
			defaultDistConfig = DistConfig(*defaultCfg.Dist)
		}
		distParamsVar, err := (*DistConfig)(cfg.Dist).ToParam(scriptIncludes, defaultDistConfig, disterFactory)
		if err != nil {
			return distgo.ProductParam{}, err
		}
		distParam = &distParamsVar
	}

	var publishParam *distgo.PublishParam
	if cfg.Publish != nil {
		defaultPublishCfg := PublishConfig{}
		if defaultCfg.Publish != nil {
			defaultPublishCfg = PublishConfig(*defaultCfg.Publish)
		}
		publishParamVar, err := (*PublishConfig)(cfg.Publish).ToParam(defaultPublishCfg)
		if err != nil {
			return distgo.ProductParam{}, err
		}
		publishParam = &publishParamVar
	}

	var dockerParam *distgo.DockerParam
	if cfg.Docker != nil {
		defaultDockerCfg := DockerConfig{}
		if defaultCfg.Docker != nil {
			defaultDockerCfg = DockerConfig(*defaultCfg.Docker)
		}
		dockerParamVar, err := (*DockerConfig)(cfg.Docker).ToParam(defaultDockerCfg, dockerBuilderFactory)
		if err != nil {
			return distgo.ProductParam{}, err
		}
		dockerParam = &dockerParamVar
	}

	var firstLevelDeps []distgo.ProductID
	seen := make(map[distgo.ProductID]struct{})
	if cfg.Dependencies != nil {
		for _, currDep := range *cfg.Dependencies {
			if _, ok := seen[currDep]; ok {
				// do not add entry if it was already seen
				continue
			}
			seen[currDep] = struct{}{}
			firstLevelDeps = append(firstLevelDeps, currDep)
		}
	}
	return distgo.ProductParam{
		ID:                     productID,
		Build:                  buildParam,
		Run:                    runParam,
		Dist:                   distParam,
		Publish:                publishParam,
		Docker:                 dockerParam,
		FirstLevelDependencies: firstLevelDeps,
	}, nil
}
