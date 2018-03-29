// Copyright 2016 Palantir Technologies, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package generator

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetInputFiles(t *testing.T) {
	for i, tc := range []struct {
		in        []string
		want      []inputFile
		wantError string
	}{
		{
			[]string{
				"0_add.tmpl",
				"1_other.txt.tmpl",
				"invalid",
			},
			[]inputFile{
				{
					TemplateFileName:  "0_add.tmpl",
					Ordering:          0,
					Name:              "add",
					OriginalExtension: "",
				},
				{
					TemplateFileName:  "1_other.txt.tmpl",
					Ordering:          1,
					Name:              "other",
					OriginalExtension: ".txt",
				},
			},
			"",
		},
		{
			[]string{
				"0_add.tmpl",
				"1_add.txt.tmpl",
			},
			nil,
			`multiple inputs have the name "add": [0_add.tmpl 1_add.txt.tmpl]`,
		},
		{
			[]string{
				"0_add.tmpl",
				"0_other.txt.tmpl",
			},
			nil,
			`multiple inputs have the ordering value 0: [0_add.tmpl 0_other.txt.tmpl]`,
		},
	} {
		got, err := getInputFiles(tc.in)
		if tc.wantError == "" {
			require.NoError(t, err)
			assert.Equal(t, tc.want, got, "Case %d", i)
		} else {
			assert.EqualError(t, err, tc.wantError, "Case %d", i)
		}
	}
}
