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

package builtintasks

import (
	"io/ioutil"
	"os"
	"path"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"

	"github.com/palantir/godel/framework/godel/config"
	"github.com/palantir/godel/framework/godellauncher"
)

func InfoTask() godellauncher.Task {
	var globalCfg godellauncher.GlobalConfig
	cmd := &cobra.Command{
		Use:   "info",
		Short: "Print information regarding gödel",
	}
	cmd.AddCommand(&cobra.Command{
		Use:   "default-tasks",
		Short: "Print configuration for default tasks",
		RunE: func(cmd *cobra.Command, args []string) error {
			projectDir, err := globalCfg.ProjectDir()
			if err != nil {
				return err
			}
			godelCfg, err := ReadGodelConfigFromProjectDir(projectDir)
			if err != nil {
				return err
			}
			bytes, err := yaml.Marshal(godelCfg.DefaultTasks)
			if err != nil {
				return errors.Wrapf(err, "failed to marshal default task configuration")
			}
			cmd.Print(string(bytes))
			return nil
		},
	})
	return godellauncher.CobraCLITask(cmd, &globalCfg)
}

// ReadGodelConfigFromProjectDir reads the gödel configuration from the "godel.yml" file in the configuration file for
// the gödel project with the specified project directory and returns it.
func ReadGodelConfigFromProjectDir(projectDir string) (config.GodelConfig, error) {
	cfgDir, err := godellauncher.ConfigDirPath(projectDir)
	if err != nil {
		return config.GodelConfig{}, err
	}
	return readGodelConfig(path.Join(cfgDir, godellauncher.GodelConfigYML))
}

func readGodelConfig(cfgFile string) (config.GodelConfig, error) {
	var godelCfg config.GodelConfig
	if _, err := os.Stat(cfgFile); err == nil {
		bytes, err := ioutil.ReadFile(cfgFile)
		if err != nil {
			return config.GodelConfig{}, errors.Wrapf(err, "failed to read file %s", cfgFile)
		}
		upgradedBytes, err := config.UpgradeConfig(bytes)
		if err != nil {
			return config.GodelConfig{}, errors.Wrapf(err, "failed to upgrade configuration")
		}
		if err := yaml.Unmarshal(upgradedBytes, &godelCfg); err != nil {
			return config.GodelConfig{}, errors.Wrapf(err, "failed to unmarshal gödel config YAML")
		}
	}
	return godelCfg, nil
}
