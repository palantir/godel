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

package dirs

import (
	"fmt"
	"io/ioutil"
	"os"
)

// TempDir creates a directory using ioutil.TempDir. If the ioutil.TempDir call is successful, returns its result and a
// function that removes the directory. The returned function is suitable for use in a defer call.
func TempDir(dir, prefix string) (string, func(), error) {
	path, err := ioutil.TempDir(dir, prefix)
	if err != nil {
		return "", nil, err
	}
	return path, RemoveAllFunc(path), nil
}

// RemoveAllFunc returns a function that calls os.RemoveAll on the specified path and prints any error that is
// encountered using fmt.Printf. Useful as a defer function to clean up directories created in tests.
func RemoveAllFunc(path string) func() {
	return func() {
		if err := os.RemoveAll(path); err != nil {
			fmt.Printf("Failed to remove directory %v in defer: %v", path, err)
		}
	}
}

// SetwdWithRestorer sets the working directory to be the specified directory and returns a function that restores the
// working directory to the value before it was changed. The returned function is suitable for use in a defer call.
func SetwdWithRestorer(dir string) (func(), error) {
	origWd, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	if err := os.Chdir(dir); err != nil {
		return nil, err
	}
	return func() {
		// restore working directory
		if err := os.Chdir(origWd); err != nil {
			fmt.Printf("failed to restore working directory to %s: %v", origWd, err)
		}
	}, nil
}
