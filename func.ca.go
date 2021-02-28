package main

import (
	"io/ioutil"
	"path/filepath"
)

// createNewRootCAFilesystem
func createNewCAFilesystem(rootSlug string) {
	//Create root CA directory
	rootCAPath, err := filepath.Abs(readConfig.Locksmith.PKIRoot + "/roots/" + rootSlug)
	check(err)
	CreateDirectory(rootCAPath)

	rootCACertsPath, err := filepath.Abs(readConfig.Locksmith.PKIRoot + "/roots/" + rootSlug + "/certs")
	check(err)
	CreateDirectory(rootCACertsPath)

	rootCACertRevListPath, err := filepath.Abs(readConfig.Locksmith.PKIRoot + "/roots/" + rootSlug + "/crls")
	check(err)
	CreateDirectory(rootCACertRevListPath)

	rootCACertKeysPath, err := filepath.Abs(readConfig.Locksmith.PKIRoot + "/roots/" + rootSlug + "/keys")
	check(err)
	CreateDirectory(rootCACertKeysPath)

	rootCACertRequestsPath, err := filepath.Abs(readConfig.Locksmith.PKIRoot + "/roots/" + rootSlug + "/reqs")
	check(err)
	CreateDirectory(rootCACertRequestsPath)

	rootCACertSerialFilePath, err := filepath.Abs(readConfig.Locksmith.PKIRoot + "/roots/" + rootSlug + "/serial.txt")
	check(err)
	// Check to see if there is a serial file
	serialCheck, err := FileExists(rootCACertSerialFilePath)
	check(err)
	// If not, create one with a starting digit
	if !serialCheck {
		d1 := []byte("1")
		err = ioutil.WriteFile(rootCACertSerialFilePath, d1, 0700)
		check(err)
	}

}
