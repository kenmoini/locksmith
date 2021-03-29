package locksmith

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

// NewRouter generates the router used in the HTTP Server
func NewRouter(basePath string) *http.ServeMux {

	var formattedBasePath string
	var apiVersionTag string

	if basePath == "" {
		basePath = "/locksmith"
	}
	formattedBasePath = strings.TrimRight(basePath, "/")

	// Create router and define routes and return that router
	router := http.NewServeMux()

	//====================================================================================
	// TEST ENDPOINT
	// Test out a random function maybe
	router.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		returnData := &ReturnGenericMessage{
			Status:   "test",
			Errors:   []string{},
			Messages: []string{"Test!"}}
		returnResponse, _ := json.Marshal(returnData)
		fmt.Fprintf(w, string(returnResponse))
	})

	//====================================================================================
	// KUBERNETES ENDPOINTS
	// Version Output - reads from variables.go
	router.HandleFunc("/version", func(w http.ResponseWriter, r *http.Request) {
		applicationVersionAPI(w, r)
	})

	// Healthz endpoint for kubernetes platforms
	router.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		healthZAPI(w, r)
	})

	//====================================================================================
	// START V1 API
	//====================================================================================
	apiVersionTag = "/v1"

	//====================================================================================
	// KEY PAIRS
	// Key Manipulation - Listing, Creating, Deleting
	router.HandleFunc(formattedBasePath+apiVersionTag+"/keys", func(w http.ResponseWriter, r *http.Request) {
		logNeworkRequestStdOut(r.Method+" "+r.RequestURI, r)
		switch r.Method {
		case "GET":
			// index - get list of keys in key store
			listKeyPairsAPI(w, r)
		default:
			methodNotAllowedAPI(w, r)
		}
	})

	router.HandleFunc(formattedBasePath+apiVersionTag+"/key", func(w http.ResponseWriter, r *http.Request) {
		logNeworkRequestStdOut(r.Method+" "+r.RequestURI, r)
		switch r.Method {
		case "GET":
			// index - get key pair in key store
			readKeyPairAPI(w, r)
		case "POST":
			// create - create new keys in key store
			createKeyPairAPI(w, r)
		default:
			methodNotAllowedAPI(w, r)
		}
	})

	//====================================================================================
	// KEY STORES
	// Key Store Manipulation - Listing, Creating
	router.HandleFunc(formattedBasePath+apiVersionTag+"/keystores", func(w http.ResponseWriter, r *http.Request) {
		logNeworkRequestStdOut(r.Method+" "+r.RequestURI, r)
		switch r.Method {
		case "GET":
			// index - get list of key stores
			listKeyStoresAPI(w, r)
		default:
			methodNotAllowedAPI(w, r)
		}
	})

	router.HandleFunc(formattedBasePath+apiVersionTag+"/keystore", func(w http.ResponseWriter, r *http.Request) {
		logNeworkRequestStdOut(r.Method+" "+r.RequestURI, r)
		switch r.Method {
		case "POST":
			// create - create new key store
			createKeyStoreAPI(w, r)
		default:
			methodNotAllowedAPI(w, r)
		}
	})

	//====================================================================================
	// ROOT CERTIFICATE AUTHORITIES
	// Root CA Manipulation - Listing, Creating, Deleting
	router.HandleFunc(formattedBasePath+apiVersionTag+"/roots", func(w http.ResponseWriter, r *http.Request) {
		logNeworkRequestStdOut(r.Method+" "+r.RequestURI, r)
		switch r.Method {
		case "GET":
			// index - get list of roots
			listRootCAsAPI(w, r)
			// http.ServeFile(w, r, "form.html")
		default:
			methodNotAllowedAPI(w, r)
		}
	})

	router.HandleFunc(formattedBasePath+apiVersionTag+"/root", func(w http.ResponseWriter, r *http.Request) {
		logNeworkRequestStdOut(r.Method+" "+r.RequestURI, r)
		switch r.Method {
		case "POST":
			// create - create new root
			createNewRootCAAPI(w, r)
		default:
			methodNotAllowedAPI(w, r)
		}
	})

	//====================================================================================
	// INTERMEDIATE CERTIFICATE AUTHORITIES
	// Intermediate CA Manipulation - Listing, Creating, Deleting
	router.HandleFunc(formattedBasePath+apiVersionTag+"/intermediates", func(w http.ResponseWriter, r *http.Request) {
		logNeworkRequestStdOut(r.Method+" "+r.RequestURI, r)
		switch r.Method {
		case "GET":
			// index - get list of intermediate CAs in parent path
			listIntermediateCAsAPI(w, r)
		default:
			methodNotAllowedAPI(w, r)
		}
	})

	router.HandleFunc(formattedBasePath+apiVersionTag+"/intermediate", func(w http.ResponseWriter, r *http.Request) {
		logNeworkRequestStdOut(r.Method+" "+r.RequestURI, r)
		switch r.Method {
		case "POST":
			// create - create new intermediate CA in parent path
			createNewIntermediateCAAPI(w, r)
		default:
			methodNotAllowedAPI(w, r)
		}
	})

	//====================================================================================
	// AUTHORITY
	// Reading a Certificate Authority's Information
	router.HandleFunc(formattedBasePath+apiVersionTag+"/authority", func(w http.ResponseWriter, r *http.Request) {
		logNeworkRequestStdOut(r.Method+" "+r.RequestURI, r)
		switch r.Method {
		case "GET":
			// index - get information for CA in parent path
			readAuthorityAPI(w, r)
		default:
			methodNotAllowedAPI(w, r)
		}
	})

	//====================================================================================
	// REVOCATIONS
	// Reading a Certificate Authority's Certificate Revocation List
	router.HandleFunc(formattedBasePath+apiVersionTag+"/revocations", func(w http.ResponseWriter, r *http.Request) {
		logNeworkRequestStdOut(r.Method+" "+r.RequestURI, r)
		switch r.Method {
		case "GET":
			// index - get CRL for CA in parent path
			readRevocationListAPI(w, r)
		default:
			methodNotAllowedAPI(w, r)
		}
	})

	//====================================================================================
	// CERTIFICATE REQUESTS
	// CSR Manipulation - Listing, Creating, Deleting
	router.HandleFunc(formattedBasePath+apiVersionTag+"/certificate-requests", func(w http.ResponseWriter, r *http.Request) {
		logNeworkRequestStdOut(r.Method+" "+r.RequestURI, r)
		switch r.Method {
		case "GET":
			// index - get list of CSRs in cert path
			listCSRsAPI(w, r)
		default:
			methodNotAllowedAPI(w, r)
		}
	})

	router.HandleFunc(formattedBasePath+apiVersionTag+"/certificate-request", func(w http.ResponseWriter, r *http.Request) {
		logNeworkRequestStdOut(r.Method+" "+r.RequestURI, r)
		switch r.Method {
		case "GET":
			// index - read a CSR in cert path
			readCSRAPI(w, r)
		case "POST":
			// create - create new csr in cert path
			createNewCSRAPI(w, r)
		default:
			methodNotAllowedAPI(w, r)
		}
	})

	//====================================================================================
	// CERTIFICATES
	// Certificate Functions - List certs, Create certs from CA slug
	router.HandleFunc(formattedBasePath+apiVersionTag+"/certificates", func(w http.ResponseWriter, r *http.Request) {
		logNeworkRequestStdOut(r.Method+" "+r.RequestURI, r)
		switch r.Method {
		case "GET":
			// index - get list of certs for a ca path
			listCertsAPI(w, r)
		default:
			methodNotAllowedAPI(w, r)
		}
	})

	router.HandleFunc(formattedBasePath+apiVersionTag+"/certificate", func(w http.ResponseWriter, r *http.Request) {
		logNeworkRequestStdOut(r.Method+" "+r.RequestURI, r)
		switch r.Method {
		case "GET":
			// index - read cert info in a ca path
			readCertificateAPI(w, r)
		case "POST":
			// create - make a new cert
			createNewCertAPI(w, r)
		default:
			methodNotAllowedAPI(w, r)
		}
	})

	return router
}

// RunHTTPServer will run the HTTP Server
func (config Config) RunHTTPServer() {
	// Set up a channel to listen to for interrupt signals
	var runChan = make(chan os.Signal, 1)

	// Set up a context to allow for graceful server shutdowns in the event
	// of an OS interrupt (defers the cancel just in case)
	ctx, cancel := context.WithTimeout(
		context.Background(),
		config.Locksmith.Server.Timeout.Server,
	)
	defer cancel()

	// Define server options
	server := &http.Server{
		Addr:         config.Locksmith.Server.Host + ":" + config.Locksmith.Server.Port,
		Handler:      NewRouter(config.Locksmith.Server.BasePath),
		ReadTimeout:  config.Locksmith.Server.Timeout.Read * time.Second,
		WriteTimeout: config.Locksmith.Server.Timeout.Write * time.Second,
		IdleTimeout:  config.Locksmith.Server.Timeout.Idle * time.Second,
	}

	// Only listen on IPV4
	l, err := net.Listen("tcp4", config.Locksmith.Server.Host+":"+config.Locksmith.Server.Port)
	check(err)

	// Handle ctrl+c/ctrl+x interrupt
	signal.Notify(runChan, os.Interrupt, syscall.SIGTSTP)

	// Alert the user that the server is starting
	log.Printf("Server is starting on %s\n", server.Addr)

	// QUIK N DUURTY TESTS
	//csrPrivKey, csrPubKey, err := GenerateRSAKeypair(4096)
	//check(err)
	//
	//pemEncodedPrivateKey, encryptedPrivateKeyBytes := pemEncodeRSAPrivateKey(csrPrivKey, "s3cr3t")
	//
	//logStdOut("privateKeyPEM: " + string(pemEncodedPrivateKey.Bytes()))
	//logStdOut("encryptedPrivateKeyBytes: " + string(encryptedPrivateKeyBytes.Bytes()))
	//base, _ := b64.StdEncoding.DecodeString(string(encryptedPrivateKeyBytes.Bytes()))
	//_, decrypt, _ := decryptBytes(base, "s3cr3t")
	//logStdOut("decrypted: " + string(decrypt))
	//logStdOut("publicKeyPEM: " + string(pemEncodeRSAPublicKey(csrPubKey).Bytes()))

	// Run the server on a new goroutine
	go func() {
		//if err := server.ListenAndServe(); err != nil {
		if err := server.Serve(l); err != nil {
			if err == http.ErrServerClosed {
				// Normal interrupt operation, ignore
			} else {
				log.Fatalf("Server failed to start due to err: %v", err)
			}
		}

	}()

	// Block on this channel listeninf for those previously defined syscalls assign
	// to variable so we can let the user know why the server is shutting down
	interrupt := <-runChan

	// If we get one of the pre-prescribed syscalls, gracefully terminate the server
	// while alerting the user
	log.Printf("Server is shutting down due to %+v\n", interrupt)
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server was unable to gracefully shutdown due to err: %+v", err)
	}
}
