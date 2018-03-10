// Copyright 2016 Palantir Technologies. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package safejson

import (
	"bytes"
)

// Unmarshal unmarshals the provided bytes (which should be valid JSON)
// into "v" using safejson.Decoder.
func Unmarshal(data []byte, v interface{}) error {
	return Decoder(bytes.NewReader(data)).Decode(v)
}
