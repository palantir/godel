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
	"gopkg.in/yaml.v2"
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

type DockerConfig struct {
	// Repository is the repository that is made available to the tag and Dockerfile templates.
	Repository *string `yaml:"repository"`

	// DockerBuilderParams contains the Docker params for this distribution.
	DockerBuildersConfig *DockerBuildersConfig `yaml:"docker-builders"`
}

func (cfg *DockerConfig) ToParam(defaultCfg DockerConfig, dockerBuilderFactory DockerBuilderFactory) (DockerParam, error) {
	dockerBuilderParams, err := cfg.DockerBuildersConfig.ToParam(cfg.DockerBuildersConfig, dockerBuilderFactory)
	if err != nil {
		return DockerParam{}, err
	}
	return DockerParam{
		Repository:          getConfigStringValue(cfg.Repository, defaultCfg.Repository, ""),
		DockerBuilderParams: dockerBuilderParams,
	}, nil
}

type DockerBuildersConfig map[DockerID]DockerBuilderConfig

func (cfgs *DockerBuildersConfig) ToParam(defaultCfg *DockerBuildersConfig, dockerBuilderFactory DockerBuilderFactory) (map[DockerID]DockerBuilderParam, error) {
	// keys that exist either only in cfgs or only in defaultCfg
	distinctCfgs := make(map[DockerID]DockerBuilderConfig)
	// keys that appear in both cfgs and defaultCfg
	commonCfgIDs := make(map[DockerID]struct{})

	if cfgs != nil {
		for dockerID, dockerCfg := range *cfgs {
			if defaultCfg != nil {
				if _, ok := (*defaultCfg)[dockerID]; ok {
					commonCfgIDs[dockerID] = struct{}{}
					continue
				}
			}
			distinctCfgs[dockerID] = dockerCfg
		}
	}
	if defaultCfg != nil {
		for distID, distCfg := range *defaultCfg {
			if cfgs != nil {
				if _, ok := (*cfgs)[distID]; ok {
					commonCfgIDs[distID] = struct{}{}
					continue
				}
			}
			distinctCfgs[distID] = distCfg
		}
	}

	dockerBuilderParamsMap := make(map[DockerID]DockerBuilderParam)
	// generate parameters for all of the distinct elements
	for dockerID, dockerBuilderCfg := range distinctCfgs {
		currParam, err := dockerBuilderCfg.ToParam(DockerBuilderConfig{}, dockerBuilderFactory)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to generate parameter for Docker configuration %s", dockerID)
		}
		dockerBuilderParamsMap[dockerID] = currParam
	}
	// merge keys that appear in both maps
	for dockerID := range commonCfgIDs {
		currCfg := (*cfgs)[dockerID]
		currParam, err := currCfg.ToParam((*defaultCfg)[dockerID], dockerBuilderFactory)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to generate parameter for dist configuration %s", dockerID)
		}
		dockerBuilderParamsMap[dockerID] = currParam
	}
	return dockerBuilderParamsMap, nil
}

type DockerBuilderConfig struct {
	// Type is the type of the DockerBuilder. This field must be non-nil and non-empty and resolve to a valid DockerBuilder.
	Type *string `yaml:"type"`
	// Config is the YAML configuration content for the DockerBuilder.
	Config *yaml.MapSlice `yaml:"config"`
	// DockerfilePath is the path to the Dockerfile that is used to build the Docker image. The path is interpreted
	// relative to ContextDir. The content of the Dockerfile supports using Go templates. The following template
	// parameters can be used in the template:
	//   * {{Product}}: the name of the product
	//   * {{Version}}: the version of the project
	//   * {{ProductTaskOutputInfo}}: the ProductTaskOutputInfo struct
	DockerfilePath *string `yaml:"dockerfile-path"`
	ContextDir     *string `yaml:"context-dir"`
	// InputProductsDir is the directory in the context dir in which input products are written.
	InputProductsDir *string `yaml:"input-products-dir"`
	// InputBuilds specifies the products whose build outputs should be made available to the Docker build task. The
	// specified products will be hard-linked into the context directory. The referenced products must be this product
	// or one of its declared dependencies.
	InputBuilds *[]ProductBuildID `yaml:"input-builds"`
	// InputDists specifies the products whose dist outputs should be made available to the Docker build task. The
	// specified dists will be hard-linked into the context directory. The referenced products must be this product
	// or one of its declared dependencies.
	InputDists   *[]ProductDistID `yaml:"input-dists"`
	TagTemplates *[]string        `yaml:"tag-templates"`
}

func (cfg *DockerBuilderConfig) ToParam(defaultCfg DockerBuilderConfig, dockerBuilderFactory DockerBuilderFactory) (DockerBuilderParam, error) {
	dockerBuilderType := getConfigStringValue(cfg.Type, defaultCfg.Type, "")
	if dockerBuilderType == "" {
		return DockerBuilderParam{}, errors.Errorf("type must be non-empty")
	}
	dockerBuilder, err := newDockerBuilder(dockerBuilderType, getConfigValue(cfg.Config, defaultCfg.Config, nil).(yaml.MapSlice), dockerBuilderFactory)
	if err != nil {
		return DockerBuilderParam{}, err
	}

	contextDir := getConfigStringValue(cfg.ContextDir, defaultCfg.ContextDir, "")
	if contextDir == "" {
		return DockerBuilderParam{}, errors.Errorf("context-dir must be non-empty")
	}
	tagTemplates := getConfigValue(cfg.TagTemplates, defaultCfg.TagTemplates, nil).([]string)
	if len(tagTemplates) == 0 {
		return DockerBuilderParam{}, errors.Errorf("tag-templates must be non-empty")
	}

	return DockerBuilderParam{
		DockerBuilder:    dockerBuilder,
		DockerfilePath:   getConfigStringValue(cfg.DockerfilePath, defaultCfg.DockerfilePath, "Dockerfile"),
		ContextDir:       contextDir,
		InputProductsDir: getConfigStringValue(cfg.InputProductsDir, defaultCfg.InputProductsDir, ""),
		InputBuilds:      getConfigValue(cfg.InputBuilds, defaultCfg.InputBuilds, nil).([]ProductBuildID),
		InputDists:       getConfigValue(cfg.InputDists, defaultCfg.InputDists, nil).([]ProductDistID),
		TagTemplates:     tagTemplates,
	}, nil
}

func newDockerBuilder(dockerBuilderType string, cfgYML yaml.MapSlice, dockerBuilderFactory DockerBuilderFactory) (DockerBuilder, error) {
	if dockerBuilderType == "" {
		return nil, errors.Errorf("type must be non-empty")
	}
	if dockerBuilderFactory == nil {
		return nil, errors.Errorf("dockerBuilderFactory must be provided")
	}
	cfgYMLBytes, err := yaml.Marshal(cfgYML)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to marshal configuration")
	}
	return dockerBuilderFactory.NewDockerBuilder(dockerBuilderType, cfgYMLBytes)
}
