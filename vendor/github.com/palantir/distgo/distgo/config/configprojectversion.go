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
	"github.com/palantir/distgo/projectversioner/git"
)

type ProjectVersionConfig v0.ProjectVersionConfig

func ToProjectVersionConfig(in *ProjectConfig) *v0.ProjectConfig {
	return (*v0.ProjectConfig)(in)
}

// ToParam returns the ProjectVersionerParam represented by the receiver *ProjectVersionConfig. If the receiver is nil,
// the git project versioner is returned.
func (cfg *ProjectVersionConfig) ToParam(projectVersionerFactory distgo.ProjectVersionerFactory) (distgo.ProjectVersionerParam, error) {
	if cfg == nil {
		// if configuration is nil, return git versioner as default
		return distgo.ProjectVersionerParam{
			ProjectVersioner: git.New(),
		}, nil
	}
	projectVersioner, err := newProjectVersioner(cfg.Type, cfg.Config, projectVersionerFactory)
	if err != nil {
		return distgo.ProjectVersionerParam{}, err
	}
	return distgo.ProjectVersionerParam{
		ProjectVersioner: projectVersioner,
	}, nil
}

func newProjectVersioner(projectVersionerType string, cfgYML yaml.MapSlice, projectVersionerFactory distgo.ProjectVersionerFactory) (distgo.ProjectVersioner, error) {
	if projectVersionerType == "" {
		return nil, errors.Errorf("project versioner type must be non-empty")
	}
	if projectVersionerFactory == nil {
		return nil, errors.Errorf("projectVersionerFactory must be provided")
	}
	cfgYMLBytes, err := yaml.Marshal(cfgYML)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to marshal configuration")
	}
	return projectVersionerFactory.NewProjectVersioner(projectVersionerType, cfgYMLBytes)
}
