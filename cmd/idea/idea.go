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

package idea

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"text/template"

	"github.com/pkg/errors"
)

const (
	defaultGoSDK       = "Go"
	imlTemplateContent = `<?xml version="1.0" encoding="UTF-8"?>
<module type="GO_MODULE" version="4">
  <component name="NewModuleRootManager" inherit-compiler-output="true">
    <exclude-output />
    <content url="file://$MODULE_DIR$" />
    <orderEntry type="jdk" jdkName="{{.GoSDK}}" jdkType="Go SDK" />
    <orderEntry type="sourceFolder" forTests="false" />
  </component>
</module>
`
	iprTemplateContent = `<?xml version="1.0" encoding="UTF-8"?>
<project version="4">
  <component name="ProjectModuleManager">
    <modules>
      <module fileurl="file://$PROJECT_DIR$/{{.ProjectName}}.iml" filepath="$PROJECT_DIR$/{{.ProjectName}}.iml" />
    </modules>
  </component>
  <component name="ProjectRootManager" version="2" default="false" assert-keyword="false" jdk-15="false" project-jdk-name="{{.GoSDK}}" project-jdk-type="Go SDK" />
  <component name="ProjectTasksOptions">
    <TaskOptions isEnabled="true">
      <option name="arguments" value="format runAll $FilePathRelativeToProjectRoot$" />
      <option name="checkSyntaxErrors" value="true" />
      <option name="description" value="" />
      <option name="exitCodeBehavior" value="ERROR" />
      <option name="fileExtension" value="go" />
      <option name="immediateSync" value="false" />
      <option name="name" value="godel" />
      <option name="output" value="" />
      <option name="outputFilters">
        <array />
      </option>
      <option name="outputFromStdout" value="false" />
      <option name="program" value="$ProjectFileDir$/godelw" />
      <option name="scopeName" value="Changed Files" />
      <option name="trackOnlyRoot" value="false" />
      <option name="workingDir" value="$ProjectFileDir$" />
      <envs />
    </TaskOptions>
  </component>
</project>
`
)

func CreateIdeaFiles(rootDir string) error {
	projectName := path.Base(rootDir)

	buffer := bytes.Buffer{}
	templateValues := map[string]string{
		"GoSDK":       defaultGoSDK,
		"ProjectName": projectName,
	}
	imlTemplate := template.Must(template.New("iml").Parse(imlTemplateContent))
	if err := imlTemplate.Execute(&buffer, templateValues); err != nil {
		return errors.Wrapf(err, "Failed to execute template %s with values %v", imlTemplate, templateValues)
	}

	imlFilePath := path.Join(rootDir, projectName+".iml")
	if err := ioutil.WriteFile(imlFilePath, buffer.Bytes(), 0644); err != nil {
		return errors.Wrapf(err, "Failed to write .iml file to %s", imlFilePath)
	}

	iprTemplate := template.Must(template.New("modules").Parse(iprTemplateContent))
	buffer = bytes.Buffer{}
	if err := iprTemplate.Execute(&buffer, templateValues); err != nil {
		return errors.Wrapf(err, "Failed to execute template %s with values %v", imlTemplate, templateValues)
	}

	iprFilePath := path.Join(rootDir, projectName+".ipr")
	if err := ioutil.WriteFile(iprFilePath, buffer.Bytes(), 0644); err != nil {
		return errors.Wrapf(err, "Failed to write .ipr file to %s", iprFilePath)
	}

	return nil
}

func CleanIdeaFiles(rootDir string) error {
	projectName := path.Base(rootDir)
	for _, ext := range []string{"iml", "ipr", "iws"} {
		currPath := path.Join(rootDir, fmt.Sprintf("%v.%v", projectName, ext))
		if err := os.Remove(currPath); err != nil && !os.IsNotExist(err) {
			return errors.Wrapf(err, "Failed to remove file %s", currPath)
		}
	}
	return nil
}
