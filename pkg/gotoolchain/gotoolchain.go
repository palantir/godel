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

// FromGoMod returns the toolchain name for the specified go.mod file. An explicit
// toolchain other than "default" takes precedence over the required go directive.
// A go directive containing only major.minor is normalized to major.minor.0.
// Custom toolchains are unsupported. This function does not inspect go.work,
// invoke Go, or compare the selected toolchain with the go directive.
func FromGoMod(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("read go.mod %s: %w", path, err)
	}
	// ParseLax ignores toolchain directives because they apply only to main modules.
	mod, err := modfile.Parse(path, data, nil)
	if err != nil {
		return "", fmt.Errorf("parse go.mod %s: %w", path, err)
	}
	if mod.Go == nil {
		return "", fmt.Errorf("go.mod %s must contain a go directive", path)
	}

	name := ""
	if mod.Toolchain != nil && mod.Toolchain.Name != "default" {
		name = mod.Toolchain.Name
	} else {
		v := mod.Go.Version
		if languageVersion.MatchString(v) {
			v += ".0"
		}
		name = "go" + v
	}
	if strings.Contains(name, "-") {
		return "", fmt.Errorf("go.mod %s specifies unsupported non-standard toolchain %q", path, name)
	}
	if !version.IsValid(name) {
		return "", fmt.Errorf("go.mod %s specifies invalid toolchain %q", path, name)
	}
	return name, nil
}
