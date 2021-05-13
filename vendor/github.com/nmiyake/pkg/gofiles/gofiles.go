// MIT License
//
// Copyright (c) 2016 Nick Miyake
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

// Package gofiles provides functions for specifying and writing Go source files.
package gofiles

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"golang.org/x/tools/go/packages"
)

// GoFileSpec represents the specification for a Go file.
type GoFileSpec struct {
	// The relative path to which the file should be written. For example, "foo/foo.go".
	RelPath string
	// Content of the file.
	Src string
}

// GoFile represents a Go file that has been written to disk.
type GoFile struct {
	// The absolute path to the Go file.
	Path string
	// The import path for the Go file. For example, "github.com/nmiyake/pkg/gofiles".
	ImportPath string
}

// Write the Go files represented by the specifications in the files parameter using the provided directory as the root
// directory.
//
// Returns a map of the written files where the key is the RelPath field of the specification that was written and the
// value is the GoFile that was written for the specification.
func Write(dir string, files []GoFileSpec) (map[string]GoFile, error) {
	dir, err := filepath.Abs(dir)
	if err != nil {
		return nil, err
	}

	// write all files
	goFiles := make(map[string]GoFile, len(files))
	for _, currFile := range files {
		filePath := filepath.Join(dir, currFile.RelPath)
		if err := os.MkdirAll(filepath.Dir(filePath), 0755); err != nil {
			return nil, err
		}
		if err := ioutil.WriteFile(filePath, []byte(currFile.Src), 0644); err != nil {
			return nil, err
		}
		goFiles[currFile.RelPath] = GoFile{
			Path: filePath,
		}
	}

	// after all files have been written, determine import path for each one.
	// Done after all files have been written because the set of files written may include "go.mod" files, which would
	// impact the import path.
	for _, currFile := range files {
		filePath := filepath.Join(dir, currFile.RelPath)
		fileDir := filepath.Dir(filePath)
		pkgs, err := packages.Load(&packages.Config{
			Dir: fileDir,
		}, ".")
		if err != nil {
			return nil, err
		}
		if len(pkgs) < 1 {
			return nil, fmt.Errorf("no packages found in %s", fileDir)
		}
		goFiles[currFile.RelPath] = GoFile{
			Path:       filePath,
			ImportPath: pkgs[0].PkgPath,
		}
	}

	return goFiles, nil
}
