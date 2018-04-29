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

package publisher

import (
	"encoding/json"
	"io"
	"os/exec"
	"strconv"

	"github.com/pkg/errors"

	"github.com/palantir/distgo/distgo"
)

type assetPublisher struct {
	assetPath string
}

func (p *assetPublisher) TypeName() (string, error) {
	nameCmd := exec.Command(p.assetPath, nameCmdName)
	outputBytes, err := nameCmd.CombinedOutput()
	if err != nil {
		return "", errors.Wrapf(err, "command %v failed with output %s", nameCmd.Args, string(outputBytes))
	}
	var typeName string
	if err := json.Unmarshal(outputBytes, &typeName); err != nil {
		return "", errors.Wrapf(err, "failed to unmarshal JSON")
	}
	return typeName, nil
}

func (p *assetPublisher) Flags() ([]distgo.PublisherFlag, error) {
	flagsCmd := exec.Command(p.assetPath, flagsCmdName)
	outputBytes, err := flagsCmd.CombinedOutput()
	if err != nil {
		return nil, errors.Wrapf(err, "command %v failed with output %s", flagsCmd.Args, string(outputBytes))
	}
	var flags []distgo.PublisherFlag
	if err := json.Unmarshal(outputBytes, &flags); err != nil {
		return nil, errors.Wrapf(err, "failed to unmarshal JSON")
	}
	return flags, nil
}

func (p *assetPublisher) RunPublish(productTaskOutputInfo distgo.ProductTaskOutputInfo, cfgYML []byte, flagVals map[distgo.PublisherFlagName]interface{}, dryRun bool, stdout io.Writer) error {
	productTaskOutputInfoJSON, err := json.Marshal(productTaskOutputInfo)
	if err != nil {
		return errors.Wrapf(err, "failed to marshal JSON for productTaskOutputInfo")
	}
	flagValsJSON, err := json.Marshal(flagVals)
	if err != nil {
		return errors.Wrapf(err, "failed to marshal JSON for flagVals")
	}

	runPublishCmd := exec.Command(p.assetPath, runPublishCmdName,
		"--"+runPublishCmdProductTaskOutputInfoFlagName, string(productTaskOutputInfoJSON),
		"--"+runPublishCmdConfigYMLFlagName, string(cfgYML),
		"--"+runPublishCmdFlagValsFlagName, string(flagValsJSON),
		"--"+runPublishCmdDryRunFlagName, strconv.FormatBool(dryRun),
	)
	runPublishCmd.Stdout = stdout

	if err := runPublishCmd.Run(); err != nil {
		return errors.Wrapf(err, "command %v failed", runPublishCmd.Args)
	}
	return nil
}
