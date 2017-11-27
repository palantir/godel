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
	"path"

	"github.com/palantir/godel/framework/godellauncher"
)

// cfgCLIArgs takes the provided inputs and returns the arguments that are suitable for using as the "os.Args" for a
// CfgCLI. The returned arguments are of the form:
//
// global.Executable [--debug] [--config <cfgDirPath>/cfgFileName] [--json <gödelCfgJSON>] [cmdPath...] [global.TaskArgs...]
//
// * The "--debug" flag will be present if global.Debug is true.
// * The "--config" flag will be present if "cfgFileName" is non-empty. The "cfgDirPath" is the parent directory of the
//   global.Wrapper path.
func cfgCLIArgs(global godellauncher.GlobalConfig, cmdPath []string, cfgFileName string) ([]string, error) {
	args := []string{global.Executable}
	if global.Debug {
		args = append(args, "--debug")
	}
	if global.Wrapper != "" {
		gödelConfigDir, err := godellauncher.ConfigDirPath(path.Dir(global.Wrapper))
		if err != nil {
			return nil, err
		}
		if cfgFileName != "" {
			args = append(args, "--config", path.Join(gödelConfigDir, cfgFileName))
		}
		cfgJSON, err := godellauncher.GodelConfigJSON(gödelConfigDir)
		if err != nil {
			return nil, err
		}
		args = append(args, "--json", string(cfgJSON))
	}
	args = append(args, cmdPath...)
	args = append(args, global.TaskArgs...)
	return args, nil
}
