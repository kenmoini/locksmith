package main

import (
	"crypto/rsa"
	"crypto/sha1"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/asn1"
	"log"
	"math/big"
	"time"
)

// setupCACert creates a Certificate resource
func setupCACert(serialNumber int64, commonName string, organization []string, organizationalUnit []string, country []string, province []string, locality []string, streetAddress []string, postalCode []string, addTime []int, sanData SANData, pubKey *rsa.PublicKey) *x509.Certificate {
	// set up our CA certificate
	actualURIs, err := bakeURIs(sanData.URIs)
	check(err)

	issuerBytes, err := marshalIANs(sanData.DNSNames, sanData.EmailAddresses, sanData.IPAddresses, actualURIs)
	check(err)

	issuerAltName := pkix.Extension{Id: asn1.ObjectIdentifier{2, 5, 29, 18}, Critical: false, Value: issuerBytes}
	currentTime := time.Now()
	yesterdayTime := time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day(), 0, 0, 0, 0, time.UTC).Add(-24 * time.Hour)

	publicKeyBytes, _, err := marshalPublicKey(pubKey)
	check(err)
	h := sha1.Sum(publicKeyBytes)
	subjectKeyID := h[:]

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
		NotBefore:             time.Date(yesterdayTime.Year(), yesterdayTime.Month(), yesterdayTime.Day(), 0, 0, 0, 0, yesterdayTime.Location()),
		NotAfter:              time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day(), 0, 0, 0, 0, time.UTC).AddDate(addTime[0], addTime[1], addTime[2]),
		IsCA:                  true,
		AuthorityKeyId:        subjectKeyID,
		DNSNames:              sanData.DNSNames,
		EmailAddresses:        sanData.EmailAddresses,
		IPAddresses:           sanData.IPAddresses,
		URIs:                  actualURIs,
		KeyUsage:              x509.KeyUsageCertSign | x509.KeyUsageCRLSign,
		BasicConstraintsValid: true,
		ExtraExtensions:       []pkix.Extension{issuerAltName},
		CRLDistributionPoints: []string{"https://ca.example.labs:443/crl/ca.example.labs_Root_Certification_Authority.crl"},
		IssuingCertificateURL: []string{"https://ca.example.labs:443/certs/ca.example.labs_Root_Certification_Authority.cert.pem"},
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
	if len(certConfig.Subject.OrganizationalUnit) == 0 {
		checkInputError = true
		checkInputErrors = append(checkInputErrors, "Missing OrganizationalUnit field")
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
	csrFileCheck, err := FileExists(certPaths.RootCACertRequestsPath + "/ca.pem")
	check(err)
	if !csrFileCheck {
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
		if !caCSR {
			check(err)
			return false, []string{"CSR Failure"}, err
		}
	}

	// Read in CSR lol
	caCSRPEM, err := readCSRFromFile(certPaths.RootCACertRequestsPath + "/ca.pem")
	log.Printf("%v", caCSRPEM.Subject.CommonName)

	// Check for certificate file
	certificateFileCheck, err := FileExists(certPaths.RootCACertsPath + "/ca.pem")
	if !certificateFileCheck {
		// Create Self-signed Certificate
		// Create CA Object
		rootCA := setupCACert(readSerialNumberAsInt64(rootSlug), caCSRPEM.Subject.CommonName, caCSRPEM.Subject.Organization, caCSRPEM.Subject.OrganizationalUnit, caCSRPEM.Subject.Country, caCSRPEM.Subject.Province, caCSRPEM.Subject.Locality, caCSRPEM.Subject.StreetAddress, caCSRPEM.Subject.PostalCode, certConfig.ExpirationDate, certConfig.SANData, pubKeyFromFile)

		// Byte Encode the Certificate - https://golang.org/pkg/crypto/x509/#CreateCertificate
		caBytes, err := CreateCert(rootCA, rootCA, pubKeyFromFile, privateKeyFromFile)
		check(err)

		// Write Certificate file
		certificateFile, err := writeCertificateFile(pemEncodeCertificate(caBytes), certPaths.RootCACertsPath+"/ca.pem")
		check(err)
		if certificateFile {
			logStdOut("Created Root CA Certificate file")
			return true, []string{"Root CA Created"}, nil
		}
	}
	return true, []string{"Finished"}, nil
}
