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

package godel

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path"
	"strings"

	"github.com/nmiyake/archiver"
	"github.com/nmiyake/pkg/dirs"
	"github.com/palantir/pkg/specdir"
	"github.com/pkg/errors"
	"gopkg.in/cheggaaa/pb.v1"

	"github.com/palantir/godel/layout"
	"github.com/palantir/godel/properties"
)

// Copies and installs the gödel package from the provided PkgSrc. If the PkgSrc includes a checksum, this
// function will check to see if a TGZ file for the version as already been downloaded and if the checksum matches. If
// it does, that file will be used. Otherwise, the TGZ will be downloaded from the specified location and the downloaded
// TGZ will be verified against the checksum. If the checksum is empty, no verification will occur. If the install
// succeeds, the following files will be created:
// "{{layout.GödelHomePath()}}/downloads/{{layout.AppName}}-{{version}}.tgz" and
// "{{layout.GödelHomePath()}}/dists/{{layout.AppName}}-{{version}}". If the downloaded distribution matches a version
// that already exists in the distribution directory and a download occurs, the existing distribution will be
// overwritten by the newly downloaded one. Returns the version of the distribution that was installed.
func install(src PkgSrc, stdout io.Writer) (string, error) {
	gödelHomeSpecDir, err := layout.GödelHomeSpecDir(specdir.Create)
	if err != nil {
		return "", errors.Wrapf(err, "failed to create SpecDir for gödel home")
	}
	gödelHome := gödelHomeSpecDir.Root()

	downloadsDir := gödelHomeSpecDir.Path(layout.DownloadsDir)
	tgzFilePath, err := getPkg(src, downloadsDir, stdout)
	if err != nil {
		return "", err
	}

	tgzVersion, err := verifyPackageTgz(tgzFilePath)
	if err != nil {
		return "", errors.Wrapf(err, "downloaded file %s is not a valid %s package", tgzFilePath, layout.AppName)
	}

	// create temporary directory in gödel home in which downloaded tgz is expanded. If verification is successful,
	// the expanded directory will be moved to the destination.
	tmpDir, cleanup, err := dirs.TempDir(gödelHome, "")
	defer cleanup()
	if err != nil {
		return "", errors.Wrapf(err, "failed to create temporary directory rooted at %s", gödelHome)
	}

	if err := archiver.UntarGz(tgzFilePath, tmpDir); err != nil {
		return "", errors.Wrapf(err, "failed to extract archive %s to %s", tgzFilePath, tmpDir)
	}

	expandedGödelDir := path.Join(tmpDir, layout.AppName+"-"+tgzVersion)
	expandedGödelApp, err := layout.AppSpecDir(expandedGödelDir, tgzVersion)
	if err != nil {
		return "", errors.Wrapf(err, "extracted archive layout did not match expected gödel layout")
	}

	version, err := getExecutableVersion(expandedGödelApp)
	if err != nil {
		return "", errors.Wrapf(err, "failed to get version of downloaded gödel package")
	}

	if version != tgzVersion {
		return "", errors.Errorf("version reported by executable does not match version specified by tgz: expected %s, was %s", tgzVersion, version)
	}

	gödelDist, err := layout.GödelDistLayout(version, specdir.Create)
	if err != nil {
		return "", errors.Wrapf(err, "failed to create distribution directory")
	}
	gödelDirDestPath := gödelDist.Path(layout.AppDir)

	// delete destination directory if it already exists
	if _, err := os.Stat(gödelDirDestPath); !os.IsNotExist(err) {
		if err != nil {
			return "", errors.Wrapf(err, "failed to stat %s", gödelDirDestPath)
		}

		if err := os.RemoveAll(gödelDirDestPath); err != nil {
			return "", errors.Wrapf(err, "failed to remove %s", gödelDirDestPath)
		}
	}

	if err := os.Rename(expandedGödelDir, gödelDirDestPath); err != nil {
		return "", errors.Wrapf(err, "failed to rename %s to %s", expandedGödelDir, gödelDirDestPath)
	}

	return version, nil
}

// GetDistPkgInfo returns the distribution URL and checksum (if it exists) from the configuration file in the provided
// directory. Returns an error if the URL cannot be read.
func GetDistPkgInfo(configDir string) (PkgWithChecksum, error) {
	propsFilePath := path.Join(configDir, fmt.Sprintf("%v.properties", layout.AppName))
	props, err := properties.Read(propsFilePath)
	if err != nil {
		return PkgWithChecksum{}, errors.Wrapf(err, "failed to read properties file %s", propsFilePath)
	}
	url, err := properties.Get(props, properties.URL)
	if err != nil {
		return PkgWithChecksum{}, errors.Wrapf(err, "failed to get URL")
	}
	checksum, _ := properties.Get(props, properties.Checksum)
	return PkgWithChecksum{
		Pkg:      url,
		Checksum: checksum,
	}, nil
}

// getExecutableVersion gets the version of gödel contained in the provided root gödel directory. Invokes the executable
// for the current platform with the "--version" flag and returns the version determined by that output.
func getExecutableVersion(gödelApp specdir.SpecDir) (string, error) {
	executablePath := gödelApp.Path(layout.AppExecutable)
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

// getPkg gets the source package from the specified source and copies it to a new file in the specified directory
// (which must already exist). Returns the path to the downloaded file.
func getPkg(src PkgSrc, destDir string, stdout io.Writer) (rPkg string, rErr error) {
	expectedChecksum := src.checksum()

	if destDirInfo, err := os.Stat(destDir); err != nil {
		if os.IsNotExist(err) {
			return "", errors.Wrapf(err, "destination directory %s does not exist", destDir)
		}
		return "", errors.WithStack(err)
	} else if !destDirInfo.IsDir() {
		return "", errors.Errorf("destination path %s exists, but is not a directory", destDir)
	}

	destFilePath := path.Join(destDir, src.name())
	if info, err := os.Stat(destFilePath); err == nil {
		if info.IsDir() {
			return "", errors.Errorf("destination path %s already exists and is a directory", destFilePath)
		}
		if expectedChecksum != "" {
			// if tgz already exists at destination and checksum is known, verify checksum of existing tgz.
			// If it matches, use existing file.
			checksum, err := sha256Checksum(destFilePath)
			if err != nil {
				// if checksum computation fails, print error but continue execution
				fmt.Fprintf(stdout, "Failed to compute checksum of %s: %v\n", destFilePath, err)
			} else if checksum == expectedChecksum {
				return destFilePath, nil
			}
		}
	}

	// create new file for package (overwrite any existing file)
	destFile, err := os.Create(destFilePath)
	if err != nil {
		return "", errors.Wrapf(err, "failed to create file %s", destFilePath)
	}
	defer func() {
		if err := destFile.Close(); err != nil && rErr == nil {
			rErr = errors.Wrapf(err, "failed to close file %s in defer", destFilePath)
		}
	}()

	r, size, err := src.getPkg()
	if err != nil {
		return "", err
	}
	defer func() {
		if err := r.Close(); err != nil && rErr == nil {
			rErr = errors.Wrapf(err, "failed to close reader for %s in defer", src.path())
		}
	}()

	h := sha256.New()
	w := io.MultiWriter(h, destFile)

	fmt.Fprintf(stdout, "Getting package from %v...\n", src.path())
	if err := copyWithProgress(w, r, size, stdout); err != nil {
		return "", errors.Wrapf(err, "failed to copy package %s to %s", src.path(), destFilePath)
	}

	// verify checksum if provided
	if expectedChecksum != "" {
		actualChecksum := hex.EncodeToString(h.Sum(nil))
		if expectedChecksum != actualChecksum {
			return "", errors.Errorf("SHA-256 checksum of downloaded package did not match expected checksum: expected %s, was %s", expectedChecksum, actualChecksum)
		}
	}

	return destFilePath, nil
}

func copyWithProgress(w io.Writer, r io.Reader, dataLen int64, stdout io.Writer) error {
	bar := pb.New64(dataLen).SetUnits(pb.U_BYTES)
	bar.SetMaxWidth(120)
	bar.Output = stdout
	bar.Start()
	defer func() {
		bar.Finish()
	}()
	mw := io.MultiWriter(w, bar)
	_, err := io.Copy(mw, r)
	return err
}

type PkgSrc interface {
	// returns a reader that can be used to read the package and the size of the package. Reader will be open and
	// ready for reads -- the caller is responsible for closing the reader when done.
	getPkg() (io.ReadCloser, int64, error)
	// returns the name of this package.
	name() string
	// returns the path to this package.
	path() string
	// returns the expected SHA-256 checksum for the package. If this function returns an empty string, then a
	// checksum will not be performed.
	checksum() string
}

type PkgWithChecksum struct {
	Pkg      string
	Checksum string
}

func (p PkgWithChecksum) ToPkgSrc() PkgSrc {
	if strings.HasPrefix(p.Pkg, "http://") || strings.HasPrefix(p.Pkg, "https://") {
		return remotePkg(p)
	}
	return localPkg(p)
}

type remotePkg PkgWithChecksum

func (p remotePkg) getPkg() (io.ReadCloser, int64, error) {
	url := p.Pkg
	response, err := http.Get(url)
	if err != nil {
		return nil, 0, errors.Wrapf(err, "get call for url %s failed", url)
	}
	if response.StatusCode >= 400 {
		return nil, 0, errors.Errorf("request for URL %s returned status code %d", url, response.StatusCode)
	}
	return response.Body, response.ContentLength, nil
}

func (p remotePkg) name() string {
	return p.Pkg[strings.LastIndex(p.Pkg, "/")+1:]
}

func (p remotePkg) path() string {
	return p.Pkg
}

func (p remotePkg) checksum() string {
	return p.Checksum
}

type localPkg PkgWithChecksum

func (p localPkg) getPkg() (io.ReadCloser, int64, error) {
	pathToLocalTgz := p.Pkg
	localTgzFileInfo, err := os.Stat(pathToLocalTgz)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, 0, errors.Errorf("%s does not exist", pathToLocalTgz)
		}
		return nil, 0, errors.WithStack(err)
	} else if localTgzFileInfo.IsDir() {
		return nil, 0, errors.Errorf("%s is a directory", pathToLocalTgz)
	}
	srcTgzFile, err := os.Open(pathToLocalTgz)
	if err != nil {
		return nil, 0, errors.Wrapf(err, "failed to open %s", pathToLocalTgz)
	}
	return srcTgzFile, localTgzFileInfo.Size(), nil
}

func (p localPkg) name() string {
	return path.Base(p.Pkg)
}

func (p localPkg) path() string {
	return p.Pkg
}

func (p localPkg) checksum() string {
	return p.Checksum
}

func sha256Checksum(filename string) (string, error) {
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return "", errors.Wrapf(err, "failed to read file %s", filename)
	}
	sha256Checksum := sha256.Sum256(bytes)
	return hex.EncodeToString(sha256Checksum[:]), nil
}
