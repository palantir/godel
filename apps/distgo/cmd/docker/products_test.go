package docker

import (
	"sort"
	"testing"

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
			"foo": {
				DockerImages: []params.DockerImage{
					{
						Repository: "test-repo",
					},
				},
			},
		},
	}
	_, _, err := productsToDistAndBuildImage([]string{"bar"}, project)
	require.Error(t, err)
}
