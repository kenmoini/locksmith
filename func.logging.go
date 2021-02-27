package main

import "log"

// check does error checking
func check(e error) {
	if e != nil {
		log.Fatalf("error: %v", e)
	}
}
