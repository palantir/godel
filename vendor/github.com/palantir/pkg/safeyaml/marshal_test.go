// Copyright 2016 Palantir Technologies. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package safeyaml

import (
	"encoding/json"
	"fmt"
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
)

type testStruct struct {
	JSONNum json.Number
	BigInt  *big.Int
	Various various
	Nils    various
}

type various struct {
	Slice     []interface{}
	Map       map[string]interface{}
	Pointer   *interface{}
	Interface fmt.Stringer
}

func TestMarshal(t *testing.T) {
	for _, testcase := range []struct {
		name     string
		input    interface{}
		expected string
	}{
		{
			name: "test that no number type is quoted",
			input: map[string]interface{}{
				"bigfloat": big.NewFloat(1234),
				"bigint":   big.NewInt(1234),
				"int":      1234,
				"json":     json.Number("1234"),
			},
			expected: "" +
				"bigfloat: 1234\n" +
				"bigint: 1234\n" +
				"int: 1234\n" +
				"json: 1234\n",
		},
		{
			name: "test that nested numbers are not quoted when possible",
			input: testStruct{
				JSONNum: json.Number("1234"),
				BigInt:  big.NewInt(1234),
				Various: various{
					Slice: []interface{}{
						json.Number("1234"),
					},
					Map: map[string]interface{}{
						"bigint": big.NewInt(1234),
					},
					Pointer:   &[]interface{}{json.Number("1234")}[0],
					Interface: big.NewInt(1234),
				},
				Nils: various{
					Slice:     nil,
					Map:       nil,
					Pointer:   nil,
					Interface: nil,
				},
			},
			expected: "" +
				"jsonnum: \"1234\"\n" + // required to stay json.Number
				"bigint: \"1234\"\n" + // required to stay *big.Int
				"various:\n" +
				"  slice:\n" +
				"  - 1234\n" +
				"  map:\n" +
				"    bigint: 1234\n" +
				"  pointer: 1234\n" +
				"  interface: \"1234\"\n" + // fmt.Stringer has more than 0 methods
				"nils:\n" +
				"  slice: []\n" +
				"  map: {}\n" +
				"  pointer: null\n" +
				"  interface: null\n",
		},
	} {
		actual, err := Marshal(testcase.input)
		if assert.NoError(t, err, testcase.name) {
			assert.Equal(t, testcase.expected, string(actual), testcase.name)
		}
	}
}
