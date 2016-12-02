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

package cli

import (
	"bytes"
	"io"
	"sort"
	"strings"
)

func (ctx Context) Sorted(f func()) {
	ctx.beginSort()
	defer ctx.endSort()
	f()
}

func (ctx Context) beginSort() {
	ctx.App.Stdout = sortedWriter{
		Buffer: new(bytes.Buffer),

		original: ctx.App.Stdout,
	}
}

func (ctx Context) endSort() {
	sw := ctx.App.Stdout.(sortedWriter) // panic if writer is wrong
	lines := strings.Split(sw.String(), "\n")
	if lines[len(lines)-1] == "" {
		lines = lines[:len(lines)-1] // account for trailing newline of last line
	}
	sort.Strings(lines)
	ctx.App.Stdout = sw.original
	for _, line := range lines {
		ctx.Println(line)
	}
}

type sortedWriter struct {
	*bytes.Buffer
	original io.Writer
}
