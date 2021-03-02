#!/bin/bash

####################################################################
## Generate PKI for Tests
####################################################################
##  This script allows for the easy creation of a PKI via OpenSSL.
##  You can use this to instantiate your own basic PKI quickly,
##  however this is used primarily to test the Golang generated
##  PKI against a known working PKI generated via OpenSSL

####################################################################
## Set up variables

CUR_DIR=$(pwd)

PKI_ROOT_DIR="${CUR_DIR}/.test_pki_root"

# Clean PKI Root
echo -e "Cleaning PKI Root..."
rm -rf $PKI_ROOT_DIR

# Create PKI Root Directory
echo -e "Creating PKI Root paths..."
mkdir -p $PKI_ROOT_DIR/{certreqs,certs,crl,newcerts,private,intermed-ca} && chmod 700 $PKI_ROOT_DIR/private

# Create initial index db
echo -e "Creating PKI Root CA Index DB file..."
touch $PKI_ROOT_DIR/ca.index

# Create initial serial index
echo -e "Creating PKI Root CA Serial Number file..."
echo "01" > $PKI_ROOT_DIR/ca.serial

#  Create initial crlnum index
echo -e "Creating PKI Root CA CRL Number file..."
echo "00" > $PKI_ROOT_DIR/ca.crlnum

# Set global Root CA OpenSSL Configuration
echo -e "\nSetting OpenSSL Configuration env var for Root CA...\n"
cp "${CUR_DIR}/openssl_extras/root-ca.cnf" "${PKI_ROOT_DIR}/ca.cnf"
export OPENSSL_CONF="${PKI_ROOT_DIR}/ca.cnf"
cd $PKI_ROOT_DIR

# Create Root CA Private Key
echo -e "Creating PKI Root CA Private Key..."
openssl genrsa -out $PKI_ROOT_DIR/private/ca.key.pem 4096
chmod 0400 $PKI_ROOT_DIR/private/ca.key.pem

# Create Root CA CSR
echo -e "\nCreating PKI Root CA CSR..."
openssl req -new -batch -out $PKI_ROOT_DIR/certreqs/ca.req.pem -key $PKI_ROOT_DIR/private/ca.key.pem

# Create Root CA Certificate
echo -e "\nCreating PKI Root CA Self-Signed Certificate..."
openssl ca -selfsign -batch -in $PKI_ROOT_DIR/certreqs/ca.req.pem -out $PKI_ROOT_DIR/ca.cert -extensions root-ca_ext -startdate $(date -u -d "-1day" "+%y%m%d000000Z" 2>/dev/null || date -u -v "-1d" "+%y%m%d000000Z") -enddate `(date -u -d "+10years+1day" "+%y%m%d000000Z" 2>/dev/null || date -u -v "+10y" -v "+1d" "+%y%m%d000000Z")`

CERT_START_LINE_NUM=$(awk '/BEGIN CERTIFICATE/{ print NR; exit }' $PKI_ROOT_DIR/ca.cert)
CERT_END_LINE_NUM=$(awk '/END CERTIFICATE/{ print NR; exit }' $PKI_ROOT_DIR/ca.cert)

tail -n +$CERT_START_LINE_NUM $PKI_ROOT_DIR/ca.cert > $PKI_ROOT_DIR/ca.cert.pem

# Create Root CA Certificate Revocation List
echo -e "\nCreating PKI Root CA Certificate Revocation List..."
openssl ca -gencrl -out $PKI_ROOT_DIR/crl/ca.crl

echo -e "\nTest PKI Created!\n"