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
	"strings"

	"github.com/palantir/godel/framework/artifactresolver"
)

const defaultResolver = "https://palantir.bintray.com/releases/{{GroupPath}}/{{Product}}/{{Version}}/{{Product}}-{{Version}}-{{OS}}-{{Arch}}.tgz"

var defaultPluginsConfig = PluginsConfig{
	DefaultResolvers: []string{defaultResolver},
}

type TasksConfigInfo struct {
	// BuiltinPluginsConfig is the configuration for built-in plugins that is built as part of gödel.
	BuiltinPluginsConfig PluginsConfig
	// TasksConfig is the fully resolved user-provided tasks configuration.
	TasksConfig TasksConfig
	// DefaultTasksPluginsConfig is the plugin configuration used to load the default tasks. It is a result of combining
	// the BuiltinPluginsConfig with the DefaultTasks config of TasksConfig.
	DefaultTasksPluginsConfig PluginsConfig
}

func BuiltinDefaultPluginsConfig() PluginsConfig {
	return defaultPluginsConfig
}

func DefaultTasksPluginsConfig(config DefaultTasksConfig) PluginsConfig {
	// start with configuration that uses default resolver
	pluginsCfg := PluginsConfig{
		DefaultResolvers: defaultPluginsConfig.DefaultResolvers,
	}
	// append default resolvers provided by the configuration
	pluginsCfg.DefaultResolvers = append(pluginsCfg.DefaultResolvers, config.DefaultResolvers...)

	for _, currPlugin := range defaultPluginsConfig.Plugins {
		currKey := locatorIDWithoutVersion(currPlugin.Locator.ID)

		cfgParam, ok := config.Tasks[currKey]
		if !ok {
			// if custom configuration is not specified, use default and continue
			pluginsCfg.Plugins = append(pluginsCfg.Plugins, currPlugin)
			continue
		}

		// custom configuration was non-empty: start it with default LocatorWithResolver configuration
		currCfg := SinglePluginConfig{
			LocatorWithResolverConfig: currPlugin.LocatorWithResolverConfig,
		}
		if cfgParam.Locator.ID != "" {
			currCfg.Locator = cfgParam.Locator
		}
		if cfgParam.Resolver != "" {
			currCfg.Resolver = cfgParam.Resolver
		}

		currCfg.Assets = append(currCfg.Assets, assetConfigFromDefault(currPlugin.Assets, cfgParam)...)
		currCfg.Assets = append(currCfg.Assets, cfgParam.Assets...)
		pluginsCfg.Plugins = append(pluginsCfg.Plugins, currCfg)
	}
	return pluginsCfg
}

func assetConfigFromDefault(baseCfg []artifactresolver.LocatorWithResolverConfig, cfg SingleDefaultTaskConfig) []artifactresolver.LocatorWithResolverConfig {
	if cfg.ExcludeAllDefaultAssets {
		return nil
	}
	exclude := make(map[string]struct{})
	for _, currExclude := range cfg.DefaultAssetsToExclude {
		exclude[currExclude] = struct{}{}
	}
	var out []artifactresolver.LocatorWithResolverConfig
	for _, asset := range baseCfg {
		if _, ok := exclude[locatorIDWithoutVersion(asset.Locator.ID)]; ok {
			continue
		}
		out = append(out, asset)
	}
	return out
}

func locatorIDWithoutVersion(locatorID string) string {
	parts := strings.Split(locatorID, ":")
	return strings.Join(parts[:2], ":")
}
