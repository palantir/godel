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

package apptasks

import (
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"

	"github.com/palantir/checks/gocd/cmd/gocd"
	"github.com/palantir/checks/gogenerate/cmd/gogenerate"
	"github.com/palantir/checks/golicense/cmd/golicense"
	"github.com/palantir/pkg/cli"
	"github.com/pkg/errors"

	"github.com/palantir/godel/apps/distgo"
	"github.com/palantir/godel/framework/godel"
	"github.com/palantir/godel/framework/godellauncher"
	"github.com/palantir/godel/framework/verifyorder"
)

func AppTasks() []godellauncher.Task {
	return []godellauncher.Task{
		createBuiltinVerifyTask("generate", gogenerate.App(), "generate.yml", verifyorder.Generate),
		createBuiltinVerifyTask("imports", gocd.App(), "imports.yml", verifyorder.Imports),
		createBuiltinVerifyTask("license", golicense.App(), "license.yml", verifyorder.License),
		createDisgtoTask("run"),
		createDisgtoTask("project-version"),
		createDisgtoTask("build"),
		createDisgtoTask("dist"),
		createDisgtoTask("clean"),
		createDisgtoTask("artifacts"),
		createDisgtoTask("products"),
		createDisgtoTask("docker"),
		createDisgtoTask("publish"),
	}
}

func createDisgtoTask(name string) godellauncher.Task {
	task := createBuiltinTaskHelper(name, distgo.App(), []string{name}, "dist.yml", nil)

	baseRunImpl := task.RunImpl
	// distgo requires working directory to be the base (project) directory, so decorate action to set the working
	// directory before invocation.
	task.RunImpl = func(t *godellauncher.Task, global godellauncher.GlobalConfig, stdout io.Writer) error {
		if global.Wrapper != "" {
			if !filepath.IsAbs(global.Wrapper) {
				absWrapperPath, err := filepath.Abs(global.Wrapper)
				if err != nil {
					return errors.Wrapf(err, "failed to convert wrapper path to absolute path")
				}
				global.Wrapper = absWrapperPath
			}
			if err := os.Chdir(path.Dir(global.Wrapper)); err != nil {
				return errors.Wrapf(err, "failed to change working directory")
			}
		}
		return baseRunImpl(t, global, stdout)
	}
	return task
}

func createBuiltinVerifyTask(name string, app *cli.App, cfgFileName string, verifyOrder int) godellauncher.Task {
	return createBuiltinTaskHelper(name, app, nil, cfgFileName, &godellauncher.VerifyOptions{
		Ordering:       verifyOrder,
		ApplyFalseArgs: []string{"--verify"},
	})
}

func createBuiltinTaskHelper(name string, app *cli.App, cmdPath []string, cfgFileName string, verify *godellauncher.VerifyOptions) godellauncher.Task {
	app.Name = godel.AppName
	currCmd := app.Command
	for _, wantSubCmdName := range cmdPath {
		for _, currSubCmd := range currCmd.Subcommands {
			if currSubCmd.Name == wantSubCmdName {
				currCmd = currSubCmd
				break
			}
		}
	}

	return godellauncher.Task{
		Name:        name,
		Description: currCmd.Usage,
		ConfigFile:  cfgFileName,
		Verify:      verify,
		RunImpl: func(t *godellauncher.Task, global godellauncher.GlobalConfig, stdout io.Writer) error {
			args, err := cfgCLIArgs(global, cmdPath, t.ConfigFile)
			if err != nil {
				return err
			}
			app.Stdout = stdout
			os.Args = args
			if exitCode := app.Run(args); exitCode != 0 {
				return fmt.Errorf("")
			}
			return nil
		},
	}
}
