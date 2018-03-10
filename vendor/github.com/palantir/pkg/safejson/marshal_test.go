// Copyright 2016 Palantir Technologies. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package safejson_test

import (
	"bytes"
	"math/big"
	"testing"

	"github.com/palantir/pkg/safejson"
)

// encodeTests are shared by TestEncoder and TestMarshal to guarantee that their
// behaviors are consistent.
var encodeTests = map[string]struct {
	in   interface{}
	want string
}{
	"unescaped HTML characters": {
		in:   `< this string contains & HTML characters that should not be escaped>`,
		want: `"< this string contains & HTML characters that should not be escaped>"`,
	},
	"big.Float as json.Number": {
		in:   big.NewFloat(3.14),
		want: `"3.14"`,
	},
	"struct containing *big.Float": {
		in:   struct{ Foo *big.Float }{Foo: big.NewFloat(3.14)},
		want: `{"Foo":"3.14"}`,
	},
	"slice of *big.Float": {
		in:   []*big.Float{big.NewFloat(3.14), big.NewFloat(8.42)},
		want: `["3.14","8.42"]`,
	},
}

func TestEncoder(t *testing.T) {
	for name, tt := range encodeTests {
		// Test Encoder
		var got bytes.Buffer
		if err := safejson.Encoder(&got).Encode(tt.in); err != nil {
			t.Errorf("failed to encode %s: %v", name, tt.in)
			continue
		}
		// json.Encoder writes single newline after writing encoded JSON
		wantEncoded := tt.want + "\n"
		if got.String() != wantEncoded {
			t.Errorf("wrong encoding for %s:\ngot:    %q\nwanted: %q", name, got.String(), wantEncoded)
		}
	}
}
