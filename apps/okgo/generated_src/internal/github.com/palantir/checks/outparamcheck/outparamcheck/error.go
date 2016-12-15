// Copyright 2013 Kamil Kisiel
// Modifications copyright 2016 Palantir Technologies, Inc.
// Licensed under the MIT License. See LICENSE in the project root
// for license information.

package outparamcheck

import (
	"fmt"
	"go/token"
	"strings"

	"github.com/dustin/go-humanize"
)

type OutParamError struct {
	Pos		token.Position
	Line		string
	Method		string
	Argument	int
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
