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

package params

import (
	"sort"
	"strings"

	"github.com/pkg/errors"

	"github.com/palantir/godel/apps/distgo/pkg/git"
	"github.com/palantir/godel/apps/distgo/pkg/osarch"
)

// ProductBuildSpec defines all of the parameters for building a specific product.
type ProductBuildSpec struct {
	Product
	ProjectDir     string
	ProductName    string
	ProductVersion string
	VersionInfo    git.ProjectInfo
}

type ProductBuildSpecWithDeps struct {
	Spec      ProductBuildSpec
	BuildDeps map[string]ProductBuildSpec
	DistDeps  map[string]ProductBuildSpec
}

func (p *ProductBuildSpecWithDeps) AllSpecs() []ProductBuildSpec {
	allSpecs := make([]ProductBuildSpec, 0, len(p.BuildDeps)+1)
	allSpecs = append(allSpecs, p.Spec)
	for _, spec := range p.BuildDeps {
		allSpecs = append(allSpecs, spec)
	}
	return allSpecs
}

const (
	defaultBuildOutputDir = "build"
	defaultDistOutputDir  = "dist"
)

func NewProductBuildSpecWithDeps(spec ProductBuildSpec, allSpecs map[string]ProductBuildSpec) (ProductBuildSpecWithDeps, error) {
	buildDeps := make(map[string]ProductBuildSpec)
	distDeps := make(map[string]ProductBuildSpec)
	for _, currDistCfg := range spec.Dist {
		for _, currDepProduct := range currDistCfg.InputProducts {
			currSpec, err := getSpec(currDepProduct, spec, allSpecs)
			if err != nil {
				return ProductBuildSpecWithDeps{}, err
			}
			buildDeps[currDepProduct] = currSpec
		}
		for _, currDepProduct := range currDistCfg.Info.Deps() {
			currSpec, err := getSpec(currDepProduct, spec, allSpecs)
			if err != nil {
				return ProductBuildSpecWithDeps{}, err
			}
			distDeps[currDepProduct] = currSpec
		}
	}

	return ProductBuildSpecWithDeps{
		Spec:      spec,
		BuildDeps: buildDeps,
		DistDeps:  distDeps,
	}, nil
}

// NewProductBuildSpec returns a fully initialized ProductBuildSpec that is a combination of the provided parameters.
// If any of the required fields in the provided configuration is blank, the returned ProjectBuildSpec will have default
// values populated in the returned object.
func NewProductBuildSpec(projectDir, productName string, gitProductInfo git.ProjectInfo, productCfg Product, projectCfg Project) ProductBuildSpec {
	buildSpec := ProductBuildSpec{
		Product:        productCfg,
		ProjectDir:     projectDir,
		ProductName:    productName,
		ProductVersion: gitProductInfo.Version,
		VersionInfo:    gitProductInfo,
	}

	if buildSpec.Build.OutputDir == "" {
		buildSpec.Build.OutputDir = firstNonEmpty(projectCfg.BuildOutputDir, defaultBuildOutputDir)
	}

	if len(buildSpec.Build.OSArchs) == 0 {
		buildSpec.Build.OSArchs = []osarch.OSArch{osarch.Current()}
	}

	if len(buildSpec.Dist) == 0 {
		// One dist with all default values.
		buildSpec.Dist = []Dist{{}}
	}
	for i := range buildSpec.Dist {
		currDistCfg := &buildSpec.Dist[i]

		if currDistCfg.OutputDir == "" {
			currDistCfg.OutputDir = firstNonEmpty(projectCfg.DistOutputDir, defaultDistOutputDir)
		}

		if currDistCfg.Info == nil || currDistCfg.Info.Type() == "" {
			currDistCfg.Info = &SLSDistInfo{}
		}

		if currDistCfg.Publish.empty() {
			currDistCfg.Publish = buildSpec.DefaultPublish
		}

		if currDistCfg.Publish.GroupID == "" {
			currDistCfg.Publish.GroupID = projectCfg.GroupID
		}

		// if distribution is SLSv2, ensure that SLSv2 tag exists for Almanac
		if currDistCfg.Info.Type() == SLSDistType {
			slsv2TagExists := false
			for _, currTag := range currDistCfg.Publish.Almanac.Tags {
				if currTag == "slsv2" {
					slsv2TagExists = true
					break
				}
			}
			if !slsv2TagExists {
				currDistCfg.Publish.Almanac.Tags = append(currDistCfg.Publish.Almanac.Tags, "slsv2")
			}
		}

		if projectCfg.DistScriptInclude != "" && currDistCfg.Script != "" {
			currDistCfg.Script = strings.Join([]string{projectCfg.DistScriptInclude, currDistCfg.Script}, "\n")
		}
	}

	return buildSpec
}

func firstNonEmpty(first, second string) string {
	if first != "" {
		return first
	}
	return second
}

func getSpec(product string, spec ProductBuildSpec, allSpecs map[string]ProductBuildSpec) (ProductBuildSpec, error) {
	currSpec, ok := allSpecs[product]
	if !ok {
		allProducts := make([]string, 0, len(allSpecs))
		for currName := range allSpecs {
			allProducts = append(allProducts, currName)
		}
		sort.Strings(allProducts)
		return ProductBuildSpec{}, errors.Errorf("Spec %v declared %v as a dependent product, but could not find configuration for that product in %v", spec.ProductName, product, allProducts)
	}
	return currSpec, nil
}
