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

package publisher

import (
	"sort"

	"github.com/pkg/errors"

	"github.com/palantir/distgo/distgo"
)

type Creator interface {
	TypeName() string
	Publisher() distgo.Publisher
}

type creatorStruct struct {
	typeName  string
	publisher func() distgo.Publisher
}

func (c *creatorStruct) TypeName() string {
	return c.typeName
}

func (c *creatorStruct) Publisher() distgo.Publisher {
	return c.publisher()
}

func NewCreator(typeName string, publisherCreator func() distgo.Publisher) Creator {
	return &creatorStruct{
		typeName:  typeName,
		publisher: publisherCreator,
	}
}

var registeredPublishers map[string]Creator

func init() {
	registeredPublishers = builtinPublishers()
}

func SetPublishers(publishers map[string]Creator) error {
	registeredPublishers = make(map[string]Creator)
	for k, v := range builtinPublishers() {
		registeredPublishers[k] = v
	}
	for k, v := range publishers {
		registeredPublishers[k] = v
	}
	return nil
}

func Publishers() map[string]distgo.Publisher {
	publishers := make(map[string]distgo.Publisher)
	for k, v := range registeredPublishers {
		publishers[k] = v.Publisher()
	}
	return publishers
}

func builtinPublishers() map[string]Creator {
	return map[string]Creator{
		LocalPublishTypeName:       NewLocalPublisherCreator(),
		ArtifactoryPublishTypeName: NewArtifactoryPublisherCreator(),
		BintrayPublishTypeName:     NewBintrayPublisherCreator(),
		GitHubPublishTypeName:      NewGitHubPublisherCreator(),
	}
}

func AssetPublisherCreators(assetPaths ...string) (map[string]Creator, error) {
	publishers := make(map[string]Creator)
	publisherNameToAssets := make(map[string][]string)
	for _, currAssetPath := range assetPaths {
		publisher := newAssetPublisher(currAssetPath)
		publisherName, err := publisher.TypeName()
		if err != nil {
			return nil, errors.Wrapf(err, "failed to determine type name for asset %s", currAssetPath)
		}
		publisherNameToAssets[publisherName] = append(publisherNameToAssets[publisherName], currAssetPath)
		publishers[publisherName] = NewCreator(publisherName, func() distgo.Publisher {
			return newAssetPublisher(currAssetPath)
		})
	}
	var sortedKeys []string
	for k := range publisherNameToAssets {
		sortedKeys = append(sortedKeys, k)
	}
	sort.Strings(sortedKeys)
	for _, k := range sortedKeys {
		if len(publisherNameToAssets[k]) <= 1 {
			continue
		}
		sort.Strings(publisherNameToAssets[k])
		return nil, errors.Errorf("publisher type %s provided by multiple assets: %v", k, publisherNameToAssets[k])
	}
	return publishers, nil
}
