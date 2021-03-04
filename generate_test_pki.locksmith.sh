#!/bin/bash

# Exits on any error
set -e

## Run the Locksmith application in the background
nohup ./locksmith -config config.yml.example &

RUN_PID=$!

# Wait a few seconds while the Locksmith server starts
sleep 5

# Generate a Root CA
curl --header "Content-Type: application/x-www-form-urlencoded" \
  --request POST \
  --data 'cert_info={"subject":{"common_name":"Example Labs Root Certificate Authority","organization":["Example Labs"],"organizationalUnit":["Example Labs Cyber and Information Security"]},"expiration_date": [10,0,0]}' \
  http://localhost:8080/locksmith/roots


kill -9 $RUN_PID

exit