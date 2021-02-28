package main

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"io/ioutil"
)

// generateRSAKeypair returns a private RSA key
func generateRSAKeypair(keySize int) (*rsa.PrivateKey, *rsa.PublicKey, error) {
	if keySize == 0 {
		keySize = 4096
	}
	// create our private and public key
	privKey, err := rsa.GenerateKey(rand.Reader, keySize)
	if err != nil {
		return nil, nil, err
	}
	return privKey, &privKey.PublicKey, nil
}

// writeRSAKeyPair creates key pairs
func writeRSAKeyPair(privKey *bytes.Buffer, pubKey *bytes.Buffer, path string) (bool, bool, error) {
	privKeyFile, err := writePrivateKey(privKey, path+".priv.pem")
	if err != nil {
		return false, false, err
	}

	pubKeyFile, err := writePublicKey(pubKey, path+".pub.pem")
	if err != nil {
		return privKeyFile, false, err
	}
	return privKeyFile, pubKeyFile, nil
}

// writePublicKey
func writePublicKey(pem *bytes.Buffer, path string) (bool, error) {
	pemByte, _ := ioutil.ReadAll(pem)
	keyFile, err := WriteByteFile(path, pemByte, 0644, false)
	if err != nil {
		return false, err
	}
	return keyFile, nil
}

// writePrivateKey
func writePrivateKey(pem *bytes.Buffer, path string) (bool, error) {
	pemByte, _ := ioutil.ReadAll(pem)
	keyFile, err := WriteByteFile(path, pemByte, 0600, false)
	if err != nil {
		return false, err
	}
	return keyFile, nil
}

// pemEncodeRSAPrivateKey
func pemEncodeRSAPrivateKey(caPrivKey *rsa.PrivateKey, rsaPrivateKeyPassword string) *bytes.Buffer {
	caPrivKeyPEM := new(bytes.Buffer)
	if rsaPrivateKeyPassword == "" {
		pem.Encode(caPrivKeyPEM, &pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: x509.MarshalPKCS1PrivateKey(caPrivKey),
		})
		return caPrivKeyPEM
	} else {
		// Eventually come back and add PEM encryption
		pem.Encode(caPrivKeyPEM, &pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: x509.MarshalPKCS1PrivateKey(caPrivKey),
		})
		return caPrivKeyPEM
	}
}

// pemEncodeRSAPublicKey
func pemEncodeRSAPublicKey(caPubKey *rsa.PublicKey) *bytes.Buffer {
	caPubKeyPEM := new(bytes.Buffer)
	pem.Encode(caPubKeyPEM, &pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: x509.MarshalPKCS1PublicKey(caPubKey),
	})
	return caPubKeyPEM
}

// LoadPublicKeyFile - loads a private key PEM file
func LoadPublicKeyFile(fileName string) []byte {
	inFile, err := ioutil.ReadFile(fileName)
	check(err)
	return inFile
}

// LoadPrivateKeyFile - loads a private key PEM file
func LoadPrivateKeyFile(fileName string) []byte {
	inFile, err := ioutil.ReadFile(fileName)
	check(err)
	return inFile
}

// DecodePublicKeyPem from file to pem struct
func DecodePublicKeyPem(inFile []byte) (*pem.Block, []byte) {
	pubPem, _ := pem.Decode(inFile)
	pubPemBytes := pubPem.Bytes
	return pubPem, pubPemBytes
}

// DecodePrivateKeyPem from file to pem struct
func DecodePrivateKeyPem(inFile []byte, rsaPrivateKeyPassword string) (*pem.Block, []byte) {
	privPem, _ := pem.Decode(inFile)
	if privPem.Type == "RSA PRIVATE KEY" {
		privPemBytes := privPem.Bytes

		if rsaPrivateKeyPassword != "" {
			privPemBytes, _ = x509.DecryptPEMBlock(privPem, []byte(rsaPrivateKeyPassword))
		} else {
			privPemBytes = privPem.Bytes
		}

		return privPem, privPemBytes
	}
	return nil, nil
}

// parsePrivateKey decodes a key from a pem
func parsePrivateKey(pemBytes []byte) *rsa.PrivateKey {
	parsedKey, err := x509.ParsePKCS1PrivateKey(pemBytes)
	check(err)
	return parsedKey
}

// parsePublicKey decodes a key from a pem
func parsePublicKey(pemBytes []byte) *rsa.PublicKey {
	parsedKey, err := x509.ParsePKCS1PublicKey(pemBytes)
	check(err)
	return parsedKey
}

// GetPrivateKey gets a private key soup to nuts
func GetPrivateKey(path string, rsaPrivateKeyPassword string) *rsa.PrivateKey {
	fileCheck, err := FileExists(path)
	check(err)
	if fileCheck {
		keyBytes := LoadPrivateKeyFile(path)
		_, keyPem := DecodePrivateKeyPem(keyBytes, rsaPrivateKeyPassword)
		return parsePrivateKey(keyPem)
	}
	return nil
}

// GetPublicKey gets a public key soup to nuts
func GetPublicKey(path string) *rsa.PublicKey {
	fileCheck, err := FileExists(path)
	check(err)
	if fileCheck {
		keyBytes := LoadPublicKeyFile(path)
		_, keyPem := DecodePublicKeyPem(keyBytes)
		return parsePublicKey(keyPem)
	}
	return nil
}
