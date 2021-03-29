package locksmith

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"crypto/sha1"
	"encoding/hex"
	"io"
	"log"
)

// passphraseToHash returns a hexadecimal string of an SHA1 checksumed passphrase
func passphraseToHash(pass string) (string, []byte) {
	// The salt is used as a unique string to defeat rainbow table attacks
	//salt := "l0ckSmithS41t"

	// This days actually it's pretty easy to make a rainbow table with a static salt - compute salt from another hash of the password
	saltHash := md5.New()
	saltHash.Write([]byte(pass))
	saltyBytes := saltHash.Sum(nil)
	salt := hex.EncodeToString(saltyBytes)

	saltyPass := []byte(pass + salt)
	hasher := sha1.New()
	hasher.Write(saltyPass)

	hash := hasher.Sum(nil)

	return hex.EncodeToString(hash), hash
}

// encryptBytes is a wrapper that takes a plain byte slice and a passphrase and returns an encrypted byte slice
func encryptBytes(bytesIn []byte, passphrase string) []byte {
	passHash, _ := passphraseToHash(passphrase)
	targetPassHash := passHash[0:32]

	// Create an AES Cipher
	block, err := aes.NewCipher([]byte(targetPassHash))
	check(err)

	// Create a new gcm block container
	gcm, err := cipher.NewGCM(block)
	check(err)

	// Never use more than 2^32 random nonces with a given key because of the risk of repeat.
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		log.Fatal(err)
	}

	// Seal will encrypt the file using the GCM mode, appending the nonce and tag (MAC value) to the final data, so we can use it to decrypt it later.
	return gcm.Seal(nonce, nonce, bytesIn, nil)
}

// decryptBytes takes in a byte slice from a file and a passphrase then returns if the encrypted byte slice was decrypted, if so the plaintext contents, and any errors
func decryptBytes(bytesIn []byte, passphrase string) (decrypted bool, plaintextBytes []byte, err error) {
	// bytesIn must be decoded from base 64 first
	// b64.StdEncoding.DecodeString(bytesIn)

	passHash, _ := passphraseToHash(passphrase)
	targetPassHash := passHash[0:32]

	// Create an AES Cipher
	block, err := aes.NewCipher([]byte(targetPassHash))
	if err != nil {
		log.Panic(err)
		return false, []byte{}, err
	}

	// Create a new gcm block container
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		log.Panic(err)
		return false, []byte{}, err
	}

	nonce := bytesIn[:gcm.NonceSize()]
	ciphertext := bytesIn[gcm.NonceSize():]
	plaintextBytes, err = gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		log.Panic(err)
		return false, []byte{}, err
	}

	// successfully decrypted
	return true, plaintextBytes, nil
}

// isPEMEncrypted Checks to see if the byte slice from a file contains a plain-text PEM file
func isPEMEncrypted(bytesIn []byte, typeOfPEM string) bool {
	return !bytes.Contains(bytesIn, []byte(typeOfPEM))
}

// isPrivateKeyEncrypted Checks to see if the byte slice from a file contains a plain-text Private PEM Key
func isPrivateKeyEncrypted(bytesIn []byte) bool {
	return !bytes.Contains(bytesIn, []byte("PRIVATE KEY"))
}

// isCertificateEncrypted Checks to see if the byte slice from a file contains a plain-text Certificate PEM file
func isCertificateEncrypted(bytesIn []byte) bool {
	return !bytes.Contains(bytesIn, []byte("CERTIFICATE"))
}

// isCertificateRequestEncrypted Checks to see if the byte slice from a file contains a plain-text Certificate Reques PEM file
func isCertificateRequestEncrypted(bytesIn []byte) bool {
	return !bytes.Contains(bytesIn, []byte("CERTIFICATE REQUEST"))
}
