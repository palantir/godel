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
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"

	"github.com/palantir/pkg/matcher"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"

	"github.com/palantir/godel/v2/framework/godel/config"
	"github.com/palantir/godel/v2/framework/godellauncher"
)

func UpgradeConfigTask(upgradeTasks []godellauncher.UpgradeConfigTask) godellauncher.Task {
	const (
		upgradeConfigCmdName = "upgrade-config"
		dryRunFlagName       = "dry-run"
		printContentFlagName = "print-content"
		legacyFlagName       = "legacy"
		backupFlagName       = "backup"
	)

	var (
		dryRunFlagVal       bool
		printContentFlagVal bool
		legacyFlagVal       bool
		backupFlagVal       bool
	)

	cmd := &cobra.Command{
		Use:   upgradeConfigCmdName,
		Short: "Upgrade configuration",
	}

	cmd.Flags().BoolVar(&dryRunFlagVal, dryRunFlagName, false, "print what the upgrade operation would do without writing changes")
	cmd.Flags().BoolVar(&printContentFlagVal, printContentFlagName, false, "print the content of the changes to stdout in addition to writing them")
	cmd.Flags().BoolVar(&legacyFlagVal, legacyFlagName, false, "upgrade pre-2.0 legacy configuration")
	cmd.Flags().BoolVar(&backupFlagVal, backupFlagName, false, "back up files before overwriting or removing them")

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
				configDirPath, err := godellauncher.ConfigDirPath(projectDir)
				if err != nil {
					return err
				}
				if legacyFlagVal {
					return runUpgradeLegacyConfig(upgradeTasks, global, projectDir, configDirPath, backupFlagVal, dryRunFlagVal, printContentFlagVal, cmd.OutOrStdout())
				}
				return runUpgradeConfig(upgradeTasks, global, projectDir, configDirPath, backupFlagVal, dryRunFlagVal, printContentFlagVal, cmd.OutOrStdout())
			}

			rootCmd := godellauncher.CobraCmdToRootCmd(cmd)
			rootCmd.SetOutput(stdout)
			return rootCmd.Execute()
		},
	}
}

func runUpgradeConfig(
	upgradeTasks []godellauncher.UpgradeConfigTask,
	global godellauncher.GlobalConfig,
	projectDir, configDirPath string,
	backup, dryRun, printContent bool,
	stdout io.Writer,
) error {

	var failedUpgrades []string
	for _, upgradeTask := range upgradeTasks {
		changed, upgradedCfgBytes, err := upgradeConfigFile(upgradeTask, global, configDirPath, backup, dryRun, stdout)
		if err != nil {
			failedUpgrades = append(failedUpgrades, upgradeError(projectDir, path.Join(configDirPath, upgradeTask.ConfigFile), err))
			continue
		}
		if !changed {
			continue
		}
		printUpgradedConfig(upgradeTask.ConfigFile, upgradedCfgBytes, dryRun, printContent, stdout)
	}

	if len(failedUpgrades) == 0 {
		return nil
	}
	dryRunPrintln(stdout, dryRun, "Failed to upgrade configuration:")
	for _, upgrade := range failedUpgrades {
		dryRunPrintln(stdout, dryRun, "\t"+upgrade)
	}
	return fmt.Errorf("")
}

func dryRunPrintln(w io.Writer, dryRun bool, content string) {
	if !dryRun {
		_, _ = fmt.Fprintln(w, content)
		return
	}
	const dryRunPrefix = "[DRY RUN] "
	_, _ = fmt.Fprintln(w, dryRunPrefix+content)
}

func printUpgradedConfig(cfgFile string, upgradedCfgBytes []byte, dryRun, printContent bool, stdout io.Writer) {
	dryRunPrintln(stdout, dryRun, fmt.Sprintf("Upgraded configuration for %s", cfgFile))
	if !printContent {
		return
	}
	dryRunPrintln(stdout, dryRun, "---")
	cfgStr := strings.TrimSuffix(string(upgradedCfgBytes), "\n")
	for _, line := range strings.Split(cfgStr, "\n") {
		dryRunPrintln(stdout, dryRun, line)
	}
	dryRunPrintln(stdout, dryRun, "---")
}

func upgradeConfigFile(task godellauncher.UpgradeConfigTask, global godellauncher.GlobalConfig, configDir string, backup, dryRun bool, stdout io.Writer) (bool, []byte, error) {
	configFile := path.Join(configDir, task.ConfigFile)
	origConfigBytes, err := ioutil.ReadFile(configFile)
	if err != nil {
		return false, nil, errors.Wrapf(err, "failed to read config file")
	}
	upgradedConfigBytes, err := task.Run(origConfigBytes, global, stdout)
	if err != nil {
		return false, nil, err
	}
	if changed := !bytes.Equal(origConfigBytes, upgradedConfigBytes); !changed {
		return false, nil, nil
	}

	if backup {
		if err := backupConfigFile(configFile, dryRun, stdout); err != nil {
			return false, nil, err
		}
	}
	if !dryRun {
		if err := ioutil.WriteFile(configFile, upgradedConfigBytes, 0644); err != nil {
			return false, nil, errors.Wrapf(err, "failed to write upgraded configuration")
		}
	}
	return true, upgradedConfigBytes, nil
}

func backupConfigFile(cfgFilePath string, dryRun bool, stdout io.Writer) error {
	// if file does not exist or it exists but is empty, no need to back up
	if fi, err := os.Stat(cfgFilePath); os.IsNotExist(err) || (err == nil && fi.Size() == 0) {
		return nil
	}

	dstPath := path.Join(path.Dir(cfgFilePath), path.Base(cfgFilePath)+".bak")
	if dryRun {
		if wd, err := os.Getwd(); err == nil {
			if filepath.IsAbs(cfgFilePath) {
				if relPath, err := filepath.Rel(wd, cfgFilePath); err == nil {
					cfgFilePath = relPath
				}
			}
			if filepath.IsAbs(dstPath) {
				if relPath, err := filepath.Rel(wd, dstPath); err == nil {
					dstPath = relPath
				}
			}
		}
		dryRunPrintln(stdout, dryRun, fmt.Sprintf("Run: mv %s %s", cfgFilePath, dstPath))
	} else {
		// file exists: move to backup location
		if err := os.Rename(cfgFilePath, dstPath); err != nil {
			return errors.Wrapf(err, "failed to rename configuration file")
		}
	}
	return nil
}

func removeConfigFile(cfgFilePath string, dryRun bool, stdout io.Writer) error {
	// if file does not exist, nothing to do
	if _, err := os.Stat(cfgFilePath); os.IsNotExist(err) {
		return nil
	}
	// remove file or print operation
	if dryRun {
		dryRunPrintln(stdout, dryRun, fmt.Sprintf("Run: rm %s", cfgFilePath))
		return nil
	}
	if err := os.Remove(cfgFilePath); err != nil {
		return errors.Wrapf(err, "failed to remove configuration file %s", cfgFilePath)
	}
	return nil
}

func upgradeError(projectDir, configFilePath string, upgradeErr error) string {
	// convert path to relative path if it is absolute. No-op if conversion to absolute path fails.
	if filepath.IsAbs(configFilePath) {
		if relPath, err := filepath.Rel(projectDir, configFilePath); err == nil {
			configFilePath = relPath
		}
	}
	return fmt.Sprintf("%s: %v", configFilePath, upgradeErr)
}

func runUpgradeLegacyConfig(
	upgradeTasks []godellauncher.UpgradeConfigTask,
	global godellauncher.GlobalConfig,
	projectDir, configDirPath string,
	backup, dryRun, printContent bool,
	stdout io.Writer,
) error {

	// record all of the original YML files in the directory
	originalYMLFiles, err := dirYMLFiles(configDirPath)
	if err != nil {
		return err
	}
	// track all of the upgraded YML files
	knownConfigFiles := make(map[string]struct{})

	upgradeTasksMap := make(map[string]godellauncher.UpgradeConfigTask)
	for _, upgradeTask := range upgradeTasks {
		upgradeTasksMap[upgradeTask.ID] = upgradeTask
	}

	var failedUpgrades []string
	// perform hard-coded one-time upgrades
	for _, currUpgrader := range hardCodedLegacyUpgraders {
		if err := currUpgrader.upgradeConfig(configDirPath, backup, dryRun, printContent, stdout); err != nil {
			failedUpgrades = append(failedUpgrades, upgradeError(projectDir, path.Join(configDirPath, currUpgrader.configFileName()), err))
		}
		knownConfigFiles[currUpgrader.configFileName()] = struct{}{}
	}

	var legacyConfigUpgraderIDs []string
	for _, upgradeTask := range upgradeTasks {
		// consider current configuration file for the plugin as known (don't warn if these files already
		// existed in config directory but were not processed by a legacy config upgrader).
		knownConfigFiles[upgradeTask.ConfigFile] = struct{}{}
		if upgradeTask.LegacyConfigFile == "" {
			continue
		}
		legacyConfigUpgraderIDs = append(legacyConfigUpgraderIDs, upgradeTask.ID)
	}
	sort.Strings(legacyConfigUpgraderIDs)
	for _, k := range legacyConfigUpgraderIDs {
		upgradeTask, ok := upgradeTasksMap[k]
		if !ok {
			// legacy task does not have an upgrader: continue
			continue
		}
		knownConfigFiles[upgradeTask.LegacyConfigFile] = struct{}{}
		if err := upgradeLegacyConfig(upgradeTask, configDirPath, global, backup, dryRun, printContent, stdout); err != nil {
			failedUpgrades = append(failedUpgrades, upgradeError(projectDir, path.Join(configDirPath, upgradeTask.ConfigFile), err))
			continue
		}
	}

	var unhandledYMLFiles []string
	for _, k := range originalYMLFiles {
		if _, ok := knownConfigFiles[k]; ok {
			continue
		}
		unhandledYMLFiles = append(unhandledYMLFiles, k)
	}
	if err := processUnhandledYMLFiles(configDirPath, unhandledYMLFiles, backup, dryRun, stdout); err != nil {
		return err
	}

	if len(failedUpgrades) == 0 {
		return nil
	}
	dryRunPrintln(stdout, dryRun, "Failed to upgrade configuration:")
	for _, upgrade := range failedUpgrades {
		dryRunPrintln(stdout, dryRun, "\t"+upgrade)
	}
	return fmt.Errorf("")
}

var hardCodedLegacyUpgraders = []hardCodedLegacyUpgrader{
	&hardCodedLegacyUpgraderImpl{
		fileName: "exclude.yml",
		upgradeConfigFn: func(configDirPath string, backup, dryRun, printContent bool, stdout io.Writer) error {
			// godel.yml itself is compatible. Only work to be performed is if "exclude.yml" exists and contains entries
			// that differ from godel.yml.
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

			currentGodelConfig, err := config.ReadGodelConfigFromFile(path.Join(configDirPath, "godel.yml"))
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
				currentGodelConfig.Exclude.Paths = append(currentGodelConfig.Exclude.Paths, legacyPath)
				modified = true
			}

			if backup {
				// back up old configuration by moving it
				if err := backupConfigFile(legacyExcludeFilePath, dryRun, stdout); err != nil {
					return errors.Wrapf(err, "failed to back up legacy configuration file")
				}
			} else {
				// remove old configuration file
				if err := removeConfigFile(legacyExcludeFilePath, dryRun, stdout); err != nil {
					return errors.Wrapf(err, "failed to remove legacy configuration file")
				}
			}

			if !modified {
				// exclude.yml did not provide any new excludes: no need to write
				return nil
			}

			upgradedCfgBytes, err := yaml.Marshal(currentGodelConfig)
			if err != nil {
				return errors.Wrapf(err, "failed to marshal upgraded godel configuration")
			}

			godelYMLPath := path.Join(configDirPath, "godel.yml")
			if backup {
				// back up godel.yml because it is about to be overwritten
				if err := backupConfigFile(godelYMLPath, dryRun, stdout); err != nil {
					return errors.Wrapf(err, "failed to back up godel.yml")
				}
			}
			if !dryRun {
				// write migrated configuration
				if err := ioutil.WriteFile(godelYMLPath, upgradedCfgBytes, 0644); err != nil {
					return errors.Wrapf(err, "failed to write upgraded configuration")
				}
			}
			printUpgradedConfig("godel.yml", upgradedCfgBytes, dryRun, printContent, stdout)
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
	upgradeConfig(configDirPath string, backup, dryRun, printContent bool, stdout io.Writer) error
}

type hardCodedLegacyUpgraderImpl struct {
	fileName        string
	upgradeConfigFn func(configDirPath string, backup, dryRun, printContent bool, stdout io.Writer) error
}

func (u *hardCodedLegacyUpgraderImpl) configFileName() string {
	return u.fileName
}

func (u *hardCodedLegacyUpgraderImpl) upgradeConfig(configDirPath string, backup, dryRun, printContent bool, stdout io.Writer) error {
	return u.upgradeConfigFn(configDirPath, backup, dryRun, printContent, stdout)
}

func upgradeLegacyConfig(upgradeTask godellauncher.UpgradeConfigTask, configDirPath string, global godellauncher.GlobalConfig, backup, dryRun, printContent bool, stdout io.Writer) error {
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

	if backup {
		// back up old configuration
		if err := backupConfigFile(legacyConfigFilePath, dryRun, stdout); err != nil {
			return errors.Wrapf(err, "failed to back up legacy configuration file")
		}
	} else {
		// remove old configuration
		if err := removeConfigFile(legacyConfigFilePath, dryRun, stdout); err != nil {
			return errors.Wrapf(err, "failed to remove legacy configuration file")
		}
	}

	dstFilePath := path.Join(configDirPath, upgradeTask.ConfigFile)
	if backup {
		// back up destination file if it already exists
		if err := backupConfigFile(dstFilePath, dryRun, stdout); err != nil {
			return errors.Wrapf(err, "failed to back up existing configuration file")
		}
	}

	// upgraded configuration is empty: no need to write
	if string(upgradedCfgBytes) == "" {
		return nil
	}

	if !dryRun {
		// write migrated configuration
		if err := ioutil.WriteFile(dstFilePath, upgradedCfgBytes, 0644); err != nil {
			return errors.Wrapf(err, "failed to write upgraded configuration")
		}
	}
	printUpgradedConfig(upgradeTask.ConfigFile, upgradedCfgBytes, dryRun, printContent, stdout)
	return nil
}

func processUnhandledYMLFiles(configDir string, unknownYMLFiles []string, backup, dryRun bool, stdout io.Writer) error {
	if len(unknownYMLFiles) == 0 {
		return nil
	}

	var unknownNonEmptyFiles []string
	for _, currUnknownFile := range unknownYMLFiles {
		currPath := path.Join(configDir, currUnknownFile)
		bytes, err := ioutil.ReadFile(currPath)
		if err != nil {
			return errors.Wrapf(err, "failed to read configuration file")
		}
		// if unknown file is empty, just remove it or back it up
		if string(bytes) == "" {
			if backup {
				if err := backupConfigFile(currPath, dryRun, stdout); err != nil {
					return err
				}
			} else {
				if err := removeConfigFile(currPath, dryRun, stdout); err != nil {
					return err
				}
			}
			continue
		}
		unknownNonEmptyFiles = append(unknownNonEmptyFiles, currUnknownFile)
	}

	if len(unknownNonEmptyFiles) == 0 {
		return nil
	}

	// if non-empty unknown files were present, print warning
	dryRunPrintln(stdout, dryRun, fmt.Sprintf(`WARNING: The following configuration file(s) were non-empty and had no known upgraders for legacy configuration: %v`, unknownNonEmptyFiles))
	dryRunPrintln(stdout, dryRun, fmt.Sprintf(`         If these configuration file(s) are for plugins, add the plugins to the configuration and rerun the upgrade task.`))
	return nil
}
