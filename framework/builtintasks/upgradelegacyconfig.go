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
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"sort"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"

	"github.com/palantir/godel/framework/godellauncher"
	"github.com/palantir/godel/framework/internal/legacyplugins"
)

func UpgradeLegacyConfigTask(upgradeTasks []godellauncher.UpgradeConfigTask) godellauncher.Task {
	const (
		upgradeLegacyConfigCmdName = "upgrade-legacy-config"
		dryRunFlagName             = "dry-run"
		printContentFlagName       = "print-content"
	)

	var (
		dryRunFlagVal       bool
		printContentFlagVal bool
	)

	cmd := &cobra.Command{
		Use:   upgradeLegacyConfigCmdName,
		Short: "Upgrade the legacy configuration",
	}

	cmd.Flags().BoolVar(&dryRunFlagVal, dryRunFlagName, false, "print what the upgrade operation would do without writing changes")
	cmd.Flags().BoolVar(&printContentFlagVal, printContentFlagName, false, "print the content of the changes to stdout in addition to writing them")

	cmd.SilenceErrors = true
	cmd.SilenceUsage = true
	return godellauncher.Task{
		Name:        cmd.Use,
		Description: cmd.Short,
		RunImpl: func(t *godellauncher.Task, global godellauncher.GlobalConfig, stdout io.Writer) error {
			args := []string{global.Executable}
			args = append(args, global.Task)
			args = append(args, global.TaskArgs...)
			os.Args = args

			cmd.RunE = func(cmd *cobra.Command, args []string) error {
				projectDir, err := global.ProjectDir()
				if err != nil {
					return err
				}
				cfgDirPath, err := godellauncher.ConfigDirPath(projectDir)
				if err != nil {
					return err
				}

				upgradeTasksMap := make(map[string]godellauncher.UpgradeConfigTask)
				for _, upgradeTask := range upgradeTasks {
					upgradeTasksMap[upgradeTask.ID] = upgradeTask
				}

				var legacyConfigUpgraderKeys []string
				for k := range legacyplugins.LegacyConfigUpgraders {
					legacyConfigUpgraderKeys = append(legacyConfigUpgraderKeys, k)
				}
				sort.Strings(legacyConfigUpgraderKeys)

				var failedUpgrades []string
				for _, k := range legacyConfigUpgraderKeys {
					upgradeTask, ok := upgradeTasksMap[k]
					if !ok {
						// legacy task does not have an upgrader: continue
						continue
					}
					if err := upgradeLegacyConfig(upgradeTask, cfgDirPath, global, dryRunFlagVal, printContentFlagVal, stdout); err != nil {
						failedUpgrades = append(failedUpgrades, fmt.Sprintf("%s: %v", path.Join(cfgDirPath, upgradeTask.ConfigFile), err))
						continue
					}
				}

				if len(failedUpgrades) == 0 {
					return nil
				}
				dryRunPrintln(stdout, dryRunFlagVal, "Failed to upgrade configuration:")
				for _, upgrade := range failedUpgrades {
					dryRunPrintln(stdout, dryRunFlagVal, "\t"+upgrade)
				}
				return fmt.Errorf("")
			}

			rootCmd := godellauncher.CobraCmdToRootCmd(cmd)
			rootCmd.SetOutput(stdout)
			return rootCmd.Execute()
		},
	}
}

func upgradeLegacyConfig(upgradeTask godellauncher.UpgradeConfigTask, configDirPath string, global godellauncher.GlobalConfig, dryRun, printContent bool, stdout io.Writer) error {
	legacyConfigFilePath := path.Join(configDirPath, legacyplugins.LegacyConfigUpgraders[upgradeTask.ID].LegacyConfigFileName)
	if _, err := os.Stat(legacyConfigFilePath); os.IsNotExist(err) {
		// if legacy file does not exist, there is no upgrade to be performed
		return nil
	}

	legacyConfigBytes, err := ioutil.ReadFile(legacyConfigFilePath)
	if err != nil {
		return errors.Wrapf(err, "failed to read legacy configuration file")
	}

	var ymlConfig yaml.MapSlice
	if err := yaml.Unmarshal(legacyConfigBytes, &ymlConfig); err != nil {
		return errors.Wrapf(err, "failed to unmarshal YAML configuration")
	}
	// add "legacy-config: true" as a key to indicate that this is a legacy configuration
	ymlConfig = append([]yaml.MapItem{{Key: "legacy-config", Value: true}}, ymlConfig...)

	ymlCfgBytes, err := yaml.Marshal(ymlConfig)
	if err != nil {
		return errors.Wrapf(err, "failed to marshal YAML")
	}
	upgradedCfgBytes, err := upgradeTask.Run(ymlCfgBytes, global, stdout)
	if err != nil {
		return errors.Wrapf(err, "failed to upgrade configuration")
	}

	// back up old configuration by moving it
	if err := backupConfigFile(legacyConfigFilePath, dryRun, stdout); err != nil {
		return errors.Wrapf(err, "failed to back up legacy configuration file")
	}

	// upgraded configuration is empty: no need to write
	if string(upgradedCfgBytes) == "" {
		return nil
	}

	if !dryRun {
		// write migrated configuration
		if err := ioutil.WriteFile(path.Join(configDirPath, upgradeTask.ConfigFile), upgradedCfgBytes, 0644); err != nil {
			return errors.Wrapf(err, "failed to write upgraded configuration")
		}
	}
	printUpgradedConfig(upgradeTask.ConfigFile, upgradedCfgBytes, dryRun, printContent, stdout)

	return nil
}
