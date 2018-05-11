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
	"gopkg.in/yaml.v2"

	"github.com/palantir/distgo/distgo"
)

type DockerConfig struct {
	// Repository is the repository that is made available to the tag and Dockerfile templates.
	Repository *string `yaml:"repository,omitempty"`

	// DockerBuilderParams contains the Docker params for this distribution.
	DockerBuildersConfig *DockerBuildersConfig `yaml:"docker-builders,omitempty"`
}

type DockerBuildersConfig map[distgo.DockerID]DockerBuilderConfig

type DockerBuilderConfig struct {
	// Type is the type of the DockerBuilder. This field must be non-nil and non-empty and resolve to a valid DockerBuilder.
	Type *string `yaml:"type,omitempty"`
	// Config is the YAML configuration content for the DockerBuilder.
	Config *yaml.MapSlice `yaml:"config,omitempty"`
	// DockerfilePath is the path to the Dockerfile that is used to build the Docker image. The path is interpreted
	// relative to ContextDir. The content of the Dockerfile supports using Go templates. The following template
	// parameters can be used in the template:
	//   * {{Product}}: the name of the product
	//   * {{Version}}: the version of the project
	//   * {{ProductTaskOutputInfo}}: the ProductTaskOutputInfo struct
	DockerfilePath *string `yaml:"dockerfile-path,omitempty"`
	ContextDir     *string `yaml:"context-dir,omitempty"`
	// InputProductsDir is the directory in the context dir in which input products are written.
	InputProductsDir *string `yaml:"input-products-dir,omitempty"`
	// InputBuilds specifies the products whose build outputs should be made available to the Docker build task. The
	// specified products will be hard-linked into the context directory. The referenced products must be this product
	// or one of its declared dependencies.
	InputBuilds *[]distgo.ProductBuildID `yaml:"input-builds,omitempty"`
	// InputDists specifies the products whose dist outputs should be made available to the Docker build task. The
	// specified dists will be hard-linked into the context directory. The referenced products must be this product
	// or one of its declared dependencies.
	InputDists *[]distgo.ProductDistID `yaml:"input-dists,omitempty"`
	// TagTemplates specifies the templates that should be used to render the tag(s) for the Docker image. If multiple
	// values are specified, the image will be tagged with all of them.
	TagTemplates *[]string `yaml:"tag-templates,omitempty"`
}
