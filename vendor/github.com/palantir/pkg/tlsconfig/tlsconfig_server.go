// Copyright 2016 Palantir Technologies. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package tlsconfig

import (
	"crypto/tls"
	"fmt"
)

type ServerParam interface {
	configureServer(*tls.Config) error
}

type serverParam func(*tls.Config) error

func (p serverParam) configureServer(cfg *tls.Config) error {
	return p(cfg)
}

// NewServerConfig returns a tls.Config that is suitable to use by a server in 2-way TLS connections configured with
// the provided parameters. The provided TLSCertProvider is used as the source for the private key and certificate that
// the server presents to clients.
func NewServerConfig(tlsCertProvider TLSCertProvider, params ...ServerParam) (*tls.Config, error) {
	if tlsCertProvider == nil {
		return nil, fmt.Errorf("tlsCertProvider provided to NewServerConfig was nil")
	}
	configurers := []configurer{authKeyPairParam(tlsCertProvider)}
	for _, p := range params {
		configurers = append(configurers, configurer(p.configureServer))
	}
	return configureTLSConfig(configurers...)
}

// ServerClientCAFiles configures the server with the CA certificates used to verify the certificates provided by
// clients. If this parameter is not provided, then the default system CAs are used.
func ServerClientCAFiles(files ...string) ServerParam {
	return ServerClientCAs(CertPoolFromCAFiles(files...))
}

// ServerClientCAs configures the server with the CA certificates used to verify the certificates provided by clients.
// If this parameter is not provided, then the default system CAs are used.
func ServerClientCAs(certPoolProvider CertPoolProvider) ServerParam {
	return serverParam(func(cfg *tls.Config) error {
		certPool, err := certPoolProvider()
		if err != nil {
			return fmt.Errorf("failed to create certificate pool: %v", err)
		}
		cfg.ClientCAs = certPool
		return nil
	})
}

// ServerClientAuthType sets the default client auth type required by the server. If this parameter is not provided,
// defaults to NoClientCert.
func ServerClientAuthType(authType tls.ClientAuthType) ServerParam {
	return serverParam(func(cfg *tls.Config) error {
		cfg.ClientAuth = authType
		return nil
	})
}

// ServerCipherSuites sets the cipher suites supported by the server. If this parameter is not provided,
// defaultCipherSuites is used.
func ServerCipherSuites(cipherSuites ...uint16) ServerParam {
	return serverParam(cipherSuitesParam(cipherSuites...))
}

// ServerNextProtos sets the list of application level protocols supported by
// the server e.g. "http/1.1" or "h2".
func ServerNextProtos(protos ...string) ServerParam {
	return serverParam(func(cfg *tls.Config) error {
		cfg.NextProtos = protos
		return nil
	})
}
