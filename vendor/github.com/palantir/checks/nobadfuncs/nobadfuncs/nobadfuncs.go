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

package nobadfuncs

import (
	"fmt"
	"go/ast"
	"go/build"
	"go/parser"
	"go/token"
	"go/types"
	"io"
	"regexp"
	"sort"

	"github.com/pkg/errors"
	"golang.org/x/tools/go/loader"
)

// FuncRef is a reference to a specific function. Matches the string representation of *types.Func, which is of the
// form "func (*net/http.Client).Do(req *net/http.Request) (*net/http.Response, error)".
type FuncRef string

func PrintAllFuncRefs(pkgs []string, stdout io.Writer) error {
	_, err := printFuncRefUsages(pkgs, nil, stdout)
	return err
}

func PrintBadFuncRefs(pkgs []string, sigs map[string]string, stdout io.Writer) (bool, error) {
	if len(sigs) == 0 {
		// if there are no signatures, there will be no output
		return true, nil
	}
	return printFuncRefUsages(pkgs, sigs, stdout)
}

func printFuncRefUsages(pkgs []string, sigs map[string]string, stdout io.Writer) (bool, error) {
	loadcfg := loader.Config{
		Build:      &build.Default,
		ParserMode: parser.ParseComments,
	}
	// add all packages to load
	for _, currPkg := range pkgs {
		loadcfg.ImportWithTests(currPkg)
	}

	// load program
	prog, err := loadcfg.Load()
	if err != nil {
		return false, errors.Wrapf(err, "failed to load program")
	}
	sort.Strings(pkgs)

	noBadRefs := true
	for _, currPkg := range pkgs {
		info := prog.Package(currPkg)
		if info == nil {
			panic(fmt.Sprintf("failed to find %s in %v; imported %v", currPkg, prog.AllPackages, prog.Imported))
		}

		for _, currFile := range info.Files {
			currOutput, err := findFuncRefUsage(info.Info, currFile, prog.Fset, sigs)
			if err != nil {
				return false, err
			}

			if len(sigs) == 0 {
				// "all" mode -- print all references
				visitInOrder(currOutput, func(pos token.Position, ref FuncRef) {
					fmt.Fprintf(stdout, "%s: %s\n", pos.String(), ref)
				})
				continue
			}

			// filter out any matches that have a whitelist comment
			filterFuncRefs(currOutput, okCommentRegxp.MatchString)

			visitInOrder(currOutput, func(pos token.Position, ref FuncRef) {
				reason, ok := sigs[string(ref)]
				if !ok {
					return
				}
				noBadRefs = false
				if reason == "" {
					reason = fmt.Sprintf("references to %q are not allowed. Remove this reference or whitelist it by adding a comment of the form '// OK: [reason]' to the line before it.", ref)
				}
				fmt.Fprintf(stdout, "%s: %s\n", pos.String(), reason)
			})
		}
	}
	return noBadRefs, nil
}

// matches a single-line comment beginning with "// OK: " followed by at least one non-whitespace character.
var okCommentRegxp = regexp.MustCompile(regexp.QuoteMeta(`// OK: `) + `\S.*`)

func filterFuncRefs(funcRefs map[FuncRef]map[token.Position]string, filter func(string) bool) {
	for _, refPosToRefComment := range funcRefs {
		for pos, comment := range refPosToRefComment {
			if filter(comment) {
				delete(refPosToRefComment, pos)
			}
		}
	}
}

func visitInOrder(funcRefs map[FuncRef]map[token.Position]string, visitor func(token.Position, FuncRef)) {
	var allPos []token.Position
	posToFuncRef := make(map[token.Position]FuncRef)

	for funcRef, posToComment := range funcRefs {
		for pos := range posToComment {
			allPos = append(allPos, pos)
			posToFuncRef[pos] = funcRef
		}
	}
	sort.Sort(posSlice(allPos))

	for _, currPos := range allPos {
		visitor(currPos, posToFuncRef[currPos])
	}
}

type posSlice []token.Position

func (a posSlice) Len() int      { return len(a) }
func (a posSlice) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a posSlice) Less(i, j int) bool {
	if a[i].Line != a[j].Line {
		return a[i].Line < a[j].Line
	}
	return a[i].Column < a[j].Column
}

// findFuncRefUsage returns all of the function references in the specified package. If "sigs" is non-empty, then only
// function signature that match a key in the "sigs" map are included; otherwise, all function references are returned.
func findFuncRefUsage(info types.Info, f *ast.File, fset *token.FileSet, sigs map[string]string) (map[FuncRef]map[token.Position]string, error) {
	rv := make(map[FuncRef]map[token.Position]string)

	// map from line to comments in file
	lineToComment := make(map[int]string)
	for _, commentGroup := range f.Comments {
		for _, comment := range commentGroup.List {
			lineToComment[fset.Position(comment.Pos()).Line] = comment.Text
		}
	}

	var keys []*ast.Ident
	for k := range info.Uses {
		keys = append(keys, k)
	}
	sort.Sort(identSlice(keys))

	for _, id := range keys {
		obj := info.Uses[id]
		funcPtr, ok := obj.(*types.Func)
		if !ok {
			continue
		}

		// transform function to a form where names are removed from receivers, params and return values
		// and package references have path to the vendor directory removed.
		funcPtr = toFuncWithNoIdentifiersRemoveVendor(funcPtr)
		currSig := FuncRef(funcPtr.String())

		if len(sigs) > 0 {
			if _, ok := sigs[string(currSig)]; !ok {
				// if sigs is non-empty, skip any entries that don't match the signature
				continue
			}
		}

		lineMap := rv[currSig]
		if lineMap == nil {
			rv[currSig] = make(map[token.Position]string)
			lineMap = rv[currSig]
		}

		currSigPos := fset.Position(id.Pos())
		lineMap[currSigPos] = lineToComment[currSigPos.Line-1]
	}
	return rv, nil
}

type identSlice []*ast.Ident

func (a identSlice) Len() int           { return len(a) }
func (a identSlice) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a identSlice) Less(i, j int) bool { return a[i].Pos() < a[j].Pos() }
