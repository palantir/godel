// Copyright 2016 Palantir Technologies. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

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
