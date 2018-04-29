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

package legacy

import (
	"github.com/palantir/godel/pkg/versionedconfig"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"

	"github.com/palantir/distgo/dister/manual/config/internal/v0"
)

type Config struct {
	versionedconfig.ConfigWithLegacy `yaml:",inline"`
	Extension                        string `yaml:"extension"`
}

func UpgradeConfig(cfgBytes []byte) ([]byte, error) {
	var legacyCfg Config
	if err := yaml.UnmarshalStrict(cfgBytes, &legacyCfg); err != nil {
		return nil, errors.Wrapf(err, "failed to unmarshal manual dister legacy configuration")
	}
	upgradedCfg := v0.Config{
		Extension: legacyCfg.Extension,
	}
	upgradedCfgBytes, err := yaml.Marshal(upgradedCfg)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to marshal manual dister v0 configuration")
	}
	return upgradedCfgBytes, nil
}
