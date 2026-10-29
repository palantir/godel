// Copyright 2026 Palantir Technologies, Inc.
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
	"fmt"
	"maps"
	"os"
	"path/filepath"

	"github.com/palantir/pkg/gotoolchain"
)

// SetEnvVars sets the environment variables based on the GodelConfig receiver c.
//
// If c.ManageGoToolchain is true, this function reads the go.mod file in the project directory and sets the value of
// the GOTOOLCHAIN environment variable based on that file. If c.ManageGoToolchain is true, this function returns an
// error if c.Environment has an entry for the key GOTOOLCHAIN.
//
// If c.ManageGoToolchain is false, the go.mod file is not read.
//
// After considering the GOTOOLCHAIN logic described above, sets all the environment variables specified in
// c.Environment.
func (c GodelConfig) SetEnvVars(projectDir string) error {
	envVarsMap := maps.Clone(c.Environment)
	if envVarsMap == nil {
		envVarsMap = make(map[string]string)
	}

	// if Go toolchain should be managed, verify that input map does not specify a GOTOOLCHAIN value,
	// determine the toolchain value based on go.mod, and include the computed value as a value to be set.
	if c.ManageGoToolchain != nil && *c.ManageGoToolchain {
		if _, exists := envVarsMap["GOTOOLCHAIN"]; exists {
			return fmt.Errorf("ManageGoToolchain is true: environment must not specify GOTOOLCHAIN")
		}
		toolchain, err := gotoolchain.FromGoMod(filepath.Join(projectDir, "go.mod"))
		if err != nil {
			return err
		}
		envVarsMap["GOTOOLCHAIN"] = toolchain
	}

	for k, v := range envVarsMap {
		if err := os.Setenv(k, v); err != nil {
			return fmt.Errorf("set environment variable %s: %w", k, err)
		}
	}

	return nil
}
