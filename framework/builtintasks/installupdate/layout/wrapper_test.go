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

package layout_test

import (
	"io/ioutil"
	"os"
	"path"
	"testing"

	"github.com/nmiyake/pkg/dirs"
	"github.com/palantir/pkg/specdir"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/palantir/godel/framework/builtintasks/installupdate/layout"
)

func TestWrapperLayoutNoValidation(t *testing.T) {
	specDir, err := specdir.New("testRoot", layout.WrapperSpec(), nil, specdir.SpecOnly)
	require.NoError(t, err)

	for i, currCase := range []struct {
		aliasName string
		want      string
	}{
		{
			aliasName: layout.WrapperConfigDir,
			want:      "testRoot/godel/config",
		},
		{
			aliasName: layout.WrapperAppDir,
			want:      "testRoot/godel",
		},
		{
			aliasName: layout.WrapperScriptFile,
			want:      "testRoot/godelw",
		},
	} {
		assert.Equal(t, currCase.want, specDir.Path(currCase.aliasName), "Case %d", i)
	}
}

func TestWrapperLayoutValidationFail(t *testing.T) {
	spec := layout.WrapperSpec()
	err := spec.Validate("testRoot", nil)
	assert.EqualError(t, err, "testRoot/godelw does not exist")
}

func TestWrapperLayoutValidation(t *testing.T) {
	tmpDir, cleanup, err := dirs.TempDir("", "")
	defer cleanup()
	require.NoError(t, err)

	err = os.MkdirAll(path.Join(tmpDir, "wrapperParent", "godel", "config"), 0755)
	require.NoError(t, err)

	err = os.MkdirAll(path.Join(tmpDir, "wrapperParent", "godel", ".src"), 0755)
	require.NoError(t, err)

	err = ioutil.WriteFile(path.Join(tmpDir, "wrapperParent", "godelw"), []byte("test file"), 0644)
	require.NoError(t, err)

	specDir, err := specdir.New(path.Join(tmpDir, "wrapperParent"), layout.WrapperSpec(), nil, specdir.Validate)
	require.NoError(t, err)

	for i, currCase := range []struct {
		aliasName string
		want      string
	}{
		{
			aliasName: layout.WrapperAppDir,
			want:      "wrapperParent/godel",
		},
		{
			aliasName: layout.WrapperConfigDir,
			want:      "wrapperParent/godel/config",
		},
		{
			aliasName: layout.WrapperScriptFile,
			want:      "wrapperParent/godelw",
		},
	} {
		assert.Equal(t, path.Join(tmpDir, currCase.want), specDir.Path(currCase.aliasName), "Case %d", i)
	}
}
