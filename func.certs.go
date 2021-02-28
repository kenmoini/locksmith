package main

import (
	"bytes"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"io/ioutil"
	"math/big"
	"net"
	"time"
)

// setupServerCert
func setupServerCert(organization string, country string, province string, locality string, streetAddress string, postalCode string, addTime []int) *x509.Certificate {
	return &x509.Certificate{
		SerialNumber: big.NewInt(2019),
		Subject: pkix.Name{
			Organization:  []string{organization},
			Country:       []string{country},
			Province:      []string{province},
			Locality:      []string{locality},
			StreetAddress: []string{streetAddress},
			PostalCode:    []string{postalCode},
		},
		IPAddresses:  []net.IP{net.IPv4(127, 0, 0, 1), net.IPv6loopback},
		NotBefore:    time.Now(),
		NotAfter:     time.Now().AddDate(addTime[0], addTime[1], addTime[2]),
		SubjectKeyId: []byte{1, 2, 3, 4, 6},
		ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:     x509.KeyUsageDigitalSignature,
	}
}

// pemEncodeCertificate
func pemEncodeCertificate(certByte []byte) *bytes.Buffer {
	pemRet := new(bytes.Buffer)
	pem.Encode(pemRet, &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: certByte,
	})
	return pemRet
}

// writeCertificateFile
func writeCertificateFile(certPem *bytes.Buffer, path string) (bool, error) {
	pemByte, _ := ioutil.ReadAll(certPem)
	keyFile, err := WriteByteFile(path, pemByte, 0600, false)
	if err != nil {
		return false, err
	}
	return keyFile, nil
}
