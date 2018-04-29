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

package manual

import (
	"os"

	"github.com/pkg/errors"

	"github.com/palantir/distgo/distgo"
)

const TypeName = "manual" // distribution that consists of a distribution whose output is created by the distribution script

type Dister struct {
	Extension string
}

func (d *Dister) TypeName() (string, error) {
	return TypeName, nil
}

func (d *Dister) Artifacts(renderedNameTemplate string) ([]string, error) {
	outputFileName := renderedNameTemplate
	if d.Extension != "" {
		outputFileName += "." + d.Extension
	}
	return []string{outputFileName}, nil
}

func (d *Dister) RunDist(distID distgo.DistID, productTaskOutputInfo distgo.ProductTaskOutputInfo) ([]byte, error) {
	// manual dister does not perform any actions (all actions are preformed by script)
	return nil, nil
}

func (d *Dister) GenerateDistArtifacts(distID distgo.DistID, productTaskOutputInfo distgo.ProductTaskOutputInfo, runDistResult []byte) error {
	outputArtifactPaths := productTaskOutputInfo.ProductDistArtifactPaths()[distID]
	if len(outputArtifactPaths) != 1 {
		return errors.Errorf("manual distribution must produce a single artifact")
	}

	// manual dister depends on the script to generate the declared output -- verify that the output exists, and fail if it does not
	fi, err := os.Stat(outputArtifactPaths[0])
	if os.IsNotExist(err) {
		return errors.Wrapf(err, "expected output does not exist at %s", outputArtifactPaths[0])
	}
	// output should not be a directory
	if fi.IsDir() {
		return errors.Errorf("output at %s is a directory", outputArtifactPaths[0])
	}
	return nil
}
