// Copyright 2016 Palantir Technologies. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package tlsconfig_test

import (
	"crypto/tls"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/palantir/pkg/tlsconfig"
)

func TestNewServerConfig(t *testing.T) {
	for currCaseNum, currCase := range []struct {
		name          string
		clientCAFiles []string
		authType      tls.ClientAuthType
		cipherSuites  []uint16
		nextProtos    []string
	}{
		{
			name: "defaults",
		},
		{
			name: "caFiles specified",
			clientCAFiles: []string{
				clientCertFile,
			},
		},
		{
			name:     "authType specified",
			authType: tls.NoClientCert,
		},
		{
			name: "cipherSuites specified",
			cipherSuites: []uint16{
				tls.TLS_RSA_WITH_AES_256_CBC_SHA,
			},
		},
		{
			name: "nextProtos specified",
			nextProtos: []string{
				"http/1.1",
			},
		},
	} {
		cfg, err := tlsconfig.NewServerConfig(
			tlsconfig.TLSCertFromFiles(serverCertFile, serverKeyFile),
			tlsconfig.ServerClientCAFiles(currCase.clientCAFiles...),
			tlsconfig.ServerClientAuthType(currCase.authType),
			tlsconfig.ServerCipherSuites(currCase.cipherSuites...),
			tlsconfig.ServerNextProtos(currCase.nextProtos...),
		)
		require.NoError(t, err)
		assert.NotNil(t, cfg, "Case %d: %s", currCaseNum, currCase.name)
	}
}

func TestNewServerConfigErrors(t *testing.T) {
	for currCaseNum, currCase := range []struct {
		name          string
		clientCAFiles []string
		wantError     string
	}{
		{
			name: "invalid CA file",
			clientCAFiles: []string{
				serverKeyFile,
			},
			wantError: "failed to create certificate pool: no certificates detected in file testdata/server-key.pem",
		},
	} {
		cfg, err := tlsconfig.NewServerConfig(
			tlsconfig.TLSCertFromFiles(serverCertFile, serverKeyFile),
			tlsconfig.ServerClientCAFiles(currCase.clientCAFiles...),
		)
		require.Error(t, err, fmt.Sprintf("Case %d: %s", currCaseNum, currCase.name))
		assert.EqualError(t, err, currCase.wantError, "Case %d: %s", currCaseNum, currCase.name)
		assert.Nil(t, cfg, "Case %d: %s", currCaseNum, currCase.name)
	}
}
