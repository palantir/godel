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
	"gopkg.in/yaml.v2"

	"github.com/palantir/distgo/dister/osarchbin"
	"github.com/palantir/distgo/distgo"
	"github.com/palantir/distgo/dockerbuilder"
	"github.com/palantir/distgo/dockerbuilder/defaultdockerbuilder"
	defaultdockerbuilderconfig "github.com/palantir/distgo/dockerbuilder/defaultdockerbuilder/config"
)

type creatorWithUpgrader struct {
	creator  dockerbuilder.CreatorFunction
	upgrader distgo.ConfigUpgrader
}

func builtinDockerBuilders() map[string]creatorWithUpgrader {
	return map[string]creatorWithUpgrader{
		defaultdockerbuilder.TypeName: {
			creator: func(cfgYML []byte) (distgo.DockerBuilder, error) {
				var cfg defaultdockerbuilderconfig.Default
				if err := yaml.UnmarshalStrict(cfgYML, &cfg); err != nil {
					return nil, errors.Wrapf(err, "failed to unmarshal YAML")
				}
				return cfg.ToDockerBuilder(), nil
			},
			upgrader: distgo.NewConfigUpgrader(osarchbin.TypeName, defaultdockerbuilderconfig.UpgradeConfig),
		},
	}
}
