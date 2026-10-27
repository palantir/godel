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
	// set environment variables
	for k, v := range c.Environment {
		if err := os.Setenv(k, v); err != nil {
			return fmt.Errorf("set environment variable %s: %w", k, err)
		}
	}

	// if GOTOOLCHAIN is not managed, done
	manageGoToolchain := c.ManageGoToolchain != nil && *c.ManageGoToolchain
	if !manageGoToolchain {
		return nil
	}

	if _, exists := c.Environment["GOTOOLCHAIN"]; exists {
		return fmt.Errorf("the map Environment cannot specify GOTOOLCHAIN when ManageGoToolchain is true")
	}

	toolchain, err := gotoolchain.FromGoMod(filepath.Join(projectDir, "go.mod"))
	if err != nil {
		return err
	}

	if err := os.Setenv("GOTOOLCHAIN", toolchain); err != nil {
		return fmt.Errorf("set GOTOOLCHAIN: %w", err)
	}

	return nil
}
