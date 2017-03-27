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

package dist_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/palantir/godel/apps/distgo/cmd/dist"
	"github.com/palantir/godel/apps/distgo/params"
)

func generateSpec(product string, deps params.DockerDeps) params.ProductBuildSpec {
	return params.ProductBuildSpec{
		ProductName: product,
		Product: params.Product{
			Dist: []params.Dist{
				{
					Info: &params.DockerDistInfo{
						DistDeps: deps,
					},
				},
			},
		},
	}
}

func TestSchedule(t *testing.T) {

	A := generateSpec("A", params.DockerDeps{
		"B": {params.SLSDistType},
		"C": {params.BinDistType},
	})
	B := generateSpec("B", params.DockerDeps{
		"D": {params.DockerDistType},
	})
	C := generateSpec("C", params.DockerDeps{
		"D": {params.SLSDistType},
	})
	D := generateSpec("D", params.DockerDeps{})
	E := generateSpec("E", params.DockerDeps{
		"DepE": {params.SLSDistType},
	})
	DepE := generateSpec("DepE", params.DockerDeps{
		"E": {params.SLSDistType},
	})

	X := generateSpec("X", params.DockerDeps{
		"Y": {params.DockerDistType},
	})
	Y := generateSpec("Y", params.DockerDeps{
		"Z": {params.DockerDistType},
	})
	Z := generateSpec("Z", params.DockerDeps{})

	for _, testcase := range []struct {
		input    []params.ProductBuildSpecWithDeps
		expected []params.ProductBuildSpecWithDeps
		valid    bool
	}{
		{
			//  (A <- B,C <- D) = D, B, C, A
			input:    []params.ProductBuildSpecWithDeps{{Spec: A}, {Spec: B}, {Spec: C}, {Spec: D}},
			expected: []params.ProductBuildSpecWithDeps{{Spec: D}, {Spec: B}, {Spec: C}, {Spec: A}},
			valid:    true,
		},
		{
			// empty
			input:    []params.ProductBuildSpecWithDeps{},
			expected: []params.ProductBuildSpecWithDeps{},
			valid:    true,
		},
		{
			//  (E <- DepE <- E) = invalid
			input:    []params.ProductBuildSpecWithDeps{{Spec: E}, {Spec: DepE}},
			expected: nil,
			valid:    false,
		},
		{
			//  (X <- Y <- Z) = Z, Y, X
			input:    []params.ProductBuildSpecWithDeps{{Spec: Y}, {Spec: X}, {Spec: Z}},
			expected: []params.ProductBuildSpecWithDeps{{Spec: Z}, {Spec: Y}, {Spec: X}},
			valid:    true,
		},
	} {
		actual, err := dist.Schedule(testcase.input)
		if !testcase.valid {
			require.Error(t, err)
			continue

		}
		require.NoError(t, err)
		for i, expectedSpec := range testcase.expected {
			require.Equal(t, expectedSpec.Spec.ProductName, actual[i].Spec.ProductName)
		}
	}

}
