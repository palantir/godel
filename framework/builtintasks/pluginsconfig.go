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
	"strings"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"

	"github.com/palantir/godel/framework/godel/config"
	"github.com/palantir/godel/framework/godellauncher"
)

func TasksConfigTask(tasksCfgInfo config.TasksConfigInfo) godellauncher.Task {
	return godellauncher.CobraCLITask(&cobra.Command{
		Use:   "tasks-config",
		Short: "Prints the full YAML configuration used to load tasks and assets",
		RunE: func(cmd *cobra.Command, args []string) error {
			return printTasksCfgInfo(tasksCfgInfo, cmd.OutOrStdout())
		},
	}, nil)
}

func printTasksCfgInfo(tasksCfgInfo config.TasksConfigInfo, stdout io.Writer) error {
	if err := printWithHeader("Built-in plugin configuration", tasksCfgInfo.BuiltinPluginsConfig, stdout); err != nil {
		return err
	}
	fmt.Fprintln(stdout)

	if err := printWithHeader("Fully resolved godel tasks configuration", tasksCfgInfo.TasksConfig, stdout); err != nil {
		return err
	}
	fmt.Fprintln(stdout)

	return printWithHeader("Plugin configuration for default tasks", tasksCfgInfo.DefaultTasksPluginsConfig, stdout)
}

func printWithHeader(header string, in interface{}, stdout io.Writer) error {
	printHeader(header, stdout)
	ymlString, err := toYAMLString(in)
	if err != nil {
		return err
	}
	fmt.Fprint(stdout, ymlString)
	return nil
}

func printHeader(header string, stdout io.Writer) {
	fmt.Fprintln(stdout, header+":")
	fmt.Fprintln(stdout, strings.Repeat("-", len(header)+1))
}

func toYAMLString(in interface{}) (string, error) {
	ymlBytes, err := yaml.Marshal(in)
	if err != nil {
		return "", errors.Wrapf(err, "failed to marshal YAML")
	}
	return string(ymlBytes), nil
}
