// Copyright 2016 Palantir Technologies, Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bufio"
	"bytes"
	"fmt"
	"go/ast"
	"go/token"
	"io"
	"regexp"
	"strconv"
	"strings"
)

// importBreakPaths returns a slice that contains the import paths before which a line breaks should be inserted.
func importBreakPaths(file *ast.File) []string {
	var output []string

	for _, decl := range file.Decls {
		if gen, ok := decl.(*ast.GenDecl); ok && gen.Tok == token.IMPORT {
			for srcIndex, currSpec := range gen.Specs {
				if srcIndex > 0 {
					// if there was a previous element, check if the group has changed
					currGroup := importGroup(importPath(currSpec))
					prevGroup := importGroup(importPath(gen.Specs[srcIndex-1]))
					if currGroup != prevGroup {
						// if group has changed, add path to the output
						output = append(output, importPath(gen.Specs[srcIndex]))
					}
				}
			}

			// assume that only one import token block exists
			break
		}
	}

	return output
}

// from github.com/palantir/go-palantir/ptimports
func importGroup(importPath string) int {
	switch {
	case inStandardLibrary(importPath):
		return 0
	default:
		return 1
	}
}

func inStandardLibrary(importPath string) bool {
	return !strings.Contains(importPath, ".")
}

type importSlice []ast.Spec

func (x importSlice) Len() int      { return len(x) }
func (x importSlice) Swap(i, j int) { x[i], x[j] = x[j], x[i] }
func (x importSlice) Less(i, j int) bool {
	ipath := importPath(x[i])
	jpath := importPath(x[j])

	igroup := importGroup(ipath)
	jgroup := importGroup(jpath)
	if igroup != jgroup {
		return igroup < jgroup
	}

	if ipath != jpath {
		return ipath < jpath
	}
	iname := importName(x[i])
	jname := importName(x[j])

	if iname != jname {
		return iname < jname
	}
	return importComment(x[i]) < importComment(x[j])
}

func importPath(s ast.Spec) string {
	t, err := strconv.Unquote(s.(*ast.ImportSpec).Path.Value)
	if err == nil {
		return t
	}
	return ""
}

func importName(s ast.Spec) string {
	n := s.(*ast.ImportSpec).Name
	if n == nil {
		return ""
	}
	return n.Name
}

func importComment(s ast.Spec) string {
	c := s.(*ast.ImportSpec).Comment
	if c == nil {
		return ""
	}
	return c.Text()
}

// from golang.org/x/tools/cmd/goimports/imports.go
var impLine = regexp.MustCompile(`^\s+(?:[\w\.]+\s+)?"(.+)"`)

func addImportSpaces(r io.Reader, breaks []string) []byte {
	var out bytes.Buffer
	sc := bufio.NewScanner(r)
	inImports := false
	done := false
	for sc.Scan() {
		s := sc.Text()

		if !inImports && !done && strings.HasPrefix(s, "import") {
			inImports = true
		}
		if inImports && (s == ")" ||
			strings.HasPrefix(s, "var") ||
			strings.HasPrefix(s, "func") ||
			strings.HasPrefix(s, "const") ||
			strings.HasPrefix(s, "type")) {
			done = true
			inImports = false
		}
		if inImports && len(breaks) > 0 {
			if m := impLine.FindStringSubmatch(s); m != nil {
				if m[1] == breaks[0] {
					if err := out.WriteByte('\n'); err != nil {
						panic(fmt.Errorf("Failed to write newline in addImportSpaces"))
					}
					breaks = breaks[1:]
				}
			}
		}
		if !inImports || s != "" {
			fmt.Fprintln(&out, s)
		}
	}
	return out.Bytes()
}
