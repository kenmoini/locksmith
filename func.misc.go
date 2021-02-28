package main

import (
	"log"
	"os"
	"path/filepath"

	"github.com/gosimple/slug"
	"gopkg.in/yaml.v2"
)

// slugger slugs a string
func slugger(textToSlug string) string {
	return slug.Make(textToSlug)
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

	log.Println("Preflight complete!")
}
