// Copyright 2016 Palantir Technologies. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build go1.7

package safejson

import (
	"encoding/json"
	"io"
)

// Encoder returns a new *json.Encoder with SetEscapeHTML(false).
func Encoder(w io.Writer) *json.Encoder {
	e := json.NewEncoder(w)
	e.SetEscapeHTML(false)
	return e
}
