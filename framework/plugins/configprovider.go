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
	"io/ioutil"
	"os"
	"path"
	"strings"

	"github.com/palantir/pkg/specdir"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"

	"github.com/palantir/godel/framework/artifactresolver"
	"github.com/palantir/godel/framework/builtintasks/installupdate/layout"
	"github.com/palantir/godel/framework/godel/config"
	"github.com/palantir/godel/framework/godellauncher"
	"github.com/palantir/godel/framework/internal/pathsinternal"
	"github.com/palantir/godel/framework/internal/pluginsinternal"
	"github.com/palantir/godel/pkg/osarch"
)

// LoadProvidedConfigurations returns all of the godellauncher.GodelConfig configurations provided by the specified
// params. Does the following:
//
// * Resolves all of the configuration providers defined in the provided params into the gödel home configs and
//   downloads directories.
// * Unmarshals all of the resolved configurations into godellauncher.TasksConfig structs.
//
// Returns all of the unmarshaled configurations.
func LoadProvidedConfigurations(taskConfigProvidersParam godellauncher.TasksConfigProvidersParam, stdout io.Writer) ([]config.TasksConfig, error) {
	gödelHomeSpecDir, err := layout.GodelHomeSpecDir(specdir.Create)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to create gödel home directory")
	}
	configsDir := gödelHomeSpecDir.Path(layout.ConfigsDir)
	downloadsDir := gödelHomeSpecDir.Path(layout.DownloadsDir)
	return resolveConfigProviders(configsDir, downloadsDir, taskConfigProvidersParam, stdout)
}

// resolveConfigProviders resolves all of the configurations provided by the specified parameters. Returns a slice that
// contains all of the resolved configurations. If errors were encountered while trying to resolve configurations,
// returns an error that summarizes the errors.
//
// For each configuration provider defined in the parameters:
//
// * If a file does not exist in the expected location in the configurations directory, resolve it
//   * If the configuration provider specifies a custom resolver, use it to resolve the configuration YML to the
//     expected location in the configurations directory
//   * Otherwise, if default resolvers are specified in the parameters, try to resolve the configuration YML to the
//     expected location from each of them in order
//   * If the configuration cannot be resolved, return an error
// * If the configuration specifies a checksum, verify that the checksum of the YML in the configurations directory
//   matches the specified checksum
// * Unmarshal the downloaded YML as godellauncher.TasksConfig
//   * If the unmarshal fails, return an error
func resolveConfigProviders(configsDir, downloadsDir string, taskConfigProvidersParam godellauncher.TasksConfigProvidersParam, stdout io.Writer) ([]config.TasksConfig, error) {
	var configs []config.TasksConfig
	providerErrors := make(map[artifactresolver.Locator]error)
	for _, currProvider := range taskConfigProvidersParam.ConfigProviders {
		currProviderLocator, ok := resolveAndVerifyConfigProvider(
			currProvider,
			providerErrors,
			configsDir,
			downloadsDir,
			taskConfigProvidersParam.DefaultResolvers,
			stdout,
		)
		if !ok {
			continue
		}

		tasksCfg, err := readConfigFromProvider(currProviderLocator, configsDir)
		if err != nil {
			providerErrors[currProviderLocator] = err
			continue
		}
		configs = append(configs, tasksCfg)
	}

	if len(providerErrors) == 0 {
		return configs, nil
	}

	// encountered errors: summarize and return
	var sortedKeys []artifactresolver.Locator
	for k := range providerErrors {
		sortedKeys = append(sortedKeys, k)
	}
	pluginsinternal.SortLocators(sortedKeys)

	errStringsParts := []string{fmt.Sprintf("failed to resolve %d configuration provider(s):", len(providerErrors))}
	for _, k := range sortedKeys {
		errStringsParts = append(errStringsParts, providerErrors[k].Error())
	}
	return nil, errors.New(strings.Join(errStringsParts, "\n"+strings.Repeat(" ", pluginsinternal.IndentSpaces)))
}

func resolveAndVerifyConfigProvider(
	currArtifact artifactresolver.LocatorWithResolverParam,
	artifactErrors map[artifactresolver.Locator]error,
	dstBaseDir, downloadsDir string,
	defaultResolvers []artifactresolver.Resolver,
	stdout io.Writer) (currLocator artifactresolver.Locator, ok bool) {

	currLocator = currArtifact.LocatorWithChecksums.Locator
	currDstPath := path.Join(dstBaseDir, pathsinternal.ConfigProviderFileName(currLocator))

	if _, err := os.Stat(currDstPath); os.IsNotExist(err) {
		downloadDstPath := path.Join(downloadsDir, pathsinternal.ConfigProviderFileName(currLocator))
		if err := artifactresolver.ResolveArtifact(currArtifact, defaultResolvers, osarch.Current(), downloadDstPath, artifactresolver.SHA256ChecksumFile, stdout); err != nil {
			artifactErrors[currLocator] = err
			return currLocator, false
		}
		if err := func() (rErr error) {
			cfgFile, err := os.OpenFile(currDstPath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
			if err != nil {
				return errors.Wrapf(err, "failed to create file %s", currDstPath)
			}
			defer func() {
				if err := cfgFile.Close(); err != nil {
					rErr = errors.Wrapf(err, "failed to close file %s", currDstPath)
				}
			}()

			downloadedFile, err := os.Open(downloadDstPath)
			if err != nil {
				return errors.Wrapf(err, "failed to open %s for reading", downloadDstPath)
			}
			if _, err := io.Copy(cfgFile, downloadedFile); err != nil {
				return err
			}
			return nil
		}(); err != nil {
			artifactErrors[currLocator] = errors.Wrapf(err, "failed to copy resolved artifact to destination")
			return currLocator, false
		}
	}

	if wantChecksum, ok := currArtifact.LocatorWithChecksums.Checksums[osarch.Current()]; ok {
		gotChecksum, err := artifactresolver.SHA256ChecksumFile(currDstPath)
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

func readConfigFromProvider(locator artifactresolver.Locator, configsDir string) (config.TasksConfig, error) {
	cfgPath := path.Join(configsDir, pathsinternal.ConfigProviderFileName(locator))

	cfgBytes, err := ioutil.ReadFile(cfgPath)
	if err != nil {
		return config.TasksConfig{}, errors.Wrapf(err, "failed to read %s", cfgPath)
	}

	var tasksCfg config.TasksConfig
	if err := yaml.Unmarshal(cfgBytes, &tasksCfg); err != nil {
		return config.TasksConfig{}, errors.Wrapf(err, "failed to unmarshal %q as godellauncher.GodelConfig", string(cfgBytes))
	}
	return tasksCfg, nil
}
