package locksmith

import (
	"encoding/json"
	"fmt"
	"net/http"
	"path/filepath"
)

// createNewCSRAPI handles the POST /v1/certificate-requests endpoint
func createNewCSRAPI(w http.ResponseWriter, r *http.Request) {
	// Load in POST JSON Data
	csrInfo := RESTPOSTCertificateRequestJSONIn{}
	err := json.NewDecoder(r.Body).Decode(&csrInfo)
	check(err)

	// Set up Parent Path
	var parentPath string
	if csrInfo.CommonNamePath != "" {
		parentPath = splitCommonNamesToPath(csrInfo.CommonNamePath)
	}
	if csrInfo.SlugPath != "" {
		parentPath = splitSlugToPath(csrInfo.SlugPath)
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

// listCSRsAPI handles the GET /v1/certificate-requests endpoint
func listCSRsAPI(w http.ResponseWriter, r *http.Request) {
	var parentPath string

	// Read in the submitted GET URL parameters
	queryParams := r.URL.Query()
	parentCNPath, presentCN := queryParams["parent_cn_path"]
	parentSlugPath, presentSlug := queryParams["parent_slug_path"]
	if presentCN {
		parentPath = splitCommonNamesToPath(parentCNPath[0])
	}
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
