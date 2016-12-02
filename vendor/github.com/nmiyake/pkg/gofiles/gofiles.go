/*
MIT License

Copyright (c) 2016 Nick Miyake

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
*/

// Package gofiles provides functions for specifying and writing Go source files.
package gofiles

import (
	"bytes"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"
	"text/template"
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
// directory. The Src field of each file is run as a Go template provided with the output of the "imports" function.
// This allows source files to import other files being written as part of the operation. For example:
//
//	"{{index . "bar/bar.go"}}" // resolves to the import path for the spec with RelPath "bar/bar.go"
//
// Returns a map of the written files where the key is the RelPath field of the specification that was written and the
// value is the GoFile that was written for the specification.
func Write(dir string, files []GoFileSpec) (map[string]GoFile, error) {
	dir, err := filepath.Abs(dir)
	if err != nil {
		return nil, err
	}

	imports, err := imports(dir, files)
	if err != nil {
		return nil, err
	}

	goFiles := make(map[string]GoFile, len(files))
	for _, currFile := range files {
		filePath := path.Join(dir, currFile.RelPath)
		buf := &bytes.Buffer{}
		t := template.Must(template.New(filePath).Parse(currFile.Src))
		if err := t.Execute(buf, imports); err != nil {
			return nil, err
		}
		if err := os.MkdirAll(path.Dir(filePath), 0755); err != nil {
			return nil, err
		}
		if err := ioutil.WriteFile(filePath, buf.Bytes(), 0644); err != nil {
			return nil, err
		}
		goFiles[currFile.RelPath] = GoFile{
			Path:       filePath,
			ImportPath: imports[currFile.RelPath],
		}
	}

	return goFiles, nil
}

// Returns a map that maps the relative path of each Go file in the provided specification to its package import path.
// For example:
//
// 	"./foo/foo.go": "github.com/nmiyake/pkg/gofiles/9012503/foo"
//
// If the relative path goes into a vendor directory, then the value will be the non-vendored import path. For example:
//
//	"./foo/vendor/github.com/nmiyake/bar/bar.go": "github.com/nmiyake/bar".
func imports(dir string, files []GoFileSpec) (map[string]string, error) {
	imports := make(map[string]string, len(files))
	for _, currFile := range files {
		fullDirPath := path.Dir(path.Join(dir, currFile.RelPath))
		importPath, err := filepath.Rel(path.Join(os.Getenv("GOPATH"), "src"), fullDirPath)
		if err != nil {
			return nil, err
		}

		vendorIndex := strings.LastIndex(importPath, "/vendor/")
		if vendorIndex != -1 {
			importPath = importPath[vendorIndex+len("/vendor/"):]
		}
		imports[currFile.RelPath] = importPath
	}
	return imports, nil
}
