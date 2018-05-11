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

package v0

import (
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"

	"github.com/palantir/distgo/distgo"
)

type DistConfig struct {
	// OutputDir specifies the default distribution output directory for product distributions created by the "dist"
	// task. The distribution output directory is written to
	// "{{OutputDir}}/{{ID}}/{{Version}}/{{NameTemplate}}", and the distribution artifacts are written to
	// "{{OutputDir}}/{{ID}}/{{Version}}".
	//
	// If a value is not specified, "out/dist" is used as the default value.
	OutputDir *string `yaml:"output-dir,omitempty"`

	// Disters is the configuration for the disters for this product. The YAML representation can be a single DisterConfig
	// or a map[DistID]DisterConfig.
	Disters *DistersConfig `yaml:"disters,omitempty"`
}

type DisterConfig struct {
	// Type is the type of the dister. This field must be non-nil and non-empty and resolve to a valid Dister.
	Type *string `yaml:"type,omitempty"`

	// Config is the YAML configuration content for the dister.
	Config *yaml.MapSlice `yaml:"config,omitempty"`

	// NameTemplate is the template used for the executable output. The following template parameters can be used in the
	// template:
	//   * {{Product}}: the name of the product.
	//   * {{Version}}: the version of the project.
	//
	// If a value is not specified, "{{Product}}-{{Version}}" is used as the default value.
	NameTemplate *string `yaml:"name-template,omitempty"`

	// Script is the content of a script that is written to file a file and run after the initial distribution
	// process but before the artifact generation process. The content of this value is written to a file and executed
	// with the project directory as the working directory. The script process inherits the environment variables of the
	// Go process and also has dist-related environment variables. Refer to the documentation for the
	// distgo.DistScriptEnvVariables function for the extra environment variables.
	Script *string `yaml:"script,omitempty"`
}

type DistersConfig map[distgo.DistID]DisterConfig

func (cfgs *DistersConfig) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var single DisterConfig
	if err := unmarshal(&single); err == nil && single.Type != nil {
		// only consider single configuration valid if it unmarshals and "type" key is explicitly specified
		*cfgs = DistersConfig{
			distgo.DistID(*single.Type): single,
		}
		return nil
	}

	var multiple map[distgo.DistID]DisterConfig
	if err := unmarshal(&multiple); err != nil {
		return errors.Errorf("failed to unmarshal configuration as single DisterConfig or as map[DistID]DisterConfig")
	}
	if len(multiple) == 0 {
		return errors.Errorf(`if "dist" key is specified, there must be at least one dist`)
	}
	*cfgs = multiple
	return nil
}
