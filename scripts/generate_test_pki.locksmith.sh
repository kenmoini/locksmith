#!/bin/bash

# Exits on any error
set -e

RED='\033[0;31m'
GREEN='\033[0;32m'
NC='\033[0m' # No Color

VERBOSITY="1" # 1=normal, 2=CMND echo

function checkStatus() {
  FUNC=$1
  STATUS=$2
  CHECK_NAME=$3
  FULL_FN_CHECK="$FUNC $CHECK_NAME"
  line='................................................................................................'

  COLOR=$RED
  if [[ "$STATUS" == "\"success\"" ]]; then
    COLOR=$GREEN
  fi
  printf "${COLOR}%s %s %s${NC}\n" "$FULL_FN_CHECK" "${line:${#FULL_FN_CHECK}}" $STATUS
}

CUR_DIR=$(pwd)

## Clean up the directory
rm -rf ${CUR_DIR}/.generated/roots/example-labs-root-certificate-authority
rm -rf ${CUR_DIR}/.generated/keystores/

## Run the Locksmith application in the background
nohup ${CUR_DIR}/dist/locksmith -config configs/config.yml.example &

RUN_PID=$!

# Wait a few seconds while the Locksmith server starts
sleep 5

CURL_GET_OPTS="-s"
CURL_POST_OPTS="-s --header \"Content-Type: application/json\" --request POST"

###################################################################################- KEY STORES
# Read the list of key stores
CMND=$(curl $CURL_GET_OPTS http://localhost:8080/locksmith/v1/keystores)
CMND_STATUS=$(echo "$CMND" | jq .status)
checkStatus "[KEY STORES][GET][LIST]" $CMND_STATUS "Listing of key store, only default expected"
if [[ $VERBOSITY == "2" ]]; then echo -e "${CMND}\n"; fi

# Create a new key stores
CMND=$(curl $CURL_POST_OPTS --data '{"key_store_name": "Example Labs"}' http://localhost:8080/locksmith/v1/keystore)
CMND_STATUS=$(echo "$CMND" | jq .status)
checkStatus "[KEY STORE][POST][CREATE]" $CMND_STATUS "Creation of 'Example Labs' Key Store"
if [[ $VERBOSITY == "2" ]]; then echo -e "${CMND}\n"; fi

# Read the list of key stores
CMND=$(curl $CURL_GET_OPTS http://localhost:8080/locksmith/v1/keystores)
CMND_STATUS=$(echo "$CMND" | jq .status)
checkStatus "[KEY STORES][GET][LIST]" $CMND_STATUS "Listing of Key Stores, default and example-labs expected"
if [[ $VERBOSITY == "2" ]]; then echo -e "${CMND}\n"; fi

echo ""

###################################################################################- KEY PAIRS
# Create MyKeyPair key pair in the default key store, no pass, no saving [works]
CMND=$(curl $CURL_POST_OPTS --data '{"key_pair_id": "MyKeyPair"}' http://localhost:8080/locksmith/v1/key)
CMND_STATUS=$(echo "$CMND" | jq .status)
checkStatus "[KEY PAIR][POST][CREATE]" $CMND_STATUS "Creating MyKeyPair in default Key Store"
if [[ $VERBOSITY == "2" ]]; then echo -e "${CMND}\n"; fi

# Create a Server Key Pair in the example-labs key store, store the private key plain text [works]
CMND=$(curl $CURL_POST_OPTS --data '{"key_pair_id": "Server Key Pair", "key_store_id": "example-labs", "store_private_key": true}' http://localhost:8080/locksmith/v1/key)
CMND_STATUS=$(echo "$CMND" | jq .status)
checkStatus "[KEY PAIR][POST][CREATE]" $CMND_STATUS "Creating MyKeyPair Key Pair and store in example-labs Key Store"
if [[ $VERBOSITY == "2" ]]; then echo -e "${CMND}\n"; fi

# Create a VDI Terminal Key Pair, in the default key store, and passphrase protect stored private pair [works]
CMND=$(curl $CURL_POST_OPTS --data '{"key_pair_id": "VDI Terminal", "store_private_key": true, "passphrase": "s3cr3t"}' http://localhost:8080/locksmith/v1/key)
CMND_STATUS=$(echo "$CMND" | jq .status)
checkStatus "[KEY PAIR][POST][CREATE]" $CMND_STATUS "Creating VDI Terminal Key Pair in default Key Store"
if [[ $VERBOSITY == "2" ]]; then echo -e "${CMND}\n"; fi

# Read the list of key pairs in the default store
CMND=$(curl $CURL_GET_OPTS http://localhost:8080/locksmith/v1/keys)
CMND_STATUS=$(echo "$CMND" | jq .status)
checkStatus "[KEY PAIRS][GET][LIST]" $CMND_STATUS "Listing Key Pairs in the default Key Store"
if [[ $VERBOSITY == "2" ]]; then echo -e "${CMND}\n"; fi

# Read the list of key pairs in the example-labs store
CMND=$(curl $CURL_GET_OPTS -G --data-urlencode "key_store_id=example-labs" http://localhost:8080/locksmith/v1/keys)
CMND_STATUS=$(echo "$CMND" | jq .status)
checkStatus "[KEY PAIRS][GET][LIST]" $CMND_STATUS "Listing Key Pairs in the example-labs Key Store"
if [[ $VERBOSITY == "2" ]]; then echo -e "${CMND}\n"; fi

# Read the public key of the MyKeyPair Key Pair in the default key store
CMND=$(curl $CURL_GET_OPTS -G --data-urlencode "key_pair_id=MyKeyPair" http://localhost:8080/locksmith/v1/key)
CMND_STATUS=$(echo "$CMND" | jq .status)
checkStatus "[KEY PAIRS][GET][READ]" $CMND_STATUS "Reading MyKeyPair Key Pair in the default Key Store"
if [[ $VERBOSITY == "2" ]]; then echo -e "${CMND}\n"; fi

# Read the private key of the Server Key Pair in the example-labs key store
CMND=$(curl $CURL_GET_OPTS -G --data-urlencode "key_pair_id=Server Key Pair" --data-urlencode "key_store_id=example-labs" --data-urlencode "passphrase=" http://localhost:8080/locksmith/v1/key)
CMND_STATUS=$(echo "$CMND" | jq .status)
checkStatus "[KEY PAIRS][GET][READ]" $CMND_STATUS "Reading Server Key Pair in the example-labs Key Store, plaintxt privKey"
if [[ $VERBOSITY == "2" ]]; then echo -e "${CMND}\n"; fi

# Read the private key of the VDI Terminal Key Pair in the default key store
CMND=$(curl $CURL_GET_OPTS -G --data-urlencode "key_pair_id=VDI Terminal" --data-urlencode "passphrase=s3cr3t" http://localhost:8080/locksmith/v1/key)
CMND_STATUS=$(echo "$CMND" | jq .status)
checkStatus "[KEY PAIRS][GET][READ]" $CMND_STATUS "Reading VDI Terminal Key Pair in the default Key Store with pass"
if [[ $VERBOSITY == "2" ]]; then echo -e "${CMND}\n"; fi

echo ""

###################################################################################- ROOT CA
# Generate a Root CA
CMND=$(curl $CURL_POST_OPTS --data '{"subject":{"common_name":"Example Labs Root Certificate Authority","organization":["Example Labs"],"organizational_unit":["Example Labs Cyber and Information Security"]},"expiration_date": [10,0,1],"san_data":{"email_addresses":["certmaster@example.labs"],"uris":["https://ca.example.labs:443/"]}}' \
  http://localhost:8080/locksmith/v1/root)
CMND_STATUS=$(echo "$CMND" | jq .status)
checkStatus "[ROOT CA][POST][CREATE]" $CMND_STATUS "Creating Example Labs Root Certificate Authority"
if [[ $VERBOSITY == "2" ]]; then echo -e "${CMND}\n"; fi

# Read the list of roots
CMND=$(curl $CURL_GET_OPTS http://localhost:8080/locksmith/v1/roots)
CMND_STATUS=$(echo "$CMND" | jq .status)
checkStatus "[ROOT CA][GET][LIST]" $CMND_STATUS "Listing Root CAs"
if [[ $VERBOSITY == "2" ]]; then echo -e "${CMND}\n"; fi

# Read the Root Certificate Authority
CMND=$(curl $CURL_GET_OPTS -G --data-urlencode "cn_path=Example Labs Root Certificate Authority" http://localhost:8080/locksmith/v1/authority)
CMND_STATUS=$(echo "$CMND" | jq .status)
checkStatus "[AUTHORITY][GET][READ]" $CMND_STATUS "Reading Example Labs Root Certificate Authority"
if [[ $VERBOSITY == "2" ]]; then echo -e "${CMND}\n"; fi

echo ""

###################################################################################- INTERMEDIATE CA
# Generate an Intermediate Certificate Authority
CMND=$(curl $CURL_POST_OPTS --data '{"parent_cn_path": "Example Labs Root Certificate Authority", "certificate_config":{"subject":{"common_name":"Example Labs Intermediate Certificate Authority","organization":["Example Labs"],"organizational_unit":["Example Labs Cyber and Information Security"]},"expiration_date": [3,0,1],"san_data":{"email_addresses":["certmaster@example.labs"],"uris":["https://ca.example.labs:443/"]}}}' \
  http://localhost:8080/locksmith/v1/intermediate)
CMND_STATUS=$(echo "$CMND" | jq .status)
checkStatus "[INTERMEDIATE][POST][CREATE]" $CMND_STATUS "Creating Example Labs Intermediate CA in Root CA"
if [[ $VERBOSITY == "2" ]]; then echo -e "${CMND}\n"; fi

# Read the Intermediate Certificate Authorities of the Root CA
CMND=$(curl $CURL_GET_OPTS -G --data-urlencode "parent_cn_path=Example Labs Root Certificate Authority" http://localhost:8080/locksmith/v1/intermediates)
CMND_STATUS=$(echo "$CMND" | jq .status)
checkStatus "[INTERMEDIATES][GET][LIST]" $CMND_STATUS "Listing Intermediate CAs in Example Labs Root CA"
if [[ $VERBOSITY == "2" ]]; then echo -e "${CMND}\n"; fi

# Read the Intermediate Certificate Authority
CMND=$(curl $CURL_GET_OPTS -G --data-urlencode "cn_path=Example Labs Root Certificate Authority/example-labs-intermediate-certificate-authority" http://localhost:8080/locksmith/v1/authority)
CMND_STATUS=$(echo "$CMND" | jq .status)
checkStatus "[AUTHORITY][GET][READ]" $CMND_STATUS "Reading Example Labs Intermediate CA in Root CA"
if [[ $VERBOSITY == "2" ]]; then echo -e "${CMND}\n"; fi

echo ""

###################################################################################- CERTIFICATE REQUESTS
# Read the list of CSRs
CMND=$(curl $CURL_GET_OPTS -G --data-urlencode "parent_cn_path=Example Labs Root Certificate Authority" http://localhost:8080/locksmith/v1/certificate-requests)
CMND_STATUS=$(echo "$CMND" | jq .status)
checkStatus "[CERTIFICATE REQUESTS][GET][LIST]" $CMND_STATUS "Listing CSRs in Example Labs Root CA"
if [[ $VERBOSITY == "2" ]]; then echo -e "${CMND}\n"; fi

echo ""

###################################################################################- CERTIFICATES
CMND=$(curl $CURL_GET_OPTS -G --data-urlencode "parent_cn_path=Example Labs Root Certificate Authority" http://localhost:8080/locksmith/v1/certificates)
CMND_STATUS=$(echo "$CMND" | jq .status)
checkStatus "[CERTIFICATES][GET][LIST]" $CMND_STATUS "Listing Certificates in Example Labs Root CA"
if [[ $VERBOSITY == "2" ]]; then echo -e "${CMND}\n"; fi

echo ""

###################################################################################- CERTIFICATE REVOCATION LISTS

###################################################################################- CERTIFICATE BUNDLE

# Generate a Server Certificate for OpenVPN
#CMND=$(curl $CURL_POST--data '{"parent_cn_path": "Example Labs Root Certificate Authority/Example Labs Intermediate Certificate Authority", "certificate_config":{"subject":{"common_name":"Example Labs OpenVPN Server","organization":["Example Labs"],"organizational_unit":["Example Labs Cyber and Information Security"]},"expiration_date": [1,0,1]}}' \
#  http://localhost:8080/locksmith/v1/certificates)
#checkStatus "[CERTIFICATE]" $CMND_STATUS "Create OpenVPN Server Certificate in Example Labs Signing CA"


kill -9 $RUN_PID

exit