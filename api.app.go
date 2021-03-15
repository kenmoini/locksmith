package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// APIApplicationVersion returns the application version
func APIApplicationVersion(w http.ResponseWriter, r *http.Request) {
	logNeworkRequestStdOut(r.Method+" "+r.RequestURI, r)
	fmt.Fprintf(w, "Locksmith version: %s\n", locksmithVersion)
}

// APIHealthZ returns the application availability & health
func APIHealthZ(w http.ResponseWriter, r *http.Request) {
	logNeworkRequestStdOut(r.Method+" "+r.RequestURI, r)
	fmt.Fprintf(w, "OK")
}

// APIMethodNotAllowed is a generic short function to return the method not allowed JSON
func APIMethodNotAllowed(w http.ResponseWriter, r *http.Request) {
	returnData := &ReturnGenericMessage{
		Status:   "method-not-allowed",
		Errors:   []string{"method not allowed"},
		Messages: []string{}}
	returnResponse, _ := json.Marshal(returnData)
	fmt.Fprintf(w, string(returnResponse))
}
