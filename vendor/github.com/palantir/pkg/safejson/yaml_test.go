// Copyright 2016 Palantir Technologies. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package safejson_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v2"

	"github.com/palantir/pkg/safejson"
)

var tests = []struct {
	input  map[interface{}]interface{}
	output map[string]interface{}
}{
	{
		input: map[interface{}]interface{}{
			"hello": "world",
			"123":   123,
			"foo": map[string]interface{}{
				"bar": 1,
				"baz": 2,
			},
		},
		output: map[string]interface{}{
			"hello": "world",
			"123":   123,
			"foo": map[string]interface{}{
				"bar": 1,
				"baz": 2,
			},
		},
	},
	{
		input: map[interface{}]interface{}{
			"1":   "one",
			"two": "2",
			"other_map": map[interface{}]interface{}{
				"sky":   "blue",
				"grass": "green",
			},
		},
		output: map[string]interface{}{
			"1":   "one",
			"two": "2",
			"other_map": map[string]interface{}{
				"sky":   "blue",
				"grass": "green",
			},
		},
	}, {
		input: map[interface{}]interface{}{
			"array": []interface{}{
				map[interface{}]interface{}{
					"a": "b",
					"b": "c",
					"c": "d",
				},
			},
		},
		output: map[string]interface{}{
			"array": []interface{}{
				map[string]interface{}{
					"a": "b",
					"b": "c",
					"c": "d",
				},
			},
		},
	},
	{
		input: map[interface{}]interface{}{
			"array": nil,
		},
		output: map[string]interface{}{
			"array": nil,
		},
	},
}

var yamlTests = []struct {
	input  string
	output map[string]interface{}
}{
	{
		input: `---
x:
  z: 0
`,
		output: map[string]interface{}{
			"x": map[string]interface{}{
				"z": 0,
			},
		},
	},
}

func TestFromYAML(t *testing.T) {
	for _, test := range tests {
		out, err := safejson.FromYAMLValue(test.input)
		require.NoError(t, err)
		assert.Equal(t, test.output, out)
	}

	invalidJSONMap := map[interface{}]interface{}{
		"two": "2",
		"other_map": map[interface{}]interface{}{
			1:       "one",
			"sky":   "blue",
			"grass": "green",
		},
	}
	out, err := safejson.FromYAMLValue(invalidJSONMap)
	assert.EqualError(t, err, "Expected map key inside other_map to be a string but was int: 1")
	assert.Nil(t, out)

	for _, test := range yamlTests {
		var y interface{}
		err := yaml.Unmarshal([]byte(test.input), &y)
		if assert.NoError(t, err) {
			j, err := safejson.FromYAMLValue(y)
			assert.NoError(t, err)
			assert.Equal(t, test.output, j)
		}
	}
}

func TestMapInStructsNotConverted(t *testing.T) {
	val := struct {
		v map[interface{}]string
	}{
		v: map[interface{}]string{
			13: "thirteen",
		},
	}

	res, err := safejson.FromYAMLValue(val)
	require.NoError(t, err)
	assert.Equal(t, val, res)
}
