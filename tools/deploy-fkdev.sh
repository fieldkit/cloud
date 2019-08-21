#!/bin/bash

uname -a

export PATH=$PATH:/opt/bin

OLD_CONTAINER=`docker-compose -f /etc/docker/compose/fk-compose.yml ps -q fk-app`

if [ -z $OLD_CONTAINER ]; then
    echo "No existing container, just starting a new one...."

    cat *.bz2 | bunzip2 | docker load

    docker-compose -f /etc/docker/compose/fk-compose.yml rm -f fk-app
    docker-compose -f /etc/docker/compose/fk-compose.yml up -d fk-app
else
    echo "Found existing container: $OLD_CONTAINER"

    docker update --restart=no $OLD_CONTAINER

    OLD_IMAGE=`docker inspect --format='{{.Image}}' $OLD_CONTAINER`

    docker rmi -f conservify/fk-cloud:old
    docker tag $OLD_IMAGE conservify/fk-cloud:old
    docker rmi conservify/fk-cloud:master

    cat *.bz2 | bunzip2 | docker load

    docker-compose -f /etc/docker/compose/fk-compose.yml stop fk-app
    docker-compose -f /etc/docker/compose/fk-compose.yml rm -f fk-app
    docker-compose -f /etc/docker/compose/fk-compose.yml up -d fk-app

    docker rmi $OLD_IMAGE
fi
