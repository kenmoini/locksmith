#!/bin/bash

# Exits on any error
set -e

CUR_DIR=$(pwd)

## Clean up the directory
rm -rf ${CUR_DIR}/.generated/roots/example-labs-root-certificate-authority
rm -rf ${CUR_DIR}/.generated/keys/

## Run the Locksmith application in the background
nohup ${CUR_DIR}/dist/locksmith -config configs/config.yml.example &

RUN_PID=$!

# Wait a few seconds while the Locksmith server starts
sleep 5

###################################################################################- KEY STORES
# Read the list of key stores
curl --request GET http://localhost:8080/locksmith/v1/keystores
echo -e "\n"
# Create a new key stores
curl --header "Content-Type: application/json" --request POST \
  --data '{"key_store_name": "Example Labs"}' http://localhost:8080/locksmith/v1/keystores
echo -e "\n"
# Read the list of key stores
curl --request GET http://localhost:8080/locksmith/v1/keystores
echo -e "\n"

###################################################################################- KEY PAIRS
# Create a key pair in the default key store
curl --header "Content-Type: application/json" --request POST \
  --data '{"key_pair_id": "MyKeyPair"}' http://localhost:8080/locksmith/v1/keys
echo -e ""
# Create a key pair in the example-labs key store
curl --header "Content-Type: application/json" --request POST \
  --data '{"key_pair_id": "Server Key Pair", "key_store_id": "example-labs"}' http://localhost:8080/locksmith/v1/keys
echo -e ""
# Create a key pair in the default key store
curl --header "Content-Type: application/json" --request POST \
  --data '{"key_pair_id": "VDI Terminal"}' http://localhost:8080/locksmith/v1/keys
echo -e "\n"
# Read the list of key pairs in the default store
curl --request GET http://localhost:8080/locksmith/v1/keys
echo -e "\n"
# Read the list of key pairs in the example-labs store
curl --request GET -G --data-urlencode "key_store_id=example-labs" http://localhost:8080/locksmith/v1/keys
echo -e "\n"

###################################################################################- ROOT CA
# Generate a Root CA
curl --header "Content-Type: application/json" \
  --request POST \
  --data '{"subject":{"common_name":"Example Labs Root Certificate Authority","organization":["Example Labs"],"organizational_unit":["Example Labs Cyber and Information Security"]},"expiration_date": [10,0,1],"san_data":{"email_addresses":["certmaster@example.labs"],"uris":["https://ca.example.labs:443/"]}}' \
  http://localhost:8080/locksmith/v1/roots
echo -e "\n"

# Read the list of roots
curl --request GET http://localhost:8080/locksmith/v1/roots
echo -e "\n"

###################################################################################- INTERMEDIATE CA
# Generate an Intermediate Certificate Authority
curl --header "Content-Type: application/json" \
  --request POST \
  --data '{"parent_cn_path": "Example Labs Root Certificate Authority", "certificate_config":{"subject":{"common_name":"Example Labs Intermediate Certificate Authority","organization":["Example Labs"],"organizational_unit":["Example Labs Cyber and Information Security"]},"expiration_date": [3,0,1],"san_data":{"email_addresses":["certmaster@example.labs"],"uris":["https://ca.example.labs:443/"]}}}' \
  http://localhost:8080/locksmith/v1/intermediates
echo -e "\n"

# Read the Intermediate Certificate Authorities of the Root CA
curl --request GET -G --data-urlencode "parent_cn_path=Example Labs Root Certificate Authority" "http://localhost:8080/locksmith/v1/intermediates"
echo -e "\n"

###################################################################################- CERTIFICATE REQUESTS
# Read the list of CSRs
curl --request GET -G --data-urlencode "parent_cn_path=Example Labs Root Certificate Authority" http://localhost:8080/locksmith/v1/certificate-requests
echo -e "\n"


###################################################################################- CERTIFICATES
curl --request GET -G --data-urlencode "parent_cn_path=Example Labs Root Certificate Authority" http://localhost:8080/locksmith/v1/certificates
echo -e "\n"

###################################################################################- CERTIFICATE REVOCATION LISTS

###################################################################################- CERTIFICATE BUNDLE

# Generate a Server Certificate for OpenVPN
#curl --header "Content-Type: application/json" \
#  --request POST \
#  --data '{"parent_cn_path": "Example Labs Root Certificate Authority/Example Labs Intermediate Certificate Authority", "certificate_config":{"subject":{"common_name":"Example Labs OpenVPN Server","organization":["Example Labs"],"organizational_unit":["Example Labs Cyber and Information Security"]},"expiration_date": [1,0,1]}}' \
#  http://localhost:8080/locksmith/v1/certificates
#echo -e "\n"


kill -9 $RUN_PID

exit