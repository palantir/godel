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

package installupdate

import (
	"io"

	"github.com/pkg/errors"
)

// RunActionAndUpgradeConfig determines the gödel version of the project before the action, performs the action,
// determines the version after the action and runs the "upgrade-config" task depending on the version before and after
// the action.
//
// If the major version after the action is <2, then no action is performed (version 2 is the first version that
// supports configuration upgrades).
//
// If the major version after the action is >=2 and the major version before the action was <2, then "upgrade-config" is
// called with the "--legacy" flag set.
//
// If the major version before and after the action are both >=2, then the "upgrade-config" task is run if the version
// after the action is >= the version before the action.
//
// If the skipUpgradeConfig variable is true, then the provided action is run without any extra work.
func RunActionAndUpgradeConfig(projectDir string, skipUpgradeConfig bool, action func() error, stdout, stderr io.Writer) error {
	// if skipUpgradeConfig is true, just run the action and return its result
	if skipUpgradeConfig {
		return action()
	}

	godelVersionBeforeUpdate, err := getGodelVersion(projectDir)
	if err != nil {
		return errors.Wrapf(err, "failed to determine version before update")
	}

	if err := action(); err != nil {
		return err
	}

	godelVersionAfterUpdate, err := getGodelVersion(projectDir)
	if err != nil {
		return errors.Wrapf(err, "failed to determine version after update")
	}

	// if version after update is <2, then nothing to be done (upgrade-config task does not exist)
	if godelVersionAfterUpdate.MajorVersionNum() < 2 {
		return nil
	}

	// if version after update is >=2 and version before update was <2, then run "upgrade-config" in "--legacy" mode
	if godelVersionBeforeUpdate.MajorVersionNum() < 2 {
		return RunUpgradeLegacyConfig(projectDir, stdout, stderr)
	}

	// if version after the action is >=2 and versions are not comparable or version after the action is newer than or
	// equal to the version before the action, upgrade configuration
	if cmp, ok := godelVersionAfterUpdate.CompareTo(godelVersionBeforeUpdate); !ok || cmp >= 0 {
		return RunUpgradeConfig(projectDir, stdout, stderr)
	}

	// versions are comparable and version after action is older than version before action: nothing to do
	return nil
}
