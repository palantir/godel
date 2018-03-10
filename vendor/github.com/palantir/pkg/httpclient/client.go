// Copyright 2016 Palantir Technologies. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package httpclient

import (
	"crypto/tls"
	"net"
	"net/http"
	"time"

	"golang.org/x/net/http2"
)

// NewHTTPClient provides you with an http client that is configured with the NewTransporter.
func NewHTTPClient(timeout time.Duration, tlsConf *tls.Config) *http.Client {
	return &http.Client{
		Transport: NewTransporter(timeout, tlsConf),
		Timeout:   timeout,
	}
}

// NewHTTP2Client provides you with an http2 client that is configured with the NewTransporter.
func NewHTTP2Client(timeout time.Duration, tlsConf *tls.Config) (*http.Client, error) {
	tr, err := NewHTTP2Transporter(timeout, tlsConf)
	if err != nil {
		return nil, err
	}

	return &http.Client{
		Transport: tr,
		Timeout:   timeout,
	}, nil
}

// NewTransporter configures a transporter that ensures you are never stuck in an infinite timeout
// and ensures that you don't leak connections. For example, your connection can get stuck
// in Dial forever even if you have a client timeout set, so this transport ensures that never happens.
func NewTransporter(timeout time.Duration, tlsConf *tls.Config) *http.Transport {
	return &http.Transport{
		Proxy:               http.ProxyFromEnvironment,
		TLSClientConfig:     tlsConf,
		MaxIdleConnsPerHost: 32,
		MaxIdleConns:        32,
		IdleConnTimeout:     timeout,
		TLSHandshakeTimeout: timeout,
		DialContext: (&net.Dialer{
			Timeout:   timeout,
			KeepAlive: timeout,
		}).DialContext,
	}
}

// NewHTTP2Transporter is the same as NewTransporter but also configures it for HTTP2 connections.
func NewHTTP2Transporter(timeout time.Duration, tlsConf *tls.Config) (*http.Transport, error) {
	tr := NewTransporter(timeout, tlsConf)
	if err := http2.ConfigureTransport(tr); err != nil {
		return nil, err
	}

	return tr, nil
}
