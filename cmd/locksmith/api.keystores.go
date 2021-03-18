package locksmith

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// listKeyStoresAPI returns the key stores for GET /keystores requests
func listKeyStoresAPI(w http.ResponseWriter, r *http.Request) {
	keyPairs := listKeyStores()

	returnData := &RESTGETKeyStoresJSONReturn{
		Status:    "success",
		Errors:    []string{},
		Messages:  []string{"Listings of Key Stores"},
		KeyStores: keyPairs}
	returnResponse, _ := json.Marshal(returnData)
	fmt.Fprintf(w, string(returnResponse))
}

// createKeyStoreAPI handles the requests for POST /keystores to create new key stores
func createKeyStoreAPI(w http.ResponseWriter, r *http.Request) {

	keyStoreInfo := RESTPOSTKeyStoresJSONIn{}
	err := json.NewDecoder(r.Body).Decode(&keyStoreInfo)
	check(err)

	if keyStoreInfo.KeyStore != "" {
		keyStoreCreated, keyStoreSlug, err := createKeyStore(keyStoreInfo.KeyStore)
		if keyStoreCreated {
			// Key store was created
			returnData := &RESTPOSTKeyStoresJSONReturn{
				Status:   "success",
				Errors:   []string{},
				Messages: []string{"Key Store successfully created!"},
				KeyStore: slugger(keyStoreSlug)}
			returnResponse, _ := json.Marshal(returnData)
			fmt.Fprintf(w, string(returnResponse))
		} else {
			// Key store creation failed
			returnData := &ReturnGenericMessage{
				Status:   "key-store-creation-failed",
				Errors:   []string{err.Error()},
				Messages: []string{}}
			returnResponse, _ := json.Marshal(returnData)
			fmt.Fprintf(w, string(returnResponse))
		}
	} else {
		// Key store name not found
		returnData := &ReturnGenericMessage{
			Status:   "key-store-name-missing",
			Errors:   []string{"Key Store name parameter missing!  Pass with `key_store_name`"},
			Messages: []string{}}
		returnResponse, _ := json.Marshal(returnData)
		fmt.Fprintf(w, string(returnResponse))
	}

}
