package main

import (
	"io/ioutil"
	"path/filepath"
)

// readSerialNumber reads the serial.txt file out
func readSerialNumber(rootSlug string) string {
	dat, err := ioutil.ReadFile(readConfig.Locksmith.PKIRoot + "/roots/" + rootSlug + "/serial.txt")
	check(err)
	return string(dat)
}

// createNewRootCAFilesystem
func createNewCAFilesystem(rootSlug string) {
	rootSlugPath := readConfig.Locksmith.PKIRoot + "/roots/" + rootSlug
	//Create root CA directory
	rootCAPath, err := filepath.Abs(rootSlugPath)
	check(err)
	CreateDirectory(rootCAPath)

	rootCACertsPath, err := filepath.Abs(rootSlugPath + "/certs")
	check(err)
	CreateDirectory(rootCACertsPath)

	rootCACertRevListPath, err := filepath.Abs(rootSlugPath + "/crls")
	check(err)
	CreateDirectory(rootCACertRevListPath)

	rootCACertKeysPath, err := filepath.Abs(rootSlugPath + "/keys")
	check(err)
	CreateDirectory(rootCACertKeysPath)

	rootCACertRequestsPath, err := filepath.Abs(rootSlugPath + "/reqs")
	check(err)
	CreateDirectory(rootCACertRequestsPath)

	rootCACertSerialFilePath, err := filepath.Abs(rootSlugPath + "/serial.txt")
	check(err)

	// Check to see if there is a serial file
	serialFile, err := WriteFile(rootCACertSerialFilePath, "1", 0600, false)
	check(err)
	if serialFile {
		logStdOut("Created serial file")
	} else {
		logStdOut("Serial file exists")
	}

	// Check for certificate authority key pair
	caKeyPath, err := filepath.Abs(rootSlugPath + "/keys/ca.priv.pem")
	check(err)
	caKeyCheck, err := FileExists(caKeyPath)
	check(err)
	if !caKeyCheck {
		rootPrivKey, rootPubKey, err := generateRSAKeypair(4096)
		check(err)

		rootPrivKeyFile, rootPubKeyFile, err := writeRSAKeyPair(pemEncodeRSAPrivateKey(rootPrivKey), pemEncodeRSAPublicKey(rootPubKey), rootCACertKeysPath+"/ca")
		check(err)
		if rootPrivKeyFile && rootPubKeyFile {
			logStdOut("Private Key Created")
		}
	}
}
