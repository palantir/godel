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
	"os"
	"strings"

	"github.com/pkg/errors"
)

const (
	propertiesURLKey      = "distributionURL"
	propertiesChecksumKey = "distributionSHA256"
)

// Reads the file at the provided path and returns the properties that it contains. The file should contain one property
// per line and the line should be of the form "key=value". Any line that starts with the character '#' is ignored.
func readPropertiesFile(path string) (map[string]string, error) {
	bytes, err := os.ReadFile(path)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to read file %s", path)
	}

	properties := make(map[string]string)
	lines := strings.Split(string(bytes), "\n")
	for _, currLine := range lines {
		currLine = strings.TrimSpace(currLine)

		if strings.HasPrefix(currLine, "#") || len(currLine) == 0 {
			continue
		}

		equalsIndex := strings.IndexAny(currLine, "=")
		if equalsIndex == -1 {
			return nil, errors.Errorf(`failed to find character "=" in line "%v" in file with lines "%v"`, currLine, lines)
		}
		properties[currLine[:equalsIndex]] = currLine[equalsIndex+1:]
	}
	return properties, nil
}
