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
	"strings"

	"github.com/pkg/errors"

	"github.com/palantir/godel/v2/framework/artifactresolver"
	v0 "github.com/palantir/godel/v2/framework/godel/config/internal/v0"
	"github.com/palantir/godel/v2/pkg/osarch"
)

type LocatorWithResolverConfig v0.LocatorWithResolverConfig

func ToLocatorWithResolverConfig(in LocatorWithResolverConfig) v0.LocatorWithResolverConfig {
	return v0.LocatorWithResolverConfig(in)
}

func ToLocatorWithResolverConfigs(in []LocatorWithResolverConfig) []v0.LocatorWithResolverConfig {
	if in == nil {
		return nil
	}
	out := make([]v0.LocatorWithResolverConfig, len(in))
	for i, v := range in {
		out[i] = ToLocatorWithResolverConfig(v)
	}
	return out
}

func (c *LocatorWithResolverConfig) ToParam() (artifactresolver.LocatorWithResolverParam, error) {
	locatorCfg := LocatorConfig(c.Locator)
	locator, err := locatorCfg.ToParam()
	if err != nil {
		return artifactresolver.LocatorWithResolverParam{}, errors.Wrapf(err, "invalid locator")
	}
	var resolver artifactresolver.Resolver
	if c.Resolver != "" {
		resolverVal, err := artifactresolver.NewTemplateResolver(c.Resolver)
		if err != nil {
			return artifactresolver.LocatorWithResolverParam{}, errors.Wrapf(err, "invalid resolver")
		}
		resolver = resolverVal
	}
	return artifactresolver.LocatorWithResolverParam{
		LocatorWithChecksums: locator,
		Resolver:             resolver,
	}, nil
}

// ConfigProviderLocatorWithResolverConfig is the configuration for a locator with resolver for a configuration
// provider. It differs from a LocatorWithResolverConfig in that the locator is a ConfigProviderLocatorConfig rather
// than a LocatorConfig.
type ConfigProviderLocatorWithResolverConfig v0.ConfigProviderLocatorWithResolverConfig

func ToConfigProviderLocatorWithResolverConfig(in ConfigProviderLocatorWithResolverConfig) v0.ConfigProviderLocatorWithResolverConfig {
	return v0.ConfigProviderLocatorWithResolverConfig(in)
}

// ToParam converts the configuration into a LocatorWithResolverParam. Any checksums that exist are put in a map where
// the key is the current OS/Arch.
func (c *ConfigProviderLocatorWithResolverConfig) ToParam() (artifactresolver.LocatorWithResolverParam, error) {
	providerLocatorCfg := ConfigProviderLocatorConfig(c.Locator)
	locatorCfg, err := providerLocatorCfg.ToLocatorConfig()
	if err != nil {
		return artifactresolver.LocatorWithResolverParam{}, err
	}
	cfg := LocatorWithResolverConfig{
		Locator:  v0.LocatorConfig(locatorCfg),
		Resolver: c.Resolver,
	}
	return cfg.ToParam()
}

type LocatorConfig v0.LocatorConfig

func ToLocatorConfig(in LocatorConfig) v0.LocatorConfig {
	return v0.LocatorConfig(in)
}

func (c *LocatorConfig) ToParam() (artifactresolver.LocatorParam, error) {
	parts := strings.Split(c.ID, ":")
	if len(parts) != 3 {
		return artifactresolver.LocatorParam{}, errors.Errorf("locator ID must consist of 3 colon-delimited components ([group]:[product]:[version]), but had %d: %q", len(parts), c.ID)
	}
	var checksums map[osarch.OSArch]string
	if c.Checksums != nil {
		checksums = make(map[osarch.OSArch]string)
		for k, v := range c.Checksums {
			osArchKey, err := osarch.New(k)
			if err != nil {
				return artifactresolver.LocatorParam{}, errors.Wrapf(err, "invalid OSArch specified in checksum key for %s", c.ID)
			}
			checksums[osArchKey] = v
		}
	}
	param := artifactresolver.LocatorParam{
		Locator: artifactresolver.Locator{
			Group:   parts[0],
			Product: parts[1],
			Version: parts[2],
		},
		Checksums: checksums,
	}
	return param, nil
}

// placeholder OS/Arch used for config provider checksums
var configProviderOSArch = osarch.Current()

// ConfigProviderLocatorConfig is the configuration for a locator for a configuration provider. It differs from a
// LocatorConfig in that only a single checksum can be specified.
type ConfigProviderLocatorConfig v0.ConfigProviderLocatorConfig

func ToConfigProviderLocatorConfig(in ConfigProviderLocatorConfig) v0.ConfigProviderLocatorConfig {
	return v0.ConfigProviderLocatorConfig(in)
}

// ToLocatorConfig translates the ConfigProviderLocatorConfig into a LocatorConfig where the checksum (if any exists) is
// keyed as the current OS/Arch.
func (c *ConfigProviderLocatorConfig) ToLocatorConfig() (LocatorConfig, error) {
	var checksums map[string]string
	if c.Checksum != "" {
		checksums = map[string]string{
			configProviderOSArch.String(): c.Checksum,
		}
	}
	return LocatorConfig{
		ID:        c.ID,
		Checksums: checksums,
	}, nil
}
