package locksmith

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/asn1"
	"encoding/pem"
)

// generateCSR takes the full lifecycle of generating and saving a CSR
func generateCSR(path string, signingKey interface{}, commonName string, organization []string, organizationalUnit []string, country []string, province []string, locality []string, streetAddress []string, postalCode []string, sanData SANData, isCA bool) (bool, error) {
	// Generate PKIX Name object
	csrSubjectName := setupCSRSubjectName(commonName, organization, organizationalUnit, country, province, locality, streetAddress, postalCode)

	// Setup CSR Template Object
	csrTemplate := setupCSR(csrSubjectName, isCA, sanData)

	// Create CSR object
	csr, err := createCSR(csrTemplate, signingKey)
	check(err)

	// Encode CSR to PEM format
	csrPEM := pemEncodeCSR(csr)

	// Write PEM to a file
	pemWriter, pemErr := writePEMFile(csrPEM, path)

	return pemWriter, pemErr
}

// setupCSRSubjectName just wraps the pkix.Name type for CSRs
func setupCSRSubjectName(commonName string, organization []string, organizationalUnit []string, country []string, province []string, locality []string, streetAddress []string, postalCode []string) pkix.Name {
	return pkix.Name{
		CommonName:         commonName,
		Organization:       organization,
		OrganizationalUnit: organizationalUnit,
		Country:            country,
		Province:           province,
		Locality:           locality,
		StreetAddress:      streetAddress,
		PostalCode:         postalCode,
	}
}

// setupCSR creates configuration information and returns a CSR Template
//func setupCSR(commonName string, organization []string, organizationalUnit []string, country []string, province []string, locality []string, streetAddress []string, postalCode []string, isCA bool) *x509.CertificateRequest {
func setupCSR(names pkix.Name, isCA bool, sanData SANData) *x509.CertificateRequest {

	// Convert string slice of URLs into actual URI objects
	actualURIs, err := bakeURIs(sanData.URIs)
	check(err)

	if isCA {
		val, err := asn1.Marshal(basicConstraints{true, 0})
		check(err)

		return &x509.CertificateRequest{
			Subject:            names,
			SignatureAlgorithm: x509.SHA512WithRSA,
			DNSNames:           sanData.DNSNames,
			EmailAddresses:     sanData.EmailAddresses,
			IPAddresses:        sanData.IPAddresses,
			URIs:               actualURIs,
			ExtraExtensions: []pkix.Extension{
				{
					// This identifies that the CSR is a CA
					Id:       asn1.ObjectIdentifier{2, 5, 29, 19},
					Value:    val,
					Critical: true,
				},
			},
		}
	}
	return &x509.CertificateRequest{
		Subject:            names,
		SignatureAlgorithm: x509.SHA512WithRSA,
		DNSNames:           sanData.DNSNames,
		EmailAddresses:     sanData.EmailAddresses,
		IPAddresses:        sanData.IPAddresses,
		URIs:               actualURIs,
	}
}

// createCSR is a wrapper for x509.CreateCertificateRequest
// template is a CSR template, priv is the CSR requester private key
func createCSR(template *x509.CertificateRequest, priv interface{}) ([]byte, error) {
	return x509.CreateCertificateRequest(rand.Reader, template, priv)
}

// pemEncodeCSR encodes a CreateCertificateRequest DER byte stream to a PEM
func pemEncodeCSR(certByte []byte) *bytes.Buffer {
	pemRet := new(bytes.Buffer)
	pem.Encode(pemRet, &pem.Block{
		Type:  "CERTIFICATE REQUEST",
		Bytes: certByte,
	})
	return pemRet
}

// readCSR converts a CSR byte stream into x509 CSR
func readCSR(asn1Data []byte) (*x509.CertificateRequest, error) {
	return x509.ParseCertificateRequest(asn1Data)
}

// readCSRFromFile wraps the functions needed to read and decode a CSR PEM
func readCSRFromFile(path string) (*x509.CertificateRequest, error) {
	// Read in file
	file, err := readPEMFile(path, "CERTIFICATE REQUEST")
	check(err)

	// Convert to x509.CSR
	return readCSR(file.Bytes)
}

// createNewCertificateRequest - creates a new Certificate Request from the API
func createNewCertificateRequest(config RESTPOSTCertificateRequestJSONIn, parentPath string) (bool, []string, *x509.CertificateRequest, RealKeyPair, error) {
	// Define needed variables
	var csrIsCA bool
	var csrCommonName string
	var csrCommonNameSlug string
	var privateKey *rsa.PrivateKey
	var publicKey *rsa.PublicKey

	// Check if the certificate configuration is valid
	certificateValid, validationMsgs, err := ValidateCertificateConfiguration(config.CertificateConfiguration)
	check(err)

	if !certificateValid {
		return false, validationMsgs, &x509.CertificateRequest{}, RealKeyPair{}, Stoerr("certificate-request-config-error")
	}

	// Define the CommonName and Slug
	csrCommonName = config.CertificateConfiguration.Subject.CommonName
	csrCommonNameSlug = slugger(csrCommonName)

	// Check to see if the csr exists, just to be safe
	absPathCSRFile := parentPath + "/certreqs/" + csrCommonNameSlug + ".req.pem"
	sluggedCSRFileExists, err := FileExists(absPathCSRFile)
	check(err)

	if sluggedCSRFileExists {
		return false, []string{"CSR exists!"}, &x509.CertificateRequest{}, RealKeyPair{}, Stoerr("certificate-request-exists")
	}

	// Check to see if there's a base64 encoded RSAPrivateKey defined
	if config.CertificateConfiguration.RSAPrivateKey != "" {
		// Check to see if we can base64 decode the Priv Key string
		//decodedPrivKey, err := b64.StdEncoding.DecodeString(config.CertificateConfiguration.RSAPrivateKey)
		//check(err)

	} else {
		// No incoming RSAPrivateKey, generate a new RSA Key Pair to generate the CSR
		// Check for CSR key pair
		privKeyCheck, err := FileExists(parentPath + "/keys/" + csrCommonNameSlug + ".priv.pem")
		check(err)

		if !privKeyCheck {
			// if there is no private key, create one
			csrPrivKey, csrPubKey, err := GenerateRSAKeypair(4096)
			check(err)

			// Save the Private Key to the file system
			var csrPrivKeyFile bool
			var csrPubKeyFile bool

			pemEncodedPrivateKey, encryptedPrivateKeyBytes := pemEncodeRSAPrivateKey(csrPrivKey, config.CertificateConfiguration.RSAPrivateKeyPassphrase)

			if config.CertificateConfiguration.RSAPrivateKeyPassphrase == "" {
				csrPrivKeyFile, csrPubKeyFile, err = writeRSAKeyPair(pemEncodedPrivateKey, pemEncodeRSAPublicKey(csrPubKey), parentPath+"/keys/"+csrCommonNameSlug)
				check(err)
				if !csrPrivKeyFile || !csrPubKeyFile {
					return false, []string{"CSR Key Pair Failure epkbNil"}, &x509.CertificateRequest{}, RealKeyPair{}, err
				}
			} else {

				encStr := B64EncodeBytesToStr(encryptedPrivateKeyBytes.Bytes())
				encBufferB := bytes.NewBufferString(encStr)

				csrPrivKeyFile, csrPubKeyFile, err = writeRSAKeyPair(encBufferB, pemEncodeRSAPublicKey(csrPubKey), parentPath+"/keys/"+csrCommonNameSlug)
				check(err)
				if !csrPrivKeyFile || !csrPubKeyFile {
					return false, []string{"CSR Key Pair Failure"}, &x509.CertificateRequest{}, RealKeyPair{}, err
				}
			}

		}

		// Read in the Private key
		privateKey = GetPrivateKey(parentPath+"/keys/"+csrCommonNameSlug+".priv.pem", config.CertificateConfiguration.RSAPrivateKeyPassphrase)

		// Read in the Public key
		publicKey = GetPublicKey(parentPath + "/keys/" + csrCommonNameSlug + ".pub.pem")
	}

	if config.CertificateConfiguration.CertificateType == "authority" || config.CertificateConfiguration.CertificateType == "authority-no-subs" {
		csrIsCA = true
	}

	// Generate the CSR
	csr, err := generateCSR(absPathCSRFile,
		privateKey,
		config.CertificateConfiguration.Subject.CommonName,
		config.CertificateConfiguration.Subject.Organization,
		config.CertificateConfiguration.Subject.OrganizationalUnit,
		config.CertificateConfiguration.Subject.Country,
		config.CertificateConfiguration.Subject.Province,
		config.CertificateConfiguration.Subject.Locality,
		config.CertificateConfiguration.Subject.StreetAddress,
		config.CertificateConfiguration.Subject.PostalCode,
		config.CertificateConfiguration.SANData,
		csrIsCA)
	check(err)

	if !csr {
		// Generation failure occurred
		return false, []string{"CSR Generation Failure"}, &x509.CertificateRequest{}, RealKeyPair{}, err
	}

	// Read in the CSR PEM file lol
	caCSRPEM, err := readCSRFromFile(absPathCSRFile)
	check(err)

	return true, []string{"Successfully generated CSR "}, caCSRPEM, RealKeyPair{PrivateKey: privateKey, PublicKey: publicKey}, err

}
