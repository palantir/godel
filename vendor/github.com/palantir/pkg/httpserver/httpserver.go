// Copyright 2016 Palantir Technologies. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package httpserver

import (
	"net"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
)

// AvailablePort returns a best-effort determination of an available port. Does so by opening a TCP listener on
// localhost, determining the port used by that listener, closing the listener and returning the address that was used
// by the listener. This is best-effort because there is no way to guarantee that another process will not take the port
// between the time when the listener is closed and the returned port is used by the caller.
func AvailablePort() (port int, rErr error) {
	l, err := net.Listen("tcp", "127.0.0.1:0")
	defer func() {
		if err := l.Close(); err != nil && rErr == nil {
			rErr = err
		}
	}()
	if err != nil {
		return 0, err
	}

	addrString := l.Addr().String()
	port, err = strconv.Atoi(addrString[strings.LastIndex(addrString, ":")+1:])
	if err != nil {
		return 0, err
	}

	return port, nil
}

// URLReady returns a channel that is sent "true" when an http.Get executed against the provided URL returns a response
// with status code http.StatusOK. This is a convenience function that calls Ready with a readyCall that consists of
// sending a GET request using the default HTTP client to the provided URL.
func URLReady(url string, params ...ReadyParam) <-chan bool {
	return Ready(func() (*http.Response, error) {
		return http.Get(url)
	}, params...)
}

// Ready returns a channel that is sent "true" when the provided readyCall returns a nil error and a response that
// returns "true" when provided to readyResp. The readyCall is invoked once every tick duration until it either returns
// a nil error and readyResp returns true for the response or the timeout duration is reached, in which case "false" is
// sent on the channel.
//
// readyCall should by a function that returns quickly. At most one readyCall will be running at a particular time.
//
// ReadyRespParam is used to specify the function that should be used to check if the response returned by the readyCall
// should be interpreted as "ready". If it is not specified, a default function that returns true if the response code
// is 200 is used.
//
// ReadyRetryIntervalParam is used to specify the retry interval for the "readyCall". If it is not specified, a default
// value of 100ms is used.
//
// WaitTimeoutParam is used to specify the timeout duration (the time after which the channel should return "false"). If
// it is not specified, a default value of 5s is used.
func Ready(readyCall func() (*http.Response, error), params ...ReadyParam) <-chan bool {
	cfg := &readyConfig{
		readyResp: func(resp *http.Response) bool {
			return resp.StatusCode == http.StatusOK
		},
		timeout:      5 * time.Second,
		tickDuration: 100 * time.Millisecond,
	}
	for _, p := range params {
		p.config(cfg)
	}

	once := &sync.Once{}
	done := make(chan struct{})

	ready := make(chan bool)
	go func() {
		timeout := time.NewTimer(cfg.timeout)
		defer timeout.Stop()

		// start a separate goroutine with the ticker. Done so that the possibly expensive action will not
		// block the timeout.
		go func() {
			ticker := time.NewTicker(cfg.tickDuration)
			defer ticker.Stop()
			for {
				select {
				case <-done:
					return
				case <-ticker.C:
					if resp, err := readyCall(); err == nil && cfg.readyResp(resp) {
						once.Do(func() {
							ready <- true
							close(done)
						})
					}
				}
			}
		}()

		// timeout channel
		for {
			select {
			case <-done:
				return
			case <-timeout.C:
				once.Do(func() {
					ready <- false
					close(done)
				})
				return
			}
		}
	}()
	return ready
}

type ReadyParam interface {
	config(*readyConfig)
}

type readyParam func(*readyConfig)

func (p readyParam) config(cfg *readyConfig) {
	p(cfg)
}

func ReadyRespParam(readyResp func(*http.Response) bool) ReadyParam {
	return readyParam(func(cfg *readyConfig) {
		cfg.readyResp = readyResp
	})
}

func WaitTimeoutParam(t time.Duration) ReadyParam {
	return readyParam(func(cfg *readyConfig) {
		cfg.timeout = t
	})
}

func ReadyRetryIntervalParam(t time.Duration) ReadyParam {
	return readyParam(func(cfg *readyConfig) {
		cfg.tickDuration = t
	})
}

type readyConfig struct {
	readyResp    func(*http.Response) bool
	timeout      time.Duration
	tickDuration time.Duration
}
