package locksmith

import (
	"encoding/json"
	"fmt"
	"net/http"
	"path/filepath"
)

// listRootCAsAPI handles the GET /v1/roots endpoint
func listRootCAsAPI(w http.ResponseWriter, r *http.Request) {
	rootListing := DirectoryListingNames(readConfig.Locksmith.PKIRoot + "/roots/")
	returnData := &ReturnGetRoots{
		Status:   "success",
		Errors:   []string{},
		Messages: []string{},
		Roots:    rootListing}
	returnResponse, _ := json.Marshal(returnData)
	fmt.Fprintf(w, string(returnResponse))
}

// createNewRootCAAPI handles the POST /v1/roots endpoint
func createNewRootCAAPI(w http.ResponseWriter, r *http.Request) {
	// Parse the JSON body into the CertificateConfiguration struct
	certInfo := CertificateConfiguration{}
	err := json.NewDecoder(r.Body).Decode(&certInfo)
	check(err)

	caName := certInfo.Subject.CommonName
	if caName == "" {
		returnData := &ReturnPostRoots{
			Status:   "root-creation-error",
			Errors:   []string{"Invalid JSON!"},
			Messages: []string{},
			Root:     RootInfo{}}
		returnResponse, _ := json.Marshal(returnData)
		fmt.Fprintf(w, string(returnResponse))
	}
	sluggedName := slugger(caName)

	// Find absolute path
	checkForRootPath, err := filepath.Abs(readConfig.Locksmith.PKIRoot + "/roots/" + sluggedName)
	check(err)

	// Check if the directory exists
	caRootPathExists, err := DirectoryExists(checkForRootPath)
	check(err)

	if caRootPathExists {
		// If the root path exists, don't regenerate a CA
		logNeworkRequestStdOut(caName+" ("+sluggedName+") root-exists", r)
		returnData := &ReturnPostRoots{
			Status:   "root-exists",
			Errors:   []string{"Root CA " + caName + " already exists!"},
			Messages: []string{},
			Root: RootInfo{
				Slug:   sluggedName,
				Serial: readSerialNumber(sluggedName)}}
		returnResponse, _ := json.Marshal(returnData)
		fmt.Fprintf(w, string(returnResponse))
	} else {
		// Generate a new Certificate Authority
		newCAState, newCA, caCert, err := createNewCA(certInfo)
		check(err)

		if newCAState {

			logNeworkRequestStdOut(caName+" ("+sluggedName+") root-created", r)
			returnData := &ReturnPostRoots{
				Status:   "root-created",
				Errors:   []string{},
				Messages: []string{"Root CA " + caName + " created!"},
				Root: RootInfo{
					Slug:     sluggedName,
					CertInfo: caCert,
					Serial:   readSerialNumber(sluggedName)}}
			returnResponse, _ := json.Marshal(returnData)
			fmt.Fprintf(w, string(returnResponse))

		} else {

			logNeworkRequestStdOut(caName+" ("+sluggedName+") root-creation-error", r)
			returnData := &ReturnPostRoots{
				Status:   "root-creation-error",
				Errors:   []string{err.Error()},
				Messages: newCA,
				Root: RootInfo{
					Slug:   sluggedName,
					Serial: readSerialNumber(sluggedName)}}
			returnResponse, _ := json.Marshal(returnData)
			fmt.Fprintf(w, string(returnResponse))

		}
	}
}
