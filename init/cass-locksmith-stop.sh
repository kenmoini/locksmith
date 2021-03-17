#!/bin/bash

set -x

source /etc/locksmith/caas_vars.sh

echo "Killing container..."
${CONTAINER_RUNTIME} kill $CONTAINER_NAME

echo "Removing container..."
${CONTAINER_RUNTIME} rm $CONTAINER_NAME -f -i