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
	"os"
	"path"
	"strings"
)

type ProjectGoFileCounter interface {
	NGoFiles(pkg string) (int, bool)
	NTotalGoFiles(pkg string) (int, bool)
}

type projectGoFileCounter struct {
	ProjectPkgInfoer
	counts map[string]goFileCount
}

type goFileCount struct {
	pkg   int
	total int
}

func NewProjectGoFileCounter(p ProjectPkgInfoer) (ProjectGoFileCounter, error) {
	counter := projectGoFileCounter{
		ProjectPkgInfoer: p,
		counts:           make(map[string]goFileCount),
	}

	// cache from pkg -> all packages imported by the package (recursive)
	importsCache := make(map[string]map[string]*PkgInfo)
	for _, v := range p.PkgInfos() {
		// determine file count by determining all of the unique packages imported by a package and then summing
		// up the package file count of each. This approach is required to avoid double-counting packages that
		// are imported multiple times.
		if _, err := counter.allImports(v, importsCache, counter.counts); err != nil {
			return nil, err
		}
	}

	return &counter, nil
}

func (p *projectGoFileCounter) NGoFiles(pkg string) (int, bool) {
	if c, ok := p.counts[pkg]; ok {
		return c.pkg, ok
	}
	return 0, false
}

func (p *projectGoFileCounter) NTotalGoFiles(pkg string) (int, bool) {
	if c, ok := p.counts[pkg]; ok {
		return c.total, ok
	}
	return 0, false
}

func (p *projectGoFileCounter) allImports(pkg *PkgInfo, cache map[string]map[string]*PkgInfo, countsMap map[string]goFileCount) (map[string]*PkgInfo, error) {
	if v, ok := cache[pkg.Path]; ok {
		return v, nil
	}

	pkgImports := make(map[string]*PkgInfo)
	for k := range pkg.Imports {
		var importPkg *PkgInfo
		if v, ok := p.PkgInfo(k); ok {
			importPkg = &v
		} else {
			if newImportPkg, empty, err := ImportPkgInfo(k, path.Join(os.Getenv("GOPATH"), "src", strings.TrimSuffix(pkg.Path, "_test")), Default); err != nil {
				return nil, err
			} else if !empty {
				importPkg = &newImportPkg
			}
		}
		pkgImports[importPkg.Path] = importPkg
		result, err := p.allImports(importPkg, cache, countsMap)
		if err != nil {
			return nil, err
		}
		for k, v := range result {
			pkgImports[k] = v
		}
	}
	cache[pkg.Path] = pkgImports

	// compute and populate counts
	counts := goFileCount{
		pkg:   pkg.NGoFiles,
		total: pkg.NGoFiles,
	}
	for _, v := range pkgImports {
		counts.total += v.NGoFiles
	}
	countsMap[pkg.Path] = counts

	return pkgImports, nil
}
