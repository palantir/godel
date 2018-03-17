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

package products_test

import (
	"os/exec"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/palantir/godel/pkg/v2/products"
)

func TestList(t *testing.T) {
	p, err := products.List()
	require.NoError(t, err)
	assert.Equal(t, []string{"godel"}, p)
}

func TestBin(t *testing.T) {
	// TOOD: re-enable after project gödel is updated to 2.0 (https://github.com/palantir/godel/issues/301)
	t.SkipNow()

	bin, err := products.Bin("godel")
	require.NoError(t, err)
	cmd := exec.Command(bin, "version")
	output, err := cmd.CombinedOutput()
	require.NoError(t, err, "Command %v failed.\nOutput:\n%s\nError:\n%v", cmd.Args, string(output), err)

	assert.True(t, strings.HasPrefix(string(output), "godel version"))
}

func TestDist(t *testing.T) {
	dist, err := products.Dist("godel")
	require.NoError(t, err)
	assert.True(t, strings.HasSuffix(dist, ".tgz"))
}
