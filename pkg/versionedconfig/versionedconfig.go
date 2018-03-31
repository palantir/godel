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

package versionedconfig

import (
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

// ConfigWithVersion is a struct with a "version" YAML field that stores the version as a string.
type ConfigWithVersion struct {
	Version string `yaml:"version,omitempty"`
}

// ConfigVersion unmarshals the provided bytes as YAML and returns the value of the top-level "version" key. The value
// of this key must be a string. If the input is valid YAML but does not contain a "version" key, returns the empty
// string. Returns an error if the input is not valid YAML or cannot be unmarshaled as YAML as specified.
func ConfigVersion(in []byte) (string, error) {
	var cfgWithVersion ConfigWithVersion
	if err := yaml.Unmarshal(in, &cfgWithVersion); err != nil {
		return "", errors.Wrapf(err, "failed to unmarshal YAML")
	}
	return cfgWithVersion.Version, nil
}
