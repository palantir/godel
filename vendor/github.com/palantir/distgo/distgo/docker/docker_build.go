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

package docker

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"sort"
	"text/template"

	"github.com/palantir/godel/pkg/osarch"
	"github.com/pkg/errors"

	"github.com/palantir/distgo/distgo"
	"github.com/palantir/distgo/distgo/dist"
)

func BuildProducts(projectInfo distgo.ProjectInfo, projectParam distgo.ProjectParam, productDockerIDs []distgo.ProductDockerID, verbose, dryRun bool, stdout io.Writer) error {
	// determine products that match specified productDockerIDs
	productParams, err := distgo.ProductParamsForDockerProductArgs(projectParam.Products, productDockerIDs...)
	if err != nil {
		return err
	}

	// create a ProductDistID for all of the products for which a Docker action will be run
	var productDistIDs []distgo.ProductDistID
	for _, currProductParam := range productParams {
		productDistIDs = append(productDistIDs, distgo.ProductDistID(currProductParam.ID))
	}
	// run dist for products (will only run dist for productDistIDs that require dist artifact generation)
	if err := dist.Products(projectInfo, projectParam, productDistIDs, dryRun, stdout); err != nil {
		return err
	}

	// sort Docker product tasks in topological order
	allProducts, _, _ := distgo.ClassifyProductParams(productParams)
	targetProducts, topoOrderedIDs, err := distgo.TopoSortProductParams(projectParam, allProducts)
	if err != nil {
		return err
	}
	for _, currID := range topoOrderedIDs {
		currProduct := targetProducts[currID]
		if err := RunBuild(projectInfo, currProduct, verbose, dryRun, stdout); err != nil {
			return err
		}
	}
	return nil
}

// RunBuild executes the Docker image build action for the specified product. The Docker outputs for all of the
// dependent products for the provided product must already exist, and the dist outputs for the current product and all
// of its dependent products must also exist in the proper locations.
func RunBuild(projectInfo distgo.ProjectInfo, productParam distgo.ProductParam, verbose, dryRun bool, stdout io.Writer) error {
	if productParam.Docker == nil {
		distgo.PrintlnOrDryRunPrintln(stdout, fmt.Sprintf("%s does not have Docker outputs; skipping build", productParam.ID), dryRun)
		return nil
	}

	var dockerIDs []distgo.DockerID
	for k := range productParam.Docker.DockerBuilderParams {
		dockerIDs = append(dockerIDs, k)
	}
	sort.Sort(distgo.ByDockerID(dockerIDs))

	productTaskOutputInfo, err := distgo.ToProductTaskOutputInfo(projectInfo, productParam)
	if err != nil {
		return err
	}

	allBuildArtifactPaths := productTaskOutputInfo.ProductDockerBuildArtifactPaths()
	allDistArtifactPaths := productTaskOutputInfo.ProductDockerDistArtifactPaths()

	for _, dockerID := range dockerIDs {
		if err := runSingleDockerBuild(
			projectInfo,
			productParam.ID,
			dockerID,
			productParam.Docker.DockerBuilderParams[dockerID],
			productTaskOutputInfo,
			allBuildArtifactPaths[dockerID],
			allDistArtifactPaths[dockerID],
			verbose,
			dryRun,
			stdout,
		); err != nil {
			return err
		}
	}
	return nil
}

func runSingleDockerBuild(
	projectInfo distgo.ProjectInfo,
	productID distgo.ProductID,
	dockerID distgo.DockerID,
	dockerBuilderParam distgo.DockerBuilderParam,
	productTaskOutputInfo distgo.ProductTaskOutputInfo,
	buildArtifactPaths map[distgo.ProductID]map[osarch.OSArch]string,
	distArtifactPaths map[distgo.ProductID]map[distgo.DistID][]string,
	verbose, dryRun bool,
	stdout io.Writer) (rErr error) {

	pathToContextDir := path.Join(projectInfo.ProjectDir, dockerBuilderParam.ContextDir)
	dockerfilePath := path.Join(pathToContextDir, dockerBuilderParam.DockerfilePath)
	originalDockerfileBytes, err := ioutil.ReadFile(dockerfilePath)
	if err != nil {
		return errors.Wrapf(err, "failed to read Dockerfile %s", dockerBuilderParam.DockerfilePath)
	}
	renderedDockerfile, err := distgo.RenderTemplate(string(originalDockerfileBytes), nil,
		distgo.ProductTemplateFunction(productID),
		distgo.VersionTemplateFunction(projectInfo.Version),
		distgo.RepositoryTemplateFunction(productTaskOutputInfo.Product.DockerOutputInfos.Repository),
		inputBuildArtifactTemplateFunction(dockerID, pathToContextDir, buildArtifactPaths),
		inputDistArtifactsTemplateFunction(dockerID, pathToContextDir, distArtifactPaths),
		tagsTemplateFunction(productTaskOutputInfo),
	)
	if err != nil {
		return err
	}

	if !dryRun {
		if renderedDockerfile != string(originalDockerfileBytes) {
			// Dockerfile contained templates and rendering them changes file: overwrite file with rendered version
			// and restore afterwards
			if err := ioutil.WriteFile(dockerfilePath, []byte(renderedDockerfile), 0644); err != nil {
				return errors.Wrapf(err, "failed to write rendered Dockerfile")
			}
			defer func() {
				if err := ioutil.WriteFile(dockerfilePath, originalDockerfileBytes, 0644); err != nil && rErr == nil {
					rErr = errors.Wrapf(err, "failed to restore original Dockerfile content")
				}
			}()
		}

		// link build artifacts into context directory
		for productID, valMap := range buildArtifactPaths {
			currOutputInfo := productTaskOutputInfo.AllProductOutputInfosMap()[productID]
			buildArtifactSrcPaths := distgo.ProductBuildArtifactPaths(projectInfo, currOutputInfo)
			for osArch, buildArtifactDstPath := range valMap {
				if err := os.MkdirAll(path.Dir(buildArtifactDstPath), 0755); err != nil {
					return errors.Wrapf(err, "failed to create directories")
				}
				if err := createNewHardLink(buildArtifactSrcPaths[osArch], buildArtifactDstPath); err != nil {
					return errors.Wrapf(err, "failed to link build artifact into context directory")
				}
			}
		}

		// link dist artifacts into context directory
		for productID, valMap := range distArtifactPaths {
			currOutputInfo := productTaskOutputInfo.AllProductOutputInfosMap()[productID]
			dstArtifactSrcPaths := distgo.ProductDistArtifactPaths(projectInfo, currOutputInfo)
			for distID, distArtifactDstPaths := range valMap {
				for i, currDstArtifactPath := range distArtifactDstPaths {
					if err := os.MkdirAll(path.Dir(currDstArtifactPath), 0755); err != nil {
						return errors.Wrapf(err, "failed to create directories")
					}
					if err := createNewHardLink(dstArtifactSrcPaths[distID][i], currDstArtifactPath); err != nil {
						return errors.Wrapf(err, "failed to link dist artifact into context directory")
					}
				}
			}
		}
	}

	distgo.PrintlnOrDryRunPrintln(stdout, fmt.Sprintf("Running Docker build for configuration %s of product %s...", dockerID, productID), dryRun)
	// run the Docker build task
	if err := dockerBuilderParam.DockerBuilder.RunDockerBuild(dockerID, productTaskOutputInfo, verbose, dryRun, stdout); err != nil {
		return err
	}
	return nil
}

func inputBuildArtifactTemplateFunction(dockerID distgo.DockerID, pathToContextDir string, buildArtifactPaths map[distgo.ProductID]map[osarch.OSArch]string) distgo.TemplateFunction {
	return func(fnMap template.FuncMap) {
		fnMap["InputBuildArtifact"] = func(productID, osArchStr string) (string, error) {
			osArchMap, ok := buildArtifactPaths[distgo.ProductID(productID)]
			if !ok {
				return "", errors.Errorf("product %s is not a build input for Docker task %s", productID, dockerID)
			}
			osArch, err := osarch.New(osArchStr)
			if err != nil {
				return "", errors.Wrapf(err, "input %s is not a valid OS/Arch", osArchStr)
			}
			dst, ok := osArchMap[osArch]
			if !ok {
				return "", errors.Errorf("OS/Arch %s for product %s is not defined as a build input for Docker task %s", osArchStr, productID, dockerID)
			}
			pathFromContextDir, err := filepath.Rel(pathToContextDir, dst)
			if err != nil {
				return "", errors.Wrapf(err, "failed to determine path")
			}
			return pathFromContextDir, nil
		}
	}
}

func inputDistArtifactsTemplateFunction(dockerID distgo.DockerID, pathToContextDir string, distArtifactPaths map[distgo.ProductID]map[distgo.DistID][]string) distgo.TemplateFunction {
	return func(fnMap template.FuncMap) {
		fnMap["InputDistArtifacts"] = func(productID, distID string) ([]string, error) {
			distIDsMap, ok := distArtifactPaths[distgo.ProductID(productID)]
			if !ok {
				return nil, errors.Errorf("product %s is not a dist input for Docker task %s", productID, dockerID)
			}
			dstArtifactPaths, ok := distIDsMap[distgo.DistID(distID)]
			if !ok {
				return nil, errors.Errorf("dist %s is not defined as a dist input for Docker task %s", distID, dockerID)
			}

			var outPaths []string
			for _, artifactPath := range dstArtifactPaths {
				pathFromContextDir, err := filepath.Rel(pathToContextDir, artifactPath)
				if err != nil {
					return nil, errors.Wrapf(err, "failed to determine path")
				}
				outPaths = append(outPaths, pathFromContextDir)
			}
			return outPaths, nil
		}
	}
}

func tagsTemplateFunction(productTaskOutputInfo distgo.ProductTaskOutputInfo) distgo.TemplateFunction {
	allOutputInfos := productTaskOutputInfo.AllProductOutputInfosMap()
	return func(fnMap template.FuncMap) {
		fnMap["Tags"] = func(productID, dockerID string) ([]string, error) {
			productOutputInfo, ok := allOutputInfos[distgo.ProductID(productID)]
			if !ok {
				return nil, errors.Errorf("product %s is not the product or a dependent product of %s", productID, productTaskOutputInfo.Product.ID)
			}
			if productOutputInfo.DockerOutputInfos == nil {
				return nil, errors.Errorf("product %s does not declare Docker outputs", productID)
			}
			dockerBuilderOutput, ok := productOutputInfo.DockerOutputInfos.DockerBuilderOutputInfos[distgo.DockerID(dockerID)]
			if !ok {
				return nil, errors.Errorf("product %s does not contain an entry for DockerID %s", productID, dockerID)
			}
			return dockerBuilderOutput.RenderedTags, nil
		}
	}
}

func createNewHardLink(src, dst string) error {
	if _, err := os.Stat(dst); err == nil {
		// ensure the target does not exists before creating a new one
		if err := os.Remove(dst); err != nil {
			return errors.Wrapf(err, "failed to remove existing file")
		}
	}
	if err := os.Link(src, dst); err != nil {
		return errors.Wrapf(err, "failed to create hard link %s from %s", dst, src)
	}
	return nil
}
