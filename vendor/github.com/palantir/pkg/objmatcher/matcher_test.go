// Copyright 2016 Palantir Technologies. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package objmatcher_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/palantir/pkg/objmatcher"
)

type testStruct struct {
	label string
	num   int
}

func TestEqualsMatcher(t *testing.T) {
	for i, currCase := range []struct {
		name        string
		matcherWant interface{}
		given       interface{}
		wantErr     string
	}{
		{
			name:        "strings match",
			matcherWant: "foo",
			given:       "foo",
		},
		{
			name:        "strings mismatch",
			matcherWant: "foo",
			given:       "bar",
			wantErr:     "want: string(foo)\ngot:  string(bar)",
		},
		{
			name:        "strings mismatch",
			matcherWant: "foo",
			given:       5,
			wantErr:     "want: string(foo)\ngot:  int(5)",
		},
		{
			name: "structs match",
			matcherWant: testStruct{
				label: "foo",
				num:   13,
			},
			given: testStruct{
				label: "foo",
				num:   13,
			},
		},
		{
			name: "structs mismatch",
			matcherWant: testStruct{
				label: "foo",
				num:   13,
			},
			given: testStruct{
				label: "bar",
				num:   13,
			},
			wantErr: "want: objmatcher_test.testStruct({label:foo num:13})\ngot:  objmatcher_test.testStruct({label:bar num:13})",
		},
		{
			name: "maps match",
			matcherWant: map[string]interface{}{
				"outer-foo": "bar",
				"outer-num": 5,
				"struct": testStruct{
					label: "inner-bar",
					num:   13,
				},
			},
			given: map[string]interface{}{
				"outer-foo": "bar",
				"outer-num": 5,
				"struct": testStruct{
					label: "inner-bar",
					num:   13,
				},
			},
		},
	} {
		gotErr := objmatcher.NewEqualsMatcher(currCase.matcherWant).Matches(currCase.given)
		if currCase.wantErr == "" {
			assert.NoError(t, gotErr, "Case %d: %v", i, currCase.name)
		} else {
			assert.EqualError(t, gotErr, currCase.wantErr, "Case %d: %v", i, currCase.name)
		}
	}
}
