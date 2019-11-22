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
	v0 "github.com/palantir/godel/v2/framework/godel/config/internal/v0"
	"github.com/palantir/godel/v2/pkg/versionedconfig"
	"github.com/pkg/errors"
)

func UpgradeConfig(cfgBytes []byte) ([]byte, error) {
	// legacy configuration is fully compatible with 2.0 configuration so no need to migrate. Configuration for godel is
	// also special because it has already been loaded by the time the program is run.
	version, err := versionedconfig.ConfigVersion(cfgBytes)
	if err != nil {
		return nil, err
	}
	switch version {
	case "", "0":
		return v0.UpgradeConfig(cfgBytes)
	default:
		return nil, errors.Errorf("unsupported version: %s", version)
	}
}
