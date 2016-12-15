// Copyright 2013 Kamil Kisiel
// Modifications copyright 2016 Palantir Technologies, Inc.
// Licensed under the MIT License. See LICENSE in the project root
// for license information.

package outparamcheck

import (
	"encoding/json"
	"fmt"
	"go/ast"
	"go/build"
	"go/token"
	"go/types"
	"io/ioutil"
	"sort"
	"strings"
	"sync"

	"github.com/kisielk/gotool"
	"github.com/pkg/errors"
	"golang.org/x/tools/go/loader"

	"github.com/palantir/godel/apps/okgo/generated_src/internal/github.com/palantir/checks/outparamcheck/exprs"
)

func Run(cfgParam string, paths []string) error {
	cfg := Config{}
	if cfgParam != "" {
		var usrCfg Config
		var err error
		if strings.HasPrefix(cfgParam, "@") {
			usrCfg, err = loadCfgFromPath(cfgParam[1:])
		} else {
			usrCfg, err = loadCfg(cfgParam)
		}
		if err != nil {
			return errors.Wrapf(err, "Failed to load configuration from parameter %s", cfgParam)
		}
		for key, val := range usrCfg {
			cfg[key] = val
		}
	}
	// add default config (values for default will override any user-supplied config for the same keys)
	for key, val := range defaultCfg {
		cfg[key] = val
	}

	prog, err := load(paths)
	if err != nil {
		return errors.WithStack(err)
	}
	errs := run(prog, cfg)
	if len(errs) > 0 {
		reportErrors(errs)
		return fmt.Errorf("%s; the parameters listed above require the use of '&', for example f(&x) instead of f(x)",
			plural(len(errs), "error", "errors"))
	}
	return nil
}

func run(prog *loader.Program, cfg Config) []OutParamError {
	var errs []OutParamError
	var mut sync.Mutex	// guards errs
	var wg sync.WaitGroup
	for _, pkgInfo := range prog.InitialPackages() {
		if pkgInfo.Pkg.Path() == "unsafe" {	// not a real package
			continue
		}

		wg.Add(1)

		go func(pkgInfo *loader.PackageInfo) {
			defer wg.Done()
			v := &visitor{
				prog:	prog,
				pkg:	pkgInfo,
				lines:	map[string][]string{},
				errors:	[]OutParamError{},
				cfg:	cfg,
			}
			for _, astFile := range pkgInfo.Files {
				exprs.Walk(v, astFile)
			}
			mut.Lock()
			defer mut.Unlock()
			errs = append(errs, v.errors...)
		}(pkgInfo)
	}
	wg.Wait()
	return errs
}

func loadCfgFromPath(cfgPath string) (Config, error) {
	cfgBytes, err := ioutil.ReadFile(cfgPath)
	if err != nil {
		return Config{}, errors.Wrapf(err, "failed to read file %s", cfgPath)
	}
	return loadCfg(string(cfgBytes))
}

func loadCfg(cfgJSON string) (Config, error) {
	var cfg Config
	if err := json.Unmarshal([]byte(cfgJSON), &cfg); err != nil {
		return Config{}, errors.Wrapf(err, "failed to unmarshal json %s", cfgJSON)
	}
	return cfg, nil
}

func load(paths []string) (*loader.Program, error) {
	loadcfg := loader.Config{
		Build: &build.Default,
	}
	includeTests := true
	rest, err := loadcfg.FromArgs(gotool.ImportPaths(paths), includeTests)
	if err != nil {
		return nil, errors.Wrapf(err, "could not parse arguments")
	}
	if len(rest) > 0 {
		return nil, errors.Errorf("unhandled extra arguments: %v", rest)
	}
	prog, err := loadcfg.Load()
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return prog, nil
}

type visitor struct {
	prog	*loader.Program
	pkg	*loader.PackageInfo
	lines	map[string][]string
	errors	[]OutParamError
	cfg	Config
}

func (v *visitor) Visit(expr ast.Expr) {
	call, ok := expr.(*ast.CallExpr)
	if !ok {
		return
	}
	key, method, ok := v.keyAndName(call)
	if !ok {
		return
	}
	for name, outs := range v.cfg {
		// Suffix-matching so they also apply to vendored packages
		if strings.HasSuffix(key, name) {
			for _, i := range outs {
				arg := call.Args[i]
				if !isAddr(arg) {
					v.errorAt(arg.Pos(), method, i)
				}
			}
		}
	}
}

func (v *visitor) keyAndName(call *ast.CallExpr) (key string, name string, ok bool) {
	switch target := call.Fun.(type) {
	case *ast.Ident:
		// Function calls without a selector; this includes calls within the
		// same package as well as calls into dot-imported packages
		if def, ok := v.pkg.Uses[target]; ok && def.Pkg() != nil {
			return fmt.Sprintf("%v.%v", def.Pkg().Path(), target.Name), target.Name, true
		}
	case *ast.SelectorExpr:
		// Function calls into other packages
		if recv, ok := target.X.(*ast.Ident); ok {
			if pkg, ok := v.pkg.Uses[recv].(*types.PkgName); ok {
				return fmt.Sprintf("%v.%v", pkg.Imported().Path(), target.Sel.Name), target.Sel.Name, true
			}
		}
		// Method calls
		if typ, ok := v.pkg.Types[target.X]; ok {
			return fmt.Sprintf("%v.%v", typ.Type.String(), target.Sel.Name), target.Sel.Name, true
		}
	}
	return "", "", false
}

func (v *visitor) errorAt(pos token.Pos, method string, argument int) {
	position := v.prog.Fset.Position(pos)
	lines, ok := v.lines[position.Filename]
	if !ok {
		contents, err := ioutil.ReadFile(position.Filename)
		if err != nil {
			contents = nil
		}
		lines = strings.Split(string(contents), "\n")
		v.lines[position.Filename] = lines
	}

	var line string
	if position.Line-1 < len(lines) {
		line = strings.TrimSpace(lines[position.Line-1])
	}
	v.errors = append(v.errors, OutParamError{position, line, method, argument})
}

func isAddr(expr ast.Expr) bool {
	switch expr := expr.(type) {
	case *ast.UnaryExpr:
		// The expected usage for output parameters, which is &x
		return expr.Op == token.AND
	case *ast.StarExpr:
		// Allow *&x as an explicit way to signal that no & is intended
		child, ok := expr.X.(*ast.UnaryExpr)
		return ok && child.Op == token.AND
	case *ast.Ident:
		// Allow passing literal nil
		return expr.Name == "nil"
	default:
		return false
	}
}

func reportErrors(errs []OutParamError) {
	sort.Sort(byLocation(errs))
	for _, err := range errs {
		fmt.Println(err)
	}
}

func plural(count int, singular, plural string) string {
	if count == 1 {
		return fmt.Sprintf("%d %s", count, singular)
	}
	return fmt.Sprintf("%d %s", count, plural)
}
