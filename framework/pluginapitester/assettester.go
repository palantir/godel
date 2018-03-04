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

// Package pluginapitester provides functions that simulate invoking a plugin from gödel. Can be used to test plugin
// implementations in plugin projects.
package pluginapitester

import (
	"bytes"
	"io"

	"github.com/pkg/errors"

	"github.com/palantir/godel/framework/artifactresolver"
	"github.com/palantir/godel/framework/godellauncher"
	"github.com/palantir/godel/framework/internal/pathsinternal"
	"github.com/palantir/godel/framework/plugins"
)

// RunAsset resolves the plugin with the provided locator and then runs it with the specified assets using RunPlugin.
// Can be used to test assets for which the plugin is published separately.
func RunAsset(
	pluginLocator artifactresolver.LocatorWithResolverParam,
	assetPaths []string,
	taskName string,
	args []string,
	projectDir string,
	debug bool,
	stdout io.Writer) (cleanup func(), rErr error) {

	pluginPath, err := resolvePlugin(pluginLocator)
	if err != nil {
		return func() {}, errors.Wrapf(err, "failed to resolve plugin")
	}
	return RunPlugin(pluginPath, assetPaths, taskName, args, projectDir, debug, stdout)
}

// resolvePlugin resolves the plugin with the provided locator and returns the path to the resolved plugin.
func resolvePlugin(pluginLocator artifactresolver.LocatorWithResolverParam) (string, error) {
	buf := &bytes.Buffer{}
	if _, err := plugins.LoadPluginsTasks(godellauncher.PluginsParam{
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
