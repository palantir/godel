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

// Based on golang.org/x/tools/imports which bears the following license:
//
// Copyright 2013 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ptimports

import (
	"go/ast"
	"go/token"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

func fixImports(fset *token.FileSet, f *ast.File, grp importGrouper, godepsPath string) {
	imports := takeImports(fset, f)
	if imports != nil && len(imports.Specs) > 0 {
		if godepsPath != "" {
			insertGodeps(imports, godepsPath)
		}
		imports.Specs = sortSpecs(fset, f, grp, imports.Specs)
		fixParens(imports)
		f.Decls = append([]ast.Decl{imports}, f.Decls...)
	}
}

func takeImports(fset *token.FileSet, f *ast.File) (imports *ast.GenDecl) {
	for len(f.Decls) > 0 {
		d, ok := f.Decls[0].(*ast.GenDecl)
		if !ok || d.Tok != token.IMPORT {
			// Not an import declaration, so we're done.
			// Import decls are always first.
			break
		}

		if imports == nil {
			imports = d
		} else {
			if imports.Doc == nil {
				imports.Doc = d.Doc
			} else if d.Doc != nil {
				imports.Doc.List = append(imports.Doc.List, d.Doc.List...)
			}
			imports.Specs = append(imports.Specs, d.Specs...)
		}

		// Put back later in a single decl
		f.Decls = f.Decls[1:]
	}
	return imports
}

func insertGodeps(d *ast.GenDecl, godepsPath string) {
	thisProject := godepsPath[:strings.Index(godepsPath, "Godeps")]
	for _, s := range d.Specs {
		path := importPath(s)
		if !strings.Contains(path, ".") {
			continue
		}
		if strings.Contains(path, "Godeps") {
			continue
		}
		if strings.HasPrefix(path, thisProject) {
			continue
		}
		path = filepath.Join(godepsPath, path)
		s.(*ast.ImportSpec).Path.Value = strconv.Quote(path)
	}
}

// All import decls require parens, even with only a single import.
func fixParens(d *ast.GenDecl) {
	if !d.Lparen.IsValid() {
		d.Lparen = d.Specs[0].Pos()
	}
}

func importPath(s ast.Spec) string {
	t, err := strconv.Unquote(s.(*ast.ImportSpec).Path.Value)
	if err == nil {
		return t
	}
	return ""
}

func importName(s ast.Spec) string {
	n := s.(*ast.ImportSpec).Name
	if n == nil {
		return ""
	}
	return n.Name
}

func importComment(s ast.Spec) string {
	c := s.(*ast.ImportSpec).Comment
	if c == nil {
		return ""
	}
	return c.Text()
}

// collapse indicates whether prev may be removed, leaving only next.
func collapse(prev, next ast.Spec) bool {
	if importPath(next) != importPath(prev) || importName(next) != importName(prev) {
		return false
	}
	return prev.(*ast.ImportSpec).Comment == nil
}

type posSpan struct {
	Start	token.Pos
	End	token.Pos
}

func sortSpecs(fset *token.FileSet, f *ast.File, grp importGrouper, specs []ast.Spec) []ast.Spec {
	// Can't short-circuit here even if specs are already sorted,
	// since they might yet need deduplication.
	// A lone import, however, may be safely ignored.
	if len(specs) <= 1 {
		return specs
	}

	// Record positions for specs.
	pos := make([]posSpan, len(specs))
	for i, s := range specs {
		pos[i] = posSpan{s.Pos(), s.End()}
	}

	// Identify comments in this range.
	// Any comment from pos[0].Start to the final line counts.
	lastLine := fset.Position(pos[len(pos)-1].End).Line
	cstart := len(f.Comments)
	cend := len(f.Comments)
	for i, g := range f.Comments {
		if g.Pos() < pos[0].Start {
			continue
		}
		if i < cstart {
			cstart = i
		}
		if fset.Position(g.End()).Line > lastLine {
			cend = i
			break
		}
	}
	comments := f.Comments[cstart:cend]

	// Assign each comment to the import spec preceding it.
	importComment := map[*ast.ImportSpec][]*ast.CommentGroup{}
	specIndex := 0
	for _, g := range comments {
		for specIndex+1 < len(specs) && pos[specIndex+1].Start <= g.Pos() {
			specIndex++
		}
		s := specs[specIndex].(*ast.ImportSpec)
		importComment[s] = append(importComment[s], g)
	}

	// Sort the import specs by import path.
	// Remove duplicates, when possible without data loss.
	// Reassign the import paths to have the same position sequence.
	// Reassign each comment to abut the end of its spec.
	// Sort the comments by new position.
	sort.Sort(byImportSpec{
		specs:	specs,
		grp:	grp,
	})

	// Dedup. Thanks to our sorting, we can just consider
	// adjacent pairs of imports.
	deduped := specs[:0]
	for i, s := range specs {
		if i == len(specs)-1 || !collapse(s, specs[i+1]) {
			deduped = append(deduped, s)
		} else {
			p := s.Pos()
			fset.File(p).MergeLine(fset.Position(p).Line)
		}
	}
	specs = deduped

	// Fix up comment positions
	for i, s := range specs {
		s := s.(*ast.ImportSpec)
		if s.Name != nil {
			s.Name.NamePos = pos[i].Start
		}
		s.Path.ValuePos = pos[i].Start
		s.EndPos = pos[i].End
		for _, g := range importComment[s] {
			for _, c := range g.List {
				c.Slash = pos[i].End
			}
		}
	}

	sort.Sort(byCommentPos(comments))

	return specs
}

type byImportSpec struct {
	specs	[]ast.Spec	// slice of *ast.ImportSpec
	grp	importGrouper
}

func (x byImportSpec) Len() int		{ return len(x.specs) }
func (x byImportSpec) Swap(i, j int)	{ x.specs[i], x.specs[j] = x.specs[j], x.specs[i] }
func (x byImportSpec) Less(i, j int) bool {
	ipath := importPath(x.specs[i])
	jpath := importPath(x.specs[j])

	igroup := x.grp.importGroup(ipath)
	jgroup := x.grp.importGroup(jpath)
	if igroup != jgroup {
		return igroup < jgroup
	}

	if ipath != jpath {
		return ipath < jpath
	}
	iname := importName(x.specs[i])
	jname := importName(x.specs[j])

	if iname != jname {
		return iname < jname
	}
	return importComment(x.specs[i]) < importComment(x.specs[j])
}

type byCommentPos []*ast.CommentGroup

func (x byCommentPos) Len() int			{ return len(x) }
func (x byCommentPos) Swap(i, j int)		{ x[i], x[j] = x[j], x[i] }
func (x byCommentPos) Less(i, j int) bool	{ return x[i].Pos() < x[j].Pos() }
