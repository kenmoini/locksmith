package locksmith

import (
	"crypto/x509"
	b64 "encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"path/filepath"
)

// readAuthorityAPI handles the GET /v1/authority endpoint
func readAuthorityAPI(w http.ResponseWriter, r *http.Request) {
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
			// Check to see if the ca.pem file exists
			caCertExists, err := FileExists(absPath + "/certs/ca.pem")
			check(err)

			if caCertExists {
				// Read in PEM file
				pem, err := readPEMFile(absPath+"/certs/ca.pem", "CERTIFICATE")
				check(err)

				// Decode to Certfificate object
				certificate, err := x509.ParseCertificate(pem.Bytes)
				check(err)

				returnData := &RESTGETAuthorityJSONReturn{
					Status:          "success",
					Errors:          []string{},
					Messages:        []string{"Certificate Authority information for '" + caPathRaw + "'"},
					Slug:            caPathRaw,
					CertificatePEM:  b64.StdEncoding.EncodeToString(pem.Bytes),
					CertificateInfo: certificate}
				returnResponse, _ := json.Marshal(returnData)
				fmt.Fprintf(w, string(returnResponse))
			} else {
				// Certificate Authority Certificate PEM File does not exists
				returnData := &ReturnGenericMessage{
					Status:   "no-ca-certificate",
					Errors:   []string{err.Error()},
					Messages: []string{"Certificate Authority Certificate PEM File does not exists!"}}
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
