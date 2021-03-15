package main

import (
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/json"
	"fmt"
	"math/big"
	"net"
	"net/http"
	"path/filepath"
	"time"
)

// setupIntermediateCACert
func setupIntermediateCACert(serialNumber int64, organization string, organizationalUnit string, country string, province string, locality string, streetAddress string, postalCode string, addTime []int) *x509.Certificate {
	return &x509.Certificate{
		SerialNumber: big.NewInt(serialNumber),
		Subject: pkix.Name{
			Organization:       []string{organization},
			OrganizationalUnit: []string{organizationalUnit},
			Country:            []string{country},
			Province:           []string{province},
			Locality:           []string{locality},
			StreetAddress:      []string{streetAddress},
			PostalCode:         []string{postalCode},
		},
		IPAddresses:           []net.IP{net.IPv4(127, 0, 0, 1), net.IPv6loopback},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().AddDate(addTime[0], addTime[1], addTime[2]),
		IsCA:                  true,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		BasicConstraintsValid: true,
	}
}

// APIListIntermediateCAs handles the GET /intermediates endpoint
func APIListIntermediateCAs(w http.ResponseWriter, r *http.Request) {
	// Parse the submitted form data
	if err := r.ParseForm(); err != nil {
		fmt.Fprintf(w, "ParseForm() err: %v", err)
		return
	}
	// Read in the submitted JSON

	intermedCAInfoRaw := r.FormValue("ica_info")
	intermedCAInfoBytes := []byte(intermedCAInfoRaw)

	intermedCAInfo := RESTGETIntermedCAJSONIn{}
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

	// Check if the directory exists
	absPath, err := filepath.Abs(readConfig.Locksmith.PKIRoot + "/roots/" + parentPath)
	checkAndFail(err)
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
			// If the intermediate doesn't exist, check the parent signing key and see if it's password protected - decrypt if needed
			logNeworkRequestStdOut(caName+" ("+sluggedName+") creating intermediate ca", r)

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
