package locksmith

import (
	"encoding/json"
	"fmt"
	"net/http"
	"path/filepath"
)

// listIntermediateCAsAPI handles the GET /v1/intermediates endpoint
func listIntermediateCAsAPI(w http.ResponseWriter, r *http.Request) {
	var parentPath string
	var parentPathRaw string

	// Read in the submitted parameters
	queryParams := r.URL.Query()
	parentCNPath, presentCN := queryParams["parent_cn_path"]
	if presentCN {
		parentPath = splitCACNChainToPath(parentCNPath[0])
		parentPathRaw = parentCNPath[0]
	}
	parentSlugPath, presentSlug := queryParams["parent_slug_path"]
	if presentSlug {
		parentPath = splitCACNChainToPath(parentSlugPath[0])
		parentPathRaw = parentSlugPath[0]
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
		check(err)

		if intermedCAParentPathExists {
			// Get listing of intermediate cas in the parent path
			intermedCAs := DirectoryListingNames(absPath + "/intermed-ca/")

			returnData := &RESTGETIntermedCAJSONReturn{
				Status:          "success",
				Errors:          []string{},
				Messages:        []string{"Listing of Intermediate Certificate Authorities under " + parentPathRaw},
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

// createNewIntermediateCAAPI handles the POST /v1/intermediates endpoint
func createNewIntermediateCAAPI(w http.ResponseWriter, r *http.Request) {
	// Parse the JSON body into the CertificateConfiguration struct
	intermedCAInfo := RESTPOSTIntermedCAJSONIn{}
	err := json.NewDecoder(r.Body).Decode(&intermedCAInfo)
	check(err)

	var parentPath string
	if intermedCAInfo.CommonNamePath != "" {
		parentPath = splitCACNChainToPath(intermedCAInfo.CommonNamePath)
	}
	if intermedCAInfo.SlugPath != "" {
		parentPath = splitCACNChainToPath(intermedCAInfo.SlugPath)
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

	if intermedCAParentPathExists {

		caName := intermedCAInfo.CertificateConfiguration.Subject.CommonName
		sluggedName := slugger(caName)

		// If the intermediate's parent path exists, check if the intermediate ca exists before (re)creating it
		logNeworkRequestStdOut(caName+" ("+sluggedName+"): Checking "+absPath+"/intermed-ca/"+sluggedName, r)
		intermedCAPathExists, err := DirectoryExists(absPath + "/intermed-ca/" + sluggedName)
		check(err)

		if intermedCAPathExists {
			// if the intermediate exists, return with an intermed-ca-exists error
			logNeworkRequestStdOut(caName+" ("+sluggedName+") intermed-ca-exists", r)
			returnData := &ReturnGenericMessage{
				Status:   "intermed-ca-exists",
				Errors:   []string{"Intermediate CA '" + caName + "' exists!"},
				Messages: []string{}}
			returnResponse, _ := json.Marshal(returnData)
			fmt.Fprintf(w, string(returnResponse))
		} else {
			// TODO:
			// If the intermediate doesn't exist, check the parent signing key and see if it's password protected - decrypt if needed

			logNeworkRequestStdOut(caName+" ("+sluggedName+") creating intermediate ca", r)
			icaCreated, _, icaCert, err := createNewIntermediateCA(intermedCAInfo, absPath)
			check(err)
			if icaCreated {
				logNeworkRequestStdOut(caName+" ("+sluggedName+") intermed-ca-created", r)
				returnData := &ReturnPostRoots{
					Status:   "intermed-ca-created",
					Errors:   []string{},
					Messages: []string{"Successfully created Intermediate CA '" + caName + "'!"},
					Root: RootInfo{
						Slug:     sluggedName,
						CertInfo: icaCert,
						Serial:   readSerialNumberAbs(absPath + "/intermed-ca/" + sluggedName)}}
				returnResponse, _ := json.Marshal(returnData)
				fmt.Fprintf(w, string(returnResponse))
			} else {
				logNeworkRequestStdOut(caName+" ("+sluggedName+") error-creating-intermed-ca", r)
				returnData := &ReturnGenericMessage{
					Status:   "intermed-ca-creation-error",
					Errors:   []string{"Error creating Intermediate CA '" + caName + "'!"},
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
