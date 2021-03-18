package locksmith

// createKeyStore creates a new key store
func createKeyStore(keyStoreName string) (bool, string, error) {
	sluggedKeyStoreName := slugger(keyStoreName)
	basePath := readConfig.Locksmith.PKIRoot + "/keystores/" + sluggedKeyStoreName

	// Check for key store
	keyStoreCheck, err := FileExists(basePath)
	check(err)

	if !keyStoreCheck {
		// if there is no key store, create one
		CreateDirectory(basePath)
		return true, sluggedKeyStoreName, nil
	}
	return false, sluggedKeyStoreName, Stoerr("Key store exists!")
}

// listKeyStores returns a list of key stores
func listKeyStores() []string {
	return DirectoryListingNames(readConfig.Locksmith.PKIRoot + "/keystores/")
}
