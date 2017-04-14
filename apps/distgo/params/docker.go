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
type DockerDeps []DockerDep

func (d DockerDeps) ToMap() map[string]map[DockerDepType]string {
	m := make(map[string]map[DockerDepType]string)
	for _, dep := range d {
		if m[dep.Product] == nil {
			m[dep.Product] = make(map[DockerDepType]string)
		}
		m[dep.Product][dep.Type] = dep.TargetFile
	}
	return m
}

type DockerImage struct {
	// Repository and Tag are the part of the image coordinates.
	// For example, in alpine:latest, alpine is the repository
	// and the latest is the tag
	Repository string
	Tag        string
	// ContextDir is the directory in which the docker build task is executed.
	ContextDir string
	// DistDeps is a slice of DockerDistDep.
	// DockerDistDep contains a product, dist type and target file.
	// For a particular product's dist type, we create a link from its output
	// inside the ContextDir with the name specified in target file.
	// This will be used to order the dist tasks such that all the dependent
	// products' dist tasks will be executed first, after which the dist tasks for the
	// current product are executed.
	Deps DockerDeps
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
		return "", nil
	}
}
