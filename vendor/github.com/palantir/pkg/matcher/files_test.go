// Copyright 2016 Palantir Technologies. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package matcher_test

import (
	"io/ioutil"
	"os"
	"path"
	"testing"

	"github.com/nmiyake/pkg/dirs"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/palantir/pkg/matcher"
)

func TestListFiles(t *testing.T) {
	tmpDir, cleanup, err := dirs.TempDir("", "")
	defer cleanup()
	require.NoError(t, err)

	cases := []struct {
		include       matcher.Matcher
		exclude       matcher.Matcher
		filesToCreate map[string]string
		want          []string
	}{
		{
			include: matcher.Name(`.+\.go`),
			filesToCreate: map[string]string{
				"notgo.txt":        "",
				"isgo.go":          "",
				".hidden.go":       "",
				".hidden/inner.go": "",
				"indir/isgo.go":    "",
			},
			want: []string{
				".hidden/inner.go",
				".hidden.go",
				"indir/isgo.go",
				"isgo.go",
			},
		},
		{
			include: matcher.Name(`.+\.go`),
			exclude: matcher.Hidden(),
			filesToCreate: map[string]string{
				"notgo.txt":        "",
				"isgo.go":          "",
				".hidden.go":       "",
				"indir/isgo.go":    "",
				".hidden/inner.go": "",
			},
			want: []string{
				"indir/isgo.go",
				"isgo.go",
			},
		},
	}

	for i, currCase := range cases {
		currCaseTmpDir := createFiles(t, tmpDir, currCase.filesToCreate)

		got, err := matcher.ListFiles(currCaseTmpDir, currCase.include, currCase.exclude)
		require.NoError(t, err, "Case %d", i)

		assert.Equal(t, currCase.want, got, "Case %d", i)
	}
}

func createFiles(t *testing.T, tmpDir string, files map[string]string) string {
	currCaseTmpDir, err := ioutil.TempDir(tmpDir, "")
	require.NoError(t, err)

	for currFile, currContent := range files {
		err := os.MkdirAll(path.Join(currCaseTmpDir, path.Dir(currFile)), 0755)
		require.NoError(t, err)

		err = ioutil.WriteFile(path.Join(currCaseTmpDir, currFile), []byte(currContent), 0644)
		require.NoError(t, err)
	}

	return currCaseTmpDir
}
