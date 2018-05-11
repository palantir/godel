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
	"gopkg.in/yaml.v2"
)

type ProjectVersionConfig struct {
	// Type is the type of the project versioner. This field must be non-empty and resolve to a valid ProjectVersioner.
	Type string `yaml:"type,omitempty"`

	// Config is the YAML configuration content for the project versioner.
	Config yaml.MapSlice `yaml:"config,omitempty"`
}
