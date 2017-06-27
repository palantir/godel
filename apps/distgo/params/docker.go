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

package params

import (
	"github.com/pkg/errors"
)

type DockerDepType string

const (
	DockerDepSLS    DockerDepType = "sls"
	DockerDepBin    DockerDepType = "bin"
	DockerDepRPM    DockerDepType = "rpm"
	DockerDepDocker DockerDepType = "docker"
)

type DockerDep struct {
	Product    string
	Type       DockerDepType
	TargetFile string
}

type DockerImage interface {
	SetRepository(repo string)
	SetDefaults(repo, tag string)
	Coordinates() (string, string)
	ContextDirectory() string
	Dependencies() []DockerDep
	Build(buildSpec ProductBuildSpecWithDeps) error
}

func ToDockerDepType(dep string) (DockerDepType, error) {
	switch dep {
	case "sls":
		return DockerDepSLS, nil
	case "rpm":
		return DockerDepRPM, nil
	case "bin":
		return DockerDepBin, nil
	case "docker":
		return DockerDepDocker, nil
	default:
		return "", errors.New("Invalid docker dependency type")
	}
}
