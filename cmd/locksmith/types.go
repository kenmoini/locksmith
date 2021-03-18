package locksmith

import (
	"crypto/x509"
	"math/big"
	"net"
	"time"
)

// errorString is a trivial implementation of error.
type errorString struct {
	s string
}

func (e *errorString) Error() string {
	return e.s
}

// Config struct for webapp config
type Config struct {
	Locksmith LocksmithYaml `yaml:"locksmith"`
}

// LocksmithYaml is what is defined for this Locksmith server
type LocksmithYaml struct {
	PKIRoot string `yaml:"pki_root"`
	Server  Server `yaml:"server"`
}

// Server configures the HTTP server
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
	Slug     string           `json:"slug"`
	Serial   string           `json:"next_serial"`
	CertInfo x509.Certificate `json:"certificate"`
}

// CertificateInformation gives a general read out of a certificate file
type CertificateInformation struct {
	CommonName     string `json:"common_name"`
	StartDate      string `json:"start_date"`
	ExpirationDate string `json:"expiration_date"`
}

// Counter for serial number
type Counter struct {
	count int64
}

// RESTGETIntermedCAJSONIn handles the data required by the GET /intermediates endpoint
type RESTGETIntermedCAJSONIn struct {
	CommonNamePath string `json:"parent_cn_path,omitempty"`
	SlugPath       string `json:"parent_slug_path,omitempty"`
}

// RESTGETIntermedCAJSONReturn handles the data returned by the GET /intermediates endpoint
type RESTGETIntermedCAJSONReturn struct {
	Status          string   `json:"status"`
	Errors          []string `json:"errors"`
	Messages        []string `json:"messages"`
	IntermediateCAs []string `json:"intermediate_certificate_authorities"`
}

// RESTGETKeyPairsJSONReturn handles the data returned by the GET /keys endpoint for key pair listings
type RESTGETKeyPairsJSONReturn struct {
	Status   string   `json:"status"`
	Errors   []string `json:"errors"`
	Messages []string `json:"messages"`
	KeyPairs []string `json:"key_pairs,omitempty"`
}

// RESTGETKeyStoresJSONReturn handles the data returned by the GET /keystores endpoint for key store listings
type RESTGETKeyStoresJSONReturn struct {
	Status    string   `json:"status"`
	Errors    []string `json:"errors"`
	Messages  []string `json:"messages"`
	KeyStores []string `json:"key_stores,omitempty"`
}

// RESTPOSTKeyStoresJSONReturn handles the data returned by the GET /keystores endpoint for key store listings
type RESTPOSTKeyStoresJSONReturn struct {
	Status   string   `json:"status"`
	Errors   []string `json:"errors"`
	Messages []string `json:"messages"`
	KeyStore string   `json:"key_store_id"`
}

// RESTPOSTKeyStoresJSONIn handles the data returned by the GET /keystores endpoint for key store listings
type RESTPOSTKeyStoresJSONIn struct {
	KeyStore string `json:"key_store_name"`
}

// RESTGETKeyPairJSONReturn handles the data returned by the GET /keys endpoint for specific key pair id data
type RESTGETKeyPairJSONReturn struct {
	Status   string   `json:"status"`
	Errors   []string `json:"errors"`
	Messages []string `json:"messages"`
	KeyPair  KeyPair  `json:"key_pair,omitempty"`
}

// KeyPair combines a string for a Public and Private Key Base64 PEM
type KeyPair struct {
	PublicKey  string `json:"public_key,omitempty"`
	PrivateKey string `json:"private_key,omitempty"`
}

// RESTPOSTIntermedCAJSONIn handles the data required by the POST /intermediates endpoint
type RESTPOSTIntermedCAJSONIn struct {
	CommonNamePath              string                   `json:"parent_cn_path,omitempty"`
	SlugPath                    string                   `json:"parent_slug_path,omitempty"`
	CertificateConfiguration    CertificateConfiguration `json:"certificate_config"`
	SigningPrivateKeyPassphrase string                   `json:"rsa_private_key_passphrase,omitempty"`
}

// CertificateConfiguration is a struct to pass Certificate Config Information into the setup functions
type CertificateConfiguration struct {
	Subject                 CertificateConfigurationSubject `json:"subject"`
	ExpirationDate          []int                           `json:"expiration_date,omitempty"`
	RSAPrivateKeyPassphrase string                          `json:"rsa_private_key_passphrase,omitempty"`
	SerialNumber            string                          `json:"serial_number,omitempty"`
	SANData                 SANData                         `json:"san_data,omitempty"`
}

// CertificateConfigurationSubject is simply a redefinition of pkix.Name
type CertificateConfigurationSubject struct {
	CommonName         string   `json:"common_name"`
	Organization       []string `json:"organization"`
	OrganizationalUnit []string `json:"organizational_unit,omitempty"`
	Country            []string `json:"country,omitempty"`
	Province           []string `json:"province,omitempty"`
	Locality           []string `json:"locality,omitempty"`
	StreetAddress      []string `json:"street_address,omitempty"`
	PostalCode         []string `json:"postal_code,omitempty"`
}

// basicConstraints is idk, something
type basicConstraints struct {
	IsCA       bool `asn1:"optional"`
	MaxPathLen int  `asn1:"optional,default:-1"`
}

// CertificateAuthorityPaths returns all the default paths generated by a new CA
type CertificateAuthorityPaths struct {
	RootCAPath               string
	RootCACertRequestsPath   string
	RootCACertsPath          string
	RootCACertRevListPath    string
	RootCANewCertsPath       string
	RootCACertKeysPath       string
	RootCAIntermediateCAPath string
	RootCACertIndexFilePath  string
	RootCACertSerialFilePath string
	RootCACrlnumFilePath     string
}

// SANData provides a collection of SANData for a certificate
type SANData struct {
	IPAddresses    []net.IP `json:"ip_addresses,omitempty"`
	EmailAddresses []string `json:"email_addresses,omitempty"`
	DNSNames       []string `json:"dns_names,omitempty"`
	URIs           []string `json:"uris,omitempty"`
	//URIs           []*url.URL `json:"uris,omitempty"`
}

// pkcs1PublicKey reflects the ASN.1 structure of a PKCS #1 public key.
type pkcs1PublicKey struct {
	N *big.Int
	E int
}

// CAIndex provides the tab-delimited structure for CA Index files
type CAIndex struct {
	State             string
	EndDate           string
	DateOfRevokation  string
	Serial            string
	PathToCertificate string
	Subject           string
}

// RESTPOSTNewKeyPairIn organizes the data required for creating a new Key Pair
type RESTPOSTNewKeyPairIn struct {
	KeyPairID       string `json:"key_pair_id"`
	KeyStoreID      string `json:"key_store_id,omitempty"`
	Passphrase      string `json:"passphrase,omitempty"`
	StorePrivateKey bool   `json:"store_private_key"`
}
