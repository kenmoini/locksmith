package main

import (
	"crypto/x509"
	"crypto/x509/pkix"
	"math/big"
	"path/filepath"
	"time"
)

func setupCACert(serialNumber int64, commonName string, organization []string, organizationalUnit []string, country []string, province []string, locality []string, streetAddress []string, postalCode []string, addTime []int) *x509.Certificate {
	// set up our CA certificate
	return &x509.Certificate{
		SerialNumber: big.NewInt(serialNumber),
		Subject: pkix.Name{
			CommonName:         commonName,
			Organization:       organization,
			OrganizationalUnit: organizationalUnit,
			Country:            country,
			Province:           province,
			Locality:           locality,
			StreetAddress:      streetAddress,
			PostalCode:         postalCode,
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().AddDate(addTime[0], addTime[1], addTime[2]),
		IsCA:                  true,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		BasicConstraintsValid: true,
	}
}

// createNewCA - creates a new Certificate Authority
func createNewCA(certConfig CertificateConfiguration) (bool, []string, error) {
	checkInputError := false
	var checkInputErrors []string
	var rootSlug string
	var caName string
	var rsaPrivateKeyPassword string

	if certConfig.Subject.CommonName == "" {
		checkInputError = true
		checkInputErrors = append(checkInputErrors, "Missing common name field")
	} else {
		caName = certConfig.Subject.CommonName
		rootSlug = slugger(caName)
	}
	if len(certConfig.Subject.Organization) == 0 {
		checkInputError = true
		checkInputErrors = append(checkInputErrors, "Missing Organization field")
	}
	if len(certConfig.ExpirationDate) != 3 {
		checkInputError = true
		checkInputErrors = append(checkInputErrors, "Missing Expiration Date field")
	}
	if checkInputError {
		return false, checkInputErrors, Stoerr("cert-config-error")
	}

	rootSlugPath := readConfig.Locksmith.PKIRoot + "/roots/" + rootSlug
	rsaPrivateKeyPassword = certConfig.RSAPrivateKeyPassphrase

	//Create root CA directory
	rootCAPath, err := filepath.Abs(rootSlugPath)
	check(err)
	CreateDirectory(rootCAPath)

	// Create certificate requests (CSR) path
	rootCACertRequestsPath := rootCAPath + "/certreqs"
	CreateDirectory(rootCACertRequestsPath)

	// Create certs path
	rootCACertsPath := rootCAPath + "/certs"
	CreateDirectory(rootCACertsPath)

	// Create crls path
	rootCACertRevListPath := rootCAPath + "/crl"
	CreateDirectory(rootCACertRevListPath)

	// Create newcerts path (wtf is newcerts for vs certs?!)
	rootCANewCertsPath := rootCAPath + "/newcerts"
	CreateDirectory(rootCANewCertsPath)

	// Create private path for CA keys
	rootCACertKeysPath := rootCAPath + "/private"
	CreateDirectory(rootCACertKeysPath)

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
		// if there is no private key, create one
		rootPrivKey, rootPubKey, err := generateRSAKeypair(4096)
		check(err)

		rootPrivKeyFile, rootPubKeyFile, err := writeRSAKeyPair(pemEncodeRSAPrivateKey(rootPrivKey, ""), pemEncodeRSAPublicKey(rootPubKey), rootCACertKeysPath+"/ca")
		check(err)
		if rootPrivKeyFile && rootPubKeyFile {
			logStdOut("RSA Key Pair Created")
		}
	}

	// Create CA Object
	rootCA := setupCACert(readSerialNumberAsInt64(rootSlug), certConfig.Subject.CommonName, certConfig.Subject.Organization, certConfig.Subject.OrganizationalUnit, certConfig.Subject.Country, certConfig.Subject.Province, certConfig.Subject.Locality, certConfig.Subject.StreetAddress, certConfig.Subject.PostalCode, certConfig.ExpirationDate)

	// Read in the Private key
	privateKeyFromFile := GetPrivateKey(rootCACertKeysPath+"/ca.priv.pem", rsaPrivateKeyPassword)

	// Read in the Public key
	pubKeyFromFile := GetPublicKey(rootCACertKeysPath + "/ca.pub.pem")

	// Byte Encode the Certificate - https://golang.org/pkg/crypto/x509/#CreateCertificate
	caBytes, err := CreateCert(rootCA, rootCA, pubKeyFromFile, privateKeyFromFile)
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
	return true, []string{"Finished"}, nil
}
