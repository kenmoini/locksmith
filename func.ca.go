package main

import (
	"crypto/x509"
	"crypto/x509/pkix"
	"math/big"
	"time"
)

// setupCACert creates a Certificate resource
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

	// Create the CA base directories and files
	certPaths := setupCAFileStructure(rootSlugPath)

	// Check for certificate authority key pair
	caKeyCheck, err := FileExists(certPaths.RootCACertKeysPath + "/ca.priv.pem")
	check(err)

	if !caKeyCheck {
		// if there is no private key, create one
		rootPrivKey, rootPubKey, err := generateRSAKeypair(4096)
		check(err)

		rootPrivKeyFile, rootPubKeyFile, err := writeRSAKeyPair(pemEncodeRSAPrivateKey(rootPrivKey, ""), pemEncodeRSAPublicKey(rootPubKey), certPaths.RootCACertKeysPath+"/ca")
		check(err)
		if rootPrivKeyFile && rootPubKeyFile {
			logStdOut("RSA Key Pair Created")
		}
	}

	// Read in the Private key
	privateKeyFromFile := GetPrivateKey(certPaths.RootCACertKeysPath+"/ca.priv.pem", rsaPrivateKeyPassword)

	// Read in the Public key
	pubKeyFromFile := GetPublicKey(certPaths.RootCACertKeysPath + "/ca.pub.pem")

	// Create Self-signed Certificate Request
	caCSR, err := generateCSR(certPaths.RootCACertRequestsPath+"/ca.pem",
		privateKeyFromFile,
		certConfig.Subject.CommonName,
		certConfig.Subject.Organization,
		certConfig.Subject.OrganizationalUnit,
		certConfig.Subject.Country,
		certConfig.Subject.Province,
		certConfig.Subject.Locality,
		certConfig.Subject.StreetAddress,
		certConfig.Subject.PostalCode,
		true)
	check(err)

	if caCSR {

		// Read in CSR lol

		// Create Self-signed Certificate
		// Create CA Object
		rootCA := setupCACert(readSerialNumberAsInt64(rootSlug), certConfig.Subject.CommonName, certConfig.Subject.Organization, certConfig.Subject.OrganizationalUnit, certConfig.Subject.Country, certConfig.Subject.Province, certConfig.Subject.Locality, certConfig.Subject.StreetAddress, certConfig.Subject.PostalCode, certConfig.ExpirationDate)

		// Byte Encode the Certificate - https://golang.org/pkg/crypto/x509/#CreateCertificate
		caBytes, err := CreateCert(rootCA, rootCA, pubKeyFromFile, privateKeyFromFile)
		check(err)

		// Check for certificate file
		certificateFileCheck, err := FileExists(certPaths.RootCACertsPath + "/ca.pem")
		if !certificateFileCheck {
			// Write Certificate file
			certificateFile, err := writeCertificateFile(pemEncodeCertificate(caBytes), certPaths.RootCACertsPath+"/ca.pem")
			check(err)
			if certificateFile {
				logStdOut("Created Root CA Certificate file")
			}
		}
		return true, []string{"Finished"}, nil
	}
	return false, []string{"CSR Failure"}, err
}
