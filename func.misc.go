package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"github.com/gosimple/slug"
	"gopkg.in/yaml.v2"
)

// slugger slugs a string
func slugger(textToSlug string) string {
	return slug.Make(textToSlug)
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

// readSerialNumber reads the serial.txt file out
func readSerialNumber(rootSlug string) string {
	/*
		dat, err := ioutil.ReadFile(readConfig.Locksmith.PKIRoot + "/roots/" + rootSlug + "/serial.txt")
		check(err)

		return strings.TrimSuffix(string(dat), "\n")
	*/
	filePath, err := filepath.Abs(readConfig.Locksmith.PKIRoot + "/roots/" + rootSlug + "/serial.txt")
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

	rootCACertSerialFilePath, err := filepath.Abs(rootSlugPath + "/serial.txt")
	check(err)

	// Update serialFile

	serialFile, err := WriteFile(rootCACertSerialFilePath, fmt.Sprintf("%v", counter.currentValue()), 0600, true)
	check(err)
	//if serialFile {
	//logStdOut("Updated serial file")
	//}
	return serialFile, err
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

	logStdOut("Preflight complete!")
}
