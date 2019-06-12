#!/bin/bash
#
# I'm sorry for this?
#

set -xe

WORKSPACE=$(realpath temp-go)

rm -rf ${WORKSPACE}/src/github.com/fieldkit/cloud/server

mkdir -p ${WORKSPACE}/src/github.com/fieldkit/cloud/server

rsync -zvua --progress --exclude=vendor server/ ${WORKSPACE}/src/github.com/fieldkit/cloud/server/

export GOPATH=${WORKSPACE}
export GO111MODULE=on

(cd ${WORKSPACE}/src/github.com/fieldkit/cloud/server && go generate)

find ${WORKSPACE} -type d -exec chmod 755 {} \;

rsync -zua ${WORKSPACE}/src/github.com/fieldkit/cloud/server/ server/
