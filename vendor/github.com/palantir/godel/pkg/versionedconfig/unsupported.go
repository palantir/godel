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

// ConfigNotSupported verifies that the provided bytes represent empty YAML. If the YAML is non-empty, return an error.
func ConfigNotSupported(name string, cfgBytes []byte) ([]byte, error) {
	var mapSlice yaml.MapSlice
	if err := yaml.Unmarshal(cfgBytes, &mapSlice); err != nil {
		return nil, errors.Wrapf(err, "failed to unmarshal %s configuration as yaml.MapSlice", name)
	}
	if len(mapSlice) != 0 {
		return nil, errors.Errorf("%s does not currently support configuration", name)
	}
	return cfgBytes, nil
}
