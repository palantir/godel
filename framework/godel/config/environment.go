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
	"os"
	"path/filepath"

	"github.com/palantir/godel/v2/pkg/gotoolchain"
)

// ApplyEnvironment sets configured environment variables in the current process.
// When manage-gotoolchain is true, it also selects GOTOOLCHAIN from go.mod in
// projectDir, overriding any inherited value. Call before loading or running
// plugins. When management is disabled, go.mod is not read.
func (c GodelConfig) ApplyEnvironment(projectDir string) error {
	manage := c.ManageGoToolchain != nil && *c.ManageGoToolchain
	var toolchain string
	if manage {
		if _, exists := c.Environment["GOTOOLCHAIN"]; exists {
			return fmt.Errorf("manage-gotoolchain is true: environment must not specify GOTOOLCHAIN")
		}
		var err error
		toolchain, err = gotoolchain.FromGoMod(filepath.Join(projectDir, "go.mod"))
		if err != nil {
			return err
		}
	}

	for k, v := range c.Environment {
		if err := os.Setenv(k, v); err != nil {
			return fmt.Errorf("set environment variable %s: %w", k, err)
		}
	}
	if manage {
		if err := os.Setenv("GOTOOLCHAIN", toolchain); err != nil {
			return fmt.Errorf("set GOTOOLCHAIN: %w", err)
		}
	}
	return nil
}
