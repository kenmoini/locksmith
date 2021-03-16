package main

import (
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/jszwec/csvutil"
)

/*
AddEntryToCAIndex adds the needed tab-separated data to the CA Index file when generating certificates
State: “V” for Valid, “E” for Expired and “R” for revoked
Enddate: in the format YYMMDDHHmmssZ (the “Z” stands for Zulu/GMT)
Date of Revocation: same format as “Enddate”
Serial: serial of the certificate
Path to Certificate: can also be “unknown”
Subject: subject of the certificate
*/
func AddEntryToCAIndex(indexPath string, certPath string, certificate *x509.Certificate) (bool, error) {
	// File is created on file structure initializtion - read the file in
	f, err := os.OpenFile(indexPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	check(err)

	defer f.Close()

	endT := certificate.NotAfter.UTC()
	endYear := fmt.Sprintf("%v", endT.Year())
	endYearShort := endYear[len(endYear)-2:]
	formattedDate := fmt.Sprintf("%s%02d%02d%02d%02d%02dZ", endYearShort, endT.Month(), endT.Day(), endT.Hour(), endT.Minute(), endT.Second())

	// create CAIndex struct
	caIndex := []CAIndex{{
		State:             "V",
		EndDate:           formattedDate,
		DateOfRevokation:  "",
		Serial:            fmt.Sprintf("%02d", certificate.SerialNumber),
		Subject:           compileSubjectString(certificate.Subject),
		PathToCertificate: certPath}}

	w := csv.NewWriter(f)
	w.Comma = '\t'

	enc := csvutil.NewEncoder(w)
	enc.AutoHeader = false

	if err := enc.Encode(caIndex); err != nil {
		log.Fatal(err)
	}

	w.Flush()
	if err := w.Error(); err != nil {
		log.Fatal(err)
	}

	return true, nil
}

// CheckCAIndexForExpiredCertificates just scans the CA Index and cycles through the lines checking the expiration date and setting E if expired

// RevokeCertInCAIndex adds the Date of Revocation and sets R for a cert based on Serial Number targeting

// NewTabDelimitedWriter just wraps an IO writer
func NewTabDelimitedWriter(w io.Writer) (writer *csv.Writer) {
	writer = csv.NewWriter(w)
	writer.Comma = '\t'

	return
}

func compileSubjectString(subject pkix.Name) (compiledString string) {
	if len(subject.Country) != 0 {
		compiledString = compiledString + "/C=" + subject.Country[0]
	}
	if len(subject.Province) != 0 {
		compiledString = compiledString + "/S=" + subject.Province[0]
	}
	if len(subject.Locality) != 0 {
		compiledString = compiledString + "/L=" + subject.Locality[0]
	}
	if len(subject.Organization) != 0 {
		compiledString = compiledString + "/O=" + subject.Organization[0]
	}
	if len(subject.OrganizationalUnit) != 0 {
		compiledString = compiledString + "/OU=" + subject.OrganizationalUnit[0]
	}
	if subject.CommonName != "" {
		compiledString = compiledString + "/CN=" + subject.CommonName
	}
	return
}
