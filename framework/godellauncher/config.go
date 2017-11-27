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
	"encoding/json"
	"io/ioutil"
	"os"
	"path"

	"github.com/palantir/pkg/matcher"
	"github.com/palantir/pkg/specdir"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"

	"github.com/palantir/godel/framework/builtintasks/installupdate/layout"
)

const (
	GodelConfigYML   = "godel.yml"
	excludeConfigYML = "exclude.yml"
)

type GodelConfig struct {
	// Plugins specifies the configuration for the plugins configured for gödel. Excluded from JSON serialization
	// because JSON serialization is only needed for legacy "exclude" back-compat (and will be removed in 2.0 release).
	Plugins PluginsConfig `yaml:"plugins" json:"-"`
	// Exclude specifies the files and directories that should be excluded from gödel operations.
	Exclude matcher.NamesPathsCfg `yaml:"exclude" json:"exclude"`
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

// GodelConfigJSON returns the JSON representation of the gödel configuration read by ReadGodelConfig.
func GodelConfigJSON(cfgDir string) ([]byte, error) {
	cfg, err := ReadGodelConfig(cfgDir)
	if err != nil {
		return nil, err
	}
	bytes, err := json.Marshal(cfg)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to marshal configuration as JSON")
	}
	return bytes, nil
}

// ReadGodelConfig reads the gödel configuration from the "gödel.yml" file in the specified directory and returns it. If
// "exclude.yml" exists in the directory, it is also read and its elements are combined with the configuration read from
// "gödel.yml".
func ReadGodelConfig(cfgDir string) (GodelConfig, error) {
	var gödelCfg GodelConfig
	gödelYML := path.Join(cfgDir, GodelConfigYML)
	if _, err := os.Stat(gödelYML); err == nil {
		bytes, err := ioutil.ReadFile(gödelYML)
		if err != nil {
			return GodelConfig{}, errors.Wrapf(err, "failed to read file %s", gödelYML)
		}
		if err := yaml.Unmarshal(bytes, &gödelCfg); err != nil {
			return GodelConfig{}, errors.Wrapf(err, "failed to unmarshal gödel config YAML")
		}
	}

	// legacy support: if "exclude.yml" exists, combine the "Exclude" configuration it defines with the new one
	excludeYML := path.Join(cfgDir, excludeConfigYML)
	if _, err := os.Stat(excludeYML); err == nil {
		var excludeCfg matcher.NamesPathsCfg
		bytes, err := ioutil.ReadFile(excludeYML)
		if err != nil {
			return GodelConfig{}, errors.Wrapf(err, "failed to read file %s", excludeYML)
		}
		if err := yaml.Unmarshal(bytes, &excludeCfg); err != nil {
			return GodelConfig{}, errors.Wrapf(err, "failed to unmarshal exclude config YAML")
		}
		gödelCfg.Exclude.Names = addNewElements(gödelCfg.Exclude.Names, excludeCfg.Names)
		gödelCfg.Exclude.Paths = addNewElements(gödelCfg.Exclude.Paths, excludeCfg.Paths)
	}

	return gödelCfg, nil
}

func addNewElements(original, new []string) []string {
	set := make(map[string]struct{})
	for _, s := range original {
		set[s] = struct{}{}
	}

	for _, s := range new {
		if _, ok := set[s]; ok {
			continue
		}
		original = append(original, s)
	}
	return original
}
