package locksmith

import (
	"bytes"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"io/ioutil"
	"math/big"
	"net"
	"time"
)

// createNewCertificateFromCSR allows the maturation of a CSR to a Certificate
func createNewCertificateFromCSR(signingCAPath string, signingCAPassphrase string, csr *x509.CertificateRequest) (certCreated bool, certificate *x509.Certificate, messages []string, err error) {
	// Check to make sure the ca.pem file exists
	signingCACertExists, err := FileExists(signingCAPath + "/certs/ca.pem")
	if !signingCACertExists {
		// Signing CA does not exist, can't sign certificate
		return false, &x509.Certificate{}, []string{"Signing CA Certificate does not exist!"}, Stoerr("no-signing-ca-certificate")
	}
	// Open Signing CA Certificate file
	//caCertFileBytes, err := ReadCertFromFile(signingCAPath + "/certs/ca.pem")
	check(err)

	// Open Signing CA Private Key file
	signingCAPrivateKeyExists, err := FileExists(signingCAPath + "/private/ca.priv.pem")
	check(err)

	if !signingCAPrivateKeyExists {
		// Signing CA private key does not exist, can't sign certificate
		return false, &x509.Certificate{}, []string{"Signing CA Private Key does not exist!"}, Stoerr("no-signing-ca-key")
	}
	// Open Signing CA Private Key
	//signingCAPrivateKey := GetPrivateKey(signingCAPath+"/private/ca.priv.pem", signingCAPassphrase)

	// Assemble certificate

	return
}

// setupServerCert
func setupServerCert(serialNumber int64, organization string, country string, province string, locality string, streetAddress string, postalCode string, addTime []int) *x509.Certificate {
	return &x509.Certificate{
		SerialNumber: big.NewInt(serialNumber),
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
