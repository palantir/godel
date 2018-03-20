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
	"strings"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/palantir/godel/framework/godellauncher"
)

func UpgradeConfigTask(upgradeTasks []godellauncher.UpgradeConfigTask) godellauncher.Task {
	const (
		upgradeConfigCmdName = "upgrade-config"
		dryRunFlagName       = "dry-run"
		printContentFlagName = "print-content"
	)

	var (
		dryRunFlagVal       bool
		printContentFlagVal bool
	)

	cmd := &cobra.Command{
		Use:   upgradeConfigCmdName,
		Short: "Upgrade configuration",
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
				configDirPath, err := godellauncher.ConfigDirPath(projectDir)
				if err != nil {
					return err
				}

				var failedUpgrades []string
				for _, upgradeTask := range upgradeTasks {
					changed, upgradedCfgBytes, err := upgradeConfigFile(upgradeTask, global, configDirPath, dryRunFlagVal, stdout)
					if err != nil {
						failedUpgrades = append(failedUpgrades, fmt.Sprintf("%s: %v", path.Join(configDirPath, upgradeTask.ConfigFile), err))
						continue
					}
					if !changed {
						continue
					}
					printUpgradedConfig(upgradeTask.ConfigFile, upgradedCfgBytes, dryRunFlagVal, printContentFlagVal, stdout)
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

func dryRunPrintln(w io.Writer, dryRun bool, content string) {
	if !dryRun {
		fmt.Fprintln(w, content)
		return
	}
	const dryRunPrefix = "[DRY RUN] "
	fmt.Fprintln(w, dryRunPrefix+content)
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

func upgradeConfigFile(task godellauncher.UpgradeConfigTask, global godellauncher.GlobalConfig, configDir string, dryRun bool, stdout io.Writer) (bool, []byte, error) {
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

	if err := backupConfigFile(configFile, dryRun, stdout); err != nil {
		return false, nil, err
	}
	if !dryRun {
		if err := ioutil.WriteFile(configFile, upgradedConfigBytes, 0644); err != nil {
			return false, nil, errors.Wrapf(err, "failed to write upgraded configuration")
		}
	}
	return true, upgradedConfigBytes, nil
}

func backupConfigFile(cfgFilePath string, dryRun bool, stdout io.Writer) error {
	// if file does not exist, no need to back up
	if _, err := os.Stat(cfgFilePath); os.IsNotExist(err) {
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
