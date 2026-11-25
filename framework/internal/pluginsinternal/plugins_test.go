// Copyright 2026 Palantir Technologies, Inc.
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

package pluginsinternal_test

import (
	"io"
	"testing"

	"github.com/palantir/godel/v2/framework/artifactresolver"
	"github.com/palantir/godel/v2/framework/internal/pluginsinternal"
	"github.com/palantir/godel/v2/pkg/osarch"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestVerifyAssetVersions(t *testing.T) {
	for _, tc := range []struct {
		name    string
		assets  []artifactresolver.LocatorWithResolverParam
		wantErr string
	}{
		{
			name: "no assets",
		},
		{
			name: "different assets",
			assets: []artifactresolver.LocatorWithResolverParam{
				assetParam("com.palantir", "foo-asset", "1.0.0"),
				assetParam("com.palantir", "bar-asset", "2.0.0"),
				assetParam("com.palantir.other", "foo-asset", "3.0.0"),
			},
		},
		{
			name: "same locator specified multiple times",
			assets: []artifactresolver.LocatorWithResolverParam{
				assetParam("com.palantir", "foo-asset", "1.0.0"),
				assetParam("com.palantir", "foo-asset", "1.0.0"),
			},
		},
		{
			name: "multiple versions of the same assets",
			assets: []artifactresolver.LocatorWithResolverParam{
				assetParam("com.palantir", "foo-asset", "1.0.0"),
				assetParam("com.palantir", "bar-asset", "1.0.0"),
				assetParam("com.palantir", "foo-asset", "2.0.0"),
				assetParam("com.palantir", "baz-asset", "1.0.0"),
				assetParam("com.palantir", "bar-asset", "0.1.0"),
				assetParam("com.palantir", "foo-asset", "1.0.0"),
			},
			wantErr: `2 asset(s) specified with multiple versions:
    com.palantir:bar-asset: 1.0.0, 0.1.0
    com.palantir:foo-asset: 1.0.0, 2.0.0`,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			err := pluginsinternal.VerifyAssetVersions(tc.assets)
			if tc.wantErr == "" {
				assert.NoError(t, err)
				return
			}
			assert.EqualError(t, err, tc.wantErr)
		})
	}
}

func TestResolveAssetsMultipleVersionsOfSameAsset(t *testing.T) {
	// The resolver and directories are invalid, so this only passes if the versions are verified before any
	// resolution is attempted.
	_, err := pluginsinternal.ResolveAssets(
		t.TempDir(),
		t.TempDir(),
		[]artifactresolver.LocatorWithResolverParam{
			assetParam("com.palantir", "foo-asset", "1.0.0"),
			assetParam("com.palantir", "foo-asset", "2.0.0"),
		},
		osarch.Current(),
		nil,
		io.Discard,
	)
	require.Error(t, err)
	assert.Equal(t, `1 asset(s) specified with multiple versions:
    com.palantir:foo-asset: 1.0.0, 2.0.0`, err.Error())
}

func assetParam(group, product, version string) artifactresolver.LocatorWithResolverParam {
	return artifactresolver.LocatorWithResolverParam{
		LocatorWithChecksums: artifactresolver.LocatorParam{
			Locator: artifactresolver.Locator{
				Group:   group,
				Product: product,
				Version: version,
			},
		},
	}
}
