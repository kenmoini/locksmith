#!/bin/bash

# Exits on any error
set -e

####################################################################
## Generate PKI for Tests
####################################################################
##  This script allows for the easy creation of a PKI via OpenSSL.
##  You can use this to instantiate your own basic PKI quickly,
##  however this is used primarily to test the Golang generated
##  PKI against a known working PKI generated via OpenSSL

function checkForProgramAndExit() {
    command -v $1
    if [[ $? -eq 0 ]]; then
        printf '%-72s %-7s\n' $1 "PASSED!";
    else
        printf '%-72s %-7s\n' $1 "FAILED!";
        exit 1
    fi
}

checkForProgramAndExit openssl

####################################################################
## Set up variables

CUR_DIR=$(pwd)

OPENSSL_PKI_ROOT_DIR="${CUR_DIR}/.test_pki_root"
LOCKSMITH_PKI_ROOT_DIR="${CUR_DIR}/.generated/roots/example-labs-root-certificate-authority"

echo -e "\n===== ISSUER COMPARISON\n"
echo "OSSL: $(openssl x509 -in ${OPENSSL_PKI_ROOT_DIR}/ca.cert.pem -noout -issuer)"
echo "Lock: $(openssl x509 -in ${LOCKSMITH_PKI_ROOT_DIR}/certs/ca.pem -noout -issuer)"

echo -e "\n===== SUBJECT COMPARISON\n"
echo "OSSL: $(openssl x509 -in ${OPENSSL_PKI_ROOT_DIR}/ca.cert.pem -noout -subject)"
echo "Lock: $(openssl x509 -in ${LOCKSMITH_PKI_ROOT_DIR}/certs/ca.pem -noout -subject)"