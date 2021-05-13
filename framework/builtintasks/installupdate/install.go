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
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/mholt/archiver/v3"
	"github.com/nmiyake/pkg/dirs"
	"github.com/palantir/godel/v2/framework/builtintasks/installupdate/layout"
	"github.com/palantir/godel/v2/godelgetter"
	"github.com/palantir/pkg/specdir"
	"github.com/pkg/errors"
	"github.com/rogpeppe/go-internal/lockedfile"
)

// Copies and installs the gödel package from the provided PkgSrc. If the PkgSrc includes a checksum, this
// function will check to see if a TGZ file for the version as already been downloaded and if the checksum matches. If
// it does, that file will be used. Otherwise, the TGZ will be downloaded from the specified location and the downloaded
// TGZ will be verified against the checksum. If the checksum is empty, no verification will occur. If the install
// succeeds, the following files will be created:
// "{{layout.GodelHomePath()}}/downloads/{{layout.AppName}}-{{version}}.tgz" and
// "{{layout.GodelHomePath()}}/dists/{{layout.AppName}}-{{version}}". If the downloaded distribution matches a version
// that already exists in the distribution directory and a download occurs, the existing distribution will be
// overwritten by the newly downloaded one. Returns the version of the distribution that was installed.
//
// Locks on a file in the Godel Home "downloads" directory with a file name derived based on the provided PkgSrc to
// ensure that this operation does not run concurrently for the same package.
func install(src godelgetter.PkgSrc, stdout io.Writer) (string, error) {
	godelHomeSpecDir, err := layout.GodelHomeSpecDir(specdir.Create)
	if err != nil {
		return "", errors.Wrapf(err, "failed to create SpecDir for gödel home")
	}
	godelHome := godelHomeSpecDir.Root()

	downloadsDir := godelHomeSpecDir.Path(layout.DownloadsDir)

	installPkgLockFilePath := filepath.Join(downloadsDir, fmt.Sprintf("install-%s.lock", src.Name()))
	installMutex := lockedfile.MutexAt(installPkgLockFilePath)
	unlockFn, err := installMutex.Lock()
	if err != nil {
		return "", errors.Wrapf(err, "failed to lock mutex for installing package")
	}
	defer unlockFn()

	tgzFilePath, err := godelgetter.DownloadIntoDirectory(src, downloadsDir, stdout)
	if err != nil {
		return "", err
	}

	tgzVersion, err := verifyPackageTgz(tgzFilePath)
	if err != nil {
		return "", errors.Wrapf(err, "downloaded file %s is not a valid %s package", tgzFilePath, layout.AppName)
	}

	// create temporary directory in gödel home in which downloaded tgz is expanded. If verification is successful,
	// the expanded directory will be moved to the destination.
	tmpDir, cleanup, err := dirs.TempDir(godelHome, "")
	defer cleanup()
	if err != nil {
		return "", errors.Wrapf(err, "failed to create temporary directory rooted at %s", godelHome)
	}

	if err := archiver.DefaultTarGz.Unarchive(tgzFilePath, tmpDir); err != nil {
		return "", errors.Wrapf(err, "failed to extract archive %s to %s", tgzFilePath, tmpDir)
	}

	expandedGodelDir := filepath.Join(tmpDir, layout.AppName+"-"+tgzVersion)
	expandedGodelApp, err := layout.AppSpecDir(expandedGodelDir, tgzVersion)
	if err != nil {
		return "", errors.Wrapf(err, "extracted archive layout did not match expected gödel layout")
	}

	version, err := getExecutableVersion(expandedGodelApp)
	if err != nil {
		return "", errors.Wrapf(err, "failed to get version of downloaded gödel package")
	}

	if version != tgzVersion {
		return "", errors.Errorf("version reported by executable does not match version specified by tgz: expected %s, was %s", tgzVersion, version)
	}

	godelDist, err := layout.GodelDistLayout(version, specdir.Create)
	if err != nil {
		return "", errors.Wrapf(err, "failed to create distribution directory")
	}
	godelDirDestPath := godelDist.Path(layout.AppDir)

	// delete destination directory if it already exists
	if _, err := os.Stat(godelDirDestPath); !os.IsNotExist(err) {
		if err != nil {
			return "", errors.Wrapf(err, "failed to stat %s", godelDirDestPath)
		}

		if err := os.RemoveAll(godelDirDestPath); err != nil {
			return "", errors.Wrapf(err, "failed to remove %s", godelDirDestPath)
		}
	}

	if err := os.Rename(expandedGodelDir, godelDirDestPath); err != nil {
		return "", errors.Wrapf(err, "failed to rename %s to %s", expandedGodelDir, godelDirDestPath)
	}

	return version, nil
}

// getExecutableVersion gets the version of gödel contained in the provided root gödel directory. Invokes the executable
// for the current platform with the "version" task and returns the version determined by that output.
func getExecutableVersion(godelApp specdir.SpecDir) (string, error) {
	executablePath := godelApp.Path(layout.AppExecutable)
	cmd := exec.Command(executablePath, "version")
	output, err := cmd.Output()
	if err != nil {
		return "", errors.Wrapf(err, "failed to execute command %v: %s", cmd.Args, string(output))
	}

	outputString := strings.TrimSpace(string(output))
	parts := strings.Split(outputString, " ")
	if len(parts) != 3 {
		return "", errors.Errorf(`expected output %s to have 3 parts when split by " ", but was %v`, outputString, parts)
	}
	return parts[2], nil
}
