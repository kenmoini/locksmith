package main

import (
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

// NewConfig returns a new decoded Config struct
func NewConfig(configPath string) (*Config, error) {
	// Create config structure
	config := &Config{}

	// Open config file
	file, err := os.Open(configPath)
	if err != nil {
		return nil, err
	}
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

	// Create generation directory
	CreateDirectory(readConfig.Locksmith.PKIRoot)

	log.Println("Preflight complete!")
}

// Func main should be as small as possible and do as little as possible by convention
func main() {
	// Generate our config based on the config supplied
	// by the user in the flags
	cfgPath, err := ParseFlags()
	if err != nil {
		log.Fatal(err)
	}
	cfg, err := NewConfig(cfgPath)
	if err != nil {
		log.Fatal(err)
	}

	// Run preflight
	PreflightSetup()

	// Run the server
	cfg.RunHTTPServer()
}
