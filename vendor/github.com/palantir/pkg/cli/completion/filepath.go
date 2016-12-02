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

package completion

import (
	"io/ioutil"
	"os"
	fp "path/filepath"
	"strings"
)

func List(partial string) []string {
	dir := fp.Dir(partial)

	list, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil
	}

	result := make([]string, 0, len(list))
	for _, finfo := range list {
		full := dirRaw(partial) + finfo.Name()
		if !strings.HasPrefix(strings.ToLower(full), strings.ToLower(partial)) {
			continue
		}
		if finfo.IsDir() {
			result = append(result, full+"/")
		} else {
			result = append(result, full)
		}
	}
	return result
}

func dirRaw(path string) string {
	i := len(path) - 1
	for i >= 0 && !os.IsPathSeparator(path[i]) {
		i--
	}
	return path[:i+1]
}
