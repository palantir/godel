// Copyright 2016 Palantir Technologies. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package httpserver_test

import (
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/palantir/pkg/httpserver"
)

func TestURLReady(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	}))
	defer server.Close()

	r := httpserver.URLReady(server.URL, httpserver.WaitTimeoutParam(200*time.Millisecond))
	assert.True(t, <-r)
}

func TestURLReadyFunctionTimeout(t *testing.T) {
	counter := int32(0)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&counter, 1)
		time.Sleep(time.Second)
	}))
	defer server.Close()

	r := httpserver.URLReady(server.URL, httpserver.WaitTimeoutParam(400*time.Millisecond))
	assert.False(t, <-r)
	assert.Equal(t, int32(1), atomic.LoadInt32(&counter))
}

func TestURLReadyTimeout(t *testing.T) {
	r := httpserver.URLReady("http://localhost:9999", httpserver.WaitTimeoutParam(200*time.Millisecond))
	assert.False(t, <-r)
}

func TestReadyPut(t *testing.T) {
	// server responds with status code http.StatusAccepted for any request with an http.MethodPost method,
	// times out for 1 second otherwise.
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			w.WriteHeader(http.StatusAccepted)
			return
		}
		time.Sleep(time.Second)
	}))
	defer server.Close()

	r := httpserver.Ready(
		func() (*http.Response, error) {
			return http.Post(server.URL, "text", nil)
		},
		httpserver.ReadyRespParam(func(resp *http.Response) bool {
			return resp.StatusCode == http.StatusAccepted
		}),
		httpserver.WaitTimeoutParam(200*time.Millisecond))
	assert.True(t, <-r)
}
