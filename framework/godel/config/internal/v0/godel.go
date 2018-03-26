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

// Copyright (c) 2016 Palantir Technologies Inc. All rights reserved.
// Use of this source code is governed by the Apache License, Version 2.0
// that can be found in the LICENSE file.

package v0

import (
	"github.com/palantir/pkg/matcher"

	"github.com/palantir/godel/pkg/versionedconfig"
)

type GodelConfig struct {
	// Version of the configuration
	versionedconfig.ConfigWithVersion `yaml:",inline"`

	// TasksConfigProviders specifies the providers used to load provided task configuration.
	TasksConfigProviders TasksConfigProvidersConfig `yaml:"tasks-config-providers"`

	// TasksConfig contains the configuration for the tasks (default and plugin).
	TasksConfig `yaml:",inline"`

	// Exclude specifies the files and directories that should be excluded from gödel operations.
	Exclude matcher.NamesPathsCfg `yaml:"exclude"`
}

type TasksConfig struct {
	// DefaultTasks specifies the configuration for the default tasks for gödel.
	DefaultTasks DefaultTasksConfig `yaml:"default-tasks"`
	// Plugins specifies the configuration for the plugins configured for gödel.
	Plugins PluginsConfig `yaml:"plugins"`
}

type DefaultTasksConfig struct {
	DefaultResolvers []string                           `yaml:"resolvers"`
	Tasks            map[string]SingleDefaultTaskConfig `yaml:"tasks"`
}

type TasksConfigProvidersConfig struct {
	DefaultResolvers []string                                  `yaml:"resolvers"`
	ConfigProviders  []ConfigProviderLocatorWithResolverConfig `yaml:"providers"`
}

type SingleDefaultTaskConfig struct {
	// LocatorWithResolverConfig contains the configuration for the locator and resolver. Any value provided here
	// overrides the default value.
	LocatorWithResolverConfig `yaml:",inline"`
	// ExcludeAllDefaultAssets specifies whether or not all of the default assets should be excluded. If this value is
	// true, then DefaultAssetsToExclude is ignored.
	ExcludeAllDefaultAssets bool `yaml:"exclude-all-default-assets"`
	// DefaultAssetsToExclude specifies the assets that should be excluded if they are provided by the default
	// configuration. Only used if ExcludeAllDefaultAssets is false.
	DefaultAssetsToExclude []string `yaml:"exclude-default-assets"`
	// Assets specifies the custom assets that should be added to the default task.
	Assets []LocatorWithResolverConfig `yaml:"assets"`
}

type PluginsConfig struct {
	DefaultResolvers []string             `yaml:"resolvers"`
	Plugins          []SinglePluginConfig `yaml:"plugins"`
}

type SinglePluginConfig struct {
	// LocatorWithResolverConfig stores the locator and the resolver for the plugin.
	LocatorWithResolverConfig `yaml:",inline"`
	// Assets stores the locators and resolvers for the assets for this plugin.
	Assets []LocatorWithResolverConfig `yaml:"assets"`
}
