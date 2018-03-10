// +build ignore

// package main writes the crypto material used by unit tests. Run using "go generate".
package main

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"math/big"
	"net"
	"os"
	"time"
)

const (
	caCertFile     = "ca-cert.pem"
	caKeyFile      = "ca-key.pem"
	serverCertFile = "server-cert.pem"
	serverKeyFile  = "server-key.pem"
	clientCertFile = "client-cert.pem"
	clientKeyFile  = "client-key.pem"

	combinedCertsFile = "combined-certs.pem"
	certWithKeyFile   = "cert-with-key.pem"
)

var (
	serverTemplate = x509.Certificate{
		NotBefore:             time.Now(),
		NotAfter:              time.Now().Add(10 * 365 * 24 * time.Hour),
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		BasicConstraintsValid: true,
		DNSNames:              []string{"localhost"},
		IPAddresses:           []net.IP{net.IPv4(127, 0, 0, 1)},
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
	}
	clientTemplate = x509.Certificate{
		NotBefore:             time.Now(),
		NotAfter:              time.Now().Add(10 * 365 * 24 * time.Hour),
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		BasicConstraintsValid: true,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth},
	}
)

func main() {
	caCertStr, caKeyStr := newCAKeyPair(1, "Root CA")
	block, _ := pem.Decode([]byte(caCertStr))
	if block == nil {
		panic("failed to parse certificate PEM")
	}
	caCert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		panic(err)
	}

	serverCertStr, serverKeyStr := newSignedKeyPair(2, "Test Org", "localhost", &serverTemplate, caCert, caKeyStr)
	clientCertStr, clientKeyStr := newSignedKeyPair(3, "Test Org", "client", &clientTemplate, caCert, caKeyStr)

	mustWriteFile(caCertFile, []byte(caCertStr), 0644)
	mustWriteFile(caKeyFile, []byte(caKeyStr), 0644)

	mustWriteFile(serverCertFile, []byte(serverCertStr), 0644)
	mustWriteFile(serverKeyFile, []byte(serverKeyStr), 0644)

	mustWriteFile(clientCertFile, []byte(clientCertStr), 0644)
	mustWriteFile(clientKeyFile, []byte(clientKeyStr), 0644)

	combinedCerts := fmt.Sprint(caCertStr, serverCertStr)
	mustWriteFile(combinedCertsFile, []byte(combinedCerts), 0644)
	serverCertWithKey := fmt.Sprint(serverCertStr, serverKeyStr)
	mustWriteFile(certWithKeyFile, []byte(serverCertWithKey), 0644)
}

func mustWriteFile(filename string, data []byte, perm os.FileMode) {
	if err := ioutil.WriteFile(filename, data, perm); err != nil {
		panic(err)
	}
}

func newCAKeyPair(serial int64, org string) (string, string) {
	caTemplate := &x509.Certificate{
		SerialNumber: big.NewInt(serial),
		Subject: pkix.Name{
			Organization: []string{org},
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().Add(10 * 365 * 24 * time.Hour),
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		BasicConstraintsValid: true,
		IsCA:        true,
		ExtKeyUsage: []x509.ExtKeyUsage{x509.ExtKeyUsageAny},
	}
	privKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		panic(err)
	}
	certDERBytes, err := x509.CreateCertificate(rand.Reader, caTemplate, caTemplate, &privKey.PublicKey, privKey)
	if err != nil {
		panic(err)
	}
	keyDERBytes := x509.MarshalPKCS1PrivateKey(privKey)
	return derBytesToStrings(certDERBytes, keyDERBytes)
}

func newSignedKeyPair(serial int64, org, cn string, template, caCert *x509.Certificate, caPrivKey string) (string, string) {
	block, _ := pem.Decode([]byte(caPrivKey))
	if block == nil {
		panic("failed to parse certificate PEM")
	}
	signingPrivKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		panic(err)
	}

	privKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		panic(err)
	}

	tmplCopy := *template
	tmplCopy.SerialNumber = big.NewInt(serial)
	tmplCopy.Subject = pkix.Name{
		Organization: []string{org},
		CommonName:   cn,
	}

	certDERBytes, err := x509.CreateCertificate(rand.Reader, &tmplCopy, caCert, &privKey.PublicKey, signingPrivKey)
	if err != nil {
		panic(err)
	}

	return derBytesToStrings(certDERBytes, x509.MarshalPKCS1PrivateKey(privKey))
}

func derBytesToStrings(certDERBytes, keyDERBytes []byte) (string, string) {
	var certBuf bytes.Buffer
	if err := pem.Encode(&certBuf, &pem.Block{Type: "CERTIFICATE", Bytes: certDERBytes}); err != nil {
		panic(err)
	}
	var keyBuf bytes.Buffer
	if err := pem.Encode(&keyBuf, &pem.Block{Type: "RSA PRIVATE KEY", Bytes: keyDERBytes}); err != nil {
		panic(err)
	}
	return certBuf.String(), keyBuf.String()
}
