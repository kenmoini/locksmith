package main

import (
	"bytes"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/asn1"
	"encoding/pem"
	"io/ioutil"
)

// generateCSR takes the full lifecycle of generating and saving a CSR
func generateCSR(path string, signingKey interface{}, organization []string, country []string, province []string, locality []string, streetAddress []string, postalCode []string, isCA bool) (bool, error) {
	csrTemplate := setupCSR(organization, country, province, locality, streetAddress, postalCode, isCA)

	csr, err := createCSR(csrTemplate, signingKey)
	check(err)

	csrPEM := pemEncodeCSR(csr)

	return writePEMFile(csrPEM, path)
}

// setupCSR creates configuration information and returns a CSR Template
func setupCSR(organization []string, country []string, province []string, locality []string, streetAddress []string, postalCode []string, isCA bool) *x509.CertificateRequest {
	names := pkix.Name{
		Organization:  organization,
		Country:       country,
		Province:      province,
		Locality:      locality,
		StreetAddress: streetAddress,
		PostalCode:    postalCode,
	}
	if isCA {
		val, err := asn1.Marshal(basicConstraints{true, 0})
		check(err)

		return &x509.CertificateRequest{
			Subject:            names,
			SignatureAlgorithm: x509.SHA512WithRSA,
			ExtraExtensions: []pkix.Extension{
				{
					// This identifies that the CSR is a CA
					Id:       asn1.ObjectIdentifier{2, 5, 29, 19},
					Value:    val,
					Critical: true,
				},
			},
		}
	} else {
		return &x509.CertificateRequest{
			Subject:            names,
			SignatureAlgorithm: x509.SHA512WithRSA,
		}
	}
}

// createCSR is a wrapper for x509.CreateCertificateRequest
// template is a CSR template, priv is the CSR requester private key
func createCSR(template *x509.CertificateRequest, priv interface{}) ([]byte, error) {
	return x509.CreateCertificateRequest(rand.Reader, template, priv)
}

// pemEncodeCSR encodes a CreateCertificateRequest DER byte stream to a PEM
func pemEncodeCSR(certByte []byte) *bytes.Buffer {
	pemRet := new(bytes.Buffer)
	pem.Encode(pemRet, &pem.Block{
		Type:  "CERTIFICATE REQUEST",
		Bytes: certByte,
	})
	return pemRet
}

// writePEMFile takes a PEM encoded bytes stream and saves it to a file
func writePEMFile(certPem *bytes.Buffer, path string) (bool, error) {
	pemByte, _ := ioutil.ReadAll(certPem)
	keyFile, err := WriteByteFile(path, pemByte, 0600, false)
	if err != nil {
		return false, err
	}
	return keyFile, nil
}
