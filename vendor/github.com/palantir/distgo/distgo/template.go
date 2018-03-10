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
	"bytes"
	"text/template"

	"github.com/pkg/errors"
)

type TemplateFunction func(fnMap template.FuncMap)

func ProductTemplateFunction(productID ProductID) TemplateFunction {
	return TemplateValueFunction("Product", string(productID))
}

func VersionTemplateFunction(version string) TemplateFunction {
	return TemplateValueFunction("Version", version)
}

func GroupIDTemplateFunction(groupID string) TemplateFunction {
	return TemplateValueFunction("GroupID", groupID)
}

func RepositoryTemplateFunction(repository string) TemplateFunction {
	return TemplateValueFunction("Repository", repository)
}

func TemplateValueFunction(key string, val interface{}) TemplateFunction {
	return func(fnMap template.FuncMap) {
		fnMap[key] = func() interface{} {
			return val
		}
	}
}

func RenderTemplate(tmplContent string, data interface{}, fns ...TemplateFunction) (string, error) {
	tmpl := template.New("distgoTemplate")
	tmplFuncs := make(map[string]interface{})
	for _, fn := range fns {
		fn(tmplFuncs)
	}
	tmpl.Funcs(tmplFuncs)
	var err error
	tmpl, err = tmpl.Parse(tmplContent)
	if err != nil {
		return "", errors.Wrapf(err, "failed to parse template %s", tmplContent)
	}
	output := &bytes.Buffer{}
	if err := tmpl.Execute(output, data); err != nil {
		return "", errors.Wrapf(err, "failed to execute template")
	}
	return output.String(), nil
}

func renderNameTemplate(nameTemplate string, productID ProductID, version string) (string, error) {
	return RenderTemplate(nameTemplate, nil,
		ProductTemplateFunction(productID),
		VersionTemplateFunction(version),
	)
}
