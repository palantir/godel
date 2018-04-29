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

package config

import (
	"fmt"
	"reflect"
)

// getConfigStringValue returns a value based on the inputs. Is similar to "getConfigValue", but if the result of
// "getConfigValue" is an empty string, the provided "defaultVal" is returned (rather than returning the the empty
// string). This has the effect of ensuring that the provided "defaultVal" is used as the string value if the value
// obtained from "primaryPtr" or "secondaryPtr" is non-nil but is the empty string.
func getConfigStringValue(primaryPtr, secondaryPtr *string, defaultVal string) string {
	strVal := getConfigValue(primaryPtr, secondaryPtr, "").(string)
	if strVal == "" {
		strVal = defaultVal
	}
	return strVal
}

// getConfigValue returns a value based on the inputs. "primaryValPtr" and "secondaryValPtr" must both be pointers that
// point to values of the same type, and "defaultVal" must be the type pointed to by "primaryValPtr" and
// "secondaryValPtr". If the pointer value of "primaryValPtr" is non-nil, its value is returned. Otherwise, if the
// pointer value of "secondaryValPtr" is non-nil, its value is returned. If the pointer values of both "primaryValPtr"
// and "secondaryValPtr" are nil, "defaultVal" is returned. As a special case, if "defaultVal" is a raw/untyped nil,
// then the defaultVal will be the zero value of the type pointed to by the valPtr inputs.
func getConfigValue(primaryValPtr, secondaryValPtr, defaultVal interface{}) interface{} {
	t1 := reflect.TypeOf(primaryValPtr)
	t2 := reflect.TypeOf(secondaryValPtr)
	if t1 != t2 {
		panic(fmt.Sprintf("types of primaryValPtr and secondaryValPtr are not equal: %v != %v", t1, t2))
	}
	if t1.Kind() != reflect.Ptr {
		panic(fmt.Sprintf("type %v is not a pointer type", t1))
	}

	// if "defaultVal" is an untyped/raw "nil" value, set it to be the zero value of the type pointed to by the pointers.
	if defaultVal == nil {
		defaultVal = reflect.Zero(t1.Elem()).Interface()
	}
	if defaultValType := reflect.TypeOf(defaultVal); t1.Elem() != defaultValType {
		panic(fmt.Sprintf("the type referenced by the pointers differs from type of default value: %v != %v", t1.Elem(), defaultValType))
	}

	if !reflect.ValueOf(primaryValPtr).IsNil() {
		return reflect.ValueOf(primaryValPtr).Elem().Interface()
	}
	if !reflect.ValueOf(secondaryValPtr).IsNil() {
		return reflect.ValueOf(secondaryValPtr).Elem().Interface()
	}
	return defaultVal
}
