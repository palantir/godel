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

package pluginapitester

import (
	"bytes"

	"github.com/pkg/errors"

	"github.com/palantir/godel/framework/artifactresolver"
	"github.com/palantir/godel/framework/godel/config"
	"github.com/palantir/godel/framework/godellauncher"
	"github.com/palantir/godel/framework/internal/pathsinternal"
	"github.com/palantir/godel/framework/internal/pluginsinternal"
	"github.com/palantir/godel/framework/plugins"
	"github.com/palantir/godel/pkg/osarch"
)

type PluginProvider interface {
	PluginFilePath() string
}

type filePluginProvider struct {
	pluginPath string
}

func (p *filePluginProvider) PluginFilePath() string {
	return p.pluginPath
}

func NewPluginProvider(pluginPath string) PluginProvider {
	return &filePluginProvider{
		pluginPath: pluginPath,
	}
}

func NewPluginProviderFromLocator(pluginLocator, pluginResolver string) (PluginProvider, error) {
	lwrConfig := config.LocatorWithResolverConfig{
		Locator: config.ToLocatorConfig(config.LocatorConfig{
			ID: pluginLocator,
		}),
		Resolver: pluginResolver,
	}
	lwrParam, err := lwrConfig.ToParam()
	if err != nil {
		return nil, err
	}
	pluginPath, err := resolvePlugin(lwrParam)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to resolve plugin")
	}
	return NewPluginProvider(pluginPath), nil
}

// resolvePlugin resolves the plugin with the provided locator and returns the path to the resolved plugin.
func resolvePlugin(pluginLocator artifactresolver.LocatorWithResolverParam) (string, error) {
	buf := &bytes.Buffer{}
	if _, _, err := plugins.LoadPluginsTasks(godellauncher.PluginsParam{
		Plugins: []godellauncher.SinglePluginParam{
			{
				LocatorWithResolverParam: pluginLocator,
			},
		},
	}, buf); err != nil {
		return "", errors.Wrapf(err, "failed to load plugin task. Output: %s", buf.String())
	}

	pluginsDir, _, _, err := pathsinternal.ResourceDirs()
	if err != nil {
		return "", errors.Wrapf(err, "failed to determine plugin directory")
	}
	return pathsinternal.PluginPath(pluginsDir, pluginLocator.LocatorWithChecksums.Locator), nil
}

type AssetProvider interface {
	AssetFilePath() string
}

type fileAssetProvider struct {
	asssetPath string
}

func (p *fileAssetProvider) AssetFilePath() string {
	return p.asssetPath
}

func NewAssetProvider(assetPath string) AssetProvider {
	return &fileAssetProvider{
		asssetPath: assetPath,
	}
}

func NewAssetProviderFromLocator(assetLocator, assetResolver string) (AssetProvider, error) {
	lwrConfig := config.LocatorWithResolverConfig{
		Locator: config.ToLocatorConfig(config.LocatorConfig{
			ID: assetLocator,
		}),
		Resolver: assetResolver,
	}
	lwrParam, err := lwrConfig.ToParam()
	if err != nil {
		return nil, err
	}
	_, assetsDir, downloadsDir, err := pathsinternal.ResourceDirs()
	if err != nil {
		return nil, err
	}
	resolver, err := artifactresolver.NewTemplateResolver(assetResolver)
	if err != nil {
		return nil, err
	}

	outputBuf := &bytes.Buffer{}
	if _, err := pluginsinternal.ResolveAssets(
		assetsDir,
		downloadsDir,
		[]artifactresolver.LocatorWithResolverParam{
			lwrParam,
		},
		osarch.Current(),
		[]artifactresolver.Resolver{
			resolver,
		},
		outputBuf,
	); err != nil {
		return nil, errors.Wrapf(err, "failed to resolve assets:\n%s", outputBuf.String())
	}
	return NewAssetProvider(pathsinternal.PluginPath(assetsDir, lwrParam.LocatorWithChecksums.Locator)), nil
}
