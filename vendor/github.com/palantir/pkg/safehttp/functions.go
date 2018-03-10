// Copyright 2016 Palantir Technologies. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package safehttp

import (
	"io"
	"net/http"
	"net/url"
)

func Get(c *http.Client, url string) (resp *http.Response, cleanup func(), err error) {
	return (*Client)(c).Get(url)
}

func Head(c *http.Client, url string) (resp *http.Response, cleanup func(), err error) {
	return (*Client)(c).Head(url)
}

func Post(c *http.Client, url string, contentType string, body io.Reader) (resp *http.Response, cleanup func(), err error) {
	return (*Client)(c).Post(url, contentType, body)
}

func PostForm(c *http.Client, url string, data url.Values) (resp *http.Response, cleanup func(), err error) {
	return (*Client)(c).PostForm(url, data)
}

func Do(c *http.Client, req *http.Request) (resp *http.Response, cleanup func(), err error) {
	return (*Client)(c).Do(req)
}
