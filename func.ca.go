package main

import (
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"math/big"
	"path/filepath"
	"time"
)

func setupCACert(serialNumber int64, organization string, country string, province string, locality string, streetAddress string, postalCode string, addTime []int) *x509.Certificate {
	// set up our CA certificate
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
		NotBefore:             time.Now(),
		NotAfter:              time.Now().AddDate(addTime[0], addTime[1], addTime[2]),
		IsCA:                  true,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		BasicConstraintsValid: true,
	}
}

// createNewRootCAFilesystem
func createNewCAFilesystem(rootSlug string) {
	rootSlugPath := readConfig.Locksmith.PKIRoot + "/roots/" + rootSlug
	//Create root CA directory
	rootCAPath, err := filepath.Abs(rootSlugPath)
	check(err)
	CreateDirectory(rootCAPath)

	rootCACertsPath, err := filepath.Abs(rootSlugPath + "/certs")
	check(err)
	CreateDirectory(rootCACertsPath)

	rootCACertRevListPath, err := filepath.Abs(rootSlugPath + "/crls")
	check(err)
	CreateDirectory(rootCACertRevListPath)

	rootCACertKeysPath, err := filepath.Abs(rootSlugPath + "/keys")
	check(err)
	CreateDirectory(rootCACertKeysPath)

	rootCACertRequestsPath, err := filepath.Abs(rootSlugPath + "/reqs")
	check(err)
	CreateDirectory(rootCACertRequestsPath)

	rootCACertSerialFilePath, err := filepath.Abs(rootSlugPath + "/serial.txt")
	check(err)

	// Check to see if there is a serial file
	serialFile, err := WriteFile(rootCACertSerialFilePath, "1", 0600, false)
	check(err)
	if serialFile {
		logStdOut("Created serial file")
	} else {
		logStdOut("Serial file exists")
	}

	rootCACrlnumFilePath, err := filepath.Abs(rootSlugPath + "/crlnumber.txt")
	check(err)

	// Check to see if there is a crlNum file
	crlNumFile, err := WriteFile(rootCACrlnumFilePath, "1", 0600, false)
	check(err)
	if crlNumFile {
		logStdOut("Created crlnum file")
	} else {
		logStdOut("crlnum file exists")
	}

	// Check for certificate authority key pair
	caKeyPath, err := filepath.Abs(rootSlugPath + "/keys/ca.priv.pem")
	check(err)
	caKeyCheck, err := FileExists(caKeyPath)
	check(err)
	if !caKeyCheck {
		rootPrivKey, rootPubKey, err := generateRSAKeypair(4096)
		check(err)

		rootPrivKeyFile, rootPubKeyFile, err := writeRSAKeyPair(pemEncodeRSAPrivateKey(rootPrivKey), pemEncodeRSAPublicKey(rootPubKey), rootCACertKeysPath+"/ca")
		check(err)
		if rootPrivKeyFile && rootPubKeyFile {
			logStdOut("RSA Key Pair Created")
		}

		// Create CA Object

		rootCA := setupCACert(readSerialNumberAsInt64(rootSlug), "Kemo Labs", "US", "NC", "Charlotte", "420 Thug Ln", "28204", []int{10, 0, 0})

		// Byte Encode the Certificate
		caBytes, err := x509.CreateCertificate(rand.Reader, rootCA, rootCA, rootPubKey, rootPrivKey)
		check(err)

		// Check for certificate file
		certificateFileCheck, err := FileExists(rootCACertsPath + "/ca.cert")
		if !certificateFileCheck {
			// Write Certificate file
			certificateFile, err := writeCertificateFile(pemEncodeCertificate(caBytes), rootCACertsPath+"/ca.cert")
			check(err)
			if certificateFile {
				logStdOut("Created Root CA Certificate file")
			}
		}
	}
}
