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
	"io"
	"os"
	"path"
	"sort"
	"strings"

	"github.com/palantir/pkg/specdir"
	"github.com/pkg/errors"

	"github.com/palantir/godel/apps/distgo/pkg/osarch"
	"github.com/palantir/godel/framework/builtintasks/installupdate/layout"
	"github.com/palantir/godel/framework/godellauncher"
	"github.com/palantir/godel/framework/pluginapi"
)

const (
	indentSpaces = 4
)

// pluginInfoWithAssets bundles a pluginapi.Info with the locators of all the assets specified for it.
type pluginInfoWithAssets struct {
	PluginInfo pluginapi.Info
	Assets     []locator
}

// LoadPluginsTasks returns all of the tasks defined by the plugins in the specified parameters. Does the following:
//
// * Resolves all of the plugins defined in the provided params for the runtime environment's OS/Architecture into the
//   gödel home plugins and downloads directories.
// * Verifies that all of the resolved plugins are valid and compatible with each other (for example, ensures that
//   multiple plugins do not provide the same task).
// * Creates runnable godellauncher.Task tasks for all of the plugins.
//
// Returns all of the tasks provided by the plugins in the provided parameters.
func LoadPluginsTasks(cfg godellauncher.PluginsConfig, stdout io.Writer) ([]godellauncher.Task, error) {
	params, err := projectParamsFromConfig(cfg)
	if err != nil {
		return nil, err
	}

	gödelHomeSpecDir, err := layout.GodelHomeSpecDir(specdir.Create)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to create gödel home directory")
	}
	pluginsDir := gödelHomeSpecDir.Path(layout.PluginsDir)
	assetsDir := gödelHomeSpecDir.Path(layout.AssetsDir)
	downloadsDir := gödelHomeSpecDir.Path(layout.DownloadsDir)

	plugins, err := resolvePlugins(pluginsDir, assetsDir, downloadsDir, osarch.Current(), params, stdout)
	if err != nil {
		return nil, err
	}
	if err := verifyPluginCompatibility(plugins); err != nil {
		return nil, err
	}

	var sortedPluginLocators []locator
	for k := range plugins {
		sortedPluginLocators = append(sortedPluginLocators, k)
	}
	sortLocators(sortedPluginLocators)

	var tasks []godellauncher.Task
	for _, pluginLoc := range sortedPluginLocators {
		pluginExecPath := pluginPath(pluginsDir, pluginLoc)
		pluginInfoWithAssets := plugins[pluginLoc]

		var assetPaths []string
		for _, assetLoc := range pluginInfoWithAssets.Assets {
			assetPaths = append(assetPaths, pluginPath(assetsDir, assetLoc))
		}
		tasks = append(tasks, pluginInfoWithAssets.PluginInfo.Tasks(pluginExecPath, assetPaths)...)
	}
	return tasks, nil
}

// resolvePlugins resolves all of the plugins defined in the provided params for the specified osArch using the provided
// plugins and downloads directories. Returns a map that contains all of the information for the valid plugins. If
// errors were encountered while trying to resolve plugins, returns an error that summarizes the errors.
//
// For each plugin defined in the parameters:
//
// * If a file does not exist in the expected location in the plugins directory, resolve it
//   * If the configuration specifies a custom resolver for the plugin, use it to resolve the plugin TGZ into the
//     downloads directory
//   * Otherwise, if default resolvers are specified in the parameters, try to resolve the plugin TGZ into the
//     downloads directory from each of them in order
//   * If the plugin TGZ cannot be resolved, return an error
//   * If the plugin TGZ was resolved, unpack the content of the TGZ (which must contain a single file) into the
//     expected location in the plugins directory
// * If the configuration specifies a checksum for the plugin and the specified osArch, verify that the checksum of
//   the plugin in the plugins directory matches the specified checksum
// * Invoke the plugin info command (specified by the InfoCommandName constant) on the plugin and parse the output
//   as the plugin information
// * If the plugin specifies assets, resolve all of the assets
//   * Asset resolution uses a process that is analogous to plugin resolution, but performs it in the assets directory
func resolvePlugins(pluginsDir, assetsDir, downloadsDir string, osArch osarch.OSArch, param projectParams, stdout io.Writer) (map[locator]pluginInfoWithAssets, error) {
	plugins := make(map[locator]pluginInfoWithAssets)
	pluginErrors := make(map[locator]error)
	for _, currPlugin := range param.Plugins {
		currPluginLocator, ok := resolveAndVerify(
			currPlugin.locatorWithResolverParam,
			pluginErrors,
			pluginsDir,
			downloadsDir,
			param.DefaultResolvers,
			osArch,
			stdout,
		)
		if !ok {
			continue
		}
		info, err := pluginapi.InfoFromPlugin(path.Join(pluginsDir, pluginFileName(currPluginLocator)))
		if err != nil {
			pluginErrors[currPluginLocator] = errors.Wrapf(err, "failed to get plugin info for plugin %+v", currPluginLocator)
			continue
		}

		// plugin has been successfully resolved: resolve assets for plugin
		assetInfoMap, err := resolveAssets(assetsDir, downloadsDir, currPlugin.Assets, osArch, param, stdout)
		if err != nil {
			pluginErrors[currPluginLocator] = errors.Wrapf(err, "failed to get asset(s) for plugin %+v", currPluginLocator)
			continue
		}

		plugins[currPluginLocator] = pluginInfoWithAssets{
			PluginInfo: info,
			Assets:     assetInfoMap,
		}
	}

	if len(pluginErrors) == 0 {
		return plugins, nil
	}

	// encountered errors: summarize and return
	var sortedKeys []locator
	for k := range pluginErrors {
		sortedKeys = append(sortedKeys, k)
	}
	sortLocators(sortedKeys)

	errStringsParts := []string{fmt.Sprintf("failed to resolve %d plugin(s):", len(pluginErrors))}
	for _, k := range sortedKeys {
		errStringsParts = append(errStringsParts, pluginErrors[k].Error())
	}
	return nil, errors.New(strings.Join(errStringsParts, "\n"+strings.Repeat(" ", indentSpaces)))
}

func resolveAssets(assetsDir, downloadsDir string, assetParams []locatorWithResolverParam, osArch osarch.OSArch, param projectParams, stdout io.Writer) ([]locator, error) {
	if len(assetParams) == 0 {
		return nil, nil
	}

	var assets []locator
	assetErrors := make(map[locator]error)
	for _, currAsset := range assetParams {
		currAssetLocator, ok := resolveAndVerify(
			currAsset,
			assetErrors,
			assetsDir,
			downloadsDir,
			param.DefaultResolvers,
			osArch,
			stdout,
		)
		if !ok {
			continue
		}
		assets = append(assets, currAssetLocator)
	}
	sortLocators(assets)

	if len(assetErrors) == 0 {
		return assets, nil
	}

	// encountered errors: summarize and return
	errStringsParts := []string{fmt.Sprintf("failed to resolve %d asset(s):", len(assetErrors))}
	for _, k := range assets {
		errStringsParts = append(errStringsParts, assetErrors[k].Error())
	}
	return nil, errors.New(strings.Join(errStringsParts, "\n"+strings.Repeat(" ", indentSpaces)))
}

func resolveAndVerify(
	currArtifact locatorWithResolverParam,
	artifactErrors map[locator]error,
	dstBaseDir, downloadsDir string,
	defaultResolvers []resolver,
	osArch osarch.OSArch,
	stdout io.Writer) (currLocator locator, ok bool) {

	currLocator = currArtifact.LocatorWithChecksums.locator
	currDstPath := path.Join(dstBaseDir, pluginFileName(currLocator))

	if _, err := os.Stat(currDstPath); os.IsNotExist(err) {
		tgzDstPath := path.Join(downloadsDir, pluginFileName(currLocator)+".tgz")
		if err := resolvePluginTGZ(currArtifact, defaultResolvers, osArch, tgzDstPath, stdout); err != nil {
			artifactErrors[currLocator] = err
			return currLocator, false
		}

		if err := func() (rErr error) {
			pluginFile, err := os.OpenFile(currDstPath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
			if err != nil {
				return errors.Wrapf(err, "failed to create file %s", currDstPath)
			}
			defer func() {
				if err := pluginFile.Close(); err != nil {
					rErr = errors.Wrapf(err, "failed to close file %s", currDstPath)
				}
			}()

			tgzFile, err := os.Open(tgzDstPath)
			if err != nil {
				return errors.Wrapf(err, "failed to open %s for reading", tgzDstPath)
			}

			if err := copyPluginTGZContent(pluginFile, tgzFile); err != nil {
				return err
			}
			return nil
		}(); err != nil {
			artifactErrors[currLocator] = errors.Wrapf(err, "failed to extract artifact from archive into destination")
			return currLocator, false
		}
	}

	if wantChecksum, ok := currArtifact.LocatorWithChecksums.Checksums[osArch]; ok {
		gotChecksum, err := sha256ChecksumFile(currDstPath)
		if err != nil {
			artifactErrors[currLocator] = errors.Wrapf(err, "failed to compute checksum for plugin")
			return currLocator, false
		}
		if gotChecksum != wantChecksum {
			artifactErrors[currLocator] = errors.Errorf("failed to verify checksum for %s: want %s, got %s", currDstPath, wantChecksum, gotChecksum)
			return currLocator, false
		}
	}
	return currLocator, true
}

// Verifies that the plugins in the provided map are compatible with one another. Specifically, ensures that there is at
// most 1 version of a given plugin (a locator with a given {group, product} pair) and that there are no conflicts
// between tasks provided by the plugins.
func verifyPluginCompatibility(plugins map[locator]pluginInfoWithAssets) error {
	// map from a plugin locator to the locators to all of the plugins that they conflict with and the error that
	// describes the conflict.
	conflicts := make(map[locator]map[locator]error)
	for currPlugin := range plugins {
		currConflicts := verifySinglePluginCompatibility(currPlugin, plugins)
		if len(currConflicts) == 0 {
			continue
		}
		conflicts[currPlugin] = currConflicts
	}

	if len(conflicts) == 0 {
		return nil
	}

	var sortedOuterKeys []locator
	for k := range conflicts {
		sortedOuterKeys = append(sortedOuterKeys, k)
	}
	sortLocators(sortedOuterKeys)

	errString := fmt.Sprintf("%d plugins had compatibility issues:", len(conflicts))
	for _, k := range sortedOuterKeys {
		errString += fmt.Sprintf("\n%s%s:", strings.Repeat(" ", indentSpaces), k.String())

		var sortedInnerKeys []locator
		for innerK := range conflicts[k] {
			sortedInnerKeys = append(sortedInnerKeys, innerK)
		}
		sortLocators(sortedInnerKeys)

		for _, innerK := range sortedInnerKeys {
			errString += fmt.Sprintf("\n%s%s", strings.Repeat(" ", indentSpaces*2), conflicts[k][innerK].Error())
		}
	}
	return errors.New(errString)
}

func verifySinglePluginCompatibility(plugin locator, plugins map[locator]pluginInfoWithAssets) map[locator]error {
	errs := make(map[locator]error)
	for otherPlugin, otherPluginInfo := range plugins {
		if otherPlugin == plugin {
			continue
		}
		if otherPlugin.Group == plugin.Group && otherPlugin.Product == plugin.Product {
			errs[otherPlugin] = fmt.Errorf("different version of the same plugin")
			continue
		}

		currPluginInfo := plugins[plugin]
		var currPluginTasks []string
		for _, currPluginTask := range currPluginInfo.PluginInfo.Tasks("", nil) {
			currPluginTasks = append(currPluginTasks, currPluginTask.Name)
		}
		sort.Strings(currPluginTasks)

		otherPluginTasks := make(map[string]struct{})
		for _, task := range otherPluginInfo.PluginInfo.Tasks("", nil) {
			otherPluginTasks[task.Name] = struct{}{}
		}

		var commonTasks []string
		for _, currPluginTask := range currPluginTasks {
			if _, ok := otherPluginTasks[currPluginTask]; !ok {
				continue
			}
			commonTasks = append(commonTasks, currPluginTask)
		}
		if len(commonTasks) != 0 {
			errs[otherPlugin] = fmt.Errorf("provides conflicting tasks: %v", commonTasks)
			continue
		}
	}
	return errs
}

func sortLocators(locs []locator) {
	sort.Slice(locs, func(i, j int) bool {
		return locs[i].String() < locs[j].String()
	})
}

func pluginPath(pluginDir string, locator locator) string {
	return path.Join(pluginDir, pluginFileName(locator))
}

func pluginFileName(locator locator) string {
	return fmt.Sprintf("%s-%s-%s", locator.Group, locator.Product, locator.Version)
}
