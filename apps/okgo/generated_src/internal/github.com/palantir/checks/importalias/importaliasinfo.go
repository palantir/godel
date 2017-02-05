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

package amalgomated

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/ioutil"
	"sort"
	"strings"
)

type ImportAliasInfo struct {
	ImportPath	string
	Alias		string
	// file -> line information for import in the file
	Occurrences	map[string]token.Position
}

type ImportAlias struct {
	ImportPath	string
	Alias		string
	Pos		token.Position
}

type projectImportAliasInfo struct {
	importInfos map[string]map[string]ImportAliasInfo
}

type ProjectImportInfo interface {
	// AddImportAliasesFromFile adds all of the import alias information from the given file.
	AddImportAliasesFromFile(filename string) error

	// ImportsWithMultipleAliases returns a map from an imported package path to all of the aliases to import the package.
	// The aliases are sorted by the number of uses of that alias.
	ImportsToAliases() map[string][]ImportAliasInfo

	// FilesToImportAliases returns a map from each file in the project to all of the alias imports in the file.
	FilesToImportAliases() map[string][]ImportAlias

	// GetAliasStatus returns the AliasStatus for the given alias used to import the package with the provided path.
	GetAliasStatus(alias, importPath string) AliasStatus
}

type AliasStatus struct {
	// true if this alias is the only alias used for a package or is the most common alias used for a package.
	OK	bool
	// recommendation for how to fix the issue if OK is false.
	Recommendation	string
}

func NewProjectImportInfo() ProjectImportInfo {
	return &projectImportAliasInfo{
		importInfos: make(map[string]map[string]ImportAliasInfo),
	}
}

func (p *projectImportAliasInfo) AddImportAliasesFromFile(filename string) error {
	src, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, filename, src, parser.ParseComments)
	if err != nil {
		return fmt.Errorf("failed to parse file %s: %v", filename, err)
	}

	var visitor visitFn
	visitor = visitFn(func(node ast.Node) ast.Visitor {
		if node == nil {
			return visitor
		}
		switch v := node.(type) {
		case *ast.ImportSpec:
			if v.Name != nil && v.Name.Name != "." && v.Name.Name != "_" {

				p.addImportAlias(filename, v.Name.Name, v.Path.Value, fset.Position(v.Pos()))
				break
			}
		}
		return visitor
	})
	ast.Walk(visitor, file)
	return nil
}

type visitFn func(node ast.Node) ast.Visitor

func (fn visitFn) Visit(node ast.Node) ast.Visitor {
	return fn(node)
}

func (p *projectImportAliasInfo) addImportAlias(file, alias, importPath string, pos token.Position) {
	if _, ok := p.importInfos[importPath]; !ok {
		p.importInfos[importPath] = make(map[string]ImportAliasInfo)
	}
	if _, ok := p.importInfos[importPath][alias]; !ok {
		p.importInfos[importPath][alias] = ImportAliasInfo{
			ImportPath:	importPath,
			Alias:		alias,
			Occurrences:	make(map[string]token.Position),
		}
	}
	p.importInfos[importPath][alias].Occurrences[file] = pos
}

func (p *projectImportAliasInfo) ImportsToAliases() map[string][]ImportAliasInfo {
	m := make(map[string][]ImportAliasInfo)
	for importPath, aliases := range p.importInfos {
		for _, aliasInfo := range aliases {
			m[importPath] = append(m[importPath], aliasInfo)
		}
	}
	for _, v := range m {
		sort.Sort(byNumOccurrencesDesc(v))
	}
	return m
}

func (p *projectImportAliasInfo) FilesToImportAliases() map[string][]ImportAlias {
	m := make(map[string][]ImportAlias)
	for importPath, aliases := range p.importInfos {
		for _, currAliasInfo := range aliases {
			for file, pos := range currAliasInfo.Occurrences {
				m[file] = append(m[file], ImportAlias{
					ImportPath:	importPath,
					Alias:		currAliasInfo.Alias,
					Pos:		pos,
				})
			}
		}
	}
	for _, v := range m {
		sort.Sort(byPos(v))
	}
	return m
}

func (p *projectImportAliasInfo) GetAliasStatus(alias, importPath string) AliasStatus {
	importsToAliases := p.ImportsToAliases()
	if aliases, ok := importsToAliases[importPath]; ok && len(aliases) > 1 {
		var mostCommonAliases []string
		for _, currAlias := range aliases {
			if len(currAlias.Occurrences) != len(aliases[0].Occurrences) {
				break
			}
			mostCommonAliases = append(mostCommonAliases, currAlias.Alias)
		}
		switch {
		case len(mostCommonAliases) > 1:
			var aliasesUsed string
			if len(mostCommonAliases) == 2 {
				aliasesUsed = fmt.Sprintf("%q and %q are both", mostCommonAliases[0], mostCommonAliases[1])
			} else {
				var quoted []string
				for _, curr := range mostCommonAliases {
					quoted = append(quoted, fmt.Sprintf("%q", curr))
				}
				aliasesUsed = strings.Join(quoted[:len(quoted)-1], ", ")
				aliasesUsed += " and " + quoted[len(quoted)-1] + " are all"
			}

			var timesUsed string
			if len(aliases[0].Occurrences) == 1 {
				timesUsed = "once"
			} else {
				timesUsed = fmt.Sprintf("%d times", len(aliases[0].Occurrences))
			}

			// there is not a single most common alias
			return AliasStatus{
				OK:		false,
				Recommendation:	fmt.Sprintf("No consensus alias exists for this import in the project (%s used %s each)", aliasesUsed, timesUsed),
			}
		case alias != mostCommonAliases[0]:
			// this is not the most common alias
			return AliasStatus{
				OK:		false,
				Recommendation:	fmt.Sprintf("Use alias %q instead", mostCommonAliases[0]),
			}
		}
	}
	return AliasStatus{
		OK: true,
	}
}

type byNumOccurrencesDesc []ImportAliasInfo

func (a byNumOccurrencesDesc) Len() int		{ return len(a) }
func (a byNumOccurrencesDesc) Swap(i, j int)	{ a[i], a[j] = a[j], a[i] }
func (a byNumOccurrencesDesc) Less(i, j int) bool {
	if len(a[i].Occurrences) == len(a[j].Occurrences) {
		// if number of occurrences are the same, do secondary sort based on name of alias
		return strings.Compare(a[i].Alias, a[j].Alias) < 0
	}
	// sort occurrences by descending order
	return len(a[i].Occurrences) > len(a[j].Occurrences)
}

type byPos []ImportAlias

func (a byPos) Len() int	{ return len(a) }
func (a byPos) Swap(i, j int)	{ a[i], a[j] = a[j], a[i] }
func (a byPos) Less(i, j int) bool {
	return a[i].Pos.Line < a[j].Pos.Line
}
