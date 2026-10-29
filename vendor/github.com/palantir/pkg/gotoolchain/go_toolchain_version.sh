#!/usr/bin/env bash

# Copyright (c) 2026 Palantir Technologies. All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

# Canonical bash implementation of the Go toolchain selection rules of the github.com/palantir/pkg/gotoolchain Go
# package. Do not modify copies of this file: they are distributed from the gotoolchain module (gotoolchain.BashScript)
# and are overwritten when updated. Changes must be made in the gotoolchain module, where tests verify that this
# implementation and the Go implementation behave identically.
#
# Source this file to define the functions below. Sourcing has no other effect: the functions do not invoke Go,
# download toolchains, modify files, or change the caller's environment or shell options. The functions are written to
# be compatible with bash 3.2 and later and can be used with "set -euo pipefail".
#
# Every function prints its result to stdout on success. On failure, it prints a diagnostic to stderr, prints nothing
# to stdout, and returns a nonzero status.
#
# The toolchain for a go.mod file is selected as follows:
#   1. If the "toolchain" directive is present and its value is not "default", its value is used. The value must name a
#      standard Go release, beta or release candidate toolchain. Release toolchains for Go 1.21 and later must include a
#      patch version ("go1.27" is rejected). Custom toolchains with suffixes ("go1.27.3-custom") are rejected.
#   2. Otherwise, the version in the "go" directive is used with a "go" prefix, appending ".0" if the "go" directive
#      specifies only a Major.Minor version ("go 1.27" selects "go1.27.0").
# A valid "go" directive is required even if the "toolchain" directive determines the result. A "toolchain" directive
# that is lower than the "go" directive is still selected. go.work files are never consulted.
#
# Only the "go" and "toolchain" directives of a go.mod file are parsed: the rest of the file is assumed to be valid.

# Prints the toolchain selected by the go.mod file at $1.
#
# Equivalent to gotoolchain.FromGoMod.
go_toolchain_version_from_go_mod() {
  local directives go_directive toolchain_directive
  directives=$(go_toolchain_version_directives "${1:-}") || return
  go_directive=${directives%%$'\n'*}
  go_directive=${go_directive#go=}
  toolchain_directive=${directives#*$'\n'toolchain=}
  go_toolchain_version_select "$go_directive" "$toolchain_directive" || {
    go_toolchain_version_error "failed to select toolchain for '${1:-}'"
    return 1
  }
}

# Prints the values of the "go" and "toolchain" directives of the go.mod file at $1 on 2 lines of the form
# "go=<value>" and "toolchain=<value>" (in that order). A value is empty if the corresponding directive is absent. Fails
# if either directive is repeated or malformed.
#
# The values can be extracted while preserving the exit status as follows:
#
#   directives=$(go_toolchain_version_directives go.mod) || exit
#   go_directive=${directives%%$'\n'*}
#   go_directive=${go_directive#go=}
#   toolchain_directive=${directives#*$'\n'toolchain=}
#
# Equivalent to gotoolchain.ReadDirectives, except that only these 2 directives are parsed and validated.
go_toolchain_version_directives() {
  local go_directive='' toolchain_directive='' found_go=false found_toolchain=false line directive value words
  [[ -f ${1:-} ]] || { go_toolchain_version_error "go.mod file '${1:-}' does not exist"; return 1; }

  while IFS= read -r line || [[ -n $line ]]; do
    # as in golang.org/x/mod/modfile, a comment starts at the first "//" (quoted strings are not supported)
    line=${line%%//*}
    words=()
    IFS=$' \t\r' read -r -a words <<< "$line"
    directive=${words[0]:-}
    [[ $directive == go || $directive == toolchain ]] || continue
    if (( ${#words[@]} != 2 )); then
      go_toolchain_version_error "$directive directive in '$1' expects exactly one argument"
      return 1
    fi
    value=${words[1]}
    case $directive in
      go)
        if [[ $found_go == true ]]; then
          go_toolchain_version_error "repeated go directive in '$1'"
          return 1
        fi
        if [[ ! $value =~ ^[1-9][0-9]*\.(0|[1-9][0-9]*)(\.(0|[1-9][0-9]*))?([a-z]+[0-9]+)?$ ]]; then
          go_toolchain_version_error "invalid go version '$value' in '$1': must match format 1.23.0"
          return 1
        fi
        found_go=true
        go_directive=$value
        ;;
      toolchain)
        if [[ $found_toolchain == true ]]; then
          go_toolchain_version_error "repeated toolchain directive in '$1'"
          return 1
        fi
        if [[ $value != default && $value != go1 && $value != go1.* ]]; then
          go_toolchain_version_error "invalid toolchain version '$value' in '$1': must match format go1.23.0 or default"
          return 1
        fi
        found_toolchain=true
        toolchain_directive=$value
        ;;
    esac
  done < "$1" || return

  printf 'go=%s\ntoolchain=%s\n' "$go_directive" "$toolchain_directive"
}

# Prints the toolchain selected by the "go" directive value $1 and the optional "toolchain" directive value $2 (an
# empty or omitted $2 indicates that the "toolchain" directive is absent).
#
# Equivalent to gotoolchain.Directives.SelectedToolchain.
go_toolchain_version_select() {
  local go_directive=${1:-} toolchain_directive=${2:-} go_toolchain
  if [[ -z $go_directive ]]; then
    go_toolchain_version_error "missing go directive"
    return 1
  fi
  if [[ $go_directive == go* ]]; then
    go_toolchain_version_error "go directive '$go_directive' must contain a numeric version"
    return 1
  fi
  go_toolchain=$(go_toolchain_version_normalize "$go_directive" language) || {
    go_toolchain_version_error "invalid go directive"
    return 1
  }
  if [[ -z $toolchain_directive || $toolchain_directive == default ]]; then
    printf '%s\n' "$go_toolchain"
    return
  fi
  if [[ $toolchain_directive != go* ]]; then
    go_toolchain_version_error "toolchain directive '$toolchain_directive' must name a Go toolchain"
    return 1
  fi
  go_toolchain_version_normalize "$toolchain_directive" exact || {
    go_toolchain_version_error "invalid toolchain directive"
    return 1
  }
}

# Prints the standard toolchain name for the version $1, which may optionally have a "go" prefix ("1.27.3" and
# "go1.27.3" both print "go1.27.3"). Patch, beta and release candidate versions are preserved. The optional mode $2
# determines the handling of Major.Minor versions:
#   * "exact" (the default if omitted or empty): requires a patch version for Go 1.21 and later ("1.27" is rejected),
#     but preserves the names of initial releases before Go 1.21 ("1.20" prints "go1.20"). Used for "toolchain"
#     directive values and explicitly specified toolchain versions.
#   * "language": appends ".0" to all Major.Minor versions ("1.27" prints "go1.27.0" and "1.20" prints "go1.20.0").
#     Used for "go" directive values.
# Fails if the version is not a standard Go release, beta or release candidate version: queries (such as "latest") and
# custom toolchain suffixes are rejected. Does not verify that the printed toolchain exists.
#
# Equivalent to gotoolchain.Normalize.
go_toolchain_version_normalize() {
  local version=${1:-} mode=${2:-exact} minor
  if [[ $mode != exact && $mode != language ]]; then
    go_toolchain_version_error "unsupported normalize mode '$mode'"
    return 1
  fi
  version=${version#go}
  if [[ $version =~ ^1\.(0|[1-9][0-9]*)$ ]]; then
    minor=${BASH_REMATCH[1]}
    if [[ $mode == language ]]; then
      version+=.0
    elif (( ${#minor} > 2 || 10#$minor >= 21 )); then
      go_toolchain_version_error "explicit version '${1:-}' needs a patch number (for example go$version.0)"
      return 1
    fi
  elif [[ ! $version =~ ^1\.(0|[1-9][0-9]*)(\.(0|[1-9][0-9]*)|rc[1-9][0-9]*|beta[1-9][0-9]*)$ ]]; then
    go_toolchain_version_error "invalid standard Go version '${1:-}': expected a release, beta, or release candidate"
    return 1
  fi
  printf 'go%s\n' "$version"
}

# Prints the error message $* to stderr so that callers can capture results from stdout. The caller determines the
# failure status.
go_toolchain_version_error() {
  printf 'gotoolchain: %s\n' "$*" >&2
}
