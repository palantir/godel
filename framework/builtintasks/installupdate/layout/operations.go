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

package layout

import (
	"crypto/sha256"
	"encoding/hex"
	"io"
	"io/ioutil"
	"os"
	"path"

	"github.com/pkg/errors"
)

// Move the file or directory at src to dst. Conceptually, this is equivalent to executing "mv src dst". dst must not
// already exist and the path up to it must already exist. Uses os.Rename to execute the move, which means that there
// may be platform-specific restrictions such as not being able to move the directory between different volumes.
func Move(src, dst string) error {
	if err := verifyDstPathSafe(dst); err != nil {
		return errors.Wrapf(err, "cannot move directory to path %s", dst)
	}
	if err := os.Rename(src, dst); err != nil {
		return errors.Wrapf(err, "failed to rename %s to %s", src, dst)
	}
	return nil
}

// CopyDir recursively copies the src directory to the path specified by dst. dst must not already exist and the path up
// to it must already exist.
func CopyDir(src, dst string) error {
	if err := verifyDstPathSafe(dst); err != nil {
		return errors.Wrapf(err, "cannot copy directory to path %s", dst)
	}

	srcInfo, err := os.Stat(src)
	if err != nil {
		return errors.Wrapf(err, "failed to stat source directory %s", src)
	}

	if err := os.Mkdir(dst, srcInfo.Mode()); err != nil {
		return errors.Wrapf(err, "failed to create destination directory %s", dst)
	}

	files, err := ioutil.ReadDir(src)
	if err != nil {
		return errors.Wrapf(err, "failed to read directory %s", src)
	}

	for _, f := range files {
		srcPath := path.Join(src, f.Name())
		dstPath := path.Join(dst, f.Name())

		if f.IsDir() {
			err = CopyDir(srcPath, dstPath)
		} else {
			err = CopyFile(srcPath, dstPath)
		}

		if err != nil {
			return errors.Wrapf(err, "failed to copy %s to %s", srcPath, dstPath)
		}
	}

	return nil
}

// CopyFile copies the file at src to dst. src must specify a regular file (non-directory, no special mode) that exists,
// and dst must specify a path that does not yet exist, but whose parent directory does exist. The copied file will have
// the same permissions as the original.
func CopyFile(src, dst string) (rErr error) {
	srcInfo, err := os.Stat(src)
	if err != nil {
		return errors.Wrapf(err, "failed to stat source file %s", src)
	}
	if !srcInfo.Mode().IsRegular() {
		return errors.Wrapf(err, "source file %s is not a regular file, had mode: %v", src, srcInfo.Mode())
	}
	srcPerms := srcInfo.Mode().Perm()

	srcFile, err := os.Open(src)
	if err != nil {
		return errors.Wrapf(err, "failed to open %s", src)
	}
	defer func() {
		if err := srcFile.Close(); err != nil {
			rErr = errors.Wrapf(err, "failed to close %s in defer", src)
		}
	}()

	if err := verifyDstPathSafe(dst); err != nil {
		return errors.Wrapf(err, "cannot copy to destination path %s", dst)
	}

	dstFile, err := os.Create(dst)
	if err != nil {
		return errors.Wrapf(err, "failed to open %s", dst)
	}
	defer func() {
		if err := dstFile.Close(); err != nil {
			rErr = errors.Wrapf(err, "failed to close %s in defer", dstFile.Name())
		}
	}()

	if _, err := io.Copy(dstFile, srcFile); err != nil {
		return errors.Wrapf(err, "failed to copy %s to %s", src, dst)
	}

	if err := dstFile.Chmod(srcPerms); err != nil {
		return errors.Wrapf(err, "failed to chmod %s to have permissions %v", dst, srcPerms)
	}

	return nil
}

// SyncDir syncs the contents of the provided directories such that the content of dstDir matches that of srcDir (except
// for files that match a name in the provided skip slice).
func SyncDir(srcDir, dstDir string, skip []string) (bool, error) {
	modified := false

	srcFiles, err := ioutil.ReadDir(srcDir)
	if err != nil {
		return modified, errors.Wrapf(err, "failed to read directory %s", srcDir)
	}
	srcFilesMap := toMap(srcFiles)

	dstFiles, err := ioutil.ReadDir(dstDir)
	if err != nil {
		return modified, errors.Wrapf(err, "failed to read directory %s", dstDir)
	}
	dstFilesMap := toMap(dstFiles)

	skipSet := toSet(skip)
	for dstFileName, dstFileInfo := range dstFilesMap {
		if _, ok := skipSet[dstFileName]; ok {
			// skip if file name is in skip list
			continue
		}

		remove := false
		srcFilePath := path.Join(srcDir, dstFileName)
		dstFilePath := path.Join(dstDir, dstFileName)

		if currSrcFileInfo, ok := srcFilesMap[dstFileName]; !ok {
			// if dst exists but src does not, remove dst
			remove = true
		} else if dstFileInfo.IsDir() != currSrcFileInfo.IsDir() {
			// if dst file and src file are different types, remove dst
			remove = true
		} else if !dstFileInfo.IsDir() {
			srcChecksum, err := Checksum(srcFilePath)
			if err != nil {
				return modified, errors.Wrapf(err, "failed to compute checksum for %s", srcFilePath)
			}
			dstChecksum, err := Checksum(dstFilePath)
			if err != nil {
				return modified, errors.Wrapf(err, "failed to compute checksum for %s", dstFilePath)
			}

			// if dst and src are both files and their checksums differ, remove dst
			if srcChecksum != dstChecksum {
				remove = true
			}
		} else {
			// if dst and src both exist and are both directories, sync them recursively
			recursiveModified, err := SyncDir(srcFilePath, dstFilePath, skip)
			if err != nil {
				return modified, errors.Wrapf(err, "failed to sync %s with %s", dstFilePath, srcFilePath)
			}
			modified = modified || recursiveModified
		}

		// remove path if it was marked for removal
		if remove {
			if err := os.RemoveAll(dstFilePath); err != nil {
				return modified, errors.Wrapf(err, "failed to remove %s", dstFilePath)
			}
			modified = true
		}
	}

	for srcFileName, srcFileInfo := range srcFilesMap {
		if _, ok := skipSet[srcFileName]; ok {
			// skip if file name is in skip list
			continue
		}

		srcFilePath := path.Join(srcDir, srcFileName)
		dstFilePath := path.Join(dstDir, srcFileName)

		if _, err := os.Stat(dstFilePath); os.IsNotExist(err) {
			// if path does not exist at destination, copy source version
			var err error
			if srcFileInfo.IsDir() {
				err = CopyDir(srcFilePath, dstFilePath)
			} else {
				err = CopyFile(srcFilePath, dstFilePath)
			}
			if err != nil {
				return modified, errors.Wrapf(err, "failed to copy %s to %s", srcFilePath, dstFilePath)
			}
			modified = true
		} else if err != nil {
			return modified, errors.Wrapf(err, "failed to stat %s", dstFilePath)
		}
	}

	return modified, nil
}

func toSet(input []string) map[string]struct{} {
	s := make(map[string]struct{}, len(input))
	for _, curr := range input {
		s[curr] = struct{}{}
	}
	return s
}

func toMap(input []os.FileInfo) map[string]os.FileInfo {
	m := make(map[string]os.FileInfo, len(input))
	for _, curr := range input {
		m[curr.Name()] = curr
	}
	return m
}

func Checksum(p string) (string, error) {
	f, err := os.Open(p)
	if err != nil {
		return "", errors.Wrapf(err, "failed to open %s", p)
	}
	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", errors.Wrapf(err, "failed to copy file to hash buffer")
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}

// SyncDirAdditive copies all of the files and directories in src that are not in dst. Directories that are present in
// both are handled recursively. Basically a recursive merge with source preservation.
func SyncDirAdditive(src, dst string) error {
	srcInfos, err := ioutil.ReadDir(src)
	if err != nil {
		return errors.Wrapf(err, "failed to open %s", src)
	}

	for _, srcInfo := range srcInfos {
		srcPath := path.Join(src, srcInfo.Name())
		dstPath := path.Join(dst, srcInfo.Name())

		if dstInfo, err := os.Stat(dstPath); os.IsNotExist(err) {
			// safe to copy
			if srcInfo.IsDir() {
				err = CopyDir(srcPath, dstPath)
			} else {
				err = CopyFile(srcPath, dstPath)
			}
			if err != nil {
				return errors.Wrapf(err, "failed to copy %s to %s", srcPath, dstPath)
			}
		} else if err != nil {
			return errors.Wrapf(err, "failed to stat %s", dstPath)
		} else if srcInfo.IsDir() && dstInfo.IsDir() {
			// if source and destination are both directories, sync recursively
			if err = SyncDirAdditive(srcPath, dstPath); err != nil {
				return errors.Wrapf(err, "failed to sync %s to %s", srcPath, dstPath)
			}
		}
	}
	return nil
}

func VerifyDirExists(dir string) error {
	return verifyPath(dir, path.Base(dir), true, false)
}

func verifyPath(p, expectedName string, isDir bool, optional bool) error {
	if path.Base(p) != expectedName {
		return errors.Errorf("%s is not a path to %s", p, expectedName)
	}

	if fi, err := os.Stat(p); err != nil {
		if os.IsNotExist(err) {
			if !optional {
				return errors.Wrapf(err, "%s does not exist", p)
			}
			// path does not exist, but it is optional so is okay
			return nil
		}
		return errors.Wrapf(err, "failed to stat %s", p)
	} else if currIsDir := fi.IsDir(); currIsDir != isDir {
		return errors.Errorf("IsDir for %s returned wrong value: expected %v, was %v", p, isDir, currIsDir)
	}

	return nil
}

// verifyDstPathSafe verifies that the provided destination path does not exist, but that the path to its parent does.
func verifyDstPathSafe(dst string) error {
	if _, err := os.Stat(dst); !os.IsNotExist(err) {
		return errors.Wrapf(err, "destination path %s already exists", dst)
	}
	if _, err := os.Stat(path.Dir(dst)); os.IsNotExist(err) {
		return errors.Wrapf(err, "parent directory of destination path %s does not exist", dst)
	}
	return nil
}
