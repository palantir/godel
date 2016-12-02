// Copyright 2016 Palantir Technologies, Inc.
// Copyright 2013 Kamil Kisiel
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

package outparamcheck

import (
	"fmt"
	"go/token"
	"strings"

	"github.com/dustin/go-humanize"
)

type OutParamError struct {
	Pos      token.Position
	Line     string
	Method   string
	Argument int
}

func (err OutParamError) Error() string {
	pos := err.Pos.String()
	// Trim prefix including /src/
	if i := strings.Index(pos, "/src/"); i != -1 {
		pos = pos[i+len("/src/"):]
	}

	line := err.Line
	comment := strings.Index(line, "//")
	if comment != -1 {
		line = line[:comment]
	}
	line = strings.TrimSpace(line)

	ord := humanize.Ordinal(err.Argument + 1)
	return fmt.Sprintf("%s\t%s  // %s argument of '%s' requires '&'", pos, line, ord, err.Method)
}

type byLocation []OutParamError

func (errs byLocation) Len() int {
	return len(errs)
}

func (errs byLocation) Swap(i, j int) {
	errs[i], errs[j] = errs[j], errs[i]
}

func (errs byLocation) Less(i, j int) bool {
	ei, ej := errs[i], errs[j]
	pi, pj := ei.Pos, ej.Pos
	if pi.Filename != pj.Filename {
		return pi.Filename < pj.Filename
	}
	if pi.Line != pj.Line {
		return pi.Line < pj.Line
	}
	if pi.Column < pj.Column {
		return pi.Column < pj.Column
	}
	return ei.Line < ej.Line
}
