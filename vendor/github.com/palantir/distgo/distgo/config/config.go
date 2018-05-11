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

package config

import (
	"fmt"
	"sort"
	"strings"

	"github.com/palantir/pkg/matcher"
	"github.com/pkg/errors"

	"github.com/palantir/distgo/distgo"
	"github.com/palantir/distgo/distgo/config/internal/v0"
)

type ProjectConfig v0.ProjectConfig

func ToProductsMap(in map[distgo.ProductID]ProductConfig) map[distgo.ProductID]v0.ProductConfig {
	if in == nil {
		return nil
	}
	out := make(map[distgo.ProductID]v0.ProductConfig, len(in))
	for k, v := range in {
		out[k] = v0.ProductConfig(v)
	}
	return out
}

func (cfg *ProjectConfig) ToParam(
	projectDir string,
	projectVersionerFactory distgo.ProjectVersionerFactory,
	disterFactory distgo.DisterFactory,
	defaultDisterCfg DisterConfig,
	dockerBuilderFactory distgo.DockerBuilderFactory,
	publisherFactory distgo.PublisherFactory) (distgo.ProjectParam, error) {

	var exclude matcher.Matcher
	if !cfg.Exclude.Empty() {
		exclude = cfg.Exclude.Matcher()
	}

	cfgProducts := cfg.Products
	if cfgProducts == nil && projectDir != "" {
		// if "products" is not specified at all, create default configuration for all main packages
		productCfgs, err := mainPkgsProductsConfig(projectDir, defaultDisterCfg, exclude)
		if err != nil {
			return distgo.ProjectParam{}, err
		}
		cfgProducts = toProductIDV0ProductConfigMap(productCfgs)
	}

	var products map[distgo.ProductID]distgo.ProductParam
	if cfgProducts != nil {
		products = make(map[distgo.ProductID]distgo.ProductParam)
	}

	var productIDs []distgo.ProductID
	for productID, productCfg := range cfgProducts {
		if strings.Contains(string(productID), ".") {
			return distgo.ProjectParam{}, errors.Errorf("ProductID cannot contain a '.': %s", productID)
		}

		productCfg := productCfg
		productParam, err := (*ProductConfig)(&productCfg).ToParam(productID, cfg.ScriptIncludes, ProductConfig(cfg.ProductDefaults), disterFactory, dockerBuilderFactory)
		if err != nil {
			return distgo.ProjectParam{}, err
		}
		products[productID] = productParam
		productIDs = append(productIDs, productID)
	}
	sort.Sort(distgo.ByProductID(productIDs))

	// compute full dependencies for all products (and error if cycles exist)
	cycleErrors := make(map[distgo.ProductID]error)
	for _, currProduct := range productIDs {
		allDeps, err := computeAllDependencies(currProduct, products, nil)
		if err != nil {
			cycleErrors[currProduct] = err
		}
		if len(allDeps) == 0 {
			continue
		}
		currProductParam := products[currProduct]
		currProductParam.AllDependencies = allDeps
		products[currProduct] = currProductParam
	}
	if len(cycleErrors) != 0 {
		// aggregate
		var sortedKeys []distgo.ProductID
		for k := range cycleErrors {
			sortedKeys = append(sortedKeys, k)
		}
		sort.Sort(distgo.ByProductID(sortedKeys))

		errOutputParts := []string{fmt.Sprintf("invalid dependencies for product(s) %v:", sortedKeys)}
		for _, currKey := range sortedKeys {
			errOutputParts = append(errOutputParts, fmt.Sprintf("%s: %v", currKey, cycleErrors[currKey]))
		}
		return distgo.ProjectParam{}, errors.Errorf("%s", strings.Join(errOutputParts, "\n  "))
	}

	// perform verification of ProductBuildID and ProductDistID dependencies in Docker outputs. Must be performed after
	// products are checked for cycles.
	for _, productID := range productIDs {
		productParam := products[productID]
		if productParam.Docker == nil {
			continue
		}
		var dockerIDs []distgo.DockerID
		for k := range productParam.Docker.DockerBuilderParams {
			dockerIDs = append(dockerIDs, k)
		}
		sort.Sort(distgo.ByDockerID(dockerIDs))
		productSubmap := newProductSubmap(productParam)
		for _, dockerID := range dockerIDs {
			dockerBuilderParam := productParam.Docker.DockerBuilderParams[dockerID]
			// verify that input builds for product are syntactically valid and specify legal products
			inputBuildProducts, err := distgo.ProductParamsForBuildProductArgs(productSubmap, dockerBuilderParam.InputBuilds...)
			if err != nil {
				return distgo.ProjectParam{}, errors.Errorf("invalid Docker input build(s) specified for DockerBuilderParam %q for product %q", dockerID, productID)
			}
			// input parameters are valid, but there may be product-level specifications. Expand all to "ProductID.OSArch" form.
			var expandedProductBuildIDs []distgo.ProductBuildID
			for _, productParam := range inputBuildProducts {
				if productParam.Build == nil {
					continue
				}
				for _, osArch := range productParam.Build.OSArchs {
					expandedProductBuildIDs = append(expandedProductBuildIDs, distgo.NewProductBuildID(productParam.ID, osArch))
				}
			}
			// assign updated slice to DockerBuilderParam and update in DockerBuilderParams map so that update is persistent
			dockerBuilderParam.InputBuilds = expandedProductBuildIDs
			productParam.Docker.DockerBuilderParams[dockerID] = dockerBuilderParam

			// verify that input dists for product are syntactically valid and specify legal products
			inputDistProducts, err := distgo.ProductParamsForDistProductArgs(productSubmap, dockerBuilderParam.InputDists...)
			if err != nil {
				return distgo.ProjectParam{}, errors.Errorf("invalid Docker input dist(s) specified for DockerBuilderParam %q for product %q", dockerID, productID)
			}
			// input parameters are valid, but there may be product-level specifications. Expand all to "ProductID.DistID" form.
			var expandedProductDistIDs []distgo.ProductDistID
			for _, productParam := range inputDistProducts {
				if productParam.Dist == nil {
					continue
				}
				var distIDs []distgo.DistID
				for distID := range productParam.Dist.DistParams {
					distIDs = append(distIDs, distID)
				}
				sort.Sort(distgo.ByDistID(distIDs))
				for _, distID := range distIDs {
					expandedProductDistIDs = append(expandedProductDistIDs, distgo.NewProductDistID(productParam.ID, distID))
				}
			}
			// assign updated slice to DockerBuilderParam and update in DockerBuilderParams map so that update is persistent
			dockerBuilderParam.InputDists = expandedProductDistIDs
			productParam.Docker.DockerBuilderParams[dockerID] = dockerBuilderParam
		}
	}

	projectVersionerCfg := (*ProjectVersionConfig)(cfg.ProjectVersioner)
	projectVersionerParam, err := projectVersionerCfg.ToParam(projectVersionerFactory)
	if err != nil {
		return distgo.ProjectParam{}, err
	}

	projectParam := distgo.ProjectParam{
		Products:              products,
		ScriptIncludes:        cfg.ScriptIncludes,
		ProjectVersionerParam: projectVersionerParam,
		Exclude:               exclude,
	}
	return projectParam, nil
}

// newProductSubmap returns a newly allocated map that contains only the provided product and all of its dependencies.
func newProductSubmap(productParam distgo.ProductParam) map[distgo.ProductID]distgo.ProductParam {
	out := make(map[distgo.ProductID]distgo.ProductParam)
	out[productParam.ID] = productParam
	for k, v := range productParam.AllDependencies {
		out[k] = v
	}
	return out
}

func computeAllDependencies(currProduct distgo.ProductID, allProducts map[distgo.ProductID]distgo.ProductParam, pathSoFar []distgo.ProductID) (map[distgo.ProductID]distgo.ProductParam, error) {
	for _, seen := range pathSoFar {
		if currProduct != seen {
			continue
		}
		pathSoFar = append(pathSoFar, seen)
		var pathStringParts []string
		for _, path := range pathSoFar {
			pathStringParts = append(pathStringParts, string(path))
		}
		return nil, errors.Errorf("cycle exists: %s", strings.Join(pathStringParts, " -> "))
	}

	currProductParam, ok := allProducts[currProduct]
	if !ok {
		return nil, errors.Errorf("%q is not a valid product", currProduct)
	}

	allDeps := make(map[distgo.ProductID]distgo.ProductParam)
	pathSoFar = append(pathSoFar, currProduct)
	for _, currDepProduct := range currProductParam.FirstLevelDependencies {
		allDeps[currDepProduct] = allProducts[currDepProduct]
		allCurDepProductDeps, err := computeAllDependencies(currDepProduct, allProducts, pathSoFar)
		if err != nil {
			return nil, err
		}
		for k, v := range allCurDepProductDeps {
			allDeps[k] = v
		}
	}
	return allDeps, nil
}

func toProductIDV0ProductConfigMap(in map[distgo.ProductID]ProductConfig) map[distgo.ProductID]v0.ProductConfig {
	if in == nil {
		return nil
	}
	out := make(map[distgo.ProductID]v0.ProductConfig, len(in))
	for k, v := range in {
		out[k] = v0.ProductConfig(v)
	}
	return out
}
