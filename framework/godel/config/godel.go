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

package config

import (
	"github.com/palantir/godel/framework/artifactresolver"
	"github.com/palantir/godel/framework/godel/config/internal/v0"
	"github.com/palantir/godel/framework/godellauncher"
	"github.com/palantir/godel/framework/internal/pluginsinternal"
)

type GodelConfig v0.GodelConfig

type TasksConfig v0.TasksConfig

type TasksConfigInfo struct {
	// BuiltinPluginsConfig is the configuration for built-in plugins that is built as part of gödel.
	BuiltinPluginsConfig PluginsConfig
	// TasksConfig is the fully resolved user-provided tasks configuration.
	TasksConfig TasksConfig
	// DefaultTasksPluginsConfig is the plugin configuration used to load the default tasks. It is a result of combining
	// the BuiltinPluginsConfig with the DefaultTasks config of TasksConfig.
	DefaultTasksPluginsConfig PluginsConfig
}

func ToTasksConfig(in TasksConfig) v0.TasksConfig {
	return v0.TasksConfig(in)
}

// Combine combines the provided TasksConfig configurations with the base configuration. In cases where values are
// overwritten, the last (most recent) values in the inputs will take precedence.
func (c *TasksConfig) Combine(configs ...TasksConfig) {
	if c.DefaultTasks.Tasks == nil {
		c.DefaultTasks.Tasks = make(map[string]v0.SingleDefaultTaskConfig)
	}

	for _, cfg := range configs {
		// DefaultTask resolvers are appended and uniquified
		c.DefaultTasks.DefaultResolvers = pluginsinternal.Uniquify(append(c.DefaultTasks.DefaultResolvers, cfg.DefaultTasks.DefaultResolvers...))

		// DefaultTask tasks key/values are simply copied (and overwritten with last writer wins for any duplicate keys)
		for k, v := range cfg.DefaultTasks.Tasks {
			c.DefaultTasks.Tasks[k] = v
		}

		// Plugin resolvers and definitions are appended and uniquified
		c.Plugins.DefaultResolvers = pluginsinternal.Uniquify(append(c.Plugins.DefaultResolvers, cfg.Plugins.DefaultResolvers...))
		c.Plugins.Plugins = append(c.Plugins.Plugins, cfg.Plugins.Plugins...)
	}
}

type DefaultTasksConfig v0.DefaultTasksConfig

func ToDefaultTasksConfig(in DefaultTasksConfig) v0.DefaultTasksConfig {
	return v0.DefaultTasksConfig(in)
}

type TasksConfigProvidersConfig v0.TasksConfigProvidersConfig

func (c *TasksConfigProvidersConfig) ToParam() (godellauncher.TasksConfigProvidersParam, error) {
	var defaultResolvers []artifactresolver.Resolver
	for _, resolverStr := range c.DefaultResolvers {
		resolver, err := artifactresolver.NewTemplateResolver(resolverStr)
		if err != nil {
			return godellauncher.TasksConfigProvidersParam{}, err
		}
		defaultResolvers = append(defaultResolvers, resolver)
	}
	var configProviders []artifactresolver.LocatorWithResolverParam
	for _, provider := range c.ConfigProviders {
		provider := ConfigProviderLocatorWithResolverConfig(provider)
		providerVal, err := provider.ToParam()
		if err != nil {
			return godellauncher.TasksConfigProvidersParam{}, err
		}
		configProviders = append(configProviders, providerVal)
	}
	return godellauncher.TasksConfigProvidersParam{
		DefaultResolvers: defaultResolvers,
		ConfigProviders:  configProviders,
	}, nil
}

type SingleDefaultTaskConfig v0.SingleDefaultTaskConfig

func ToSingleDefaultTaskConfig(in SingleDefaultTaskConfig) v0.SingleDefaultTaskConfig {
	return v0.SingleDefaultTaskConfig(in)
}

func ToTasks(in map[string]SingleDefaultTaskConfig) map[string]v0.SingleDefaultTaskConfig {
	if in == nil {
		return nil
	}
	out := make(map[string]v0.SingleDefaultTaskConfig, len(in))
	for k, v := range in {
		out[k] = ToSingleDefaultTaskConfig(v)
	}
	return out
}

type PluginsConfig v0.PluginsConfig

func ToPluginsConfig(in PluginsConfig) v0.PluginsConfig {
	return v0.PluginsConfig(in)
}

func (c *PluginsConfig) ToParam() (godellauncher.PluginsParam, error) {
	var defaultResolvers []artifactresolver.Resolver
	for _, resolverStr := range c.DefaultResolvers {
		resolver, err := artifactresolver.NewTemplateResolver(resolverStr)
		if err != nil {
			return godellauncher.PluginsParam{}, err
		}
		defaultResolvers = append(defaultResolvers, resolver)
	}
	var plugins []godellauncher.SinglePluginParam
	for _, plugin := range c.Plugins {
		plugin := SinglePluginConfig(plugin)
		pluginParam, err := plugin.ToParam()
		if err != nil {
			return godellauncher.PluginsParam{}, err
		}
		plugins = append(plugins, pluginParam)
	}
	return godellauncher.PluginsParam{
		DefaultResolvers: defaultResolvers,
		Plugins:          plugins,
	}, nil
}

type SinglePluginConfig v0.SinglePluginConfig

func ToSinglePluginConfig(in SinglePluginConfig) v0.SinglePluginConfig {
	return v0.SinglePluginConfig(in)
}

func ToSinglePluginConfigs(in []SinglePluginConfig) []v0.SinglePluginConfig {
	if in == nil {
		return nil
	}
	out := make([]v0.SinglePluginConfig, len(in))
	for i, v := range in {
		out[i] = ToSinglePluginConfig(v)
	}
	return out
}

func (c *SinglePluginConfig) ToParam() (godellauncher.SinglePluginParam, error) {
	locatorWithResolverConfig := LocatorWithResolverConfig(c.LocatorWithResolverConfig)
	locatorWithResolverParam, err := locatorWithResolverConfig.ToParam()
	if err != nil {
		return godellauncher.SinglePluginParam{}, err
	}
	var assets []artifactresolver.LocatorWithResolverParam
	for _, assetCfg := range c.Assets {
		assetCfg := LocatorWithResolverConfig(assetCfg)
		assetParamVal, err := assetCfg.ToParam()
		if err != nil {
			return godellauncher.SinglePluginParam{}, err
		}
		assets = append(assets, assetParamVal)
	}
	return godellauncher.SinglePluginParam{
		LocatorWithResolverParam: locatorWithResolverParam,
		Assets: assets,
	}, nil
}
