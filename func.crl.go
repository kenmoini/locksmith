package main

import (
	"crypto"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"time"
)

// CreateCRLObject will create the CRL Object
func CreateCRLObject(certList []pkix.RevokedCertificate, key crypto.Signer, issuingCert *x509.Certificate, expiryTime time.Time) ([]byte, error) {
	if certList == nil {
		return nil, newError("Missing certificate list required to create CRL.")
	}
	if key == nil {
		return nil, newError("Missing issuing key required to create CRL.")
	}
	if issuingCert == nil {
		return nil, newError("Missing issuing certificate required to create CRL.")
	}
	if expiryTime == nil {
		expiryTime = time.Now().AddDate(1, 0, 0)
	}
	crlBytes, err := issuingCert.CreateCRL(rand.Reader, key, certList, time.Now(), expiryTime)
	check(err)

	return crlBytes, err
}
