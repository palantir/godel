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

	"github.com/pkg/errors"

	"github.com/palantir/godel/apps/distgo/cmd/build"
	"github.com/palantir/godel/apps/distgo/cmd/dist"
	"github.com/palantir/godel/apps/distgo/params"
)

func Build(products []string, cfg params.Project, wd, baseRepo string, verbose bool, stdout io.Writer) error {
	// the docker build tasks first runs dist task on the products
	// on which the docker images have a dependency. after building the dists,
	// the images are built in ordered way since the images can have dependencies among themselves.
	productsToDist, productsToBuildImage, err := productsToDistAndBuildImage(products, cfg)
	if err != nil {
		return err
	}

	productsToDistRequired, err := dist.RequiresDist(productsToDist, cfg, wd)
	if err != nil {
		return err
	}
	if len(productsToDistRequired) != 0 {
		// run the dist task
		if err := dist.Products(productsToDistRequired, cfg, false, wd, stdout); err != nil {
			return err
		}
	}

	// build docker images
	buildSpecsWithDeps, err := build.SpecsWithDepsForArgs(cfg, productsToBuildImage, wd)
	if err != nil {
		return err
	}
	orderedSpecs, err := OrderBuildSpecs(buildSpecsWithDeps)
	if err != nil {
		return err
	}
	if baseRepo != "" {
		// if base repo is specified, join it to each image's repo
		for i := range orderedSpecs {
			for j := range orderedSpecs[i].Spec.DockerImages {
				orderedSpecs[i].Spec.DockerImages[j].Repository = path.Join(baseRepo,
					orderedSpecs[i].Spec.DockerImages[j].Repository)
			}
		}
	}
	return RunBuild(orderedSpecs, verbose, stdout)
}

func RunBuild(buildSpecsWithDeps []params.ProductBuildSpecWithDeps, verbose bool, stdout io.Writer) error {
	specsMap := buildSpecsMap(buildSpecsWithDeps)
	for i := range buildSpecsWithDeps {
		for _, image := range buildSpecsWithDeps[i].Spec.DockerImages {
			if err := buildImage(image, buildSpecsWithDeps[i], specsMap, verbose, stdout); err != nil {
				return err
			}
		}
	}
	return nil
}

func buildImage(image params.DockerImage, buildSpecsWithDeps params.ProductBuildSpecWithDeps, specsMap map[string]params.ProductBuildSpecWithDeps, verbose bool, stdout io.Writer) error {
	fmt.Fprintf(stdout, "Building docker image for %s and tagging it as %s:%s\n", buildSpecsWithDeps.Spec.ProductName, image.Repository, image.Tag)

	contextDir := path.Join(buildSpecsWithDeps.Spec.ProjectDir, image.ContextDir)

	// link dependent dist artifacts into the context directory
	for depProduct, depTypes := range dockerDepsToMap(image.Deps) {
		for depType, targetFile := range depTypes {
			if !isDist(depType) {
				continue
			}
			if _, ok := specsMap[depProduct]; !ok {
				return errors.Errorf("Unable to find the dependent product %v for %v",
					depProduct, buildSpecsWithDeps.Spec.ProductName)
			}
			depSpec := specsMap[depProduct].Spec
			distsMap := buildDistsMap(specsMap[depProduct].Spec)
			if _, ok := distsMap[string(depType)]; !ok {
				return errors.Errorf("Unable to find dist type %v on the dependent product %v for %v",
					depType, depProduct, buildSpecsWithDeps.Spec.ProductName)
			}
			distCfg := distsMap[string(depType)]

			for _, artifactLocation := range dist.FullArtifactsPaths(dist.ToDister(distCfg.Info), depSpec, distCfg) {
				if targetFile == "" {
					targetFile = path.Base(artifactLocation)
				}
				target := path.Join(contextDir, targetFile)
				if _, err := os.Stat(target); err == nil {
					// ensure the target does not exists before creating a new one
					if err := os.Remove(target); err != nil {
						return err
					}
				}
				if err := os.Link(artifactLocation, target); err != nil {
					return err
				}
			}
		}
	}
	builder := GetBuilder(image)
	buildWriter := ioutil.Discard
	if verbose {
		buildWriter = stdout
	}
	return builder.build(buildSpecsWithDeps, buildWriter)
}

func dockerDepsToMap(deps []params.DockerDep) map[string]map[params.DockerDepType]string {
	m := make(map[string]map[params.DockerDepType]string)
	for _, dep := range deps {
		if m[dep.Product] == nil {
			m[dep.Product] = make(map[params.DockerDepType]string)
		}
		m[dep.Product][dep.Type] = dep.TargetFile
	}
	return m
}

func buildSpecsMap(specs []params.ProductBuildSpecWithDeps) map[string]params.ProductBuildSpecWithDeps {
	specMap := make(map[string]params.ProductBuildSpecWithDeps)
	for _, spec := range specs {
		specMap[spec.Spec.ProductName] = spec
	}
	return specMap
}

func buildDistsMap(spec params.ProductBuildSpec) map[string]params.Dist {
	distMap := make(map[string]params.Dist)
	for _, dist := range spec.Dist {
		distMap[string(dist.Info.Type())] = dist
	}
	return distMap
}
