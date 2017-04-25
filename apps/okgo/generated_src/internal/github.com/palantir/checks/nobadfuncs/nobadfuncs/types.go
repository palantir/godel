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
	"go/types"
	"strings"
)

// returns a new version of the provided *types.Func where all of the identifier names have been removed and all
// package references have their vendor references removed.
//
// Without this, the "String()" function for a function returns output of the form:
//  func (github.com/palantir/checks/vendor/github.com/Foo).Foo(paramVarName github.com/palantir/checks/vendor/github.com/foo.FooType) (namedReturnVar github.com/palantir/checks/vendor/github.com/foo.FooType)
//
// The "String()" function for the function returned by this function for the above would be:
//  func (github.com/Foo).Foo(github.com/foo.FooType) github.com/foo.FooType
func toFuncWithNoIdentifiersRemoveVendor(in *types.Func) *types.Func {
	sig, ok := in.Type().(*types.Signature)
	if !ok {
		return in
	}
	newSig := types.NewSignature(newVarNoName(sig.Recv()), newTupleNoNames(sig.Params()), newTupleNoNames(sig.Results()), sig.Variadic())
	newSig = toTypeRemoveVendor(newSig).(*types.Signature)
	return types.NewFunc(in.Pos(), pkgNoVendor(in.Pkg()), in.Name(), newSig)
}

func newTupleNoNames(in *types.Tuple) *types.Tuple {
	if in == nil || in.Len() == 0 {
		return in
	}
	var newVars []*types.Var
	for i := 0; i < in.Len(); i++ {
		newVars = append(newVars, newVarNoName(in.At(i)))
	}
	return types.NewTuple(newVars...)
}

func newVarNoName(in *types.Var) *types.Var {
	if in == nil {
		return in
	}
	return types.NewVar(in.Pos(), pkgNoVendor(in.Pkg()), "", in.Type())
}

func pkgNoVendor(in *types.Package) *types.Package {
	if in == nil {
		return nil
	}
	return types.NewPackage(removeVendor(in.Path()), in.Name())
}

func removeVendor(in string) string {
	out := in
	if vendorIdx := strings.LastIndex(out, "vendor/"); vendorIdx != -1 {
		out = out[vendorIdx+len("vendor/"):]
	}
	return out
}

func toTypeRemoveVendor(in types.Type) types.Type {
	switch typ := in.(type) {
	default:
		panic(fmt.Errorf("unrecognized type: %v", in))
	case *types.Basic:
		return in
	case *types.Array:
		return types.NewArray(toTypeRemoveVendor(typ.Elem()), typ.Len())
	case *types.Slice:
		return types.NewSlice(toTypeRemoveVendor(typ.Elem()))
	case *types.Struct:
		return in
	case *types.Pointer:
		return types.NewPointer(toTypeRemoveVendor(typ.Elem()))
	case *types.Tuple:
		return newTupleRemoveVendor(typ)
	case *types.Signature:
		return types.NewSignature(newVarRemoveVendor(typ.Recv()), newTupleRemoveVendor(typ.Params()), newTupleRemoveVendor(typ.Results()), typ.Variadic())
	case *types.Interface:
		return in
	case *types.Map:
		return types.NewMap(toTypeRemoveVendor(typ.Key()), toTypeRemoveVendor(typ.Elem()))
	case *types.Chan:
		return types.NewChan(typ.Dir(), toTypeRemoveVendor(typ.Elem()))
	case *types.Named:
		var methods []*types.Func
		for i := 0; i < typ.NumMethods(); i++ {
			methods = append(methods, typ.Method(i))
		}
		// this is the crux of the function: for all type names, transform the "package" parameter such that the
		// path to the vendor directory is removed.
		typName := types.NewTypeName(typ.Obj().Pos(), pkgNoVendor(typ.Obj().Pkg()), typ.Obj().Name(), typ.Obj().Type())
		return types.NewNamed(typName, typ.Underlying(), methods)
	}
}

func newTupleRemoveVendor(in *types.Tuple) *types.Tuple {
	var newVars []*types.Var
	for i := 0; i < in.Len(); i++ {
		newVars = append(newVars, newVarRemoveVendor(in.At(i)))
	}
	return types.NewTuple(newVars...)
}

func newVarRemoveVendor(in *types.Var) *types.Var {
	if in == nil {
		return in
	}
	return types.NewVar(in.Pos(), in.Pkg(), in.Name(), toTypeRemoveVendor(in.Type()))
}
