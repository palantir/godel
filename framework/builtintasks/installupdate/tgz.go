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
	"archive/tar"
	"compress/gzip"
	"io"
	"os"
	"path"
	"sort"
	"strings"

	"github.com/pkg/errors"

	"github.com/palantir/godel/v2/framework/builtintasks/installupdate/layout"
)

func verifyPackageTgz(tgzFile string) (string, error) {
	entries, err := getPathsInTgz(tgzFile)
	if err != nil {
		return "", errors.Wrapf(err, "failed to get directories in tgz file %s", tgzFile)
	}
	actualPaths := sortedKeys(entries)
	version, err := versionFromEntries(actualPaths)
	if err != nil {
		return "", errors.Wrapf(err, "could not determine version from tgz file entries")
	}

	expectedPaths := layout.AppSpec().Paths(layout.AppSpecTemplate(version), false)
	for _, currExpectedPath := range expectedPaths {
		if _, ok := entries[currExpectedPath]; !ok {
			return "", errors.Errorf("tgz %s does not contain a valid package: failed to find %s.\nRequired: %v\nActual:   %v", tgzFile, currExpectedPath, expectedPaths, actualPaths)
		}
	}

	return version, nil
}

func versionFromEntries(sortedEntries []string) (string, error) {
	dirEntry := sortedEntries[0]
	expectedPrefix := layout.AppName + "-"
	if !strings.HasPrefix(dirEntry, expectedPrefix) {
		return "", errors.Errorf("entry %s in %v did not have expected prefix %s", dirEntry, sortedEntries, expectedPrefix)
	}
	return dirEntry[len(expectedPrefix):], nil
}

// getPathsInTgz returns a map that contains the paths to the entries present in the specified tgz file. The returned
// map will also contain subdirectories as independent entries -- that is, if the archive contains
// "root/intermediate/leaf.txt", the returned map will have "root", "root/intermediate" and "root/intermediate/leaf.txt"
// as entries.
func getPathsInTgz(tgzFile string) (rPaths map[string]bool, rErr error) {
	file, err := os.Open(tgzFile)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to open file %s", tgzFile)
	}
	defer func() {
		if err := file.Close(); err != nil && rErr == nil {
			rErr = errors.Wrapf(err, "failed to close file %s in defer", tgzFile)
		}
	}()

	gzf, err := gzip.NewReader(file)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to create gzip reader for file %s", tgzFile)
	}

	tarReader := tar.NewReader(gzf)
	dirs := make(map[string]bool)
	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		} else if err != nil {
			return nil, errors.Wrapf(err, "failed to read entry in file %s", tgzFile)
		}

		switch header.Typeflag {
		case tar.TypeDir:
			dirs[path.Dir(header.Name)] = true
		case tar.TypeReg:
			dirs[header.Name] = true
		default:
		}
	}
	return dirs, nil
}

func sortedKeys(input map[string]bool) []string {
	output := make([]string, 0, len(input))
	for key := range input {
		output = append(output, key)
	}
	sort.Strings(output)
	return output
}
