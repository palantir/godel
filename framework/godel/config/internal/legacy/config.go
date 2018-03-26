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

package legacy

import (
	"sort"

	"github.com/palantir/pkg/matcher"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"

	"github.com/palantir/godel/framework/godel/config/internal/v0"
)

type GodelConfig struct {
	// DefaultTasks specifies the configuration for the default tasks for gödel.
	DefaultTasks DefaultTasksConfig `yaml:"default-tasks"`
	// Plugins specifies the configuration for the plugins configured for gödel.
	Plugins PluginsConfig `yaml:"plugins"`
	// Exclude specifies the files and directories that should be excluded from gödel operations.
	Exclude matcher.NamesPathsCfg `yaml:"exclude"`
}

type DefaultTasksConfig map[string]SingleDefaultTaskConfig

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

type LocatorWithResolverConfig struct {
	Locator  LocatorConfig `yaml:"locator"`
	Resolver string        `yaml:"resolver"`
}

type LocatorConfig struct {
	ID        string            `yaml:"id"`
	Checksums map[string]string `yaml:"checksums"`
}

func UpgradeConfig(cfgBytes []byte) ([]byte, error) {
	var legacyCfg GodelConfig
	if err := yaml.UnmarshalStrict(cfgBytes, &legacyCfg); err != nil {
		return nil, errors.Wrapf(err, "failed to unmarshal license-plugin legacy configuration")
	}

	v0Cfg := v0.GodelConfig{}

	// DefaultTasks
	var sortedKeys []string
	for k := range legacyCfg.DefaultTasks {
		sortedKeys = append(sortedKeys, k)
	}
	sort.Strings(sortedKeys)
	if len(sortedKeys) > 0 {
		v0Cfg.TasksConfig.DefaultTasks.Tasks = make(map[string]v0.SingleDefaultTaskConfig)
		for _, k := range sortedKeys {
			legacyDefaultTaskCfg := legacyCfg.DefaultTasks[k]
			v0Cfg.TasksConfig.DefaultTasks.Tasks[k] = v0.SingleDefaultTaskConfig{
				LocatorWithResolverConfig: toV0LocatorWithResolverConfig(legacyDefaultTaskCfg.LocatorWithResolverConfig),
				ExcludeAllDefaultAssets:   legacyDefaultTaskCfg.ExcludeAllDefaultAssets,
				DefaultAssetsToExclude:    legacyDefaultTaskCfg.DefaultAssetsToExclude,
				Assets:                    toV0LocatorWithResolverConfigs(legacyDefaultTaskCfg.Assets),
			}
		}

	}

	// Plugins
	for _, legacyPluginConfig := range legacyCfg.Plugins.Plugins {
		v0Cfg.Plugins.Plugins = append(v0Cfg.Plugins.Plugins, v0.SinglePluginConfig{
			LocatorWithResolverConfig: toV0LocatorWithResolverConfig(legacyPluginConfig.LocatorWithResolverConfig),
			Assets: toV0LocatorWithResolverConfigs(legacyPluginConfig.Assets),
		})
	}

	// Exclude
	v0Cfg.Exclude = legacyCfg.Exclude

	upgradedBytes, err := yaml.Marshal(v0Cfg)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to marshal godel v0 configuration")
	}
	return upgradedBytes, nil
}

func toV0LocatorWithResolverConfig(in LocatorWithResolverConfig) v0.LocatorWithResolverConfig {
	return v0.LocatorWithResolverConfig{
		Locator: v0.LocatorConfig{
			ID:        in.Locator.ID,
			Checksums: in.Locator.Checksums,
		},
		Resolver: in.Resolver,
	}
}

func toV0LocatorWithResolverConfigs(in []LocatorWithResolverConfig) []v0.LocatorWithResolverConfig {
	if in == nil {
		return nil
	}
	out := make([]v0.LocatorWithResolverConfig, len(in))
	for i, v := range in {
		out[i] = toV0LocatorWithResolverConfig(v)
	}
	return out
}
