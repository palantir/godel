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

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/palantir/godel/apps/distgo/cmd/dist"
	"github.com/palantir/godel/apps/distgo/params"
)

func generateSpec(product string, deps params.DockerDistDeps) params.ProductBuildSpec {
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

func TestOrderBuildSpecs(t *testing.T) {
	A := generateSpec("A", params.DockerDistDeps{
		{Product: "B", DistType: params.SLSDistType, TargetFile: ""},
		{Product: "C", DistType: params.BinDistType, TargetFile: ""},
	})
	B := generateSpec("B", params.DockerDistDeps{
		{Product: "D", DistType: params.DockerDistType, TargetFile: ""},
	})
	C := generateSpec("C", params.DockerDistDeps{
		{Product: "D", DistType: params.SLSDistType, TargetFile: ""},
	})
	D := generateSpec("D", params.DockerDistDeps{})
	E := generateSpec("E", params.DockerDistDeps{
		{Product: "DepE", DistType: params.SLSDistType, TargetFile: ""},
	})
	DepE := generateSpec("DepE", params.DockerDistDeps{
		{Product: "E", DistType: params.SLSDistType, TargetFile: ""},
	})

	X := generateSpec("X", params.DockerDistDeps{
		{Product: "Y", DistType: params.DockerDistType, TargetFile: ""},
	})
	Y := generateSpec("Y", params.DockerDistDeps{
		{Product: "Z", DistType: params.DockerDistType, TargetFile: ""},
	})
	Z := generateSpec("Z", params.DockerDistDeps{})

	for _, testcase := range []struct {
		input     []params.ProductBuildSpecWithDeps
		expected  []params.ProductBuildSpecWithDeps
		expectErr string
	}{
		{
			//  (A <- B,C <- D) = D, B, C, A
			input:    []params.ProductBuildSpecWithDeps{{Spec: A}, {Spec: B}, {Spec: C}, {Spec: D}},
			expected: []params.ProductBuildSpecWithDeps{{Spec: D}, {Spec: B}, {Spec: C}, {Spec: A}},
		},
		{
			// empty
			input:    []params.ProductBuildSpecWithDeps{},
			expected: []params.ProductBuildSpecWithDeps{},
		},
		{
			//  (E <- DepE <- E) = invalid
			input:     []params.ProductBuildSpecWithDeps{{Spec: E}, {Spec: DepE}},
			expected:  nil,
			expectErr: "Failed to generate ordering among the products.",
		},
		{
			//  (X <- Y <- Z) = Z, Y, X
			input:    []params.ProductBuildSpecWithDeps{{Spec: Y}, {Spec: X}, {Spec: Z}},
			expected: []params.ProductBuildSpecWithDeps{{Spec: Z}, {Spec: Y}, {Spec: X}},
		},
	} {
		actual, err := dist.OrderBuildSpecs(testcase.input)
		if testcase.expectErr != "" {
			require.Contains(t, err.Error(), testcase.expectErr)
			continue
		}
		require.NoError(t, err)
		for i, expectedSpec := range testcase.expected {
			assert.Equal(t, expectedSpec.Spec.ProductName, actual[i].Spec.ProductName)
		}
	}
}
