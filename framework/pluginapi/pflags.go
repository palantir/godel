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

	"github.com/spf13/pflag"
)

func AddAllPFlags(fset *pflag.FlagSet) (debug *bool, projectDir *string, godelConfig *string, config *string, assets *[]string) {
	goFlagSet := &flag.FlagSet{}
	debug, projectDir, godelConfig, config = AddAllFlags(goFlagSet)
	fset.AddGoFlagSet(goFlagSet)
	assets = AddAssetsPFlag(fset)
	return
}

func AddAllPFlagsPtrs(fset *pflag.FlagSet, debug *bool, projectDir *string, godelConfig *string, config *string, assets *[]string) {
	goFlagSet := &flag.FlagSet{}
	AddAllFlagsPtrs(goFlagSet, debug, projectDir, godelConfig, config)
	fset.AddGoFlagSet(goFlagSet)
	AddAssetsPFlagPtr(fset, assets)
}

func AddDebugPFlag(fset *pflag.FlagSet) *bool {
	gofset := &flag.FlagSet{}
	debug := AddDebugFlag(gofset)
	fset.AddGoFlagSet(gofset)
	return debug
}

func AddDebugPFlagPtr(fset *pflag.FlagSet, debug *bool) {
	if debug == nil {
		return
	}
	fset.BoolVar(debug, DebugFlagName, false, "run in debug mode")
}

func AddProjectDirPFlag(fset *pflag.FlagSet) *string {
	gofset := &flag.FlagSet{}
	projectDir := AddProjectDirFlag(gofset)
	fset.AddGoFlagSet(gofset)
	return projectDir
}

func AddProjectDirPFlagPtr(fset *pflag.FlagSet, projectDir *string) {
	if projectDir == nil {
		return
	}
	fset.StringVar(projectDir, ProjectDirFlagName, "", "path to project directory")
}

func AddGodelConfigPFlag(fset *pflag.FlagSet) *string {
	gofset := &flag.FlagSet{}
	godelConfig := AddGodelConfigFlag(gofset)
	fset.AddGoFlagSet(gofset)
	return godelConfig
}

func AddGodelConfigPFlagPtr(fset *pflag.FlagSet, godelConfig *string) {
	if godelConfig == nil {
		return
	}
	fset.StringVar(godelConfig, GodelConfigFlagName, "", "path to the godel.yml configuration file")
}

func AddConfigPFlag(fset *pflag.FlagSet) *string {
	gofset := &flag.FlagSet{}
	config := AddConfigFlag(gofset)
	fset.AddGoFlagSet(gofset)
	return config
}

func AddConfigPFlagPtr(fset *pflag.FlagSet, config *string) {
	if config == nil {
		return
	}
	fset.StringVar(config, ConfigFlagName, "", "path to the plugin configuration file")
}

func AddAssetsPFlag(fset *pflag.FlagSet) *[]string {
	var assets []string
	AddAssetsPFlagPtr(fset, &assets)
	return &assets
}

func AddAssetsPFlagPtr(fset *pflag.FlagSet, assets *[]string) {
	if assets == nil {
		return
	}
	fset.StringSliceVar(assets, AssetsFlagName, nil, "path(s) to the plugin asset(s)")
}
