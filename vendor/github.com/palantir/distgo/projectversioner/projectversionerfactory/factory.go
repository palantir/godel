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

package projectversionerfactory

import (
	"github.com/pkg/errors"

	"github.com/palantir/distgo/distgo"
	"github.com/palantir/distgo/projectversioner"
)

func New(providedProjectVersionerCreators []projectversioner.Creator, providedConfigUpgraders []distgo.ConfigUpgrader) (distgo.ProjectVersionerFactory, error) {
	var types []string
	seenTypes := make(map[string]struct{})
	projectVersionerCreators := make(map[string]projectversioner.CreatorFunction)
	configUpgraders := make(map[string]distgo.ConfigUpgrader)
	for k, v := range builtinProjectVersioners() {
		types = append(types, k)
		seenTypes[k] = struct{}{}
		projectVersionerCreators[k] = v.creator
		configUpgraders[k] = v.upgrader
	}
	for _, currCreator := range providedProjectVersionerCreators {
		if _, ok := seenTypes[currCreator.TypeName()]; ok {
			return nil, errors.Errorf("dister creator with type %q specified more than once", currCreator.TypeName())
		}
		seenTypes[currCreator.TypeName()] = struct{}{}
		types = append(types, currCreator.TypeName())
		projectVersionerCreators[currCreator.TypeName()] = currCreator.Creator()
	}
	for _, currUpgrader := range providedConfigUpgraders {
		currUpgrader := currUpgrader
		configUpgraders[currUpgrader.TypeName()] = currUpgrader
	}
	return &projectVersionerFactoryImpl{
		types: types,
		projectVersionerCreators:        projectVersionerCreators,
		projectVersionerConfigUpgraders: configUpgraders,
	}, nil
}

type projectVersionerFactoryImpl struct {
	types                           []string
	projectVersionerCreators        map[string]projectversioner.CreatorFunction
	projectVersionerConfigUpgraders map[string]distgo.ConfigUpgrader
}

func (f *projectVersionerFactoryImpl) Types() []string {
	return f.types
}

func (f *projectVersionerFactoryImpl) NewProjectVersioner(typeName string, cfgYMLBytes []byte) (distgo.ProjectVersioner, error) {
	creatorFn, ok := f.projectVersionerCreators[typeName]
	if !ok {
		return nil, errors.Errorf("no project versioner registered for project versioner type %q (registered project versioner(s): %v)", typeName, f.types)
	}
	return creatorFn(cfgYMLBytes)
}

func (f *projectVersionerFactoryImpl) ConfigUpgrader(typeName string) (distgo.ConfigUpgrader, error) {
	if _, ok := f.projectVersionerCreators[typeName]; !ok {
		return nil, errors.Errorf("no project versioner registered for project versioner type %q (registered project versioner(s): %v)", typeName, f.types)
	}
	upgrader, ok := f.projectVersionerConfigUpgraders[typeName]
	if !ok {
		return nil, errors.Errorf("%s is a valid project versioner but does not have a config upgrader", typeName)
	}
	return upgrader, nil
}
