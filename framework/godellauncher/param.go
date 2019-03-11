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

package godellauncher

import (
	"github.com/palantir/pkg/specdir"
	"github.com/pkg/errors"

	"github.com/palantir/godel/framework/artifactresolver"
	"github.com/palantir/godel/framework/builtintasks/installupdate/layout"
)

const (
	GodelConfigYML = "godel.yml"
)

type TasksConfigProvidersParam struct {
	DefaultResolvers []artifactresolver.Resolver
	ConfigProviders  []artifactresolver.LocatorWithResolverParam
}

type PluginsParam struct {
	DefaultResolvers []artifactresolver.Resolver
	Plugins          []SinglePluginParam
}

type SinglePluginParam struct {
	artifactresolver.LocatorWithResolverParam
	Assets           []artifactresolver.LocatorWithResolverParam
	Override         bool
	FromPluginConfig bool
}

// ConfigDirPath returns the path to the gödel configuration directory given the path to the project directory. Returns
// an error if the directory structure does not match what is expected.
func ConfigDirPath(projectDirPath string) (string, error) {
	if projectDirPath == "" {
		return "", errors.Errorf("projectDirPath was empty")
	}
	wrapper, err := specdir.New(projectDirPath, layout.WrapperSpec(), nil, specdir.Validate)
	if err != nil {
		return "", err
	}
	return wrapper.Path(layout.WrapperConfigDir), nil
}
