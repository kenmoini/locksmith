package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"
)

// NewRouter generates the router used in the HTTP Server
func NewRouter(basePath string) *http.ServeMux {
	if basePath == "" {
		basePath = "/locksmith"
	}
	// Create router and define routes and return that router
	router := http.NewServeMux()

	router.HandleFunc(basePath+"/version", func(w http.ResponseWriter, r *http.Request) {
		logNeworkRequestStdOut(r.Method+" "+basePath+"/version", r)
		fmt.Fprintf(w, "Locksmith version: %s\n", locksmithVersion)
	})

	router.HandleFunc(basePath+"/healthz", func(w http.ResponseWriter, r *http.Request) {
		logNeworkRequestStdOut(r.Method+" "+basePath+"/healthz", r)
		fmt.Fprintf(w, "OK")
	})

	router.HandleFunc(basePath+"/roots", func(w http.ResponseWriter, r *http.Request) {
		logNeworkRequestStdOut(r.Method+" "+basePath+"/roots", r)
		switch r.Method {
		// index - get list of roots
		case "GET":
			rootListing := DirectoryListingNames(readConfig.Locksmith.PKIRoot + "/roots/")
			returnData := &ReturnGetRoots{
				Status:   "success",
				Errors:   []string{},
				Messages: []string{},
				Roots:    rootListing}
			returnResponse, _ := json.Marshal(returnData)
			fmt.Fprintf(w, string(returnResponse))
			// http.ServeFile(w, r, "form.html")
		case "POST":
			// create - create new root

			if err := r.ParseForm(); err != nil {
				fmt.Fprintf(w, "ParseForm() err: %v", err)
				return
			}

			caName := r.FormValue("name")
			sluggedName := slugger(caName)
			rsaPrivateKeyPassword := r.FormValue("rsaPrivateKeyPassword")

			checkForRootPath, err := filepath.Abs(readConfig.Locksmith.PKIRoot + "/roots/" + sluggedName)
			check(err)

			caRootPathExists, err := DirectoryExists(checkForRootPath)
			check(err)

			if caRootPathExists {
				logNeworkRequestStdOut(caName+" ("+sluggedName+") root-exists", r)
				returnData := &ReturnPostRoots{
					Status:   "root-exists",
					Errors:   []string{},
					Messages: []string{},
					Root: RootInfo{
						Slug:   sluggedName,
						Serial: readSerialNumber(sluggedName)}}
				returnResponse, _ := json.Marshal(returnData)
				fmt.Fprintf(w, string(returnResponse))
			} else {
				createNewCAFilesystem(sluggedName, caName, rsaPrivateKeyPassword)
				logNeworkRequestStdOut(caName+" ("+sluggedName+") root-created", r)
				returnData := &ReturnPostRoots{
					Status:   "root-created",
					Errors:   []string{},
					Messages: []string{},
					Root: RootInfo{
						Slug:   sluggedName,
						Serial: readSerialNumber(sluggedName)}}
				returnResponse, _ := json.Marshal(returnData)
				fmt.Fprintf(w, string(returnResponse))
			}
		default:
			returnData := &ReturnGenericMessage{
				Status:   "method-not-allowed",
				Errors:   []string{},
				Messages: []string{"method not allowed"}}
			returnResponse, _ := json.Marshal(returnData)
			fmt.Fprintf(w, string(returnResponse))
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
