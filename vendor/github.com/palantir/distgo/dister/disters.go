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

package dister

import (
	"sort"

	"github.com/pkg/errors"

	"github.com/palantir/distgo/distgo"
)

type CreatorFunction func(cfgYML []byte) (distgo.Dister, error)

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

type disterFactory struct {
	disterCreators map[string]CreatorFunction
}

func (f *disterFactory) NewDister(typeName string, cfgYMLBytes []byte) (distgo.Dister, error) {
	creatorFn, ok := f.disterCreators[typeName]
	if !ok {
		var disterNames []string
		for k := range f.disterCreators {
			disterNames = append(disterNames, k)
		}
		sort.Strings(disterNames)
		return nil, errors.Errorf("no dister registered for dister type %q (registered disters: %v)", typeName, disterNames)
	}
	return creatorFn(cfgYMLBytes)
}

func NewDisterFactory(providedDisterCreators ...Creator) (distgo.DisterFactory, error) {
	disterCreators := make(map[string]CreatorFunction)
	for k, v := range builtinDisters() {
		disterCreators[k] = v
	}
	for _, currCreator := range providedDisterCreators {
		disterCreators[currCreator.TypeName()] = currCreator.Creator()
	}
	return &disterFactory{
		disterCreators: disterCreators,
	}, nil
}

func builtinDisters() map[string]CreatorFunction {
	return map[string]CreatorFunction{
		OSArchBinDistTypeName: NewOSArchBinDisterFromConfig,
		ManualDistTypeName:    NewManualDisterFromConfig,
	}
}

func AssetDisterCreators(assetPaths ...string) ([]Creator, error) {
	var disterCreators []Creator
	disterNameToAssets := make(map[string][]string)
	for _, currAssetPath := range assetPaths {
		currDister := assetDister{
			assetPath: currAssetPath,
		}
		disterName, err := currDister.TypeName()
		if err != nil {
			return nil, errors.Wrapf(err, "failed to determine dister type name for asset %s", currAssetPath)
		}
		disterNameToAssets[disterName] = append(disterNameToAssets[disterName], currAssetPath)
		disterCreators = append(disterCreators, NewCreator(disterName,
			func(cfgYML []byte) (distgo.Dister, error) {
				currDister.cfgYML = string(cfgYML)
				if err := currDister.VerifyConfig(); err != nil {
					return nil, err
				}
				return &currDister, nil
			}))
	}
	var sortedKeys []string
	for k := range disterNameToAssets {
		sortedKeys = append(sortedKeys, k)
	}
	sort.Strings(sortedKeys)
	for _, k := range sortedKeys {
		if len(disterNameToAssets[k]) <= 1 {
			continue
		}
		sort.Strings(disterNameToAssets[k])
		return nil, errors.Errorf("dister type %s provided by multiple assets: %v", k, disterNameToAssets[k])
	}
	return disterCreators, nil
}
