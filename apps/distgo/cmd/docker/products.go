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

package docker

import (
	"sort"

	"github.com/pkg/errors"

	"github.com/palantir/godel/apps/distgo/params"
)

func productsToDistAndBuildImage(
	products []string,
	cfg params.Project,
) (
	productsToDist []string,
	productsToBuildImage []string,
	err error,
) {
	if len(products) == 0 {
		for product := range cfg.FilteredProducts() {
			products = append(products, product)
		}
	}
	sort.Strings(products)
	if err := validateProducts(products, cfg); err != nil {
		return nil, nil, err
	}

	visited := make(map[string]struct{})
	distProducts := make(map[string]struct{})
	imageProducts := make(map[string]struct{})
	productQueue := make([]string, len(products))
	copy(productQueue, products)
	for len(productQueue) > 0 {
		productName := productQueue[0]
		productQueue = productQueue[1:]
		visited[productName] = struct{}{}
		productSpec := cfg.Products[productName]
		if len(productSpec.DockerImages) > 0 {
			imageProducts[productName] = struct{}{}
		}
		for _, image := range productSpec.DockerImages {
			for _, dep := range image.Deps {
				if isDist(dep.Type) {
					distProducts[productName] = struct{}{}
				}
				if dep.Type == params.DockerDepDocker {
					// has a docker image dependency. Add to product queue if not visited
					if _, ok := visited[dep.Product]; ok {
						continue
					}
					productQueue = append(productQueue, dep.Product)
				}
			}
		}
	}
	return setToSlice(distProducts), setToSlice(imageProducts), nil
}

func validateProducts(products []string, cfg params.Project) error {
	var unknownProducts []string
	for _, product := range products {
		if _, ok := cfg.Products[product]; !ok {
			unknownProducts = append(unknownProducts, product)
		}
	}
	sort.Strings(unknownProducts)
	if len(unknownProducts) != 0 {
		var filteredProducts []string
		for product := range cfg.FilteredProducts() {
			filteredProducts = append(filteredProducts, product)
		}
		sort.Strings(filteredProducts)
		return errors.Errorf(
			"Invalid products: %v. Valid products are: %v",
			unknownProducts,
			filteredProducts,
		)
	}
	return nil
}

func isDist(dep params.DockerDepType) bool {
	switch dep {
	case params.DockerDepSLS, params.DockerDepBin, params.DockerDepRPM:
		return true
	default:
		return false
	}
}

func setToSlice(s map[string]struct{}) []string {
	var result []string
	for item := range s {
		result = append(result, item)
	}
	sort.Strings(result)
	return result
}
