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

package distgo

import (
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

type PublishID string

type PublishParam struct {
	// GroupID is the Maven group ID used for the publish operation.
	GroupID string

	// PublishInfo contains extra configuration for the publish operation. The key is the type of publish and the value
	// is the raw YAML configuration for that publish operation type.
	PublishInfo map[PublishID][]byte
}

type PublishOutputInfo struct {
	GroupID string `json:"groupId"`
}

func (p *PublishParam) ToPublishOutputInfo() PublishOutputInfo {
	return PublishOutputInfo{
		GroupID: p.GroupID,
	}
}

type PublishConfig struct {
	// GroupID is the product-specific configuration equivalent to the global GroupID configuration.
	GroupID *string `yaml:"group-id"`

	// PublishInfo contains extra configuration for the publish operation. The key is the type of publish and the value
	// is the configuration for that publish operation type.
	PublishInfo *map[PublishID]yaml.MapSlice `yaml:"info"`
}

func (cfg *PublishConfig) ToParam(defaultCfg PublishConfig) (PublishParam, error) {
	publishInfo, err := mergePublishInfos(cfg.PublishInfo, defaultCfg.PublishInfo)
	if err != nil {
		return PublishParam{}, err
	}
	return PublishParam{
		GroupID:     getConfigStringValue(cfg.GroupID, defaultCfg.GroupID, ""),
		PublishInfo: publishInfo,
	}, nil
}

func mergePublishInfos(cfgPublishInfo, defaultCfgPublishInfo *map[PublishID]yaml.MapSlice) (map[PublishID][]byte, error) {
	distinctVals := make(map[PublishID]yaml.MapSlice)
	commonKeys := make(map[PublishID]struct{})

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

	var publishInfo map[PublishID][]byte
	if len(distinctVals) > 0 || len(commonKeys) > 0 {
		publishInfo = make(map[PublishID][]byte)
		for publishID, mapSlice := range distinctVals {
			bytes, err := yaml.Marshal(mapSlice)
			if err != nil {
				return nil, errors.Wrapf(err, "failed to marshal publish config info")
			}
			publishInfo[publishID] = bytes
		}
		for publishID := range commonKeys {
			bytes, err := yaml.Marshal((*cfgPublishInfo)[publishID])
			if err != nil {
				return nil, errors.Wrapf(err, "failed to marshal publish config info")
			}
			publishInfo[publishID] = bytes
		}
	}
	return publishInfo, nil
}
