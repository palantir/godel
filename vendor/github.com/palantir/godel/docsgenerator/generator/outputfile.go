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

package generator

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"regexp"
	"strings"

	"github.com/pkg/errors"
)

func writeOutputFiles(outputDir string, inFile inputFileWithParsedContent, fromImage string) (string, error) {
	currOutputDir := path.Join(outputDir, inFile.FileInfo.OutputDirName())
	if err := os.MkdirAll(currOutputDir, 0755); err != nil {
		return "", errors.Wrapf(err, "failed to create output directory")
	}
	scriptContent := bashScript(inFile.ParsedContent.TutorialCodeParts)
	if err := ioutil.WriteFile(path.Join(currOutputDir, inFile.FileInfo.OutputScriptFileName()), []byte(scriptContent), 0755); err != nil {
		return "", errors.Wrapf(err, "failed to write script file")
	}
	dockerfileContent := dockerFile(fromImage, inFile.FileInfo.OutputScriptFileName())
	if err := ioutil.WriteFile(path.Join(currOutputDir, "Dockerfile"), []byte(dockerfileContent), 0644); err != nil {
		return "", errors.Wrapf(err, "failed to write Dockerfile")
	}
	return currOutputDir, nil
}

func dockerFile(fromImage, scriptFileName string) string {
	dockerFileContent := `FROM {{FROM_IMAGE}}

ADD {{SCRIPT_FILE}} /scripts/
RUN /scripts/{{SCRIPT_FILE}} 2>&1
`
	dockerFileContent = strings.Replace(dockerFileContent, "{{FROM_IMAGE}}", fromImage, -1)
	dockerFileContent = strings.Replace(dockerFileContent, "{{SCRIPT_FILE}}", scriptFileName, -1)
	return dockerFileContent
}

func bashScript(codeParts []tutorialCodePart) string {
	outputCodeParts := []string{bashScriptCommonCode()}
	for _, currPart := range codeParts {
		currCode := currPart.Code
		if currPart.WantFail {
			currCode += " || true"
		}
		outputCodeParts = append(outputCodeParts, fmt.Sprintf(bashScriptSingleCmdCode, currCode))
	}
	return strings.Join(outputCodeParts, "")
}

const (
	bashScriptCommonCodeTmpl = `#!/usr/bin/env bash
print_then_run () {
    echo "%s"
    echo "$1"
    echo "%s"

    echo "%s"
    eval "$1"
    echo "%s"
}
`
	bashScriptSingleCmdCode = `
set +e
read -d '' ACTION <<"EOF"
%s
EOF
set -e
print_then_run "$ACTION"
`

	bashRunStart = "BASH_RUN:-------------"
	outputStart  = "OUTPUT:---------------"
	endDelimiter = "----------------------"
)

var bashOutputRegexp = regexp.MustCompile(
	`(?sm)` +
		`^` + regexp.QuoteMeta(bashRunStart) + `\n` +
		`(.*?)\n` +
		regexp.QuoteMeta(endDelimiter) + `\n` +
		regexp.QuoteMeta(outputStart) + `\n` +
		`(.*?)\n?` +
		regexp.QuoteMeta(endDelimiter) + `$`)

func bashScriptCommonCode() string {
	return fmt.Sprintf(bashScriptCommonCodeTmpl,
		bashRunStart,
		endDelimiter,
		outputStart,
		endDelimiter,
	)
}

type bashRunCmd struct {
	cmd    string
	output string
}

func (c *bashRunCmd) String() string {
	output := `➜ ` + c.cmd
	if c.output != "" {
		output += "\n" + c.output
	}
	return output
}

func parseBashRunCmdFromOutput(output string) []bashRunCmd {
	var cmds []bashRunCmd
	for _, submatch := range bashOutputRegexp.FindAllStringSubmatch(output, -1) {
		cmds = append(cmds, bashRunCmd{
			cmd:    submatch[1],
			output: submatch[2],
		})
	}
	return cmds
}
