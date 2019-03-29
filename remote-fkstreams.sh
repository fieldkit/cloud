#!/bin/bash

set -xe

pushd ~/conservify/dev-ops/provisioning
if [ -z "$APP_SERVER_ADDRESS" ]; then
    source ./setup-env.sh
fi
popd

if [ -z "$APP_SERVER_ADDRESS" ]; then
    echo No cloud configuration.
    exit 2
fi

echo "Starting ssh-agent..."
eval $(ssh-agent)
trap "ssh-agent -k" exit
ssh-add ~/.ssh/cfy.pem

pushd server/tools/fkstreams

go build -o fkstreams *.go
scp fkstreams core@$APP_SERVER_ADDRESS:

popd

ssh core@$APP_SERVER_ADDRESS ./fkstreams --scheme https --host api.fkdev.org --locate
