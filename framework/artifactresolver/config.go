package artifactresolver

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"

	"github.com/palantir/godel/apps/distgo/pkg/osarch"
)

type LocatorWithResolverParam struct {
	LocatorWithChecksums LocatorParam
	Resolver             Resolver
}

type LocatorWithResolverConfig struct {
	Locator  LocatorConfig `yaml:"locator"`
	Resolver string        `yaml:"resolver"`
}

func (c *LocatorWithResolverConfig) ToParam() (LocatorWithResolverParam, error) {
	locator, err := c.Locator.ToParam()
	if err != nil {
		return LocatorWithResolverParam{}, errors.Wrapf(err, "invalid locator")
	}
	var resolver Resolver
	if c.Resolver != "" {
		resolverVal, err := NewTemplateResolver(c.Resolver)
		if err != nil {
			return LocatorWithResolverParam{}, errors.Wrapf(err, "invalid resolver")
		}
		resolver = resolverVal
	}
	return LocatorWithResolverParam{
		LocatorWithChecksums: locator,
		Resolver:             resolver,
	}, nil
}

// ConfigProviderLocatorWithResolverConfig is the configuration for a locator with resolver for a configuration
// provider. It differs from a LocatorWithResolverConfig in that the locator is a ConfigProviderLocatorConfig rather
// than a LocatorConfig.
type ConfigProviderLocatorWithResolverConfig struct {
	Locator  ConfigProviderLocatorConfig `yaml:"locator"`
	Resolver string                      `yaml:"resolver"`
}

// ToParam converts the configuration into a LocatorWithResolverParam. Any checksums that exist are put in a map where
// the key is the current OS/Arch.
func (c *ConfigProviderLocatorWithResolverConfig) ToParam() (LocatorWithResolverParam, error) {
	locatorCfg, err := c.Locator.ToLocatorConfig()
	if err != nil {
		return LocatorWithResolverParam{}, err
	}
	cfg := LocatorWithResolverConfig{
		Locator:  locatorCfg,
		Resolver: c.Resolver,
	}
	return cfg.ToParam()
}

type LocatorParam struct {
	Locator
	Checksums map[osarch.OSArch]string
}

type LocatorConfig struct {
	ID        string            `yaml:"id"`
	Checksums map[string]string `yaml:"checksums"`
}

func (c *LocatorConfig) ToParam() (LocatorParam, error) {
	parts := strings.Split(c.ID, ":")
	if len(parts) != 3 {
		return LocatorParam{}, errors.Errorf("locator ID must consist of 3 colon-delimited components ([group]:[product]:[version]), but had %d: %q", len(parts), c.ID)
	}
	var checksums map[osarch.OSArch]string
	if c.Checksums != nil {
		checksums = make(map[osarch.OSArch]string)
		for k, v := range c.Checksums {
			osArchKey, err := osarch.New(k)
			if err != nil {
				return LocatorParam{}, errors.Wrapf(err, "invalid OSArch specified in checksum key for %s", c.ID)
			}
			checksums[osArchKey] = v
		}
	}
	param := LocatorParam{
		Locator: Locator{
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
type ConfigProviderLocatorConfig struct {
	ID       string `yaml:"id"`
	Checksum string `yaml:"checksum"`
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

type Locator struct {
	Group   string
	Product string
	Version string
}

func (l Locator) String() string {
	return fmt.Sprintf("%s:%s:%s", l.Group, l.Product, l.Version)
}
