package main

import (
	"crypto"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"time"
)

// NewCRLFromFile takes in a list of serial numbers, one per line, as well as the issuing certificate
// of the CRL, and the private key. This function is then used to parse the list and generate a CRL
/*
func NewCRLFromFile(serialList, issuerFile, keyFile []byte, expiryTime string) ([]byte, error) {

	var revokedCerts []pkix.RevokedCertificate
	var oneWeek = time.Duration(604800) * time.Second

	expiryInt, err := strconv.ParseInt(expiryTime, 0, 32)
	if err != nil {
		return nil, err
	}
	newDurationFromInt := time.Duration(expiryInt) * time.Second
	newExpiryTime := time.Now().Add(newDurationFromInt)
	if expiryInt == 0 {
		newExpiryTime = time.Now().Add(oneWeek)
	}

	// Parse the PEM encoded certificate
	issuerCert, err := helpers.ParseCertificatePEM(issuerFile)
	if err != nil {
		return nil, err
	}

	// Split input file by new lines
	individualCerts := strings.Split(string(serialList), "\n")

	// For every new line, create a new revokedCertificate and add it to slice
	for _, value := range individualCerts {
		if len(strings.TrimSpace(value)) == 0 {
			continue
		}

		tempBigInt := new(big.Int)
		tempBigInt.SetString(value, 10)
		tempCert := pkix.RevokedCertificate{
			SerialNumber:   tempBigInt,
			RevocationTime: time.Now(),
		}
		revokedCerts = append(revokedCerts, tempCert)
	}

	strPassword := os.Getenv("CFSSL_CA_PK_PASSWORD")
	password := []byte(strPassword)
	if strPassword == "" {
		password = nil
	}

	// Parse the key given
	key, err := helpers.ParsePrivateKeyPEMWithPassword(keyFile, password)
	if err != nil {
		log.Debug("Malformed private key %v", err)
		return nil, err
	}

	return CreateCRLObject(revokedCerts, key, issuerCert, newExpiryTime)
}
*/

// CreateCRLObject will create the CRL Object
func CreateCRLObject(certList []pkix.RevokedCertificate, key crypto.Signer, issuingCert *x509.Certificate, expiryTime time.Time) ([]byte, error) {
	if certList == nil {
		return nil, Stoerr("Missing certificate list required to create CRL.")
	}
	if key == nil {
		return nil, Stoerr("Missing issuing key required to create CRL.")
	}
	if issuingCert == nil {
		return nil, Stoerr("Missing issuing certificate required to create CRL.")
	}
	if expiryTime.IsZero() {
		expiryTime = time.Now().AddDate(1, 0, 0)
	}
	crlBytes, err := issuingCert.CreateCRL(rand.Reader, key, certList, time.Now(), expiryTime)
	check(err)

	return crlBytes, err
}
