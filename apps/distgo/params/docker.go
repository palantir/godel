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
type DockerImageType string

const (
	DockerDepSLS           DockerDepType   = "sls"
	DockerDepBin           DockerDepType   = "bin"
	DockerDepRPM           DockerDepType   = "rpm"
	DockerDepDocker        DockerDepType   = "docker"
	DefaultDockerImageType DockerImageType = "default"
	SLSDockerImageType     DockerImageType = "sls"
)

type DockerDep struct {
	Product    string
	Type       DockerDepType
	TargetFile string
}

type DockerImage struct {
	// Repository and Tag are the part of the image coordinates.
	// For example, in alpine:latest, alpine is the repository
	// and the latest is the tag
	Repository string
	Tag        string
	// ContextDir is the directory in which the docker build task is executed.
	ContextDir      string
	BuildArgsScript string
	// DistDeps is a slice of DockerDistDep.
	// DockerDistDep contains a product, dist type and target file.
	// For a particular product's dist type, we create a link from its output
	// inside the ContextDir with the name specified in target file.
	// This will be used to order the dist tasks such that all the dependent
	// products' dist tasks will be executed first, after which the dist tasks for the
	// current product are executed.
	Deps []DockerDep
	Info DockerImageInfo
}

type DockerImageInfo interface {
	Type() DockerImageType
}

type DefaultDockerImageInfo struct{}

func (info *DefaultDockerImageInfo) Type() DockerImageType {
	return DefaultDockerImageType
}

type SLSDockerImageInfo struct {
	GroupID      string
	ProuductType string
	Extensions   map[string]interface{}
}

func (info *SLSDockerImageInfo) Type() DockerImageType {
	return SLSDockerImageType
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
