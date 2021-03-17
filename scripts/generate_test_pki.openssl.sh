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

PKI_ROOT_DIR="${CUR_DIR}/.test_pki_root"

echo -e "\nCleaning PKI Root..."
rm -rf $PKI_ROOT_DIR

echo -e "\n####################################################################"
echo -e "## START CREATING ROOT CA...\n"

echo -e "########### Creating PKI Root paths..."
mkdir -p $PKI_ROOT_DIR/{certreqs,certs,crl,newcerts,private,intermed-ca} && chmod 700 $PKI_ROOT_DIR/private

echo -e "########### Creating PKI Root CA Index DB file..."
touch $PKI_ROOT_DIR/ca.index

echo -e "########### Creating PKI Root CA Serial Number file..."
echo "01" > $PKI_ROOT_DIR/ca.serial

echo -e "########### Creating PKI Root CA CRL Number file..."
echo "00" > $PKI_ROOT_DIR/ca.crlnum

# Set global Root CA OpenSSL Configuration
echo -e "\n## Setting OpenSSL Configuration env var for Root CA...\n"
cp "${CUR_DIR}/openssl_extras/root-ca.cnf" "${PKI_ROOT_DIR}/ca.cnf"
export OPENSSL_CONF="${PKI_ROOT_DIR}/ca.cnf"
cd $PKI_ROOT_DIR

echo -e "########### Creating PKI Root CA Private Key..."
openssl genrsa -out $PKI_ROOT_DIR/private/ca.key.pem 4096 &>/dev/null
chmod 0400 $PKI_ROOT_DIR/private/ca.key.pem

echo -e "########### Creating PKI Root CA CSR..."
openssl req -new -batch -out $PKI_ROOT_DIR/certreqs/ca.req.pem -key $PKI_ROOT_DIR/private/ca.key.pem

echo -e "########### Creating PKI Root CA Self-Signed Certificate..."
openssl ca -selfsign -batch -in $PKI_ROOT_DIR/certreqs/ca.req.pem -out $PKI_ROOT_DIR/ca.cert -extensions root-ca_ext -startdate $(date -u -d "-1day" "+%y%m%d000000Z" 2>/dev/null || date -u -v "-1d" "+%y%m%d000000Z") -enddate `(date -u -d "+10years+1day" "+%y%m%d000000Z" 2>/dev/null || date -u -v "+10y" -v "+1d" "+%y%m%d000000Z")`

CERT_START_LINE_NUM=$(awk '/BEGIN CERTIFICATE/{ print NR; exit }' $PKI_ROOT_DIR/ca.cert)
CERT_END_LINE_NUM=$(awk '/END CERTIFICATE/{ print NR; exit }' $PKI_ROOT_DIR/ca.cert)

echo -e "########### Creating PKI Root CA Self-Signed Certificate PEM file..."
tail -n +$CERT_START_LINE_NUM $PKI_ROOT_DIR/ca.cert > $PKI_ROOT_DIR/ca.cert.pem

echo -e "########### Creating PKI Root CA Certificate Revocation List..."
openssl ca -gencrl -out $PKI_ROOT_DIR/crl/ca.crl

echo -e "\n####################################################################"
echo -e "########### FINISHED CREATING ROOT CA!"

####################################################################
## Setup Intermediate Certificate Authority

echo -e "\n####################################################################"
echo -e "########### START CREATING INTERMEDIATE CA...\n"

INTERMED_CA_PKI_ROOT_DIR="${PKI_ROOT_DIR}/intermed-ca"

echo -e "########### Creating PKI Intermediate CA Certificate Path..."
mkdir -p $INTERMED_CA_PKI_ROOT_DIR/{certreqs,certs,crl,newcerts,private,intermed-ca} && chmod 700 $INTERMED_CA_PKI_ROOT_DIR/private

echo -e "########### Creating PKI Intermediate CA Index DB file..."
touch $INTERMED_CA_PKI_ROOT_DIR/ca.index

echo -e "########### Creating PKI Intermediate CA Serial Number file..."
echo "01" > $INTERMED_CA_PKI_ROOT_DIR/ca.serial

echo -e "########### Creating PKI Intermediate CA CRL Number file..."
echo "00" > $INTERMED_CA_PKI_ROOT_DIR/ca.crlnum

# Set Intermediate CA OpenSSL Configuration
echo -e "\n## Setting OpenSSL Configuration env var for Intermediate CA...\n"
cp "${CUR_DIR}/openssl_extras/intermed-ca.cnf" "${INTERMED_CA_PKI_ROOT_DIR}/ca.cnf"
export OPENSSL_CONF="${INTERMED_CA_PKI_ROOT_DIR}/ca.cnf"
cd $INTERMED_CA_PKI_ROOT_DIR

echo -e "########### Creating PKI Intermediate CA Private Key..."
openssl genrsa -out $INTERMED_CA_PKI_ROOT_DIR/private/ca.key.pem 4096 &>/dev/null
chmod 0400 $INTERMED_CA_PKI_ROOT_DIR/private/ca.key.pem

echo -e "########### Creating PKI Intermediate CA CSR..."
openssl req -new -batch -out $INTERMED_CA_PKI_ROOT_DIR/certreqs/ca.req.pem -key $INTERMED_CA_PKI_ROOT_DIR/private/ca.key.pem

####################################################################
## Create Intermediate CA Certificate
echo -e "\n## Setting OpenSSL Configuration env var for Root CA...\n"
export OPENSSL_CONF="${PKI_ROOT_DIR}/ca.cnf"
cd $PKI_ROOT_DIR

echo -e "########### Creating PKI Intermediate CA Certificate, signed by Root CA..."
openssl ca -batch -in $INTERMED_CA_PKI_ROOT_DIR/certreqs/ca.req.pem -out $INTERMED_CA_PKI_ROOT_DIR/ca.cert -extensions intermed-ca_ext -startdate $(date -u -d "-1day" "+%y%m%d000000Z" 2>/dev/null || date -u -v "-1d" "+%y%m%d000000Z") -enddate `(date -u -d "+3years+1day" "+%y%m%d000000Z" 2>/dev/null || date -u -v "+3y" -v "+1d" "+%y%m%d000000Z")`

CERT_START_LINE_NUM=$(awk '/BEGIN CERTIFICATE/{ print NR; exit }' $INTERMED_CA_PKI_ROOT_DIR/ca.cert)
CERT_END_LINE_NUM=$(awk '/END CERTIFICATE/{ print NR; exit }' $INTERMED_CA_PKI_ROOT_DIR/ca.cert)

echo -e "########### Creating PKI Intermediate CA Self-Signed Certificate PEM file..."
tail -n +$CERT_START_LINE_NUM $INTERMED_CA_PKI_ROOT_DIR/ca.cert > $INTERMED_CA_PKI_ROOT_DIR/ca.cert.pem

####################################################################
# Set Intermediate CA OpenSSL Configuration
echo -e "\n## Setting OpenSSL Configuration env var for Intermediate CA...\n"
export OPENSSL_CONF="${INTERMED_CA_PKI_ROOT_DIR}/ca.cnf"
cd $INTERMED_CA_PKI_ROOT_DIR

echo -e "########### Creating PKI Intermediate CA Certificate Revocation List..."
openssl ca -gencrl -out $INTERMED_CA_PKI_ROOT_DIR/crl/ca.crl

echo -e "\n####################################################################"
echo -e "########### FINISHED CREATING INTERMEDIATE CA!"

echo -e "\n####################################################################"
echo -e "########### START CREATING SERVER CERTIFICATE..."

echo -e "########### Creating Server Certificate Private Key..."
openssl genrsa -out $INTERMED_CA_PKI_ROOT_DIR/private/test.server.key.pem 4096 &>/dev/null
chmod 0400 $INTERMED_CA_PKI_ROOT_DIR/private/test.server.key.pem

####################################################################
# Set Intermediate CA OpenSSL Configuration
echo -e "\n## Setting OpenSSL Configuration env var for Server Certificate...\n"
cp "${CUR_DIR}/openssl_extras/server-cert.cnf" "${INTERMED_CA_PKI_ROOT_DIR}/server.cnf"
export OPENSSL_CONF="${INTERMED_CA_PKI_ROOT_DIR}/server.cnf"
cd $INTERMED_CA_PKI_ROOT_DIR

echo -e "########### Creating Server Certificate Request..."
openssl req -new -batch -out $INTERMED_CA_PKI_ROOT_DIR/certreqs/test.server.req.pem -key $INTERMED_CA_PKI_ROOT_DIR/private/test.server.key.pem

echo -e "########### Creating Server Certificate signed with Intermediate CA..."
openssl ca -batch -in $INTERMED_CA_PKI_ROOT_DIR/certreqs/test.server.req.pem -out $INTERMED_CA_PKI_ROOT_DIR/certs/test.server.cert -extensions server_ext

CERT_START_LINE_NUM=$(awk '/BEGIN CERTIFICATE/{ print NR; exit }' $INTERMED_CA_PKI_ROOT_DIR/certs/test.server.cert)
CERT_END_LINE_NUM=$(awk '/END CERTIFICATE/{ print NR; exit }' $INTERMED_CA_PKI_ROOT_DIR/certs/test.server.cert)

echo -e "########### Creating Server Certificate PEM file..."
tail -n +$CERT_START_LINE_NUM $INTERMED_CA_PKI_ROOT_DIR/certs/test.server.cert > $INTERMED_CA_PKI_ROOT_DIR/certs/test.server.cert.pem

echo -e "########### Creating Server DH Params..."
openssl dhparam -out $INTERMED_CA_PKI_ROOT_DIR/private/test.server.dhparams-1024.pem 1024 &>/dev/null

echo -e "\n####################################################################"
echo -e "########### FINISHED CREATING SERVER CERTIFICATE!"

echo -e "\n####################################################################"
echo -e "## Creating Certificate Bundle File..."
cat ${PKI_ROOT_DIR}/ca.cert.pem > ${INTERMED_CA_PKI_ROOT_DIR}/ca-bundle.cert.pem
cat ${INTERMED_CA_PKI_ROOT_DIR}/ca.cert.pem >> ${INTERMED_CA_PKI_ROOT_DIR}/ca-bundle.cert.pem

echo -e "\n####################################################################"
echo -e "########### Test PKI Created!"
echo -e "####################################################################\n"