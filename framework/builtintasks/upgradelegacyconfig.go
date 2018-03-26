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
	"strings"

	"github.com/palantir/pkg/matcher"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"

	"github.com/palantir/godel/framework/godellauncher"
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

				// record all of the original YML files in the directory
				originalYMLFiles, err := dirYMLFiles(cfgDirPath)
				if err != nil {
					return err
				}
				// track all of the upgraded YML files
				upgradedYMLFiles := make(map[string]struct{})

				upgradeTasksMap := make(map[string]godellauncher.UpgradeConfigTask)
				for _, upgradeTask := range upgradeTasks {
					upgradeTasksMap[upgradeTask.ID] = upgradeTask
				}

				var legacyConfigUpgraderIDs []string
				for _, upgradeTask := range upgradeTasks {
					if upgradeTask.LegacyConfigFile == "" {
						continue
					}
					legacyConfigUpgraderIDs = append(legacyConfigUpgraderIDs, upgradeTask.ID)
				}
				sort.Strings(legacyConfigUpgraderIDs)

				var failedUpgrades []string
				// perform hard-coded one-time upgraders
				for _, currUpgrader := range hardCodedLegacyUpgraders {
					if err := currUpgrader.upgradeConfig(cfgDirPath, dryRunFlagVal, printContentFlagVal, stdout); err != nil {
						failedUpgrades = append(failedUpgrades, upgradeError(projectDir, path.Join(cfgDirPath, currUpgrader.configFileName()), err))
					}
					upgradedYMLFiles[currUpgrader.configFileName()] = struct{}{}
				}
				for _, k := range legacyConfigUpgraderIDs {
					upgradeTask, ok := upgradeTasksMap[k]
					if !ok {
						// legacy task does not have an upgrader: continue
						continue
					}
					upgradedYMLFiles[upgradeTask.ConfigFile] = struct{}{}
					if err := upgradeLegacyConfig(upgradeTask, cfgDirPath, global, dryRunFlagVal, printContentFlagVal, stdout); err != nil {
						failedUpgrades = append(failedUpgrades, upgradeError(projectDir, path.Join(cfgDirPath, upgradeTask.ConfigFile), err))
						continue
					}
				}

				var unknownYMLFiles []string
				for _, k := range originalYMLFiles {
					if _, ok := upgradedYMLFiles[k]; ok {
						continue
					}
					unknownYMLFiles = append(unknownYMLFiles, k)
				}
				if len(unknownYMLFiles) > 0 {
					// if unknown files were present, print warning
					dryRunPrintln(stdout, dryRunFlagVal, fmt.Sprintf(`WARNING: the following configuration files had no known upgraders for legacy configuration: %v
If these configuration files are plugins, add the plugin to the configuration and re-run the upgrade task.`, unknownYMLFiles))
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

var hardCodedLegacyUpgraders = []hardCodedLegacyUpgrader{
	&hardCodedLegacyUpgraderImpl{
		fileName: "exclude.yml",
		upgradeConfigFn: func(configDirPath string, dryRun, printContent bool, stdout io.Writer) error {
			// godel.yml itself is compatible. Only work to be performed is if "exclude.yml" exists and contains entries that
			// differ from godel.yml.
			legacyExcludeFilePath := path.Join(configDirPath, "exclude.yml")
			if _, err := os.Stat(legacyExcludeFilePath); os.IsNotExist(err) {
				// if legacy file does not exist, there is no upgrade to be performed
				return nil
			}
			legacyConfigBytes, err := ioutil.ReadFile(legacyExcludeFilePath)
			if err != nil {
				return errors.Wrapf(err, "failed to read legacy configuration file")
			}
			var excludeCfg matcher.NamesPathsCfg
			if err := yaml.UnmarshalStrict(legacyConfigBytes, &excludeCfg); err != nil {
				return errors.Wrapf(err, "failed to unmarshal legacy exclude configuration")
			}

			currentGodelConfig, err := readGodelConfig(path.Join(configDirPath, "godel.yml"))
			if err != nil {
				return errors.Wrapf(err, "failed to read godel configuration")
			}

			existingNames := make(map[string]struct{})
			for _, name := range currentGodelConfig.Exclude.Names {
				existingNames[name] = struct{}{}
			}
			existingPaths := make(map[string]struct{})
			for _, path := range currentGodelConfig.Exclude.Paths {
				existingPaths[path] = struct{}{}
			}

			modified := false
			for _, legacyName := range excludeCfg.Names {
				if _, ok := existingNames[legacyName]; ok {
					continue
				}
				currentGodelConfig.Exclude.Names = append(currentGodelConfig.Exclude.Names, legacyName)
				modified = true
			}
			for _, legacyPath := range excludeCfg.Paths {
				if _, ok := existingPaths[legacyPath]; ok {
					continue
				}
				currentGodelConfig.Exclude.Names = append(currentGodelConfig.Exclude.Paths, legacyPath)
				modified = true
			}

			// back up old configuration by moving it
			if err := backupConfigFile(legacyExcludeFilePath, dryRun, stdout); err != nil {
				return errors.Wrapf(err, "failed to back up legacy configuration file")
			}

			if !modified {
				// exclude.yml did not provide any new excludes: no need to write
				return nil
			}

			upgradedCfgBytes, err := yaml.Marshal(currentGodelConfig)
			if err != nil {
				return errors.Wrapf(err, "failed to marshal upgraded godel configuration")
			}

			if !dryRun {
				// write migrated configuration
				if err := ioutil.WriteFile(path.Join(configDirPath, "godel.yml"), upgradedCfgBytes, 0644); err != nil {
					return errors.Wrapf(err, "failed to write upgraded configuration")
				}
			}
			printUpgradedConfig("godel.yml", upgradedCfgBytes, dryRun, printContent, stdout)
			return nil
		},
	},
	&hardCodedLegacyUpgraderImpl{
		fileName: "imports.yml",
		upgradeConfigFn: func(configDirPath string, dryRun, printContent bool, stdout io.Writer) error {
			legacyConfigFilePath := path.Join(configDirPath, "imports.yml")
			if _, err := os.Stat(legacyConfigFilePath); os.IsNotExist(err) {
				// if legacy file does not exist, there is no upgrade to be performed
				return nil
			}
			legacyConfigBytes, err := ioutil.ReadFile(legacyConfigFilePath)
			if err != nil {
				return errors.Wrapf(err, "failed to read legacy configuration file")
			}

			// back up old configuration by moving it
			if err := backupConfigFile(legacyConfigFilePath, dryRun, stdout); err != nil {
				return errors.Wrapf(err, "failed to back up legacy configuration file")
			}

			if string(legacyConfigBytes) != "" {
				// if legacy file had content, warn that the functionality is no longer supported
				dryRunPrintln(stdout, dryRun, `WARNING: the functionality provided by "imports.yml" is no longer supported: the file has been backed up without being upgraded or migrated`)
			}
			return nil
		},
	},
}

func dirYMLFiles(inputDir string) ([]string, error) {
	fis, err := ioutil.ReadDir(inputDir)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to read input directory")
	}
	var ymlFiles []string
	for _, fi := range fis {
		if fi.IsDir() {
			continue
		}
		if strings.HasSuffix(fi.Name(), ".yml") {
			ymlFiles = append(ymlFiles, fi.Name())
		}
	}
	return ymlFiles, nil
}

type hardCodedLegacyUpgrader interface {
	configFileName() string
	upgradeConfig(configDirPath string, dryRun, printContent bool, stdout io.Writer) error
}

type hardCodedLegacyUpgraderImpl struct {
	fileName        string
	upgradeConfigFn func(configDirPath string, dryRun, printContent bool, stdout io.Writer) error
}

func (u *hardCodedLegacyUpgraderImpl) configFileName() string {
	return u.fileName
}

func (u *hardCodedLegacyUpgraderImpl) upgradeConfig(configDirPath string, dryRun, printContent bool, stdout io.Writer) error {
	return u.upgradeConfigFn(configDirPath, dryRun, printContent, stdout)
}

func upgradeLegacyConfig(upgradeTask godellauncher.UpgradeConfigTask, configDirPath string, global godellauncher.GlobalConfig, dryRun, printContent bool, stdout io.Writer) error {
	legacyConfigFilePath := path.Join(configDirPath, upgradeTask.LegacyConfigFile)
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
