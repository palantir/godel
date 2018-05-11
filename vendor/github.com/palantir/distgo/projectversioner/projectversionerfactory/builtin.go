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
	"gopkg.in/yaml.v2"

	"github.com/palantir/distgo/distgo"
	"github.com/palantir/distgo/projectversioner"
	"github.com/palantir/distgo/projectversioner/git"
	gitconfig "github.com/palantir/distgo/projectversioner/git/config"
	"github.com/palantir/distgo/projectversioner/script"
	scriptconfig "github.com/palantir/distgo/projectversioner/script/config"
)

type creatorWithUpgrader struct {
	creator  projectversioner.CreatorFunction
	upgrader distgo.ConfigUpgrader
}

func builtinProjectVersioners() map[string]creatorWithUpgrader {
	return map[string]creatorWithUpgrader{
		git.TypeName: {
			creator: func(cfgYML []byte) (distgo.ProjectVersioner, error) {
				return git.New(), nil
			},
			upgrader: distgo.NewConfigUpgrader(git.TypeName, gitconfig.UpgradeConfig),
		},
		script.TypeName: {
			creator: func(cfgYML []byte) (distgo.ProjectVersioner, error) {
				var cfg scriptconfig.Script
				if err := yaml.UnmarshalStrict(cfgYML, &cfg); err != nil {
					return nil, errors.Wrapf(err, "failed to unmarshal YAML")
				}
				return cfg.ToProjectVersioner(), nil
			},
			upgrader: distgo.NewConfigUpgrader(script.TypeName, scriptconfig.UpgradeConfig),
		},
	}
}
