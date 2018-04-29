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
	"regexp"
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
	// Go process and also has dist-related environment variables. Refer to the documentation for the
	// distgo.DistScriptEnvVariables function for the extra environment variables.
	Script string

	// InputDir specifies the configuration for copying files from an input directory.
	InputDir InputDirParam

	// Dister is the Dister that performs the dist operation for this parameter.
	Dister Dister
}

type InputDirParam struct {
	Path   string
	Ignore []*regexp.Regexp
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
