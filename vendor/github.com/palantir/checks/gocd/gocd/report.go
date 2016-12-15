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
	"sort"
	"strings"

	"github.com/pkg/errors"
)

type ImportReport struct {
	Imports         []ImportReportPkg `json:"imports"`
	MainOnlyImports []ImportReportPkg `json:"mainOnlyImports"`
	TestOnlyImports []ImportReportPkg `json:"testOnlyImports"`
}

type importReportPkgByPath []ImportReportPkg

func (p importReportPkgByPath) Len() int           { return len(p) }
func (p importReportPkgByPath) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
func (p importReportPkgByPath) Less(i, j int) bool { return p[i].Path < p[j].Path }

type ImportReportPkg struct {
	// path for the package
	Path string `json:"path"`
	// number of Go files in the package
	NGoFiles int `json:"numGoFiles"`
	// number of Go files imported by this package (does not include files in package itself)
	NImportedGoFiles int `json:"numImportedGoFiles"`
	// package path of the packages that import this package
	ImportSrc []string `json:"importedFrom"`
}

func CreateImportReport(rootDir string) (ImportReport, error) {
	project, err := NewProjectPkgInfoer(rootDir)
	if err != nil {
		return ImportReport{}, err
	}

	pkgs, err := importReportPkgs(project)
	if err != nil {
		return ImportReport{}, err
	}

	report := ImportReport{
		Imports:         make([]ImportReportPkg, 0),
		MainOnlyImports: make([]ImportReportPkg, 0),
		TestOnlyImports: make([]ImportReportPkg, 0),
	}

	for _, v := range pkgs {
		switch {
		case importedByTestOnly(&v):
			report.TestOnlyImports = append(report.TestOnlyImports, v)
		case importedByMainOnly(&v, project):
			report.MainOnlyImports = append(report.MainOnlyImports, v)
		default:
			report.Imports = append(report.Imports, v)
		}
	}

	sort.Sort(importReportPkgByPath(report.Imports))
	sort.Sort(importReportPkgByPath(report.MainOnlyImports))
	sort.Sort(importReportPkgByPath(report.TestOnlyImports))
	return report, nil
}

func importedByMainOnly(pkg *ImportReportPkg, project ProjectPkgInfoer) bool {
	for _, p := range pkg.ImportSrc {
		if pkgInfo, ok := project.PkgInfo(p); ok {
			if !strings.HasSuffix(pkgInfo.Path, "_test") && pkgInfo.Name != "main" {
				// if non-test, non-main package imports this package, it is not imported by main only
				return false
			}
		}
	}
	return true
}

func importedByTestOnly(pkg *ImportReportPkg) bool {
	for _, p := range pkg.ImportSrc {
		if !strings.HasSuffix(p, "_test") {
			// if any non-test package imports this package, it is not imported by test only
			return false
		}
	}
	return true
}

func importReportPkgs(project ProjectPkgInfoer) (map[string]ImportReportPkg, error) {
	counter, err := NewProjectGoFileCounter(project)
	if err != nil {
		return nil, err
	}
	impProvs := make(map[string]ImportReportPkg)
	for _, pkg := range project.PkgInfos() {
		for k := range pkg.Imports {
			// skip intra-project imports
			if !strings.Contains(k, "/vendor/") && strings.HasPrefix(k, project.RootDirImportPath()) {
				continue
			}

			if _, ok := impProvs[k]; !ok {
				// first time import has been seen -- add to map
				nGoFiles, ok := counter.NGoFiles(k)
				if !ok {
					return nil, errors.Errorf("could not determine number of Go files in %s", k)
				}
				nTotalGoFiles, ok := counter.NTotalGoFiles(k)
				if !ok {
					return nil, errors.Errorf("could not determine number of Go files in %s", k)
				}
				impProvs[k] = ImportReportPkg{
					Path:             k,
					NGoFiles:         nGoFiles,
					NImportedGoFiles: nTotalGoFiles - nGoFiles,
				}
			}

			// known to exist because of statement above
			impProv := impProvs[k]
			impProv.ImportSrc = append(impProv.ImportSrc, pkg.Path)
			impProvs[k] = impProv
		}
	}
	return impProvs, nil
}
