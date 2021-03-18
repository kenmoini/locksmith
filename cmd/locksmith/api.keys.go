package locksmith

import (
	b64 "encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"path/filepath"
)

// listKeyPairsAPI handles the GET /v1/keys endpoint
func listKeyPairsAPI(w http.ResponseWriter, r *http.Request) {
	var sluggedKeyStoreID string

	// Read in the submitted parameters
	queryParams := r.URL.Query()
	keyPairID, presentKPID := queryParams["key_pair_id"]
	keyStoreID, presentKSID := queryParams["key_store_id"]
	passphrase, presentPassphrase := queryParams["passphrase"]

	if presentKSID {
		sluggedKeyStoreID = slugger(keyStoreID[0])
	} else {
		sluggedKeyStoreID = "default"
	}

	keyStorePath, err := filepath.Abs(readConfig.Locksmith.PKIRoot + "/keystores/" + sluggedKeyStoreID + "/")
	check(err)

	checkKeyStorePath, err := DirectoryExists(keyStorePath)
	check(err)

	if !checkKeyStorePath {
		// No valid key store
		returnData := &ReturnGenericMessage{
			Status:   "invalid-key-store",
			Errors:   []string{"Invalid Key Store '" + sluggedKeyStoreID + "'!"},
			Messages: []string{}}
		returnResponse, _ := json.Marshal(returnData)
		fmt.Fprintf(w, string(returnResponse))
	} else {
		// Key store exists, proceed
		if presentKPID {
			pubKeyPath := readConfig.Locksmith.PKIRoot + "/keystores/" + sluggedKeyStoreID + "/" + slugger(keyPairID[0]) + "/rsa.pub.pem"
			if presentPassphrase {
				// Targeting a specific Key Pair ID for a Private and Public Key Pair
				//The Passphrase is not empty, open and decrypt the private key if the passphrase is valid
				privKeyPath := readConfig.Locksmith.PKIRoot + "/keystores/" + sluggedKeyStoreID + "/" + slugger(keyPairID[0]) + "/rsa.priv.pem"
				keyBytes := LoadKeyFile(readConfig.Locksmith.PKIRoot + "/keystores/" + sluggedKeyStoreID + "/" + slugger(keyPairID[0]) + "/rsa.pub.pem")

				fileCheck, err := FileExists(privKeyPath)
				check(err)
				if fileCheck {
					// Now check the password against the key
					privKeyBytes := LoadKeyFile(privKeyPath)
					if isPrivateKeyEncrypted(privKeyBytes) {
						// Test decoding
						decodedPrivKey, err := b64.StdEncoding.DecodeString(string(privKeyBytes))
						check(err)

						bit, byted, err := decryptBytes(decodedPrivKey, passphrase[0])
						check(err)
						if bit {
							// Successfully decoded, return bytes
							returnData := &RESTGETKeyPairJSONReturn{
								Status:   "success",
								Errors:   []string{},
								Messages: []string{"Loaded Key Pair ID '" + keyPairID[0] + "' (" + slugger(keyPairID[0]) + ") in Key Store '" + sluggedKeyStoreID + "'"},
								KeyPair:  KeyPair{PublicKey: b64.StdEncoding.EncodeToString(keyBytes), PrivateKey: b64.StdEncoding.EncodeToString(byted)}}
							returnResponse, _ := json.Marshal(returnData)
							fmt.Fprintf(w, string(returnResponse))

						} else {
							// Decryption failed for some reason
							returnData := &ReturnGenericMessage{
								Status:   "private-key-decryption-error",
								Errors:   []string{"Private Key decryption failed for Key Pair ID '" + keyPairID[0] + "' (" + slugger(keyPairID[0]) + ") in Key Store '" + sluggedKeyStoreID + "'!"},
								Messages: []string{}}
							returnResponse, _ := json.Marshal(returnData)
							fmt.Fprintf(w, string(returnResponse))
						}
					} else {
						// Plain text key, send it on through
						returnData := &RESTGETKeyPairJSONReturn{
							Status:   "success",
							Errors:   []string{},
							Messages: []string{"Loaded Key Pair ID '" + keyPairID[0] + "' (" + slugger(keyPairID[0]) + ") in Key Store '" + sluggedKeyStoreID + "'"},
							KeyPair:  KeyPair{PublicKey: b64.StdEncoding.EncodeToString(keyBytes), PrivateKey: b64.StdEncoding.EncodeToString(privKeyBytes)}}
						returnResponse, _ := json.Marshal(returnData)
						fmt.Fprintf(w, string(returnResponse))
					}
				} else {
					// Private Key does not exist
					returnData := &ReturnGenericMessage{
						Status:   "no-private-key",
						Errors:   []string{"No Private Key is stored for Key Pair ID '" + keyPairID[0] + "' (" + slugger(keyPairID[0]) + ") in Key Store '" + sluggedKeyStoreID + "'!"},
						Messages: []string{}}
					returnResponse, _ := json.Marshal(returnData)
					fmt.Fprintf(w, string(returnResponse))
				}

			} else {

				// Targeting a specific Key Pair ID for a Public Key
				fileCheck, err := FileExists(pubKeyPath)
				check(err)
				if fileCheck {
					keyBytes := LoadKeyFile(pubKeyPath)

					returnData := &RESTGETKeyPairJSONReturn{
						Status:   "success",
						Errors:   []string{},
						Messages: []string{"Public Key for Key Pair ID '" + keyPairID[0] + "' (" + slugger(keyPairID[0]) + ") in Key Store '" + sluggedKeyStoreID + "'"},
						KeyPair:  KeyPair{PublicKey: b64.StdEncoding.EncodeToString(keyBytes)}}
					returnResponse, _ := json.Marshal(returnData)
					fmt.Fprintf(w, string(returnResponse))
				} else {
					// Key Pair does not exist
					returnData := &ReturnGenericMessage{
						Status:   "invalid-key-pair-id",
						Errors:   []string{"Invalid Key Pair ID in Key Store '" + sluggedKeyStoreID + "'!"},
						Messages: []string{}}
					returnResponse, _ := json.Marshal(returnData)
					fmt.Fprintf(w, string(returnResponse))
				}

			}
		} else {
			keyPairs := DirectoryListingNames(readConfig.Locksmith.PKIRoot + "/keystores/" + sluggedKeyStoreID + "/")
			if len(keyPairs) > 0 {
				// Return list of key pair ids (dirs lol) in the key store
				returnData := &RESTGETKeyPairsJSONReturn{
					Status:   "success",
					Errors:   []string{},
					Messages: []string{"Listing of Key Pair IDs in Key Store '" + sluggedKeyStoreID + "'"},
					KeyPairs: keyPairs}
				returnResponse, _ := json.Marshal(returnData)
				fmt.Fprintf(w, string(returnResponse))
			} else {
				// Key store is empty
				returnData := &ReturnGenericMessage{
					Status:   "empty-key-store",
					Errors:   []string{"Key Store '" + sluggedKeyStoreID + "' is empty!"},
					Messages: []string{}}
				returnResponse, _ := json.Marshal(returnData)
				fmt.Fprintf(w, string(returnResponse))
			}
		}
	}
}

// createKeyPairAPI handles the POST /v1/keys endpoint
func createKeyPairAPI(w http.ResponseWriter, r *http.Request) {
	var sluggedKeyStoreID string

	keyPairInfo := RESTPOSTNewKeyPairIn{}
	err := json.NewDecoder(r.Body).Decode(&keyPairInfo)
	check(err)

	if keyPairInfo.KeyStoreID != "" {
		sluggedKeyStoreID = slugger(keyPairInfo.KeyStoreID)
	} else {
		sluggedKeyStoreID = "default"
	}

	if keyPairInfo.KeyPairID != "" {
		basePath := readConfig.Locksmith.PKIRoot + "/keystores/" + sluggedKeyStoreID + "/" + slugger(keyPairInfo.KeyPairID)
		pubKeyPath := basePath + "/rsa.pub.pem"

		// Check for certificate authority key pair
		keyCheck, err := FileExists(pubKeyPath)
		check(err)

		if !keyCheck {
			// if there is no private key, create one

			// Create the directory
			CreateDirectory(basePath)

			privKey, pubKey, err := GenerateRSAKeypair(4096)
			check(err)

			if keyPairInfo.StorePrivateKey {
				// Save the Private Key to the file system
				pemEncodedPrivateKey := pemEncodeRSAPrivateKey(privKey, keyPairInfo.Passphrase)
				privKeyFile, pubKeyFile, err := writeRSAKeyPair(pemEncodedPrivateKey, pemEncodeRSAPublicKey(pubKey), basePath+"/rsa")
				check(err)

				if !privKeyFile || !pubKeyFile {
					// Something messed up...
					returnData := &ReturnGenericMessage{
						Status:   "key-pair-generation-error",
						Errors:   []string{err.Error()},
						Messages: []string{"Key Pair generation error!"}}
					returnResponse, _ := json.Marshal(returnData)
					fmt.Fprintf(w, string(returnResponse))
				} else {
					// All clear - pass keys
					returnData := &RESTGETKeyPairJSONReturn{
						Status:   "success",
						Errors:   []string{},
						Messages: []string{"Successfully created Key Pair '" + slugger(keyPairInfo.KeyPairID) + "' in Key Store '" + sluggedKeyStoreID + "'!"},
						KeyPair:  KeyPair{PublicKey: b64.StdEncoding.EncodeToString(pemEncodeRSAPublicKey(pubKey).Bytes()), PrivateKey: b64.StdEncoding.EncodeToString(pemEncodedPrivateKey.Bytes())}}
					returnResponse, _ := json.Marshal(returnData)
					fmt.Fprintf(w, string(returnResponse))
				}

			} else {
				// Do NOT save the Private Key to the file system (the default case)
				pemEncodedPrivateKey := pemEncodeRSAPrivateKey(privKey, keyPairInfo.Passphrase)
				pubKeyFile, err := writeKeyFile(pemEncodeRSAPublicKey(pubKey), basePath+"/rsa.pub.pem", 0644)

				if !pubKeyFile {
					// Something messed up...
					returnData := &ReturnGenericMessage{
						Status:   "key-pair-generation-error",
						Errors:   []string{err.Error()},
						Messages: []string{"Key Pair generation error!"}}
					returnResponse, _ := json.Marshal(returnData)
					fmt.Fprintf(w, string(returnResponse))
				} else {
					// All clear - pass keys
					returnData := &RESTGETKeyPairJSONReturn{
						Status:   "success",
						Errors:   []string{},
						Messages: []string{"Successfully created Key Pair '" + slugger(keyPairInfo.KeyPairID) + "' in Key Store '" + sluggedKeyStoreID + "'!"},
						KeyPair:  KeyPair{PublicKey: b64.StdEncoding.EncodeToString(pemEncodeRSAPublicKey(pubKey).Bytes()), PrivateKey: b64.StdEncoding.EncodeToString(pemEncodedPrivateKey.Bytes())}}
					returnResponse, _ := json.Marshal(returnData)
					fmt.Fprintf(w, string(returnResponse))
				}

			}

		} else {
			// Key pair exists already
			returnData := &ReturnGenericMessage{
				Status:   "key-pair-exists",
				Errors:   []string{"Key Pair '" + slugger(keyPairInfo.KeyPairID) + "' in Key Store '" + sluggedKeyStoreID + "' already exists!"},
				Messages: []string{}}
			returnResponse, _ := json.Marshal(returnData)
			fmt.Fprintf(w, string(returnResponse))
		}
	}

}
