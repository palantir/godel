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

type LocatorWithResolverConfig struct {
	Locator  LocatorConfig `yaml:"locator,omitempty"`
	Resolver string        `yaml:"resolver,omitempty"`
}

// ConfigProviderLocatorWithResolverConfig is the configuration for a locator with resolver for a configuration
// provider. It differs from a LocatorWithResolverConfig in that the locator is a ConfigProviderLocatorConfig rather
// than a LocatorConfig.
type ConfigProviderLocatorWithResolverConfig struct {
	Locator  ConfigProviderLocatorConfig `yaml:"locator,omitempty"`
	Resolver string                      `yaml:"resolver,omitempty"`
}

type LocatorConfig struct {
	ID        string            `yaml:"id,omitempty"`
	Checksums map[string]string `yaml:"checksums,omitempty"`
}

// ConfigProviderLocatorConfig is the configuration for a locator for a configuration provider. It differs from a
// LocatorConfig in that only a single checksum can be specified.
type ConfigProviderLocatorConfig struct {
	ID       string `yaml:"id,omitempty"`
	Checksum string `yaml:"checksum,omitempty"`
}
