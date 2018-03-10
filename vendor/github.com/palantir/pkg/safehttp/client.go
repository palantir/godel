// Copyright 2016 Palantir Technologies. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package safehttp

import (
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
)

// Client is a type alias for http.Client that redefines the request functions (Get, Head, Post, PostForm, Do) to return
// an additional cleanup function. The returned cleanup function should be deferred after the call is made, and ensures
// that the response body is fully drained and closed so that subsequent calls that are made using the same client will
// properly reuse connections.
type Client http.Client

func (c *Client) Get(url string) (resp *http.Response, cleanup func(), err error) {
	resp, err = (*http.Client)(c).Get(url)
	return resp, responseCloser(resp), err
}

func (c *Client) Head(url string) (resp *http.Response, cleanup func(), err error) {
	resp, err = (*http.Client)(c).Head(url)
	return resp, responseCloser(resp), err
}

func (c *Client) Post(url string, contentType string, body io.Reader) (resp *http.Response, cleanup func(), err error) {
	resp, err = (*http.Client)(c).Post(url, contentType, body)
	return resp, responseCloser(resp), err
}

func (c *Client) PostForm(url string, data url.Values) (resp *http.Response, cleanup func(), err error) {
	resp, err = (*http.Client)(c).PostForm(url, data)
	return resp, responseCloser(resp), err
}

func (c *Client) Do(req *http.Request) (resp *http.Response, cleanup func(), err error) {
	resp, err = (*http.Client)(c).Do(req)
	return resp, responseCloser(resp), err
}

func responseCloser(resp *http.Response) func() {
	if resp == nil {
		// if response is nil, return a no-op cleanup function
		return func() {}
	}
	return func() {
		// drain and close treated as best-effort and errors are not reported
		_, _ = io.Copy(ioutil.Discard, resp.Body)
		_ = resp.Body.Close()
	}
}
