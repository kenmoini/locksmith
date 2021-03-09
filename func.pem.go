package main

import (
	"bytes"
	"encoding/pem"
	"io/ioutil"
	"log"
)

// readPEMFile reads a PEM file and decodes it, along with a type check
// Types can include CERTIFICATE REQUEST, CERTIFICATE, PRIVATE KEY, PUBLIC KEY
func readPEMFile(path string, matchType string) (*pem.Block, error) {
	fileBytes, err := ReadFileToBytes(path)
	check(err)

	block, rest := pem.Decode(fileBytes)
	if block == nil || block.Type != matchType {
		log.Fatal("failed to decode PEM block containing a " + matchType + ": " + string(rest))
	}
	return block, nil
}

// writePEMFile takes a PEM encoded bytes stream and saves it to a file
func writePEMFile(certPem *bytes.Buffer, path string) (bool, error) {
	pemByte, _ := ioutil.ReadAll(certPem)
	keyFile, err := WriteByteFile(path, pemByte, 0600, false)
	if err != nil {
		return false, err
	}
	return keyFile, nil
}
