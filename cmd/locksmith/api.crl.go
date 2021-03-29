package locksmith

import (
	"crypto/x509"
	"encoding/json"
	"fmt"
	"net/http"
	"path/filepath"
)

// readRevocationListAPI handles the GET /v1/revocations endpoint
func readRevocationListAPI(w http.ResponseWriter, r *http.Request) {
	var caPath string
	var caPathRaw string

	// Read in the submitted parameters
	queryParams := r.URL.Query()
	parentCNPath, presentCN := queryParams["cn_path"]
	if presentCN {
		caPath = splitCACNChainToPath(parentCNPath[0])
		caPathRaw = parentCNPath[0]
	}
	parentSlugPath, presentSlug := queryParams["slug_path"]
	if presentSlug {
		caPath = splitCACNChainToPath(parentSlugPath[0])
		caPathRaw = parentSlugPath[0]
	}

	// Neither options are submitted - error
	if caPath == "" {
		returnData := &ReturnGenericMessage{
			Status:   "missing-parent-path",
			Errors:   []string{"Missing parent path!  Must supply either `cn_path` or `slug_path`"},
			Messages: []string{}}
		returnResponse, _ := json.Marshal(returnData)
		fmt.Fprintf(w, string(returnResponse))
	} else {

		// Check if the directory exists
		absPath, err := filepath.Abs(readConfig.Locksmith.PKIRoot + "/roots/" + caPath)
		checkAndFail(err)

		caParentPathExists, err := DirectoryExists(absPath)
		check(err)

		if caParentPathExists {
			// Check to see if the ca.crl file exists
			caCertExists, err := FileExists(absPath + "/crl/ca.crl")
			check(err)

			if caCertExists {
				// Read in PEM file
				pem, err := readPEMFile(absPath+"/crl/ca.crl", "X509 CRL")
				check(err)

				// Decode to Certfificate object
				certificateList, err := x509.ParseCRL(pem.Bytes)
				check(err)

				returnData := &RESTGETRevocationListJSONReturn{
					Status:          "success",
					Errors:          []string{},
					Messages:        []string{"Certificate Revocation List for '" + caPathRaw + "'"},
					Slug:            caPathRaw,
					CertificatePEM:  B64EncodeBytesToStr(pem.Bytes),
					CertificateList: certificateList}
				returnResponse, _ := json.Marshal(returnData)
				fmt.Fprintf(w, string(returnResponse))
			} else {
				// Certificate Authority CRL File does not exists
				returnData := &ReturnGenericMessage{
					Status:   "no-ca-crl",
					Errors:   []string{err.Error()},
					Messages: []string{"Certificate Authority Revocation List File does not exists!"}}
				returnResponse, _ := json.Marshal(returnData)
				fmt.Fprintf(w, string(returnResponse))
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
