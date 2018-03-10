// Copyright 2016 Palantir Technologies. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package safejson provides functions that allows JSON to be marshaled
// and unmarshaled in a safe and consistent manner. This package exists
// to standardize and normalize some of the default behavior of the json
// package that is unintuitive.
//
// Marshal:
//
// The default encoder returned by json.NewEncoder has "SetEscapeHTML" set
// to "true", which makes sense for HTML environments, but results in output
// that is hard to read in non-HTML environments. The default behavior of
// json.Marshal also appends a newline to the end of the generated JSON which,
// although this is technically legal from a JSON perspective, is often unexpected.
//
// Unmarshal:
//
// The default decoder returned by json.NewDecoder does not have the "UseNumber"
// behavior enabled. This means that all numeric values are unmarshaled as a float64.
// This behavior is generally less flexible, so safejson sets "UseNumber" to "true",
// which ensures that all numbers are unmarshaled as a json.Number.
package safejson
