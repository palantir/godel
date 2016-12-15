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

// Based on golang.org/x/tools/imports which bears the following license:
//
// Copyright 2013 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ptimports

import (
	"bufio"
	"bytes"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/printer"
	"go/token"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/palantir/pkg/pkgpath"
	"golang.org/x/tools/go/ast/astutil"
)

// Process formats and adjusts imports for the provided file.
func Process(filename string, src []byte) ([]byte, error) {
	fileSet := token.NewFileSet()
	file, adjust, err := parse(fileSet, filename, src)
	if err != nil {
		return nil, err
	}

	if err := removeUnusedImports(fileSet, file, filename); err != nil {
		return nil, err
	}

	abs, err := filepath.Abs(filename)
	if err != nil {
		return nil, err
	}
	relative := abs
	if goPathSrcRel, err := pkgpath.NewAbsPkgPath(abs).GoPathSrcRel(); err == nil {
		relative = goPathSrcRel
	}
	segments := strings.Split(relative, "/")
	if len(segments) < 3 {
		return nil, fmt.Errorf("expected repo to be located under at least 3 subdirectories but received relative filepath: %v", relative)
	}
	// append trailing / to prevent matches on repos with superstring names
	repoPath := filepath.Join(segments[:3]...) + "/"
	grp := newVendoredGrouper(repoPath)

	fixImports(fileSet, file, grp, godepsPath(filename))
	imps := astutil.Imports(fileSet, file)

	var spacesBefore []string
	lastGroup := -1
	for _, impSection := range imps {
		// Within each block of contiguous imports, see if any

		// we'll need to put a space between them so it's
		// compatible with gofmt.
		for _, importSpec := range impSection {
			importPath, _ := strconv.Unquote(importSpec.Path.Value)
			groupNum := grp.importGroup(importPath)
			if groupNum != lastGroup && lastGroup != -1 {
				spacesBefore = append(spacesBefore, importPath)
			}
			lastGroup = groupNum
		}
	}

	printerMode := printer.UseSpaces | printer.TabIndent
	printConfig := &printer.Config{Mode: printerMode, Tabwidth: 8}

	var buf bytes.Buffer
	err = printConfig.Fprint(&buf, fileSet, file)
	if err != nil {
		return nil, err
	}
	out := buf.Bytes()
	if adjust != nil {
		out = adjust(src, out)
	}
	out = addImportSpaces(bytes.NewReader(out), spacesBefore)

	out, err = format.Source(out)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func godepsPath(filename string) string {
	// Start with absolute path of file being formatted
	abs, err := filepath.Abs(filename)
	if err != nil {
		return ""
	}
	dir := filepath.Dir(abs)

	// Go up directories looking for Godeps
	for {
		fi, err := os.Stat(filepath.Join(dir, "Godeps"))
		if err == nil && fi.IsDir() {
			break	// found
		}
		if err != nil && !os.IsNotExist(err) {
			return ""
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return ""
		}
		dir = parent
	}

	// Have absolute path, just need import path
	// Assume last path element containing dot is the first part of import path
	godeps := filepath.Join(dir, "Godeps", "_workspace", "src")
	lastDot := strings.LastIndex(godeps, ".")
	if lastDot == -1 {
		return ""
	}
	lastSlashBeforeDot := strings.LastIndex(godeps[:lastDot], "/")
	if lastSlashBeforeDot == -1 {
		return ""
	}
	return godeps[lastSlashBeforeDot+1:]
}

// parse parses src, which was read from filename,
// as a Go source file or statement list.
func parse(fset *token.FileSet, filename string, src []byte) (*ast.File, func(orig, src []byte) []byte, error) {
	parserMode := parser.ParseComments

	// Try as whole source file.
	file, err := parser.ParseFile(fset, filename, src, parserMode)
	if err == nil {
		return file, nil, nil
	}
	// If the error is that the source file didn't begin with a
	// package line, fall through to try as a source fragment.
	// Stop and return on any other error.
	if !strings.Contains(err.Error(), "expected 'package'") {
		return nil, nil, err
	}

	// If this is a declaration list, make it a source file
	// by inserting a package clause.
	// Insert using a ;, not a newline, so that the line numbers
	// in psrc match the ones in src.
	psrc := append([]byte("package main;"), src...)
	file, err = parser.ParseFile(fset, filename, psrc, parserMode)
	if err == nil {
		// If a main function exists, we will assume this is a main
		// package and leave the file.
		if containsMainFunc(file) {
			return file, nil, nil
		}

		adjust := func(orig, src []byte) []byte {
			// Remove the package clause.
			// Gofmt has turned the ; into a \n.
			src = src[len("package main\n"):]
			return matchSpace(orig, src)
		}
		return file, adjust, nil
	}
	// If the error is that the source file didn't begin with a
	// declaration, fall through to try as a statement list.
	// Stop and return on any other error.
	if !strings.Contains(err.Error(), "expected declaration") {
		return nil, nil, err
	}

	// If this is a statement list, make it a source file
	// by inserting a package clause and turning the list
	// into a function body.  This handles expressions too.
	// Insert using a ;, not a newline, so that the line numbers
	// in fsrc match the ones in src.
	fsrc := append(append([]byte("package p; func _() {"), src...), '}')
	file, err = parser.ParseFile(fset, filename, fsrc, parserMode)
	if err == nil {
		adjust := func(orig, src []byte) []byte {
			// Remove the wrapping.
			// Gofmt has turned the ; into a \n\n.
			src = src[len("package p\n\nfunc _() {"):]
			src = src[:len(src)-len("}\n")]
			// Gofmt has also indented the function body one level.
			// Remove that indent.
			src = bytes.Replace(src, []byte("\n\t"), []byte("\n"), -1)
			return matchSpace(orig, src)
		}
		return file, adjust, nil
	}

	// Failed, and out of options.
	return nil, nil, err
}

// containsMainFunc checks if a file contains a function declaration with the
// function signature 'func main()'
func containsMainFunc(file *ast.File) bool {
	for _, decl := range file.Decls {
		if f, ok := decl.(*ast.FuncDecl); ok {
			if f.Name.Name != "main" {
				continue
			}

			if len(f.Type.Params.List) != 0 {
				continue
			}

			if f.Type.Results != nil && len(f.Type.Results.List) != 0 {
				continue
			}

			return true
		}
	}

	return false
}

func cutSpace(b []byte) (before, middle, after []byte) {
	i := 0
	for i < len(b) && (b[i] == ' ' || b[i] == '\t' || b[i] == '\n') {
		i++
	}
	j := len(b)
	for j > 0 && (b[j-1] == ' ' || b[j-1] == '\t' || b[j-1] == '\n') {
		j--
	}
	if i <= j {
		return b[:i], b[i:j], b[j:]
	}
	return nil, nil, b[j:]
}

// matchSpace reformats src to use the same space context as orig.
// 1) If orig begins with blank lines, matchSpace inserts them at the beginning of src.
// 2) matchSpace copies the indentation of the first non-blank line in orig
//    to every non-blank line in src.
// 3) matchSpace copies the trailing space from orig and uses it in place
//   of src's trailing space.
func matchSpace(orig []byte, src []byte) []byte {
	before, _, after := cutSpace(orig)
	i := bytes.LastIndex(before, []byte{'\n'})
	before, indent := before[:i+1], before[i+1:]

	_, src, _ = cutSpace(src)

	var b bytes.Buffer
	_, _ = b.Write(before)
	for len(src) > 0 {
		line := src
		if i := bytes.IndexByte(line, '\n'); i >= 0 {
			line, src = line[:i+1], line[i+1:]
		} else {
			src = nil
		}
		if len(line) > 0 && line[0] != '\n' {	// not blank
			_, _ = b.Write(indent)
		}
		_, _ = b.Write(line)
	}
	_, _ = b.Write(after)
	return b.Bytes()
}

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
					_ = out.WriteByte('\n')
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
