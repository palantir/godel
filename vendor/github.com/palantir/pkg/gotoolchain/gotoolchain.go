// Copyright (c) 2026 Palantir Technologies. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package gotoolchain determines the Go toolchain that should be used for a module based on the content of its go.mod
// file. The package provides a Go implementation and an equivalent bash implementation (see BashScript) of the same
// selection rules so that tools written in Go and environments that can only run shell scripts (such as CI
// configuration) select the same toolchain.
//
// The toolchain for a go.mod file is selected as follows:
//
//  1. If the "toolchain" directive is present and its value is not "default", its value is used. The value must name a
//     standard Go release, beta or release candidate toolchain (for example, "go1.27.3", "go1.28rc1" or "go1.28beta1").
//     Release toolchains for Go 1.21 and later must include a patch version: "go1.27" is rejected because no such
//     toolchain exists. Custom toolchains with suffixes (for example, "go1.27.3-custom") are rejected.
//  2. Otherwise, the version in the "go" directive is used with a "go" prefix. If the "go" directive specifies only a
//     Major.Minor version, ".0" is appended (for example, "go 1.27" selects "go1.27.0").
//
// A valid "go" directive is required even if a "toolchain" directive determines the result. The selected value is
// returned as-is even if a "toolchain" directive specifies a version that is lower than the "go" directive (the go
// command rejects such a module): consumers that need different behavior should implement it on top of the values
// returned by this package.
//
// Only the provided go.mod file is read: go.work files are never consulted. Selection does not invoke the go command,
// does not perform any network operations and does not verify that the selected toolchain exists.
//
// The Go implementation parses go.mod using golang.org/x/mod/modfile and thus rejects go.mod files that are not
// syntactically valid. The bash implementation only parses the "go" and "toolchain" directives. Both implementations
// produce the same result for any syntactically valid go.mod file, and both reject missing, repeated or malformed "go"
// and "toolchain" directives.
package gotoolchain

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"

	"golang.org/x/mod/modfile"
)

// NormalizeMode determines how Normalize handles Major.Minor versions.
type NormalizeMode int

const (
	// NormalizeExact requires a patch version for releases of Go 1.21 and later: "1.27" is rejected. Major.Minor
	// versions before Go 1.21 are preserved ("1.20" becomes "go1.20") because those were the names of initial releases.
	// This is the mode used for "toolchain" directive values and explicitly specified toolchain versions.
	NormalizeExact NormalizeMode = iota
	// NormalizeLanguage appends ".0" to all Major.Minor versions, including versions before Go 1.21 ("1.27" becomes
	// "go1.27.0" and "1.20" becomes "go1.20.0"). This is the mode used for "go" directive values.
	NormalizeLanguage
)

// String returns the name of the mode as used by the bash implementation.
func (m NormalizeMode) String() string {
	switch m {
	case NormalizeExact:
		return "exact"
	case NormalizeLanguage:
		return "language"
	default:
		return "NormalizeMode(" + strconv.Itoa(int(m)) + ")"
	}
}

var (
	majorMinorRegexp      = regexp.MustCompile(`^1\.(0|[1-9][0-9]*)$`)
	standardVersionRegexp = regexp.MustCompile(`^1\.(0|[1-9][0-9]*)(\.(0|[1-9][0-9]*)|rc[1-9][0-9]*|beta[1-9][0-9]*)$`)
)

// Normalize returns the standard toolchain name for the provided version, which may optionally have a "go" prefix (for
// example, "1.27.3" and "go1.27.3" both become "go1.27.3"). Patch, beta and release candidate versions are preserved.
// The handling of Major.Minor versions is determined by mode. Returns an error if the version is not a standard Go
// release, beta or release candidate version: queries (such as "latest") and custom toolchain suffixes are rejected.
// Normalization does not verify that the returned toolchain exists.
func Normalize(version string, mode NormalizeMode) (string, error) {
	if mode != NormalizeExact && mode != NormalizeLanguage {
		return "", fmt.Errorf("unsupported normalize mode %v", mode)
	}
	v := strings.TrimPrefix(version, "go")
	if m := majorMinorRegexp.FindStringSubmatch(v); m != nil {
		if mode == NormalizeLanguage {
			return "go" + v + ".0", nil
		}
		// minor versions with more than 2 digits are compared by length so that the comparison cannot overflow
		if minor, _ := strconv.Atoi(m[1]); len(m[1]) > 2 || minor >= 21 {
			return "", fmt.Errorf("explicit version %q needs a patch number (for example go%s.0)", version, v)
		}
		return "go" + v, nil
	}
	if !standardVersionRegexp.MatchString(v) {
		return "", fmt.Errorf("invalid standard Go version %q: expected a release, beta, or release candidate", version)
	}
	return "go" + v, nil
}

// Directives are the values of the "go" and "toolchain" directives of a go.mod file. A field is empty if the
// corresponding directive is absent.
type Directives struct {
	Go        string
	Toolchain string
}

// ParseDirectives parses the provided go.mod content and returns its "go" and "toolchain" directives. The path is only
// used in error messages. Returns an error if the content is not a syntactically valid go.mod file (including if the
// "go" or "toolchain" directive is repeated or malformed). Module versions in other directives are not required to be
// canonical, since they do not affect toolchain selection. Does not require the "go" directive to be present: use
// Directives.SelectedToolchain to select a toolchain.
func ParseDirectives(path string, data []byte) (Directives, error) {
	// modfile.ParseLax is not used because it ignores the toolchain directive. The identity version fixer accepts
	// non-canonical versions in other directives, which modfile.Parse rejects if no fixer is provided.
	f, err := modfile.Parse(path, data, func(_, version string) (string, error) { return version, nil })
	if err != nil {
		return Directives{}, fmt.Errorf("failed to parse %s: %w", path, err)
	}
	var d Directives
	if f.Go != nil {
		d.Go = f.Go.Version
	}
	if f.Toolchain != nil {
		d.Toolchain = f.Toolchain.Name
	}
	return d, nil
}

// ReadDirectives reads the go.mod file at the provided path and returns its "go" and "toolchain" directives. See
// ParseDirectives for details.
func ReadDirectives(goModPath string) (Directives, error) {
	data, err := os.ReadFile(goModPath)
	if err != nil {
		return Directives{}, fmt.Errorf("failed to read go.mod file: %w", err)
	}
	return ParseDirectives(goModPath, data)
}

// SelectedToolchain returns the toolchain (for example, "go1.27.3") selected by the directives using the rules
// described in the package documentation. Returns an error if the "go" directive is absent or invalid or if the
// "toolchain" directive does not name a standard toolchain.
func (d Directives) SelectedToolchain() (string, error) {
	if d.Go == "" {
		return "", fmt.Errorf("missing go directive")
	}
	if strings.HasPrefix(d.Go, "go") {
		return "", fmt.Errorf("go directive %q must contain a numeric version", d.Go)
	}
	goToolchain, err := Normalize(d.Go, NormalizeLanguage)
	if err != nil {
		return "", fmt.Errorf("invalid go directive: %w", err)
	}
	if d.Toolchain == "" || d.Toolchain == "default" {
		return goToolchain, nil
	}
	if !strings.HasPrefix(d.Toolchain, "go") {
		return "", fmt.Errorf("toolchain directive %q must name a Go toolchain", d.Toolchain)
	}
	toolchain, err := Normalize(d.Toolchain, NormalizeExact)
	if err != nil {
		return "", fmt.Errorf("invalid toolchain directive: %w", err)
	}
	return toolchain, nil
}

// FromGoModContent returns the toolchain selected by the provided go.mod content. The path is only used in error
// messages.
func FromGoModContent(path string, data []byte) (string, error) {
	d, err := ParseDirectives(path, data)
	if err != nil {
		return "", err
	}
	toolchain, err := d.SelectedToolchain()
	if err != nil {
		return "", fmt.Errorf("failed to select toolchain for %s: %w", path, err)
	}
	return toolchain, nil
}

// FromGoMod returns the toolchain selected by the go.mod file at the provided path.
func FromGoMod(goModPath string) (string, error) {
	data, err := os.ReadFile(goModPath)
	if err != nil {
		return "", fmt.Errorf("failed to read go.mod file: %w", err)
	}
	return FromGoModContent(goModPath, data)
}
