package main

import "time"

// Config struct for webapp config
type Config struct {
	Locksmith LocksmithYaml `yaml:"locksmith"`
}

// LocksmithYaml is what is defined for this Pilot Light server
type LocksmithYaml struct {
	PKIRoot string `yaml:"pki_root"`
	Server  Server `yaml:"server"`
}

// ReturnGenericMessage - Generic message
type ReturnGenericMessage struct {
	Status   string   `json:"status"`
	Errors   []string `json:"errors"`
	Messages []string `json:"messages"`
}

// ReturnGetRoots - GET /roots
type ReturnGetRoots struct {
	Status   string   `json:"status"`
	Errors   []string `json:"errors"`
	Messages []string `json:"messages"`
	Roots    []string `json:"roots"`
}

// ReturnPostRoots - GET /roots
type ReturnPostRoots struct {
	Status   string   `json:"status"`
	Errors   []string `json:"errors"`
	Messages []string `json:"messages"`
	Root     RootInfo `json:"root"`
}

// RootInfo provides general root informations
type RootInfo struct {
	Slug   string `json:"slug"`
	Serial string `json:"serial"`
}

// Server configures the HTTP server providing Ignition files
type Server struct {
	// Host is the local machine IP Address to bind the HTTP Server to
	Host string `yaml:"host"`

	BasePath string `yaml:"base_path"`

	// Port is the local machine TCP Port to bind the HTTP Server to
	Port    string `yaml:"port"`
	Timeout struct {
		// Server is the general server timeout to use
		// for graceful shutdowns
		Server time.Duration `yaml:"server"`

		// Write is the amount of time to wait until an HTTP server
		// write opperation is cancelled
		Write time.Duration `yaml:"write"`

		// Read is the amount of time to wait until an HTTP server
		// read operation is cancelled
		Read time.Duration `yaml:"read"`

		// Read is the amount of time to wait
		// until an IDLE HTTP session is closed
		Idle time.Duration `yaml:"idle"`
	} `yaml:"timeout"`
}

// errorString is a trivial implementation of error.
type errorString struct {
	s string
}

func (e *errorString) Error() string {
	return e.s
}
