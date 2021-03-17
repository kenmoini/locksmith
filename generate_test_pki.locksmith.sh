#!/bin/bash

# Exits on any error
set -e

## Clean up the directory
rm -rf ./.generated/roots/example-labs-root-certificate-authority

## Run the Locksmith application in the background
nohup ./locksmith -config config.yml.example &

RUN_PID=$!

# Wait a few seconds while the Locksmith server starts
sleep 5

# Generate a Root CA
curl --header "Content-Type: application/json" \
  --request POST \
  --data '{"subject":{"common_name":"Example Labs Root Certificate Authority","organization":["Example Labs"],"organizational_unit":["Example Labs Cyber and Information Security"]},"expiration_date": [10,0,1],"san_data":{"email_addresses":["certmaster@example.labs"],"uris":["https://ca.example.labs:443/"]}}' \
  http://localhost:8080/locksmith/v1/roots
echo -e "\n"

# Read the list of roots
curl --request GET http://localhost:8080/locksmith/v1/roots
echo -e "\n"

# Generate an Intermediate Certificate Authority
curl --header "Content-Type: application/json" \
  --request POST \
  --data '{"parent_cn_path": "Example Labs Root Certificate Authority", "certificate_config":{"subject":{"common_name":"Example Labs Intermediate Certificate Authority","organization":["Example Labs"],"organizational_unit":["Example Labs Cyber and Information Security"]},"expiration_date": [3,0,1],"san_data":{"email_addresses":["certmaster@example.labs"],"uris":["https://ca.example.labs:443/"]}}}' \
  http://localhost:8080/locksmith/v1/intermediates
echo -e "\n"

# Read the Intermediate Certificate Authorities of the Root CA
curl --request GET -G --data-urlencode "parent_cn_path=Example Labs Root Certificate Authority" "http://localhost:8080/locksmith/v1/intermediates"
echo -e "\n"

# Generate a Server Certificate for OpenVPN
#curl --header "Content-Type: application/json" \
#  --request POST \
#  --data '{"parent_cn_path": "Example Labs Root Certificate Authority/Example Labs Intermediate Certificate Authority", "certificate_config":{"subject":{"common_name":"Example Labs OpenVPN Server","organization":["Example Labs"],"organizational_unit":["Example Labs Cyber and Information Security"]},"expiration_date": [1,0,1]}}' \
#  http://localhost:8080/locksmith/v1/certificates
#echo -e "\n"


kill -9 $RUN_PID

exit