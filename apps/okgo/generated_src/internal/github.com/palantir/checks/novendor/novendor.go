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

package amalgomated

import (
	"fmt"
	"go/build"
	"io"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"

	"github.com/nmiyake/pkg/dirs"
	"github.com/nmiyake/pkg/errorstringer"
	"github.com/palantir/pkg/cli"
	"github.com/palantir/pkg/cli/flag"
	"github.com/palantir/pkg/matcher"
	"github.com/palantir/pkg/pkgpath"
	"github.com/pkg/errors"
)

const (
	pkgsFlagName		= "pkgs"
	projectPkgFlagName	= "project-package"
	fullPathFlagName	= "full"
)

var (
	pkgsFlag	= flag.StringSlice{
		Name:	pkgsFlagName,
		Usage:	"paths to the pacakges to check",
	}
	projectPkgFlag	= flag.BoolFlag{
		Name:	projectPkgFlagName,
		Usage:	"use the 'project' paradigm to interpret packages and only output projects that are unused",
		Value:	true,
	}
	fullPathFlag	= flag.BoolFlag{
		Name:	fullPathFlagName,
		Alias:	"f",
		Usage:	"include full path of unused packages (default omits path to vendor directory)",
	}
)

func AmalgomatedMain() {
	app := cli.NewApp(cli.DebugHandler(errorstringer.SingleStack))
	app.Flags = append(app.Flags,
		projectPkgFlag,
		fullPathFlag,
		pkgsFlag,
	)
	app.Action = func(ctx cli.Context) error {
		wd, err := dirs.GetwdEvalSymLinks()
		if err != nil {
			return errors.Wrapf(err, "Failed to get working directory")
		}
		return doNovendor(wd, ctx.Slice(pkgsFlagName), ctx.Bool(projectPkgFlagName), ctx.Bool(fullPathFlagName), ctx.App.Stdout)
	}
	os.Exit(app.Run(os.Args))
}

type pkgWithSrc struct {
	pkg	string
	src	string
}

func doNovendor(projectDir string, pkgPaths []string, groupPkgsByProject, fullPath bool, w io.Writer) error {
	if !path.IsAbs(projectDir) {
		return errors.Errorf("projectDir %s must be an absolute path", projectDir)
	}

	gopath := os.Getenv("GOPATH")
	if gopath == "" {
		return errors.Errorf("GOPATH environment variable must be set")
	}

	if relPath, err := filepath.Rel(path.Join(gopath, "src"), projectDir); err != nil || strings.HasPrefix(relPath, "../") {
		return errors.Errorf("Project directory %s must be a subdirectory of $GOPATH/src (%s)", projectDir, path.Join(gopath, "src"))
	}

	if len(pkgPaths) == 0 {
		// exclude vendor directories
		matcher := matcher.Any(pkgpath.DefaultGoPkgExcludeMatcher(), matcher.Name("vendor"))
		pkgs, err := pkgpath.PackagesInDir(projectDir, matcher)
		if err != nil {
			return errors.Wrapf(err, "Failed to list packages")
		}

		pkgPaths, err = pkgs.Paths(pkgpath.Relative)
		if err != nil {
			return errors.Wrapf(err, "Failed to convert package paths")
		}
	}

	pkgsToProcess := make([]pkgWithSrc, len(pkgPaths))
	for i, pkgPath := range pkgPaths {
		pkgsToProcess[i] = pkgWithSrc{
			pkg:	".",
			src:	path.Join(projectDir, pkgPath),
		}
	}

	unusedPkgs, err := getUnusedVendoredPkgs(projectDir, pkgsToProcess, groupPkgsByProject, fullPath)
	if err != nil {
		return errors.Wrapf(err, "Failed to determine unused packages")
	}
	if len(unusedPkgs) > 0 {
		fmt.Fprintln(w, strings.Join(unusedPkgs, "\n"))
		return fmt.Errorf("")
	}

	return nil
}

func getUnusedVendoredPkgs(projectDir string, pkgsToProcess []pkgWithSrc, groupPkgsByProject, fullPath bool) ([]string, error) {
	allVendoredPkgs, err := getAllVendoredPkgs(projectDir)
	if err != nil {
		return nil, err
	}

	allProjectPkgs := make(map[string]bool)
	for _, currPkg := range pkgsToProcess {
		examinedImports := make(map[string]bool)
		imps, err := getAllImports(currPkg.pkg, currPkg.src, projectDir, examinedImports, true)
		if err != nil {
			return nil, errors.Wrapf(err, "getAllFailed")
		}
		for k, v := range imps {
			allProjectPkgs[k] = v
		}
	}

	var unusedVendorPkgs []string
	if groupPkgsByProject {
		// do package-level grouping
		allProjectPkgsGrouped := make(map[string]bool)
		for k := range allProjectPkgs {
			vendorPath, nonVendorFullPath := splitPathOnVendor(k)
			vendoredRepoOrgProjectPath := path.Join(vendorPath, repoOrgProjectPath(nonVendorFullPath))
			allProjectPkgsGrouped[vendoredRepoOrgProjectPath] = true
		}

		usedKeys := make(map[string]bool)
		for k := range allVendoredPkgs {
			vendorPath, nonVendorFullPath := splitPathOnVendor(k)
			vendoredRepoOrgProjectPath := path.Join(vendorPath, repoOrgProjectPath(nonVendorFullPath))
			if !allProjectPkgsGrouped[vendoredRepoOrgProjectPath] && !usedKeys[vendoredRepoOrgProjectPath] {
				unusedVendorPkgs = append(unusedVendorPkgs, vendoredRepoOrgProjectPath)
				usedKeys[vendoredRepoOrgProjectPath] = true
			}
		}
	} else {
		for k := range allVendoredPkgs {
			if !allProjectPkgs[k] {
				unusedVendorPkgs = append(unusedVendorPkgs, k)
			}
		}
	}

	if !fullPath {
		// if fullPath is false, remove vendor portion from output
		for i, pkgName := range unusedVendorPkgs {
			_, pkgName = splitPathOnVendor(pkgName)
			unusedVendorPkgs[i] = pkgName
		}
	}
	sort.Strings(unusedVendorPkgs)
	return unusedVendorPkgs, nil
}

func getAllVendoredPkgs(projectRoot string) (map[string]bool, error) {
	vendoredPkgs := make(map[string]bool)
	err := filepath.Walk(projectRoot, func(currPath string, info os.FileInfo, err error) error {
		if info.IsDir() {
			rel, err := filepath.Rel(projectRoot, currPath)
			if err != nil {
				return err
			}
			inVendorDir := false
			for _, currPart := range strings.Split(rel, "/") {
				if currPart == "vendor" {
					inVendorDir = true
					break
				}
			}
			if inVendorDir {
				// if this is a directory in a vendor directory, attempt to parse as a package
				pkg, err := doImport(".", currPath, build.ImportComment)
				// record import path if package could be parsed and import path is not "." (which can
				// happen for some directories like testdata which cannot be imported)
				if err == nil && pkg.ImportPath != "." {
					vendoredPkgs[pkg.ImportPath] = true
				}
			}
		}
		return nil
	})
	if err != nil {
		return nil, errors.Wrapf(err, "Failed to determine vendored packages")
	}
	return vendoredPkgs, nil
}

// getAllImports takes an import and returns all of the packages that it imports (excluding standard library packages).
// Includes all transitive imports and the package of the import itself. Assumes that the import occurs in a package in
// "srcDir". If the "test" parameter is "true", considers all imports in the test files for the package as well.
func getAllImports(importPkgPath, srcDir, projectRoot string, examinedImports map[string]bool, includeTests bool) (map[string]bool, error) {
	importedPkgs := make(map[string]bool)
	if !strings.Contains(importPkgPath, ".") {
		// if package is a standard package, return empty
		return nil, nil
	}

	// ignore error because doImport returns partial object even on error. As long as an ImportPath is present,
	// proceed with determining imports.
	pkg, _ := doImport(importPkgPath, srcDir, build.ImportComment)
	if pkg.ImportPath == "" {
		return nil, nil
	}

	// skip import if package has already been examined
	if examinedImports[pkg.ImportPath] {
		return importedPkgs, nil
	}

	importedPkgs[pkg.ImportPath] = true
	examinedImports[pkg.ImportPath] = true

	currPkgImports := pkg.Imports
	if rel, err := filepath.Rel(projectRoot, pkg.Dir); err == nil && !strings.HasPrefix(rel, "../") {
		// if import is internal, update "srcDir" to be pkg.Dir to ensure that resolution is done against the
		// last internal package that was encountered
		srcDir = pkg.Dir
		if includeTests {
			// if import is internal and includeTests is true, consider imports from test files
			currPkgImports = append(currPkgImports, pkg.TestImports...)
			currPkgImports = append(currPkgImports, pkg.XTestImports...)
		}
	}

	// add packages from imports (don't examine transitive test dependencies)
	for _, currImport := range currPkgImports {
		if examinedImports[currImport] {
			continue
		}

		currImportedPkgs, err := getAllImports(currImport, srcDir, projectRoot, examinedImports, false)
		if err != nil {
			return nil, errors.Wrapf(err, "isExternalImport failed for %v", currImport)
		}
		examinedImports[currImport] = true

		for k, v := range currImportedPkgs {
			importedPkgs[k] = v
		}
	}
	return importedPkgs, nil
}

// takes the provided input, splits it on the path separator and returns the path up to the last "vendor" directory as
// the first return value and the path after the last "vendor" directory as the second return value. For example, if
// "foo/bar/vendor/inner/vendor/github.com/org/repo" is provided as input, the output is ("foo/bar/vendor/inner/vendor",
// "github.com/org/repo").
func splitPathOnVendor(pkgPath string) (string, string) {
	// get last index of "vendor"
	pathParts := strings.Split(pkgPath, "/")
	vendorIndex := -1
	for i := len(pathParts) - 1; i >= 0; i-- {
		if pathParts[i] == "vendor" {
			vendorIndex = i
			break
		}
	}
	return strings.Join(pathParts[:vendorIndex+1], "/"), strings.Join(pathParts[vendorIndex+1:], "/")
}

// returns the path that contains at most the first 3 elements of the package path. In most schemes, this will
// correspond to the source repository, organization and project ("gibhub.com/user/repo", "golang.org/x/crypto").
// If the path is shorter than 3 portions (for example, "gopkg.in/yaml.v2"), the path will be returned as-is. Does not
// do any semantic analysis, so 3 portions will be returned even if logically the repository is only 2 levels. For
// example, if "gopkg.in/project/subpackage" is provided as input, the first 3 parts of the path will be returned even
// though conceptually it represents a subpackage.
func repoOrgProjectPath(pkgPath string) string {
	_, pkgPath = splitPathOnVendor(pkgPath)
	pathParts := strings.Split(pkgPath, "/")
	lastIdx := len(pathParts)
	if lastIdx > 3 {
		lastIdx = 3
	}
	return strings.Join(pathParts[:lastIdx], "/")
}

// allContext is a build.Context based on build.Default that has "UseAllFiles" set to true. Makes it such that analysis
// is done on all Go files rather than on just those that match the default build context.
var allContext = getAllContext()

func getAllContext() build.Context {
	ctx := build.Default
	ctx.UseAllFiles = true
	return ctx
}

func doImport(path, srcDir string, mode build.ImportMode) (*build.Package, error) {
	return allContext.Import(path, srcDir, mode)
}
