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

package gocd

import (
	"fmt"
	"go/build"
	"go/token"
	"io/ioutil"
	"strings"

	"github.com/pkg/errors"
)

type PkgInfo struct {
	// import path for package ("_test" is appended if the package represents the tests files)
	Path string
	// name of package
	Name string
	// number of .go files in the package directory
	NGoFiles int
	// importPath of all of the packages imported by the package. If usage information was retrieved, the value is
	// a set that contains the files in the package that imported the package; otherwise, it is nil.
	Imports map[string]map[string]struct{}
}

type PkgMode bool

const (
	Default         = false
	Test    PkgMode = true
)

func (m PkgMode) empty(pkg *build.Package) bool {
	var numFiles int
	switch m {
	case Default:
		numFiles = len(pkg.GoFiles)
	case Test:
		numFiles = len(pkg.TestGoFiles) + len(pkg.XTestGoFiles)
	default:
		panic(fmt.Sprintf("unhandled mode: %v", m))
	}
	return numFiles == 0
}

func (m PkgMode) imports(pkg *build.Package) []string {
	switch m {
	case Default:
		return pkg.Imports
	case Test:
		return append(pkg.TestImports, pkg.XTestImports...)
	default:
		panic(fmt.Sprintf("unhandled mode: %v", m))
	}
}

func (m PkgMode) importPos(pkg *build.Package) map[string][]token.Position {
	switch m {
	case Default:
		return pkg.ImportPos
	case Test:
		return combine(pkg.TestImportPos, pkg.XTestImportPos)
	default:
		panic(fmt.Sprintf("unhandled mode: %v", m))
	}
}

// DirPkgInfo returns a PkgInfo for the package in the specified srcDir using the specified mode. If the mode is
// Default, the package information is that of the non-test files in the package, while if it is Test, it is the
// information for the test files (internal and external) in the package. The package information is obtained by running
// a local import (".") for the package from its own directory. If the mode is Test, the path of the returned package
// will have "_test" appended to it to differentiate it from the non-test package.
func DirPkgInfo(srcDir string, mode PkgMode) (PkgInfo, bool, error) {
	return ImportPkgInfo(".", srcDir, mode)
}

// ImportPkgInfo returns a PkgInfo for the package specified by importPkgPath imported from srcPkgDir using the
// specified mode. If the mode is Default, the package information is that of the non-test files in the package, while
// if it is Test, it is the information for the test files (internal and external) in the package. The package
// information is obtained by running an import for importPkgPath from the srcPkgDir directory, which is equivalent to
// an import statement `import "importPkgPath"` in a package located in srcPkgDir. If the package resolved from that
// location is a vendored package, the path will be the vendored import path. If the mode is Test, the path of the
// returned package will have "_test" appended to it to differentiate it from the non-test package.
func ImportPkgInfo(importPkgPath, srcPkgDir string, mode PkgMode) (PkgInfo, bool, error) {
	// get information for package
	pkg, err := doImport(importPkgPath, srcPkgDir)
	if err != nil {
		return PkgInfo{}, false, err
	}

	pkgImportPath := pkg.ImportPath
	// if test package info, append "_test" to the import path to differentiate the test package from the non-test
	// package
	if mode == Test {
		pkgImportPath += "_test"
	}

	// get number of Go files in this package
	nGoFiles, err := nGoFiles(pkg)
	if err != nil {
		return PkgInfo{}, false, err
	}

	imports := make(map[string]map[string]struct{})
	for k, v := range importsWithLocs(mode.importPos(pkg)) {
		// translate import path to actual path used by project (for example, may be in a vendor directory)
		pkg, err := doImport(k, srcPkgDir)
		if err != nil {
			return PkgInfo{}, false, err
		}
		imports[pkg.ImportPath] = v
	}

	pi := PkgInfo{
		Path:     pkgImportPath,
		Name:     pkg.Name,
		NGoFiles: nGoFiles,
		Imports:  imports,
	}

	return pi, mode.empty(pkg), nil
}

func importsWithLocs(posMap map[string][]token.Position) map[string]map[string]struct{} {
	info := make(map[string]map[string]struct{})
	for k, v := range posMap {
		if isStdLibImport(k) {
			continue
		}
		files := make(map[string]struct{}, len(v))
		for i := range v {
			files[v[i].Filename] = struct{}{}
		}
		info[k] = files
	}
	return info
}

func combine(maps ...map[string][]token.Position) map[string][]token.Position {
	combined := make(map[string][]token.Position)
	for _, m := range maps {
		for k, v := range m {
			combined[k] = append(combined[k], v...)
		}
	}
	return combined
}

// nGoFiles returns the number of Go files in the provided package. Returns the number of files in the package directory
// whose name has the suffix ".go".
func nGoFiles(pkg *build.Package) (int, error) {
	fis, err := ioutil.ReadDir(pkg.Dir)
	if err != nil {
		return 0, errors.Errorf("failed to determine number of Go files in %s: %v", pkg.Dir, err)
	}
	nGoFiles := 0
	for _, fi := range fis {
		if !fi.IsDir() && strings.HasSuffix(fi.Name(), ".go") {
			nGoFiles++
		}
	}
	return nGoFiles, nil
}

// allContext is a build.Context based on build.Default that has "UseAllFiles" set to true. Makes it such that analysis
// is done on all Go files rather than on just those that match the default build context.
var allContext = getAllContext()

func getAllContext() build.Context {
	ctx := build.Default
	ctx.UseAllFiles = true
	return ctx
}

func doImport(path, srcDir string) (*build.Package, error) {
	pkg, err := allContext.Import(path, srcDir, build.ImportComment)
	if err != nil {
		if _, ok := err.(*build.MultiplePackageError); ok {
			// if error is multiple packages, re-try using default context (build tags may be used to
			// exclude packages)
			if pkg, err := build.Import(path, srcDir, build.ImportComment); err == nil {
				return pkg, nil
			}
		}
		return pkg, errors.Wrapf(err, "failed to import package %s using srcDir %s", path, srcDir)
	}
	return pkg, nil
}

func isStdLibImport(pkg string) bool {
	return !strings.Contains(pkg, ".")
}
