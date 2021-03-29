package locksmith

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/x509"
	"encoding/pem"
	"io/ioutil"
	"math/big"
	"time"
)

// createNewCertificateFromCSR allows the maturation of a CSR to a Certificate
func createNewCertificateFromCSR(signingCAPath string, signingCAPassphrase string, csr *x509.CertificateRequest, certificateType string, csrPublicKey *rsa.PublicKey) (certCreated bool, certificate *x509.Certificate, messages []string, err error) {
	// Check to make sure the ca.pem file exists
	signingCACertExists, err := FileExists(signingCAPath + "/certs/ca.pem")
	check(err)

	if !signingCACertExists {
		// Signing CA does not exist, can't sign certificate
		return false, &x509.Certificate{}, []string{"Signing CA Certificate does not exist!"}, Stoerr("no-signing-ca-certificate")
	}
	// Open Signing CA Certificate file
	signingCACertFileBytes, err := ReadCACertificate(signingCAPath)
	check(err)

	// Check for Signing CA Private Key file
	signingCAPrivateKeyExists, err := FileExists(signingCAPath + "/private/ca.priv.pem")
	check(err)

	if !signingCAPrivateKeyExists {
		// Signing CA private key does not exist, can't sign certificate
		return false, &x509.Certificate{}, []string{"Signing CA Private Key does not exist!"}, Stoerr("no-signing-ca-key")
	}
	// Open Signing CA Key Pair
	signingCAPrivateKey := GetPrivateKey(signingCAPath+"/private/ca.priv.pem", signingCAPassphrase)
	signingCAPublicKey := GetPublicKey(signingCAPath + "/private/ca.pub.pem")

	// Check for the Signing CA's Serial file
	signingCASerialFileExists, err := FileExists(signingCAPath + "/ca.serial")
	check(err)

	if !signingCASerialFileExists {
		// Signing CA serial file does not exist, can't sign certificate
		return false, &x509.Certificate{}, []string{"Signing CA Serial file does not exist!"}, Stoerr("no-signing-ca-serial-file")
	}
	// Read in the current serial number
	currentSerial := readSerialNumberAsInt64Abs(signingCAPath + "/ca.serial")

	// Check for the Signing CA's Index DB
	signingCAIndexDBExists, err := FileExists(signingCAPath + "/ca.index")
	check(err)

	if !signingCAIndexDBExists {
		// Signing CA Index DB file does not exist, can't sign certificate
		return false, &x509.Certificate{}, []string{"Signing CA Index DB does not exist!"}, Stoerr("no-signing-ca-index-db")
	}

	// Assemble certificate
	switch certificateType {
	case "authority":
		// do something for CA certs
	case "server":
	default:
		// by default, we'll generate a server type certificate
		certificate = setupServerCert(currentSerial, csr, []int{1, 0, 1}, signingCAPublicKey)
	}

	// Sign Certificate
	certBytes, err := CreateCert(certificate, signingCACertFileBytes, csrPublicKey, signingCAPrivateKey)
	check(err)

	// Save Signed Certificate File
	certificatePath := signingCAPath + "/certs/" + slugger(certificate.Subject.CommonName) + ".pem"
	certificateFile, err := writeCertificateFile(pemEncodeCertificate(certBytes), certificatePath)
	check(err)
	if !certificateFile {
		return false, &x509.Certificate{}, []string{"Certificate Creation Failure!"}, err
	}

	cert, err := ReadCertFromFile(certificatePath)
	check(err)

	// Increment Signing CA Serial Number
	increaseSerial, err := IncreaseSerialNumberAbs(signingCAPath + "/ca.serial")
	check(err)

	if !increaseSerial {
		return false, &x509.Certificate{}, []string{"Signing CA Serial Increment Error"}, err
	}

	// Add Certificate to Signing CA Index DB
	addedEntry, err := AddEntryToCAIndex(signingCAPath+"/ca.index", certificatePath)
	check(err)

	if !addedEntry {
		return false, &x509.Certificate{}, []string{"Signing CA Index Entry Error"}, err
	}

	// Finally, return the certificate
	return true, cert, []string{"Certificate created successfully!"}, nil
}

// setupServerCert creates the Certificate structure of a Server type certificate
func setupServerCert(serialNumber int64, csr *x509.CertificateRequest, addTime []int, signingPubKey *rsa.PublicKey) *x509.Certificate {

	// Set time for UTC format
	currentTime := time.Now()
	yesterdayTime := time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day(), 0, 0, 0, 0, time.UTC).Add(-24 * time.Hour)

	// Set up SubjectKeyID from Public Key
	publicKeyBytes, _, err := marshalPublicKey(signingPubKey)
	check(err)
	h := sha1.Sum(publicKeyBytes)
	subjectKeyID := h[:]

	return &x509.Certificate{
		SignatureAlgorithm:    x509.SHA512WithRSA,
		SerialNumber:          big.NewInt(serialNumber),
		Subject:               csr.Subject,
		IPAddresses:           csr.IPAddresses,
		URIs:                  csr.URIs,
		DNSNames:              csr.DNSNames,
		NotBefore:             time.Date(yesterdayTime.Year(), yesterdayTime.Month(), yesterdayTime.Day(), 0, 0, 0, 0, yesterdayTime.Location()),
		NotAfter:              time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day(), 0, 0, 0, 0, time.UTC).AddDate(addTime[0], addTime[1], addTime[2]),
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		KeyUsage:              x509.KeyUsageDigitalSignature,
		AuthorityKeyId:        subjectKeyID,
		BasicConstraintsValid: true,
		IsCA:                  false,
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

// CreateCert is a wrapper for x509.CreateCertificate to switch between parent certificates through the chain
func CreateCert(certTemplate *x509.Certificate, signingCert *x509.Certificate, certPubkey, signingPrivKey interface{}) (cert []byte, err error) {
	return x509.CreateCertificate(rand.Reader, certTemplate, signingCert, certPubkey, signingPrivKey)
}

// ReadCertFromFile wraps the needed functions to safely read a PEM certificate
func ReadCertFromFile(path string) (*x509.Certificate, error) {
	// Check if the file exists
	certificateFileCheck, err := FileExists(path)
	if !certificateFileCheck {
		return nil, err
	}

	// Read in PEM file
	pem, err := readPEMFile(path, "CERTIFICATE")
	check(err)

	// Decode to Certfificate object
	return x509.ParseCertificate(pem.Bytes)
}
