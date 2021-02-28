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
	privKeyFile, err := writePrivateKey(privKey, path+".priv.key")
	if err != nil {
		return false, false, err
	}

	pubKeyFile, err := writePublicKey(pubKey, path+".pub.key")
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
func pemEncodeRSAPrivateKey(caPrivKey *rsa.PrivateKey) *bytes.Buffer {
	caPrivKeyPEM := new(bytes.Buffer)
	pem.Encode(caPrivKeyPEM, &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(caPrivKey),
	})
	return caPrivKeyPEM
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
