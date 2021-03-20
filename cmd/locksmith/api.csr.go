package locksmith

import (
	b64 "encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"path/filepath"
)

// listCSRsAPI handles the GET /v1/certificate-requests endpoint
func listCSRsAPI(w http.ResponseWriter, r *http.Request) {
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
			certificateRequests := DirectoryListingNamesNoExt(absPath + "/certreqs/")

			certificateRequests = rmStrFromStrSlice("ca", certificateRequests)

			if len(certificateRequests) > 0 {
				// Got some hits
				returnData := &RESTGETCertificateRequestsJSONReturn{
					Status:              "success",
					Errors:              []string{},
					Messages:            []string{"Listing of CSRs for CA Path '" + parentPathRaw + "'"},
					CertificateRequests: certificateRequests}
				returnResponse, _ := json.Marshal(returnData)
				fmt.Fprintf(w, string(returnResponse))
			} else {
				// No CSRs in this path
				returnData := &ReturnGenericMessage{
					Status:   "no-csrs-found",
					Errors:   []string{"No Certificate Requests found in this CA Path!"},
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

// createNewCSRAPI handles the POST /v1/certificate-requests endpoint
func createNewCSRAPI(w http.ResponseWriter, r *http.Request) {
	// Load in POST JSON Data
	csrInfo := RESTPOSTCertificateRequestJSONIn{}
	err := json.NewDecoder(r.Body).Decode(&csrInfo)
	check(err)

	// Set up Parent Path
	var parentPath string
	var parentPathRaw string

	if csrInfo.CommonNamePath != "" {
		parentPath = splitCACNChainToPath(csrInfo.CommonNamePath)
		parentPathRaw = csrInfo.CommonNamePath
	}
	if csrInfo.SlugPath != "" {
		parentPath = splitCACNChainToPath(csrInfo.SlugPath)
		parentPathRaw = csrInfo.SlugPath
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

		// Check if the parent path directory exists
		absPath, err := filepath.Abs(readConfig.Locksmith.PKIRoot + "/roots/" + parentPath)
		checkAndFail(err)
		csrCAParentPathExists, err := DirectoryExists(absPath)
		check(err)

		if csrCAParentPathExists {

			csrName := csrInfo.CertificateConfiguration.Subject.CommonName
			sluggedCSRCommonName := slugger(csrName)

			// If the CSR's parent path exists, check if the CSR exists before (re)creating it
			logNeworkRequestStdOut(csrName+" ("+sluggedCSRCommonName+"): Checking "+absPath+"/certreqs/"+sluggedCSRCommonName+".req.pem", r)
			sluggedCSRFileExists, err := DirectoryExists(absPath + "/certreqs/" + sluggedCSRCommonName + ".req.pem")
			check(err)

			if sluggedCSRFileExists {
				// if the intermediate exists, return with an intermed-ca-exists error
				logNeworkRequestStdOut(csrName+" ("+sluggedCSRCommonName+") csr-exists in '"+parentPathRaw+"'", r)
				returnData := &ReturnGenericMessage{
					Status:   "certificate-request-exists",
					Errors:   []string{"Certificate " + csrName + " exists in '" + parentPathRaw + "'!"},
					Messages: []string{}}
				returnResponse, _ := json.Marshal(returnData)
				fmt.Fprintf(w, string(returnResponse))
			} else {
				// CSR does not exist, go ahead with creation

				logNeworkRequestStdOut(csrName+" ("+sluggedCSRCommonName+") creating certificate request in '"+parentPathRaw+"'", r)
				csrCreated, messages, csrCert, keyPair, err := createNewCertificateRequest(csrInfo, absPath)
				check(err)
				if csrCreated {
					logNeworkRequestStdOut(csrName+" ("+sluggedCSRCommonName+") csr-created in '"+parentPathRaw+"'", r)
					returnData := &RESTPOSTCertificateRequestJSONReturn{
						Status:   "certificate-request-created",
						Errors:   []string{},
						Messages: []string{"Successfully created Certificate Request " + csrName + " in '" + parentPathRaw + "'!"},
						CSRInfo: CertificateRequestInfo{
							Slug:               sluggedCSRCommonName,
							CertificateRequest: csrCert,
							KeyPair: KeyPair{
								PublicKey:  b64.StdEncoding.EncodeToString(pemEncodeRSAPublicKey(keyPair.PublicKey).Bytes()),
								PrivateKey: b64.StdEncoding.EncodeToString(pemEncodeRSAPrivateKey(keyPair.PrivateKey, csrInfo.CertificateConfiguration.RSAPrivateKeyPassphrase).Bytes())}}}
					returnResponse, _ := json.Marshal(returnData)
					fmt.Fprintf(w, string(returnResponse))
				} else {
					logNeworkRequestStdOut(csrName+" ("+sluggedCSRCommonName+") error-creating-csr", r)
					returnData := &ReturnGenericMessage{
						Status:   "certificate-request-creation-error",
						Errors:   []string{"Error creating Certificate Request " + csrName + " in '" + parentPathRaw + "'!"},
						Messages: messages}
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

}
