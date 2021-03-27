package locksmith

import (
	"crypto/x509"
	b64 "encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"path/filepath"
)

// readCertificateAPI handles the GET /v1/certificate endpoint
func readCertificateAPI(w http.ResponseWriter, r *http.Request) {
	var parentPath string
	var parentPathRaw string
	var certificateID string

	// Read in the submitted GET URL parameters
	queryParams := r.URL.Query()
	parentCNPath, presentCN := queryParams["cn_path"]
	parentSlugPath, presentSlug := queryParams["slug_path"]
	certificateIn, presentCertificateID := queryParams["certificate_id"]

	if presentCN {
		parentPath = splitCACNChainToPath(parentCNPath[0])
		parentPathRaw = parentCNPath[0]
	}
	if presentSlug {
		parentPath = splitCACNChainToPath(parentSlugPath[0])
		parentPathRaw = parentSlugPath[0]
	}
	if presentCertificateID {
		certificateID = slugger(certificateIn[0])
	}

	// Neither options are submitted - error
	if parentPath == "" {
		returnData := &ReturnGenericMessage{
			Status:   "missing-parent-path",
			Errors:   []string{"Missing parent path!  Must supply either `cn_path` or `slug_path`"},
			Messages: []string{}}
		returnResponse, _ := json.Marshal(returnData)
		fmt.Fprintf(w, string(returnResponse))
	} else {
		// Check if the CA Path directory exists
		absPath, err := filepath.Abs(readConfig.Locksmith.PKIRoot + "/roots/" + parentPath)
		checkAndFail(err)

		certParentPathExists, err := DirectoryExists(absPath)
		check(err)

		if certParentPathExists {
			// certificateID has to be present and not null
			if certificateID == "" {
				returnData := &ReturnGenericMessage{
					Status:   "missing-certificate-id",
					Errors:   []string{"Missing Certificate ID!  Must supply `certificate_id`"},
					Messages: []string{}}
				returnResponse, _ := json.Marshal(returnData)
				fmt.Fprintf(w, string(returnResponse))
			} else {
				// certificateID is defined - validate that it exists

				certFileExists, err := FileExists(absPath + "/certs/" + certificateID + ".pem")
				check(err)

				if certFileExists {
					// Certificate exists, read it in and spit it out

					// Read in PEM file
					pem, err := readPEMFile(absPath+"/certs/"+certificateID+".pem", "CERTIFICATE")
					check(err)

					// Decode to Certfificate object
					certificate, err := x509.ParseCertificate(pem.Bytes)
					check(err)

					returnData := &RESTGETCertificateInformationJSONReturn{
						Status:          "success",
						Errors:          []string{},
						Messages:        []string{"Certificate information for '" + parentPathRaw + "'"},
						Slug:            certificateID,
						CertificatePEM:  b64.StdEncoding.EncodeToString(pem.Bytes),
						CertificateInfo: certificate}
					returnResponse, _ := json.Marshal(returnData)
					fmt.Fprintf(w, string(returnResponse))

				} else {
					// Certificate does not exist
					returnData := &ReturnGenericMessage{
						Status:   "no-certificate",
						Errors:   []string{},
						Messages: []string{"Certificate '" + certificateIn[0] + "' PEM File does not exists in '" + parentPathRaw + "'!"}}
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

// listCertsAPI handles the GET /v1/certificates endpoint
func listCertsAPI(w http.ResponseWriter, r *http.Request) {
	var parentPath string
	var parentPathRaw string

	// Read in the submitted GET URL parameters
	queryParams := r.URL.Query()
	parentCNPath, presentCN := queryParams["cn_path"]
	parentSlugPath, presentSlug := queryParams["slug_path"]
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
			Errors:   []string{"Missing parent path!  Must supply either `cn_path` or `slug_path`"},
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

// createNewCertAPI handles the POST /v1/certificate endpoint
func createNewCertAPI(w http.ResponseWriter, r *http.Request) {
	// Load in POST JSON Data
	certInfo := RESTPOSTCertificateJSONIn{}
	err := json.NewDecoder(r.Body).Decode(&certInfo)
	check(err)

	// Set up Parent Path
	var parentPath string
	var parentPathRaw string

	if certInfo.CommonNamePath != "" {
		parentPath = splitCACNChainToPath(certInfo.CommonNamePath)
		parentPathRaw = certInfo.CommonNamePath
	}
	if certInfo.SlugPath != "" {
		parentPath = splitCACNChainToPath(certInfo.SlugPath)
		parentPathRaw = certInfo.SlugPath
	}

	// Neither options are submitted - error
	if parentPath == "" {
		returnData := &ReturnGenericMessage{
			Status:   "missing-parent-path",
			Errors:   []string{"Missing parent path!  Must supply either `cn_path` or `slug_path`"},
			Messages: []string{}}
		returnResponse, _ := json.Marshal(returnData)
		fmt.Fprintf(w, string(returnResponse))
	} else {

		// Check if the parent path directory exists
		absPath, err := filepath.Abs(readConfig.Locksmith.PKIRoot + "/roots/" + parentPath)
		checkAndFail(err)

		certCAParentPathExists, err := DirectoryExists(absPath)
		check(err)

		if certCAParentPathExists {

			certName := certInfo.CertificateConfiguration.Subject.CommonName
			sluggedCertCommonName := slugger(certName)

			// If the Cert's parent path exists, check if the cert exists before (re)creating it
			logNeworkRequestStdOut(certName+" ("+sluggedCertCommonName+"): Checking "+absPath+"/certs/"+sluggedCertCommonName+".pem", r)
			sluggedCertFileExists, err := DirectoryExists(absPath + "/certs/" + sluggedCertCommonName + ".pem")
			check(err)

			if sluggedCertFileExists {
				// if the cert exists, return with an certificate-exists error
				logNeworkRequestStdOut(certName+" ("+sluggedCertCommonName+") certificate exists in '"+parentPathRaw+"'", r)
				returnData := &ReturnGenericMessage{
					Status:   "certificate-exists",
					Errors:   []string{"Certificate " + certName + " exists in '" + parentPathRaw + "'!"},
					Messages: []string{}}
				returnResponse, _ := json.Marshal(returnData)
				fmt.Fprintf(w, string(returnResponse))
			} else {
				// Cert does not exist, go ahead with creation

				logNeworkRequestStdOut(certName+" ("+sluggedCertCommonName+") creating certificate in '"+parentPathRaw+"'", r)
				//csrCreated, messages, csrCert, keyPair, err := createNewCertificateRequest(csrInfo, absPath)
				//check(err)
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
