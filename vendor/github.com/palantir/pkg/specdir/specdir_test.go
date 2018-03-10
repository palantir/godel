// Copyright 2016 Palantir Technologies. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package specdir_test

import (
	"io/ioutil"
	"os"
	"path"
	"regexp"
	"testing"

	"github.com/nmiyake/pkg/dirs"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/palantir/pkg/specdir"
)

func TestSpecDirConstruction(t *testing.T) {
	tmpDir, cleanup, err := dirs.TempDir("", "")
	defer cleanup()
	require.NoError(t, err)

	for i, currCase := range []struct {
		rootDir       string
		spec          specdir.LayoutSpec
		pathsToCreate map[string]specdir.PathType
		expectedError string
	}{
		{
			rootDir: "root",
			spec:    specdir.NewLayoutSpec(specdir.Dir(specdir.LiteralName("root"), ""), true),
			pathsToCreate: map[string]specdir.PathType{
				"root": specdir.DirPath,
			},
		},
		{
			rootDir:       "missing",
			spec:          specdir.NewLayoutSpec(specdir.Dir(specdir.LiteralName("root"), ""), true),
			expectedError: `^.+/missing is not a path to root$`,
		},
	} {
		currCaseTmpDir, err := ioutil.TempDir(tmpDir, "")
		require.NoError(t, err)

		rootForCreation := path.Join(currCaseTmpDir, currCase.rootDir)
		err = os.Mkdir(rootForCreation, 0755)
		require.NoError(t, err)

		createDirectoryStructure(t, currCaseTmpDir, currCase.pathsToCreate)

		_, err = specdir.New(rootForCreation, currCase.spec, nil, specdir.Validate)
		if currCase.expectedError == "" {
			assert.NoError(t, err, "Case %d", i)
		} else {
			assert.Regexp(t, regexp.MustCompile(currCase.expectedError), err.Error(), "Case %d", i)
		}
	}
}

func TestSpecDirCreateMode(t *testing.T) {
	tmpDir, cleanup, err := dirs.TempDir("", "")
	defer cleanup()
	require.NoError(t, err)

	spec := specdir.NewLayoutSpec(
		specdir.Dir(specdir.LiteralName("root"), "",
			specdir.Dir(specdir.LiteralName("child"), ""),
		), true)
	rootForCreation := path.Join(tmpDir, "root")
	_, err = specdir.New(rootForCreation, spec, nil, specdir.Create)
	require.NoError(t, err)

	// use SpecDir to do validation of creation
	_, err = specdir.New(rootForCreation, spec, nil, specdir.Validate)
	assert.NoError(t, err)
}

func TestSpecDirGetPath(t *testing.T) {
	tmpDir, cleanup, err := dirs.TempDir("", "")
	defer cleanup()
	require.NoError(t, err)

	for i, currCase := range []struct {
		rootDir       string
		spec          specdir.LayoutSpec
		pathsToCreate map[string]specdir.PathType
		alias         string
		expectedPath  string
	}{
		{
			rootDir: "root",
			spec: specdir.NewLayoutSpec(specdir.Dir(specdir.LiteralName("root"), "",
				specdir.Dir(specdir.LiteralName("child"), "",
					specdir.Dir(specdir.LiteralName("grandchild"), "VeryInnerDir"))), true),
			pathsToCreate: map[string]specdir.PathType{
				"root":                  specdir.DirPath,
				"root/child":            specdir.DirPath,
				"root/child/grandchild": specdir.DirPath,
			},
			expectedPath: "root/child/grandchild",
		},
		{
			rootDir: "root",
			spec: specdir.NewLayoutSpec(specdir.Dir(specdir.LiteralName("root"), "",
				specdir.Dir(specdir.LiteralName("child"), "",
					specdir.Dir(specdir.LiteralName("grandchild"), "VeryInnerDir"))), false),
			pathsToCreate: map[string]specdir.PathType{
				"root":                  specdir.DirPath,
				"root/child":            specdir.DirPath,
				"root/child/grandchild": specdir.DirPath,
			},
			expectedPath: "root/child/grandchild",
		},
	} {
		currCaseTmpDir, err := ioutil.TempDir(tmpDir, "")
		require.NoError(t, err)

		rootForCreation := path.Join(currCaseTmpDir, currCase.rootDir)
		err = os.Mkdir(rootForCreation, 0755)
		require.NoError(t, err)

		createDirectoryStructure(t, currCaseTmpDir, currCase.pathsToCreate)

		specDir, err := specdir.New(rootForCreation, currCase.spec, nil, specdir.Validate)
		require.NoError(t, err)

		actualPath := specDir.Path("VeryInnerDir")

		assert.Equal(t, path.Join(currCaseTmpDir, currCase.expectedPath), actualPath, "Case %d", i)
	}
}
