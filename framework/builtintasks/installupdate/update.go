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
	"io/ioutil"
	"os"
	"path"

	"github.com/nmiyake/pkg/dirs"
	"github.com/palantir/pkg/specdir"
	"github.com/pkg/errors"

	"github.com/palantir/godel/framework/builtintasks/installupdate/layout"
	"github.com/palantir/godel/godelgetter"
)

// NewInstall performs a new installation of gödel in the specified directory using the specified file as the source.
// Calls "Install" to install the package provided as a parameter. Once the package is installed, the wrapper and
// settings files are copied from the newly downloaded distribution to the specified path. If there was a previous
// installation of gödel in the path, it is overwritten by the new file. However, changes in the "var" directory are
// purely additive -- files that have been added in this directory in the new distribution will be added, but existing
// files will not be modified or removed.
func NewInstall(dstDirPath, srcPkgPath string, stdout io.Writer) error {
	if err := layout.VerifyDirExists(dstDirPath); err != nil {
		return errors.Wrapf(err, "path %s does not specify an existing directory", dstDirPath)
	}
	if err := update(dstDirPath, godelgetter.NewPkgSrc(srcPkgPath, ""), true, stdout); err != nil {
		return errors.Wrapf(err, "failed to install from %s into %s", srcPkgPath, dstDirPath)
	}
	return nil
}

// Update updates gödel. Calls "Install" to download and install the package specified in the "{{properties.Url}}"
// property of the properties file for the provided gödel wrapper script. Once the package is installed, the provided
// wrapper file and its directory are overwritten with the files provided by the package that was downloaded. However,
// changes in the "var" directory are purely additive -- files that have been added in this directory in the new
// distribution will be added, but existing files will not be modified or removed.
func Update(wrapperScriptPath string, stdout io.Writer) error {
	wrapperScriptDir := path.Dir(wrapperScriptPath)
	wrapper, err := specdir.New(wrapperScriptDir, layout.WrapperSpec(), nil, specdir.Validate)
	if err != nil {
		return errors.Wrapf(err, "wrapper script %s is not in a valid location", wrapperScriptPath)
	}
	pkg, err := distPkgInfo(wrapper.Path(layout.WrapperConfigDir))
	if err != nil {
		return errors.Wrapf(err, "failed to get URL from properties file")
	}
	if err := update(wrapperScriptDir, pkg, false, stdout); err != nil {
		return errors.Wrapf(err, "failed to update")
	}
	return nil
}

// Returns the distribution URL and checksum (if it exists) from the configuration file in the provided directory.
// Returns an error if the URL cannot be read.
func distPkgInfo(configDir string) (godelgetter.PkgSrc, error) {
	propsFilePath := path.Join(configDir, fmt.Sprintf("%s.properties", layout.AppName))
	props, err := readPropertiesFile(propsFilePath)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to read properties file %s", propsFilePath)
	}
	url, ok := props[propertiesURLKey]
	if !ok {
		return nil, errors.Wrapf(err, "properties file %s does not contain key %s", propsFilePath, propertiesURLKey)
	}
	checksum := props[propertiesChecksumKey]
	return godelgetter.NewPkgSrc(url, checksum), nil
}

func update(wrapperScriptDir string, pkg godelgetter.PkgSrc, newInstall bool, stdout io.Writer) error {
	mode := specdir.Validate
	if newInstall {
		mode = specdir.SpecOnly
	}
	wrapper, err := specdir.New(wrapperScriptDir, layout.WrapperSpec(), nil, mode)
	if err != nil {
		return errors.Wrapf(err, "%s is not a valid wrapper directory", wrapperScriptDir)
	}

	version, err := install(pkg, stdout)
	if err != nil {
		return err
	}

	gödelDist, err := layout.GodelDistLayout(version, specdir.Validate)
	if err != nil {
		return errors.Wrapf(err, "unable to get gödel home directory")
	}

	tmpDir, cleanup, err := dirs.TempDir(wrapperScriptDir, "")
	defer cleanup()
	if err != nil {
		return errors.Wrapf(err, "failed to create temporary directory in %s", wrapperScriptDir)
	}

	// copy new wrapper script to temp directory on same device and then move to destination
	installedGödelWrapper := gödelDist.Path(layout.WrapperScriptFile)
	tmpGödelWrapper := path.Join(tmpDir, "godelw")
	if err := layout.CopyFile(installedGödelWrapper, tmpGödelWrapper); err != nil {
		return errors.Wrapf(err, "failed to copy %s to %s", installedGödelWrapper, tmpGödelWrapper)
	}

	if err := os.Rename(tmpGödelWrapper, wrapper.Path(layout.WrapperScriptFile)); err != nil {
		return errors.Wrapf(err, "failed to move wrapper script into place")
	}

	if newInstall {
		// if this is a new install, ensure that required destination paths exist
		if err := layout.WrapperSpec().CreateDirectoryStructure(wrapperScriptDir, nil, false); err != nil {
			return errors.Wrapf(err, "failed to ensure that required paths exist")
		}
	}

	// additively sync config directory
	if err := layout.SyncDirAdditive(gödelDist.Path(layout.WrapperConfigDir), wrapper.Path(layout.WrapperConfigDir)); err != nil {
		return errors.Wrapf(err, "failed to additively sync from %s to %s", gödelDist.Path(layout.WrapperConfigDir), wrapper.Path(layout.WrapperConfigDir))
	}

	// overlay all directories except "config"
	installedGödelWrapperDir := gödelDist.Path(layout.WrapperAppDir)
	wrapperDirFiles, err := ioutil.ReadDir(installedGödelWrapperDir)
	if err != nil {
		return errors.Wrapf(err, "failed to list files in directory %s", installedGödelWrapperDir)
	}

	for _, currWrapperFile := range wrapperDirFiles {
		syncSrcPath := path.Join(installedGödelWrapperDir, currWrapperFile.Name())
		syncDestPath := path.Join(wrapperScriptDir, layout.AppName, currWrapperFile.Name())

		if currWrapperFile.IsDir() && currWrapperFile.Name() == layout.WrapperConfigDir {
			// do not sync "config" directory
			continue
		} else {
			// if destination file exists, remove it
			if _, err := os.Stat(syncDestPath); err == nil || !os.IsNotExist(err) {
				if err := os.RemoveAll(syncDestPath); err != nil {
					return errors.Wrapf(err, "failed to remove %s", syncDestPath)
				}
			}

			// safe to copy
			var err error
			if currWrapperFile.IsDir() {
				err = layout.CopyDir(syncSrcPath, syncDestPath)
			} else {
				err = layout.CopyFile(syncSrcPath, syncDestPath)
			}
			if err != nil {
				return errors.Wrapf(err, "failed to copy %s to %s", syncSrcPath, syncDestPath)
			}
		}
	}
	return nil
}
