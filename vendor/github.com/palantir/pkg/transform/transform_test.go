// Copyright 2016 Palantir Technologies. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package transform

import (
	"encoding/json"
	"fmt"
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Satisfied by *big.Int and *big.Float
type Signer interface {
	Sign() int
}

type testStruct struct {
	BigFloat  *big.Float
	Stringer  fmt.Stringer
	Signer    Signer
	Interface interface{}
}

func TestTransform(t *testing.T) {
	for i, test := range []struct {
		name     string
		rules    Rules
		input    interface{}
		expected interface{}
	}{
		{
			name:  "basic transforms",
			rules: Rules{jsonNumberToString, bigFloatToJSONNumber},
			input: map[string]interface{}{
				"number": json.Number("0"),
				"int":    big.NewInt(1),
				"float":  big.NewFloat(2),
			},
			expected: map[string]interface{}{
				"number": "0",
				"int":    big.NewInt(1),
				"float":  json.Number("2"),
			},
		},
		{
			name:     "top-level nil interface",
			rules:    Rules{jsonNumberToNil},
			input:    fmt.Stringer(nil),
			expected: fmt.Stringer(nil),
		},
		{
			name:  "nil values in map",
			rules: Rules{jsonNumberToNil},
			input: map[string]interface{}{
				"nil":    nil,
				"number": json.Number("0"),
				"slice of interface{}": []interface{}{
					json.Number("0"),
				},
				"slice of json.Number": []json.Number{
					json.Number("0"),
				},
			},
			expected: map[string]interface{}{
				"nil":    nil,
				"number": nil,
				"slice of interface{}": []interface{}{
					nil,
				},
				"slice of json.Number": []json.Number{
					json.Number("0"),
				},
			},
		},
		{
			name:  "transformed value assignable to field type",
			rules: Rules{bigFloatToJSONNumber},
			input: &testStruct{
				BigFloat:  big.NewFloat(1),
				Signer:    big.NewFloat(2),
				Stringer:  big.NewFloat(3),
				Interface: big.NewFloat(4),
			},
			expected: &testStruct{
				BigFloat:  big.NewFloat(1),
				Signer:    big.NewFloat(2),
				Stringer:  json.Number("3"),
				Interface: json.Number("4"),
			},
		},
		{
			name:  "array",
			rules: Rules{bigFloatToJSONNumber},
			input: [...]fmt.Stringer{
				nil,
				big.NewFloat(1),
				json.Number("2"),
			},
			expected: [...]fmt.Stringer{
				nil,
				json.Number("1"),
				json.Number("2"),
			},
		},
	} {
		result := test.rules.Apply(test.input)
		assert.Equal(t, test.expected, result, "case %d: %s", i, test.name)
	}
}

func jsonNumberToString(jn json.Number) string {
	return string(jn)
}

func bigFloatToJSONNumber(bf *big.Float) json.Number {
	return json.Number(bf.String())
}

func jsonNumberToNil(jn json.Number) interface{} {
	return nil
}
