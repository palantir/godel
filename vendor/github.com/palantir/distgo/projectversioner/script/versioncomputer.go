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

package script

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"strings"

	"github.com/pkg/errors"

	"github.com/palantir/distgo/distgo"
)

const TypeName = "script"

type ProjectVersioner struct {
	ScriptContent string
}

func New(scriptContent string) distgo.ProjectVersioner {
	return &ProjectVersioner{
		ScriptContent: scriptContent,
	}
}

func (v *ProjectVersioner) TypeName() (string, error) {
	return TypeName, nil
}

func (v *ProjectVersioner) ProjectVersion(projectDir string) (rVersion string, rErr error) {
	tmpDir, err := ioutil.TempDir("", "godel-distgo-project-versioner-script")
	if err != nil {
		return "", errors.Wrapf(err, "failed to create temporary directory")
	}
	defer func() {
		if err := os.RemoveAll(tmpDir); rErr == nil && err != nil {
			rErr = errors.Wrapf(err, "failed to remove temporary directory")
		}
	}()

	versionScript := path.Join(tmpDir, "version")
	if err := ioutil.WriteFile(versionScript, []byte(v.ScriptContent), 0755); err != nil {
		return "", errors.Wrapf(err, "failed to write version script to %s", versionScript)
	}
	versionScriptCmd := exec.Command(versionScript)
	versionScriptCmd.Dir = projectDir
	versionScriptCmd.Env = append(os.Environ(), fmt.Sprintf("PROJECT_DIR=%s", projectDir))
	outputBytes, err := versionScriptCmd.CombinedOutput()
	output := string(outputBytes)
	if err != nil {
		return "", errors.Wrapf(err, "command %v failed with output %s", versionScriptCmd.Args, output)
	}
	return strings.TrimSpace(output), nil
}
