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

package publisherfactory

import (
	"github.com/palantir/distgo/distgo"
	"github.com/palantir/distgo/publisher"
	"github.com/palantir/distgo/publisher/artifactory"
	artifactoryconfig "github.com/palantir/distgo/publisher/artifactory/config"
	"github.com/palantir/distgo/publisher/bintray"
	bintrayconfig "github.com/palantir/distgo/publisher/bintray/config"
	"github.com/palantir/distgo/publisher/github"
	githubconfig "github.com/palantir/distgo/publisher/github/config"
	"github.com/palantir/distgo/publisher/mavenlocal"
	mavenlocalconfig "github.com/palantir/distgo/publisher/mavenlocal/config"
)

type creatorWithUpgrader struct {
	Creator  publisher.Creator
	Upgrader distgo.ConfigUpgrader
}

func builtinPublishers() map[string]creatorWithUpgrader {
	return map[string]creatorWithUpgrader{
		mavenlocal.TypeName: {
			Creator:  mavenlocal.PublisherCreator(),
			Upgrader: distgo.NewConfigUpgrader(mavenlocal.TypeName, mavenlocalconfig.UpgradeConfig),
		},
		artifactory.TypeName: {
			Creator:  artifactory.PublisherCreator(),
			Upgrader: distgo.NewConfigUpgrader(artifactory.TypeName, artifactoryconfig.UpgradeConfig),
		},
		bintray.TypeName: {
			Creator:  bintray.PublisherCreator(),
			Upgrader: distgo.NewConfigUpgrader(bintray.TypeName, bintrayconfig.UpgradeConfig),
		},
		github.TypeName: {
			Creator:  github.PublisherCreator(),
			Upgrader: distgo.NewConfigUpgrader(github.TypeName, githubconfig.UpgradeConfig),
		},
	}
}
