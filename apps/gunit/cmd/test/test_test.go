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

package test

import (
	"io/ioutil"
	"testing"

	"github.com/nmiyake/pkg/dirs"
	"github.com/nmiyake/pkg/gofiles"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreatePlaceholderUsesBuildConstraints(t *testing.T) {
	testDir, cleanup, err := dirs.TempDir(".", "")
	defer cleanup()
	require.NoError(t, err)

	specs := []gofiles.GoFileSpec{
		{
			RelPath: "foo/main.go",
			Src: `// +build ignore

package main`,
		},
		{
			RelPath: "foo/zoo.go",
			Src:     `package zoo`,
		},
	}
	_, err = gofiles.Write(testDir, specs)
	require.NoError(t, err)

	writtenFiles, err := createPlaceholderTestFiles([]string{"foo"}, testDir)
	require.NoError(t, err)

	content, err := ioutil.ReadFile(writtenFiles[0])
	require.NoError(t, err)

	// generated placeholder should be of package "zoo" rather than package "main" because latter is ignored using
	// build constraint.
	want := `package zoo
// temporary placeholder test file created by gunit
`
	assert.Equal(t, want, string(content))
}
