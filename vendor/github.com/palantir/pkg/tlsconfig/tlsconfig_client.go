// Copyright 2016 Palantir Technologies. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package tlsconfig

import (
	"crypto/tls"
	"fmt"
)

// TLSCertProvider is a function that returns a tls.Certificate used for TLS communication.
type TLSCertProvider func() (tls.Certificate, error)

// NewClientConfig returns a tls.Config that is suitable to use by a client in 2-way TLS connections configured with
// the provided parameters.
func NewClientConfig(params ...ClientParam) (*tls.Config, error) {
	configurers := make([]configurer, len(params))
	for i, p := range params {
		configurers[i] = configurer(p.configureClient)
	}
	return configureTLSConfig(configurers...)
}

type ClientParam interface {
	configureClient(*tls.Config) error
}

type clientParam func(*tls.Config) error

func (p clientParam) configureClient(cfg *tls.Config) error {
	return p(cfg)
}

// ClientKeyPairFiles configures the client with the key pair that it should present to servers when communicating using
// TLS with client authentication (2-way SSL). If this parameter is not provided, the client will not present a
// certificate.
func ClientKeyPairFiles(certFile, keyFile string) ClientParam {
	return ClientKeyPair(TLSCertFromFiles(certFile, keyFile))
}

// ClientKeyPair configures the client with the key pair that it should present to servers when communicating using TLS
// with client authentication (2-way SSL). If this parameter is not provided, the client will not present a certificate.
func ClientKeyPair(certProvider TLSCertProvider) ClientParam {
	return clientParam(authKeyPairParam(certProvider))
}

// ClientRootCAFiles configures the client with the CA certificates used to verify the certificates provided by servers.
// If this parameter is not provided, then the default system CAs are used.
func ClientRootCAFiles(files ...string) ClientParam {
	return ClientRootCAs(CertPoolFromCAFiles(files...))
}

// ClientRootCAs configures the client with the CA certificates used to verify the certificates provided by servers.
// If this parameter is not provided, then the default system CAs are used.
func ClientRootCAs(certPoolProvider CertPoolProvider) ClientParam {
	return clientParam(func(cfg *tls.Config) error {
		if certPoolProvider == nil {
			return fmt.Errorf("certPoolProvider provided to ClientRootCAs was nil")
		}
		certPool, err := certPoolProvider()
		if err != nil {
			return fmt.Errorf("failed to create certificate pool: %v", err)
		}
		cfg.RootCAs = certPool
		return nil
	})
}

// ClientCipherSuites sets the cipher suites supported by the client. If this parameter is not provided,
// defaultCipherSuites is used.
func ClientCipherSuites(cipherSuites ...uint16) ClientParam {
	return clientParam(cipherSuitesParam(cipherSuites...))
}
