package main

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/asn1"
	"encoding/pem"
	"errors"
	"fmt"
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
	keyFile, err := WriteByteFile(path, pemByte, 0400, false)
	if err != nil {
		return false, err
	}
	return keyFile, nil
}

// pemEncodeRSAPrivateKey
func pemEncodeRSAPrivateKey(caPrivKey *rsa.PrivateKey, rsaPrivateKeyPassword string) *bytes.Buffer {
	caPrivKeyPEM := new(bytes.Buffer)

	privateKeyBlock := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(caPrivKey),
	}

	if rsaPrivateKeyPassword != "" {
		privateKeyBlock, _ = x509.EncryptPEMBlock(rand.Reader, privateKeyBlock.Type, privateKeyBlock.Bytes, []byte(rsaPrivateKeyPassword), x509.PEMCipherAES256)
	}

	pem.Encode(caPrivKeyPEM, privateKeyBlock)
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

// marshallPublicKey converts a key into byte types
func marshalPublicKey(pub interface{}) (publicKeyBytes []byte, publicKeyAlgorithm pkix.AlgorithmIdentifier, err error) {
	switch pub := pub.(type) {
	case *rsa.PublicKey:
		publicKeyBytes, err = asn1.Marshal(pkcs1PublicKey{
			N: pub.N,
			E: pub.E,
		})
		if err != nil {
			return nil, pkix.AlgorithmIdentifier{}, err
		}
		publicKeyAlgorithm.Algorithm = oidPublicKeyRSA
		// This is a NULL parameters value which is required by
		// RFC 3279, Section 2.3.1.
		publicKeyAlgorithm.Parameters = asn1.NullRawValue
	case *ecdsa.PublicKey:
		publicKeyBytes = elliptic.Marshal(pub.Curve, pub.X, pub.Y)
		oid, ok := oidFromNamedCurve(pub.Curve)
		if !ok {
			return nil, pkix.AlgorithmIdentifier{}, errors.New("x509: unsupported elliptic curve")
		}
		publicKeyAlgorithm.Algorithm = oidPublicKeyECDSA
		var paramBytes []byte
		paramBytes, err = asn1.Marshal(oid)
		if err != nil {
			return
		}
		publicKeyAlgorithm.Parameters.FullBytes = paramBytes
	case ed25519.PublicKey:
		publicKeyBytes = pub
		publicKeyAlgorithm.Algorithm = oidPublicKeyEd25519
	default:
		return nil, pkix.AlgorithmIdentifier{}, fmt.Errorf("x509: unsupported public key type: %T", pub)
	}

	return publicKeyBytes, publicKeyAlgorithm, nil
}
