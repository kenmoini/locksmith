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
curl --header "Content-Type: application/x-www-form-urlencoded" \
  --request POST \
  --data 'cert_info={"subject":{"common_name":"Example Labs Root Certificate Authority","organization":["Example Labs"],"organizational_unit":["Example Labs Cyber and Information Security"]},"expiration_date": [10,0,1],"san_data":{"email_addresses":["certmaster@example.labs"],"uris":["https://ca.example.labs:443/"]}}' \
  http://localhost:8080/locksmith/roots

# Generate an Intermediate Certificate Authority
curl --header "Content-Type: application/x-www-form-urlencoded" \
  --request POST \
  --data 'ica_info={"parent_cn_path": "Example Labs Root Certificate Authority", "certificate_config":{"subject":{"common_name":"Example Labs Intermediate Certificate Authority","organization":["Example Labs"],"organizational_unit":["Example Labs Cyber and Information Security"]},"expiration_date": [1,1,1],"san_data":{"email_addresses":["certmaster@example.labs"],"uris":["https://ca.example.labs:443/"]}}}' \
  http://localhost:8080/locksmith/intermediates

# Read the Intermediate Certificate Authorities of the Root CA
curl --header "Content-Type: application/x-www-form-urlencoded" \
  --request GET \
  --data 'ica_info={"parent_cn_path": "Example Labs Root Certificate Authority"}' \
  http://localhost:8080/locksmith/intermediates


kill -9 $RUN_PID

exit