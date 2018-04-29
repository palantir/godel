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
	"fmt"
	"io"
	"os"
	"time"

	"github.com/pkg/errors"

	"github.com/palantir/distgo/distgo"
	"github.com/palantir/distgo/distgo/dist"
)

func Products(projectInfo distgo.ProjectInfo, projectParam distgo.ProjectParam, configModTime *time.Time, productDistIDs []distgo.ProductDistID, publisher distgo.Publisher, flagVals map[distgo.PublisherFlagName]interface{}, dryRun bool, stdout io.Writer) error {
	// run dist for products (will only run dist for productDistIDs that require dist artifact generation)
	if err := dist.Products(projectInfo, projectParam, configModTime, productDistIDs, dryRun, stdout); err != nil {
		return err
	}

	productParams, err := distgo.ProductParamsForDistProductArgs(projectParam.Products, productDistIDs...)
	if err != nil {
		return err
	}
	for _, currProduct := range productParams {
		if err := Run(projectInfo, currProduct, publisher, flagVals, dryRun, stdout); err != nil {
			return err
		}
	}
	return nil
}

// Run executes the publish action for the specified product. Produces both the dist output directory and the dist
// artifacts for the product. The outputs for the dependent products for the provided product must already exist in the
// proper locations.
func Run(projectInfo distgo.ProjectInfo, productParam distgo.ProductParam, publisher distgo.Publisher, flagVals map[distgo.PublisherFlagName]interface{}, dryRun bool, stdout io.Writer) error {
	if productParam.Dist == nil {
		distgo.PrintlnOrDryRunPrintln(stdout, fmt.Sprintf("%s does not have dist outputs; skipping publish", productParam.ID), dryRun)
		return nil
	}

	// verify that distribution artifacts to publish exists
	productOutputInfo, err := productParam.ToProductOutputInfo(projectInfo.Version)
	if err != nil {
		return errors.Wrapf(err, "failed to compute output info")
	}
	for _, currDistID := range productOutputInfo.DistOutputInfos.DistIDs {
		for _, currArtifactPath := range distgo.ProductDistArtifactPaths(projectInfo, productOutputInfo)[currDistID] {
			if _, err := os.Stat(currArtifactPath); os.IsNotExist(err) {
				return errors.Errorf("distribution artifact for product %s with dist %s does not exist at %s", productParam.ID, currDistID, currArtifactPath)
			}
		}
	}

	// run publish
	productTaskOutputInfo, err := distgo.ToProductTaskOutputInfo(projectInfo, productParam)
	if err != nil {
		return err
	}
	publisherType, err := publisher.TypeName()
	if err != nil {
		return errors.Wrapf(err, "failed to determine type of publisher")
	}
	var publishCfgBytes []byte
	if productParam.Publish != nil {
		publishCfgBytes = productParam.Publish.PublishInfo[distgo.PublisherTypeID(publisherType)].ConfigBytes
	}
	if err := publisher.RunPublish(productTaskOutputInfo, publishCfgBytes, flagVals, dryRun, stdout); err != nil {
		return errors.Wrapf(err, "failed to publish %s using %s publisher", productParam.ID, publisherType)
	}

	return nil
}
