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
	"github.com/palantir/pkg/matcher"
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

	// ProjectVersionerParam provides the operation for determining the project version.
	ProjectVersionerParam ProjectVersionerParam

	// Exclude is a matcher that matches any directories that should be ignored as main files. Only relevant if products
	// are not specified.
	Exclude matcher.Matcher
}

func (p *ProjectParam) ProjectInfo(projectDir string) (ProjectInfo, error) {
	version, err := p.ProjectVersionerParam.ProjectVersioner.ProjectVersion(projectDir)
	if err != nil {
		return ProjectInfo{}, err
	}
	return ProjectInfo{
		ProjectDir: projectDir,
		Version:    version,
	}, nil
}
