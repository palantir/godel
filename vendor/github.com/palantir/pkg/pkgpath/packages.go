// Copyright (c) 2016 Palantir Technologies. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package pkgpath

import (
	"fmt"
	"go/build"
	"go/parser"
	"go/token"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"

	"github.com/palantir/pkg/matcher"
)

// DefaultGoPkgExcludeMatcher returns a matcher that matches names that standard Go tools generally exclude as Go
// packages. This includes hidden directories, directories named "testdata" and directories that start with an
// underscore.
func DefaultGoPkgExcludeMatcher() matcher.Matcher {
	return matcher.Any(matcher.Hidden(), matcher.Name("testdata"), matcher.Name("_.+"))
}

type Type int

const (
	// Absolute is the absolute path to a package, e.g.: /Volumes/git/go/src/github.com/org/project
	Absolute Type = iota
	// GoPathSrcRelative is the path to a package relative to "$GOPATH/src", e.g.: github.com/org/project. Is the
	// file path rather than the import path, so includes vendor directories in path, e.g.:
	// github.com/org/project/vendor/github.com/other/project
	GoPathSrcRelative
	// Relative is the relative path to a package relative to a directory. Always includes the "./" prefix, e.g.:
	// ./., ./app/main.
	Relative
)

func (t Type) String() string {
	switch t {
	case Absolute:
		return "Absolute"
	case GoPathSrcRelative:
		return "GoPathSrcRelative"
	case Relative:
		return "Relative"
	default:
		return fmt.Sprintf("%d", int(t))
	}
}

type PkgPather interface {
	Abs() string
	GoPathSrcRel() (string, error)
	Rel(root string) (string, error)
}

func NewAbsPkgPath(absPath string) PkgPather {
	return &pkgPath{
		pathType: Absolute,
		path:     absPath,
	}
}

func NewGoPathSrcRelPkgPath(goPathSrcRelPath string) PkgPather {
	return &pkgPath{
		pathType: GoPathSrcRelative,
		path:     goPathSrcRelPath,
	}
}

func NewRelPkgPath(relPath, baseDir string) PkgPather {
	return &pkgPath{
		pathType: Relative,
		path:     relPath,
		baseDir:  baseDir,
	}
}

type pkgPath struct {
	pathType Type
	path     string
	baseDir  string // only present if Type is Relative
}

func (p *pkgPath) Abs() string {
	switch p.pathType {
	case Absolute:
		return p.path
	case GoPathSrcRelative:
		return path.Join(os.Getenv("GOPATH"), "src", p.path)
	case Relative:
		return path.Join(p.baseDir, p.path)
	default:
		panic(fmt.Sprintf("unhandled case: %v", p.path))
	}
}

func (p *pkgPath) GoPathSrcRel() (string, error) {
	return relPathNoParentDir(p.Abs(), path.Join(os.Getenv("GOPATH"), "src"), "")
}

func (p *pkgPath) Rel(baseDir string) (string, error) {
	return relPathNoParentDir(p.Abs(), baseDir, "./")
}

func relPathNoParentDir(absPath, baseDir, prepend string) (string, error) {
	const parentDirPath = ".." + string(filepath.Separator)
	relPath, err := filepath.Rel(baseDir, absPath)
	if err != nil {
		return "", err
	}
	if strings.HasPrefix(relPath, parentDirPath) {
		return "", fmt.Errorf("resolving %s against base %s produced relative path starting with %s: %s", absPath, baseDir, parentDirPath, relPath)
	}
	return prepend + relPath, nil
}

type packages struct {
	rootDir string
	// key is absolute package path, value is package name
	pkgs map[string]string
}

type Packages interface {
	RootDir() string
	Packages(pathType Type) (map[string]string, error)
	Paths(pathType Type) ([]string, error)
	Filter(exclude matcher.Matcher) (Packages, error)
}

// Filter returns a Packages object that contains all of the packages that do not match the provided matcher.
func (p *packages) Filter(exclude matcher.Matcher) (Packages, error) {
	allPkgsRelPaths, err := p.Packages(Relative)
	if err != nil {
		return nil, fmt.Errorf("failed to get relative paths for packages: %v", err)
	}

	filteredAbsPathPkgs := make(map[string]string)
	for currPkgRelPath, currPkg := range allPkgsRelPaths {
		if exclude == nil || !exclude.Match(currPkgRelPath) {
			filteredAbsPathPkgs[path.Join(p.rootDir, currPkgRelPath)] = currPkg
		}
	}

	return createPkgsWithValidation(p.rootDir, filteredAbsPathPkgs)
}

func (p *packages) RootDir() string {
	return p.rootDir
}

func (p *packages) Packages(pathType Type) (map[string]string, error) {
	pkgs := make(map[string]string, len(p.pkgs))
	for currPath, currPkg := range p.pkgs {
		pkgs[currPath] = currPkg
	}

	var f func(string) (string, error)
	switch pathType {
	case Absolute:
		return pkgs, nil
	case GoPathSrcRelative:
		f = func(absPath string) (string, error) {
			return NewAbsPkgPath(absPath).GoPathSrcRel()
		}
	case Relative:
		f = func(absPath string) (string, error) {
			return NewAbsPkgPath(absPath).Rel(p.rootDir)
		}
	default:
		return nil, fmt.Errorf("unrecognized path type: %v", pathType)
	}

	relPathsMap := make(map[string]string, len(pkgs))
	for currAbsPath, currPkg := range pkgs {
		currRelPath, err := f(currAbsPath)
		if err != nil {
			return nil, fmt.Errorf("unable to get relative path for %s: %v", currAbsPath, err)
		}
		relPathsMap[currRelPath] = currPkg
	}
	return relPathsMap, nil
}

func (p *packages) Paths(pathType Type) ([]string, error) {
	pkgs, err := p.Packages(pathType)
	if err != nil {
		return nil, err
	}
	pkgPaths := make([]string, 0, len(pkgs))
	for currPath := range pkgs {
		pkgPaths = append(pkgPaths, currPath)
	}
	sort.Strings(pkgPaths)
	return pkgPaths, nil
}

// PackagesFromPaths creates a Packages using the provided relative paths. If any of the relative paths end in a splat
// ("/..."), then all of the sub-directories of that directory are also considered.
func PackagesFromPaths(rootDir string, relPaths []string) (Packages, error) {
	absoluteRoot, err := filepath.Abs(rootDir)
	if err != nil {
		return nil, fmt.Errorf("failed to convert %s to absolute path: %v", rootDir, err)
	}

	expandedRelPaths, err := expandPaths(rootDir, relPaths)
	if err != nil {
		return nil, fmt.Errorf("failed to expand paths %v: %v", relPaths, err)
	}

	pkgs := make(map[string]string, len(expandedRelPaths))
	for _, currPath := range expandedRelPaths {
		currAbsPath := path.Join(absoluteRoot, currPath)
		currPkg, err := getPrimaryPkgForDir(currAbsPath, nil)
		if err != nil {
			return nil, fmt.Errorf("unable to determine package for directory %s: %v", currAbsPath, err)
		}
		pkgs[currAbsPath] = currPkg
	}

	return createPkgsWithValidation(absoluteRoot, pkgs)
}

// PackagesInDir creates a Packages that contains all of the packages rooted at the provided directory. Every directory
// rooted in the provided directory whose path does not match the provided exclude matcher is considered as a package.
func PackagesInDir(rootDir string, exclude matcher.Matcher) (Packages, error) {
	dirAbsolutePath, err := filepath.Abs(rootDir)
	if err != nil {
		return nil, fmt.Errorf("failed to convert %s to absolute path: %v", rootDir, err)
	}

	allPkgs := make(map[string]string)
	if err := filepath.Walk(dirAbsolutePath, func(currPath string, currInfo os.FileInfo, err error) error {
		currRelPath, currRelPathErr := filepath.Rel(dirAbsolutePath, currPath)

		// skip current path if it matches an exclude
		if currRelPathErr == nil && exclude != nil && exclude.Match(currRelPath) {
			return nil
		}

		if err != nil {
			return err
		}

		if !currInfo.IsDir() {
			return nil
		}

		if currRelPathErr != nil {
			return currRelPathErr
		}

		// create a filter for processing package files that only passes if it does not match an exclude
		filter := func(info os.FileInfo) bool {
			// if exclude exists and matches the file, skip it
			if exclude != nil && exclude.Match(path.Join(currRelPath, info.Name())) {
				return false
			}
			// process file if it would be included in build context (handles things like build tags)
			match, _ := build.Default.MatchFile(currPath, info.Name())
			return match
		}

		pkgName, err := getPrimaryPkgForDir(currPath, filter)
		if err != nil {
			return fmt.Errorf("unable to determine package for directory %s: %v", currPath, err)
		}

		if pkgName != "" {
			allPkgs[currPath] = pkgName
		}

		return nil
	}); err != nil {
		return nil, err
	}

	return createPkgsWithValidation(dirAbsolutePath, allPkgs)
}

func createPkgsWithValidation(rootDir string, pkgs map[string]string) (*packages, error) {
	if !path.IsAbs(rootDir) {
		return nil, fmt.Errorf("rootDir %s is not an absolute path", rootDir)
	}

	for currAbsPkgPath := range pkgs {
		if !path.IsAbs(currAbsPkgPath) {
			return nil, fmt.Errorf("package %s in packages %v is not an absolute path", currAbsPkgPath, pkgs)
		}
	}

	return &packages{
		rootDir: rootDir,
		pkgs:    pkgs,
	}, nil
}

func expandPaths(rootDir string, relPaths []string) ([]string, error) {
	var expandedRelPaths []string
	for _, currRelPath := range relPaths {
		if strings.HasSuffix(currRelPath, "/...") {
			// expand splatted paths
			splatBaseDir := currRelPath[:len(currRelPath)-len("/...")]
			baseDirAbsPath := path.Join(rootDir, splatBaseDir)
			err := filepath.Walk(baseDirAbsPath, func(path string, info os.FileInfo, err error) error {
				if err != nil {
					return err
				}
				if info.IsDir() {
					relPath, err := filepath.Rel(rootDir, path)
					if err != nil {
						return err
					}
					expandedRelPaths = append(expandedRelPaths, relPath)
				}
				return nil
			})
			if err != nil {
				return nil, err
			}
		} else {
			expandedRelPaths = append(expandedRelPaths, currRelPath)
		}
	}
	return expandedRelPaths, nil
}

func getPrimaryPkgForDir(dir string, filter func(os.FileInfo) bool) (string, error) {
	pkgs, err := parser.ParseDir(token.NewFileSet(), dir, filter, parser.PackageClauseOnly)
	if err != nil {
		return "", fmt.Errorf("failed to parse directory %s as a package: %v", dir, err)
	}

	switch len(pkgs) {
	case 0:
		return "", nil
	case 1:
		// if only one entry exists, return its package
		for _, value := range pkgs {
			return value.Name, nil
		}
	default:
		// more than 1 entry exists: filter down to unique packages (if a package ends in "_test", remove suffix)
		uniquePkgs := make(map[string]struct{})
		for _, value := range pkgs {
			uniquePkgs[strings.TrimSuffix(value.Name, "_test")] = struct{}{}
		}

		// if there is only a single package, return it
		if len(uniquePkgs) == 1 {
			for pkg := range uniquePkgs {
				return pkg, nil
			}
		}

		// more than one package exists: return error
		pkgs := make([]string, 0, len(uniquePkgs))
		for pkg := range uniquePkgs {
			pkgs = append(pkgs, pkg)
		}
		sort.Strings(pkgs)
		return "", fmt.Errorf("directory %s contains more than 1 package: %v", dir, pkgs)
	}

	return "", nil
}
