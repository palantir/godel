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

package docker

import (
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/palantir/godel/apps/distgo/params"
)

func Test_ProductsToDistAndBuild_ValidInput(t *testing.T) {
	project := params.Project{
		Products: map[string]params.Product{
			"foo": {
				DockerImages: []params.DockerImage{
					{
						Deps: []params.DockerDep{
							{
								Product: "bar",
								Type:    params.DockerDepDocker,
							},
						},
					},
				},
			},
			"bar": {
				DockerImages: []params.DockerImage{
					{
						Deps: []params.DockerDep{
							{
								Product: "bar",
								Type:    params.DockerDepSLS,
							},
							{
								Product: "baz",
								Type:    params.DockerDepDocker,
							},
						},
					},
				},
			},
			"baz": {
				DockerImages: []params.DockerImage{
					{
						Deps: []params.DockerDep{
							{
								Product: "baz",
								Type:    params.DockerDepSLS,
							},
						},
					},
				},
			},
		},
	}
	distProducts, imageProducts, err := productsToDistAndBuildImage([]string{"foo"}, project)
	require.NoError(t, err)
	sort.Strings(distProducts)
	sort.Strings(imageProducts)
	require.Equal(t, []string{"bar", "baz"}, distProducts)
	require.Equal(t, []string{"bar", "baz", "foo"}, imageProducts)
}

func Test_ProductsToDistAndBuild_InvalidProduct(t *testing.T) {
	project := params.Project{
		Products: map[string]params.Product{
			"abc": {
				DockerImages: []params.DockerImage{
					{
						Repository: "test-repo",
					},
				},
			},
			"foo": {
				DockerImages: []params.DockerImage{
					{
						Repository: "test-repo",
					},
				},
			},
			"xyz": {
				DockerImages: []params.DockerImage{
					{
						Repository: "test-repo",
					},
				},
			},
		},
	}
	_, _, err := productsToDistAndBuildImage([]string{"baz", "bar", "abc"}, project)
	require.Error(t, err)
	assert.EqualError(t, err, `Invalid products: [bar baz]. Valid products are: [abc foo xyz]`)
}
