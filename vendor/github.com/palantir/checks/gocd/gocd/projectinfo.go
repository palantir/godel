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
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"

	"github.com/pkg/errors"
)

type PkgInfos []*PkgInfo

type ProjectPkgInfoer interface {
	RootDirImportPath() string
	PkgInfo(pkg string) (PkgInfo, bool)
	PkgInfos() PkgInfos
}

type projectPkgInfo struct {
	// import path to the "root" of the project
	rootDirImportPath string
	// stores packages that have been retrieved
	pkgs map[string]PkgInfo
}

func (p *projectPkgInfo) RootDirImportPath() string {
	return p.rootDirImportPath
}

func (p *projectPkgInfo) PkgInfo(pkg string) (PkgInfo, bool) {
	v, ok := p.pkgs[pkg]
	return v, ok
}

func (p *projectPkgInfo) PkgInfos() PkgInfos {
	var pi []*PkgInfo
	for _, v := range p.pkgs {
		v := v // intentional -- create separate variable that can be addressed
		pi = append(pi, &v)
	}
	sort.Sort(pkgInfoByPath(pi))
	return pi
}

type pkgInfoByPath []*PkgInfo

func (p pkgInfoByPath) Len() int           { return len(p) }
func (p pkgInfoByPath) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
func (p pkgInfoByPath) Less(i, j int) bool { return p[i].Path < p[j].Path }

func NewProjectPkgInfoer(rootDir string) (ProjectPkgInfoer, error) {
	rootDirImportPath, err := dirImportPath(rootDir)
	if err != nil {
		return nil, err
	}

	pkgs := make(map[string]PkgInfo)
	if err := filepath.Walk(rootDir, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			return nil
		}

		// skip any paths in a vendor directory
		if strings.Contains(path, "/vendor/") {
			return nil
		}

		fis, err := ioutil.ReadDir(path)
		if err != nil {
			return err
		}
		goFileExists := false
		for _, fi := range fis {
			if !fi.IsDir() && strings.HasSuffix(fi.Name(), ".go") {
				goFileExists = true
				break
			}
		}

		// skip directory if it does not contain at least one file ending in ".go"
		if !goFileExists {
			return nil
		}

		if pkg, empty, err := DirPkgInfo(path, Default); err != nil {
			return err
		} else if !empty {
			pkgs[pkg.Path] = pkg
		}

		if pkg, empty, err := DirPkgInfo(path, Test); err != nil {
			return err
		} else if !empty {
			pkgs[pkg.Path] = pkg
		}
		return nil
	}); err != nil {
		return nil, err
	}

	return &projectPkgInfo{
		rootDirImportPath: rootDirImportPath,
		pkgs:              pkgs,
	}, nil
}

func dirImportPath(dir string) (string, error) {
	// attempt to import
	if pkg, err := doImport(".", dir); err == nil {
		return pkg.ImportPath, nil
	}

	// import may fail if directory does not contain buildable Go files. In that case, determine import path
	// relative to GOPATH/src.
	if dirPath, err := filepath.EvalSymlinks(dir); err == nil {
		if importPath, err := filepath.Rel(path.Join(os.Getenv("GOPATH"), "src"), dirPath); err == nil {
			return importPath, nil
		}
	}

	return "", errors.Errorf("failed to determine import path for %s", dir)
}
