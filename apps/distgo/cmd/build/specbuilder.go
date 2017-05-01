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

package build

import (
	"fmt"
	"io"
	"os"
	"path"
	"sort"

	"github.com/palantir/pkg/matcher"
	"github.com/palantir/pkg/pkgpath"
	"github.com/pkg/errors"

	"github.com/palantir/godel/apps/distgo/cmd"
	"github.com/palantir/godel/apps/distgo/params"
	"github.com/palantir/godel/apps/distgo/pkg/git"
	"github.com/palantir/godel/apps/distgo/pkg/imports"
	"github.com/palantir/godel/apps/distgo/pkg/osarch"
)

// RequiresBuild returns a slice that contains the ProductBuildSpecs that have not been built for the provided
// ProductBuildSpecWithDeps matching the provided osArchs filter. A product is considered to require building if its
// output executable does not exist or if the output executable's modification date is older than any of the Go files
// required to build the product.
func RequiresBuild(specWithDeps params.ProductBuildSpecWithDeps, osArchs cmd.OSArchFilter) RequiresBuildInfo {
	info := newRequiresBuildInfo(specWithDeps, osArchs)
	for _, currSpec := range specWithDeps.AllSpecs() {
		if currSpec.Build.Skip {
			continue
		}
		paths := ArtifactPaths(currSpec)
		for _, currOSArch := range currSpec.Build.OSArchs {
			if osArchs.Matches(currOSArch) {
				if fi, err := os.Stat(paths[currOSArch]); err == nil {
					if goFiles, err := imports.AllFiles(path.Join(currSpec.ProjectDir, currSpec.Build.MainPkg)); err == nil {
						if newerThan, err := goFiles.NewerThan(fi); err == nil && !newerThan {
							// if the build artifact for the product already exists and none of the source files for the
							// product are newer than the build artifact, consider spec up-to-date
							continue
						}
					}
				}
				// spec/osArch combination requires build
				info.addInfo(currSpec, currOSArch)
			}
		}
	}
	return info
}

type RequiresBuildInfo interface {
	Specs() []params.ProductBuildSpec
	RequiresBuild(product string, osArch osarch.OSArch) bool
	addInfo(spec params.ProductBuildSpec, osArch osarch.OSArch)
}

type requiresBuildInfo struct {
	// ordered slice of product names
	orderedProducts []string
	// map from product name to build spec for the product
	products map[string]params.ProductBuildSpec
	// map from product name to OS/Archs for which product requires build
	productsRequiresBuildOSArch map[string][]osarch.OSArch
	// the products that were examined in creating this requiresBuildInfo
	examinedProducts map[string]struct{}
	// the OSArchFilter used when creating this requiresBuildInfo
	examinedOSArchs cmd.OSArchFilter
}

func newRequiresBuildInfo(specWithDeps params.ProductBuildSpecWithDeps, osArchs cmd.OSArchFilter) RequiresBuildInfo {
	examinedProducts := make(map[string]struct{})
	for _, spec := range specWithDeps.AllSpecs() {
		examinedProducts[spec.ProductName] = struct{}{}
	}

	return &requiresBuildInfo{
		products:                    make(map[string]params.ProductBuildSpec),
		productsRequiresBuildOSArch: make(map[string][]osarch.OSArch),
		examinedProducts:            examinedProducts,
		examinedOSArchs:             osArchs,
	}
}

func (b *requiresBuildInfo) addInfo(spec params.ProductBuildSpec, osArch osarch.OSArch) {
	k := spec.ProductName
	_, productSeen := b.products[k]
	if !productSeen {
		b.orderedProducts = append(b.orderedProducts, k)
	}
	b.products[k] = spec
	b.productsRequiresBuildOSArch[k] = append(b.productsRequiresBuildOSArch[k], osArch)
}

func (b *requiresBuildInfo) RequiresBuild(product string, osArch osarch.OSArch) bool {
	// if required product/OSArch was not considered, return true (assume it needs to be built)
	if _, ok := b.examinedProducts[product]; !ok || !b.examinedOSArchs.Matches(osArch) {
		return true
	}
	for _, v := range b.productsRequiresBuildOSArch[product] {
		if v == osArch {
			return true
		}
	}
	return false
}

func (b *requiresBuildInfo) Specs() []params.ProductBuildSpec {
	specs := make([]params.ProductBuildSpec, len(b.orderedProducts))
	for i, k := range b.orderedProducts {
		specs[i] = b.products[k]
	}
	return specs
}

func RunBuildFunc(buildActionFunc cmd.BuildFunc, cfg params.Project, products []string, wd string, stdout io.Writer) error {
	buildSpecsWithDeps, err := SpecsWithDepsForArgs(cfg, products, wd)
	if err != nil {
		return err
	}
	if err := buildActionFunc(buildSpecsWithDeps, stdout); err != nil {
		return err
	}
	return nil
}

func SpecsWithDepsForArgs(cfg params.Project, products []string, wd string) ([]params.ProductBuildSpecWithDeps, error) {
	// if configuration is empty, default to all main pkgs
	if len(cfg.Products) == 0 {
		cfg.Products = make(map[string]params.Product)
		if err := addMainPkgsToConfig(cfg, wd); err != nil {
			return nil, errors.Wrapf(err, "failed to get main packages from %v", wd)
		}
	}

	// determine version for git directory
	productInfo, err := git.NewProjectInfo(wd)
	if err != nil {
		// if version could not be determined, use "unspecified"
		productInfo.Version = "unspecified"
	}

	// create BuildSpec for all products
	allBuildSpecs := make(map[string]params.ProductBuildSpec)
	for currProduct, currProductCfg := range cfg.Products {
		allBuildSpecs[currProduct] = params.NewProductBuildSpec(wd, currProduct, productInfo, currProductCfg, cfg)
	}

	// get products that aren't excluded by configuration
	filteredProducts := cfg.FilteredProducts()

	if len(filteredProducts) == 0 {
		return nil, fmt.Errorf("No products found.")
	}

	// if arguments are provided, filter to only build products named in arguments
	if len(products) != 0 {
		var unknownProducts []string
		// create map of provided products
		argProducts := make(map[string]bool, len(products))
		for _, currArg := range products {
			argProducts[currArg] = true
			if _, ok := filteredProducts[currArg]; !ok {
				unknownProducts = append(unknownProducts, currArg)
			}
		}

		// throw error if any of the specified products were unknown
		if len(unknownProducts) > 0 {
			sort.Strings(unknownProducts)
			sortedKnownProducts := make([]string, 0, len(filteredProducts))
			for currProduct := range filteredProducts {
				sortedKnownProducts = append(sortedKnownProducts, currProduct)
			}
			sort.Strings(sortedKnownProducts)
			return nil, fmt.Errorf("Invalid products: %v. Valid products are %v.", unknownProducts, sortedKnownProducts)
		}

		// iterate over filteredProducts map and remove any keys not present in provided arguments
		for k := range filteredProducts {
			if _, ok := argProducts[k]; !ok {
				delete(filteredProducts, k)
			}
		}
	}

	sortedFilteredProducts := make([]string, 0, len(filteredProducts))
	for currProduct := range filteredProducts {
		sortedFilteredProducts = append(sortedFilteredProducts, currProduct)
	}
	sort.Strings(sortedFilteredProducts)

	var buildSpecsWithDeps []params.ProductBuildSpecWithDeps
	for _, currProduct := range sortedFilteredProducts {
		currSpec, err := params.NewProductBuildSpecWithDeps(allBuildSpecs[currProduct], allBuildSpecs)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to create build spec for %v", currProduct)
		}

		buildSpecsWithDeps = append(buildSpecsWithDeps, currSpec)
	}
	return buildSpecsWithDeps, nil
}

func addMainPkgsToConfig(cfg params.Project, projectDir string) error {
	mainPkgPaths, err := mainPkgPaths(projectDir)
	if err != nil {
		return errors.Wrapf(err, "failed to determine paths to main packages in %v", projectDir)
	}

	for _, currMainPkgPath := range mainPkgPaths {
		currMainPkgAbsPath := path.Join(projectDir, currMainPkgPath)
		productName := path.Base(currMainPkgAbsPath)
		cfg.Products[productName] = params.Product{
			Build: params.Build{
				MainPkg: currMainPkgPath,
			},
		}
	}
	return nil
}

func mainPkgPaths(projectDir string) ([]string, error) {
	// TODO: this should use Exclude specified in config to determine directories to examine
	pkgs, err := pkgpath.PackagesInDir(projectDir, matcher.Name("vendor"))
	if err != nil {
		return nil, errors.Wrapf(err, "failed to list packages in project %v", projectDir)
	}

	pkgsMap, err := pkgs.Packages(pkgpath.Relative)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get paths for packges")
	}

	var mainPkgPaths []string
	for currPath, currPkg := range pkgsMap {
		if currPkg == "main" {
			mainPkgPaths = append(mainPkgPaths, currPath)
		}
	}
	return mainPkgPaths, nil
}
