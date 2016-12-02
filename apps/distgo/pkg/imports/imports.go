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

package imports

import (
	"go/build"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
)

// GoFiles is a map from package paths to the names of the buildable .go source files (.go files excluding Cgo and test
// files) in the package.
type GoFiles map[string][]string

// NewerThan returns true if the modification time of any of the GoFiles is newer than that of the provided file.
func (g GoFiles) NewerThan(fi os.FileInfo) (bool, error) {
	for pkg, files := range g {
		for _, goFile := range files {
			currPath := path.Join(pkg, goFile)
			currFi, err := os.Stat(currPath)
			if err != nil {
				return false, errors.Wrapf(err, "Failed to stat file %v", currPath)
			}
			if currFi.ModTime().After(fi.ModTime()) {
				return true, nil
			}
		}
	}
	return false, nil
}

// AllFiles returns a map that contains all of the non-standard library Go files that are imported (and thus required to
// build) the specified package (including the package itself). The keys in the returned map are the paths to the
// packages and the values are a slice of the names of the .go source files in the package (excluding Cgo and test
// files).
func AllFiles(pkgPath string) (GoFiles, error) {
	// package name to all non-test Go files in the package
	pkgFiles := make(map[string][]string)

	absPkgPath, err := filepath.Abs(pkgPath)
	if err != nil {
		return nil, errors.Wrapf(err, "Failed to convert %v to absolute path", pkgPath)
	}

	pkgsToProcess := []string{
		absPkgPath,
	}

	for len(pkgsToProcess) > 0 {
		currPkg := pkgsToProcess[0]
		pkgsToProcess = pkgsToProcess[1:]
		if _, ok := pkgFiles[currPkg]; ok {
			continue
		}

		// parse current package
		pkg, err := build.Import(".", currPkg, build.ImportComment)
		if err != nil {
			return nil, errors.Wrapf(err, "Failed to import package %v", currPkg)
		}

		// add all files for the current package to output
		pkgFiles[currPkg] = pkg.GoFiles

		// convert all non-built-in imports into packages and add to packages to process
		for importPath := range pkg.ImportPos {
			if !strings.Contains(importPath, ".") {
				// if import is a standard package, skip
				continue
			}
			importPkg, err := build.Import(importPath, currPkg, build.ImportComment)
			if err != nil {
				return nil, errors.Wrapf(err, "Failed to import package %v using srcDir %v", importPath, currPkg)
			}
			pkgsToProcess = append(pkgsToProcess, importPkg.Dir)
		}
	}
	return GoFiles(pkgFiles), nil
}
