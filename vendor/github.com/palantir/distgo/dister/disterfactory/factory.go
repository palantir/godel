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

package disterfactory

import (
	"github.com/pkg/errors"

	"github.com/palantir/distgo/dister"
	"github.com/palantir/distgo/distgo"
)

func New(providedDisterCreators []dister.Creator, providedConfigUpgraders []distgo.ConfigUpgrader) (distgo.DisterFactory, error) {
	var types []string
	seenTypes := make(map[string]struct{})
	disterCreators := make(map[string]dister.CreatorFunction)
	configUpgraders := make(map[string]distgo.ConfigUpgrader)
	for k, v := range builtinDisters() {
		types = append(types, k)
		seenTypes[k] = struct{}{}
		disterCreators[k] = v.creator
		configUpgraders[k] = v.upgrader
	}
	for _, currCreator := range providedDisterCreators {
		if _, ok := seenTypes[currCreator.TypeName()]; ok {
			return nil, errors.Errorf("dister creator with type %q specified more than once", currCreator.TypeName())
		}
		seenTypes[currCreator.TypeName()] = struct{}{}
		types = append(types, currCreator.TypeName())
		disterCreators[currCreator.TypeName()] = currCreator.Creator()
	}
	for _, currUpgrader := range providedConfigUpgraders {
		currUpgrader := currUpgrader
		configUpgraders[currUpgrader.TypeName()] = currUpgrader
	}
	return &disterFactoryImpl{
		types:                 types,
		disterCreators:        disterCreators,
		disterConfigUpgraders: configUpgraders,
	}, nil
}

type disterFactoryImpl struct {
	types                 []string
	disterCreators        map[string]dister.CreatorFunction
	disterConfigUpgraders map[string]distgo.ConfigUpgrader
}

func (f *disterFactoryImpl) Types() []string {
	return f.types
}

func (f *disterFactoryImpl) NewDister(typeName string, cfgYMLBytes []byte) (distgo.Dister, error) {
	creatorFn, ok := f.disterCreators[typeName]
	if !ok {
		return nil, errors.Errorf("no dister registered for dister type %q (registered disters: %v)", typeName, f.types)
	}
	return creatorFn(cfgYMLBytes)
}

func (f *disterFactoryImpl) ConfigUpgrader(typeName string) (distgo.ConfigUpgrader, error) {
	if _, ok := f.disterCreators[typeName]; !ok {
		return nil, errors.Errorf("no disters registered for dister type %q (registered disters: %v)", typeName, f.types)
	}
	upgrader, ok := f.disterConfigUpgraders[typeName]
	if !ok {
		return nil, errors.Errorf("%s is a valid dister but does not have a config upgrader", typeName)
	}
	return upgrader, nil
}
