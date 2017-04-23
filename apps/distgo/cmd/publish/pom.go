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

package publish

import (
	"bytes"
	"text/template"

	"github.com/pkg/errors"

	"github.com/palantir/godel/apps/distgo/templating"
)

const pomTemplate = `<?xml version="1.0" encoding="UTF-8"?>
<project xsi:schemaLocation="http://maven.apache.org/POM/4.0.0 http://maven.apache.org/xsd/maven-4.0.0.xsd" xmlns="http://maven.apache.org/POM/4.0.0"
xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance">
<modelVersion>4.0.0</modelVersion>
<groupId>{{.Publish.GroupID}}</groupId>
<artifactId>{{.ProductName}}</artifactId>
<version>{{.ProductVersion}}</version>
<packaging>{{packagingType}}</packaging>
</project>
`

func generatePOM(cfg templating.Config, distType string) ([]byte, error) {
	funcs := template.FuncMap{
		"packagingType": func() string { return distType },
	}
	t := template.Must(template.New("pom").Funcs(funcs).Parse(pomTemplate))

	pomFileBuf := bytes.Buffer{}
	if err := t.Execute(&pomFileBuf, cfg); err != nil {
		return nil, errors.Wrapf(err, "failed to execute template")
	}
	return pomFileBuf.Bytes(), nil
}
