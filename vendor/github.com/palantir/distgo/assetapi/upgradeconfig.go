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

package assetapi

import (
	"os/exec"
	"strings"

	"github.com/pkg/errors"
)

func UpgradeConfigError(err error, output []byte) error {
	if err == nil {
		return nil
	}
	if _, ok := err.(*exec.ExitError); ok {
		output := strings.TrimSuffix(strings.Trim(string(output), "Error: "), "\n")
		return errors.Errorf("failed to upgrade asset configuration: %s", output)
	}
	return errors.Wrapf(err, "failed to upgrade asset configuration")
}
