package main

import (
	"log"
	"net/http"
)

// logNetworkRequestStdOut adds a logger wrapper to add extra network client information to the log
func logNeworkRequestStdOut(s string, r *http.Request) {
	logStdOut("IP[" + ReadUserIP(r) + "] UA[" + r.UserAgent() + "] " + string(s))
	//log.Printf("[%s] %s\n", ReadUserIP(r), string(s))
}

// logStdOut just logs something to stdout
func logStdOut(s string) {
	log.Printf("%s\n", string(s))
}

// logStdErr just logs to stderr
func logStdErr(s string) {
	log.Fatalf("%s\n", string(s))
}

// check does error checking
func check(e error) {
	if e != nil {
		log.Printf("error: %v", e)
	}
}

// checkAndFail checks for an error type and fails
func checkAndFail(e error) {
	if e != nil {
		log.Fatalf("error: %v", e)
	}
}
