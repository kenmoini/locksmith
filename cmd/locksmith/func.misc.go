package locksmith

import (
	"bufio"
	"errors"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/gosimple/slug"
	"gopkg.in/yaml.v2"
)

// slugger slugs a string
func slugger(textToSlug string) string {
	return slug.Make(textToSlug)
}

// currentValue returns the current count of the Counter
func (counter Counter) currentValue() int64 {
	return counter.count
}

// increment increases an int/int64 value in a Counter
func (counter *Counter) increment() {
	counter.count++
}

// readSerialNumberAsInt is a wrapper that converts the string serial number in a serial file to an int
func readSerialNumberAsInt(rootSlugPath string) int {
	i, _ := strconv.Atoi(readSerialNumber(rootSlugPath))
	return i
}

// readSerialNumberAsInt64 converts an int converted serial number to int64
func readSerialNumberAsInt64(rootSlugPath string) int64 {
	return int64(readSerialNumberAsInt(rootSlugPath))
}

// readSerialNumber reads the ca.serial file out
func readSerialNumber(rootSlug string) string {
	filePath, err := filepath.Abs(readConfig.Locksmith.PKIRoot + "/roots/" + rootSlug + "/ca.serial")
	check(err)
	file, err := os.Open(filePath)
	check(err)
	defer file.Close()

	s := bufio.NewScanner(file)
	var serial string
	for s.Scan() {
		serial = s.Text()
		break
	}
	return serial
}

// readSerialNumberAsIntAbs is a wrapper that converts the string serial number in a serial file to an int
func readSerialNumberAsIntAbs(path string) int {
	i, _ := strconv.Atoi(readSerialNumberAbs(path))
	return i
}

// readSerialNumberAsInt64Abs converts an int converted serial number to int64
func readSerialNumberAsInt64Abs(path string) int64 {
	return int64(readSerialNumberAsIntAbs(path))
}

// readSerialNumberAbs reads the ca.serial file out from an absolute path
func readSerialNumberAbs(path string) string {
	file, err := os.Open(path)
	check(err)
	defer file.Close()

	s := bufio.NewScanner(file)
	var serial string
	for s.Scan() {
		serial = s.Text()
		break
	}
	return serial
}

// IncreaseSerialNumber just updates a root CAs serial
func IncreaseSerialNumber(rootSlug string) (bool, error) {
	serNum := readSerialNumberAsInt64(rootSlug)

	//currentSerialNumString := readSerialNumber(rootSlug)
	//logStdOut("currentSerialNumString: " + currentSerialNumString)

	//currentSerialNumber, _ := strconv.Atoi(currentSerialNumString)
	//log.Printf("i=%d, type: %T\n", currentSerialNumber, currentSerialNumber)
	//serNum = int64(currentSerialNumber)

	counter := Counter{serNum}
	//log.Printf("i=%d, type: %T\n", counter.currentValue(), counter.currentValue())

	counter.increment()

	//log.Printf("i=%d, type: %T\n", counter.currentValue(), counter.currentValue())

	rootSlugPath := readConfig.Locksmith.PKIRoot + "/roots/" + rootSlug

	rootCACertSerialFilePath, err := filepath.Abs(rootSlugPath + "/ca.serial")
	check(err)

	// Update serialFile

	serialFile, err := WriteFile(rootCACertSerialFilePath, fmt.Sprintf("%v", counter.currentValue()), 0600, true)
	check(err)
	//if serialFile {
	//logStdOut("Updated serial file")
	//}
	return serialFile, err
}

// IncreaseSerialNumberAbs just updates a root CAs serial via absolute path to the serial file
func IncreaseSerialNumberAbs(path string) (bool, error) {
	//logStdOut("incrementing " + path)
	serNum := readSerialNumberAsInt64Abs(path)

	counter := Counter{serNum}
	counter.increment()

	return WriteFile(path, fmt.Sprintf("%v", counter.currentValue()), 0600, true)
}

// bakeURIs converts URL strings to actual URI slices
func bakeURIs(uris []string) ([]*url.URL, error) {
	actualURIs := []*url.URL{}
	for _, s := range uris {
		if err := isIA5String(s); err != nil {
			return nil, errors.New("x509: SAN uniformResourceIdentifier is malformed")
		}
		u, err := url.Parse(s)
		check(err)
		actualURIs = append(actualURIs, u)
	}
	return actualURIs, nil
}

// NewConfig returns a new decoded Config struct
func NewConfig(configPath string) (*Config, error) {
	// Create config structure
	config := &Config{}

	// Open config file
	file, err := os.Open(configPath)
	checkAndFail(err)
	defer file.Close()

	// Init new YAML decode
	d := yaml.NewDecoder(file)

	// Start YAML decoding from file
	if err := d.Decode(&config); err != nil {
		return nil, err
	}

	readConfig = config

	return config, nil
}

// PreflightSetup just makes sure the stage is set
func PreflightSetup() {

	// Create PKI Root directory
	PKIRootPath, err := filepath.Abs(readConfig.Locksmith.PKIRoot)
	checkAndFail(err)
	CreateDirectory(PKIRootPath)

	// Create PKI Root root directory
	PKIRootRootsPath, err := filepath.Abs(readConfig.Locksmith.PKIRoot + "/roots")
	checkAndFail(err)
	CreateDirectory(PKIRootRootsPath)

	// Create PKI Extra Keys directory
	PKIKeysRootsPath, err := filepath.Abs(readConfig.Locksmith.PKIRoot + "/keystores")
	checkAndFail(err)
	CreateDirectory(PKIKeysRootsPath)
	CreateDirectory(PKIKeysRootsPath + "/default")

	logStdOut("Preflight complete!")
}

/*

// splitSlugToPath takes a slug string and splits it into the relative path
// eg converts "example-labs-root-certificate-authority/example-labs-ica/signing-ca" to "example-labs-root-certificate-authority/intermed-ca/example-labs-ica/intermed-ca/signing-ca/"
func splitSlugToPath(slug string) string {
	splitPath := strings.Split(strings.TrimSuffix(strings.TrimPrefix(strings.ToLower(slug), "/"), "/"), "/")
	var path string
	for i, part := range splitPath {
		path = path + slugger(part) + "/"
		if i != (len(splitPath) - 1) {
			path = path + "intermed-ca/"
		}
	}
	return path
}

// splitCommonNamesToPath takes a CN string and splits it into the relative path while slugging
// eg, converts "Example Labs Root Certificate Authority/Example Labs ICA/Signing CA" to "example-labs-root-certificate-authority/intermed-ca/example-labs-ica/intermed-ca/signing-ca/"
func splitCommonNamesToPath(cnPath string) string {
	splitPath := strings.Split(strings.TrimSuffix(strings.TrimPrefix(strings.ToLower(cnPath), "/"), "/"), "/")
	var path string
	for i, part := range splitPath {
		path = path + slugger(part) + "/"
		if i != (len(splitPath) - 1) {
			path = path + "intermed-ca/"
		}
	}
	return path
}

*/

// splitCACNChainToPath takes a CN string and splits it into the relative path while slugging
// eg, converts "Example Labs Root Certificate Authority/Example Labs ICA/Signing CA" to "example-labs-root-certificate-authority/intermed-ca/example-labs-ica/intermed-ca/signing-ca/"
// or converts "example-labs-root-certificate-authority/example-labs-ica/signing-ca" to "example-labs-root-certificate-authority/intermed-ca/example-labs-ica/intermed-ca/signing-ca/"
// you can even mix and match - the parts are slugged regardless of input and produce the same result
func splitCACNChainToPath(cnPath string) string {
	splitPath := strings.Split(strings.TrimSuffix(strings.TrimPrefix(strings.ToLower(cnPath), "/"), "/"), "/")
	var path string
	for i, part := range splitPath {
		path = path + slugger(part) + "/"
		if i != (len(splitPath) - 1) {
			path = path + "intermed-ca/"
		}
	}
	return path
}

// rmStrFromStrSlice removes a string from a string slice
func rmStrFromStrSlice(r string, s []string) []string {
	for i, v := range s {
		if v == r {
			return append(s[:i], s[i+1:]...)
		}
	}
	return s
}
