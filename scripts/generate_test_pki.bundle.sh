#!/bin/bash

# Exits on any error
set -e

./generate_test_pki.openssl.sh

./generate_test_pki.locksmith.sh

./generate_test_pki.compare.sh