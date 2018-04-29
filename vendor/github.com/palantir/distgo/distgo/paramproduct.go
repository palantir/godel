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
	"sort"
)

type ProductParam struct {
	// ID is the unique identifier for this product. Its value comes from the key for the product in the Products map in
	// the configuration.
	ID ProductID

	// Build specifies the build configuration for the product.
	Build *BuildParam

	// Run specifies the run configuration for the product.
	Run *RunParam

	// Dist specifies the dist configuration for the product.
	Dist *DistParam

	// Publish specifies the publish configuration for the product.
	Publish *PublishParam

	// Docker specifies the Docker configuration for the product.
	Docker *DockerParam

	// FirstLevelDependencies stores the IDs of the products that are declared as dependencies of this product.
	FirstLevelDependencies []ProductID

	// AllDependencies stores all of the dependent products of this product. It is a result of expanding all of the
	// dependencies in FirstLevelDependencies.
	AllDependencies map[ProductID]ProductParam
}

func (p *ProductParam) AllProductParams() []ProductParam {
	allProductParams := []ProductParam{*p}
	for _, currParam := range p.AllDependencies {
		allProductParams = append(allProductParams, currParam)
	}
	sort.Slice(allProductParams, func(i, j int) bool {
		return allProductParams[i].ID < allProductParams[j].ID
	})
	return allProductParams
}

func (p *ProductParam) AllDependenciesSortedIDs() []ProductID {
	var sortedKeys []ProductID
	for k := range p.AllDependencies {
		sortedKeys = append(sortedKeys, k)
	}
	sort.Sort(ByProductID(sortedKeys))
	return sortedKeys
}

type ProductOutputInfo struct {
	ID                ProductID          `json:"productId"`
	BuildOutputInfo   *BuildOutputInfo   `json:"buildOutputInfo"`
	DistOutputInfos   *DistOutputInfos   `json:"distOutputInfos"`
	PublishOutputInfo *PublishOutputInfo `json:"publishOutputInfo"`
	DockerOutputInfos *DockerOutputInfos `json:"dockerOutputInfos"`
}

func (p *ProductParam) ToProductOutputInfo(version string) (ProductOutputInfo, error) {
	var buildOutputInfo *BuildOutputInfo
	if p.Build != nil {
		buildOutputInfoVar, err := p.Build.ToBuildOutputInfo(p.ID, version)
		if err != nil {
			return ProductOutputInfo{}, err
		}
		buildOutputInfo = &buildOutputInfoVar
	}
	var distOutputInfos *DistOutputInfos
	if p.Dist != nil {
		distOutputInfosVar, err := p.Dist.ToDistOutputInfos(p.ID, version)
		if err != nil {
			return ProductOutputInfo{}, err
		}
		distOutputInfos = &distOutputInfosVar
	}
	var publishOutputInfo *PublishOutputInfo
	if p.Publish != nil {
		publishOutputInfoVar := p.Publish.ToPublishOutputInfo()
		publishOutputInfo = &publishOutputInfoVar
	}
	var dockerOutputInfos *DockerOutputInfos
	if p.Docker != nil {
		dockerOutputInfosVar, err := p.Docker.ToDockerOutputInfos(p.ID, version)
		if err != nil {
			return ProductOutputInfo{}, err
		}
		dockerOutputInfos = &dockerOutputInfosVar
	}
	return ProductOutputInfo{
		ID:                p.ID,
		BuildOutputInfo:   buildOutputInfo,
		DistOutputInfos:   distOutputInfos,
		PublishOutputInfo: publishOutputInfo,
		DockerOutputInfos: dockerOutputInfos,
	}, nil
}
