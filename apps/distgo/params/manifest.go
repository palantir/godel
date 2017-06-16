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

package params

import (
	"strings"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

type slsManifest struct {
	ManifestVersion string                 `yaml:"manifest-version"`
	ProductGroup    string                 `yaml:"product-group"`
	ProductName     string                 `yaml:"product-name"`
	ProductVersion  string                 `yaml:"product-version"`
	ProductType     string                 `yaml:"product-type,omitempty"`
	Extensions      map[string]interface{} `yaml:"extensions,omitempty"`
}

func GetManifest(groupID, name, version, productType string, extensions map[string]interface{}) (string, error) {
	var missingRequired []string
	if groupID == "" {
		missingRequired = append(missingRequired, "group-id")
	}
	if name == "" {
		missingRequired = append(missingRequired, "product-name")
	}
	if version == "" {
		missingRequired = append(missingRequired, "product-version")
	}
	if len(missingRequired) > 0 {
		return "", errors.Errorf("required properties were missing: %s", strings.Join(missingRequired, ", "))
	}

	m := slsManifest{
		ManifestVersion: "1.0",
		ProductGroup:    groupID,
		ProductName:     name,
		ProductVersion:  version,
		Extensions:      extensions,
	}
	if productType != "" {
		m.ProductType = productType
	}
	manifestBytes, err := yaml.Marshal(m)
	if err != nil {
		return "", errors.Wrapf(err, "failed to marshal %v as YAML", m)
	}
	return string(manifestBytes), nil
}
