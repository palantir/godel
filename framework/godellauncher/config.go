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
	"io/ioutil"
	"os"
	"path"

	"github.com/palantir/pkg/matcher"
	"github.com/palantir/pkg/specdir"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"

	"github.com/palantir/godel/framework/artifactresolver"
	"github.com/palantir/godel/framework/builtintasks/installupdate/layout"
	"github.com/palantir/godel/pkg/versionedconfig"
)

const (
	GodelConfigYML = "godel.yml"
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

// Combine combines the provided TasksConfig configurations with the base configuration. In cases where values are
// overwritten, the last (most recent) values in the inputs will take precedence.
func (c *TasksConfig) Combine(configs ...TasksConfig) {
	if c.DefaultTasks.Tasks == nil {
		c.DefaultTasks.Tasks = make(map[string]SingleDefaultTaskConfig)
	}

	for _, cfg := range configs {
		// DefaultTask resolvers are appended
		c.DefaultTasks.DefaultResolvers = append(c.DefaultTasks.DefaultResolvers, cfg.DefaultTasks.DefaultResolvers...)
		// DefaultTask tasks key/values are simply copied (and overwritten with last writer wins for any duplicate keys)
		for k, v := range cfg.DefaultTasks.Tasks {
			c.DefaultTasks.Tasks[k] = v
		}

		// Plugin resolvers and definitions are appended
		c.Plugins.DefaultResolvers = append(c.Plugins.DefaultResolvers, cfg.Plugins.DefaultResolvers...)
		c.Plugins.Plugins = append(c.Plugins.Plugins, cfg.Plugins.Plugins...)
	}
}

type DefaultTasksConfig struct {
	DefaultResolvers []string                           `yaml:"resolvers"`
	Tasks            map[string]SingleDefaultTaskConfig `yaml:"tasks"`
}

type TasksConfigProvidersParam struct {
	DefaultResolvers []artifactresolver.Resolver
	ConfigProviders  []artifactresolver.LocatorWithResolverParam
}

type TasksConfigProvidersConfig struct {
	DefaultResolvers []string                                                   `yaml:"resolvers"`
	ConfigProviders  []artifactresolver.ConfigProviderLocatorWithResolverConfig `yaml:"providers"`
}

func (c *TasksConfigProvidersConfig) ToParam() (TasksConfigProvidersParam, error) {
	var defaultResolvers []artifactresolver.Resolver
	for _, resolverStr := range c.DefaultResolvers {
		resolver, err := artifactresolver.NewTemplateResolver(resolverStr)
		if err != nil {
			return TasksConfigProvidersParam{}, err
		}
		defaultResolvers = append(defaultResolvers, resolver)
	}
	var configProviders []artifactresolver.LocatorWithResolverParam
	for _, provider := range c.ConfigProviders {
		providerVal, err := provider.ToParam()
		if err != nil {
			return TasksConfigProvidersParam{}, err
		}
		configProviders = append(configProviders, providerVal)
	}
	return TasksConfigProvidersParam{
		DefaultResolvers: defaultResolvers,
		ConfigProviders:  configProviders,
	}, nil
}

type SingleDefaultTaskConfig struct {
	// LocatorWithResolverConfig contains the configuration for the locator and resolver. Any value provided here
	// overrides the default value.
	artifactresolver.LocatorWithResolverConfig `yaml:",inline"`
	// ExcludeAllDefaultAssets specifies whether or not all of the default assets should be excluded. If this value is
	// true, then DefaultAssetsToExclude is ignored.
	ExcludeAllDefaultAssets bool `yaml:"exclude-all-default-assets"`
	// DefaultAssetsToExclude specifies the assets that should be excluded if they are provided by the default
	// configuration. Only used if ExcludeAllDefaultAssets is false.
	DefaultAssetsToExclude []string `yaml:"exclude-default-assets"`
	// Assets specifies the custom assets that should be added to the default task.
	Assets []artifactresolver.LocatorWithResolverConfig `yaml:"assets"`
}

type PluginsParam struct {
	DefaultResolvers []artifactresolver.Resolver
	Plugins          []SinglePluginParam
}

type PluginsConfig struct {
	DefaultResolvers []string             `yaml:"resolvers"`
	Plugins          []SinglePluginConfig `yaml:"plugins"`
}

func (c *PluginsConfig) ToParam() (PluginsParam, error) {
	var defaultResolvers []artifactresolver.Resolver
	for _, resolverStr := range c.DefaultResolvers {
		resolver, err := artifactresolver.NewTemplateResolver(resolverStr)
		if err != nil {
			return PluginsParam{}, err
		}
		defaultResolvers = append(defaultResolvers, resolver)
	}
	var plugins []SinglePluginParam
	for _, plugin := range c.Plugins {
		pluginParam, err := plugin.ToParam()
		if err != nil {
			return PluginsParam{}, err
		}
		plugins = append(plugins, pluginParam)
	}
	return PluginsParam{
		DefaultResolvers: defaultResolvers,
		Plugins:          plugins,
	}, nil
}

type SinglePluginParam struct {
	artifactresolver.LocatorWithResolverParam
	Assets []artifactresolver.LocatorWithResolverParam
}

type SinglePluginConfig struct {
	// LocatorWithResolverConfig stores the locator and the resolver for the plugin.
	artifactresolver.LocatorWithResolverConfig `yaml:",inline"`
	// Assets stores the locators and resolvers for the assets for this plugin.
	Assets []artifactresolver.LocatorWithResolverConfig `yaml:"assets"`
}

func (c *SinglePluginConfig) ToParam() (SinglePluginParam, error) {
	locatorWithResolverParam, err := c.LocatorWithResolverConfig.ToParam()
	if err != nil {
		return SinglePluginParam{}, err
	}
	var assets []artifactresolver.LocatorWithResolverParam
	for _, assetCfg := range c.Assets {
		assetParamVal, err := assetCfg.ToParam()
		if err != nil {
			return SinglePluginParam{}, err
		}
		assets = append(assets, assetParamVal)
	}
	return SinglePluginParam{
		LocatorWithResolverParam: locatorWithResolverParam,
		Assets: assets,
	}, nil
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

// ReadGodelConfigFromProjectDir reads the gödel configuration from the "godel.yml" file in the configuration file for
// the gödel project with the specified project directory and returns it.
func ReadGodelConfigFromProjectDir(projectDir string) (GodelConfig, error) {
	cfgDir, err := ConfigDirPath(projectDir)
	if err != nil {
		return GodelConfig{}, err
	}
	cfg, err := ReadGodelConfig(cfgDir)
	if err != nil {
		return GodelConfig{}, err
	}
	return cfg, nil
}

// ReadGodelConfig reads the gödel configuration from the "godel.yml" file in the specified directory and returns it.
func ReadGodelConfig(cfgDir string) (GodelConfig, error) {
	var godelCfg GodelConfig
	godelYML := path.Join(cfgDir, GodelConfigYML)
	if _, err := os.Stat(godelYML); err == nil {
		bytes, err := ioutil.ReadFile(godelYML)
		if err != nil {
			return GodelConfig{}, errors.Wrapf(err, "failed to read file %s", godelYML)
		}
		upgradedBytes, err := UpgradeGodelConfig(bytes)
		if err != nil {
			return GodelConfig{}, errors.Wrapf(err, "failed to upgrade configuration")
		}
		if err := yaml.Unmarshal(upgradedBytes, &godelCfg); err != nil {
			return GodelConfig{}, errors.Wrapf(err, "failed to unmarshal gödel config YAML")
		}
	}
	return godelCfg, nil
}

func UpgradeGodelConfig(cfgBytes []byte) ([]byte, error) {
	version, err := versionedconfig.ConfigVersion(cfgBytes)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to unmarshal version")
	}
	switch version {
	case "", "0":
		// verify that configuration is valid
		var cfg GodelConfig
		if err := yaml.UnmarshalStrict(cfgBytes, &cfg); err != nil {
			return nil, errors.Wrapf(err, "failed to unmarshal configuration")
		}
		// configuration is valid and is current version: return input bytes
		return cfgBytes, nil
	default:
		return nil, errors.Errorf("unsupported version: %s", version)
	}
}
