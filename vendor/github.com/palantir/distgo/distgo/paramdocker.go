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

	"github.com/palantir/godel/pkg/osarch"
	"github.com/pkg/errors"
)

type DockerID string

type ByDockerID []DockerID

func (a ByDockerID) Len() int           { return len(a) }
func (a ByDockerID) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByDockerID) Less(i, j int) bool { return a[i] < a[j] }

type DockerParam struct {
	// Repository is the Docker repository. This value is made available to TagTemplates as {{Repository}}.
	Repository string

	// DockerBuilderParams contains the Docker params for this distribution.
	DockerBuilderParams map[DockerID]DockerBuilderParam
}

type DockerOutputInfos struct {
	DockerIDs                []DockerID                           `json:"dockerIds"`
	Repository               string                               `json:"repository"`
	DockerBuilderOutputInfos map[DockerID]DockerBuilderOutputInfo `json:"dockerBuilderOutputInfos"`
}

type OSArchID string

type ByOSArchID []OSArchID

func (a ByOSArchID) Len() int           { return len(a) }
func (a ByOSArchID) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByOSArchID) Less(i, j int) bool { return a[i] < a[j] }

type DockerBuilderOutputInfo struct {
	ContextDir       string                              `json:"contextDir"`
	DockerfilePath   string                              `json:"dockerfilePath"`
	InputProductsDir string                              `json:"inputProductsDir"`
	RenderedTags     []string                            `json:"renderedDockerTags"`
	InputBuilds      map[ProductID]map[OSArchID]struct{} `json:"inputBuilds"`
	InputDists       map[ProductID]map[DistID]struct{}   `json:"inputDists"`
}

func (doi *DockerBuilderOutputInfo) InputBuildProductIDs() []ProductID {
	var productIDs []ProductID
	for k := range doi.InputBuilds {
		productIDs = append(productIDs, k)
	}
	sort.Sort(ByProductID(productIDs))
	return productIDs
}

func (doi *DockerBuilderOutputInfo) InputBuildOSArchs(productID ProductID) []OSArchID {
	var osArchIDs []OSArchID
	for k := range doi.InputBuilds[productID] {
		osArchIDs = append(osArchIDs, k)
	}
	sort.Sort(ByOSArchID(osArchIDs))
	return osArchIDs
}

func (doi *DockerBuilderOutputInfo) InputDistProductIDs() []ProductID {
	var productIDs []ProductID
	for k := range doi.InputDists {
		productIDs = append(productIDs, k)
	}
	sort.Sort(ByProductID(productIDs))
	return productIDs
}

func (doi *DockerBuilderOutputInfo) InputDistDistIDs(productID ProductID) []DistID {
	var distIDs []DistID
	for k := range doi.InputDists[productID] {
		distIDs = append(distIDs, k)
	}
	sort.Sort(ByDistID(distIDs))
	return distIDs
}

func (p *DockerParam) ToDockerOutputInfos(productID ProductID, version string) (DockerOutputInfos, error) {
	var dockerIDs []DockerID
	var dockerOutputInfos map[DockerID]DockerBuilderOutputInfo
	if len(p.DockerBuilderParams) > 0 {
		dockerOutputInfos = make(map[DockerID]DockerBuilderOutputInfo)
		for dockerID, dockerBuilderParam := range p.DockerBuilderParams {
			dockerIDs = append(dockerIDs, dockerID)
			currDockerOutputInfo, err := dockerBuilderParam.ToDockerBuilderOutputInfo(productID, version, p.Repository)
			if err != nil {
				return DockerOutputInfos{}, err
			}
			dockerOutputInfos[dockerID] = currDockerOutputInfo
		}
	}
	sort.Sort(ByDockerID(dockerIDs))
	return DockerOutputInfos{
		DockerIDs:                dockerIDs,
		Repository:               p.Repository,
		DockerBuilderOutputInfos: dockerOutputInfos,
	}, nil
}

type DockerBuilderParam struct {
	// DockerBuilder is the builder used to build the Docker image.
	DockerBuilder DockerBuilder

	// DockerfilePath is the path to the Dockerfile that is used to build the Docker image. The path is interpreted
	// relative to ContextDir. The content of the Dockerfile supports using Go templates. The following template
	// parameters can be used in the template:
	//   * {{Product}}: the name of the product
	//   * {{Version}}: the version of the project
	//   * {{InputBuildArtifact(productID, osArch string) (string, error)}}: the path to the build artifact for the specified input product
	//   * {{InputDistArtifacts(productID, distID string) ([]string, error)}}: the paths to the dist artifacts for the specified input product
	//   * {{Tags(productID, dockerID string) ([]string, error)}}: the tags for the specified Docker image
	DockerfilePath string

	// ContextDir is the Docker context directory for building the Docker image.
	ContextDir string

	// Name of directory within ContextDir in which dependencies are linked.
	InputProductsDir string

	// InputBuilds stores the ProductBuildIDs for the input builds. The IDs must be unique and in expanded form.
	InputBuilds []ProductBuildID

	// InputDists stores the ProductDistIDs for the input dists. The IDs must be unique and in expanded form.
	InputDists []ProductDistID

	// TagTemplates contains the templates for the tags that will be used to tat the image generated by this builder.
	// The tag should be the form that would be provided to the "docker tag" command -- for example,
	// "fedora/httpd:version1.0" or "myregistryhost:5000/fedora/httpd:version1.0".
	//
	// The tag templates are rendered using Go templates. The following template parameters can be used in the template:
	//   * {{Product}}: the name of the product
	//   * {{Version}}: the version of the project
	//   * {{Repository}}: the Docker repository
	TagTemplates []string
}

func (p *DockerBuilderParam) ToDockerBuilderOutputInfo(productID ProductID, version, repository string) (DockerBuilderOutputInfo, error) {
	var renderedTags []string
	for _, currTagTemplate := range p.TagTemplates {
		currRenderedTag, err := RenderTemplate(currTagTemplate, nil,
			ProductTemplateFunction(productID),
			VersionTemplateFunction(version),
			RepositoryTemplateFunction(repository),
		)
		if err != nil {
			return DockerBuilderOutputInfo{}, err
		}
		renderedTags = append(renderedTags, currRenderedTag)
	}
	var inputBuilds map[ProductID]map[OSArchID]struct{}
	if len(p.InputBuilds) > 0 {
		inputBuilds = make(map[ProductID]map[OSArchID]struct{})
		for _, productBuildID := range p.InputBuilds {
			productID, buildID, err := productBuildID.Parse()
			if err != nil {
				return DockerBuilderOutputInfo{}, err
			}
			if buildID == (osarch.OSArch{}) {
				return DockerBuilderOutputInfo{}, errors.Errorf("BuildID cannot be empty")
			}
			if _, ok := inputBuilds[productID]; !ok {
				inputBuilds[productID] = make(map[OSArchID]struct{})
			}
			inputBuilds[productID][OSArchID(buildID.String())] = struct{}{}
		}
	}
	var inputDists map[ProductID]map[DistID]struct{}
	if len(p.InputDists) > 0 {
		inputDists = make(map[ProductID]map[DistID]struct{})
		for _, productDistID := range p.InputDists {
			productID, distID := productDistID.Parse()
			if distID == "" {
				return DockerBuilderOutputInfo{}, errors.Errorf("DistID cannot be empty")
			}
			if _, ok := inputDists[productID]; !ok {
				inputDists[productID] = make(map[DistID]struct{})
			}
			inputDists[productID][distID] = struct{}{}
		}
	}
	return DockerBuilderOutputInfo{
		ContextDir:       p.ContextDir,
		DockerfilePath:   p.DockerfilePath,
		InputProductsDir: p.InputProductsDir,
		RenderedTags:     renderedTags,
		InputBuilds:      inputBuilds,
		InputDists:       inputDists,
	}, nil
}
