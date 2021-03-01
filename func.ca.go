package main

import (
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"math/big"
	"path/filepath"
	"time"
)

func setupCACert(serialNumber int64, commonName string, organization string, organizationalUnit string, country string, province string, locality string, streetAddress string, postalCode string, addTime []int) *x509.Certificate {
	// set up our CA certificate
	return &x509.Certificate{
		SerialNumber: big.NewInt(serialNumber),
		Subject: pkix.Name{
			CommonName:         commonName,
			Organization:       []string{organization},
			OrganizationalUnit: []string{organizationalUnit},
			Country:            []string{country},
			Province:           []string{province},
			Locality:           []string{locality},
			StreetAddress:      []string{streetAddress},
			PostalCode:         []string{postalCode},
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().AddDate(addTime[0], addTime[1], addTime[2]),
		IsCA:                  true,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		BasicConstraintsValid: true,
	}
}

// createNewRootCAFilesystem
func createNewCAFilesystem(rootSlug string, caName string, rsaPrivateKeyPassword string) {
	rootSlugPath := readConfig.Locksmith.PKIRoot + "/roots/" + rootSlug
	//Create root CA directory
	rootCAPath, err := filepath.Abs(rootSlugPath)
	check(err)
	CreateDirectory(rootCAPath)

	// Create certs path
	rootCACertsPath := rootCAPath + "/certs"
	CreateDirectory(rootCACertsPath)

	// Create newcerts path (wtf is newcerts for vs certs?!)
	rootCANewCertsPath := rootCAPath + "/newcerts"
	CreateDirectory(rootCANewCertsPath)

	// Create crls path
	rootCACertRevListPath := rootCAPath + "/crl"
	CreateDirectory(rootCACertRevListPath)

	// Create private path for CA keys
	rootCACertKeysPath := rootCAPath + "/private"
	CreateDirectory(rootCACertKeysPath)

	// Create certificate requests (CSR) path
	rootCACertRequestsPath := rootCAPath + "/certreqs"
	CreateDirectory(rootCACertRequestsPath)

	// Create intermediate CA path
	rootCAIntermediateCAPath := rootCAPath + "/intermed-ca"
	CreateDirectory(rootCAIntermediateCAPath)

	//  CREATE INDEX DATABASE FILE
	rootCACertIndexFilePath := rootCAPath + "/ca.index"
	// Check to see if there is an Index file
	IndexFile, err := WriteFile(rootCACertIndexFilePath, "", 0600, false)
	check(err)
	if IndexFile {
		logStdOut("Created Index file")
	} else {
		logStdOut("Index file exists")
	}

	//  CREATE SERIAL FILE
	rootCACertSerialFilePath := rootCAPath + "/serial.txt"
	// Check to see if there is a serial file
	serialFile, err := WriteFile(rootCACertSerialFilePath, "01", 0600, false)
	check(err)
	if serialFile {
		logStdOut("Created serial file")
	} else {
		logStdOut("Serial file exists")
	}

	//  CREATE CERTIFICATE REVOKATION NUMBER FILE
	rootCACrlnumFilePath := rootCAPath + "/crlnumber.txt"
	// Check to see if there is a crlNum file
	crlNumFile, err := WriteFile(rootCACrlnumFilePath, "00", 0600, false)
	check(err)
	if crlNumFile {
		logStdOut("Created crlnum file")
	} else {
		logStdOut("crlnum file exists")
	}

	// Check for certificate authority key pair
	caKeyCheck, err := FileExists(rootCACertKeysPath + "/ca.priv.pem")
	check(err)

	if !caKeyCheck {
		rootPrivKey, rootPubKey, err := generateRSAKeypair(4096)
		check(err)

		rootPrivKeyFile, rootPubKeyFile, err := writeRSAKeyPair(pemEncodeRSAPrivateKey(rootPrivKey, ""), pemEncodeRSAPublicKey(rootPubKey), rootCACertKeysPath+"/ca")
		check(err)
		if rootPrivKeyFile && rootPubKeyFile {
			logStdOut("RSA Key Pair Created")
		}
	}

	// Create CA Object
	rootCA := setupCACert(readSerialNumberAsInt64(rootSlug), "Kemo Labs Root Certificate Authority", "Kemo Labs", "Kemo Labs Cyber and Information Security", "US", "NC", "Charlotte", "420 Thug Ln", "28204", []int{10, 0, 0})

	// Read in the Private key
	privateKeyFromFile := GetPrivateKey(rootCACertKeysPath+"/ca.priv.pem", rsaPrivateKeyPassword)

	// Read in the Public key
	pubKeyFromFile := GetPublicKey(rootCACertKeysPath + "/ca.pub.pem")

	// Byte Encode the Certificate
	caBytes, err := x509.CreateCertificate(rand.Reader, rootCA, rootCA, pubKeyFromFile, privateKeyFromFile)
	check(err)

	// Check for certificate file
	certificateFileCheck, err := FileExists(rootCACertsPath + "/ca.pem")
	if !certificateFileCheck {
		// Write Certificate file
		certificateFile, err := writeCertificateFile(pemEncodeCertificate(caBytes), rootCACertsPath+"/ca.pem")
		check(err)
		if certificateFile {
			logStdOut("Created Root CA Certificate file")
		}
	}
}
