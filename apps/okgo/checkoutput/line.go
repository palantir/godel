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

package checkoutput

import (
	"fmt"

	"github.com/palantir/pkg/pkgpath"
	"github.com/pkg/errors"
)

type Issue struct {
	path    string
	line    int
	column  int
	message string
	baseDir string
}

func (i Issue) Path(pathType pkgpath.Type) (string, error) {
	if pathType == pkgpath.Relative {
		return i.path, nil
	}

	relPath := pkgpath.NewRelPkgPath(i.path, i.baseDir)
	switch pathType {
	case pkgpath.Absolute:
		return relPath.Abs(), nil
	case pkgpath.GoPathSrcRelative:
		return relPath.GoPathSrcRel()
	default:
		return "", errors.Errorf("unhandled type: %v", pathType)
	}
}

func (i Issue) Message() string {
	return i.message
}

func (i Issue) BaseDir() string {
	return i.baseDir
}

func (i Issue) String() string {
	filePart := ""
	if i.path != "" || i.line != 0 || i.column != 0 {
		// only include filePart if path, line or column is a non-default value
		filePart = fmt.Sprintf("%v:%v", i.path, i.line)
		if i.column != 0 {
			filePart = fmt.Sprintf("%v:%v", filePart, i.column)
		}
		filePart += ": "
	}
	return filePart + i.message
}
