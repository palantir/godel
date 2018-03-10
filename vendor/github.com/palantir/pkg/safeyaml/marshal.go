// Copyright 2016 Palantir Technologies. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package safeyaml

import (
	"encoding/json"
	"math/big"
	"strconv"

	"gopkg.in/yaml.v2"

	"github.com/palantir/pkg/transform"
)

func Marshal(in interface{}) (out []byte, err error) {
	converted := numbersToPrimitives(in)
	return yaml.Marshal(converted)
}

func numbersToPrimitives(in interface{}) interface{} {
	// turn instances of json.Number, big.Int, big.Float into primitive numbers
	rules := transform.Rules{
		func(jn json.Number) interface{} {
			if primitive, ok := parsePrimitive(string(jn)); ok {
				return primitive
			}
			return jn
		},
		func(bi *big.Int) interface{} {
			if primitive, ok := parsePrimitive(bi.String()); ok {
				return primitive
			}
			return bi
		},
		func(bf *big.Float) interface{} {
			if primitive, ok := parsePrimitive(bf.Text('g', -1)); ok {
				return primitive
			}
			return bf
		},
	}
	return rules.Apply(in)
}

func parsePrimitive(str string) (prim interface{}, ok bool) {
	i64, err := strconv.ParseInt(str, 0, 64)
	if err == nil {
		return i64, true
	}
	f64, err := strconv.ParseFloat(str, 64)
	if err == nil {
		return f64, true
	}
	return nil, false
}
