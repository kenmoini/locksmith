package locksmith

import (
	"bytes"
	"crypto"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/asn1"
	"encoding/pem"
	"math/big"
	"time"
)

// SetupNewCRLTemplate wraps a RevokationList type with a bit of pre-processing
func SetupNewCRLTemplate(SignatureAlgorithm x509.SignatureAlgorithm, nextUpdate time.Time) *x509.RevocationList {
	if nextUpdate.IsZero() {
		nextUpdate = time.Now().AddDate(1, 0, 0)
	}

	return &x509.RevocationList{
		SignatureAlgorithm:  SignatureAlgorithm,
		RevokedCertificates: nil,
		Number:              big.NewInt(0),
		ThisUpdate:          time.Now(),
		NextUpdate:          nextUpdate,
		ExtraExtensions:     []pkix.Extension{}}
}

// NewCRL basically just wraps CreateRevocationList in order to create a new blank CRL
func NewCRL(template *x509.RevocationList, issuer *x509.Certificate, priv crypto.Signer) ([]byte, error) {
	return x509.CreateRevocationList(rand.Reader, template, issuer, priv)
}

// PEMEncodeCRL encodes a CreateCertificateRequest DER byte stream to a PEM
func PEMEncodeCRL(certByte []byte) *bytes.Buffer {
	pemRet := new(bytes.Buffer)
	pem.Encode(pemRet, &pem.Block{
		Type:  "X509 CRL",
		Bytes: certByte,
	})
	return pemRet
}

// CreateNewCRLForCA wraps all the processes needed to create a new CRL for a CA
func CreateNewCRLForCA(certificate *x509.Certificate, privateKey crypto.Signer, path string) (bool, error) {
	// Create the template
	crlTemplate := SetupNewCRLTemplate(certificate.SignatureAlgorithm, time.Now().AddDate(1, 0, 0))

	// Take the SAN data from the Certificate and format for IAN
	issuerBytes, err := marshalIANs(certificate.DNSNames, certificate.EmailAddresses, certificate.IPAddresses, certificate.URIs)
	check(err)
	issuerAltName := pkix.Extension{Id: asn1.ObjectIdentifier{2, 5, 29, 18}, Critical: false, Value: issuerBytes}

	// Add IAN Extension to CRL Template
	crlTemplate.ExtraExtensions = []pkix.Extension{issuerAltName}

	// Create an actual CRL object
	crlObject, err := NewCRL(crlTemplate, certificate, privateKey)
	check(err)

	// PEM Encode the object
	pemBytes := PEMEncodeCRL(crlObject)

	// Save the PEM to a file
	return writePEMFile(pemBytes, path)
}

// ReadCRLFromFile just wraps a byte reader and CRL Decoder
func ReadCRLFromFile(path string) (*x509.Certificate, error) {
	// Check if the file exists
	certificateFileCheck, err := FileExists(path)
	if !certificateFileCheck {
		return nil, err
	}

	// Read in PEM file
	pem, err := readPEMFile(path, "X509 CRL")
	check(err)

	// Decode to Certfificate object
	return x509.ParseCertificate(pem.Bytes)
}

/*

Wtf is all this legacy code?

*/

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

	return nil, nil
}
