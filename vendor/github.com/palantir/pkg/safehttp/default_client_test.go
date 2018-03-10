// Copyright 2016 Palantir Technologies. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package safehttp_test

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/http/httptrace"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDefaultClientLeaksNoClose(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Hello, client")
	}))
	defer ts.Close()

	var reused []bool
	reusedTracker := func(info httptrace.GotConnInfo) {
		reused = append(reused, info.Reused)
	}

	for i := 0; i <= 2; i++ {
		// create "GET" request
		req, _ := http.NewRequest("GET", ts.URL, nil)
		trace := &httptrace.ClientTrace{
			GotConn: reusedTracker,
		}
		req = req.WithContext(httptrace.WithClientTrace(req.Context(), trace))

		// execute request, but do not drain the body
		_, err := http.DefaultClient.Do(req)
		require.NoError(t, err)
	}

	// none of the connections are reused
	assert.Equal(t, []bool{false, false, false}, reused)
}

func TestDefaultClientDoesNotLeakWhenClosed(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Hello, client")
	}))
	defer ts.Close()

	var reused []bool
	reusedTracker := func(info httptrace.GotConnInfo) {
		reused = append(reused, info.Reused)
	}

	for i := 0; i <= 2; i++ {
		// create "GET" request
		req, _ := http.NewRequest("GET", ts.URL, nil)
		trace := &httptrace.ClientTrace{
			GotConn: reusedTracker,
		}
		req = req.WithContext(httptrace.WithClientTrace(req.Context(), trace))

		// execute request
		resp, err := http.DefaultClient.Do(req)
		require.NoError(t, err)

		// drain and close response body
		_, err = ioutil.ReadAll(resp.Body)
		require.NoError(t, err)
		err = resp.Body.Close()
		require.NoError(t, err)
	}

	// first connection is not reused, but subsequent connections are
	assert.Equal(t, []bool{false, true, true}, reused)
}
