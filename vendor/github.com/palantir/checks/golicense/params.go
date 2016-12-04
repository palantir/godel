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

package golicense

import (
	"fmt"
	"sort"
	"strings"

	"github.com/palantir/pkg/matcher"
	"github.com/pkg/errors"
)

type LicenseParams struct {
	Header        string
	CustomHeaders CustomLicenseParams
	Exclude       matcher.Matcher
}

type CustomLicenseParams interface {
	Len() int
	headers() []CustomLicenseParam
}

type customLicenseParams []CustomLicenseParam

func (p customLicenseParams) Len() int {
	return len(p)
}

func (p customLicenseParams) headers() []CustomLicenseParam {
	return p
}

func (p customLicenseParams) validate() error {
	var emptyNameParams []CustomLicenseParam
	nameToParams := make(map[string][]CustomLicenseParam)

	for _, v := range p {
		if v.Name == "" {
			emptyNameParams = append(emptyNameParams, v)
		}
		nameToParams[v.Name] = append(nameToParams[v.Name], v)
	}

	if len(emptyNameParams) > 0 {
		return errors.Errorf("custom header entries have blank names: %+v", emptyNameParams)
	}

	var nameCollisionMsgs []string
	for k, v := range nameToParams {
		if len(v) > 1 {
			nameCollisionMsgs = append(nameCollisionMsgs, fmt.Sprintf("%s: %+v", k, v))
		}
	}
	if len(nameCollisionMsgs) > 0 {
		return errors.Errorf(strings.Join(append([]string{"multiple custom header entries have the same name:"}, nameCollisionMsgs...), "\n\t"))
	}

	// map from path to custom header entries that have the path
	pathsToCustomEntries := make(map[string][]string)
	for _, ch := range p {
		for _, path := range ch.IncludePaths {
			pathsToCustomEntries[path] = append(pathsToCustomEntries[path], ch.Name)
		}
	}
	var customPathCollisionMsgs []string
	sortedKeys := make([]string, 0, len(pathsToCustomEntries))
	for k := range pathsToCustomEntries {
		sortedKeys = append(sortedKeys, k)
	}
	sort.Strings(sortedKeys)
	for _, k := range sortedKeys {
		v := pathsToCustomEntries[k]
		if len(v) > 1 {
			customPathCollisionMsgs = append(customPathCollisionMsgs, fmt.Sprintf("%s: %s", k, strings.Join(v, ", ")))
		}
	}
	if len(customPathCollisionMsgs) > 0 {
		return errors.Errorf(strings.Join(append([]string{"the same path is defined by multiple custom header entries:"}, customPathCollisionMsgs...), "\n\t"))
	}

	return nil
}

func NewCustomLicenseParams(customHeaders []CustomLicenseParam) (CustomLicenseParams, error) {
	params := customLicenseParams(customHeaders)
	if err := params.validate(); err != nil {
		return nil, err
	}
	return params, nil
}

type CustomLicenseParam struct {
	Name         string
	Header       string
	IncludePaths []string
}
