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

package installupdate

import (
	"io/ioutil"
	"testing"

	"github.com/nmiyake/pkg/dirs"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	propertiesFileContent = `# Lines starting with '#' should be ignored
distributionURL=https://github.com/palantir/godel/godel-0.0.1.tgz
distributionSHA256=871ee79691aee47301214ab0ae10bd851e5bc0d48042ba8d33ac85cfcc7eb6cc
`
	urlValue      = "https://github.com/palantir/godel/godel-0.0.1.tgz"
	checksumValue = "871ee79691aee47301214ab0ae10bd851e5bc0d48042ba8d33ac85cfcc7eb6cc"
)

func TestReadProperties(t *testing.T) {
	tmpDir, cleanup, err := dirs.TempDir("", "")
	defer cleanup()
	require.NoError(t, err)

	for i, currCase := range []struct {
		propertiesFileWriter func(dir string) string
		want                 map[string]string
	}{
		{
			propertiesFileWriter: func(dir string) string {
				dest, err := ioutil.TempFile(dir, "")
				require.NoError(t, err)
				err = dest.Close()
				require.NoError(t, err)
				err = writePropertiesFile(dest.Name())
				require.NoError(t, err)
				return dest.Name()
			},
			want: map[string]string{
				propertiesURLKey:      urlValue,
				propertiesChecksumKey: checksumValue,
			},
		},
	} {
		propertiesFile := currCase.propertiesFileWriter(tmpDir)
		got, err := readPropertiesFile(propertiesFile)
		require.NoError(t, err)

		assert.Equal(t, currCase.want, got, "Case %d", i)
	}
}

func writePropertiesFile(filePath string) error {
	err := ioutil.WriteFile(filePath, []byte(propertiesFileContent), 0644)
	if err != nil {
		return errors.Wrapf(err, "failed to write file to path %s", filePath)
	}
	return nil
}
