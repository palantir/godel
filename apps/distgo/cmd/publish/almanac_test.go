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

package publish_test

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/palantir/godel/apps/distgo/cmd/publish"
)

type errorRoundTripper struct{}

func (s *errorRoundTripper) RoundTrip(r *http.Request) (*http.Response, error) {
	return nil, errors.Errorf("Unable to connect")
}

func TestAlmanacConnectionInfo(t *testing.T) {
	var handlerFunc func(w http.ResponseWriter, r *http.Request)
	handlerFuncPtr := &handlerFunc
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		localHandlerFunc := *handlerFuncPtr
		localHandlerFunc(w, r)
	}))
	defer ts.Close()
	a := publish.AlmanacInfo{
		URL: ts.URL,
	}

	for i, currCase := range []struct {
		action   func(a publish.AlmanacInfo) error
		handler  func(w http.ResponseWriter, r *http.Request)
		verifier func(caseNum int, err error)
	}{
		{
			action: func(a publish.AlmanacInfo) error {
				return a.CheckConnectivity(http.DefaultClient)
			},
			handler: func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, "/v1/units", r.URL.String())
				assert.Equal(t, http.MethodGet, r.Method)
				_, ok := r.Header[http.CanonicalHeaderKey("X-timestamp")]
				assert.True(t, ok)
				_, ok = r.Header[http.CanonicalHeaderKey("X-authorization")]
				assert.True(t, ok)
				_, err := w.Write([]byte("hello"))
				require.NoError(t, err)
			},
			verifier: func(caseNum int, err error) {
				assert.NoError(t, err, "Case %d", caseNum)
			},
		},
		{
			action: func(a publish.AlmanacInfo) error {
				return a.CheckConnectivity(http.DefaultClient)
			},
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusNotFound)
			},
			verifier: func(caseNum int, err error) {
				assert.Regexp(t, regexp.MustCompile("Received non-success status code: 404 Not Found."), err.Error(), "Case %d", caseNum)
			},
		},
		{
			action: func(a publish.AlmanacInfo) error {
				client := &http.Client{Transport: &errorRoundTripper{}}
				return a.CreateProduct(client, "foo")
			},
			verifier: func(caseNum int, err error) {
				assert.Regexp(t, `Almanac request failed: .+`, err.Error(), "Case %d", caseNum)
				// error should not contain authorization header
				assert.NotRegexp(t, `X-Authorization`, err.Error(), "Case %d", caseNum)
			},
		},
		{
			action: func(a publish.AlmanacInfo) error {
				return a.CreateUnit(http.DefaultClient, publish.AlmanacUnit{
					Product: "testProduct",
					Tags:    []string{"tag-1", "tag2"},
				}, "0.0.1", ioutil.Discard)
			},
			handler: func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, "/v1/units", r.URL.String())
				assert.Equal(t, http.MethodPost, r.Method)
				_, ok := r.Header[http.CanonicalHeaderKey("X-timestamp")]
				assert.True(t, ok)
				_, ok = r.Header[http.CanonicalHeaderKey("X-authorization")]
				assert.True(t, ok)

				bytes, err := ioutil.ReadAll(r.Body)
				require.NoError(t, err)
				assert.JSONEq(t, `{"product":"testProduct","branch":"","revision":"","url":"","metadata":{"version":"0.0.1"},"tags":["tag-1","tag2"]}`, string(bytes))

				_, err = w.Write([]byte("hello"))
				require.NoError(t, err)
			},
			verifier: func(caseNum int, err error) {
				assert.NoError(t, err, "Case %d", caseNum)
			},
		},
		{
			action: func(a publish.AlmanacInfo) error {
				return a.ReleaseProduct(http.DefaultClient, "testProduct", "testBranch", "testRevision")
			},
			handler: func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, "/v1/units/testProduct/testBranch/testRevision/releases", r.URL.String())
				assert.Equal(t, http.MethodPost, r.Method)
				_, ok := r.Header[http.CanonicalHeaderKey("X-timestamp")]
				assert.True(t, ok)
				_, ok = r.Header[http.CanonicalHeaderKey("X-authorization")]
				assert.True(t, ok)

				bytes, err := ioutil.ReadAll(r.Body)
				require.NoError(t, err)
				assert.JSONEq(t, `{"name":"GA"}`, string(bytes))

				_, err = w.Write([]byte("hello"))
				require.NoError(t, err)
			},
			verifier: func(caseNum int, err error) {
				assert.NoError(t, err, "Case %d", caseNum)
			},
		},
	} {
		handlerFunc = currCase.handler
		err := currCase.action(a)
		currCase.verifier(i, err)
	}
}
