#!/bin/bash

set -x

source /etc/locksmith/caas_vars.sh

/etc/locksmith/caas_stop.sh

sleep 3

echo "Checking for stale network lock file..."
FILE_CHECK="/var/lib/cni/networks/${NETWORK_NAME}/${IP_ADDRESS}"
if [[ -f "$FILE_CHECK" ]]; then
    rm $FILE_CHECK
fi

echo "Starting container ${CONTAINER_NAME}..."
${CONTAINER_RUNTIME} run -d --name "${CONTAINER_NAME}" --network "${NETWORK_NAME}" --ip "${IP_ADDRESS}" -p "${CONTAINER_PORT}" -v ${VOLUME_MOUNT_ONE} ${RESOURCE_LIMITS} ${CONTAINER_SOURCE}