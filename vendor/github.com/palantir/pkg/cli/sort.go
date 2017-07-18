// Copyright 2016 Palantir Technologies. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

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
