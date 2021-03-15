package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"path/filepath"
)

// APIListRootCAs handles the GET /roots endpoint
func APIListRootCAs(w http.ResponseWriter, r *http.Request) {
	rootListing := DirectoryListingNames(readConfig.Locksmith.PKIRoot + "/roots/")
	returnData := &ReturnGetRoots{
		Status:   "success",
		Errors:   []string{},
		Messages: []string{},
		Roots:    rootListing}
	returnResponse, _ := json.Marshal(returnData)
	fmt.Fprintf(w, string(returnResponse))
}

// APICreateNewRootCA handles the POST /roots endpoint
func APICreateNewRootCA(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		fmt.Fprintf(w, "ParseForm() err: %v", err)
		return
	}

	certInfoRaw := r.FormValue("cert_info")
	certInfoBytes := []byte(certInfoRaw)

	certInfo := CertificateConfiguration{}
	err := json.Unmarshal(certInfoBytes, &certInfo)
	check(err)

	caName := certInfo.Subject.CommonName
	sluggedName := slugger(caName)
	logStdOut("caName " + caName)
	logStdOut("sluggedName " + sluggedName)

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
			Errors:   []string{},
			Messages: []string{},
			Root: RootInfo{
				Slug:   sluggedName,
				Serial: readSerialNumber(sluggedName)}}
		returnResponse, _ := json.Marshal(returnData)
		fmt.Fprintf(w, string(returnResponse))
	} else {
		// Generate a new Certificate Authority
		newCAState, newCA, err := createNewCA(certInfo)
		check(err)
		returnData := &ReturnPostRoots{}
		if newCAState {
			logNeworkRequestStdOut(caName+" ("+sluggedName+") root-created", r)
			returnData = &ReturnPostRoots{
				Status:   "root-created",
				Errors:   []string{},
				Messages: []string{},
				Root: RootInfo{
					Slug:   sluggedName,
					Serial: readSerialNumber(sluggedName)}}
		} else {
			logNeworkRequestStdOut(caName+" ("+sluggedName+") root-creation-error", r)
			returnData = &ReturnPostRoots{
				Status:   "root-creation-error",
				Errors:   newCA,
				Messages: []string{},
				Root: RootInfo{
					Slug:   sluggedName,
					Serial: readSerialNumber(sluggedName)}}
		}
		returnResponse, _ := json.Marshal(returnData)
		fmt.Fprintf(w, string(returnResponse))
	}
}
