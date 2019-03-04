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
	"crypto/sha256"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"

	"github.com/palantir/pkg/matcher"
)

type ChecksumSet struct {
	RootDir   string
	Checksums map[string]FileChecksumInfo
}

func (c *ChecksumSet) SortedKeys() []string {
	var sorted []string
	for k := range c.Checksums {
		sorted = append(sorted, k)
	}
	sort.Strings(sorted)
	return sorted
}

func (c *ChecksumSet) Diff(other ChecksumSet) ChecksumsDiff {
	diffs := make(map[string]string)

	// determine missing and extra entries
	for k := range c.Checksums {
		if _, ok := other.Checksums[k]; !ok {
			diffs[k] = "missing"
		}
	}
	for k := range other.Checksums {
		if _, ok := c.Checksums[k]; !ok {
			diffs[k] = "extra"
		}
	}

	// Diff content
	for k, v := range c.Checksums {
		otherV, ok := other.Checksums[k]
		if !ok {
			continue
		}

		if v.IsDir != otherV.IsDir {
			if v.IsDir {
				diffs[k] = "changed from directory to file"
			} else {
				diffs[k] = "changed from file to directory"
			}
			continue
		}
		if v.SHA256checksum != otherV.SHA256checksum {
			diffs[k] = fmt.Sprintf("checksum changed from %s to %s", v.SHA256checksum, otherV.SHA256checksum)
		}
	}
	return ChecksumsDiff{
		RootDir: c.RootDir,
		Diffs:   diffs,
	}
}

type ChecksumsDiff struct {
	RootDir string
	Diffs   map[string]string
}

func (c *ChecksumsDiff) String() string {
	var sortedKeys []string
	for k := range c.Diffs {
		sortedKeys = append(sortedKeys, k)
	}
	sort.Strings(sortedKeys)

	var parts []string
	for _, k := range sortedKeys {
		parts = append(parts, fmt.Sprintf("%s: %s", path.Join(c.RootDir, k), c.Diffs[k]))
	}
	return strings.Join(parts, "\n")
}

type FileChecksumInfo struct {
	Path           string
	IsDir          bool
	SHA256checksum string
}

func ChecksumsForMatchingPaths(rootDir string, m matcher.Matcher) (ChecksumSet, error) {
	pathsToChecksums := make(map[string]FileChecksumInfo)
	if err := filepath.Walk(rootDir, func(path string, info os.FileInfo, err error) error {
		relPath, err := filepath.Rel(rootDir, path)
		if err != nil {
			return err
		}
		if m == nil || m.Match(relPath) {
			checksum, err := newChecksum(path, info)
			if err != nil {
				return err
			}
			checksum.Path = relPath
			pathsToChecksums[relPath] = checksum
		}
		return nil
	}); err != nil {
		return ChecksumSet{}, fmt.Errorf("failed to walk directory %q: %v", rootDir, err)
	}
	return ChecksumSet{
		RootDir:   rootDir,
		Checksums: pathsToChecksums,
	}, nil
}

func newChecksum(filePath string, info os.FileInfo) (FileChecksumInfo, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return FileChecksumInfo{}, err
	}
	defer func() {
		// file is opened for reading only, so safe to ignore errors on close
		_ = f.Close()
	}()

	if info.IsDir() {
		return FileChecksumInfo{
			Path:  filePath,
			IsDir: true,
		}, nil
	}

	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return FileChecksumInfo{}, err
	}
	return FileChecksumInfo{
		Path:           filePath,
		SHA256checksum: fmt.Sprintf("%x", h.Sum(nil)),
	}, nil
}
