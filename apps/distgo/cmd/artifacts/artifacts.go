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
	"path/filepath"
	"strconv"

	"github.com/palantir/godel/apps/distgo/cmd"
	"github.com/palantir/godel/apps/distgo/cmd/build"
	"github.com/palantir/godel/apps/distgo/cmd/dist"
	"github.com/palantir/godel/apps/distgo/params"
	"github.com/palantir/godel/apps/distgo/pkg/osarch"
)

// DockerArtifacts returns a map from product name to a slice that contains all of the Docker repository:tag labels
// defined for the products.
func DockerArtifacts(buildSpecsWithDeps []params.ProductBuildSpecWithDeps) map[string][]string {
	output := make(map[string][]string)
	for _, spec := range buildSpecsWithDeps {
		for _, currImage := range spec.Spec.DockerImages {
			imageName := fmt.Sprint(currImage.Repository, ":", currImage.Tag)
			output[spec.Spec.ProductName] = append(output[spec.Spec.ProductName], imageName)
		}
	}
	return output
}

// DistArtifacts returns a map from product name to OrderedStringSliceMap, where the values of the OrderedStringSliceMap
// contains the mapping from the String representation of the index of the dist type ("0", "1", etc.) to the paths for
// the artifact for that type.
func DistArtifacts(buildSpecsWithDeps []params.ProductBuildSpecWithDeps, absPath bool) (map[string]OrderedStringSliceMap, error) {
	return artifacts(buildSpecsWithDeps, func(spec params.ProductBuildSpec) buildSpecWithPaths {
		distTypeToPathMap := newOrderedStringSliceMap()
		for i, currDistCfg := range spec.Dist {
			artifactPaths := dist.FullArtifactsPaths(dist.ToDister(currDistCfg.Info), spec, currDistCfg)
			distTypeToPathMap.PutValues(strconv.Itoa(i), artifactPaths)
		}
		return buildSpecWithPaths{spec: &spec, paths: distTypeToPathMap}
	}, absPath)
}

type BuildArtifactsParams struct {
	AbsPath       bool
	RequiresBuild bool
	OSArchs       cmd.OSArchFilter
}

// BuildArtifacts returns a map from product name to OrderedStringMap, where the values of the OrderedStringMap contains
// the mapping from the OSArch to the path for the artifact for that OSArch.
func BuildArtifacts(buildSpecsWithDeps []params.ProductBuildSpecWithDeps, buildParams BuildArtifactsParams) (map[string]OrderedStringSliceMap, error) {
	artifacts, err := artifacts(buildSpecsWithDeps, func(spec params.ProductBuildSpec) buildSpecWithPaths {
		osArchToPathMap := newOrderedStringSliceMap()
		buildPaths := build.ArtifactPaths(spec)

		for _, osArch := range spec.Build.OSArchs {
			if v, ok := buildPaths[osArch]; ok && buildParams.OSArchs.Matches(osArch) {
				osArchToPathMap.Add(osArch.String(), v)
			}
		}
		return buildSpecWithPaths{spec: &spec, paths: osArchToPathMap}
	}, buildParams.AbsPath)

	// if error occurred or requiresBuild is not true, return
	if err != nil || !buildParams.RequiresBuild {
		return artifacts, err
	}

	// otherwise, filter to artifacts that require build
	for _, spec := range buildSpecsWithDeps {
		requiresBuildInfo := build.RequiresBuild(spec, buildParams.OSArchs)
		for product := range artifacts {
			// copy keys over which iteration occurs because keys are removed during traversal
			src := artifacts[product].Keys()
			origKeys := make([]string, len(src))
			copy(origKeys, src)

			// remove any OSArch values that do not need to be built
			for _, osArchStr := range origKeys {
				osArch, err := osarch.New(osArchStr)
				if err != nil {
					return nil, err
				}
				if !requiresBuildInfo.RequiresBuild(product, osArch) {
					artifacts[product].Remove(osArchStr)
				}
			}
			// if product no longer has any OSArch values after filtering, remove it from the map
			if len(artifacts[product].Keys()) == 0 {
				delete(artifacts, product)
			}
		}
	}

	return artifacts, err
}

func artifacts(buildSpecsWithDeps []params.ProductBuildSpecWithDeps, f artifactPathsFunc, absPath bool) (map[string]OrderedStringSliceMap, error) {
	artifacts := make(map[string]OrderedStringSliceMap)
	for _, currBuildSpecWithDeps := range buildSpecsWithDeps {
		specWithPaths := f(currBuildSpecWithDeps.Spec)
		for _, k := range specWithPaths.paths.Keys() {
			if !absPath {
				var newVals []string
				for _, currPath := range specWithPaths.paths.Get(k) {
					relPathValue, err := filepath.Rel(specWithPaths.spec.ProjectDir, currPath)
					if err != nil {
						return nil, err
					}
					newVals = append(newVals, relPathValue)
				}
				specWithPaths.paths.PutValues(k, newVals)
			}
		}
		if len(specWithPaths.paths.Keys()) > 0 {
			artifacts[specWithPaths.spec.ProductName] = specWithPaths.paths
		}
	}
	return artifacts, nil
}

type buildSpecWithPaths struct {
	paths OrderedStringSliceMap
	spec  *params.ProductBuildSpec
}

type artifactPathsFunc func(spec params.ProductBuildSpec) buildSpecWithPaths

// OrderedStringSliceMap represents an ordered map with strings as the keys and a string slice as the values.
type OrderedStringSliceMap interface {
	Keys() []string
	Get(k string) []string
	Add(k, v string)
	PutValues(k string, v []string)
	Remove(k string) bool
}

type orderedStringSliceMap struct {
	innerMap    map[string][]string
	orderedKeys []string
}

func newOrderedStringSliceMap() OrderedStringSliceMap {
	return &orderedStringSliceMap{
		innerMap: make(map[string][]string),
	}
}

func (m *orderedStringSliceMap) Get(k string) []string {
	return m.innerMap[k]
}

func (m *orderedStringSliceMap) Add(k, v string) {
	prevVal, keyPreviouslyInMap := m.innerMap[k]
	m.innerMap[k] = append(prevVal, v)
	if !keyPreviouslyInMap {
		m.orderedKeys = append(m.orderedKeys, k)
	}
}

func (m *orderedStringSliceMap) PutValues(k string, v []string) {
	_, keyPreviouslyInMap := m.innerMap[k]
	m.innerMap[k] = v
	if !keyPreviouslyInMap {
		m.orderedKeys = append(m.orderedKeys, k)
	}
}

func (m *orderedStringSliceMap) Remove(k string) bool {
	if _, keyInMap := m.innerMap[k]; !keyInMap {
		return false
	}

	// remove key from map
	delete(m.innerMap, k)
	// remove key from slice
	for i, v := range m.orderedKeys {
		if v == k {
			m.orderedKeys = append(m.orderedKeys[:i], m.orderedKeys[i+1:]...)
			break
		}
	}
	return true
}

func (m *orderedStringSliceMap) Keys() []string {
	return m.orderedKeys
}
