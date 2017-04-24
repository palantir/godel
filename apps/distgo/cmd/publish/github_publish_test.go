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

package publish

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUploadURIForProduct(t *testing.T) {
	for i, currCase := range []struct {
		templateURL string
		name        string
		want        string
	}{
		{
			"https://uploads.github.com/repos/octocat/Hello-World/releases/1/assets{?name,label}",
			"foo.tgz",
			"https://uploads.github.com/repos/octocat/Hello-World/releases/1/assets?name=foo.tgz",
		},
		{
			"https://github.enterprise.host/api/uploads/repos/octocat/Hello-World/releases/1/assets{?name,label}",
			"foo.tgz",
			"https://github.enterprise.host/api/uploads/repos/octocat/Hello-World/releases/1/assets?name=foo.tgz",
		},
		{
			"https://github.enterprise.host/api/uploads/repos/octocat/Hello-World/releases/1/assets{?name,label}",
			"unsafe-!-chars-@-name.tgz",
			"https://github.enterprise.host/api/uploads/repos/octocat/Hello-World/releases/1/assets?name=unsafe-%21-chars-%40-name.tgz",
		},
	} {
		got, err := uploadURIForProduct(currCase.templateURL, currCase.name)
		require.NoError(t, err, "Case %d", i)
		assert.Equal(t, currCase.want, got, "Case %d", i)
	}
}
