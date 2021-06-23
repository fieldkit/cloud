#!/bin/bash

source aws.env

export FIELDKIT_POSTGRES_URL=postgres://fieldkit:password@127.0.0.1/fieldkit?sslmode=disable
export FIELDKIT_STREAMS_BUCKET_NAME=fk-streams
export FIELDKIT_ARCHIVER=default

set -xe

make setup fkdata

build/fkdata "$@"
