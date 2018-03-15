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

package artifacts

import (
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"sort"

	"github.com/palantir/godel/pkg/osarch"
	"github.com/pkg/errors"

	"github.com/palantir/distgo/distgo"
	"github.com/palantir/distgo/distgo/build"
)

func PrintBuildArtifacts(projectInfo distgo.ProjectInfo, projectParam distgo.ProjectParam, productBuildIDs []distgo.ProductBuildID, absPath, requiresBuild bool, stdout io.Writer) error {
	productParams, err := distgo.ProductParamsForBuildProductArgs(projectParam.Products, productBuildIDs...)
	if err != nil {
		return err
	}
	artifacts, err := Build(projectInfo, productParams, requiresBuild)
	if err != nil {
		return err
	}
	return printArtifacts(artifacts, &printArtifactOptions{
		projectDir: projectInfo.ProjectDir,
		wantAbs:    absPath,
	}, stdout)
}

// Build returns a map from product name to build artifact paths. If requiresBuild is true, only returns the artifacts
// that need to be built.
func Build(projectInfo distgo.ProjectInfo, productParams []distgo.ProductParam, requiresBuild bool) (map[distgo.ProductID][]string, error) {
	outputPaths := make(map[distgo.ProductID]map[osarch.OSArch]string)
	for _, currProductParam := range productParams {
		if requiresBuild {
			requiresBuildParam, err := build.RequiresBuild(projectInfo, currProductParam)
			if err != nil {
				return nil, errors.Wrapf(err, "failed to determine if product %s needs to be built", currProductParam.ID)
			}
			if requiresBuildParam == nil {
				continue
			}
			currProductParam = *requiresBuildParam
		}
		outputInfo, err := distgo.ToProductTaskOutputInfo(projectInfo, currProductParam)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to compute output info for %s", currProductParam.ID)
		}
		currOutputPaths := outputInfo.ProductBuildArtifactPaths()
		if len(currOutputPaths) == 0 {
			continue
		}
		outputPaths[currProductParam.ID] = currOutputPaths
	}

	buildArtifacts := make(map[distgo.ProductID][]string)
	for k, v := range outputPaths {
		for _, currPath := range v {
			buildArtifacts[k] = append(buildArtifacts[k], currPath)
		}
	}
	for _, v := range buildArtifacts {
		sort.Strings(v)
	}
	return buildArtifacts, nil
}

func PrintDistArtifacts(projectInfo distgo.ProjectInfo, projectParam distgo.ProjectParam, productDistIDs []distgo.ProductDistID, absPath bool, stdout io.Writer) error {
	productParams, err := distgo.ProductParamsForDistProductArgs(projectParam.Products, productDistIDs...)
	if err != nil {
		return err
	}
	artifacts, err := Dist(projectInfo, productParams)
	if err != nil {
		return err
	}
	return printArtifacts(artifacts, &printArtifactOptions{
		projectDir: projectInfo.ProjectDir,
		wantAbs:    absPath,
	}, stdout)
}

// Dist returns a map from ProductID all of the dist artifact paths for the product.
func Dist(projectInfo distgo.ProjectInfo, productParams []distgo.ProductParam) (map[distgo.ProductID][]string, error) {
	outputPaths := make(map[distgo.ProductID][]string)
	for _, currProductParam := range productParams {
		currOutputInfo, err := distgo.ToProductTaskOutputInfo(projectInfo, currProductParam)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to compute output info for %s", currProductParam.ID)
		}
		var currPaths []string
		for _, currDistArtifactPaths := range currOutputInfo.ProductDistArtifactPaths() {
			currPaths = append(currPaths, currDistArtifactPaths...)
		}
		sort.Strings(currPaths)
		outputPaths[currProductParam.ID] = currPaths
	}
	return outputPaths, nil
}

func PrintDockerArtifacts(projectInfo distgo.ProjectInfo, projectParam distgo.ProjectParam, productDockerIDs []distgo.ProductDockerID, stdout io.Writer) error {
	productParams, err := distgo.ProductParamsForDockerProductArgs(projectParam.Products, productDockerIDs...)
	if err != nil {
		return err
	}
	artifacts, err := Docker(projectInfo, productParams)
	if err != nil {
		return err
	}
	return printArtifacts(artifacts, nil, stdout)
}

// Docker returns a map from ProductID all of the tags for the Docker images for the product.
func Docker(projectInfo distgo.ProjectInfo, productParams []distgo.ProductParam) (map[distgo.ProductID][]string, error) {
	outputPaths := make(map[distgo.ProductID][]string)
	for _, currProductParam := range productParams {
		currOutputInfo, err := distgo.ToProductTaskOutputInfo(projectInfo, currProductParam)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to compute output info for %s", currProductParam.ID)
		}
		currDockerOutputInfos := currOutputInfo.Product.DockerOutputInfos.DockerBuilderOutputInfos
		var dockerIDs []distgo.DockerID
		for k := range currDockerOutputInfos {
			dockerIDs = append(dockerIDs, k)
		}
		sort.Sort(distgo.ByDockerID(dockerIDs))
		for _, dockerID := range dockerIDs {
			outputPaths[currProductParam.ID] = append(outputPaths[currProductParam.ID], currDockerOutputInfos[dockerID].RenderedTags...)
		}
	}
	return outputPaths, nil
}

func printArtifacts(artifacts map[distgo.ProductID][]string, opts *printArtifactOptions, stdout io.Writer) error {
	var wd string
	var outputs []string
	for _, productToOutputs := range artifacts {
		for _, currPath := range productToOutputs {
			if opts != nil && filepath.IsAbs(currPath) != opts.wantAbs {
				if !filepath.IsAbs(currPath) {
					if wd == "" {
						gotWd, err := os.Getwd()
						if err != nil {
							return errors.Wrapf(err, "failed to determine working directory")
						}
						wd = gotWd
					}
					currPath = path.Join(wd, currPath)
				} else {
					relPath, err := filepath.Rel(opts.projectDir, currPath)
					if err != nil {
						return errors.Wrapf(err, "failed to convert path to relative path")
					}
					currPath = relPath
				}
			}
			outputs = append(outputs, currPath)
		}
	}
	sort.Strings(outputs)

	for _, output := range outputs {
		fmt.Fprintln(stdout, output)
	}
	return nil
}

type printArtifactOptions struct {
	projectDir string
	wantAbs    bool
}
