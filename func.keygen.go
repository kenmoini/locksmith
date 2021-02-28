package main

import (
	"crypto/rand"
	"crypto/rsa"
)

// generatePrivKey returns a private RSA key
func generatePrivKey() (*rsa.PrivateKey, error) {
	// create our private and public key
	privKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		return nil, err
	}
	return privKey, nil
}
