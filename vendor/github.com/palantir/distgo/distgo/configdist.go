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
	"path"
	"sort"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

type DistID string

type ByDistID []DistID

func (a ByDistID) Len() int           { return len(a) }
func (a ByDistID) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByDistID) Less(i, j int) bool { return a[i] < a[j] }

type DistParam struct {
	// OutputDir specifies the default distribution output directory for product distributions created by the "dist"
	// task. The distribution output directory is written to
	// "{{OutputDir}}/{{ID}}/{{Version}}/{{DistID}}/{{NameTemplate}}", and the distribution artifacts are written to
	// "{{OutputDir}}/{{ID}}/{{Version}}/{{DistID}}".
	OutputDir string

	// DistParams contains the dist params for this distribution.
	DistParams map[DistID]DisterParam
}

type DistOutputInfos struct {
	DistOutputDir string                    `json:"distOutputDir"`
	DistIDs       []DistID                  `json:"distIds"`
	DistInfos     map[DistID]DistOutputInfo `json:"distInfos"`
}

func (p *DistParam) ToDistOutputInfos(productID ProductID, version string) (DistOutputInfos, error) {
	var distIDs []DistID
	var distInfos map[DistID]DistOutputInfo
	if len(p.DistParams) > 0 {
		distInfos = make(map[DistID]DistOutputInfo)
		for distID, distParam := range p.DistParams {
			distIDs = append(distIDs, distID)
			distOutputInfo, err := distParam.ToDistOutputInfo(productID, version)
			if err != nil {
				return DistOutputInfos{}, err
			}
			distInfos[distID] = distOutputInfo
		}
		sort.Sort(ByDistID(distIDs))
	}
	return DistOutputInfos{
		DistOutputDir: p.OutputDir,
		DistIDs:       distIDs,
		DistInfos:     distInfos,
	}, nil
}

type DisterParam struct {
	// NameTemplate is the template used for the dist output. The following template parameters can be used in the
	// template:
	//   * {{Product}}: the name of the product
	//   * {{Version}}: the version of the project
	NameTemplate string

	// Script is the content of a script that is written to file a file and run after the initial distribution
	// process but before the artifact generation process. The contents of this value are written to a file and executed
	// with the project directory as the working directory. The script process inherits the environment variables of the
	// Go process and also has the following environment variables defined:
	//
	//   PROJECT_DIR: the root directory of project.
	//   VERSION: the version of the project.
	//   PRODUCT: the name of the product.
	//   DEP_PRODUCT_IDS: the IDs of the dependent products for the product, where each item is delimited by a colon.
	//
	// The following environment variables are defined if the build configuration for the product is non-nil:
	//   BUILD_DIR: the build output directory for the product ("{{OutputDir}}/{{ID}}/{{Version}}").
	//   BUILD_NAME: the rendered NameTemplate for the build for this project.
	//   OS_ARCHS: the OS/Arch combinations for this product, where each item is delimited by a colon.
	//
	// The following environment variables are defined if the dist configuration for the product is non-nil:
	//   DIST_DIR: the distribution output directory for the product ("{{OutputDir}}/{{ID}}/{{Version}}").
	//   DIST_NAME: the rendered NameTemplate for the distribution.
	//   DIST_ARTIFACTS: the dist artifacts for the product, where each item is delimited by a colon.
	//
	// Each dependent product adds the following set of environment variables that start with "{{ID}}_":
	//
	// The following environment variables are defined if the build configuration for the product is non-nil:
	//   {{ID}}_BUILD_DIR: the build output directory for the product ("{{OutputDir}}/{{ID}}/{{Version}}").
	//   {{ID}}_BUILD_NAME: the rendered NameTemplate for the build for this project.
	//   {{ID}}_OS_ARCHS: the OS/Arch combinations for this product, where each item is delimited by a colon.
	//
	// The following environment variables are defined if the dist configuration for the product is non-nil:
	//   {{ID}}_DIST_DIR: the distribution output directory for the product ("{{OutputDir}}/{{ID}}/{{Version}}").
	//   {{ID}}_DIST_NAME: the rendered NameTemplate for the distribution.
	//   {{ID}}_DIST_ARTIFACTS: the dist artifacts for the product, where each item is delimited by a colon.
	Script string

	// Dister is the Dister that performs the dist operation for this parameter.
	Dister Dister
}

type DistOutputInfo struct {
	DistNameTemplateRendered string   `json:"distNameTemplateRendered"`
	DistArtifactNames        []string `json:"distArtifactNames"`
}

func (p *DisterParam) ToDistOutputInfo(productID ProductID, version string) (DistOutputInfo, error) {
	renderedName, err := renderNameTemplate(p.NameTemplate, productID, version)
	if err != nil {
		return DistOutputInfo{}, errors.Wrapf(err, "failed to render name template")
	}
	artifactNames, err := p.Dister.Artifacts(renderedName)
	if err != nil {
		return DistOutputInfo{}, errors.Wrapf(err, "failed to determine artifact names")
	}
	return DistOutputInfo{
		DistNameTemplateRendered: renderedName,
		DistArtifactNames:        artifactNames,
	}, nil
}

type DistConfig struct {
	// OutputDir specifies the default distribution output directory for product distributions created by the "dist"
	// task. The distribution output directory is written to
	// "{{OutputDir}}/{{ID}}/{{Version}}/{{NameTemplate}}", and the distribution artifacts are written to
	// "{{OutputDir}}/{{ID}}/{{Version}}".
	//
	// If a value is not specified, "out/dist" is used as the default value.
	OutputDir *string `yaml:"output-dir"`

	// Disters is the configuration for the disters for this product. The YAML representation can be a single DisterConfig
	// or a map[DistID]DisterConfig.
	Disters *DistersConfig `yaml:"disters"`
}

// ToParam returns the DistParam represented by the receiver *DisterConfig and the provided default DisterConfig. If a
// config value is specified (non-nil) in the receiver config, it is used. If a config value is not specified in the
// receiver config but is specified in the default config, the default config value is used. If a value is not specified
// in either configuration, the program-specified default value (if any) is used.
func (cfg *DistConfig) ToParam(scriptIncludes string, defaultCfg DistConfig, disterFactory DisterFactory) (DistParam, error) {
	outputDir := getConfigStringValue(cfg.OutputDir, defaultCfg.OutputDir, "out/dist")
	if path.IsAbs(outputDir) {
		return DistParam{}, errors.Errorf("output-dir cannot be specified as an absolute path")
	}
	disters, err := cfg.Disters.ToParam(defaultCfg.Disters, scriptIncludes, disterFactory)
	if err != nil {
		return DistParam{}, err
	}
	return DistParam{
		OutputDir:  outputDir,
		DistParams: disters,
	}, nil
}

type DisterConfig struct {
	// Type is the type of the dister. This field must be non-nil and non-empty and resolve to a valid Dister.
	Type *string `yaml:"type"`

	// Config is the YAML configuration content for the dister.
	Config *yaml.MapSlice `yaml:"config"`

	// NameTemplate is the template used for the executable output. The following template parameters can be used in the
	// template:
	//   * {{Product}}: the name of the product.
	//   * {{Version}}: the version of the project.
	//
	// If a value is not specified, "{{Product}}-{{Version}}" is used as the default value.
	NameTemplate *string `yaml:"name-template"`

	// Script is the content of a script that is written to file a file and run after the initial distribution
	// process but before the artifact generation process. The content of this value is written to a file and executed
	// with the project directory as the working directory. The script process inherits the environment variables of the
	// Go process and also has the following environment variables defined:
	//
	//   PROJECT_DIR: the root directory of project.
	//   VERSION: the version of the project.
	//   PRODUCT: the name of the product.
	//   DEP_PRODUCT_IDS: the IDs of the dependent products for the product, where each item is delimited by a colon.
	//
	// The following environment variables are defined if the build configuration for the product is non-nil:
	//   BUILD_DIR: the build output directory for the product ("{{OutputDir}}/{{ID}}/{{Version}}").
	//   BUILD_NAME: the rendered NameTemplate for the build for this project.
	//   OS_ARCHS: the OS/Arch combinations for this product, where each item is delimited by a colon.
	//
	// The following environment variables are defined if the dist configuration for the product is non-nil:
	//   DIST_DIR: the distribution output directory for the product ("{{OutputDir}}/{{ID}}/{{Version}}").
	//   DIST_NAME: the rendered NameTemplate for the distribution.
	//   DIST_ARTIFACTS: the dist artifacts for the product, where each item is delimited by a colon.
	//
	// Each dependent product adds the following set of environment variables that start with "{{ID}}_":
	//
	// The following environment variables are defined if the build configuration for the product is non-nil:
	//   {{ID}}_BUILD_DIR: the build output directory for the product ("{{OutputDir}}/{{ID}}/{{Version}}").
	//   {{ID}}_BUILD_NAME: the rendered NameTemplate for the build for this project.
	//   {{ID}}_OS_ARCHS: the OS/Arch combinations for this product, where each item is delimited by a colon.
	//
	// The following environment variables are defined if the dist configuration for the product is non-nil:
	//   {{ID}}_DIST_DIR: the distribution output directory for the product ("{{OutputDir}}/{{ID}}/{{Version}}").
	//   {{ID}}_DIST_NAME: the rendered NameTemplate for the distribution.
	//   {{ID}}_DIST_ARTIFACTS: the dist artifacts for the product, where each item is delimited by a colon.
	Script *string `yaml:"script"`
}

func (cfg *DisterConfig) ToParam(defaultCfg DisterConfig, scriptIncludes string, disterFactory DisterFactory) (DisterParam, error) {
	disterType := getConfigStringValue(cfg.Type, defaultCfg.Type, "")
	if disterType == "" {
		return DisterParam{}, errors.Errorf("dister type must be specified for DisterConfig")
	}
	dister, err := newDister(disterType, getConfigValue(cfg.Config, defaultCfg.Config, nil).(yaml.MapSlice), disterFactory)
	if err != nil {
		return DisterParam{}, err
	}

	return DisterParam{
		NameTemplate: getConfigStringValue(cfg.NameTemplate, defaultCfg.NameTemplate, "{{Product}}-{{Version}}"),
		Script:       createScriptContent(getConfigStringValue(cfg.Script, defaultCfg.Script, ""), scriptIncludes),
		Dister:       dister,
	}, nil
}

func newDister(disterType string, cfgYML yaml.MapSlice, disterFactory DisterFactory) (Dister, error) {
	if disterType == "" {
		return nil, errors.Errorf("dister type must be non-empty")
	}
	if disterFactory == nil {
		return nil, errors.Errorf("disterFactory must be provided")
	}
	cfgYMLBytes, err := yaml.Marshal(cfgYML)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to marshal configuration")
	}
	return disterFactory.NewDister(disterType, cfgYMLBytes)
}

type DistersConfig map[DistID]DisterConfig

func (cfgs *DistersConfig) ToParam(defaultCfg *DistersConfig, scriptIncludes string, disterFactory DisterFactory) (map[DistID]DisterParam, error) {
	// keys that exist either only in cfgs or only in defaultCfg
	distinctCfgs := make(map[DistID]DisterConfig)
	// keys that appear in both cfgs and defaultCfg
	commonCfgIDs := make(map[DistID]struct{})

	if cfgs != nil {
		for distID, distCfg := range *cfgs {
			if defaultCfg != nil {
				if _, ok := (*defaultCfg)[distID]; ok {
					commonCfgIDs[distID] = struct{}{}
					continue
				}
			}
			distinctCfgs[distID] = distCfg
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

	distParamsMap := make(map[DistID]DisterParam)
	// generate parameters for all of the distinct elements
	for distID, distCfg := range distinctCfgs {
		currParam, err := distCfg.ToParam(DisterConfig{}, scriptIncludes, disterFactory)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to generate parameter for dist configuration %s", distID)
		}
		distParamsMap[distID] = currParam
	}
	// merge keys that appear in both maps
	for distID := range commonCfgIDs {
		currCfg := (*cfgs)[distID]
		currParam, err := currCfg.ToParam((*defaultCfg)[distID], scriptIncludes, disterFactory)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to generate parameter for dist configuration %s", distID)
		}
		distParamsMap[distID] = currParam
	}
	return distParamsMap, nil
}

func (cfgs *DistersConfig) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var single DisterConfig
	if err := unmarshal(&single); err == nil && single.Type != nil {
		// only consider single configuration valid if it unmarshals and "type" key is explicitly specified
		*cfgs = DistersConfig{
			DistID(*single.Type): single,
		}
		return nil
	}

	var multiple map[DistID]DisterConfig
	if err := unmarshal(&multiple); err != nil {
		return errors.Errorf("failed to unmarshal configuration as single DisterConfig or as map[DistID]DisterConfig")
	}
	if len(multiple) == 0 {
		return errors.Errorf(`if "dist" key is specified, there must be at least one dist`)
	}
	*cfgs = multiple
	return nil
}

func ToMapSlice(in interface{}) (yaml.MapSlice, error) {
	bytes, err := yaml.Marshal(in)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to marshal input as YAML")
	}
	var mapSlice yaml.MapSlice
	if err = yaml.Unmarshal(bytes, &mapSlice); err != nil {
		return nil, errors.Wrapf(err, "failed to unmarshal bytes as MapSlice")
	}
	return mapSlice, nil
}
