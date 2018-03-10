// Copyright 2016 Palantir Technologies. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package transform

import (
	"reflect"
)

// Rules is a slice where the elements are unary functions like func(*big.Float)json.Number.
type Rules []interface{}

func (rules Rules) Apply(in interface{}) interface{} {
	result := rules.apply(reflect.TypeOf(in), reflect.ValueOf(in))
	if result.IsValid() {
		return result.Interface()
	}
	return nil
}

func (rules Rules) apply(target reflect.Type, in reflect.Value) reflect.Value {
	if !in.IsValid() {
		return in
	}

	// find applicable rule
	for _, rule := range rules {
		if in.Type().AssignableTo(reflect.TypeOf(rule).In(0)) {
			return applyRule(rule, target, in)
		}
	}

	switch in.Kind() {
	case reflect.Array:
		return rules.applyArray(in)
	case reflect.Interface:
		return rules.applyInterface(target, in)
	case reflect.Map:
		return rules.applyMap(in)
	case reflect.Ptr:
		return rules.applyPointer(in)
	case reflect.Slice:
		return rules.applySlice(in)
	case reflect.Struct:
		return rules.applyStruct(in)
	default:
		return in
	}
}

func (rules Rules) applyArray(in reflect.Value) reflect.Value {
	result := reflect.New(reflect.ArrayOf(in.Len(), in.Type().Elem())).Elem()
	for i := 0; i < in.Len(); i++ {
		newV := rules.apply(in.Type().Elem(), in.Index(i))
		result.Index(i).Set(newV)
	}
	return result
}

func (rules Rules) applyInterface(target reflect.Type, in reflect.Value) reflect.Value {
	if in.IsNil() {
		return in
	}
	return rules.apply(target, in.Elem())
}

func (rules Rules) applyMap(in reflect.Value) reflect.Value {
	if in.IsNil() {
		return in
	}
	result := reflect.MakeMap(in.Type())
	for _, k := range in.MapKeys() {
		newK := rules.apply(in.Type().Key(), k)
		newV := rules.apply(in.Type().Elem(), in.MapIndex(k))
		result.SetMapIndex(newK, newV)
	}
	return result
}

func (rules Rules) applyPointer(in reflect.Value) reflect.Value {
	if in.IsNil() {
		return in
	}
	ptr := reflect.New(in.Type()).Elem()   // create a pointer of the right type
	ptr.Set(reflect.New(in.Type().Elem())) // create a value for it to point to
	ptr.Elem().Set(rules.apply(in.Type().Elem(), in.Elem()))
	return ptr
}

func (rules Rules) applySlice(in reflect.Value) reflect.Value {
	if in.IsNil() {
		return in
	}
	result := reflect.MakeSlice(in.Type(), in.Len(), in.Cap())
	for i := 0; i < in.Len(); i++ {
		newV := rules.apply(in.Type().Elem(), in.Index(i))
		result.Index(i).Set(newV)
	}
	return result
}

func (rules Rules) applyStruct(in reflect.Value) reflect.Value {
	result := reflect.New(in.Type()).Elem()
	for i := 0; i < in.NumField(); i++ {
		structField := in.Type().Field(i)
		if structField.PkgPath != "" && !structField.Anonymous {
			return in // unexported field, cannot safely copy struct
		}
		newV := rules.apply(structField.Type, in.Field(i))
		result.Field(i).Set(newV)
	}
	return result
}

func applyRule(rule interface{}, target reflect.Type, in reflect.Value) reflect.Value {
	result := reflect.ValueOf(rule).Call([]reflect.Value{in})[0]
	if !result.IsValid() && !canBeNil(target) {
		return in
	}
	if result.IsValid() && !result.Type().AssignableTo(target) {
		return in
	}
	return result
}

// canBeNil reports whether an untyped nil can be assigned to the type. See reflect.Zero.
func canBeNil(typ reflect.Type) bool {
	switch typ.Kind() {
	case reflect.Chan, reflect.Func, reflect.Interface, reflect.Map, reflect.Ptr, reflect.Slice:
		return true
	default:
		return false
	}
}
