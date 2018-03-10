// Copyright 2016 Palantir Technologies. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build !go1.7

package safejson

import (
	"bytes"
	"encoding/json"
	"io"
)

func Encoder(w io.Writer) *json.Encoder {
	return json.NewEncoder(unescapedHTMLWriter{w})
}

type unescapedHTMLWriter struct {
	w io.Writer
}

// Write unescapes the HTML contents of p and writes it to the underlying writer.
func (u unescapedHTMLWriter) Write(p []byte) (n int, err error) {
	return u.w.Write(htmlUnescape(p))
}

// htmlUnescape returns a copy of the slice s with unescaped special HTML
// characters like <, >, and &.
//
// Warning: this allocates 3 additional copies of the slice s.
func htmlUnescape(s []byte) []byte {
	s = bytes.Replace(s, []byte("\\u003c"), []byte("<"), -1)
	s = bytes.Replace(s, []byte("\\u003e"), []byte(">"), -1)
	return bytes.Replace(s, []byte("\\u0026"), []byte("&"), -1)
}
