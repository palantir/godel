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

package plugins

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"

	"github.com/palantir/godel/apps/distgo/pkg/osarch"
	"github.com/palantir/godel/framework/godellauncher"
)

type projectParams struct {
	DefaultResolvers []resolver
	Plugins          []singlePluginParam
}

func projectParamsFromConfig(cfg godellauncher.PluginsConfig) (projectParams, error) {
	var resolvers []resolver
	for _, resolverTmpl := range cfg.DefaultResolvers {
		resolver, err := newTemplateResolver(resolverTmpl)
		if err != nil {
			return projectParams{}, err
		}
		resolvers = append(resolvers, resolver)
	}
	var singlePluginParams []singlePluginParam
	for _, p := range cfg.Plugins {
		singlePluginParam, err := singlePluginParamFromConfig(p)
		if err != nil {
			return projectParams{}, err
		}
		singlePluginParams = append(singlePluginParams, singlePluginParam)
	}
	return projectParams{
		DefaultResolvers: resolvers,
		Plugins:          singlePluginParams,
	}, nil
}

type singlePluginParam struct {
	locatorWithResolverParam
	Assets []locatorWithResolverParam
}

func singlePluginParamFromConfig(c godellauncher.SinglePluginConfig) (singlePluginParam, error) {
	pluginLocWithResolverParam, err := locatorWithResolverParamFromConfig(c.LocatorWithResolverConfig)
	if err != nil {
		return singlePluginParam{}, err
	}

	var assetParams []locatorWithResolverParam
	for _, asset := range c.Assets {
		assetLocWithResolverParam, err := locatorWithResolverParamFromConfig(asset)
		if err != nil {
			return singlePluginParam{}, err
		}
		assetParams = append(assetParams, assetLocWithResolverParam)
	}

	return singlePluginParam{
		locatorWithResolverParam: pluginLocWithResolverParam,
		Assets: assetParams,
	}, nil
}

type locatorWithResolverParam struct {
	LocatorWithChecksums locatorWithChecksumsParam
	Resolver             resolver
}

func locatorWithResolverParamFromConfig(c godellauncher.LocatorWithResolverConfig) (locatorWithResolverParam, error) {
	locator, err := locatorWithChecksumsParamFromConfig(c.Locator)
	if err != nil {
		return locatorWithResolverParam{}, errors.Wrapf(err, "invalid locator")
	}
	var resolver resolver
	if c.Resolver != "" {
		var err error
		resolver, err = newTemplateResolver(c.Resolver)
		if err != nil {
			return locatorWithResolverParam{}, errors.Wrapf(err, "invalid resolver")
		}
	}
	return locatorWithResolverParam{
		LocatorWithChecksums: locator,
		Resolver:             resolver,
	}, nil
}

type locatorWithChecksumsParam struct {
	locator
	Checksums map[osarch.OSArch]string
}

type locator struct {
	Group   string
	Product string
	Version string
}

func (l locator) String() string {
	return fmt.Sprintf("%s:%s:%s", l.Group, l.Product, l.Version)
}

func locatorWithChecksumsParamFromConfig(cfg godellauncher.LocatorConfig) (locatorWithChecksumsParam, error) {
	parts := strings.Split(cfg.ID, ":")
	if len(parts) != 3 {
		return locatorWithChecksumsParam{}, errors.Errorf("locator ID must consist of 3 colon-delimited components ([group]:[product]:[version]), but had %d: %q", len(parts), cfg.ID)
	}
	var checksums map[osarch.OSArch]string
	if cfg.Checksums != nil {
		checksums = make(map[osarch.OSArch]string)
		for k, v := range cfg.Checksums {
			osArchKey, err := osarch.New(k)
			if err != nil {
				return locatorWithChecksumsParam{}, errors.Wrapf(err, "invalid OSArch specified in checksum key for %s", cfg.ID)
			}
			checksums[osArchKey] = v
		}
	}
	param := locatorWithChecksumsParam{
		locator: locator{
			Group:   parts[0],
			Product: parts[1],
			Version: parts[2],
		},
		Checksums: checksums,
	}
	return param, nil
}
