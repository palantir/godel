// Copyright 2016 Palantir Technologies. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package tlsconfig_test

import (
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/palantir/pkg/tlsconfig"
)

const (
	caCertFile     = "testdata/ca-cert.pem"
	serverCertFile = "testdata/server-cert.pem"
	serverKeyFile  = "testdata/server-key.pem"
	clientCertFile = "testdata/client-cert.pem"
	clientKeyFile  = "testdata/client-key.pem"
)

func TestUseTLSConfigClientAuthConnection(t *testing.T) {
	for i, tc := range []struct {
		name              string
		serverTLSProvider tlsconfig.TLSCertProvider
		serverParams      []tlsconfig.ServerParam
		clientParams      []tlsconfig.ClientParam
	}{
		{
			name:              "TLS with client cert required",
			serverTLSProvider: tlsconfig.TLSCertFromFiles(serverCertFile, serverKeyFile),
			serverParams: []tlsconfig.ServerParam{
				tlsconfig.ServerClientAuthType(tls.RequireAndVerifyClientCert),
				tlsconfig.ServerClientCAs(tlsconfig.CertPoolFromCAFiles(caCertFile)),
			},
			clientParams: []tlsconfig.ClientParam{
				tlsconfig.ClientKeyPairFiles(clientCertFile, clientKeyFile),
				tlsconfig.ClientRootCAs(tlsconfig.CertPoolFromCAFiles(caCertFile)),
			},
		},
		{
			name:              "TLS with no client cert",
			serverTLSProvider: tlsconfig.TLSCertFromFiles(serverCertFile, serverKeyFile),
			serverParams: []tlsconfig.ServerParam{
				tlsconfig.ServerClientAuthType(tls.NoClientCert),
				tlsconfig.ServerClientCAs(tlsconfig.CertPoolFromCAFiles(caCertFile)),
			},
			clientParams: []tlsconfig.ClientParam{
				tlsconfig.ClientRootCAs(tlsconfig.CertPoolFromCAFiles(caCertFile)),
			},
		},
	} {
		func() {
			server := httptest.NewUnstartedServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
				fmt.Fprintf(rw, "OK: %s", req.URL.Path)
			}))
			serverCfg, err := tlsconfig.NewServerConfig(tc.serverTLSProvider, tc.serverParams...)
			require.NoError(t, err)
			server.TLS = serverCfg
			server.StartTLS()
			defer server.Close()

			clientCfg, err := tlsconfig.NewClientConfig(tc.clientParams...)
			require.NoError(t, err)

			client := &http.Client{
				Transport: &http.Transport{
					TLSClientConfig: clientCfg,
				},
			}

			resp, err := client.Get(server.URL + "/hello")
			require.NoError(t, err)
			bytes, err := ioutil.ReadAll(resp.Body)
			require.NoError(t, err)

			assert.Equal(t, "OK: /hello", string(bytes), "Case %d: %s", i, tc.name)
		}()
	}
}
