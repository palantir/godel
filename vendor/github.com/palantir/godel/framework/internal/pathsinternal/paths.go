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

package pathsinternal

import (
	"fmt"
	"path"

	"github.com/palantir/pkg/specdir"
	"github.com/pkg/errors"

	"github.com/palantir/godel/framework/artifactresolver"
	"github.com/palantir/godel/framework/builtintasks/installupdate/layout"
)

func PluginPath(pluginDir string, locator artifactresolver.Locator) string {
	return path.Join(pluginDir, PluginFileName(locator))
}

func PluginFileName(locator artifactresolver.Locator) string {
	return fmt.Sprintf("%s-%s-%s", locator.Group, locator.Product, locator.Version)
}

func ConfigProviderFileName(locator artifactresolver.Locator) string {
	return PluginFileName(locator) + ".yml"
}

func ResourceDirs() (pluginsDir string, assetsDir string, downloadsDir string, rErr error) {
	godelHomeSpecDir, err := layout.GodelHomeSpecDir(specdir.Create)
	if err != nil {
		return "", "", "", errors.Wrapf(err, "failed to create gödel home directory")
	}
	return godelHomeSpecDir.Path(layout.PluginsDir), godelHomeSpecDir.Path(layout.AssetsDir), godelHomeSpecDir.Path(layout.DownloadsDir), nil
}
