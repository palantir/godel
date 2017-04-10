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
	"sort"

	"github.com/palantir/pkg/matcher"
)

type Dist struct {
	// OutputDir is the directory to which the distribution is written.
	OutputDir string

	// InputDir is the path (from the project root) to a directory whose contents will be copied into the output
	// distribution directory at the beginning of the "dist" command. Can be used to include static resources and
	// other files required in a distribution.
	InputDir string

	// InputProducts is a slice of the names of products in the project (other than the current one) whose binaries
	// are required for the "dist" task. The "dist" task will ensure that the outputs of "build" exist for all of
	// the products specified in this slice (and will build the products as part of the task if necessary) and make
	// the outputs available to the "dist" script as environment variables. Note that the "dist" task only
	// guarantees that the products will be built and their locations will be available in the environment variables
	// provided to the script -- it is the responsibility of the user to write logic in the dist script to copy the
	// generated binaries.
	InputProducts []string

	// Script is the content of a script that is written to file a file and run after the initial distribution
	// process but before the artifact generation process. The contents of this value are written to a file with a
	// header `#!/bin/bash` with the contents of the global `dist-script-include` prepended and executed. The script
	// process inherits the environment variables of the Go process and also has the following environment variables
	// defined:
	//
	//   DIST_DIR: the absolute path to the root directory of the distribution created for the current product
	//   PROJECT_DIR: the root directory of project
	//   PRODUCT: product name,
	//   VERSION: product version
	//   IS_SNAPSHOT: 1 if the version contains a git hash as part of the string, 0 otherwise
	Script string

	// Info specifies the type of the distribution to be built and configuration for it. If unspecified, defaults to
	// a DistInfo of type SLSDistType.
	Info DistInfo

	// Publish is the configuration for the "publish" task.
	Publish Publish
}

type DistInfoType string

const (
	SLSDistType    DistInfoType = "sls"    // distribution that uses the Standard Layout Specification
	BinDistType    DistInfoType = "bin"    // distribution that includes all of the binaries for a product
	RPMDistType    DistInfoType = "rpm"    // RPM distribution
	DockerDistType DistInfoType = "docker" // docker image
)

type DistInfo interface {
	Type() DistInfoType
	// returns a list of products the dist depends on
	Deps() []string
}

type BinDistInfo struct {
	// OmitInitSh specifies whether or not the distribution should omit the auto-generated "init.sh" invocation
	// script. If true, the "init.sh" script will not be generated and included in the output distribution.
	OmitInitSh bool

	// InitShTemplateFile is the relative path to the template that should be used to generate the "init.sh" script.
	// If the value is absent, the default template will be used.
	InitShTemplateFile string
}

func (i *BinDistInfo) Type() DistInfoType {
	return BinDistType
}

func (i *BinDistInfo) Deps() []string {
	// no deps for bin type
	return nil
}

type SLSDistInfo struct {
	// InitShTemplateFile is the path to a template file that is used as the basis for the init.sh script of the
	// distribution. The path is relative to the project root directory. The contents of the file is processed using
	// Go templates and is provided with a distgo.ProductBuildSpec struct. If omitted, the default init.sh script
	// is used.
	InitShTemplateFile string

	// ManifestTemplateFile is the path to a template file that is used as the basis for the manifest.yml file of
	// the distribution. The path is relative to the project root directory. The contents of the file is processed
	// using Go templates and is provided with a distgo.ProductBuildSpec struct.
	ManifestTemplateFile string

	// ServiceArgs is the string provided as the service arguments for the default init.sh file generated for the distribution.
	ServiceArgs string

	// ProductType is the SLS product type for the distribution.
	ProductType string

	// ManifestExtensions contain the SLS manifest extensions for the distribution.
	ManifestExtensions map[string]interface{}

	// Reloadable will enable the `init.sh reload` command which sends SIGHUP to the process.
	Reloadable bool

	// YMLValidationExclude specifies a matcher used to specify YML files or paths that should not be validated as
	// part of creating the distribution. By default, the SLS distribution task verifies that all "*.yml" and
	// "*.yaml" files in the distribution are syntactically valid. If a distribution is known to ship with YML files
	// that are not valid YML, this parameter can be used to exclude those files from validation.
	YMLValidationExclude matcher.Matcher
}

func (i *SLSDistInfo) Type() DistInfoType {
	return SLSDistType
}

func (i *SLSDistInfo) Deps() []string {
	// no deps for sls type
	return nil
}

type RPMDistInfo struct {
	// Release is the release identifier that forms part of the name/version/release/architecture quadruplet
	// uniquely identifying the RPM package. Default is "1".
	Release string
	// ConfigFiles is a slice of absolute paths within the RPM that correspond to configuration files. RPM
	// identifies these as mutable. Default is no files.
	ConfigFiles []string
	// BeforeInstallScript is the content of shell script to run before this RPM is installed. Optional.
	BeforeInstallScript string
	// AfterInstallScript is the content of shell script to run immediately after this RPM is installed. Optional.
	AfterInstallScript string
	// AfterRemoveScript is the content of shell script to clean up after this RPM is removed. Optional.
	AfterRemoveScript string
}

func (i *RPMDistInfo) Type() DistInfoType {
	return RPMDistType
}

func (i *RPMDistInfo) Deps() []string {
	// no deps for rpm type
	return nil
}

type DockerDistDep struct {
	Product    string
	DistType   DistInfoType
	TargetFile string
}
type DockerDistDeps []DockerDistDep

func (d DockerDistDeps) ToMap() map[string]map[DistInfoType]string {
	m := make(map[string]map[DistInfoType]string)
	for _, dep := range d {
		if m[dep.Product] == nil {
			m[dep.Product] = make(map[DistInfoType]string)
		}
		m[dep.Product][dep.DistType] = dep.TargetFile
	}
	return m
}

type DockerDistInfo struct {
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
	DistDeps DockerDistDeps
}

func (d *DockerDistInfo) Type() DistInfoType {
	return DockerDistType
}

func (d *DockerDistInfo) Deps() []string {
	var deps []string
	for product := range d.DistDeps.ToMap() {
		deps = append(deps, product)
	}
	sort.Strings(deps)
	return deps
}
