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

package v0

import (
	"github.com/palantir/distgo/distgo"
)

// ProductConfig represents user-specified configuration on how to build a specific product.
type ProductConfig struct {
	// Build specifies the build configuration for the product.
	Build *BuildConfig `yaml:"build"`

	// Run specifies the run configuration for the product.
	Run *RunConfig `yaml:"run"`

	// Dist specifies the dist configuration for the product.
	Dist *DistConfig `yaml:"dist"`

	// Publish specifies the publish configuration for the product.
	Publish *PublishConfig `yaml:"publish"`

	// Docker specifies the Docker configuration for the product.
	Docker *DockerConfig `yaml:"docker"`

	// Dependencies specifies the first-level dependencies of this product. Stores the IDs of the products.
	Dependencies *[]distgo.ProductID `yaml:"dependencies"`
}
