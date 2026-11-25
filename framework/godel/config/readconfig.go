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
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/palantir/godel/v2/framework/godellauncher"
	"github.com/palantir/pkg/matcher"
	"github.com/pkg/errors"
	yamlv3 "go.yaml.in/yaml/v3"
	"gopkg.in/yaml.v2"
)

// ReadGodelConfigFromProjectDir reads the gödel configuration from the "godel.yml" file in the configuration directory
// for the gödel project with the specified project directory and returns it. Returns an empty configuration if the
// configuration file does not exist.
func ReadGodelConfigFromProjectDir(projectDir string) (GodelConfig, error) {
	cfgDir, err := godellauncher.ConfigDirPath(projectDir)
	if err != nil {
		return GodelConfig{}, err
	}
	return ReadGodelConfigFromFile(filepath.Join(cfgDir, godellauncher.GodelConfigYML))
}

// ReadGodelConfigFromFile reads the gödel configuration from the provided file and returns the loaded configuration.
// Returns an empty configuration if the file does not exist. Returns an error if any mapping in the configuration
// specifies the same key more than once.
func ReadGodelConfigFromFile(cfgFile string) (GodelConfig, error) {
	if _, err := os.Stat(cfgFile); os.IsNotExist(err) {
		return GodelConfig{}, nil
	}

	var godelCfg GodelConfig
	bytes, err := os.ReadFile(cfgFile)
	if err != nil {
		return GodelConfig{}, errors.Wrapf(err, "failed to read file %s", cfgFile)
	}
	if err := verifyNoDuplicateKeys(bytes); err != nil {
		return GodelConfig{}, errors.Wrapf(err, "invalid gödel configuration in %s", cfgFile)
	}
	upgradedBytes, err := UpgradeConfig(bytes)
	if err != nil {
		return GodelConfig{}, errors.Wrapf(err, "failed to upgrade configuration")
	}
	if err := yaml.Unmarshal(upgradedBytes, &godelCfg); err != nil {
		return GodelConfig{}, errors.Wrapf(err, "failed to unmarshal gödel config YAML")
	}
	return godelCfg, nil
}

// ReadGodelConfigExcludesFromFile reads the excludes specified in the gödel godel configuration from the provided file
// and returns the loaded configuration. Returns an empty configuration if the file does not exist. Callers that only
// require the exclude configuration should prefer this function to using the ReadGodelConfigFrom* functions and
// accessing the exclude configuration there, as this function is more robust to configuration changes (for example, the
// ReadGodelConfigFrom* functions will return an error if a portion of the configuration that is unrelated to excludes
// cannot be unmarshaled into the expected type or contains duplicate keys, while this function only considers the
// "exclude" portion of configuration). Note that neither function returns an error for unrecognized keys.
func ReadGodelConfigExcludesFromFile(cfgFilePath string) (matcher.NamesPathsCfg, error) {
	if _, err := os.Stat(cfgFilePath); os.IsNotExist(err) {
		return matcher.NamesPathsCfg{}, nil
	}

	type excludeConfig struct {
		Exclude matcher.NamesPathsCfg `yaml:"exclude,omitempty"`
	}
	cfgBytes, err := os.ReadFile(cfgFilePath)
	if err != nil {
		return matcher.NamesPathsCfg{}, errors.WithStack(err)
	}
	var exclude excludeConfig
	if err := yaml.Unmarshal(cfgBytes, &exclude); err != nil {
		return matcher.NamesPathsCfg{}, errors.WithStack(err)
	}
	return exclude.Exclude, nil
}

// verifyNoDuplicateKeys returns an error if any mapping in the provided YAML specifies the same key more than once.
//
// This check is needed because the configuration is unmarshaled using the non-strict yaml.v2 decoder, which silently
// uses the last value for a duplicate key. This hides editing mistakes such as a locator with two "id" keys, which
// pairs the version from the last "id" with whatever checksums were not updated. yaml.v2 strict decoding is not used
// because it also rejects unrecognized keys, which would break existing configurations that are otherwise valid.
//
// Content that cannot be parsed by yaml.v3 is not reported as an error by this function: it only rejects duplicate
// keys, and malformed content is reported by the yaml.v2 unmarshal that follows.
func verifyNoDuplicateKeys(cfgBytes []byte) error {
	var root yamlv3.Node
	if err := yamlv3.Unmarshal(cfgBytes, &root); err != nil {
		return nil
	}
	var duplicates []duplicateKey
	collectDuplicateKeys(&root, &duplicates)
	if len(duplicates) == 0 {
		return nil
	}
	sort.Slice(duplicates, func(i, j int) bool {
		return duplicates[i].line < duplicates[j].line
	})
	var duplicateStrs []string
	for _, duplicate := range duplicates {
		duplicateStrs = append(duplicateStrs, fmt.Sprintf("line %d: key %q already defined at line %d", duplicate.line, duplicate.key, duplicate.firstLine))
	}
	return errors.Errorf("configuration contains duplicate keys:\n    %s", strings.Join(duplicateStrs, "\n    "))
}

type duplicateKey struct {
	key       string
	line      int
	firstLine int
}

func collectDuplicateKeys(node *yamlv3.Node, duplicates *[]duplicateKey) {
	if node.Kind == yamlv3.MappingNode {
		keyLines := make(map[string]int)
		for i := 0; i+1 < len(node.Content); i += 2 {
			key := node.Content[i]
			// "<<" is the YAML merge key, which yaml.v2 processes for every occurrence rather than overwriting
			if key.Kind != yamlv3.ScalarNode || key.Value == "<<" {
				continue
			}
			if line, ok := keyLines[key.Value]; ok {
				*duplicates = append(*duplicates, duplicateKey{key: key.Value, line: key.Line, firstLine: line})
				continue
			}
			keyLines[key.Value] = key.Line
		}
	}
	for _, child := range node.Content {
		collectDuplicateKeys(child, duplicates)
	}
}
