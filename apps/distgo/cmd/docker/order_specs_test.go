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

package docker_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/palantir/godel/apps/distgo/cmd/docker"
	"github.com/palantir/godel/apps/distgo/params"
)

func generateSpec(product string, deps []params.DockerDep) params.ProductBuildSpec {
	return params.ProductBuildSpec{
		ProductName: product,
		Product: params.Product{
			DockerImages: []params.DockerImage{
				{
					Deps: deps,
				},
			},
		},
	}
}

func TestOrderBuildSpecs(t *testing.T) {
	A := generateSpec("A", []params.DockerDep{
		{Product: "B", Type: params.DockerDepDocker, TargetFile: ""},
		{Product: "C", Type: params.DockerDepDocker, TargetFile: ""},
	})
	B := generateSpec("B", []params.DockerDep{
		{Product: "D", Type: params.DockerDepDocker, TargetFile: ""},
	})
	C := generateSpec("C", []params.DockerDep{
		{Product: "D", Type: params.DockerDepDocker, TargetFile: ""},
	})
	D := generateSpec("D", []params.DockerDep{})
	E := generateSpec("E", []params.DockerDep{
		{Product: "DepE", Type: params.DockerDepDocker, TargetFile: ""},
	})
	DepE := generateSpec("DepE", []params.DockerDep{
		{Product: "E", Type: params.DockerDepDocker, TargetFile: ""},
	})

	X := generateSpec("X", []params.DockerDep{
		{Product: "Y", Type: params.DockerDepDocker, TargetFile: ""},
	})
	Y := generateSpec("Y", []params.DockerDep{
		{Product: "Z", Type: params.DockerDepDocker, TargetFile: ""},
	})
	Z := generateSpec("Z", []params.DockerDep{})

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
		actual, err := docker.OrderBuildSpecs(testcase.input)
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
