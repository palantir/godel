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

package dist

import (
	"fmt"
	"io"
	"os"
	"sort"
	"strings"

	"github.com/pkg/errors"

	"github.com/palantir/distgo/distgo"
	"github.com/palantir/distgo/distgo/build"
)

func Products(projectInfo distgo.ProjectInfo, projectParam distgo.ProjectParam, productDistIDs []distgo.ProductDistID, dryRun bool, stdout io.Writer) error {
	productParams, err := distgo.ProductParamsForDistProductArgs(projectParam.Products, productDistIDs...)
	if err != nil {
		return err
	}

	allProducts, _, dependentProducts := distgo.ClassifyProductParams(productParams)
	var productParamsToBuild []distgo.ProductParam
	for _, currProductID := range sortedMapKeys(allProducts) {
		currProduct := projectParam.Products[currProductID]
		if _, ok := dependentProducts[currProductID]; !ok && currProduct.Dist == nil {
			// current product is not a dependency of any specified product and doesn't declare a dist output. In this
			// case, no need to build the build outputs because they will not be used.
			continue
		}
		requiresBuildParam, err := build.RequiresBuild(projectInfo, projectParam.Products[currProductID])
		if err != nil {
			return err
		}
		if requiresBuildParam == nil {
			continue
		}
		productParamsToBuild = append(productParamsToBuild, *requiresBuildParam)
	}
	if len(productParamsToBuild) != 0 {
		if err := build.Run(projectInfo, productParamsToBuild, build.Options{
			Parallel: true,
			DryRun:   dryRun,
		}, stdout); err != nil {
			return err
		}
	}

	// sort dist product tasks in topological order
	targetProducts, topoOrderedIDs, err := distgo.TopoSortProductParams(projectParam, allProducts)
	if err != nil {
		return err
	}
	for _, currProductID := range topoOrderedIDs {
		requiresDistParam, err := RequiresDist(projectInfo, targetProducts[currProductID])
		if err != nil {
			return err
		}
		if requiresDistParam == nil {
			continue
		}
		if err := Run(projectInfo, *requiresDistParam, dryRun, stdout); err != nil {
			return errors.Wrapf(err, "dist failed for %s", currProductID)
		}
	}
	return nil
}

// Run executes the Dist action for the specified product. Produces both the dist output directory and the dist
// artifacts for all of the disters for the product. The outputs for the dependent products for the provided product
// must already exist in the proper locations.
func Run(projectInfo distgo.ProjectInfo, productParam distgo.ProductParam, dryRun bool, stdout io.Writer) error {
	if productParam.Dist == nil {
		distgo.PrintlnOrDryRunPrintln(stdout, fmt.Sprintf("%s does not define a dist configuration; skipping dist", productParam.ID), dryRun)
		return nil
	}

	productOutputInfo, err := productParam.ToProductOutputInfo(projectInfo.Version)
	if err != nil {
		return err
	}

	productTaskOutputInfo, err := distgo.ToProductTaskOutputInfo(projectInfo, productParam)
	if err != nil {
		return err
	}
	distWorkDirs := distgo.ProductDistWorkDirs(projectInfo, productOutputInfo)

	for _, currDistID := range productTaskOutputInfo.Product.DistOutputInfos.DistIDs {
		// create empty output directory
		if !dryRun {
			distWorkDir := distWorkDirs[currDistID]
			// remove output directory if it already exists
			if err := os.RemoveAll(distWorkDir); err != nil {
				return errors.Wrapf(err, "failed to remove dist output directory %s", distWorkDir)
			}
			// create output directory
			if err := os.MkdirAll(distWorkDir, 0755); err != nil {
				return errors.Wrapf(err, "failed to create dist output directory %s", distWorkDir)
			}
		}

		fmt.Fprintln(stdout, distgo.ProductDistArtifactPaths(projectInfo, productOutputInfo))
		distgo.PrintlnOrDryRunPrintln(stdout, fmt.Sprintf("Creating distribution for %s at %v", productParam.ID, strings.Join(distgo.ProductDistArtifactPaths(projectInfo, productOutputInfo)[currDistID], ", ")), dryRun)
		if !dryRun {
			// run dist task
			runDistOutput, err := productParam.Dist.DistParams[currDistID].Dister.RunDist(currDistID, productTaskOutputInfo)
			if err != nil {
				return err
			}
			// execute dist script
			if err := distgo.WriteAndExecuteScript(projectInfo, productParam.Dist.DistParams[currDistID].Script, distgo.DistScriptEnvVariables(currDistID, productTaskOutputInfo), stdout); err != nil {
				return errors.Wrapf(err, "failed to execute dist script")
			}
			// generate dist artifacts
			if err := productParam.Dist.DistParams[currDistID].Dister.GenerateDistArtifacts(currDistID, productTaskOutputInfo, runDistOutput); err != nil {
				return err
			}
		}
		distgo.PrintlnOrDryRunPrintln(stdout, fmt.Sprintf("Finished creating %s distribution for %s", currDistID, productParam.ID), dryRun)
	}
	return nil
}

func sortedMapKeys(m map[distgo.ProductID]struct{}) []distgo.ProductID {
	var out []distgo.ProductID
	for k := range m {
		out = append(out, k)
	}
	sort.Sort(distgo.ByProductID(out))
	return out
}
