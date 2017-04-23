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
	"path/filepath"

	"github.com/palantir/godel/apps/distgo/cmd"
	"github.com/palantir/godel/apps/distgo/cmd/build"
	"github.com/palantir/godel/apps/distgo/cmd/dist"
	"github.com/palantir/godel/apps/distgo/params"
	"github.com/palantir/godel/apps/distgo/pkg/osarch"
)

// DistArtifacts returns a map from product name to OrderedStringMap, where the values of the OrderedStringMap contains
// the mapping from the DistType to the path for the artifact for that type.
func DistArtifacts(buildSpecsWithDeps []params.ProductBuildSpecWithDeps, absPath bool) (map[string]OrderedStringMap, error) {
	return artifacts(buildSpecsWithDeps, func(spec params.ProductBuildSpec) buildSpecWithPaths {
		distTypeToPathMap := newOrderedStringMap()

		for _, currDistCfg := range spec.Dist {
			distTypeToPathMap.Put(string(currDistCfg.Info.Type()), dist.FullArtifactPath(currDistCfg.Info.Type(), spec, currDistCfg))
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
func BuildArtifacts(buildSpecsWithDeps []params.ProductBuildSpecWithDeps, buildParams BuildArtifactsParams) (map[string]OrderedStringMap, error) {
	artifacts, err := artifacts(buildSpecsWithDeps, func(spec params.ProductBuildSpec) buildSpecWithPaths {
		osArchToPathMap := newOrderedStringMap()
		buildPaths := build.ArtifactPaths(spec)

		for _, osArch := range spec.Build.OSArchs {
			if v, ok := buildPaths[osArch]; ok && buildParams.OSArchs.Matches(osArch) {
				osArchToPathMap.Put(osArch.String(), v)
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

func artifacts(buildSpecsWithDeps []params.ProductBuildSpecWithDeps, f artifactPathsFunc, absPath bool) (map[string]OrderedStringMap, error) {
	artifacts := make(map[string]OrderedStringMap)
	for _, currBuildSpecWithDeps := range buildSpecsWithDeps {
		specWithPaths := f(currBuildSpecWithDeps.Spec)
		for _, k := range specWithPaths.paths.Keys() {
			if !absPath {
				absPathValue, err := filepath.Rel(specWithPaths.spec.ProjectDir, specWithPaths.paths.Get(k))
				if err != nil {
					return nil, err
				}
				specWithPaths.paths.Put(k, absPathValue)
			}
		}
		if len(specWithPaths.paths.Keys()) > 0 {
			artifacts[specWithPaths.spec.ProductName] = specWithPaths.paths
		}
	}
	return artifacts, nil
}

type buildSpecWithPaths struct {
	paths OrderedStringMap
	spec  *params.ProductBuildSpec
}

type artifactPathsFunc func(spec params.ProductBuildSpec) buildSpecWithPaths

// OrderedStringMap represents an ordered map with strings as the keys and values.
type OrderedStringMap interface {
	Keys() []string
	Get(k string) string
	Put(k string, v string)
	Remove(k string) bool
}

type orderedStringMap struct {
	innerMap    map[string]string
	orderedKeys []string
}

func newOrderedStringMap() OrderedStringMap {
	return &orderedStringMap{
		innerMap: make(map[string]string),
	}
}

func (m *orderedStringMap) Get(k string) string {
	return m.innerMap[k]
}

func (m *orderedStringMap) Put(k string, v string) {
	_, keyPreviouslyInMap := m.innerMap[k]
	m.innerMap[k] = v
	if !keyPreviouslyInMap {
		m.orderedKeys = append(m.orderedKeys, k)
	}
}

func (m *orderedStringMap) Remove(k string) bool {
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

func (m *orderedStringMap) Keys() []string {
	return m.orderedKeys
}
