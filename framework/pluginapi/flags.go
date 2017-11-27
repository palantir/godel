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

package pluginapi

import (
	"flag"
)

const (
	// DebugFlagName is the name of the boolean flag that will be provided as "--<DebugFlagName>" if "debug" is true.
	DebugFlagName = "debug"
	// ProjectDirFlagName is the name of the string flag that is provided as "--<ProjectDirFlagName> <dir>", where
	// "<dir>" is the
	// path to the project directory.
	ProjectDirFlagName = "project-dir"
	// GodelConfigFlagName is the name of the string flag that is provided as "--<GodelConfigFlagName> <config>", where
	// "<config>" is the path to the configuration file for gödel.
	GodelConfigFlagName = "godel-config"
	// ConfigFlagName is the name of the string flag that is provided as "--<ConfigFlagName> <config>", where "<config>"
	// is the path to the configuration file for the plugin.
	ConfigFlagName = "config"
	// AssetsFlagName is the name of the assets flag that is provided as "--<AssetsFlagName> <assets>", where "<assets>"
	// is a comma-delimited list of the paths to the assets for the plugin.
	AssetsFlagName = "assets"
)

func AddDebugFlag(fset *flag.FlagSet) *bool {
	var debug bool
	AddDebugFlagPtr(fset, &debug)
	return &debug
}

func AddDebugFlagPtr(fset *flag.FlagSet, debug *bool) {
	if debug == nil {
		return
	}
	fset.BoolVar(debug, DebugFlagName, false, "run in debug mode")
}

func AddProjectDirFlag(fset *flag.FlagSet) *string {
	return addStringFlag(fset, AddProjectDirFlagPtr)
}

func AddProjectDirFlagPtr(fset *flag.FlagSet, projectDir *string) {
	if projectDir == nil {
		return
	}
	fset.StringVar(projectDir, ProjectDirFlagName, "", "path to project directory")
}

func AddGodelConfigFlag(fset *flag.FlagSet) *string {
	return addStringFlag(fset, AddGodelConfigFlagPtr)
}

func AddGodelConfigFlagPtr(fset *flag.FlagSet, godelConfig *string) {
	if godelConfig == nil {
		return
	}
	fset.StringVar(godelConfig, GodelConfigFlagName, "", "path to the godel.yml configuration file")
}

func AddConfigFlag(fset *flag.FlagSet) *string {
	return addStringFlag(fset, AddConfigFlagPtr)
}

func AddConfigFlagPtr(fset *flag.FlagSet, config *string) {
	if config == nil {
		return
	}
	fset.String(ConfigFlagName, "", "path to the plugin configuration file")
}

func addStringFlag(fset *flag.FlagSet, fn func(*flag.FlagSet, *string)) *string {
	var str string
	fn(fset, &str)
	return &str
}

func AddAllFlags(fset *flag.FlagSet) (debug *bool, projectDir *string, godelConfig *string, config *string) {
	debug = AddDebugFlag(fset)
	projectDir = AddProjectDirFlag(fset)
	godelConfig = AddGodelConfigFlag(fset)
	config = AddConfigFlag(fset)
	return
}

func AddAllFlagsPtrs(fset *flag.FlagSet, debug *bool, projectDir *string, godelConfig *string, config *string) {
	AddDebugFlagPtr(fset, debug)
	AddProjectDirFlagPtr(fset, projectDir)
	AddGodelConfigFlagPtr(fset, godelConfig)
	AddConfigFlagPtr(fset, config)
}
