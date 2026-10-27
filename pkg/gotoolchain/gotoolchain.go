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

// Package gotoolchain resolves the exact Go toolchain specified by a module.
package gotoolchain

import (
	"fmt"
	"go/version"
	"os"
	"regexp"
	"strings"

	"golang.org/x/mod/modfile"
)

var languageVersion = regexp.MustCompile(`^[0-9]+\.[0-9]+$`)

// FromGoMod returns the toolchain name for the specified go.mod file. The toolchain that should be used for a go.mod
// file is determined as follows:
//
//  1. If the toolchain directive has any value other than "default", that value is used.
//  2. If the toolchain directive is absent or is "default", the go directive's version is used. If that version consists
//     of only major.minor (for example, 1.27), it is normalized to "major.minor.0".
//
// A go directive is required even when an explicit toolchain directive exists. Custom toolchains are unsupported. This
// function does not inspect go.work, invoke Go, or compare the selected toolchain with the go directive.
//
// Returns an error if the provided go.mod file is invalid or does not specify a valid toolchain based on the outlined
// rules.
func FromGoMod(goModFilePath string) (string, error) {
	data, err := os.ReadFile(goModFilePath)
	if err != nil {
		return "", fmt.Errorf("read go.mod %s: %w", goModFilePath, err)
	}
	// ParseLax ignores toolchain directives because they apply only to main modules.
	mod, err := modfile.Parse(goModFilePath, data, nil)
	if err != nil {
		return "", fmt.Errorf("parse go.mod %s: %w", goModFilePath, err)
	}
	if mod.Go == nil {
		return "", fmt.Errorf("go.mod %s must contain a go directive", goModFilePath)
	}

	var toolchainValue string
	if mod.Toolchain != nil && mod.Toolchain.Name != "default" {
		toolchainValue = mod.Toolchain.Name
	} else {
		v := mod.Go.Version
		if languageVersion.MatchString(v) {
			v += ".0"
		}
		toolchainValue = "go" + v
	}
	if strings.Contains(toolchainValue, "-") {
		return "", fmt.Errorf("go.mod %s specifies unsupported non-standard toolchain %q", goModFilePath, toolchainValue)
	}
	if !version.IsValid(toolchainValue) {
		return "", fmt.Errorf("go.mod %s specifies invalid toolchain %q", goModFilePath, toolchainValue)
	}
	return toolchainValue, nil
}
