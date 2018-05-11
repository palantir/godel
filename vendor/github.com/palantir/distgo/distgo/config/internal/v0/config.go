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

package v0

import (
	"bytes"
	"sort"

	"github.com/palantir/pkg/matcher"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"

	"github.com/palantir/distgo/distgo"
)

type ProjectConfig struct {
	// Products maps product names to configurations.
	Products map[distgo.ProductID]ProductConfig `yaml:"products,omitempty"`

	// ProductDefaults specifies the default values that should be used for unspecified values in the products map. If a
	// field in a top-level key in a "ProductConfig" value in the "Products" map is nil and the corresponding value in
	// ProductDefaults is non-nil, the value in ProductDefaults is used.
	ProductDefaults ProductConfig `yaml:"product-defaults,omitempty"`

	// ScriptIncludes specifies a string that is appended to every script that is written out. Can be used to define
	// functions or constants for all scripts.
	ScriptIncludes string `yaml:"script-includes,omitempty"`

	// ProjectVersioner specifies the operation that is used to compute the version for the project. If unspecified,
	// defaults to using the git project versioner (refer to the "projectversioner/git" package for details on the
	// implementation of this operation).
	ProjectVersioner *ProjectVersionConfig `yaml:"project-versioner,omitempty"`

	// Exclude matches the paths to exclude when determining the projects to build.
	Exclude matcher.NamesPathsCfg `yaml:"exclude,omitempty"`
}

func UpgradeConfig(
	cfgBytes []byte,
	projectVersionerFactory distgo.ProjectVersionerFactory,
	disterFactory distgo.DisterFactory,
	dockerBuilderFactory distgo.DockerBuilderFactory,
	publisherFactory distgo.PublisherFactory) ([]byte, error) {

	var cfg ProjectConfig
	if err := yaml.UnmarshalStrict(cfgBytes, &cfg); err != nil {
		return nil, errors.Wrapf(err, "failed to unmarshal dist-plugin v0 configuration")
	}

	changed := false
	projectVerionerChanged, err := upgradeProjectVersioner(&cfg, projectVersionerFactory)
	if err != nil {
		return nil, err
	}
	changed = changed || projectVerionerChanged

	assetsChanged, err := upgradeAssets(&cfg, disterFactory, dockerBuilderFactory, publisherFactory)
	if err != nil {
		return nil, err
	}
	changed = changed || assetsChanged

	if !changed {
		return cfgBytes, nil
	}
	upgradedBytes, err := yaml.Marshal(cfg)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to marshal dist-plugin v0 configuration")
	}
	return upgradedBytes, nil
}

// upgradeProjectVersioner upgrades the project versioner for the provided configuration. Returns true if any changes
// were made by the upgrade. If any upgrade operations are performed, the provided configuration is modified directly.
func upgradeProjectVersioner(cfg *ProjectConfig, projectVersionerFactory distgo.ProjectVersionerFactory) (changed bool, rErr error) {
	if cfg.ProjectVersioner == nil {
		return false, nil
	}

	upgrader, err := projectVersionerFactory.ConfigUpgrader(cfg.ProjectVersioner.Type)
	if err != nil {
		return false, errors.Wrapf(err, "failed to upgrade project versioner of type %q", cfg.ProjectVersioner.Type)
	}
	originalCfgBytes, err := yaml.Marshal(cfg.ProjectVersioner.Config)
	if err != nil {
		return false, errors.Wrapf(err, "failed to marshal configuration for project versioner of type %q", cfg.ProjectVersioner.Type)
	}
	upgradedCfgBytes, err := upgrader.UpgradeConfig(originalCfgBytes)
	if err != nil {
		return false, errors.Wrapf(err, "failed to upgrade configuration for project versioner of type %q", cfg.ProjectVersioner.Type)
	}

	if bytes.Equal(originalCfgBytes, upgradedCfgBytes) {
		// upgrade was a no-op: do not modify configuration and continue
		return false, nil
	}

	var yamlRep yaml.MapSlice
	if err := yaml.Unmarshal(upgradedCfgBytes, &yamlRep); err != nil {
		return false, errors.Wrapf(err, "failed to unmarshal YAML of upgraded configuration for project versioner of type %q", cfg.ProjectVersioner.Type)
	}

	cfg.ProjectVersioner.Config = yamlRep
	return true, nil
}

// upgradeAssets upgrades the assets for the provided configuration. Returns true if any upgrade operations were
// performed. If any upgrade operations were performed, the provided configuration is modified directly.
func upgradeAssets(
	cfg *ProjectConfig,
	disterFactory distgo.DisterFactory,
	dockerBuilderFactory distgo.DockerBuilderFactory,
	publisherFactory distgo.PublisherFactory) (changed bool, rErr error) {

	defaultsChanged, err := upgradeProductAssets(&cfg.ProductDefaults, disterFactory, dockerBuilderFactory, publisherFactory)
	if err != nil {
		return false, errors.Wrapf(err, "failed to upgrade assets for product defaults")
	}
	changed = changed || defaultsChanged

	for k, v := range cfg.Products {
		v := v
		currProductChanged, err := upgradeProductAssets(&v, disterFactory, dockerBuilderFactory, publisherFactory)
		if err != nil {
			return false, errors.Wrapf(err, "failed to upgrade assets for product %q", k)
		}
		changed = changed || currProductChanged
		cfg.Products[k] = v
	}
	return changed, nil
}

func upgradeProductAssets(
	cfg *ProductConfig,
	disterFactory distgo.DisterFactory,
	dockerBuilderFactory distgo.DockerBuilderFactory,
	publisherFactory distgo.PublisherFactory) (changed bool, rErr error) {

	// upgrade dister assets
	if cfg.Dist != nil && cfg.Dist.Disters != nil {
		var sortedDistIDs []distgo.DistID
		for k := range *cfg.Dist.Disters {
			sortedDistIDs = append(sortedDistIDs, k)
		}
		sort.Sort(distgo.ByDistID(sortedDistIDs))

		for _, distID := range sortedDistIDs {
			dister := (*cfg.Dist.Disters)[distID]
			if dister.Config == nil {
				continue
			}

			upgrader, err := disterFactory.ConfigUpgrader(*dister.Type)
			if err != nil {
				return false, errors.Wrapf(err, "failed to upgrade dist %s of type %q", distID, *dister.Type)
			}
			assetCfgBytes, err := yaml.Marshal(*dister.Config)
			if err != nil {
				return false, errors.Wrapf(err, "failed to marshal configuration for dist %s of type %q", distID, *dister.Type)
			}

			upgradedBytes, err := upgrader.UpgradeConfig(assetCfgBytes)
			if err != nil {
				return false, errors.Wrapf(err, "failed to upgrade configuration for dist %s of type %q", distID, *dister.Type)
			}

			if bytes.Equal(assetCfgBytes, upgradedBytes) {
				// upgrade was a no-op: do not modify configuration and continue
				continue
			}
			changed = true

			var yamlRep yaml.MapSlice
			if err := yaml.Unmarshal(upgradedBytes, &yamlRep); err != nil {
				return false, errors.Wrapf(err, "failed to unmarshal YAML of upgraded configuration for dist %s of type %q", distID, *dister.Type)
			}

			dister.Config = &yamlRep
			(*cfg.Dist.Disters)[distID] = dister
		}
	}

	// upgrade docker builder assets
	if cfg.Docker != nil && cfg.Docker.DockerBuildersConfig != nil {
		var sortedDockerIDs []distgo.DockerID
		for k := range *cfg.Docker.DockerBuildersConfig {
			sortedDockerIDs = append(sortedDockerIDs, k)
		}
		sort.Sort(distgo.ByDockerID(sortedDockerIDs))

		for _, dockerID := range sortedDockerIDs {
			dockerBuilder := (*cfg.Docker.DockerBuildersConfig)[dockerID]
			if dockerBuilder.Config == nil {
				continue
			}

			upgrader, err := dockerBuilderFactory.ConfigUpgrader(*dockerBuilder.Type)
			if err != nil {
				return false, errors.Wrapf(err, "failed to upgrade docker builder %s of type %q", dockerID, *dockerBuilder.Type)
			}
			assetCfgBytes, err := yaml.Marshal(*dockerBuilder.Config)
			if err != nil {
				return false, errors.Wrapf(err, "failed to upgrade docker builder %s of type %q", dockerID, *dockerBuilder.Type)
			}

			upgradedBytes, err := upgrader.UpgradeConfig(assetCfgBytes)
			if err != nil {
				return false, errors.Wrapf(err, "failed to upgrade docker builder %s of type %q", dockerID, *dockerBuilder.Type)
			}

			if bytes.Equal(assetCfgBytes, upgradedBytes) {
				// upgrade was a no-op: do not modify configuration and continue
				continue
			}
			changed = true

			var yamlRep yaml.MapSlice
			if err := yaml.Unmarshal(upgradedBytes, &yamlRep); err != nil {
				return false, errors.Wrapf(err, "failed to unmarshal YAML of upgraded configuration for dist %s of type %q", dockerID, *dockerBuilder.Type)
			}

			dockerBuilder.Config = &yamlRep
			(*cfg.Docker.DockerBuildersConfig)[dockerID] = dockerBuilder
		}
	}

	// upgrade publisher assets
	if cfg.Publish != nil && cfg.Publish.PublishInfo != nil {
		var sortedPublisherTypeIDs []distgo.PublisherTypeID
		for k := range *cfg.Publish.PublishInfo {
			sortedPublisherTypeIDs = append(sortedPublisherTypeIDs, k)
		}
		sort.Sort(distgo.ByPublisherTypeID(sortedPublisherTypeIDs))

		for _, publisherTypeID := range sortedPublisherTypeIDs {
			publisher := (*cfg.Publish.PublishInfo)[publisherTypeID]
			if publisher.Config == nil {
				continue
			}

			upgrader, err := publisherFactory.ConfigUpgrader(string(publisherTypeID))
			if err != nil {
				return false, errors.Wrapf(err, "failed to upgrade publisher %q", publisherTypeID)
			}
			assetCfgBytes, err := yaml.Marshal(*publisher.Config)
			if err != nil {
				return false, errors.Wrapf(err, "failed to upgrade publisher %q", publisherTypeID)
			}

			upgradedBytes, err := upgrader.UpgradeConfig(assetCfgBytes)
			if err != nil {
				return false, errors.Wrapf(err, "failed to upgrade publisher %q", publisherTypeID)
			}

			if bytes.Equal(assetCfgBytes, upgradedBytes) {
				// upgrade was a no-op: do not modify configuration and continue
				continue
			}
			changed = true

			var yamlRep yaml.MapSlice
			if err := yaml.Unmarshal(upgradedBytes, &yamlRep); err != nil {
				return false, errors.Wrapf(err, "failed to unmarshal YAML of upgraded configuration for publisher %q", publisherTypeID)
			}

			publisher.Config = &yamlRep
			(*cfg.Publish.PublishInfo)[publisherTypeID] = publisher
		}
	}
	return changed, nil
}
