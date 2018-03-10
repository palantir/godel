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

package dockerbuilder

import (
	"sort"

	"github.com/pkg/errors"

	"github.com/palantir/distgo/distgo"
)

type CreatorFunction func(cfgYML []byte) (distgo.DockerBuilder, error)

type Creator interface {
	TypeName() string
	Creator() CreatorFunction
}

type creatorStruct struct {
	typeName string
	creator  CreatorFunction
}

func (c *creatorStruct) TypeName() string {
	return c.typeName
}

func (c *creatorStruct) Creator() CreatorFunction {
	return c.creator
}

func NewCreator(typeName string, creatorFn CreatorFunction) Creator {
	return &creatorStruct{
		typeName: typeName,
		creator:  creatorFn,
	}
}

type dockerBuilderFactory struct {
	dockerBuilderCreators map[string]CreatorFunction
}

func (f *dockerBuilderFactory) NewDockerBuilder(typeName string, cfgYMLBytes []byte) (distgo.DockerBuilder, error) {
	creatorFn, ok := f.dockerBuilderCreators[typeName]
	if !ok {
		var dockerBuilderNames []string
		for k := range f.dockerBuilderCreators {
			dockerBuilderNames = append(dockerBuilderNames, k)
		}
		sort.Strings(dockerBuilderNames)
		return nil, errors.Errorf("no DockerBuilder registered for DockerBuilder type %q (registered disters: %v)", typeName, dockerBuilderNames)
	}
	return creatorFn(cfgYMLBytes)
}

func NewDockerBuilderFactory(providedDockerBuilderCreators ...Creator) (distgo.DockerBuilderFactory, error) {
	disterCreators := make(map[string]CreatorFunction)
	for k, v := range builtinDisters() {
		disterCreators[k] = v
	}
	for _, currCreator := range providedDockerBuilderCreators {
		disterCreators[currCreator.TypeName()] = currCreator.Creator()
	}
	return &dockerBuilderFactory{
		dockerBuilderCreators: disterCreators,
	}, nil
}

func builtinDisters() map[string]CreatorFunction {
	return map[string]CreatorFunction{
		DefaultBuilderTypeName: NewDefaultDockerBuilderFromConfig,
	}
}

func AssetDockerBuilderCreators(assetPaths ...string) ([]Creator, error) {
	var dockerBuilderCreators []Creator
	dockerBuilderNameToAssets := make(map[string][]string)
	for _, currAssetPath := range assetPaths {
		currDockerBuilder := assetDockerBuilder{
			assetPath: currAssetPath,
		}
		dockerBuilderName, err := currDockerBuilder.TypeName()
		if err != nil {
			return nil, errors.Wrapf(err, "failed to determine DockerBuilder type name for asset %s", currAssetPath)
		}
		dockerBuilderNameToAssets[dockerBuilderName] = append(dockerBuilderNameToAssets[dockerBuilderName], currAssetPath)
		dockerBuilderCreators = append(dockerBuilderCreators, NewCreator(dockerBuilderName,
			func(cfgYML []byte) (distgo.DockerBuilder, error) {
				currDockerBuilder.cfgYML = string(cfgYML)
				if err := currDockerBuilder.VerifyConfig(); err != nil {
					return nil, err
				}
				return &currDockerBuilder, nil
			}))
	}
	var sortedKeys []string
	for k := range dockerBuilderNameToAssets {
		sortedKeys = append(sortedKeys, k)
	}
	sort.Strings(sortedKeys)
	for _, k := range sortedKeys {
		if len(dockerBuilderNameToAssets[k]) <= 1 {
			continue
		}
		sort.Strings(dockerBuilderNameToAssets[k])
		return nil, errors.Errorf("DockerBuilder type %s provided by multiple assets: %v", k, dockerBuilderNameToAssets[k])
	}
	return dockerBuilderCreators, nil
}
