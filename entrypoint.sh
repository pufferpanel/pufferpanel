#!/usr/bin/env sh

if [ -z "${PUFFER_DOCKER_ROOT}" ]; then
    echo "PUFFER_DOCKER_ROOT must be set to use this image"
    echo "Please set this ENV option to the server directory volume, including the _data in the path"
    echo "Example: /var/lib/docker/volumes/pufferpanel-servers/_data"
    exit 1
fi

/pufferpanel/bin/pufferpanel run