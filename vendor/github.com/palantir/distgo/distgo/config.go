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

package distgo

import (
	"fmt"
	"io/ioutil"
	"sort"
	"strings"

	"github.com/palantir/pkg/matcher"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

type ProductID string

func ToProductIDs(in []string) []ProductID {
	var ids []ProductID
	for _, id := range in {
		ids = append(ids, ProductID(id))
	}
	return ids
}

func ProductIDsToStrings(in []ProductID) []string {
	var ids []string
	for _, id := range in {
		ids = append(ids, string(id))
	}
	return ids
}

type ByProductID []ProductID

func (a ByProductID) Len() int           { return len(a) }
func (a ByProductID) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByProductID) Less(i, j int) bool { return a[i] < a[j] }

type ProjectParam struct {
	// Products contains the parameters for the defined products.
	Products map[ProductID]ProductParam

	// ScriptIncludes specifies a string that is appended to every script that is written out. Can be used to define
	// functions or constants for all scripts.
	ScriptIncludes string

	// VersionFunc is a shell script that specifies how the "version" string is generated for a project. If specified,
	// the content of the string is written to a temporary file and executed with the project directory as the working
	// directory. The first line of the output is used as the version for the project (and if the output is blank,
	// "unspecified" is used). If VersionFunc is blank, the output of the git.ProjectVersion function is used.
	VersionScript string

	// Exclude is a matcher that matches any directories that should be ignored as main files. Only relevant if products
	// are not specified.
	Exclude matcher.Matcher
}

func (p *ProjectParam) ProjectInfo(projectDir string) (ProjectInfo, error) {
	version, err := ProjectVersion(projectDir, p.VersionScript)
	if err != nil {
		return ProjectInfo{}, err
	}
	return ProjectInfo{
		ProjectDir: projectDir,
		Version:    version,
	}, nil
}

type ProjectConfig struct {
	// Products maps product names to configurations.
	Products map[ProductID]ProductConfig `yaml:"products"`

	// ProductDefaults specifies the default values that should be used for unspecified values in the products map. If a
	// field in a top-level key in a "ProductConfig" value in the "Products" map is nil and the corresponding value in
	// ProductDefaults is non-nil, the value in ProductDefaults is used.
	ProductDefaults ProductConfig `yaml:"product-defaults"`

	// ScriptIncludes specifies a string that is appended to every script that is written out. Can be used to define
	// functions or constants for all scripts.
	ScriptIncludes string `yaml:"script-includes"`

	// VersionFunc is a shell script that specifies how the "version" string is generated for a project. If specified,
	// the content of the string is written to a temporary file and executed with the project directory as the working
	// directory. The first line of the output is used as the version for the project (and if the output is blank,
	// "unspecified" is used). If VersionFunc is blank, the output of the git.ProjectVersion function is used.
	VersionScript string `yaml:"version-script"`

	// Exclude matches the paths to exclude when determining the projects to build.
	Exclude matcher.NamesPathsCfg `yaml:"exclude"`
}

func (cfg *ProjectConfig) ToParam(projectDir string, disterFactory DisterFactory, defaultDisterCfg DisterConfig, dockerBuilderFactory DockerBuilderFactory) (ProjectParam, error) {
	var exclude matcher.Matcher
	if !cfg.Exclude.Empty() {
		exclude = cfg.Exclude.Matcher()
	}

	cfgProducts := cfg.Products
	if cfgProducts == nil && projectDir != "" {
		// if "products" is not specified at all, create default configuration for all main packages
		productCfgs, err := mainPkgsProductsConfig(projectDir, defaultDisterCfg, exclude)
		if err != nil {
			return ProjectParam{}, err
		}
		cfgProducts = productCfgs
	}

	var products map[ProductID]ProductParam
	if cfgProducts != nil {
		products = make(map[ProductID]ProductParam)
	}

	var productIDs []ProductID
	for productID, productCfg := range cfgProducts {
		if strings.Contains(string(productID), ".") {
			return ProjectParam{}, errors.Errorf("ProductID cannot contain a '.': %s", productID)
		}

		productParam, err := productCfg.ToParam(productID, cfg.ScriptIncludes, cfg.ProductDefaults, disterFactory, dockerBuilderFactory)
		if err != nil {
			return ProjectParam{}, err
		}
		products[productID] = productParam
		productIDs = append(productIDs, productID)
	}
	sort.Sort(ByProductID(productIDs))

	// compute full dependencies for all products (and error if cycles exist)
	cycleErrors := make(map[ProductID]error)
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
		var sortedKeys []ProductID
		for k := range cycleErrors {
			sortedKeys = append(sortedKeys, k)
		}
		sort.Sort(ByProductID(sortedKeys))

		errOutputParts := []string{fmt.Sprintf("invalid dependencies for product(s) %v:", sortedKeys)}
		for _, currKey := range sortedKeys {
			errOutputParts = append(errOutputParts, fmt.Sprintf("%s: %v", currKey, cycleErrors[currKey]))
		}
		return ProjectParam{}, errors.Errorf("%s", strings.Join(errOutputParts, "\n  "))
	}

	// perform verification of ProductBuildID and ProductDistID dependencies in Docker outputs. Must be performed after
	// products are checked for cycles.
	for _, productID := range productIDs {
		productParam := products[productID]
		if productParam.Docker == nil {
			continue
		}
		var dockerIDs []DockerID
		for k := range productParam.Docker.DockerBuilderParams {
			dockerIDs = append(dockerIDs, k)
		}
		sort.Sort(ByDockerID(dockerIDs))
		productSubmap := newProductSubmap(productParam)
		for _, dockerID := range dockerIDs {
			dockerBuilderParam := productParam.Docker.DockerBuilderParams[dockerID]
			// verify that input builds for product are syntactically valid and specify legal products
			inputBuildProducts, err := ProductParamsForBuildProductArgs(productSubmap, dockerBuilderParam.InputBuilds...)
			if err != nil {
				return ProjectParam{}, errors.Errorf("invalid Docker input build(s) specified for DockerBuilderParam %q for product %q", dockerID, productID)
			}
			// input parameters are valid, but there may be product-level specifications. Expand all to "ProductID.OSArch" form.
			var expandedProductBuildIDs []ProductBuildID
			for _, productParam := range inputBuildProducts {
				if productParam.Build == nil {
					continue
				}
				for _, osArch := range productParam.Build.OSArchs {
					expandedProductBuildIDs = append(expandedProductBuildIDs, NewProductBuildID(productParam.ID, osArch))
				}
			}
			// assign updated slice to DockerBuilderParam and update in DockerBuilderParams map so that update is persistent
			dockerBuilderParam.InputBuilds = expandedProductBuildIDs
			productParam.Docker.DockerBuilderParams[dockerID] = dockerBuilderParam

			// verify that input dists for product are syntactically valid and specify legal products
			inputDistProducts, err := ProductParamsForDistProductArgs(productSubmap, dockerBuilderParam.InputDists...)
			if err != nil {
				return ProjectParam{}, errors.Errorf("invalid Docker input dist(s) specified for DockerBuilderParam %q for product %q", dockerID, productID)
			}
			// input parameters are valid, but there may be product-level specifications. Expand all to "ProductID.DistID" form.
			var expandedProductDistIDs []ProductDistID
			for _, productParam := range inputDistProducts {
				if productParam.Dist == nil {
					continue
				}
				var distIDs []DistID
				for distID := range productParam.Dist.DistParams {
					distIDs = append(distIDs, distID)
				}
				sort.Sort(ByDistID(distIDs))
				for _, distID := range distIDs {
					expandedProductDistIDs = append(expandedProductDistIDs, NewProductDistID(productParam.ID, distID))
				}
			}
			// assign updated slice to DockerBuilderParam and update in DockerBuilderParams map so that update is persistent
			dockerBuilderParam.InputDists = expandedProductDistIDs
			productParam.Docker.DockerBuilderParams[dockerID] = dockerBuilderParam
		}
	}

	projectParam := ProjectParam{
		Products:       products,
		ScriptIncludes: cfg.ScriptIncludes,
		VersionScript:  createScriptContent(cfg.VersionScript, cfg.ScriptIncludes),
		Exclude:        exclude,
	}
	return projectParam, nil
}

// newProductSubmap returns a newly allocated map that contains only the provided product and all of its dependencies.
func newProductSubmap(productParam ProductParam) map[ProductID]ProductParam {
	out := make(map[ProductID]ProductParam)
	out[productParam.ID] = productParam
	for k, v := range productParam.AllDependencies {
		out[k] = v
	}
	return out
}

func computeAllDependencies(currProduct ProductID, allProducts map[ProductID]ProductParam, pathSoFar []ProductID) (map[ProductID]ProductParam, error) {
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

	allDeps := make(map[ProductID]ProductParam)
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

func LoadConfigFromFile(cfgFile string) (ProjectConfig, error) {
	cfgBytes, err := ioutil.ReadFile(cfgFile)
	if err != nil {
		return ProjectConfig{}, errors.Wrapf(err, "failed to read configuration file")
	}
	return LoadConfig(cfgBytes)
}

func LoadConfig(cfgBytes []byte) (ProjectConfig, error) {
	var cfg ProjectConfig
	if err := yaml.Unmarshal(cfgBytes, &cfg); err != nil {
		return ProjectConfig{}, errors.Wrapf(err, "failed to unmarshal configuration")
	}
	return cfg, nil
}
