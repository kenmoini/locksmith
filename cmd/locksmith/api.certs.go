package locksmith

import (
	"encoding/json"
	"fmt"
	"net/http"
	"path/filepath"
)

// createNewCertAPI handles the POST /v1/certificates endpoint
func createNewCertAPI(w http.ResponseWriter, r *http.Request) {
	// Load in POST JSON Data
	csrInfo := RESTPOSTCertificateJSONIn{}
	err := json.NewDecoder(r.Body).Decode(&csrInfo)
	check(err)

	// Set up Parent Path
	var parentPath string
	if csrInfo.CommonNamePath != "" {
		parentPath = splitCACNChainToPath(csrInfo.CommonNamePath)
	}
	if csrInfo.SlugPath != "" {
		parentPath = splitCACNChainToPath(csrInfo.SlugPath)
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

}

// listCertsAPI handles the GET /v1/certificates endpoint
func listCertsAPI(w http.ResponseWriter, r *http.Request) {
	var parentPath string
	var parentPathRaw string

	// Read in the submitted GET URL parameters
	queryParams := r.URL.Query()
	parentCNPath, presentCN := queryParams["parent_cn_path"]
	parentSlugPath, presentSlug := queryParams["parent_slug_path"]
	if presentCN {
		parentPath = splitCACNChainToPath(parentCNPath[0])
		parentPathRaw = parentCNPath[0]
	}
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

		caParentPathExists, err := DirectoryExists(absPath)
		check(err)

		if caParentPathExists {
			// Get listing of CSRs in the parent path
			certificates := DirectoryListingNamesNoExt(absPath + "/certs/")

			certificates = rmStrFromStrSlice("ca", certificates)

			if len(certificates) > 0 {
				// Got some hits
				returnData := &RESTGETCertificatesJSONReturn{
					Status:       "success",
					Errors:       []string{},
					Messages:     []string{"Listing of Certificates for CA Path '" + parentPathRaw + "'"},
					Certificates: certificates}
				returnResponse, _ := json.Marshal(returnData)
				fmt.Fprintf(w, string(returnResponse))
			} else {
				// No Certificates in this path
				returnData := &ReturnGenericMessage{
					Status:   "no-certs-found",
					Errors:   []string{"No Certificates found in this CA Path '" + parentPathRaw + "'!"},
					Messages: []string{}}
				returnResponse, _ := json.Marshal(returnData)
				fmt.Fprintf(w, string(returnResponse))
			}
		} else {
			// Parent path does not exist, return invalid-parent-path
			returnData := &ReturnGenericMessage{
				Status:   "invalid-parent-path",
				Errors:   []string{"Invalid parent path, no such chain exists!"},
				Messages: []string{}}
			returnResponse, _ := json.Marshal(returnData)
			fmt.Fprintf(w, string(returnResponse))
		}
	}
}
