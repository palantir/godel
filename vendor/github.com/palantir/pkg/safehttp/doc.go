// Copyright 2016 Palantir Technologies. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package safehttp provides types and functions that allow http.Client functions to be used in a manner that ensures
// that HTTP connections are not leaked. The default implementation of http.Client functions such as "Get" and "Post"
// return a reader as part of the response object. The connection used by the call can be reused only if the body of the
// response is fully drained and closed. In practice, it is easy to forget that both of these actions are necessary,
// which can lead to a large number of leaked/persistent http connections.
//
// The safehttp package provides functions that wrap the http.Client functionality for functions that return responses
// in a manner that also returns a cleanup function that drains and closes the body of the response. Callers can simply
// defer the returned cleanup function to ensure that the connections are properly relinquished. It is safe for the
// cleanup function to execute even if the body has already been drained or closed.
package safehttp
