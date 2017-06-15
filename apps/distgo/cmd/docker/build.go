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
	"os"
	"path"

	"github.com/pkg/errors"

	"github.com/palantir/godel/apps/distgo/cmd/build"
	"github.com/palantir/godel/apps/distgo/cmd/dist"
	"github.com/palantir/godel/apps/distgo/params"
)

func Build(cfg params.Project, wd string, baseRepo string, stdout io.Writer) error {
	// the docker build tasks first runs dist task on the products
	// on which the docker images have a dependency. after building the dists,
	// the images are built in ordered way since the images can have dependencies among themselves.
	productsToDist := make(map[string]struct{})
	productsToBuildImage := make(map[string]struct{})
	for productName, productSpec := range cfg.Products {
		if len(productSpec.DockerImages) > 0 {
			productsToBuildImage[productName] = struct{}{}
		}
		for _, image := range productSpec.DockerImages {
			for _, dep := range image.Dependencies() {
				if isDist(dep.Type) {
					productsToDist[dep.Product] = struct{}{}
				}
			}
		}
	}

	// run the dist task
	if err := dist.Products(setToSlice(productsToDist), cfg, false, wd, stdout); err != nil {
		return err
	}

	// build docker images
	buildSpecsWithDeps, err := build.SpecsWithDepsForArgs(cfg, setToSlice(productsToBuildImage), wd)
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
				repo, _ := orderedSpecs[i].Spec.DockerImages[j].Coordinates()
				orderedSpecs[i].Spec.DockerImages[j].SetRepository(path.Join(baseRepo, repo))
			}
		}
	}
	return RunBuild(orderedSpecs, stdout)
}

func RunBuild(buildSpecsWithDeps []params.ProductBuildSpecWithDeps, stdout io.Writer) error {
	specsMap := buildSpecsMap(buildSpecsWithDeps)
	for i := range buildSpecsWithDeps {
		for _, image := range buildSpecsWithDeps[i].Spec.DockerImages {
			if err := buildImage(image, buildSpecsWithDeps[i], specsMap, stdout); err != nil {
				return err
			}
		}
	}
	return nil
}

func buildImage(image params.DockerImage, buildSpecsWithDeps params.ProductBuildSpecWithDeps, specsMap map[string]params.ProductBuildSpecWithDeps, stdout io.Writer) error {
	repo, tag := image.Coordinates()
	fmt.Fprintf(stdout, "Building docker image for %s and tagging it as %s:%s\n", buildSpecsWithDeps.Spec.ProductName, repo, tag)

	contextDir := path.Join(buildSpecsWithDeps.Spec.ProjectDir, image.ContextDirectory())

	// link dependent dist artifacts into the context directory
	for depProduct, depTypes := range dockerDepsToMap(image.Dependencies()) {
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

	return image.Build(buildSpecsWithDeps)
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

func isDist(dep params.DockerDepType) bool {
	switch dep {
	case params.DockerDepSLS, params.DockerDepBin, params.DockerDepRPM:
		return true
	default:
		return false
	}
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

func setToSlice(s map[string]struct{}) []string {
	var result []string
	for item := range s {
		result = append(result, item)
	}
	return result
}
