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

package slsspec

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/palantir/pkg/matcher"
	"github.com/palantir/pkg/specdir"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

const (
	Deployment             = "deployment"
	Manifest               = "manifest"
	Service                = "service"
	ServiceBin             = "service/bin"
	InitSh                 = "init.sh"
	serviceNameTemplate    = "ServiceName"
	serviceVersionTemplate = "ServiceVersion"
)

func New() specdir.LayoutSpec {
	return specdir.NewLayoutSpec(
		specdir.Dir(specdir.CompositeName(specdir.TemplateName(serviceNameTemplate), specdir.LiteralName("-"), specdir.TemplateName(serviceVersionTemplate)), "",
			specdir.Dir(specdir.LiteralName("deployment"), Deployment,
				specdir.File(specdir.LiteralName("manifest.yml"), Manifest),
			),
			specdir.Dir(specdir.LiteralName("service"), Service,
				specdir.Dir(specdir.LiteralName("bin"), ServiceBin,
					specdir.File(specdir.LiteralName("init.sh"), InitSh),
				),
			),
		),
		true,
	)
}

func TemplateValues(product, version string) specdir.TemplateValues {
	return specdir.TemplateValues{
		serviceNameTemplate:    product,
		serviceVersionTemplate: version,
	}
}

func Validate(rootDir string, values specdir.TemplateValues, excludeYML matcher.Matcher) error {
	_, err := specdir.New(rootDir, New(), values, specdir.Validate)
	if err != nil {
		return err
	}
	// check validity of all YML files except those in service directory (binaries may contain invalid YML)
	invalidYMLFiles, err := invalidYMLFiles(rootDir, func(path string) bool {
		relPath, err := filepath.Rel(rootDir, path)
		return err == nil && excludeYML != nil && excludeYML.Match(relPath)
	})
	if err != nil {
		return errors.Wrapf(err, "failed to determine validity of YML files")
	}
	if len(invalidYMLFiles) > 0 {
		msg := fmt.Sprintf("invalid YML files: %v\n", invalidYMLFiles)
		msg += "If these files are known to be correct, exclude them from validation using the SLS YML validation exclude matcher."
		return errors.Errorf(msg)
	}
	return nil
}

func invalidYMLFiles(rootDir string, skip func(path string) bool) ([]string, error) {
	var invalidYMLFiles []string
	if err := filepath.Walk(rootDir, func(currFile string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() || skip(currFile) {
			return nil
		}
		if strings.HasSuffix(info.Name(), ".yml") || strings.HasSuffix(info.Name(), ".yaml") {
			bytes, err := ioutil.ReadFile(currFile)
			if err != nil {
				return err
			}
			var m interface{}
			if err := yaml.Unmarshal(bytes, &m); err != nil {
				rootRelPath := currFile
				if relPath, err := filepath.Rel(rootDir, currFile); err == nil {
					rootRelPath = path.Join(path.Base(rootDir), relPath)
				}
				// if YML file fails unmarshal, treat it as invalid
				invalidYMLFiles = append(invalidYMLFiles, rootRelPath)
			}
		}
		return nil
	}); err != nil {
		return nil, err
	}
	return invalidYMLFiles, nil
}
