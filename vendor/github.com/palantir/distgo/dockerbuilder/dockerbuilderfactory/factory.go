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

package dockerbuilderfactory

import (
	"github.com/pkg/errors"

	"github.com/palantir/distgo/distgo"
	"github.com/palantir/distgo/dockerbuilder"
)

func New(providedDockerBuilderCreators []dockerbuilder.Creator, providedConfigUpgraders []distgo.ConfigUpgrader) (distgo.DockerBuilderFactory, error) {
	var types []string
	seenTypes := make(map[string]struct{})
	dockerBuilderCreators := make(map[string]dockerbuilder.CreatorFunction)
	configUpgraders := make(map[string]distgo.ConfigUpgrader)
	for k, v := range builtinDockerBuilders() {
		types = append(types, k)
		seenTypes[k] = struct{}{}
		dockerBuilderCreators[k] = v.creator
		configUpgraders[k] = v.upgrader
	}
	for _, currCreator := range providedDockerBuilderCreators {
		if _, ok := seenTypes[currCreator.TypeName()]; ok {
			return nil, errors.Errorf("docker builder creator with type %q specified more than once", currCreator.TypeName())
		}
		seenTypes[currCreator.TypeName()] = struct{}{}
		types = append(types, currCreator.TypeName())
		dockerBuilderCreators[currCreator.TypeName()] = currCreator.Creator()
	}
	for _, currUpgrader := range providedConfigUpgraders {
		currUpgrader := currUpgrader
		configUpgraders[currUpgrader.TypeName()] = currUpgrader
	}
	return &dockerBuilderFactory{
		types: types,
		dockerBuilderCreators:        dockerBuilderCreators,
		dockerBuilderConfigUpgraders: configUpgraders,
	}, nil
}

type dockerBuilderFactory struct {
	types                        []string
	dockerBuilderCreators        map[string]dockerbuilder.CreatorFunction
	dockerBuilderConfigUpgraders map[string]distgo.ConfigUpgrader
}

func (f *dockerBuilderFactory) Types() []string {
	return f.types
}

func (f *dockerBuilderFactory) NewDockerBuilder(typeName string, cfgYMLBytes []byte) (distgo.DockerBuilder, error) {
	creatorFn, ok := f.dockerBuilderCreators[typeName]
	if !ok {
		return nil, errors.Errorf("no DockerBuilder registered for DockerBuilder type %q (registered disters: %v)", typeName, f.types)
	}
	return creatorFn(cfgYMLBytes)
}

func (f *dockerBuilderFactory) ConfigUpgrader(typeName string) (distgo.ConfigUpgrader, error) {
	if _, ok := f.dockerBuilderCreators[typeName]; !ok {
		return nil, errors.Errorf("no docker builder registered for docker builder type %q (registered docker builders: %v)", typeName, f.types)
	}
	upgrader, ok := f.dockerBuilderConfigUpgraders[typeName]
	if !ok {
		return nil, errors.Errorf("%s is a valid docker builder but does not have a config upgrader", typeName)
	}
	return upgrader, nil
}
