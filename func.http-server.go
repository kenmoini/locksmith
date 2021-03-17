package main

import (
	"context"
	b64 "encoding/base64"
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
	router.HandleFunc(formattedBasePath+"/test", func(w http.ResponseWriter, r *http.Request) {

		logStdOut(b64.StdEncoding.EncodeToString(encryptBytes([]byte("This is my secret text!"), "passw0rd")))

		encoded := "R59BHVuy9FYc0jFkb8pbMseZn/98tnXf/0zl7lqOJMxtW494z7NLpp9I23faiMv7hRiw"
		encodedBytes, err := b64.StdEncoding.DecodeString(encoded)
		check(err)
		bit, byted, _ := decryptBytes(encodedBytes, "passw0rd")
		if bit {
			logStdOut(string(byted))
		}

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
	router.HandleFunc(formattedBasePath+"/version", func(w http.ResponseWriter, r *http.Request) {
		APIApplicationVersion(w, r)
	})

	// Healthz endpoint for kubernetes platforms
	router.HandleFunc(formattedBasePath+"/healthz", func(w http.ResponseWriter, r *http.Request) {
		APIHealthZ(w, r)
	})

	//====================================================================================
	// START V1 API
	//====================================================================================
	apiVersionTag = "/v1"

	//====================================================================================
	// CERTIFICATES
	// Certificate Functions - List certs, Create certs from CA slug
	router.HandleFunc(formattedBasePath+apiVersionTag+"/certs", func(w http.ResponseWriter, r *http.Request) {
		logNeworkRequestStdOut(r.Method+" "+r.RequestURI, r)
		switch r.Method {
		case "GET":
			// index - get list of certs for a ca path
			queryParams := r.URL.Query()
			caPath, present := queryParams["ca_path"] //ca_path=["root-ca/intermed-ca/sub-ca"]
			if !present || len(caPath) == 0 {
				returnData := &ReturnGenericMessage{
					Status:   "no-ca-path",
					Errors:   []string{},
					Messages: []string{"No Certificate Authority Path!"}}
				returnResponse, _ := json.Marshal(returnData)
				fmt.Fprintf(w, string(returnResponse))
			} else {
				// Split the path along the path delimiter
				splitPath := strings.Split(caPath[0], "/")
				logStdOut(splitPath[0])
			}
		case "POST":
			// create - make a new cert and CSR

		default:
			APIMethodNotAllowed(w, r)
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
			APIListRootCAs(w, r)
			// http.ServeFile(w, r, "form.html")
		case "POST":
			// create - create new root
			APICreateNewRootCA(w, r)
		default:
			APIMethodNotAllowed(w, r)
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
			APIListIntermediateCAs(w, r)
		case "POST":
			// create - create new intermedate CA in parent path
			APICreateNewIntermediateCA(w, r)
		default:
			APIMethodNotAllowed(w, r)
		}
	})

	//====================================================================================
	// KEY PAIRS
	// Key Manipulation - Listing, Creating, Deleting
	router.HandleFunc(formattedBasePath+apiVersionTag+"/keys", func(w http.ResponseWriter, r *http.Request) {
		logNeworkRequestStdOut(r.Method+" "+r.RequestURI, r)
		switch r.Method {
		case "GET":
			// index - get list of keys in key store
			APIListKeyPairs(w, r)
		case "POST":
			// create - create new keys in key store
			APICreateKeyPair(w, r)
		default:
			APIMethodNotAllowed(w, r)
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
			APIListCSRs(w, r)
		case "POST":
			// create - create new csr in cert path
			APICreateNewCSR(w, r)
		default:
			APIMethodNotAllowed(w, r)
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
