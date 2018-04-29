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

package config

import (
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"

	"github.com/palantir/distgo/distgo"
	"github.com/palantir/distgo/distgo/config/internal/v0"
)

type PublishConfig v0.PublishConfig

func ToPublishConfig(in *PublishConfig) *v0.PublishConfig {
	return (*v0.PublishConfig)(in)
}

func (cfg *PublishConfig) ToParam(defaultCfg PublishConfig) (distgo.PublishParam, error) {
	publishInfo, err := mergePublishInfos(fromPublishIDPublisherConfig(cfg.PublishInfo), fromPublishIDPublisherConfig(defaultCfg.PublishInfo))
	if err != nil {
		return distgo.PublishParam{}, err
	}
	return distgo.PublishParam{
		GroupID:     getConfigStringValue(cfg.GroupID, defaultCfg.GroupID, ""),
		PublishInfo: publishInfo,
	}, nil
}

type PublisherConfig v0.PublisherConfig

func ToPublisherConfig(in *PublisherConfig) *v0.PublisherConfig {
	return (*v0.PublisherConfig)(in)
}

func (cfg *PublisherConfig) ToParam() (distgo.PublisherParam, error) {
	var cfgBytes []byte
	if cfg.Config != nil {
		bytes, err := yaml.Marshal(cfg.Config)
		if err != nil {
			return distgo.PublisherParam{}, errors.Wrapf(err, "failed to marshal publisher configuration")
		}
		cfgBytes = bytes
	}

	return distgo.PublisherParam{
		ConfigBytes: cfgBytes,
	}, nil
}

func mergePublishInfos(cfgPublishInfo, defaultCfgPublishInfo *map[distgo.PublisherTypeID]PublisherConfig) (map[distgo.PublisherTypeID]distgo.PublisherParam, error) {
	distinctVals := make(map[distgo.PublisherTypeID]PublisherConfig)
	commonKeys := make(map[distgo.PublisherTypeID]struct{})

	if cfgPublishInfo != nil {
		for publishID, val := range *cfgPublishInfo {
			if defaultCfgPublishInfo != nil {
				if _, ok := (*defaultCfgPublishInfo)[publishID]; ok {
					commonKeys[publishID] = struct{}{}
					continue
				}
			}
			distinctVals[publishID] = val
		}
	}

	if defaultCfgPublishInfo != nil {
		for publishID, val := range *defaultCfgPublishInfo {
			if cfgPublishInfo != nil {
				if _, ok := (*cfgPublishInfo)[publishID]; ok {
					commonKeys[publishID] = struct{}{}
					continue
				}
			}
			distinctVals[publishID] = val
		}
	}

	var publishInfo map[distgo.PublisherTypeID]distgo.PublisherParam
	if len(distinctVals) > 0 || len(commonKeys) > 0 {
		publishInfo = make(map[distgo.PublisherTypeID]distgo.PublisherParam)
		for publishID, publisherConfig := range distinctVals {
			publisherParam, err := publisherConfig.ToParam()
			if err != nil {
				return nil, errors.Wrapf(err, "failed to create publisher param for %s", publishID)
			}
			publishInfo[publishID] = publisherParam
		}
		for publishID := range commonKeys {
			publisherCfg := (*cfgPublishInfo)[publishID]
			publisherParam, err := publisherCfg.ToParam()
			if err != nil {
				return nil, errors.Wrapf(err, "failed to create publisher param for %s", publishID)
			}
			publishInfo[publishID] = publisherParam
		}
	}
	return publishInfo, nil
}

func ToPublishInfo(in *map[distgo.PublisherTypeID]PublisherConfig) *map[distgo.PublisherTypeID]v0.PublisherConfig {
	if in == nil {
		return nil
	}
	out := make(map[distgo.PublisherTypeID]v0.PublisherConfig, len(*in))
	for k, v := range *in {
		out[k] = v0.PublisherConfig(v)
	}
	return &out
}

func fromPublishIDPublisherConfig(in *map[distgo.PublisherTypeID]v0.PublisherConfig) *map[distgo.PublisherTypeID]PublisherConfig {
	if in == nil {
		return nil
	}
	out := make(map[distgo.PublisherTypeID]PublisherConfig, len(*in))
	for k, v := range *in {
		out[k] = PublisherConfig(v)
	}
	return &out
}
