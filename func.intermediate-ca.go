package main

import (
	"crypto/rsa"
	"crypto/sha1"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/asn1"
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"net/http"
	"path/filepath"
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
		SerialNumber: big.NewInt(serialNumber),
		Subject: pkix.Name{
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
func createNewIntermediateCA(configWrapper RESTPOSTIntermedCAJSONIn, parentPath string) (bool, []string, error) {
	checkInputError := false
	var checkInputErrors []string
	var rootSlug string
	var caName string

	logStdOut("Parent Path: " + parentPath)

	if configWrapper.CertificateConfiguration.Subject.CommonName == "" {
		checkInputError = true
		checkInputErrors = append(checkInputErrors, "Missing common name field")
	} else {
		caName = configWrapper.CertificateConfiguration.Subject.CommonName
		rootSlug = slugger(caName)
	}
	if len(configWrapper.CertificateConfiguration.Subject.Organization) == 0 {
		checkInputError = true
		checkInputErrors = append(checkInputErrors, "Missing Organization field")
	}
	if len(configWrapper.CertificateConfiguration.Subject.OrganizationalUnit) == 0 {
		checkInputError = true
		checkInputErrors = append(checkInputErrors, "Missing OrganizationalUnit field")
	}
	if len(configWrapper.CertificateConfiguration.ExpirationDate) != 3 {
		checkInputError = true
		checkInputErrors = append(checkInputErrors, "Missing Expiration Date field")
	}
	if checkInputError {
		return false, checkInputErrors, Stoerr("cert-config-error")
	}

	rootSlugPath := parentPath + "/intermed-ca/" + rootSlug
	rsaPrivateKeyPassword := configWrapper.CertificateConfiguration.RSAPrivateKeyPassphrase
	signingCARSAPrivateKeyPassword := configWrapper.SigningPrivateKeyPassphrase

	// Create the Intermediate CA base directories and files
	certPaths := setupCAFileStructure(rootSlugPath)

	// Check for Intermediate CA key pair
	caKeyCheck, err := FileExists(certPaths.RootCACertKeysPath + "/ca.priv.pem")
	check(err)

	if !caKeyCheck {
		// if there is no private key, create one
		rootPrivKey, rootPubKey, err := generateRSAKeypair(4096)
		check(err)

		rootPrivKeyFile, rootPubKeyFile, err := writeRSAKeyPair(pemEncodeRSAPrivateKey(rootPrivKey, ""), pemEncodeRSAPublicKey(rootPubKey), certPaths.RootCACertKeysPath+"/ca")
		check(err)
		if rootPrivKeyFile && rootPubKeyFile {
			logStdOut("Intermediate CA RSA Key Pair Created")
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
			configWrapper.CertificateConfiguration.Subject.CommonName,
			configWrapper.CertificateConfiguration.Subject.Organization,
			configWrapper.CertificateConfiguration.Subject.OrganizationalUnit,
			configWrapper.CertificateConfiguration.Subject.Country,
			configWrapper.CertificateConfiguration.Subject.Province,
			configWrapper.CertificateConfiguration.Subject.Locality,
			configWrapper.CertificateConfiguration.Subject.StreetAddress,
			configWrapper.CertificateConfiguration.Subject.PostalCode,
			true)
		if !caCSR {
			check(err)
			return false, []string{"Intermediate CA CSR Failure"}, err
		}
	}

	// Read in CSR lol
	caCSRPEM, err := readCSRFromFile(certPaths.RootCACertRequestsPath + "/ca.pem")
	log.Printf("Created CSR with CN: %v", caCSRPEM.Subject.CommonName)

	// Check for certificate file
	certificateFileCheck, err := FileExists(certPaths.RootCACertsPath + "/ca.pem")
	if !certificateFileCheck {
		// Create Parent Signed Certificate
		// Create Intermediate CA Object
		intermedCA := setupIntermediateCACert(readSerialNumberAsInt64Abs(certPaths.RootCACertSerialFilePath), caCSRPEM.Subject.CommonName, caCSRPEM.Subject.Organization, caCSRPEM.Subject.OrganizationalUnit, caCSRPEM.Subject.Country, caCSRPEM.Subject.Province, caCSRPEM.Subject.Locality, caCSRPEM.Subject.StreetAddress, caCSRPEM.Subject.PostalCode, configWrapper.CertificateConfiguration.ExpirationDate, configWrapper.CertificateConfiguration.SANData, pubKeyFromFile)

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
			return false, []string{"Intermediate CA Certificate Creation Failure!"}, err
		}
	}

	// Read in Certificate File lol
	caCert, err := ReadCertFromFile(certPaths.RootCACertsPath + "/ca.pem")
	check(err)

	// Create CRL with CA Cert
	caCRL, err := CreateNewCRLForCA(caCert, privateKeyFromFile, certPaths.RootCACertRevListPath+"/ca.crl")
	if !caCRL {
		logStdOut("Intermediate CA CRL ERROR!")
		return false, []string{"Intermediate CA CRL Creation Error"}, err
	}

	return true, []string{"Finished creating Intermediate CA: " + caCert.Subject.CommonName}, nil

}

// APIListIntermediateCAs handles the GET /intermediates endpoint
func APIListIntermediateCAs(w http.ResponseWriter, r *http.Request) {
	var parentPath string

	// Parse the submitted form data
	if err := r.ParseForm(); err != nil {
		fmt.Fprintf(w, "ParseForm() err: %v", err)
		return
	}
	// Read in the submitted JSON
	queryParams := r.URL.Query()
	parentCNPath, presentCN := queryParams["parent_cn_path"]
	if presentCN {
		parentPath = splitCommonNamesToPath(parentCNPath[0])
	}
	parentSlugPath, presentSlug := queryParams["parent_slug_path"]
	if presentSlug {
		parentPath = splitSlugToPath(parentSlugPath[0])
	}

	// Neither options are submitted - error
	if parentPath == "" {
		returnData := &ReturnGenericMessage{
			Status:   "missing-parent-path",
			Errors:   []string{"Missing parent path!  Must supply either `parent_cn_path` or `parent_slug_path`"},
			Messages: []string{}}
		returnResponse, _ := json.Marshal(returnData)
		fmt.Fprintf(w, string(returnResponse))
	} else {
		logStdOut("parentPath " + parentPath)

		// Check if the directory exists
		absPath, err := filepath.Abs(readConfig.Locksmith.PKIRoot + "/roots/" + parentPath)
		checkAndFail(err)
		logStdOut("absPath: " + absPath)
		intermedCAParentPathExists, err := DirectoryExists(absPath)

		if intermedCAParentPathExists {
			// Get listing of intermediate cas in the parent path
			intermedCAs := DirectoryListingNames(absPath + "/intermed-ca/")

			returnData := &RESTGETIntermedCAJSONReturn{
				Status:          "success",
				Errors:          []string{},
				Messages:        []string{"listing of intermed cas"},
				IntermediateCAs: intermedCAs}
			returnResponse, _ := json.Marshal(returnData)
			fmt.Fprintf(w, string(returnResponse))
		} else {
			// Parent path does not exist, return invalid-parent-path
			returnData := &ReturnGenericMessage{
				Status:   "invalid-parent-path",
				Errors:   []string{"Invalid parent path, no chain exists!"},
				Messages: []string{}}
			returnResponse, _ := json.Marshal(returnData)
			fmt.Fprintf(w, string(returnResponse))
		}
	}
}

// APICreateNewIntermediateCA handles the POST /intermediates endpoint
func APICreateNewIntermediateCA(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		fmt.Fprintf(w, "ParseForm() err: %v", err)
		return
	}

	intermedCAInfoRaw := r.FormValue("ica_info")
	intermedCAInfoBytes := []byte(intermedCAInfoRaw)

	intermedCAInfo := RESTPOSTIntermedCAJSONIn{}
	err := json.Unmarshal(intermedCAInfoBytes, &intermedCAInfo)
	check(err)

	var parentPath string
	if intermedCAInfo.CommonNamePath != "" {
		parentPath = splitCommonNamesToPath(intermedCAInfo.CommonNamePath)
	}
	if intermedCAInfo.SlugPath != "" {
		parentPath = splitSlugToPath(intermedCAInfo.SlugPath)
	}
	// Neither options are submitted - error
	if parentPath == "" {
		returnData := &ReturnGenericMessage{
			Status:   "missing-parent-path",
			Errors:   []string{"Missing parent path!  Must supply either `parent_cn_path` or `parent_slug_path`"},
			Messages: []string{}}
		returnResponse, _ := json.Marshal(returnData)
		fmt.Fprintf(w, string(returnResponse))
	} else {
		logStdOut("parentPath " + parentPath)
	}

	// Check if the parent path directory exists
	absPath, err := filepath.Abs(readConfig.Locksmith.PKIRoot + "/roots/" + parentPath)
	checkAndFail(err)
	intermedCAParentPathExists, err := DirectoryExists(absPath)
	check(err)

	caName := intermedCAInfo.CertificateConfiguration.Subject.CommonName
	sluggedName := slugger(caName)

	if intermedCAParentPathExists {
		// If the intermediate's parent path exists, check if the intermediate ca exists before (re)creating it
		logNeworkRequestStdOut(caName+" ("+sluggedName+"): Checking "+absPath+"/intermed-ca/"+sluggedName, r)
		intermedCAPathExists, err := DirectoryExists(absPath + "/intermed-ca/" + sluggedName)
		check(err)

		if intermedCAPathExists {
			// if the intermediate exists, return with an intermed-ca-exists error
			logNeworkRequestStdOut(caName+" ("+sluggedName+") intermed-ca-exists", r)
			returnData := &ReturnGenericMessage{
				Status:   "intermed-ca-exists",
				Errors:   []string{"Intermediate CA " + caName + " exists!"},
				Messages: []string{}}
			returnResponse, _ := json.Marshal(returnData)
			fmt.Fprintf(w, string(returnResponse))
		} else {
			// TODO:
			// If the intermediate doesn't exist, check the parent signing key and see if it's password protected - decrypt if needed

			logNeworkRequestStdOut(caName+" ("+sluggedName+") creating intermediate ca", r)
			icaCreated, _, err := createNewIntermediateCA(intermedCAInfo, absPath)
			check(err)
			if icaCreated {
				logNeworkRequestStdOut(caName+" ("+sluggedName+") intermed-ca-created", r)
				returnData := &ReturnGenericMessage{
					Status:   "intermed-ca-created",
					Errors:   []string{},
					Messages: []string{"Successfully created Intermediate CA " + caName + "!"}}
				returnResponse, _ := json.Marshal(returnData)
				fmt.Fprintf(w, string(returnResponse))
			} else {
				logNeworkRequestStdOut(caName+" ("+sluggedName+") error-creating-intermed-ca", r)
				returnData := &ReturnGenericMessage{
					Status:   "error-creating-intermed-ca",
					Errors:   []string{"Error creating Intermediate CA " + caName + "!"},
					Messages: []string{}}
				returnResponse, _ := json.Marshal(returnData)
				fmt.Fprintf(w, string(returnResponse))
			}
		}
	} else {
		// Parent path does not exist, return invalid-parent-path
		returnData := &ReturnGenericMessage{
			Status:   "invalid-parent-path",
			Errors:   []string{"Invalid parent path, no chain exists!"},
			Messages: []string{}}
		returnResponse, _ := json.Marshal(returnData)
		fmt.Fprintf(w, string(returnResponse))
	}
}
