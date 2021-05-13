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

package dirchecksum

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/termie/go-shutil"
)

// ChecksumsForDirAfterAction returns the Checksumset for the state of the target directory after running the provided
// action on it. The target directory is returned to its original state after the action is run. Uses the following
// process:
//
// * Copies the target directory to a temporary location (unique location in the same parent directory)
// * Moves the target directory to another temporary location (unique location in the same parent directory)
// * Moves the copied directory to the original location
// * Runs the provided action
// * Computes the checksums for the directory
// * Removes the directory
// * Moves the original target directory from the temporary location back to its original location
//
// The result of the above is that the checksums are computed for the directory after the action is run, but the target
// directory stays in its original state. This function registers a signal handler that restores the state to the
// original state on SIGINT or SIGTERM signals and then calls os.Exit(1).
func ChecksumsForDirAfterAction(dir string, action func(dir string) error) (ChecksumSet, error) {
	var origDirCopy string
	var origDirMoved string

	cleanupFn := mustDefer(func() {
		// remove copied directory
		if origDirCopy != "" {
			_ = os.RemoveAll(origDirCopy)
		}
		// move original directory back to original location
		if origDirMoved != "" {
			_ = os.RemoveAll(dir)
			_ = os.Rename(origDirMoved, dir)
		}
	})
	defer cleanupFn()

	// copy original directory to temporary location
	var err error
	origDirCopy, err = createTmpDirPath(filepath.Dir(dir))
	if err != nil {
		return ChecksumSet{}, err
	}
	if err := shutil.CopyTree(dir, origDirCopy, nil); err != nil {
		return ChecksumSet{}, fmt.Errorf("failed to copy directory: %v", err)
	}

	// move original directory to temporary location
	origDirMoved, err = createTmpDirPath(filepath.Dir(dir))
	if err != nil {
		return ChecksumSet{}, err
	}
	if err := os.Rename(dir, origDirMoved); err != nil {
		return ChecksumSet{}, fmt.Errorf("failed to move original output directory to temporary location: %v", err)
	}

	// move copied directory to original location
	if err := os.Rename(origDirCopy, dir); err != nil {
		return ChecksumSet{}, fmt.Errorf("failed to move copied output directory to original location: %v", err)
	}
	origDirCopy = ""

	if err := action(dir); err != nil {
		return ChecksumSet{}, err
	}

	newChecksums, err := ChecksumsForMatchingPaths(dir, nil)
	if err != nil {
		return ChecksumSet{}, fmt.Errorf("failed to compute new checksums: %v", err)
	}
	return newChecksums, nil
}

func createTmpDirPath(parentDir string) (string, error) {
	tmpDir, err := os.MkdirTemp(parentDir, "amalgomate-verify-")
	if err != nil {
		return "", fmt.Errorf("failed to create temporary directory: %v", err)
	}
	if err := os.RemoveAll(tmpDir); err != nil {
		return "", fmt.Errorf("failed to remove temporary directory: %v", err)
	}
	return tmpDir, nil
}
