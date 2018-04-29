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
	"github.com/palantir/distgo/distgo"
	"github.com/palantir/distgo/distgo/config/internal/v0"
)

type RunConfig v0.RunConfig

func ToRunConfig(in *RunConfig) *v0.RunConfig {
	return (*v0.RunConfig)(in)
}

// ToParam returns the RunParam represented by the receiver *RunConfig and the provided default RunConfig. If a config
// value is specified (non-nil) in the receiver config, it is used. If a config value is not specified in the receiver
// config but is specified in the default config, the default config value is used. If a value is not specified in
// either configuration, the program-specified default value (if any) is used.
func (cfg *RunConfig) ToParam(defaultCfg RunConfig) distgo.RunParam {
	return distgo.RunParam{
		Args: getConfigValue(cfg.Args, defaultCfg.Args, nil).([]string),
	}
}
