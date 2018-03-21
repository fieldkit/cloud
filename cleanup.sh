#!/bin/bash

set -x

docker rmi fieldkit-server-build
docker rmi fieldkit-admin-build
docker rmi fieldkit-frontend-build
docker rmi fieldkit-landing-build
docker rmi conservify/fk-cloud:latest
DOCKER_TAG=`git rev-parse --abbrev-ref HEAD`
docker rmi conservify/fk-cloud:$DOCKER_TAG
