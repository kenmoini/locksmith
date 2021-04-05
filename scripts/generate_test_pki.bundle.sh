#!/bin/bash

# Exits on any error
set -e

./scripts/generate_test_pki.openssl.sh

./scripts/generate_test_pki.locksmith.sh

./scripts/generate_test_pki.compare.sh