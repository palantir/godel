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
	"path"
	"strings"

	"github.com/palantir/godel/pkg/osarch"
	"github.com/pkg/errors"

	"github.com/palantir/distgo/distgo"
	"github.com/palantir/distgo/distgo/config/internal/v0"
)

type BuildConfig v0.BuildConfig

func ToBuildConfig(in *BuildConfig) *v0.BuildConfig {
	return (*v0.BuildConfig)(in)
}

// ToParam returns the BuildParam represented by the receiver *BuildConfig and the provided default BuildConfig. If a
// config value is specified (non-nil) in the receiver config, it is used. If a config value is not specified in the
// receiver config but is specified in the default config, the default config value is used. If a value is not specified
// in either configuration, the program-specified default value (if any) is used.
func (cfg *BuildConfig) ToParam(scriptIncludes string, defaultCfg BuildConfig) (distgo.BuildParam, error) {
	outputDir := getConfigStringValue(cfg.OutputDir, defaultCfg.OutputDir, "out/build")
	if path.IsAbs(outputDir) {
		return distgo.BuildParam{}, errors.Errorf("output-dir cannot be specified as an absolute path")
	}
	mainPkg := getConfigStringValue(cfg.MainPkg, defaultCfg.MainPkg, "")
	if mainPkg != "" && !strings.HasPrefix(mainPkg, "./") {
		mainPkg = "./" + mainPkg
	}

	return distgo.BuildParam{
		NameTemplate:    getConfigStringValue(cfg.NameTemplate, defaultCfg.NameTemplate, "{{Product}}"),
		OutputDir:       outputDir,
		MainPkg:         mainPkg,
		BuildArgsScript: distgo.CreateScriptContent(getConfigStringValue(cfg.BuildArgsScript, defaultCfg.BuildArgsScript, ""), scriptIncludes),
		VersionVar:      getConfigStringValue(cfg.VersionVar, defaultCfg.VersionVar, ""),
		Environment:     getConfigValue(cfg.Environment, defaultCfg.Environment, nil).(map[string]string),
		OSArchs:         getConfigValue(cfg.OSArchs, defaultCfg.OSArchs, []osarch.OSArch{osarch.Current()}).([]osarch.OSArch),
	}, nil
}
