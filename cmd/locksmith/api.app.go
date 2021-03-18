package locksmith

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// applicationVersionAPI returns the application version
func applicationVersionAPI(w http.ResponseWriter, r *http.Request) {
	logNeworkRequestStdOut(r.Method+" "+r.RequestURI, r)
	fmt.Fprintf(w, "Locksmith version: %s\n", locksmithVersion)
}

// healthZAPI returns the application availability & health
func healthZAPI(w http.ResponseWriter, r *http.Request) {
	logNeworkRequestStdOut(r.Method+" "+r.RequestURI, r)
	fmt.Fprintf(w, "OK")
}

// methodNotAllowedAPI is a generic short function to return the method not allowed JSON
func methodNotAllowedAPI(w http.ResponseWriter, r *http.Request) {
	returnData := &ReturnGenericMessage{
		Status:   "method-not-allowed",
		Errors:   []string{"method not allowed"},
		Messages: []string{}}
	returnResponse, _ := json.Marshal(returnData)
	fmt.Fprintf(w, string(returnResponse))
}
