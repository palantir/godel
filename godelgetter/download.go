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

package godelgetter

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"

	"github.com/pkg/errors"
	"gopkg.in/cheggaaa/pb.v1"
)

// DownloadIntoDirectory downloads the provided package into the specified output directory. The output directory must
// already exist. The download progress is written to the provided writer. Returns the path to the downloaded file.
// If the path of pkgSrc and the download destination both refer to the same file (and the file exists), then if a
// checksum is provided it is verified (otherwise it is a no-op).
func DownloadIntoDirectory(pkgSrc PkgSrc, dstDir string, w io.Writer) (rPkg string, rErr error) {
	if dstDirInfo, err := os.Stat(dstDir); err != nil {
		if os.IsNotExist(err) {
			return "", errors.Wrapf(err, "destination directory %s does not exist", dstDir)
		}
		return "", errors.Wrapf(err, "failed to stat download directory")
	} else if !dstDirInfo.IsDir() {
		return "", errors.Errorf("destination path %s exists, but is not a directory", dstDir)
	}

	dstFilePath := path.Join(dstDir, pkgSrc.Name())
	if !pkgSrc.Same(dstFilePath) {
		// download the source package to the destination
		if err := Download(pkgSrc, dstFilePath, w); err != nil {
			return "", err
		}
	} else if wantChecksum := pkgSrc.Checksum(); wantChecksum != "" {
		// destination file and source file are the same -- if expected checksum was provided, verify that the checksum
		// for the existing file matches the expected checksum.
		checksum, err := computeSHA256Checksum(dstFilePath)
		if err != nil {
			return "", errors.Wrapf(err, "failed to compute checksum of %s", dstFilePath)
		}
		if checksum != wantChecksum {
			return "", errors.Errorf("checksum of %s does not match provided checksum: expected %s, was %s", dstFilePath, wantChecksum, checksum)
		}
	}
	return dstFilePath, nil
}

// Download downloads the provided package to the specified path. The parent directory of the path must exist. If the
// destination file already exists, it is overwritten. The download progress is written to the provided writer.
func Download(pkgSrc PkgSrc, dstFilePath string, w io.Writer) (rErr error) {
	wantChecksum := pkgSrc.Checksum()
	if info, err := os.Stat(dstFilePath); err == nil {
		if info.IsDir() {
			return errors.Errorf("destination path %s already exists and is a directory", dstFilePath)
		}
		if wantChecksum != "" {
			// if file already exists at destination and checksum is known, verify checksum of existing file.
			// If it matches, use existing file.
			checksum, err := computeSHA256Checksum(dstFilePath)
			if err != nil {
				return errors.Wrapf(err, "failed to compute checksum of %s", dstFilePath)
			}
			if checksum == wantChecksum {
				return nil
			}
		}
	}

	// open reader from source
	r, size, err := pkgSrc.Reader()
	if err != nil {
		return err
	}
	defer func() {
		if err := r.Close(); err != nil && rErr == nil {
			rErr = errors.Wrapf(err, "failed to close reader for %s in defer", pkgSrc.Path())
		}
	}()

	// create new file for package (overwrite any existing file)
	dstFile, err := os.Create(dstFilePath)
	if err != nil {
		return errors.Wrapf(err, "failed to create file %s", dstFilePath)
	}
	defer func() {
		if err := dstFile.Close(); err != nil && rErr == nil {
			rErr = errors.Wrapf(err, "failed to close file %s in defer", dstFilePath)
		}
	}()

	h := sha256.New()
	mw := io.MultiWriter(h, dstFile)

	fmt.Fprintf(w, "Getting package from %v...\n", pkgSrc.Path())
	if err := copyWithProgress(mw, r, size, w); err != nil {
		return errors.Wrapf(err, "failed to copy package %s to %s", pkgSrc.Path(), dstFilePath)
	}

	// verify checksum if provided
	if wantChecksum != "" {
		actualChecksum := hex.EncodeToString(h.Sum(nil))
		if wantChecksum != actualChecksum {
			return errors.Errorf("SHA-256 checksum of downloaded package did not match expected checksum: expected %s, was %s", wantChecksum, actualChecksum)
		}
	}
	return nil
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

func computeSHA256Checksum(filename string) (string, error) {
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return "", errors.Wrapf(err, "failed to read file %s", filename)
	}
	sha256Checksum := sha256.Sum256(bytes)
	return hex.EncodeToString(sha256Checksum[:]), nil
}
