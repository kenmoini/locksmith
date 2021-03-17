#!/bin/bash

CONTAINER_RUNTIME="/usr/bin/podman"

CONTAINER_NAME="locksmith"
NETWORK_NAME="lanBridge"
IP_ADDRESS="192.168.42.9"
CONTAINER_PORT="8080"

VOLUME_MOUNT_ONE="/etc/locksmith:/etc/locksmith/"
CONTAINER_SOURCE="quay.io/kenmoini/locksmith:latest"

RESOURCE_LIMITS="-m 512m"