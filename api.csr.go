package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"path/filepath"
)

// APICreateNewCSR handles the POST /certificate-requests endpoint
func APICreateNewCSR(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		fmt.Fprintf(w, "ParseForm() err: %v", err)
		return
	}

	certInfoRaw := r.FormValue("csr_info")
	certInfoBytes := []byte(certInfoRaw)

	certInfo := CertificateConfiguration{}
	err := json.Unmarshal(certInfoBytes, &certInfo)
	check(err)

	caName := certInfo.Subject.CommonName
	sluggedName := slugger(caName)
	logStdOut("caName " + caName)
	logStdOut("sluggedName " + sluggedName)
}

// APIListCSRs handles the GET /certificate-requests endpoint
func APIListCSRs(w http.ResponseWriter, r *http.Request) {
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
		// Check if the directory exists
		absPath, err := filepath.Abs(readConfig.Locksmith.PKIRoot + "/roots/" + parentPath)
		checkAndFail(err)
		intermedCAParentPathExists, err := DirectoryExists(absPath)

		if intermedCAParentPathExists {
			// Get listing of intermediate cas in the parent path
			intermedCAs := DirectoryListingNames(absPath + "/certreqs/")

			returnData := &RESTGETIntermedCAJSONReturn{
				Status:          "success",
				Errors:          []string{},
				Messages:        []string{"listing of CSRs"},
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
