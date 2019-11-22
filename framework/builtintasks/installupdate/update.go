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
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"strings"
	"time"

	"github.com/google/go-github/v28/github"
	"github.com/nmiyake/pkg/dirs"
	"github.com/palantir/godel/v2/framework/builtintasks/installupdate/layout"
	"github.com/palantir/godel/v2/godelgetter"
	"github.com/palantir/pkg/specdir"
	"github.com/pkg/errors"
)

// NewInstall performs a new installation of gödel in the specified directory using the specified package as the source.
// Calls "Install" to install the package provided as a parameter. Once the package is installed, the wrapper and
// settings files are copied from the newly downloaded distribution to the specified path. If there was a previous
// installation of gödel in the path, it is overwritten by the new file. However, changes in the "var" directory are
// purely additive -- files that have been added in this directory in the new distribution will be added, but existing
// files will not be modified or removed.
func NewInstall(dstDirPath string, srcPkg godelgetter.PkgSrc, stdout io.Writer) error {
	if err := layout.VerifyDirExists(dstDirPath); err != nil {
		return errors.Wrapf(err, "path %s does not specify an existing directory", dstDirPath)
	}
	if err := update(dstDirPath, srcPkg, true, stdout); err != nil {
		return errors.Wrapf(err, "failed to install from %s into %s", srcPkg.Path(), dstDirPath)
	}
	return nil
}

// Update updates the gödel installation in the specified directory to be the distribution specified by the provided
// package source. Once the package is installed, the existing wrapper script and its directory are overwritten with the
// files provided by the package that was downloaded. However, changes in the "godel" directory are purely additive --
// files that have been added in this directory in the new distribution will be added, but existing files will not be
// modified or removed.
func Update(projectDirPath string, srcPkg godelgetter.PkgSrc, stdout io.Writer) error {
	if err := update(projectDirPath, srcPkg, false, stdout); err != nil {
		return errors.Wrapf(err, "update failed")
	}
	return nil
}

// InstallVersion installs the specified version of gödel in the provided project directory. If targetVersion is the
// empty string, the latest version is determined and used.
func InstallVersion(projectDir, targetVersion, wantChecksum string, cacheValidDuration time.Duration, newInstall bool, stdout io.Writer) error {
	if targetVersion == "" {
		version, err := latestGodelVersion(cacheValidDuration)
		if err != nil {
			return err
		}
		targetVersion = version
	}
	pkgSrc, err := pkgSrcForVersion(targetVersion, wantChecksum)
	if err != nil {
		return err
	}

	var installFn func(string, godelgetter.PkgSrc, io.Writer) error
	if newInstall {
		installFn = NewInstall
	} else {
		installFn = Update
	}
	return installFn(projectDir, pkgSrc, stdout)
}

// pkgSrcForVersion returns a package source for the provided version. If the distribution for the provided version has
// been downloaded locally (and its checksum matches the expected checksum if one is provided), the package source uses
// the filesystem path. Otherwise, the package source specifies the Bintray download URL. Sets the provided checksum as
// the expected checksum for the package.
func pkgSrcForVersion(version, wantChecksum string) (godelgetter.PkgSrc, error) {
	if version == "" {
		return nil, errors.Errorf("version for package must be specified")
	}

	// consider distribution URL to be canonical source
	canonicalSrcPkgPath := fmt.Sprintf("https://palantir.bintray.com/releases/com/palantir/godel/godel/%s/godel-%s.tgz", version, version)

	pkgPath, checksum, err := downloadedTGZForVersion(version)
	if err != nil || (wantChecksum != "" && checksum != wantChecksum) {
		// if downloaded version was not present locally, fall back on canonical source
		pkgPath = canonicalSrcPkgPath
	}
	return godelgetter.NewPkgSrc(pkgPath, wantChecksum, godelgetter.PkgSrcCanonicalSourceParam(canonicalSrcPkgPath)), nil
}

// downloadedTGZForVersion returns the path and checksum for the downloaded TGZ for the specified version. Returns an
// error if the TGZ for the specified version does not exist (has not been downloaded).
func downloadedTGZForVersion(version string) (string, string, error) {
	godelHomeSpecDir, err := layout.GodelHomeSpecDir(specdir.Create)
	if err != nil {
		return "", "", errors.Wrapf(err, "failed to create SpecDir for gödel home")
	}
	downloadsDirPath := godelHomeSpecDir.Path(layout.DownloadsDir)
	downloadedTGZ := path.Join(downloadsDirPath, fmt.Sprintf("%s-%s.tgz", layout.AppName, version))
	if _, err := os.Stat(downloadedTGZ); err != nil {
		return "", "", errors.Wrapf(err, "failed to stat downloaded TGZ file")
	}
	checksum, err := layout.Checksum(downloadedTGZ)
	if err != nil {
		return "", "", errors.Wrapf(err, "failed to compute checksum")
	}
	return downloadedTGZ, checksum, nil
}

// latestGodelVersion returns the latest version of gödel. Does so by querying the GitHub API or looking up the value
// from cache. If a cache value is within the timeframe of the provided duration (time.Now - cacheExpiration), it is
// returned.
func latestGodelVersion(cacheExpiration time.Duration) (string, error) {
	if cacheExpiration != 0 {
		versionCfg, err := readLatestCachedVersion()
		if err == nil && storedLatestVersionValid(versionCfg, cacheExpiration) {
			return versionCfg.LatestVersion, nil
		}
	}
	client := github.NewClient(http.DefaultClient)
	rel, _, err := client.Repositories.GetLatestRelease(context.Background(), "palantir", "godel")
	if err != nil {
		return "", errors.Wrap(err, "failed to determine latest release")
	}
	latestVersion := *rel.TagName
	if len(latestVersion) >= 2 && latestVersion[0] == 'v' && latestVersion[1] >= '0' && latestVersion[1] <= '9' {
		// if version begins with 'v' and is followed by a digit, trim the leading 'v'
		latestVersion = latestVersion[1:]
	}
	if err := writeLatestCachedVersion(latestVersion); err != nil {
		return "", errors.Wrapf(err, "failed to write latest version to cache")
	}
	return latestVersion, nil
}

const latestVersionFileName = "latest-version.json"

func readLatestCachedVersion() (versionConfig, error) {
	godelHomeSpecDir, err := layout.GodelHomeSpecDir(specdir.Create)
	if err != nil {
		return versionConfig{}, errors.Wrapf(err, "failed to create SpecDir for gödel home")
	}
	cacheDirPath := godelHomeSpecDir.Path(layout.CacheDir)
	latestVersionFile := path.Join(cacheDirPath, latestVersionFileName)

	bytes, err := ioutil.ReadFile(latestVersionFile)
	if err != nil {
		return versionConfig{}, errors.Wrapf(err, "failed to read version file")
	}
	var versionCfg versionConfig
	if err := json.Unmarshal(bytes, &versionCfg); err != nil {
		return versionConfig{}, errors.Wrapf(err, "failed to unmarshal version file")
	}
	return versionCfg, nil
}

func writeLatestCachedVersion(version string) error {
	godelHomeSpecDir, err := layout.GodelHomeSpecDir(specdir.Create)
	if err != nil {
		return errors.Wrapf(err, "failed to create SpecDir for gödel home")
	}
	cacheDirPath := godelHomeSpecDir.Path(layout.CacheDir)
	latestVersionFile := path.Join(cacheDirPath, latestVersionFileName)

	bytes, err := json.Marshal(versionConfig{
		LatestVersion: version,
		Timestamp:     time.Now().Unix(),
	})
	if err != nil {
		return errors.Wrapf(err, "failed to marshal version config as JSON")
	}

	if err := ioutil.WriteFile(latestVersionFile, bytes, 0644); err != nil {
		return errors.Wrap(err, "failed to write version file")
	}
	return nil
}

func storedLatestVersionValid(cfg versionConfig, cacheExpiration time.Duration) bool {
	storedTime := time.Unix(cfg.Timestamp, 0)
	return !storedTime.Before(time.Now().Add(-1 * cacheExpiration))
}

type versionConfig struct {
	LatestVersion string `json:"latestVersion"`
	Timestamp     int64  `json:"timestamp"`
}

func setGodelPropertyKey(projectDir, key, val string) error {
	wrapperSpec, err := specdir.New(projectDir, layout.WrapperSpec(), nil, specdir.Validate)
	if err != nil {
		return errors.Wrapf(err, "unable to create wrapper spec")
	}
	configDir := wrapperSpec.Path(layout.WrapperConfigDir)

	propsFilePath := path.Join(configDir, fmt.Sprintf("%s.properties", layout.AppName))
	bytes, err := ioutil.ReadFile(propsFilePath)
	if err != nil {
		return errors.Wrapf(err, "failed to read properties file")
	}
	lines := strings.Split(string(bytes), "\n")
	for i, currLine := range lines {
		if !strings.HasPrefix(currLine, key+"=") {
			continue
		}
		lines[i] = key + "=" + val
	}
	output := strings.Join(lines, "\n")
	if err := ioutil.WriteFile(propsFilePath, []byte(output), 0644); err != nil {
		return errors.Wrapf(err, "failed to write properties file")
	}
	return nil
}

// GodelPropsDistPkgInfo returns a package that consists of the distribution URL and checksum (if it exists) from the
// gödel configuration file for the gödel installation in the provided project directory.
func GodelPropsDistPkgInfo(projectDir string) (godelgetter.PkgSrc, error) {
	wrapperSpec, err := specdir.New(projectDir, layout.WrapperSpec(), nil, specdir.Validate)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to create wrapper spec")
	}
	configDir := wrapperSpec.Path(layout.WrapperConfigDir)

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

	godelDist, err := layout.GodelDistLayout(version, specdir.Validate)
	if err != nil {
		return errors.Wrapf(err, "unable to get gödel home directory")
	}

	tmpDir, cleanup, err := dirs.TempDir(wrapperScriptDir, "")
	defer cleanup()
	if err != nil {
		return errors.Wrapf(err, "failed to create temporary directory in %s", wrapperScriptDir)
	}

	// copy new wrapper script to temp directory on same device and then move to destination
	installedGodelWrapper := godelDist.Path(layout.WrapperScriptFile)
	tmpGodelWrapper := path.Join(tmpDir, "godelw")
	if err := layout.CopyFile(installedGodelWrapper, tmpGodelWrapper); err != nil {
		return errors.Wrapf(err, "failed to copy %s to %s", installedGodelWrapper, tmpGodelWrapper)
	}

	if err := os.Rename(tmpGodelWrapper, wrapper.Path(layout.WrapperScriptFile)); err != nil {
		return errors.Wrapf(err, "failed to move wrapper script into place")
	}

	if newInstall {
		// if this is a new install, ensure that required destination paths exist
		if err := layout.WrapperSpec().CreateDirectoryStructure(wrapperScriptDir, nil, false); err != nil {
			return errors.Wrapf(err, "failed to ensure that required paths exist")
		}
	}

	// additively sync config directory
	if err := layout.SyncDirAdditive(godelDist.Path(layout.WrapperConfigDir), wrapper.Path(layout.WrapperConfigDir)); err != nil {
		return errors.Wrapf(err, "failed to additively sync from %s to %s", godelDist.Path(layout.WrapperConfigDir), wrapper.Path(layout.WrapperConfigDir))
	}

	// overlay all directories except "config"
	installedGodelWrapperDir := godelDist.Path(layout.WrapperAppDir)
	wrapperDirFiles, err := ioutil.ReadDir(installedGodelWrapperDir)
	if err != nil {
		return errors.Wrapf(err, "failed to list files in directory %s", installedGodelWrapperDir)
	}

	for _, currWrapperFile := range wrapperDirFiles {
		syncSrcPath := path.Join(installedGodelWrapperDir, currWrapperFile.Name())
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

	// update values in godel.properties
	canonicalSrc := pkg.CanonicalSource()
	if canonicalSrc == "" {
		canonicalSrc = pkg.Path()
	}
	if err := setGodelPropertyKey(wrapperScriptDir, propertiesURLKey, canonicalSrc); err != nil {
		return errors.Wrap(err, "failed to update URL in godel properties file")
	}
	if err := setGodelPropertyKey(wrapperScriptDir, propertiesChecksumKey, pkg.Checksum()); err != nil {
		return errors.Wrap(err, "failed to update checksum in godel properties file")
	}
	return nil
}
