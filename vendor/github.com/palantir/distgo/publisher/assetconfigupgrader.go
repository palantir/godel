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

package publisher

import (
	"encoding/base64"
	"os/exec"

	"github.com/pkg/errors"

	"github.com/palantir/distgo/assetapi"
)

type assetConfigUpgrader struct {
	typeName  string
	assetPath string
}

func (u *assetConfigUpgrader) TypeName() string {
	return u.typeName
}

func (u *assetConfigUpgrader) UpgradeConfig(config []byte) ([]byte, error) {
	upgradeConfigCmd := exec.Command(u.assetPath, "upgrade-config", base64.StdEncoding.EncodeToString(config))
	output, err := upgradeConfigCmd.CombinedOutput()
	if err != nil {
		return nil, assetapi.UpgradeConfigError(err, output)
	}
	decodedBytes, err := base64.StdEncoding.DecodeString(string(output))
	if err != nil {
		return nil, errors.Wrapf(err, "failed to decode base64")
	}
	return decodedBytes, nil
}
