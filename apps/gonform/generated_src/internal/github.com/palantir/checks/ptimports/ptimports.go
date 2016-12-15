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

// Based on golang.org/x/tools/cmd/goimports which bears the following license:
//
// Copyright 2013 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package amalgomated

import (
	"bytes"
	"github.com/palantir/godel/apps/gonform/generated_src/internal/github.com/palantir/checks/ptimports/amalgomated_flag"
	"fmt"
	"go/scanner"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/palantir/godel/apps/gonform/generated_src/internal/github.com/palantir/checks/ptimports/ptimports"
)

var (
	exitCode	= 0
	list		= flag.Bool("l", false, "list files whose formatting differs from ptimport's")
	write		= flag.Bool("w", false, "Do not print reformatted sources to standard output. If a file's formatting is different from ptimports's, overwrite it with ptimports's version.")
)

func report(err error) {
	scanner.PrintError(os.Stderr, err)
	exitCode = 2
}

func usage() {
	fmt.Fprintf(os.Stderr, "usage: ptimports [flags] [path...]\n")
	flag.PrintDefaults()
	os.Exit(2)
}

func isGoFile(f os.FileInfo) bool {
	// ignore non-Go files
	name := f.Name()
	return !f.IsDir() && !strings.HasPrefix(name, ".") && strings.HasSuffix(name, ".go")
}

func processFile(filename string, in io.Reader) error {
	if in == nil {
		f, err := os.Open(filename)
		if err != nil {
			return err
		}
		defer func() {
			_ = f.Close()
		}()
		in = f
	}

	src, err := ioutil.ReadAll(in)
	if err != nil {
		return err
	}

	res, err := ptimports.Process(filename, src)
	if err != nil {
		return err
	}

	if *list {
		if !bytes.Equal(src, res) {
			fmt.Println(filename)
		}
		return nil
	}

	if *write {
		// only write when file changed
		if !bytes.Equal(src, res) {
			return ioutil.WriteFile(filename, res, 0)
		}
	} else {
		// print regardless of whether they are equal
		fmt.Print(string(res))
	}
	return nil
}

func visitFile(path string, f os.FileInfo, err error) error {
	if err == nil && isGoFile(f) {
		err = processFile(path, nil)
	}
	if err != nil {
		report(err)
	}
	if f.IsDir() && shouldSkipDir(f.Name()) {
		return filepath.SkipDir
	}
	return nil
}

func shouldSkipDir(name string) bool {
	return name == "Godeps" || name == "vendor"
}

func AmalgomatedMain() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	// call gofmtMain in a separate function
	// so that it can use defer and have them
	// run before the exit.
	gofmtMain()
	os.Exit(exitCode)
}

func gofmtMain() {
	flag.Usage = usage
	flag.Parse()
	paths := flag.Args()

	if len(paths) == 0 {
		usage()
	}

	for _, path := range paths {
		switch dir, err := os.Stat(path); {
		case err != nil:
			report(err)
		case dir.IsDir():
			if err := filepath.Walk(path, visitFile); err != nil {
				report(err)
			}
		default:
			if err := processFile(path, nil); err != nil {
				report(err)
			}
		}
	}
}
