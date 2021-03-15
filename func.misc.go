package main

import (
	"bufio"
	"crypto/elliptic"
	"crypto/sha1"
	"encoding/asn1"
	"errors"
	"fmt"
	"math/big"
	"net"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"unicode"

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

// readSerialNumberAsIntAbs is a wrapper that converts the string serial number in a serial file to an int
func readSerialNumberAsIntAbs(path string) int {
	i, _ := strconv.Atoi(readSerialNumberAbs(path))
	return i
}

// readSerialNumberAsInt64Abs converts an int converted serial number to int64
func readSerialNumberAsInt64Abs(path string) int64 {
	return int64(readSerialNumberAsIntAbs(path))
}

// readSerialNumberAbs reads the serial.txt file out from an absolute path
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

// IncreaseSerialNumberAbs just updates a root CAs serial via absolute path to the serial file
func IncreaseSerialNumberAbs(path string) (bool, error) {
	logStdOut("incrementing " + path)
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

// isIA5String checks to ensure the string is ASCII
func isIA5String(s string) error {
	for _, r := range s {
		// Per RFC5280 "IA5String is limited to the set of ASCII characters"
		if r > unicode.MaxASCII {
			return fmt.Errorf("x509: %q cannot be encoded as an IA5String", s)
		}
	}
	return nil
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

func forEachIAN(extension []byte, callback func(tag int, data []byte) error) error {
	// RFC 5280, 4.2.1.6

	// SubjectAltName ::= GeneralNames
	//
	// GeneralNames ::= SEQUENCE SIZE (1..MAX) OF GeneralName
	//
	// GeneralName ::= CHOICE {
	//      otherName                       [0]     OtherName,
	//      rfc822Name                      [1]     IA5String,
	//      dNSName                         [2]     IA5String,
	//      x400Address                     [3]     ORAddress,
	//      directoryName                   [4]     Name,
	//      ediPartyName                    [5]     EDIPartyName,
	//      uniformResourceIdentifier       [6]     IA5String,
	//      iPAddress                       [7]     OCTET STRING,
	//      registeredID                    [8]     OBJECT IDENTIFIER }
	var seq asn1.RawValue
	rest, err := asn1.Unmarshal(extension, &seq)
	if err != nil {
		return err
	} else if len(rest) != 0 {
		return errors.New("x509: trailing data after X.509 extension")
	}
	if !seq.IsCompound || seq.Tag != 16 || seq.Class != 0 {
		return asn1.StructuralError{Msg: "bad SAN sequence"}
	}

	rest = seq.Bytes
	for len(rest) > 0 {
		var v asn1.RawValue
		rest, err = asn1.Unmarshal(rest, &v)
		if err != nil {
			return err
		}

		if err := callback(v.Tag, v.Bytes); err != nil {
			return err
		}
	}

	return nil
}

// domainToReverseLabels converts a textual domain name like foo.example.com to
// the list of labels in reverse order, e.g. ["com", "example", "foo"].
func domainToReverseLabels(domain string) (reverseLabels []string, ok bool) {
	for len(domain) > 0 {
		if i := strings.LastIndexByte(domain, '.'); i == -1 {
			reverseLabels = append(reverseLabels, domain)
			domain = ""
		} else {
			reverseLabels = append(reverseLabels, domain[i+1:])
			domain = domain[:i]
		}
	}

	if len(reverseLabels) > 0 && len(reverseLabels[0]) == 0 {
		// An empty label at the end indicates an absolute value.
		return nil, false
	}

	for _, label := range reverseLabels {
		if len(label) == 0 {
			// Empty labels are otherwise invalid.
			return nil, false
		}

		for _, c := range label {
			if c < 33 || c > 126 {
				// Invalid character.
				return nil, false
			}
		}
	}

	return reverseLabels, true
}

// marshalIANs marshals a list of addresses into a the contents of an X.509
// IssuerAlternativeName extension.
func marshalIANs(dnsNames, emailAddresses []string, ipAddresses []net.IP, uris []*url.URL) (derBytes []byte, err error) {
	var rawValues []asn1.RawValue
	for _, name := range dnsNames {
		if err := isIA5String(name); err != nil {
			return nil, err
		}
		rawValues = append(rawValues, asn1.RawValue{Tag: nameTypeDNS, Class: 2, Bytes: []byte(name)})
	}
	for _, email := range emailAddresses {
		if err := isIA5String(email); err != nil {
			return nil, err
		}
		rawValues = append(rawValues, asn1.RawValue{Tag: nameTypeEmail, Class: 2, Bytes: []byte(email)})
	}
	for _, rawIP := range ipAddresses {
		// If possible, we always want to encode IPv4 addresses in 4 bytes.
		ip := rawIP.To4()
		if ip == nil {
			ip = rawIP
		}
		rawValues = append(rawValues, asn1.RawValue{Tag: nameTypeIP, Class: 2, Bytes: ip})
	}
	for _, uri := range uris {
		uriStr := uri.String()
		if err := isIA5String(uriStr); err != nil {
			return nil, err
		}
		rawValues = append(rawValues, asn1.RawValue{Tag: nameTypeURI, Class: 2, Bytes: []byte(uriStr)})
	}
	return asn1.Marshal(rawValues)
}

// parseIANExtension parses Issuer Alternative Name extension data structs from a raw byte slice
func parseIANExtension(value []byte) (dnsNames, emailAddresses []string, ipAddresses []net.IP, uris []*url.URL, err error) {
	err = forEachIAN(value, func(tag int, data []byte) error {
		switch tag {
		case nameTypeEmail:
			email := string(data)
			if err := isIA5String(email); err != nil {
				return errors.New("x509: IAN rfc822Name is malformed")
			}
			emailAddresses = append(emailAddresses, email)
		case nameTypeDNS:
			name := string(data)
			if err := isIA5String(name); err != nil {
				return errors.New("x509: IAN dNSName is malformed")
			}
			dnsNames = append(dnsNames, string(name))
		case nameTypeURI:
			uriStr := string(data)
			if err := isIA5String(uriStr); err != nil {
				return errors.New("x509: IAN uniformResourceIdentifier is malformed")
			}
			uri, err := url.Parse(uriStr)
			if err != nil {
				return fmt.Errorf("x509: cannot parse URI %q: %s", uriStr, err)
			}
			if len(uri.Host) > 0 {
				if _, ok := domainToReverseLabels(uri.Host); !ok {
					return fmt.Errorf("x509: cannot parse URI %q: invalid domain", uriStr)
				}
			}
			uris = append(uris, uri)
		case nameTypeIP:
			switch len(data) {
			case net.IPv4len, net.IPv6len:
				ipAddresses = append(ipAddresses, data)
			default:
				return errors.New("x509: cannot parse IP address of length " + strconv.Itoa(len(data)))
			}
		}

		return nil
	})

	return
}

// bigIntHash creates a SHA1 bytes array from an int
func bigIntHash(n *big.Int) []byte {
	h := sha1.New()
	h.Write(n.Bytes())
	return h.Sum(nil)
}

// oidFromNamedCurve takes an EC and returns the ASN1 OID
func oidFromNamedCurve(curve elliptic.Curve) (asn1.ObjectIdentifier, bool) {
	switch curve {
	case elliptic.P224():
		return oidNamedCurveP224, true
	case elliptic.P256():
		return oidNamedCurveP256, true
	case elliptic.P384():
		return oidNamedCurveP384, true
	case elliptic.P521():
		return oidNamedCurveP521, true
	}

	return nil, false
}

// splitSlugToPath takes a slug string and splits it into the relative path
// eg converts "example-labs-root-certificate-authority/example-labs-ica/server-signing-ca" to "example-labs-root-certificate-authority/intermed-ca/example-labs-ica/intermed-ca/server-signing-ca/"
func splitSlugToPath(slug string) string {
	splitPath := strings.Split(strings.TrimSuffix(strings.TrimPrefix(strings.ToLower(slug), "/"), "/"), "/")
	var path string
	for i, part := range splitPath {
		path = path + part + "/"
		if i != (len(splitPath) - 1) {
			path = path + "intermed-ca/"
		}
	}
	return path
}

// splitCommonNamesToPath takes a CN string and splits it into the relative path while slugging
// eg, converts "Example Labs Root Certificate Authority/Example Labs ICA/Server Signing CA" to "example-labs-root-certificate-authority/intermed-ca/example-labs-ica/intermed-ca/server-signing-ca/"
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
