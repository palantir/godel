// Copyright 2016 Palantir Technologies. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package safejson_test

import (
	"encoding/json"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/palantir/pkg/safejson"
)

func TestMarshal(t *testing.T) {
	for name, tt := range encodeTests {
		// Test Marshal
		got, err := safejson.Marshal(tt.in)
		if err != nil {
			t.Errorf("failed to marshal %s: %v", name, tt.in)
			continue
		}
		if string(got) != tt.want {
			t.Errorf("wrong encoding for %s:\ngot:    %q\nwanted: %q", name, string(got), tt.want)
		}
	}
}

func TestConcurrentMarshal(t *testing.T) {
	var (
		old = 532173
		new = 23589217
	)
	got, err := safejson.Marshal(old)
	if err != nil {
		t.Fatalf("failed to marshal %d", old)
	}
	if _, err := safejson.Marshal(new); err != nil {
		t.Fatalf("failed to marshal %d", new)
	}
	want := strconv.Itoa(old)
	if string(got) != want {
		t.Errorf("buffer reuse:\ngot:    %q\nwanted: %q", string(got), want)
	}
}

func TestDecodeNumber(t *testing.T) {
	in := map[string]interface{}{
		"a": json.Number("12"),
		"b": json.Number("34"),
		"c": json.Number("56"),
	}

	var out struct {
		First  int         `json:"a"`
		Second json.Number `json:"b"`
		Third  interface{} `json:"c"`
	}

	jsonBytes, err := safejson.Marshal(in)
	require.NoError(t, err)

	err = safejson.Unmarshal(jsonBytes, &out)
	assert.NoError(t, err)

	assert.Equal(t, 12, out.First)
	assert.Equal(t, json.Number("34"), out.Second)
	assert.Equal(t, json.Number("56"), out.Third)
}
