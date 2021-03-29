package locksmith

import (
	"bytes"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/asn1"
	"fmt"
	"math/big"
	"time"
)

// setupIntermediateCACert
func setupIntermediateCACert(serialNumber int64, commonName string, organization []string, organizationalUnit []string, country []string, province []string, locality []string, streetAddress []string, postalCode []string, addTime []int, sanData SANData, pubKey *rsa.PublicKey) *x509.Certificate {
	// set up our Intermediate CA certificate

	// Convert string slice of URLs into actual URI objects
	actualURIs, err := bakeURIs(sanData.URIs)
	check(err)

	// Take the SAN data and format for IAN
	issuerBytes, err := marshalIANs(sanData.DNSNames, sanData.EmailAddresses, sanData.IPAddresses, actualURIs)
	check(err)
	issuerAltName := pkix.Extension{Id: asn1.ObjectIdentifier{2, 5, 29, 18}, Critical: false, Value: issuerBytes}

	// Set time for UTC format
	currentTime := time.Now()
	yesterdayTime := time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day(), 0, 0, 0, 0, time.UTC).Add(-24 * time.Hour)

	// Set up SubjectKeyID from Public Key
	publicKeyBytes, _, err := marshalPublicKey(pubKey)
	check(err)
	h := sha1.Sum(publicKeyBytes)
	subjectKeyID := h[:]

	return &x509.Certificate{
		SignatureAlgorithm: x509.SHA512WithRSA,
		SerialNumber:       big.NewInt(serialNumber),
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

// createNewIntermediateCA - creates a new Intermediate Certificate Authority
func createNewIntermediateCA(configWrapper RESTPOSTIntermedCAJSONIn, parentPath string) (bool, []string, x509.Certificate, error) {
	// Define needed variables
	var rootSlug string
	var caName string

	// Check if the certificate configuration is valid
	certificateValid, validationMsgs, err := ValidateCertificateConfiguration(configWrapper.CertificateConfiguration)
	check(err)

	if !certificateValid {
		return false, validationMsgs, x509.Certificate{}, Stoerr("cert-config-error")
	}

	caName = configWrapper.CertificateConfiguration.Subject.CommonName
	rootSlug = slugger(caName)

	rootSlugPath := parentPath + "/intermed-ca/" + rootSlug
	rsaPrivateKeyPassword := configWrapper.CertificateConfiguration.RSAPrivateKeyPassphrase
	signingCARSAPrivateKeyPassword := configWrapper.SigningPrivateKeyPassphrase

	// Create the Intermediate CA base directories and files
	certPaths := setupCAFileStructure(rootSlugPath)

	// Check for Intermediate CA key pair
	caKeyCheck, err := FileExists(certPaths.RootCAKeysPath + "/ca.priv.pem")
	check(err)

	if !caKeyCheck {
		// if there is no private key, create one
		rootPrivKey, rootPubKey, err := GenerateRSAKeypair(4096)
		check(err)

		pemEncodedPrivateKey, encryptedPrivateKeyBytes := pemEncodeRSAPrivateKey(rootPrivKey, rsaPrivateKeyPassword)

		if rsaPrivateKeyPassword == "" {
			rootPrivKeyFile, rootPubKeyFile, err := writeRSAKeyPair(pemEncodedPrivateKey, pemEncodeRSAPublicKey(rootPubKey), certPaths.RootCAKeysPath+"/ca")
			check(err)
			if !rootPrivKeyFile || !rootPubKeyFile {
				return false, []string{"Root CA Private Key Failure"}, x509.Certificate{}, err
			}
		} else {

			encStr := B64EncodeBytesToStr(encryptedPrivateKeyBytes.Bytes())
			encBufferB := bytes.NewBufferString(encStr)

			rootPrivKeyFile, rootPubKeyFile, err := writeRSAKeyPair(encBufferB, pemEncodeRSAPublicKey(rootPubKey), certPaths.RootCAKeysPath+"/ca")
			check(err)
			if !rootPrivKeyFile || !rootPubKeyFile {
				return false, []string{"Root CA Private Key Failure"}, x509.Certificate{}, err
			}
		}

	}

	// Read in the Private key
	privateKeyFromFile := GetPrivateKey(certPaths.RootCAKeysPath+"/ca.priv.pem", rsaPrivateKeyPassword)

	// Read in the Public key
	pubKeyFromFile := GetPublicKey(certPaths.RootCAKeysPath + "/ca.pub.pem")

	// Create Self-signed Certificate Request
	csrFileCheck, err := FileExists(certPaths.RootCACertRequestsPath + "/ca.pem")
	check(err)
	if !csrFileCheck {
		caCSR, err := generateCSR(certPaths.RootCACertRequestsPath+"/ca.pem",
			privateKeyFromFile,
			configWrapper.CertificateConfiguration.Subject.CommonName,
			configWrapper.CertificateConfiguration.Subject.Organization,
			configWrapper.CertificateConfiguration.Subject.OrganizationalUnit,
			configWrapper.CertificateConfiguration.Subject.Country,
			configWrapper.CertificateConfiguration.Subject.Province,
			configWrapper.CertificateConfiguration.Subject.Locality,
			configWrapper.CertificateConfiguration.Subject.StreetAddress,
			configWrapper.CertificateConfiguration.Subject.PostalCode,
			configWrapper.CertificateConfiguration.SANData,
			true)
		if !caCSR {
			check(err)
			return false, []string{"Intermediate CA CSR Failure"}, x509.Certificate{}, err
		}
	}

	// Read in CSR lol
	caCSRPEM, err := readCSRFromFile(certPaths.RootCACertRequestsPath + "/ca.pem")
	check(err)
	//log.Printf("Created CSR with CN: %v", caCSRPEM.Subject.CommonName)

	// Copy Intermediate CA Certificate Request File to the Signing CA's certreqs folder
	copyCSRErr := CopyFile(certPaths.RootCACertRequestsPath+"/ca.pem", parentPath+"/certreqs/"+slugger(caCSRPEM.Subject.CommonName)+".pem", 4096)
	check(copyCSRErr)

	// Check for certificate file
	certificateFileCheck, err := FileExists(certPaths.RootCACertsPath + "/ca.pem")
	check(err)
	if !certificateFileCheck {
		// Create Parent Signed Certificate
		// Create Intermediate CA Object
		// Serial number should come from the signing CA's serial
		intermedCA := setupIntermediateCACert(readSerialNumberAsInt64Abs(parentPath+"/ca.serial"), caCSRPEM.Subject.CommonName, caCSRPEM.Subject.Organization, caCSRPEM.Subject.OrganizationalUnit, caCSRPEM.Subject.Country, caCSRPEM.Subject.Province, caCSRPEM.Subject.Locality, caCSRPEM.Subject.StreetAddress, caCSRPEM.Subject.PostalCode, configWrapper.CertificateConfiguration.ExpirationDate, configWrapper.CertificateConfiguration.SANData, pubKeyFromFile)

		// Read in the Signing CA
		rootCA, err := ReadCACertificate(parentPath)
		check(err)

		rootCAPrivateKeyFromFile := GetPrivateKey(parentPath+"/private/ca.priv.pem", signingCARSAPrivateKeyPassword)

		// Byte Encode the Certificate - https://golang.org/pkg/crypto/x509/#CreateCertificate
		caBytes, err := CreateCert(intermedCA, rootCA, pubKeyFromFile, rootCAPrivateKeyFromFile)
		check(err)

		// Write Certificate file
		certificateFile, err := writeCertificateFile(pemEncodeCertificate(caBytes), certPaths.RootCACertsPath+"/ca.pem")
		check(err)
		if !certificateFile {
			return false, []string{"Intermediate CA Certificate Creation Failure!"}, x509.Certificate{}, err
		}

		// Increase the serial number in the Intermediate CA Serial file
		increaseSerial, err := IncreaseSerialNumberAbs(certPaths.RootCACertSerialFilePath)
		check(err)
		if !increaseSerial {
			logStdOut("Serial Increment ERROR!")
			return false, []string{"Intermediate CA Serial Increment Error"}, x509.Certificate{}, err
		}
		// Increase the Signing CA serial number
		increaseSerial, err = IncreaseSerialNumberAbs(parentPath + "/ca.serial")
		check(err)
		if !increaseSerial {
			logStdOut("Serial Increment ERROR!")
			return false, []string{"Signing CA Serial Increment Error"}, x509.Certificate{}, err
		}
	}

	// Read in Certificate File lol
	caCert, err := ReadCertFromFile(certPaths.RootCACertsPath + "/ca.pem")
	check(err)

	// Copy Intermediate CA Certificate File to the Signing CA's certs folder
	copyCertErr := CopyFile(certPaths.RootCACertsPath+"/ca.pem", parentPath+"/certs/"+slugger(caCert.Subject.CommonName)+".pem", 4096)
	check(copyCertErr)

	// Copy Intermediate CA Certificate File to the Signing CA's newcerts folder
	serialNumber := fmt.Sprintf("%02d", caCert.SerialNumber)
	copyNewCertsErr := CopyFile(certPaths.RootCACertsPath+"/ca.pem", parentPath+"/newcerts/"+serialNumber+".pem", 4096)
	check(copyNewCertsErr)

	// Add Certificate to Signing CA Index
	addedEntry, err := AddEntryToCAIndex(parentPath+"/ca.index", certPaths.RootCACertsPath+"/ca.pem")
	check(err)
	if !addedEntry {
		logStdOut("Signing CA Index ERROR!")
		return false, []string{"Signing CA Index Entry Error"}, x509.Certificate{}, err
	}

	// Create CRL with CA Cert
	caCRL, err := CreateNewCRLForCA(caCert, privateKeyFromFile, certPaths.RootCACertRevListPath+"/ca.crl")
	check(err)
	if !caCRL {
		logStdOut("Intermediate CA CRL ERROR!")
		return false, []string{"Intermediate CA CRL Creation Error"}, x509.Certificate{}, err
	}

	return true, []string{"Finished creating Intermediate CA: " + caCert.Subject.CommonName}, *caCert, nil

}
