package main

import (
	"crypto/x509"
	"crypto/x509/pkix"
	"math/big"
	"net"
	"strings"
	"time"
)

// splitSlugToPath takes a slug string and splits it into the relative path
// eg converts "example-labs-root-certificate-authority/example-labs-ica/server-signing-ca" to "example-labs-root-certificate-authority/intermed-ca/example-labs-ica/intermed-ca/server-signing-ca/"
func splitSlugToPath(slug string) string {
	splitPath := strings.Split(strings.TrimSuffix(strings.TrimPrefix(strings.ToLower(slug), "/"), "/"), "/")
	var path string
	for i, part := range splitPath {
		path = path + part + "/"
		if i != (len(splitPath) - 1) {
			path = path + "intermed-ca/"
		}
	}
	return path
}

// splitCommonNamesToPath takes a CN string and splits it into the relative path while slugging
// eg, converts "Example Labs Root Certificate Authority/Example Labs ICA/Server Signing CA" to "example-labs-root-certificate-authority/intermed-ca/example-labs-ica/intermed-ca/server-signing-ca/"
func splitCommonNamesToPath(cnPath string) string {
	splitPath := strings.Split(strings.TrimSuffix(strings.TrimPrefix(strings.ToLower(cnPath), "/"), "/"), "/")
	var path string
	for i, part := range splitPath {
		path = path + slugger(part) + "/"
		if i != (len(splitPath) - 1) {
			path = path + "intermed-ca/"
		}
	}
	return path
}

// setupIntermediateCACert
func setupIntermediateCACert(serialNumber int64, organization string, organizationalUnit string, country string, province string, locality string, streetAddress string, postalCode string, addTime []int) *x509.Certificate {
	return &x509.Certificate{
		SerialNumber: big.NewInt(serialNumber),
		Subject: pkix.Name{
			Organization:       []string{organization},
			OrganizationalUnit: []string{organizationalUnit},
			Country:            []string{country},
			Province:           []string{province},
			Locality:           []string{locality},
			StreetAddress:      []string{streetAddress},
			PostalCode:         []string{postalCode},
		},
		IPAddresses:           []net.IP{net.IPv4(127, 0, 0, 1), net.IPv6loopback},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().AddDate(addTime[0], addTime[1], addTime[2]),
		IsCA:                  true,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		BasicConstraintsValid: true,
	}
}
