package main

import (
	"log"
	"net/http"
)

func logNeworkRequestStdOut(s string, r *http.Request) {
	logStdOut("IP[" + ReadUserIP(r) + "] UA[" + r.UserAgent() + "] " + string(s))
	//log.Printf("[%s] %s\n", ReadUserIP(r), string(s))
}
func logStdOut(s string) {
	log.Printf("%s\n", string(s))
}

// check does error checking
func check(e error) {
	if e != nil {
		log.Printf("error: %v", e)
	}
}

// checkAndFail
func checkAndFail(e error) {
	if e != nil {
		log.Fatalf("error: %v", e)
	}
}
