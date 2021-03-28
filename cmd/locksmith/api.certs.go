package locksmith

import (
	"crypto/x509"
	b64 "encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"path/filepath"
	"strings"
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
	var csr *x509.CertificateRequest

	if certInfo.CommonNamePath != "" {
		parentPath = splitCACNChainToPath(certInfo.CommonNamePath)
		parentPathRaw = certInfo.CommonNamePath
	}
	if certInfo.SlugPath != "" {
		parentPath = splitCACNChainToPath(certInfo.SlugPath)
		parentPathRaw = certInfo.SlugPath
	}

	// Check to see if the CSR is passed in via base64 encoded input
	if certInfo.CertificateRequestInput.FromPEM != "" {
		csrSource, err := b64.StdEncoding.DecodeString(certInfo.CertificateRequestInput.FromPEM)
		check(err)

		csrPEM, err := decodeByteSliceToPEM(csrSource, "CERTIFICATE REQUEST")
		check(err)

		csr, err = readCSR(csrPEM.Bytes)
		check(err)
	}
	// Check to see if we can retreive the CSR from the file system
	if certInfo.CertificateRequestInput.FromCAPath.Target != "" && certInfo.CertificateRequestInput.FromCAPath.CNPath != "" {
		// See if the CAPath is valid
		csrCAPath := splitCACNChainToPath(certInfo.CertificateRequestInput.FromCAPath.CNPath)

		csrCAAbsPath, err := filepath.Abs(readConfig.Locksmith.PKIRoot + "/roots/" + csrCAPath)
		checkAndFail(err)

		csrCAParentPathExists, err := DirectoryExists(csrCAAbsPath)
		check(err)
		if !csrCAParentPathExists {
			// Path invalid, return error
			returnData := &ReturnGenericMessage{
				Status:   "invalid-csr-parent-path",
				Errors:   []string{"Invalid parent path for CSR, no chain exists!"},
				Messages: []string{}}
			returnResponse, _ := json.Marshal(returnData)
			fmt.Fprintf(w, string(returnResponse))
			return
		}
		// Explode Target - See if the Target is a valid CSR target type
		csrTarget := strings.Split(certInfo.CertificateRequestInput.FromCAPath.Target, "/")
		if csrTarget[0] != "certreqs" {
			// Target type invalid, return error
			returnData := &ReturnGenericMessage{
				Status:   "invalid-csr-target-type",
				Errors:   []string{"Invalid target type for CSR, expecting 'certreqs/target_id'!"},
				Messages: []string{}}
			returnResponse, _ := json.Marshal(returnData)
			fmt.Fprintf(w, string(returnResponse))
			return
		}

		// See if the Target is a valid CSR
		csrFileExists, err := FileExists(csrCAAbsPath + "/certreqs/" + slugger(csrTarget[1]) + ".req.pem")
		if csrFileExists {
			csr, err = readCSRFromFile(csrCAAbsPath + "/certreqs/" + slugger(csrTarget[1]) + ".req.pem")
			check(err)
		} else {
			// Target doesn't exist, return error
			returnData := &ReturnGenericMessage{
				Status:   "invalid-csr-target",
				Errors:   []string{"Invalid target for CSR, no such target exists!"},
				Messages: []string{}}
			returnResponse, _ := json.Marshal(returnData)
			fmt.Fprintf(w, string(returnResponse))
			return
		}
	}
	if csr == nil {
		// No CSR, return error
		returnData := &ReturnGenericMessage{
			Status:   "no-csr-supplied",
			Errors:   []string{"No CSR was provided, no certificate to build from!"},
			Messages: []string{}}
		returnResponse, _ := json.Marshal(returnData)
		fmt.Fprintf(w, string(returnResponse))
		return
	}

	// Neither options are submitted - error
	if parentPath == "" {
		returnData := &ReturnGenericMessage{
			Status:   "missing-parent-path",
			Errors:   []string{"Missing parent path!  Must supply either `cn_path` or `slug_path`"},
			Messages: []string{}}
		returnResponse, _ := json.Marshal(returnData)
		fmt.Fprintf(w, string(returnResponse))
		return
	} else {

		// Check if the parent path directory exists
		absPath, err := filepath.Abs(readConfig.Locksmith.PKIRoot + "/roots/" + parentPath)
		checkAndFail(err)

		certCAParentPathExists, err := DirectoryExists(absPath)
		check(err)

		if certCAParentPathExists {

			certName := csr.Subject.CommonName
			sluggedCertCommonName := slugger(certName)

			// If the Cert's parent path exists, check if the cert exists before (re)creating it
			logNeworkRequestStdOut(certName+" ("+sluggedCertCommonName+"): Checking "+absPath+"/certs/"+sluggedCertCommonName+".pem", r)
			sluggedCertFileExists, err := FileExists(absPath + "/certs/" + sluggedCertCommonName + ".pem")
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
				return
			} else {
				// Cert does not exist, go ahead with creation

				logNeworkRequestStdOut(certName+" ("+sluggedCertCommonName+") Creating certificate in '"+parentPathRaw+"'", r)
				certCreated, certificate, messages, err := createNewCertificateFromCSR(absPath, certInfo.SigningPrivateKeyPassphrase, csr)
				check(err)

				if certCreated {
					logNeworkRequestStdOut(certName+" ("+sluggedCertCommonName+") cert-created in '"+parentPathRaw+"'", r)
					returnData := &RESTPOSTCertificateJSONReturn{
						Status:   "success",
						Errors:   []string{},
						Messages: []string{"Successfully created Certificate " + certName + " in '" + parentPathRaw + "'!"},
						CertInfo: CertificateInfo{
							Slug:           sluggedCertCommonName,
							Certificate:    certificate,
							CertificatePEM: b64.StdEncoding.EncodeToString(pemEncodeCSR(certificate.Raw).Bytes())}}
					returnResponse, _ := json.Marshal(returnData)
					fmt.Fprintf(w, string(returnResponse))
					return

				} else {
					// Certificate wasn't created, return error
					logNeworkRequestStdOut(certName+" ("+sluggedCertCommonName+") error-creating-cert", r)
					returnData := &ReturnGenericMessage{
						Status:   "certificate-creation-error",
						Errors:   []string{"Error creating Certificate " + certName + " in '" + parentPathRaw + "'!"},
						Messages: messages}
					returnResponse, _ := json.Marshal(returnData)
					fmt.Fprintf(w, string(returnResponse))
					return
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
			return
		}
	}

}
