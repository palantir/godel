// Copyright (c) 2026 Palantir Technologies. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gotoolchain

import (
	"bytes"
	_ "embed"
)

// BashScriptFileName is the name of the file that contains the bash implementation in this module.
const BashScriptFileName = "go_toolchain_version.sh"

//go:embed go_toolchain_version.sh
var bashScript []byte

// BashScript returns the content of a bash script that implements the same toolchain selection rules as this package.
// Sourcing the script defines the following functions, which print their result to stdout on success and print a
// diagnostic to stderr and return a nonzero status on failure:
//
//   - go_toolchain_version_from_go_mod <go.mod path>: equivalent to FromGoMod
//   - go_toolchain_version_directives <go.mod path>: equivalent to ReadDirectives; prints "go=<value>" and
//     "toolchain=<value>" lines
//   - go_toolchain_version_select <go directive> [<toolchain directive>]: equivalent to Directives.SelectedToolchain
//   - go_toolchain_version_normalize <version> [exact|language]: equivalent to Normalize
//
// The script is intended to be written to repositories and environments that cannot use the Go implementation so that
// they use the same selection logic. A new slice is returned on every call, so callers may modify it.
func BashScript() []byte {
	return bytes.Clone(bashScript)
}
