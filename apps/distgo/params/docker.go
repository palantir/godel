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
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"

	"github.com/pkg/errors"
)

type DockerDepType string
type DockerImageType string

const (
	DockerDepSLS          DockerDepType = "sls"
	DockerDepBin          DockerDepType = "bin"
	DockerDepRPM          DockerDepType = "rpm"
	DockerDepDocker       DockerDepType = "docker"
	ManifestLabel         string        = "com.palantir.sls.manifest"
	ConfigurationLabel    string        = "com.palantir.sls.configuration"
	ConfigurationFileName string        = "configuration.yml"
)

type DockerDep struct {
	Product    string
	Type       DockerDepType
	TargetFile string
}

type DockerImage interface {
	SetRepository(repo string)
	SetDefaults(repo string, tag string)
	Coordinates() (string, string)
	ContextDirectory() string
	Dependencies() []DockerDep
	Build(buildSpec ProductBuildSpecWithDeps) error
}

type DefaultDockerImage struct {
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
	Deps []DockerDep
}

type SLSDockerImage struct {
	DefaultDockerImage
	GroupID      string
	ProuductType string
	Extensions   map[string]interface{}
}

func (sdi *SLSDockerImage) Build(buildSpec ProductBuildSpecWithDeps) error {
	contextDir := path.Join(buildSpec.Spec.ProjectDir, sdi.ContextDir)
	configFile := path.Join(contextDir, ConfigurationFileName)
	var args []string
	args = append(args, "build")
	args = append(args, "--tag", fmt.Sprintf("%s:%s", sdi.Repository, sdi.Tag))
	manifest, err := GetManifest(sdi.GroupID, buildSpec.Spec.ProductName, buildSpec.Spec.ProductVersion, sdi.ProuductType, sdi.Extensions)
	if err != nil {
		return errors.Wrap(err, "Failed to get manifest for the image")
	}
	args = append(args, "--label", fmt.Sprintf("%s=%s", ManifestLabel, manifest))
	if _, err := os.Stat(configFile); err == nil {
		content, err := ioutil.ReadFile(configFile)
		if err != nil {
			return errors.Wrapf(err, "Failed to read the file %s", configFile)
		}
		args = append(args, "--label", fmt.Sprintf("%s=%s", ConfigurationLabel, string(content)))
	}
	args = append(args, contextDir)

	buildCmd := exec.Command("docker", args...)
	if output, err := buildCmd.CombinedOutput(); err != nil {
		return errors.Wrap(err, fmt.Sprintf("docker build failed with error:\n%s\n", string(output)))
	}
	return nil
}

func (di *DefaultDockerImage) ContextDirectory() string {
	return di.ContextDir
}

func (di *DefaultDockerImage) Dependencies() []DockerDep {
	return di.Deps
}

func (di *DefaultDockerImage) Build(buildSpec ProductBuildSpecWithDeps) error {
	contextDir := path.Join(buildSpec.Spec.ProjectDir, di.ContextDir)
	var args []string
	args = append(args, "build")
	args = append(args, "--tag", fmt.Sprintf("%s:%s", di.Repository, di.Tag))
	args = append(args, contextDir)

	buildCmd := exec.Command("docker", args...)
	if output, err := buildCmd.CombinedOutput(); err != nil {
		return errors.Wrap(err, fmt.Sprintf("docker build failed with error:\n%s\n", string(output)))
	}
	return nil
}

func (di *DefaultDockerImage) Coordinates() (string, string) {
	return di.Repository, di.Tag
}

func (di *DefaultDockerImage) SetDefaults(repo, tag string) {
	if di.Repository == "" {
		di.Repository = repo
	}
	if di.Tag == "" {
		di.Tag = tag
	}
}

func (di *DefaultDockerImage) SetRepository(repo string) {
	di.Repository = repo
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
