// Copyright 2016 Palantir Technologies. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package safejson

import (
	"encoding/json"
	"io"
)

// Decoder returns a new *json.Decoder with UseNumber enabled.
func Decoder(r io.Reader) *json.Decoder {
	decoder := json.NewDecoder(r)
	decoder.UseNumber()
	return decoder
}
