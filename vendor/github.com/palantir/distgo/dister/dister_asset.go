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

package dister

import (
	"encoding/json"
	"os/exec"
	"strings"

	"github.com/pkg/errors"

	"github.com/palantir/distgo/distgo"
)

type assetDister struct {
	assetPath string
	cfgYML    string
}

func (d *assetDister) TypeName() (string, error) {
	nameCmd := exec.Command(d.assetPath, nameCmdName)
	outputBytes, err := runCommand(nameCmd)
	if err != nil {
		return "", err
	}
	var typeName string
	if err := json.Unmarshal(outputBytes, &typeName); err != nil {
		return "", errors.Wrapf(err, "failed to unmarshal JSON")
	}
	return typeName, nil
}

func (d *assetDister) Artifacts(renderedName string) ([]string, error) {
	artifactsCmd := exec.Command(d.assetPath, artifactPathsCmdName,
		"--"+commonCmdConfigYMLFlagName, d.cfgYML,
		"--"+artifactPathsCmdRenderedNameFlagName, renderedName,
	)
	outputBytes, err := runCommand(artifactsCmd)
	if err != nil {
		return nil, err
	}
	var artifactPaths []string
	if err := json.Unmarshal(outputBytes, &artifactPaths); err != nil {
		return nil, errors.Wrapf(err, "failed to unmarshal JSON")
	}
	return artifactPaths, nil
}

func (d *assetDister) VerifyConfig() error {
	verifyConfigCmd := exec.Command(d.assetPath, verifyConfigCmdName,
		"--"+commonCmdConfigYMLFlagName, d.cfgYML,
	)
	if _, err := runCommand(verifyConfigCmd); err != nil {
		return err
	}
	return nil
}

func (d *assetDister) RunDist(distID distgo.DistID, productTaskOutputInfo distgo.ProductTaskOutputInfo) ([]byte, error) {
	productTaskOutputInfoJSON, err := json.Marshal(productTaskOutputInfo)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to marshal JSON")
	}
	runDistCmd := exec.Command(d.assetPath, runDistCmdName,
		"--"+commonCmdConfigYMLFlagName, d.cfgYML,
		"--"+runDistCmdDistIDFlagName, string(distID),
		"--"+runDistCmdProductTaskOutputInfoFlagName, string(productTaskOutputInfoJSON),
	)
	outputBytes, err := runCommand(runDistCmd)
	if err != nil {
		return nil, err
	}
	return outputBytes, nil
}

func (d *assetDister) GenerateDistArtifacts(distID distgo.DistID, productTaskOutputInfo distgo.ProductTaskOutputInfo, runDistResult []byte) error {
	productTaskOutputInfoJSON, err := json.Marshal(productTaskOutputInfo)
	if err != nil {
		return errors.Wrapf(err, "failed to marshal JSON")
	}
	generateDistArtifactsCmd := exec.Command(d.assetPath, generateDistArtifactsCmdName,
		"--"+commonCmdConfigYMLFlagName, d.cfgYML,
		"--"+generateDistArtifactsCmdDistIDFlagName, string(distID),
		"--"+generateDistArtifactsCmdProductTaskOutputInfoFlagName, string(productTaskOutputInfoJSON),
		"--"+generateDistArtifactsCmdRunDistResultFlagName, string(runDistResult),
	)
	if _, err := runCommand(generateDistArtifactsCmd); err != nil {
		return err
	}
	return nil
}

func runCommand(cmd *exec.Cmd) ([]byte, error) {
	outputBytes, err := cmd.CombinedOutput()
	if err != nil {
		return outputBytes, errors.New(strings.TrimSpace(strings.TrimPrefix(string(outputBytes), "Error: ")))
	}
	return outputBytes, nil
}
